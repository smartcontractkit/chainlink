package offchainreporting

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/offchain_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	ocr "github.com/smartcontractkit/libocr/offchainreporting"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

type Config interface {
	DefaultChainID() *big.Int
	Dev() bool
	EvmGasLimitDefault() uint64
	JobPipelineResultWriteQueueDepth() uint64
	OCRBlockchainTimeout(time.Duration) time.Duration
	OCRContractConfirmations(uint16) uint16
	OCRContractPollInterval(time.Duration) time.Duration
	OCRContractSubscribeInterval(time.Duration) time.Duration
	OCRContractTransmitterTransmitTimeout() time.Duration
	OCRDatabaseTimeout() time.Duration
	OCRDefaultTransactionQueueDepth() uint32
	OCRKeyBundleID(*models.Sha256Hash) (models.Sha256Hash, error)
	OCRObservationGracePeriod() time.Duration
	OCRObservationTimeout(time.Duration) time.Duration
	OCRTraceLogging() bool
	OCRTransmitterAddress(*ethkey.EIP55Address) (ethkey.EIP55Address, error)
	P2PBootstrapPeers([]string) ([]string, error)
	P2PPeerID(*p2pkey.PeerID) (p2pkey.PeerID, error)
	P2PV2Bootstrappers() []ocrtypes.BootstrapperLocator
}

type Delegate struct {
	db                    *gorm.DB
	jobORM                job.ORM
	keyStore              *keystore.OCR
	pipelineRunner        pipeline.Runner
	peerWrapper           *SingletonPeerWrapper
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	chainSet              evm.ChainSet
}

var _ job.Delegate = (*Delegate)(nil)

func NewDelegate(
	db *gorm.DB,
	jobORM job.ORM,
	keyStore *keystore.OCR,
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
	return job.OffchainReporting
}

func (Delegate) AfterJobCreated(spec job.Job)  {}
func (Delegate) BeforeJobDeleted(spec job.Job) {}

func (d Delegate) ServicesForSpec(jobSpec job.Job) (services []job.Service, err error) {
	if jobSpec.OffchainreportingOracleSpec == nil {
		return nil, errors.Errorf("offchainreporting.Delegate expects an *job.OffchainreportingOracleSpec to be present, got %v", jobSpec)
	}
	concreteSpec := jobSpec.OffchainreportingOracleSpec
	chain, err := d.chainSet.Get(jobSpec.OffchainreportingOracleSpec.EVMChainID.ToInt())
	if err != nil {
		return nil, err
	}

	contract, err := offchain_aggregator_wrapper.NewOffchainAggregator(concreteSpec.ContractAddress.Address(), chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregator")
	}

	contractFilterer, err := offchainaggregator.NewOffchainAggregatorFilterer(concreteSpec.ContractAddress.Address(), chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregatorFilterer")
	}

	contractCaller, err := offchainaggregator.NewOffchainAggregatorCaller(concreteSpec.ContractAddress.Address(), chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregatorCaller")
	}

	gormdb, errdb := d.db.DB()
	if errdb != nil {
		return nil, errors.Wrap(errdb, "unable to open sql db")
	}
	ocrdb := NewDB(gormdb, concreteSpec.ID)

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

	peerID, err := chain.Config().P2PPeerID(concreteSpec.P2PPeerID)
	if err != nil {
		return nil, err
	}
	peerWrapper := d.peerWrapper
	if peerWrapper == nil {
		return nil, errors.New("cannot setup OCR job service, libp2p peer was missing")
	} else if !peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR jobs require a started and running peer. Did you forget to specify P2P_LISTEN_PORT?")
	} else if peerWrapper.PeerID != peerID {
		return nil, errors.Errorf("given peer with ID '%s' does not match OCR configured peer with ID: %s", peerWrapper.PeerID.String(), peerID.String())
	}
	bootstrapPeers, err := chain.Config().P2PBootstrapPeers(concreteSpec.P2PBootstrapPeers)
	if err != nil {
		return nil, err
	}
	v2BootstrapPeers := chain.Config().P2PV2Bootstrappers()

	loggerWith := logger.Default.With(
		"contractAddress", concreteSpec.ContractAddress,
		"jobName", jobSpec.Name.ValueOrZero(),
		"jobID", jobSpec.ID,
	)
	ocrLogger := NewLogger(loggerWith, chain.Config().OCRTraceLogging(), func(msg string) {
		d.jobORM.RecordError(context.Background(), jobSpec.ID, msg)
	})

	lc := ocrtypes.LocalConfig{
		BlockchainTimeout:                      chain.Config().OCRBlockchainTimeout(time.Duration(concreteSpec.BlockchainTimeout)),
		ContractConfigConfirmations:            chain.Config().OCRContractConfirmations(concreteSpec.ContractConfigConfirmations),
		SkipContractConfigConfirmations:        chain.IsL2(),
		ContractConfigTrackerPollInterval:      chain.Config().OCRContractPollInterval(time.Duration(concreteSpec.ContractConfigTrackerPollInterval)),
		ContractConfigTrackerSubscribeInterval: chain.Config().OCRContractSubscribeInterval(time.Duration(concreteSpec.ContractConfigTrackerSubscribeInterval)),
		ContractTransmitterTransmitTimeout:     chain.Config().OCRContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                        chain.Config().OCRDatabaseTimeout(),
		DataSourceTimeout:                      chain.Config().OCRObservationTimeout(time.Duration(concreteSpec.ObservationTimeout)),
		DataSourceGracePeriod:                  chain.Config().OCRObservationGracePeriod(),
	}
	if chain.Config().Dev() {
		// Skips config validation so we can use any config parameters we want.
		// For example to lower contractConfigTrackerPollInterval to speed up tests.
		lc.DevelopmentMode = ocrtypes.EnableDangerousDevelopmentMode
	}
	if err := ocr.SanityCheckLocalConfig(lc); err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("OCR job using local config %+v", lc))

	if concreteSpec.IsBootstrapPeer {
		bootstrapper, err := ocr.NewBootstrapNode(ocr.BootstrapNodeArgs{
			BootstrapperFactory:   peerWrapper.Peer,
			V1Bootstrappers:       bootstrapPeers,
			ContractConfigTracker: tracker,
			Database:              ocrdb,
			LocalConfig:           lc,
			Logger:                ocrLogger,
		})
		if err != nil {
			return nil, errors.Wrap(err, "error calling NewBootstrapNode")
		}
		services = append(services, bootstrapper)
	} else {
		if len(bootstrapPeers) < 1 {
			return nil, errors.New("need at least one bootstrap peer")
		}
		kb, err := chain.Config().OCRKeyBundleID(concreteSpec.EncryptedOCRKeyBundleID)
		if err != nil {
			return nil, err
		}
		ocrkey, exists := d.keyStore.DecryptedOCRKey(kb)
		if !exists {
			return nil, errors.Errorf("OCR key '%v' does not exist", concreteSpec.EncryptedOCRKeyBundleID)
		}
		contractABI, err := abi.JSON(strings.NewReader(offchainaggregator.OffchainAggregatorABI))
		if err != nil {
			return nil, errors.Wrap(err, "could not get contract ABI JSON")
		}

		ta, err := chain.Config().OCRTransmitterAddress(concreteSpec.TransmitterAddress)
		if err != nil {
			return nil, err
		}

		strategy := bulletprooftxmanager.NewQueueingTxStrategy(jobSpec.ExternalJobID, chain.Config().OCRDefaultTransactionQueueDepth())

		contractTransmitter := NewOCRContractTransmitter(
			concreteSpec.ContractAddress.Address(),
			contractCaller,
			contractABI,
			NewTransmitter(chain.TxManager(), d.db, ta.Address(), chain.Config().EvmGasLimitDefault(), strategy),
			chain.LogBroadcaster(),
			tracker,
			chain.ID(),
		)

		runResults := make(chan pipeline.Run, chain.Config().JobPipelineResultWriteQueueDepth())
		jobSpec.PipelineSpec.JobName = jobSpec.Name.ValueOrZero()
		jobSpec.PipelineSpec.JobID = jobSpec.ID
		oracle, err := ocr.NewOracle(ocr.OracleArgs{
			Database: ocrdb,
			Datasource: &dataSource{
				pipelineRunner: d.pipelineRunner,
				ocrLogger:      *loggerWith,
				jobSpec:        jobSpec,
				spec:           *jobSpec.PipelineSpec,
				runResults:     runResults,
			},
			LocalConfig:                  lc,
			ContractTransmitter:          contractTransmitter,
			ContractConfigTracker:        tracker,
			PrivateKeys:                  &ocrkey,
			BinaryNetworkEndpointFactory: peerWrapper.Peer,
			Logger:                       ocrLogger,
			V1Bootstrappers:              bootstrapPeers,
			V2Bootstrappers:              v2BootstrapPeers,
			MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(concreteSpec.ContractAddress.Address()),
		})
		if err != nil {
			return nil, errors.Wrap(err, "error calling NewOracle")
		}
		services = append(services, oracle)

		// RunResultSaver needs to be started first so its available
		// to read db writes. It is stopped last after the Oracle is shut down
		// so no further runs are enqueued and we can drain the queue.
		services = append([]job.Service{NewResultRunSaver(
			postgres.UnwrapGormDB(d.db),
			runResults,
			d.pipelineRunner,
			make(chan struct{}),
			*loggerWith,
		)}, services...)
	}

	return services, nil
}
