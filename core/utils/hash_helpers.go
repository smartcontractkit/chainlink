package utils

import (
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
)

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b) // Assignment for errcheck. Only used in tests so we can ignore.
	return b
}

// NewHash return random Keccak256
func NewHash() common.Hash {
	return common.BytesToHash(randomBytes(32))
}

func PadByteToHash(b byte) common.Hash {
	var h [32]byte
	h[31] = b
	return h
}
