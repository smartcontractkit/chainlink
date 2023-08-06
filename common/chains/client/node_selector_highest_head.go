package client

import (
	"math"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type highestHeadNodeSelector[
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
] []Node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]

func NewHighestHeadNodeSelector[
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
](nodes []Node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) NodeSelector[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB] {
	return highestHeadNodeSelector[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB](nodes)
}

func (s highestHeadNodeSelector[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) Select() Node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB] {
	var highestHeadNumber int64 = math.MinInt64
	var highestHeadNodes []Node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]
	for _, n := range s {
		state, currentHeadNumber, _ := n.StateAndLatest()
		if state == NodeStateAlive && currentHeadNumber >= highestHeadNumber {
			if highestHeadNumber < currentHeadNumber {
				highestHeadNumber = currentHeadNumber
				highestHeadNodes = nil
			}
			highestHeadNodes = append(highestHeadNodes, n)
		}
	}
	return firstOrHighestPriority(highestHeadNodes)
}

func (s highestHeadNodeSelector[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) Name() string {
	return NodeSelectionMode_HighestHead
}
