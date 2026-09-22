package v2

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var _ EventSink = (*CoordinatedEngine)(nil)
var _ WorkflowEngine = (*CoordinatedEngine)(nil)
var _ Subscriber = (*CoordinatedEngine)(nil)

// CoordinatedEngine is the execution-only workflow engine: it registers no
// triggers itself and instead relies on a TriggerCoordinator to Subscribe,
// register, and deliver events to it.
//
// All execution machinery lives on the embedded baseEngine, including the single
// services.Engine. CoordinatedEngine adds only its own lifecycle.
type CoordinatedEngine struct {
	*baseEngine
}

// NewCoordinatedEngine constructs the execution-only engine. cfg.TriggerAcknowledger
// is required: the engine holds no handles, so it cannot acknowledge by itself.
func NewCoordinatedEngine(cfg *EngineConfig) (*CoordinatedEngine, error) {
	if cfg.TriggerAcknowledger == nil {
		return nil, errors.New("trigger acknowledger not set")
	}

	base, lggr, err := newBaseEngine(cfg)
	if err != nil {
		return nil, err
	}

	e := &CoordinatedEngine{baseEngine: base}
	base.attachService(lggr, "WorkflowCoordinatedEngine", e.start, e.close)
	return e, nil
}

func (e *CoordinatedEngine) start(ctx context.Context) error {
	return e.startWith(ctx, e.init, nil)
}

// init is the execution-only initialization: DON sync -> OnInitialized.
func (e *CoordinatedEngine) init(ctx context.Context) {
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

func (e *CoordinatedEngine) close() error {
	ctx, cancel := e.shutdownCtx()
	defer cancel()

	e.closeCommon(ctx)
	return nil
}
