package offchainreporting2

import (
	"time"

	"github.com/smartcontractkit/sqlx"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/relay"
	"github.com/smartcontractkit/chainlink/core/services/telemetry"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type Delegate struct {
	db                    *sqlx.DB
	jobORM                job.ORM
	pipelineRunner        pipeline.Runner
	peerWrapper           *ocrcommon.SingletonPeerWrapper
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	chainSet              evm.ChainSet
	lggr                  logger.Logger
	relayer               relay.Relayer
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
	relayer relay.Relayer,
) *Delegate {
	return &Delegate{
		db,
		jobORM,
		pipelineRunner,
		peerWrapper,
		monitoringEndpointGen,
		chainSet,
		lggr,
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

	chain, err := d.chainSet.Get(spec.EVMChainID.ToInt())
	if err != nil {
		return nil, err
	}

	ocrDB := NewDB(d.db.DB, spec.ID, d.lggr)

	ocr2Provider, err := d.relayer.NewOCR2Provider(jobSpec.ExternalJobID, spec)
	if err != nil {
		return nil, errors.Wrap(err, "error calling 'relayer.NewOCR2Provider'")
	}
	services = append(services, ocr2Provider)

	var peerID p2pkey.PeerID
	if spec.P2PPeerID != nil {
		peerID = *spec.P2PPeerID
	} else {
		peerID = chain.Config().P2PPeerID()
	}
	peerWrapper := d.peerWrapper
	if peerWrapper == nil {
		return nil, errors.New("cannot setup OCR2 job service, libp2p peer was missing")
	} else if !peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR2 jobs require a started and running peer. Did you forget to specify P2P_LISTEN_PORT?")
	} else if peerWrapper.PeerID != peerID {
		return nil, errors.Errorf("given peer with ID '%s' does not match OCR2 configured peer with ID: %s", peerWrapper.PeerID.String(), peerID.String())
	}
	bootstrapPeers, err := ocrcommon.GetValidatedBootstrapPeers(spec.P2PBootstrapPeers, chain)
	if err != nil {
		return nil, err
	}
	d.lggr.Debugw("Using bootstrap peers", "peers", bootstrapPeers)

	loggerWith := d.lggr.With(
		"OCRLogger", "true",
		"contractAddress", spec.ContractAddress,
		"jobName", jobSpec.Name.ValueOrZero(),
		"jobID", jobSpec.ID,
	)
	ocrLogger := logger.NewOCRWrapper(loggerWith, true, func(msg string) {
		d.lggr.ErrorIf(d.jobORM.RecordError(jobSpec.ID, msg), "unable to record error")
	})

	lc := toLocalConfig(chain.Config(), *spec)
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
		if len(bootstrapPeers) < 1 {
			return nil, errors.New("need at least one bootstrap peer")
		}

		runResults := make(chan pipeline.Run, chain.Config().JobPipelineResultWriteQueueDepth())
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
			MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractAddress.Address()),
			OffchainConfigDigester:       offchainConfigDigester,
			OffchainKeyring:              ocr2Provider.OffchainKeyring(),
			OnchainKeyring:               ocr2Provider.OnchainKeyring(),
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

func toLocalConfig(config ValidationConfig, spec job.OffchainReporting2OracleSpec) ocrtypes.LocalConfig {
	var blockchainTimeout time.Duration
	if spec.BlockchainTimeout != 0 {
		blockchainTimeout = time.Duration(spec.BlockchainTimeout)
	} else {
		blockchainTimeout = config.OCR2BlockchainTimeout()
	}

	var ccConfirmations uint16
	if spec.ContractConfigConfirmations != 0 {
		ccConfirmations = spec.ContractConfigConfirmations
	} else {
		ccConfirmations = config.OCR2ContractConfirmations()
	}

	var ccTrackerPollInterval time.Duration
	if spec.ContractConfigTrackerPollInterval != 0 {
		ccTrackerPollInterval = time.Duration(spec.ContractConfigTrackerPollInterval)
	} else {
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
