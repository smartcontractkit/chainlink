package client

import (
	"sync/atomic"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type roundRobinSelector[
	CHAIN_ID types.ID,
	RPC any,
] struct {
	nodes           []Node[CHAIN_ID, RPC]
	roundRobinCount atomic.Uint32
}

func NewRoundRobinSelector[
	CHAIN_ID types.ID,
	RPC any,
](nodes []Node[CHAIN_ID, RPC]) NodeSelector[CHAIN_ID, RPC] {
	return &roundRobinSelector[CHAIN_ID, RPC]{
		nodes: nodes,
	}
}

func (s *roundRobinSelector[CHAIN_ID, RPC]) Select() Node[CHAIN_ID, RPC] {
	var liveNodes []Node[CHAIN_ID, RPC]
	for _, n := range s.nodes {
		if n.State() == nodeStateAlive {
			liveNodes = append(liveNodes, n)
		}
	}

	nNodes := len(liveNodes)
	if nNodes == 0 {
		return nil
	}

	// NOTE: Inc returns the number after addition, so we must -1 to get the "current" counter
	count := s.roundRobinCount.Add(1) - 1
	idx := int(count % uint32(nNodes))

	return liveNodes[idx]
}

func (s *roundRobinSelector[CHAIN_ID, RPC]) Name() string {
	return NodeSelectionModeRoundRobin
}
