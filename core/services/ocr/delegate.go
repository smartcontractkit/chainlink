package ocr

import (
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	ocr "github.com/smartcontractkit/libocr/offchainreporting"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/offchain_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/core/utils"
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
	cfg                   Config
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
	cfg Config,
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
		cfg,
	}
}

func (d Delegate) JobType() job.Type {
	return job.OffchainReporting
}

func (Delegate) AfterJobCreated(spec job.Job)  {}
func (Delegate) BeforeJobDeleted(spec job.Job) {}

// ServicesForSpec returns the OCR services that need to run for this job
func (d Delegate) ServicesForSpec(jb job.Job) (services []job.ServiceCtx, err error) {
	if jb.OCROracleSpec == nil {
		return nil, errors.Errorf("offchainreporting.Delegate expects an *job.OffchainreportingOracleSpec to be present, got %v", jb)
	}
	chain, err := d.chainSet.Get(jb.OCROracleSpec.EVMChainID.ToInt())
	if err != nil {
		return nil, err
	}
	concreteSpec, err := job.LoadEnvConfigVarsOCR(chain.Config(), d.keyStore.P2P(), *jb.OCROracleSpec)
	if err != nil {
		return nil, err
	}
	lggr := d.lggr.With(
		"contractAddress", concreteSpec.ContractAddress,
		"jobName", jb.Name.ValueOrZero(),
		"jobID", jb.ID,
		"externalJobID", jb.ExternalJobID)

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

	ocrDB := NewDB(d.db, concreteSpec.ID, lggr, d.cfg)

	tracker := NewOCRContractTracker(
		contract,
		contractFilterer,
		contractCaller,
		chain.Client(),
		chain.LogBroadcaster(),
		jb.ID,
		lggr,
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

	ocrLogger := logger.NewOCRWrapper(lggr, chain.Config().OCRTraceLogging(), func(msg string) {
		d.jobORM.TryRecordError(jb.ID, msg)
	})

	lc := toLocalConfig(chain.Config(), *concreteSpec)
	if err = ocr.SanityCheckLocalConfig(lc); err != nil {
		return nil, err
	}
	lggr.Info(fmt.Sprintf("OCR job using local config %+v", lc))

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
		bootstrapperCtx := job.NewServiceAdapter(bootstrapper)
		services = append(services, bootstrapperCtx)
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

		strategy := txmgr.NewQueueingTxStrategy(jb.ExternalJobID, chain.Config().OCRDefaultTransactionQueueDepth())

		var checker txmgr.TransmitCheckerSpec
		if chain.Config().OCRSimulateTransactions() {
			checker.CheckerType = txmgr.TransmitCheckerTypeSimulate
		}

		if concreteSpec.TransmitterAddress == nil {
			return nil, errors.New("TransmitterAddress is missing")
		}

		contractTransmitter := NewOCRContractTransmitter(
			concreteSpec.ContractAddress.Address(),
			contractCaller,
			contractABI,
			ocrcommon.NewTransmitter(chain.TxManager(), concreteSpec.TransmitterAddress.Address(), chain.Config().EvmGasLimitDefault(), strategy, checker),
			chain.LogBroadcaster(),
			tracker,
			chain.ID(),
		)

		runResults := make(chan pipeline.Run, chain.Config().JobPipelineResultWriteQueueDepth())
		jb.PipelineSpec.JobName = jb.Name.ValueOrZero()
		jb.PipelineSpec.JobID = jb.ID

		var configOverrider ocrtypes.ConfigOverrider
		configOverriderService, err := d.maybeCreateConfigOverrider(lggr, chain, concreteSpec.ContractAddress)
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
				jb,
				*jb.PipelineSpec,
				lggr,
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
			MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(concreteSpec.ContractAddress.String()),
			ConfigOverrider:              configOverrider,
		})
		if err != nil {
			return nil, errors.Wrap(err, "error calling NewOracle")
		}
		oracleCtx := job.NewServiceAdapter(oracle)
		services = append(services, oracleCtx)

		// RunResultSaver needs to be started first so its available
		// to read db writes. It is stopped last after the Oracle is shut down
		// so no further runs are enqueued and we can drain the queue.
		services = append([]job.ServiceCtx{ocrcommon.NewResultRunSaver(
			runResults,
			d.pipelineRunner,
			make(chan struct{}),
			lggr,
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
