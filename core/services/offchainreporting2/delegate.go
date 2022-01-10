package offchainreporting2

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/services/relay/types"

	"github.com/smartcontractkit/chainlink/core/services/keystore"

	"github.com/smartcontractkit/chainlink/core/config"

	"github.com/smartcontractkit/sqlx"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/telemetry"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type Config interface {
	config.OCR2Config
	Dev() bool
	JobPipelineResultWriteQueueDepth() uint64
}

type Delegate struct {
	db                    *sqlx.DB
	jobORM                job.ORM
	pipelineRunner        pipeline.Runner
	peerWrapper           *ocrcommon.SingletonPeerWrapper
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	chainSet              evm.ChainSet
	cfg                   Config
	lggr                  logger.Logger
	ks                    keystore.OCR2
	relayer               types.Relayer
}

var _ job.Delegate = (*Delegate)(nil)

func NewDelegate(
	db *sqlx.DB,
	jobORM job.ORM,
	pipelineRunner pipeline.Runner,
	peerWrapper *ocrcommon.SingletonPeerWrapper,
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator,
	chainSet evm.ChainSet,
	lggr logger.Logger,
	cfg Config,
	ks keystore.OCR2,
	relayer types.Relayer,
) *Delegate {
	return &Delegate{
		db,
		jobORM,
		pipelineRunner,
		peerWrapper,
		monitoringEndpointGen,
		chainSet,
		cfg,
		lggr,
		ks,
		relayer,
	}
}

func (d Delegate) JobType() job.Type {
	return job.OffchainReporting2
}

func (Delegate) OnJobCreated(spec job.Job) {}
func (Delegate) OnJobDeleted(spec job.Job) {}

func (Delegate) AfterJobCreated(spec job.Job)  {}
func (Delegate) BeforeJobDeleted(spec job.Job) {}

func (d Delegate) ServicesForSpec(jobSpec job.Job) (services []job.Service, err error) {
	spec := jobSpec.Offchainreporting2OracleSpec
	if spec == nil {
		return nil, errors.Errorf("offchainreporting.Delegate expects an *job.Offchainreporting2OracleSpec to be present, got %v", jobSpec)
	}

	ocr2Provider, err := d.relayer.NewOCR2Provider(jobSpec.ExternalJobID, spec)
	if err != nil {
		return nil, errors.Wrap(err, "error calling 'relayer.NewOCR2Provider'")
	}
	services = append(services, ocr2Provider)

	ocrDB := NewDB(d.db.DB, spec.ID, d.lggr)
	peerWrapper := d.peerWrapper
	if peerWrapper == nil {
		return nil, errors.New("cannot setup OCR2 job service, libp2p peer was missing")
	} else if !peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR2 jobs require a started and running peer. Did you forget to specify P2P_LISTEN_PORT?")
	}

	loggerWith := d.lggr.With(
		"OCRLogger", "true",
		"contractID", spec.ContractID,
		"jobName", jobSpec.Name.ValueOrZero(),
		"jobID", jobSpec.ID,
	)
	ocrLogger := ocrcommon.NewOCRWrapper(loggerWith, true, func(msg string) {
		d.lggr.ErrorIf(d.jobORM.RecordError(jobSpec.ID, msg), "unable to record error")
	})

	lc := toLocalConfig(d.cfg, *spec)
	if err = libocr2.SanityCheckLocalConfig(lc); err != nil {
		return nil, err
	}
	d.lggr.Infow("OCR2 job using local config",
		"BlockchainTimeout", lc.BlockchainTimeout,
		"ContractConfigConfirmations", lc.ContractConfigConfirmations,
		"ContractConfigTrackerPollInterval", lc.ContractConfigTrackerPollInterval,
		"ContractTransmitterTransmitTimeout", lc.ContractTransmitterTransmitTimeout,
		"DatabaseTimeout", lc.DatabaseTimeout,
	)

	tracker := ocr2Provider.ContractConfigTracker()
	offchainConfigDigester := ocr2Provider.OffchainConfigDigester()

	if spec.IsBootstrapPeer {
		bootstrapNodeArgs := libocr2.BootstrapperArgs{
			BootstrapperFactory:    peerWrapper.Peer2,
			ContractConfigTracker:  tracker,
			Database:               ocrDB,
			LocalConfig:            lc,
			Logger:                 ocrLogger,
			OffchainConfigDigester: offchainConfigDigester,
		}
		var bootstrapper *libocr2.Bootstrapper
		d.lggr.Debugw("Launching new bootstrap node", "args", bootstrapNodeArgs)
		bootstrapper, err = libocr2.NewBootstrapper(bootstrapNodeArgs)
		if err != nil {
			return nil, errors.Wrap(err, "error calling NewBootstrapNode")
		}
		services = append(services, bootstrapper)
	} else {
		bootstrapPeers, err := ocrcommon.GetValidatedBootstrapPeers(spec.P2PBootstrapPeers, peerWrapper.Config().P2PV2Bootstrappers())
		if err != nil {
			return nil, err
		}
		d.lggr.Debugw("Using bootstrap peers", "peers", bootstrapPeers)
		// Fetch the specified OCR2 key bundle
		var kbID string
		if spec.OCRKeyBundleID.Valid {
			kbID = spec.OCRKeyBundleID.String
		} else if kbID, err = d.cfg.OCR2KeyBundleID(); err != nil {
			return nil, err
		}
		kb, err := d.ks.Get(kbID)
		if err != nil {
			return nil, err
		}

		runResults := make(chan pipeline.Run, d.cfg.JobPipelineResultWriteQueueDepth())
		juelsPerFeeCoinPipelineSpec := pipeline.Spec{
			ID:           jobSpec.ID,
			DotDagSource: spec.JuelsPerFeeCoinPipeline,
			CreatedAt:    time.Now(),
		}
		numericalMedianFactory := median.NumericalMedianFactory{
			ContractTransmitter: ocr2Provider.MedianContract(),
			DataSource: ocrcommon.NewDataSourceV2(d.pipelineRunner,
				jobSpec,
				*jobSpec.PipelineSpec,
				loggerWith,
				runResults,
			),
			JuelsPerFeeCoinDataSource: ocrcommon.NewInMemoryDataSource(d.pipelineRunner, jobSpec, juelsPerFeeCoinPipelineSpec, loggerWith),
			ReportCodec:               ocr2Provider.ReportCodec(),
			Logger:                    ocrLogger,
		}

		jobSpec.PipelineSpec.JobName = jobSpec.Name.ValueOrZero()
		jobSpec.PipelineSpec.JobID = jobSpec.ID
		oracle, err := libocr2.NewOracle(libocr2.OracleArgs{
			BinaryNetworkEndpointFactory: peerWrapper.Peer2,
			V2Bootstrappers:              bootstrapPeers,
			ContractTransmitter:          ocr2Provider.ContractTransmitter(),
			ContractConfigTracker:        tracker,
			Database:                     ocrDB,
			LocalConfig:                  lc,
			Logger:                       ocrLogger,
			MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID),
			OffchainConfigDigester:       offchainConfigDigester,
			OffchainKeyring:              kb,
			OnchainKeyring:               kb,
			ReportingPluginFactory:       numericalMedianFactory,
		})
		if err != nil {
			return nil, errors.Wrap(err, "error calling NewOracle")
		}
		services = append(services, oracle)

		// RunResultSaver needs to be started first so its available
		// to read odb writes. It is stopped last after the Oracle is shut down
		// so no further runs are enqueued and we can drain the queue.
		services = append([]job.Service{ocrcommon.NewResultRunSaver(
			runResults,
			d.pipelineRunner,
			make(chan struct{}),
			loggerWith,
		)}, services...)
	}

	return services, nil
}

func toLocalConfig(config Config, spec job.OffchainReporting2OracleSpec) ocrtypes.LocalConfig {
	var (
		blockchainTimeout     = time.Duration(spec.BlockchainTimeout)
		ccConfirmations       = spec.ContractConfigConfirmations
		ccTrackerPollInterval = time.Duration(spec.ContractConfigTrackerPollInterval)
	)
	if blockchainTimeout == 0 {
		blockchainTimeout = config.OCR2BlockchainTimeout()
	}
	if ccConfirmations == 0 {
		ccConfirmations = config.OCR2ContractConfirmations()
	}
	if ccTrackerPollInterval == 0 {
		ccTrackerPollInterval = config.OCR2ContractPollInterval()
	}
	lc := ocrtypes.LocalConfig{
		BlockchainTimeout:                  blockchainTimeout,
		ContractConfigConfirmations:        ccConfirmations,
		ContractConfigTrackerPollInterval:  ccTrackerPollInterval,
		ContractTransmitterTransmitTimeout: config.OCR2ContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                    config.OCR2DatabaseTimeout(),
	}
	if config.Dev() {
		// Skips config validation so we can use any config parameters we want.
		// For example to lower contractConfigTrackerPollInterval to speed up tests.
		lc.DevelopmentMode = ocrtypes.EnableDangerousDevelopmentMode
	}
	return lc
}
