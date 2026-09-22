package v2

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
)

var _ EventSink = (*ExecutionEngine)(nil)
var _ WorkflowEngine = (*ExecutionEngine)(nil)

// ExecutionEngine is the execution-only workflow engine.
//
// All execution machinery lives on the embedded baseEngine, including the single
// services.Engine. ExecutionEngine adds only its own lifecycle.
type ExecutionEngine struct {
	*baseEngine

	// subs is what the WASM Subscribe call declared, retained for the
	// TriggerCoordinator to register on this engine's behalf. Written once by
	// init, strictly before OnInitialized fires; see Subscriptions.
	subs []*sdkpb.TriggerSubscription
}

// NewExecutionEngine constructs the execution-only engine. cfg.TriggerAcknowledger
// is required: the engine holds no handles, so it cannot acknowledge by itself.
func NewExecutionEngine(cfg *EngineConfig) (*ExecutionEngine, error) {
	if cfg.TriggerAcknowledger == nil {
		return nil, errors.New("trigger acknowledger not set")
	}

	base, lggr, err := newBaseEngine(cfg)
	if err != nil {
		return nil, err
	}

	e := &ExecutionEngine{baseEngine: base}
	base.attachService(lggr, "WorkflowExecutionEngine", e.start, e.close)
	return e, nil
}

func (e *ExecutionEngine) start(ctx context.Context) error {
	return e.startWith(ctx, e.init, nil)
}

// init is the execution-only initialization: DON sync -> Subscribe ->
// retain subscriptions -> OnInitialized.
func (e *ExecutionEngine) init(ctx context.Context) {
	// Tracer is no-op if DebugMode is false
	ctx, span := e.tracer.Start(ctx, "workflow_engine_init",
		trace.WithAttributes(
			attribute.String("version", "v2"),
			attribute.String("component", "workflow_engine"),
		))
	defer span.End()

	if err := e.initDONSubscribe(ctx); err != nil {
		e.cfg.Hooks.OnInitialized(err)
		return
	}
	subs, err := e.initSubscriptions(ctx)
	if err != nil {
		e.cfg.Hooks.OnInitialized(err)
		return
	}
	e.subs = subs
	e.initDone(ctx)
}

// Subscriptions returns the trigger subscriptions this engine declared during
// init, with the tenant identity they belong to. The engine registers nothing
// itself; the TriggerCoordinator registers these on its behalf.
//
// Only valid after OnInitialized has fired with a nil error. No lock is needed:
// init writes these before firing that hook, and the caller reads them only
// after observing it, which is the happens-before edge.
func (e *ExecutionEngine) Subscriptions() ([]*sdkpb.TriggerSubscription, contexts.CRE) {
	return e.subs, e.cre()
}

func (e *ExecutionEngine) close() error {
	ctx, cancel := e.shutdownCtx()
	defer cancel()

	e.closeCommon(ctx)
	return nil
}
