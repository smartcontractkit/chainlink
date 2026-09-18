package v2

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

// TriggerCoordinator owns trigger registration, the trigger handle map, event
// channel reading, and acknowledgement for all workflows on this node.
type TriggerCoordinator interface {
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
	_ TriggerCoordinator = (*triggerCoordinator)(nil)
	_ v2.Acknowledger    = (*triggerCoordinator)(nil)
)

// RegistrationParams carries the workflow-scoped metadata stamped into every
// TriggerRegistrationRequest for one workflow. contexts.CRE only holds tenant
// identity (org/owner/workflow); the syncer, which owns these values, supplies
// the rest alongside the subscriptions.
type RegistrationParams struct {
	// WorkflowName is nil-safe: a nil value sends both the hex-encoded and
	// decoded workflow name empty, rather than panicking on a nil interface
	// call (see v2.RegisterWorkflowTriggers).
	WorkflowName                  types.WorkflowName
	WorkflowTag                   string
	WorkflowDonID                 uint32
	WorkflowRegistryAddress       string
	WorkflowRegistryChainSelector string
}

// eventSink is the delivery surface the coordinator uses to feed engines. It
// is the engine's transitional Put method; the coordinator never holds an
// engine reference, it resolves one from the registry at delivery time and
// only needs this narrow interface.
type eventSink interface {
	Put(ctx context.Context, event v2.RoutedTriggerEvent) error
}

// activeExecutionsReporter is how UnregisterTriggers waits for a workflow's
// in-flight executions to finish before releasing its trigger handles. Every
// engine implementation already exposes this (it backs Drain/DrainableService
// on the syncer side); the coordinator only needs this one method of it, and
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

// triggerCoordinator is the implementation of TriggerCoordinator.
type triggerCoordinator struct {
	services.Service
	eng *services.Engine

	lggr     logger.Logger
	capReg   core.CapabilitiesRegistry
	registry *EngineRegistry    // resolve engines at delivery time
	regTime  limits.TimeLimiter // bounds total trigger registration time
	metrics  *monitoring.WorkflowsMetricLabeler
	clock    clockwork.Clock

	mu        sync.RWMutex
	workflows map[types.WorkflowID]*workflowTriggers
}

// workflowTriggers is everything the coordinator owns for one workflow: the
// tenant context, registration metadata, and the handle map.
type workflowTriggers struct {
	cre     contexts.CRE
	params  RegistrationParams
	handles map[string]*v2.TriggerHandle // registrationID -> handle
}

// NewTriggerCoordinator returns a coordinator wired to the given engine
// registry. regTime bounds the total time spent registering a workflow's
// triggers.
func NewTriggerCoordinator(lggr logger.Logger, capReg core.CapabilitiesRegistry, registry *EngineRegistry, regTime limits.TimeLimiter, metrics *monitoring.WorkflowsMetricLabeler, clock clockwork.Clock) TriggerCoordinator {
	d := &triggerCoordinator{
		lggr:      logger.Named(lggr, "TriggerCoordinator"),
		capReg:    capReg,
		registry:  registry,
		regTime:   regTime,
		metrics:   metrics,
		clock:     clock,
		workflows: make(map[types.WorkflowID]*workflowTriggers),
	}
	d.Service, d.eng = services.Config{
		Name:  "TriggerCoordinator",
		Start: d.start,
		Close: d.close,
	}.NewServiceEngine(d.lggr)
	return d
}

// start is a no-op: the coordinator has no background work of its own — reader
// goroutines are started per subscription by RegisterTriggers and tracked by
// the embedded services.Engine. The method exists to satisfy the
// services.Service Start hook contract.
func (d *triggerCoordinator) start(context.Context) error { return nil }

func (d *triggerCoordinator) close() error {
	// Reader goroutines are owned by d.eng and stopped when it stops. Handles
	// are dropped with the coordinator; per-workflow cleanup (unregister, drain,
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
func (d *triggerCoordinator) RegisterTriggers(ctx context.Context, cre contexts.CRE, params RegistrationParams, subs []*sdkpb.TriggerSubscription) ([]string, error) {
	wid, err := types.WorkflowIDFromHex(cre.Workflow)
	if err != nil {
		return nil, fmt.Errorf("invalid workflowID in CRE context: %w", err)
	}

	triggerCapIDs, handles, eventChans, err := v2.RegisterWorkflowTriggers(ctx, v2.TriggerRegistrationDeps{
		CapRegistry: d.capReg,
		RegTimeout:  d.regTime,
		// ChainAllowed and Settings are both nil for now: chain access
		// enforcement moves with the limiter split (CRE-6177), and there's
		// no dynamic settings source wired up here yet. RegisterWorkflowTriggers
		// treats both as nil-safe, matching this coordinator's prior behavior.
		Logger:  d.lggr,
		Metrics: d.metrics,
	}, v2.TriggerRegistrationMetadata{
		WorkflowID:                    cre.Workflow,
		WorkflowOwner:                 cre.Owner,
		WorkflowName:                  params.WorkflowName,
		WorkflowTag:                   params.WorkflowTag,
		WorkflowDonID:                 params.WorkflowDonID,
		WorkflowRegistryChainSelector: params.WorkflowRegistryChainSelector,
		WorkflowRegistryAddress:       params.WorkflowRegistryAddress,
		OrgID:                         cre.Org,
	}, subs)
	if err != nil {
		return nil, err
	}

	d.mu.Lock()
	// Replace any prior state (e.g. re-registration after a config update).
	d.workflows[wid] = &workflowTriggers{cre: cre, params: params, handles: handles}
	d.mu.Unlock()

	// start listening for trigger events only if all registrations succeeded
	for idx, triggerEventCh := range eventChans {
		d.startReader(ctx, wid, idx, subs[idx], triggerEventCh)
	}
	return triggerCapIDs, nil
}

// startReader runs one reader goroutine per subscription. It resolves the
// engine from the registry on every event rather than holding a reference —
// see the deliver closure below for why. A missing or non-conforming engine
// is just another delivery failure to RunTriggerReader, not a reason for the
// reader itself to exit.
func (d *triggerCoordinator) startReader(ctx context.Context, wid types.WorkflowID, idx int, sub *sdkpb.TriggerSubscription, triggerEventCh <-chan capabilities.TriggerResponse) {
	// deliver resolves the engine from the registry on every call — never a
	// held reference — since it's invoked once per event by RunTriggerReader.
	// An engine that's gone, or that doesn't accept trigger events, is just
	// another delivery failure to RunTriggerReader: logged and retried on the
	// next event, not a reason for the reader itself to exit early.
	deliver := func(ctx context.Context, event v2.RoutedTriggerEvent) error {
		svc, found := d.registry.Get(wid)
		if !found {
			return fmt.Errorf("no engine registered for workflow %s", wid)
		}
		sink, ok := svc.Service.(eventSink)
		if !ok {
			return fmt.Errorf("engine for workflow %s does not accept trigger events", wid)
		}
		return sink.Put(ctx, event)
	}
	d.eng.GoCtx(context.WithoutCancel(ctx), func(ctx context.Context) {
		v2.RunTriggerReader(ctx, d.lggr, d.metrics, d.clock, wid.String(), sub.Id, idx, triggerEventCh, deliver)
	})
}

// Ack acknowledges a trigger event by resolving the registration's handle
// and calling AckEvent on the trigger capability. Engines call this after
// execution starts, on duplicate executions, and on shard-ownership denials —
// the point at which the event is fully handled and must not be redelivered.
func (d *triggerCoordinator) Ack(ctx context.Context, workflowID, triggerCapID, triggerRegistrationID, eventID string) error {
	wid, err := types.WorkflowIDFromHex(workflowID)
	if err != nil {
		d.metrics.With(platform.KeyTriggerID, triggerCapID).IncrementTriggerEventAckFailureCounter(ctx)
		return fmt.Errorf("invalid workflowID: %w", err)
	}

	d.mu.RLock()
	var handle *v2.TriggerHandle
	if wt, ok := d.workflows[wid]; ok {
		handle = wt.handles[triggerRegistrationID]
	}
	d.mu.RUnlock()

	// handle is resolved above (rather than left to v2.AckTriggerHandle) so
	// the not-found error here can keep the extra "for workflow %s" context
	// Engine's equivalent has no use for (an Engine only ever has one
	// workflow).
	if err := v2.AckTriggerHandle(ctx, d.lggr, d.metrics, triggerCapID, triggerRegistrationID, eventID, handle); err != nil {
		if handle == nil {
			return fmt.Errorf("failed to find trigger %s for workflow %s", triggerRegistrationID, workflowID)
		}
		return err
	}
	return nil
}

// UnregisterTriggers unregisters the workflow's triggers with the capability
// registry, stopping event ingress immediately. The handles are retained —
// so an execution already in flight can still Ack — until the engine reports
// no more active executions (or releaseHandlesTimeout elapses), at which
// point they're released in the background. Safe to call even if the
// workflow was never registered.
func (d *triggerCoordinator) UnregisterTriggers(workflowID string) error {
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
	failCount := v2.UnregisterTriggerHandles(ctx, d.lggr, workflowID, wt.params.WorkflowDonID, wt.handles)
	d.lggr.Infow("All triggers unregistered", "numTriggers", len(wt.handles), "failed", failCount)
	d.metrics.IncrementWorkflowUnregisteredCounter(ctx)

	d.eng.GoCtx(context.WithoutCancel(ctx), func(ctx context.Context) {
		d.releaseHandlesWhenDrained(ctx, wid)
	})
	return nil
}

// releaseHandlesWhenDrained waits for the workflow's engine to report zero
// active executions — resolving it from the registry the same way the reader
// resolves eventSink, since the coordinator holds no engine reference — and
// then drops the retained handles. It gives up and releases anyway, with a
// warning, if releaseHandlesTimeout elapses or the engine is no longer in the
// registry (already popped, so nothing is left to wait for).
func (d *triggerCoordinator) releaseHandlesWhenDrained(ctx context.Context, wid types.WorkflowID) {
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
