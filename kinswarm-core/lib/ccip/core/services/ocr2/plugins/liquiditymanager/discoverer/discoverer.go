package discoverer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

var (
	ErrNotFound = errors.New("not found")
)

type Factory interface {
	NewDiscoverer(selector models.NetworkSelector, rebalancerAddress models.Address) (Discoverer, error)
}

type Discoverer interface {
	// Discover fetches the entire graph
	Discover(ctx context.Context) (graph.Graph, error)
	// DiscoverBalances fetch only the balances rather building the entire graph
	DiscoverBalances(context.Context, graph.Graph) error
}

type evmDep struct {
	ethClient client.Client
}

type factory struct {
	evmDeps           map[models.NetworkSelector]evmDep
	cachedDiscoverers sync.Map
	lggr              logger.Logger
}

type Opt func(f *factory)

func NewFactory(lggr logger.Logger, opts ...Opt) *factory {
	f := &factory{
		evmDeps: make(map[models.NetworkSelector]evmDep),
		lggr:    lggr,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func WithEvmDep(networkID models.NetworkSelector, ethClient client.Client) Opt {
	return func(f *factory) {
		f.evmDeps[networkID] = evmDep{
			ethClient: ethClient,
		}
	}
}

// NewDiscoverer implements Factory.
func (f *factory) NewDiscoverer(selector models.NetworkSelector, rebalancerAddress models.Address) (Discoverer, error) {
	d, err := f.getDiscoverer(selector, rebalancerAddress)
	if errors.Is(err, ErrNotFound) {
		return f.initDiscoverer(selector, rebalancerAddress)
	}
	return d, err
}

func (f *factory) initDiscoverer(selector models.NetworkSelector, lmAddress models.Address) (Discoverer, error) {
	var d Discoverer

	switch typ := selector.Type(); typ {
	case models.NetworkTypeEvm:
		_, exists := f.evmDeps[selector]
		if !exists {
			return nil, fmt.Errorf("evm dependencies not found for selector %d", selector)
		}
		d = newEvmDiscoverer(f.lggr, f.evmDeps, lmAddress, selector)
		f.lggr.Debugw("Created EVM Discoverer", "selector", selector, "lmAddress", lmAddress)
	}
	f.cachedDiscoverers.Store(f.cacheKey(selector, lmAddress), d)
	return d, nil
}

func (f *factory) getDiscoverer(selector models.NetworkSelector, rebalancerAddress models.Address) (Discoverer, error) {
	if d, ok := f.cachedDiscoverers.Load(f.cacheKey(selector, rebalancerAddress)); ok {
		return d.(Discoverer), nil
	}
	return nil, ErrNotFound
}

func (f *factory) cacheKey(selector models.NetworkSelector, rebalancerAddress models.Address) string {
	return fmt.Sprintf("%d-%s", selector, rebalancerAddress.String())
}
