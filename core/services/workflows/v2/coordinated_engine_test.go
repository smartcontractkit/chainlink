package v2_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
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
// a CoordinatedEngine registers no triggers of its own, so this is the only way an event reaches it.
func TestCoordinatedEngine_ExecuteTrigger(t *testing.T) {
	t.Parallel()

	module := modulemocks.NewModuleV2(t)
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)

	executionFinishedCh := make(chan string, 1)
	executionErrorCh := make(chan string, 1)

	cfg := defaultTestConfig(t, nil)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = setupMockBillingClient(t)
	cfg.TriggerAcknowledger = noopAcknowledger{}
	cfg.Hooks = v2.LifecycleHooks{
		OnExecutionFinished: func(_ string, status string) {
			executionFinishedCh <- status
		},
		OnExecutionError: func(msg string) {
			executionErrorCh <- msg
		},
	}

	engine, err := v2.NewCoordinatedEngine(cfg)
	require.NoError(t, err)

	// The execution must reach WASM and return a value.
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		Return(&sdkpb.ExecutionResult{
			Result: &sdkpb.ExecutionResult_Value{},
		}, nil).
		Once()

	ctx := contexts.WithCRE(t.Context(), contexts.CRE{
		Owner:    cfg.WorkflowOwner,
		Workflow: cfg.WorkflowID,
	})
	event := v2.RoutedTriggerEvent{
		WorkflowID:   cfg.WorkflowID,
		TriggerCapID: "id_0",
		TriggerIndex: 0,
		ObservedAt:   time.Now(),
		Event: capabilities.TriggerResponse{
			Event: capabilities.TriggerEvent{
				TriggerType: "basic-trigger@1.0.0",
				ID:          "coordinated_happy_event",
			},
		},
	}

	require.NoError(t, engine.ExecuteTrigger(ctx, event))

	require.Equal(t, "completed", <-executionFinishedCh)
	select {
	case msg := <-executionErrorCh:
		t.Fatalf("unexpected OnExecutionError: %s", msg)
	default:
	}
	require.Equal(t, int32(0), engine.ActiveExecutions())
}
