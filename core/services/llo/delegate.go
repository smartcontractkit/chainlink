package llo

import (
	"context"
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	ocr2plus "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
)

var _ job.ServiceCtx = &delegate{}

type Closer interface {
	Close() error
}

type delegate struct {
	services.StateMachine

	cfg    DelegateConfig
	codecs map[llotypes.ReportFormat]llo.ReportCodec

	prrc llo.PredecessorRetirementReportCache
	src  llo.ShouldRetireCache
	ds   llo.DataSource

	oracle Closer
}

type DelegateConfig struct {
	Logger     logger.Logger
	DataSource sqlutil.DataSource
	Runner     streams.Runner
	Registry   Registry
	JobName    null.String

	// LLO
	ChannelDefinitionCache llotypes.ChannelDefinitionCache

	// OCR3
	BinaryNetworkEndpointFactory ocr2types.BinaryNetworkEndpointFactory
	V2Bootstrappers              []ocrcommontypes.BootstrapperLocator
	ContractConfigTracker        ocr2types.ContractConfigTracker
	ContractTransmitter          ocr3types.ContractTransmitter[llotypes.ReportInfo]
	Database                     ocr3types.Database
	OCRLogger                    ocrcommontypes.Logger
	MonitoringEndpoint           ocrcommontypes.MonitoringEndpoint
	OffchainConfigDigester       ocr2types.OffchainConfigDigester
	OffchainKeyring              ocr2types.OffchainKeyring
	OnchainKeyring               ocr3types.OnchainKeyring[llotypes.ReportInfo]
	LocalConfig                  ocr2types.LocalConfig
}

func NewDelegate(cfg DelegateConfig) (job.ServiceCtx, error) {
	if cfg.DataSource == nil {
		return nil, errors.New("DataSource must not be nil")
	}
	if cfg.Runner == nil {
		return nil, errors.New("Runner must not be nil")
	}
	if cfg.Registry == nil {
		return nil, errors.New("Registry must not be nil")
	}
	codecs := make(map[llotypes.ReportFormat]llo.ReportCodec)

	// NOTE: All codecs must be specified here
	codecs[llotypes.ReportFormatJSON] = llo.JSONReportCodec{}
	codecs[llotypes.ReportFormatEVM] = evm.ReportCodec{}

	// TODO: Do these services need starting?
	// https://smartcontract-it.atlassian.net/browse/MERC-3386
	prrc := llo.NewPredecessorRetirementReportCache()
	src := llo.NewShouldRetireCache()
	ds := newDataSource(cfg.Logger.Named("DataSource"), cfg.Registry)

	return &delegate{services.StateMachine{}, cfg, codecs, prrc, src, ds, nil}, nil
}

func (d *delegate) Start(ctx context.Context) error {
	return d.StartOnce("LLODelegate", func() error {
		// create the oracle from config values
		oracle, err := ocr2plus.NewOracle(ocr2plus.OCR3OracleArgs[llotypes.ReportInfo]{
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
			ReportingPluginFactory: llo.NewPluginFactory(
				d.prrc, d.src, d.cfg.ChannelDefinitionCache, d.ds, d.cfg.Logger.Named("LLOReportingPlugin"), d.codecs,
			),
			MetricsRegisterer: prometheus.WrapRegistererWith(map[string]string{"job_name": d.cfg.JobName.ValueOrZero()}, prometheus.DefaultRegisterer),
		})

		if err != nil {
			return fmt.Errorf("%w: failed to create new OCR oracle", err)
		}

		d.oracle = oracle

		return oracle.Start()
	})
}

func (d *delegate) Close() error {
	return d.StopOnce("LLODelegate", d.oracle.Close)
}
