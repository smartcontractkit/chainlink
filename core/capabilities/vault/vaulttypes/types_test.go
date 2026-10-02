package vaulttypes

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestIsGatewaySecretsMethod(t *testing.T) {
	t.Parallel()

	for _, method := range GatewaySecretsMethods {
		assert.True(t, IsGatewaySecretsMethod(method), method)
	}
	assert.False(t, IsGatewaySecretsMethod(MethodPublicKeyGet))
	assert.False(t, IsGatewaySecretsMethod(MethodSecretsGet))
	assert.False(t, IsGatewaySecretsMethod("vault.unsupported"))
}

func TestRecoverNodeSigner(t *testing.T) {
	t.Parallel()

	digest := "6acc256f1ba5bf28f834040ecd80ff28d316aa1806ef21ace295f45226b0a66d"
	requestID := "req-1"

	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	addr := crypto.PubkeyToAddress(key.PublicKey)
	sig, err := utils.GenerateEthSignature(key, NodeSignaturePayload(requestID, digest))
	require.NoError(t, err)

	t.Run("recovers signer", func(t *testing.T) {
		t.Parallel()
		got, err := RecoverNodeSigner(requestID, digest, sig)
		require.NoError(t, err)
		assert.Equal(t, addr, got)
	})

	t.Run("different request ID recovers a different address", func(t *testing.T) {
		t.Parallel()
		got, err := RecoverNodeSigner("req-2", digest, sig)
		require.NoError(t, err)
		assert.NotEqual(t, addr, got)
	})

	t.Run("different digest recovers a different address", func(t *testing.T) {
		t.Parallel()
		got, err := RecoverNodeSigner(requestID, "a8e6e0e131f05fc98fedc4826550a13bbdeb93b01bd9b97aa462a2411860e7d9", sig)
		require.NoError(t, err)
		assert.NotEqual(t, addr, got)
	})

	t.Run("malformed signature rejected", func(t *testing.T) {
		t.Parallel()
		_, err := RecoverNodeSigner(requestID, digest, []byte{0x01, 0x02})
		require.Error(t, err)
	})

	t.Run("does not mutate the signature", func(t *testing.T) {
		t.Parallel()
		legacyV := append([]byte{}, sig...)
		legacyV[64] += 27
		before := append([]byte{}, legacyV...)
		got, err := RecoverNodeSigner(requestID, digest, legacyV)
		require.NoError(t, err)
		assert.Equal(t, addr, got)
		assert.Equal(t, before, legacyV)
	})
}

func TestNodeSignatureRequestID(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "req-1", NodeSignatureRequestID(MethodSecretsList, "0xabc"+RequestIDSeparator+"req-1"))
	assert.Equal(t, "req-1", NodeSignatureRequestID(MethodSecretsList, "req-1"))
	// The gateway replays cached public key signatures across requests.
	assert.Empty(t, NodeSignatureRequestID(MethodPublicKeyGet, "req-1"))
}

func TestNodeSignaturePayload_LengthPrefixed(t *testing.T) {
	t.Parallel()

	// Without a length prefix, ("ab", "c...") and ("a", "bc...") would collide.
	assert.NotEqual(t, NodeSignaturePayload("ab", "cd"), NodeSignaturePayload("a", "bcd"))
}

// TestNodeSignaturePayload_GoldenVector pins the node signature wire format.
// Clients (cre-cli, SDKs) verify against these bytes; if this test fails, the
// format changed and every verifier must change with it.
func TestNodeSignaturePayload_GoldenVector(t *testing.T) {
	t.Parallel()

	const (
		requestID = "req-1"
		digest    = "6acc256f1ba5bf28f834040ecd80ff28d316aa1806ef21ace295f45226b0a66d"
		payload   = "0eedc4a332bf0e745ce2d39a7fd9f39e5911bcb289a095ac23b2a8cacd55b654"
		// Signed by private key 4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318.
		sig    = "d104b8b4b8d9199c770d74f83a01b8450a5f7ee4ef6313b1acb39fff7c01790e499cda7c60e5bcb7b0d5ffde2b53a9b14398acdeb45df44c92f7fcd98b7d0c5500"
		signer = "0x2c7536E3605D9C16a7a3D7b1898e529396a65c23"
	)

	assert.Equal(t, payload, hex.EncodeToString(NodeSignaturePayload(requestID, digest)))

	sigBytes, err := hex.DecodeString(sig)
	require.NoError(t, err)
	got, err := RecoverNodeSigner(requestID, digest, sigBytes)
	require.NoError(t, err)
	assert.Equal(t, signer, got.Hex())
}
