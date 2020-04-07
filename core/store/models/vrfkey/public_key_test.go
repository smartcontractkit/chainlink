package vrfkey

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/signatures/cryptotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValueScanIdentityPointSet(t *testing.T) {
	randomStream := cryptotest.NewStream(t, 0)
	for i := 0; i < 10; i++ {
		p := suite.Point().Pick(randomStream)
		var pk, nPk, nnPk PublicKey
		marshaledKey, err := p.MarshalBinary()
		require.NoError(t, err, "failed to marshal public key")
		require.Equal(t, copy(pk[:], marshaledKey),
			CompressedPublicKeyLength, "failed to copy marshaled key to pk")
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
