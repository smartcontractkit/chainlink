package vaultutils

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
)

// EncryptedDecryptionShareB64Prefix marks a vault-encrypted decryption key share
// encoded with StdEncoding base64 (after the prefix). Unprefixed strings are hex-encoded.
const EncryptedDecryptionShareB64Prefix = "b64:"

// DecodeEncryptedDecryptionShareString decodes a single encrypted decryption share string
// from the vault GetSecrets response. It supports the b64-prefixed form emitted by current
// vault nodes and hex for backwards compatibility.
func DecodeEncryptedDecryptionShareString(s string) ([]byte, error) {
	if strings.HasPrefix(s, EncryptedDecryptionShareB64Prefix) {
		return base64.StdEncoding.DecodeString(s[len(EncryptedDecryptionShareB64Prefix):])
	}
	return hex.DecodeString(s)
}
