package streams

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-data-streams/streams"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

type ChannelDefinitionCacheFactory interface {
	NewCache(addr common.Address, fromBlock int64) streams.ChannelDefinitionCache
}

var _ ChannelDefinitionCacheFactory = &channelDefinitionCacheFactory{}

type channelDefinitionCacheFactory struct {
	lggr logger.Logger
	orm  ChannelDefinitionCacheORM // TODO: pass in a pre-scoped ORM (to EVM chain ID)
	lp   logpoller.LogPoller

	caches map[common.Address]struct{}
	mu     sync.Mutex
}

func (f *channelDefinitionCacheFactory) NewCache(addr common.Address, fromBlock int64) streams.ChannelDefinitionCache {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, exists := f.caches[common.Address]; exists {
		// TODO: can we do better?
		panic("cannot create duplicate cache")
	}
	f.caches[addr] = struct{}{}
	return NewChannelDefinitionCache(f.lggr, f.orm, f.lp, addr, fromBlock)
}
