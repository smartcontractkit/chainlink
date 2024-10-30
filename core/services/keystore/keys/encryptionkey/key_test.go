package encryptionkey

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	key, err := New()
	require.NoError(t, err)

	assert.NotNil(t, key.PublicKey)
	assert.NotNil(t, key.privateKey)
}

func TestPublicKey(t *testing.T) {
	key, err := New()
	require.NoError(t, err)

	assert.Equal(t, *key.publicKey, key.PublicKey())
}

func TestPublicKeyString(t *testing.T) {
	key := "my-test-public-key"
	var pubkey [32]byte
	copy(pubkey[:], key)
	k := Key{
		publicKey: &pubkey,
	}

	expected := hex.EncodeToString([]byte(key))
	// given the key is a [32]byte we need to ensure the encoded string is 64 character long
	for len(expected) < 64 {
		expected += "0"
	}

	assert.Equal(t, expected, k.PublicKeyString())
}

func TestDecrypt(t *testing.T) {
	key, err := New()
	require.NoError(t, err)

	secret := []byte("my-secret")

	ciphertext, err := key.Encrypt(secret)
	require.NoError(t, err)

	plaintext, err := key.Decrypt(ciphertext)
	require.NoError(t, err)

	assert.Equal(t, secret, plaintext)
}
