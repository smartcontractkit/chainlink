package common

import (
	"encoding/binary"

	"golang.org/x/exp/slices"
)

func Uint32ToBytes(val uint32) []byte {
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, val)
	return result
}

func BytesToUint32(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}

// input string can't have any 0x0 characters
func StringToAlignedBytes(input string, size int) []byte {
	aligned := make([]byte, size)
	copy(aligned, input)
	return aligned
}

func AlignedBytesToString(data []byte) string {
	idx := slices.IndexFunc(data, func(b byte) bool { return b == 0 })
	return string(data[:idx])
}
