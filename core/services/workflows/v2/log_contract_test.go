package v2_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"

	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"

	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

// TestEngine_ExecutionLogContract pins the engine log lines that external
// integrity detectors (and the fault-injection reconciler) parse from node
// logs. The message text and the "executionID" structured field key are a
// CONTRACT: renaming either silently blinds the detector, so these assertions
// must break CI first.
//
// Contract lines:
//   - "Workflow execution starting ..."
//   - "Workflow execution finished successfully"
//   - "Workflow execution failed"
//
// each carrying a structured "executionID" field (attached via the
// execution-scoped logger in engine.go).
func TestEngine_ExecutionLogContract(t *testing.T) {
	t.Parallel()

	t.Run("successful execution", func(t *testing.T) {
		t.Parallel()
		logs, wantExecID := runExecutionForLogContract(t, nil, nil)
		requireContractLog(t, logs, "Workflow execution starting ...", wantExecID)
		requireContractLog(t, logs, "Workflow execution finished successfully", wantExecID)
	})

	t.Run("failed execution", func(t *testing.T) {
		t.Parallel()
		userErrResult := &sdkpb.ExecutionResult{
			Result: &sdkpb.ExecutionResult_Error{Error: "user workflow error"},
		}
		logs, wantExecID := runExecutionForLogContract(t, userErrResult, nil)
		requireContractLog(t, logs, "Workflow execution starting ...", wantExecID)
		requireContractLog(t, logs, "Workflow execution failed", wantExecID)
	})
}

// runExecutionForLogContract runs one engine execution end-to-end with an
// observed logger and returns the captured logs plus the expected execution ID.
// execResult/execErr are what the module returns for the triggered execution
// (nil/nil = success).
func runExecutionForLogContract(t *testing.T, execResult *sdkpb.ExecutionResult, execErr error) (*observer.ObservedLogs, string) {
	t.Helper()

	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start()
	module.EXPECT().Close()
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)
	billingClient := setupMockBillingClient(t)

	initDoneCh := make(chan error)
	subscribedToTriggersCh := make(chan []string, 1)
	executionFinishedCh := make(chan string)

	cfg := defaultTestConfig(t, nil)
	var logs *observer.ObservedLogs
	cfg.Lggr, logs = logger.TestObserved(t, zapcore.InfoLevel)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = billingClient
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized:          func(err error) { initDoneCh <- err },
		OnSubscribedToTriggers: func(triggerIDs []string) { subscribedToTriggersCh <- triggerIDs },
		OnExecutionFinished:    func(executionID string, _ string) { executionFinishedCh <- executionID },
	}

	engine, err := v2.NewEngine(cfg)
	require.NoError(t, err)

	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(1), nil).Once()
	trigger := capmocks.NewTriggerCapability(t)
	capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger, nil)
	eventCh := make(chan capabilities.TriggerResponse)
	trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(eventCh, nil).Once()
	trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()
	trigger.EXPECT().AckEvent(matches.AnyContext, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(execResult, execErr).Once()

	require.NoError(t, engine.Start(t.Context()))
	require.NoError(t, <-initDoneCh)
	require.Equal(t, []string{"id_0"}, <-subscribedToTriggersCh)

	mockTriggerEvent := capabilities.TriggerEvent{
		TriggerType: "basic-trigger@1.0.0",
		ID:          "log_contract_event",
		Payload:     nil,
	}
	eventCh <- capabilities.TriggerResponse{Event: mockTriggerEvent}
	gotExecID := <-executionFinishedCh

	require.NoError(t, engine.Close())

	wantExecID := wantExecutionID(t, cfg.WorkflowID, mockTriggerEvent.ID, 0)
	require.Equal(t, wantExecID, gotExecID)
	return logs, wantExecID
}

// requireContractLog asserts that a log entry with exactly msg was emitted and
// carries the structured field key "executionID" with the expected value.
func requireContractLog(t *testing.T, logs *observer.ObservedLogs, msg, wantExecutionID string) {
	t.Helper()

	entries := logs.FilterMessage(msg).All()
	require.NotEmpty(t, entries, "contract log line %q was not emitted — external detectors parse this exact message", msg)

	for _, entry := range entries {
		got, ok := entry.ContextMap()["executionID"]
		require.True(t, ok, "contract log line %q is missing the structured field key %q — external detectors parse this exact key", msg, "executionID")
		require.Equal(t, wantExecutionID, got, "contract log line %q carries an unexpected executionID", msg)
	}
}
