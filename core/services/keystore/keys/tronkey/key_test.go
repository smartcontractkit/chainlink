package tronkey

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTronKeyRawPrivateKey(t *testing.T) {
	t.Run("Create from raw bytes and check string representation", func(t *testing.T) {
		// Generate a private key
		privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
		require.NoError(t, err, "Failed to generate ECDSA key")

		// Create TronKey from raw bytes
		tronKey := Raw(privateKeyECDSA.D.Bytes())

		// Check string representation
		expectedStr := "<Tron Raw Private Key>"
		assert.Equal(t, expectedStr, tronKey.String(), "Unexpected string representation")
		assert.Equal(t, expectedStr, tronKey.GoString(), "String() and GoString() should return the same value")
	})
}

func TestTronKeyNewKeyGeneration(t *testing.T) {
	t.Run("Generate new key and verify its components", func(t *testing.T) {
		// Generate a new key
		key, err := New()
		require.NoError(t, err, "Failed to generate new TronKey")

		// Verify key components
		assert.NotNil(t, key.pubKey, "Public key should not be nil")
		assert.NotNil(t, key.privKey, "Private key should not be nil")
	})

	t.Run("Multiple key generations produce unique keys", func(t *testing.T) {
		key1, err := New()
		require.NoError(t, err, "Failed to generate first key")

		key2, err := New()
		require.NoError(t, err, "Failed to generate second key")

		assert.NotEqual(t, key1.privKey, key2.privKey, "Generated private keys should be unique")
		assert.NotEqual(t, key1.pubKey, key2.pubKey, "Generated public keys should be unique")
	})
}

func TestKeyAddress(t *testing.T) {
	t.Run("Known private key and expected address", func(t *testing.T) {
		// Tests cases from https://developers.tron.network/docs/account
		privateKeyHex := "b406adb115b43e103c7b1dc8b5931f63279a5b6b2cf7328638814c43171a2908"
		expectedAddress := "TDdcf5iMDkB61oGM27TNak55eVX214thBG"

		privateKeyBytes, err := hex.DecodeString(privateKeyHex)
		require.NoError(t, err, "Failed to decode private key hex")

		privateKey, err := crypto.ToECDSA(privateKeyBytes)
		require.NoError(t, err, "Failed to convert private key to ECDSA")

		key := Key{
			privKey: privateKey,
			pubKey:  &privateKey.PublicKey,
		}
		require.NotNil(t, key.privKey, "Private key is nil")

		address := key.Base58Address()
		require.Equal(t, expectedAddress, address, "Generated address does not match expected address")
	})

	t.Run("Generate new key and check address format", func(t *testing.T) {
		newKey, err := New()
		if err != nil {
			t.Fatalf("Failed to generate new key: %v", err)
		}

		newAddress := newKey.Base58Address()
		isValid := isValidBase58Address(newAddress)
		require.True(t, isValid, "Generated address is not valid")
	})
}
