package triggers

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
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

	t.Run("subscribe failure frees the limit and closes the queue", func(t *testing.T) {
		t.Parallel()
		admission, queues := newTrackedAdmission()
		f := newCoordinatorFixture(t, withAdmission(admission))
		sub := newTestSubscriber()
		sub.err = errors.New("wasm boom")

		_, err := f.c.RegisterTriggers(t.Context(), sub, testParams(t))
		require.ErrorContains(t, err, "failed to subscribe to triggers")
		assert.Equal(t, 0, f.limits.inUse())
		assert.False(t, f.registered())
		assert.True(t, queues.get(t, 0).closed.Load(), "the registration's queue is closed")
	})

	t.Run("queue creation failure frees the limit and never subscribes", func(t *testing.T) {
		t.Parallel()
		admission := newTestAdmission()
		admission.NewTriggerEventQueue = func() (limits.QueueLimiter[CoordinatedEvent], error) {
			return nil, errors.New("queue boom")
		}
		f := newCoordinatorFixture(t, withAdmission(admission))
		sub := newTestSubscriber()

		_, err := f.c.RegisterTriggers(t.Context(), sub, testParams(t))
		require.ErrorContains(t, err, "failed to create trigger event queue")
		assert.Equal(t, 0, f.limits.inUse())
		assert.Zero(t, sub.calls)
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

// TestCoordinator_MetricsIncludeSdkLabel makes sure the coordinator's metrics
// carry the sdk label.
func TestCoordinator_MetricsIncludeSdkLabel(t *testing.T) {
	t.Parallel()
	wfMetrics, reader := newTestMeterMetrics(t)
	f := newCoordinatorFixture(t, withMetrics(wfMetrics))
	trigger, _ := f.expectTrigger(t)
	trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()

	_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
	require.NoError(t, err)
	require.NoError(t, f.c.UnregisterTriggers(t.Context(), validWorkflowID))

	sdk := map[string]string{platform.KeySDK: testSdkName}
	require.Eventually(t, func() bool {
		return counterTotal(t, reader, "platform_engine_workflow_unregistered_count", sdk) == 1
	}, 5*time.Second, 5*time.Millisecond)
}

// TestCoordinator_NoDeliver covers dropped delivery: events must not reach an
// engine that is gone from the registry or not coordinated, and each case
// must be countable the same way every other drop reason is.
func TestCoordinator_NoDeliver(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name       string
		mutate     func(f *coordinatorFixture)
		dropReason string
	}{
		{"drops events once the engine is gone from the registry", func(f *coordinatorFixture) {
			f.engines.remove(f.wid)
		}, monitoring.TriggerDropReasonEngineGone},
		{"drops events for an engine that is not coordinated", func(f *coordinatorFixture) {
			f.engine.coordinated = false
		}, monitoring.TriggerDropReasonEngineNotCoordinated},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			wfMetrics, reader := newTestMeterMetrics(t)
			f := newCoordinatorFixture(t, withMetrics(wfMetrics))
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

			dropped := map[string]string{platform.KeyTriggerDropReason: tc.dropReason}
			require.Eventually(t, func() bool {
				return counterTotal(t, reader, "platform_engine_trigger_event_dropped_total", dropped) == 2
			}, 5*time.Second, 5*time.Millisecond)
		})
	}
}

// TestCoordinator_DeliverAdmission covers what the coordinator does with an
// admission error from the engine (e.g. the shard failover manager on a sharded
// node). A denied event is ACKed and counted as an admission_denied drop, while a cached event
// is left un-ACKed so the secondary can replay it.
func TestCoordinator_DeliverAdmission(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		engineErr   error
		wantAck     bool
		wantDropped int64
	}{
		{
			name:        "denied is ACKed and counted",
			engineErr:   fmt.Errorf("shard ownership denied: not owner: %w", ErrAdmissionDenied),
			wantAck:     true,
			wantDropped: 1,
		},
		{name: "cached is left for replay", engineErr: ErrAdmissionCached},
		{name: "other error is not ACKed", engineErr: errors.New("metering reserve failed")},
		{name: "executed is not ACKed by the coordinator"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			wfMetrics, reader := newTestMeterMetrics(t)
			f := newCoordinatorFixture(t, withMetrics(wfMetrics))
			f.engine.err = tc.engineErr
			trigger, eventCh := f.expectTrigger(t)
			trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()
			if tc.wantAck {
				trigger.EXPECT().AckEvent(mock.Anything, RegistrationID(validWorkflowID, 0), "evt-1", testMethod).Return(nil).Once()
			}

			_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
			require.NoError(t, err)
			eventCh <- triggerEvent("evt-1")
			// Wait for the admitter to deliver the event: unregistering first
			// would drop it from the queue.
			require.Eventually(t, func() bool { return len(f.engine.executed()) == 1 }, 5*time.Second, 10*time.Millisecond)

			// The registration is only released once its executions have
			// returned, so the delivery, including any ACK and drop metric, is complete.
			require.NoError(t, f.c.UnregisterTriggers(t.Context(), validWorkflowID))
			require.Eventually(t, func() bool { return !f.registered() }, 5*time.Second, 10*time.Millisecond)

			if !tc.wantAck {
				trigger.AssertNotCalled(t, "AckEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}
			denied := map[string]string{platform.KeyTriggerDropReason: monitoring.TriggerDropReasonAdmissionDenied}
			assert.Equal(t, tc.wantDropped, counterTotal(t, reader, "platform_engine_trigger_event_dropped_total", denied),
				"admission_denied drops")
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

	t.Run("registration state is gone, the slot is freed and the queue is closed when no execution is in flight", func(t *testing.T) {
		t.Parallel()
		admission, queues := newTrackedAdmission()
		f := newCoordinatorFixture(t, withAdmission(admission))
		trigger, _ := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.MatchedBy(func(req capabilities.TriggerRegistrationRequest) bool {
			return req.TriggerID == RegistrationID(validWorkflowID, 0) && req.Method == testMethod
		})).Return(nil).Once()

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		require.NoError(t, f.c.UnregisterTriggers(t.Context(), validWorkflowID))

		require.Eventually(t, func() bool { return !f.registered() && f.limits.inUse() == 0 }, 5*time.Second, 10*time.Millisecond)
		assert.Equal(t, 1, f.limits.freeCount())
		assert.True(t, queues.get(t, 0).closed.Load(), "the registration's queue is closed")
		f.limits.mu.Lock()
		assert.Equal(t, testTenant, f.limits.tenants[len(f.limits.tenants)-1], "the limit is freed under the workflow's tenant")
		f.limits.mu.Unlock()
	})

	t.Run("keeps registration state while an execution is in flight so it can still ACK", func(t *testing.T) {
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

		// The in-flight execution still resolves its handle, even though the
		// workflow slot is already free (see the slot-release test below).
		require.NoError(t, f.c.Ack(t.Context(), testTriggerCapID, regID, "evt-1"))
		assert.True(t, f.registered())
		assert.Equal(t, 0, f.limits.inUse())

		// Ingress is stopped: an event sent after unregister must never be
		// delivered, including after the in-flight delivery returns. The send
		// stays parked on the unbuffered channel until the reader exits.
		go func() { eventCh <- triggerEvent("evt-2") }()

		close(release)
		require.Eventually(t, func() bool { return !f.registered() && f.limits.inUse() == 0 }, 5*time.Second, 10*time.Millisecond)

		executed := f.engine.executed()
		require.Len(t, executed, 1, "evt-2 must not be delivered after unregister")
		assert.Equal(t, "evt-1", executed[0].Event.Event.ID)
	})

	// UnregisterTriggers frees the workflow slot right away, even while an
	// execution is still in flight and its handles are kept around so it can
	// ACK. Slot release and handle retention are separate concerns.
	t.Run("frees the workflow slot immediately, even while an execution is still in flight to ACK", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		trigger, eventCh := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()
		trigger.EXPECT().AckEvent(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
		started, release := f.blockingExecution()
		defer close(release)

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)
		eventCh <- triggerEvent("evt-1")
		<-started

		require.NoError(t, f.c.UnregisterTriggers(t.Context(), validWorkflowID))

		// No Eventually: the execution is still held, so the drain cannot have
		// finished yet. If the slot is only free once drained, this is 1, not 0.
		assert.Equal(t, 0, f.limits.inUse(), "the workflow slot must be freed by UnregisterTriggers itself, not by the drain that follows it")
		assert.Equal(t, 1, f.limits.freeCount())
	})

	t.Run("repeated calls unregister with the capability once", func(t *testing.T) {
		t.Parallel()
		f := newCoordinatorFixture(t)
		trigger, _ := f.expectTrigger(t)
		// require that UnregisterTrigger on the trigger capability is called exactly once
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

	t.Run("a failed unregistration does not prevent a subsequent successful unregistration", func(t *testing.T) {
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
	admission, queues := newTrackedAdmission()
	f := newCoordinatorFixture(t, withAdmission(admission))

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
	assert.Equal(t, 1, f.limits.inUse(), "UnregisterTriggers already freed the old slot, so the new registration takes the only slot in use")

	close(release)
	oldWT.readers.Wait() // the old registration's own readers have now exited

	// The old release waiter, woken by the same Wait(), may still run and
	// (wrongly, if buggy) drop the new registration or free its slot. Poll
	// instead of a single sleep-then-check: this fails the instant the bug
	// appears rather than only if it happens to land inside a guessed window.
	require.Never(t, func() bool {
		cw, _ := f.c.workflows.get(validWorkflowID)
		return cw != newWT
	}, 100*time.Millisecond, 2*time.Millisecond, "the old waiter must not drop the new registration")
	assert.Equal(t, 1, f.limits.inUse(), "the old waiter must not free the new registration's slot")
	assert.True(t, queues.get(t, 0).closed.Load(), "the old waiter closes its own queue")
	assert.False(t, queues.get(t, 1).closed.Load(), "the old waiter must not close the new registration's queue")

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

	c := NewCoordinator(newRegisterDeps(t, capReg, limits.NewGateLimiter(true)), newTestAdmission(), engines, &fakeWorkflowLimits{}, clockwork.NewFakeClock())
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

// TestCoordinator_Admission covers the path between a workflow's readers and its
// executions. Each case holds executions so that it controls exactly which events
// wait on the queue.
func TestCoordinator_Admission(t *testing.T) {
	t.Parallel()
	// singleSlot returns admission limiters that run one execution at a time,
	// and the pool, so the test can see when the admitter waits for the slot.
	singleSlot := func(mutate func(*AdmissionLimiters)) (AdmissionLimiters, *countingPool) {
		pool := &countingPool{ResourcePoolLimiter: limits.GlobalResourcePoolLimiter(1)}
		admission := newTestAdmission()
		admission.ExecutionConcurrency = pool
		if mutate != nil {
			mutate(&admission)
		}
		return admission, pool
	}
	// queueLen reports how many events wait on the registration's queue.
	queueLen := func(t *testing.T, f *coordinatorFixture) int {
		cw, ok := f.c.workflows.get(validWorkflowID)
		require.True(t, ok)
		n, err := cw.queue.Len(t.Context())
		require.NoError(t, err)
		return n
	}
	t.Run("a full queue drops the event and the reader keeps reading", func(t *testing.T) {
		t.Parallel()
		admission, pool := singleSlot(func(a *AdmissionLimiters) { a.NewTriggerEventQueue = newTestQueue(1) })
		wfMetrics, reader := newTestMeterMetrics(t)
		f := newCoordinatorFixture(t, withAdmission(admission), withMetrics(wfMetrics))
		trigger, eventCh := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()
		release := f.holdExecutions(t)
		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)

		// Move one event to each stage, waiting for each before sending the next:
		//   evt-1: executing, holding the only slot.
		//   evt-2: taken off the queue by the admitter, which waits for that slot.
		//   evt-3: in the queue, which is now full (1/1).
		eventCh <- triggerEvent("evt-1")
		require.Eventually(t, func() bool { return len(f.engine.executed()) == 1 }, 5*time.Second, 5*time.Millisecond)
		eventCh <- triggerEvent("evt-2")
		// The admitter's 2nd slot request: it now holds evt-2 and waits.
		require.Eventually(t, func() bool { return pool.waits.Load() == 2 }, 5*time.Second, 5*time.Millisecond)
		eventCh <- triggerEvent("evt-3")
		// The admitter is busy waiting, so evt-3 stays in the queue.
		require.Eventually(t, func() bool { return queueLen(t, f) == 1 }, 5*time.Second, 5*time.Millisecond)

		// The queue is full: evt-4 is dropped. The send only completes because
		// the reader keeps reading. The reader handles one event at a time, so
		// once the probe after it is taken, evt-4 has been dropped. ReadLoop
		// drops the probe itself, as an errored trigger response, without
		// enqueueing it.
		eventCh <- triggerEvent("evt-4")
		eventCh <- capabilities.TriggerResponse{Err: errors.New("probe")}

		release()
		require.Eventually(t, func() bool { return len(f.engine.executed()) == 3 }, 5*time.Second, 5*time.Millisecond)

		// The queue is FIFO and there is one slot, so had evt-4 been queued it
		// would run before evt-5. Once evt-5 has run, the list is final.
		eventCh <- triggerEvent("evt-5")
		require.Eventually(t, func() bool { return len(f.engine.executed()) >= 4 }, 5*time.Second, 5*time.Millisecond)
		assert.Equal(t, []string{"evt-1", "evt-2", "evt-3", "evt-5"}, f.engine.executedIDs())

		// The drop is counted as on the legacy engine (engine.go put).
		queueFull := map[string]string{platform.KeyTriggerDropReason: monitoring.TriggerDropReasonEnqueueQueueFull}
		assert.Equal(t, int64(1), counterTotal(t, reader, "platform_engine_trigger_event_dropped_total", queueFull))
		assert.Equal(t, int64(1), counterTotal(t, reader, "platform_engine_workflow_trigger_event_queue_full", nil))
		assert.Equal(t, int64(1), counterTotal(t, reader, "platform_engine_trigger_event_enqueue_dropped_total", nil))
		assert.Equal(t, int64(4), counterTotal(t, reader, "platform_engine_trigger_event_enqueued_total", nil), "evt-1, evt-2, evt-3 and evt-5")
	})

	unregistered := map[string]string{platform.KeyTriggerDropReason: monitoring.TriggerDropReasonUnregistered}
	t.Run("unregistering counts events not yet admitted as dropped", func(t *testing.T) {
		t.Parallel()
		queues := &queueTracker{}
		admission, pool := singleSlot(func(a *AdmissionLimiters) { a.NewTriggerEventQueue = queues.newQueue })
		wfMetrics, reader := newTestMeterMetrics(t)
		f := newCoordinatorFixture(t, withAdmission(admission), withMetrics(wfMetrics))
		trigger, eventCh := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()
		release := f.holdExecutions(t)
		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)

		// evt-1 executes, evt-2 is held by the admitter waiting for the slot,
		// and evt-3 waits on the queue.
		eventCh <- triggerEvent("evt-1")
		require.Eventually(t, func() bool { return len(f.engine.executed()) == 1 }, 5*time.Second, 5*time.Millisecond)
		eventCh <- triggerEvent("evt-2")
		require.Eventually(t, func() bool { return pool.waits.Load() == 2 }, 5*time.Second, 5*time.Millisecond)
		eventCh <- triggerEvent("evt-3")
		require.Eventually(t, func() bool { return queueLen(t, f) == 1 }, 5*time.Second, 5*time.Millisecond)

		require.NoError(t, f.c.UnregisterTriggers(t.Context(), validWorkflowID))
		release()
		// The queue is closed only after its leftover events are counted.
		require.Eventually(t, func() bool { return queues.get(t, 0).closed.Load() }, 5*time.Second, 5*time.Millisecond)

		assert.Equal(t, []string{"evt-1"}, f.engine.executedIDs())
		assert.Equal(t, int64(2), counterTotal(t, reader, "platform_engine_trigger_event_dropped_total", unregistered), "evt-2 and evt-3")
		assert.Equal(t, int64(2), counterTotal(t, reader, "platform_engine_trigger_event_dequeue_dropped_total", nil))
	})
	t.Run("unregistering with nothing queued drops nothing", func(t *testing.T) {
		t.Parallel()
		admission, queues := newTrackedAdmission()
		wfMetrics, reader := newTestMeterMetrics(t)
		f := newCoordinatorFixture(t, withAdmission(admission), withMetrics(wfMetrics))
		trigger, eventCh := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()
		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)

		eventCh <- triggerEvent("evt-1")
		require.Eventually(t, func() bool { return len(f.engine.executed()) == 1 }, 5*time.Second, 5*time.Millisecond)

		require.NoError(t, f.c.UnregisterTriggers(t.Context(), validWorkflowID))
		require.Eventually(t, func() bool { return queues.get(t, 0).closed.Load() }, 5*time.Second, 5*time.Millisecond)

		assert.Equal(t, int64(0), counterTotal(t, reader, "platform_engine_trigger_event_dropped_total", unregistered))
	})

	t.Run("an event that expires while still queued is dropped and never executed", func(t *testing.T) {
		t.Parallel()
		admission, pool := singleSlot(func(a *AdmissionLimiters) {
			a.TriggerEventQueueTimeout = limits.NewUpperBoundLimiter(time.Minute)
		})
		wfMetrics, reader := newTestMeterMetrics(t)
		f := newCoordinatorFixture(t, withAdmission(admission), withMetrics(wfMetrics))
		trigger, eventCh := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()
		release := f.holdExecutions(t)
		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)

		// evt-1 and evt-2 move off the queue as before, leaving evt-3 queued.
		eventCh <- triggerEvent("evt-1")
		require.Eventually(t, func() bool { return len(f.engine.executed()) == 1 }, 5*time.Second, 5*time.Millisecond)
		eventCh <- triggerEvent("evt-2")
		require.Eventually(t, func() bool { return pool.waits.Load() == 2 }, 5*time.Second, 5*time.Millisecond)
		eventCh <- triggerEvent("evt-3")
		require.Eventually(t, func() bool { return queueLen(t, f) == 1 }, 5*time.Second, 5*time.Millisecond)

		// evt-3's 1-minute deadline passes while it is still in the queue.
		f.clock.Advance(2 * time.Minute)

		release()
		expired := map[string]string{platform.KeyTriggerDropReason: monitoring.TriggerDropReasonExpired}
		require.Eventually(t, func() bool {
			return counterTotal(t, reader, "platform_engine_trigger_event_dropped_total", expired) == 1
		},
			5*time.Second, 5*time.Millisecond, "evt-3 dropped as expired")
		assert.Equal(t, int64(1), counterTotal(t, reader, "platform_engine_trigger_event_expired_total", nil))
		assert.Equal(t, []string{"evt-1", "evt-2"}, f.engine.executedIDs(), "evt-3 must never reach ExecuteTrigger")
	})

	t.Run("the deadline is stamped at dispatch, so a later timeout change does not move it", func(t *testing.T) {
		t.Parallel()
		timeout := &mutableTimeoutLimiter{value: time.Minute}
		admission, pool := singleSlot(func(a *AdmissionLimiters) { a.TriggerEventQueueTimeout = timeout })
		f := newCoordinatorFixture(t, withAdmission(admission))
		trigger, eventCh := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()
		release := f.holdExecutions(t)
		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)

		eventCh <- triggerEvent("evt-1")
		require.Eventually(t, func() bool { return len(f.engine.executed()) == 1 }, 5*time.Second, 5*time.Millisecond)
		eventCh <- triggerEvent("evt-2")
		require.Eventually(t, func() bool { return pool.waits.Load() == 2 }, 5*time.Second, 5*time.Millisecond)

		// evt-3 is enqueued while the setting is still 1 minute.
		eventCh <- triggerEvent("evt-3")
		require.Eventually(t, func() bool { return queueLen(t, f) == 1 }, 5*time.Second, 5*time.Millisecond)

		// Shrink the setting after evt-3 is stamped. If the admitter recomputed the
		// deadline now, evt-3 would look expired.
		timeout.set(time.Millisecond)
		f.clock.Advance(2 * time.Second)

		release()
		require.Eventually(t, func() bool { return len(f.engine.executed()) == 3 }, 5*time.Second, 5*time.Millisecond)
		assert.Equal(t, []string{"evt-1", "evt-2", "evt-3"}, f.engine.executedIDs(), "evt-3 keeps the deadline stamped at dispatch")
	})
}

// TestCoordinator_AvailableCapacity checks that AvailableCapacity reports how many
// executions can still start at each scope: per workflow, per owner, per org, and globally.
func TestCoordinator_AvailableCapacity(t *testing.T) {
	t.Parallel()

	t.Run("reports free capacity at every scope across registered workflows", func(t *testing.T) {
		t.Parallel()
		pools := limits.MultiResourcePoolLimiter[int]{
			limits.GlobalResourcePoolLimiter(3), // workflow
			limits.GlobalResourcePoolLimiter(5), // owner
			limits.GlobalResourcePoolLimiter(7), // org
			limits.GlobalResourcePoolLimiter(9), // global
		}
		admission := newTestAdmission()
		admission.ExecutionConcurrency = pools
		f := newCoordinatorFixture(t, withAdmission(admission))
		trigger, _ := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()

		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)

		got, err := f.c.AvailableCapacity(t.Context())
		require.NoError(t, err)
		assert.Equal(t, 9, got.Global, "global has no tenant key, so it is read without a per-workflow context")
		assert.Equal(t, map[string]int{validWorkflowID: 3}, got.PerWorkflow)
		assert.Equal(t, map[string]int{testOwner: 5}, got.PerOwner)
		assert.Equal(t, map[string]int{testOrg: 7}, got.PerOrg)
	})

	t.Run("a running execution takes one slot at every scope until it finishes", func(t *testing.T) {
		t.Parallel()
		admission := newTestAdmission()
		admission.ExecutionConcurrency = limits.MultiResourcePoolLimiter[int]{
			limits.GlobalResourcePoolLimiter(3), // workflow
			limits.GlobalResourcePoolLimiter(5), // owner
			limits.GlobalResourcePoolLimiter(7), // org
			limits.GlobalResourcePoolLimiter(9), // global
		}
		f := newCoordinatorFixture(t, withAdmission(admission))
		trigger, eventCh := f.expectTrigger(t)
		trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()
		release := f.holdExecutions(t)
		_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
		require.NoError(t, err)

		capacity := func(wf, owner, org, global int) Capacity {
			return Capacity{
				Global:      global,
				PerWorkflow: map[string]int{validWorkflowID: wf},
				PerOwner:    map[string]int{testOwner: owner},
				PerOrg:      map[string]int{testOrg: org},
			}
		}
		requireCapacity := func(want Capacity, msg string) {
			require.EventuallyWithT(t, func(c *assert.CollectT) {
				got, err := f.c.AvailableCapacity(t.Context())
				if assert.NoError(c, err) {
					assert.Equal(c, want, got)
				}
			}, 5*time.Second, 5*time.Millisecond, msg)
		}

		requireCapacity(capacity(3, 5, 7, 9), "registering takes no execution slot")

		eventCh <- triggerEvent("evt-1")
		require.Eventually(t, func() bool { return len(f.engine.executed()) == 1 }, 5*time.Second, 5*time.Millisecond)
		requireCapacity(capacity(2, 4, 6, 8), "the running execution holds one slot at every scope")

		release()
		requireCapacity(capacity(3, 5, 7, 9), "the slot is given back at every scope when the execution finishes")
	})

	t.Run("reads each owner and org pool once, however many workflows share it", func(t *testing.T) {
		t.Parallel()
		wf, owner, org, global := newCountingPool(3), newCountingPool(5), newCountingPool(7), newCountingPool(9)
		c := &coordinator{
			admission: AdmissionLimiters{ExecutionConcurrency: limits.MultiResourcePoolLimiter[int]{wf, owner, org, global}},
			workflows: newWorkflowRegistry(),
		}
		// Three workflows of one owner in one org.
		for _, id := range []string{"wf-a", "wf-b", "wf-c"} {
			c.workflows.set(id, &coordinatedWorkflow{cre: contexts.CRE{Org: testOrg, Owner: testOwner, Workflow: id}})
		}

		got, err := c.AvailableCapacity(t.Context())
		require.NoError(t, err)
		assert.Len(t, got.PerWorkflow, 3)
		assert.Equal(t, int32(3), wf.availables.Load(), "one read per workflow")
		assert.Equal(t, int32(1), owner.availables.Load(), "the shared owner is read once")
		assert.Equal(t, int32(1), org.availables.Load(), "the shared org is read once")
		assert.Equal(t, int32(1), global.availables.Load())
	})

	t.Run("returns an error when the shared limiter does not break down by scope", func(t *testing.T) {
		t.Parallel()
		// The default fixture uses one pool with no scopes, so there is nothing per scope
		// to report. AvailableCapacity must return an error, not empty or wrong numbers.
		f := newCoordinatorFixture(t)

		_, err := f.c.AvailableCapacity(t.Context())
		require.ErrorIs(t, err, ErrCapacityNotSupported)
	})
}

// BenchmarkCoordinator_AvailableCapacity measures how long one AvailableCapacity
// call takes on a node with 1,000 workflows, 100 owners and 10 orgs. The pools are
// built like the production ones, so each read does a real settings lookup.
//
//	go test -run '^$' -bench AvailableCapacity -benchmem ./core/services/workflows/v2/triggers/
func BenchmarkCoordinator_AvailableCapacity(b *testing.B) {
	const workflows, owners, orgs = 1000, 100, 10

	getter, err := settings.NewJSONGetter([]byte(`{}`))
	require.NoError(b, err)
	lf := limits.Factory{Logger: logger.Nop(), Settings: getter}
	scopeLimits := []settings.Setting[int]{
		cresettings.Default.PerWorkflow.ExecutionConcurrencyLimit,
		cresettings.Default.PerOwner.WorkflowExecutionConcurrencyLimit,
		cresettings.Default.PerOrg.WorkflowExecutionConcurrencyLimit,
		cresettings.Default.WorkflowExecutionConcurrencyLimit,
	}
	pools := make(limits.MultiResourcePoolLimiter[int], 0, len(scopeLimits))
	for _, setting := range scopeLimits {
		pool, makeErr := limits.MakeResourcePoolLimiter(lf, setting)
		require.NoError(b, makeErr)
		b.Cleanup(func() { _ = pool.Close() })
		pools = append(pools, pool)
	}

	c := &coordinator{admission: AdmissionLimiters{ExecutionConcurrency: pools}, workflows: newWorkflowRegistry()}
	for i := range workflows {
		id := fmt.Sprintf("%064x", i)
		c.workflows.set(id, &coordinatedWorkflow{cre: contexts.CRE{
			Org:      fmt.Sprintf("org-%d", i%orgs),
			Owner:    fmt.Sprintf("%040x", i%owners),
			Workflow: id,
		}})
	}

	ctx := b.Context()
	// The first read of each tenant creates its pool entry; measure steady state.
	_, err = c.AvailableCapacity(ctx)
	require.NoError(b, err)

	for b.Loop() {
		if _, err := c.AvailableCapacity(ctx); err != nil {
			b.Fatal(err)
		}
	}
}

// TestCoordinator_DropObservability covers drops that happen while resolving
// the engine, and makes sure expected execution outcomes (duplicate,
// draining, closed) are never logged as errors. The drop-counting side of
// this is covered by TestCoordinator_NoDeliver.
func TestCoordinator_DropObservability(t *testing.T) {
	t.Parallel()

	t.Run("expected execution outcomes are logged at info, not error", func(t *testing.T) {
		t.Parallel()
		for _, tc := range []struct {
			name string
			err  error
		}{
			{"duplicate execution", ErrDuplicateExecution},
			{"engine draining", ErrEngineDraining},
			{"engine closed", ErrEngineClosed},
			{"admission cached", ErrAdmissionCached},
			{"admission denied", fmt.Errorf("not owner: %w", ErrAdmissionDenied)},
		} {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				lggr, errorLogs := logger.TestObserved(t, zapcore.ErrorLevel)
				f := newCoordinatorFixture(t, withLogger(lggr))
				trigger, eventCh := f.expectTrigger(t)
				trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()
				trigger.EXPECT().AckEvent(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
				f.engine.err = tc.err
				_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
				require.NoError(t, err)

				eventCh <- triggerEvent("evt-1")
				require.Eventually(t, func() bool { return len(f.engine.executed()) == 1 }, 5*time.Second, 5*time.Millisecond)

				assert.Empty(t, errorLogs.All(), "an expected execution outcome must not log at error, the way legacy's engine.go does")
			})
		}

		t.Run("an unexpected execution error still logs at error", func(t *testing.T) {
			t.Parallel()
			lggr, errorLogs := logger.TestObserved(t, zapcore.ErrorLevel)
			f := newCoordinatorFixture(t, withLogger(lggr))
			trigger, eventCh := f.expectTrigger(t)
			trigger.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Maybe()
			f.engine.err = errors.New("boom")
			_, err := f.c.RegisterTriggers(t.Context(), newTestSubscriber(), testParams(t))
			require.NoError(t, err)

			eventCh <- triggerEvent("evt-1")
			require.Eventually(t, func() bool { return len(errorLogs.All()) == 1 }, 5*time.Second, 5*time.Millisecond,
				"an unanticipated execution error must still be visible as an error")
		})
	})
}
