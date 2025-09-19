package pkg

import (
	"encoding/hex"
	"fmt"
)

// HexStringTo32Bytes converts a hex string (with or without 0x prefix) to [32]byte
func HexStringTo32Bytes(hexStr string) ([32]byte, error) {
	var result [32]byte

	// Remove 0x prefix if present
	if len(hexStr) >= 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}

	// Validate length
	if len(hexStr) != 64 {
		return result, fmt.Errorf("invalid hex string length: expected 64 hex characters, got %d", len(hexStr))
	}

	// Decode hex string
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return result, fmt.Errorf("invalid hex string: %w", err)
	}

	// Copy to fixed-size array
	copy(result[:], bytes)
	return result, nil
}
