package syncer

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// eventHandler is a handler for WorkflowRegistryEvent events.  Each event type has a corresponding
// method that handles the event.
type eventHandler struct {
	lggr    logger.Logger
	orm     ORM
	fetcher FetcherFunc
}

// newEventHandler returns a new eventHandler instance.
func newEventHandler(
	lggr logger.Logger,
	orm ORM,
	gateway FetcherFunc,
) *eventHandler {
	return &eventHandler{
		lggr:    lggr,
		orm:     orm,
		fetcher: gateway,
	}
}

func (h *eventHandler) Handle(ctx context.Context, event WorkflowRegistryEvent) error {
	switch event.EventType {
	case ForceUpdateSecretsEvent:
		return h.forceUpdateSecretsEvent(ctx, event)
	default:
		return fmt.Errorf("event type unsupported: %v", event.EventType)
	}
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
