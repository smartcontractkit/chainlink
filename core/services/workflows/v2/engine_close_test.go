package v2_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

func legacyCloseTestConfig(t *testing.T) (*v2.EngineConfig, *regmocks.CapabilitiesRegistry) {
	t.Helper()
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)
	cfg := defaultTestConfig(t, nil)
	cfg.CapRegistry = capreg
	return cfg, capreg
}

// TestEngine_CloseCancelsAndWaitsForItsOwnExecution covers an execution the
// legacy engine starts itself, from its trigger queue: Close cancels it and
// waits for it before closing the module.
func TestEngine_CloseCancelsAndWaitsForItsOwnExecution(t *testing.T) {
	t.Parallel()

	cfg, capreg := legacyCloseTestConfig(t)
	cfg.BillingClient = setupMockBillingClient(t)
	trigger := capmocks.NewTriggerCapability(t)
	eventCh := make(chan capabilities.TriggerResponse)
	trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(eventCh, nil).Once()
	trigger.EXPECT().AckEvent(matches.AnyContext, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()
	capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger, nil).Once()

	exec := newBlockingExecution()
	re := newTestEngine(t, cfg, v2.NewEngine, func(module *modulemocks.ModuleV2) {
		module.EXPECT().Start().Once()
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(1), nil).Once()
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).RunAndReturn(exec.execute).Once()
		exec.expectModuleClose(module)
	})
	require.NoError(t, re.engine.Start(t.Context()))
	require.NoError(t, <-re.initializedCh)
	require.Equal(t, []string{"id_0"}, <-re.subscribedToTriggersCh)

	eventCh <- capabilities.TriggerResponse{Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "own_event"}}
	<-exec.started

	closeWithin(t, re.engine, 10*time.Second)
	assert.True(t, exec.cancelled.Load(), "Close must cancel an execution the engine started")
	assert.False(t, exec.moduleClosedEarly.Load(), "the module was closed while an execution was still running on it")
}

// TestEngine_CloseCancelsAndWaitsForExternallyInvokedExecution covers an
// execution invoked through ExecuteTrigger from a caller's goroutine, as
// ShardFailoverManager replays are (on context.Background, from the shard
// communicator): Close cancels it and waits for it, rather than blocking until
// it finishes on its own.
func TestEngine_CloseCancelsAndWaitsForExternallyInvokedExecution(t *testing.T) {
	t.Parallel()

	cfg, _ := legacyCloseTestConfig(t)
	cfg.BillingClient = setupMockBillingClient(t)
	exec := newBlockingExecution()
	re := newTestEngine(t, cfg, v2.NewEngine, func(module *modulemocks.ModuleV2) {
		module.EXPECT().Start().Once()
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(0), nil).Once()
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).RunAndReturn(exec.execute).Once()
		exec.expectModuleClose(module)
	})
	require.NoError(t, re.engine.Start(t.Context()))
	require.NoError(t, <-re.initializedCh)

	_, event := coordinatedEvent(t, cfg, "replayed_event")
	replayCtx := contexts.WithCRE(context.Background(), contexts.CRE{Owner: cfg.WorkflowOwner, Workflow: cfg.WorkflowID})
	execDone := make(chan struct{})
	go func() {
		defer close(execDone)
		_ = re.engine.ExecuteTrigger(replayCtx, event)
	}()
	<-exec.started

	closeWithin(t, re.engine, 10*time.Second)
	assert.True(t, exec.cancelled.Load(), "Close must cancel an execution running on a caller's goroutine")
	assert.False(t, exec.moduleClosedEarly.Load(), "the module was closed while an execution was still running on it")
	<-execDone
}

func TestEngine_ExecuteTriggerAfterCloseIsRejected(t *testing.T) {
	t.Parallel()

	cfg, _ := legacyCloseTestConfig(t)
	// No execution expectation: a rejected call must never reach the module.
	re := newTestEngine(t, cfg, v2.NewEngine, func(module *modulemocks.ModuleV2) {
		module.EXPECT().Start().Once()
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(0), nil).Once()
		module.EXPECT().Close().Once()
	})
	require.NoError(t, re.engine.Start(t.Context()))
	require.NoError(t, <-re.initializedCh)
	require.NoError(t, re.engine.Close())

	ctx, event := coordinatedEvent(t, cfg, "after_close_event")
	require.ErrorIs(t, re.engine.ExecuteTrigger(ctx, event), v2.ErrEngineClosed)
}
