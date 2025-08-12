package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	vault_api "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
)

var (
	_ connector.GatewayConnectorHandler = (*GatewayHandler)(nil)

	HandlerName = "VaultHandler"
)

type metrics struct {
	// Given that all requests are coming from the gateway, we can assume that all errors are internal errors
	requestInternalError metric.Int64Counter
	requestSuccess       metric.Int64Counter
}

func newMetrics() (*metrics, error) {
	requestInternalError, err := beholder.GetMeter().Int64Counter("vault_node_request_internal_error")
	if err != nil {
		return nil, fmt.Errorf("failed to register internal error counter: %w", err)
	}

	requestSuccess, err := beholder.GetMeter().Int64Counter("vault_node_request_success")
	if err != nil {
		return nil, fmt.Errorf("failed to register success counter: %w", err)
	}

	return &metrics{
		requestInternalError: requestInternalError,
		requestSuccess:       requestSuccess,
	}, nil
}

type GatewayHandler struct {
	secretsService SecretsService
	gwConnector    core.GatewayConnector
	lggr           logger.Logger
	metrics        *metrics
}

func NewGatewayHandler(secretsService SecretsService, gwConnector core.GatewayConnector, lggr logger.Logger) (*GatewayHandler, error) {
	metrics, err := newMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics: %w", err)
	}

	return &GatewayHandler{
		secretsService: secretsService,
		gwConnector:    gwConnector,
		lggr:           lggr.Named(HandlerName),
		metrics:        metrics,
	}, nil
}

func (h *GatewayHandler) Start(ctx context.Context) error {
	return nil
}

func (h *GatewayHandler) Close() error {
	return nil
}

func (h *GatewayHandler) ID(ctx context.Context) (string, error) {
	return HandlerName, nil
}

func (h *GatewayHandler) HandleGatewayMessage(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) (err error) {
	h.lggr.Debugf("Received message from gateway %s: %v", gatewayID, req)

	var response *jsonrpc.Response[json.RawMessage]
	switch req.Method {
	case vault_api.MethodSecretsCreate:
		response = h.handleSecretsCreate(ctx, gatewayID, req)
	default:
		response = h.errorResponse(ctx, gatewayID, req, api.UnsupportedMethodError, errors.New("unsupported method: "+req.Method))
	}

	if err = h.gwConnector.SendToGateway(ctx, gatewayID, response); err != nil {
		h.lggr.Errorf("Failed to send message to gateway %s: %v", gatewayID, err)
		return err
	}

	h.lggr.Infof("Sent message to gateway %s: %v", gatewayID, response)
	h.metrics.requestSuccess.Add(ctx, 1, metric.WithAttributes(
		attribute.String("gateway_id", gatewayID),
	))
	return nil
}

func (h *GatewayHandler) handleSecretsCreate(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) *jsonrpc.Response[json.RawMessage] {
	var requestData vault_api.SecretsCreateRequest
	if err := json.Unmarshal(*req.Params, &requestData); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.UserMessageParseError, err)
	}
	h.lggr.Infof("Debugging: handleSecretsCreate 1 %s: %v", gatewayID, req)
	vaultCapRequest := vault.CreateSecretsRequest{
		EncryptedSecrets: []*vault.EncryptedSecret{
			{
				Id: &vault.SecretIdentifier{
					Owner:     requestData.Owner,
					Namespace: "", // TBD
					Key:       requestData.ID,
				},
				EncryptedValue: requestData.Value,
			},
		},
	}
	vaultCapResponse, err := h.secretsService.CreateSecrets(ctx, &vaultCapRequest)
	if err != nil {
		h.lggr.Infof("Debugging: h.secretsService.CreateSecrets failed, erro: %s", err.Error())
		return h.errorResponse(ctx, gatewayID, req, api.FatalError, err)
	}
	h.lggr.Infof("Debugging: handleSecretsCreate 2 %s: %v", gatewayID, req)

	resultBytes, err := json.Marshal(vaultCapResponse)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.NodeReponseEncodingError, err)
	}
	h.lggr.Infof("Debugging: handleSecretsCreate 3 %s: %v", gatewayID, req)

	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      req.ID,
		Method:  req.Method,
		Result:  (*json.RawMessage)(&resultBytes),
	}
}

func (h *GatewayHandler) errorResponse(
	ctx context.Context,
	gatewayID string,
	req *jsonrpc.Request[json.RawMessage],
	errorCode api.ErrorCode,
	err error,
) *jsonrpc.Response[json.RawMessage] {
	h.lggr.Infof("GatewayHandler error code: %d, err: %s", errorCode, err.Error())
	h.metrics.requestInternalError.Add(ctx, 1, metric.WithAttributes(
		attribute.String("gateway_id", gatewayID),
		attribute.String("error", errorCode.String()),
	))

	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      req.ID,
		Error: &jsonrpc.WireError{
			Code:    api.ToJSONRPCErrorCode(errorCode),
			Message: err.Error(),
		},
	}
}
