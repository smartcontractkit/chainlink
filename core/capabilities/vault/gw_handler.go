package vault

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	vault_api "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
	vault2 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/vault"
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

type gatewaySender interface {
	SendToGateway(ctx context.Context, gatewayID string, resp *jsonrpc.Response[json.RawMessage]) error
}

type GatewayHandler struct {
	capRegistry    core.CapabilitiesRegistry
	secretsService SecretsService
	gatewaySender  gatewaySender
	lggr           logger.Logger
	metrics        *metrics
}

func NewGatewayHandler(capabilitiesRegistry core.CapabilitiesRegistry, secretsService SecretsService, gwsender gatewaySender, lggr logger.Logger) (*GatewayHandler, error) {
	metrics, err := newMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics: %w", err)
	}

	return &GatewayHandler{
		capRegistry:    capabilitiesRegistry,
		secretsService: secretsService,
		gatewaySender:  gwsender,
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
	h.lggr.Debugw("received message from gateway", "gatewayID", gatewayID, "req", req, "requestID", req.ID)

	var response *jsonrpc.Response[json.RawMessage]
	switch req.Method {
	case vault_api.MethodSecretsCreate:
		response = h.handleSecretsCreate(ctx, gatewayID, req)
	case vault_api.MethodSecretsGet:
		response = h.handleSecretsGet(ctx, gatewayID, req)
	case vault_api.MethodSecretsUpdate:
		response = h.handleSecretsUpdate(ctx, gatewayID, req)
	default:
		response = h.errorResponse(ctx, gatewayID, req, api.UnsupportedMethodError, errors.New("unsupported method: "+req.Method))
	}

	if err = h.gatewaySender.SendToGateway(ctx, gatewayID, response); err != nil {
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
		return h.errorResponse(ctx, gatewayID, req, api.FatalError, err)
	}

	jsonResponse, err := toJSONResponse(vaultCapResponse, req.Method)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.NodeReponseEncodingError, err)
	}
	return jsonResponse
}

func (h *GatewayHandler) handleSecretsUpdate(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) *jsonrpc.Response[json.RawMessage] {
	r := &vault.UpdateSecretsRequest{}
	if err := json.Unmarshal(*req.Params, r); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.UserMessageParseError, err)
	}

	vaultCapResponse, err := h.secretsService.UpdateSecrets(ctx, r)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.FatalError, err)
	}

	jsonResponse, err := toJSONResponse(vaultCapResponse, req.Method)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.NodeReponseEncodingError, err)
	}
	return jsonResponse
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
		return h.errorResponse(ctx, gatewayID, req, api.FatalError, err)
	}
	vaultResponseProto := &vault.GetSecretsResponse{}
	err = proto.Unmarshal(vaultCapResponse.Payload, vaultResponseProto)
	if err != nil {
		h.lggr.Errorf("Debugging: handleSecretsCreate failed to unmarshal response: %s. Payload was: %s", err.Error(), string(vaultCapResponse.Payload))
		return h.errorResponse(ctx, gatewayID, req, api.NodeReponseEncodingError, err)
	}
	if len(vaultResponseProto.GetResponses()) != 1 {
		return h.errorResponse(ctx, gatewayID, req, api.FatalError, errors.New("unexpected number of responses in CreateSecretsResponse: expected 1, got "+strconv.Itoa(len(vaultResponseProto.GetResponses()))))
	}
	secretResponse := vaultResponseProto.GetResponses()[0]
	vaultAPIResponse := vault_api.SecretsGetResponse{
		SecretID: vault_api.SecretIdentifier{
			Key:       secretResponse.Id.GetKey(),
			Namespace: secretResponse.Id.GetNamespace(),
			Owner:     secretResponse.Id.GetOwner(),
		},
	}

	switch method := secretResponse.Result.(type) {
	case *vault.SecretResponse_Data:
		vaultAPIResponse.SecretValue = vault_api.SecretData{
			EncryptedValue:               method.Data.GetEncryptedValue(),
			EncryptedDecryptionKeyShares: make([]*vault_api.EncryptedShares, 0, len(method.Data.GetEncryptedDecryptionKeyShares())),
		}
		for _, decryptionShare := range method.Data.GetEncryptedDecryptionKeyShares() {
			encryptedShare := vault_api.EncryptedShares{
				EncryptionKey: decryptionShare.GetEncryptionKey(),
				Shares:        make([]string, 0, len(decryptionShare.Shares)),
			}
			encryptedShare.Shares = append(encryptedShare.Shares, decryptionShare.GetShares()...)

			vaultAPIResponse.SecretValue.EncryptedDecryptionKeyShares = append(vaultAPIResponse.SecretValue.EncryptedDecryptionKeyShares, &encryptedShare)
		}
	case *vault.SecretResponse_Error:
		vaultAPIResponse.Error = method.Error
	}

	vaultAPIResponseBytes, err := json.Marshal(vaultAPIResponse)
	if err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.NodeReponseEncodingError, err)
	}
	vaultAPIResponseJSON := json.RawMessage(vaultAPIResponseBytes)
	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      req.ID,
		Method:  req.Method,
		Result:  &vaultAPIResponseJSON,
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

func toJSONResponse(vaultCapResponse *vault2.Response, method string) (*jsonrpc.Response[json.RawMessage], error) {
	vaultResponse := &vault_api.ResponseBase{
		ID:    vaultCapResponse.ID,
		Error: vaultCapResponse.Error,
		Response: vault_api.SignedResponse{
			Payload:    vaultCapResponse.Payload,
			Context:    vaultCapResponse.Context,
			Signatures: vaultCapResponse.Signatures,
		},
	}

	vaultResponseBytes, err := json.Marshal(vaultResponse)
	if err != nil {
		return nil, errors.New("failed to marshal vault response: " + err.Error())
	}
	vaultResponseJSON := json.RawMessage(vaultResponseBytes)
	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      vaultResponse.ID,
		Method:  method,
		Result:  &vaultResponseJSON,
	}, nil
}
