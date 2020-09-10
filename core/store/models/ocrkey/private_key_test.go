package ocrkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOCRKeys_Encrypt_Decrypt(t *testing.T) {
	pk := OCRPrivateKeys{
		onChainSignging:    nil,
		offChainSigning:    nil,
		offChainEncryption: nil,
	}
	encryptedPKs, err := pk.Encrypt("password")
	require.NoError(t, err)
	pk2, err := encryptedPKs.Decrypt("password")
	require.NoError(t, err)
	assert.Equal(t, pk.onChainSignging.D, pk2.onChainSignging.D)
	assert.Equal(t, pk.onChainSignging.X, pk2.onChainSignging.X)
	assert.Equal(t, pk.onChainSignging.Y, pk2.onChainSignging.Y)
	assert.Equal(t, pk.offChainSigning.PublicKey, pk2.offChainSigning.PublicKey)
	assert.Equal(t, pk.offChainEncryption, pk2.offChainEncryption)
}
