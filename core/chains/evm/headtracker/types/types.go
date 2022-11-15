package types

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services"
)

// HeadSaver maintains chains persisted in DB. All methods are thread-safe.
type HeadSaver interface {
	// Save updates the latest block number, if indeed the latest, and persists
	// this number in case of reboot.
	Save(ctx context.Context, head *evmtypes.Head) error
	// LoadFromDB loads latest EvmHeadTrackerHistoryDepth heads, returns the latest chain.
	LoadFromDB(ctx context.Context) (*evmtypes.Head, error)
	// LatestHeadFromDB returns the highest seen head from DB.
	LatestHeadFromDB(ctx context.Context) (*evmtypes.Head, error)
	// LatestChain returns the block header with the highest number that has been seen, or nil.
	LatestChain() *evmtypes.Head
	// Chain returns a head for the specified hash, or nil.
	Chain(hash common.Hash) *evmtypes.Head
}

// HeadTracker holds and stores the latest block number experienced by this particular node in a thread safe manner.
// Reconstitutes the last block number from the data store on reboot.
//
//go:generate mockery --quiet --name HeadTracker --output ../mocks/ --case=underscore
type HeadTracker interface {
	services.ServiceCtx
	// Backfill given a head will fill in any missing heads up to the given depth
	// (used for testing)
	Backfill(ctx context.Context, headWithChain *evmtypes.Head, depth uint) (err error)
	LatestChain() *evmtypes.Head
}

// HeadTrackable represents any object that wishes to respond to ethereum events,
// after being subscribed to HeadBroadcaster
//
//go:generate mockery --quiet --name HeadTrackable --output ../mocks/ --case=underscore
type HeadTrackable interface {
	OnNewLongestChain(ctx context.Context, head *evmtypes.Head)
}

type HeadBroadcasterRegistry interface {
	Subscribe(callback HeadTrackable) (currentLongestChain *evmtypes.Head, unsubscribe func())
}

// HeadBroadcaster relays heads from the head tracker to subscribed jobs, it is less robust against
// congestion than the head tracker, and missed heads should be expected by consuming jobs
//
//go:generate mockery --quiet --name HeadBroadcaster --output ../mocks/ --case=underscore
type HeadBroadcaster interface {
	services.ServiceCtx
	BroadcastNewLongestChain(head *evmtypes.Head)
	HeadBroadcasterRegistry
}

// NewHeadHandler is a callback that handles incoming heads
type NewHeadHandler func(ctx context.Context, header *evmtypes.Head) error

// HeadListener manages evmclient.Client connection that receives heads from the eth node
type HeadListener interface {
	// ListenForNewHeads kicks off the listen loop (not thread safe)
	// done() must be executed upon leaving ListenForNewHeads()
	ListenForNewHeads(handleNewHead NewHeadHandler, done func())
	// ReceivingHeads returns true if the listener is receiving heads (thread safe)
	ReceivingHeads() bool
	// Connected returns true if the listener is connected (thread safe)
	Connected() bool
}
