package syncer

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
)

var ErrNotImplemented = errors.New("not implemented")

// WorkflowRegistryrEventType is the type of event that is emitted by the WorkflowRegistry
type WorkflowRegistryEventType string

var (
	// ForceUpdateSecretsEvent is emitted when a request to force update a workflows secrets is made
	ForceUpdateSecretsEvent WorkflowRegistryEventType = "WorkflowForceUpdateSecretsRequestedV1"

	// WorkflowRegisteredEvent is emitted when a workflow is registered
	WorkflowRegisteredEvent WorkflowRegistryEventType = "WorkflowRegisteredV1"

	// WorkflowUpdatedEvent is emitted when a workflow is updated
	WorkflowUpdatedEvent WorkflowRegistryEventType = "WorkflowUpdatedV1"

	// WorkflowPausedEvent is emitted when a workflow is paused
	WorkflowPausedEvent WorkflowRegistryEventType = "WorkflowPausedV1"

	// WorkflowActivatedEvent is emitted when a workflow is activated
	WorkflowActivatedEvent WorkflowRegistryEventType = "WorkflowActivatedV1"

	// WorkflowDeletedEvent is emitted when a workflow is deleted
	WorkflowDeletedEvent WorkflowRegistryEventType = "WorkflowDeletedV1"
)

// WorkflowRegistryForceUpdateSecretsRequestedV1 is a chain agnostic definition of the WorkflowRegistry
// ForceUpdateSecretsRequested event.
type WorkflowRegistryForceUpdateSecretsRequestedV1 struct {
	SecretsURLHash []byte
	Owner          []byte
	WorkflowName   string
}

type WorkflowRegistryWorkflowRegisteredV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner []byte
	DonID         uint32
	Status        uint8
	WorkflowName  string
	BinaryURL     string
	ConfigURL     string
	SecretsURL    string
}

type WorkflowRegistryWorkflowUpdatedV1 struct {
	OldWorkflowID [32]byte
	WorkflowOwner []byte
	DonID         uint32
	NewWorkflowID [32]byte
	WorkflowName  string
	BinaryURL     string
	ConfigURL     string
	SecretsURL    string
}

type WorkflowRegistryWorkflowPausedV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner []byte
	DonID         uint32
	WorkflowName  string
}

type WorkflowRegistryWorkflowActivatedV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner []byte
	DonID         uint32
	WorkflowName  string
}

type WorkflowRegistryWorkflowDeletedV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner []byte
	DonID         uint32
	WorkflowName  string
}

// eventHandler is a handler for WorkflowRegistryEvent events.  Each event type has a corresponding
// method that handles the event.
type eventHandler struct {
	lggr           logger.Logger
	orm            WorkflowSecretsDS
	fetcher        FetcherFunc
	workflowStore  store.Store
	capRegistry    core.CapabilitiesRegistry
	engineRegistry *engineRegistry
}

// newEventHandler returns a new eventHandler instance.
func newEventHandler(
	lggr logger.Logger,
	orm ORM,
	gateway FetcherFunc,
	workflowStore store.Store,
	capRegistry core.CapabilitiesRegistry,
	engineRegistry *engineRegistry,
) *eventHandler {
	return &eventHandler{
		lggr:           lggr,
		orm:            orm,
		fetcher:        gateway,
		workflowStore:  workflowStore,
		capRegistry:    capRegistry,
		engineRegistry: engineRegistry,
	}
}

func (h *eventHandler) Handle(ctx context.Context, event WorkflowRegistryEvent) error {
	switch event.EventType {
	case ForceUpdateSecretsEvent:
		return h.forceUpdateSecretsEvent(ctx, event)
	case WorkflowRegisteredEvent:
		return h.workflowRegisteredEvent(ctx, event)
	case WorkflowUpdatedEvent:
		return h.workflowUpdatedEvent(ctx, event)
	case WorkflowPausedEvent:
		return h.workflowPausedEvent(ctx, event)
	case WorkflowActivatedEvent:
		return h.workflowActivatedEvent(ctx, event)
	default:
		return fmt.Errorf("event type unsupported: %v", event.EventType)
	}
}

// workflowRegisteredEvent handles the WorkflowRegisteredEvent event type.
// TODO: Implement this method
func (h *eventHandler) workflowRegisteredEvent(
	_ context.Context,
	_ WorkflowRegistryEvent,
) error {
	return ErrNotImplemented
}

// workflowUpdatedEvent handles the WorkflowUpdatedEvent event type.
func (h *eventHandler) workflowUpdatedEvent(
	_ context.Context,
	_ WorkflowRegistryEvent,
) error {
	return ErrNotImplemented
}

// workflowPausedEvent handles the WorkflowPausedEvent event type.
func (h *eventHandler) workflowPausedEvent(
	_ context.Context,
	_ WorkflowRegistryEvent,
) error {
	return ErrNotImplemented
}

// workflowActivatedEvent handles the WorkflowActivatedEvent event type.
func (h *eventHandler) workflowActivatedEvent(
	_ context.Context,
	_ WorkflowRegistryEvent,
) error {
	return ErrNotImplemented
}

// forceUpdateSecretsEvent handles the ForceUpdateSecretsEvent event type.
func (h *eventHandler) forceUpdateSecretsEvent(
	ctx context.Context,
	event WorkflowRegistryEvent,
) error {
	// Get the URL of the secrets file from the event data
	data, ok := event.Data.(WorkflowRegistryForceUpdateSecretsRequestedV1)
	if !ok {
		return fmt.Errorf("invalid data type %T for event", event.Data)
	}

	hash := hex.EncodeToString(data.SecretsURLHash)

	url, err := h.orm.GetSecretsURLByHash(ctx, hash)
	if err != nil {
		h.lggr.Errorf("failed to get URL by hash %s : %s", hash, err)
		return err
	}

	// Fetch the contents of the secrets file from the url via the fetcher
	secrets, err := h.fetcher(ctx, url)
	if err != nil {
		return err
	}

	// Update the secrets in the ORM
	if _, err := h.orm.Update(ctx, hash, string(secrets)); err != nil {
		return err
	}

	return nil
}
