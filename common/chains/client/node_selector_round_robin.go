package client

import (
	"sync/atomic"

	nodetypes "github.com/smartcontractkit/chainlink/v2/common/chains/client/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type roundRobinSelector[
	CHAIN_ID types.ID,
	HEAD nodetypes.Head,
	RPC_CLIENT nodetypes.NodeClient[CHAIN_ID, HEAD],
] struct {
	nodes           []Node[CHAIN_ID, HEAD, RPC_CLIENT]
	roundRobinCount atomic.Uint32
}

func NewRoundRobinSelector[
	CHAIN_ID types.ID,
	HEAD nodetypes.Head,
	RPC_CLIENT nodetypes.NodeClient[CHAIN_ID, HEAD],
](nodes []Node[CHAIN_ID, HEAD, RPC_CLIENT]) NodeSelector[CHAIN_ID, HEAD, RPC_CLIENT] {
	return &roundRobinSelector[CHAIN_ID, HEAD, RPC_CLIENT]{
		nodes: nodes,
	}
}

func (s *roundRobinSelector[CHAIN_ID, HEAD, RPC_CLIENT]) Select() Node[CHAIN_ID, HEAD, RPC_CLIENT] {
	var liveNodes []Node[CHAIN_ID, HEAD, RPC_CLIENT]
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

func (s *roundRobinSelector[CHAIN_ID, HEAD, RPC_CLIENT]) Name() string {
	return NodeSelectionMode_RoundRobin
}
