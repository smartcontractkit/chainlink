package vrfkey

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
)

func TestVRFKeys_KeyV2(t *testing.T) {
	k, err := NewV2()
	require.NoError(t, err)

	assert.Equal(t, hexutil.Encode(k.PublicKey[:]), k.ID())
	assert.Equal(t, internal.NewRaw(secp256k1.ToInt(*k.k).Bytes()), k.Raw())

	t.Run("generates proof", func(t *testing.T) {
		p, err := k.GenerateProof(big.NewInt(1))

		assert.NotZero(t, p)
		assert.NoError(t, err)
	})
}
