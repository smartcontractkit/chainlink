package evmliquiditymanager

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrInternalCacheIssue = errors.New("internal cache issue")
)

// Factory initializes a new liquidity manager instance.
type Factory interface {
	// NewLiquidityManager will initialize a new rebalancer instance based on the provided params.
	NewLiquidityManager(networkID models.NetworkSelector, address models.Address) (LiquidityManager, error)

	// GetLiquidityManager returns an already initialized (via NewLiquidityManager) liquidity manager instance.
	// If it does not exist returns ErrNotFound.
	GetLiquidityManager(networkID models.NetworkSelector, address models.Address) (LiquidityManager, error)
}

type evmDep struct {
	ethClient client.Client
}

type BaseLiquidityManagerFactory struct {
	evmDeps map[models.NetworkSelector]evmDep
	cache   sync.Map
	lggr    logger.Logger
}

type Opt func(f *BaseLiquidityManagerFactory)

func NewBaseLiquidityManagerFactory(lggr logger.Logger, opts ...Opt) *BaseLiquidityManagerFactory {
	f := &BaseLiquidityManagerFactory{
		evmDeps: make(map[models.NetworkSelector]evmDep),
		lggr:    lggr,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func WithEvmDep(networkID models.NetworkSelector, ethClient client.Client) Opt {
	return func(f *BaseLiquidityManagerFactory) {
		f.evmDeps[networkID] = evmDep{
			ethClient: ethClient,
		}
	}
}

func (b *BaseLiquidityManagerFactory) NewLiquidityManager(networkSel models.NetworkSelector, address models.Address) (LiquidityManager, error) {
	rb, err := b.GetLiquidityManager(networkSel, address)
	if errors.Is(err, ErrNotFound) {
		return b.initLiquidityManager(networkSel, address)
	}
	return rb, err
}

func (b *BaseLiquidityManagerFactory) initLiquidityManager(networkSel models.NetworkSelector, address models.Address) (LiquidityManager, error) {
	var rb LiquidityManager
	var err error

	switch typ := networkSel.Type(); typ {
	case models.NetworkTypeEvm:
		evmDeps, exists := b.evmDeps[networkSel]
		if !exists {
			return nil, fmt.Errorf("evm dependencies not found for selector %d", networkSel)
		}

		rb, err = NewEvmLiquidityManager(address, networkSel, evmDeps.ethClient, b.lggr)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("liquidity manager of type %v (network selector: %d) is not supported", typ, networkSel)
	}

	b.cache.Store(b.cacheKey(networkSel, address), rb)
	return rb, nil
}

func (b *BaseLiquidityManagerFactory) GetLiquidityManager(networkSel models.NetworkSelector, address models.Address) (LiquidityManager, error) {
	k := b.cacheKey(networkSel, address)

	rawVal, exists := b.cache.Load(k)
	if !exists {
		return nil, ErrNotFound
	}

	rb, is := rawVal.(LiquidityManager)
	if !is {
		return nil, ErrInternalCacheIssue
	}

	return rb, nil
}

func (b *BaseLiquidityManagerFactory) cacheKey(networkSel models.NetworkSelector, address models.Address) string {
	return fmt.Sprintf("rebalancer-%d-%s", networkSel, address.String())
}
