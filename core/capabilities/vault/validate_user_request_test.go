package vault_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vault "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	vaultcapmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestGatewayVaultRequestProcessor_ProcessRequest_RejectsNilParams(t *testing.T) {
	t.Parallel()

	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	req := jsonrpc.Request[json.RawMessage]{
		ID:     "req-1",
		Method: vaulttypes.MethodSecretsCreate,
	}

	authorizer := vaultcapmocks.NewAuthorizer(t)
	processor := mustNewGatewayVaultRequestProcessor(t, validator, authorizer, false)
	err = processRequestErr(processor, t, &req)
	require.Error(t, err)
	require.True(t, vault.IsInvalidVaultParamsError(err))
}

func TestGatewayVaultRequestProcessor_ProcessRequest_LeavesParamsUnchangedBeforeAuth(t *testing.T) {
	t.Parallel()

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

	authorizer := vaultcapmocks.NewAuthorizer(t)
	authorizer.EXPECT().AuthorizeRequest(t.Context(), mock.MatchedBy(func(got jsonrpc.Request[json.RawMessage]) bool {
		return got.Params != nil && string(*got.Params) == originalParams
	})).Return(vault.NewAuthResult("", "0xabc", "digest", 0), nil)

	processor := mustNewGatewayVaultRequestProcessor(t, validator, authorizer, false)
	_, err = processor.ProcessRequest(t.Context(), &req, nil)
	require.NoError(t, err)
}

func TestGatewayVaultRequestProcessor_ProcessRequest_NormalizesAllSecretsMethods(t *testing.T) {
	t.Parallel()

	owner := "0xabc"
	prefixedRequestID := owner + vaulttypes.RequestIDSeparator + "req-1"

	for _, method := range vaulttypes.GatewaySecretsMethods {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			validator := mustNewTestRequestValidator(t)
			params := gatewaySecretsMethodParamsForValidation(t, method, owner)
			req := jsonrpc.Request[json.RawMessage]{
				ID:     "req-1",
				Method: method,
				Params: params,
			}

			authorizer := vaultcapmocks.NewAuthorizer(t)
			authorizer.EXPECT().AuthorizeRequest(t.Context(), mock.Anything).Return(vault.NewAuthResult("", owner, "digest", 0), nil)

			processor := mustNewGatewayVaultRequestProcessor(t, validator, authorizer, false)
			authorized, err := processor.ProcessRequest(t.Context(), &req, nil)
			require.NoError(t, err)
			require.Equal(t, prefixedRequestID, authorized.Req.ID)

			switch method {
			case vaulttypes.MethodSecretsCreate:
				var parsed vaultcommon.CreateSecretsRequest
				require.NoError(t, json.Unmarshal(*authorized.Req.Params, &parsed))
				require.Equal(t, prefixedRequestID, parsed.RequestId)
				require.Equal(t, vaulttypes.DefaultNamespace, parsed.EncryptedSecrets[0].Id.Namespace)
			case vaulttypes.MethodSecretsUpdate:
				var parsed vaultcommon.UpdateSecretsRequest
				require.NoError(t, json.Unmarshal(*authorized.Req.Params, &parsed))
				require.Equal(t, prefixedRequestID, parsed.RequestId)
				require.Equal(t, vaulttypes.DefaultNamespace, parsed.EncryptedSecrets[0].Id.Namespace)
			case vaulttypes.MethodSecretsDelete:
				var parsed vaultcommon.DeleteSecretsRequest
				require.NoError(t, json.Unmarshal(*authorized.Req.Params, &parsed))
				require.Equal(t, prefixedRequestID, parsed.RequestId)
				require.Equal(t, vaulttypes.DefaultNamespace, parsed.Ids[0].Namespace)
			case vaulttypes.MethodSecretsList:
				var parsed vaultcommon.ListSecretIdentifiersRequest
				require.NoError(t, json.Unmarshal(*authorized.Req.Params, &parsed))
				require.Equal(t, prefixedRequestID, parsed.RequestId)
				require.Equal(t, vaulttypes.DefaultNamespace, parsed.Namespace)
			}
		})
	}
}

func TestGatewayVaultRequestProcessor_ProcessRequest_StripOwnerPrefixPreservesJWTDigest_AllSecretsMethods(t *testing.T) {
	t.Parallel()

	owner, err := vault.DeriveJWTAuthorizedVaultWorkflowOwner("org-test", 1, "")
	require.NoError(t, err)

	for _, method := range vaulttypes.GatewaySecretsMethods {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			validator := mustNewTestRequestValidator(t)

			originalRequestID := "req-1"
			params := gatewaySecretsMethodParamsForStripPrefixDigest(t, method, owner, originalRequestID)
			req := jsonrpc.Request[json.RawMessage]{
				Version: jsonrpc.JsonRpcVersion,
				ID:      originalRequestID,
				Method:  method,
				Params:  params,
			}

			jwtDigest, err := req.Digest()
			require.NoError(t, err)

			prefixedRequestID := owner + vaulttypes.RequestIDSeparator + originalRequestID
			prefixedParams := gatewaySecretsMethodParamsForStripPrefixDigest(t, method, owner, prefixedRequestID)
			req.ID = prefixedRequestID
			req.Params = prefixedParams

			authorizer := vaultcapmocks.NewAuthorizer(t)
			authorizer.EXPECT().AuthorizeRequest(t.Context(), mock.MatchedBy(func(got jsonrpc.Request[json.RawMessage]) bool {
				gotDigest, digestErr := got.Digest()
				return digestErr == nil && gotDigest == jwtDigest && got.ID == originalRequestID
			})).Return(vault.NewAuthResult("org-test", owner, jwtDigest, 0), nil)

			processor := mustNewGatewayVaultRequestProcessor(t, validator, authorizer, true)
			authorized, err := processor.ProcessRequest(t.Context(), &req, nil)
			require.NoError(t, err)
			require.Equal(t, prefixedRequestID, authorized.Req.ID)
		})
	}
}

func TestRequestValidator_ValidateCreateSecretsRequest_AcceptsBatchAtLimit(t *testing.T) {
	t.Parallel()

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

	err = validator.ValidateCreateSecretsRequest(t.Context(), nil, &vaultcommon.CreateSecretsRequest{
		RequestId:        "req-create-limit",
		EncryptedSecrets: encryptedSecrets,
	}, true)
	require.NoError(t, err)
}

func TestRequestValidator_ValidateUpdateSecretsRequest_AcceptsBatchAtLimit(t *testing.T) {
	t.Parallel()

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

	err = validator.ValidateUpdateSecretsRequest(t.Context(), nil, &vaultcommon.UpdateSecretsRequest{
		RequestId:        "req-update-limit",
		EncryptedSecrets: encryptedSecrets,
	}, true)
	require.NoError(t, err)
}

func TestRequestValidator_ValidateCreateSecretsRequest_RejectsBatchAboveLimit(t *testing.T) {
	t.Parallel()

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

	err = validator.ValidateCreateSecretsRequest(t.Context(), nil, &vaultcommon.CreateSecretsRequest{
		RequestId:        "req-create-over-limit",
		EncryptedSecrets: encryptedSecrets,
	}, true)
	require.Error(t, err)
	require.ErrorContains(t, err, fmt.Sprintf("request batch size exceeds maximum of %d", vaulttypes.MaxBatchSize))
}

func TestRequestValidator_ValidateUpdateSecretsRequest_RejectsBatchAboveLimit(t *testing.T) {
	t.Parallel()

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

	err = validator.ValidateUpdateSecretsRequest(t.Context(), nil, &vaultcommon.UpdateSecretsRequest{
		RequestId:        "req-update-over-limit",
		EncryptedSecrets: encryptedSecrets,
	}, true)
	require.Error(t, err)
	require.ErrorContains(t, err, fmt.Sprintf("request batch size exceeds maximum of %d", vaulttypes.MaxBatchSize))
}

func TestRequestValidator_ValidateDeleteSecretsRequest_AcceptsBatchAtLimit(t *testing.T) {
	t.Parallel()

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

	err = validator.ValidateDeleteSecretsRequest(t.Context(), &vaultcommon.DeleteSecretsRequest{
		RequestId: "req-delete-limit",
		Ids:       ids,
	})
	require.NoError(t, err)
}

func TestRequestValidator_ValidateDeleteSecretsRequest_RejectsBatchAboveLimit(t *testing.T) {
	t.Parallel()

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

	err = validator.ValidateDeleteSecretsRequest(t.Context(), &vaultcommon.DeleteSecretsRequest{
		RequestId: "req-delete-over-limit",
		Ids:       ids,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, fmt.Sprintf("request batch size exceeds maximum of %d", vaulttypes.MaxBatchSize))
}

func TestGatewayVaultRequestProcessor_ProcessRequest_InvalidParamsParityWithStripPrefix(t *testing.T) {
	t.Parallel()

	for _, stripOwnerPrefix := range []bool{false, true} {
		t.Run(fmt.Sprintf("stripOwnerPrefix=%t", stripOwnerPrefix), func(t *testing.T) {
			t.Parallel()

			validator := mustNewTestRequestValidator(t)
			authorizer := vaultcapmocks.NewAuthorizer(t)
			processor := mustNewGatewayVaultRequestProcessor(t, validator, authorizer, stripOwnerPrefix)

			req := jsonrpc.Request[json.RawMessage]{
				ID:     "req-1",
				Method: vaulttypes.MethodSecretsCreate,
			}
			err := processRequestErr(processor, t, &req)
			require.Error(t, err)
			require.True(t, vault.IsInvalidVaultParamsError(err))
			require.ErrorContains(t, err, "request params must not be nil")
		})
	}
}

func TestGatewayVaultRequestProcessor_ProcessRequest_PreservesJWTDigestBeforeAuth(t *testing.T) {
	t.Parallel()

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

	authorizer := vaultcapmocks.NewAuthorizer(t)
	authorizer.EXPECT().AuthorizeRequest(t.Context(), mock.MatchedBy(func(got jsonrpc.Request[json.RawMessage]) bool {
		gotDigest, digestErr := got.Digest()
		return digestErr == nil && gotDigest == digestBefore
	})).Return(vault.NewAuthResult("org-test", owner, digestBefore, 0), nil)

	processor := mustNewGatewayVaultRequestProcessor(t, validator, authorizer, false)
	_, err = processor.ProcessRequest(t.Context(), &req, nil)
	require.NoError(t, err)
}

func TestGatewayVaultRequestProcessor_ProcessRequest_CoversAllGatewaySecretsMethods(t *testing.T) {
	t.Parallel()

	owner := "0xabc"
	for _, method := range vaulttypes.GatewaySecretsMethods {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			validator := mustNewTestRequestValidator(t)
			params := gatewaySecretsMethodParamsForValidation(t, method, owner)
			req := jsonrpc.Request[json.RawMessage]{
				ID:     "req-1",
				Method: method,
				Params: params,
			}

			authorizer := vaultcapmocks.NewAuthorizer(t)
			authorizer.EXPECT().AuthorizeRequest(t.Context(), mock.Anything).Return(vault.NewAuthResult("", owner, "digest", 0), nil)

			processor := mustNewGatewayVaultRequestProcessor(t, validator, authorizer, false)
			authorized, err := processor.ProcessRequest(t.Context(), &req, nil)
			require.NoError(t, err)
			require.NotNil(t, authorized.Req.Params)
		})
	}
}

func mustNewGatewayVaultRequestProcessor(t *testing.T, validator *vault.RequestValidator, authorizer vault.Authorizer, stripOwnerPrefix bool) *vault.GatewayVaultRequestProcessor {
	t.Helper()
	processor, err := vault.NewGatewayVaultRequestProcessor(validator, authorizer, stripOwnerPrefix, logger.TestLogger(t))
	require.NoError(t, err)
	return processor
}

func processRequestErr(processor *vault.GatewayVaultRequestProcessor, t *testing.T, req *jsonrpc.Request[json.RawMessage]) error {
	t.Helper()
	_, err := processor.ProcessRequest(t.Context(), req, nil)
	return err
}

func mustNewTestRequestValidator(t *testing.T) *vault.RequestValidator {
	t.Helper()
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)
	return validator
}

func gatewaySecretsMethodParamsForStripPrefixDigest(t *testing.T, method, owner, requestID string) *json.RawMessage {
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

func gatewaySecretsMethodParamsForValidation(t *testing.T, method, owner string) *json.RawMessage {
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
