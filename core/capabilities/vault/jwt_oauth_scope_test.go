package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func TestOAuthScopeForVaultRPCMethod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		wantScope  string
		wantErrMsg string
	}{
		{
			name:      "create",
			method:    vaulttypes.MethodSecretsCreate,
			wantScope: OAuthScopeVaultSecretsCreate,
		},
		{
			name:      "update",
			method:    vaulttypes.MethodSecretsUpdate,
			wantScope: OAuthScopeVaultSecretsUpdate,
		},
		{
			name:      "delete",
			method:    vaulttypes.MethodSecretsDelete,
			wantScope: OAuthScopeVaultSecretsDelete,
		},
		{
			name:      "list",
			method:    vaulttypes.MethodSecretsList,
			wantScope: OAuthScopeVaultSecretsList,
		},
		{
			name:       "unknown method",
			method:     "vault.unknown.op",
			wantErrMsg: `no OAuth scope mapping for Vault method "vault.unknown.op"`,
		},
		{
			name:       "secrets get not JWT-scoped in map",
			method:     vaulttypes.MethodSecretsGet,
			wantErrMsg: `no OAuth scope mapping for Vault method "vault.secrets.get"`,
		},
		{
			name:       "public key get not in map",
			method:     vaulttypes.MethodPublicKeyGet,
			wantErrMsg: `no OAuth scope mapping for Vault method "vault.publicKey.get"`,
		},
		{
			name:       "empty method",
			method:     "",
			wantErrMsg: `no OAuth scope mapping for Vault method ""`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := OAuthScopeForVaultRPCMethod(tc.method)
			if tc.wantErrMsg != "" {
				require.Error(t, err)
				assert.Equal(t, tc.wantErrMsg, err.Error())
				assert.Empty(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantScope, got)
		})
	}
}
