package cltest

import "encoding/hex"

func MustHexDecodeString(s string) []byte {
	a, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return a
}

func MustHexDecode32ByteString(s string) [32]byte {
	a, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	if len(a) != 32 {
		panic("not 32 bytes")
	}
	var res [32]byte
	copy(res[:], a[:])
	return res
}
