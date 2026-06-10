package vault

import (
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

func testOCRContext(t *testing.T) []byte {
	t.Helper()
	ctx, err := hex.DecodeString("000ec4f6a2ba011e909eccf64628855b848e08876a1edd938a1372a9e51adff100000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	return ctx
}

func makeSignedCreateSecretsResponse(t *testing.T, requestID string, numSigners int) (jsonrpc.Response[json.RawMessage], []capabilities.Node) {
	t.Helper()

	createResp := &vaultcommon.CreateSecretsResponse{
		RequestId: requestID,
		Responses: []*vaultcommon.CreateSecretResponse{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:       "W",
					Namespace: "",
					Owner:     "foo",
				},
				Success: false,
				Error:   "failed to verify ciphertext: cannot unmarshal data: unexpected end of JSON input",
			},
		},
	}
	payload, err := vaultutils.ToCanonicalJSON(createResp)
	require.NoError(t, err)

	return makeSignedVaultResponse(t, vaulttypes.MethodSecretsCreate, requestID, json.RawMessage(payload), numSigners)
}

func makeSignedVaultResponse(t *testing.T, method, requestID string, payload json.RawMessage, numSigners int) (jsonrpc.Response[json.RawMessage], []capabilities.Node) {
	t.Helper()
	require.GreaterOrEqual(t, numSigners, 1)

	ctx := testOCRContext(t)
	keys := make([]*ecdsa.PrivateKey, numSigners)
	nodes := make([]capabilities.Node, numSigners)
	for i := range keys {
		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		keys[i] = key

		addr := crypto.PubkeyToAddress(key.PublicKey)
		var signer [32]byte
		copy(signer[:20], addr.Bytes())
		pid := p2ptypes.PeerID{0: uint8(i)}
		nodes[i] = capabilities.Node{PeerID: &pid, Signer: signer}
	}

	configDigest, err := ocr2types.BytesToConfigDigest(ctx[:32])
	require.NoError(t, err)
	epoch := binary.BigEndian.Uint32(ctx[32+27 : 32+31])
	round := ctx[63]

	fullHash := ocr2key.ReportToSigData(ocr2types.ReportContext{
		ReportTimestamp: ocr2types.ReportTimestamp{
			ConfigDigest: configDigest,
			Epoch:        epoch,
			Round:        round,
		},
	}, []byte(payload))

	signatures := make([][]byte, numSigners)
	for i, key := range keys {
		sig, err := crypto.Sign(fullHash, key)
		require.NoError(t, err)
		signatures[i] = sig
	}

	sor := vaulttypes.SignedOCRResponse{
		Payload:    payload,
		Context:    ctx,
		Signatures: signatures,
	}
	rawResp, err := json.Marshal(sor)
	require.NoError(t, err)
	rm := json.RawMessage(rawResp)

	return jsonrpc.Response[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      requestID,
		Method:  method,
		Result:  &rm,
	}, nodes
}
