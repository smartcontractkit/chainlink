package triggers

import (
	"context"
	"errors"
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
		f := newCoordinatorFixtureWithAdmission(t, admission)
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
		f := newCoordinatorFixtureWithAdmission(t, admission)
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

	t.Run("registration state is gone, the slot is freed and the queue is closed when no execution is in flight", func(t *testing.T) {
		t.Parallel()
		admission, queues := newTrackedAdmission()
		f := newCoordinatorFixtureWithAdmission(t, admission)
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

		// The in-flight execution still resolves its handle.
		require.NoError(t, f.c.Ack(t.Context(), testTriggerCapID, regID, "evt-1"))
		assert.True(t, f.registered())
		assert.Equal(t, 1, f.limits.inUse())

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
	f := newCoordinatorFixtureWithAdmission(t, admission)

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
//
//nolint:paralleltest // its cases read metrics through useTestMeter, which swaps a global
func TestCoordinator_Admission(t *testing.T) {
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
	queueLen := func(f *coordinatorFixture) int {
		cw, ok := f.c.workflows.get(validWorkflowID)
		require.True(t, ok)
		n, err := cw.queue.Len(t.Context())
		require.NoError(t, err)
		return n
	}
	t.Run("a full queue drops the event and the reader keeps reading", func(t *testing.T) { //nolint:paralleltest // see above
		reader := useTestMeter(t) // before the fixture, so its metrics are read back
		admission, pool := singleSlot(func(a *AdmissionLimiters) { a.NewTriggerEventQueue = newTestQueue(1) })
		f := newCoordinatorFixtureWithAdmission(t, admission)
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
		require.Eventually(t, func() bool { return queueLen(f) == 1 }, 5*time.Second, 5*time.Millisecond)

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
}
