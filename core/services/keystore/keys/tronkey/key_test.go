package tronkey

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTronKeyNewKeyGeneration(t *testing.T) {
	t.Run("Generate new key and verify its components", func(t *testing.T) {
		// Generate a new key
		key, err := New()
		require.NoError(t, err, "Failed to generate new TronKey")

		// Verify key components
		assert.NotNil(t, key.pubKey, "Public key should not be nil")
		assert.NotNil(t, key.raw, "Private key should not be nil")
	})

	t.Run("Multiple key generations produce unique keys", func(t *testing.T) {
		key1, err := New()
		require.NoError(t, err, "Failed to generate first key")

		key2, err := New()
		require.NoError(t, err, "Failed to generate second key")

		assert.NotEqual(t, key1.raw, key2.raw, "Generated private keys should be unique")
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

		key := Key{pubKey: &privateKey.PublicKey}

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
