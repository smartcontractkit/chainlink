package typconv

import (
	"encoding/hex"
)

// HexEncode converts a byte slice to a hex representation
func HexEncode(addr []byte) string {
	return "0x" + hex.EncodeToString(addr)
}
