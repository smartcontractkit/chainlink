package vault

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func testAggregator(t *testing.T, mcr *mockCapabilitiesRegistry) *baseAggregator {
	t.Helper()
	m, err := newMetrics()
	require.NoError(t, err)
	return &baseAggregator{
		capabilitiesRegistry: mcr,
		metrics:              m,
	}
}

func makeNodes(t *testing.T, signers []string) []capabilities.Node {
	nodes := make([]capabilities.Node, 0, len(signers))
	for idx, s := range signers {
		b, err := hex.DecodeString(s)
		require.NoError(t, err)
		nodes = append(nodes, capabilities.Node{PeerID: &p2ptypes.PeerID{0: uint8(idx)}, Signer: [32]byte(b)})
	}
	return nodes
}

func TestAggregator_Valid_Signatures(t *testing.T) {
	currResp, nodes := makeSignedCreateSecretsResponse(t, "1", 2)
	mcr := &mockCapabilitiesRegistry{F: 1, Nodes: nodes}
	agg := testAggregator(t, mcr)

	responses := map[string]jsonrpc.Response[json.RawMessage]{
		"a": currResp,
	}
	resp, err := agg.Aggregate(t.Context(), logger.Test(t), currResp.ID, responses, &currResp)
	require.NoError(t, err)
	assert.Equal(t, &currResp, resp)
}

func TestAggregator_SignedResponseRequestIDMismatch(t *testing.T) {
	t.Parallel()
	currResp, nodes := makeSignedCreateSecretsResponse(t, "signed-for-other-request", 2)
	mcr := &mockCapabilitiesRegistry{F: 1, Nodes: nodes}
	agg := testAggregator(t, mcr)

	responses := map[string]jsonrpc.Response[json.RawMessage]{
		"a": currResp,
	}

	_, err := agg.Aggregate(t.Context(), logger.Test(t), "expected-request", responses, &currResp)
	require.ErrorContains(t, err, "insufficient valid responses to reach quorum")
}

func TestAggregator_PublicKeyGet_SkipsSignatureValidation(t *testing.T) {
	t.Parallel()
	currResp, nodes := makeSignedCreateSecretsResponse(t, "1", 2)
	currResp.Method = vaulttypes.MethodPublicKeyGet

	mcr := &mockCapabilitiesRegistry{F: 1, Nodes: nodes}
	agg := testAggregator(t, mcr)

	responses := map[string]jsonrpc.Response[json.RawMessage]{
		"a": currResp,
		"b": currResp,
		"c": currResp,
	}
	resp, err := agg.Aggregate(t.Context(), logger.Test(t), currResp.ID, responses, &currResp)
	require.NoError(t, err)
	assert.Equal(t, &currResp, resp)
}

func TestAggregator_SignedResponseMissingRequestID_Accepted(t *testing.T) {
	t.Parallel()
	payload := json.RawMessage([]byte(`{"responses":[{"error":"failed to verify ciphertext: cannot unmarshal data: unexpected end of JSON input","id":{"key":"W","namespace":"","owner":"foo"},"success":false}]}`))
	currResp, nodes := makeSignedVaultResponse(t, vaulttypes.MethodSecretsCreate, "expected-request", payload, 2)
	mcr := &mockCapabilitiesRegistry{F: 1, Nodes: nodes}
	agg := testAggregator(t, mcr)

	responses := map[string]jsonrpc.Response[json.RawMessage]{
		"a": currResp,
	}
	resp, err := agg.Aggregate(t.Context(), logger.Test(t), "expected-request", responses, &currResp)
	require.NoError(t, err)
	assert.Equal(t, &currResp, resp)
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
		Method:  vaulttypes.MethodSecretsCreate,
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
	agg := testAggregator(t, mcr)

	currResp := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsCreate,
		Result:  (*json.RawMessage)(nil),
		Error: &jsonrpc.WireError{
			Code:    123,
			Message: "some error",
		},
	}
	responses := map[string]jsonrpc.Response[json.RawMessage]{
		"a": currResp,
		"b": currResp,
		"c": currResp,
	}
	resp, err := agg.Aggregate(t.Context(), logger.Test(t), currResp.ID, responses, &currResp)
	require.NoError(t, err)
	assert.Equal(t, &currResp, resp)
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
	agg := testAggregator(t, mcr)

	oldResp1 := newMessage(t)
	oldResp2 := newMessage(t)
	currResp := newMessage(t)
	responses := map[string]jsonrpc.Response[json.RawMessage]{
		"a": *oldResp1,
		"b": *oldResp2,
		"c": *currResp,
	}
	resp, err := agg.Aggregate(t.Context(), logger.Test(t), currResp.ID, responses, currResp)
	require.NoError(t, err)

	respDigests := make([]string, 0, 3)
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

func makeUnsignedVaultRPCResponse(t *testing.T, payloadJSON string) jsonrpc.Response[json.RawMessage] {
	t.Helper()
	sor := vaulttypes.SignedOCRResponse{
		Payload:    json.RawMessage(payloadJSON),
		Context:    []byte{},
		Signatures: [][]byte{},
	}
	raw, err := json.Marshal(sor)
	require.NoError(t, err)
	rm := json.RawMessage(raw)
	return jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "quorum-tie",
		Method:  vaulttypes.MethodSecretsCreate,
		Result:  &rm,
	}
}

func TestValidateUsingQuorum_tiedMajoritiesPickDigestDeterministically(t *testing.T) {
	t.Parallel()

	a := &baseAggregator{}
	lggr := logger.Test(t)
	don := capabilities.DON{
		F:       1,
		Members: make([]p2ptypes.PeerID, 6),
	}

	ra := makeUnsignedVaultRPCResponse(t, `{"v":"aaa"}`)
	rb := makeUnsignedVaultRPCResponse(t, `{"v":"zzz"}`)
	digestA, err := a.sha(&ra)
	require.NoError(t, err)
	digestB, err := a.sha(&rb)
	require.NoError(t, err)
	require.NotEqual(t, digestA, digestB)

	wantWinner := min(digestB, digestA)

	for range 300 {
		m := map[string]jsonrpc.Response[json.RawMessage]{
			"n0": ra, "n1": ra, "n2": ra,
			"n3": rb, "n4": rb, "n5": rb,
		}
		got, err := a.validateUsingQuorum(t.Context(), don, nil, m, lggr)
		require.NoError(t, err)
		require.NotNil(t, got)
		gotDigest, derr := a.sha(got)
		require.NoError(t, derr)
		require.Equal(t, wantWinner, gotDigest,
			"with two disjoint 2F+1 majorities the winning digest must not depend on map iteration order")
	}
}

func TestAggregator_InsufficientResponses(t *testing.T) {
	mcr := &mockCapabilitiesRegistry{F: 1}
	agg := testAggregator(t, mcr)

	rm := json.RawMessage([]byte(`{}`))
	currResp := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsCreate,
		Result:  &rm,
	}
	responses := map[string]jsonrpc.Response[json.RawMessage]{
		"a": currResp,
	}
	_, err := agg.Aggregate(t.Context(), logger.Test(t), currResp.ID, responses, &currResp)
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
	agg := testAggregator(t, mcr)

	rm1 := json.RawMessage([]byte(`{}`))
	resp1 := &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsCreate,
		Result:  &rm1,
	}
	rm2 := json.RawMessage([]byte(`{"foo": "bar"}`))
	resp2 := &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsCreate,
		Result:  &rm2,
	}
	rm3 := json.RawMessage([]byte(`{"baz": "qux"}`))
	resp3 := &jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsCreate,
		Result:  &rm3,
	}
	responses := map[string]jsonrpc.Response[json.RawMessage]{
		"a": *resp1,
		"b": *resp2,
		"c": *resp3,
	}
	_, err := agg.Aggregate(t.Context(), logger.Test(t), resp3.ID, responses, resp3)
	require.ErrorContains(t, err, "failed to validate using quorum: quorum unobtainable")
}

func makeDONWithNodesForTest(t *testing.T, name string, id uint32, f uint8, memberOffset byte, nodeCount int) capabilities.DONWithNodes {
	t.Helper()
	nodes := make([]capabilities.Node, nodeCount)
	members := make([]p2ptypes.PeerID, nodeCount)
	for i := range nodeCount {
		pid := p2ptypes.PeerID{}
		pid[0] = memberOffset + byte(i)
		pid[1] = byte(i)
		nodes[i] = capabilities.Node{PeerID: &pid, Signer: [32]byte{}}
		members[i] = pid
	}
	return capabilities.DONWithNodes{
		DON: capabilities.DON{
			Name:    name,
			ID:      id,
			F:       f,
			Members: members,
		},
		Nodes: nodes,
	}
}

func TestAggregator_MultipleRegistryDONs_SelectsByVaultHandlerDonName(t *testing.T) {
	donOther := makeDONWithNodesForTest(t, "staging-vault", 1, 2, 0x10, 7)
	donMine := makeDONWithNodesForTest(t, "cre-reliability-vault", 2, 1, 0x20, 4)
	mcr := &mockCapabilitiesRegistry{DONs: []capabilities.DONWithNodes{donOther, donMine}}
	agg := &baseAggregator{
		capabilitiesRegistry: mcr,
		vaultHandlerDonID:    "cre-reliability-vault",
	}

	rm := json.RawMessage([]byte(`{}`))
	currResp := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsCreate,
		Result:  &rm,
	}
	responses := map[string]jsonrpc.Response[json.RawMessage]{
		"a": currResp,
		"b": currResp,
		"c": currResp,
	}
	resp, err := agg.Aggregate(t.Context(), logger.Test(t), currResp.ID, responses, &currResp)
	require.NoError(t, err)
	require.Equal(t, currResp.ID, resp.ID)
}

func TestAggregator_MultipleRegistryDONs_SelectsByIDWhenNameEmpty(t *testing.T) {
	donOther := makeDONWithNodesForTest(t, "", 1, 2, 0x10, 7)
	donMine := makeDONWithNodesForTest(t, "", 99, 1, 0x20, 4)
	mcr := &mockCapabilitiesRegistry{DONs: []capabilities.DONWithNodes{donOther, donMine}}
	agg := &baseAggregator{
		capabilitiesRegistry: mcr,
		vaultHandlerDonID:    "99",
	}

	rm := json.RawMessage([]byte(`{}`))
	currResp := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsCreate,
		Result:  &rm,
	}
	responses := map[string]jsonrpc.Response[json.RawMessage]{
		"a": currResp,
		"b": currResp,
		"c": currResp,
	}
	resp, err := agg.Aggregate(t.Context(), logger.Test(t), currResp.ID, responses, &currResp)
	require.NoError(t, err)
	require.Equal(t, currResp.ID, resp.ID)
}

func TestAggregator_MultipleRegistryDONs_NoMatchingVaultHandlerDonId(t *testing.T) {
	donA := makeDONWithNodesForTest(t, "don-a", 1, 1, 0x10, 4)
	donB := makeDONWithNodesForTest(t, "don-b", 2, 1, 0x20, 4)
	mcr := &mockCapabilitiesRegistry{DONs: []capabilities.DONWithNodes{donA, donB}}
	agg := &baseAggregator{
		capabilitiesRegistry: mcr,
		vaultHandlerDonID:    "unknown-vault",
	}

	rm := json.RawMessage([]byte(`{}`))
	currResp := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsCreate,
		Result:  &rm,
	}
	responses := map[string]jsonrpc.Response[json.RawMessage]{"a": currResp}
	_, err := agg.Aggregate(t.Context(), logger.Test(t), currResp.ID, responses, &currResp)
	require.ErrorContains(t, err, "none match vault handler DonId")
}

func TestAggregator_MultipleRegistryDONs_AmbiguousMatchingVaultHandlerDonId(t *testing.T) {
	donA := makeDONWithNodesForTest(t, "same-name", 1, 1, 0x10, 4)
	donB := makeDONWithNodesForTest(t, "same-name", 2, 1, 0x20, 4)
	mcr := &mockCapabilitiesRegistry{DONs: []capabilities.DONWithNodes{donA, donB}}
	agg := &baseAggregator{
		capabilitiesRegistry: mcr,
		vaultHandlerDonID:    "same-name",
	}

	rm := json.RawMessage([]byte(`{}`))
	currResp := jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsCreate,
		Result:  &rm,
	}
	responses := map[string]jsonrpc.Response[json.RawMessage]{"a": currResp}
	_, err := agg.Aggregate(t.Context(), logger.Test(t), currResp.ID, responses, &currResp)
	require.ErrorContains(t, err, "2 DONs match vault handler DonId")
}

func makeQuorumTestNode(t *testing.T, idx int) (*ecdsa.PrivateKey, capabilities.Node, common.Address) {
	t.Helper()
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	addr := crypto.PubkeyToAddress(key.PublicKey)
	signer := [32]byte{}
	copy(signer[:20], addr.Bytes())
	node := capabilities.Node{Signer: signer}
	return key, node, addr
}

func makeQuorumErrorResponse() jsonrpc.Response[json.RawMessage] {
	return jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "1",
		Method:  vaulttypes.MethodSecretsDelete,
		Error:   &jsonrpc.WireError{Code: -32000, Message: "boom"},
	}
}

func signEnvelope(t *testing.T, key *ecdsa.PrivateKey, resp *jsonrpc.Response[json.RawMessage]) {
	t.Helper()
	signEnvelopeForRequestID(t, key, resp, vaulttypes.NodeSignatureRequestID(resp.Method, resp.ID))
}

func signEnvelopeForRequestID(t *testing.T, key *ecdsa.PrivateKey, resp *jsonrpc.Response[json.RawMessage], requestID string) {
	t.Helper()
	digest, err := resp.Digest()
	require.NoError(t, err)
	sig, err := utils.GenerateEthSignature(key, vaulttypes.NodeSignaturePayload(requestID, digest))
	require.NoError(t, err)
	resp.NodeSignatures = [][]byte{sig}
}

func makeQuorumSignedOCRResponse(ocrSig string) jsonrpc.Response[json.RawMessage] {
	raw := json.RawMessage(`{"payload":"cGF5bG9hZA==","context":"","signatures":["` + ocrSig + `"]}`)
	return jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "0xabc" + vaulttypes.RequestIDSeparator + "1",
		Method:  vaulttypes.MethodSecretsList,
		Result:  &raw,
	}
}

func makeQuorumTestNodes(t *testing.T, n int) ([]*ecdsa.PrivateKey, []capabilities.Node, []common.Address) {
	t.Helper()
	keys := make([]*ecdsa.PrivateKey, n)
	nodes := make([]capabilities.Node, n)
	addrs := make([]common.Address, n)
	for i := range keys {
		keys[i], nodes[i], addrs[i] = makeQuorumTestNode(t, i)
	}
	return keys, nodes, addrs
}

func requireNodeSigsVerify(t *testing.T, resp *jsonrpc.Response[json.RawMessage], members []common.Address, minRequired int) {
	t.Helper()
	digest, err := resp.Digest()
	require.NoError(t, err)
	requestID := vaulttypes.NodeSignatureRequestID(resp.Method, resp.ID)
	seen := map[common.Address]struct{}{}
	for _, sig := range resp.NodeSignatures {
		signer, err := vaulttypes.RecoverNodeSigner(requestID, digest, sig)
		require.NoError(t, err)
		require.Contains(t, members, signer)
		require.NotContains(t, seen, signer, "duplicate node signer")
		seen[signer] = struct{}{}
	}
	require.GreaterOrEqual(t, len(seen), minRequired)
}

// Node signatures cover the full envelope, including OCR report signatures that
// a.sha strips; they must verify against the full digest, not the quorum key.
func TestValidateUsingQuorum_NodeSignaturesEnforced_SignedOCRResult(t *testing.T) {
	t.Parallel()

	a := &baseAggregator{nodeSignaturesEnabled: limits.NewGateLimiter(true)}
	don := capabilities.DON{F: 1, Members: make([]p2ptypes.PeerID, 4)}
	keys, nodes, addrs := makeQuorumTestNodes(t, 4)

	resps := map[string]jsonrpc.Response[json.RawMessage]{}
	for i := range 3 {
		r := makeQuorumSignedOCRResponse("c2ln")
		signEnvelope(t, keys[i], &r)
		resps[strings.ToLower(addrs[i].Hex())] = r
	}

	got, err := a.validateUsingQuorum(t.Context(), don, nodes, resps, logger.Test(t))
	require.NoError(t, err)
	require.Len(t, got.NodeSignatures, 3)
	requireNodeSigsVerify(t, got, addrs, 3)
}

// When OCR report signatures differ between nodes, responses still agree under
// a.sha, but only signatures over the returned envelope may be attached to it.
func TestValidateUsingQuorum_NodeSignaturesEnforced_DifferingOCRSignatures(t *testing.T) {
	t.Parallel()

	a := &baseAggregator{nodeSignaturesEnabled: limits.NewGateLimiter(true)}
	don := capabilities.DON{F: 1, Members: make([]p2ptypes.PeerID, 4)}
	keys, nodes, addrs := makeQuorumTestNodes(t, 4)

	resps := map[string]jsonrpc.Response[json.RawMessage]{}
	for i, ocrSig := range []string{"c2lnMQ==", "c2lnMg==", "c2lnMw=="} {
		r := makeQuorumSignedOCRResponse(ocrSig)
		signEnvelope(t, keys[i], &r)
		resps[strings.ToLower(addrs[i].Hex())] = r
	}

	got, err := a.validateUsingQuorum(t.Context(), don, nodes, resps, logger.Test(t))
	require.NoError(t, err)
	require.Len(t, got.NodeSignatures, 1)
	requireNodeSigsVerify(t, got, addrs, 1)
}

// A response that fails node-signature verification must not be returned, even
// when its content matches the quorum digest and its key sorts first.
func TestValidateUsingQuorum_NodeSignaturesEnforced_UnverifiedResponseNotReturned(t *testing.T) {
	t.Parallel()

	a := &baseAggregator{nodeSignaturesEnabled: limits.NewGateLimiter(true)}
	don := capabilities.DON{F: 1, Members: make([]p2ptypes.PeerID, 4)}
	keys, nodes, addrs := makeQuorumTestNodes(t, 4)

	resps := map[string]jsonrpc.Response[json.RawMessage]{}
	for i := range 3 {
		r := makeQuorumSignedOCRResponse("c2ln")
		signEnvelope(t, keys[i], &r)
		resps[strings.ToLower(addrs[i].Hex())] = r
	}
	junk := makeQuorumSignedOCRResponse("anVuaw==")
	junk.NodeSignatures = [][]byte{make([]byte, 65)}
	resps["0x0000000000000000000000000000000000000000"] = junk

	got, err := a.validateUsingQuorum(t.Context(), don, nodes, resps, logger.Test(t))
	require.NoError(t, err)
	assert.JSONEq(t, string(*resps[strings.ToLower(addrs[0].Hex())].Result), string(*got.Result))
	requireNodeSigsVerify(t, got, addrs, 3)
}

// A node signature bound to a different request ID must not count: signed
// responses can't be replayed as the answer to another request.
func TestValidateUsingQuorum_NodeSignaturesEnforced_ReplayedRequestIDRejected(t *testing.T) {
	t.Parallel()

	a := &baseAggregator{nodeSignaturesEnabled: limits.NewGateLimiter(true)}
	don := capabilities.DON{F: 1, Members: make([]p2ptypes.PeerID, 4)}
	keys, nodes, addrs := makeQuorumTestNodes(t, 4)

	resps := map[string]jsonrpc.Response[json.RawMessage]{}
	for i := range 2 {
		r := makeQuorumErrorResponse()
		signEnvelope(t, keys[i], &r)
		resps[strings.ToLower(addrs[i].Hex())] = r
	}
	replayed := makeQuorumErrorResponse()
	signEnvelopeForRequestID(t, keys[2], &replayed, "some-earlier-request")
	resps[strings.ToLower(addrs[2].Hex())] = replayed

	_, err := a.validateUsingQuorum(t.Context(), don, nodes, resps, logger.Test(t))
	require.ErrorIs(t, err, errInsufficientResponsesForQuorum)
}

// A signer counts once: a node replaying another member's signature must not
// pad quorum.
func TestValidateUsingQuorum_NodeSignaturesEnforced_DuplicateSignerCountsOnce(t *testing.T) {
	t.Parallel()

	a := &baseAggregator{nodeSignaturesEnabled: limits.NewGateLimiter(true)}
	don := capabilities.DON{F: 1, Members: make([]p2ptypes.PeerID, 4)}
	keys, nodes, addrs := makeQuorumTestNodes(t, 4)

	resps := map[string]jsonrpc.Response[json.RawMessage]{}
	for i := range 2 {
		r := makeQuorumErrorResponse()
		signEnvelope(t, keys[i], &r)
		resps[strings.ToLower(addrs[i].Hex())] = r
	}
	resps[strings.ToLower(addrs[2].Hex())] = resps[strings.ToLower(addrs[0].Hex())]

	_, err := a.validateUsingQuorum(t.Context(), don, nodes, resps, logger.Test(t))
	require.ErrorIs(t, err, errInsufficientResponsesForQuorum)
}

func TestValidateUsingQuorum_NodeSignaturesEnforced(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	a := &baseAggregator{nodeSignaturesEnabled: limits.NewGateLimiter(true)}
	don := capabilities.DON{F: 1, Members: make([]p2ptypes.PeerID, 4)}

	keys := make([]*ecdsa.PrivateKey, 4)
	nodes := make([]capabilities.Node, 4)
	addrs := make([]common.Address, 4)
	for i := range keys {
		keys[i], nodes[i], addrs[i] = makeQuorumTestNode(t, i)
	}

	resps := map[string]jsonrpc.Response[json.RawMessage]{}
	for i := range 3 {
		r := makeQuorumErrorResponse()
		signEnvelope(t, keys[i], &r)
		resps[strings.ToLower(addrs[i].Hex())] = r
	}
	// Fourth node replies unsigned; it must not count toward quorum.
	resps[strings.ToLower(addrs[3].Hex())] = makeQuorumErrorResponse()

	got, err := a.validateUsingQuorum(t.Context(), don, nodes, resps, lggr)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.Len(t, got.NodeSignatures, 3)
	requireNodeSigsVerify(t, got, addrs[:3], 3)
}

func TestValidateUsingQuorum_NodeSignaturesEnforced_NonMemberSignerDropped(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	a := &baseAggregator{nodeSignaturesEnabled: limits.NewGateLimiter(true)}
	don := capabilities.DON{F: 1, Members: make([]p2ptypes.PeerID, 4)}

	keys := make([]*ecdsa.PrivateKey, 4)
	nodes := make([]capabilities.Node, 4)
	addrs := make([]common.Address, 4)
	for i := range keys {
		keys[i], nodes[i], addrs[i] = makeQuorumTestNode(t, i)
	}

	outsider, _, _ := makeQuorumTestNode(t, 4)

	resps := map[string]jsonrpc.Response[json.RawMessage]{}
	for i := range 2 {
		r := makeQuorumErrorResponse()
		signEnvelope(t, keys[i], &r)
		resps[strings.ToLower(addrs[i].Hex())] = r
	}
	// Third node's response is signed by a key outside the DON signer set:
	// the signature recovers to a non-member address and must be dropped.
	rSignedByOutsider := makeQuorumErrorResponse()
	signEnvelope(t, outsider, &rSignedByOutsider)
	resps[strings.ToLower(addrs[2].Hex())] = rSignedByOutsider

	_, err := a.validateUsingQuorum(t.Context(), don, nodes, resps, lggr)
	require.ErrorIs(t, err, errInsufficientResponsesForQuorum)
}

func TestValidateUsingQuorum_NodeSignaturesDisabled(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	a := &baseAggregator{nodeSignaturesEnabled: limits.NewGateLimiter(false)}
	don := capabilities.DON{F: 1, Members: make([]p2ptypes.PeerID, 4)}

	keys := make([]*ecdsa.PrivateKey, 4)
	nodes := make([]capabilities.Node, 4)
	addrs := make([]common.Address, 4)
	for i := range keys {
		keys[i], nodes[i], addrs[i] = makeQuorumTestNode(t, i)
	}

	// Unsigned responses still reach quorum when enforcement is off.
	resps := map[string]jsonrpc.Response[json.RawMessage]{}
	for i := range 3 {
		resps[strings.ToLower(addrs[i].Hex())] = makeQuorumErrorResponse()
	}

	got, err := a.validateUsingQuorum(t.Context(), don, nodes, resps, lggr)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Empty(t, got.NodeSignatures)
}
