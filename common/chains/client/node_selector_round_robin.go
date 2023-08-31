package client

import (
	"sync/atomic"

	nodetypes "github.com/smartcontractkit/chainlink/v2/common/chains/client/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type roundRobinSelector[
	CHAIN_ID types.ID,
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
	RPC_CLIENT nodetypes.NodeClientAPI[CHAIN_ID, BLOCK_HASH, HEAD, SUB],
] struct {
	nodes           []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]
	roundRobinCount atomic.Uint32
}

func NewRoundRobinSelector[
	CHAIN_ID types.ID,
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
	RPC_CLIENT nodetypes.NodeClientAPI[CHAIN_ID, BLOCK_HASH, HEAD, SUB],
](nodes []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) NodeSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT] {
	return &roundRobinSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]{
		nodes: nodes,
	}
}

func (s *roundRobinSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Select() Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT] {
	var liveNodes []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]
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
	count := s.roundRobinCount.Add(1) - 1
	idx := int(count % uint32(nNodes))

	return liveNodes[idx]
}

func (s *roundRobinSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Name() string {
	return NodeSelectionMode_RoundRobin
}
