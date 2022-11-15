package client

import (
	"sync"

	"github.com/smartcontractkit/chainlink/core/utils"
)

type totalDifficultyNodeSelector struct {
	nodes          []Node
	lastBestNodeMu sync.Mutex
	lastBestNode   Node
}

func NewTotalDifficultyNodeSelector(nodes []Node) NodeSelector {
	return &totalDifficultyNodeSelector{nodes: nodes}
}

func (s *totalDifficultyNodeSelector) Select() Node {
	s.lastBestNodeMu.Lock()
	defer s.lastBestNodeMu.Unlock()

	var node Node
	// NodeNoNewHeadsThreshold may not be enabled, in this case all nodes have td == nil
	var maxTD *utils.Big
	if s.lastBestNode != nil {
		state, _, td := s.lastBestNode.StateAndLatest()
		if state == NodeStateAlive {
			node = s.lastBestNode
			maxTD = td
		}
	}

	for _, n := range s.nodes {
		if n == s.lastBestNode {
			continue
		}
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

	s.lastBestNode = node

	return node
}

func (s *totalDifficultyNodeSelector) Name() string {
	return NodeSelectionMode_TotalDifficulty
}
