package workflowkey

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/nacl/box"
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

func TestEncryptKeyRawPrivateKey(t *testing.T) {
	privKey, err := New()
	require.NoError(t, err)

	privateKey := privKey.Raw()

	assert.Equal(t, "<Workflow Raw Private Key>", privateKey.String())
	assert.Equal(t, privateKey.String(), privateKey.GoString())
}

func TestEncryptKeyFromRawPrivateKey(t *testing.T) {
	boxPubKey, boxPrivKey, err := box.GenerateKey(cryptorand.Reader)
	require.NoError(t, err)

	privKey := make([]byte, 32)
	copy(privKey, boxPrivKey[:])
	key := Raw(privKey).Key()

	assert.Equal(t, boxPubKey, key.publicKey)
	assert.Equal(t, boxPrivKey, key.privateKey)
	assert.Equal(t, key.String(), key.GoString())

	byteBoxPubKey := make([]byte, 32)
	copy(byteBoxPubKey, boxPubKey[:])

	assert.Equal(t, hex.EncodeToString(byteBoxPubKey), key.PublicKeyString())
	assert.Equal(t, fmt.Sprintf("WorkflowKey{PrivateKey: <redacted>, PublicKey: %s}", byteBoxPubKey), key.String())
}

func TestPublicKeyStringAndID(t *testing.T) {
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
	assert.Equal(t, expected, k.ID())
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
