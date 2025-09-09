package tonkey

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

func TestTONKey(t *testing.T) {
	t.Run("Generate new key and verify its components", func(t *testing.T) {
		// Generate a new key
		key, err := New()
		require.NoError(t, err, "Failed to generate new TONKey")

		// Verify key components
		assert.NotNil(t, key.pubKey, "Public key should not be nil")
		assert.NotNil(t, key.raw, "Private key should not be nil")

		addressBase64 := key.AddressBase64()

		assert.Len(t, addressBase64, 48, "Address in base64 should be 48 chars")

		decodedAddress, err := base64.RawURLEncoding.DecodeString(addressBase64)
		if err != nil {
			require.NoError(t, err, "Failed to decode addressBase64")
		}

		assert.Len(t, decodedAddress, 36, "Decoded address should be 36 bytes")
	})

	t.Run("Generate new key and verify its components", func(t *testing.T) {
		// Generate a new key
		key, err := New()
		require.NoError(t, err, "Failed to generate new TONKey")

		// Verify key components
		assert.NotNil(t, key.pubKey, "Public key should not be nil")
		assert.NotNil(t, key.raw, "Private key should not be nil")
	})

	t.Run("Signature is valid and verifiable", func(t *testing.T) {
		key, err := New()
		require.NoError(t, err)

		msg := []byte("test message")
		sig, err := key.Sign(msg)
		require.NoError(t, err)

		valid := ed25519.Verify(key.GetPublic(), msg, sig)
		assert.True(t, valid, "Signature should be valid")
	})

	t.Run("Public key string encoding is consistent", func(t *testing.T) {
		key, err := New()
		require.NoError(t, err)

		pubKeyStr := key.PublicKeyStr()
		decoded, err := hex.DecodeString(pubKeyStr)
		require.NoError(t, err)

		assert.Equal(t, []byte(key.GetPublic()), decoded, "Decoded public key should match original")
	})

	t.Run("MustNewInsecure should not panic on valid reader", func(t *testing.T) {
		assert.NotPanics(t, func() {
			MustNewInsecure(strings.NewReader("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ012345")) // 64 bytes
		})
	})

	t.Run("Custom wallet version and workchain gives valid address", func(t *testing.T) {
		key, err := New()
		require.NoError(t, err)

		addr, err := key.PubkeyToAddressWith(wallet.V4R1, -1)
		require.NoError(t, err)
		assert.NotNil(t, addr)

		// Check that the base64 address is decodable and has expected length
		decoded, err := base64.RawURLEncoding.DecodeString(addr.String())
		require.NoError(t, err, "base64 address should be valid base64")
		assert.Len(t, decoded, 36, "Decoded address should be 36 bytes")
	})

	t.Run("Same key produces different addresses for different wallet versions and workchains", func(t *testing.T) {
		key, err := New()
		require.NoError(t, err)

		addrV3WC0, err := key.PubkeyToAddressWith(wallet.V3, 0)
		require.NoError(t, err)

		addrV4WC0, err := key.PubkeyToAddressWith(wallet.V4R1, 0)
		require.NoError(t, err)

		addrV3WC1, err := key.PubkeyToAddressWith(wallet.V3, 1)
		require.NoError(t, err)

		addrV3WCMinus1, err := key.PubkeyToAddressWith(wallet.V3, -1)
		require.NoError(t, err)

		// Assert addresses are all different
		assert.NotEqual(t, addrV3WC0.String(), addrV4WC0.String(), "V3 and V4 addresses should differ")
		assert.NotEqual(t, addrV3WC0.String(), addrV3WC1.String(), "Workchain 0 and 1 addresses should differ")
		assert.NotEqual(t, addrV3WC0.String(), addrV3WCMinus1.String(), "Workchain 0 and -1 addresses should differ")
	})

	t.Run("RawAddress returns correct non-base64 address", func(t *testing.T) {
		key, err := New()
		require.NoError(t, err)

		rawAddress := key.RawAddress()
		addressBase64 := key.AddressBase64()

		assert.NotEmpty(t, rawAddress, "Raw address should not be empty")
		assert.NotEqual(t, addressBase64, rawAddress, "Raw and base64 addresses should differ")

		// Parse raw address and validate
		parsed, err := address.ParseRawAddr(rawAddress)
		require.NoError(t, err, "Raw address should be parseable")
		assert.Equal(t, rawAddress, parsed.StringRaw(), "Parsed raw address should match original")

		// Default workchain is 0
		assert.True(t, strings.HasPrefix(rawAddress, "0:"), "Expected raw address to start with '0:' for default workchain")
	})

	t.Run("KeyFor and Raw produce consistent keys", func(t *testing.T) {
		key1, err := New()
		require.NoError(t, err)

		raw := key1.Raw()
		key2 := KeyFor(raw)

		assert.Equal(t, key1.PublicKeyStr(), key2.PublicKeyStr(), "Restored key should have same public key")
		assert.Equal(t, key1.AddressBase64(), key2.AddressBase64(), "Restored key should have same address")
	})
}
