package vault

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func TestBindVaultOwners_RejectsUnregisteredMethods(t *testing.T) {
	params := json.RawMessage(`{}`)
	err := bindVaultOwners(jsonrpc.Request[json.RawMessage]{
		Method: "vault.unsupported",
		Params: &params,
	}, "0xauthorized")
	require.ErrorContains(t, err, "owner binding not implemented for method \"vault.unsupported\"")
}

func TestBindVaultOwners_RejectsEmptyBatch(t *testing.T) {
	owner := "0xauthorized"
	for _, tc := range []struct {
		method string
		params *json.RawMessage
	}{
		{
			method: vaulttypes.MethodSecretsCreate,
			params: mustMarshalParams(t, vaultcommon.CreateSecretsRequest{EncryptedSecrets: []*vaultcommon.EncryptedSecret{}}),
		},
		{
			method: vaulttypes.MethodSecretsUpdate,
			params: mustMarshalParams(t, vaultcommon.UpdateSecretsRequest{EncryptedSecrets: []*vaultcommon.EncryptedSecret{}}),
		},
		{
			method: vaulttypes.MethodSecretsDelete,
			params: mustMarshalParams(t, vaultcommon.DeleteSecretsRequest{Ids: []*vaultcommon.SecretIdentifier{}}),
		},
	} {
		t.Run(tc.method, func(t *testing.T) {
			err := bindVaultOwners(jsonrpc.Request[json.RawMessage]{
				Method: tc.method,
				Params: tc.params,
			}, owner)
			require.ErrorContains(t, err, "request batch must contain at least 1 item")
		})
	}
}

func TestBindVaultOwners_RejectsMalformedBatches(t *testing.T) {
	owner := "0xauthorized"
	tests := []struct {
		name        string
		method      string
		params      *json.RawMessage
		errContains string
	}{
		{
			name:   "create rejects nil encrypted secret",
			method: vaulttypes.MethodSecretsCreate,
			params: mustMarshalParams(t, vaultcommon.CreateSecretsRequest{
				EncryptedSecrets: []*vaultcommon.EncryptedSecret{nil},
			}),
			errContains: "encrypted secret must not be nil at index 0",
		},
		{
			name:   "create rejects nil secret id",
			method: vaulttypes.MethodSecretsCreate,
			params: mustMarshalParams(t, vaultcommon.CreateSecretsRequest{
				EncryptedSecrets: []*vaultcommon.EncryptedSecret{
					{Id: nil, EncryptedValue: "ab"},
				},
			}),
			errContains: "secret ID must not be nil at index 0",
		},
		{
			name:   "update rejects nil encrypted secret",
			method: vaulttypes.MethodSecretsUpdate,
			params: mustMarshalParams(t, vaultcommon.UpdateSecretsRequest{
				EncryptedSecrets: []*vaultcommon.EncryptedSecret{nil},
			}),
			errContains: "encrypted secret must not be nil at index 0",
		},
		{
			name:   "delete rejects nil secret identifier",
			method: vaulttypes.MethodSecretsDelete,
			params: mustMarshalParams(t, vaultcommon.DeleteSecretsRequest{
				Ids: []*vaultcommon.SecretIdentifier{nil},
			}),
			errContains: "secret ID must not be nil at index 0",
		},
		{
			name:   "create rejects owner mismatch in batch",
			method: vaulttypes.MethodSecretsCreate,
			params: mustMarshalParams(t, vaultcommon.CreateSecretsRequest{
				EncryptedSecrets: []*vaultcommon.EncryptedSecret{
					{Id: &vaultcommon.SecretIdentifier{Owner: owner, Namespace: "ns", Key: "k0"}, EncryptedValue: "ab"},
					{Id: &vaultcommon.SecretIdentifier{Owner: "0xother", Namespace: "ns", Key: "k1"}, EncryptedValue: "ab"},
				},
			}),
			errContains: "encrypted secret owner at index 1",
		},
		{
			name:   "delete rejects owner mismatch in batch",
			method: vaulttypes.MethodSecretsDelete,
			params: mustMarshalParams(t, vaultcommon.DeleteSecretsRequest{
				Ids: []*vaultcommon.SecretIdentifier{
					{Owner: owner, Namespace: "ns", Key: "k0"},
					{Owner: "0xother", Namespace: "ns", Key: "k1"},
				},
			}),
			errContains: "secret identifier owner at index 1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := bindVaultOwners(jsonrpc.Request[json.RawMessage]{
				Method: tc.method,
				Params: tc.params,
			}, owner)
			require.ErrorContains(t, err, tc.errContains)
		})
	}
}

func TestBindVaultOwners_CoversAllUserSecretsMethods(t *testing.T) {
	owner := "0xauthorized"
	for _, method := range vaulttypes.UserSecretsMethods {
		t.Run(method, func(t *testing.T) {
			params := userSecretsMethodParamsForOwnerBinding(t, method, owner)
			err := bindVaultOwners(jsonrpc.Request[json.RawMessage]{
				Method: method,
				Params: params,
			}, owner)
			require.NoError(t, err)
		})
	}
}

func userSecretsMethodParamsForOwnerBinding(t *testing.T, method, owner string) *json.RawMessage {
	t.Helper()

	var payload any
	switch method {
	case vaulttypes.MethodSecretsCreate:
		payload = vaultcommon.CreateSecretsRequest{
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{Id: &vaultcommon.SecretIdentifier{Owner: owner, Namespace: "ns", Key: "k"}, EncryptedValue: "ab"},
			},
		}
	case vaulttypes.MethodSecretsUpdate:
		payload = vaultcommon.UpdateSecretsRequest{
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{Id: &vaultcommon.SecretIdentifier{Owner: owner, Namespace: "ns", Key: "k"}, EncryptedValue: "ab"},
			},
		}
	case vaulttypes.MethodSecretsDelete:
		payload = vaultcommon.DeleteSecretsRequest{
			Ids: []*vaultcommon.SecretIdentifier{{Owner: owner, Namespace: "ns", Key: "k"}},
		}
	case vaulttypes.MethodSecretsList:
		payload = vaultcommon.ListSecretIdentifiersRequest{Owner: owner, Namespace: "ns"}
	default:
		t.Fatalf("add owner-binding test params for %s", method)
	}

	raw, err := json.Marshal(payload)
	require.NoError(t, err)
	return (*json.RawMessage)(&raw)
}

func mustMarshalParams(t *testing.T, payload any) *json.RawMessage {
	t.Helper()
	raw, err := json.Marshal(payload)
	require.NoError(t, err)
	return (*json.RawMessage)(&raw)
}
