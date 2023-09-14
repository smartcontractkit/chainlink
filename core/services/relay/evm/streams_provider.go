package evm

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	datastreams "github.com/smartcontractkit/chainlink-data-streams/llo"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo"
)

var _ commontypes.LLOProvider = (*streamsProvider)(nil)

type streamsProvider struct {
	configWatcher          *configWatcher
	transmitter            llo.Transmitter
	logger                 logger.Logger
	channelDefinitionCache commontypes.ChannelDefinitionCache

	ms services.MultiStart
}

func NewLLOProvider(
	configWatcher *configWatcher,
	transmitter llo.Transmitter,
	lggr logger.Logger,
	channelDefinitionCache commontypes.ChannelDefinitionCache,
) relaytypes.LLOProvider {
	return &streamsProvider{
		configWatcher,
		transmitter,
		lggr,
		channelDefinitionCache,
		services.MultiStart{},
	}
}

func (p *streamsProvider) Start(ctx context.Context) error {
	return p.ms.Start(ctx, p.configWatcher, p.transmitter, p.channelDefinitionCache)
}

func (p *streamsProvider) Close() error {
	return p.ms.Close()
}

func (p *streamsProvider) Ready() error {
	return errors.Join(p.configWatcher.Ready(), p.transmitter.Ready(), p.channelDefinitionCache.Ready())
}

func (p *streamsProvider) Name() string {
	return p.logger.Name()
}

func (p *streamsProvider) HealthReport() map[string]error {
	report := map[string]error{}
	services.CopyHealth(report, p.configWatcher.HealthReport())
	services.CopyHealth(report, p.transmitter.HealthReport())
	services.CopyHealth(report, p.channelDefinitionCache.HealthReport())
	return report
}

func (p *streamsProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return p.configWatcher.ContractConfigTracker()
}

func (p *streamsProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return p.configWatcher.OffchainConfigDigester()
}

func (p *streamsProvider) OnchainConfigCodec() datastreams.OnchainConfigCodec {
	// TODO: This should probably be moved to core since its chain-specific
	return &datastreams.JSONOnchainConfigCodec{}
}

func (p *streamsProvider) ContractTransmitter() commontypes.LLOTransmitter {
	return p.transmitter
}

func (p *streamsProvider) ChannelDefinitionCache() commontypes.ChannelDefinitionCache {
	return p.channelDefinitionCache
}
