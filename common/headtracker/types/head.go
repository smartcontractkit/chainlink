package types

import (
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type Head[BLOCK_HASH types.Hashable, CHAIN_ID txmgrtypes.ID] interface {
	types.Head[BLOCK_HASH]

	// ChainID returns the chain ID that the head is for
	ChainId() CHAIN_ID

	// Check if ChainId is nil
	IsChainIdNil() bool

	// Check if the two heads are of the same chain
	IsSameChain(CHAIN_ID) bool
}
