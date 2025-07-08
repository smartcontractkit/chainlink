package vault

import (
	"context"
	"encoding/json"
	"errors"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	vault_api "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
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

func (h *Handler) HandleGatewayMessage(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) (err error) {
	// TODO: authorize the request
	// TODO: check method and handle accordingly
	h.lggr.Infof("Received message from gateway %s: %v", gatewayID, req)
	// TODO: do something with the request and send a proper response
	// TODO: Add prom counters

	var response json.RawMessage
	switch req.Method {
	case vault_api.MethodSecretsCreate:
		response, err = h.handleSecretsCreate(ctx, gatewayID, req)

	default:
		// This should never happen because the gateway should not send a message with an unsupported method
		err = errors.New("unsupported method: " + req.Method)
	}

	if err != nil {
		// TODO: Return an error responseto the Gateway
		return err
	}

	jsonResponse := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      req.ID,
		Result:  &response,
	}

	if err = h.gatewaySender.SendToGateway(ctx, gatewayID, &jsonResponse); err != nil {
		h.lggr.Errorf("Failed to send message to gateway %s: %v", gatewayID, err)
		return err
	}

	return nil
}

func (h *Handler) handleSecretsCreate(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) (json.RawMessage, error) {
	var requestData vault_api.SecretsCreateRequest
	err := json.Unmarshal(*req.Params, &requestData)
	if err != nil {
		h.lggr.Errorf("Failed to unmarshal request: %v", err)
		return nil, err
	}

	// DUMMY RESPONSE
	responseData := vault_api.SecretsCreateResponse{
		ResponseBase: vault_api.ResponseBase{
			Success: true,
		},
		SecretID: requestData.ID,
	}

	var resultBytes json.RawMessage
	resultBytes, err = json.Marshal(responseData)
	if err != nil {
		h.lggr.Errorf("Failed to marshal response: %v", err)
		return nil, err
	}
	return resultBytes, nil
}
