package ocrkey_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertKeyBundlesNotEqual(t *testing.T, pk1 ocrkey.KeyV2, pk2 ocrkey.KeyV2) {
	assert.NotEqual(t, pk1.ID(), pk2.ID())
	assert.NotEqual(t, pk1.ExportedOnChainSigning().X, pk2.ExportedOnChainSigning().X)
	assert.NotEqual(t, pk1.ExportedOnChainSigning().Y, pk2.ExportedOnChainSigning().Y)
	assert.NotEqual(t, pk1.ExportedOnChainSigning().D, pk2.ExportedOnChainSigning().D)
	assert.NotEqual(t, pk1.ExportedOffChainSigning().PublicKey(), pk2.ExportedOffChainSigning().PublicKey())
	assert.NotEqual(t, pk1.ExportedOffChainEncryption(), pk2.ExportedOffChainEncryption())
}

func TestOCRKeys_NewKeyBundle(t *testing.T) {
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
