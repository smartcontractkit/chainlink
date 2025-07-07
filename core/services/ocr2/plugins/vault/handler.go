package vault

import (
	"context"
	"encoding/json"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
)

var _ connector.GatewayConnectorHandler = (*Handler)(nil)

const HandlerName = "VaultHandler"

type gatewaySender interface {
	SendToGateway(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) error
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
		lggr:          lggr.Named(HandlerName),
	}
}

func (h *Handler) Start(ctx context.Context) error {
	return nil
}

func (h *Handler) Close() error {
	return nil
}

func (h *Handler) ID(ctx context.Context) (string, error) {
	return HandlerName, nil
}

func (h *Handler) HandleGatewayMessage(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) error {
	// TODO: do something with the request
	err := h.gatewaySender.SendToGateway(ctx, gatewayID, nil)
	if err != nil {
		h.lggr.Errorf("Failed to send message to gateway %s: %v", gatewayID, err)
		return err
	}
	return nil
}
