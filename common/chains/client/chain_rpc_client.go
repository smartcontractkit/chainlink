package client

import (
	"context"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type ChainRPCClient[
	CHAIN_ID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK_HASH types.Hashable,
	TX any,
	TX_HASH types.Hashable,
	EVENT any,
	EVENT_OPS any, // event filter query options
	TX_RECEIPT any,
	FEE feetypes.Fee,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
] interface {
	RPCClient[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]

	Close() error
	ClientChainID(context.Context) (CHAIN_ID, error)
	Dial(callerCtx context.Context) error
	DisconnectAll()
	SetState(state NodeState)
}
