package client

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type ChainRPCClient[
	CHAIN_ID types.ID,
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
] interface {
	Close() error
	ClientChainID(context.Context) (CHAIN_ID, error)
	Dial(callerCtx context.Context) error
	DisconnectAll()
	SetState(state NodeState)
	Subscribe(ctx context.Context, channel chan<- HEAD, args ...interface{}) (SUB, error)
	ClientVersion(context.Context) (string, error)
}
