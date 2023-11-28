package client

import (
	"math"
	"sort"
	"sync/atomic"
)

type priorityLevelNodeSelector struct {
	nodes           []Node
	roundRobinCount []atomic.Uint32
}

type nodeWithPriority struct {
	node     Node
	priority int32
}

// Deprecated: use [pkg/github.com/smartcontractkit/chainlink/v2/common/client.NewPriorityLevelNodeSelector]
func NewPriorityLevelNodeSelector(nodes []Node) NodeSelector {
	return &priorityLevelNodeSelector{
		nodes:           nodes,
		roundRobinCount: make([]atomic.Uint32, nrOfPriorityTiers(nodes)),
	}
}

func (s priorityLevelNodeSelector) Select() Node {
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

func (s priorityLevelNodeSelector) Name() string {
	return NodeSelectionMode_PriorityLevel
}

// getHighestPriorityAliveTier filters nodes that are not in state NodeStateAlive and
// returns only the highest tier of alive nodes
func (s priorityLevelNodeSelector) getHighestPriorityAliveTier() []nodeWithPriority {
	var nodes []nodeWithPriority
	for _, n := range s.nodes {
		if n.State() == NodeStateAlive {
			nodes = append(nodes, nodeWithPriority{n, n.Order()})
		}
	}

	if len(nodes) == 0 {
		return nil
	}

	return removeLowerTiers(nodes)
}

// removeLowerTiers take a slice of nodeWithPriority and keeps only the highest tier
func removeLowerTiers(nodes []nodeWithPriority) []nodeWithPriority {
	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i].priority > nodes[j].priority
	})

	var nodes2 []nodeWithPriority
	currentPriority := nodes[len(nodes)-1].priority

	for _, n := range nodes {
		if n.priority == currentPriority {
			nodes2 = append(nodes2, n)
		}
	}

	return nodes2
}

// nrOfPriorityTiers calculates the total number of priority tiers
func nrOfPriorityTiers(nodes []Node) int32 {
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
func firstOrHighestPriority(nodes []Node) Node {
	hp := int32(math.MaxInt32)
	var node Node
	for _, n := range nodes {
		if n.Order() < hp {
			hp = n.Order()
			node = n
		}
	}
	return node
}
