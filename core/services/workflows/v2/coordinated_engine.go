package v2

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var _ EventSink = (*coordinatedEngine)(nil)
var _ WorkflowEngine = (*coordinatedEngine)(nil)
var _ Subscriber = (*coordinatedEngine)(nil)

// coordinatedEngine is an execution-only workflow engine: it registers no
// triggers itself and instead relies on an external manager to Subscribe,
// register triggers and call its ExecuteTrigger method
//
// All execution machinery lives on the embedded baseEngine, including the single
// services.Engine.
type coordinatedEngine struct {
	*baseEngine
}

// NewCoordinatedEngine constructs the execution-only engine. cfg.TriggerAcknowledger
// is required: the engine holds no handles, so it cannot acknowledge by itself.
func NewCoordinatedEngine(cfg *EngineConfig) (WorkflowEngine, error) {
	if cfg.TriggerAcknowledger == nil {
		return nil, errors.New("trigger acknowledger not set")
	}

	base, lggr, err := newBaseEngine(cfg)
	if err != nil {
		return nil, err
	}

	e := &coordinatedEngine{baseEngine: base}
	base.attachService(lggr, "WorkflowCoordinatedEngine", e.start, e.close)
	return e, nil
}

func (e *coordinatedEngine) start(ctx context.Context) error {
	return e.startWith(ctx, e.init, nil)
}

func (e *coordinatedEngine) init(ctx context.Context) {
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

	e.initDone(ctx)
}

func (e *coordinatedEngine) close() error {
	ctx, cancel := e.shutdownCtx()
	defer cancel()

	e.closeCommon(ctx)
	return nil
}

// IsCoordinated indicates whether the engine needs an external trigger coordinator.
func (e *coordinatedEngine) IsCoordinated() bool {
	return true
}
