package vault_test

import (
	"context"
	"encoding/json"
	"errors"
	"maps"
	"slices"
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
		RequestId: "req-1",
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

// digestRecordingAuthorizer authorizes every request and records the digest it was asked to authorize.
type digestRecordingAuthorizer struct {
	owner   string
	digests []string
}

func (a *digestRecordingAuthorizer) AuthorizeRequest(_ context.Context, req jsonrpc.Request[json.RawMessage]) (*vault.AuthResult, error) {
	digest, err := req.Digest()
	if err != nil {
		return nil, err
	}
	a.digests = append(a.digests, digest)
	return vault.NewAuthResult("", a.owner, digest, 0), nil
}

// Each Vault node re-derives the caller's params from the gateway-forwarded request and authorizes
// that digest. Any params the gateway accepts must therefore give the nodes the digest the gateway
// authorized; anything else must be rejected at the gateway before authorization.
func TestGatewayVaultRequestProcessor_ProcessRequest_GatewayAndNodesAuthorizeSameDigest(t *testing.T) {
	t.Parallel()

	const owner = "0x00000000000000000000000000000000000000aa"
	secretID := `{"key":"k","owner":"` + owner + `","namespace":"main"}`
	secretIDNoNamespace := `{"key":"k","owner":"` + owner + `"}`

	// Params fields other than request_id, with and without every namespace set.
	bodies := map[string]struct{ withNamespace, withoutNamespace string }{
		vaulttypes.MethodSecretsCreate: {
			`"encrypted_secrets":[{"id":` + secretID + `,"encrypted_value":"ab"}]`,
			`"encrypted_secrets":[{"id":` + secretIDNoNamespace + `,"encrypted_value":"ab"}]`,
		},
		vaulttypes.MethodSecretsUpdate: {
			`"encrypted_secrets":[{"id":` + secretID + `,"encrypted_value":"ab"}]`,
			`"encrypted_secrets":[{"id":` + secretIDNoNamespace + `,"encrypted_value":"ab"}]`,
		},
		vaulttypes.MethodSecretsDelete: {
			`"ids":[` + secretID + `]`,
			`"ids":[` + secretIDNoNamespace + `]`,
		},
		vaulttypes.MethodSecretsList: {
			`"owner":"` + owner + `","namespace":"main"`,
			`"owner":"` + owner + `"`,
		},
	}
	require.ElementsMatch(t, vaulttypes.GatewaySecretsMethods, slices.Collect(maps.Keys(bodies)))

	for method, body := range bodies {
		variants := []struct {
			name     string
			params   string
			accepted bool
		}{
			{"canonical", `{"request_id":"req-1",` + body.withNamespace + `}`, true},
			{"reordered keys and whitespace", `{ ` + body.withNamespace + ` , "request_id" : "req-1" }`, true},
			{"extra field", `{"request_id":"req-1",` + body.withNamespace + `,"x":1}`, false},
			{"missing namespace", `{"request_id":"req-1",` + body.withoutNamespace + `}`, false},
			{"missing request_id", `{` + body.withNamespace + `}`, false},
			{"request_id other than the request ID", `{"request_id":"req-2",` + body.withNamespace + `}`, false},
			{"differently cased key", `{"REQUEST_ID":"req-1",` + body.withNamespace + `}`, false},
		}
		for _, v := range variants {
			t.Run(method+"/"+v.name, func(t *testing.T) {
				t.Parallel()

				raw := json.RawMessage(v.params)
				req := jsonrpc.Request[json.RawMessage]{
					Version: jsonrpc.JsonRpcVersion,
					ID:      "req-1",
					Method:  method,
					Params:  &raw,
				}

				gatewayAuth := &digestRecordingAuthorizer{owner: owner}
				gateway := mustNewGatewayVaultRequestProcessor(t, mustNewTestRequestValidator(t), gatewayAuth, false)
				authorized, err := gateway.ProcessRequest(t.Context(), &req, nil)
				if !v.accepted {
					require.ErrorContains(t, err, "params must be in canonical form")
					require.True(t, vault.IsInvalidVaultParamsError(err))
					require.Empty(t, gatewayAuth.digests, "gateway must reject before authorizing")
					return
				}
				require.NoError(t, err)

				forwarded := authorized.Req
				nodeAuth := &digestRecordingAuthorizer{owner: owner}
				node := mustNewGatewayVaultRequestProcessor(t, mustNewTestRequestValidator(t), nodeAuth, true)
				_, err = node.ProcessRequest(t.Context(), &forwarded, nil)
				require.NoError(t, err)

				require.Len(t, gatewayAuth.digests, 1)
				require.Equal(t, gatewayAuth.digests, nodeAuth.digests)
			})
		}
	}
}

func mustListNamespace(t *testing.T, params *json.RawMessage) string {
	t.Helper()
	var parsed vaultcommon.ListSecretIdentifiersRequest
	require.NoError(t, json.Unmarshal(*params, &parsed))
	return parsed.Namespace
}
