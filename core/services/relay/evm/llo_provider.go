package evm

import (
	"context"
	"errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo"
)

var _ commontypes.LLOProvider = (*lloProvider)(nil)

type lloProvider struct {
	cp                     commontypes.ConfigProvider
	transmitter            llo.Transmitter
	logger                 logger.Logger
	channelDefinitionCache llotypes.ChannelDefinitionCache

	ms services.MultiStart
}

func NewLLOProvider(
	cp commontypes.ConfigProvider,
	transmitter llo.Transmitter,
	lggr logger.Logger,
	channelDefinitionCache llotypes.ChannelDefinitionCache,
) relaytypes.LLOProvider {
	return &lloProvider{
		cp,
		transmitter,
		lggr.Named("LLOProvider"),
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

func (p *lloProvider) OnchainConfigCodec() datastreamsllo.OnchainConfigCodec {
	// TODO: This should probably be moved to core since its chain-specific
	// https://smartcontract-it.atlassian.net/browse/MERC-3661
	return &datastreamsllo.JSONOnchainConfigCodec{}
}

func (p *lloProvider) ContractTransmitter() llotypes.Transmitter {
	return p.transmitter
}

func (p *lloProvider) ChannelDefinitionCache() llotypes.ChannelDefinitionCache {
	return p.channelDefinitionCache
}
