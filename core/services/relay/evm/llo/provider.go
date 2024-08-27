package llo

import (
	"context"
	"errors"

	pkgerrors "github.com/pkg/errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo"
	lloconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/llo/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var _ commontypes.LLOProvider = (*lloProvider)(nil)

type lloProvider struct {
	cp                     commontypes.ConfigProvider
	transmitter            llo.Transmitter
	logger                 logger.Logger
	channelDefinitionCache llotypes.ChannelDefinitionCache

	ms services.MultiStart
}

func NewProvider(
	lggr logger.Logger,
	chain legacyevm.Chain,
	relayCfg types.RelayConfig,
	relayOpts *types.RelayOpts,
	lloCfg lloconfig.PluginConfig,
	cdc llotypes.ChannelDefinitionCache,
	transmitter llo.Transmitter,
) (relaytypes.LLOProvider, error) {
	lggr = logger.Named(lggr, "LLOProvider")
	cp, err := NewConfigProvider(lggr, chain, relayCfg, relayOpts)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	if !relayCfg.EffectiveTransmitterID.Valid {
		return nil, pkgerrors.New("EffectiveTransmitterID must be specified")
	}

	return &lloProvider{
		cp,
		transmitter,
		lggr,
		cdc,
		services.MultiStart{},
	}, nil
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
