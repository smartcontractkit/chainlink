package vaultutils

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeEncryptedDecryptionShareString(t *testing.T) {
	raw := []byte{0xde, 0xad, 0xbe, 0xef}

	t.Run("hex", func(t *testing.T) {
		got, err := DecodeEncryptedDecryptionShareString(hex.EncodeToString(raw))
		require.NoError(t, err)
		require.Equal(t, raw, got)
	})

	t.Run("b64 prefix", func(t *testing.T) {
		s := EncryptedDecryptionShareB64Prefix + base64.StdEncoding.EncodeToString(raw)
		got, err := DecodeEncryptedDecryptionShareString(s)
		require.NoError(t, err)
		require.Equal(t, raw, got)
	})
}
