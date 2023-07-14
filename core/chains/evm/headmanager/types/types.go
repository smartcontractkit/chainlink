package types

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

// Saver maintains chains persisted in DB. All methods are thread-safe.
type Saver interface {
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

// Tracker holds and stores the latest block number experienced by this particular node in a thread safe manner.
// Reconstitutes the last block number from the data store on reboot.
//
//go:generate mockery --quiet --name Tracker --output ../mocks/ --case=underscore
type Tracker interface {
	services.ServiceCtx
	// Backfill given a head will fill in any missing heads up to the given depth
	// (used for testing)
	Backfill(ctx context.Context, headWithChain *evmtypes.Head, depth uint) (err error)
	LatestChain() *evmtypes.Head
}

// HeadTrackable represents any object that wishes to respond to ethereum events,
// after being subscribed to Broadcaster
type HeadTrackable = commontypes.HeadTrackable[*evmtypes.Head, common.Hash]

type BroadcasterRegistry interface {
	Subscribe(callback HeadTrackable) (currentLongestChain *evmtypes.Head, unsubscribe func())
}

// Broadcaster relays heads from the head tracker to subscribed jobs, it is less robust against
// congestion than the head tracker, and missed heads should be expected by consuming jobs
//
//go:generate mockery --quiet --name Broadcaster --output ../mocks/ --case=underscore
type Broadcaster interface {
	services.ServiceCtx
	BroadcastNewLongestChain(head *evmtypes.Head)
	BroadcasterRegistry
}
