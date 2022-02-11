package testutils

import (
	"math/big"
	mrand "math/rand"

	"github.com/ethereum/go-ethereum/common"
)

// NOTE: To avoid circular dependencies, this package may not import anything
// from "github.com/smartcontractkit/chainlink/core"

// FixtureChainID matches the chain always added by fixtures.sql
// It is set to 0 since no real chain ever has this ID and allows a virtual
// "test" chain ID to be used without clashes
var FixtureChainID = big.NewInt(0)

// NewAddress return a random new address
func NewAddress() common.Address {
	return common.BytesToAddress(randomBytes(20))
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = mrand.Read(b) // Assignment for errcheck. Only used in tests so we can ignore.
	return b
}

// Random32Byte returns a random [32]byte
func Random32Byte() (b [32]byte) {
	copy(b[:], randomBytes(32))
	return b
}
