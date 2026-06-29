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
//	ValidateStructureBeforeAuth → AuthorizeRequest → Prefix ID → StampAuthorizedParams
//	    (no param mutation)        (on raw bytes)               (namespace + request_id)
//
// AuthorizeRequest runs while params are still digest-safe. It also applies the replay guard
// (digest deduplication) and validates that payload owners match the authorized workflow owner
// before this processor rewrites the request ID or stamps params.
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

// ProcessRequest runs validate → authorize → prefix ID → stamp params.
func (p *GatewayVaultRequestProcessor) ProcessRequest(
	ctx context.Context,
	req *jsonrpc.Request[json.RawMessage],
	publicKey *tdh2easy.PublicKey,
) (*AuthorizedGatewayVaultRequest, error) {
	if p.stripOwnerPrefixForAuth {
		originalRequestID, _ := stripPrefixedVaultRequestID(req.ID)
		req.ID = originalRequestID
	}

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
	if p.stripOwnerPrefixForAuth {
		createReq.RequestId = req.ID
		if err := marshalVaultParams(req, &createReq); err != nil {
			return nil, InvalidVaultParamsError{Method: req.Method, Err: err}
		}
	} else {
		createReq.RequestId = coalesceRequestID(createReq.RequestId, req.ID)
	}

	skipLabelValidation := publicKey == nil
	if err := p.validator.ValidateCreateSecretsRequest(ctx, publicKey, &createReq, skipLabelValidation); err != nil {
		return nil, p.validationError(req, err)
	}

	return p.authorizeAndStamp(ctx, req, func(prefixedRequestID string) error {
		createReq.RequestId = prefixedRequestID
		vaultutils.ApplyEncryptedSecretNamespaceDefaults(createReq.EncryptedSecrets)
		return marshalVaultParams(req, &createReq)
	})
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
	if p.stripOwnerPrefixForAuth {
		updateReq.RequestId = req.ID
		if err := marshalVaultParams(req, &updateReq); err != nil {
			return nil, InvalidVaultParamsError{Method: req.Method, Err: err}
		}
	} else {
		updateReq.RequestId = coalesceRequestID(updateReq.RequestId, req.ID)
	}

	skipLabelValidation := publicKey == nil
	if err := p.validator.ValidateUpdateSecretsRequest(ctx, publicKey, &updateReq, skipLabelValidation); err != nil {
		return nil, p.validationError(req, err)
	}

	return p.authorizeAndStamp(ctx, req, func(prefixedRequestID string) error {
		updateReq.RequestId = prefixedRequestID
		vaultutils.ApplyEncryptedSecretNamespaceDefaults(updateReq.EncryptedSecrets)
		return marshalVaultParams(req, &updateReq)
	})
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
	if p.stripOwnerPrefixForAuth {
		deleteReq.RequestId = req.ID
		if err := marshalVaultParams(req, &deleteReq); err != nil {
			return nil, InvalidVaultParamsError{Method: req.Method, Err: err}
		}
	} else {
		deleteReq.RequestId = coalesceRequestID(deleteReq.RequestId, req.ID)
	}

	if err := p.validator.ValidateDeleteSecretsRequest(ctx, &deleteReq); err != nil {
		return nil, p.validationError(req, err)
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
	if p.stripOwnerPrefixForAuth {
		listReq.RequestId = req.ID
		if err := marshalVaultParams(req, &listReq); err != nil {
			return nil, InvalidVaultParamsError{Method: req.Method, Err: err}
		}
	} else {
		listReq.RequestId = coalesceRequestID(listReq.RequestId, req.ID)
	}

	if err := p.validator.ValidateListSecretIdentifiersRequest(ctx, &listReq); err != nil {
		return nil, p.validationError(req, err)
	}

	return p.authorizeAndStamp(ctx, req, func(prefixedRequestID string) error {
		listReq.RequestId = prefixedRequestID
		if listReq.Namespace == "" {
			listReq.Namespace = vaulttypes.DefaultNamespace
		}
		return marshalVaultParams(req, &listReq)
	})
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

func coalesceRequestID(paramsRequestID, envelopeRequestID string) string {
	if paramsRequestID != "" {
		return paramsRequestID
	}
	return envelopeRequestID
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
