package syncer

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type handler interface {
	Handle(ctx context.Context, event WorkflowRegistryEvent) error
}

// eventHandler is a map of event types to their respective handlers that implements Handler.
type eventHandler struct {
	lggr    logger.Logger
	orm     ORM
	fetcher FetcherFunc
}

// newEventHandler returns a new eventHandler with a map of event types to their respective handlers.
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

// Handle processes the ForceUpdateSecretsEvent by fetching the secrets from the URL for a given event
// and updating the local state.
func (h *eventHandler) forceUpdateSecretsEvent(
	ctx context.Context,
	event WorkflowRegistryEvent,
) error {
	// Get the URL of the secrets file from the event data
	url, err := getSecretsURL(event)
	if err != nil {
		h.lggr.Errorf("failed to get URL hash", err)
		return err
	}

	// Fetch the contents of the secrets file from the url via the fetcher
	secrets, err := h.fetcher(ctx, url)
	if err != nil {
		return err
	}

	// Update the secrets in the ORM
	if _, err := h.orm.Update(ctx, url, string(secrets)); err != nil {
		return err
	}

	return nil
}

// getSecretsURL returns the URL of the secrets contents from the event data and fails
// if the URL is not found or is not a string.
func getSecretsURL(event WorkflowRegistryEvent) (string, error) {
	raw, found := event.Data["SecretsURL"]
	if !found {
		return "", fmt.Errorf("failed to fetch secrets hash from event : %+v", event.Data)
	}

	url, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("failed to fetch secrets hash from event : %+v", event.Data)
	}
	return url, nil
}
