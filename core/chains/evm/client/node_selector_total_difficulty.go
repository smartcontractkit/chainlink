package client

import (
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type totalDifficultyNodeSelector []Node

func NewTotalDifficultyNodeSelector(nodes []Node) NodeSelector {
	return totalDifficultyNodeSelector(nodes)
}

func (s totalDifficultyNodeSelector) Select() Node {
	var node Node
	// NodeNoNewHeadsThreshold may not be enabled, in this case all nodes have td == nil
	var maxTD *utils.Big

	for _, n := range s {
		state, _, td := n.StateAndLatest()
		if state != NodeStateAlive {
			continue
		}
		// first, or td > max (which may be nil)
		if node == nil || td != nil && (maxTD == nil || td.Cmp(maxTD) > 0) {
			node = n
			maxTD = td
		}
	}

	return node
}

func (s totalDifficultyNodeSelector) Name() string {
	return NodeSelectionMode_TotalDifficulty
}
