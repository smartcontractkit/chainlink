package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	workflowsyncerv2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/v2"
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

// GatewayHandler serves Vault requests received from the gateway on the node side.
type GatewayHandler struct {
	services.Service
	eng *services.Engine

	secretsService    vaulttypes.SecretsService
	gatewayConnector  gatewayConnector
	authorizer        Authorizer
	requestValidator  *RequestValidator
	jwtAuthService    services.Service
	lggr              logger.Logger
	metrics           *metrics

	// TODO add org resolver? https://smartcontract-it.atlassian.net/browse/CRE-1707
}

// NewGatewayHandler creates a Vault gateway connector handler with internal auth wiring.
// Pass a non-nil authorizer only in tests or other cases that need to override the default
// allowlist/JWT authorization chain.
func NewGatewayHandler(
	secretsService vaulttypes.SecretsService,
	connector gatewayConnector,
	workflowRegistrySyncer workflowsyncerv2.WorkflowRegistrySyncer,
	lggr logger.Logger,
	limitsFactory limits.Factory,
	authorizer Authorizer,
	auth0 *Auth0Config,
) (*GatewayHandler, error) {
	var jwtAuthService services.Service
	var jwtBasedAuth Authorizer
	if auth0 != nil {
		var err error
		jwtAuthService, err = NewJWTBasedAuth(JWTBasedAuthConfig{
			IssuerURL: auth0.IssuerURL,
			Audience:  auth0.Audience,
			TenantID:  auth0.TenantID,
		}, limitsFactory, lggr)
		if err != nil {
			return nil, fmt.Errorf("failed to create JWTBasedAuth: %w", err)
		}
		jwtBasedAuth = jwtAuthService.(Authorizer)
	}

	if authorizer == nil {
		allowListBasedAuth := NewAllowListBasedAuth(lggr, workflowRegistrySyncer)
		authorizer = NewAuthorizer(allowListBasedAuth, jwtBasedAuth, lggr)
	}

	requestValidator, err := NewRequestValidatorFromLimitsFactory(limitsFactory)
	if err != nil {
		return nil, fmt.Errorf("failed to create request validator: %w", err)
	}

	metrics, err := newMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics: %w", err)
	}

	gh := &GatewayHandler{
		secretsService:   secretsService,
		gatewayConnector: connector,
		authorizer:       authorizer,
		requestValidator: requestValidator,
		jwtAuthService:   jwtAuthService,
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
	if h.jwtAuthService != nil {
		if err := h.jwtAuthService.Start(ctx); err != nil {
			return fmt.Errorf("failed to start JWTBasedAuth: %w", err)
		}
	}
	if gwerr := h.gatewayConnector.AddHandler(ctx, h.Methods(), h); gwerr != nil {
		return fmt.Errorf("failed to add vault handler to connector: %w", gwerr)
	}
	return nil
}

func (h *GatewayHandler) close() error {
	var jwtAuthErr error
	if h.jwtAuthService != nil {
		jwtAuthErr = h.jwtAuthService.Close()
	}
	var validatorErr error
	if h.requestValidator != nil {
		validatorErr = h.requestValidator.Close()
	}
	if gwerr := h.gatewayConnector.RemoveHandler(context.Background(), h.Methods()); gwerr != nil {
		return errors.Join(fmt.Errorf("failed to remove vault handler from connector: %w", gwerr), jwtAuthErr, validatorErr)
	}
	return errors.Join(jwtAuthErr, validatorErr)
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
		if _, prepErr := h.validateAuthorizeAndPrefix(ctx, req); prepErr != nil {
			response = h.gatewayErrorResponse(ctx, gatewayID, req, prepErr)
			break
		}
		response = h.handleSecretsCreate(ctx, gatewayID, req)
	case vaulttypes.MethodSecretsUpdate:
		if _, prepErr := h.validateAuthorizeAndPrefix(ctx, req); prepErr != nil {
			response = h.gatewayErrorResponse(ctx, gatewayID, req, prepErr)
			break
		}
		response = h.handleSecretsUpdate(ctx, gatewayID, req)
	case vaulttypes.MethodSecretsDelete:
		if _, prepErr := h.validateAuthorizeAndPrefix(ctx, req); prepErr != nil {
			response = h.gatewayErrorResponse(ctx, gatewayID, req, prepErr)
			break
		}
		response = h.handleSecretsDelete(ctx, gatewayID, req)
	case vaulttypes.MethodSecretsList:
		authResult, prepErr := h.validateAuthorizeAndPrefix(ctx, req)
		if prepErr != nil {
			response = h.gatewayErrorResponse(ctx, gatewayID, req, prepErr)
			break
		}
		response = h.handleSecretsList(ctx, gatewayID, req, authResult)
	case vaulttypes.MethodPublicKeyGet:
		response = h.handlePublicKeyGet(ctx, gatewayID, req)
	default:
		response = h.errorResponse(ctx, gatewayID, req, api.UnsupportedMethodError, errors.New("unsupported method: "+req.Method))
	}

	if err = h.gatewayConnector.SendToGateway(ctx, gatewayID, response); err != nil {
		h.lggr.Errorw("Failed to send message to gateway", "gatewayID", gatewayID, "error", err)
		return err
	}

	h.lggr.Infow("Sent message to gateway", "gatewayID", gatewayID, "resp", response, "requestID", req.ID)
	h.metrics.requestSuccess.Add(ctx, 1, metric.WithAttributes(
		attribute.String("gateway_id", gatewayID),
	))
	return nil
}

func (h *GatewayHandler) validateAuthorizeAndPrefix(ctx context.Context, req *jsonrpc.Request[json.RawMessage]) (*AuthResult, error) {
	if h.authorizer == nil {
		err := errors.New("authorizer is nil")
		h.lggr.Errorw("failed to authorize gateway request", "method", req.Method, "requestID", req.ID, "error", err)
		return nil, err
	}

	incomingOwner := ""
	if idx := strings.Index(req.ID, vaulttypes.RequestIDSeparator); idx != -1 {
		incomingOwner = req.ID[:idx]
	}

	if err := h.requestValidator.PrepareUserJSONRPCRequest(ctx, req, UserJSONRPCValidationOptions{
		SkipLabelValidation: true,
	}, true); err != nil {
		if IsInvalidVaultParamsError(err) {
			h.lggr.Warnw("gateway request validation failed", "method", req.Method, "requestID", req.ID, "error", err)
		} else {
			h.lggr.Errorw("failed to prepare gateway request for authorization", "method", req.Method, "requestID", req.ID, "error", err)
		}
		return nil, err
	}

	h.lggr.Debugw("authorizing gateway request", "method", req.Method, "requestID", req.ID)
	authResult, err := h.authorizer.AuthorizeRequest(ctx, *req)
	if err != nil {
		authErr := fmt.Errorf("request not authorized: %w", err)
		h.lggr.Errorw("gateway request authorization failed", "method", req.Method, "requestID", req.ID, "hasAuth", req.Auth != "", "incomingOwner", incomingOwner, "error", authErr)
		return nil, authErr
	}

	originalRequestID := req.ID
	authorizedOwner := authResult.AuthorizedOwner()
	req.ID = authorizedOwner + vaulttypes.RequestIDSeparator + originalRequestID
	if err := h.requestValidator.FinalizeAuthorizedJSONRPCRequest(req, req.ID); err != nil {
		h.lggr.Errorw("failed to stamp prefixed request ID in params", "method", req.Method, "requestID", req.ID, "error", err)
		return nil, fmt.Errorf("failed to stamp prefixed request ID in params: %w", err)
	}

	h.lggr.Debugw("authorized gateway request", "method", req.Method, "requestID", req.ID, "owner", authorizedOwner, "orgID", authResult.OrgID(), "workflowOwner", authResult.WorkflowOwner())
	return authResult, nil
}

func (h *GatewayHandler) gatewayErrorResponse(
	ctx context.Context,
	gatewayID string,
	req *jsonrpc.Request[json.RawMessage],
	err error,
) *jsonrpc.Response[json.RawMessage] {
	if IsInvalidVaultParamsError(err) {
		return h.errorResponse(ctx, gatewayID, req, api.InvalidParamsError, errors.New("invalid params error: "+err.Error()))
	}
	return h.errorResponse(ctx, gatewayID, req, api.HandlerError, err)
}

func (h *GatewayHandler) handleSecretsCreate(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) *jsonrpc.Response[json.RawMessage] {
	vaultCapRequest := vaultcommon.CreateSecretsRequest{}
	if err := json.Unmarshal(*req.Params, &vaultCapRequest); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.UserMessageParseError, err)
	}

	h.lggr.Debugw("Processing authorized create secrets request", "request", vaultCapRequest.String())
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

	h.lggr.Debugw("Processing authorized update secrets request", "request", vaultCapRequest.String())
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

	h.lggr.Debugw("Processing authorized delete secrets request", "request", r.String())
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

func (h *GatewayHandler) handleSecretsList(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage], authResult *AuthResult) *jsonrpc.Response[json.RawMessage] {
	r := &vaultcommon.ListSecretIdentifiersRequest{}
	if err := json.Unmarshal(*req.Params, r); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.UserMessageParseError, err)
	}
	r.Owner = authResult.AuthorizedOwner()

	h.lggr.Debugw("Processing authorized list secrets request", "request", r.String())
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
	h.lggr.Errorw("gateway handler error response", "gatewayID", gatewayID, "requestID", req.ID, "method", req.Method, "errorCode", errorCode, "error", err)
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
