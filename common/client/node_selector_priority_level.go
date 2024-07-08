package client

import (
	"math"
	"sort"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type priorityLevelNodeSelector[
	CHAIN_ID types.ID,
	RPC any,
] struct {
	nodes           []Node[CHAIN_ID, RPC]
	roundRobinCount []atomic.Uint32
}

type nodeWithPriority[
	CHAIN_ID types.ID,
	RPC any,
] struct {
	node     Node[CHAIN_ID, RPC]
	priority int32
}

func NewPriorityLevelNodeSelector[
	CHAIN_ID types.ID,
	RPC any,
](nodes []Node[CHAIN_ID, RPC]) NodeSelector[CHAIN_ID, RPC] {
	return &priorityLevelNodeSelector[CHAIN_ID, RPC]{
		nodes:           nodes,
		roundRobinCount: make([]atomic.Uint32, nrOfPriorityTiers(nodes)),
	}
}

func (s priorityLevelNodeSelector[CHAIN_ID, RPC]) Select() Node[CHAIN_ID, RPC] {
	nodes := s.getHighestPriorityAliveTier()

	if len(nodes) == 0 {
		return nil
	}
	priorityLevel := nodes[len(nodes)-1].priority

	// NOTE: Inc returns the number after addition, so we must -1 to get the "current" counter
	count := s.roundRobinCount[priorityLevel].Add(1) - 1
	idx := int(count % uint32(len(nodes)))

	return nodes[idx].node
}

func (s priorityLevelNodeSelector[CHAIN_ID, RPC]) Name() string {
	return NodeSelectionModePriorityLevel
}

// getHighestPriorityAliveTier filters nodes that are not in state NodeStateAlive and
// returns only the highest tier of alive nodes
func (s priorityLevelNodeSelector[CHAIN_ID, RPC]) getHighestPriorityAliveTier() []nodeWithPriority[CHAIN_ID, RPC] {
	var nodes []nodeWithPriority[CHAIN_ID, RPC]
	for _, n := range s.nodes {
		if n.State() == NodeStateAlive {
			nodes = append(nodes, nodeWithPriority[CHAIN_ID, RPC]{n, n.Order()})
		}
	}

	if len(nodes) == 0 {
		return nil
	}

	return removeLowerTiers(nodes)
}

// removeLowerTiers take a slice of nodeWithPriority[CHAIN_ID, BLOCK_HASH, HEAD, RPC] and keeps only the highest tier
func removeLowerTiers[
	CHAIN_ID types.ID,
	RPC any,
](nodes []nodeWithPriority[CHAIN_ID, RPC]) []nodeWithPriority[CHAIN_ID, RPC] {
	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i].priority > nodes[j].priority
	})

	var nodes2 []nodeWithPriority[CHAIN_ID, RPC]
	currentPriority := nodes[len(nodes)-1].priority

	for _, n := range nodes {
		if n.priority == currentPriority {
			nodes2 = append(nodes2, n)
		}
	}

	return nodes2
}

// nrOfPriorityTiers calculates the total number of priority tiers
func nrOfPriorityTiers[
	CHAIN_ID types.ID,
	RPC any,
](nodes []Node[CHAIN_ID, RPC]) int32 {
	highestPriority := int32(0)
	for _, n := range nodes {
		priority := n.Order()
		if highestPriority < priority {
			highestPriority = priority
		}
	}
	return highestPriority + 1
}

// firstOrHighestPriority takes a list of nodes and returns the first one with the highest priority
func firstOrHighestPriority[
	CHAIN_ID types.ID,
	RPC any,
](nodes []Node[CHAIN_ID, RPC]) Node[CHAIN_ID, RPC] {
	hp := int32(math.MaxInt32)
	var node Node[CHAIN_ID, RPC]
	for _, n := range nodes {
		if n.Order() < hp {
			hp = n.Order()
			node = n
		}
	}
	return node
}
