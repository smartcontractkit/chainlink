package utils

import (
	"crypto/rand"

	"github.com/ethereum/go-ethereum/common"
)

// NewHash return random Keccak256
func NewHash() common.Hash {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return common.BytesToHash(b)
}

// PadByteToHash returns a hash with zeros padded on the left of the given byte.
func PadByteToHash(b byte) common.Hash {
	var h [32]byte
	h[31] = b
	return h
}
