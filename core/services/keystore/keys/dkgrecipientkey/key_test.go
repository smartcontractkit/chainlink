package dkgrecipientkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

func TestUnit_New(t *testing.T) {
	key, err := New()
	require.NoError(t, err)

	assert.NotNil(t, key.PublicKey())
	assert.NotNil(t, key.Raw())
	assert.Equal(t, key.ID(), key.PublicKeyString())
}

func TestUnit_PublicKey(t *testing.T) {
	key, err := New()
	require.NoError(t, err)

	pubKey := key.PublicKey()
	assert.NotNil(t, pubKey)
	assert.Len(t, pubKey, 33)
}

func TestUnit_PublicKeyString(t *testing.T) {
	key, err := New()
	require.NoError(t, err)

	pubKeyStr := key.PublicKeyString()
	assert.NotEmpty(t, pubKeyStr)
	assert.Len(t, pubKeyStr, 66)
}

func TestUnit_ECDH(t *testing.T) {
	key1, err := New()
	require.NoError(t, err)

	key2, err := New()
	require.NoError(t, err)

	secret1, err := key1.ECDH(key2.PublicKey())
	require.NoError(t, err)

	secret2, err := key2.ECDH(key1.PublicKey())
	require.NoError(t, err)

	assert.Equal(t, secret1, secret2)
	assert.NotEmpty(t, secret1)
}

func TestUnit_Raw(t *testing.T) {
	key, err := New()
	require.NoError(t, err)

	raw := key.Raw()
	assert.NotNil(t, raw)

	rawBytes := internal.Bytes(raw)
	assert.NotEmpty(t, rawBytes)

	key2 := KeyFor(raw)
	assert.Equal(t, key.PublicKeyString(), key2.PublicKeyString())
}

func TestUnit_KeyUniqueness(t *testing.T) {
	key1, err := New()
	require.NoError(t, err)

	key2, err := New()
	require.NoError(t, err)

	key3, err := New()
	require.NoError(t, err)

	assert.NotEqual(t, key1.PublicKey(), key2.PublicKey())
	assert.NotEqual(t, key1.PublicKey(), key3.PublicKey())
	assert.NotEqual(t, key2.PublicKey(), key3.PublicKey())

	assert.NotEqual(t, key1.PublicKeyString(), key2.PublicKeyString())
	assert.NotEqual(t, key1.PublicKeyString(), key3.PublicKeyString())
	assert.NotEqual(t, key2.PublicKeyString(), key3.PublicKeyString())
}
