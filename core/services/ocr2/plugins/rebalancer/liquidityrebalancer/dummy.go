package liquidityrebalancer

import (
	"fmt"
	"sort"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditygraph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

// DummyRebalancer selects the node with the highest balance
// and moves all the liquidity from the other nodes to it.
// inflightTransfers are ignored. This implementation is just an example.
type DummyRebalancer struct{}

func NewDummyRebalancer() *DummyRebalancer {
	return &DummyRebalancer{}
}

func (r *DummyRebalancer) ComputeTransfersToBalance(g liquiditygraph.LiquidityGraph, inflightTransfers []models.PendingTransfer, _ []models.NetworkLiquidity) ([]models.Transfer, error) {
	if g.IsEmpty() {
		return nil, fmt.Errorf("empty graph")
	}

	keys := g.GetNetworks()
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	luckyNode := keys[0]
	maxV, err := g.GetLiquidity(luckyNode)
	if err != nil {
		return nil, fmt.Errorf("get weight %v: %w", luckyNode, err)
	}

	for _, k := range keys {
		w, err := g.GetLiquidity(k)
		if err != nil {
			return nil, fmt.Errorf("get weight %v: %w", k, err)
		}

		if w.Cmp(maxV) > 0 {
			luckyNode = k
			maxV = w
		}
	}

	transfers := make([]models.Transfer, 0)
	for _, node := range g.GetNetworks() {
		if node == luckyNode {
			continue
		}

		w, err := g.GetLiquidity(node)
		if err != nil {
			return nil, fmt.Errorf("get weight %v: %w", node, err)
		}

		if w.BitLen() == 0 {
			continue
		}

		transfers = append(transfers, models.Transfer{
			From:   node,
			To:     luckyNode,
			Amount: w,
		})
	}

	return transfers, nil
}
