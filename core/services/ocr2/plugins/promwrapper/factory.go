package promwrapper

import (
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ types.ReportingPluginFactory = &promFactory{}

type promFactory struct {
	wrapped    types.ReportingPluginFactory
	pluginName string
	evmChainID *big.Int
}

func (p *promFactory) NewReportingPlugin(config types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	plugin, info, err := p.wrapped.NewReportingPlugin(config)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}

	prom := New(plugin, p.pluginName, p.evmChainID)
	return prom, info, nil
}

func NewPromFactory(wrapped types.ReportingPluginFactory, pluginName string, evmChainID *big.Int) types.ReportingPluginFactory {
	return &promFactory{
		wrapped:    wrapped,
		pluginName: pluginName,
		evmChainID: evmChainID,
	}
}
