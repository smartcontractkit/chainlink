package vault

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"testing"

	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func makeNodes(t *testing.T, signers []string) []capabilities.Node {
	nodes := []capabilities.Node{}
	for idx, s := range signers {
		b, err := hex.DecodeString(s)
		require.NoError(t, err)
		nodes = append(nodes, capabilities.Node{PeerID: &p2ptypes.PeerID{0: uint8(idx)}, Signer: [32]byte(b)}) //nolint:gosec // G115
	}
	return nodes
}

func TestAggregator_Valid_Signatures(t *testing.T) {
	signers := []string{
		"d6da96fe596705b32bc3a0e11cdefad77feaad79000000000000000000000000",
		"327aa349c9718cd36c877d1e90458fe1929768ad000000000000000000000000",
		"e9bf394856d73402b30e160d0e05c847796f0e29000000000000000000000000",
		"efd5bdb6c3256f04489a6ca32654d547297f48b9000000000000000000000000",
	}
	nodes := makeNodes(t, signers)
	mcr := &mockCapabilitiesRegistry{F: 1, Nodes: nodes}
	agg := &baseAggregator{capabilitiesRegistry: mcr}

	ctx, err := hex.DecodeString("000ec4f6a2ba011e909eccf64628855b848e08876a1edd938a1372a9e51adff100000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	sig1, err := hex.DecodeString("d1067844e2849b404d903730c4cae19f090d53a578a1e8dc16ecbdc0285c1f186599108abbe0073b78bc148a6504907474ed3a6881df917e6d142cff70acfb5900")
	require.NoError(t, err)
	sig2, err := hex.DecodeString("c7517c188d297093a6f602046fad7feafe19454ee9dc269b19c8e6c01268037d1f7b423eeecbc495dd2d9a65e106bc3eab849ddfd74a10cbd4ad50c7d953bd4b01")
	require.NoError(t, err)

	rm := json.RawMessage([]byte(`{"responses":[{"error":"failed to verify ciphertext: cannot unmarshal data: unexpected end of JSON input","id":{"key":"W","namespace":"","owner":"foo"},"success":false}]}`))
	sor := vaulttypes.SignedOCRResponse{
		Payload: rm,
		Context: ctx,
		Signatures: [][]byte{
			sig1,
			sig2,
		},
	}
	rawResp, err := json.Marshal(sor)
	require.NoError(t, err)

	currResp := &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsCreate,
		Result:  (*json.RawMessage)(&rawResp),
	}
	responses := map[string]*jsonrpc.Response[json.RawMessage]{
		"a": currResp,
	}
	resp, err := agg.Aggregate(t.Context(), logger.Test(t), responses, currResp)
	require.NoError(t, err)
	assert.Equal(t, currResp, resp)
}

func mustRandom(length int) []byte {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	return randomBytes
}

func newMessage(t *testing.T) *jsonrpc.Response[json.RawMessage] {
	ctx, err := hex.DecodeString("000ec4f6a2ba011e909eccf64628855b848e08876a1edd938a1372a9e51adff100000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)

	rm := json.RawMessage([]byte(`{"responses":[{"error":"failed to verify ciphertext: cannot unmarshal data: unexpected end of JSON input","id":{"key":"W","namespace":"","owner":"foo"},"success":false}]}`))
	sor := vaulttypes.SignedOCRResponse{
		Payload: rm,
		Context: ctx,
		Signatures: [][]byte{
			mustRandom(65),
			mustRandom(65),
		},
	}
	rawResp, err := json.Marshal(sor)
	require.NoError(t, err)

	return &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsGet,
		Result:  (*json.RawMessage)(&rawResp),
	}
}

func TestAggregator_Valid_FallsBackToQuorum(t *testing.T) {
	// No valid signers
	signers := []string{
		hex.EncodeToString(mustRandom(64)),
		hex.EncodeToString(mustRandom(64)),
		hex.EncodeToString(mustRandom(64)),
		hex.EncodeToString(mustRandom(64)),
	}
	nodes := makeNodes(t, signers)
	mcr := &mockCapabilitiesRegistry{F: 1, Nodes: nodes}
	agg := &baseAggregator{capabilitiesRegistry: mcr}

	currResp := &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsGet,
		Result:  (*json.RawMessage)(nil),
		Error: &jsonrpc.WireError{
			Code:    123,
			Message: "some error",
		},
	}
	responses := map[string]*jsonrpc.Response[json.RawMessage]{
		"a": currResp,
		"b": currResp,
		"c": currResp,
	}
	resp, err := agg.Aggregate(t.Context(), logger.Test(t), responses, currResp)
	require.NoError(t, err)
	assert.Equal(t, currResp, resp)
}

func TestAggregator_Valid_FallsBackToQuorum_ExcludesSignaturesInSha(t *testing.T) {
	// No valid signers
	signers := []string{
		hex.EncodeToString(mustRandom(64)),
		hex.EncodeToString(mustRandom(64)),
		hex.EncodeToString(mustRandom(64)),
		hex.EncodeToString(mustRandom(64)),
	}
	nodes := makeNodes(t, signers)
	mcr := &mockCapabilitiesRegistry{F: 1, Nodes: nodes}
	agg := &baseAggregator{capabilitiesRegistry: mcr}

	oldResp1 := newMessage(t)
	oldResp2 := newMessage(t)
	currResp := newMessage(t)
	responses := map[string]*jsonrpc.Response[json.RawMessage]{
		"a": oldResp1,
		"b": oldResp2,
		"c": currResp,
	}
	resp, err := agg.Aggregate(t.Context(), logger.Test(t), responses, currResp)
	require.NoError(t, err)

	respDigests := []string{}
	for _, r := range []*jsonrpc.Response[json.RawMessage]{oldResp1, oldResp2, currResp} {
		dig, ierr := r.Digest()
		require.NoError(t, ierr)
		respDigests = append(respDigests, dig)
	}

	// The response is one of the responses we received.
	digest, err := resp.Digest()
	require.NoError(t, err)
	assert.Contains(t, respDigests, digest)
}

func TestAggregator_InsufficientResponses(t *testing.T) {
	mcr := &mockCapabilitiesRegistry{F: 1}
	agg := &baseAggregator{capabilitiesRegistry: mcr}

	rm := json.RawMessage([]byte(`{}`))
	currResp := &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsGet,
		Result:  &rm,
	}
	responses := map[string]*jsonrpc.Response[json.RawMessage]{
		"a": currResp,
	}
	_, err := agg.Aggregate(t.Context(), logger.Test(t), responses, currResp)
	require.ErrorContains(t, err, "insufficient valid responses to reach quorum")
}

func TestAggregator_QuorumUnobtainable(t *testing.T) {
	// No valid signers
	signers := []string{
		hex.EncodeToString(mustRandom(64)),
		hex.EncodeToString(mustRandom(64)),
		hex.EncodeToString(mustRandom(64)),
		hex.EncodeToString(mustRandom(64)),
	}
	nodes := makeNodes(t, signers)
	mcr := &mockCapabilitiesRegistry{F: 1, Nodes: nodes}
	agg := &baseAggregator{capabilitiesRegistry: mcr}

	rm1 := json.RawMessage([]byte(`{}`))
	resp1 := &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsGet,
		Result:  &rm1,
	}
	rm2 := json.RawMessage([]byte(`{"foo": "bar"}`))
	resp2 := &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsGet,
		Result:  &rm2,
	}
	rm3 := json.RawMessage([]byte(`{"baz": "qux"}`))
	resp3 := &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsGet,
		Result:  &rm3,
	}
	responses := map[string]*jsonrpc.Response[json.RawMessage]{
		"a": resp1,
		"b": resp2,
		"c": resp3,
	}
	_, err := agg.Aggregate(t.Context(), logger.Test(t), responses, resp3)
	require.ErrorContains(t, err, "failed to validate using quorum: quorum unobtainable")
}
