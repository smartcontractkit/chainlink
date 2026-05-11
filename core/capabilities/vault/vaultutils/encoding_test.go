package vaultutils

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeEncryptedValue_Hex(t *testing.T) {
	raw := []byte{0xde, 0xad, 0xbe, 0xef}
	enc := hex.EncodeToString(raw)
	got, err := DecodeEncryptedValue(enc, CiphertextStringEncodingHex)
	require.NoError(t, err)
	require.Equal(t, raw, got)
}

func TestDecodeEncryptedValue_Hex_oddLengthRejected(t *testing.T) {
	_, err := DecodeEncryptedValue("abc", CiphertextStringEncodingHex)
	require.Error(t, err)
}

func TestDecodeEncryptedValue_Base64(t *testing.T) {
	raw := []byte{0x01, 0x02, 0x03, 0xff}
	enc := base64.StdEncoding.EncodeToString(raw)
	got, err := DecodeEncryptedValue(enc, CiphertextStringEncodingBase64)
	require.NoError(t, err)
	require.Equal(t, raw, got)
}

func TestDecodeEncryptedValue_Base64_ambiguousHexStringUsesBase64(t *testing.T) {
	// "aabb" is valid hex ([0xaa,0xbb]) and valid base64 (different bytes).
	s := "aabb"
	gotHex, err := DecodeEncryptedValue(s, CiphertextStringEncodingHex)
	require.NoError(t, err)
	gotB64, err := DecodeEncryptedValue(s, CiphertextStringEncodingBase64)
	require.NoError(t, err)
	require.NotEqual(t, gotHex, gotB64)
}

func TestDecodeEncryptedValue_InvalidEncodingConst(t *testing.T) {
	_, err := DecodeEncryptedValue("00", CiphertextStringEncoding(99))
	require.Error(t, err)
}

func TestDecodeEncryptedValue_Invalid(t *testing.T) {
	_, err := DecodeEncryptedValue("not-valid-hex!", CiphertextStringEncodingHex)
	require.Error(t, err)

	_, err = DecodeEncryptedValue("!!!not-base64!!!", CiphertextStringEncodingBase64)
	require.Error(t, err)
}
