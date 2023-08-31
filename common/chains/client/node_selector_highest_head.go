package client

import (
	"math"

	nodetypes "github.com/smartcontractkit/chainlink/v2/common/chains/client/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type highestHeadNodeSelector[
	CHAIN_ID types.ID,
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
	RPC_CLIENT nodetypes.NodeClientAPI[CHAIN_ID, BLOCK_HASH, HEAD, SUB],
] []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]

func NewHighestHeadNodeSelector[
	CHAIN_ID types.ID,
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
	RPC_CLIENT nodetypes.NodeClientAPI[CHAIN_ID, BLOCK_HASH, HEAD, SUB],
](nodes []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) NodeSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT] {
	return highestHeadNodeSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT](nodes)
}

func (s highestHeadNodeSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Select() Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT] {
	var highestHeadNumber int64 = math.MinInt64
	var highestHeadNodes []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]
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

func (s highestHeadNodeSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Name() string {
	return NodeSelectionMode_HighestHead
}
