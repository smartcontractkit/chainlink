package vault_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vault "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func TestRequestValidator_PrepareUserJSONRPCRequest_RejectsNilParams(t *testing.T) {
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-1",
		Method: vaulttypes.MethodSecretsCreate,
	}
	err = validator.PrepareUserJSONRPCRequest(t.Context(), &req, vault.UserJSONRPCValidationOptions{
		SkipLabelValidation: true,
	}, false)
	require.Error(t, err)
	require.True(t, vault.IsInvalidVaultParamsError(err))
}

func TestRequestValidator_PrepareUserJSONRPCRequest_LeavesParamsUnchanged(t *testing.T) {
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	params, err := json.Marshal(vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:   "secret",
					Owner: "0xabc",
				},
				EncryptedValue: "ab",
			},
		},
	})
	require.NoError(t, err)
	raw := json.RawMessage(params)
	originalParams := string(raw)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-1",
		Method: vaulttypes.MethodSecretsCreate,
		Params: &raw,
	}
	err = validator.PrepareUserJSONRPCRequest(t.Context(), &req, vault.UserJSONRPCValidationOptions{
		SkipLabelValidation: true,
	}, false)
	require.NoError(t, err)
	require.Equal(t, originalParams, string(*req.Params))
}

func TestRequestValidator_FinalizeAuthorizedJSONRPCRequest_NormalizesCreateRequest(t *testing.T) {
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	params, err := json.Marshal(vaultcommon.CreateSecretsRequest{
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:   "secret",
					Owner: "0xabc",
				},
				EncryptedValue: "ab",
			},
		},
	})
	require.NoError(t, err)
	raw := json.RawMessage(params)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-1",
		Method: vaulttypes.MethodSecretsCreate,
		Params: &raw,
	}
	err = validator.PrepareUserJSONRPCRequest(t.Context(), &req, vault.UserJSONRPCValidationOptions{
		SkipLabelValidation: true,
	}, false)
	require.NoError(t, err)

	prefixedRequestID := "0xabc" + vaulttypes.RequestIDSeparator + "req-1"
	require.NoError(t, validator.FinalizeAuthorizedJSONRPCRequest(&req, prefixedRequestID))

	var parsed vaultcommon.CreateSecretsRequest
	require.NoError(t, json.Unmarshal(*req.Params, &parsed))
	require.Equal(t, prefixedRequestID, parsed.RequestId)
	require.Equal(t, vaulttypes.DefaultNamespace, parsed.EncryptedSecrets[0].Id.Namespace)
	require.Equal(t, prefixedRequestID, req.ID)
}

func TestRequestValidator_PrepareWithStripOwnerPrefixPreservesJWTDigest(t *testing.T) {
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	owner, err := vault.DeriveJWTAuthorizedVaultWorkflowOwner("org-test", 1, "")
	require.NoError(t, err)

	originalRequestID := "req-1"
	params, err := json.Marshal(vaultcommon.ListSecretIdentifiersRequest{
		Namespace: "main",
		Owner:     owner,
		RequestId: originalRequestID,
	})
	require.NoError(t, err)
	raw := json.RawMessage(params)

	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      originalRequestID,
		Method:  vaulttypes.MethodSecretsList,
		Params:  &raw,
	}
	err = validator.PrepareUserJSONRPCRequest(t.Context(), &req, vault.UserJSONRPCValidationOptions{
		SkipLabelValidation: true,
	}, false)
	require.NoError(t, err)

	jwtDigest, err := req.Digest()
	require.NoError(t, err)

	prefixedRequestID := owner + vaulttypes.RequestIDSeparator + originalRequestID
	require.NoError(t, validator.FinalizeAuthorizedJSONRPCRequest(&req, prefixedRequestID))

	err = validator.PrepareUserJSONRPCRequest(t.Context(), &req, vault.UserJSONRPCValidationOptions{
		SkipLabelValidation: true,
	}, true)
	require.NoError(t, err)

	digestAfterStrip, err := req.Digest()
	require.NoError(t, err)
	require.Equal(t, jwtDigest, digestAfterStrip)
	require.Equal(t, originalRequestID, req.ID)
}

func TestRequestValidator_PrepareThenAuthorizePreservesJWTDigest(t *testing.T) {
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	owner, err := vault.DeriveJWTAuthorizedVaultWorkflowOwner("org-test", 1, "")
	require.NoError(t, err)

	params, err := json.Marshal(vaultcommon.ListSecretIdentifiersRequest{
		Namespace: "main",
		Owner:     owner,
	})
	require.NoError(t, err)
	raw := json.RawMessage(params)

	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-1",
		Method:  vaulttypes.MethodSecretsList,
		Params:  &raw,
	}
	digestBefore, err := req.Digest()
	require.NoError(t, err)

	err = validator.PrepareUserJSONRPCRequest(t.Context(), &req, vault.UserJSONRPCValidationOptions{
		SkipLabelValidation: true,
	}, false)
	require.NoError(t, err)

	digestAfterPrepare, err := req.Digest()
	require.NoError(t, err)
	require.Equal(t, digestBefore, digestAfterPrepare)
}

func TestRequestValidator_PrepareUserJSONRPCRequest_CoversAllUserSecretsMethods(t *testing.T) {
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	owner := "0xabc"
	for _, method := range vaulttypes.UserSecretsMethods {
		t.Run(method, func(t *testing.T) {
			params := userSecretsMethodParamsForValidation(t, method, owner)
			req := jsonrpc.Request[json.RawMessage]{
				ID:     "req-1",
				Method: method,
				Params: params,
			}
			err := validator.PrepareUserJSONRPCRequest(t.Context(), &req, vault.UserJSONRPCValidationOptions{
				SkipLabelValidation: true,
			}, false)
			require.NoError(t, err)
			require.NotNil(t, req.Params)
		})
	}
}

func userSecretsMethodParamsForValidation(t *testing.T, method, owner string) *json.RawMessage {
	t.Helper()

	var payload any
	switch method {
	case vaulttypes.MethodSecretsCreate:
		payload = vaultcommon.CreateSecretsRequest{
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{Id: &vaultcommon.SecretIdentifier{Owner: owner, Key: "k"}, EncryptedValue: "ab"},
			},
		}
	case vaulttypes.MethodSecretsUpdate:
		payload = vaultcommon.UpdateSecretsRequest{
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{Id: &vaultcommon.SecretIdentifier{Owner: owner, Key: "k"}, EncryptedValue: "ab"},
			},
		}
	case vaulttypes.MethodSecretsDelete:
		payload = vaultcommon.DeleteSecretsRequest{
			Ids: []*vaultcommon.SecretIdentifier{{Owner: owner, Key: "k"}},
		}
	case vaulttypes.MethodSecretsList:
		payload = vaultcommon.ListSecretIdentifiersRequest{Owner: owner}
	default:
		t.Fatalf("add validation test params for %s", method)
	}

	raw, err := json.Marshal(payload)
	require.NoError(t, err)
	return (*json.RawMessage)(&raw)
}
