package vault

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/encoding/protojson"

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
	capRegistry    core.CapabilitiesRegistry
	secretsService SecretsService
	gwConnector    core.GatewayConnector
	lggr           logger.Logger
	metrics        *metrics
}

func NewGatewayHandler(capabilitiesRegistry core.CapabilitiesRegistry, secretsService SecretsService, gwConnector core.GatewayConnector, lggr logger.Logger) (*GatewayHandler, error) {
	metrics, err := newMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics: %w", err)
	}

	return &GatewayHandler{
		capRegistry:    capabilitiesRegistry,
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
	case vault_api.MethodSecretsGet:
		response = h.handleSecretsGet(ctx, gatewayID, req)
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
		RequestId: req.ID,
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
	h.lggr.Infof("Debugging: handleSecretsCreate got response. GatewayId: %s, req: %v, Response: %s", gatewayID, req, vaultCapResponse.String())

	vaultResponseProto := &vault.CreateSecretsResponse{}
	err = protojson.Unmarshal(vaultCapResponse.Payload, vaultResponseProto)
	if err != nil {
		h.lggr.Errorf("Debugging: handleSecretsCreate failed to unmarshal response: %s. Payload was: %s", err.Error(), string(vaultCapResponse.Payload))
		return h.errorResponse(ctx, gatewayID, req, api.NodeReponseEncodingError, err)
	}
	if len(vaultResponseProto.GetResponses()) != 1 {
		return h.errorResponse(ctx, gatewayID, req, api.FatalError, errors.New("unexpected number of responses in CreateSecretsResponse: expected 1, got "+fmt.Sprint(len(vaultResponseProto.GetResponses()))))
	}
	secretResponse := vaultResponseProto.GetResponses()[0]
	vaultApiResponse := vault_api.SecretsCreateResponse{
		ResponseBase: vault_api.ResponseBase{
			ID:         vaultCapResponse.ID,
			Error:      vaultCapResponse.Error,
			Format:     vaultCapResponse.Format,
			Context:    vaultCapResponse.Context,
			Signatures: vaultCapResponse.Signatures,
		},
		SecretID: vault_api.SecretIdentifier{
			Key:       secretResponse.Id.GetKey(),
			Namespace: secretResponse.Id.GetNamespace(),
			Owner:     secretResponse.Id.GetOwner(),
		},
	}

	h.lggr.Infof("Debugging: handleSecretsCreate 3 %s: %v", gatewayID, req)

	vaultApiResponseBytes, err := json.Marshal(vaultApiResponse)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.NodeReponseEncodingError, err)
	}
	vaultApiResponseJson := json.RawMessage(vaultApiResponseBytes)
	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      vaultApiResponse.ID,
		Method:  req.Method,
		Result:  &vaultApiResponseJson,
	}
}

func (h *GatewayHandler) handleSecretsGet(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) *jsonrpc.Response[json.RawMessage] {
	var requestData vault_api.SecretsGetRequest
	if err := json.Unmarshal(*req.Params, &requestData); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.UserMessageParseError, err)
	}
	h.lggr.Infof("Debugging: handleSecretsGet 1 %s: %v", gatewayID, req)
	encryptionKeys, err := h.getEncryptionKeys(ctx)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.FatalError, err)
	}
	getSecretsRequest := vault.GetSecretsRequest{
		Requests: []*vault.SecretRequest{
			{
				Id: &vault.SecretIdentifier{
					Owner:     requestData.Owner,
					Namespace: "", // TBD
					Key:       requestData.ID,
				},
				EncryptionKeys: encryptionKeys,
			},
		},
	}
	vaultCapResponse, err := h.secretsService.GetSecrets(ctx, req.ID, &getSecretsRequest)
	if err != nil {
		h.lggr.Infof("Debugging: h.secretsService.GetSecrets failed, erro: %s", err.Error())
		return h.errorResponse(ctx, gatewayID, req, api.FatalError, err)
	}
	h.lggr.Infof("Debugging: handleSecretsGet 2 %s: %v", gatewayID, req)

	resultBytes, err := json.Marshal(vaultCapResponse)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.NodeReponseEncodingError, err)
	}
	h.lggr.Infof("Debugging: handleSecretsGet 3 %s: %v", gatewayID, req)

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

// getEncryptionKeys retrieves the encryption keys of all members in the Workflow DON.
func (h *GatewayHandler) getEncryptionKeys(ctx context.Context) ([]string, error) {
	myNode, err := h.capRegistry.LocalNode(ctx)
	if err != nil {
		return nil, errors.New("failed to get local node from registry" + err.Error())
	}

	encryptionKeys := make([]string, 0, len(myNode.WorkflowDON.Members))
	for _, peerID := range myNode.WorkflowDON.Members {
		peerNode, err := h.capRegistry.NodeByPeerID(ctx, peerID)
		if err != nil {
			return nil, errors.New("failed to get node info for peerID: " + peerID.String() + " - " + err.Error())
		}
		encryptionKeys = append(encryptionKeys, hex.EncodeToString(peerNode.EncryptionPublicKey[:]))
	}
	// Sort the encryption keys to ensure consistent ordering across all nodes.
	sort.Strings(encryptionKeys)
	return encryptionKeys, nil
}
