package llo

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	ocr2plus "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"

	corelogger "github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/observation"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/retirement"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/promwrapper"
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
	reportCodecs map[llotypes.ReportFormat]datastreamsllo.ReportCodec

	src   datastreamsllo.ShouldRetireCache
	ds    datastreamsllo.DataSource
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
	ReportingPluginConfig    datastreamsllo.Config
	RetirementReportCache    retirement.RetirementReportCache
	RetirementReportCodec    datastreamsllo.RetirementReportCodec
	ShouldRetireCache        datastreamsllo.ShouldRetireCache
	PluginMonitoringEndpoint telemetry.MultitypeMonitoringEndpoint
	DonID                    uint32
	ChainID                  string

	// OCR3
	TraceLogging                 bool
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
	})
	ds := observation.NewDataSource(logger.Named(lggr, "DataSource"), cfg.Registry, t)

	return &delegate{services.StateMachine{}, cfg, reportCodecs, cfg.ShouldRetireCache, ds, t, []Closer{}}, nil
}

func (d *delegate) Start(ctx context.Context) error {
	return d.StartOnce("LLODelegate", func() error {
		// create the oracle from config values
		if !(len(d.cfg.ContractConfigTrackers) == 1 || len(d.cfg.ContractConfigTrackers) == 2) {
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

			oracle, err := ocr2plus.NewOracle(ocr2plus.OCR3OracleArgs[llotypes.ReportInfo]{
				BinaryNetworkEndpointFactory: d.cfg.BinaryNetworkEndpointFactory,
				V2Bootstrappers:              d.cfg.V2Bootstrappers,
				ContractConfigTracker:        configTracker,
				ContractTransmitter:          d.cfg.ContractTransmitter,
				Database:                     d.cfg.NewOCR3DB(int32(i)), // //nolint:gosec // G115 // impossible due to check on line 119
				LocalConfig:                  d.cfg.LocalConfig,
				Logger:                       ocrLogger,
				MonitoringEndpoint:           d.cfg.OCR3MonitoringEndpoint,
				OffchainConfigDigester:       d.cfg.OffchainConfigDigester,
				OffchainKeyring:              d.cfg.OffchainKeyring,
				OnchainKeyring:               d.cfg.OnchainKeyring,
				ReportingPluginFactory: promwrapper.NewReportingPluginFactory(
					datastreamsllo.NewPluginFactory(
						datastreamsllo.PluginFactoryParams{
							Config:                           d.cfg.ReportingPluginConfig,
							PredecessorRetirementReportCache: psrrc,
							ShouldRetireCache:                d.src,
							RetirementReportCodec:            d.cfg.RetirementReportCodec,
							ChannelDefinitionCache:           d.cfg.ChannelDefinitionCache,
							DataSource:                       d.ds,
							Logger:                           logger.Named(lggr, "ReportingPlugin"),
							OnchainConfigCodec:               llo.EVMOnchainConfigCodec{},
							ReportCodecs:                     d.reportCodecs,
							OutcomeTelemetryCh:               d.telem.GetOutcomeTelemetryCh(),
							ReportTelemetryCh:                d.telem.GetReportTelemetryCh(),
							DonID:                            d.cfg.DonID,
						},
					),
					lggr,
					d.cfg.ChainID,
					"llo",
				),
				MetricsRegisterer: prometheus.WrapRegistererWith(map[string]string{"job_name": d.cfg.JobName.ValueOrZero()}, prometheus.DefaultRegisterer),
			})

			if err != nil {
				return fmt.Errorf("%w: failed to create new OCR oracle", err)
			}

			d.oracles = append(d.oracles, oracle)

			merr = errors.Join(merr, oracle.Start())
		}

		return merr
	})
}

func (d *delegate) Close() error {
	return d.StopOnce("LLODelegate", func() (merr error) {
		for _, oracle := range d.oracles {
			merr = errors.Join(merr, oracle.Close())
		}
		merr = errors.Join(merr, d.telem.Close())
		return merr
	})
}
