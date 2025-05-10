package utils

import "encoding/hex"

// HexEncodeBytes converts bytes to "0x" prefixed hex string
func HexEncodeBytes(data []byte) string {
	return "0x" + hex.EncodeToString(data)
}

// HexEncodeBytes32 converts [32]byte to "0x" prefixed hex string
func HexEncodeBytes32(data [32]byte) string {
	return HexEncodeBytes(data[:])
}

// HexEncodeByteSlices returns "0x"-prefixed hex strings from a [][]byte
func HexEncodeByteSlices(v [][]byte) []string {
	signersHex := make([]string, len(v))
	for i, signer := range v {
		signersHex[i] = HexEncodeBytes(signer)
	}
	return signersHex
}

// HexEncodeBytes32Slice returns "0x"-prefixed hex strings from a [][32]byte
func HexEncodeBytes32Slice(v [][32]byte) []string {
	signersHex := make([]string, len(v))
	for i, signer := range v {
		signersHex[i] = HexEncodeBytes32(signer)
	}
	return signersHex
}
