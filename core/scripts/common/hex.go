package common

import (
	"encoding/hex"
	"strings"
)

// EncodeToHex generates a hex string based on the byte representation, with the '0x' prefix
func EncodeToHex(str []byte) string {
	return "0x" + hex.EncodeToString(str)
}

// EncodeToString is a wrapper method for hex.EncodeToString
func EncodeToString(str []byte) string {
	return hex.EncodeToString(str)
}

// DecodeString returns the byte representation of the hexadecimal string
func DecodeString(str string) ([]byte, error) {
	return hex.DecodeString(str)
}

// DecodeHex converts a hex string to a byte array
func DecodeHex(str string) ([]byte, error) {
	str = strings.TrimPrefix(str, "0x")

	return hex.DecodeString(str)
}
