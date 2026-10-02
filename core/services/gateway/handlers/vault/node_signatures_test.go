package vault

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	vaultcapmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	vaulttypesmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	connector_mocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	gwcommon "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

// requireClientVerifiesNodeSigs verifies raw (a user-facing gateway response)
// the way a client would, against the ID the client sees.
func requireClientVerifiesNodeSigs(t *testing.T, raw []byte, wantID string, members []common.Address, minRequired int) {
	t.Helper()
	var resp jsonrpc.Response[json.RawMessage]
	require.NoError(t, json.Unmarshal(raw, &resp))
	require.Equal(t, wantID, resp.ID)
	requireNodeSigsVerify(t, &resp, members, minRequired)
}

// Contract test: responses produced by real vault node handlers, signed with
// real OCR2 keys, must pass gateway enforcement and verify for the client. The
// response carries OCR report signatures that fail the fast path, so this also
// covers the quorum fallback over OCR-signed results.
func TestNodeSignatures_NodeToGatewayContract(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	const clientRequestID = "req-1"
	prefixedID := owner + vaulttypes.RequestIDSeparator + clientRequestID
	params, err := json.Marshal(vaultcommon.DeleteSecretsRequest{
		RequestId: clientRequestID,
		Ids:       []*vaultcommon.SecretIdentifier{{Key: "Foo", Namespace: "main", Owner: owner}},
	})
	require.NoError(t, err)

	nodes := make([]capabilities.Node, 4)
	members := make([]common.Address, 4)
	nodeResps := make([]*jsonrpc.Response[json.RawMessage], 3)
	for i := range nodes {
		var kb ocr2key.KeyBundle
		kb, err = ocr2key.New(corekeys.EVM)
		require.NoError(t, err)
		members[i] = common.HexToAddress(kb.OnChainPublicKey())
		var signer [32]byte
		copy(signer[:20], members[i].Bytes())
		nodes[i] = capabilities.Node{PeerID: &p2ptypes.PeerID{0: uint8(i)}, Signer: signer}
		if i >= len(nodeResps) {
			continue // fourth member never responds
		}

		var envelopeSigner connector.Signer
		envelopeSigner, err = vaultcap.NewOCR2KeySigner(kb)
		require.NoError(t, err)
		secretsService := vaulttypesmocks.NewSecretsService(t)
		secretsService.EXPECT().DeleteSecrets(mock.Anything, mock.Anything).Return(&vaulttypes.Response{
			ID:         clientRequestID,
			Payload:    []byte(`{}`),
			Context:    []byte("ctx"),
			Signatures: [][]byte{[]byte("not-a-valid-ocr-signature")},
		}, nil)
		auth := vaultcapmocks.NewAuthorizer(t)
		auth.EXPECT().AuthorizeRequest(mock.Anything, mock.Anything).
			Return(vaultcap.NewAuthResult("", owner, "digest", time.Now().Add(time.Minute).Unix()), nil)
		gwConnector := connector_mocks.NewGatewayConnector(t)
		gwConnector.On("SendToGateway", mock.Anything, "gateway-1", mock.Anything).
			Run(func(args mock.Arguments) { nodeResps[i] = args.Get(2).(*jsonrpc.Response[json.RawMessage]) }).
			Return(nil).Once()

		var nodeHandler *vaultcap.GatewayHandler
		nodeHandler, err = vaultcap.NewGatewayHandler(secretsService, gwConnector, envelopeSigner, nil, lggr,
			limits.Factory{Settings: cresettings.DefaultGetter}, vaultcap.NewAuthorizer(auth, nil, lggr), nil)
		require.NoError(t, err)

		raw := json.RawMessage(params)
		require.NoError(t, nodeHandler.HandleGatewayMessage(t.Context(), "gateway-1", &jsonrpc.Request[json.RawMessage]{
			Version: jsonrpc.JsonRpcVersion, ID: prefixedID, Method: vaulttypes.MethodSecretsDelete, Params: &raw,
		}))
		require.NotNil(t, nodeResps[i])
		require.Equal(t, prefixedID, nodeResps[i].ID)
		require.Nil(t, nodeResps[i].Error, "node must return an OCR-signed result, not an error")
		require.Contains(t, string(*nodeResps[i].Result), `"signatures"`)
	}

	h, callback, _, _ := setupHandler(t)
	h.(*handler).aggregator = &baseAggregator{
		capabilitiesRegistry:  &mockCapabilitiesRegistry{F: 1, Nodes: nodes},
		metrics:               h.(*handler).metrics,
		vaultHandlerDonID:     h.(*handler).donConfig.DonID,
		nodeSignaturesEnabled: limits.NewGateLimiter(true),
	}
	_, err = h.(*handler).newActiveRequest(jsonrpc.Request[json.RawMessage]{ID: prefixedID, Method: vaulttypes.MethodSecretsDelete}, callback)
	require.NoError(t, err)
	for i, r := range nodeResps {
		require.NoError(t, h.HandleNodeMessage(t.Context(), r, fmt.Sprintf("0xnode%d", i)))
	}

	resp, err := callback.Wait(t.Context())
	require.NoError(t, err)
	requireClientVerifiesNodeSigs(t, resp.RawResponse, clientRequestID, members, 3)
}

// When the signed-OCR fast path fails, Aggregate must fall back to quorum and
// still accept validly node-signed responses.
func TestAggregate_NodeSignaturesEnforced_FallsBackToQuorum(t *testing.T) {
	t.Parallel()

	keys, nodes, addrs := makeQuorumTestNodes(t, 4)
	for i := range nodes {
		nodes[i].PeerID = &p2ptypes.PeerID{0: uint8(i)}
	}
	a := &baseAggregator{
		capabilitiesRegistry:  &mockCapabilitiesRegistry{F: 1, Nodes: nodes},
		metrics:               testAggregator(t, nil).metrics,
		nodeSignaturesEnabled: limits.NewGateLimiter(true),
	}

	resps := map[string]jsonrpc.Response[json.RawMessage]{}
	for i := range 3 {
		r := makeQuorumSignedOCRResponse("c2ln") // not a valid OCR signature
		signEnvelope(t, keys[i], &r)
		resps[fmt.Sprintf("0xnode%d", i)] = r
	}
	curr := resps["0xnode2"]

	got, err := a.Aggregate(t.Context(), logger.Test(t), curr.ID, resps, &curr)
	require.NoError(t, err)
	requireNodeSigsVerify(t, got, addrs, 3)
}

// Public key signatures are cached and replayed to later callers; they must
// still verify for a different request ID.
func TestVaultHandler_PublicKeyGet_CachedNodeSignatures(t *testing.T) {
	t.Parallel()

	h, callback, don, _ := setupHandler(t)
	keys, nodes, addrs := makeQuorumTestNodes(t, 4)
	for i := range nodes {
		nodes[i].PeerID = &p2ptypes.PeerID{0: uint8(i)}
	}
	h.(*handler).aggregator = &baseAggregator{
		capabilitiesRegistry:  &mockCapabilitiesRegistry{F: 1, Nodes: nodes},
		metrics:               h.(*handler).metrics,
		vaultHandlerDonID:     h.(*handler).donConfig.DonID,
		nodeSignaturesEnabled: limits.NewGateLimiter(true),
	}
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ar, err := h.(*handler).newActiveRequest(jsonrpc.Request[json.RawMessage]{ID: "request_id", Method: vaulttypes.MethodPublicKeyGet}, callback)
	require.NoError(t, err)
	require.NoError(t, h.(*handler).handlePublicKeyGet(t.Context(), ar))

	_, pk, _, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	pkBytes, err := pk.Marshal()
	require.NoError(t, err)
	resultBytes, err := json.Marshal(&vaultcommon.GetPublicKeyResponse{PublicKey: hex.EncodeToString(pkBytes)})
	require.NoError(t, err)
	for i := range 3 {
		r := jsonrpc.Response[json.RawMessage]{
			Version: jsonrpc.JsonRpcVersion,
			ID:      "request_id",
			Method:  vaulttypes.MethodPublicKeyGet,
			Result:  (*json.RawMessage)(&resultBytes),
		}
		signEnvelope(t, keys[i], &r)
		require.NoError(t, h.HandleNodeMessage(t.Context(), &r, fmt.Sprintf("0xnode%d", i)))
	}

	resp, err := callback.Wait(t.Context())
	require.NoError(t, err)
	requireClientVerifiesNodeSigs(t, resp.RawResponse, "request_id", addrs, 3)

	// Served from cache.
	callback = gwcommon.NewCallback()
	require.NoError(t, h.HandleJSONRPCUserMessage(t.Context(), jsonrpc.Request[json.RawMessage]{
		ID: "another_request_id", Method: vaulttypes.MethodPublicKeyGet,
	}, callback))
	resp, err = callback.Wait(t.Context())
	require.NoError(t, err)
	requireClientVerifiesNodeSigs(t, resp.RawResponse, "another_request_id", addrs, 3)
}
