package discoverer

import (
	"context"
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/rebalancer"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

type evmDiscoverer struct {
	evmClients       map[models.NetworkSelector]evmDep
	masterRebalancer models.Address
	masterSelector   models.NetworkSelector
}

func (e *evmDiscoverer) Discover(ctx context.Context) (graph.Graph, error) {
	getVertexInfo := func(ctx context.Context, selector models.NetworkSelector, rebalancerAddress models.Address) (graph.Data, []dataItem, error) {
		dep, ok := e.evmClients[selector]
		if !ok {
			return graph.Data{}, nil, fmt.Errorf("no client for master chain %+v", selector)
		}
		rebal, err := rebalancer.NewRebalancer(common.Address(rebalancerAddress), dep.ethClient)
		if err != nil {
			return graph.Data{}, nil, fmt.Errorf("new rebalancer: %w", err)
		}
		liquidity, err := rebal.GetLiquidity(&bind.CallOpts{
			Context: ctx,
		})
		if err != nil {
			return graph.Data{}, nil, fmt.Errorf("get liquidity: %w", err)
		}
		token, err := rebal.ILocalToken(&bind.CallOpts{
			Context: ctx,
		})
		if err != nil {
			return graph.Data{}, nil, fmt.Errorf("get token: %w", err)
		}
		xchainRebalancers, err := rebal.GetAllCrossChainRebalancers(&bind.CallOpts{
			Context: ctx,
		})
		if err != nil {
			return graph.Data{}, nil, fmt.Errorf("get all cross chain rebalancers: %w", err)
		}
		var (
			neighbors            []dataItem
			xchainRebalancerData = make(map[models.NetworkSelector]graph.XChainRebalancerData)
		)
		for _, v := range xchainRebalancers {
			neighbors = append(neighbors, dataItem{
				networkSelector:   models.NetworkSelector(v.RemoteChainSelector),
				rebalancerAddress: models.Address(v.RemoteRebalancer),
			})
			xchainRebalancerData[models.NetworkSelector(v.RemoteChainSelector)] = graph.XChainRebalancerData{
				RemoteRebalancerAddress:   models.Address(v.RemoteRebalancer),
				LocalBridgeAdapterAddress: models.Address(v.LocalBridge),
				RemoteTokenAddress:        models.Address(v.RemoteToken),
			}
		}

		configDigestAndEpoch, err := rebal.LatestConfigDigestAndEpoch(&bind.CallOpts{Context: ctx})
		if err != nil {
			return graph.Data{}, nil, fmt.Errorf("latest config digest and epoch: %w", err)
		}

		return graph.Data{
			Liquidity:         liquidity,
			TokenAddress:      models.Address(token),
			RebalancerAddress: rebalancerAddress,
			XChainRebalancers: xchainRebalancerData,
			ConfigDigest:      models.ConfigDigest{ConfigDigest: configDigestAndEpoch.ConfigDigest},
			NetworkSelector:   selector,
		}, neighbors, nil
	}

	return discover(ctx, e.masterSelector, e.masterRebalancer, getVertexInfo)
}

type dataItem struct {
	networkSelector   models.NetworkSelector
	rebalancerAddress models.Address
}

func discover(
	ctx context.Context,
	startNetwork models.NetworkSelector,
	startAddress models.Address,
	getVertexInfo func(
		ctx context.Context,
		network models.NetworkSelector,
		rebalancerAddress models.Address,
	) (graph.Data, []dataItem, error),
) (graph.Graph, error) {
	g := graph.NewGraph()

	seen := mapset.NewSet[dataItem]()
	queue := mapset.NewSet[dataItem]()

	start := dataItem{
		networkSelector:   startNetwork,
		rebalancerAddress: startAddress,
	}
	queue.Add(start)
	seen.Add(start)
	for queue.Cardinality() > 0 {
		elem, ok := queue.Pop()
		if !ok {
			return nil, fmt.Errorf("unexpected internal error")
		}

		val, neighbors, err := getVertexInfo(ctx, elem.networkSelector, elem.rebalancerAddress)
		if err != nil {
			return nil, fmt.Errorf("could not get value for vertex %+v: %w", elem, err)
		}
		g.AddNetwork(elem.networkSelector, val)

		for _, neighbor := range neighbors {
			if !g.HasNetwork(neighbor.networkSelector) {
				val2, _, err := getVertexInfo(ctx, neighbor.networkSelector, neighbor.rebalancerAddress)
				if err != nil {
					return nil, fmt.Errorf("could not get value for vertex %+v: %w", elem, err)
				}
				g.AddNetwork(neighbor.networkSelector, val2)
			}

			if err := g.AddConnection(elem.networkSelector, neighbor.networkSelector); err != nil {
				return nil, fmt.Errorf("error adding connection from %+v to %+v: %w", elem, neighbor, err)
			}

			if !seen.Contains(neighbor) {
				queue.Add(neighbor)
				seen.Add(neighbor)
			}
		}
	}

	return g, nil
}
