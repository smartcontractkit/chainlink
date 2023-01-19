package evm

import (
	"context"

	"go.uber.org/multierr"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm/mercury"
)

var _ relaytypes.MercuryProvider = (*mercuryProvider)(nil)

type mercuryProvider struct {
	configWatcher *configWatcher
	transmitter   mercury.Transmitter
	reportCodec   relaymercury.ReportCodec

	ms services.MultiStart
}

func (p *mercuryProvider) Start(ctx context.Context) error {
	return p.ms.Start(ctx, p.configWatcher, p.transmitter)
}

func (p *mercuryProvider) Close() error {
	return p.ms.Close()
}

func (p *mercuryProvider) Ready() error {
	return multierr.Combine(p.configWatcher.Ready(), p.transmitter.Ready())
}

func (p *mercuryProvider) Healthy() error {
	return multierr.Combine(p.configWatcher.Healthy(), p.transmitter.Healthy())
}

func (p *mercuryProvider) Name() string {
	return "EVM.MercuryProvider"
}

func (p *mercuryProvider) HealthReport() map[string]error {
	return map[string]error{p.Name(): p.Healthy()}
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
