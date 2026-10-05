package triggers

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

// TestCoordinator_RegisterTriggers covers registration: the happy path with event
// delivery, and each failure path that must free the workflow-count limit slot.
func TestCoordinator_RegisterTriggers(t *testing.T) {
	t.Parallel()

	t.Run("registers triggers and delivers events to the engine", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		var gotReq capabilities.TriggerRegistrationRequest
		eventCh := make(chan capabilities.TriggerResponse)
		trigger := capmocks.NewTriggerCapability(t)
		trigger.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).
			Run(func(_ context.Context, req capabilities.TriggerRegistrationRequest) { gotReq = req }).
			Return((<-chan capabilities.TriggerResponse)(eventCh), nil).Once()
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()
		f.capReg.EXPECT().GetTrigger(mock.Anything, testTriggerCapID).Return(trigger, nil).Once()

		triggerIDs, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		assert.Equal(t, []string{testTriggerCapID}, triggerIDs)
		assert.Equal(t, 1, f.limits.inUse())
		assert.True(t, f.registered())

		assert.Equal(t, RegistrationID(validWorkflowID, 0), gotReq.TriggerID)
		assert.Equal(t, validWorkflowID, gotReq.Metadata.WorkflowID)
		assert.Equal(t, testOwner, gotReq.Metadata.WorkflowOwner)
		assert.Equal(t, testDonID, gotReq.Metadata.WorkflowDonID)
		assert.Equal(t, uint32(pinnedWorkflowDonConfigVersion), gotReq.Metadata.WorkflowDonConfigVersion)

		eventCh <- triggerEvent("evt-1")
		require.Eventually(t, func() bool { return len(f.engine.executed()) == 1 }, 5*time.Second, 10*time.Millisecond)

		got := f.engine.executed()[0]
		assert.Equal(t, validWorkflowID, got.WorkflowID)
		assert.Equal(t, testTriggerCapID, got.TriggerCapID)
		assert.Equal(t, 0, got.TriggerIndex)
		assert.Equal(t, "evt-1", got.Event.Event.ID)
		assert.Equal(t, f.clock.Now(), got.ObservedAt)

		f.engine.mu.Lock()
		assert.Equal(t, testTenant, f.engine.tenants[0], "execution ctx must carry the workflow's CRE")
		f.engine.mu.Unlock()
	})

	t.Run("invalid workflow ID", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		sub := newTestSubscriber()
		sub.tenant.Workflow = "not-hex"

		_, err := f.c.RegisterTriggers(t.Context(), sub, testParams(t))
		require.ErrorContains(t, err, "invalid workflow id")
		assert.Equal(t, 0, sub.calls)
		assert.Equal(t, 0, f.limits.inUse())
	})

	t.Run("workflow limit reached maps to the scoped sentinel and never subscribes", func(t *testing.T) {
		t.Parallel()
		for _, tc := range []struct {
			scope settings.Scope
			want  error
		}{
			{settings.ScopeOwner, types.ErrPerOwnerWorkflowCountLimitReached},
			{settings.ScopeGlobal, types.ErrGlobalWorkflowCountLimitReached},
		} {
			f := newCoordinatorFixture(t)
			f.limits.useErr = limits.ErrorResourceLimited[int]{Scope: tc.scope, Limit: 1, Used: 1, Amount: 1}
			sub := newTestSubscriber()

			_, err := f.c.RegisterTriggers(t.Context(), sub, testParams(t))
			require.ErrorIs(t, err, tc.want)
			assert.Equal(t, 0, sub.calls)
			assert.False(t, f.registered())
		}
	})

	t.Run("the limit is acquired under the workflow's tenant", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		f.limits.useErr = errors.New("boom")

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.Error(t, err)
		f.limits.mu.Lock()
		defer f.limits.mu.Unlock()
		require.Len(t, f.limits.tenants, 1)
		assert.Equal(t, testTenant, f.limits.tenants[0])
	})

	t.Run("subscribe failure frees the limit", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		sub := newTestSubscriber()
		sub.err = errors.New("wasm boom")

		_, err := f.c.RegisterTriggers(t.Context(), sub, testParams(t))
		require.ErrorContains(t, err, "failed to subscribe to triggers")
		assert.Equal(t, 0, f.limits.inUse())
		assert.False(t, f.registered())
	})

	t.Run("registration failure frees the limit and leaves nothing registered", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		f.capReg.EXPECT().GetTrigger(mock.Anything, testTriggerCapID).Return(nil, errors.New("not found")).Once()

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.ErrorContains(t, err, "trigger capability not found")
		assert.Equal(t, 0, f.limits.inUse())
		assert.False(t, f.registered())
	})
}

// TestCoordinator_NoDeliver covers dropped delivery: events must not reach an
// engine that is gone from the registry or not coordinated.
func TestCoordinator_NoDeliver(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		mutate func(f *coordinatorFixture)
	}{
		{"drops events once the engine is gone from the registry", func(f *coordinatorFixture) {
			f.engines.mu.Lock()
			delete(f.engines.engines, f.wid)
			f.engines.mu.Unlock()
		}},
		{"drops events for an engine that is not coordinated", func(f *coordinatorFixture) {
			f.engine.coordinated = false
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := newCoordinatorFixture(t)
			trigger, eventCh := f.expectTrigger(t)
			trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()

			_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
			require.NoError(t, err)
			tc.mutate(f)

			// The channel is unbuffered: the second send only completes once the
			// reader has finished handling the first.
			eventCh <- triggerEvent("evt-1")
			eventCh <- triggerEvent("evt-2")
			assert.Empty(t, f.engine.executed())
		})
	}
}

// TestCoordinator_Ack covers acknowledgement: routing to the registration's
// handle, and the failure paths for unknown or malformed registration IDs.
func TestCoordinator_Ack(t *testing.T) {
	t.Parallel()

	t.Run("routes to the registration's handle", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		trigger, _ := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()
		regID := RegistrationID(validWorkflowID, 0)
		trigger.EXPECT().AckEvent(mock.Anything, regID, "evt-1", testMethod).Return(nil).Once()

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		require.NoError(t, f.c.Ack(t.Context(), testTriggerCapID, regID, "evt-1"))
	})

	t.Run("unknown workflow", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		err := f.c.Ack(t.Context(), testTriggerCapID, RegistrationID(validWorkflowID, 0), "evt-1")
		require.ErrorContains(t, err, "failed to find trigger")
	})

	t.Run("malformed registration ID", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		err := f.c.Ack(t.Context(), testTriggerCapID, "not-a-registration", "evt-1")
		require.ErrorContains(t, err, "invalid trigger registration ID")
	})
}

// TestCoordinator_UnregisterTriggers covers unregistration: ingress stops
// immediately, handles stay available until in-flight executions drain, and the
// limit slot is freed exactly once, also on retry and on drain timeout.
func TestCoordinator_UnregisterTriggers(t *testing.T) {
	t.Parallel()

	t.Run("unknown workflow", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		require.ErrorIs(t, f.c.UnregisterTriggers(t.Context(), validWorkflowID), ErrWorkflowNotCoordinated)
	})

	t.Run("registration state is gone and the slot is freed when no execution is in flight", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		trigger, _ := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.MatchedBy(func(req capabilities.TriggerRegistrationRequest) bool {
			return req.TriggerID == RegistrationID(validWorkflowID, 0) && req.Method == testMethod
		})).Return(nil).Once()

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		require.NoError(t, f.c.UnregisterTriggers(t.Context(), validWorkflowID))

		require.Eventually(t, func() bool { return !f.registered() && f.limits.inUse() == 0 }, 5*time.Second, 10*time.Millisecond)
		assert.Equal(t, 1, f.limits.freeCount())
		f.limits.mu.Lock()
		assert.Equal(t, testTenant, f.limits.tenants[len(f.limits.tenants)-1], "the limit is freed under the workflow's tenant")
		f.limits.mu.Unlock()
	})

	t.Run("keeps registration state and slot while an execution is in flight so it can still ACK", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		trigger, eventCh := f.expectTrigger(t)
		regID := RegistrationID(validWorkflowID, 0)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.MatchedBy(func(req capabilities.TriggerRegistrationRequest) bool {
			return req.TriggerID == regID && req.Method == testMethod
		})).Return(nil).Once()
		trigger.EXPECT().AckEvent(mock.Anything, regID, "evt-1", testMethod).Return(nil).Once()
		started, release := f.blockingExecution()

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		eventCh <- triggerEvent("evt-1")
		<-started

		require.NoError(t, f.c.UnregisterTriggers(t.Context(), validWorkflowID))

		// Ingress is stopped: an event sent after unregister is never delivered
		// to the engine. The send runs in a goroutine since the channel is
		// unbuffered and no reader may remain to receive it.
		go func() { eventCh <- triggerEvent("evt-2") }()
		require.Never(t, func() bool {
			return slices.ContainsFunc(f.engine.executed(), func(e CoordinatedEvent) bool {
				return e.Event.Event.ID == "evt-2"
			})
		}, 100*time.Millisecond, 10*time.Millisecond, "event sent after unregister must not be delivered")

		// The in-flight execution still resolves its handle.
		require.NoError(t, f.c.Ack(t.Context(), testTriggerCapID, regID, "evt-1"))
		assert.True(t, f.registered())
		assert.Equal(t, 1, f.limits.inUse())

		close(release)
		require.Eventually(t, func() bool { return !f.registered() && f.limits.inUse() == 0 }, 5*time.Second, 10*time.Millisecond)
	})

	t.Run("repeated calls unregister with the capability once", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		trigger, _ := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		require.NoError(t, f.c.UnregisterTriggers(t.Context(), validWorkflowID))
		// The second call may land before or after the release; both are fine.
		err = f.c.UnregisterTriggers(t.Context(), validWorkflowID)
		if err != nil {
			require.ErrorIs(t, err, ErrWorkflowNotCoordinated)
		}
		require.Eventually(t, func() bool { return f.limits.inUse() == 0 }, 5*time.Second, 10*time.Millisecond)
		assert.Equal(t, 1, f.limits.freeCount())
	})

	t.Run("a failed unregister is returned and retried on the next call", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		trigger, eventCh := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(errors.New("capability down")).Once()
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()
		// Keep a reader busy so the retry happens before the release.
		started, release := f.blockingExecution()

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		eventCh <- triggerEvent("evt-1")
		<-started

		require.ErrorContains(t, f.c.UnregisterTriggers(t.Context(), validWorkflowID), "failed to unregister 1 of 1 triggers")
		require.NoError(t, f.c.UnregisterTriggers(t.Context(), validWorkflowID))

		close(release)
		require.Eventually(t, func() bool { return f.limits.inUse() == 0 }, 5*time.Second, 10*time.Millisecond)
		assert.Equal(t, 1, f.limits.freeCount(), "only one release waiter may free the slot")
	})

	t.Run("releases after the drain timeout even if an execution hangs", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		trigger, eventCh := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()
		started, release := f.blockingExecution()
		defer close(release)

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		eventCh <- triggerEvent("evt-1")
		<-started

		require.NoError(t, f.c.UnregisterTriggers(t.Context(), validWorkflowID))
		require.NoError(t, f.clock.BlockUntilContext(t.Context(), 1))
		assert.True(t, f.registered())

		f.clock.Advance(defaultDrainTimeout)
		require.Eventually(t, func() bool { return !f.registered() && f.limits.inUse() == 0 }, 5*time.Second, 10*time.Millisecond)
	})
}

// TestCoordinator_ReRegisterWhileDraining covers the overlap of two registrations
// for the same workflow: the old registration drains while the new one is active.
//
// The test registers, starts an execution that blocks, and unregisters. It then
// registers again before the old registration drains. The new registration must
// replace the old state, and the old release waiter must not delete the new state
// or free the new slot when the old readers exit. The new registration must still
// ACK against its own handles.
func TestCoordinator_ReRegisterWhileDraining(t *testing.T) {
	t.Parallel()
	f := newCoordinatorFixture(t)

	oldTrigger, oldCh := f.expectTrigger(t)
	oldTrigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()
	started, release := f.blockingExecution()

	_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
	require.NoError(t, err)
	oldCh <- triggerEvent("evt-1")
	<-started
	require.NoError(t, f.c.UnregisterTriggers(t.Context(), validWorkflowID))

	oldWT, _ := f.c.workflows.get(validWorkflowID)

	// The replacement engine registers while the old registration is still draining.
	f.engine.execute = nil
	newTrigger, _ := f.expectTrigger(t)
	newTrigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()
	_, err = f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
	require.NoError(t, err)

	newWT, _ := f.c.workflows.get(validWorkflowID)
	require.NotSame(t, oldWT, newWT, "registering again must replace, not reuse, the old state")
	assert.Equal(t, 2, f.limits.inUse(), "the old slot is held until its drain completes")

	close(release)
	oldWT.readers.Wait() // the old registration's own readers have now exited

	// The old release waiter, woken by the same Wait(), may still run and
	// (wrongly, if buggy) drop the new registration. Poll instead of a single
	// sleep-then-check: this fails the instant the bug appears rather than
	// only if it happens to land inside a guessed window.
	require.Never(t, func() bool {
		cw, _ := f.c.workflows.get(validWorkflowID)
		return cw != newWT
	}, 100*time.Millisecond, 2*time.Millisecond, "the old waiter must not drop the new registration")
	require.Eventually(t, func() bool { return f.limits.inUse() == 1 }, time.Second, time.Millisecond,
		"the old waiter must still free its own slot")

	regID := RegistrationID(validWorkflowID, 0)
	newTrigger.EXPECT().AckEvent(mock.Anything, regID, "evt-2", testMethod).Return(nil).Once()
	require.NoError(t, f.c.Ack(t.Context(), testTriggerCapID, regID, "evt-2"))
}

// TestCoordinator_CloseUnregistersRemainingWorkflows covers shutdown: Close
// unregisters the triggers that remain registered when the coordinator stops.
func TestCoordinator_CloseUnregistersRemainingWorkflows(t *testing.T) {
	t.Parallel()
	capReg := regmocks.NewCapabilitiesRegistry(t)
	engines := newFakeEngineRegistry()
	wid, err := types.WorkflowIDFromHex(validWorkflowID)
	require.NoError(t, err)
	engines.set(wid, &fakeEngine{coordinated: true})

	c := NewCoordinator(newRegisterDeps(t, capReg, limits.NewGateLimiter(true)), engines, &fakeWorkflowLimits{}, clockwork.NewFakeClock())
	require.NoError(t, c.Start(t.Context()))

	trigger := capmocks.NewTriggerCapability(t)
	trigger.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).
		Return((<-chan capabilities.TriggerResponse)(make(chan capabilities.TriggerResponse)), nil).Once()

	// The assertion: Close must unregister the remaining registration with the
	// capability, or this Once() expectation fails the test.
	trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.MatchedBy(func(req capabilities.TriggerRegistrationRequest) bool {
		return req.TriggerID == RegistrationID(validWorkflowID, 0) && req.Method == testMethod
	})).Return(nil).Once()
	capReg.EXPECT().GetTrigger(mock.Anything, testTriggerCapID).Return(trigger, nil).Once()

	_, err = c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
	require.NoError(t, err)
	require.NoError(t, c.Close())
}
