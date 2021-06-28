package csakey

import (
	"crypto/ed25519"
	"testing"

	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	passphrase := "passphrase"
	key, err := New(passphrase, utils.FastScryptParams)
	require.NoError(t, err)

	rawprivkey, err := key.EncryptedPrivateKey.Decrypt("passphrase")
	require.NoError(t, err)

	privkey := ed25519.PrivateKey(rawprivkey)
	assert.Equal(t, ed25519.PublicKey(key.PublicKey), privkey.Public())
}

func Test_Unlock(t *testing.T) {
	passphrase := "passphrase"
	key, err := New(passphrase, utils.FastScryptParams)
	require.NoError(t, err)

	err = key.Unlock(passphrase)
	require.NoError(t, err)

	expected, err := key.EncryptedPrivateKey.Decrypt(passphrase)
	require.NoError(t, err)

	assert.Equal(t, expected, key.privateKey)
}

func Test_GetPrivateKey(t *testing.T) {
	passphrase := "passphrase"
	key, err := New(passphrase, utils.FastScryptParams)
	require.NoError(t, err)

	privkey, err := key.Unsafe_GetPrivateKey()
	require.NoError(t, err)
	assert.Equal(t, key.privateKey, privkey)
}
