package evm

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type pluginProvider struct {
	services.Service
	chainReader         types.ContractReader
	codec               types.Codec
	contractTransmitter ocrtypes.ContractTransmitter
	configWatcher       *configWatcher
	lggr                logger.Logger
	ms                  services.MultiStart
}

var _ types.PluginProvider = (*pluginProvider)(nil)

func NewPluginProvider(
	chainReader types.ContractReader,
	codec types.Codec,
	contractTransmitter ocrtypes.ContractTransmitter,
	configWatcher *configWatcher,
	lggr logger.Logger,
) *pluginProvider {
	return &pluginProvider{
		chainReader:         chainReader,
		codec:               codec,
		contractTransmitter: contractTransmitter,
		configWatcher:       configWatcher,
		lggr:                lggr,
		ms:                  services.MultiStart{},
	}
}

func (p *pluginProvider) Name() string { return p.lggr.Name() }

func (p *pluginProvider) Ready() error { return nil }

func (p *pluginProvider) HealthReport() map[string]error {
	hp := map[string]error{p.Name(): p.Ready()}
	services.CopyHealth(hp, p.configWatcher.HealthReport())
	return hp
}

func (p *pluginProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return p.contractTransmitter
}

func (p *pluginProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return p.configWatcher.OffchainConfigDigester()
}

func (p *pluginProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return p.configWatcher.configPoller
}

func (p *pluginProvider) ChainReader() types.ContractReader {
	return p.chainReader
}

func (p *pluginProvider) Codec() types.Codec {
	return p.codec
}

func (p *pluginProvider) Start(ctx context.Context) error {
	return p.configWatcher.Start(ctx)
}

func (p *pluginProvider) Close() error {
	return p.configWatcher.Close()
}
