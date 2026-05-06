package vault

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

var (
	// ErrMissingVaultOAuthScope is returned when a Vault JWT carries no OAuth scope
	// or permissions usable for Vault JSON-RPC authorization.
	ErrMissingVaultOAuthScope = errors.New("missing OAuth scope for Vault JWT authorization")
	// ErrVaultJWTMultipleOAuthScopes is returned when more than one Vault secret scope
	// is present (after filtering to known Vault OAuth scopes). Authorization requires
	// a single unambiguous Vault scope.
	ErrVaultJWTMultipleOAuthScopes = errors.New("vault JWT must carry exactly one Vault secret OAuth scope")
	// ErrVaultJWTScopeDenied is returned when the token's scopes do not authorize the requested Vault method.
	ErrVaultJWTScopeDenied = errors.New("jwt OAuth scope does not authorize this Vault method")
)

// OAuth scopes issued by the CRE authorization server for Vault secret operations.
// These must stay aligned with cre-platform-graphql permission→scope mapping.
const (
	OAuthScopeVaultSecretsCreate = "create:secrets"
	OAuthScopeVaultSecretsUpdate = "update:secrets"
	OAuthScopeVaultSecretsDelete = "delete:secrets"
	OAuthScopeVaultSecretsList   = "list:secrets"
)

var vaultMethodOAuthScopes = map[string]string{
	vaulttypes.MethodSecretsCreate: OAuthScopeVaultSecretsCreate,
	vaulttypes.MethodSecretsUpdate: OAuthScopeVaultSecretsUpdate,
	vaulttypes.MethodSecretsDelete: OAuthScopeVaultSecretsDelete,
	vaulttypes.MethodSecretsList:   OAuthScopeVaultSecretsList,
}

// canonicalVaultOAuthScopes lists every Vault secret scope issued for JWT authorization.
var canonicalVaultOAuthScopes = []string{
	OAuthScopeVaultSecretsCreate,
	OAuthScopeVaultSecretsUpdate,
	OAuthScopeVaultSecretsDelete,
	OAuthScopeVaultSecretsList,
}

// filterToCanonicalVaultOAuthScopes returns a deduplicated list of claims that match one
// of the known Vault secret OAuth scopes (case-insensitive). Other scopes (e.g. openid)
// are ignored so typical OAuth access tokens still work.
func filterToCanonicalVaultOAuthScopes(scopes []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, s := range scopes {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		for _, canon := range canonicalVaultOAuthScopes {
			if !strings.EqualFold(s, canon) {
				continue
			}
			key := strings.ToLower(canon)
			if _, dup := seen[key]; dup {
				break
			}
			seen[key] = struct{}{}
			out = append(out, canon)
			break
		}
	}
	return out
}

// OAuthScopeForVaultRPCMethod returns the OAuth scope required to authorize the given
// Vault JSON-RPC method over the JWT path.
func OAuthScopeForVaultRPCMethod(method string) (string, error) {
	scope, ok := vaultMethodOAuthScopes[method]
	if !ok || scope == "" {
		return "", fmt.Errorf("no OAuth scope mapping for Vault method %q", method)
	}
	return scope, nil
}

func extractOAuthScopesFromClaims(claims jwt.MapClaims) []string {
	seen := map[string]struct{}{}
	var out []string
	add := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		key := strings.ToLower(s)
		if _, dup := seen[key]; dup {
			return
		}
		seen[key] = struct{}{}
		out = append(out, s)
	}

	if raw, ok := claims["scope"]; ok {
		switch v := raw.(type) {
		case string:
			for _, part := range strings.Fields(v) {
				add(part)
			}
		case []interface{}:
			for _, item := range v {
				if s, ok := item.(string); ok {
					add(s)
				}
			}
		case []string:
			for _, s := range v {
				add(s)
			}
		}
	}

	// Auth0 API Authorization often emits custom permissions as a string array claim.
	if raw, ok := claims["permissions"]; ok {
		switch v := raw.(type) {
		case []interface{}:
			for _, item := range v {
				if s, ok := item.(string); ok {
					add(s)
				}
			}
		case []string:
			for _, s := range v {
				add(s)
			}
		}
	}

	return out
}

// enforceVaultJWTOAuthScopes ensures the token carries exactly one known Vault secret
// OAuth scope (after collecting scope and permissions claims) and that it matches the
// JSON-RPC method. Non-Vault scopes in the same claims (e.g. openid) are ignored.
func enforceVaultJWTOAuthScopes(method string, tokenScopes []string) error {
	expected, err := OAuthScopeForVaultRPCMethod(method)
	if err != nil {
		return fmt.Errorf("%w: unsupported Vault JSON-RPC method %q", ErrVaultJWTScopeDenied, method)
	}

	vaultScopes := filterToCanonicalVaultOAuthScopes(tokenScopes)
	switch len(vaultScopes) {
	case 0:
		return ErrMissingVaultOAuthScope
	case 1:
		if strings.EqualFold(vaultScopes[0], expected) {
			return nil
		}
		return fmt.Errorf("%w: method %q requires scope %q, got %q", ErrVaultJWTScopeDenied, method, expected, vaultScopes[0])
	default:
		return fmt.Errorf("%w: found %d Vault secret scopes, want exactly one", ErrVaultJWTMultipleOAuthScopes, len(vaultScopes))
	}
}
