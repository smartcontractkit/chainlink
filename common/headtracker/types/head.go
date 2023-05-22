package types

import (
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type Head[BLOCK_HASH types.Hashable, CHAIN_ID ID[C], C comparable] interface {
	types.Head[BLOCK_HASH]
	// ChainID returns the chain ID that the head is for
	ChainID() CHAIN_ID
}
