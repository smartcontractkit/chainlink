package types

import (
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type Head[BLOCK_HASH types.Hashable, CHAIN_ID types.ID] interface {
	types.Head[BLOCK_HASH]
	// ChainID returns the chain ID that the head is for
	ChainID() CHAIN_ID
	// Returns true if the head has a chain Id
	HasChainID() bool
	// IsValid returns true if the head is valid.
	IsValid() bool
}
