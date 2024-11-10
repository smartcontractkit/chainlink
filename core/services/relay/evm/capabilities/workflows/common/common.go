package common

import (
	"github.com/ethereum/go-ethereum/crypto"
)

func Keccak256Hash(b []byte) string {
	return crypto.Keccak256Hash(b).Hex()
}
