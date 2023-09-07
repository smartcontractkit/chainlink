package client

import (
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type totalDifficultyNodeSelector[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC_CLIENT NodeClient[CHAIN_ID, HEAD],
] []Node[CHAIN_ID, HEAD, RPC_CLIENT]

func NewTotalDifficultyNodeSelector[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC_CLIENT NodeClient[CHAIN_ID, HEAD],
](nodes []Node[CHAIN_ID, HEAD, RPC_CLIENT]) NodeSelector[CHAIN_ID, HEAD, RPC_CLIENT] {
	return totalDifficultyNodeSelector[CHAIN_ID, HEAD, RPC_CLIENT](nodes)
}

func (s totalDifficultyNodeSelector[CHAIN_ID, HEAD, RPC_CLIENT]) Select() Node[CHAIN_ID, HEAD, RPC_CLIENT] {
	// NodeNoNewHeadsThreshold may not be enabled, in this case all nodes have td == nil
	var highestTD *utils.Big
	var nodes []Node[CHAIN_ID, HEAD, RPC_CLIENT]
	var aliveNodes []Node[CHAIN_ID, HEAD, RPC_CLIENT]

	for _, n := range s {
		state, _, currentTD := n.StateAndLatest()
		if state != nodeStateAlive {
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

func (s totalDifficultyNodeSelector[CHAIN_ID, HEAD, RPC_CLIENT]) Name() string {
	return NodeSelectionMode_TotalDifficulty
}
