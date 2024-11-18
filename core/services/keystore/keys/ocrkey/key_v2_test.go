package ocrkey_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocrkey"
)

func assertKeyBundlesNotEqual(t *testing.T, pk1 ocrkey.KeyV2, pk2 ocrkey.KeyV2) {
	assert.NotEqual(t, pk1.ID(), pk2.ID())
	assert.NotEqual(t, pk1.ExportedOnChainSigning().X, pk2.ExportedOnChainSigning().X)
	assert.NotEqual(t, pk1.ExportedOnChainSigning().Y, pk2.ExportedOnChainSigning().Y)
	assert.NotEqual(t, pk1.ExportedOnChainSigning().D, pk2.ExportedOnChainSigning().D)
	assert.NotEqual(t, pk1.ExportedOffChainSigning().PublicKey(), pk2.ExportedOffChainSigning().PublicKey())
	assert.NotEqual(t, pk1.ExportedOffChainEncryption(), pk2.ExportedOffChainEncryption())
}

func TestOCRKeys_New(t *testing.T) {
	t.Parallel()
	pk1, err := ocrkey.NewV2()
	require.NoError(t, err)
	pk2, err := ocrkey.NewV2()
	require.NoError(t, err)
	pk3, err := ocrkey.NewV2()
	require.NoError(t, err)
	assertKeyBundlesNotEqual(t, pk1, pk2)
	assertKeyBundlesNotEqual(t, pk1, pk3)
	assertKeyBundlesNotEqual(t, pk2, pk3)
}
func TestOCRKeys_Raw_Key(t *testing.T) {
	t.Parallel()
	key := ocrkey.MustNewV2XXXTestingOnly(big.NewInt(1))
	require.Equal(t, key.ID(), key.Raw().Key().ID())
}
