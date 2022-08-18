package client

import "sync"

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
	highestHeadNumber := int64(-1)
	if s.lastBestNode != nil {
		highestHeadNumber = s.lastBestNode.LatestReceivedBlockNumber()
		node = s.lastBestNode
	}

	for _, n := range s.nodes {
		if n == s.lastBestNode {
			continue
		}
		latestReceivedBlockNumber := n.LatestReceivedBlockNumber()
		if n.State() == NodeStateAlive && latestReceivedBlockNumber > highestHeadNumber {
			node = n
			highestHeadNumber = latestReceivedBlockNumber
		}
	}

	if s.lastBestNode == nil {
		s.lastBestNode = node
	}

	return node
}
