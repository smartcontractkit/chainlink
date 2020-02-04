package vrfkey

import (
	"testing"

	"chainlink/core/services/signatures/cryptotest"

	"github.com/stretchr/testify/require"
)

func TestValueScanIdentityPointSet(t *testing.T) {
	randomStream := cryptotest.NewStream(t, 0)
	for i := 0; i < 10; i++ {
		p := suite.Point().Pick(randomStream)
		var pk, nPk, nnPk PublicKey
		marshaledKey, err := p.MarshalBinary()
		require.NoError(t, err)
		require.Equal(t, copy(pk[:], marshaledKey),
			CompressedPublicKeyLength)
		require.NotEqual(t, pk, nPk)
		np, err := pk.Point()
		require.NoError(t, err)
		require.True(t, p.Equal(np))
		value, err := pk.Value()
		require.NoError(t, err)
		nPk.Scan(value)
		require.Equal(t, pk, nPk)
		nnPk.Set(pk)
		require.Equal(t, pk, nnPk)
	}
}
