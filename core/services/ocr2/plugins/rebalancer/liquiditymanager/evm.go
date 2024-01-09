package liquiditymanager

import (
	"context"
	"fmt"
	"math/big"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditygraph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

type EvmLiquidityManager struct {
	address   models.Address
	networkID models.NetworkID
}

func NewEvmLiquidityManager(address models.Address) *EvmLiquidityManager {
	return &EvmLiquidityManager{}
}

func (e EvmLiquidityManager) MoveLiquidity(ctx context.Context, chainID models.NetworkID, amount *big.Int) error {
	return nil
}

func (e EvmLiquidityManager) GetLiquidityManagers(ctx context.Context) (map[models.NetworkID]models.Address, error) {
	return nil, nil
}

func (e EvmLiquidityManager) GetBalance(ctx context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (e EvmLiquidityManager) GetPendingTransfers(ctx context.Context) ([]models.PendingTransfer, error) {
	return nil, nil
}

func (e EvmLiquidityManager) Discover(ctx context.Context, lmFactory Factory) (*Registry, liquiditygraph.LiquidityGraph, error) {
	g := liquiditygraph.NewGraph()
	lms := NewRegistry()

	type qItem struct {
		networkID models.NetworkID
		lmAddress models.Address
	}

	seen := mapset.NewSet[qItem]()
	queue := mapset.NewSet[qItem]()

	elem := qItem{networkID: e.networkID, lmAddress: e.address}
	queue.Add(elem)
	seen.Add(elem)

	for queue.Cardinality() > 0 {
		elem, ok := queue.Pop()
		if !ok {
			return nil, nil, fmt.Errorf("unexpected internal error, there is a bug in the algorithm")
		}

		// TODO: investigate fetching the balance here.
		g.AddNetwork(elem.networkID, big.NewInt(0))

		lm, err := lmFactory.NewLiquidityManager(elem.networkID, elem.lmAddress)
		if err != nil {
			return nil, nil, fmt.Errorf("init liquidity manager: %w", err)
		}

		destinationLMs, err := lm.GetLiquidityManagers(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("get %v destination liquidity managers: %w", elem.networkID, err)
		}

		for destNetworkID, lmAddr := range destinationLMs {
			g.AddConnection(elem.networkID, destNetworkID)

			newElem := qItem{networkID: destNetworkID, lmAddress: lmAddr}
			if !seen.Contains(newElem) {
				queue.Add(newElem)
				seen.Add(newElem)

				if _, exists := lms.Get(destNetworkID); !exists {
					lms.Add(destNetworkID, lmAddr)
				}
			}
		}
	}

	return lms, g, nil
}

func (e EvmLiquidityManager) Close(ctx context.Context) error {
	return nil
}
