package vault

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

// GatewayVaultRequestProcessor orchestrates the shared gateway-routed vault JSON-RPC pipeline
// used by the gateway public handler and the node-side gateway connector handler.
//
// Pipeline invariant:
//
//	ValidateStructureBeforeAuth → NormalizeParamsForAuth → AuthorizeRequest → Prefix ID → StampAuthorizedParams → ValidateOwnerScopedLimits
//	    (no param mutation)        (gateway: check form;     (on raw bytes)               (namespace + request_id)      (ciphertext size)
//	                                node: re-derive params)
//
// AuthorizeRequest runs while params are still digest-safe. It also applies the replay guard
// (digest deduplication) and validates that payload owners match the authorized workflow owner
// before this processor rewrites the request ID or stamps params.
//
// Every Vault node authorizes the request again. The params it receives are the gateway's stamped
// params, so the node re-derives the caller's params from their typed struct and authorizes the
// digest of that. The gateway therefore rejects, before authorizing, params that the nodes would
// re-derive differently (unknown fields, a missing namespace, a request_id other than the request
// ID). Otherwise the gateway would authorize a digest that no node can reproduce, and each node
// would hold the request for the whole allowlist retry window.
//
// Owner-scoped limit checks are deferred until after authorization: each new owner tenant
// registered by a scoped limiter spawns a persistent background updater, so checking them
// pre-auth would let unauthenticated callers create unbounded limiter tenants.
type GatewayVaultRequestProcessor struct {
	validator               *RequestValidator
	authorizer              Authorizer
	stripOwnerPrefixForAuth bool
	lggr                    logger.Logger
}

// AuthorizedGatewayVaultRequest is a gateway-routed vault request after validation, authorization,
// owner-prefixed ID stamping, and param normalization.
type AuthorizedGatewayVaultRequest struct {
	Req        jsonrpc.Request[json.RawMessage]
	AuthResult *AuthResult
}

// NewGatewayVaultRequestProcessor constructs the shared gateway vault request processor.
func NewGatewayVaultRequestProcessor(validator *RequestValidator, authorizer Authorizer, stripOwnerPrefixForAuth bool, lggr logger.Logger) (*GatewayVaultRequestProcessor, error) {
	if validator == nil || authorizer == nil {
		return nil, errors.New("validator and authorizer are required to construct a gateway vault request processor")
	}
	return &GatewayVaultRequestProcessor{
		validator:               validator,
		authorizer:              authorizer,
		stripOwnerPrefixForAuth: stripOwnerPrefixForAuth,
		lggr:                    logger.Named(lggr, "GatewayVaultRequestProcessor"),
	}, nil
}

// Close releases validator-owned limiter resources.
func (p *GatewayVaultRequestProcessor) Close() error {
	return p.validator.Close()
}

// ProcessRequest runs validate structure → normalize params → authorize → prefix ID → stamp params → validate owner-scoped limits.
// In node mode (stripOwnerPrefixForAuth), the envelope ID received from the gateway (owner::id)
// is restored on error so error responses remain matchable by the gateway, which tracks active
// requests by the exact envelope ID it forwarded to the nodes.
func (p *GatewayVaultRequestProcessor) ProcessRequest(
	ctx context.Context,
	req *jsonrpc.Request[json.RawMessage],
	publicKey *tdh2easy.PublicKey,
) (*AuthorizedGatewayVaultRequest, error) {
	if !p.stripOwnerPrefixForAuth {
		return p.processRequest(ctx, req, publicKey)
	}

	envelopeID := req.ID
	strippedID, _ := stripPrefixedVaultRequestID(req.ID)
	req.ID = strippedID
	authorized, err := p.processRequest(ctx, req, publicKey)
	if err != nil {
		req.ID = envelopeID
		return nil, err
	}
	return authorized, nil
}

func (p *GatewayVaultRequestProcessor) processRequest(
	ctx context.Context,
	req *jsonrpc.Request[json.RawMessage],
	publicKey *tdh2easy.PublicKey,
) (*AuthorizedGatewayVaultRequest, error) {
	switch req.Method {
	case vaulttypes.MethodSecretsCreate:
		return p.processCreateSecretsRequest(ctx, req, publicKey)
	case vaulttypes.MethodSecretsUpdate:
		return p.processUpdateSecretsRequest(ctx, req, publicKey)
	case vaulttypes.MethodSecretsDelete:
		return p.processDeleteSecretsRequest(ctx, req)
	case vaulttypes.MethodSecretsList:
		return p.processListSecretIdentifiersRequest(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported gateway vault method: %s", req.Method)
	}
}

func (p *GatewayVaultRequestProcessor) processCreateSecretsRequest(
	ctx context.Context,
	req *jsonrpc.Request[json.RawMessage],
	publicKey *tdh2easy.PublicKey,
) (*AuthorizedGatewayVaultRequest, error) {
	if req.Params == nil {
		return nil, InvalidVaultParamsError{Method: req.Method, Err: errors.New("request params must not be nil")}
	}

	var createReq vaultcommon.CreateSecretsRequest
	if err := json.Unmarshal(*req.Params, &createReq); err != nil {
		return nil, InvalidVaultParamsError{Method: req.Method, Err: err}
	}
	createReq.RequestId = req.ID

	skipLabelValidation := publicKey == nil
	if err := p.validator.ValidateEncryptedSecretsStructure(ctx, publicKey, createReq.RequestId, createReq.EncryptedSecrets, skipLabelValidation); err != nil {
		return nil, p.validationError(req, err)
	}

	if err := p.normalizeParamsForAuth(req, &createReq, func() {
		vaultutils.ApplyEncryptedSecretNamespaceDefaults(createReq.EncryptedSecrets)
	}); err != nil {
		return nil, InvalidVaultParamsError{Method: req.Method, Err: err}
	}

	authorized, err := p.authorizeAndStamp(ctx, req, func(prefixedRequestID string) error {
		createReq.RequestId = prefixedRequestID
		vaultutils.ApplyEncryptedSecretNamespaceDefaults(createReq.EncryptedSecrets)
		return marshalVaultParams(req, &createReq)
	})
	if err != nil {
		return nil, err
	}

	if err := p.validator.ValidateCiphertextSizes(ctx, authorized.AuthResult.AuthorizedOwner(), createReq.EncryptedSecrets); err != nil {
		return nil, p.validationError(req, err)
	}
	return authorized, nil
}

func (p *GatewayVaultRequestProcessor) processUpdateSecretsRequest(
	ctx context.Context,
	req *jsonrpc.Request[json.RawMessage],
	publicKey *tdh2easy.PublicKey,
) (*AuthorizedGatewayVaultRequest, error) {
	if req.Params == nil {
		return nil, InvalidVaultParamsError{Method: req.Method, Err: errors.New("request params must not be nil")}
	}

	var updateReq vaultcommon.UpdateSecretsRequest
	if err := json.Unmarshal(*req.Params, &updateReq); err != nil {
		return nil, InvalidVaultParamsError{Method: req.Method, Err: err}
	}
	updateReq.RequestId = req.ID

	skipLabelValidation := publicKey == nil
	if err := p.validator.ValidateEncryptedSecretsStructure(ctx, publicKey, updateReq.RequestId, updateReq.EncryptedSecrets, skipLabelValidation); err != nil {
		return nil, p.validationError(req, err)
	}

	if err := p.normalizeParamsForAuth(req, &updateReq, func() {
		vaultutils.ApplyEncryptedSecretNamespaceDefaults(updateReq.EncryptedSecrets)
	}); err != nil {
		return nil, InvalidVaultParamsError{Method: req.Method, Err: err}
	}

	authorized, err := p.authorizeAndStamp(ctx, req, func(prefixedRequestID string) error {
		updateReq.RequestId = prefixedRequestID
		vaultutils.ApplyEncryptedSecretNamespaceDefaults(updateReq.EncryptedSecrets)
		return marshalVaultParams(req, &updateReq)
	})
	if err != nil {
		return nil, err
	}

	if err := p.validator.ValidateCiphertextSizes(ctx, authorized.AuthResult.AuthorizedOwner(), updateReq.EncryptedSecrets); err != nil {
		return nil, p.validationError(req, err)
	}
	return authorized, nil
}

func (p *GatewayVaultRequestProcessor) processDeleteSecretsRequest(
	ctx context.Context,
	req *jsonrpc.Request[json.RawMessage],
) (*AuthorizedGatewayVaultRequest, error) {
	if req.Params == nil {
		return nil, InvalidVaultParamsError{Method: req.Method, Err: errors.New("request params must not be nil")}
	}

	var deleteReq vaultcommon.DeleteSecretsRequest
	if err := json.Unmarshal(*req.Params, &deleteReq); err != nil {
		return nil, InvalidVaultParamsError{Method: req.Method, Err: err}
	}
	deleteReq.RequestId = req.ID

	if err := p.validator.ValidateDeleteSecretsRequest(ctx, &deleteReq); err != nil {
		return nil, p.validationError(req, err)
	}

	if err := p.normalizeParamsForAuth(req, &deleteReq, func() {
		vaultutils.ApplySecretIdentifierNamespaceDefaults(deleteReq.Ids)
	}); err != nil {
		return nil, InvalidVaultParamsError{Method: req.Method, Err: err}
	}

	return p.authorizeAndStamp(ctx, req, func(prefixedRequestID string) error {
		deleteReq.RequestId = prefixedRequestID
		vaultutils.ApplySecretIdentifierNamespaceDefaults(deleteReq.Ids)
		return marshalVaultParams(req, &deleteReq)
	})
}

func (p *GatewayVaultRequestProcessor) processListSecretIdentifiersRequest(
	ctx context.Context,
	req *jsonrpc.Request[json.RawMessage],
) (*AuthorizedGatewayVaultRequest, error) {
	if req.Params == nil {
		return nil, InvalidVaultParamsError{Method: req.Method, Err: errors.New("request params must not be nil")}
	}

	var listReq vaultcommon.ListSecretIdentifiersRequest
	if err := json.Unmarshal(*req.Params, &listReq); err != nil {
		return nil, InvalidVaultParamsError{Method: req.Method, Err: err}
	}
	listReq.RequestId = req.ID

	if err := p.validator.ValidateListSecretIdentifiersRequest(ctx, &listReq); err != nil {
		return nil, p.validationError(req, err)
	}

	if err := p.normalizeParamsForAuth(req, &listReq, func() {
		applyListNamespaceDefault(&listReq)
	}); err != nil {
		return nil, InvalidVaultParamsError{Method: req.Method, Err: err}
	}

	return p.authorizeAndStamp(ctx, req, func(prefixedRequestID string) error {
		listReq.RequestId = prefixedRequestID
		applyListNamespaceDefault(&listReq)
		return marshalVaultParams(req, &listReq)
	})
}

func applyListNamespaceDefault(listReq *vaultcommon.ListSecretIdentifiersRequest) {
	if listReq.Namespace == "" {
		listReq.Namespace = vaulttypes.DefaultNamespace
	}
}

// errParamsNotNodeDerivable is returned when the gateway receives params that the Vault nodes
// would re-derive differently, so the digest the gateway authorizes would never match theirs.
var errParamsNotNodeDerivable = errors.New("params must be in canonical form: request_id equal to the request ID, every namespace set, and no fields the request type does not define")

// normalizeParamsForAuth brings params to the form every Vault node authorizes: the typed struct,
// with RequestId already set to the envelope ID by the caller. In node mode req.Params is rewritten
// to that form. In gateway mode applyDefaults first fills in what the gateway stamps before
// forwarding (default namespaces), and the request is rejected unless the caller sent exactly that
// form. req.Params is left as sent, so the gateway still authorizes the caller's raw bytes.
func (p *GatewayVaultRequestProcessor) normalizeParamsForAuth(req *jsonrpc.Request[json.RawMessage], params any, applyDefaults func()) error {
	if p.stripOwnerPrefixForAuth {
		return marshalVaultParams(req, params)
	}
	applyDefaults()

	derived := *req
	if err := marshalVaultParams(&derived, params); err != nil {
		return err
	}
	sentDigest, err := req.Digest()
	if err != nil {
		return err
	}
	derivedDigest, err := derived.Digest()
	if err != nil {
		return err
	}
	if sentDigest != derivedDigest {
		return errParamsNotNodeDerivable
	}
	return nil
}

func (p *GatewayVaultRequestProcessor) authorizeAndStamp(
	ctx context.Context,
	req *jsonrpc.Request[json.RawMessage],
	stamp func(prefixedRequestID string) error,
) (*AuthorizedGatewayVaultRequest, error) {
	incomingOwner := ""
	if idx := strings.Index(req.ID, vaulttypes.RequestIDSeparator); idx != -1 {
		incomingOwner = req.ID[:idx]
	}

	p.lggr.Debugw("authorizing gateway vault request", "method", req.Method, "requestID", req.ID)
	authResult, err := p.authorizer.AuthorizeRequest(ctx, *req)
	if err != nil {
		authErr := fmt.Errorf("request not authorized: %w", err)
		p.lggr.Errorw("gateway vault request authorization failed", "method", req.Method, "requestID", req.ID, "hasAuth", req.Auth != "", "incomingOwner", incomingOwner, "error", authErr)
		return nil, authErr
	}

	originalRequestID := req.ID
	authorizedOwner := authResult.AuthorizedOwner()
	prefixedRequestID := authorizedOwner + vaulttypes.RequestIDSeparator + originalRequestID
	req.ID = prefixedRequestID

	if err := stamp(prefixedRequestID); err != nil {
		p.lggr.Errorw("failed to stamp authorized request params", "method", req.Method, "requestID", req.ID, "error", err)
		return nil, fmt.Errorf("failed to stamp authorized request params: %w", err)
	}

	p.lggr.Debugw("authorized gateway vault request", "method", req.Method, "requestID", req.ID, "owner", authorizedOwner, "orgID", authResult.OrgID(), "workflowOwner", authResult.WorkflowOwner())
	return &AuthorizedGatewayVaultRequest{
		Req:        *req,
		AuthResult: authResult,
	}, nil
}

func (p *GatewayVaultRequestProcessor) validationError(req *jsonrpc.Request[json.RawMessage], err error) error {
	invalidErr := InvalidVaultParamsError{Method: req.Method, Err: err}
	if IsInvalidVaultParamsError(err) {
		p.lggr.Warnw("gateway vault request validation failed", "method", req.Method, "requestID", req.ID, "error", invalidErr)
	} else {
		p.lggr.Errorw("failed to validate gateway vault request before authorization", "method", req.Method, "requestID", req.ID, "error", invalidErr)
	}
	return invalidErr
}

func stripPrefixedVaultRequestID(requestID string) (originalRequestID, prefixedOwner string) {
	prefixedOwner, originalRequestID, ok := strings.Cut(requestID, vaulttypes.RequestIDSeparator)
	if !ok {
		return requestID, ""
	}
	return originalRequestID, prefixedOwner
}

func marshalVaultParams(req *jsonrpc.Request[json.RawMessage], params any) error {
	marshaled, err := json.Marshal(params)
	if err != nil {
		return err
	}
	raw := json.RawMessage(marshaled)
	req.Params = &raw
	return nil
}

// MasterPublicKeyFromSecretsService loads the vault master public key from a secrets service.
func MasterPublicKeyFromSecretsService(ctx context.Context, secretsService vaulttypes.SecretsService) (*tdh2easy.PublicKey, error) {
	resp, err := secretsService.GetPublicKey(ctx, &vaultcommon.GetPublicKeyRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get vault public key: %w", err)
	}
	if resp == nil || resp.PublicKey == "" {
		return nil, errors.New("vault public key is unavailable")
	}

	masterPublicKeyBytes, err := hex.DecodeString(resp.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode vault public key: %w", err)
	}

	masterPublicKey := &tdh2easy.PublicKey{}
	if err := masterPublicKey.Unmarshal(masterPublicKeyBytes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vault public key: %w", err)
	}
	return masterPublicKey, nil
}
