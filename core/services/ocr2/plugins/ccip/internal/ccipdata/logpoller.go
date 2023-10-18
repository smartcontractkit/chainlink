package ccipdata

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var _ Reader = &LogPollerReader{}

// LogPollerReader implements the Reader interface by using a logPoller instance to fetch the events.
type LogPollerReader struct {
	lp     logpoller.LogPoller
	lggr   logger.Logger
	client evmclient.Client

	dependencyCache sync.Map
}

func NewLogPollerReader(lp logpoller.LogPoller, lggr logger.Logger, client evmclient.Client) *LogPollerReader {
	return &LogPollerReader{
		lp:              lp,
		lggr:            lggr,
		client:          client,
		dependencyCache: sync.Map{},
	}
}

func (c *LogPollerReader) LatestBlock(ctx context.Context) (int64, error) {
	return c.lp.LatestBlock(pg.WithParentCtx(ctx))
}

func parseLogs[T any](logs []logpoller.Log, lggr logger.Logger, parseFunc func(log types.Log) (*T, error)) ([]Event[T], error) {
	reqs := make([]Event[T], 0, len(logs))
	for _, log := range logs {
		data, err := parseFunc(log.ToGethLog())
		if err == nil {
			reqs = append(reqs, Event[T]{
				Data: *data,
				Meta: Meta{
					BlockTimestamp: log.BlockTimestamp,
					BlockNumber:    log.BlockNumber,
					TxHash:         log.TxHash,
					LogIndex:       uint(log.LogIndex),
				},
			})
		}
	}

	if len(logs) != len(reqs) {
		lggr.Warnw("Some logs were not parsed", "logs", len(logs), "requests", len(reqs))
	}
	return reqs, nil
}

func (c *LogPollerReader) loadOnRamp(addr common.Address) (*evm_2_evm_onramp.EVM2EVMOnRampFilterer, error) {
	onRamp, exists := loadCachedDependency[*evm_2_evm_onramp.EVM2EVMOnRampFilterer](&c.dependencyCache, addr)
	if exists {
		return onRamp, nil
	}

	onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRampFilterer(addr, c.client)
	if err != nil {
		return nil, err
	}

	c.dependencyCache.Store(addr, onRamp)
	return onRamp, nil
}

func (c *LogPollerReader) loadPriceRegistry(addr common.Address) (*price_registry.PriceRegistryFilterer, error) {
	priceRegistry, exists := loadCachedDependency[*price_registry.PriceRegistryFilterer](&c.dependencyCache, addr)
	if exists {
		return priceRegistry, nil
	}

	priceRegistry, err := price_registry.NewPriceRegistryFilterer(addr, c.client)
	if err != nil {
		return nil, err
	}

	c.dependencyCache.Store(addr, priceRegistry)
	return priceRegistry, nil
}

func (c *LogPollerReader) loadOffRamp(addr common.Address) (*evm_2_evm_offramp.EVM2EVMOffRampFilterer, error) {
	offRamp, exists := loadCachedDependency[*evm_2_evm_offramp.EVM2EVMOffRampFilterer](&c.dependencyCache, addr)
	if exists {
		return offRamp, nil
	}

	offRamp, err := evm_2_evm_offramp.NewEVM2EVMOffRampFilterer(addr, c.client)
	if err != nil {
		return nil, err
	}

	c.dependencyCache.Store(addr, offRamp)
	return offRamp, nil
}

func (c *LogPollerReader) loadCommitStore(addr common.Address) (*commit_store.CommitStoreFilterer, error) {
	commitStore, exists := loadCachedDependency[*commit_store.CommitStoreFilterer](&c.dependencyCache, addr)
	if exists {
		return commitStore, nil
	}

	commitStore, err := commit_store.NewCommitStoreFilterer(addr, c.client)
	if err != nil {
		return nil, err
	}

	c.dependencyCache.Store(addr, commitStore)
	return commitStore, nil
}

func loadCachedDependency[T any](cache *sync.Map, addr common.Address) (T, bool) {
	var empty T

	if rawVal, exists := cache.Load(addr); exists {
		if dep, is := rawVal.(T); is {
			return dep, true
		}
	}

	return empty, false
}
