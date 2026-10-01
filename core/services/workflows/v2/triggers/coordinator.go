package triggers

import (
	"context"
	"errors"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

// ErrWorkflowNotCoordinated is returned by UnregisterTriggers for
// a workflowID the coordinator never registered.
var ErrWorkflowNotCoordinated = errors.New("workflow not registered with the trigger coordinator")

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
	// can still resolve its handle to ACK. It also ensures that any resources associated with the workflow are properly released.
	// This includes managing the workflowLimits by calling workflowLimits.Free.
	//
	// Returns ErrWorkflowNotCoordinated if workflowID was never registered here
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

// noopCoordinator is a no-op implementation type.  It logs the registration
// and teardown calls on Register/UnregisterTriggers, and does nothing else. No trigger is
// registered with the capability registry and no event is ever delivered.
type noopCoordinator struct {
	services.Service
	eng  *services.Engine
	lggr logger.Logger
}

func NewCoordinator(capReg registry.CapabilitiesRegistry, engineRegistry EngineRegistry, clock clockwork.Clock, lggr logger.Logger) Coordinator {
	c := &noopCoordinator{
		lggr: logger.Named(lggr, "TriggerCoordinator"),
	}

	c.Service, c.eng = services.Config{
		Name: "TriggerCoordinator",
		Start: func(context.Context) error {
			c.lggr.Warnw("No-op trigger coordinator started: workflows on the coordinated engine will register no triggers and receive no events")
			return nil
		},
		Close: func() error { return nil },
	}.NewServiceEngine(c.lggr)

	return c
}

func (c *noopCoordinator) RegisterTriggers(ctx context.Context, subscriber Subscriber, params RegistrationParams) ([]string, error) {
	c.lggr.Infow("No-op RegisterTriggers",
		"workflowID", subscriber.Tenant().Workflow,
		"workflowOwner", params.WorkflowOwner,
		"workflowName", params.WorkflowName.String(),
		"workflowDonID", params.WorkflowDonID)
	return []string{}, nil
}

func (c *noopCoordinator) Ack(ctx context.Context, triggerCapID, triggerRegistrationID, eventID string) error {
	return nil
}

func (c *noopCoordinator) UnregisterTriggers(workflowID string) error {
	c.lggr.Infow("No-op UnregisterTriggers", "workflowID", workflowID)
	return nil
}
