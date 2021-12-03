package ocr2key_test

import (
	"testing"

	"golang.org/x/crypto/curve25519"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertKeyBundlesEqual(t *testing.T, pk1 *ocr2key.KeyBundle, pk2 *ocr2key.KeyBundle) {
	assert.Equal(t, pk1.ID(), pk2.ID())
	assert.EqualValues(t, pk1.OffchainKeyring, pk2.OffchainKeyring)
	assert.EqualValues(t, pk1.OnchainKeyring, pk2.OnchainKeyring)
}

func assertKeyBundlesNotEqual(t *testing.T, pk1 *ocr2key.KeyBundle, pk2 *ocr2key.KeyBundle) {
	assert.NotEqual(t, pk1.ID(), pk2.ID())
	assert.NotEqualValues(t, pk1.OffchainKeyring, pk2.OffchainKeyring)
	assert.NotEqualValues(t, pk1.OnchainKeyring, pk2.OnchainKeyring)
}

func TestOCR2keys_NewKeyBundle(t *testing.T) {
	t.Parallel()
	pk1, err := ocr2key.NewKeyBundle()
	require.NoError(t, err)
	pk2, err := ocr2key.NewKeyBundle()
	require.NoError(t, err)
	pk3, err := ocr2key.NewKeyBundle()
	require.NoError(t, err)
	assertKeyBundlesNotEqual(t, pk1, pk2)
	assertKeyBundlesNotEqual(t, pk1, pk3)
	assertKeyBundlesNotEqual(t, pk2, pk3)
}

// TestOCR2Keys_Encrypt_Decrypt tests that keys are identical after encrypting
// and then decrypting
func TestOCR2Keys_Encrypt_Decrypt(t *testing.T) {
	t.Parallel()
	pk, err := ocr2key.NewKeyBundle()
	require.NoError(t, err)
	pkEncrypted, err := pk.Encrypt(cltest.Password, utils.FastScryptParams)
	require.NoError(t, err)
	// check that properties on encrypted key match those on ocr2key
	require.Equal(t, pk.ID(), pkEncrypted.ID.String())
	assert.EqualValues(t, pk.PublicKeyAddressOnChain(), pkEncrypted.OnchainSigningAddress)
	assert.EqualValues(t, pk.PublicKeyOffChain(), pkEncrypted.OffchainSigningPublicKey)
	// XXX: EqualValues will treat [32]byte{} and []byte{... 32 items} as different?
	encryptedPublicKeyConfig := [curve25519.PointSize]byte{}
	copy(encryptedPublicKeyConfig[:], pkEncrypted.OffchainEncryptionPublicKey)
	assert.EqualValues(t, pk.PublicKeyConfig(), encryptedPublicKeyConfig)
	pkDecrypted, err := pkEncrypted.Decrypt(cltest.Password)
	require.NoError(t, err)
	assertKeyBundlesEqual(t, pk, pkDecrypted)
}
