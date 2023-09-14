package streams

import (
	"context"
	"fmt"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-data-streams/streams"
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocr2plus "github.com/smartcontractkit/libocr/offchainreporting2plus"
	ocr3types "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var _ job.ServiceCtx = &delegate{}

type delegate struct {
	cfg    DelegateConfig
	codecs map[commontypes.StreamsReportFormat]streams.ReportCodec
}

type DelegateConfig struct {
	Logger  logger.Logger
	Queryer pg.Queryer
	Runner  Runner

	// OCR3
	BinaryNetworkEndpointFactory ocr2types.BinaryNetworkEndpointFactory
	V2Bootstrappers              []ocrcommontypes.BootstrapperLocator
	ContractConfigTracker        ocr2types.ContractConfigTracker
	ContractTransmitter          ocr3types.ContractTransmitter[commontypes.StreamsReportInfo]
	Database                     ocr3types.Database
	OCRLogger                    ocrcommontypes.Logger
	MonitoringEndpoint           ocrcommontypes.MonitoringEndpoint
	OffchainConfigDigester       ocr2types.OffchainConfigDigester
	OffchainKeyring              ocr2types.OffchainKeyring
	OnchainKeyring               ocr3types.OnchainKeyring[commontypes.StreamsReportInfo]
	LocalConfig                  ocr2types.LocalConfig
}

func NewDelegate(cfg DelegateConfig) job.ServiceCtx {
	// TODO: add the chain codecs here
	// TODO: nil checks?
	codecs := make(map[commontypes.StreamsReportFormat]streams.ReportCodec)
	return &delegate{cfg, codecs}
}

func (d *delegate) Start(ctx context.Context) error {
	// create the oracle from config values
	// TODO: Do these services need starting?
	prrc := streams.NewPredecessorRetirementReportCache()
	src := streams.NewShouldRetireCache()
	cdc := streams.NewChannelDefinitionCache()
	orm := NewORM(d.cfg.Queryer)
	sc := NewStreamCache(orm)
	if err := sc.Load(ctx, d.cfg.Logger.Named("StreamCache"), d.cfg.Runner); err != nil {
		return err
	}
	ds := NewDataSource(d.cfg.Logger.Named("DataSource"), sc)
	llo, err := ocr2plus.NewOracle(ocr2plus.OCR3OracleArgs[commontypes.StreamsReportInfo]{
		BinaryNetworkEndpointFactory: d.cfg.BinaryNetworkEndpointFactory,
		V2Bootstrappers:              d.cfg.V2Bootstrappers,
		ContractConfigTracker:        d.cfg.ContractConfigTracker,
		ContractTransmitter:          d.cfg.ContractTransmitter,
		Database:                     d.cfg.Database,
		LocalConfig:                  d.cfg.LocalConfig,
		Logger:                       d.cfg.OCRLogger,
		MonitoringEndpoint:           d.cfg.MonitoringEndpoint,
		OffchainConfigDigester:       d.cfg.OffchainConfigDigester,
		OffchainKeyring:              d.cfg.OffchainKeyring,
		OnchainKeyring:               d.cfg.OnchainKeyring,
		ReportingPluginFactory: streams.NewPluginFactory(
			prrc, src, cdc, ds, d.cfg.Logger.Named("StreamsReportingPlugin"), d.codecs,
		),
	})

	if err != nil {
		return fmt.Errorf("%w: failed to create new OCR oracle", err)
	}

	return llo.Start()
}

func (d *delegate) Close() error {
	panic("TODO")
}
