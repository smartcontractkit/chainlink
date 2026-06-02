package confidentialrelay

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialrelay"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// chainlink-common's confidentialrelay.Validate (introduced in chainlink-common#2021)
// rejects request params missing any field the canonical hash binds to and rejects
// malformed structured fields (Owner must be a 0x-prefixed 20-byte hex Ethereum address;
// ExecutionID must be 32-byte hex with no prefix). Hash returns an error on Validate
// failure, so test params that exercise the real aggregator must satisfy these formats.
const (
	testOwner       = "0x0000000000000000000000000000000000000001"
	testExecutionID = "0000000000000000000000000000000000000000000000000000000000000001"
	testEnclavePK   = "aabbcc"
)

// validEnclaveConfig satisfies chainlink-common's confidentialrelay.Validate,
// which requires a non-empty Signers list, non-zero F, and a non-empty
// MasterPublicKey. Without it Hash() errors and the real aggregator never
// reaches quorum.
func validEnclaveConfig() relaytypes.EnclaveConfig {
	return relaytypes.EnclaveConfig{
		Signers:         [][]byte{{0x01}, {0x02}, {0x03}, {0x04}},
		MasterPublicKey: []byte("test-master-public-key"),
		T:               3,
		F:               1,
	}
}

func validEnclaveConfigPtr() *relaytypes.EnclaveConfig {
	c := validEnclaveConfig()
	return &c
}

func validCapParams(workflowID string) relaytypes.CapabilityRequestParams {
	return relaytypes.CapabilityRequestParams{
		WorkflowID:    workflowID,
		Owner:         testOwner,
		ExecutionID:   testExecutionID,
		ReferenceID:   "ref-1",
		CapabilityID:  "cap-1",
		Payload:       "in",
		EnclaveConfig: validEnclaveConfigPtr(),
	}
}

func validSecretsParams(workflowID string) relaytypes.SecretsRequestParams {
	return relaytypes.SecretsRequestParams{
		WorkflowID:       workflowID,
		Owner:            testOwner,
		ExecutionID:      testExecutionID,
		Secrets:          []relaytypes.SecretIdentifier{{Key: "k1", Namespace: "ns"}},
		EnclavePublicKey: testEnclavePK,
		EnclaveConfig:    validEnclaveConfigPtr(),
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
		Signatures: []relaytypes.RelayResponseSignature{{
			Signer:    signer,
			Signature: append([]byte("sig-"), signer...),
		}},
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

func decodeSignedCapability(t *testing.T, resp *jsonrpc.Response[json.RawMessage]) relaytypes.SignedCapabilityResponseResult {
	t.Helper()
	require.NotNil(t, resp)
	require.NotNil(t, resp.Result)
	var out relaytypes.SignedCapabilityResponseResult
	require.NoError(t, json.Unmarshal(*resp.Result, &out))
	return out
}

func TestAggregator_capabilityExec_quorumMergesSignatures(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	agg := &aggregator{}

	params := validCapParams("wf-1")
	req := capExecRequest(t, "req-1", params)
	result := relaytypes.CapabilityResponseResult{Payload: "out-A"}

	resps := map[string]jsonrpc.Response[json.RawMessage]{
		"n0": capExecSignedResponse(t, "req-1", result, []byte("signer-0")),
		"n1": capExecSignedResponse(t, "req-1", result, []byte("signer-1")),
		"n2": capExecSignedResponse(t, "req-1", result, []byte("signer-2")),
		"n3": capExecSignedResponse(t, "req-1", result, []byte("signer-3")),
	}

	out, err := agg.Aggregate(req, resps, 1, 4, lggr)
	require.NoError(t, err)
	got := decodeSignedCapability(t, out)
	require.Equal(t, result, got.Result)
	require.Len(t, got.Signatures, 4)
	for i := 0; i < len(got.Signatures)-1; i++ {
		require.Negative(t, bytes.Compare(got.Signatures[i].Signer, got.Signatures[i+1].Signer),
			"signatures should be sorted by signer for determinism")
	}
}

func TestAggregator_capabilityExec_dedupesDuplicateSigner(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	agg := &aggregator{}

	params := validCapParams("wf-1")
	req := capExecRequest(t, "req-dup", params)
	result := relaytypes.CapabilityResponseResult{Payload: "out-A"}

	// Three responses, two carrying the same signer bytes (defensive dedupe path).
	resps := map[string]jsonrpc.Response[json.RawMessage]{
		"n0": capExecSignedResponse(t, "req-dup", result, []byte("signer-shared")),
		"n1": capExecSignedResponse(t, "req-dup", result, []byte("signer-shared")),
		"n2": capExecSignedResponse(t, "req-dup", result, []byte("signer-2")),
	}

	out, err := agg.Aggregate(req, resps, 1, 3, lggr)
	require.NoError(t, err)
	got := decodeSignedCapability(t, out)
	require.Len(t, got.Signatures, 2, "duplicates by signer must collapse")
}

func TestAggregator_capabilityExec_tiedMajoritiesDeterministic(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	agg := &aggregator{}

	params := validCapParams("wf-1")
	req := capExecRequest(t, "req-tie", params)

	resultA := relaytypes.CapabilityResponseResult{Payload: "out-A"}
	resultB := relaytypes.CapabilityResponseResult{Payload: "out-B"}
	hashA, err := resultA.Hash(params)
	require.NoError(t, err)
	hashB, err := resultB.Hash(params)
	require.NoError(t, err)
	require.NotEqual(t, hashA, hashB)

	wantResult := resultA
	if bytes.Compare(hashB[:], hashA[:]) < 0 {
		wantResult = resultB
	}

	// Each pair has F+1 sigs at F=1; map iteration order must not change the winner.
	for range 100 {
		resps := map[string]jsonrpc.Response[json.RawMessage]{
			"n0": capExecSignedResponse(t, "req-tie", resultA, []byte("signer-0")),
			"n1": capExecSignedResponse(t, "req-tie", resultA, []byte("signer-1")),
			"n2": capExecSignedResponse(t, "req-tie", resultB, []byte("signer-2")),
			"n3": capExecSignedResponse(t, "req-tie", resultB, []byte("signer-3")),
		}
		out, err := agg.Aggregate(req, resps, 1, 4, lggr)
		require.NoError(t, err)
		got := decodeSignedCapability(t, out)
		require.Equal(t, wantResult, got.Result,
			"tied majorities must resolve to the lex-smallest logical hash regardless of map order")
		require.Len(t, got.Signatures, 2)
	}
}

func TestAggregator_skipsTransportErrorsFromQuorum(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	agg := &aggregator{}

	params := validCapParams("wf-1")
	req := capExecRequest(t, "req-err", params)
	result := relaytypes.CapabilityResponseResult{Payload: "out-A"}

	errResp := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-err",
		Method:  MethodCapabilityExec,
		Error:   &jsonrpc.WireError{Code: -32000, Message: "node failure"},
	}

	// 2 transport errors + 2 valid signed responses. F=1 -> F+1=2; signed pair forms quorum.
	resps := map[string]jsonrpc.Response[json.RawMessage]{
		"n0": errResp,
		"n1": errResp,
		"n2": capExecSignedResponse(t, "req-err", result, []byte("signer-2")),
		"n3": capExecSignedResponse(t, "req-err", result, []byte("signer-3")),
	}

	out, err := agg.Aggregate(req, resps, 1, 4, lggr)
	require.NoError(t, err)
	got := decodeSignedCapability(t, out)
	require.Equal(t, result, got.Result)
	require.Len(t, got.Signatures, 2)
}

func TestAggregator_quorumUnobtainable(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	agg := &aggregator{}

	params := validCapParams("wf-1")
	req := capExecRequest(t, "req-uno", params)

	// 4 distinct logical results means no bucket reaches F+1=2; remaining=0; unattainable.
	resps := map[string]jsonrpc.Response[json.RawMessage]{
		"n0": capExecSignedResponse(t, "req-uno", relaytypes.CapabilityResponseResult{Payload: "p-0"}, []byte("s-0")),
		"n1": capExecSignedResponse(t, "req-uno", relaytypes.CapabilityResponseResult{Payload: "p-1"}, []byte("s-1")),
		"n2": capExecSignedResponse(t, "req-uno", relaytypes.CapabilityResponseResult{Payload: "p-2"}, []byte("s-2")),
		"n3": capExecSignedResponse(t, "req-uno", relaytypes.CapabilityResponseResult{Payload: "p-3"}, []byte("s-3")),
	}

	out, err := agg.Aggregate(req, resps, 1, 4, lggr)
	require.Nil(t, out)
	require.ErrorIs(t, err, errQuorumUnobtainable)
}

func TestAggregator_secretsGet_quorumMergesSignatures(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	agg := &aggregator{}

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
			Signatures: []relaytypes.RelayResponseSignature{{
				Signer:    signer,
				Signature: append([]byte("sig-"), signer...),
			}},
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

	out, err := agg.Aggregate(req, resps, 1, 3, lggr)
	require.NoError(t, err)
	require.NotNil(t, out)
	require.NotNil(t, out.Result)

	var got relaytypes.SignedSecretsResponseResult
	require.NoError(t, json.Unmarshal(*out.Result, &got))
	require.Equal(t, result, got.Result)
	require.Len(t, got.Signatures, 3)
}
