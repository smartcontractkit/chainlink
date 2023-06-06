package evm

import (
	"context"
	"errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/exp/maps"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
)

var _ relaytypes.MercuryProvider = (*mercuryProvider)(nil)

type mercuryProvider struct {
	configWatcher *configWatcher
	transmitter   mercury.Transmitter
	reportCodec   relaymercury.ReportCodec
	logger        logger.Logger

	ms services.MultiStart
}

func NewMercuryProvider(
	configWatcher *configWatcher,
	transmitter mercury.Transmitter,
	reportCodec relaymercury.ReportCodec,
	lggr logger.Logger,
) *mercuryProvider {
	return &mercuryProvider{
		configWatcher,
		transmitter,
		reportCodec,
		lggr,
		services.MultiStart{},
	}
}

func (p *mercuryProvider) Start(ctx context.Context) error {
	return p.ms.Start(ctx, p.configWatcher, p.transmitter)
}

func (p *mercuryProvider) Close() error {
	return p.ms.Close()
}

func (p *mercuryProvider) Ready() error {
	return errors.Join(p.configWatcher.Ready(), p.transmitter.Ready())
}

func (p *mercuryProvider) Name() string {
	return p.logger.Name()
}

func (p *mercuryProvider) HealthReport() map[string]error {
	report := map[string]error{}
	maps.Copy(report, p.configWatcher.HealthReport())
	maps.Copy(report, p.transmitter.HealthReport())
	return report
}

func (p *mercuryProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return p.configWatcher.ContractConfigTracker()
}

func (p *mercuryProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return p.configWatcher.OffchainConfigDigester()
}

func (p *mercuryProvider) OnchainConfigCodec() relaymercury.OnchainConfigCodec {
	return relaymercury.StandardOnchainConfigCodec{}
}

func (p *mercuryProvider) ReportCodec() relaymercury.ReportCodec {
	return p.reportCodec
}

func (p *mercuryProvider) ContractTransmitter() relaymercury.Transmitter {
	return p.transmitter
}
