package vault_test

import (
	"encoding/json"
	"fmt"
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

func TestRequestValidator_FinalizeAuthorizedJSONRPCRequest_NormalizesAllSecretsMethods(t *testing.T) {
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	owner := "0xabc"
	prefixedRequestID := owner + vaulttypes.RequestIDSeparator + "req-1"

	for _, method := range vaulttypes.UserSecretsMethods {
		t.Run(method, func(t *testing.T) {
			params := userSecretsMethodParamsForValidation(t, method, owner)
			req := jsonrpc.Request[json.RawMessage]{
				ID:     "req-1",
				Method: method,
				Params: params,
			}
			err = validator.PrepareUserJSONRPCRequest(t.Context(), &req, vault.UserJSONRPCValidationOptions{
				SkipLabelValidation: true,
			}, false)
			require.NoError(t, err)

			require.NoError(t, validator.FinalizeAuthorizedJSONRPCRequest(&req, prefixedRequestID))
			require.Equal(t, prefixedRequestID, req.ID)

			switch method {
			case vaulttypes.MethodSecretsCreate:
				var parsed vaultcommon.CreateSecretsRequest
				require.NoError(t, json.Unmarshal(*req.Params, &parsed))
				require.Equal(t, prefixedRequestID, parsed.RequestId)
				require.Equal(t, vaulttypes.DefaultNamespace, parsed.EncryptedSecrets[0].Id.Namespace)
			case vaulttypes.MethodSecretsUpdate:
				var parsed vaultcommon.UpdateSecretsRequest
				require.NoError(t, json.Unmarshal(*req.Params, &parsed))
				require.Equal(t, prefixedRequestID, parsed.RequestId)
				require.Equal(t, vaulttypes.DefaultNamespace, parsed.EncryptedSecrets[0].Id.Namespace)
			case vaulttypes.MethodSecretsDelete:
				var parsed vaultcommon.DeleteSecretsRequest
				require.NoError(t, json.Unmarshal(*req.Params, &parsed))
				require.Equal(t, prefixedRequestID, parsed.RequestId)
				require.Equal(t, vaulttypes.DefaultNamespace, parsed.Ids[0].Namespace)
			case vaulttypes.MethodSecretsList:
				var parsed vaultcommon.ListSecretIdentifiersRequest
				require.NoError(t, json.Unmarshal(*req.Params, &parsed))
				require.Equal(t, prefixedRequestID, parsed.RequestId)
				require.Equal(t, vaulttypes.DefaultNamespace, parsed.Namespace)
			}
		})
	}
}

func TestRequestValidator_PrepareWithStripOwnerPrefixPreservesJWTDigest_AllSecretsMethods(t *testing.T) {
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	owner, err := vault.DeriveJWTAuthorizedVaultWorkflowOwner("org-test", 1, "")
	require.NoError(t, err)

	for _, method := range vaulttypes.UserSecretsMethods {
		t.Run(method, func(t *testing.T) {
			originalRequestID := "req-1"
			params := userSecretsMethodParamsForStripPrefixDigest(t, method, owner, originalRequestID)
			req := jsonrpc.Request[json.RawMessage]{
				Version: jsonrpc.JsonRpcVersion,
				ID:      originalRequestID,
				Method:  method,
				Params:  params,
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
		})
	}
}

func TestRequestValidator_PrepareUserJSONRPCRequest_AcceptsCreateBatchAtLimit(t *testing.T) {
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	owner := "0xabc"
	encryptedSecrets := make([]*vaultcommon.EncryptedSecret, vaulttypes.MaxBatchSize)
	for i := range encryptedSecrets {
		encryptedSecrets[i] = &vaultcommon.EncryptedSecret{
			Id: &vaultcommon.SecretIdentifier{
				Key:   fmt.Sprintf("key%d", i),
				Owner: owner,
			},
			EncryptedValue: "ab",
		}
	}

	params, err := json.Marshal(vaultcommon.CreateSecretsRequest{EncryptedSecrets: encryptedSecrets})
	require.NoError(t, err)
	raw := json.RawMessage(params)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-create-limit",
		Method: vaulttypes.MethodSecretsCreate,
		Params: &raw,
	}
	err = validator.PrepareUserJSONRPCRequest(t.Context(), &req, vault.UserJSONRPCValidationOptions{
		SkipLabelValidation: true,
	}, false)
	require.NoError(t, err)
}

func TestRequestValidator_PrepareUserJSONRPCRequest_AcceptsUpdateBatchAtLimit(t *testing.T) {
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	owner := "0xabc"
	encryptedSecrets := make([]*vaultcommon.EncryptedSecret, vaulttypes.MaxBatchSize)
	for i := range encryptedSecrets {
		encryptedSecrets[i] = &vaultcommon.EncryptedSecret{
			Id: &vaultcommon.SecretIdentifier{
				Key:   fmt.Sprintf("key%d", i),
				Owner: owner,
			},
			EncryptedValue: "ab",
		}
	}

	params, err := json.Marshal(vaultcommon.UpdateSecretsRequest{EncryptedSecrets: encryptedSecrets})
	require.NoError(t, err)
	raw := json.RawMessage(params)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-update-limit",
		Method: vaulttypes.MethodSecretsUpdate,
		Params: &raw,
	}
	err = validator.PrepareUserJSONRPCRequest(t.Context(), &req, vault.UserJSONRPCValidationOptions{
		SkipLabelValidation: true,
	}, false)
	require.NoError(t, err)
}

func TestRequestValidator_PrepareUserJSONRPCRequest_RejectsCreateBatchAboveLimit(t *testing.T) {
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	owner := "0xabc"
	encryptedSecrets := make([]*vaultcommon.EncryptedSecret, vaulttypes.MaxBatchSize+1)
	for i := range encryptedSecrets {
		encryptedSecrets[i] = &vaultcommon.EncryptedSecret{
			Id: &vaultcommon.SecretIdentifier{
				Key:   fmt.Sprintf("key%d", i),
				Owner: owner,
			},
			EncryptedValue: "ab",
		}
	}

	params, err := json.Marshal(vaultcommon.CreateSecretsRequest{EncryptedSecrets: encryptedSecrets})
	require.NoError(t, err)
	raw := json.RawMessage(params)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-create-over-limit",
		Method: vaulttypes.MethodSecretsCreate,
		Params: &raw,
	}
	err = validator.PrepareUserJSONRPCRequest(t.Context(), &req, vault.UserJSONRPCValidationOptions{
		SkipLabelValidation: true,
	}, false)
	require.Error(t, err)
	require.True(t, vault.IsInvalidVaultParamsError(err))
	require.ErrorContains(t, err, fmt.Sprintf("request batch size exceeds maximum of %d", vaulttypes.MaxBatchSize))
}

func TestRequestValidator_PrepareUserJSONRPCRequest_RejectsUpdateBatchAboveLimit(t *testing.T) {
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	owner := "0xabc"
	encryptedSecrets := make([]*vaultcommon.EncryptedSecret, vaulttypes.MaxBatchSize+1)
	for i := range encryptedSecrets {
		encryptedSecrets[i] = &vaultcommon.EncryptedSecret{
			Id: &vaultcommon.SecretIdentifier{
				Key:   fmt.Sprintf("key%d", i),
				Owner: owner,
			},
			EncryptedValue: "ab",
		}
	}

	params, err := json.Marshal(vaultcommon.UpdateSecretsRequest{EncryptedSecrets: encryptedSecrets})
	require.NoError(t, err)
	raw := json.RawMessage(params)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-update-over-limit",
		Method: vaulttypes.MethodSecretsUpdate,
		Params: &raw,
	}
	err = validator.PrepareUserJSONRPCRequest(t.Context(), &req, vault.UserJSONRPCValidationOptions{
		SkipLabelValidation: true,
	}, false)
	require.Error(t, err)
	require.True(t, vault.IsInvalidVaultParamsError(err))
	require.ErrorContains(t, err, fmt.Sprintf("request batch size exceeds maximum of %d", vaulttypes.MaxBatchSize))
}

func TestRequestValidator_PrepareUserJSONRPCRequest_AcceptsDeleteBatchAtLimit(t *testing.T) {
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	owner := "0xabc"
	ids := make([]*vaultcommon.SecretIdentifier, vaulttypes.MaxBatchSize)
	for i := range ids {
		ids[i] = &vaultcommon.SecretIdentifier{
			Key:   fmt.Sprintf("key%d", i),
			Owner: owner,
		}
	}

	params, err := json.Marshal(vaultcommon.DeleteSecretsRequest{Ids: ids})
	require.NoError(t, err)
	raw := json.RawMessage(params)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-delete-limit",
		Method: vaulttypes.MethodSecretsDelete,
		Params: &raw,
	}
	err = validator.PrepareUserJSONRPCRequest(t.Context(), &req, vault.UserJSONRPCValidationOptions{
		SkipLabelValidation: true,
	}, false)
	require.NoError(t, err)
}

func TestRequestValidator_PrepareUserJSONRPCRequest_RejectsDeleteBatchAboveLimit(t *testing.T) {
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	owner := "0xabc"
	ids := make([]*vaultcommon.SecretIdentifier, vaulttypes.MaxBatchSize+1)
	for i := range ids {
		ids[i] = &vaultcommon.SecretIdentifier{
			Key:   fmt.Sprintf("key%d", i),
			Owner: owner,
		}
	}

	params, err := json.Marshal(vaultcommon.DeleteSecretsRequest{Ids: ids})
	require.NoError(t, err)
	raw := json.RawMessage(params)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-delete-over-limit",
		Method: vaulttypes.MethodSecretsDelete,
		Params: &raw,
	}
	err = validator.PrepareUserJSONRPCRequest(t.Context(), &req, vault.UserJSONRPCValidationOptions{
		SkipLabelValidation: true,
	}, false)
	require.Error(t, err)
	require.True(t, vault.IsInvalidVaultParamsError(err))
	require.ErrorContains(t, err, fmt.Sprintf("request batch size exceeds maximum of %d", vaulttypes.MaxBatchSize))
}

func TestRequestValidator_PrepareUserJSONRPCRequest_InvalidParamsParityWithStripPrefix(t *testing.T) {
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	for _, stripOwnerPrefix := range []bool{false, true} {
		t.Run(fmt.Sprintf("stripOwnerPrefix=%t", stripOwnerPrefix), func(t *testing.T) {
			req := jsonrpc.Request[json.RawMessage]{
				ID:     "req-1",
				Method: vaulttypes.MethodSecretsCreate,
			}
			err := validator.PrepareUserJSONRPCRequest(t.Context(), &req, vault.UserJSONRPCValidationOptions{
				SkipLabelValidation: true,
			}, stripOwnerPrefix)
			require.Error(t, err)
			require.True(t, vault.IsInvalidVaultParamsError(err))
			require.ErrorContains(t, err, "request params must not be nil")
		})
	}
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

func userSecretsMethodParamsForStripPrefixDigest(t *testing.T, method, owner, requestID string) *json.RawMessage {
	t.Helper()

	var payload any
	switch method {
	case vaulttypes.MethodSecretsCreate:
		payload = vaultcommon.CreateSecretsRequest{
			RequestId: requestID,
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{Id: &vaultcommon.SecretIdentifier{Owner: owner, Key: "k", Namespace: "main"}, EncryptedValue: "ab"},
			},
		}
	case vaulttypes.MethodSecretsUpdate:
		payload = vaultcommon.UpdateSecretsRequest{
			RequestId: requestID,
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{Id: &vaultcommon.SecretIdentifier{Owner: owner, Key: "k", Namespace: "main"}, EncryptedValue: "ab"},
			},
		}
	case vaulttypes.MethodSecretsDelete:
		payload = vaultcommon.DeleteSecretsRequest{
			RequestId: requestID,
			Ids:       []*vaultcommon.SecretIdentifier{{Owner: owner, Key: "k", Namespace: "main"}},
		}
	case vaulttypes.MethodSecretsList:
		payload = vaultcommon.ListSecretIdentifiersRequest{
			RequestId: requestID,
			Namespace: "main",
			Owner:     owner,
		}
	default:
		t.Fatalf("add strip-prefix digest test params for %s", method)
	}

	raw, err := json.Marshal(payload)
	require.NoError(t, err)
	return (*json.RawMessage)(&raw)
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
