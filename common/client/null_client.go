package client

import (
	"math/big"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// NullClient satisfies the Client but has no side effects
type NullClient[
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	TXRECEIPT any,
	EVENT any,
	EVENTOPS any, // event filter query options
	R txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
] struct {
	cid  *big.Int
	lggr logger.Logger
}

func NewNullClient[
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	TXRECEIPT any,
	EVENT any,
	EVENTOPS any, // event filter query options
	R txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
](cid *big.Int, lggr logger.Logger) *NullClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, TXRECEIPT, EVENT, EVENTOPS, R, FEE] {
	return &NullClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, TXRECEIPT, EVENT, EVENTOPS, R, FEE]{cid: cid, lggr: lggr.Named("NullClient")}
}
