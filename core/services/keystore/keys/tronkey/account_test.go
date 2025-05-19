package tronkey

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DecodeBase58(t *testing.T) {
	invalidAddresses := []string{
		"TronEnergyioE1Z3ukeRv38sYkv5Jn55bL",
		"TronEnergyioNijNo8g3LF2ABKUAae6D2Z",
		"TronEnergyio3ZMcXA5hSjrTxaioKGgqyr",
	}

	validAddresses := []string{
		"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
		"TVj7RNVHy6thbM7BWdSe9G6gXwKhjhdNZS",
		"THPvaUhoh2Qn2y9THCZML3H815hhFhn5YC",
	}

	for _, addr := range invalidAddresses {
		_, err := DecodeCheck(addr)
		require.Error(t, err)
	}

	for _, addr := range validAddresses {
		_, err := DecodeCheck(addr)
		require.NoError(t, err)
	}
}

func TestAddress(t *testing.T) {
	t.Run("Valid Addresses", func(t *testing.T) {
		validAddresses := []string{
			"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
			"TVj7RNVHy6thbM7BWdSe9G6gXwKhjhdNZS",
			"THPvaUhoh2Qn2y9THCZML3H815hhFhn5YC",
		}

		for _, addrStr := range validAddresses {
			t.Run(addrStr, func(t *testing.T) {
				addr, err := Base58ToAddress(addrStr)
				require.NoError(t, err)
				require.Equal(t, addrStr, addr.String())

				decoded, err := DecodeCheck(addrStr)
				require.NoError(t, err)
				require.True(t, bytes.Equal(decoded, addr.Bytes()))
			})
		}
	})

	t.Run("Invalid Addresses", func(t *testing.T) {
		invalidAddresses := []string{
			"TronEnergyioE1Z3ukeRv38sYkv5Jn55bL",
			"TronEnergyioNijNo8g3LF2ABKUAae6D2Z",
			"TronEnergyio3ZMcXA5hSjrTxaioKGgqyr",
		}

		for _, addrStr := range invalidAddresses {
			t.Run(addrStr, func(t *testing.T) {
				_, err := Base58ToAddress(addrStr)
				require.Error(t, err)

				_, err = DecodeCheck(addrStr)
				require.Error(t, err)
			})
		}
	})

	t.Run("Address Conversion", func(t *testing.T) {
		addrStr := "TSvT6Bg3siokv3dbdtt9o4oM1CTXmymGn1"
		addr, err := Base58ToAddress(addrStr)
		require.NoError(t, err)

		t.Run("To Bytes", func(t *testing.T) {
			bytes := addr.Bytes()
			require.Len(t, bytes, 21)
		})

		t.Run("To Hex", func(t *testing.T) {
			hex := addr.Hex()
			require.Equal(t, "0x", hex[:2])
			require.Len(t, hex, 44)
		})
	})

	t.Run("Address Validity", func(t *testing.T) {
		t.Run("Valid Address", func(t *testing.T) {
			addr, err := Base58ToAddress("TSvT6Bg3siokv3dbdtt9o4oM1CTXmymGn1")
			require.NoError(t, err)
			require.True(t, isValid(addr))
		})

		t.Run("Zero Address", func(t *testing.T) {
			addr := Address{}
			require.False(t, isValid(addr))
		})
	})
}

func TestHexToAddress(t *testing.T) {
	t.Run("Valid Hex Addresses", func(t *testing.T) {
		validHexAddresses := []string{
			"41a614f803b6fd780986a42c78ec9c7f77e6ded13c",
			"41b2a2e1b2e1b2e1b2e1b2e1b2e1b2e1b2e1b2e1b2",
			"41c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3",
		}

		for _, hexStr := range validHexAddresses {
			t.Run(hexStr, func(t *testing.T) {
				addr, err := HexToAddress(hexStr)
				require.NoError(t, err)
				require.Equal(t, "0x"+hexStr, addr.Hex())
			})
		}
	})

	t.Run("Invalid Hex Addresses", func(t *testing.T) {
		invalidHexAddresses := []string{
			"41a614f803b6fd780986a42c78ec9c7f77e6ded13",      // Too short
			"41b2a2e1b2e1b2e1b2e1b2e1b2e1b2e1b2e1b2e1b2e1b2", // Too long
			"41g3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3",     // Invalid character 'g'
			"c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3",     // Missing prefix '41'
		}

		for _, hexStr := range invalidHexAddresses {
			t.Run(hexStr, func(t *testing.T) {
				_, err := HexToAddress(hexStr)
				require.Error(t, err)
			})
		}
	})
}

// Helper Functions for testing

// isValid checks if the address is a valid TRON address
func isValid(a Address) bool {
	// Check if it's a valid Base58 address
	base58Str := a.String()
	if isValidBase58Address(base58Str) {
		return true
	}

	// Check if it's a valid hex address
	hexStr := a.Hex()
	return isValidHexAddress(strings.TrimPrefix(hexStr, "0x"))
}

// isValidBase58Address check if a string is a valid Base58 TRON address
func isValidBase58Address(address string) bool {
	// Check if the address starts with 'T' and is 34 characters long
	if len(address) != 34 || address[0] != 'T' {
		return false
	}

	// Check if the address contains only valid Base58 characters
	validChars := regexp.MustCompile("^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]+$")
	return validChars.MatchString(address)
}

// isValidHexAddressto check if a string is a valid hex TRON address
func isValidHexAddress(address string) bool {
	// Check if the address starts with '41' and is 42 characters long
	if len(address) != 42 || address[:2] != "41" {
		return false
	}

	// Check if the address contains only valid hexadecimal characters
	validChars := regexp.MustCompile("^[0-9A-Fa-f]+$")
	return validChars.MatchString(address[2:]) // Check the part after '41'
}
