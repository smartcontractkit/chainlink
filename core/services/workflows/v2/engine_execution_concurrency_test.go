package v2_test

import (
	"context"
	"runtime"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
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
