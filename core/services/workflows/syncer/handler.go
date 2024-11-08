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
	h.lggr.Debugf("got event %+v", event)
	h.lggr.Debugf("got URL %x", event.Output.Data["SecretsURL"])
	hash, err := getURLHash(event)
	if err != nil {
		h.lggr.Errorf("failed to get URL hash", err)
		return err
	}

	// Fetch the secrets url from ORM
	url, err := h.orm.GetSecretsURL(ctx, hash)
	if err != nil {
		h.lggr.Errorf("failed to get secrets URL for hash %x", hash)
		return err
	}

	// Fetch the contents of the secrets file from the url via the fetcher
	secrets, err := h.fetcher(ctx, url)
	h.lggr.Debugf("fetched these contents %s", secrets)
	if err != nil {
		return err
	}

	// Update the secrets in the ORM
	h.lggr.Debugf("calling update with url %s and secrets %s", url, secrets)
	if _, err := h.orm.Update(ctx, url, string(secrets)); err != nil {
		return err
	}

	return nil
}

func getURLHash(event WorkflowRegistryEvent) (string, error) {
	hash, ok := event.Data["SecretsURL"].([]byte)
	if !ok {
		return "", fmt.Errorf("failed to fetch secrets hash from event : %+v", event.Data)
	}
	return string(hash), nil
}
