package v2

import (
	"context"
	"errors"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

// ErrWorkflowNotCoordinated is returned by UnregisterTriggers for
// a workflowID the coordinator never registered.
var ErrWorkflowNotCoordinated = errors.New("workflow not registered with the trigger coordinator")

// RegistrationParams is the per-workflow metadata RegisterTriggers needs to
// build each capability's RequestMetadata.
type RegistrationParams struct {
	WorkflowOwner                 string
	WorkflowName                  string
	DecodedWorkflowName           string
	WorkflowTag                   string
	WorkflowDonID                 uint32
	WorkflowDonConfigVersion      uint32
	WorkflowRegistryChainSelector string
	WorkflowRegistryAddress       string
}

// TriggerCoordinator owns trigger registration, the trigger handle map, event
// delivery, and acknowledgement for every workflow running the coordinated engine.
// It is a node-level singleton, started and stopped with the syncer.
type TriggerCoordinator interface {
	services.Service
	v2.Acknowledger

	// RegisterTriggers calls subscriber.Subscribe to obtain the engine's trigger
	// subscriptions, registers them with the capability registry, retains the
	// resulting handles and reader goroutines, and returns the registered
	// trigger capability IDs. On any failure it unregisters what it already
	// registered for this call and returns the error.
	// Partial registration is never left behind.
	RegisterTriggers(ctx context.Context, subscriber v2.Subscriber, params RegistrationParams) ([]string, error)

	// UnregisterTriggers stops ingress for workflowID immediately (unregisters with the capability registry)
	// and cleans up the handle map once the engine has been drained and closed, so an execution already in flight
	// can still resolve its handle to ACK.
	//
	// Returns ErrWorkflowNotCoordinated if workflowID was never registered here
	UnregisterTriggers(workflowID string) error
}

type noopTriggerCoordinator struct {
	services.Service
	eng  *services.Engine
	lggr logger.Logger
}

func NewTriggerCoordinator(capReg registry.CapabilitiesRegistry, engineRegistry *EngineRegistry, clock clockwork.Clock, lggr logger.Logger) TriggerCoordinator {
	c := &noopTriggerCoordinator{
		lggr: logger.Named(lggr, "TriggerCoordinator"),
	}

	c.Service, c.eng = services.Config{
		Name:  "TriggerCoordinator",
		Start: func(context.Context) error { return nil },
		Close: func() error { return nil },
	}.NewServiceEngine(c.lggr)

	return c
}

func (c *noopTriggerCoordinator) RegisterTriggers(ctx context.Context, subscriber v2.Subscriber, params RegistrationParams) ([]string, error) {
	return []string{}, nil
}

func (c *noopTriggerCoordinator) Ack(ctx context.Context, triggerCapID, triggerRegistrationID, eventID string) error {
	return nil
}

func (c *noopTriggerCoordinator) UnregisterTriggers(workflowID string) error {
	return nil
}
