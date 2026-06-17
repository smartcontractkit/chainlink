package vault_test

import (
	"encoding/json"
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

func TestGatewayVaultRequestPipeline_ProcessGatewayVaultRequest_GatewayModePreservesDigest(t *testing.T) {
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

	pipeline := vault.NewGatewayVaultRequestPipeline(validator, authorizer, logger.TestLogger(t))
	authorized, err := pipeline.ProcessGatewayVaultRequest(t.Context(), &req, vault.GatewayVaultRequestPipelineOptions{
		StripOwnerPrefixForAuth: false,
		SkipLabelValidation:     true,
	})
	require.NoError(t, err)
	require.Equal(t, owner+vaulttypes.RequestIDSeparator+"req-1", authorized.Req.ID)
	require.Equal(t, vaulttypes.DefaultNamespace, mustListNamespace(t, authorized.Req.Params))
}

func TestGatewayVaultRequestPipeline_ProcessGatewayVaultRequest_NodeReauthPreservesDigest(t *testing.T) {
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

	pipeline := vault.NewGatewayVaultRequestPipeline(validator, authorizer, logger.TestLogger(t))
	authorized, err := pipeline.ProcessGatewayVaultRequest(t.Context(), &req, vault.GatewayVaultRequestPipelineOptions{
		StripOwnerPrefixForAuth: true,
		SkipLabelValidation:     true,
	})
	require.NoError(t, err)
	require.Equal(t, prefixedRequestID, authorized.Req.ID)
}

func mustListNamespace(t *testing.T, params *json.RawMessage) string {
	t.Helper()
	var parsed vaultcommon.ListSecretIdentifiersRequest
	require.NoError(t, json.Unmarshal(*params, &parsed))
	return parsed.Namespace
}
