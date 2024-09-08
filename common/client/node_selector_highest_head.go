package client

import (
	"math"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type highestHeadNodeSelector[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC NodeClient[CHAIN_ID, HEAD],
] []Node[CHAIN_ID, HEAD, RPC]

func NewHighestHeadNodeSelector[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC NodeClient[CHAIN_ID, HEAD],
](nodes []Node[CHAIN_ID, HEAD, RPC]) NodeSelector[CHAIN_ID, HEAD, RPC] {
	return highestHeadNodeSelector[CHAIN_ID, HEAD, RPC](nodes)
}

func (s highestHeadNodeSelector[CHAIN_ID, HEAD, RPC]) Select() Node[CHAIN_ID, HEAD, RPC] {
	var highestHeadNumber int64 = math.MinInt64
	var highestHeadNodes []Node[CHAIN_ID, HEAD, RPC]
	for _, n := range s {
		state, currentChainInfo := n.StateAndLatest()
		currentHeadNumber := currentChainInfo.BlockNumber
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

func (s highestHeadNodeSelector[CHAIN_ID, HEAD, RPC]) Name() string {
	return NodeSelectionModeHighestHead
}
