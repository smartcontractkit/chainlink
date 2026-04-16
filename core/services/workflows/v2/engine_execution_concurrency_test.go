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
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"

	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	workflowEvents "github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
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

	wantExecID1, err := workflowEvents.GenerateExecutionID(cfg.WorkflowID, "event_concurrency_1")
	require.NoError(t, err)
	wantExecID2, err := workflowEvents.GenerateExecutionID(cfg.WorkflowID, "event_concurrency_2")
	require.NoError(t, err)

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

	engine, err := v2.NewEngine(cfg)
	require.NoError(t, err)

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

	for i := 0; i < 10_000; i++ {
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
		id, err := workflowEvents.GenerateExecutionID(cfg.WorkflowID, eid)
		require.NoError(t, err)
		wantExecIDs[id] = struct{}{}
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

	engine, err := v2.NewEngine(cfg)
	require.NoError(t, err)

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
