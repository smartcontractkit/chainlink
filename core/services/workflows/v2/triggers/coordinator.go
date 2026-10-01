package triggers

import (
	"context"
	"errors"
	"fmt"
	"sync"
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

	mu        sync.Mutex
	workflows map[string]*workflowTriggers // workflowID (hex) -> state
}

type workflowTriggers struct {
	wid   types.WorkflowID
	cre   contexts.CRE
	donID uint32
	// cancel stops only this registration's readers. A workflow can be
	// re-registered while its old registration is still draining, so this
	// must not affect the new one.
	cancel  context.CancelFunc
	handles map[string]*Handle

	// readers tracks the reader goroutines for this registration. A reader
	// only exits once its own delivery has returned, so once Wait() returns,
	// nothing started by this registration can still ACK.
	readers sync.WaitGroup

	// guarded by coordinator.mu
	unregistered bool // capability-side unregistration succeeded
	releasing    bool // release waiter spawned

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
		workflows:      make(map[string]*workflowTriggers),
	}

	c.Service, c.eng = services.Config{
		Name:  "TriggerCoordinator",
		Close: c.close,
	}.NewServiceEngine(c.lggr)

	return c
}

// close runs after the reader and release goroutines have exited, and
// unregisters whatever the syncer did not tear down before shutdown.
func (c *coordinator) close() error {
	c.mu.Lock()
	pending := make(map[string]*workflowTriggers, len(c.workflows))
	for workflowID, wt := range c.workflows {
		if !wt.unregistered {
			pending[workflowID] = wt
		}
	}
	c.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	var errs error
	for workflowID, wt := range pending {
		wt.cancel()
		if failCount := Unregister(contexts.WithCRE(ctx, wt.cre), wt.lggr, workflowID, wt.donID, wt.handles); failCount > 0 {
			errs = errors.Join(errs, fmt.Errorf("workflow %s: failed to unregister %d of %d triggers", workflowID, failCount, len(wt.handles)))
		}
	}
	return errs
}

func (c *coordinator) RegisterTriggers(ctx context.Context, subscriber Subscriber, params RegistrationParams) ([]string, error) {
	cre := subscriber.Tenant()
	workflowID := cre.Workflow
	wid, err := types.WorkflowIDFromHex(workflowID)
	if err != nil {
		return nil, fmt.Errorf("invalid workflow id %q: %w", workflowID, err)
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
	if err := c.UnregisterTriggers(workflowID); err != nil && !errors.Is(err, ErrWorkflowNotCoordinated) {
		lggr.Errorw("Failed to unregister previous trigger registration", "err", err)
	}

	if err := c.useWorkflowLimit(ctx, lggr, wfMetrics); err != nil {
		return nil, err
	}
	freeLimit := true
	defer func() {
		if freeLimit {
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
	wt := &workflowTriggers{
		wid:     wid,
		cre:     cre,
		donID:   params.WorkflowDonID,
		cancel:  cancel,
		handles: handles,
		lggr:    lggr,
		metrics: wfMetrics,
	}

	// Set before publishing wt: an UnregisterTriggers racing in right after
	// must not see zero readers and release before they've even started.
	wt.readers.Add(len(eventChans))

	c.mu.Lock()
	c.workflows[workflowID] = wt
	c.mu.Unlock()

	deliver := c.deliver(wt)
	for idx, eventCh := range eventChans {
		triggerCapID := triggerCapIDs[idx]
		c.eng.GoCtx(readerCtx, func(ctx context.Context) {
			defer wt.readers.Done()
			ReadLoop(ctx, lggr, wfMetrics, c.clock, workflowID, triggerCapID, idx, eventCh, deliver)
		})
	}

	freeLimit = false
	return triggerCapIDs, nil
}

// deliver looks up the engine fresh on every event instead of caching it,
// since the engine for a workflow can change while this reader is running.
func (c *coordinator) deliver(wt *workflowTriggers) func(context.Context, CoordinatedEvent) {
	return func(ctx context.Context, event CoordinatedEvent) {
		engine, found := c.engines.Get(wt.wid)
		if !found {
			wt.lggr.Infow("Engine gone, dropping trigger event", "triggerID", event.TriggerCapID)
			return
		}
		if !engine.IsCoordinated() {
			wt.lggr.Errorw("Engine is not coordinated, dropping trigger event", "triggerID", event.TriggerCapID)
			return
		}

		// WithoutCancel: unregistering stops ingress, it must not kill an execution already running.
		if err := engine.ExecuteTrigger(context.WithoutCancel(ctx), event); err != nil {
			wt.lggr.Errorw("Failed to execute trigger event", "triggerID", event.TriggerCapID, "err", err)
		}
	}
}

func (c *coordinator) Ack(ctx context.Context, triggerCapID, triggerRegistrationID, eventID string) error {
	workflowID, err := ParseWorkflowID(triggerRegistrationID)
	if err != nil {
		return err
	}

	c.mu.Lock()
	var handle *Handle
	wt, found := c.workflows[workflowID]
	if found {
		handle = wt.handles[triggerRegistrationID]
	}
	c.mu.Unlock()

	lggr, wfMetrics := c.lggr, c.deps.Metrics
	if found {
		lggr, wfMetrics = wt.lggr, wt.metrics
	}
	return Ack(ctx, lggr, wfMetrics, triggerCapID, triggerRegistrationID, eventID, handle)
}

func (c *coordinator) UnregisterTriggers(workflowID string) error {
	c.mu.Lock()
	wt, ok := c.workflows[workflowID]
	if !ok {
		c.mu.Unlock()
		return ErrWorkflowNotCoordinated
	}
	if wt.unregistered {
		c.mu.Unlock()
		return nil
	}
	spawnRelease := !wt.releasing
	wt.releasing = true
	c.mu.Unlock()

	// Cancel first: not every capability closes its event channel on
	// Unregister, so this is the only guaranteed way to stop delivery.
	wt.cancel()
	if spawnRelease {
		c.eng.Go(func(ctx context.Context) { c.releaseWhenDrained(ctx, workflowID, wt) })
	}

	ctx, cancel := c.eng.NewCtx()
	defer cancel()
	ctx = contexts.WithCRE(ctx, wt.cre)

	if failCount := Unregister(ctx, wt.lggr, workflowID, wt.donID, wt.handles); failCount > 0 {
		return fmt.Errorf("failed to unregister %d of %d triggers", failCount, len(wt.handles))
	}

	c.mu.Lock()
	wt.unregistered = true
	c.mu.Unlock()

	wt.lggr.Infow("Unregistered triggers, retaining handles until drained", "numTriggers", len(wt.handles))
	wt.metrics.IncrementWorkflowUnregisteredCounter(ctx)
	return nil
}

// releaseWhenDrained drops the handle map and frees the workflow-count limit
// once no execution can still need it to ACK. Waiting for this registration's
// readers to exit is enough: an ACK only ever happens while a reader is still
// running its delivery.
func (c *coordinator) releaseWhenDrained(ctx context.Context, workflowID string, wt *workflowTriggers) {
	ctx = contexts.WithCRE(ctx, wt.cre)

	drained := make(chan struct{})
	// Exits with the readers; outlives this waiter only on a timeout.
	go func() {
		wt.readers.Wait()
		close(drained)
	}()

	select {
	case <-drained:
	case <-c.clock.After(c.drainTimeout):
		wt.lggr.Errorw("Timed out waiting for drain, releasing trigger handles anyway")
	case <-ctx.Done():
		return
	}

	c.mu.Lock()
	// The workflow may have been re-registered while draining; its state is not ours to drop.
	if c.workflows[workflowID] == wt {
		delete(c.workflows, workflowID)
	}
	c.mu.Unlock()

	c.freeWorkflowLimit(ctx, wt.lggr)
	wt.lggr.Infow("Released trigger handles")
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
