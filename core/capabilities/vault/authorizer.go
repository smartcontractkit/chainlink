package vault

import (
	"context"
	"encoding/json"
	"errors"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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

// AuthorizedOwner returns the canonical owner to use for request scoping.
func (a *AuthResult) AuthorizedOwner() string {
	if a == nil {
		return ""
	}
	if a.orgID != "" {
		return a.orgID
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
