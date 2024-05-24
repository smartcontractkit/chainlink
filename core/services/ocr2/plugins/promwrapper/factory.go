package promwrapper

import (
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ types.ReportingPluginFactory = &promFactory{}

type promFactory struct {
	wrapped   types.ReportingPluginFactory
	name      string
	chainType string
	chainID   *big.Int
}

func (p *promFactory) NewReportingPlugin(config types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	plugin, info, err := p.wrapped.NewReportingPlugin(config)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}

	prom := New(plugin, p.name, p.chainType, p.chainID, config, nil)
	return prom, info, nil
}

func NewPromFactory(wrapped types.ReportingPluginFactory, name, chainType string, chainID *big.Int) types.ReportingPluginFactory {
	return &promFactory{
		wrapped:   wrapped,
		name:      name,
		chainType: chainType,
		chainID:   chainID,
	}
}
