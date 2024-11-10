package syncer

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type Handler interface {
	// Handle processes the event.  Fails if the event type is not supported.
	Handle(context.Context, WorkflowRegistryEvent) error
}

// eventHandler is a map of event types to their respective handlers that implements Handler.
type eventHandler struct {
	handlers map[WorkflowRegistryEventType]Handler
}

// newEventHandler returns a new eventHandler with a map of event types to their respective handlers.
func newEventHandler(
	lggr logger.Logger,
	orm ORM,
	gateway FetcherFunc,
) *eventHandler {
	return &eventHandler{
		handlers: map[WorkflowRegistryEventType]Handler{
			ForceUpdateSecretsEvent: newForceUpdateSecretsHandler(lggr, orm, gateway),
		},
	}
}

func (h *eventHandler) Handle(ctx context.Context, event WorkflowRegistryEvent) error {
	handler, found := h.handlers[event.EventType]
	if !found {
		return fmt.Errorf("event type unsupported : %s", event.EventType)
	}
	return handler.Handle(ctx, event)
}

type forceUpdateSecretsHandler struct {
	lggr    logger.Logger
	orm     ORM
	fetcher FetcherFunc
}

func newForceUpdateSecretsHandler(lggr logger.Logger, orm ORM, gateway FetcherFunc) *forceUpdateSecretsHandler {
	return &forceUpdateSecretsHandler{
		lggr:    lggr,
		orm:     orm,
		fetcher: gateway,
	}
}

// Handle processes the ForceUpdateSecretsEvent by fetching the secrets from the URL for a given event
// and updating the local state.
func (h *forceUpdateSecretsHandler) Handle(
	ctx context.Context,
	event WorkflowRegistryEvent,
) error {
	if event.EventType != ForceUpdateSecretsEvent {
		return fmt.Errorf("event type unsupported : %s", event.EventType)
	}

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
