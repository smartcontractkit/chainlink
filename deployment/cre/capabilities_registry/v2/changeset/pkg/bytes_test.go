package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHexStringTo32Bytes(t *testing.T) {
	t.Run("valid hex string with 0x prefix", func(t *testing.T) {
		hexStr := "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
		expected := [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20}

		result, err := HexStringTo32Bytes(hexStr)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("valid hex string without 0x prefix", func(t *testing.T) {
		hexStr := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
		expected := [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20}

		result, err := HexStringTo32Bytes(hexStr)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("valid hex string with uppercase letters", func(t *testing.T) {
		hexStr := "0x0102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F20"
		expected := [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20}

		result, err := HexStringTo32Bytes(hexStr)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("all zeros", func(t *testing.T) {
		hexStr := "0x0000000000000000000000000000000000000000000000000000000000000000"
		expected := [32]byte{}

		result, err := HexStringTo32Bytes(hexStr)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("all ones", func(t *testing.T) {
		hexStr := "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
		expected := [32]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

		result, err := HexStringTo32Bytes(hexStr)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("error cases", func(t *testing.T) {
		testCases := []struct {
			name        string
			hexStr      string
			expectedErr string
		}{
			{
				name:        "too short",
				hexStr:      "0x0102030405060708090a0b0c0d0e0f",
				expectedErr: "invalid hex string length: expected 64 hex characters, got 30",
			},
			{
				name:        "too long",
				hexStr:      "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f2021",
				expectedErr: "invalid hex string length: expected 64 hex characters, got 66",
			},
			{
				name:        "empty string",
				hexStr:      "",
				expectedErr: "invalid hex string length: expected 64 hex characters, got 0",
			},
			{
				name:        "only 0x prefix",
				hexStr:      "0x",
				expectedErr: "invalid hex string length: expected 64 hex characters, got 0",
			},
			{
				name:        "invalid hex characters",
				hexStr:      "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fzz",
				expectedErr: "invalid hex string",
			},
			{
				name:        "contains non-hex characters",
				hexStr:      "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fgh",
				expectedErr: "invalid hex string",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := HexStringTo32Bytes(tc.hexStr)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
			})
		}
	})
}
