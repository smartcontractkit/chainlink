package types

import (
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

//go:generate mockery --quiet --name Head --output ./mocks/ --case=underscore
type Head[H types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable, CHAIN_ID txmgrtypes.ID] interface {
	types.Head[BLOCK_HASH]
	// Equals returns true if the two heads are equal
	Equals(H) bool
	// ChainID returns the chain ID that the head is for
	ChainId() CHAIN_ID
	// Returns true if the head has a chain Id
	HasChainId() bool
	// Check if the two heads are of the same chain
	IsSameChain(CHAIN_ID) bool
}
