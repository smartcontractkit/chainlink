package headtracker

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// HeadSaver is an chain agnostic interface for saving and loading heads
// Different chains will instantiate generic HeadSaver type with their native Head and BlockHash types.
type HeadSaver[H types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable] interface {
	// Save updates the latest block number, if indeed the latest, and persists
	// this number in case of reboot.
	Save(ctx context.Context, head H) error
	// Load loads latest heads up to latestFinalized - historyDepth, returns the latest chain.
	Load(ctx context.Context, latestFinalized int64) (H, error)
	// LatestChain returns the block header with the highest number that has been seen, or nil.
	LatestChain() H
	// Chain returns a head for the specified hash, or nil.
	Chain(hash BLOCK_HASH) H
	// MarkFinalized - marks matching block and all it's direct ancestors as finalized
	MarkFinalized(ctx context.Context, latestFinalized H) error
}
