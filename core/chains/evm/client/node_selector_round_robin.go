package client

import "go.uber.org/atomic"

type roundRobinSelector struct {
	nodes           []Node
	roundRobinCount atomic.Uint32
}

func NewRoundRobinSelector(nodes []Node) NodeSelector {
	return &roundRobinSelector{
		nodes: nodes,
	}
}

func (s *roundRobinSelector) Select() Node {
	var liveNodes []Node
	for _, n := range s.nodes {
		if n.State() == NodeStateAlive {
			liveNodes = append(liveNodes, n)
		}
	}

	nNodes := len(liveNodes)
	if nNodes == 0 {
		return nil
	}

	// NOTE: Inc returns the number after addition, so we must -1 to get the "current" counter
	count := s.roundRobinCount.Inc() - 1
	idx := int(count % uint32(nNodes))

	return liveNodes[idx]
}

func (s roundRobinSelector) Name() string {
	return NodeSelectionMode_RoundRobin
}
