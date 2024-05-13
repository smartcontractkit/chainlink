package discoverer

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

const (
	discoverGoroutines = 4
)

type evmLiquidityGetter func(ctx context.Context, selector models.NetworkSelector, lmAddress common.Address) (*big.Int, error)

type evmDiscoverer struct {
	lock             sync.RWMutex
	evmClients       map[models.NetworkSelector]evmDep
	masterRebalancer models.Address
	masterSelector   models.NetworkSelector
	liquidityGetter  evmLiquidityGetter
}

func (e *evmDiscoverer) Discover(ctx context.Context) (graph.Graph, error) {
	return graph.NewGraphWithData(ctx, graph.Vertex{
		NetworkSelector:  e.masterSelector,
		LiquidityManager: e.masterRebalancer,
	}, e.getVertexData)
}

// DiscoverBalances discovers the balances of all networks in the graph.
// Up to discovererGoroutines goroutines are used to fetch the liquidities concurrently.
func (e *evmDiscoverer) DiscoverBalances(ctx context.Context, g graph.Graph) error {
	networks := g.GetNetworks()
	liquidityGetter := e.liquidityGetter
	if liquidityGetter == nil {
		liquidityGetter = e.defaultLiquidityGetter
	}
	running := make(chan struct{}, discoverGoroutines)
	results := make(chan error, len(networks))
	go func() {
		for _, selector := range networks {
			running <- struct{}{}
			go func(c context.Context, selector models.NetworkSelector) {
				defer func() { <-running }()
				err := e.updateLiquidity(c, selector, g, liquidityGetter)
				if err != nil {
					err = fmt.Errorf("get liquidity: %w", err)
				}
				results <- err
			}(ctx, selector)
		}
	}()

	// wait for results, we expect the same number of results as networks
	var errs error
	for range networks {
		err := <-results
		if err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (e *evmDiscoverer) getVertexData(ctx context.Context, v graph.Vertex) (graph.Data, []graph.Vertex, error) {
	selector, lmAddress := v.NetworkSelector, v.LiquidityManager
	dep, ok := e.getDep(selector)
	if !ok {
		return graph.Data{}, nil, fmt.Errorf("no client for master chain %+v", selector)
	}
	rebal, err := liquiditymanager.NewLiquidityManager(common.Address(lmAddress), dep.ethClient)
	if err != nil {
		return graph.Data{}, nil, fmt.Errorf("new liquiditymanager: %w", err)
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
		neighbors            []graph.Vertex
		xchainRebalancerData = make(map[models.NetworkSelector]graph.XChainLiquidityManagerData)
	)
	for _, v := range xchainRebalancers {
		neighbors = append(neighbors, graph.Vertex{
			NetworkSelector:  models.NetworkSelector(v.RemoteChainSelector),
			LiquidityManager: models.Address(v.RemoteRebalancer),
		})
		xchainRebalancerData[models.NetworkSelector(v.RemoteChainSelector)] = graph.XChainLiquidityManagerData{
			RemoteLiquidityManagerAddress: models.Address(v.RemoteRebalancer),
			LocalBridgeAdapterAddress:     models.Address(v.LocalBridge),
			RemoteTokenAddress:            models.Address(v.RemoteToken),
		}
	}

	configDigestAndEpoch, err := rebal.LatestConfigDigestAndEpoch(&bind.CallOpts{Context: ctx})
	if err != nil {
		return graph.Data{}, nil, fmt.Errorf("latest config digest and epoch: %w", err)
	}

	minimumLiquidity, err := rebal.GetMinimumLiquidity(&bind.CallOpts{Context: ctx})
	if err != nil {
		return graph.Data{}, nil, fmt.Errorf("get target balance: %w", err)
	}

	return graph.Data{
		Liquidity:               liquidity,
		TokenAddress:            models.Address(token),
		LiquidityManagerAddress: lmAddress,
		XChainLiquidityManagers: xchainRebalancerData,
		ConfigDigest:            models.ConfigDigest{ConfigDigest: configDigestAndEpoch.ConfigDigest},
		NetworkSelector:         selector,
		MinimumLiquidity:        minimumLiquidity,
	}, neighbors, nil
}

func (e *evmDiscoverer) updateLiquidity(ctx context.Context, selector models.NetworkSelector, g graph.Graph, liquidityGetter evmLiquidityGetter) error {
	lmAddress, err := g.GetLiquidityManagerAddress(selector)
	if err != nil {
		return fmt.Errorf("get rebalancer address: %w", err)
	}
	liquidity, err := liquidityGetter(ctx, selector, common.Address(lmAddress))
	if err != nil {
		return fmt.Errorf("get liquidity: %w", err)
	}
	_ = g.SetLiquidity(selector, liquidity) // TODO: handle non-existing network
	return nil
}

func (e *evmDiscoverer) getDep(selector models.NetworkSelector) (*evmDep, bool) {
	e.lock.RLock()
	defer e.lock.RUnlock()

	dep, ok := e.evmClients[selector]
	if !ok {
		return nil, false
	}
	return &dep, true
}

func (e *evmDiscoverer) defaultLiquidityGetter(ctx context.Context, selector models.NetworkSelector, lmAddress common.Address) (*big.Int, error) {
	dep, ok := e.getDep(selector)
	if !ok {
		return nil, fmt.Errorf("no client for master chain %+v", selector)
	}
	rebal, err := liquiditymanager.NewLiquidityManager(lmAddress, dep.ethClient)
	if err != nil {
		return nil, fmt.Errorf("new liquiditymanager: %w", err)
	}
	return rebal.GetLiquidity(&bind.CallOpts{
		Context: ctx,
	})
}
