package vault

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
)

var _ connector.GatewayConnectorHandler = (*Handler)(nil)

type gatewaySender interface {
	SendToGateway(ctx context.Context, gatewayID string, msg *api.Message) error
}

type Handler struct {
	vault         *Service
	gatewaySender gatewaySender
	lggr          logger.Logger
}

var ConnectorMethod = "vault"

func NewHandler(vault *Service, gwsender gatewaySender, lggr logger.Logger) *Handler {
	return &Handler{
		vault:         vault,
		gatewaySender: gwsender,
		lggr:          lggr.Named("VaultHandler"),
	}
}

func (h *Handler) Start(ctx context.Context) error {
	return nil
}

func (h *Handler) Close() error {
	return nil
}

func (h *Handler) HandleGatewayMessage(ctx context.Context, gatewayID string, msg *api.Message) {
	// TODO: do something with the request
	err := h.gatewaySender.SendToGateway(ctx, gatewayID, msg)
	if err != nil {
		h.lggr.Errorf("Failed to send message to gateway %s: %v", gatewayID, err)
	}
}
