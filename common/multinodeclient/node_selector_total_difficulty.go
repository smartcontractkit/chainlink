package multinodeclient

import (
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type totalDifficultyNodeSelector[
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	EVENT any,
	EVENTOPS any,
	TXRECEIPT txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
] []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]

func NewTotalDifficultyNodeSelector[
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	EVENT any,
	EVENTOPS any,
	TXRECEIPT txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
](nodes []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) NodeSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE] {
	return totalDifficultyNodeSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE](nodes)
}

func (s totalDifficultyNodeSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) Select() Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE] {
	// NodeNoNewHeadsThreshold may not be enabled, in this case all nodes have td == nil
	var highestTD *utils.Big
	var nodes []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]
	var aliveNodes []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]

	for _, n := range s {
		state, _, currentTD := n.StateAndLatest()
		if state != NodeStateAlive {
			continue
		}

		aliveNodes = append(aliveNodes, n)
		if currentTD != nil && (highestTD == nil || currentTD.Cmp(highestTD) >= 0) {
			if highestTD == nil || currentTD.Cmp(highestTD) > 0 {
				highestTD = currentTD
				nodes = nil
			}
			nodes = append(nodes, n)
		}
	}

	//If all nodes have td == nil pick one from the nodes that are alive
	if len(nodes) == 0 {
		return firstOrHighestPriority(aliveNodes)
	}
	return firstOrHighestPriority(nodes)
}

func (s totalDifficultyNodeSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) Name() string {
	return NodeSelectionMode_TotalDifficulty
}
