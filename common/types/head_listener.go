package types

import "context"

// HeadListener manages evmclient.Client connection that receives heads from the eth node
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
