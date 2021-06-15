package vrfkey

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"

	"github.com/smartcontractkit/chainlink/core/services/signatures/cryptotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValueScanIdentityPointSet(t *testing.T) {
	randomStream := cryptotest.NewStream(t, 0)
	for i := 0; i < 10; i++ {
		p := suite.Point().Pick(randomStream)
		var pk, nPk, nnPk secp256k1.PublicKey
		marshaledKey, err := p.MarshalBinary()
		require.NoError(t, err, "failed to marshal public key")
		require.Equal(t, copy(pk[:], marshaledKey),
			secp256k1.CompressedPublicKeyLength, "failed to copy marshaled key to pk")
		assert.NotEqual(t, pk, nPk, "equality test succeeds on different keys!")
		np, err := pk.Point()
		require.NoError(t, err, "failed to marshal public key")
		assert.True(t, p.Equal(np), "Point should give the point we constructed pk from")
		value, err := pk.Value()
		require.NoError(t, err, "failed to serialize public key for database")
		nPk.Scan(value)
		assert.Equal(t, pk, nPk,
			"roundtripping public key through db Value/Scan gave different key!")
		nnPk.Set(pk)
		assert.Equal(t, pk, nnPk,
			"setting one PubliKey to another should result in equal keys")
	}

}

// Tests that PublicKey.Hash gives the same result as the VRFCoordinator's
func TestHash(t *testing.T) {
	pk, err := secp256k1.NewPublicKeyFromHex("0x9dc09a0f898f3b5e8047204e7ce7e44b587920932f08431e29c9bf6923b8450a01")
	assert.NoError(t, err)
	assert.Equal(t, "0xc4406d555db624837188b91514a5f47e34d825d935ab887a35c06a3e7c41de69", pk.MustHash().String())
}
