package client

import (
	"math"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type highestHeadNodeSelector[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC_CLIENT NodeClient[CHAIN_ID, HEAD],
] []Node[CHAIN_ID, HEAD, RPC_CLIENT]

func NewHighestHeadNodeSelector[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC_CLIENT NodeClient[CHAIN_ID, HEAD],
](nodes []Node[CHAIN_ID, HEAD, RPC_CLIENT]) NodeSelector[CHAIN_ID, HEAD, RPC_CLIENT] {
	return highestHeadNodeSelector[CHAIN_ID, HEAD, RPC_CLIENT](nodes)
}

func (s highestHeadNodeSelector[CHAIN_ID, HEAD, RPC_CLIENT]) Select() Node[CHAIN_ID, HEAD, RPC_CLIENT] {
	var highestHeadNumber int64 = math.MinInt64
	var highestHeadNodes []Node[CHAIN_ID, HEAD, RPC_CLIENT]
	for _, n := range s {
		state, currentHeadNumber, _ := n.StateAndLatest()
		if state == nodeStateAlive && currentHeadNumber >= highestHeadNumber {
			if highestHeadNumber < currentHeadNumber {
				highestHeadNumber = currentHeadNumber
				highestHeadNodes = nil
			}
			highestHeadNodes = append(highestHeadNodes, n)
		}
	}
	return firstOrHighestPriority(highestHeadNodes)
}

func (s highestHeadNodeSelector[CHAIN_ID, HEAD, RPC_CLIENT]) Name() string {
	return NodeSelectionMode_HighestHead
}
