package types

import (
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// Heads is a collection of heads. All methods are thread-safe.
type Heads[H types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable] interface {
	// LatestHead returns the block header with the highest number that has been seen, or nil.
	LatestHead() H
	// HeadByHash returns a head for the specified hash, or nil.
	HeadByHash(hash BLOCK_HASH) H
	// AddHeads adds newHeads to the collection, eliminates duplicates,
	// sorts by head number, fixes parents and cuts off old heads (historyDepth).
	AddHeads(historyDepth uint, newHeads ...H)
	// Count returns number of heads in the collection.
	Count() int
}
