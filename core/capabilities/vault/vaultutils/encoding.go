package vaultutils

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// CiphertextStringEncoding selects how an encrypted value string is encoded on the wire.
// Callers must choose the encoding that matches how the value was produced; there is no
// auto-detection in DecodeEncryptedValue, so ambiguous strings cannot be mis-decoded.
type CiphertextStringEncoding int

const (
	// CiphertextStringEncodingHex is hex-encoded raw ciphertext bytes (legacy gateway input).
	CiphertextStringEncodingHex CiphertextStringEncoding = iota
	// CiphertextStringEncodingBase64 is standard base64 (StdEncoding) without a prefix.
	CiphertextStringEncodingBase64
)

// DecodeEncryptedValue decodes ciphertext using exactly one wire encoding.
func DecodeEncryptedValue(s string, enc CiphertextStringEncoding) ([]byte, error) {
	switch enc {
	case CiphertextStringEncodingHex:
		return hex.DecodeString(s)
	case CiphertextStringEncodingBase64:
		if s == "" {
			return nil, fmt.Errorf("empty base64 ciphertext")
		}
		return base64.StdEncoding.DecodeString(s)
	default:
		return nil, fmt.Errorf("unknown ciphertext string encoding: %d", enc)
	}
}
