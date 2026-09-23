package v2

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

var _ Acknowledger = (*Engine)(nil)
var _ EventSink = (*Engine)(nil)
var _ WorkflowEngine = (*Engine)(nil)

// Engine is the legacy trigger-owning workflow engine: it embeds the shared
// execution machinery and adds the extra responsibilities of trigger registration,
// handle ownership, acknowledgement, and the node's workflow-count limit.
//
// Fields declared here must not duplicate baseEngine's: a shadowed field would
// leave these methods reading a zero value while the execution path reads the
// real one, and the compiler will not catch it.
type Engine struct {
	*baseEngine

	workflowLimitUsed atomic.Bool // true if GlobalWorkflowLimit must be freed

	// registration ID -> trigger capability
	triggers map[string]*triggerCapability
	// used to separate registration and unregistration phases
	triggersRegMu sync.Mutex

	allTriggerEventsQueueCh limits.QueueLimiter[RoutedTriggerEvent]
	executionsSemaphore     limits.ResourcePoolLimiter[int]
}

type triggerCapability struct {
	capabilities.TriggerCapability
	payload *anypb.Any
	method  string
}

// NewEngine constructs the legacy trigger-owning engine: it registers its own
// triggers, holds the handles, acknowledges through itself, and owns the
// workflow-count limit.
func NewEngine(cfg *EngineConfig) (*Engine, error) {
	base, lggr, err := newBaseEngine(cfg)
	if err != nil {
		return nil, err
	}

	e := &Engine{
		baseEngine:              base,
		triggers:                make(map[string]*triggerCapability),
		allTriggerEventsQueueCh: cfg.LocalLimiters.TriggerEventQueue,
		executionsSemaphore:     cfg.LocalLimiters.ExecutionConcurrency,
	}

	// Self-inject: the engine is its own acknowledger in M1, because it holds the
	// trigger handles. In M2 the OCR reporting plugin implements the Acknowledger.
	if cfg.TriggerAcknowledger == nil {
		cfg.TriggerAcknowledger = e
	}

	// The lifecycle is Engine's; the single services.Engine lives on the base.
	base.attachService(lggr, "WorkflowEngineV2", e.start, e.close)
	return e, nil
}

func (e *Engine) start(ctx context.Context) error {
	return e.startWith(ctx, e.init, e.handleAllTriggerEvents)
}

// init is the legacy initialization: it acquires the workflow-count limit before
// anything else and registers the workflow's triggers before reporting success.
func (e *Engine) init(ctx context.Context) {
	// Tracer is no-op if DebugMode is false
	ctx, span := e.tracer.Start(ctx, "workflow_engine_init",
		trace.WithAttributes(
			attribute.String("version", "v2"),
			attribute.String("component", "workflow_engine"),
		))
	defer span.End()

	if err := e.useWorkflowLimit(ctx); err != nil {
		e.cfg.Hooks.OnInitialized(err)
		return
	}
	if err := e.initDONSubscribe(ctx); err != nil {
		e.cfg.Hooks.OnInitialized(err)
		return
	}

	subscriptions, err := e.Subscribe(ctx)
	if err != nil {
		e.logger().Errorw("failed to subscribe to triggers", "err", err)
		e.cfg.Hooks.OnInitialized(err)
		return
	}
	if err := e.runTriggerSubscriptionPhase(ctx, subscriptions); err != nil {
		e.logger().Errorw("Workflow Engine initialization failed", "err", err)
		e.cfg.Hooks.OnInitialized(err)
		return
	}
	e.initDone(ctx)
}

// useWorkflowLimit acquires one slot of the node's workflow-count limit. The
// returned error is the one that must reach OnInitialized: the scope-specific
// sentinel for a limit breach, the raw error otherwise.
func (e *Engine) useWorkflowLimit(ctx context.Context) error {
	err := e.cfg.GlobalWorkflowLimit.Use(ctx, 1)
	if err == nil {
		e.workflowLimitUsed.Store(true)
		return nil
	}

	errLimited, ok := errors.AsType[limits.ErrorResourceLimited[int]](err)
	if !ok {
		return err
	}
	switch errLimited.Scope {
	case settings.ScopeOwner:
		e.logger().Infow("Per owner workflow count limit reached", "err", err)
		e.metrics.IncrementWorkflowLimitPerOwnerCounter(ctx)
		return types.ErrPerOwnerWorkflowCountLimitReached
	case settings.ScopeGlobal:
		e.logger().Infow("Global workflow count limit reached", "err", err)
		e.metrics.IncrementWorkflowLimitGlobalCounter(ctx)
		return types.ErrGlobalWorkflowCountLimitReached
	default:
		e.logger().Errorw("Workflow count limit reached for unexpected scope", "scope", errLimited.Scope, "err", err)
		return err
	}
}

func (e *Engine) close() error {
	ctx, cancel := e.shutdownCtx()
	defer cancel()

	e.triggersRegMu.Lock()
	e.unregisterAllTriggers(ctx)
	e.triggersRegMu.Unlock()
	e.metrics.IncrementWorkflowUnregisteredCounter(ctx)

	e.closeCommon(ctx)

	if e.workflowLimitUsed.Load() {
		return e.cfg.GlobalWorkflowLimit.Free(ctx, 1)
	}
	return nil
}

func (e *Engine) runTriggerSubscriptionPhase(ctx context.Context, subscriptions []*sdkpb.TriggerSubscription) error {
	// check if all requested triggers exist in the registry
	triggers := make([]capabilities.TriggerCapability, 0, len(subscriptions))
	for _, sub := range subscriptions {
		_, labels, _ := capabilities.ParseID(sub.Id)
		chainSelector, err2 := capabilities.ChainSelectorLabel(labels)
		if err2 != nil {
			return fmt.Errorf("invalid chain selector for ID %s: %w", sub.Id, err2)
		}
		if chainSelector != nil {
			err2 := e.cfg.LocalLimiters.ChainAllowed.AllowErr(contexts.WithChainSelector(ctx, *chainSelector))
			if err2 != nil {
				if errors.Is(err2, limits.ErrorNotAllowed{}) {
					return fmt.Errorf("unable to subscribe to capability %s: ChainSelector %d: %w", sub.Id, *chainSelector, err2)
				}
				return fmt.Errorf("failed to check access for ChainSelector %d: %w", *chainSelector, err2)
			}
		}
		triggerCap, triggerErr := e.cfg.CapRegistry.GetTrigger(ctx, sub.Id)
		if triggerErr != nil {
			return fmt.Errorf("trigger capability not found: %w", triggerErr)
		}
		triggers = append(triggers, triggerCap)
	}

	// register to all triggers concurrently
	regCtx, regCancel, err := e.cfg.LocalLimiters.TriggerRegistrationsTime.WithTimeout(ctx)
	if err != nil {
		return err
	}
	defer regCancel()

	// trigger registration results for use in concurrent trigger subscriptions
	type triggerRegResult struct {
		index          int
		registrationID string
		triggerCap     capabilities.TriggerCapability
		eventCh        <-chan capabilities.TriggerResponse
		payload        *anypb.Any
		method         string
		triggerCapID   string
	}

	resultsCh := make(chan triggerRegResult, len(subscriptions))
	g, gCtx := errgroup.WithContext(regCtx)

	// Launch concurrent trigger registrations
	for i, sub := range subscriptions {
		triggerCap := triggers[i]
		g.Go(func() error {
			registrationID := TriggerRegistrationID(e.cfg.WorkflowID, i)
			args := []any{"triggerID", sub.Id, "method", sub.Method}
			if sub.Payload != nil {
				args = append(args, "payload", protojson.Format(sub.Payload))
			}
			e.logger().Infow("Registering trigger", args...)
			metadata := capabilities.RequestMetadata{
				WorkflowID:                    e.cfg.WorkflowID,
				WorkflowOwner:                 e.cfg.WorkflowOwner,
				WorkflowName:                  e.cfg.WorkflowName.Hex(),
				WorkflowTag:                   e.cfg.WorkflowTag,
				DecodedWorkflowName:           e.cfg.WorkflowName.String(),
				WorkflowDonID:                 e.localNode.Load().WorkflowDON.ID,
				WorkflowDonConfigVersion:      pinnedWorkflowDonConfigVersion,
				ReferenceID:                   fmt.Sprintf("trigger_%d", i),
				WorkflowRegistryChainSelector: e.cfg.WorkflowRegistryChainSelector,
				WorkflowRegistryAddress:       e.cfg.WorkflowRegistryAddress,
				EngineVersion:                 platform.ValueWorkflowVersionV2,
				// no WorkflowExecutionID needed (or available at this stage)
			}
			var creGetter settings.Getter
			if e.cfg.LocalLimiters != nil {
				creGetter = e.cfg.LocalLimiters.Settings
			}
			propagateOrgIDMeta, _ := cresettings.Default.PropagateOrgIDInRequestMetadata.GetOrDefault(gCtx, creGetter)
			if propagateOrgIDMeta && e.orgID != "" {
				metadata.OrgID = e.orgID
			}
			triggerEventCh, regErr := triggerCap.RegisterTrigger(gCtx, capabilities.TriggerRegistrationRequest{
				TriggerID: registrationID,
				Metadata:  metadata,
				Payload:   sub.Payload,
				Method:    sub.Method,
				// no Config needed - NoDAG uses Payload
			})
			if regErr != nil {
				e.logger().Errorw("Trigger registration failed", "triggerID", sub.Id, "err", regErr)
				e.metrics.With(platform.KeyTriggerID, sub.Id).IncrementRegisterTriggerFailureCounter(gCtx)
				return fmt.Errorf("failed to register trigger %s: %w", sub.Id, regErr)
			}
			// Send successful result
			resultsCh <- triggerRegResult{
				index:          i,
				registrationID: registrationID,
				triggerCap:     triggerCap,
				eventCh:        triggerEventCh,
				payload:        sub.Payload,
				method:         sub.Method,
				triggerCapID:   sub.Id,
			}
			return nil
		})
	}

	// wait for all registrations to complete.
	// returns first non-nil error.
	registrationErr := g.Wait()
	close(resultsCh)

	// Collect results into e.triggers map
	e.triggersRegMu.Lock()
	defer e.triggersRegMu.Unlock()

	eventChans := make([]<-chan capabilities.TriggerResponse, len(subscriptions))
	triggerCapIDs := make([]string, len(subscriptions))

	for result := range resultsCh {
		e.triggers[result.registrationID] = &triggerCapability{
			TriggerCapability: result.triggerCap,
			payload:           result.payload,
			method:            result.method,
		}
		eventChans[result.index] = result.eventCh
		triggerCapIDs[result.index] = result.triggerCapID
	}

	// If any registration failed, unregister successful ones and return error
	if registrationErr != nil {
		e.logger().Errorw("One or more trigger registrations failed - reverting all", "err", registrationErr)
		e.unregisterAllTriggers(ctx) // needs to be called under e.triggersRegMu lock
		return registrationErr
	}

	// start listening for trigger events only if all registrations succeeded
	for idx, triggerEventCh := range eventChans {
		e.srvcEng.GoCtx(context.WithoutCancel(ctx), func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case event, isOpen := <-triggerEventCh:
					if !isOpen {
						return
					}
					triggerID := subscriptions[idx].Id
					eventID := event.Event.ID
					e.metrics.With(platform.KeyTriggerID, triggerID).IncrementTriggerEventReceivedCounter(ctx)
					e.logger().Debugw("Processing trigger event", "triggerID", triggerID, "eventID", eventID)
					if event.Err != nil {
						e.logger().Errorw("Received a trigger event with error, dropping", "triggerID", triggerID, "err", event.Err)
						tm := e.metrics.With(platform.KeyTriggerID, triggerID)
						tm.IncrementWorkflowTriggerEventErrorCounter(ctx)
						tm.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonTriggerResponseError)
						continue
					}

					routed := RoutedTriggerEvent{
						WorkflowID:     e.cfg.WorkflowID,
						TriggerCapID:   triggerID,
						TriggerIndex:   idx,
						ObservedAt:     e.cfg.Clock.Now(),
						SequenceNumber: 0,
						Event:          event,
					}

					if err := e.put(ctx, routed); err != nil {
						// Draining is expected during workflow deletion, so it logs at info rather than error level.
						if errors.Is(err, ErrEngineDraining) {
							e.logger().Infow("Dropping trigger event: engine draining", "triggerID", triggerID, "eventID", eventID)
						} else {
							e.logger().Errorw("Failed to put routed trigger event", "triggerID", triggerID, "eventID", eventID, "err", err)
						}
					}
				}
			}
		})
	}
	e.logger().Infow("All triggers registered successfully", "numTriggers", len(subscriptions), "triggerIDs", triggerCapIDs)
	e.metrics.IncrementWorkflowRegisteredCounter(ctx)
	e.cfg.Hooks.OnSubscribedToTriggers(triggerCapIDs)
	return nil
}

// NOTE: needs to be called under the triggersRegMu lock
func (e *Engine) unregisterAllTriggers(ctx context.Context) {
	failCount := 0
	for registrationID, trigger := range e.triggers {
		err := trigger.UnregisterTrigger(ctx, capabilities.TriggerRegistrationRequest{
			TriggerID: registrationID,
			Metadata: capabilities.RequestMetadata{
				WorkflowID:    e.cfg.WorkflowID,
				WorkflowDonID: e.localNode.Load().WorkflowDON.ID,
			},
			Payload: trigger.payload,
			Method:  trigger.method,
		})
		if err != nil {
			e.logger().Errorw("Failed to unregister trigger", "registrationId", registrationID, "err", err)
			failCount++
		}
	}
	e.logger().Infow("All triggers unregistered", "numTriggers", len(e.triggers), "failed", failCount)
	e.triggers = make(map[string]*triggerCapability)
}

func (e *Engine) Ack(ctx context.Context, triggerCapID, triggerRegistrationID, eventID string) error {
	e.logger().Infow("ACKing trigger event", "triggerRegistrationID", triggerRegistrationID, "eventID", eventID)

	tm := e.metrics.With(platform.KeyTriggerID, triggerCapID)

	e.triggersRegMu.Lock()
	trigger, ok := e.triggers[triggerRegistrationID]
	e.triggersRegMu.Unlock()

	if !ok {
		tm.IncrementTriggerEventAckFailureCounter(ctx)
		return fmt.Errorf("failed to find trigger %s", triggerRegistrationID)
	}
	err := trigger.AckEvent(ctx, triggerRegistrationID, eventID, trigger.method)
	if err != nil {
		tm.IncrementTriggerEventAckFailureCounter(ctx)
		return err
	}
	tm.IncrementTriggerEventAckSuccessCounter(ctx)
	return nil
}

// errObservedAtMissing guards put's deadline derivation: put computes the queue
// deadline from ObservedAt, so an unstamped event is rejected rather than given
// a garbage deadline.
var errObservedAtMissing = errors.New("trigger event ObservedAt not set")

// put enqueues a trigger event into the engine's internal queue.
func (e *Engine) put(ctx context.Context, event RoutedTriggerEvent) error {
	triggerID := event.TriggerCapID
	eventID := event.Event.Event.ID
	idx := event.TriggerIndex

	if e.Draining() {
		e.logger().Infow("Engine is draining, dropping trigger event before enqueue", "triggerID", triggerID, "eventID", eventID)
		tm := e.metrics.With(platform.KeyTriggerID, triggerID)
		tm.IncrementTriggerEventEnqueueDroppedCounter(ctx)
		tm.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonEnqueueDraining)
		e.cfg.Hooks.OnTriggerEventDropped(triggerID, eventID, "draining")
		return ErrEngineDraining
	}

	// Admission check: an external management layer (e.g. ShardFailoverManager)
	// decides whether this event should be processed by the engine.
	if err := e.cfg.Hooks.OnTriggerAdmission(ctx, event); err != nil {
		if errors.Is(err, ErrAdmissionCache) {
			e.logger().Infow("Trigger event cached for failover, not enqueuing", "triggerID", triggerID, "eventID", eventID)
			return err
		}
		// Denied: ACK and drop
		registrationID := TriggerRegistrationID(e.cfg.WorkflowID, event.TriggerIndex)
		if ackErr := e.cfg.TriggerAcknowledger.Ack(ctx, event.TriggerCapID, registrationID, eventID); ackErr != nil {
			e.logger().Errorw("failed to ACK trigger after admission denial", "eventID", eventID, "err", ackErr)
		}
		e.metrics.With(platform.KeyTriggerID, triggerID).IncrementTriggerEventDroppedTotal(ctx, "admission_denied")
		return err
	}

	// The deadline below is derived from ObservedAt, so reject events that
	// arrive without it rather than silently backfill a garbage deadline.
	if event.ObservedAt.IsZero() {
		e.logger().Errorw("Trigger event missing ObservedAt, dropping", "triggerID", triggerID, "eventID", eventID)
		e.metrics.With(platform.KeyTriggerID, triggerID).IncrementTriggerEventDroppedTotal(ctx, "observed_at_missing")
		return errObservedAtMissing
	}
	queueTimeout, err := e.cfg.LocalLimiters.TriggerEventQueueTimeout.Limit(ctx)
	if err != nil {
		if !limits.IsErrRecoverable(err) {
			// No value was resolved, so the deadline below would already be expired.
			e.logger().Errorw("Failed to get trigger event queue time limit with no usable value", "err", err)
			e.metrics.With(platform.KeyTriggerID, triggerID).
				IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonQueueAgeLimitReadFailed)
			return ErrEnqueueFailed
		}
		// A settings read failure is not a reason to drop a customer's trigger event:
		// the limiter still returns a usable timeout, so stamp the deadline and continue.
		e.logger().Errorw("Failed to get trigger event queue time limit; continuing with the value the limiter returned", "err", err)
		e.metrics.IncrementLimitReadFallbackCounter(ctx, cresettings.Default.PerWorkflow.TriggerEventQueueTimeout.Key)
	}
	event.Deadline = event.ObservedAt.Add(queueTimeout)

	if err := e.allTriggerEventsQueueCh.Put(ctx, event); err != nil {
		tm := e.metrics.With(platform.KeyTriggerID, triggerID)
		tm.IncrementTriggerEventEnqueueDroppedCounter(ctx)
		retErr := ErrEnqueueFailed
		if _, ok := errors.AsType[limits.ErrorQueueFull](err); ok {
			// queue full, drop the event
			e.logger().Errorw("Trigger event queue is full, dropping event", "triggerID", triggerID, "triggerIndex", idx, "err", err)
			tm.IncrementWorkflowTriggerEventQueueFullCounter(ctx)
			tm.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonEnqueueQueueFull)
			retErr = ErrQueueFull
		} else {
			tm.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonEnqueueFailed)
		}
		e.logger().Errorw("Failed to enqueue trigger event", "triggerID", triggerID, "triggerIndex", idx, "err", err)
		tm.IncrementWorkflowTriggerEventErrorCounter(ctx)
		return retErr
	}
	e.metrics.With(platform.KeyTriggerID, triggerID).IncrementTriggerEventEnqueuedCounter(ctx)
	e.logger().Debugw("Enqueued trigger event", "triggerID", triggerID, "eventID", eventID)

	return nil
}

// handleAllTriggerEvents drains the engine's trigger-event queue (populated by put method)
// and executes each event in turn.
func (e *Engine) handleAllTriggerEvents(ctx context.Context) {
	for {
		queueHead, err := e.allTriggerEventsQueueCh.Wait(ctx)
		if err != nil {
			return
		}
		eventID := queueHead.Event.Event.ID
		triggerMetricLabels := e.metrics.With(platform.KeyTriggerID, queueHead.TriggerCapID)
		if e.Draining() {
			triggerMetricLabels.IncrementTriggerEventDequeueDroppedCounter(ctx)
			triggerMetricLabels.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonDequeueDraining)
			e.cfg.Hooks.OnTriggerEventDropped(queueHead.TriggerCapID, eventID, "draining")
			e.logger().Infow("Engine is draining, stopping trigger handling loop", "eventID", eventID, "triggerID", queueHead.TriggerCapID)
			return
		}

		now := e.cfg.Clock.Now()
		eventAge := now.Sub(queueHead.ObservedAt)
		e.logger().Debugw("Popped a trigger event from the queue", "eventID", eventID, "eventAgeMs", eventAge.Milliseconds())
		triggerMetricLabels.RecordTriggerEventQueueWaitSeconds(ctx, eventAge.Seconds())
		if now.After(queueHead.Deadline) {
			e.logger().Warnw("Trigger event is too old, skipping execution", "triggerID", queueHead.TriggerCapID, "eventID", eventID, "eventAgeMs", eventAge.Milliseconds())
			triggerMetricLabels.IncrementTriggerEventExpiredCounter(ctx)
			triggerMetricLabels.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonExpired)
			continue
		}

		semWaitStart := e.cfg.Clock.Now()
		free, err := e.executionsSemaphore.Wait(ctx, 1) // block if too many concurrent workflow executions
		triggerMetricLabels.RecordExecutionSemaphoreWaitSeconds(ctx, e.cfg.Clock.Now().Sub(semWaitStart).Seconds())
		if err != nil {
			e.logger().Errorw("Failed to acquire executions semaphore", "err", err)
			triggerMetricLabels.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonExecutionSemaphoreWaitFailed)
			continue
		}

		e.srvcEng.GoCtx(context.WithoutCancel(ctx), func(ctx context.Context) {
			defer free()

			if err := e.ExecuteTrigger(ctx, queueHead); err != nil {
				// Dedup is an expected outcome (the event is handled, just not executed here), so it's logged at info rather than error level.
				if errors.Is(err, ErrDuplicateExecution) {
					e.logger().Infow("Skipping trigger event execution", "triggerID", queueHead.TriggerCapID, "eventID", queueHead.Event.Event.ID, "err", err)
				} else {
					e.logger().Errorw("Failed to execute trigger event", "triggerID", queueHead.TriggerCapID, "eventID", queueHead.Event.Event.ID, "err", err)
				}
			}
		})
	}
}
