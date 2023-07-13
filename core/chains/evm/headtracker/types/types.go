package types

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// HeadSaver maintains chains persisted in DB. All methods are thread-safe.
type HeadSaver interface {
	// Save updates the latest block number, if indeed the latest, and persists
	// this number in case of reboot.
	Save(ctx context.Context, head *evmtypes.Head) error
	// Load loads latest EvmHeadTrackerHistoryDepth heads from the DB, returns the latest chain.
	Load(ctx context.Context) (*evmtypes.Head, error)
	// LatestHeadFromDB returns the highest seen head from DB.
	LatestHeadFromDB(ctx context.Context) (*evmtypes.Head, error)
	// LatestChain returns the block header with the highest number that has been seen, or nil.
	LatestChain() *evmtypes.Head
	// Chain returns a head for the specified hash, or nil.
	Chain(hash common.Hash) *evmtypes.Head
}

// Type Alias for EVM Head Tracker Components
type (
	HeadBroadcasterRegistry = commontypes.HeadBroadcasterRegistry[*evmtypes.Head, common.Hash]
	HeadTracker             = commontypes.HeadTracker[*evmtypes.Head, common.Hash]
	HeadTrackable           = commontypes.HeadTrackable[*evmtypes.Head, common.Hash]
	HeadListener            = commontypes.HeadListener[*evmtypes.Head, common.Hash]
	HeadBroadcaster         = commontypes.HeadBroadcaster[*evmtypes.Head, common.Hash]
)
