package channeldefinitions

import (
	"net/http"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
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
	}
}

type channelDefinitionCacheFactory struct {
	lggr   logger.Logger
	orm    ChannelDefinitionCacheORM
	lp     logpoller.LogPoller
	client *http.Client
}

func (f *channelDefinitionCacheFactory) NewCache(cfg lloconfig.PluginConfig) (llotypes.ChannelDefinitionCache, error) {
	if cfg.ChannelDefinitions != "" {
		return NewStaticChannelDefinitionCache(f.lggr, cfg.ChannelDefinitions)
	}

	addr := cfg.ChannelDefinitionsContractAddress
	fromBlock := cfg.ChannelDefinitionsContractFromBlock
	donID := cfg.DonID

	return NewChannelDefinitionCache(f.lggr, f.orm, f.client, f.lp, addr, donID, fromBlock), nil
}
