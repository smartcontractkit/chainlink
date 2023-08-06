package client

import (
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type totalDifficultyNodeSelector[
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

func NewTotalDifficultyNodeSelector[
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
	return totalDifficultyNodeSelector[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB](nodes)
}

func (s totalDifficultyNodeSelector[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) Select() Node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB] {
	// NodeNoNewHeadsThreshold may not be enabled, in this case all nodes have td == nil
	var highestTD *utils.Big
	var nodes []Node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]
	var aliveNodes []Node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]

	for _, n := range s {
		state, _, currentTD := n.StateAndLatest()
		if state != NodeStateAlive {
			continue
		}

		aliveNodes = append(aliveNodes, n)
		if currentTD != nil && (highestTD == nil || currentTD.Cmp(highestTD) >= 0) {
			if highestTD == nil || currentTD.Cmp(highestTD) > 0 {
				highestTD = currentTD
				nodes = nil
			}
			nodes = append(nodes, n)
		}
	}

	//If all nodes have td == nil pick one from the nodes that are alive
	if len(nodes) == 0 {
		return firstOrHighestPriority(aliveNodes)
	}
	return firstOrHighestPriority(nodes)
}

func (s totalDifficultyNodeSelector[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) Name() string {
	return NodeSelectionMode_TotalDifficulty
}
