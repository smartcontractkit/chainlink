package workflowkey

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/nacl/box"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

func TestNew(t *testing.T) {
	key, err := New()
	require.NoError(t, err)

	assert.NotNil(t, key.PublicKey)
	assert.NotNil(t, key.raw)
}

func TestPublicKey(t *testing.T) {
	key, err := New()
	require.NoError(t, err)

	assert.Equal(t, *key.publicKey, key.PublicKey())
}

func TestEncryptKeyFromRawPrivateKey(t *testing.T) {
	boxPubKey, boxPrivKey, err := box.GenerateKey(cryptorand.Reader)
	require.NoError(t, err)

	privKey := make([]byte, 32)
	copy(privKey, boxPrivKey[:])
	key := KeyFor(internal.NewRaw(privKey))

	assert.Equal(t, boxPubKey, key.publicKey)
	assert.Equal(t, boxPrivKey[:], internal.Bytes(key.raw))

	byteBoxPubKey := make([]byte, 32)
	copy(byteBoxPubKey, boxPubKey[:])

	assert.Equal(t, hex.EncodeToString(byteBoxPubKey), key.PublicKeyString())
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

func TestMustNewXXXTestingOnly(t *testing.T) {
	tests := []struct {
		name        string
		k           *big.Int
		wantSuccess bool
	}{
		{
			name:        "generates valid key from big.Int",
			k:           big.NewInt(1),
			wantSuccess: true,
		},
		{
			name:        "panics on nil input",
			k:           nil,
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantSuccess {
				require.Panics(t, func() { MustNewXXXTestingOnly(tt.k) })
				return
			}

			key := MustNewXXXTestingOnly(tt.k)
			require.NotNil(t, key.raw)
			require.NotNil(t, key.publicKey)

			// Verify key generation is deterministic
			if tt.k.Cmp(big.NewInt(1)) != 0 {
				key1 := MustNewXXXTestingOnly(tt.k)
				require.Equal(t, key1.PublicKey(), key.PublicKey())
			}
		})
	}
}
