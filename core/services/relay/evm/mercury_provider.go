package evm

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	mercurytypes "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	v1 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	v2 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
	v4 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v4"
	"github.com/smartcontractkit/chainlink-data-streams/mercury"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury/config"
	evmmercury "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	reportcodecv1 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1/reportcodec"
	reportcodecv2 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v2/reportcodec"
	reportcodecv3 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
	reportcodecv4 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v4/reportcodec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var _ commontypes.MercuryProvider = (*mercuryProvider)(nil)

type mercuryProvider struct {
	cp                 commontypes.ConfigProvider
	codec              commontypes.Codec
	csaSigner          *coretypes.Ed25519Signer
	transmitter        evmmercury.Transmitter
	reportCodecV1      v1.ReportCodec
	reportCodecV2      v2.ReportCodec
	reportCodecV3      v3.ReportCodec
	reportCodecV4      v4.ReportCodec
	mercuryChainReader mercurytypes.ChainReader
	logger             logger.Logger
	ms                 services.MultiStart
}

func NewMercuryProvider(
	ctx context.Context,
	jobID int32,
	relayConfig types.RelayConfig,
	cfg config.PluginConfig,
	transmitterCfg evmmercury.TransmitterConfig,
	cp commontypes.ConfigProvider,
	codec commontypes.Codec,
	mercuryChainReader mercurytypes.ChainReader,
	lggr logger.SugaredLogger,
	csaKeystore coretypes.Keystore,
	mercuryPool wsrpc.Pool,
	mercuryORM evmmercury.ORM,
	triggerCapability *triggers.MercuryTriggerService,
) (*mercuryProvider, error) {
	reportCodecV1 := reportcodecv1.NewReportCodec(*relayConfig.FeedID, lggr.Named("ReportCodecV1"))
	reportCodecV2 := reportcodecv2.NewReportCodec(*relayConfig.FeedID, lggr.Named("ReportCodecV2"))
	reportCodecV3 := reportcodecv3.NewReportCodec(*relayConfig.FeedID, lggr.Named("ReportCodecV3"))
	reportCodecV4 := reportcodecv4.NewReportCodec(*relayConfig.FeedID, lggr.Named("ReportCodecV4"))

	getCodecForFeed := func(feedID mercuryutils.FeedID) (evmmercury.TransmitterReportDecoder, error) {
		var transmitterCodec evmmercury.TransmitterReportDecoder
		switch feedID.Version() {
		case 1:
			transmitterCodec = reportCodecV1
		case 2:
			transmitterCodec = reportCodecV2
		case 3:
			transmitterCodec = reportCodecV3
		case 4:
			transmitterCodec = reportCodecV4
		default:
			return nil, fmt.Errorf("invalid feed version %d", feedID.Version())
		}
		return transmitterCodec, nil
	}

	benchmarkPriceDecoder := func(ctx context.Context, feedID mercuryutils.FeedID, report ocrtypes.Report) (*big.Int, error) {
		benchmarkPriceCodec, benchmarkPriceErr := getCodecForFeed(feedID)
		if benchmarkPriceErr != nil {
			return nil, benchmarkPriceErr
		}
		return benchmarkPriceCodec.BenchmarkPriceFromReport(ctx, report)
	}

	csaPub := relayConfig.EffectiveTransmitterID.String
	csaSigner, err := coretypes.NewEd25519Signer(csaPub, csaKeystore.Sign)
	if err != nil {
		return nil, fmt.Errorf("failed to create ed25519 signer: %w", err)
	}
	clients := make(map[string]wsrpc.Client)
	for _, server := range cfg.GetServers() {
		clients[server.URL], err = mercuryPool.Checkout(ctx, csaPub, csaSigner, server.PubKey, server.URL)
		if err != nil {
			return nil, err
		}
	}
	transmitterCodec, err := getCodecForFeed(mercuryutils.FeedID(*relayConfig.FeedID))
	if err != nil {
		return nil, err
	}
	transmitter := evmmercury.NewTransmitter(lggr, transmitterCfg, clients, csaPub, jobID, *relayConfig.FeedID, mercuryORM, transmitterCodec, benchmarkPriceDecoder, triggerCapability)
	return &mercuryProvider{
		cp,
		codec,
		csaSigner,
		transmitter,
		reportCodecV1,
		reportCodecV2,
		reportCodecV3,
		reportCodecV4,
		mercuryChainReader,
		lggr,
		services.MultiStart{},
	}, nil
}

func (p *mercuryProvider) Start(ctx context.Context) error {
	return p.ms.Start(ctx, p.cp, p.transmitter, p.csaSigner)
}

func (p *mercuryProvider) Close() error {
	return p.ms.Close()
}

func (p *mercuryProvider) Ready() error {
	return errors.Join(p.cp.Ready(), p.transmitter.Ready())
}

func (p *mercuryProvider) Name() string {
	return p.logger.Name()
}

func (p *mercuryProvider) HealthReport() map[string]error {
	report := map[string]error{}
	services.CopyHealth(report, p.cp.HealthReport())
	services.CopyHealth(report, p.transmitter.HealthReport())
	return report
}

func (p *mercuryProvider) MercuryChainReader() mercurytypes.ChainReader {
	return p.mercuryChainReader
}

func (p *mercuryProvider) Codec() commontypes.Codec {
	return p.codec
}

func (p *mercuryProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return p.cp.ContractConfigTracker()
}

func (p *mercuryProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return p.cp.OffchainConfigDigester()
}

func (p *mercuryProvider) OnchainConfigCodec() mercurytypes.OnchainConfigCodec {
	return mercury.StandardOnchainConfigCodec{}
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

func (p *mercuryProvider) ReportCodecV4() v4.ReportCodec {
	return p.reportCodecV4
}

func (p *mercuryProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return p.transmitter
}

func (p *mercuryProvider) MercuryServerFetcher() mercurytypes.ServerFetcher {
	return p.transmitter
}

func (p *mercuryProvider) ContractReader() commontypes.ContractReader {
	return nil
}

var _ mercurytypes.ChainReader = (*mercuryChainReader)(nil)

type mercuryChainReader struct {
	tracker heads.Tracker
}

func NewChainReader(h heads.Tracker) mercurytypes.ChainReader {
	return &mercuryChainReader{h}
}

func NewMercuryChainReader(h heads.Tracker) mercurytypes.ChainReader {
	return &mercuryChainReader{
		tracker: h,
	}
}

func (r *mercuryChainReader) LatestHeads(ctx context.Context, k int) ([]mercurytypes.Head, error) {
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
