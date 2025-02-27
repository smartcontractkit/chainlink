package keystore_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

func TestTestKeystore(t *testing.T) {
	t.Parallel()
	ks := keystore.NewTestKeyStore(t)
	t.Run("GenerateEthKeys", func(t *testing.T) {
		// Define 20 chain IDs
		chainIDs := make([]*big.Int, 20)
		for i := 0; i < 20; i++ {
			chainIDs[i] = big.NewInt(int64(i + 1))
		}

		// Generate 20 Ethereum keys
		keys := ks.GenerateEthKeys(chainIDs...)

		// Validate the generated keys
		require.Len(t, keys, 20)
		dedup := make(map[string]struct{})
		for i, key := range keys {
			assert.Equal(t, chainIDs[i].Uint64(), key.EVMChainID)
			assert.NotEmpty(t, key.JSON)
			assert.Equal(t, "password", key.Password)
			assert.NotContains(t, dedup, key.JSON)
			dedup[key.JSON] = struct{}{}
		}
	})
	t.Run("GenerateP2PKeys", func(t *testing.T) {
		key := ks.GenerateP2PKey()
		assert.NotEmpty(t, key.JSON)
		assert.Equal(t, "password", key.Password)
	})
}
