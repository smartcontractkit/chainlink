package llo

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lloconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/llo/config"
)

type ChannelDefinitionCacheFactory interface {
	NewCache(cfg lloconfig.PluginConfig) (llotypes.ChannelDefinitionCache, error)
}

var _ ChannelDefinitionCacheFactory = &channelDefinitionCacheFactory{}

func NewChannelDefinitionCacheFactory(lggr logger.Logger, orm ChannelDefinitionCacheORM, lp logpoller.LogPoller, client *http.Client) ChannelDefinitionCacheFactory {
	return &channelDefinitionCacheFactory{
		lggr,
		orm,
		lp,
		client,
		make(map[common.Address]map[uint32]struct{}),
		sync.Mutex{},
	}
}

type channelDefinitionCacheFactory struct {
	lggr   logger.Logger
	orm    ChannelDefinitionCacheORM
	lp     logpoller.LogPoller
	client *http.Client

	caches map[common.Address]map[uint32]struct{}
	mu     sync.Mutex
}

// TODO: Test this
// MERC-3653
func (f *channelDefinitionCacheFactory) NewCache(cfg lloconfig.PluginConfig) (llotypes.ChannelDefinitionCache, error) {
	if cfg.ChannelDefinitions != "" {
		return NewStaticChannelDefinitionCache(f.lggr, cfg.ChannelDefinitions)
	}

	addr := cfg.ChannelDefinitionsContractAddress
	fromBlock := cfg.ChannelDefinitionsContractFromBlock
	donID := cfg.DonID

	f.mu.Lock()
	defer f.mu.Unlock()

	if _, exists := f.caches[addr][donID]; exists {
		// This shouldn't really happen and isn't supported
		return nil, fmt.Errorf("cache already exists for contract address %s and don ID %d", addr.Hex(), donID)
	}
	if _, exists := f.caches[addr]; !exists {
		f.caches[addr] = make(map[uint32]struct{})
	}
	f.caches[addr][donID] = struct{}{}
	return NewChannelDefinitionCache(f.lggr, f.orm, f.client, f.lp, addr, donID, fromBlock), nil
}
