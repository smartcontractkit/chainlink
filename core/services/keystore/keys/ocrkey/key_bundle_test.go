package ocrkey_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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

func TestOCRKeys_NewBundleIDMatchesOld(t *testing.T) {
	t.Parallel()
	oldKey, err := ocrkey.New()
	require.NoError(t, err)
	newKey := oldKey.ToV2()
	require.Equal(t, oldKey.ID.String(), newKey.ID())
}

func TestOCRKeys_Raw_Key(t *testing.T) {
	t.Parallel()
	key := ocrkey.MustNewV2XXXTestingOnly(big.NewInt(1))
	require.Equal(t, key.ID(), key.Raw().Key().ID())
}

func TestOCRKeys_BundleSetID(t *testing.T) {
	t.Parallel()

	k, err := ocrkey.New()
	require.NoError(t, err)
	ek, err := k.Encrypt("test", utils.FastScryptParams)
	require.NoError(t, err)

	oldId := ek.GetID()
	err = ek.SetID("48656c6c6f20476f7068657221")
	require.NoError(t, err)

	assert.NotEqual(t, oldId, ek.GetID())

	err = ek.SetID("invalid id")
	assert.Error(t, err)
}

func TestOCRKeys_BundleDecrypt(t *testing.T) {
	t.Parallel()

	k, err := ocrkey.New()
	require.NoError(t, err)
	ek, err := k.Encrypt("test", utils.FastScryptParams)
	require.NoError(t, err)

	_, err = ek.Decrypt("wrongpass")
	assert.Error(t, err)

	dk, err := ek.Decrypt("test")
	require.NoError(t, err)

	dk.GoString()
	assert.Equal(t, k.GoString(), dk.GoString())
	assert.Equal(t, k.ID.String(), dk.ID.String())
}

func TestOCRKeys_BundleMarshalling(t *testing.T) {
	t.Parallel()

	k, err := ocrkey.New()
	require.NoError(t, err)
	k2, err := ocrkey.New()
	require.NoError(t, err)

	mk, err := k.MarshalJSON()
	require.NoError(t, err)

	err = k2.UnmarshalJSON(mk)
	require.NoError(t, err)

	assert.Equal(t, k.String(), k2.String())
}
