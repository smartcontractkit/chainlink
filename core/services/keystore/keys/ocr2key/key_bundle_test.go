package ocr2key_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
)

func assertKeyBundlesNotEqual(t *testing.T, pk1 ocr2key.KeyBundle, pk2 ocr2key.KeyBundle) {
	assert.NotEqual(t, pk1.ID(), pk2.ID())
	assert.NotEqualValues(t, pk1.OffchainPublicKey(), pk2.OffchainPublicKey())
	assert.NotEqualValues(t, pk1.OnChainPublicKey(), pk2.OnChainPublicKey())
}

func TestOCR2Keys_New(t *testing.T) {
	t.Parallel()
	var keys []ocr2key.KeyBundle

	// create two keys for each chain type
	for _, chain := range chaintype.SupportedChainTypes {
		pk0, err := ocr2key.New(chain)
		require.NoError(t, err)
		pk1, err := ocr2key.New(chain)
		require.NoError(t, err)

		keys = append(keys, pk0)
		keys = append(keys, pk1)
	}

	// validate keys are unique
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			assertKeyBundlesNotEqual(t, keys[i], keys[j])
		}
	}

	// validate chain types
	for i := 0; i < len(keys); i += 2 {
		// check key for same chain
		require.Equal(t, keys[i].ChainType(), keys[i+1].ChainType())

		// check 1 key for each chain
		for j := i + 2; j < len(keys); j += 2 {
			require.NotEqual(t, keys[i].ChainType(), keys[j].ChainType())
		}
	}
}

func TestOCR2KeyBundle_RawToKey(t *testing.T) {
	t.Parallel()

	for _, chain := range chaintype.SupportedChainTypes {
		pk, err := ocr2key.New(chain)
		require.NoError(t, err)

		pkFromRaw := pk.Raw().Key()
		assert.NotNil(t, pkFromRaw)
	}
}

func TestOCR2KeyBundle_BundleBase(t *testing.T) {
	t.Parallel()

	for _, chain := range chaintype.SupportedChainTypes {
		kb, err := ocr2key.New(chain)
		require.NoError(t, err)

		assert.NotNil(t, kb.ID())
		assert.Equal(t, fmt.Sprintf(`bundle: KeyBundle{chainType: %s, id: %s}`, chain, kb.ID()), fmt.Sprintf(`bundle: %s`, kb))
	}
}
