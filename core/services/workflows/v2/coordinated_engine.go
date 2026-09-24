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

var _ WorkflowEngine = (*coordinatedEngine)(nil)

// coordinatedEngine is an execution-only workflow engine: it registers no
// triggers itself and instead relies on an external manager to Subscribe,
// register/unregister triggers, call its ExecuteTrigger method and handle trigger acknowledgement.
// coordinatedEngine makes no trigger admission decisions. OnTriggerAdmission is not called.

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

func (e *coordinatedEngine) Start(ctx context.Context) error {
	return e.base.Start(ctx)
}

func (e *coordinatedEngine) Close() error {
	return e.base.Close()
}

func (e *coordinatedEngine) Ready() error {
	return e.base.Ready()
}

func (e *coordinatedEngine) HealthReport() map[string]error {
	return e.base.HealthReport()
}

func (e *coordinatedEngine) Name() string {
	return e.base.Name()
}

func (e *coordinatedEngine) ExecuteTrigger(ctx context.Context, event RoutedTriggerEvent) error {
	return e.base.ExecuteTrigger(ctx, event)
}

func (e *coordinatedEngine) Drain() bool {
	return e.base.Drain()
}

func (e *coordinatedEngine) ActiveExecutions() int32 {
	return e.base.ActiveExecutions()
}

func (e *coordinatedEngine) DrainStartedAt() (time.Time, bool) {
	return e.base.DrainStartedAt()
}

func (e *coordinatedEngine) Subscribe(ctx context.Context) ([]*sdkpb.TriggerSubscription, error) {
	return e.base.Subscribe(ctx)
}

func (e *coordinatedEngine) Tenant() contexts.CRE {
	return e.base.Tenant()
}

// IsCoordinated reports that this engine makes no trigger admission decisions.
// An external coordinator must own this
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
