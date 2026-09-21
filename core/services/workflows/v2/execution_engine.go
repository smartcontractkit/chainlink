package v2

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var _ EventSink = (*ExecutionEngine)(nil)
var _ WorkflowEngine = (*ExecutionEngine)(nil)

// ExecutionEngine is the execution-only workflow engine.
//
// All execution machinery lives on the embedded baseEngine, including the single
// services.Engine. ExecutionEngine adds only its own lifecycle.
type ExecutionEngine struct {
	*baseEngine
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
	base.attachService(lggr, e.start, e.close)
	return e, nil
}

func (e *ExecutionEngine) start(ctx context.Context) error {
	return e.startWith(ctx, e.init, nil)
}

// init is the execution-only initialization: DON sync -> Subscribe ->
// OnSubscriptionsReady -> OnInitialized.
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
	if _, err := e.initSubscriptions(ctx); err != nil {
		e.cfg.Hooks.OnInitialized(err)
		return
	}
	e.initDone(ctx)
}

func (e *ExecutionEngine) close() error {
	ctx, cancel := e.shutdownCtx()
	defer cancel()

	e.closeCommon(ctx)
	return nil
}
