package syncer

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type Handler interface {
	Handle(context.Context, WorkflowRegistryEvent) error
}

type eventHandler struct {
	handlers map[WorkflowRegistryEventType]Handler
}

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

func (h *forceUpdateSecretsHandler) Handle(
	ctx context.Context,
	event WorkflowRegistryEvent,
) error {
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

func getSecretsURL(event WorkflowRegistryEvent) (string, error) {
	url, ok := event.Data["SecretsURL"].(string)
	if !ok {
		return "", fmt.Errorf("failed to fetch secrets hash from event : %+v", event.Data)
	}
	return url, nil
}
