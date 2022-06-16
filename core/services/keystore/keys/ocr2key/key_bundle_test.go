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
	pk1, err := ocr2key.New(chaintype.EVM)
	require.NoError(t, err)
	pk2, err := ocr2key.New(chaintype.EVM)
	require.NoError(t, err)
	pk3, err := ocr2key.New(chaintype.Solana)
	require.NoError(t, err)
	pk4, err := ocr2key.New(chaintype.Solana)
	require.NoError(t, err)
	pk5, err := ocr2key.New(chaintype.Terra)
	require.NoError(t, err)
	pk6, err := ocr2key.New(chaintype.Terra)
	require.NoError(t, err)
	_, err = ocr2key.New("invalid")
	assert.Error(t, err)
	assertKeyBundlesNotEqual(t, pk1, pk2)
	assertKeyBundlesNotEqual(t, pk3, pk4)
	assertKeyBundlesNotEqual(t, pk1, pk3)
	assertKeyBundlesNotEqual(t, pk5, pk6)
	assertKeyBundlesNotEqual(t, pk1, pk5)
	assertKeyBundlesNotEqual(t, pk3, pk5)
	assert.Equal(t, pk1.ChainType(), pk2.ChainType())
	assert.Equal(t, pk3.ChainType(), pk4.ChainType())
	assert.Equal(t, pk5.ChainType(), pk6.ChainType())
	assert.NotEqual(t, pk1.ChainType(), pk3.ChainType())
	assert.NotEqual(t, pk1.ChainType(), pk5.ChainType())
	assert.NotEqual(t, pk3.ChainType(), pk5.ChainType())
}

func TestOCR2KeyBundle_RawToKey(t *testing.T) {
	t.Parallel()

	pk1, err := ocr2key.New(chaintype.EVM)
	require.NoError(t, err)
	pk2, err := ocr2key.New(chaintype.Solana)
	require.NoError(t, err)
	pk3, err := ocr2key.New(chaintype.Terra)
	require.NoError(t, err)

	pk1FromRaw := pk1.Raw().Key()
	pk2FromRaw := pk2.Raw().Key()
	pk3FromRaw := pk3.Raw().Key()

	assert.NotNil(t, pk1FromRaw)
	assert.NotNil(t, pk2FromRaw)
	assert.NotNil(t, pk3FromRaw)
}

func TestOCR2KeyBundle_BundleBase(t *testing.T) {
	t.Parallel()

	kb, err := ocr2key.New(chaintype.EVM)
	require.NoError(t, err)

	assert.NotNil(t, kb.ID())
	assert.Equal(t, fmt.Sprintf(`bundle: KeyBundle{chainType: evm, id: %s}`, kb.ID()), fmt.Sprintf(`bundle: %s`, kb))
}
