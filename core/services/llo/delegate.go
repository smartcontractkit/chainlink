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
	lloconfig "github.com/smartcontractkit/chainlink-data-streams/llo/pluginconfig"
	lloprotocol "github.com/smartcontractkit/chainlink-data-streams/llo/protocol"
	"github.com/smartcontractkit/chainlink-data-streams/llo/retirement"
	"github.com/smartcontractkit/chainlink-data-streams/llo/transmitter"
	llov30 "github.com/smartcontractkit/chainlink-data-streams/llo/v30"
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

	// PluginVersions selects the plugin per protocol instance, positionally
	// aligned with ContractConfigTrackers (see chainlink-data-streams
	// llo/pluginconfig.PluginConfig.PluginVersions). The two entries differ only
	// during blue/green handover across different plugin versions , where one
	// instance hands over to the other in a different version.
	PluginVersions []lloconfig.PluginVersion
	// V31Config carries the v31 plugin knobs from the job's plugin config. Only
	// read by v31 instances; zero fields fall through to the plugin defaults.
	V31Config lloconfig.V31Config
	// BinaryNetworkEndpoint2Factory is the OCR3.1 ("2") network endpoint factory
	// (peerWrapper.Peer3_1). Required when any instance is v31.
	BinaryNetworkEndpoint2Factory ocr2types.BinaryNetworkEndpoint2Factory
	// KeyValueDatabaseFactory provides the replicated per-configDigest key-value
	// store the OCR3.1 protocol requires. One factory serves both instances: it
	// keys the database by config digest, so they get separate keyspaces.
	// Required when any instance is v31.
	KeyValueDatabaseFactory ocr3_1types.KeyValueDatabaseFactory
}

// anyV31 reports whether any protocol instance runs the v31 plugin. The
// OCR3.1-only dependencies are per job, so one v31 instance requires them.
func (cfg DelegateConfig) anyV31() bool {
	for _, v := range cfg.PluginVersions {
		if v == lloconfig.PluginVersionV31 {
			return true
		}
	}
	return false
}

// validateInstances checks the per-instance plugin selection and the
// dependencies it implies.
func (cfg DelegateConfig) validateInstances() error {
	if len(cfg.PluginVersions) != len(cfg.ContractConfigTrackers) {
		return fmt.Errorf("expected one PluginVersions entry per ContractConfigTracker, got %d entries for %d trackers", len(cfg.PluginVersions), len(cfg.ContractConfigTrackers))
	}
	for i, v := range cfg.PluginVersions {
		switch v {
		case lloconfig.PluginVersionV30, lloconfig.PluginVersionV31:
		default:
			return fmt.Errorf("unsupported plugin version for instance %d: %q", i, v)
		}
	}
	if cfg.anyV31() {
		if cfg.KeyValueDatabaseFactory == nil {
			return errors.New("KeyValueDatabaseFactory must not be nil when running OCR3.1")
		}
		if cfg.BinaryNetworkEndpoint2Factory == nil {
			return errors.New("BinaryNetworkEndpoint2Factory must not be nil when running OCR3.1")
		}
	}
	return nil
}

func NewDelegate(cfg DelegateConfig) (job.ServiceCtx, error) {
	lggr := logger.Sugared(cfg.Logger).With("jobName", cfg.JobName.ValueOrZero(), "donID", cfg.DonID)
	if cfg.DataSource == nil {
		return nil, errors.New("DataSource must not be nil")
	}
	if cfg.Runner == nil {
		return nil, errors.New("runner must not be nil")
	}
	if cfg.Registry == nil {
		return nil, errors.New("registry must not be nil")
	}
	if cfg.RetirementReportCache == nil {
		return nil, errors.New("RetirementReportCache must not be nil")
	}
	if cfg.ShouldRetireCache == nil {
		return nil, errors.New("ShouldRetireCache must not be nil")
	}
	if err := cfg.validateInstances(); err != nil {
		return nil, err
	}
	var codecLggr logger.Logger
	if cfg.ReportingPluginConfig.VerboseLogging {
		codecLggr = logger.Named(lggr, "ReportCodecs")
	} else {
		codecLggr = logger.Nop()
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

			// NewDelegate rejected any version this switch does not handle, so
			// a new one added upstream fails at startup rather than silently
			// running v30.
			var oracle ocr2plus.Oracle
			var err error
			switch version := d.cfg.PluginVersions[i]; version {
			case lloconfig.PluginVersionV31:
				oracle, err = d.newOracleV31(i, configTracker, lggr, ocrLogger, psrrc)
			default:
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

// v31FactoryParams assembles the v31 plugin factory params, mapping the job's
// V31Config knobs onto it. Knobs left at zero are forwarded as zero, which the
// factory reads as "apply the plugin default".
func (d *delegate) v31FactoryParams(lggr logger.Logger, psrrc lloprotocol.PredecessorRetirementReportCache) llov31.PluginFactoryParams {
	return llov31.PluginFactoryParams{
		VerboseLogging:                   d.cfg.ReportingPluginConfig.VerboseLogging || d.cfg.V31Config.VerboseLogging,
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
		MaxSnapshotRounds:                d.cfg.V31Config.MaxSnapshotRounds,
		BlobLifetimeRounds:               d.cfg.V31Config.BlobLifetimeRounds,
		MaxDurationBlobObservation:       d.cfg.V31Config.MaxDurationBlobObservation.Duration(),
		BlobInFlightWaitFactor:           d.cfg.V31Config.BlobInFlightWaitFactor,
		MaxBlobSnapshotAge:               d.cfg.V31Config.MaxBlobSnapshotAge.Duration(),
		MaxRoundPeriod:                   d.cfg.V31Config.MaxRoundPeriod.Duration(),
	}
}

// newOracleV31 builds an OCR3.1 oracle running the llo/v31 reporting plugin. It
// differs from v30 by the OCR3.1 oracle args (OCR3_1OracleArgs2), the "2"
// network endpoint factory, and the required replicated KeyValueDatabaseFactory.
func (d *delegate) newOracleV31(i int, configTracker ocr2types.ContractConfigTracker, lggr logger.Logger, ocrLogger ocrcommontypes.Logger, psrrc lloprotocol.PredecessorRetirementReportCache) (ocr2plus.Oracle, error) {
	factory := promwrapper31.NewReportingPluginFactory(
		llov31.NewPluginFactory(d.v31FactoryParams(lggr, psrrc)),
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
