package utils

import (
	"crypto/rand"

	"github.com/ethereum/go-ethereum/common"
)

// NewHash return random Keccak256
func NewHash() common.Hash {
	return common.BytesToHash(randomBytes(32))
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}
