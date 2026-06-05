package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

// AuthResult is the normalized authorization output shared by
// AllowListBasedAuth and JWTBasedAuth.
type AuthResult struct {
	orgID         string
	workflowOwner string
	digest        string
	expiresAt     int64
}

// NewAuthResult remains exported for cross-package tests that cannot construct
// AuthResult directly because its fields are intentionally private.
func NewAuthResult(orgID, workflowOwner, digest string, expiresAt int64) *AuthResult {
	return &AuthResult{
		orgID:         orgID,
		workflowOwner: workflowOwner,
		digest:        digest,
		expiresAt:     expiresAt,
	}
}

// OrgID returns the authorized org ID, if present.
func (a *AuthResult) OrgID() string {
	if a == nil {
		return ""
	}
	return a.orgID
}

// WorkflowOwner returns the authorized workflow owner, if present.
func (a *AuthResult) WorkflowOwner() string {
	if a == nil {
		return ""
	}
	return a.workflowOwner
}

// AuthorizedOwner returns the canonical workflow-owner address used for Vault request ID prefixing
// and secret ownership (JWT-derived owner or allowlisted workflow owner).
func (a *AuthResult) AuthorizedOwner() string {
	if a == nil {
		return ""
	}
	return a.workflowOwner
}

// Digest returns the request digest used for replay protection.
func (a *AuthResult) Digest() string {
	if a == nil {
		return ""
	}
	return a.digest
}

// ExpiresAt returns the unix timestamp (UTC) after which this
// authorization is no longer valid.
func (a *AuthResult) ExpiresAt() int64 {
	if a == nil {
		return 0
	}
	return a.expiresAt
}

// Authorizer selects the applicable auth mechanism for a Vault request.
type Authorizer interface {
	AuthorizeRequest(ctx context.Context, req jsonrpc.Request[json.RawMessage]) (*AuthResult, error)
}

type authorizer struct {
	allowListBasedAuth Authorizer
	jwtBasedAuth       Authorizer
	replayGuard        *RequestReplayGuard
	lggr               logger.Logger
}

func NewAuthorizer(allowListBasedAuth Authorizer, jwtBasedAuth Authorizer, lggr logger.Logger) Authorizer {
	return &authorizer{
		allowListBasedAuth: allowListBasedAuth,
		jwtBasedAuth:       jwtBasedAuth,
		replayGuard:        NewRequestReplayGuard(),
		lggr:               logger.Named(lggr, "VaultAuthorizer"),
	}
}

func (a *authorizer) AuthorizeRequest(ctx context.Context, req jsonrpc.Request[json.RawMessage]) (*AuthResult, error) {
	authResult, err := a.authorizeRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if authResult == nil {
		err = errors.New("auth mechanism returned nil auth result")
		a.lggr.Errorw("auth mechanism returned nil auth result", "method", req.Method, "requestID", req.ID, "hasAuth", req.Auth != "")
		return nil, err
	}
	if err := a.replayGuard.CheckAndRecord(authResult.Digest(), authResult.ExpiresAt()); err != nil {
		a.lggr.Debugw("replay guard rejected request", "method", req.Method, "requestID", req.ID, "owner", authResult.AuthorizedOwner(), "digest", authResult.Digest(), "expiresAt", authResult.ExpiresAt(), "hasAuth", req.Auth != "", "error", err)
		return nil, err
	}
	if ownerErr := validatePreparedVaultOwners(req, authResult.AuthorizedOwner()); ownerErr != nil {
		a.lggr.Errorw("owner binding rejected request", "method", req.Method, "requestID", req.ID, "owner", authResult.AuthorizedOwner(), "hasAuth", req.Auth != "", "error", ownerErr)
		return nil, ownerErr
	}
	a.lggr.Debugw("request authorized", "method", req.Method, "requestID", req.ID, "owner", authResult.AuthorizedOwner(), "digest", authResult.Digest(), "expiresAt", authResult.ExpiresAt(), "hasAuth", req.Auth != "")
	return authResult, nil
}

func (a *authorizer) authorizeRequest(ctx context.Context, req jsonrpc.Request[json.RawMessage]) (*AuthResult, error) {
	// Requests without req.Auth continue using the allowlist-based path for backwards compatibility.
	// Existing clients do not populate the auth field yet, so treating an empty value as JWT would break them.
	if req.Auth == "" {
		return a.authorizeAllowListBasedAuth(ctx, req)
	}
	return a.authorizeJWTBasedAuth(ctx, req)
}

func (a *authorizer) authorizeAllowListBasedAuth(ctx context.Context, req jsonrpc.Request[json.RawMessage]) (*AuthResult, error) {
	if a.allowListBasedAuth == nil {
		err := errors.New("AllowListBasedAuth authorizer is nil")
		a.lggr.Errorw("AllowListBasedAuth unavailable", "method", req.Method, "requestID", req.ID, "error", err)
		return nil, err
	}
	return a.allowListBasedAuth.AuthorizeRequest(ctx, req)
}

func (a *authorizer) authorizeJWTBasedAuth(ctx context.Context, req jsonrpc.Request[json.RawMessage]) (*AuthResult, error) {
	if a.jwtBasedAuth == nil {
		err := errors.New("JWTBasedAuth is nil")
		a.lggr.Errorw("JWTBasedAuth unavailable", "method", req.Method, "requestID", req.ID, "error", err)
		return nil, err
	}
	return a.jwtBasedAuth.AuthorizeRequest(ctx, req)
}

// validatePreparedVaultOwners checks that every secret identifier in the request
// payload belongs to workflowOwner. It runs after allowlist or JWT authentication
// so neither path can be exploited to mutate another owner's secrets.
func validatePreparedVaultOwners(req jsonrpc.Request[json.RawMessage], workflowOwner string) error {
	switch req.Method {
	case vaulttypes.MethodSecretsCreate:
		parsed := &vaultcommon.CreateSecretsRequest{}
		if err := json.Unmarshal(*req.Params, parsed); err != nil {
			return fmt.Errorf("failed to parse create secrets request: %w", err)
		}
		return validateEncryptedSecretOwnersMatchWorkflowOwner(parsed.EncryptedSecrets, workflowOwner)
	case vaulttypes.MethodSecretsUpdate:
		parsed := &vaultcommon.UpdateSecretsRequest{}
		if err := json.Unmarshal(*req.Params, parsed); err != nil {
			return fmt.Errorf("failed to parse update secrets request: %w", err)
		}
		return validateEncryptedSecretOwnersMatchWorkflowOwner(parsed.EncryptedSecrets, workflowOwner)
	case vaulttypes.MethodSecretsDelete:
		parsed := &vaultcommon.DeleteSecretsRequest{}
		if err := json.Unmarshal(*req.Params, parsed); err != nil {
			return fmt.Errorf("failed to parse delete secrets request: %w", err)
		}
		return validateSecretIdentifierOwnersMatchWorkflowOwner(parsed.Ids, workflowOwner)
	case vaulttypes.MethodSecretsList:
		parsed := &vaultcommon.ListSecretIdentifiersRequest{}
		if err := json.Unmarshal(*req.Params, parsed); err != nil {
			return fmt.Errorf("failed to parse list secrets request: %w", err)
		}
		if normalizeOwner(parsed.Owner) != normalizeOwner(workflowOwner) {
			return fmt.Errorf("list secrets owner %q does not match authorized workflow owner %q", parsed.Owner, workflowOwner)
		}
		return nil
	default:
		return nil
	}
}

func validateEncryptedSecretOwnersMatchWorkflowOwner(encryptedSecrets []*vaultcommon.EncryptedSecret, workflowOwner string) error {
	if len(encryptedSecrets) == 0 {
		return errors.New("encrypted secrets must contain at least one identifier")
	}
	for idx, encryptedSecret := range encryptedSecrets {
		if encryptedSecret == nil || encryptedSecret.Id == nil {
			return fmt.Errorf("encrypted secret at index %d must include an identifier", idx)
		}
		if normalizeOwner(encryptedSecret.Id.Owner) != normalizeOwner(workflowOwner) {
			return fmt.Errorf("encrypted secret owner at index %d %q does not match authorized workflow owner %q", idx, encryptedSecret.Id.Owner, workflowOwner)
		}
	}
	return nil
}

func validateSecretIdentifierOwnersMatchWorkflowOwner(ids []*vaultcommon.SecretIdentifier, workflowOwner string) error {
	if len(ids) == 0 {
		return errors.New("secret identifiers must not be empty")
	}
	for idx, id := range ids {
		if id == nil {
			return fmt.Errorf("secret identifier at index %d must not be nil", idx)
		}
		if normalizeOwner(id.Owner) != normalizeOwner(workflowOwner) {
			return fmt.Errorf("secret identifier owner at index %d %q does not match authorized workflow owner %q", idx, id.Owner, workflowOwner)
		}
	}
	return nil
}
