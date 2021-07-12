package cltest

import (
	"encoding/hex"
)

func MustHexDecodeString(s string) []byte {
	a, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return a
}
