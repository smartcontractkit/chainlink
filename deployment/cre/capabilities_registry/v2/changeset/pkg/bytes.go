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

// BytesTo32 converts a []byte of length 32 into a [32]byte.
func BytesTo32(b []byte) ([32]byte, error) {
	var out [32]byte
	if len(b) != 32 {
		return out, fmt.Errorf("expected 32 bytes, got %d", len(b))
	}
	copy(out[:], b)
	return out, nil
}
