package ocrkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertPrivateKeysEqual(t *testing.T, pk1 *OCRPrivateKeys, pk2 *OCRPrivateKeys) {
	assert.Equal(t, pk1.onChainSigning.X, pk2.onChainSigning.X)
	assert.Equal(t, pk1.onChainSigning.Y, pk2.onChainSigning.Y)
	assert.Equal(t, pk1.onChainSigning.D, pk2.onChainSigning.D)
	assert.Equal(t, pk1.offChainSigning.PublicKey(), pk2.offChainSigning.PublicKey())
	assert.Equal(t, pk1.offChainEncryption, pk2.offChainEncryption)
}

func assertPrivateKeysNotEqual(t *testing.T, pk1 *OCRPrivateKeys, pk2 *OCRPrivateKeys) {
	assert.NotEqual(t, pk1.onChainSigning.X, pk2.onChainSigning.X)
	assert.NotEqual(t, pk1.onChainSigning.Y, pk2.onChainSigning.Y)
	assert.NotEqual(t, pk1.onChainSigning.D, pk2.onChainSigning.D)
	assert.NotEqual(t, pk1.offChainSigning.PublicKey(), pk2.offChainSigning.PublicKey())
	assert.NotEqual(t, pk1.offChainEncryption, pk2.offChainEncryption)
}

// Tests that NewDeterministicOCRPrivateKeysXXXTestingOnly creates deterministic
// OCRPrivateKeys
func TestOCRKeys_NewDeterministicOCRPrivateKeysXXXTestingOnly(t *testing.T) {
	t.Parallel()
	pk := NewDeterministicOCRPrivateKeysXXXTestingOnly(1)
	pkSameSeed := NewDeterministicOCRPrivateKeysXXXTestingOnly(1)
	pkDifferentSeed := NewDeterministicOCRPrivateKeysXXXTestingOnly(2)
	assertPrivateKeysEqual(t, pk, pkSameSeed)
	assertPrivateKeysNotEqual(t, pk, pkDifferentSeed)
}

func TestOCRKeys_NewOCRPrivateKeys(t *testing.T) {
	t.Parallel()
	pk1 := NewOCRPrivateKeys()
	pk2 := NewOCRPrivateKeys()
	pk3 := NewOCRPrivateKeys()
	assertPrivateKeysNotEqual(t, pk1, pk2)
	assertPrivateKeysNotEqual(t, pk1, pk3)
	assertPrivateKeysNotEqual(t, pk2, pk3)
}

func TestOCRKeys_Encrypt_Decrypt(t *testing.T) {
	t.Parallel()
	pkEncrypted := NewDeterministicOCRPrivateKeysXXXTestingOnly(1)
	encryptedPKs, err := pkEncrypted.Encrypt("password")
	require.NoError(t, err)
	pkDecrypted, err := encryptedPKs.Decrypt("password")
	require.NoError(t, err)
	assertPrivateKeysEqual(t, pkEncrypted, pkDecrypted)
}
