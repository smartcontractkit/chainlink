package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
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

type gatewayConnector interface {
	SendToGateway(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) error
	AddHandler(ctx context.Context, methods []string, handler core.GatewayConnectorHandler) error
	RemoveHandler(ctx context.Context, methods []string) error
}

type GatewayHandler struct {
	services.Service
	eng *services.Engine

	secretsService   vaulttypes.SecretsService
	gatewayConnector gatewayConnector
	lggr             logger.Logger
	metrics          *metrics
}

func NewGatewayHandler(secretsService vaulttypes.SecretsService, connector gatewayConnector, lggr logger.Logger) (*GatewayHandler, error) {
	metrics, err := newMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics: %w", err)
	}

	gh := &GatewayHandler{
		secretsService:   secretsService,
		gatewayConnector: connector,
		lggr:             lggr.Named(HandlerName),
		metrics:          metrics,
	}
	gh.Service, gh.eng = services.Config{
		Name:  "GatewayHandler",
		Start: gh.start,
		Close: gh.close,
	}.NewServiceEngine(lggr)
	return gh, nil
}

func (h *GatewayHandler) start(ctx context.Context) error {
	if gwerr := h.gatewayConnector.AddHandler(ctx, h.Methods(), h); gwerr != nil {
		return fmt.Errorf("failed to add vault handler to connector: %w", gwerr)
	}
	return nil
}

func (h *GatewayHandler) close() error {
	if gwerr := h.gatewayConnector.RemoveHandler(context.Background(), h.Methods()); gwerr != nil {
		return fmt.Errorf("failed to remove vault handler from connector: %w", gwerr)
	}
	return nil
}

func (h *GatewayHandler) ID(ctx context.Context) (string, error) {
	return HandlerName, nil
}

func (h *GatewayHandler) Methods() []string {
	return vaulttypes.Methods
}

func (h *GatewayHandler) HandleGatewayMessage(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) (err error) {
	h.lggr.Debugw("received message from gateway", "gatewayID", gatewayID, "req", req, "requestID", req.ID)

	var response *jsonrpc.Response[json.RawMessage]
	switch req.Method {
	case vaulttypes.MethodSecretsCreate:
		response = h.handleSecretsCreate(ctx, gatewayID, req)
	case vaulttypes.MethodSecretsUpdate:
		response = h.handleSecretsUpdate(ctx, gatewayID, req)
	case vaulttypes.MethodSecretsDelete:
		response = h.handleSecretsDelete(ctx, gatewayID, req)
	case vaulttypes.MethodSecretsList:
		response = h.handleSecretsList(ctx, gatewayID, req)
	case vaulttypes.MethodPublicKeyGet:
		response = h.handlePublicKeyGet(ctx, gatewayID, req)
	default:
		response = h.errorResponse(ctx, gatewayID, req, api.UnsupportedMethodError, errors.New("unsupported method: "+req.Method))
	}

	if err = h.gatewayConnector.SendToGateway(ctx, gatewayID, response); err != nil {
		h.lggr.Errorf("Failed to send message to gateway %s: %v", gatewayID, err)
		return err
	}

	h.lggr.Infow("Sent message to gateway", "gatewayID", gatewayID, "resp", response, "requestID", req.ID)
	h.metrics.requestSuccess.Add(ctx, 1, metric.WithAttributes(
		attribute.String("gateway_id", gatewayID),
	))
	return nil
}

func (h *GatewayHandler) handleSecretsCreate(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) *jsonrpc.Response[json.RawMessage] {
	vaultCapRequest := vaultcommon.CreateSecretsRequest{}
	if err := json.Unmarshal(*req.Params, &vaultCapRequest); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.UserMessageParseError, err)
	}

	vaultCapRequest.RequestId = req.ID

	vaultCapResponse, err := h.secretsService.CreateSecrets(ctx, &vaultCapRequest)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.FatalError, err)
	}

	jsonResponse, err := toJSONResponse(vaultCapResponse, req.Method)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.NodeReponseEncodingError, err)
	}
	return jsonResponse
}

func (h *GatewayHandler) handleSecretsUpdate(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) *jsonrpc.Response[json.RawMessage] {
	vaultCapRequest := vaultcommon.UpdateSecretsRequest{}
	if err := json.Unmarshal(*req.Params, &vaultCapRequest); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.UserMessageParseError, err)
	}

	vaultCapResponse, err := h.secretsService.UpdateSecrets(ctx, &vaultCapRequest)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.FatalError, err)
	}

	jsonResponse, err := toJSONResponse(vaultCapResponse, req.Method)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.NodeReponseEncodingError, err)
	}
	return jsonResponse
}

func (h *GatewayHandler) handleSecretsDelete(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) *jsonrpc.Response[json.RawMessage] {
	r := &vaultcommon.DeleteSecretsRequest{}
	if err := json.Unmarshal(*req.Params, r); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.UserMessageParseError, err)
	}

	resp, err := h.secretsService.DeleteSecrets(ctx, r)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.HandlerError, fmt.Errorf("failed to delete secrets: %w", err))
	}

	resultBytes, err := resp.ToJSONRPCResult()
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.NodeReponseEncodingError, err)
	}

	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      req.ID,
		Method:  req.Method,
		Result:  (*json.RawMessage)(&resultBytes),
	}
}

func (h *GatewayHandler) handleSecretsList(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) *jsonrpc.Response[json.RawMessage] {
	r := &vaultcommon.ListSecretIdentifiersRequest{}
	if err := json.Unmarshal(*req.Params, r); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.UserMessageParseError, err)
	}

	resp, err := h.secretsService.ListSecretIdentifiers(ctx, r)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.HandlerError, fmt.Errorf("failed to list secret identifiers: %w", err))
	}

	resultBytes, err := resp.ToJSONRPCResult()
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.NodeReponseEncodingError, err)
	}

	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      req.ID,
		Method:  req.Method,
		Result:  (*json.RawMessage)(&resultBytes),
	}
}

func (h *GatewayHandler) handlePublicKeyGet(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) *jsonrpc.Response[json.RawMessage] {
	r := &vaultcommon.GetPublicKeyRequest{}
	if err := json.Unmarshal(*req.Params, r); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.UserMessageParseError, err)
	}

	resp, err := h.secretsService.GetPublicKey(ctx, r)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.HandlerError, fmt.Errorf("failed to get public key: %w", err))
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.NodeReponseEncodingError, err)
	}

	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      req.ID,
		Method:  req.Method,
		Result:  (*json.RawMessage)(&b),
	}
}

func (h *GatewayHandler) errorResponse(
	ctx context.Context,
	gatewayID string,
	req *jsonrpc.Request[json.RawMessage],
	errorCode api.ErrorCode,
	err error,
) *jsonrpc.Response[json.RawMessage] {
	h.lggr.Errorf("error code: %d, err: %s", errorCode, err.Error())
	h.metrics.requestInternalError.Add(ctx, 1, metric.WithAttributes(
		attribute.String("gateway_id", gatewayID),
		attribute.String("error", errorCode.String()),
	))

	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      req.ID,
		Method:  req.Method,
		Error: &jsonrpc.WireError{
			Code:    api.ToJSONRPCErrorCode(errorCode),
			Message: err.Error(),
		},
	}
}

func toJSONResponse(vaultCapResponse *vaulttypes.Response, method string) (*jsonrpc.Response[json.RawMessage], error) {
	vaultResponseBytes, err := vaultCapResponse.ToJSONRPCResult()
	if err != nil {
		return nil, errors.New("failed to marshal vault capability response: " + err.Error())
	}
	var vaultResponseJSON json.RawMessage = vaultResponseBytes
	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      vaultCapResponse.ID,
		Method:  method,
		Result:  &vaultResponseJSON,
	}, nil
}
