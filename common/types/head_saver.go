package types

import (
	"context"
)

// HeadSaver is an chain agnostic interface for saving and loading heads
// Different chains will instantiate generic HeadSaver type with their native Head and BlockHash types.
type HeadSaver[H Head[BLOCK_HASH], BLOCK_HASH Hashable] interface {
	// Save updates the latest block number, if indeed the latest, and persists
	// this number in case of reboot.
	Save(ctx context.Context, head H) error
	// Load loads latest EvmHeadTrackerHistoryDepth heads, returns the latest chain.
	Load(ctx context.Context) (H, error)
	// LatestChain returns the block header with the highest number that has been seen, or nil.
	LatestChain() H
	// Chain returns a head for the specified hash, or nil.
	Chain(hash BLOCK_HASH) H
}

// TODO: This will be a temporary interface until we move the Head Tracker component to make use of HeadTrackerHead
type InMemoryHeadSaver[H HeadTrackerHead[BLOCK_HASH, CHAIN_ID], BLOCK_HASH Hashable, CHAIN_ID ID] interface {
	HeadSaver[H, BLOCK_HASH]
}
