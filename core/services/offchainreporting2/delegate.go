package offchainreporting2

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	offchain_aggregator_wrapper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers2/generated/offchainaggregator"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/telemetry"
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocr "github.com/smartcontractkit/libocr/offchainreporting2"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type DelegateConfig interface {
	Chain() ocrcommon.Chain
	ChainID() *big.Int
	Dev() bool
	EvmGasLimitDefault() uint64
	JobPipelineResultWriteQueueDepth() uint64
	OCR2BlockchainTimeout() time.Duration
	OCR2ContractConfirmations() uint16
	OCR2ContractPollInterval() time.Duration
	OCR2ContractTransmitterTransmitTimeout() time.Duration
	OCR2DatabaseTimeout() time.Duration
	OCR2DefaultTransactionQueueDepth() uint32
	OCR2KeyBundleID() (string, error)
	OCR2ObservationTimeout() time.Duration
	OCR2P2PBootstrapPeers() ([]string, error)
	OCR2P2PPeerID() (p2pkey.PeerID, error)
	OCR2P2PV2Bootstrappers() []ocrcommontypes.BootstrapperLocator
	OCR2TraceLogging() bool
	OCR2TransmitterAddress() (ethkey.EIP55Address, error)
}

type Delegate struct {
	db                    *gorm.DB
	jobORM                job.ORM
	keyStore              keystore.OCR2
	pipelineRunner        pipeline.Runner
	peerWrapper           *SingletonPeerWrapper
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	chainSet              evm.ChainSet
}

var _ job.Delegate = (*Delegate)(nil)

func NewDelegate(
	db *gorm.DB,
	jobORM job.ORM,
	keyStore keystore.OCR2,
	pipelineRunner pipeline.Runner,
	peerWrapper *SingletonPeerWrapper,
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator,
	chainSet evm.ChainSet,
) *Delegate {
	return &Delegate{
		db,
		jobORM,
		keyStore,
		pipelineRunner,
		peerWrapper,
		monitoringEndpointGen,
		chainSet,
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
	if jobSpec.Offchainreporting2OracleSpec == nil {
		return nil, errors.Errorf("offchainreporting.Delegate expects an *job.Offchainreporting2OracleSpec to be present, got %v", jobSpec)
	}
	spec := jobSpec.Offchainreporting2OracleSpec

	chain, err := d.chainSet.Get(jobSpec.OffchainreportingOracleSpec.EVMChainID.ToInt())
	if err != nil {
		return nil, err
	}

	contract, err := offchain_aggregator_wrapper.NewOffchainAggregator(spec.ContractAddress.Address(), chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregator")
	}

	contractFilterer, err := ocr2aggregator.NewOCR2AggregatorFilterer(spec.ContractAddress.Address(), chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregatorFilterer")
	}

	contractCaller, err := ocr2aggregator.NewOCR2AggregatorCaller(spec.ContractAddress.Address(), chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregatorCaller")
	}

	gormdb, errdb := d.db.DB()
	if errdb != nil {
		return nil, errors.Wrap(errdb, "unable to open sql db")
	}
	ocrdb := NewDB(gormdb, spec.ID)

	tracker := NewOCRContractTracker(
		contract,
		contractFilterer,
		contractCaller,
		chain.Client(),
		chain.LogBroadcaster(),
		jobSpec.ID,
		*logger.Default,
		d.db,
		ocrdb,
		chain,
		chain.HeadBroadcaster(),
	)
	services = append(services, tracker)

	var peerID p2pkey.PeerID
	if spec.P2PPeerID != nil {
		peerID = *spec.P2PPeerID
	} else if cPeerID, err := chain.Config().OCR2P2PPeerID(); err == nil {
		peerID = cPeerID
	} else {
		return nil, err
	}
	peerWrapper := d.peerWrapper
	if peerWrapper == nil {
		return nil, errors.New("cannot setup OCR2 job service, libp2p peer was missing")
	} else if !peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR2 jobs require a started and running peer. Did you forget to specify OCR2_P2P_LISTEN_PORT?")
	} else if peerWrapper.PeerID != peerID {
		return nil, errors.Errorf("given peer with ID '%s' does not match OCR2 configured peer with ID: %s", peerWrapper.PeerID.String(), peerID.String())
	}
	bootstrapPeers := spec.P2PBootstrapPeers
	if bootstrapPeers == nil {
		bootstrapPeers, err = chain.Config().OCR2P2PBootstrapPeers()
		if err != nil {
			return nil, err
		}
	}
	v2BootstrapPeers := chain.Config().OCR2P2PV2Bootstrappers()
	logger.Debugw("Using bootstrap peers", "v1", bootstrapPeers, "v2", v2BootstrapPeers)

	loggerWith := logger.Default.With(
		"contractAddress", spec.ContractAddress,
		"jobName", jobSpec.Name.ValueOrZero(),
		"jobID", jobSpec.ID,
	)
	ocrLogger := ocrcommon.NewLogger(loggerWith, chain.Config().OCR2TraceLogging(), func(msg string) {
		d.jobORM.RecordError(context.Background(), jobSpec.ID, msg)
	})

	lc := computeLocalConfig(chain.Config(), *spec)

	if err := ocr.SanityCheckLocalConfig(lc); err != nil {
		return nil, err
	}
	logger.Infow("OCR2 job using local config",
		"BlockchainTimeout", lc.BlockchainTimeout,
		"ContractConfigConfirmations", lc.ContractConfigConfirmations,
		"ContractConfigTrackerPollInterval", lc.ContractConfigTrackerPollInterval,
		"ContractTransmitterTransmitTimeout", lc.ContractTransmitterTransmitTimeout,
		"DatabaseTimeout", lc.DatabaseTimeout,
	)

	offchainConfigDigester := evmutil.EVMOffchainConfigDigester{
		ChainID:         chain.Config().ChainID().Uint64(),
		ContractAddress: spec.ContractAddress.Address(),
	}

	if spec.IsBootstrapPeer {
		bootstrapper, err := ocr.NewBootstrapNode(ocr.BootstrapNodeArgs{
			BootstrapperFactory:    peerWrapper.Peer,
			ContractConfigTracker:  tracker,
			Database:               ocrdb,
			LocalConfig:            lc,
			Logger:                 ocrLogger,
			OffchainConfigDigester: offchainConfigDigester,
		})
		if err != nil {
			return nil, errors.Wrap(err, "error calling NewBootstrapNode")
		}
		services = append(services, bootstrapper)
	} else {

		if len(bootstrapPeers)+len(v2BootstrapPeers) < 1 {
			return nil, errors.New("need at least one bootstrap peer")
		}

		var kb string
		if spec.EncryptedOCRKeyBundleID.Valid {
			kb = spec.EncryptedOCRKeyBundleID.String
		} else if kb, err = chain.Config().OCR2KeyBundleID(); err != nil {
			return nil, err
		}

		ocrkey, err := d.keyStore.Get(kb)
		if err != nil {
			return nil, err
		}
		contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
		if err != nil {
			return nil, errors.Wrap(err, "could not get contract ABI JSON")
		}

		var ta ethkey.EIP55Address
		if spec.TransmitterAddress != nil {
			ta = *spec.TransmitterAddress
		} else if ta, err = chain.Config().OCR2TransmitterAddress(); err != nil {
			return nil, err
		}

		strategy := bulletprooftxmanager.NewQueueingTxStrategy(jobSpec.ExternalJobID, chain.Config().OCR2DefaultTransactionQueueDepth())

		contractTransmitter := NewOCRContractTransmitter(
			contract.Address(),
			contractCaller,
			contractABI,
			ocrcommon.NewTransmitter(chain.TxManager(), d.db, ta.Address(), chain.Config().EvmGasLimitDefault(), strategy),
			chain.LogBroadcaster(),
			tracker,
		)

		runResults := make(chan pipeline.Run, chain.Config().JobPipelineResultWriteQueueDepth())
		numericalMedianFactory := median.NumericalMedianFactory{
			ContractTransmitter: contractTransmitter,
			DataSource: &dataSource{
				pipelineRunner: d.pipelineRunner,
				ocrLogger:      *loggerWith,
				jobSpec:        jobSpec,
				spec:           *jobSpec.PipelineSpec,
				runResults:     runResults,
			},
			Logger: ocrLogger,
		}

		jobSpec.PipelineSpec.JobName = jobSpec.Name.ValueOrZero()
		jobSpec.PipelineSpec.JobID = jobSpec.ID
		oracle, err := ocr.NewOracle(ocr.OracleArgs{
			BinaryNetworkEndpointFactory: peerWrapper.Peer,
			V2Bootstrappers:              v2BootstrapPeers,
			ContractTransmitter:          contractTransmitter,
			ContractConfigTracker:        tracker,
			Database:                     ocrdb,
			LocalConfig:                  lc,
			Logger:                       ocrLogger,
			MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractAddress.Address()),
			OffchainConfigDigester:       offchainConfigDigester,
			OffchainKeyring:              &ocrkey.OffchainKeyring,
			OnchainKeyring:               &ocrkey.OnchainKeyring,
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
			postgres.UnwrapGormDB(d.db),
			runResults,
			d.pipelineRunner,
			make(chan struct{}),
			*loggerWith,
		)}, services...)
	}

	return services, nil
}

func computeLocalConfig(config ValidationConfig, spec job.OffchainReporting2OracleSpec) ocrtypes.LocalConfig {
	var blockchainTimeout time.Duration
	if spec.BlockchainTimeout != 0 {
		blockchainTimeout = time.Duration(spec.BlockchainTimeout)
	} else {
		blockchainTimeout = config.OCR2BlockchainTimeout()
	}

	var contractConfirmations uint16
	if spec.ContractConfigConfirmations != 0 {
		contractConfirmations = spec.ContractConfigConfirmations
	} else {
		contractConfirmations = config.OCR2ContractConfirmations()
	}

	var contractConfigTrackerPollInterval time.Duration
	if spec.ContractConfigTrackerPollInterval != 0 {
		contractConfigTrackerPollInterval = time.Duration(spec.ContractConfigTrackerPollInterval)
	} else {
		contractConfigTrackerPollInterval = config.OCR2ContractPollInterval()
	}

	lc := ocrtypes.LocalConfig{
		BlockchainTimeout:                  blockchainTimeout,
		ContractConfigConfirmations:        contractConfirmations,
		ContractConfigTrackerPollInterval:  contractConfigTrackerPollInterval,
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
