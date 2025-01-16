package types

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink-framework/chains"
)

type Client[H chains.Head[BLOCK_HASH], S chains.Subscription, ID chains.ID, BLOCK_HASH chains.Hashable] interface {
	HeadByNumber(ctx context.Context, number *big.Int) (head H, err error)
	HeadByHash(ctx context.Context, hash BLOCK_HASH) (head H, err error)
	// ConfiguredChainID returns the chain ID that the node is configured to connect to
	ConfiguredChainID() (id ID)
	// SubscribeToHeads is the method in which the client receives new Head.
	// It can be implemented differently for each chain i.e websocket, polling, etc
	SubscribeToHeads(ctx context.Context) (<-chan H, S, error)
	// LatestFinalizedBlock - returns the latest block that was marked as finalized
	LatestFinalizedBlock(ctx context.Context) (head H, err error)
}
