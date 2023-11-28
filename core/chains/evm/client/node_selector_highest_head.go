package client

import (
	"math"
)

type highestHeadNodeSelector []Node

// Deprecated: use [pkg/github.com/smartcontractkit/chainlink/v2/common/client.NewHighestHeadNodeSelector]
func NewHighestHeadNodeSelector(nodes []Node) NodeSelector {
	return highestHeadNodeSelector(nodes)
}

func (s highestHeadNodeSelector) Select() Node {
	var highestHeadNumber int64 = math.MinInt64
	var highestHeadNodes []Node
	for _, n := range s {
		state, currentHeadNumber, _ := n.StateAndLatest()
		if state == NodeStateAlive && currentHeadNumber >= highestHeadNumber {
			if highestHeadNumber < currentHeadNumber {
				highestHeadNumber = currentHeadNumber
				highestHeadNodes = nil
			}
			highestHeadNodes = append(highestHeadNodes, n)
		}
	}
	return firstOrHighestPriority(highestHeadNodes)
}

func (s highestHeadNodeSelector) Name() string {
	return NodeSelectionMode_HighestHead
}
