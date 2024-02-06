package liquiditymanager

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

// Factory initializes a new liquidity manager instance.
//
//go:generate mockery --quiet --name Factory --output ../rebalancermocks --filename lm_factory_mock.go --case=underscore
type Factory interface {
	NewRebalancer(networkID models.NetworkSelector, address models.Address) (Rebalancer, error)
}

type evmDep struct {
	lp        logpoller.LogPoller
	ethClient client.Client
}

type BaseRebalancerFactory struct {
	evmDeps map[models.NetworkSelector]evmDep
	lggr    logger.Logger
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

func WithEvmDep(networkID models.NetworkSelector, lp logpoller.LogPoller, ethClient client.Client) Opt {
	return func(f *BaseRebalancerFactory) {
		f.evmDeps[networkID] = evmDep{
			lp:        lp,
			ethClient: ethClient,
		}
	}
}

func (b *BaseRebalancerFactory) NewRebalancer(networkSel models.NetworkSelector, address models.Address) (Rebalancer, error) {
	switch typ := networkSel.Type(); typ {
	case models.NetworkTypeEvm:
		evmDeps, exists := b.evmDeps[networkSel]
		if !exists {
			return nil, fmt.Errorf("evm dependencies not found for selector %d", networkSel)
		}
		return NewEvmRebalancer(address, networkSel, evmDeps.ethClient, evmDeps.lp, b.lggr)
	default:
		return nil, fmt.Errorf("liquidity manager of type %v (network selector: %d) is not supported", typ, networkSel)
	}
}
