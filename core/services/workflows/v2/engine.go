package v2

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
)

var (
	_ Acknowledger   = (*engine)(nil)
	_ EventSink      = (*engine)(nil)
	_ WorkflowEngine = (*engine)(nil)
)

// engine is the legacy trigger-owning workflow engine: it wraps the shared
// execution machinery and adds the extra responsibilities of trigger registration,
// handle ownership, acknowledgement, and the node's workflow-count limit.
//
// base is a named field rather than embedded so engine's own method set only
// contains what it explicitly delegates below.
type engine struct {
	base *baseEngine

	workflowLimitUsed atomic.Bool // true if GlobalWorkflowLimit must be freed

	// registration ID -> trigger handle
	triggers map[string]*triggers.Handle
	// used to separate registration and unregistration phases
	triggersRegMu sync.Mutex

	allTriggerEventsQueueCh limits.QueueLimiter[triggers.CoordinatedEvent]
	executionsSemaphore     limits.ResourcePoolLimiter[int]
}

// NewEngine constructs the legacy trigger-owning engine: it registers its own
// triggers, holds the handles, acknowledges through itself, and owns the
// workflow-count limit.
func NewEngine(cfg *EngineConfig) (WorkflowEngine, error) {
	base, lggr, err := newBaseEngine(cfg)
	if err != nil {
		return nil, err
	}

	e := &engine{
		base:                    base,
		triggers:                make(map[string]*triggers.Handle),
		allTriggerEventsQueueCh: cfg.LocalLimiters.TriggerEventQueue,
		executionsSemaphore:     cfg.LocalLimiters.ExecutionConcurrency,
	}

	// Self-inject: the engine is its own acknowledger because it holds the
	// trigger handles.
	if cfg.TriggerAcknowledger == nil {
		cfg.TriggerAcknowledger = e
	}

	base.initServiceEngine(lggr, "WorkflowEngineV2", e.start, e.close)
	return e, nil
}

func (e *engine) Start(ctx context.Context) error {
	return e.base.Start(ctx)
}

func (e *engine) Close() error {
	return e.base.Close()
}

func (e *engine) Ready() error {
	return e.base.Ready()
}

func (e *engine) HealthReport() map[string]error {
	return e.base.HealthReport()
}

func (e *engine) Name() string {
	return e.base.Name()
}

func (e *engine) ExecuteTrigger(ctx context.Context, event triggers.CoordinatedEvent) error {
	return e.base.ExecuteTrigger(ctx, event)
}

func (e *engine) Drain() bool {
	return e.base.Drain()
}

func (e *engine) ActiveExecutions() int32 {
	return e.base.ActiveExecutions()
}

func (e *engine) DrainStartedAt() (time.Time, bool) {
	return e.base.DrainStartedAt()
}

func (e *engine) Subscribe(ctx context.Context) ([]*sdkpb.TriggerSubscription, error) {
	return e.base.Subscribe(ctx)
}

func (e *engine) Tenant() contexts.CRE {
	return e.base.Tenant()
}

// IsCoordinated reports that this engine manages its own trigger registration and acknowledgement.
func (e *engine) IsCoordinated() bool {
	return false
}

func (e *engine) start(ctx context.Context) error {
	return e.base.startWith(ctx, e.init, e.handleAllTriggerEvents)
}

// init is the legacy initialization: it acquires the workflow-count limit before
// anything else and registers the workflow's triggers before reporting success.
func (e *engine) init(ctx context.Context) {
	// Tracer is no-op if DebugMode is false
	ctx, span := e.base.tracer.Start(ctx, "workflow_engine_init",
		trace.WithAttributes(
			attribute.String("version", "v2"),
			attribute.String("component", "workflow_engine"),
		))
	defer span.End()

	if err := e.useWorkflowLimit(ctx); err != nil {
		e.base.cfg.Hooks.OnInitialized(err)
		return
	}
	if err := e.base.initDONSubscribe(ctx); err != nil {
		e.base.cfg.Hooks.OnInitialized(err)
		return
	}

	subscriptions, err := e.base.Subscribe(ctx)
	if err != nil {
		e.base.logger().Errorw("failed to subscribe to triggers", "err", err)
		e.base.cfg.Hooks.OnInitialized(err)
		return
	}
	if err := e.runTriggerSubscriptionPhase(ctx, subscriptions); err != nil {
		e.base.logger().Errorw("Workflow Engine initialization failed", "err", err)
		e.base.cfg.Hooks.OnInitialized(err)
		return
	}
	e.base.initDone(ctx)
}

// useWorkflowLimit acquires one slot of the node's workflow-count limit. The
// returned error is the one that must reach OnInitialized: the scope-specific
// sentinel for a limit breach, the raw error otherwise.
func (e *engine) useWorkflowLimit(ctx context.Context) error {
	if err := e.base.cfg.GlobalWorkflowLimit.Use(ctx, 1); err != nil {
		errLimited, ok := errors.AsType[limits.ErrorResourceLimited[int]](err)
		if !ok {
			return err
		}
		switch errLimited.Scope {
		case settings.ScopeOwner:
			e.base.logger().Infow("Per owner workflow count limit reached", "err", err)
			e.base.metrics.IncrementWorkflowLimitPerOwnerCounter(ctx)
			return types.ErrPerOwnerWorkflowCountLimitReached
		case settings.ScopeGlobal:
			e.base.logger().Infow("Global workflow count limit reached", "err", err)
			e.base.metrics.IncrementWorkflowLimitGlobalCounter(ctx)
			return types.ErrGlobalWorkflowCountLimitReached
		default:
			e.base.logger().Errorw("Workflow count limit reached for unexpected scope", "scope", errLimited.Scope, "err", err)
			return err
		}
	}

	e.workflowLimitUsed.Store(true)
	return nil
}

func (e *engine) close() error {
	ctx, cancel := e.base.shutdownCtx()
	defer cancel()

	e.triggersRegMu.Lock()
	e.unregisterAllTriggers(ctx)
	e.triggersRegMu.Unlock()
	e.base.metrics.IncrementWorkflowUnregisteredCounter(ctx)

	e.base.closeCommon(ctx)

	if e.workflowLimitUsed.Load() {
		return e.base.cfg.GlobalWorkflowLimit.Free(ctx, 1)
	}
	return nil
}

func (e *engine) runTriggerSubscriptionPhase(ctx context.Context, subscriptions []*sdkpb.TriggerSubscription) error {
	var creGetter settings.Getter
	if e.base.cfg.LocalLimiters != nil {
		creGetter = e.base.cfg.LocalLimiters.Settings
	}

	deps := triggers.RegisterDeps{
		CapRegistry:  e.base.cfg.CapRegistry,
		RegTimeout:   e.base.cfg.LocalLimiters.TriggerRegistrationsTime,
		ChainAllowed: e.base.cfg.LocalLimiters.ChainAllowed,
		Settings:     creGetter,
		Logger:       e.base.logger(),
		Metrics:      e.base.metrics,
	}
	meta := triggers.RegisterMetadata{
		WorkflowID:                    e.base.cfg.WorkflowID,
		WorkflowOwner:                 e.base.cfg.WorkflowOwner,
		WorkflowName:                  e.base.cfg.WorkflowName,
		WorkflowTag:                   e.base.cfg.WorkflowTag,
		WorkflowDonID:                 e.base.localNode.Load().WorkflowDON.ID,
		WorkflowDonConfigVersion:      pinnedWorkflowDonConfigVersion,
		WorkflowRegistryChainSelector: e.base.cfg.WorkflowRegistryChainSelector,
		WorkflowRegistryAddress:       e.base.cfg.WorkflowRegistryAddress,
		OrgID:                         e.base.orgID,
	}

	triggerCapIDs, handles, eventChans, err := triggers.Register(ctx, deps, meta, subscriptions)
	if err != nil {
		return err
	}

	e.triggersRegMu.Lock()
	e.triggers = handles
	e.triggersRegMu.Unlock()

	// start listening for trigger events now that all registrations succeeded
	for idx, triggerEventCh := range eventChans {
		triggerID := subscriptions[idx].Id
		deliver := func(ctx context.Context, routed triggers.CoordinatedEvent) {
			eventID := routed.Event.Event.ID
			if err := e.put(ctx, routed); err != nil {
				// Draining is expected during workflow deletion, so it logs at info rather than error level.
				if errors.Is(err, ErrEngineDraining) {
					e.base.logger().Infow("Dropping trigger event: engine draining", "triggerID", triggerID, "eventID", eventID)
				} else {
					e.base.logger().Errorw("Failed to put routed trigger event", "triggerID", triggerID, "eventID", eventID, "err", err)
				}
			}
		}
		e.base.srvcEng.GoCtx(context.WithoutCancel(ctx), func(ctx context.Context) {
			triggers.ReadLoop(ctx, e.base.logger(), e.base.metrics, e.base.cfg.Clock, e.base.cfg.WorkflowID, triggerID, idx, triggerEventCh, deliver)
		})
	}
	e.base.cfg.Hooks.OnSubscribedToTriggers(triggerCapIDs)
	return nil
}

// NOTE: needs to be called under the triggersRegMu lock
func (e *engine) unregisterAllTriggers(ctx context.Context) {
	failCount := triggers.Unregister(ctx, e.base.logger(), e.base.cfg.WorkflowID, e.base.localNode.Load().WorkflowDON.ID, e.triggers)
	e.base.logger().Infow("All triggers unregistered", "numTriggers", len(e.triggers), "failed", failCount)
	e.triggers = make(map[string]*triggers.Handle)
}

func (e *engine) Ack(ctx context.Context, triggerCapID, triggerRegistrationID, eventID string) error {
	e.triggersRegMu.Lock()
	handle := e.triggers[triggerRegistrationID]
	e.triggersRegMu.Unlock()

	return triggers.Ack(ctx, e.base.logger(), e.base.metrics, triggerCapID, triggerRegistrationID, eventID, handle)
}

// errObservedAtMissing guards put's deadline derivation: put computes the queue
// deadline from ObservedAt, so an unstamped event is rejected rather than given
// a garbage deadline.
var errObservedAtMissing = errors.New("trigger event ObservedAt not set")

// put enqueues a trigger event into the engine's internal queue.
func (e *engine) put(ctx context.Context, event triggers.CoordinatedEvent) error {
	triggerID := event.TriggerCapID
	eventID := event.Event.Event.ID
	idx := event.TriggerIndex

	if e.base.Draining() {
		e.base.logger().Infow("Engine is draining, dropping trigger event before enqueue", "triggerID", triggerID, "eventID", eventID)
		tm := e.base.metrics.With(platform.KeyTriggerID, triggerID)
		tm.IncrementTriggerEventEnqueueDroppedCounter(ctx)
		tm.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonEnqueueDraining)
		e.base.cfg.Hooks.OnTriggerEventDropped(triggerID, eventID, "draining")
		return ErrEngineDraining
	}

	// Admission check: an external management layer (e.g. ShardFailoverManager)
	// decides whether this event should be processed by the engine.
	if err := e.base.cfg.Hooks.OnTriggerAdmission(ctx, event); err != nil {
		if errors.Is(err, ErrAdmissionCache) {
			e.base.logger().Infow("Trigger event cached for failover, not enqueuing", "triggerID", triggerID, "eventID", eventID)
			return err
		}
		// Denied: ACK and drop
		registrationID := triggers.RegistrationID(e.base.cfg.WorkflowID, event.TriggerIndex)
		if ackErr := e.base.cfg.TriggerAcknowledger.Ack(ctx, event.TriggerCapID, registrationID, eventID); ackErr != nil {
			e.base.logger().Errorw("failed to ACK trigger after admission denial", "eventID", eventID, "err", ackErr)
		}
		e.base.metrics.With(platform.KeyTriggerID, triggerID).IncrementTriggerEventDroppedTotal(ctx, "admission_denied")
		return err
	}

	// The deadline below is derived from ObservedAt, so reject events that
	// arrive without it rather than silently backfill a garbage deadline.
	if event.ObservedAt.IsZero() {
		e.base.logger().Errorw("Trigger event missing ObservedAt, dropping", "triggerID", triggerID, "eventID", eventID)
		e.base.metrics.With(platform.KeyTriggerID, triggerID).IncrementTriggerEventDroppedTotal(ctx, "observed_at_missing")
		return errObservedAtMissing
	}
	queueTimeout, err := e.base.cfg.LocalLimiters.TriggerEventQueueTimeout.Limit(ctx)
	if err != nil {
		if !limits.IsErrRecoverable(err) {
			// No value was resolved, so the deadline below would already be expired.
			e.base.logger().Errorw("Failed to get trigger event queue time limit with no usable value", "err", err)
			e.base.metrics.With(platform.KeyTriggerID, triggerID).
				IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonQueueAgeLimitReadFailed)
			return ErrEnqueueFailed
		}
		// A settings read failure is not a reason to drop a customer's trigger event:
		// the limiter still returns a usable timeout, so stamp the deadline and continue.
		e.base.logger().Errorw("Failed to get trigger event queue time limit; continuing with the value the limiter returned", "err", err)
		e.base.metrics.IncrementLimitReadFallbackCounter(ctx, cresettings.Default.PerWorkflow.TriggerEventQueueTimeout.Key)
	}
	event.Deadline = event.ObservedAt.Add(queueTimeout)

	if err := e.allTriggerEventsQueueCh.Put(ctx, event); err != nil {
		tm := e.base.metrics.With(platform.KeyTriggerID, triggerID)
		tm.IncrementTriggerEventEnqueueDroppedCounter(ctx)
		retErr := ErrEnqueueFailed
		if _, ok := errors.AsType[limits.ErrorQueueFull](err); ok {
			// queue full, drop the event
			e.base.logger().Errorw("Trigger event queue is full, dropping event", "triggerID", triggerID, "triggerIndex", idx, "err", err)
			tm.IncrementWorkflowTriggerEventQueueFullCounter(ctx)
			tm.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonEnqueueQueueFull)
			retErr = ErrQueueFull
		} else {
			tm.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonEnqueueFailed)
		}
		e.base.logger().Errorw("Failed to enqueue trigger event", "triggerID", triggerID, "triggerIndex", idx, "err", err)
		tm.IncrementWorkflowTriggerEventErrorCounter(ctx)
		return retErr
	}
	e.base.metrics.With(platform.KeyTriggerID, triggerID).IncrementTriggerEventEnqueuedCounter(ctx)
	e.base.logger().Debugw("Enqueued trigger event", "triggerID", triggerID, "eventID", eventID)

	return nil
}

// handleAllTriggerEvents drains the engine's trigger-event queue (populated by put method)
// and executes each event in turn.
func (e *engine) handleAllTriggerEvents(ctx context.Context) {
	for {
		queueHead, err := e.allTriggerEventsQueueCh.Wait(ctx)
		if err != nil {
			return
		}
		eventID := queueHead.Event.Event.ID
		triggerMetricLabels := e.base.metrics.With(platform.KeyTriggerID, queueHead.TriggerCapID)
		if e.base.Draining() {
			triggerMetricLabels.IncrementTriggerEventDequeueDroppedCounter(ctx)
			triggerMetricLabels.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonDequeueDraining)
			e.base.cfg.Hooks.OnTriggerEventDropped(queueHead.TriggerCapID, eventID, "draining")
			e.base.logger().Infow("Engine is draining, stopping trigger handling loop", "eventID", eventID, "triggerID", queueHead.TriggerCapID)
			return
		}

		now := e.base.cfg.Clock.Now()
		eventAge := now.Sub(queueHead.ObservedAt)
		e.base.logger().Debugw("Popped a trigger event from the queue", "eventID", eventID, "eventAgeMs", eventAge.Milliseconds())
		triggerMetricLabels.RecordTriggerEventQueueWaitSeconds(ctx, eventAge.Seconds())
		if now.After(queueHead.Deadline) {
			e.base.logger().Warnw("Trigger event is too old, skipping execution", "triggerID", queueHead.TriggerCapID, "eventID", eventID, "eventAgeMs", eventAge.Milliseconds())
			triggerMetricLabels.IncrementTriggerEventExpiredCounter(ctx)
			triggerMetricLabels.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonExpired)
			continue
		}

		semWaitStart := e.base.cfg.Clock.Now()
		free, err := e.executionsSemaphore.Wait(ctx, 1) // block if too many concurrent workflow executions
		triggerMetricLabels.RecordExecutionSemaphoreWaitSeconds(ctx, e.base.cfg.Clock.Now().Sub(semWaitStart).Seconds())
		if err != nil {
			e.base.logger().Errorw("Failed to acquire executions semaphore", "err", err)
			triggerMetricLabels.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonExecutionSemaphoreWaitFailed)
			continue
		}

		e.base.srvcEng.GoCtx(context.WithoutCancel(ctx), func(ctx context.Context) {
			defer free()

			if err := e.ExecuteTrigger(ctx, queueHead); err != nil {
				// Dedup is an expected outcome (the event is handled, just not executed here), so it's logged at info rather than error level.
				if errors.Is(err, ErrDuplicateExecution) {
					e.base.logger().Infow("Skipping trigger event execution", "triggerID", queueHead.TriggerCapID, "eventID", queueHead.Event.Event.ID, "err", err)
				} else {
					e.base.logger().Errorw("Failed to execute trigger event", "triggerID", queueHead.TriggerCapID, "eventID", queueHead.Event.Event.ID, "err", err)
				}
			}
		})
	}
}
