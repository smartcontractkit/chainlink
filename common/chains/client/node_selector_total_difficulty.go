package client

import (
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	nodetypes "github.com/smartcontractkit/chainlink/v2/common/chains/client/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type totalDifficultyNodeSelector[
	CHAIN_ID types.ID,
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
	RPC_CLIENT nodetypes.NodeClientAPI[CHAIN_ID, BLOCK_HASH, HEAD, SUB],
] []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]

func NewTotalDifficultyNodeSelector[
	CHAIN_ID types.ID,
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
	RPC_CLIENT nodetypes.NodeClientAPI[CHAIN_ID, BLOCK_HASH, HEAD, SUB],
](nodes []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) NodeSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT] {
	return totalDifficultyNodeSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT](nodes)
}

func (s totalDifficultyNodeSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Select() Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT] {
	// NodeNoNewHeadsThreshold may not be enabled, in this case all nodes have td == nil
	var highestTD *utils.Big
	var nodes []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]
	var aliveNodes []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]

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

func (s totalDifficultyNodeSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Name() string {
	return NodeSelectionMode_TotalDifficulty
}
