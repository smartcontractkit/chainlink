package v2_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

type noopAcknowledger struct{}

func (noopAcknowledger) Ack(_ context.Context, _, _, _ string) error { return nil }

func TestNewCoordinatedEngine_RequiresAcknowledger(t *testing.T) {
	t.Parallel()

	cfg := defaultTestConfig(t, nil)
	cfg.TriggerAcknowledger = nil

	_, err := v2.NewCoordinatedEngine(cfg)
	require.EqualError(t, err, "trigger acknowledger not set")
}

// TestCoordinatedEngine_ExecuteTrigger drives an external trigger event straight
// into ExecuteTrigger and asserts the execution runs to completion. Unlike the legacy Engine,
// a coordinatedEngine registers no triggers of its own, so this is the only way an event reaches it.
func TestCoordinatedEngine_ExecuteTrigger(t *testing.T) {
	t.Parallel()

	baseCfg := coordinatedTestConfig(t)
	baseCfg.BillingClient = setupMockBillingClient(t)

	// The execution must reach WASM and return a value.
	re := newTestEngine(t, baseCfg, v2.NewCoordinatedEngine, func(module *modulemocks.ModuleV2) {
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
			Return(&sdkpb.ExecutionResult{
				Result: &sdkpb.ExecutionResult_Value{},
			}, nil).
			Once()
	})

	ctx, event := coordinatedEvent(t, baseCfg, "coordinated_happy_event")
	require.NoError(t, re.engine.ExecuteTrigger(ctx, event))

	require.Equal(t, "completed", <-re.executionFinishedCh)
	require.Equal(t, int32(0), re.errorCalls.Load(), "OnExecutionError should not fire on the happy path")
	require.Equal(t, int32(0), re.engine.ActiveExecutions())
}

// TestCoordinatedEngine_CloseWaitsForInFlightExecution runs an execution on a
// caller's goroutine, as the trigger coordinator does. Close must cancel it and
// wait for it to return before closing the module it runs on.
func TestCoordinatedEngine_CloseWaitsForInFlightExecution(t *testing.T) {
	t.Parallel()

	baseCfg := coordinatedTestConfig(t)
	baseCfg.BillingClient = setupMockBillingClient(t)
	exec := newBlockingExecution()
	re := newTestEngine(t, baseCfg, v2.NewCoordinatedEngine, func(module *modulemocks.ModuleV2) {
		module.EXPECT().Start().Once()
		module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).RunAndReturn(exec.execute).Once()
		exec.expectModuleClose(module)
	})
	require.NoError(t, re.engine.Start(t.Context()))
	require.NoError(t, <-re.initializedCh)

	ctx, event := coordinatedEvent(t, baseCfg, "in_flight_event")
	execDone := make(chan struct{})
	go func() {
		defer close(execDone)
		_ = re.engine.ExecuteTrigger(ctx, event)
	}()
	<-exec.started

	closeWithin(t, re.engine, 10*time.Second)
	assert.True(t, exec.cancelled.Load(), "Close must cancel an execution running on the caller's goroutine")
	assert.False(t, exec.moduleClosedEarly.Load(), "the module was closed while an execution was still running on it")
	<-execDone
}

func TestCoordinatedEngine_ExecuteTriggerAfterCloseIsRejected(t *testing.T) {
	t.Parallel()

	baseCfg := coordinatedTestConfig(t)
	// No Execute expectation: a rejected call must never reach the module.
	re := newTestEngine(t, baseCfg, v2.NewCoordinatedEngine, func(module *modulemocks.ModuleV2) {
		module.EXPECT().Start().Once()
		module.EXPECT().Close().Once()
	})
	require.NoError(t, re.engine.Start(t.Context()))
	require.NoError(t, <-re.initializedCh)
	require.NoError(t, re.engine.Close())

	ctx, event := coordinatedEvent(t, baseCfg, "after_close_event")
	require.ErrorIs(t, re.engine.ExecuteTrigger(ctx, event), v2.ErrEngineClosed)
}

func coordinatedTestConfig(t *testing.T) *v2.EngineConfig {
	t.Helper()
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)

	cfg := defaultTestConfig(t, nil)
	cfg.CapRegistry = capreg
	cfg.TriggerAcknowledger = noopAcknowledger{}
	return cfg
}

// coordinatedEvent builds an event as the trigger coordinator delivers it, with
// the workflow's tenant on the ctx.
func coordinatedEvent(t *testing.T, cfg *v2.EngineConfig, eventID string) (context.Context, triggers.CoordinatedEvent) {
	t.Helper()
	ctx := contexts.WithCRE(t.Context(), contexts.CRE{Owner: cfg.WorkflowOwner, Workflow: cfg.WorkflowID})
	return ctx, triggers.CoordinatedEvent{
		WorkflowID:   cfg.WorkflowID,
		TriggerCapID: "id_0",
		TriggerIndex: 0,
		ObservedAt:   time.Now(),
		Event: capabilities.TriggerResponse{
			Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: eventID},
		},
	}
}

// blockingExecution parks a module Execute call until its ctx is cancelled or
// release is closed, and records whether the module was closed under it.
type blockingExecution struct {
	started, returned, release chan struct{}
	cancelled                  atomic.Bool
	moduleClosedEarly          atomic.Bool
}

func newBlockingExecution() *blockingExecution {
	return &blockingExecution{started: make(chan struct{}), returned: make(chan struct{}), release: make(chan struct{})}
}

func (b *blockingExecution) execute(ctx context.Context, _ *sdkpb.ExecuteRequest, _ host.ExecutionHelper) (*sdkpb.ExecutionResult, error) {
	close(b.started)
	defer close(b.returned)
	select {
	case <-ctx.Done():
		b.cancelled.Store(true)
		return nil, ctx.Err()
	case <-b.release:
		return &sdkpb.ExecutionResult{Result: &sdkpb.ExecutionResult_Value{}}, nil
	}
}

// expectModuleClose registers the module's Close and flags it if it runs
// while the execution is still in progress.
func (b *blockingExecution) expectModuleClose(module *modulemocks.ModuleV2) {
	module.EXPECT().Close().Run(func() {
		select {
		case <-b.returned:
		default:
			b.moduleClosedEarly.Store(true)
		}
	}).Once()
}

// closeWithin closes the engine, failing the test if Close does not return in time.
func closeWithin(t *testing.T, engine v2.WorkflowEngine, timeout time.Duration) {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- engine.Close() }()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(timeout):
		t.Fatalf("Close did not return within %s", timeout)
	}
}
