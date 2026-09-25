package vault_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
)

func TestGatewayVaultRequestProcessor_ProcessRequest_GatewayModePreservesDigest(t *testing.T) {
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
		return digestErr == nil && gotDigest == digestBefore && got.ID == "req-1"
	})).Return(vault.NewAuthResult("org-test", owner, digestBefore, 0), nil)

	processor := mustNewGatewayVaultRequestProcessor(t, validator, authorizer, false)
	authorized, err := processor.ProcessRequest(t.Context(), &req, nil)
	require.NoError(t, err)
	require.Equal(t, owner+vaulttypes.RequestIDSeparator+"req-1", authorized.Req.ID)
	require.Equal(t, vaulttypes.DefaultNamespace, mustListNamespace(t, authorized.Req.Params))
}

func TestGatewayVaultRequestProcessor_ProcessRequest_NodeReauthPreservesDigest(t *testing.T) {
	t.Parallel()

	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	owner := "0xabc"
	originalRequestID := "req-1"
	prefixedRequestID := owner + vaulttypes.RequestIDSeparator + originalRequestID
	params := gatewaySecretsMethodParamsForStripPrefixDigest(t, vaulttypes.MethodSecretsList, owner, prefixedRequestID)

	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      prefixedRequestID,
		Method:  vaulttypes.MethodSecretsList,
		Params:  params,
	}

	authorizer := vaultcapmocks.NewAuthorizer(t)
	authorizer.EXPECT().AuthorizeRequest(t.Context(), mock.MatchedBy(func(got jsonrpc.Request[json.RawMessage]) bool {
		return got.ID == originalRequestID
	})).Return(vault.NewAuthResult("", owner, "digest", 0), nil)

	processor := mustNewGatewayVaultRequestProcessor(t, validator, authorizer, true)
	authorized, err := processor.ProcessRequest(t.Context(), &req, nil)
	require.NoError(t, err)
	require.Equal(t, prefixedRequestID, authorized.Req.ID)
}

func TestGatewayVaultRequestProcessor_ProcessRequest_NodeModeRestoresEnvelopeIDOnError(t *testing.T) {
	t.Parallel()

	owner := "0xabc"
	originalRequestID := "req-1"
	prefixedRequestID := owner + vaulttypes.RequestIDSeparator + originalRequestID

	t.Run("authorization failure", func(t *testing.T) {
		t.Parallel()

		validator := mustNewTestRequestValidator(t)
		params := gatewaySecretsMethodParamsForStripPrefixDigest(t, vaulttypes.MethodSecretsList, owner, prefixedRequestID)

		req := jsonrpc.Request[json.RawMessage]{
			ID:     prefixedRequestID,
			Method: vaulttypes.MethodSecretsList,
			Params: params,
		}

		authorizer := vaultcapmocks.NewAuthorizer(t)
		authorizer.EXPECT().AuthorizeRequest(t.Context(), mock.MatchedBy(func(got jsonrpc.Request[json.RawMessage]) bool {
			return got.ID == originalRequestID
		})).Return(nil, errors.New("not allowlisted"))

		processor := mustNewGatewayVaultRequestProcessor(t, validator, authorizer, true)
		_, err := processor.ProcessRequest(t.Context(), &req, nil)
		require.ErrorContains(t, err, "request not authorized")
		require.Equal(t, prefixedRequestID, req.ID)
	})

	t.Run("pre-auth validation failure", func(t *testing.T) {
		t.Parallel()

		validator := mustNewTestRequestValidator(t)
		raw := json.RawMessage(`{"request_id":"req-1","ids":[]}`)

		req := jsonrpc.Request[json.RawMessage]{
			ID:     prefixedRequestID,
			Method: vaulttypes.MethodSecretsDelete,
			Params: &raw,
		}

		authorizer := vaultcapmocks.NewAuthorizer(t)

		processor := mustNewGatewayVaultRequestProcessor(t, validator, authorizer, true)
		_, err := processor.ProcessRequest(t.Context(), &req, nil)
		require.ErrorContains(t, err, "request batch must contain at least 1 item")
		require.True(t, vault.IsInvalidVaultParamsError(err))
		require.Equal(t, prefixedRequestID, req.ID)
	})
}

func mustListNamespace(t *testing.T, params *json.RawMessage) string {
	t.Helper()
	var parsed vaultcommon.ListSecretIdentifiersRequest
	require.NoError(t, json.Unmarshal(*params, &parsed))
	return parsed.Namespace
}

// TestGatewayVaultRequestProcessor_ProcessRequest_RejectsOversizedBlobPayload confirms
// the CRIT-1 ingress fix at the earliest gate: a create batch whose queued representation
// (StoredPendingQueueItem blob payload) exceeds VaultMaxBlobPayloadSizeLimit is rejected
// before authorization and before fanning out to the vault DON. EncryptedValue is
// hex-encoded on the wire (2x its decoded size), which is what pushes a full batch of
// cap-size ciphertexts over the limit.
func TestGatewayVaultRequestProcessor_ProcessRequest_RejectsOversizedBlobPayload(t *testing.T) {
	t.Parallel()

	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	// 10 secrets, each 2KB decoded = 4KB hex on the wire: ~40KB queued representation,
	// over the 25.6KB blob cap, while every per-secret limit (2KB decoded) passes.
	oversizedValue := strings.Repeat("ab", 2048)
	secrets := make([]*vaultcommon.EncryptedSecret, 10)
	for i := range secrets {
		secrets[i] = &vaultcommon.EncryptedSecret{
			Id:             &vaultcommon.SecretIdentifier{Owner: "0xabc", Key: fmt.Sprintf("k%d", i)},
			EncryptedValue: oversizedValue,
		}
	}

	for _, mode := range []struct {
		name             string
		stripOwnerPrefix bool
	}{
		{name: "gateway mode", stripOwnerPrefix: false},
		{name: "node mode", stripOwnerPrefix: true},
	} {
		t.Run(mode.name, func(t *testing.T) {
			t.Parallel()
			// No authorizer expectation: the blob-size check runs pre-auth, so an
			// AuthorizeRequest call here fails the mock.
			authorizer := vaultcapmocks.NewAuthorizer(t)
			processor := mustNewGatewayVaultRequestProcessor(t, validator, authorizer, mode.stripOwnerPrefix)

			req := mustWriteRequest(t, vaulttypes.MethodSecretsCreate, secrets)
			_, err := processor.ProcessRequest(t.Context(), &req, nil)
			require.Error(t, err)
			require.True(t, vault.IsInvalidVaultParamsError(err))
			require.ErrorContains(t, err, "request exceeds maximum pending queue blob payload size")
		})
	}
}
