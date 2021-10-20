package offchainreporting

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/offchain_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	ocr "github.com/smartcontractkit/libocr/offchainreporting"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

type DelegateConfig interface {
	Chain() *chains.Chain
	ChainID() *big.Int
	Dev() bool
	EthGasLimitDefault() uint64
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
	FlagsContractAddress() string
}

type Delegate struct {
	db                 *gorm.DB
	txm                txManager
	jobORM             job.ORM
	config             DelegateConfig
	keyStore           *keystore.OCR
	pipelineRunner     pipeline.Runner
	ethClient          eth.Client
	logBroadcaster     log.Broadcaster
	peerWrapper        *SingletonPeerWrapper
	monitoringEndpoint ocrtypes.MonitoringEndpoint
	chain              *chains.Chain
	headBroadcaster    httypes.HeadBroadcaster
}

var _ job.Delegate = (*Delegate)(nil)

const ConfigOverriderPollInterval = 30 * time.Second

func NewDelegate(
	db *gorm.DB,
	txm txManager,
	jobORM job.ORM,
	config DelegateConfig,
	keyStore *keystore.OCR,
	pipelineRunner pipeline.Runner,
	ethClient eth.Client,
	logBroadcaster log.Broadcaster,
	peerWrapper *SingletonPeerWrapper,
	monitoringEndpoint ocrtypes.MonitoringEndpoint,
	chain *chains.Chain,
	headBroadcaster httypes.HeadBroadcaster,
) *Delegate {
	return &Delegate{
		db,
		txm,
		jobORM,
		config,
		keyStore,
		pipelineRunner,
		ethClient,
		logBroadcaster,
		peerWrapper,
		monitoringEndpoint,
		chain,
		headBroadcaster,
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

	contract, err := offchain_aggregator_wrapper.NewOffchainAggregator(concreteSpec.ContractAddress.Address(), d.ethClient)
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregator")
	}

	contractFilterer, err := offchainaggregator.NewOffchainAggregatorFilterer(concreteSpec.ContractAddress.Address(), d.ethClient)
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregatorFilterer")
	}

	contractCaller, err := offchainaggregator.NewOffchainAggregatorCaller(concreteSpec.ContractAddress.Address(), d.ethClient)
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
		d.ethClient,
		d.logBroadcaster,
		jobSpec.ID,
		*logger.Default,
		d.db,
		ocrdb,
		d.chain,
		d.headBroadcaster,
	)
	services = append(services, tracker)

	peerID, err := d.config.P2PPeerID(concreteSpec.P2PPeerID)
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
	bootstrapPeers, err := d.config.P2PBootstrapPeers(concreteSpec.P2PBootstrapPeers)
	if err != nil {
		return nil, err
	}
	v2BootstrapPeers := d.config.P2PV2Bootstrappers()

	loggerWith := logger.CreateLogger(logger.Default.With(
		"contractAddress", concreteSpec.ContractAddress,
		"jobName", jobSpec.Name.ValueOrZero(),
		"jobID", jobSpec.ID))
	ocrLogger := NewLogger(loggerWith, d.config.OCRTraceLogging(), func(msg string) {
		d.jobORM.RecordError(context.Background(), jobSpec.ID, msg)
	})

	lc := ocrtypes.LocalConfig{
		BlockchainTimeout:                      d.config.OCRBlockchainTimeout(time.Duration(concreteSpec.BlockchainTimeout)),
		ContractConfigConfirmations:            d.config.OCRContractConfirmations(concreteSpec.ContractConfigConfirmations),
		SkipContractConfigConfirmations:        d.config.Chain().IsL2(),
		ContractConfigTrackerPollInterval:      d.config.OCRContractPollInterval(time.Duration(concreteSpec.ContractConfigTrackerPollInterval)),
		ContractConfigTrackerSubscribeInterval: d.config.OCRContractSubscribeInterval(time.Duration(concreteSpec.ContractConfigTrackerSubscribeInterval)),
		ContractTransmitterTransmitTimeout:     d.config.OCRContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                        d.config.OCRDatabaseTimeout(),
		DataSourceTimeout:                      d.config.OCRObservationTimeout(time.Duration(concreteSpec.ObservationTimeout)),
		DataSourceGracePeriod:                  d.config.OCRObservationGracePeriod(),
	}
	if d.config.Dev() {
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
		kb, err := d.config.OCRKeyBundleID(concreteSpec.EncryptedOCRKeyBundleID)
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

		ta, err := d.config.OCRTransmitterAddress(concreteSpec.TransmitterAddress)
		if err != nil {
			return nil, err
		}

		strategy := bulletprooftxmanager.NewQueueingTxStrategy(jobSpec.ExternalJobID, d.config.OCRDefaultTransactionQueueDepth())

		contractTransmitter := NewOCRContractTransmitter(
			concreteSpec.ContractAddress.Address(),
			contractCaller,
			contractABI,
			NewTransmitter(d.txm, d.db, ta.Address(), d.config.EthGasLimitDefault(), strategy),
			d.logBroadcaster,
			tracker,
			d.config.ChainID(),
		)

		runResults := make(chan pipeline.RunWithResults, d.config.JobPipelineResultWriteQueueDepth())
		jobSpec.PipelineSpec.JobName = jobSpec.Name.ValueOrZero()
		jobSpec.PipelineSpec.JobID = jobSpec.ID

		var configOverrider ocrtypes.ConfigOverrider
		configOverriderService, err := d.maybeCreateConfigOverrider(loggerWith, concreteSpec.ContractAddress)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create ConfigOverrider")
		}

		// NOTE: conditional assigning to `configOverrider` is necessary due to the unfortunate fact that assigning `nil` to an
		// interface variable causes `x == nil` checks to always return false, so methods on the interface cannot be safely called then.
		//
		// the problematic case would be:
		//    configOverriderService, err := d.maybeCreateConfigOverrider(...)
		//	  if err != nil { return ... }
		//	  configOverrider = configOverriderService // contract might be `nil`
		//    assert.False(configOverrider != nil) // even if 'contract' was nil, this check will return true, unexpectedly
		if configOverriderService != nil {
			services = append(services, configOverriderService)
			configOverrider = configOverriderService
		}

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
			MonitoringEndpoint:           d.monitoringEndpoint,
			ConfigOverrider:              configOverrider,
		})
		if err != nil {
			return nil, errors.Wrap(err, "error calling NewOracle")
		}
		services = append(services, oracle)

		// RunResultSaver needs to be started first so its available
		// to read db writes. It is stopped last after the Oracle is shut down
		// so no further runs are enqueued and we can drain the queue.
		services = append([]job.Service{NewResultRunSaver(
			d.db,
			runResults,
			d.pipelineRunner,
			make(chan struct{}),
			*loggerWith,
		)}, services...)
	}

	return services, nil
}

func (d *Delegate) maybeCreateConfigOverrider(logger *logger.Logger, contractAddress ethkey.EIP55Address) (*ConfigOverriderImpl, error) {
	flagsContractAddress := d.config.FlagsContractAddress()
	if flagsContractAddress != "" {
		flags, err := NewFlags(flagsContractAddress, d.ethClient)
		if err != nil {
			return nil, errors.Wrapf(err,
				"OCR: unable to create Flags contract instance, check address: %s or remove FLAGS_CONTRACT_ADDRESS configuration variable",
				flagsContractAddress,
			)
		}

		ticker := utils.NewPausableTicker(ConfigOverriderPollInterval)
		return NewConfigOverriderImpl(logger, contractAddress, flags, &ticker)
	}
	return nil, nil
}
