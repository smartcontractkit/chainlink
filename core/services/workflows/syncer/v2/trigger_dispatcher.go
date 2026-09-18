package v2

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

// TriggerDispatcher owns trigger registration, the trigger handle map, event
// channel reading, and acknowledgement for all workflows on this node.
type TriggerDispatcher interface {
	services.Service
	RegisterTriggers(ctx context.Context, cre contexts.CRE, params RegistrationParams, subs []*sdkpb.TriggerSubscription) ([]string, error)
	// UnregisterTriggers stops event ingress for the workflow immediately, then
	// releases its retained trigger handles once the engine reports no more
	// active executions (or a bounded timeout elapses). Safe to call even if
	// the workflow was never registered.
	UnregisterTriggers(workflowID string) error
	// Ack acknowledges a trigger event by resolving the registration's handle
	// and calling AckEvent on the trigger capability. Engines call this after
	// execution starts, on duplicate executions, and on shard-ownership
	// denials — the point at which the event is fully handled and must not be
	// redelivered.
	Ack(ctx context.Context, workflowID, triggerCapID, triggerRegistrationID, eventID string) error
}

var (
	_ TriggerDispatcher     = (*triggerDispatcher)(nil)
	_ v2.Acknowledger       = (*triggerDispatcher)(nil)
	_ WorkflowLimitReporter = (*triggerDispatcher)(nil)
)

// WorkflowLimitReporter records workflow count limit rejections, surfaced by
// the syncer when it acquires the limit on an engine's behalf.
type WorkflowLimitReporter interface {
	ReportWorkflowLimitPerOwner(ctx context.Context)
	ReportWorkflowLimitGlobal(ctx context.Context)
}

// RegistrationParams carries the workflow-scoped metadata stamped into every
// TriggerRegistrationRequest for one workflow. contexts.CRE only holds tenant
// identity (org/owner/workflow); the syncer, which owns these values, supplies
// the rest alongside the subscriptions.
type RegistrationParams struct {
	WorkflowName                  string // decoded workflow name
	WorkflowTag                   string
	WorkflowDonID                 uint32
	WorkflowRegistryAddress       string
	WorkflowRegistryChainSelector string
}

// eventSink is the delivery surface the dispatcher uses to feed engines. It
// is the engine's transitional Put method; the dispatcher never holds an
// engine reference, it resolves one from the registry at delivery time and
// only needs this narrow interface.
type eventSink interface {
	Put(ctx context.Context, event v2.RoutedTriggerEvent) error
}

// activeExecutionsReporter is how UnregisterTriggers waits for a workflow's
// in-flight executions to finish before releasing its trigger handles. Every
// engine implementation already exposes this (it backs Drain/DrainableService
// on the syncer side); the dispatcher only needs this one method of it, and
// resolves it from the registry the same way eventSink is resolved.
type activeExecutionsReporter interface {
	ActiveExecutions() int32
}

const (
	// releaseHandlesPollInterval is how often UnregisterTriggers rechecks
	// ActiveExecutions while waiting to release a workflow's trigger handles.
	releaseHandlesPollInterval = 250 * time.Millisecond
	// releaseHandlesTimeout bounds that wait so a workflow whose engine never
	// reaches zero active executions (e.g. it was already popped from the
	// registry mid-drain) doesn't leak its trigger handles forever.
	releaseHandlesTimeout = 5 * time.Minute
)

// triggerDispatcher is the implementation of TriggerDispatcher.
type triggerDispatcher struct {
	services.Service
	eng *services.Engine

	lggr     logger.Logger
	capReg   core.CapabilitiesRegistry
	registry *EngineRegistry    // resolve engines at delivery time
	regTime  limits.TimeLimiter // bounds total trigger registration time
	metrics  *monitoring.WorkflowsMetricLabeler
	clock    clockwork.Clock
	// pinnedConfigVersion is the DON config version stamped into registration
	// metadata. It is pinned to 1 so that config updates on the registry do not
	// force forwarder contract updates: the version is included in every report
	// but is irrelevant to forwarder validation, which only checks DON ID and
	// the set of signer public keys.
	pinnedConfigVersion uint32

	mu        sync.RWMutex
	workflows map[types.WorkflowID]*workflowTriggers
}

// workflowTriggers is everything the dispatcher owns for one workflow: the
// tenant context, registration metadata, and the handle map.
type workflowTriggers struct {
	cre     contexts.CRE
	params  RegistrationParams
	handles map[string]*triggerHandle // registrationID -> handle
}

// triggerHandle is a registered trigger capability plus the registration
// payload/method needed to unregister it.
type triggerHandle struct {
	capabilities.TriggerCapability
	payload *anypb.Any
	method  string
}

// NewTriggerDispatcher returns a dispatcher wired to the given engine
// registry. regTime bounds the total time spent registering a workflow's
// triggers.
func NewTriggerDispatcher(lggr logger.Logger, capReg core.CapabilitiesRegistry, registry *EngineRegistry, regTime limits.TimeLimiter, metrics *monitoring.WorkflowsMetricLabeler, clock clockwork.Clock) TriggerDispatcher {
	d := &triggerDispatcher{
		lggr:                logger.Named(lggr, "TriggerDispatcher"),
		capReg:              capReg,
		registry:            registry,
		regTime:             regTime,
		metrics:             metrics,
		clock:               clock,
		pinnedConfigVersion: 1,
		workflows:           make(map[types.WorkflowID]*workflowTriggers),
	}
	d.Service, d.eng = services.Config{
		Name:  "TriggerDispatcher",
		Start: d.start,
		Close: d.close,
	}.NewServiceEngine(d.lggr)
	return d
}

// ReportWorkflowLimitPerOwner records a per-owner workflow count limit rejection.
func (d *triggerDispatcher) ReportWorkflowLimitPerOwner(ctx context.Context) {
	d.metrics.IncrementWorkflowLimitPerOwnerCounter(ctx)
}

// ReportWorkflowLimitGlobal records a global workflow count limit rejection.
func (d *triggerDispatcher) ReportWorkflowLimitGlobal(ctx context.Context) {
	d.metrics.IncrementWorkflowLimitGlobalCounter(ctx)
}

// start is a no-op: the dispatcher has no background work of its own — reader
// goroutines are started per subscription by RegisterTriggers and tracked by
// the embedded services.Engine. The method exists to satisfy the
// services.Service Start hook contract.
func (d *triggerDispatcher) start(context.Context) error { return nil }

func (d *triggerDispatcher) close() error {
	// Reader goroutines are owned by d.eng and stopped when it stops. Handles
	// are dropped with the dispatcher; per-workflow cleanup (unregister, drain,
	// handle release) is driven by the syncer via UnregisterTriggers.
	d.mu.Lock()
	defer d.mu.Unlock()
	d.workflows = make(map[types.WorkflowID]*workflowTriggers)
	return nil
}

// RegisterTriggers validates and registers the given subscriptions with the
// capability registry, retains the handles and the workflow's tenant context,
// starts one reader goroutine per subscription, and returns the registered
// trigger capability IDs. On any registration failure it rolls back the
// successful registrations.
func (d *triggerDispatcher) RegisterTriggers(ctx context.Context, cre contexts.CRE, params RegistrationParams, subs []*sdkpb.TriggerSubscription) ([]string, error) {
	wid, err := types.WorkflowIDFromHex(cre.Workflow)
	if err != nil {
		return nil, fmt.Errorf("invalid workflowID in CRE context: %w", err)
	}

	// check if all requested triggers exist in the registry
	triggers := make([]capabilities.TriggerCapability, 0, len(subs))
	for _, sub := range subs {
		_, labels, _ := capabilities.ParseID(sub.Id)
		chainSelector, err2 := capabilities.ChainSelectorLabel(labels)
		if err2 != nil {
			return nil, fmt.Errorf("invalid chain selector for ID %s: %w", sub.Id, err2)
		}
		_ = chainSelector // chain access check moves with the limiter split (CRE-6177)

		triggerCap, triggerErr := d.capReg.GetTrigger(ctx, sub.Id)
		if triggerErr != nil {
			return nil, fmt.Errorf("trigger capability not found: %w", triggerErr)
		}
		triggers = append(triggers, triggerCap)
	}

	// register to all triggers concurrently
	regCtx, regCancel, err := d.regTime.WithTimeout(ctx)
	if err != nil {
		return nil, err
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

	resultsCh := make(chan triggerRegResult, len(subs))
	g, gCtx := errgroup.WithContext(regCtx)

	// Launch concurrent trigger registrations
	for i, sub := range subs {
		triggerCap := triggers[i]
		g.Go(func() error {
			registrationID := v2.TriggerRegistrationID(cre.Workflow, i)
			args := []any{"triggerID", sub.Id, "method", sub.Method}
			if sub.Payload != nil {
				args = append(args, "payload", protojson.Format(sub.Payload))
			}
			d.lggr.Infow("Registering trigger", args...)
			metadata := capabilities.RequestMetadata{
				WorkflowID:                    cre.Workflow,
				WorkflowOwner:                 cre.Owner,
				WorkflowName:                  params.WorkflowName,
				DecodedWorkflowName:           params.WorkflowName,
				WorkflowTag:                   params.WorkflowTag,
				WorkflowDonID:                 params.WorkflowDonID,
				WorkflowDonConfigVersion:      d.pinnedConfigVersion,
				ReferenceID:                   fmt.Sprintf("trigger_%d", i),
				WorkflowRegistryChainSelector: params.WorkflowRegistryChainSelector,
				WorkflowRegistryAddress:       params.WorkflowRegistryAddress,
				EngineVersion:                 platform.ValueWorkflowVersionV2,
				// no WorkflowExecutionID needed (or available at this stage)
			}
			var creGetter settings.Getter
			propagateOrgIDMeta, _ := cresettings.Default.PropagateOrgIDInRequestMetadata.GetOrDefault(gCtx, creGetter)
			if propagateOrgIDMeta && cre.Org != "" {
				metadata.OrgID = cre.Org
			}
			triggerEventCh, regErr := triggerCap.RegisterTrigger(gCtx, capabilities.TriggerRegistrationRequest{
				TriggerID: registrationID,
				Metadata:  metadata,
				Payload:   sub.Payload,
				Method:    sub.Method,
			})
			if regErr != nil {
				d.lggr.Errorw("Trigger registration failed", "triggerID", sub.Id, "err", regErr)
				d.metrics.With(platform.KeyTriggerID, sub.Id).IncrementRegisterTriggerFailureCounter(gCtx)
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

	// Collect results into the per-workflow state and the flat ACK index.
	eventChans := make([]<-chan capabilities.TriggerResponse, len(subs))
	triggerCapIDs := make([]string, len(subs))
	wt := &workflowTriggers{cre: cre, params: params, handles: make(map[string]*triggerHandle, len(subs))}

	for result := range resultsCh {
		wt.handles[result.registrationID] = &triggerHandle{
			TriggerCapability: result.triggerCap,
			payload:           result.payload,
			method:            result.method,
		}
		eventChans[result.index] = result.eventCh
		triggerCapIDs[result.index] = result.triggerCapID
	}

	// If any registration failed, unregister successful ones and return error
	if registrationErr != nil {
		d.lggr.Errorw("One or more trigger registrations failed - reverting all", "err", registrationErr)
		d.unregisterAll(ctx, wid, wt)
		return nil, registrationErr
	}

	d.mu.Lock()
	// Replace any prior state (e.g. re-registration after a config update).
	d.workflows[wid] = wt
	d.mu.Unlock()

	// start listening for trigger events only if all registrations succeeded
	for idx, triggerEventCh := range eventChans {
		d.startReader(ctx, wid, idx, subs[idx], triggerEventCh)
	}
	d.lggr.Infow("All triggers registered successfully", "numTriggers", len(subs), "triggerIDs", triggerCapIDs)
	d.metrics.IncrementWorkflowRegisteredCounter(ctx)
	return triggerCapIDs, nil
}

// startReader runs one reader goroutine per subscription. It resolves the
// engine from the registry at delivery time and never holds an engine
// reference; if the engine is gone the reader exits.
func (d *triggerDispatcher) startReader(ctx context.Context, wid types.WorkflowID, idx int, sub *sdkpb.TriggerSubscription, triggerEventCh <-chan capabilities.TriggerResponse) {
	d.eng.GoCtx(context.WithoutCancel(ctx), func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case event, isOpen := <-triggerEventCh:
				if !isOpen {
					return
				}
				triggerID := sub.Id
				eventID := event.Event.ID
				d.metrics.With(platform.KeyTriggerID, triggerID).IncrementTriggerEventReceivedCounter(ctx)
				d.lggr.Debugw("Processing trigger event", "triggerID", triggerID, "eventID", eventID)
				if event.Err != nil {
					d.lggr.Errorw("Received a trigger event with error, dropping", "triggerID", triggerID, "err", event.Err)
					tm := d.metrics.With(platform.KeyTriggerID, triggerID)
					tm.IncrementWorkflowTriggerEventErrorCounter(ctx)
					tm.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonTriggerResponseError)
					continue
				}

				// resolve the engine at DELIVERY time — never hold a reference.
				svc, found := d.registry.Get(wid)
				if !found {
					return // engine gone; reader exits
				}
				sink, ok := svc.Service.(eventSink)
				if !ok {
					d.lggr.Errorw("Engine does not accept trigger events", "workflowID", wid)
					return
				}

				routed := v2.RoutedTriggerEvent{
					WorkflowID:   wid.String(),
					TriggerCapID: triggerID,
					TriggerIndex: idx,
					ObservedAt:   d.clock.Now(),
					Event:        event,
				}
				if err := sink.Put(ctx, routed); err != nil {
					// Draining is expected during workflow deletion, so it logs at info rather than error level.
					if errors.Is(err, v2.ErrEngineDraining) {
						d.lggr.Infow("Dropping trigger event: engine draining", "triggerID", triggerID, "eventID", eventID)
					} else {
						d.lggr.Errorw("Failed to put routed trigger event", "triggerID", triggerID, "eventID", eventID, "err", err)
					}
				}
			}
		}
	})
}

// Ack acknowledges a trigger event by resolving the registration's handle
// and calling AckEvent on the trigger capability. Engines call this after
// execution starts, on duplicate executions, and on shard-ownership denials —
// the point at which the event is fully handled and must not be redelivered.
func (d *triggerDispatcher) Ack(ctx context.Context, workflowID, triggerCapID, triggerRegistrationID, eventID string) error {
	d.lggr.Infow("ACKing trigger event", "triggerRegistrationID", triggerRegistrationID, "eventID", eventID)

	tm := d.metrics.With(platform.KeyTriggerID, triggerCapID)

	wid, err := types.WorkflowIDFromHex(workflowID)
	if err != nil {
		tm.IncrementTriggerEventAckFailureCounter(ctx)
		return fmt.Errorf("invalid workflowID: %w", err)
	}

	d.mu.RLock()
	var handle *triggerHandle
	if wt, ok := d.workflows[wid]; ok {
		handle = wt.handles[triggerRegistrationID]
	}
	d.mu.RUnlock()

	if handle == nil {
		tm.IncrementTriggerEventAckFailureCounter(ctx)
		return fmt.Errorf("failed to find trigger %s for workflow %s", triggerRegistrationID, workflowID)
	}
	if err := handle.AckEvent(ctx, triggerRegistrationID, eventID, handle.method); err != nil {
		tm.IncrementTriggerEventAckFailureCounter(ctx)
		return err
	}
	tm.IncrementTriggerEventAckSuccessCounter(ctx)
	return nil
}

// UnregisterTriggers unregisters the workflow's triggers with the capability
// registry, stopping event ingress immediately. The handles are retained —
// so an execution already in flight can still Ack — until the engine reports
// no more active executions (or releaseHandlesTimeout elapses), at which
// point they're released in the background. Safe to call even if the
// workflow was never registered.
func (d *triggerDispatcher) UnregisterTriggers(workflowID string) error {
	wid, err := types.WorkflowIDFromHex(workflowID)
	if err != nil {
		return fmt.Errorf("invalid workflowID: %w", err)
	}

	d.mu.Lock()
	wt, ok := d.workflows[wid]
	d.mu.Unlock()
	if !ok {
		return fmt.Errorf("no triggers registered for workflow %s", workflowID)
	}

	// Unregister with the capability registry outside the lock.
	ctx := context.Background()
	for registrationID, handle := range wt.handles {
		if unregErr := handle.UnregisterTrigger(ctx, capabilities.TriggerRegistrationRequest{
			TriggerID: registrationID,
			Metadata: capabilities.RequestMetadata{
				WorkflowID:    workflowID,
				WorkflowDonID: wt.params.WorkflowDonID,
			},
			Payload: handle.payload,
			Method:  handle.method,
		}); unregErr != nil {
			d.lggr.Errorw("Failed to unregister trigger", "registrationId", registrationID, "err", unregErr)
		}
	}
	d.metrics.IncrementWorkflowUnregisteredCounter(ctx)

	d.eng.GoCtx(context.WithoutCancel(ctx), func(ctx context.Context) {
		d.releaseHandlesWhenDrained(ctx, wid)
	})
	return nil
}

// releaseHandlesWhenDrained waits for the workflow's engine to report zero
// active executions — resolving it from the registry the same way the reader
// resolves eventSink, since the dispatcher holds no engine reference — and
// then drops the retained handles. It gives up and releases anyway, with a
// warning, if releaseHandlesTimeout elapses or the engine is no longer in the
// registry (already popped, so nothing is left to wait for).
func (d *triggerDispatcher) releaseHandlesWhenDrained(ctx context.Context, wid types.WorkflowID) {
	ctx, cancel := context.WithTimeout(ctx, releaseHandlesTimeout)
	defer cancel()
	ticker := time.NewTicker(releaseHandlesPollInterval)
	defer ticker.Stop()

waitForDrain:
	for {
		svc, found := d.registry.Get(wid)
		if !found {
			break waitForDrain
		}
		reporter, ok := svc.Service.(activeExecutionsReporter)
		if !ok || reporter.ActiveExecutions() == 0 {
			break waitForDrain
		}
		select {
		case <-ctx.Done():
			d.lggr.Warnw("Timed out waiting for active executions to drain; releasing trigger handles anyway", "workflowID", wid)
			break waitForDrain
		case <-ticker.C:
		}
	}

	d.mu.Lock()
	delete(d.workflows, wid)
	d.mu.Unlock()
}

// unregisterAll rolls back every successful registration when one or more
// registrations failed. Unlike UnregisterTriggers it also drops the handles,
// because a failed RegisterTriggers leaves no executions in flight.
func (d *triggerDispatcher) unregisterAll(ctx context.Context, wid types.WorkflowID, wt *workflowTriggers) {
	for registrationID, handle := range wt.handles {
		if err := handle.UnregisterTrigger(ctx, capabilities.TriggerRegistrationRequest{
			TriggerID: registrationID,
			Metadata: capabilities.RequestMetadata{
				WorkflowID:    wid.String(),
				WorkflowDonID: wt.params.WorkflowDonID,
			},
			Payload: handle.payload,
			Method:  handle.method,
		}); err != nil {
			d.lggr.Errorw("Failed to unregister trigger", "registrationId", registrationID, "err", err)
		}
	}
	d.mu.Lock()
	delete(d.workflows, wid)
	d.mu.Unlock()
}
