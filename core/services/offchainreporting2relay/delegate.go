package offchainreporting2relay

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/sqlx"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting2"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/relay"
	"github.com/smartcontractkit/chainlink/core/services/relay/solana"
	"github.com/smartcontractkit/chainlink/core/services/telemetry"
	ocr "github.com/smartcontractkit/libocr/offchainreporting2"
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
	relayers              relay.Relayers
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
	relayers relay.Relayers,
) *Delegate {
	return &Delegate{
		db,
		jobORM,
		pipelineRunner,
		peerWrapper,
		monitoringEndpointGen,
		chainSet,
		lggr,
		relayers,
	}
}

func (d Delegate) JobType() job.Type {
	return job.OffchainReporting2Relay
}

func (Delegate) OnJobCreated(spec job.Job) {}
func (Delegate) OnJobDeleted(spec job.Job) {}

func (Delegate) AfterJobCreated(spec job.Job)  {}
func (Delegate) BeforeJobDeleted(spec job.Job) {}

func (d Delegate) ServicesForSpec(jobSpec job.Job) (services []job.Service, err error) {
	if jobSpec.Offchainreporting2OracleSpec == nil {
		return nil, errors.Errorf("offchainreporting.Delegate expects an *job.Offchainreporting2OracleSpec to be present, got %v", jobSpec)
	}
	spec := jobSpec.Offchainreporting2OracleSpec

	chain, err := d.chainSet.Get(jobSpec.Offchainreporting2OracleSpec.EVMChainID.ToInt())
	if err != nil {
		return nil, err
	}

	ocrdb := offchainreporting2.NewDB(d.db.DB, spec.ID)

	var kbID string
	if spec.EncryptedOCRKeyBundleID.Valid {
		kbID = spec.EncryptedOCRKeyBundleID.String
	} else if kbID, err = chain.Config().OCRKeyBundleID(); err != nil {
		return nil, err
	}

	// TODO [relay]: make a relay choice depending on job spec
	relayer, ok := d.relayers["solana"]
	if !ok {
		return nil, fmt.Errorf("unknown relayer type: %s", "TODO [relay]: spec.relayerType")
	}

	ocr2Provider, err := relayer.NewOCR2Provider(solana.OCR2ProviderConfig{
		NodeURL:     "", // TODO [relay]: add validator url from job spec
		Address:     spec.ContractAddress.String(),
		JobID:       spec.ID,
		KeyBundleID: kbID,
	})

	if err != nil {
		return nil, errors.Wrap(err, "error calling 'relay.NewOCR2Provider'")
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
	bootstrapPeers := spec.P2PBootstrapPeers
	if bootstrapPeers == nil {
		bootstrapPeers, err = chain.Config().P2PBootstrapPeers()
		if err != nil {
			return nil, err
		}
	}
	v2BootstrapPeers := chain.Config().P2PV2Bootstrappers()
	d.lggr.Debugw("Using bootstrap peers", "v1", bootstrapPeers, "v2", v2BootstrapPeers)

	loggerWith := d.lggr.With(
		"OCRLogger", "true",
		"contractAddress", spec.ContractAddress,
		"jobName", jobSpec.Name.ValueOrZero(),
		"jobID", jobSpec.ID,
	)
	ocrLogger := logger.NewOCRWrapper(loggerWith, true, func(msg string) {
		d.jobORM.RecordError(jobSpec.ID, msg)
	})

	lc := computeLocalConfig(chain.Config(), *spec)

	if cerr := ocr.SanityCheckLocalConfig(lc); cerr != nil {
		return nil, cerr
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
		bootstrapNodeArgs := ocr.BootstrapperArgs{
			BootstrapperFactory:    peerWrapper.Peer2,
			ContractConfigTracker:  tracker,
			Database:               ocrdb,
			LocalConfig:            lc,
			Logger:                 ocrLogger,
			OffchainConfigDigester: offchainConfigDigester,
		}
		var bootstrapper *ocr.Bootstrapper
		d.lggr.Debugw("Launching new bootstrap node", "args", bootstrapNodeArgs)
		bootstrapper, err = ocr.NewBootstrapper(bootstrapNodeArgs)
		if err != nil {
			return nil, errors.Wrap(err, "error calling NewBootstrapNode")
		}
		services = append(services, bootstrapper)
	} else {
		if len(bootstrapPeers)+len(v2BootstrapPeers) < 1 {
			return nil, errors.New("need at least one bootstrap peer")
		}

		runResults := make(chan pipeline.Run, chain.Config().JobPipelineResultWriteQueueDepth())
		juelsPerFeeCoinPipelineSpec := pipeline.Spec{
			ID:           jobSpec.ID,
			DotDagSource: jobSpec.Offchainreporting2OracleSpec.JuelsPerFeeCoinPipeline,
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
		oracle, err := ocr.NewOracle(ocr.OracleArgs{
			BinaryNetworkEndpointFactory: peerWrapper.Peer2,
			V2Bootstrappers:              v2BootstrapPeers,
			ContractTransmitter:          ocr2Provider.ContractTransmitter(),
			ContractConfigTracker:        tracker,
			Database:                     ocrdb,
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
		// to read db writes. It is stopped last after the Oracle is shut down
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

func computeLocalConfig(config offchainreporting2.ValidationConfig, spec job.OffchainReporting2OracleSpec) ocrtypes.LocalConfig {
	var blockchainTimeout time.Duration
	if spec.BlockchainTimeout != 0 {
		blockchainTimeout = time.Duration(spec.BlockchainTimeout)
	} else {
		blockchainTimeout = config.OCRBlockchainTimeout()
	}

	var contractConfirmations uint16
	if spec.ContractConfigConfirmations != 0 {
		contractConfirmations = spec.ContractConfigConfirmations
	} else {
		contractConfirmations = config.OCRContractConfirmations()
	}

	var contractConfigTrackerPollInterval time.Duration
	if spec.ContractConfigTrackerPollInterval != 0 {
		contractConfigTrackerPollInterval = time.Duration(spec.ContractConfigTrackerPollInterval)
	} else {
		contractConfigTrackerPollInterval = config.OCRContractPollInterval()
	}

	lc := ocrtypes.LocalConfig{
		BlockchainTimeout:                  blockchainTimeout,
		ContractConfigConfirmations:        contractConfirmations,
		ContractConfigTrackerPollInterval:  contractConfigTrackerPollInterval,
		ContractTransmitterTransmitTimeout: config.OCRContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                    config.OCRDatabaseTimeout(),
	}
	if config.Dev() {
		// Skips config validation so we can use any config parameters we want.
		// For example to lower contractConfigTrackerPollInterval to speed up tests.
		lc.DevelopmentMode = ocrtypes.EnableDangerousDevelopmentMode
	}
	return lc
}
