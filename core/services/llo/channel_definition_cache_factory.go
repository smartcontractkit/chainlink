package llo

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type ChannelDefinitionCacheFactory interface {
	NewCache(addr common.Address, fromBlock int64) commontypes.ChannelDefinitionCache
}

var _ ChannelDefinitionCacheFactory = &channelDefinitionCacheFactory{}

func NewChannelDefinitionCacheFactory(lggr logger.Logger, orm ChannelDefinitionCacheORM, lp logpoller.LogPoller) ChannelDefinitionCacheFactory {
	return &channelDefinitionCacheFactory{
		lggr,
		orm,
		lp,
		make(map[common.Address]struct{}),
		sync.Mutex{},
	}
}

type channelDefinitionCacheFactory struct {
	lggr logger.Logger
	orm  ChannelDefinitionCacheORM // TODO: pass in a pre-scoped ORM (to EVM chain ID)
	lp   logpoller.LogPoller

	caches map[common.Address]struct{}
	mu     sync.Mutex
}

func (f *channelDefinitionCacheFactory) NewCache(addr common.Address, fromBlock int64) commontypes.ChannelDefinitionCache {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, exists := f.caches[addr]; exists {
		// TODO: can we do better?
		panic("cannot create duplicate cache")
	}
	f.caches[addr] = struct{}{}
	return NewChannelDefinitionCache(f.lggr, f.orm, f.lp, addr, fromBlock)
}
