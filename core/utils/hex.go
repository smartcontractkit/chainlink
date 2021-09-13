package utils

import (
	"encoding/hex"
	"strings"
)

// IsHexBytes returns true if the given bytes array is basically HEX encoded value.
func IsHexBytes(arr []byte) bool {
	_, err := hex.DecodeString(strings.TrimPrefix(string(arr), "0x"))
	return err == nil
}
