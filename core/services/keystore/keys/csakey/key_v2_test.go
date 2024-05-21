package csakey

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSAKeyV2_RawPrivateKey(t *testing.T) {
	_, privKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	privateKey := Raw(privKey)

	assert.Equal(t, "<CSA Raw Private Key>", privateKey.String())
	assert.Equal(t, privateKey.String(), privateKey.GoString())
}

func TestCSAKeyV2_FromRawPrivateKey(t *testing.T) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	keyV2 := Raw(privKey).Key()

	assert.Equal(t, pubKey, keyV2.PublicKey)
	assert.Equal(t, privKey, *keyV2.privateKey)
	assert.Equal(t, keyV2.String(), keyV2.GoString())
	assert.Equal(t, hex.EncodeToString(pubKey), keyV2.PublicKeyString())
	assert.Equal(t, fmt.Sprintf("CSAKeyV2{PrivateKey: <redacted>, PublicKey: %s}", pubKey), keyV2.String())
}

func TestCSAKeyV2_NewV2(t *testing.T) {
	keyV2, err := NewV2()
	require.NoError(t, err)

	assert.Equal(t, 2, keyV2.Version)
	assert.NotNil(t, keyV2.PublicKey)
	assert.NotNil(t, keyV2.privateKey)
}
