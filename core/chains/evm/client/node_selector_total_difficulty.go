package client

import "math/big"

type totalDifficultyNodeSelector []Node

// Deprecated: use [pkg/github.com/smartcontractkit/chainlink/v2/common/client.NewTotalDifficultyNodeSelector]
func NewTotalDifficultyNodeSelector(nodes []Node) NodeSelector {
	return totalDifficultyNodeSelector(nodes)
}

func (s totalDifficultyNodeSelector) Select() Node {
	// NodeNoNewHeadsThreshold may not be enabled, in this case all nodes have td == nil
	var highestTD *big.Int
	var nodes []Node
	var aliveNodes []Node

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

func (s totalDifficultyNodeSelector) Name() string {
	return NodeSelectionMode_TotalDifficulty
}
