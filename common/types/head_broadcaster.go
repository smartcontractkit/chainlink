package types

import "github.com/smartcontractkit/chainlink/v2/core/services"

type HeadBroadcaster[H Head[BLOCK_HASH], BLOCK_HASH Hashable] interface {
	services.ServiceCtx
	BroadcastNewLongestChain(H)
	HeadBroadcasterRegistry[H, BLOCK_HASH]
}

type HeadBroadcasterRegistry[H Head[BLOCK_HASH], BLOCK_HASH Hashable] interface {
	Subscribe(callback HeadTrackable[H, BLOCK_HASH]) (currentLongestChain H, unsubscribe func())
}
