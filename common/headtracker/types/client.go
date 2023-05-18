package types

import (
	"context"
	"math/big"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
)

type HeadTrackerClient[H commontypes.Head[BLOCK_HASH], S Subscription, ID txmgrtypes.ID, BLOCK_HASH commontypes.Hashable] interface {
	HeadByNumber(ctx context.Context, number *big.Int) (head H, err error)
	HeadByHash(ctx context.Context, hash BLOCK_HASH) (head H, err error)
	ConfiguredChainID() (id ID)

	// SubscribeNewHead is the method in which the client receives new Head.
	// It can be implemented differently for each chain i.e websocket, polling, etc
	SubscribeNewHead(ctx context.Context, ch chan<- H) (S, error)
}
