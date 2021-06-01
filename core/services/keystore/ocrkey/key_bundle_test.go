package ocrkey_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"golang.org/x/crypto/curve25519"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keystore/ocrkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertKeyBundlesEqual(t *testing.T, pk1 *ocrkey.KeyBundle, pk2 *ocrkey.KeyBundle) {
	assert.Equal(t, pk1.ID, pk2.ID)
	assert.Equal(t, pk1.ExportedOnChainSigning().Curve, pk2.ExportedOnChainSigning().Curve)
	assert.Equal(t, pk1.ExportedOnChainSigning().X, pk2.ExportedOnChainSigning().X)
	assert.Equal(t, pk1.ExportedOnChainSigning().Y, pk2.ExportedOnChainSigning().Y)
	assert.Equal(t, pk1.ExportedOnChainSigning().D, pk2.ExportedOnChainSigning().D)
	assert.Equal(t, pk1.ExportedOffChainSigning(), pk2.ExportedOffChainSigning())
	assert.Equal(t, pk1.ExportedOffChainEncryption(), pk2.ExportedOffChainEncryption())
}

func assertKeyBundlesNotEqual(t *testing.T, pk1 *ocrkey.KeyBundle, pk2 *ocrkey.KeyBundle) {
	assert.NotEqual(t, pk1.ID, pk2.ID)
	assert.NotEqual(t, pk1.ExportedOnChainSigning().X, pk2.ExportedOnChainSigning().X)
	assert.NotEqual(t, pk1.ExportedOnChainSigning().Y, pk2.ExportedOnChainSigning().Y)
	assert.NotEqual(t, pk1.ExportedOnChainSigning().D, pk2.ExportedOnChainSigning().D)
	assert.NotEqual(t, pk1.ExportedOffChainSigning().PublicKey(), pk2.ExportedOffChainSigning().PublicKey())
	assert.NotEqual(t, pk1.ExportedOffChainEncryption(), pk2.ExportedOffChainEncryption())
}

func TestOCRKeys_NewKeyBundle(t *testing.T) {
	t.Parallel()
	pk1, err := ocrkey.NewKeyBundle()
	require.NoError(t, err)
	pk2, err := ocrkey.NewKeyBundle()
	require.NoError(t, err)
	pk3, err := ocrkey.NewKeyBundle()
	require.NoError(t, err)
	assertKeyBundlesNotEqual(t, pk1, pk2)
	assertKeyBundlesNotEqual(t, pk1, pk3)
	assertKeyBundlesNotEqual(t, pk2, pk3)
}

// TestOCRKeys_Encrypt_Decrypt tests that keys are identical after encrypting
// and then decrypting
func TestOCRKeys_Encrypt_Decrypt(t *testing.T) {
	t.Parallel()
	pk, err := ocrkey.NewKeyBundle()
	require.NoError(t, err)
	pkEncrypted, err := pk.Encrypt(cltest.Password, utils.FastScryptParams)
	require.NoError(t, err)
	// check that properties on encrypted key match those on OCRkey
	require.Equal(t, pk.ID, pkEncrypted.ID)
	require.Equal(t, ocrkey.OnChainSigningAddress(pk.PublicKeyAddressOnChain()), pkEncrypted.OnChainSigningAddress)
	require.Equal(t, ocrkey.OffChainPublicKey(pk.PublicKeyOffChain()), pkEncrypted.OffChainPublicKey)
	pkDecrypted, err := pkEncrypted.Decrypt(cltest.Password)
	require.NoError(t, err)
	assertKeyBundlesEqual(t, pk, pkDecrypted)
}

func TestOCRKeys_ScalarTooBig(t *testing.T) {
	t.Parallel()
	tooBig := new(big.Int)
	buf := make([]byte, curve25519.PointSize+1)
	buf[0] = 0x01
	tooBig.SetBytes(buf)
	kbr := ocrkey.KeyBundleRawData{
		EcdsaD: *tooBig,
	}
	jb, err := json.Marshal(&kbr)
	require.NoError(t, err)

	kb := ocrkey.KeyBundle{}
	err = kb.UnmarshalJSON(jb)
	assert.Equal(t, ocrkey.ErrScalarTooBig, errors.Cause(err))
}
