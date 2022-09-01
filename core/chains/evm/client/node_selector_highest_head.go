package client

import (
	"math"
	"sync"
)

type highestHeadNodeSelector struct {
	nodes          []Node
	lastBestNodeMu sync.Mutex
	lastBestNode   Node
}

func NewHighestHeadNodeSelector(nodes []Node) NodeSelector {
	return &highestHeadNodeSelector{
		nodes:          nodes,
		lastBestNodeMu: sync.Mutex{},
		lastBestNode:   nil,
	}
}

func (s *highestHeadNodeSelector) Select() Node {
	s.lastBestNodeMu.Lock()
	defer s.lastBestNodeMu.Unlock()

	var node Node
	// NodeNoNewHeadsThreshold may not be enabled, in this case all nodes have latestReceivedBlockNumber == -1
	var highestHeadNumber int64 = math.MinInt64
	if s.lastBestNode != nil {
		state, latestReceivedBlockNumber := s.lastBestNode.StateAndLatestBlockNumber()
		if state == NodeStateAlive {
			node = s.lastBestNode
			highestHeadNumber = latestReceivedBlockNumber
		}
	}

	for _, n := range s.nodes {
		if n == s.lastBestNode {
			continue
		}
		state, latestReceivedBlockNumber := n.StateAndLatestBlockNumber()
		if state == NodeStateAlive && latestReceivedBlockNumber > highestHeadNumber {
			node = n
			highestHeadNumber = latestReceivedBlockNumber
		}
	}

	s.lastBestNode = node

	return node
}

func (s *highestHeadNodeSelector) Name() string {
	return NodeSelectionMode_HighestHead
}
