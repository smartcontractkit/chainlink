package types

import (
	"github.com/smartcontractkit/chainlink-framework/chains"
)

type Head[BLOCK_HASH chains.Hashable, CHAIN_ID chains.ID] interface {
	chains.Head[BLOCK_HASH]
	// ChainID returns the chain ID that the head is for
	ChainID() CHAIN_ID
	// Returns true if the head has a chain Id
	HasChainID() bool
	// IsValid returns true if the head is valid.
	IsValid() bool
}
