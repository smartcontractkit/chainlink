package discoverer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/graph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

var (
	ErrNotFound = errors.New("not found")
)

//go:generate mockery --quiet --name Factory --output ./mocks --filename factory_mock.go --case=underscore
type Factory interface {
	NewDiscoverer(selector models.NetworkSelector, rebalancerAddress models.Address) (Discoverer, error)
}

//go:generate mockery --quiet --name Discoverer --output ./mocks --filename discoverer_mock.go --case=underscore
type Discoverer interface {
	Discover(ctx context.Context) (graph.Graph, error)
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

func NewFactory(lggr logger.Logger, opts ...Opt) Factory {
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

func (f *factory) initDiscoverer(selector models.NetworkSelector, rebalancerAddress models.Address) (Discoverer, error) {
	var d Discoverer

	switch typ := selector.Type(); typ {
	case models.NetworkTypeEvm:
		_, exists := f.evmDeps[selector]
		if !exists {
			return nil, fmt.Errorf("evm dependencies not found for selector %d", selector)
		}
		d = &evmDiscoverer{
			evmClients:       f.evmDeps,
			masterRebalancer: rebalancerAddress,
			masterSelector:   selector,
		}
	}

	f.cachedDiscoverers.Store(f.cacheKey(selector, rebalancerAddress), d)
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
