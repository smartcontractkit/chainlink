package utils

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
)

// NewHash return random Keccak256
func NewHash() common.Hash {
	b := make([]byte, 32)
	// #nosec this method is only used in tests
	_, _ = rand.Read(b) // Assignment for errcheck. Only used in tests so we can ignore.
	return common.BytesToHash(b)
}

// PadByteToHash returns a hash with zeros padded on the left of the given byte.
func PadByteToHash(b byte) common.Hash {
	var h [32]byte
	h[31] = b
	return h
}
