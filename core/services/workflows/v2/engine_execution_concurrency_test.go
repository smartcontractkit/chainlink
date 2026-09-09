package v2_test

import (
	"context"
	"fmt"
	"runtime"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

// TestEngine_ExecutionConcurrencySerializesOverlappingRuns proves that when PerWorkflow
// ExecutionConcurrencyLimit is 1, a second trigger cannot start Module.Execute until the first
// run completes (executionsSemaphore.Wait blocks in handleAllTriggerEvents).
func TestEngine_ExecutionConcurrencySerializesOverlappingRuns(t *testing.T) {
	t.Parallel()

	continueFirst := make(chan struct{})
	var execMu sync.Mutex
	var execOrder []string

	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start().Once()
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(1), nil).Once()
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Run(
		func(_ context.Context, _ *sdkpb.ExecuteRequest, eh host.ExecutionHelper) {
			execMu.Lock()
			execOrder = append(execOrder, eh.GetWorkflowExecutionID())
			n := len(execOrder)
			execMu.Unlock()
			if n == 1 {
				<-continueFirst
			}
		}).Return(nil, nil).Times(2)
	module.EXPECT().Close().Once()

	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil).Once()

	initDoneCh := make(chan error, 1)
	subscribedToTriggersCh := make(chan []string, 1)
	executionFinishedCh := make(chan string, 2)

	cfg := defaultTestConfig(t, func(cfg *cresettings.Workflows) {
		cfg.ExecutionConcurrencyLimit.DefaultValue = 1
	})
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = setupMockBillingClient(t)

	wantExecID1 := wantExecutionID(t, cfg.WorkflowID, "event_concurrency_1", 0)
	wantExecID2 := wantExecutionID(t, cfg.WorkflowID, "event_concurrency_2", 0)

	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
		OnSubscribedToTriggers: func(triggerIDs []string) {
			subscribedToTriggersCh <- triggerIDs
		},
		OnExecutionFinished: func(executionID string, _ string) {
			executionFinishedCh <- executionID
			if executionID == wantExecID2 {
				close(executionFinishedCh)
			}
		},
	}

	engine := newDispatchedEngine(t, cfg)

	trigger := capmocks.NewTriggerCapability(t)
	capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger, nil).Once()
	eventCh := make(chan capabilities.TriggerResponse)
	trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(eventCh, nil).Once()
	trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()
	trigger.EXPECT().AckEvent(matches.AnyContext, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	require.NoError(t, engine.Start(t.Context()))
	require.NoError(t, <-initDoneCh)
	require.Equal(t, []string{"id_0"}, <-subscribedToTriggersCh)

	eventCh <- capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{
			TriggerType: "basic-trigger@1.0.0",
			ID:          "event_concurrency_1",
			Payload:     nil,
		},
	}

	require.Eventually(t, func() bool {
		execMu.Lock()
		defer execMu.Unlock()
		return len(execOrder) == 1 && execOrder[0] == wantExecID1
	}, 2*time.Second, 5*time.Millisecond, "first execution should start")

	eventCh <- capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{
			TriggerType: "basic-trigger@1.0.0",
			ID:          "event_concurrency_2",
			Payload:     nil,
		},
	}

	for range 10_000 {
		runtime.Gosched()
	}
	execMu.Lock()
	gotMid := slices.Clone(execOrder)
	execMu.Unlock()
	require.Equal(t, []string{wantExecID1}, gotMid,
		"second execution must not start while the first holds the executions semaphore")

	continueFirst <- struct{}{}

	require.Eventually(t, func() bool {
		execMu.Lock()
		defer execMu.Unlock()
		return slices.Equal(execOrder, []string{wantExecID1, wantExecID2})
	}, 2*time.Second, 5*time.Millisecond, "second execution should start after the first completes")

	finishedIDs := make([]string, 0, 2)
	for id := range executionFinishedCh {
		finishedIDs = append(finishedIDs, id)
	}
	require.Equal(t, []string{wantExecID1, wantExecID2}, finishedIDs)

	require.NoError(t, engine.Close())
}

// TestEngine_StaleTriggerEventIsSkipped proves that trigger events older than
// TriggerEventQueueTimeout are dropped and never reach Module.Execute.
//
// Strategy: with ExecutionConcurrencyLimit=1 and a FakeClock we can control
// exactly which events age out. We send 7 "early" events. event_0's Execute
// blocks (holding the semaphore), event_1 gets popped by the handler but
// stalls at the semaphore (already past its age check), and events 2-6 remain
// in the queue. After advancing the fake clock past the timeout, we unblock
// event_0. event_1 resumes and executes (it already cleared the age check),
// while events 2-6 are popped and detected as stale (5 skipped). Then 3
// fresh events are sent and all execute. Total: 10 events, 5 expire, 5
// execute.
func TestEngine_StaleTriggerEventIsSkipped(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()

	const queueTimeout = 5 * time.Second

	fakeClock := clockwork.NewFakeClock()
	blockerStarted := make(chan struct{})
	blockerRelease := make(chan struct{})

	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start().Once()
	// init call → return trigger subscriptions
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		Return(newTriggerSubs(1), nil).Once()
	// event_0 → block until we release it (holds the semaphore)
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		Run(func(_ context.Context, _ *sdkpb.ExecuteRequest, _ host.ExecutionHelper) {
			close(blockerStarted)
			<-blockerRelease
		}).Return(nil, nil).Once()
	// event_1 + 3 fresh events = 4 fast executions
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		Return(nil, nil).Times(4)
	module.EXPECT().Close().Once()

	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil).Once()

	initDoneCh := make(chan error, 1)
	subscribedToTriggersCh := make(chan []string, 1)
	executionFinishedCh := make(chan string, 10)

	var lggr logger.Logger
	var logs *observer.ObservedLogs

	cfg := defaultTestConfig(t, func(cfg *cresettings.Workflows) {
		cfg.TriggerEventQueueTimeout.DefaultValue = queueTimeout
		cfg.ExecutionConcurrencyLimit.DefaultValue = 1
	})
	lggr, logs = logger.TestObserved(t, zapcore.WarnLevel)
	cfg.Lggr = lggr
	cfg.Clock = fakeClock
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = setupMockBillingClient(t)

	wantExecIDs := make(map[string]struct{}, 5)
	for _, eid := range []string{"event_0", "event_1", "fresh_0", "fresh_1", "fresh_2"} {
		wantExecIDs[wantExecutionID(t, cfg.WorkflowID, eid, 0)] = struct{}{}
	}

	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
		OnSubscribedToTriggers: func(triggerIDs []string) {
			subscribedToTriggersCh <- triggerIDs
		},
		OnExecutionFinished: func(executionID string, _ string) {
			executionFinishedCh <- executionID
		},
	}

	engine := newDispatchedEngine(t, cfg)

	trigger := capmocks.NewTriggerCapability(t)
	capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger, nil).Once()
	eventCh := make(chan capabilities.TriggerResponse)
	trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(eventCh, nil).Once()
	trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()
	trigger.EXPECT().AckEvent(matches.AnyContext, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	require.NoError(t, engine.Start(t.Context()))
	require.NoError(t, <-initDoneCh)
	require.Equal(t, []string{"id_0"}, <-subscribedToTriggersCh)

	// Send 7 events that will be timestamped at the current fake clock time.
	// event_0 will execute and block; event_1 will be popped but stall at the
	// semaphore; events 2-6 will sit in the queue.
	for i := range 7 {
		eventCh <- capabilities.TriggerResponse{
			Event: capabilities.TriggerEvent{
				TriggerType: "basic-trigger@1.0.0",
				ID:          fmt.Sprintf("event_%d", i),
			},
		}
	}

	// Wait for event_0's Execute to start (semaphore is now held).
	<-blockerStarted
	// Give the handler goroutine time to pop event_1 and block at the semaphore.
	time.Sleep(200 * time.Millisecond)

	// Advance the fake clock so events still in the queue become stale.
	fakeClock.Advance(queueTimeout + time.Second)

	// Unblock event_0 → event_1 will resume (already past age check),
	// events 2-6 will be detected as too old.
	close(blockerRelease)

	// Send 3 fresh events — timestamped at the advanced clock time.
	for i := range 3 {
		eventCh <- capabilities.TriggerResponse{
			Event: capabilities.TriggerEvent{
				TriggerType: "basic-trigger@1.0.0",
				ID:          fmt.Sprintf("fresh_%d", i),
			},
		}
	}

	gotIDs := make(map[string]struct{}, 5)
	for range 5 {
		gotIDs[<-executionFinishedCh] = struct{}{}
	}
	require.Equal(t, wantExecIDs, gotIDs,
		"expected exactly 5 executions: event_0, event_1, and 3 fresh events")

	require.Eventually(t, func() bool {
		return logs.FilterMessage("Trigger event is too old, skipping execution").Len() >= 5
	}, 2*time.Second, 50*time.Millisecond,
		"expected 5 stale-event warnings for events 2-6")

	require.NoError(t, engine.Close())
}

// mutableSettingsGetter is a settings.Getter whose values can be edited mid-test.
// It ignores scope/tenant resolution and answers by bare setting key, which is
// sufficient to drive a single setting (TriggerEventQueueTimeout) for one workflow.
type mutableSettingsGetter struct {
	mu     sync.Mutex
	values map[string]string
}

var _ settings.Getter = (*mutableSettingsGetter)(nil)

func (m *mutableSettingsGetter) GetScoped(_ context.Context, _ settings.Scope, key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.values[key], nil
}

func (m *mutableSettingsGetter) set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[key] = value
}

// TestEngine_TriggerEventDeadlineUsesDispatchTimeSetting proves that each event
// keeps the expiry it was given when it entered the queue. Changing the
// TriggerEventQueueTimeout setting later does not retroactively expire events
// that are already queued — only events sent after the change get the new
// timeout.
//
// Flow (FakeClock, ExecutionConcurrencyLimit=1 so executions run one at a time):
//  1. Send E0/E1/E2 while the timeout is 10s. Each gets Deadline = now+10s.
//     E0 starts executing and blocks; E1 waits for the execution slot; E2 waits
//     in the queue.
//  2. Change the timeout to 1s and advance the clock 2s. The clock is now past
//     the new 1s timeout, but still before the queued events' 10s deadlines.
//  3. Let E0 finish. E1 executes, then E2 is pulled from the queue and must
//     still EXECUTE — its 10s deadline has not passed. (The old code compared
//     the event's age against the current setting and would have dropped it:
//     2s > 1s.)
//  4. E1 is still executing, so send E3 — it gets Deadline = now+1s from the
//     new setting. Advance the clock past that, let E1 finish: E2 executes and
//     E3 EXPIRES, showing the new setting does apply to newly sent events.
//
// Note on running this test against the pre-deadline code (which compared
// event age against the current setting): it fails the execution-set
// assertion with E2 missing and E3 present. E2 expires because its 2s age
// exceeds the shrunk 1s setting; E3 then executes because the handler, freed
// by E2's early drop, pops it fresh off the queue while its age is still
// under 1s. E0 and E1 execute in both versions — E0 ran before the change,
// and E1 had already passed its age check when it was popped.
func TestEngine_TriggerEventDeadlineUsesDispatchTimeSetting(t *testing.T) {
	t.Parallel()

	const (
		longTimeout  = 10 * time.Second
		shortTimeout = 1 * time.Second
	)

	fakeClock := clockwork.NewFakeClock()
	blocker1Started := make(chan struct{})
	blocker1Release := make(chan struct{})
	blocker2Started := make(chan struct{})
	blocker2Release := make(chan struct{})

	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start().Once()
	// init call → trigger subscriptions
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		Return(newTriggerSubs(1), nil).Once()
	// Executions acquire the semaphore in pop order, so Execute calls are
	// serialized: E0 and E1 block in turn (each holding the semaphore), E2 runs
	// fast. E3 must never reach Module.Execute (asserted via execution IDs).
	var execMu sync.Mutex
	execCount := 0
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		Run(func(_ context.Context, _ *sdkpb.ExecuteRequest, _ host.ExecutionHelper) {
			execMu.Lock()
			execCount++
			n := execCount
			execMu.Unlock()
			switch n {
			case 1:
				close(blocker1Started)
				<-blocker1Release
			case 2:
				close(blocker2Started)
				<-blocker2Release
			}
		}).Return(nil, nil)
	module.EXPECT().Close().Once()

	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil).Once()

	initDoneCh := make(chan error, 1)
	subscribedToTriggersCh := make(chan []string, 1)
	executionFinishedCh := make(chan string, 3)

	var lggr logger.Logger
	var logs *observer.ObservedLogs

	limitCfgFn := func(cfg *cresettings.Workflows) {
		cfg.ExecutionConcurrencyLimit.DefaultValue = 1
	}
	cfg := defaultTestConfig(t, limitCfgFn)
	lggr, logs = logger.TestObserved(t, zapcore.DebugLevel)
	cfg.Lggr = lggr
	cfg.Clock = fakeClock
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = setupMockBillingClient(t)

	// Rebuild the engine limiters on a factory wired to a mutable settings
	// getter so TriggerEventQueueTimeout can be edited mid-test. Every other
	// setting falls back to its default (the getter returns "" for it).
	queueTimeoutKey := cresettings.Default.PerWorkflow.TriggerEventQueueTimeout.Key
	getter := &mutableSettingsGetter{values: map[string]string{
		queueTimeoutKey: longTimeout.String(),
	}}
	mutableLimiters, err := v2.NewLimiters(limits.Factory{
		Logger:   logger.Test(t),
		Settings: getter,
	}, limitCfgFn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, mutableLimiters.Close()) })
	cfg.LocalLimiters = mutableLimiters

	// Map each test event to the execution ID the engine generates for it, and
	// keep the reverse mapping so the assertions below report event names
	// (event_long_0) instead of opaque execution-ID hashes.
	longEventIDs := []string{"event_long_0", "event_long_1", "event_long_2"}
	shortEventID := "event_short_3"
	eventByExecID := make(map[string]string, len(longEventIDs)+1)
	for _, eid := range append(slices.Clone(longEventIDs), shortEventID) {
		eventByExecID[wantExecutionID(t, cfg.WorkflowID, eid, 0)] = eid
	}
	wantEvents := make(map[string]struct{}, len(longEventIDs))
	for _, eid := range longEventIDs {
		wantEvents[eid] = struct{}{}
	}

	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
		OnSubscribedToTriggers: func(triggerIDs []string) {
			subscribedToTriggersCh <- triggerIDs
		},
		OnExecutionFinished: func(executionID string, _ string) {
			executionFinishedCh <- executionID
		},
	}

	engine := newDispatchedEngine(t, cfg)

	trigger := capmocks.NewTriggerCapability(t)
	capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger, nil).Once()
	eventCh := make(chan capabilities.TriggerResponse)
	trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(eventCh, nil).Once()
	trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()
	trigger.EXPECT().AckEvent(matches.AnyContext, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	// Tear the engine down even when an assertion fails mid-test, so the test
	// fails fast instead of hanging: a running engine holds limiter tenants and
	// the limiter cleanups (registered earlier) would block on them.
	// Registered AFTER the trigger mock above so LIFO runs it BEFORE the mock's
	// own AssertExpectations cleanup — engine.Close() calls UnregisterTrigger,
	// which the mock expects.
	t.Cleanup(func() { require.NoError(t, engine.Close()) })

	require.NoError(t, engine.Start(t.Context()))
	require.NoError(t, <-initDoneCh)
	require.Equal(t, []string{"id_0"}, <-subscribedToTriggersCh)

	sendEvent := func(eventID string) {
		eventCh <- capabilities.TriggerResponse{
			Event: capabilities.TriggerEvent{
				TriggerType: "basic-trigger@1.0.0",
				ID:          eventID,
			},
		}
	}

	// Release stalled executions on failure paths so engine.Close() can complete.
	var releaseBlocker1, releaseBlocker2 sync.Once
	t.Cleanup(func() {
		releaseBlocker1.Do(func() { close(blocker1Release) })
		releaseBlocker2.Do(func() { close(blocker2Release) })
	})

	// Phase 1: dispatch three events under the 10s timeout. E0's execution
	// blocks holding the semaphore; E1 stalls at the semaphore; E2 sits in the
	// queue. All three deadlines are stamped at dispatch: now+10s.
	for _, eid := range longEventIDs {
		sendEvent(eid)
	}
	// All three Puts must complete (deadlines stamped) before editing the setting.
	require.Eventually(t, func() bool {
		return logs.FilterMessage("Enqueued trigger event").Len() >= 3
	}, 2*time.Second, 50*time.Millisecond, "expected 3 events enqueued before editing the setting")
	<-blocker1Started

	// Edit the setting while events are queued: 10s → 1s. Advance the clock
	// past the NEW setting (2s > 1s) but before the dispatched deadlines.
	getter.set(queueTimeoutKey, shortTimeout.String())
	fakeClock.Advance(2 * time.Second)

	// Release E0: E1 executes, then E2 is dequeued. E2's deadline (stamped at
	// dispatch under the 10s timeout) is still in the future, so it must
	// EXECUTE — the previous age-vs-current-setting check would have dropped it.
	releaseBlocker1.Do(func() { close(blocker1Release) })
	<-blocker2Started

	// Phase 2: E1's execution still holds the semaphore, so E3 — dispatched
	// under the 1s setting (deadline now+1s) — sits in the queue. Advance past
	// E3's deadline, then release E1: E2 executes and E3 must EXPIRE.
	sendEvent("event_short_3")
	require.Eventually(t, func() bool {
		return logs.FilterMessage("Enqueued trigger event").Len() >= 4
	}, 2*time.Second, 50*time.Millisecond, "expected E3 enqueued before advancing the clock")
	fakeClock.Advance(2 * time.Second)
	releaseBlocker2.Do(func() { close(blocker2Release) })

	gotEvents := make(map[string]struct{}, 3)
	for range 3 {
		execID := <-executionFinishedCh
		eid, ok := eventByExecID[execID]
		if !ok {
			t.Fatalf("execution finished with unknown execution ID %s", execID)
		}
		gotEvents[eid] = struct{}{}
	}
	require.Equal(t, wantEvents, gotEvents,
		"events sent while the timeout was 10s must all execute despite the setting shrinking to 1s")

	require.Eventually(t, func() bool {
		return logs.FilterMessage("Trigger event is too old, skipping execution").Len() >= 1
	}, 2*time.Second, 50*time.Millisecond,
		"expected E3, dispatched under the 1s setting, to expire after its deadline")

	// E3 must never execute.
	select {
	case execID := <-executionFinishedCh:
		eid, _ := eventByExecID[execID]
		t.Fatalf("unexpected execution of expired event %s (%s)", eid, execID)
	case <-time.After(200 * time.Millisecond):
	}
}
