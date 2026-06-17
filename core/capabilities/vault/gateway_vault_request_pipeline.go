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
)

// GatewayVaultRequestPipeline orchestrates the shared gateway-routed vault JSON-RPC pipeline
// used by the gateway public handler and the node-side gateway connector handler.
//
// Pipeline invariant:
//
//	ValidateStructureBeforeAuth → AuthorizeRequest → Prefix ID → StampAuthorizedParams
//	    (no param mutation)        (on raw bytes)               (namespace + request_id)
//
// AuthorizeRequest runs while params are still digest-safe. It also applies the replay guard
// (digest deduplication) and validates that payload owners match the authorized workflow owner
// before this pipeline rewrites the request ID or stamps params.
type GatewayVaultRequestPipeline struct {
	validator  *RequestValidator
	authorizer Authorizer
	lggr       logger.Logger
}

// GatewayVaultRequestPipelineOptions configures pipeline behavior for gateway vs node re-auth paths.
type GatewayVaultRequestPipelineOptions struct {
	// StripOwnerPrefixForAuth removes gateway-added owner prefixes before digest verification.
	StripOwnerPrefixForAuth bool
	PublicKey               *tdh2easy.PublicKey
	SkipLabelValidation     bool
}

// AuthorizedGatewayVaultRequest is a gateway-routed vault request after validation, authorization,
// owner-prefixed ID stamping, and param normalization.
type AuthorizedGatewayVaultRequest struct {
	Req        jsonrpc.Request[json.RawMessage]
	AuthResult *AuthResult
}

// NewGatewayVaultRequestPipeline constructs the shared gateway vault request pipeline.
func NewGatewayVaultRequestPipeline(validator *RequestValidator, authorizer Authorizer, lggr logger.Logger) *GatewayVaultRequestPipeline {
	return &GatewayVaultRequestPipeline{
		validator:  validator,
		authorizer: authorizer,
		lggr:       logger.Named(lggr, "GatewayVaultRequestPipeline"),
	}
}

// Validator returns the pipeline request validator for lifecycle management.
func (p *GatewayVaultRequestPipeline) Validator() *RequestValidator {
	return p.validator
}

// Close releases validator-owned limiter resources.
func (p *GatewayVaultRequestPipeline) Close() error {
	if p == nil || p.validator == nil {
		return nil
	}
	return p.validator.Close()
}

// ProcessGatewayVaultRequest runs validate → authorize → prefix ID → stamp params.
func (p *GatewayVaultRequestPipeline) ProcessGatewayVaultRequest(
	ctx context.Context,
	req *jsonrpc.Request[json.RawMessage],
	opts GatewayVaultRequestPipelineOptions,
) (*AuthorizedGatewayVaultRequest, error) {
	if p == nil || p.authorizer == nil || p.validator == nil {
		err := errors.New("gateway vault request pipeline is not configured")
		p.lggr.Errorw("gateway vault request pipeline unavailable", "method", req.Method, "requestID", req.ID, "error", err)
		return nil, err
	}

	incomingOwner := ""
	if idx := strings.Index(req.ID, vaulttypes.RequestIDSeparator); idx != -1 {
		incomingOwner = req.ID[:idx]
	}

	validationOpts := UserJSONRPCValidationOptions{
		PublicKey:           opts.PublicKey,
		SkipLabelValidation: opts.SkipLabelValidation,
	}

	if err := p.validator.ValidateStructureBeforeAuth(ctx, req, validationOpts, opts.StripOwnerPrefixForAuth); err != nil {
		if IsInvalidVaultParamsError(err) {
			p.lggr.Warnw("gateway vault request validation failed", "method", req.Method, "requestID", req.ID, "error", err)
		} else {
			p.lggr.Errorw("failed to validate gateway vault request before authorization", "method", req.Method, "requestID", req.ID, "error", err)
		}
		return nil, err
	}

	p.lggr.Debugw("authorizing gateway vault request", "method", req.Method, "requestID", req.ID)
	// Authorize on the pre-stamp request bytes. Authorizer also runs replay-guard and owner checks.
	authResult, err := p.authorizer.AuthorizeRequest(ctx, *req)
	if err != nil {
		authErr := fmt.Errorf("request not authorized: %w", err)
		p.lggr.Errorw("gateway vault request authorization failed", "method", req.Method, "requestID", req.ID, "hasAuth", req.Auth != "", "incomingOwner", incomingOwner, "error", authErr)
		return nil, authErr
	}

	originalRequestID := req.ID
	authorizedOwner := authResult.AuthorizedOwner()
	req.ID = authorizedOwner + vaulttypes.RequestIDSeparator + originalRequestID
	if err := p.validator.StampAuthorizedParams(req, req.ID); err != nil {
		p.lggr.Errorw("failed to stamp authorized request params", "method", req.Method, "requestID", req.ID, "error", err)
		return nil, fmt.Errorf("failed to stamp authorized request params: %w", err)
	}

	p.lggr.Debugw("authorized gateway vault request", "method", req.Method, "requestID", req.ID, "owner", authorizedOwner, "orgID", authResult.OrgID(), "workflowOwner", authResult.WorkflowOwner())
	return &AuthorizedGatewayVaultRequest{
		Req:        *req,
		AuthResult: authResult,
	}, nil
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
