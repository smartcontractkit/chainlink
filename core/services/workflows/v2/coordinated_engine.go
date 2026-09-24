package v2

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
)

var _ EventSink = (*coordinatedEngine)(nil)
var _ WorkflowEngine = (*coordinatedEngine)(nil)
var _ Subscriber = (*coordinatedEngine)(nil)

// coordinatedEngine is an execution-only workflow engine: it registers no
// triggers itself and instead relies on an external manager to Subscribe,
// register triggers and call its ExecuteTrigger method
//
// base is a named field rather than embedded so coordinatedEngine's own
// method set only contains what it explicitly delegates below
type coordinatedEngine struct {
	base *baseEngine
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

	e := &coordinatedEngine{base: base}
	base.attachService(lggr, "WorkflowCoordinatedEngine", e.start, e.close)
	return e, nil
}

// Start starts the engine.
func (e *coordinatedEngine) Start(ctx context.Context) error {
	return e.base.Start(ctx)
}

// Close stops the engine.
func (e *coordinatedEngine) Close() error {
	return e.base.Close()
}

// Ready reports readiness.
func (e *coordinatedEngine) Ready() error {
	return e.base.Ready()
}

// HealthReport reports health.
func (e *coordinatedEngine) HealthReport() map[string]error {
	return e.base.HealthReport()
}

// Name returns the service name.
func (e *coordinatedEngine) Name() string {
	return e.base.Name()
}

// ExecuteTrigger runs the workflow for a routed trigger event.
// This is how the TriggerCoordinator delivers events.
func (e *coordinatedEngine) ExecuteTrigger(ctx context.Context, event RoutedTriggerEvent) error {
	return e.base.ExecuteTrigger(ctx, event)
}

// Drain marks the engine as draining.
func (e *coordinatedEngine) Drain() bool {
	return e.base.Drain()
}

// ActiveExecutions returns the number of in-flight executions.
func (e *coordinatedEngine) ActiveExecutions() int32 {
	return e.base.ActiveExecutions()
}

// DrainStartedAt returns when draining began, if it has.
func (e *coordinatedEngine) DrainStartedAt() (time.Time, bool) {
	return e.base.DrainStartedAt()
}

// Subscribe issues the WASM Subscribe call.
// Used by the TriggerCoordinator to obtain this engine's subscriptions on demand.
func (e *coordinatedEngine) Subscribe(ctx context.Context) ([]*sdkpb.TriggerSubscription, error) {
	return e.base.Subscribe(ctx)
}

// Tenant is the engine's tenant identity.
func (e *coordinatedEngine) Tenant() contexts.CRE {
	return e.base.Tenant()
}

// IsCoordinated reports that an external TriggerCoordinator owns this
// engine's trigger registration and acknowledgement.
func (e *coordinatedEngine) IsCoordinated() bool {
	return true
}

func (e *coordinatedEngine) start(ctx context.Context) error {
	return e.base.startWith(ctx, e.init, nil)
}

func (e *coordinatedEngine) init(ctx context.Context) {
	// Tracer is no-op if DebugMode is false
	ctx, span := e.base.tracer.Start(ctx, "workflow_engine_init",
		trace.WithAttributes(
			attribute.String("version", "v2"),
			attribute.String("component", "workflow_engine"),
		))
	defer span.End()

	if err := e.base.initDONSubscribe(ctx); err != nil {
		e.base.cfg.Hooks.OnInitialized(err)
		return
	}

	e.base.initDone(ctx)
}

func (e *coordinatedEngine) close() error {
	ctx, cancel := e.base.shutdownCtx()
	defer cancel()

	e.base.closeCommon(ctx)
	return nil
}
