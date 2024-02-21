package liquiditymanager

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrInternalCacheIssue = errors.New("internal cache issue")
)

// Factory initializes a new liquidity manager instance.
//
//go:generate mockery --quiet --name Factory --output ../rebalancermocks --filename lm_factory_mock.go --case=underscore
type Factory interface {
	// NewRebalancer will initialize a new rebalancer instance based on the provided params.
	NewRebalancer(networkID models.NetworkSelector, address models.Address) (Rebalancer, error)

	// GetRebalancer returns an already initialized (via NewRebalancer) rebalancer instance.
	// If it does not exist returns ErrNotFound.
	GetRebalancer(networkID models.NetworkSelector, address models.Address) (Rebalancer, error)
}

type evmDep struct {
	ethClient client.Client
}

type BaseRebalancerFactory struct {
	evmDeps           map[models.NetworkSelector]evmDep
	cachedRebalancers sync.Map
	lggr              logger.Logger
}

type Opt func(f *BaseRebalancerFactory)

func NewBaseRebalancerFactory(lggr logger.Logger, opts ...Opt) *BaseRebalancerFactory {
	f := &BaseRebalancerFactory{
		evmDeps: make(map[models.NetworkSelector]evmDep),
		lggr:    lggr,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func WithEvmDep(networkID models.NetworkSelector, ethClient client.Client) Opt {
	return func(f *BaseRebalancerFactory) {
		f.evmDeps[networkID] = evmDep{
			ethClient: ethClient,
		}
	}
}

func (b *BaseRebalancerFactory) NewRebalancer(networkSel models.NetworkSelector, address models.Address) (Rebalancer, error) {
	rb, err := b.GetRebalancer(networkSel, address)
	if errors.Is(err, ErrNotFound) {
		return b.initRebalancer(networkSel, address)
	}
	return rb, err
}

func (b *BaseRebalancerFactory) initRebalancer(networkSel models.NetworkSelector, address models.Address) (Rebalancer, error) {
	var rb Rebalancer
	var err error

	switch typ := networkSel.Type(); typ {
	case models.NetworkTypeEvm:
		evmDeps, exists := b.evmDeps[networkSel]
		if !exists {
			return nil, fmt.Errorf("evm dependencies not found for selector %d", networkSel)
		}

		rb, err = NewEvmRebalancer(address, networkSel, evmDeps.ethClient, b.lggr)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("liquidity manager of type %v (network selector: %d) is not supported", typ, networkSel)
	}

	b.cachedRebalancers.Store(b.cacheKey(networkSel, address), rb)
	return rb, nil
}

func (b *BaseRebalancerFactory) GetRebalancer(networkSel models.NetworkSelector, address models.Address) (Rebalancer, error) {
	k := b.cacheKey(networkSel, address)

	rawVal, exists := b.cachedRebalancers.Load(k)
	if !exists {
		return nil, ErrNotFound
	}

	rb, is := rawVal.(Rebalancer)
	if !is {
		return nil, ErrInternalCacheIssue
	}

	return rb, nil
}

func (b *BaseRebalancerFactory) cacheKey(networkSel models.NetworkSelector, address models.Address) string {
	return fmt.Sprintf("rebalancer-%d-%s", networkSel, address.String())
}
