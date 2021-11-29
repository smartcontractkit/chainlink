package offchainreporting

import (
	"fmt"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/offchain_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	ocr "github.com/smartcontractkit/libocr/offchainreporting"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

type Delegate struct {
	db                    *sqlx.DB
	jobORM                job.ORM
	keyStore              keystore.Master
	pipelineRunner        pipeline.Runner
	peerWrapper           *ocrcommon.SingletonPeerWrapper
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	chainSet              evm.ChainSet
	lggr                  logger.Logger
}

var _ job.Delegate = (*Delegate)(nil)

const ConfigOverriderPollInterval = 30 * time.Second

func NewDelegate(
	db *sqlx.DB,
	jobORM job.ORM,
	keyStore keystore.Master,
	pipelineRunner pipeline.Runner,
	peerWrapper *ocrcommon.SingletonPeerWrapper,
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator,
	chainSet evm.ChainSet,
	lggr logger.Logger,
) *Delegate {
	return &Delegate{
		db,
		jobORM,
		keyStore,
		pipelineRunner,
		peerWrapper,
		monitoringEndpointGen,
		chainSet,
		lggr.Named("OCR"),
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
	chain, err := d.chainSet.Get(jobSpec.OffchainreportingOracleSpec.EVMChainID.ToInt())
	if err != nil {
		return nil, err
	}
	concreteSpec, err := job.LoadEnvConfigVarsOCR(chain.Config(), d.keyStore.P2P(), *jobSpec.OffchainreportingOracleSpec)
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

	ocrDB := NewDB(d.db.DB, concreteSpec.ID, d.lggr)

	tracker := NewOCRContractTracker(
		contract,
		contractFilterer,
		contractCaller,
		chain.Client(),
		chain.LogBroadcaster(),
		jobSpec.ID,
		d.lggr,
		d.db,
		ocrDB,
		chain.Config(),
		chain.HeadBroadcaster(),
	)
	services = append(services, tracker)

	peerWrapper := d.peerWrapper
	if peerWrapper == nil {
		return nil, errors.New("cannot setup OCR job service, libp2p peer was missing")
	} else if !peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR jobs require a started and running peer. Did you forget to specify P2P_LISTEN_PORT?")
	} else if peerWrapper.PeerID.String() != concreteSpec.P2PPeerID.String() {
		return nil, errors.Errorf("given peer with ID '%s' does not match OCR configured peer with ID: %s", peerWrapper.PeerID.String(), concreteSpec.P2PPeerID.String())
	}
	var bootstrapPeers []string
	if concreteSpec.P2PBootstrapPeers != nil {
		bootstrapPeers = concreteSpec.P2PBootstrapPeers
	} else {
		bootstrapPeers, err = chain.Config().P2PBootstrapPeers()
		if err != nil {
			return nil, err
		}
	}
	// TODO: May want to follow up with spec override support for v2 bootstrappers?
	v2BootstrapPeers := chain.Config().P2PV2Bootstrappers()

	loggerWith := d.lggr.With(
		"contractAddress", concreteSpec.ContractAddress,
		"jobName", jobSpec.Name.ValueOrZero(),
		"jobID", jobSpec.ID,
	)
	ocrLogger := logger.NewOCRWrapper(loggerWith, chain.Config().OCRTraceLogging(), func(msg string) {
		d.jobORM.TryRecordError(jobSpec.ID, msg)
	})

	lc := NewLocalConfig(chain.Config(), *concreteSpec)
	if err = ocr.SanityCheckLocalConfig(lc); err != nil {
		return nil, err
	}
	loggerWith.Info(fmt.Sprintf("OCR job using local config %+v", lc))

	if concreteSpec.IsBootstrapPeer {
		var bootstrapper *ocr.BootstrapNode
		bootstrapper, err = ocr.NewBootstrapNode(ocr.BootstrapNodeArgs{
			BootstrapperFactory:   peerWrapper.Peer,
			V1Bootstrappers:       bootstrapPeers,
			V2Bootstrappers:       v2BootstrapPeers,
			ContractConfigTracker: tracker,
			Database:              ocrDB,
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

		ocrkey, err := d.keyStore.OCR().Get(concreteSpec.EncryptedOCRKeyBundleID.String())
		if err != nil {
			return nil, err
		}
		contractABI, err := abi.JSON(strings.NewReader(offchainaggregator.OffchainAggregatorABI))
		if err != nil {
			return nil, errors.Wrap(err, "could not get contract ABI JSON")
		}

		strategy := bulletprooftxmanager.NewQueueingTxStrategy(jobSpec.ExternalJobID, chain.Config().OCRDefaultTransactionQueueDepth(), chain.Config().OCRSimulateTransactions())

		contractTransmitter := NewOCRContractTransmitter(
			concreteSpec.ContractAddress.Address(),
			contractCaller,
			contractABI,
			ocrcommon.NewTransmitter(chain.TxManager(), concreteSpec.TransmitterAddress.Address(), chain.Config().EvmGasLimitDefault(), strategy),
			chain.LogBroadcaster(),
			tracker,
			chain.ID(),
		)

		runResults := make(chan pipeline.Run, chain.Config().JobPipelineResultWriteQueueDepth())
		jobSpec.PipelineSpec.JobName = jobSpec.Name.ValueOrZero()
		jobSpec.PipelineSpec.JobID = jobSpec.ID

		var configOverrider ocrtypes.ConfigOverrider
		configOverriderService, err := d.maybeCreateConfigOverrider(loggerWith, chain, concreteSpec.ContractAddress)
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
			Database: ocrDB,
			Datasource: ocrcommon.NewDataSourceV1(
				d.pipelineRunner,
				jobSpec,
				*jobSpec.PipelineSpec,
				loggerWith,
				runResults,
			),
			LocalConfig:                  lc,
			ContractTransmitter:          contractTransmitter,
			ContractConfigTracker:        tracker,
			PrivateKeys:                  ocrkey,
			BinaryNetworkEndpointFactory: peerWrapper.Peer,
			Logger:                       ocrLogger,
			V1Bootstrappers:              bootstrapPeers,
			V2Bootstrappers:              v2BootstrapPeers,
			MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(concreteSpec.ContractAddress.Address()),
			ConfigOverrider:              configOverrider,
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

func (d *Delegate) maybeCreateConfigOverrider(logger logger.Logger, chain evm.Chain, contractAddress ethkey.EIP55Address) (*ConfigOverriderImpl, error) {
	flagsContractAddress := chain.Config().FlagsContractAddress()
	if flagsContractAddress != "" {
		flags, err := NewFlags(flagsContractAddress, chain.Client())
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
