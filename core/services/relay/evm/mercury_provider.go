package evm

import (
	"context"
	"errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	relaymercuryv1 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v1"
	relaymercuryv2 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v2"
	relaymercuryv3 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v3"
	"github.com/smartcontractkit/chainlink-relay/pkg/services"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
)

var _ relaytypes.MercuryProvider = (*mercuryProvider)(nil)

type mercuryProvider struct {
	configWatcher *configWatcher
	transmitter   mercury.Transmitter
	reportCodecV1 relaymercuryv1.ReportCodec
	reportCodecV2 relaymercuryv2.ReportCodec
	reportCodecV3 relaymercuryv3.ReportCodec
	logger        logger.Logger

	ms services.MultiStart
}

func NewMercuryProvider(
	configWatcher *configWatcher,
	transmitter mercury.Transmitter,
	reportCodecV1 relaymercuryv1.ReportCodec,
	reportCodecV2 relaymercuryv2.ReportCodec,
	reportCodecV3 relaymercuryv3.ReportCodec,
	lggr logger.Logger,
) *mercuryProvider {
	return &mercuryProvider{
		configWatcher,
		transmitter,
		reportCodecV1,
		reportCodecV2,
		reportCodecV3,
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
	services.CopyHealth(report, p.configWatcher.HealthReport())
	services.CopyHealth(report, p.transmitter.HealthReport())
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

func (p *mercuryProvider) ReportCodecV1() relaymercuryv1.ReportCodec {
	return p.reportCodecV1
}

func (p *mercuryProvider) ReportCodecV2() relaymercuryv2.ReportCodec {
	return p.reportCodecV2
}

func (p *mercuryProvider) ReportCodecV3() relaymercuryv3.ReportCodec {
	return p.reportCodecV3
}

func (p *mercuryProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return p.transmitter
}

func (p *mercuryProvider) MercuryServerFetcher() relaymercury.MercuryServerFetcher {
	return p.transmitter
}
