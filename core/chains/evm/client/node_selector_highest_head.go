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
		state, latestReceivedBlockNumber := s.lastBestNode.State()
		if state == NodeStateAlive {
			node = s.lastBestNode
			highestHeadNumber = latestReceivedBlockNumber
		}
	}

	for _, n := range s.nodes {
		if n == s.lastBestNode {
			continue
		}
		state, latestReceivedBlockNumber := n.State()
		if state == NodeStateAlive && latestReceivedBlockNumber > highestHeadNumber {
			node = n
			highestHeadNumber = latestReceivedBlockNumber
		}
	}

	if s.lastBestNode == nil {
		s.lastBestNode = node
	}

	return node
}
