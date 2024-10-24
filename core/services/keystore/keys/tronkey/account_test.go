package tronkey

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
		assert.NotNil(t, err)
	}

	for _, addr := range validAddresses {
		_, err := DecodeCheck(addr)
		assert.Nil(t, err)
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
				assert.Nil(t, err)
				assert.Equal(t, addrStr, addr.String())

				decoded, err := DecodeCheck(addrStr)
				assert.Nil(t, err)
				assert.True(t, bytes.Equal(decoded, addr.Bytes()))
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
				assert.NotNil(t, err)

				_, err = DecodeCheck(addrStr)
				assert.NotNil(t, err)
			})
		}
	})

	t.Run("Address Conversion", func(t *testing.T) {
		addrStr := "TSvT6Bg3siokv3dbdtt9o4oM1CTXmymGn1"
		addr, err := Base58ToAddress(addrStr)
		assert.Nil(t, err)

		t.Run("To Bytes", func(t *testing.T) {
			bytes := addr.Bytes()
			assert.Equal(t, 21, len(bytes))
		})

		t.Run("To Hex", func(t *testing.T) {
			hex := addr.Hex()
			assert.True(t, hex[:2] == "0x")
			assert.Equal(t, 44, len(hex)) // first 2 bytes are 0x
		})
	})

	t.Run("Address Validity", func(t *testing.T) {
		t.Run("Valid Address", func(t *testing.T) {
			addr, _ := Base58ToAddress("TSvT6Bg3siokv3dbdtt9o4oM1CTXmymGn1")
			assert.True(t, isValid(addr))
		})

		t.Run("Zero Address", func(t *testing.T) {
			addr := Address{}
			assert.False(t, isValid(addr))
		})
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
