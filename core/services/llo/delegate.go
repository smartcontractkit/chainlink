package llo

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/guregu/null.v4"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	ocr2plus "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3shims"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	llodatasource "github.com/smartcontractkit/chainlink-data-streams/llo/datasource"
	llov31 "github.com/smartcontractkit/chainlink-data-streams/llo/dev/v31"
	lloprotocol "github.com/smartcontractkit/chainlink-data-streams/llo/protocol"
	"github.com/smartcontractkit/chainlink-data-streams/llo/retirement"
	"github.com/smartcontractkit/chainlink-data-streams/llo/transmitter"
	llov30 "github.com/smartcontractkit/chainlink-data-streams/llo/v30"
	corelogger "github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/observation"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/promwrapper"
	promwrapper31 "github.com/smartcontractkit/chainlink/v2/core/services/ocr3_1/promwrapper"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

var _ job.ServiceCtx = &delegate{}

type Closer interface {
	Close() error
}

type delegate struct {
	services.StateMachine

	cfg          DelegateConfig
	reportCodecs map[llotypes.ReportFormat]lloprotocol.ReportCodec

	// src is the shared ShouldRetireCache. llov30.ShouldRetireCache and
	// llov31.ShouldRetireCache have identical method sets, so this value serves
	// both versions.
	src llov30.ShouldRetireCache
	// ds is the shared LLO data source (llodatasource.DataSource); v30 and v31 both
	// consume it, lifecycle gating is driven by the round's DSOpts.
	ds    llodatasource.DataSource
	telem telem.TelemeterService

	oracles []Closer
}

type DelegateConfig struct {
	Logger                      logger.Logger
	DataSource                  sqlutil.DataSource
	Runner                      streams.Runner
	Registry                    observation.Registry
	JobName                     null.String
	CaptureEATelemetry          bool
	CaptureObservationTelemetry bool
	CaptureOutcomeTelemetry     bool
	CaptureReportTelemetry      bool

	// LLO
	ChannelDefinitionCache   llotypes.ChannelDefinitionCache
	ReportingPluginConfig    llov30.Config
	RetirementReportCache    retirement.RetirementReportCache
	RetirementReportCodec    lloprotocol.RetirementReportCodec
	ShouldRetireCache        llov30.ShouldRetireCache
	PluginMonitoringEndpoint telemetry.MultitypeMonitoringEndpoint
	DonID                    uint32
	ChainID                  string

	// OCR3
	TraceLogging                 bool
	SampleTelemetry              bool
	BinaryNetworkEndpointFactory ocr2types.BinaryNetworkEndpointFactory
	V2Bootstrappers              []ocrcommontypes.BootstrapperLocator
	// One Oracle will be started for each ContractConfigTracker
	ContractConfigTrackers []ocr2types.ContractConfigTracker
	ContractTransmitter    ocr3types.ContractTransmitter[llotypes.ReportInfo]
	OCR3MonitoringEndpoint ocrcommontypes.MonitoringEndpoint
	OffchainConfigDigester ocr2types.OffchainConfigDigester
	OffchainKeyring        ocr2types.OffchainKeyring
	OnchainKeyring         ocr3types.OnchainKeyring[llotypes.ReportInfo]
	LocalConfig            ocr2types.LocalConfig
	NewOCR3DB              func(pluginID int32) ocr3types.Database

	// OCR3.1 (only required when OCR31 is true; see chainlink-data-streams
	// llo/config.PluginConfig.OCRVersion)
	OCR31 bool
	// BinaryNetworkEndpoint2Factory is the OCR3.1 ("2") network endpoint factory
	// (peerWrapper.Peer3_1). Required when OCR31 is true.
	BinaryNetworkEndpoint2Factory ocr2types.BinaryNetworkEndpoint2Factory
	// KeyValueDatabaseFactory provides the replicated per-configDigest key-value
	// store the OCR3.1 protocol requires. Required when OCR31 is true.
	KeyValueDatabaseFactory ocr3_1types.KeyValueDatabaseFactory
}

func NewDelegate(cfg DelegateConfig) (job.ServiceCtx, error) {
	lggr := logger.Sugared(cfg.Logger).With("jobName", cfg.JobName.ValueOrZero(), "donID", cfg.DonID)
	if cfg.DataSource == nil {
		return nil, errors.New("DataSource must not be nil")
	}
	if cfg.Runner == nil {
		return nil, errors.New("Runner must not be nil")
	}
	if cfg.Registry == nil {
		return nil, errors.New("Registry must not be nil")
	}
	if cfg.RetirementReportCache == nil {
		return nil, errors.New("RetirementReportCache must not be nil")
	}
	if cfg.ShouldRetireCache == nil {
		return nil, errors.New("ShouldRetireCache must not be nil")
	}
	if cfg.OCR31 {
		if cfg.KeyValueDatabaseFactory == nil {
			return nil, errors.New("KeyValueDatabaseFactory must not be nil when running OCR3.1")
		}
		if cfg.BinaryNetworkEndpoint2Factory == nil {
			return nil, errors.New("BinaryNetworkEndpoint2Factory must not be nil when running OCR3.1")
		}
	}
	var codecLggr logger.Logger
	if cfg.ReportingPluginConfig.VerboseLogging {
		codecLggr = logger.Named(lggr, "ReportCodecs")
	} else {
		codecLggr = corelogger.NullLogger
	}
	reportCodecs := NewReportCodecs(codecLggr, cfg.DonID)

	t := telem.NewTelemeterService(telem.TelemeterParams{
		Logger:                      lggr,
		MonitoringEndpoint:          cfg.PluginMonitoringEndpoint,
		DonID:                       cfg.DonID,
		CaptureEATelemetry:          cfg.CaptureEATelemetry,
		CaptureObservationTelemetry: cfg.CaptureObservationTelemetry,
		CaptureOutcomeTelemetry:     cfg.CaptureOutcomeTelemetry,
		CaptureReportTelemetry:      cfg.CaptureReportTelemetry,
		SampleTelemetry:             cfg.SampleTelemetry,
	})

	ds := observation.NewDataSource(logger.Named(lggr, "DataSource"), cfg.Registry, t)

	notifier, ok := cfg.ContractTransmitter.(transmitter.TransmitNotifier)
	if ok {
		notifier.OnTransmit(t.TrackSeqNr)
	}

	return &delegate{services.StateMachine{}, cfg, reportCodecs, cfg.ShouldRetireCache, ds, t, []Closer{}}, nil
}

func (d *delegate) Start(ctx context.Context) error {
	return d.StartOnce("LLODelegate", func() error {
		// create the oracle from config values
		if len(d.cfg.ContractConfigTrackers) != 1 && len(d.cfg.ContractConfigTrackers) != 2 {
			return fmt.Errorf("expected either 1 or 2 ContractConfigTrackers, got: %d", len(d.cfg.ContractConfigTrackers))
		}

		d.cfg.Logger.Debugw("Starting LLO job", "instances", len(d.cfg.ContractConfigTrackers), "jobName", d.cfg.JobName.ValueOrZero(), "captureEATelemetry", d.cfg.CaptureEATelemetry, "donID", d.cfg.DonID)

		var merr error

		merr = errors.Join(merr, d.telem.Start(ctx))

		psrrc := retirement.NewPluginScopedRetirementReportCache(d.cfg.RetirementReportCache, d.cfg.OnchainKeyring, d.cfg.RetirementReportCodec)
		for i, configTracker := range d.cfg.ContractConfigTrackers {
			lggr := logger.Named(d.cfg.Logger, strconv.Itoa(i))
			switch i {
			case 0:
				lggr = logger.With(lggr, "instanceType", "Blue")
			case 1:
				lggr = logger.With(lggr, "instanceType", "Green")
			}
			ocrLogger := logger.NewOCRWrapper(NewSuppressedLogger(lggr, d.cfg.TraceLogging, d.cfg.TraceLogging || d.cfg.ReportingPluginConfig.VerboseLogging), d.cfg.TraceLogging, func(msg string) {
				// NOTE: Some OCR loggers include a DB-persist here
				// We do not DB persist errors in LLO, since they could be quite voluminous and ought to be present in logs anyway.
				// This is a performance optimization
			})

			var oracle ocr2plus.Oracle
			var err error
			if d.cfg.OCR31 {
				oracle, err = d.newOracleV31(i, configTracker, lggr, ocrLogger, psrrc)
			} else {
				oracle, err = d.newOracleV30(i, configTracker, lggr, ocrLogger, psrrc)
			}
			if err != nil {
				return fmt.Errorf("%w: failed to create new OCR oracle", err)
			}

			d.oracles = append(d.oracles, oracle)

			merr = errors.Join(merr, oracle.Start())
		}

		return merr
	})
}

// newOracleV30 builds an OCR3.0 oracle running the llo/v30 reporting plugin.
func (d *delegate) newOracleV30(i int, configTracker ocr2types.ContractConfigTracker, lggr logger.Logger, ocrLogger ocrcommontypes.Logger, psrrc lloprotocol.PredecessorRetirementReportCache) (ocr2plus.Oracle, error) {
	return ocr2plus.NewOracle(ocr2plus.OCR3OracleArgs2[llotypes.ReportInfo]{
		BinaryNetworkEndpointFactory: d.cfg.BinaryNetworkEndpointFactory,
		V2Bootstrappers:              d.cfg.V2Bootstrappers,
		ContractConfigTracker:        configTracker,
		ContractTransmitter:          d.cfg.ContractTransmitter,
		Database:                     d.cfg.NewOCR3DB(int32(i)), //nolint:gosec // G115 // impossible due to ContractConfigTrackers length check
		LocalConfig:                  d.cfg.LocalConfig,
		Logger:                       ocrLogger,
		MonitoringEndpoint:           d.cfg.OCR3MonitoringEndpoint,
		OffchainConfigDigester:       d.cfg.OffchainConfigDigester,
		OffchainKeyring:              d.cfg.OffchainKeyring,
		OnchainKeyring:               ocr3shims.OnchainKeyringAsOnchainKeyring2(d.cfg.OnchainKeyring),
		ReportingPluginFactory: promwrapper.NewReportingPluginFactory(
			llov30.NewPluginFactory(
				llov30.PluginFactoryParams{
					Config:                           d.cfg.ReportingPluginConfig,
					PredecessorRetirementReportCache: psrrc,
					ShouldRetireCache:                d.src,
					RetirementReportCodec:            d.cfg.RetirementReportCodec,
					ChannelDefinitionCache:           d.cfg.ChannelDefinitionCache,
					DataSource:                       d.ds,
					Logger:                           logger.Named(lggr, "ReportingPlugin"),
					OnchainConfigCodec:               lloprotocol.EVMOnchainConfigCodec{},
					ReportCodecs:                     d.reportCodecs,
					OutcomeTelemetryCh:               d.telem.GetOutcomeTelemetryCh(),
					ReportTelemetryCh:                d.telem.GetReportTelemetryCh(),
					DonID:                            d.cfg.DonID,
				},
			),
			lggr,
			"",
			d.cfg.ChainID,
			"llo",
		),
		MetricsRegisterer: prometheus.WrapRegistererWith(map[string]string{"job_name": d.cfg.JobName.ValueOrZero()}, prometheus.DefaultRegisterer),
	})
}

// newOracleV31 builds an OCR3.1 oracle running the llo/v31 reporting plugin. It
// differs from v30 by the OCR3.1 oracle args (OCR3_1OracleArgs2), the "2"
// network endpoint factory, and the required replicated KeyValueDatabaseFactory.
func (d *delegate) newOracleV31(i int, configTracker ocr2types.ContractConfigTracker, lggr logger.Logger, ocrLogger ocrcommontypes.Logger, psrrc lloprotocol.PredecessorRetirementReportCache) (ocr2plus.Oracle, error) {
	factory := promwrapper31.NewReportingPluginFactory(
		llov31.NewPluginFactory(llov31.PluginFactoryParams{
			Config:                           llov31.Config{VerboseLogging: d.cfg.ReportingPluginConfig.VerboseLogging},
			PredecessorRetirementReportCache: psrrc,
			ShouldRetireCache:                d.src,
			RetirementReportCodec:            d.cfg.RetirementReportCodec,
			ChannelDefinitionCache:           d.cfg.ChannelDefinitionCache,
			DataSource:                       d.ds,
			Logger:                           logger.Named(lggr, "ReportingPlugin"),
			OnchainConfigCodec:               lloprotocol.EVMOnchainConfigCodec{},
			ReportCodecs:                     d.reportCodecs,
			OutcomeTelemetryCh:               d.telem.GetOutcomeTelemetryCh(),
			ReportTelemetryCh:                d.telem.GetReportTelemetryCh(),
			DonID:                            d.cfg.DonID,
			BlobThreshold:                    0, // 0 => llov31.DefaultBlobThreshold
		}),
		lggr,
		"",
		d.cfg.ChainID,
		"llo",
	)
	return ocr2plus.NewOracle(ocr2plus.OCR3_1OracleArgs2[llotypes.ReportInfo]{
		BinaryNetworkEndpointFactory: d.cfg.BinaryNetworkEndpoint2Factory,
		V2Bootstrappers:              d.cfg.V2Bootstrappers,
		ContractConfigTracker:        configTracker,
		ContractTransmitter:          d.cfg.ContractTransmitter,
		Database:                     d.cfg.NewOCR3DB(int32(i)), //nolint:gosec // G115 // impossible due to ContractConfigTrackers length check
		KeyValueDatabaseFactory:      d.cfg.KeyValueDatabaseFactory,
		LocalConfig:                  d.cfg.LocalConfig,
		Logger:                       ocrLogger,
		MonitoringEndpoint:           d.cfg.OCR3MonitoringEndpoint,
		OffchainConfigDigester:       d.cfg.OffchainConfigDigester,
		OffchainKeyring:              d.cfg.OffchainKeyring,
		OnchainKeyring:               ocr3shims.OnchainKeyringAsOnchainKeyring2(d.cfg.OnchainKeyring),
		ReportingPluginFactory:       factory,
		MetricsRegisterer:            prometheus.WrapRegistererWith(map[string]string{"job_name": d.cfg.JobName.ValueOrZero()}, prometheus.DefaultRegisterer),
	})
}

func (d *delegate) Close() error {
	return d.StopOnce("LLODelegate", func() (merr error) {
		for _, oracle := range d.oracles {
			merr = errors.Join(merr, oracle.Close())
		}
		if closer, ok := d.ds.(Closer); ok {
			merr = errors.Join(merr, closer.Close())
		}
		merr = errors.Join(merr, d.telem.Close())
		return merr
	})
}
