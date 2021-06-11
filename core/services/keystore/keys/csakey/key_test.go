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
