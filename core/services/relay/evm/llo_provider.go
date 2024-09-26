package evm

import (
	"context"
	"errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

var _ commontypes.LLOProvider = (*lloProvider)(nil)

type LLOTransmitter interface {
	services.Service
	llotypes.Transmitter
}

type lloProvider struct {
	cp                     commontypes.ConfigProvider
	transmitter            LLOTransmitter
	logger                 logger.Logger
	channelDefinitionCache llotypes.ChannelDefinitionCache

	ms services.MultiStart
}

func NewLLOProvider(
	cp commontypes.ConfigProvider,
	transmitter LLOTransmitter,
	lggr logger.Logger,
	channelDefinitionCache llotypes.ChannelDefinitionCache,
) relaytypes.LLOProvider {
	return &lloProvider{
		cp,
		transmitter,
		logger.Named(lggr, "LLOProvider"),
		channelDefinitionCache,
		services.MultiStart{},
	}
}

func (p *lloProvider) Start(ctx context.Context) error {
	err := p.ms.Start(ctx, p.cp, p.transmitter, p.channelDefinitionCache)
	return err
}

func (p *lloProvider) Close() error {
	return p.ms.Close()
}

func (p *lloProvider) Ready() error {
	return errors.Join(p.cp.Ready(), p.transmitter.Ready(), p.channelDefinitionCache.Ready())
}

func (p *lloProvider) Name() string {
	return p.logger.Name()
}

func (p *lloProvider) HealthReport() map[string]error {
	report := map[string]error{}
	services.CopyHealth(report, p.cp.HealthReport())
	services.CopyHealth(report, p.transmitter.HealthReport())
	services.CopyHealth(report, p.channelDefinitionCache.HealthReport())
	return report
}

func (p *lloProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return p.cp.ContractConfigTracker()
}

func (p *lloProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return p.cp.OffchainConfigDigester()
}

func (p *lloProvider) ContractTransmitter() llotypes.Transmitter {
	return p.transmitter
}

func (p *lloProvider) ChannelDefinitionCache() llotypes.ChannelDefinitionCache {
	return p.channelDefinitionCache
}
