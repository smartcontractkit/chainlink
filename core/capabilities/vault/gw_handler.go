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

type gatewayHandlerConfig struct {
	authorizer Authorizer
}

// GatewayHandlerOption customizes GatewayHandler construction for tests and future auth extensions.
type GatewayHandlerOption func(*gatewayHandlerConfig)

// WithAuthorizer overrides the default Vault request authorizer.
func WithAuthorizer(authorizer Authorizer) GatewayHandlerOption {
	return func(cfg *gatewayHandlerConfig) {
		cfg.authorizer = authorizer
	}
}

// GatewayHandler serves Vault requests received from the gateway on the node side.
type GatewayHandler struct {
	services.Service
	eng *services.Engine

	secretsService   vaulttypes.SecretsService
	gatewayConnector gatewayConnector
	authorizer       Authorizer
	jwtAuthService   services.Service
	lggr             logger.Logger
	metrics          *metrics
}

// NewGatewayHandler creates a Vault gateway connector handler with internal auth wiring.
func NewGatewayHandler(secretsService vaulttypes.SecretsService, connector gatewayConnector, workflowRegistrySyncer workflowsyncerv2.WorkflowRegistrySyncer, lggr logger.Logger, limitsFactory limits.Factory, opts ...GatewayHandlerOption) (*GatewayHandler, error) {
	cfg := gatewayHandlerConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.authorizer == nil {
		allowListBasedAuth := NewAllowListBasedAuth(lggr, workflowRegistrySyncer)
		jwtBasedAuth, err := NewJWTBasedAuth(JWTBasedAuthConfig{}, limitsFactory, lggr, WithDisabledJWTBasedAuth())
		if err != nil {
			return nil, fmt.Errorf("failed to create JWTBasedAuth: %w", err)
		}
		cfg.authorizer = NewAuthorizer(allowListBasedAuth, jwtBasedAuth, lggr)
		return newGatewayHandlerWithAuthorizer(secretsService, connector, cfg.authorizer, jwtBasedAuth, lggr)
	}
	return newGatewayHandlerWithAuthorizer(secretsService, connector, cfg.authorizer, nil, lggr)
}

func newGatewayHandlerWithAuthorizer(secretsService vaulttypes.SecretsService, connector gatewayConnector, authorizer Authorizer, jwtAuthService services.Service, lggr logger.Logger) (*GatewayHandler, error) {
	metrics, err := newMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics: %w", err)
	}

	gh := &GatewayHandler{
		secretsService:   secretsService,
		gatewayConnector: connector,
		authorizer:       authorizer,
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
	if gwerr := h.gatewayConnector.RemoveHandler(context.Background(), h.Methods()); gwerr != nil {
		return errors.Join(fmt.Errorf("failed to remove vault handler from connector: %w", gwerr), jwtAuthErr)
	}
	return jwtAuthErr
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
		authResult, authErr := h.authorizeAndPrefixRequest(ctx, req)
		if authErr != nil {
			response = h.errorResponse(ctx, gatewayID, req, api.HandlerError, authErr)
			break
		}
		response = h.handleSecretsCreate(ctx, gatewayID, req, authResult)
	case vaulttypes.MethodSecretsUpdate:
		authResult, authErr := h.authorizeAndPrefixRequest(ctx, req)
		if authErr != nil {
			response = h.errorResponse(ctx, gatewayID, req, api.HandlerError, authErr)
			break
		}
		response = h.handleSecretsUpdate(ctx, gatewayID, req, authResult)
	case vaulttypes.MethodSecretsDelete:
		authResult, authErr := h.authorizeAndPrefixRequest(ctx, req)
		if authErr != nil {
			response = h.errorResponse(ctx, gatewayID, req, api.HandlerError, authErr)
			break
		}
		response = h.handleSecretsDelete(ctx, gatewayID, req, authResult)
	case vaulttypes.MethodSecretsList:
		authResult, authErr := h.authorizeAndPrefixRequest(ctx, req)
		if authErr != nil {
			response = h.errorResponse(ctx, gatewayID, req, api.HandlerError, authErr)
			break
		}
		response = h.handleSecretsList(ctx, gatewayID, req, authResult)
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

func (h *GatewayHandler) authorizeAndPrefixRequest(ctx context.Context, req *jsonrpc.Request[json.RawMessage]) (*AuthResult, error) {
	if h.authorizer == nil {
		err := errors.New("authorizer is nil")
		h.lggr.Errorw("failed to authorize gateway request", "method", req.Method, "requestID", req.ID, "error", err)
		return nil, err
	}

	originalRequestID := req.ID
	incomingOwner := ""
	if idx := strings.Index(req.ID, vaulttypes.RequestIDSeparator); idx != -1 {
		incomingOwner = req.ID[:idx]
		originalRequestID = req.ID[idx+len(vaulttypes.RequestIDSeparator):]
	}

	authReq := *req
	authReq.ID = originalRequestID
	if err := stripPrefixedRequestIDFromParams(&authReq, originalRequestID); err != nil {
		h.lggr.Errorw("failed to normalize gateway request for authorization", "method", req.Method, "requestID", originalRequestID, "error", err)
		return nil, err
	}

	h.lggr.Debugw("authorizing gateway request", "method", req.Method, "requestID", originalRequestID)
	authResult, err := h.authorizer.AuthorizeRequest(ctx, authReq)
	if err != nil {
		authErr := fmt.Errorf("request not authorized: %w", err)
		h.lggr.Errorw("gateway request authorization failed", "method", req.Method, "requestID", originalRequestID, "hasAuth", req.Auth != "", "incomingOwner", incomingOwner, "error", authErr)
		return nil, authErr
	}
	authorizedOwner := authResult.AuthorizedOwner()
	if incomingOwner != "" && normalizeOwner(incomingOwner) != normalizeOwner(authorizedOwner) {
		prefixErr := fmt.Errorf("request owner prefix %q does not match authorized owner %q", incomingOwner, authorizedOwner)
		h.lggr.Errorw("gateway request owner prefix mismatch", "method", req.Method, "requestID", originalRequestID, "incomingOwner", incomingOwner, "authorizedOwner", authorizedOwner, "error", prefixErr)
		return nil, prefixErr
	}

	req.ID = authorizedOwner + vaulttypes.RequestIDSeparator + originalRequestID
	h.lggr.Debugw("authorized gateway request", "method", req.Method, "requestID", req.ID, "owner", authorizedOwner, "orgID", authResult.OrgID(), "workflowOwner", authResult.WorkflowOwner())
	return authResult, nil
}

func stripPrefixedRequestIDFromParams(req *jsonrpc.Request[json.RawMessage], originalRequestID string) error {
	if req.Params == nil {
		return nil
	}

	switch req.Method {
	case vaulttypes.MethodSecretsCreate:
		parsed := &vaultcommon.CreateSecretsRequest{}
		if err := json.Unmarshal(*req.Params, parsed); err != nil {
			return err
		}
		parsed.RequestId = originalRequestID
		return rewriteRequestParams(req, parsed)
	case vaulttypes.MethodSecretsUpdate:
		parsed := &vaultcommon.UpdateSecretsRequest{}
		if err := json.Unmarshal(*req.Params, parsed); err != nil {
			return err
		}
		parsed.RequestId = originalRequestID
		return rewriteRequestParams(req, parsed)
	case vaulttypes.MethodSecretsDelete:
		parsed := &vaultcommon.DeleteSecretsRequest{}
		if err := json.Unmarshal(*req.Params, parsed); err != nil {
			return err
		}
		parsed.RequestId = originalRequestID
		return rewriteRequestParams(req, parsed)
	case vaulttypes.MethodSecretsList:
		parsed := &vaultcommon.ListSecretIdentifiersRequest{}
		if err := json.Unmarshal(*req.Params, parsed); err != nil {
			return err
		}
		parsed.RequestId = originalRequestID
		return rewriteRequestParams(req, parsed)
	default:
		return nil
	}
}

func rewriteRequestParams(req *jsonrpc.Request[json.RawMessage], payload any) error {
	params, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	raw := json.RawMessage(params)
	req.Params = &raw
	return nil
}

func setAuthorizedIdentityFields(req any, authResult *AuthResult) error {
	if authResult == nil {
		return errors.New("auth result is nil")
	}

	// Critical: the Vault capability trusts OrgId and WorkflowOwner to be set by
	// the Vault node only after authorization and request validation succeed. We
	// must overwrite any JSON-provided values here; otherwise a malicious request
	// could smuggle mismatched identity fields into the capability call.
	switch r := req.(type) {
	case *vaultcommon.CreateSecretsRequest:
		r.OrgId = authResult.OrgID()
		r.WorkflowOwner = authResult.WorkflowOwner()
	case *vaultcommon.UpdateSecretsRequest:
		r.OrgId = authResult.OrgID()
		r.WorkflowOwner = authResult.WorkflowOwner()
	case *vaultcommon.DeleteSecretsRequest:
		r.OrgId = authResult.OrgID()
		r.WorkflowOwner = authResult.WorkflowOwner()
	case *vaultcommon.ListSecretIdentifiersRequest:
		r.OrgId = authResult.OrgID()
		r.WorkflowOwner = authResult.WorkflowOwner()
	default:
		return fmt.Errorf("unsupported vault request type %T", req)
	}

	return nil
}

func (h *GatewayHandler) handleSecretsCreate(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage], authResult *AuthResult) *jsonrpc.Response[json.RawMessage] {
	vaultCapRequest := vaultcommon.CreateSecretsRequest{}
	if err := json.Unmarshal(*req.Params, &vaultCapRequest); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.UserMessageParseError, err)
	}

	authorizedOwner := authResult.AuthorizedOwner()
	vaultCapRequest.RequestId = req.ID
	if err := setAuthorizedIdentityFields(&vaultCapRequest, authResult); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.FatalError, err)
	}
	for idx, encryptedSecret := range vaultCapRequest.EncryptedSecrets {
		if encryptedSecret != nil && encryptedSecret.Id != nil && normalizeOwner(encryptedSecret.Id.Owner) != normalizeOwner(authorizedOwner) {
			h.lggr.Debugw("create secrets request owner mismatch", "requestID", req.ID, "secretOwner", encryptedSecret.Id.Owner, "authorizedOwner", authorizedOwner, "index", idx)
			return h.errorResponse(ctx, gatewayID, req, api.FatalError, fmt.Errorf("secret ID owner %q does not match authorized owner %q at index %d", encryptedSecret.Id.Owner, authorizedOwner, idx))
		}
	}

	h.lggr.Debugf("Processing authorized and normalized create secrets request [%s]", vaultCapRequest.String())
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

func (h *GatewayHandler) handleSecretsUpdate(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage], authResult *AuthResult) *jsonrpc.Response[json.RawMessage] {
	vaultCapRequest := vaultcommon.UpdateSecretsRequest{}
	if err := json.Unmarshal(*req.Params, &vaultCapRequest); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.UserMessageParseError, err)
	}
	authorizedOwner := authResult.AuthorizedOwner()
	vaultCapRequest.RequestId = req.ID
	if err := setAuthorizedIdentityFields(&vaultCapRequest, authResult); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.FatalError, err)
	}
	for idx, encryptedSecret := range vaultCapRequest.EncryptedSecrets {
		if encryptedSecret != nil && encryptedSecret.Id != nil && normalizeOwner(encryptedSecret.Id.Owner) != normalizeOwner(authorizedOwner) {
			h.lggr.Debugw("update secrets request owner mismatch", "requestID", req.ID, "secretOwner", encryptedSecret.Id.Owner, "authorizedOwner", authorizedOwner, "index", idx)
			return h.errorResponse(ctx, gatewayID, req, api.FatalError, fmt.Errorf("secret ID owner %q does not match authorized owner %q at index %d", encryptedSecret.Id.Owner, authorizedOwner, idx))
		}
	}

	h.lggr.Debugf("Processing authorized and normalized update secrets request [%s]", vaultCapRequest.String())
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

func (h *GatewayHandler) handleSecretsDelete(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage], authResult *AuthResult) *jsonrpc.Response[json.RawMessage] {
	r := &vaultcommon.DeleteSecretsRequest{}
	if err := json.Unmarshal(*req.Params, r); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.UserMessageParseError, err)
	}
	authorizedOwner := authResult.AuthorizedOwner()
	r.RequestId = req.ID
	if err := setAuthorizedIdentityFields(r, authResult); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.FatalError, err)
	}
	for idx, secretID := range r.Ids {
		if secretID != nil && normalizeOwner(secretID.Owner) != normalizeOwner(authorizedOwner) {
			h.lggr.Debugw("delete secrets request owner mismatch", "requestID", req.ID, "secretOwner", secretID.Owner, "authorizedOwner", authorizedOwner, "index", idx)
			return h.errorResponse(ctx, gatewayID, req, api.FatalError, fmt.Errorf("secret ID owner %q does not match authorized owner %q at index %d", secretID.Owner, authorizedOwner, idx))
		}
	}

	h.lggr.Debugf("Processing authorized and normalized delete secrets request [%s]", r.String())
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
	r.RequestId = req.ID
	r.Owner = authResult.AuthorizedOwner()
	if err := setAuthorizedIdentityFields(r, authResult); err != nil {
		return h.errorResponse(ctx, gatewayID, req, api.FatalError, err)
	}

	h.lggr.Debugf("Processing authorized and normalized list secrets request [%s]", r.String())
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
