package dkgrecipientkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

func TestNew(t *testing.T) {
	key, err := New()
	require.NoError(t, err)

	assert.NotNil(t, key.PublicKey())
	assert.NotNil(t, key.Raw())
}

func TestPublicKey(t *testing.T) {
	key, err := New()
	require.NoError(t, err)

	pubKey := key.PublicKey()
	assert.NotNil(t, pubKey)
	assert.Len(t, pubKey, 33) // P256 public key should be 33 bytes
}

func TestPublicKeyString(t *testing.T) {
	key, err := New()
	require.NoError(t, err)

	pubKeyStr := key.PublicKeyString()
	assert.NotEmpty(t, pubKeyStr)
	assert.Len(t, pubKeyStr, 66) // 33 bytes * 2 hex chars per byte
}

func TestECDH(t *testing.T) {
	// Create two keys
	key1, err := New()
	require.NoError(t, err)

	key2, err := New()
	require.NoError(t, err)

	// Perform ECDH key exchange
	secret1, err := key1.ECDH(key2.PublicKey())
	require.NoError(t, err)

	secret2, err := key2.ECDH(key1.PublicKey())
	require.NoError(t, err)

	// The shared secrets should be the same
	assert.Equal(t, secret1, secret2)
	assert.NotEmpty(t, secret1)
}

func TestRaw(t *testing.T) {
	key, err := New()
	require.NoError(t, err)

	raw := key.Raw()
	assert.NotNil(t, raw)

	// Raw should contain some data
	rawBytes := internal.Bytes(raw)
	assert.NotEmpty(t, rawBytes)
}

func TestKeyUniqueness(t *testing.T) {
	// Create multiple keys and ensure they're all different
	key1, err := New()
	require.NoError(t, err)

	key2, err := New()
	require.NoError(t, err)

	key3, err := New()
	require.NoError(t, err)

	// All public keys should be different
	assert.NotEqual(t, key1.PublicKey(), key2.PublicKey())
	assert.NotEqual(t, key1.PublicKey(), key3.PublicKey())
	assert.NotEqual(t, key2.PublicKey(), key3.PublicKey())

	// All public key strings should be different
	assert.NotEqual(t, key1.PublicKeyString(), key2.PublicKeyString())
	assert.NotEqual(t, key1.PublicKeyString(), key3.PublicKeyString())
	assert.NotEqual(t, key2.PublicKeyString(), key3.PublicKeyString())
}
