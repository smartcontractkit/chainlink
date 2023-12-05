package evm

import (
	"context"
	"errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	mercurytypes "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	v1 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	v2 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
	relaymercury "github.com/smartcontractkit/chainlink-data-streams/mercury"

	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
)

var _ types.MercuryProvider = (*mercuryProvider)(nil)

type mercuryProvider struct {
	configWatcher *configWatcher
	transmitter   mercury.Transmitter
	reportCodecV1 v1.ReportCodec
	reportCodecV2 v2.ReportCodec
	reportCodecV3 v3.ReportCodec
	chainReader   mercurytypes.ChainReader
	logger        logger.Logger

	ms services.MultiStart
}

func NewMercuryProvider(
	configWatcher *configWatcher,
	transmitter mercury.Transmitter,
	reportCodecV1 v1.ReportCodec,
	reportCodecV2 v2.ReportCodec,
	reportCodecV3 v3.ReportCodec,
	chainReader mercurytypes.ChainReader,
	lggr logger.Logger,
) *mercuryProvider {
	return &mercuryProvider{
		configWatcher,
		transmitter,
		reportCodecV1,
		reportCodecV2,
		reportCodecV3,
		chainReader,
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

func (p *mercuryProvider) OnchainConfigCodec() mercurytypes.OnchainConfigCodec {
	return relaymercury.StandardOnchainConfigCodec{}
}

func (p *mercuryProvider) ReportCodecV1() v1.ReportCodec {
	return p.reportCodecV1
}

func (p *mercuryProvider) ReportCodecV2() v2.ReportCodec {
	return p.reportCodecV2
}

func (p *mercuryProvider) ReportCodecV3() v3.ReportCodec {
	return p.reportCodecV3
}

func (p *mercuryProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return p.transmitter
}

func (p *mercuryProvider) MercuryServerFetcher() mercurytypes.ServerFetcher {
	return p.transmitter
}

func (p *mercuryProvider) ChainReader() mercurytypes.ChainReader {
	return p.chainReader
}

var _ mercurytypes.ChainReader = (*chainReader)(nil)

type chainReader struct {
	tracker httypes.HeadTracker
}

func NewChainReader(h httypes.HeadTracker) mercurytypes.ChainReader {
	return &chainReader{
		tracker: h,
	}
}

func (r *chainReader) LatestHeads(ctx context.Context, k int) ([]mercurytypes.Head, error) {
	evmBlocks := r.tracker.LatestChain().AsSlice(k)
	if len(evmBlocks) == 0 {
		return nil, nil
	}

	blocks := make([]mercurytypes.Head, len(evmBlocks))
	for x := 0; x < len(evmBlocks); x++ {
		blocks[x] = mercurytypes.Head{
			Number:    uint64(evmBlocks[x].BlockNumber()),
			Hash:      evmBlocks[x].Hash.Bytes(),
			Timestamp: uint64(evmBlocks[x].Timestamp.Unix()),
		}
	}

	return blocks, nil
}
