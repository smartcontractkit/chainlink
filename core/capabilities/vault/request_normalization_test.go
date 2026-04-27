package vault

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func TestNormalizeRequestWithIdentity_RoundTripsCreateRequestIdentity(t *testing.T) {
	req := mustVaultJSONRPCRequest(t, vaulttypes.MethodSecretsCreate, &vaultcommon.CreateSecretsRequest{
		RequestId:     "req-1",
		OrgId:         "",
		WorkflowOwner: "",
	})

	withIdentity, err := NormalizeRequestWithIdentity(req, "org-123", "0xowner")
	require.NoError(t, err)

	parsedWithIdentity := mustParseCreateSecretsRequest(t, withIdentity)
	require.Equal(t, "org-123", parsedWithIdentity.OrgId)
	require.Equal(t, "0xowner", parsedWithIdentity.WorkflowOwner)

	stripped, err := StripRequestIdentity(withIdentity)
	require.NoError(t, err)

	parsedStripped := mustParseCreateSecretsRequest(t, stripped)
	require.Empty(t, parsedStripped.OrgId)
	require.Empty(t, parsedStripped.WorkflowOwner)
	require.Equal(t, "req-1", parsedStripped.RequestId)
}

func TestStripRequestIdentity_LeavesMalformedParamsUnchanged(t *testing.T) {
	raw := json.RawMessage(`{"not":"valid"`)
	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-1",
		Method:  vaulttypes.MethodSecretsCreate,
		Params:  &raw,
	}

	stripped, err := StripRequestIdentity(req)
	require.NoError(t, err)
	require.Equal(t, string(raw), string(*stripped.Params))
}

func mustVaultJSONRPCRequest(t *testing.T, method string, payload any) jsonrpc.Request[json.RawMessage] {
	t.Helper()

	b, err := json.Marshal(payload)
	require.NoError(t, err)

	raw := json.RawMessage(b)
	return jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-1",
		Method:  method,
		Params:  &raw,
	}
}

func mustParseCreateSecretsRequest(t *testing.T, req jsonrpc.Request[json.RawMessage]) *vaultcommon.CreateSecretsRequest {
	t.Helper()

	parsed := &vaultcommon.CreateSecretsRequest{}
	require.NotNil(t, req.Params)
	require.NoError(t, json.Unmarshal(*req.Params, parsed))
	return parsed
}
