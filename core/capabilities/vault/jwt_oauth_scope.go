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
	// ErrVaultJWTScopeDenied is returned when the token's scopes do not authorize the requested Vault method.
	ErrVaultJWTScopeDenied = errors.New("JWT OAuth scope does not authorize this Vault method")
)

// OAuth scopes issued by the CRE authorization server for Vault secret operations.
// These must stay aligned with cre-platform-graphql permission→scope mapping.
const (
	OAuthScopeVaultSecretsCreate = "create:secrets"
	OAuthScopeVaultSecretsUpdate = "update:secrets"
	OAuthScopeVaultSecretsDelete = "delete:secrets"
	OAuthScopeVaultSecretsList   = "list:secrets"
)

var vaultMethodOAuthScopes = map[string][]string{
	vaulttypes.MethodSecretsCreate: {OAuthScopeVaultSecretsCreate},
	vaulttypes.MethodSecretsUpdate: {OAuthScopeVaultSecretsUpdate},
	vaulttypes.MethodSecretsDelete: {OAuthScopeVaultSecretsDelete},
	vaulttypes.MethodSecretsList:   {OAuthScopeVaultSecretsList},
}

// OAuthScopeForVaultRPCMethod returns the OAuth scope required to authorize the given
// Vault JSON-RPC method over the JWT path.
func OAuthScopeForVaultRPCMethod(method string) (string, error) {
	scopes, ok := vaultMethodOAuthScopes[method]
	if !ok || len(scopes) == 0 {
		return "", fmt.Errorf("no OAuth scope mapping for Vault method %q", method)
	}
	return scopes[0], nil
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

func enforceVaultJWTOAuthScopes(method string, tokenScopes []string) error {
	required, ok := vaultMethodOAuthScopes[method]
	if !ok {
		return fmt.Errorf("%w: unsupported Vault JSON-RPC method %q", ErrVaultJWTScopeDenied, method)
	}
	if len(tokenScopes) == 0 {
		return ErrMissingVaultOAuthScope
	}
	for _, need := range required {
		for _, granted := range tokenScopes {
			if strings.EqualFold(strings.TrimSpace(granted), need) {
				return nil
			}
		}
	}
	return fmt.Errorf("%w: method %q requires scope %q", ErrVaultJWTScopeDenied, method, required[0])
}
