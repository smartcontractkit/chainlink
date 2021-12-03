package ocr2key_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertKeyBundlesNotEqual(t *testing.T, pk1 *ocr2key.KeyBundle, pk2 *ocr2key.KeyBundle) {
	assert.NotEqual(t, pk1.ID(), pk2.ID())
	assert.NotEqualValues(t, pk1.OffchainKeyring, pk2.OffchainKeyring)
	assert.NotEqualValues(t, pk1.OnchainKeyring(), pk2.OnchainKeyring())
}

func TestOCR2keys_New(t *testing.T) {
	t.Parallel()
	pk1, err := ocr2key.New("evm")
	require.NoError(t, err)
	pk2, err := ocr2key.New("evm")
	require.NoError(t, err)
	pk3, err := ocr2key.New("evm")
	require.NoError(t, err)
	assertKeyBundlesNotEqual(t, pk1, pk2)
	assertKeyBundlesNotEqual(t, pk1, pk3)
	assertKeyBundlesNotEqual(t, pk2, pk3)
}
