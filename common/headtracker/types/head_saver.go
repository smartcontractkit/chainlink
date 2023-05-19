package types

import (
	"context"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
)

// HeadSaver is an chain agnostic interface for saving and loading heads
// Each chain that are not EVM will have its own custom implementation of HeadSaver interface.
type HeadSaver[H commontypes.Head[BLOCK_HASH], BLOCK_HASH commontypes.Hashable] interface {
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
