package ocrkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fastScryptParamsXXXTestingOnly = scryptParams{N: 2, P: 1}

func assertPrivateKeysEqual(t *testing.T, pk1 *KeyBundle, pk2 *KeyBundle) {
	assert.Equal(t, pk1.ID, pk2.ID)
	assert.Equal(t, pk1.onChainSigning.X, pk2.onChainSigning.X)
	assert.Equal(t, pk1.onChainSigning.Y, pk2.onChainSigning.Y)
	assert.Equal(t, pk1.onChainSigning.D, pk2.onChainSigning.D)
	assert.Equal(t, pk1.offChainSigning, pk2.offChainSigning)
	assert.Equal(t, pk1.offChainEncryption, pk2.offChainEncryption)
}

func assertPrivateKeysNotEqual(t *testing.T, pk1 *KeyBundle, pk2 *KeyBundle) {
	assert.NotEqual(t, pk1.ID, pk2.ID)
	assert.NotEqual(t, pk1.onChainSigning.X, pk2.onChainSigning.X)
	assert.NotEqual(t, pk1.onChainSigning.Y, pk2.onChainSigning.Y)
	assert.NotEqual(t, pk1.onChainSigning.D, pk2.onChainSigning.D)
	assert.NotEqual(t, pk1.offChainSigning.PublicKey(), pk2.offChainSigning.PublicKey())
	assert.NotEqual(t, pk1.offChainEncryption, pk2.offChainEncryption)
}

func TestOCRKeys_NewKeyBundle(t *testing.T) {
	t.Parallel()
	pk1, err := NewKeyBundle()
	require.NoError(t, err)
	pk2, err := NewKeyBundle()
	require.NoError(t, err)
	pk3, err := NewKeyBundle()
	require.NoError(t, err)
	assertPrivateKeysNotEqual(t, pk1, pk2)
	assertPrivateKeysNotEqual(t, pk1, pk3)
	assertPrivateKeysNotEqual(t, pk2, pk3)
}

// TestOCRKeys_Encrypt_Decrypt tests that keys are identical after encrypting
// and then decrypting
func TestOCRKeys_Encrypt_Decrypt(t *testing.T) {
	t.Parallel()
	pk, err := NewKeyBundle()
	require.NoError(t, err)
	pkEncrypted, err := pk.encrypt("password", fastScryptParamsXXXTestingOnly)
	require.NoError(t, err)
	// check that properties on encrypted key match those on OCRkey
	require.Equal(t, pk.ID, pkEncrypted.ID)
	require.Equal(t, pk.onChainSigning.Address(), pkEncrypted.OnChainSigningAddress)
	require.Equal(t, pk.offChainSigning.PublicKey(), pkEncrypted.OffChainPublicKey)
	pkDecrypted, err := pkEncrypted.Decrypt("password")
	require.NoError(t, err)
	assertPrivateKeysEqual(t, pk, pkDecrypted)
}
