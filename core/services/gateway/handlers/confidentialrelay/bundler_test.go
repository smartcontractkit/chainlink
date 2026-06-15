package confidentialrelay

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialrelay"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// chainlink-common's confidentialrelay.Validate rejects request params missing any
// field the canonical hash binds to (Owner must be a 0x-prefixed 20-byte hex address;
// ExecutionID must be 32-byte hex with no prefix). Test params satisfy these formats.
const (
	testOwner       = "0x0000000000000000000000000000000000000001"
	testExecutionID = "0000000000000000000000000000000000000000000000000000000000000001"
	testEnclavePK   = "aabbcc"
)

func validCapParams(workflowID string) relaytypes.CapabilityRequestParams {
	return relaytypes.CapabilityRequestParams{
		WorkflowID:   workflowID,
		Owner:        testOwner,
		ExecutionID:  testExecutionID,
		ReferenceID:  "ref-1",
		CapabilityID: "cap-1",
		Payload:      "in",
	}
}

func validSecretsParams(workflowID string) relaytypes.SecretsRequestParams {
	return relaytypes.SecretsRequestParams{
		WorkflowID:       workflowID,
		Owner:            testOwner,
		ExecutionID:      testExecutionID,
		Secrets:          []relaytypes.SecretIdentifier{{Key: "k1", Namespace: "ns"}},
		EnclavePublicKey: testEnclavePK,
	}
}

func capExecRequest(t *testing.T, id string, params relaytypes.CapabilityRequestParams) jsonrpc.Request[json.RawMessage] {
	t.Helper()
	raw, err := json.Marshal(params)
	require.NoError(t, err)
	rm := json.RawMessage(raw)
	return jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      id,
		Method:  MethodCapabilityExec,
		Params:  &rm,
	}
}

func capExecSignedResponse(t *testing.T, id string, result relaytypes.CapabilityResponseResult, signer []byte) jsonrpc.Response[json.RawMessage] {
	t.Helper()
	signed := relaytypes.SignedCapabilityResponseResult{
		Result: result,
		Signature: relaytypes.RelayResponseSignature{
			Signer:    signer,
			Signature: append([]byte("sig-"), signer...),
		},
	}
	raw, err := json.Marshal(signed)
	require.NoError(t, err)
	rm := json.RawMessage(raw)
	return jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      id,
		Method:  MethodCapabilityExec,
		Result:  &rm,
	}
}

func decodeCapabilityBundle(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) relaytypes.SignedCapabilityResponseBundle {
	t.Helper()
	require.NotNil(t, resp)
	require.NotNil(t, resp.Result)
	var out relaytypes.SignedCapabilityResponseBundle
	require.NoError(t, json.Unmarshal(*resp.Result, &out))
	return out
}

// The gateway is a dumb fan-in: it forwards every collected per-node response in
// one bundle, untouched, without verifying signatures or counting signers.
func TestBundler_capabilityExec_forwardsAllResponses(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	b := &bundler{}

	req := capExecRequest(t, "req-1", validCapParams("wf-1"))
	resultA := relaytypes.CapabilityResponseResult{Payload: "out-A"}
	resultB := relaytypes.CapabilityResponseResult{Payload: "out-B"}

	// Divergent results are NOT resolved by the gateway; both are forwarded for
	// the enclave to verify and pick from.
	resps := map[string]jsonrpc.Response[json.RawMessage]{
		"n0": capExecSignedResponse(t, "req-1", resultA, []byte("signer-2")),
		"n1": capExecSignedResponse(t, "req-1", resultA, []byte("signer-0")),
		"n2": capExecSignedResponse(t, "req-1", resultB, []byte("signer-1")),
	}

	out, count, err := b.Bundle(req, resps, lggr)
	require.NoError(t, err)
	require.Equal(t, 3, count)
	bundle := decodeCapabilityBundle(t, out)
	require.Len(t, bundle.Responses, 3, "every collected response is forwarded; order is not significant")
}

// A single response carries exactly one signature (the contract no longer permits
// an array), so a single node cannot stuff multiple signer identities. The bundler
// makes no quorum decision regardless of how few or many responses it holds.
func TestBundler_capabilityExec_singleResponseIsOneEntry(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	b := &bundler{}

	req := capExecRequest(t, "req-solo", validCapParams("wf-1"))
	resps := map[string]jsonrpc.Response[json.RawMessage]{
		"n0": capExecSignedResponse(t, "req-solo", relaytypes.CapabilityResponseResult{Payload: "x"}, []byte("signer-0")),
	}

	out, count, err := b.Bundle(req, resps, lggr)
	require.NoError(t, err)
	require.Equal(t, 1, count)
	bundle := decodeCapabilityBundle(t, out)
	require.Len(t, bundle.Responses, 1)
}

func TestBundler_skipsTransportErrorsAndUndecodable(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	b := &bundler{}

	req := capExecRequest(t, "req-err", validCapParams("wf-1"))
	result := relaytypes.CapabilityResponseResult{Payload: "out-A"}

	errResp := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-err",
		Method:  MethodCapabilityExec,
		Error:   &jsonrpc.WireError{Code: -32000, Message: "node failure"},
	}
	garbage := json.RawMessage(`{"not":"a signed response"`)
	garbageResp := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-err",
		Method:  MethodCapabilityExec,
		Result:  &garbage,
	}

	resps := map[string]jsonrpc.Response[json.RawMessage]{
		"n0": errResp,
		"n1": garbageResp,
		"n2": capExecSignedResponse(t, "req-err", result, []byte("signer-2")),
		"n3": capExecSignedResponse(t, "req-err", result, []byte("signer-3")),
	}

	out, count, err := b.Bundle(req, resps, lggr)
	require.NoError(t, err)
	require.Equal(t, 2, count, "transport errors and undecodable results are dropped from the bundle")
	bundle := decodeCapabilityBundle(t, out)
	require.Len(t, bundle.Responses, 2)
}

func TestBundler_emptyResponses(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	b := &bundler{}

	req := capExecRequest(t, "req-empty", validCapParams("wf-1"))
	out, count, err := b.Bundle(req, map[string]jsonrpc.Response[json.RawMessage]{}, lggr)
	require.NoError(t, err)
	require.Equal(t, 0, count)
	bundle := decodeCapabilityBundle(t, out)
	require.Empty(t, bundle.Responses)
}

func TestBundler_secretsGet_forwardsAllResponses(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	b := &bundler{}

	params := validSecretsParams("wf-1")
	rawParams, err := json.Marshal(params)
	require.NoError(t, err)
	rm := json.RawMessage(rawParams)
	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-secrets",
		Method:  MethodSecretsGet,
		Params:  &rm,
	}

	result := relaytypes.SecretsResponseResult{
		Secrets: []relaytypes.SecretEntry{{
			ID:              relaytypes.SecretIdentifier{Key: "k1", Namespace: "ns"},
			Ciphertext:      "ct",
			EncryptedShares: []string{"sh-0", "sh-1"},
		}},
	}
	signed := func(signer []byte) jsonrpc.Response[json.RawMessage] {
		raw, mErr := json.Marshal(relaytypes.SignedSecretsResponseResult{
			Result: result,
			Signature: relaytypes.RelayResponseSignature{
				Signer:    signer,
				Signature: append([]byte("sig-"), signer...),
			},
		})
		require.NoError(t, mErr)
		rmsg := json.RawMessage(raw)
		return jsonrpc.Response[json.RawMessage]{
			Version: jsonrpc.JsonRpcVersion,
			ID:      "req-secrets",
			Method:  MethodSecretsGet,
			Result:  &rmsg,
		}
	}

	resps := map[string]jsonrpc.Response[json.RawMessage]{
		"n0": signed([]byte("signer-0")),
		"n1": signed([]byte("signer-1")),
		"n2": signed([]byte("signer-2")),
	}

	out, count, err := b.Bundle(req, resps, lggr)
	require.NoError(t, err)
	require.Equal(t, 3, count)

	var bundle relaytypes.SignedSecretsResponseBundle
	require.NoError(t, json.Unmarshal(*out.Result, &bundle))
	require.Len(t, bundle.Responses, 3)
	require.Equal(t, result, bundle.Responses[0].Result)
}

func TestBundler_unknownMethod(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	b := &bundler{}

	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-x",
		Method:  "confidential.unknown",
	}
	out, _, err := b.Bundle(req, map[string]jsonrpc.Response[json.RawMessage]{}, lggr)
	require.Nil(t, out)
	require.ErrorIs(t, err, errUnknownMethod)
}
