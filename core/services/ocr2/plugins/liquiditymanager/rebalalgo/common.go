package rebalalgo

import (
	"fmt"
	"math/big"

	big2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

type Funds struct {
	AvailableAmount *big.Int
}

// getExpectedGraph returns a copy of the graph instance with all the non-executed transfers applied.
func getExpectedGraph(g graph.Graph, nonExecutedTransfers []UnexecutedTransfer) (graph.Graph, error) {
	expG := g.Clone()

	for _, tr := range nonExecutedTransfers {
		liqTo, err := expG.GetLiquidity(tr.ToNetwork())
		if err != nil {
			return nil, err
		}
		expG.SetLiquidity(tr.ToNetwork(), big.NewInt(0).Add(liqTo, tr.TransferAmount()))

		// we only subtract from the sender if the transfer is still in progress, otherwise the source value would have already been updated
		switch tr.TransferStatus() {
		case models.TransferStatusProposed, models.TransferStatusInflight:
			liqFrom, err := expG.GetLiquidity(tr.FromNetwork())
			if err != nil {
				return nil, err
			}
			expG.SetLiquidity(tr.FromNetwork(), big.NewInt(0).Sub(liqFrom, tr.TransferAmount()))
		}
	}

	return expG, nil
}

// mergeProposedTransfers merges multiple transfers with the same sender and recipient into a single transfer.
func mergeProposedTransfers(transfers []models.ProposedTransfer) []models.ProposedTransfer {
	sums := make(map[[2]models.NetworkSelector]*big.Int)
	for _, tr := range transfers {
		k := [2]models.NetworkSelector{tr.From, tr.To}
		if _, exists := sums[k]; !exists {
			sums[k] = tr.TransferAmount()
			continue
		}
		sums[k] = big.NewInt(0).Add(sums[k], tr.TransferAmount())
	}

	merged := make([]models.ProposedTransfer, 0, len(transfers))
	for k, v := range sums {
		merged = append(merged, models.ProposedTransfer{From: k[0], To: k[1], Amount: big2.New(v)})
	}
	return merged
}

func minBigInt(a, b *big.Int) *big.Int {
	switch a.Cmp(b) {
	case -1: // a < b
		return a
	case 0: // a == b
		return a
	case 1: // a > b
		return b
	}
	return nil
}

// availableTransferableAmount calculates the available transferable amount of liquidity for a given network
// at two different time points (graphNow and graphLater).
// It takes a graph.Graph instance for the current time point (graphNow), a graph.Graph instance for the future time point (graphLater),
// and a models.NetworkSelector that represents the network for which to calculate the transferable amount.
// It returns the minimum of the available transferable amounts calculated from graphNow and graphLater as a *big.Int
func availableTransferableAmount(graphNow, graphLater graph.Graph, net models.NetworkSelector) (*big.Int, error) {
	nowData, err := graphNow.GetData(net)
	if err != nil {
		return nil, fmt.Errorf("error during GetData for %d in graphNow: %v", net, err)
	}
	availableAmountNow := big.NewInt(0).Sub(nowData.Liquidity, nowData.MinimumLiquidity)
	laterData, err := graphLater.GetData(net)
	if err != nil {
		return nil, fmt.Errorf("error during GetData for %d in graphLater: %v", net, err)
	}
	availableAmountLater := big.NewInt(0).Sub(laterData.Liquidity, laterData.MinimumLiquidity)
	return minBigInt(availableAmountNow, availableAmountLater), nil
}

// getTargetLiquidityDifferences calculates the liquidity differences between two graph instances.
// It returns two maps, liqDiffsNow and liqDiffsLater, where each map contains the liquidity differences for each network.
// The function iterates over the networks in graphNow and graphLater and compares their target liquidity and liquidity values.
// If the target liquidity is set to 0, automated rebalancing is disabled for that network.
// The liquidity differences are calculated by subtracting the liquidity from the target liquidity.
// The function uses the models.NetworkSelector type to identify networks.
func getTargetLiquidityDifferences(graphNow, graphLater graph.Graph) (liqDiffsNow, liqDiffsLater map[models.NetworkSelector]*big.Int, err error) {
	liqDiffsNow = make(map[models.NetworkSelector]*big.Int)
	liqDiffsLater = make(map[models.NetworkSelector]*big.Int)

	for _, net := range graphNow.GetNetworks() {
		dataNow, err := graphNow.GetData(net)
		if err != nil {
			return nil, nil, fmt.Errorf("get data now of net %v: %w", net, err)
		}

		dataLater, err := graphLater.GetData(net)
		if err != nil {
			return nil, nil, fmt.Errorf("get data later of net %v: %w", net, err)
		}

		if dataNow.TargetLiquidity == nil {
			return nil, nil, fmt.Errorf("target liquidity is nil for network %v", net)
		}
		if dataNow.TargetLiquidity.Cmp(big.NewInt(0)) == 0 {
			// automated rebalancing is disabled if target is set to 0
			liqDiffsNow[net] = big.NewInt(0)
			liqDiffsLater[net] = big.NewInt(0)
			continue
		}

		liqDiffsNow[net] = big.NewInt(0).Sub(dataNow.TargetLiquidity, dataNow.Liquidity)
		liqDiffsLater[net] = big.NewInt(0).Sub(dataLater.TargetLiquidity, dataLater.Liquidity)
	}

	return liqDiffsNow, liqDiffsLater, nil
}

// filterUnexecutedTransfers filters out transfers that have already been executed.
func filterUnexecutedTransfers(nonExecutedTransfers []UnexecutedTransfer) []UnexecutedTransfer {
	filtered := make([]UnexecutedTransfer, 0, len(nonExecutedTransfers))
	for _, tr := range nonExecutedTransfers {
		if tr.TransferStatus() != models.TransferStatusExecuted {
			filtered = append(filtered, tr)
		}
	}
	return filtered
}
