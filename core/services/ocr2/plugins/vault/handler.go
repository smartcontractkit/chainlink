package vault

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	vault_api "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
)

var (
	_ connector.GatewayConnectorHandler = (*Handler)(nil)

	ConnectorMethod = "vault"
	HandlerName     = "VaultHandler"

	promRequestInternalError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "vault_node_request_internal_error",
		Help: "Vault node: Metric to track internal errors",
	}, []string{"error"})

	promRequestSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "vault_node_request_success",
		Help: "Vault node: Metric to track successful requests",
	}, []string{})
)

type gatewaySender interface {
	SendToGateway(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) error
}

type Handler struct {
	vault         *Service
	gatewaySender gatewaySender
	lggr          logger.Logger
}

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
	h.lggr.Debugf("Received message from gateway %s: %v", gatewayID, req)

	var response *jsonrpc.Response[json.RawMessage]
	switch req.Method {
	case vault_api.MethodSecretsCreate:
		response = h.handleSecretsCreate(req)
	default:
		response = h.errorResponse(req, api.UnsupportedMethodError, errors.New("unsupported method: "+req.Method))
	}

	if err = h.gatewaySender.SendToGateway(ctx, gatewayID, response); err != nil {
		h.lggr.Errorf("Failed to send message to gateway %s: %v", gatewayID, err)
		return err
	}

	h.lggr.Infof("Sent message to gateway %s: %v", gatewayID, response)
	promRequestSuccess.WithLabelValues().Inc()
	return nil
}

func (h *Handler) handleSecretsCreate(req *jsonrpc.Request[json.RawMessage]) *jsonrpc.Response[json.RawMessage] {
	var requestData vault_api.SecretsCreateRequest
	if err := json.Unmarshal(*req.Params, &requestData); err != nil {
		return h.errorResponse(req, api.UserMessageParseError, err)
	}

	// DUMMY RESPONSE
	responseData := vault_api.SecretsCreateResponse{
		ResponseBase: vault_api.ResponseBase{
			Success: true,
		},
		SecretID: requestData.ID,
	}

	resultBytes, err := json.Marshal(responseData)
	if err != nil {
		return h.errorResponse(req, api.NodeReponseEncodingError, err)
	}

	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      req.ID,
		Result:  (*json.RawMessage)(&resultBytes),
	}
}

func (h *Handler) errorResponse(
	req *jsonrpc.Request[json.RawMessage],
	errorCode api.ErrorCode,
	err error,
) *jsonrpc.Response[json.RawMessage] {
	h.lggr.Errorf("error code: %d, err: %s", errorCode, err.Error())
	// Given that all requests are coming from the gateway, we can assume that all errors are internal errors
	promRequestInternalError.WithLabelValues(errorCode.String()).Inc()

	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      req.ID,
		Error: &jsonrpc.WireError{
			Code:    api.ToJSONRPCErrorCode(errorCode),
			Message: err.Error(),
		},
	}
}
