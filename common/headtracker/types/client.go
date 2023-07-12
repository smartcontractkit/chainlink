package types

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type Client[H types.Head[BLOCK_HASH], S types.Subscription, ID types.ID, BLOCK_HASH types.Hashable] interface {
	HeadByNumber(ctx context.Context, number *big.Int) (head H, err error)
	// ConfiguredChainID returns the chain ID that the node is configured to connect to
	ConfiguredChainID() (id ID)
	// SubscribeNewHead is the method in which the client receives new Head.
	// It can be implemented differently for each chain i.e websocket, polling, etc
	SubscribeNewHead(ctx context.Context, ch chan<- H) (S, error)
}
