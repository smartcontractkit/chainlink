package triggers

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

// ErrWorkflowNotCoordinated is returned by UnregisterTriggers for
// a workflowID the coordinator never registered.
var ErrWorkflowNotCoordinated = errors.New("workflow not registered with the trigger coordinator")

const (
	defaultDrainTimeout = 10 * time.Minute
	shutdownTimeout     = 5 * time.Second

	// pinnedWorkflowDonConfigVersion mirrors v2's pin to 1, so config updates on
	// the registry don't force forwarder contract updates.
	pinnedWorkflowDonConfigVersion = 1
)

// RegistrationParams is the per-workflow metadata RegisterTriggers needs to
// build each capability's RequestMetadata.
type RegistrationParams struct {
	WorkflowOwner                 string
	WorkflowName                  types.WorkflowName
	WorkflowTag                   string
	WorkflowDonID                 uint32
	WorkflowRegistryChainSelector string
	WorkflowRegistryAddress       string
}

// Coordinator owns trigger registration, the trigger handle map, event
// delivery, and acknowledgement for every workflow running the coordinated engine.
// It is a node-level singleton, started and stopped with the syncer.
type Coordinator interface {
	services.Service
	Acknowledger

	// RegisterTriggers calls subscriber.Subscribe to obtain the engine's trigger
	// subscriptions, registers them with the capability registry, retains the
	// resulting handles and reader goroutines, and returns the registered
	// trigger capability IDs. On any failure it unregisters what it already
	// registered for this call and returns the error.
	// Partial registration is never left behind.
	RegisterTriggers(ctx context.Context, subscriber Subscriber, params RegistrationParams) ([]string, error)

	// UnregisterTriggers stops ingress for workflowID immediately (unregisters with the capability registry)
	// and cleans up the handle map once the engine has been drained and closed, so an execution already in flight
	// can still resolve its handle to ACK. It also frees the workflow-count limit slot acquired at registration.
	//
	// Returns ErrWorkflowNotCoordinated if workflowID was never registered here.
	// A failed capability unregistration is returned and retried on the next call.
	UnregisterTriggers(workflowID string) error
}

// RegisteredEngine is what EngineRegistry.Get returns: a sink the coordinator
// can deliver events to, plus the flag callers must check before doing so.
type RegisteredEngine interface {
	EventSink

	// IsCoordinated is true if the engine does not manage its own
	// trigger registration, trigger dequeuing, execution or acknowledgement.
	IsCoordinated() bool
}

// EngineRegistry is the coordinator's read-only view of running engines. It is
// how the coordinator resolves a workflow's engine at delivery time; the
// coordinator keeps no engine map of its own.
type EngineRegistry interface {
	// Get returns the engine for workflowID, or false if none is registered
	// (e.g. the workflow was unregistered while events were still queued).
	Get(workflowID types.WorkflowID) (RegisteredEngine, bool)
}

type coordinator struct {
	services.Service
	eng  *services.Engine
	lggr logger.Logger

	// deps carries the node-level logger and metrics; RegisterTriggers
	// overrides both with the workflow-scoped ones.
	deps           RegisterDeps
	engines        EngineRegistry
	workflowLimits limits.ResourceLimiter[int]
	clock          clockwork.Clock

	drainTimeout time.Duration

	workflows *workflowRegistry
}

// coordinatedWorkflow represents the state of a workflow's triggers within the coordinator.
// A successful call to RegisterTriggers creates a coordinatedWorkflow.
// It tracks the workflow ID, the context for cancellation, the handles for each trigger,
// and the status of unregistration and release.
type coordinatedWorkflow struct {
	wid   types.WorkflowID
	cre   contexts.CRE
	donID uint32
	// cancel stops reading all trigger events for this coordinated workflow instance.
	// A workflow can be re-registered to the same triggers multiple times. This cancellation,
	// is independent of subsequent trigger registrations.
	cancel  context.CancelFunc
	handles map[string]*Handle

	// readers tracks the reader goroutines for this registration. A reader
	// only exits once its own delivery has returned, so once Wait() returns,
	// nothing started by this registration can still ACK.
	readers sync.WaitGroup

	unregistered atomic.Bool // capability-side unregistration succeeded
	releasing    atomic.Bool // release waiter spawned

	lggr    logger.Logger
	metrics *monitoring.WorkflowsMetricLabeler
}

func NewCoordinator(
	deps RegisterDeps,
	engineRegistry EngineRegistry,
	workflowLimits limits.ResourceLimiter[int],
	clock clockwork.Clock,
) Coordinator {
	c := &coordinator{
		lggr:           logger.Named(deps.Logger, "TriggerCoordinator"),
		deps:           deps,
		engines:        engineRegistry,
		workflowLimits: workflowLimits,
		clock:          clock,
		drainTimeout:   defaultDrainTimeout,
		workflows:      newWorkflowRegistry(),
	}

	c.Service, c.eng = services.Config{
		Name:  "TriggerCoordinator",
		Close: c.close,
	}.NewServiceEngine(c.lggr)

	return c
}

// close runs after the reader and release goroutines exit.
// It unregisters any triggers still in a pending state before returning.
func (c *coordinator) close() error {
	pending := c.workflows.pending()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	var errs error
	for workflowID, cw := range pending {
		cw.cancel()
		if failCount := Unregister(contexts.WithCRE(ctx, cw.cre), cw.lggr, workflowID, cw.donID, cw.handles); failCount > 0 {
			errs = errors.Join(errs, fmt.Errorf("workflow %s: failed to unregister %d of %d triggers", workflowID, failCount, len(cw.handles)))
		}
	}
	return errs
}

func (c *coordinator) RegisterTriggers(ctx context.Context, subscriber Subscriber, params RegistrationParams) (triggerCapIDs []string, err error) {
	cre := subscriber.Tenant()
	workflowID := cre.Workflow
	wid, idErr := types.WorkflowIDFromHex(workflowID)
	if idErr != nil {
		return nil, fmt.Errorf("invalid workflow id %q: %w", workflowID, idErr)
	}
	ctx = contexts.WithCRE(ctx, cre)

	lggr := logger.With(c.lggr, "workflowID", workflowID)
	wfMetrics := c.deps.Metrics.With(
		platform.KeyWorkflowID, workflowID,
		platform.KeyWorkflowOwner, params.WorkflowOwner,
		platform.KeyWorkflowName, params.WorkflowName.String(),
		platform.KeyOrganizationID, cre.Org,
	)

	// Registration IDs derive from the workflowID, so a leftover registration
	// must be unregistered first: unregistering it later would remove this one.
	if err = c.UnregisterTriggers(workflowID); err != nil && !errors.Is(err, ErrWorkflowNotCoordinated) {
		lggr.Errorw("Failed to unregister previous trigger registration", "err", err)
	}

	if err = c.useWorkflowLimit(ctx, lggr, wfMetrics); err != nil {
		return nil, err
	}
	// Any failure return frees the limit slot because of the named err.
	// Keep returns explicit (no naked return) so the defer always sees the final err.
	defer func() {
		if err != nil {
			c.freeWorkflowLimit(ctx, lggr)
		}
	}()

	subs, err := subscriber.Subscribe(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to triggers: %w", err)
	}

	deps := c.deps
	deps.Logger = lggr
	deps.Metrics = wfMetrics
	triggerCapIDs, handles, eventChans, err := Register(ctx, deps, RegisterMetadata{
		WorkflowID:                    workflowID,
		WorkflowOwner:                 params.WorkflowOwner,
		WorkflowName:                  params.WorkflowName,
		WorkflowTag:                   params.WorkflowTag,
		WorkflowDonID:                 params.WorkflowDonID,
		WorkflowDonConfigVersion:      pinnedWorkflowDonConfigVersion,
		WorkflowRegistryChainSelector: params.WorkflowRegistryChainSelector,
		WorkflowRegistryAddress:       params.WorkflowRegistryAddress,
		OrgID:                         cre.Org,
	}, subs)
	if err != nil {
		return nil, err
	}

	readerCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	cw := &coordinatedWorkflow{
		wid:     wid,
		cre:     cre,
		donID:   params.WorkflowDonID,
		cancel:  cancel,
		handles: handles,
		lggr:    lggr,
		metrics: wfMetrics,
	}

	// Set before publishing cw: an UnregisterTriggers racing in right after
	// must not see zero readers and release before they've even started.
	cw.readers.Add(len(eventChans))

	c.workflows.set(workflowID, cw)

	deliverFn := c.buildDeliverFn(cw)
	for idx, eventCh := range eventChans {
		triggerCapID := triggerCapIDs[idx]
		c.eng.GoCtx(readerCtx, func(ctx context.Context) {
			defer cw.readers.Done()
			ReadLoop(ctx, lggr, wfMetrics, c.clock, workflowID, triggerCapID, idx, eventCh, deliverFn)
		})
	}

	return triggerCapIDs, nil
}

// buildDeliverFn creates a deliver function that looks up the engine fresh on every event
// instead of caching it, since the engine for a workflow can change while this reader is running.
func (c *coordinator) buildDeliverFn(cw *coordinatedWorkflow) func(context.Context, CoordinatedEvent) {
	return func(ctx context.Context, event CoordinatedEvent) {
		engine, found := c.engines.Get(cw.wid)
		if !found {
			cw.lggr.Infow("Engine gone, dropping trigger event", "triggerID", event.TriggerCapID)
			return
		}
		if !engine.IsCoordinated() {
			cw.lggr.Errorw("Engine is not coordinated, dropping trigger event", "triggerID", event.TriggerCapID)
			return
		}

		// WithoutCancel: unregistering stops ingress, it must not kill an execution already running.
		if err := engine.ExecuteTrigger(context.WithoutCancel(ctx), event); err != nil {
			cw.lggr.Errorw("Failed to execute trigger event", "triggerID", event.TriggerCapID, "err", err)
		}
	}
}

func (c *coordinator) Ack(ctx context.Context, triggerCapID, triggerRegistrationID, eventID string) error {
	workflowID, err := ParseWorkflowID(triggerRegistrationID)
	if err != nil {
		return err
	}

	var handle *Handle
	lggr, wfMetrics := c.lggr, c.deps.Metrics
	if cw, ok := c.workflows.get(workflowID); ok {
		handle = cw.handles[triggerRegistrationID]
		lggr, wfMetrics = cw.lggr, cw.metrics
	}
	return Ack(ctx, lggr, wfMetrics, triggerCapID, triggerRegistrationID, eventID, handle)
}

func (c *coordinator) UnregisterTriggers(workflowID string) error {
	cw, ok := c.workflows.get(workflowID)
	if !ok {
		return ErrWorkflowNotCoordinated
	}
	if cw.unregistered.Load() {
		return nil
	}
	// CAS: exactly one caller spawns the release waiter.
	spawnRelease := cw.releasing.CompareAndSwap(false, true)

	// call cancel to stop all reader goroutines associated with this workflow.
	cw.cancel()
	if spawnRelease {
		c.eng.Go(func(ctx context.Context) { c.releaseWhenDrained(ctx, workflowID, cw) })
	}

	ctx, cancel := c.eng.NewCtx()
	defer cancel()
	ctx = contexts.WithCRE(ctx, cw.cre)

	if failCount := Unregister(ctx, cw.lggr, workflowID, cw.donID, cw.handles); failCount > 0 {
		return fmt.Errorf("failed to unregister %d of %d triggers", failCount, len(cw.handles))
	}

	cw.unregistered.Store(true)

	cw.lggr.Infow("Unregistered triggers, retaining handles until drained", "numTriggers", len(cw.handles))
	cw.metrics.IncrementWorkflowUnregisteredCounter(ctx)
	return nil
}

// releaseWhenDrained drops the handle map and frees the workflow-count limit
// once no execution can still need it to ACK. Waiting for this registration's
// readers to exit is enough: an ACK only ever happens while a reader is still
// running its delivery.
func (c *coordinator) releaseWhenDrained(ctx context.Context, workflowID string, cw *coordinatedWorkflow) {
	ctx = contexts.WithCRE(ctx, cw.cre)

	select {
	case <-cw.drained():
	case <-c.clock.After(c.drainTimeout):
		cw.lggr.Errorw("Timed out waiting for drain, releasing trigger handles anyway")
	case <-ctx.Done():
		return
	}

	c.workflows.deleteIf(workflowID, cw)

	c.freeWorkflowLimit(ctx, cw.lggr)
	cw.lggr.Infow("Released trigger handles")
}

// useWorkflowLimit acquires one slot of the node's workflow-count limit,
// mapping a limit breach to the scope-specific sentinel the syncer expects.
func (c *coordinator) useWorkflowLimit(ctx context.Context, lggr logger.Logger, wfMetrics *monitoring.WorkflowsMetricLabeler) error {
	err := c.workflowLimits.Use(ctx, 1)
	if err == nil {
		return nil
	}

	errLimited, ok := errors.AsType[limits.ErrorResourceLimited[int]](err)
	if !ok {
		return err
	}
	switch errLimited.Scope {
	case settings.ScopeOwner:
		lggr.Infow("Per owner workflow count limit reached", "err", err)
		wfMetrics.IncrementWorkflowLimitPerOwnerCounter(ctx)
		return types.ErrPerOwnerWorkflowCountLimitReached
	case settings.ScopeGlobal:
		lggr.Infow("Global workflow count limit reached", "err", err)
		wfMetrics.IncrementWorkflowLimitGlobalCounter(ctx)
		return types.ErrGlobalWorkflowCountLimitReached
	default:
		lggr.Errorw("Workflow count limit reached for unexpected scope", "scope", errLimited.Scope, "err", err)
		return err
	}
}

// freeWorkflowLimit expects ctx to carry the workflow's CRE: the limit is keyed on it.
func (c *coordinator) freeWorkflowLimit(ctx context.Context, lggr logger.Logger) {
	if err := c.workflowLimits.Free(ctx, 1); err != nil {
		lggr.Errorw("Failed to free workflow count limit", "err", err)
	}
}

// drained returns a channel that will be closed once all readers have finished, indicating that the workflow triggers have been fully drained.
func (cw *coordinatedWorkflow) drained() chan struct{} {
	drained := make(chan struct{})
	// Exits with the readers; outlives this waiter only on a timeout.
	go func() {
		cw.readers.Wait()
		close(drained)
	}()
	return drained
}

// workflowRegistry guards the workflowID -> triggers map so coordinator
// methods never handle the mutex directly.
type workflowRegistry struct {
	mu        sync.Mutex
	workflows map[string]*coordinatedWorkflow // workflowID (hex) -> state
}

func newWorkflowRegistry() *workflowRegistry {
	return &workflowRegistry{workflows: make(map[string]*coordinatedWorkflow)}
}

func (r *workflowRegistry) get(workflowID string) (*coordinatedWorkflow, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cw, ok := r.workflows[workflowID]
	return cw, ok
}

func (r *workflowRegistry) set(workflowID string, cw *coordinatedWorkflow) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.workflows[workflowID] = cw
}

// deleteIf removes workflowID only if it still maps to cw: the workflow may
// have been re-registered while draining, and that state is not ours to drop.
func (r *workflowRegistry) deleteIf(workflowID string, cw *coordinatedWorkflow) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.workflows[workflowID] == cw {
		delete(r.workflows, workflowID)
	}
}

// pending returns the workflows whose capability-side unregistration has not
// completed yet.
func (r *workflowRegistry) pending() map[string]*coordinatedWorkflow {
	r.mu.Lock()
	defer r.mu.Unlock()
	pending := make(map[string]*coordinatedWorkflow, len(r.workflows))
	for workflowID, cw := range r.workflows {
		if !cw.unregistered.Load() {
			pending[workflowID] = cw
		}
	}
	return pending
}
