package csakey

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

func TestCSAKeyV2_FromRawPrivateKey(t *testing.T) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	keyV2 := KeyFor(internal.NewRaw(privKey))

	assert.Equal(t, pubKey, keyV2.PublicKey)
	assert.Equal(t, []byte(privKey), internal.Bytes(keyV2.raw))
	assert.Equal(t, hex.EncodeToString(pubKey), keyV2.PublicKeyString())
}

func TestCSAKeyV2_NewV2(t *testing.T) {
	keyV2, err := NewV2()
	require.NoError(t, err)

	assert.Equal(t, 2, keyV2.Version)
	assert.NotNil(t, keyV2.PublicKey)
	assert.NotNil(t, keyV2.raw)
}
