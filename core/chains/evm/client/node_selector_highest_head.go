package client

import (
	"math"
)

type highestHeadNodeSelector []Node

func NewHighestHeadNodeSelector(nodes []Node) NodeSelector {
	return highestHeadNodeSelector(nodes)
}

func (s highestHeadNodeSelector) Select() Node {
	var node Node
	// NodeNoNewHeadsThreshold may not be enabled, in this case all nodes have latestReceivedBlockNumber == -1
	var highestHeadNumber int64 = math.MinInt64

	for _, n := range s {
		state, latestReceivedBlockNumber, _ := n.StateAndLatest()
		if state == NodeStateAlive && latestReceivedBlockNumber > highestHeadNumber {
			node = n
			highestHeadNumber = latestReceivedBlockNumber
		}
	}

	return node
}

func (s highestHeadNodeSelector) Name() string {
	return NodeSelectionMode_HighestHead
}
