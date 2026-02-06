package fakes

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	confidentialhttp "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialhttp"
)

func TestAESGCMEncryption(t *testing.T) {
	plaintext := []byte("Hello, World! This is a test message.")

	encrypted, err := aesGCMEncrypt(plaintext, FakeAESGCMEncryptionKey)
	require.NoError(t, err)
	assert.NotEqual(t, plaintext, encrypted)
	assert.Greater(t, len(encrypted), len(plaintext)) // nonce + ciphertext + tag

	// Verify we can decrypt it
	decrypted, err := aesGCMDecrypt(encrypted, FakeAESGCMEncryptionKey)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestTDH2Encryption(t *testing.T) {
	plaintext := []byte("Hello, World! This is a TDH2 test message.")

	// Encrypt using the fake TDH2 public key
	encrypted, err := tdh2Encrypt(plaintext)
	require.NoError(t, err)
	assert.NotEqual(t, plaintext, encrypted)

	// Decrypt using the helper function
	decrypted, err := DecryptFakeTDH2Ciphertext(encrypted)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestTDH2KeysInitialization(t *testing.T) {
	// Get public key
	pubKey, err := GetFakeTDH2PublicKey()
	require.NoError(t, err)
	require.NotNil(t, pubKey)

	// Get private shares
	privShares, err := GetFakeTDH2PrivateShares()
	require.NoError(t, err)
	require.Len(t, privShares, FakeTDH2TotalShares)

	// Verify keys are consistent across calls (lazy init)
	pubKey2, err := GetFakeTDH2PublicKey()
	require.NoError(t, err)
	assert.Equal(t, pubKey, pubKey2)
}

func TestHasEncryptionSecret(t *testing.T) {
	t.Run("returns true when magic key exists", func(t *testing.T) {
		secrets := []*confidentialhttp.SecretIdentifier{
			{Key: "other-key"},
			{Key: AESGCMEncryptionKeyName},
		}
		assert.True(t, hasEncryptionSecret(secrets))
	})

	t.Run("returns false when magic key does not exist", func(t *testing.T) {
		secrets := []*confidentialhttp.SecretIdentifier{
			{Key: "other-key"},
			{Key: "another-key"},
		}
		assert.False(t, hasEncryptionSecret(secrets))
	})

	t.Run("returns false for empty secrets", func(t *testing.T) {
		assert.False(t, hasEncryptionSecret(nil))
		assert.False(t, hasEncryptionSecret([]*confidentialhttp.SecretIdentifier{}))
	})
}

func TestTDH2EncryptionRoundTrip(t *testing.T) {
	testCases := []struct {
		name      string
		plaintext []byte
	}{
		{"short message", []byte("Hi")},
		{"medium message", []byte("This is a medium length test message for TDH2 encryption.")},
		{"binary data", []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := tdh2Encrypt(tc.plaintext)
			require.NoError(t, err)

			decrypted, err := DecryptFakeTDH2Ciphertext(encrypted)
			require.NoError(t, err)
			assert.Equal(t, tc.plaintext, decrypted)
		})
	}
}

func TestTDH2EncryptionEmptyMessage(t *testing.T) {
	// Empty messages are a special case - TDH2 may return nil instead of empty slice
	encrypted, err := tdh2Encrypt([]byte{})
	require.NoError(t, err)

	decrypted, err := DecryptFakeTDH2Ciphertext(encrypted)
	require.NoError(t, err)
	assert.Empty(t, decrypted) // Use Empty instead of Equal to handle nil vs []byte{}
}

// aesGCMDecrypt is a helper for testing - decrypts AES-GCM encrypted data
func aesGCMDecrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
