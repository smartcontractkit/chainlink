package types

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/services"
)

//go:generate mockery --quiet --name HeadTracker --output ../mocks/ --case=underscore
type HeadTracker[H Head[BLOCK_HASH], BLOCK_HASH Hashable] interface {
	services.ServiceCtx
	// Backfill given a head will fill in any missing heads up to the given depth
	// (used for testing)
	Backfill(ctx context.Context, headWithChain H, depth uint) (err error)
	LatestChain() H
}

// HeadTrackable is implemented by the core txm,
// to be able to receive head events from any chain.
// Chain implementations should notify head events to the core txm via this interface.
//
//go:generate mockery --quiet --name HeadTrackable --output ./mocks/ --case=underscore
type HeadTrackable[H Head[BLOCK_HASH], BLOCK_HASH Hashable] interface {
	OnNewLongestChain(ctx context.Context, head H)
}

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

// HeadListener is a chain agnostic interface that manages connection of Client that receives heads from the blockchain node
type HeadListener[H Head[BLOCK_HASH], BLOCK_HASH Hashable] interface {
	// ListenForNewHeads kicks off the listen loop (not thread safe)
	// done() must be executed upon leaving ListenForNewHeads()
	ListenForNewHeads(handleNewHead NewHeadHandler[H, BLOCK_HASH], done func())

	// ReceivingHeads returns true if the listener is receiving heads (thread safe)
	ReceivingHeads() bool

	// Connected returns true if the listener is connected (thread safe)
	Connected() bool

	// HealthReport returns report of errors within HeadListener
	HealthReport() map[string]error
}

// NewHeadHandler is a callback that handles incoming heads
type NewHeadHandler[H Head[BLOCK_HASH], BLOCK_HASH Hashable] func(ctx context.Context, header H) error

type HeadBroadcaster[H Head[BLOCK_HASH], BLOCK_HASH Hashable] interface {
	services.ServiceCtx
	BroadcastNewLongestChain(H)
	HeadBroadcasterRegistry[H, BLOCK_HASH]
}

type HeadBroadcasterRegistry[H Head[BLOCK_HASH], BLOCK_HASH Hashable] interface {
	Subscribe(callback HeadTrackable[H, BLOCK_HASH]) (currentLongestChain H, unsubscribe func())
}
