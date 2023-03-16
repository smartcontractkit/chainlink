package ocr2

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	dkgpkg "github.com/smartcontractkit/ocr2vrf/dkg"
	"github.com/smartcontractkit/ocr2vrf/ocr2vrf"
	"github.com/smartcontractkit/sqlx"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/dkg/persistence"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/utils"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	drocr_service "github.com/smartcontractkit/chainlink/core/services/directrequestocr"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/directrequestocr"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/dkg"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/median"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2keeper"
	ocr2vrfconfig "github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/config"
	ocr2coordinator "github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/coordinator"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/juelsfeecoin"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/reasonablegasprice"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/reportserializer"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/promwrapper"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/relay"
	evmrelay "github.com/smartcontractkit/chainlink/core/services/relay/evm"
	evmrelaytypes "github.com/smartcontractkit/chainlink/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/telemetry"
)

type Delegate struct {
	db                    *sqlx.DB
	jobORM                job.ORM
	pipelineRunner        pipeline.Runner
	peerWrapper           *ocrcommon.SingletonPeerWrapper
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	chainSet              evm.ChainSet
	cfg                   validate.Config
	lggr                  logger.Logger
	ks                    keystore.OCR2
	dkgSignKs             keystore.DKGSign
	dkgEncryptKs          keystore.DKGEncrypt
	ethKs                 keystore.Eth
	relayers              map[relay.Network]types.Relayer
	isNewlyCreatedJob     bool // Set to true if this is a new job freshly added, false if job was present already on node boot.
	mailMon               *utils.MailboxMonitor
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
	cfg validate.Config,
	ks keystore.OCR2,
	dkgSignKs keystore.DKGSign,
	dkgEncryptKs keystore.DKGEncrypt,
	ethKs keystore.Eth,
	relayers map[relay.Network]types.Relayer,
	mailMon *utils.MailboxMonitor,
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
		dkgSignKs,
		dkgEncryptKs,
		ethKs,
		relayers,
		false,
		mailMon,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.OffchainReporting2
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {
	// This is only called first time the job is created
	d.isNewlyCreatedJob = true
}
func (d *Delegate) AfterJobCreated(spec job.Job)  {}
func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

// ServicesForSpec returns the OCR2 services that need to run for this job
func (d *Delegate) ServicesForSpec(jb job.Job) ([]job.ServiceCtx, error) {
	spec := jb.OCR2OracleSpec
	if spec == nil {
		return nil, errors.Errorf("offchainreporting2.Delegate expects an *job.Offchainreporting2OracleSpec to be present, got %v", jb)
	}
	if !spec.TransmitterID.Valid {
		return nil, errors.Errorf("expected a transmitterID to be specified")
	}
	relayer, exists := d.relayers[spec.Relay]
	if !exists {
		return nil, errors.Errorf("%s relay does not exist is it enabled?", spec.Relay)
	}

	lggr := logger.Sugared(d.lggr.Named("OCR").With(
		"contractID", spec.ContractID,
		"jobName", jb.Name.ValueOrZero(),
		"jobID", jb.ID,
	))

	if spec.Relay == relay.EVM {
		chainID, err2 := spec.RelayConfig.EVMChainID()
		if err2 != nil {
			return nil, err2
		}
		chain, err2 := d.chainSet.Get(big.NewInt(chainID))
		if err2 != nil {
			return nil, errors.Wrap(err2, "get chainset")
		}

		spec.RelayConfig["sendingKeys"] = []string{spec.TransmitterID.String}

		// effectiveTransmitterAddress is the transmitter address registered on the ocr contract. This is by default the EOA account on the node.
		// In the case of forwarding, the transmitter address is the forwarder contract deployed onchain between EOA and OCR contract.
		effectiveTransmitterAddress := spec.TransmitterID
		if jb.ForwardingAllowed {
			fwdrAddress, fwderr := chain.TxManager().GetForwarderForEOA(common.HexToAddress(spec.TransmitterID.String))
			if fwderr == nil {
				effectiveTransmitterAddress = null.StringFrom(fwdrAddress.String())
			} else {
				lggr.Warnw("Skipping forwarding for job, will fallback to default behavior", "job", jb.Name, "err", fwderr)
			}
		}
		spec.RelayConfig["effectiveTransmitterAddress"] = effectiveTransmitterAddress
	}

	ocrDB := NewDB(d.db, spec.ID, lggr, d.cfg)
	peerWrapper := d.peerWrapper
	if peerWrapper == nil {
		return nil, errors.New("cannot setup OCR2 job service, libp2p peer was missing")
	} else if !peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR2 jobs require a started and running p2p v2 peer")
	}

	ocrLogger := logger.NewOCRWrapper(lggr, true, func(msg string) {
		lggr.ErrorIf(d.jobORM.RecordError(jb.ID, msg), "unable to record error")
	})

	lc := validate.ToLocalConfig(d.cfg, *spec)
	if err := libocr2.SanityCheckLocalConfig(lc); err != nil {
		return nil, err
	}
	lggr.Infow("OCR2 job using local config",
		"BlockchainTimeout", lc.BlockchainTimeout,
		"ContractConfigConfirmations", lc.ContractConfigConfirmations,
		"ContractConfigTrackerPollInterval", lc.ContractConfigTrackerPollInterval,
		"ContractTransmitterTransmitTimeout", lc.ContractTransmitterTransmitTimeout,
		"DatabaseTimeout", lc.DatabaseTimeout,
	)

	bootstrapPeers, err := ocrcommon.GetValidatedBootstrapPeers(spec.P2PV2Bootstrappers, peerWrapper.Config().P2PV2Bootstrappers())
	if err != nil {
		return nil, err
	}
	lggr.Debugw("Using bootstrap peers", "peers", bootstrapPeers)
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

	var pluginOracle plugins.OraclePlugin
	var ocr2Provider types.Plugin
	switch spec.PluginType {
	case job.Median:
		medianProvider, err2 := relayer.NewMedianProvider(
			types.RelayArgs{
				ExternalJobID: jb.ExternalJobID,
				JobID:         spec.ID,
				ContractID:    spec.ContractID,
				New:           d.isNewlyCreatedJob,
				RelayConfig:   spec.RelayConfig.Bytes(),
			}, types.PluginArgs{
				TransmitterID: spec.TransmitterID.String,
				PluginConfig:  spec.PluginConfig.Bytes(),
			})
		if err2 != nil {
			return nil, err2
		}
		oracleArgsNoPlugin := libocr2.OracleArgs{
			BinaryNetworkEndpointFactory: peerWrapper.Peer2,
			V2Bootstrappers:              bootstrapPeers,
			ContractTransmitter:          medianProvider.ContractTransmitter(),
			ContractConfigTracker:        medianProvider.ContractConfigTracker(),
			Database:                     ocrDB,
			LocalConfig:                  lc,
			Logger:                       ocrLogger,
			MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID, synchronization.OCR2Median),
			OffchainConfigDigester:       medianProvider.OffchainConfigDigester(),
			OffchainKeyring:              kb,
			OnchainKeyring:               kb,
		}
		return median.NewMedianServices(jb, medianProvider, d.pipelineRunner, runResults, lggr, ocrLogger, oracleArgsNoPlugin, d.cfg)
	case job.DKG:
		chainID, err2 := spec.RelayConfig.EVMChainID()
		if err2 != nil {
			return nil, err2
		}
		chain, err2 := d.chainSet.Get(big.NewInt(chainID))
		if err2 != nil {
			return nil, errors.Wrap(err2, "get chainset")
		}
		ocr2vrfRelayer := evmrelay.NewOCR2VRFRelayer(d.db, chain, lggr.Named("OCR2VRFRelayer"), d.ethKs)
		dkgProvider, err2 := ocr2vrfRelayer.NewDKGProvider(
			types.RelayArgs{
				ExternalJobID: jb.ExternalJobID,
				JobID:         spec.ID,
				ContractID:    spec.ContractID,
				New:           d.isNewlyCreatedJob,
				RelayConfig:   spec.RelayConfig.Bytes(),
			}, types.PluginArgs{
				TransmitterID: spec.TransmitterID.String,
				PluginConfig:  spec.PluginConfig.Bytes(),
			})
		if err2 != nil {
			return nil, err2
		}
		noopMonitoringEndpoint := telemetry.NoopAgent{}
		oracleArgsNoPlugin := libocr2.OracleArgs{
			BinaryNetworkEndpointFactory: peerWrapper.Peer2,
			V2Bootstrappers:              bootstrapPeers,
			ContractTransmitter:          dkgProvider.ContractTransmitter(),
			ContractConfigTracker:        dkgProvider.ContractConfigTracker(),
			Database:                     ocrDB,
			LocalConfig:                  lc,
			Logger:                       ocrLogger,
			// Telemetry ingress for DKG is currently not supported so a noop monitoring endpoint is being used
			MonitoringEndpoint:     &noopMonitoringEndpoint,
			OffchainConfigDigester: dkgProvider.OffchainConfigDigester(),
			OffchainKeyring:        kb,
			OnchainKeyring:         kb,
		}
		return dkg.NewDKGServices(
			jb,
			dkgProvider,
			lggr,
			ocrLogger,
			d.dkgSignKs,
			d.dkgEncryptKs,
			chain.Client(),
			oracleArgsNoPlugin,
			d.db,
			d.cfg,
			big.NewInt(chainID),
			spec.Relay,
		)
	case job.OCR2VRF:
		chainID, err2 := spec.RelayConfig.EVMChainID()
		if err2 != nil {
			return nil, err2
		}

		// Automatically provide the node's local sending keys to the job spec for OCR2VRF.
		var sendingKeys []string
		ethSendingKeys, err2 := d.ethKs.EnabledKeysForChain(big.NewInt(chainID))
		if err2 != nil {
			return nil, errors.Wrap(err2, "get eth sending keys")
		}
		for _, s := range ethSendingKeys {
			sendingKeys = append(sendingKeys, s.Address.String())
		}
		spec.RelayConfig["sendingKeys"] = sendingKeys

		chain, err2 := d.chainSet.Get(big.NewInt(chainID))
		if err2 != nil {
			return nil, errors.Wrap(err2, "get chainset")
		}
		if jb.ForwardingAllowed != chain.Config().EvmUseForwarders() {
			return nil, errors.New("transaction forwarding settings must be consistent for ocr2vrf")
		}

		var cfg ocr2vrfconfig.PluginConfig
		err2 = json.Unmarshal(spec.PluginConfig.Bytes(), &cfg)
		if err2 != nil {
			return nil, errors.Wrap(err2, "unmarshal ocr2vrf plugin config")
		}

		err2 = ocr2vrfconfig.ValidatePluginConfig(cfg, d.dkgSignKs, d.dkgEncryptKs)
		if err2 != nil {
			return nil, errors.Wrap(err2, "validate ocr2vrf plugin config")
		}

		ocr2vrfRelayer := evmrelay.NewOCR2VRFRelayer(d.db, chain, lggr.Named("OCR2VRFRelayer"), d.ethKs)

		vrfProvider, err2 := ocr2vrfRelayer.NewOCR2VRFProvider(
			types.RelayArgs{
				ExternalJobID: jb.ExternalJobID,
				JobID:         spec.ID,
				ContractID:    spec.ContractID,
				New:           d.isNewlyCreatedJob,
				RelayConfig:   spec.RelayConfig.Bytes(),
			}, types.PluginArgs{
				TransmitterID: spec.TransmitterID.String,
				PluginConfig:  spec.PluginConfig.Bytes(),
			})
		if err2 != nil {
			return nil, errors.Wrap(err2, "new vrf provider")
		}

		dkgProvider, err2 := ocr2vrfRelayer.NewDKGProvider(
			types.RelayArgs{
				ExternalJobID: jb.ExternalJobID,
				JobID:         spec.ID,
				ContractID:    cfg.DKGContractAddress,
				RelayConfig:   spec.RelayConfig.Bytes(),
			}, types.PluginArgs{
				TransmitterID: spec.TransmitterID.String,
				PluginConfig:  spec.PluginConfig.Bytes(),
			})
		if err2 != nil {
			return nil, errors.Wrap(err2, "new dkg provider")
		}

		dkgContract, err2 := dkg.NewOnchainDKGClient(cfg.DKGContractAddress, chain.Client())
		if err2 != nil {
			return nil, errors.Wrap(err2, "new onchain dkg client")
		}

		timeout := 1 * time.Second
		juelsPerFeeCoin, err2 := juelsfeecoin.NewLinkEthPriceProvider(
			common.HexToAddress(cfg.LinkEthFeedAddress), chain.Client(), timeout)
		if err2 != nil {
			return nil, errors.Wrap(err2, "new link eth price provider")
		}

		reasonableGasPrice := reasonablegasprice.NewReasonableGasPriceProvider(
			chain.TxManager().GetGasEstimator(),
			timeout,
			chain.Config().EvmMaxGasPriceWei(),
			chain.Config().EvmEIP1559DynamicFees(),
		)

		encryptionSecretKey, err2 := d.dkgEncryptKs.Get(cfg.DKGEncryptionPublicKey)
		if err2 != nil {
			return nil, errors.Wrap(err2, "get DKG encryption key")
		}
		signingSecretKey, err2 := d.dkgSignKs.Get(cfg.DKGSigningPublicKey)
		if err2 != nil {
			return nil, errors.Wrap(err2, "get DKG signing key")
		}
		keyID, err2 := dkg.DecodeKeyID(cfg.DKGKeyID)
		if err2 != nil {
			return nil, errors.Wrap(err2, "decode DKG key ID")
		}

		coordinator, err2 := ocr2coordinator.New(
			lggr.Named("OCR2VRFCoordinator"),
			common.HexToAddress(spec.ContractID),
			common.HexToAddress(cfg.VRFCoordinatorAddress),
			common.HexToAddress(cfg.DKGContractAddress),
			chain.Client(),
			chain.LogPoller(),
			chain.Config().EvmFinalityDepth(),
		)
		if err2 != nil {
			return nil, errors.Wrap(err2, "create ocr2vrf coordinator")
		}
		l := lggr.Named("OCR2VRF").With(
			"jobName", jb.Name.ValueOrZero(),
			"jobID", jb.ID,
		)
		vrfLogger := logger.NewOCRWrapper(l.With(
			"vrfContractID", spec.ContractID), true, func(msg string) {
			lggr.ErrorIf(d.jobORM.RecordError(jb.ID, msg), "unable to record error")
		})
		dkgLogger := logger.NewOCRWrapper(l.With(
			"dkgContractID", cfg.DKGContractAddress), true, func(msg string) {
			lggr.ErrorIf(d.jobORM.RecordError(jb.ID, msg), "unable to record error")
		})
		dkgReportingPluginFactoryDecorator := func(wrapped ocr2types.ReportingPluginFactory) ocr2types.ReportingPluginFactory {
			return promwrapper.NewPromFactory(wrapped, "DKG", string(relay.EVM), chain.ID())
		}
		vrfReportingPluginFactoryDecorator := func(wrapped ocr2types.ReportingPluginFactory) ocr2types.ReportingPluginFactory {
			return promwrapper.NewPromFactory(wrapped, "OCR2VRF", string(relay.EVM), chain.ID())
		}
		noopMonitoringEndpoint := telemetry.NoopAgent{}
		oracles, err2 := ocr2vrf.NewOCR2VRF(ocr2vrf.DKGVRFArgs{
			VRFLogger:                    vrfLogger,
			DKGLogger:                    dkgLogger,
			BinaryNetworkEndpointFactory: peerWrapper.Peer2,
			V2Bootstrappers:              bootstrapPeers,
			OffchainKeyring:              kb,
			OnchainKeyring:               kb,
			VRFOffchainConfigDigester:    vrfProvider.OffchainConfigDigester(),
			VRFContractConfigTracker:     vrfProvider.ContractConfigTracker(),
			VRFContractTransmitter:       vrfProvider.ContractTransmitter(),
			VRFDatabase:                  ocrDB,
			VRFLocalConfig:               lc,
			VRFMonitoringEndpoint:        d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID, synchronization.OCR2VRF),
			DKGContractConfigTracker:     dkgProvider.ContractConfigTracker(),
			DKGOffchainConfigDigester:    dkgProvider.OffchainConfigDigester(),
			DKGContract:                  dkgpkg.NewOnchainContract(dkgContract, &altbn_128.G2{}),
			DKGContractTransmitter:       dkgProvider.ContractTransmitter(),
			DKGDatabase:                  ocrDB,
			DKGLocalConfig:               lc,
			// Telemetry ingress for DKG is currently not supported so a noop monitoring endpoint is being used
			DKGMonitoringEndpoint:              &noopMonitoringEndpoint,
			Serializer:                         reportserializer.NewReportSerializer(&altbn_128.G1{}),
			JuelsPerFeeCoin:                    juelsPerFeeCoin,
			ReasonableGasPrice:                 reasonableGasPrice,
			Coordinator:                        coordinator,
			Esk:                                encryptionSecretKey.KyberScalar(),
			Ssk:                                signingSecretKey.KyberScalar(),
			KeyID:                              keyID,
			DKGReportingPluginFactoryDecorator: dkgReportingPluginFactoryDecorator,
			VRFReportingPluginFactoryDecorator: vrfReportingPluginFactoryDecorator,
			DKGSharePersistence:                persistence.NewShareDB(d.db, lggr.Named("DKGShareDB"), d.cfg, big.NewInt(chainID), spec.Relay),
		})
		if err2 != nil {
			return nil, errors.Wrap(err2, "new ocr2vrf")
		}

		// RunResultSaver needs to be started first, so it's available
		// to read odb writes. It is stopped last after the OraclePlugin is shut down
		// so no further runs are enqueued, and we can drain the queue.
		runResultSaver := ocrcommon.NewResultRunSaver(
			runResults,
			d.pipelineRunner,
			make(chan struct{}),
			lggr,
			d.cfg.JobPipelineMaxSuccessfulRuns(),
		)

		// NOTE: we return from here with the services because the OCR2VRF oracles are defined
		// and exported from the ocr2vrf library. It takes care of running the DKG and OCR2VRF
		// oracles under the hood together.
		oracleCtx := job.NewServiceAdapter(oracles)
		return []job.ServiceCtx{runResultSaver, vrfProvider, dkgProvider, oracleCtx}, nil
	case job.OCR2Keeper:
		keeperProvider, rgstry, encoder, logProvider, err2 := ocr2keeper.EVMDependencies(jb, d.db, lggr, d.chainSet, d.pipelineRunner)
		if err2 != nil {
			return nil, errors.Wrap(err2, "could not build dependencies for ocr2 keepers")
		}

		var cfg ocr2keeper.PluginConfig
		err2 = json.Unmarshal(spec.PluginConfig.Bytes(), &cfg)
		if err2 != nil {
			return nil, errors.Wrap(err2, "unmarshal ocr2keepers plugin config")
		}

		err2 = ocr2keeper.ValidatePluginConfig(cfg)
		if err2 != nil {
			return nil, errors.Wrap(err2, "ocr2keepers plugin config validation failure")
		}

		conf := ocr2keepers.DelegateConfig{
			BinaryNetworkEndpointFactory: peerWrapper.Peer2,
			V2Bootstrappers:              bootstrapPeers,
			ContractTransmitter:          keeperProvider.ContractTransmitter(),
			ContractConfigTracker:        keeperProvider.ContractConfigTracker(),
			KeepersDatabase:              ocrDB,
			LocalConfig:                  lc,
			Logger:                       ocrLogger,
			MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID, synchronization.OCR2Automation),
			OffchainConfigDigester:       keeperProvider.OffchainConfigDigester(),
			OffchainKeyring:              kb,
			OnchainKeyring:               kb,
			HeadSubscriber:               rgstry,
			Registry:                     rgstry,
			ReportEncoder:                encoder,
			PerformLogProvider:           logProvider,
			CacheExpiration:              cfg.CacheExpiration.Value(),
			CacheEvictionInterval:        cfg.CacheEvictionInterval.Value(),
			MaxServiceWorkers:            cfg.MaxServiceWorkers,
			ServiceQueueLength:           cfg.ServiceQueueLength,
		}
		pluginService, err2 := ocr2keepers.NewDelegate(conf)
		if err2 != nil {
			return nil, errors.Wrap(err, "could not create new keepers ocr2 delegate")
		}

		// RunResultSaver needs to be started first, so it's available
		// to read odb writes. It is stopped last after the OraclePlugin is shut down
		// so no further runs are enqueued, and we can drain the queue.
		runResultSaver := ocrcommon.NewResultRunSaver(
			runResults,
			d.pipelineRunner,
			make(chan struct{}),
			lggr,
			d.cfg.JobPipelineMaxSuccessfulRuns(),
		)

		return []job.ServiceCtx{
			runResultSaver,
			keeperProvider,
			rgstry,
			logProvider,
			pluginService,
		}, nil
	case job.OCR2Functions:
		if spec.Relay != relay.EVM {
			return nil, fmt.Errorf("unsupported relay: %s", spec.Relay)
		}
		drProvider, err2 := evmrelay.NewOCR2DRProvider(
			d.chainSet,
			types.RelayArgs{
				ExternalJobID: jb.ExternalJobID,
				JobID:         spec.ID,
				ContractID:    spec.ContractID,
				RelayConfig:   spec.RelayConfig.Bytes(),
				New:           d.isNewlyCreatedJob,
			},
			types.PluginArgs{
				TransmitterID: spec.TransmitterID.String,
				PluginConfig:  spec.PluginConfig.Bytes(),
			},
			lggr.Named("OCR2DRRelayer"),
			d.ethKs,
		)
		if err2 != nil {
			return nil, err2
		}
		ocr2Provider = drProvider

		var relayConfig evmrelaytypes.RelayConfig
		err2 = json.Unmarshal(spec.RelayConfig.Bytes(), &relayConfig)
		if err2 != nil {
			return nil, err2
		}
		chain, err2 := d.chainSet.Get(relayConfig.ChainID.ToInt())
		if err2 != nil {
			return nil, err2
		}
		pluginORM := drocr_service.NewORM(d.db, lggr, d.cfg, common.HexToAddress(spec.ContractID))
		pluginOracle, err2 = directrequestocr.NewDROracle(jb, d.pipelineRunner, d.jobORM, pluginORM, chain, lggr, ocrLogger, d.mailMon)
		if err2 != nil {
			return nil, err2
		}
	default:
		return nil, errors.Errorf("plugin type %s not supported", spec.PluginType)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialise plugin")
	}

	pluginFactory, err := pluginOracle.GetPluginFactory()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get plugin factory")
	}

	pluginServices, err := pluginOracle.GetServices()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get plugin services")
	}

	oracle, err := libocr2.NewOracle(libocr2.OracleArgs{
		BinaryNetworkEndpointFactory: peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		ContractTransmitter:          ocr2Provider.ContractTransmitter(),
		ContractConfigTracker:        ocr2Provider.ContractConfigTracker(),
		Database:                     ocrDB,
		LocalConfig:                  lc,
		Logger:                       ocrLogger,
		MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID, synchronization.OCR2Functions),
		OffchainConfigDigester:       ocr2Provider.OffchainConfigDigester(),
		OffchainKeyring:              kb,
		OnchainKeyring:               kb,
		ReportingPluginFactory:       pluginFactory,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error calling NewOracle")
	}

	// RunResultSaver needs to be started first, so it's available
	// to read odb writes. It is stopped last after the OraclePlugin is shut down
	// so no further runs are enqueued, and we can drain the queue.
	runResultSaver := ocrcommon.NewResultRunSaver(
		runResults,
		d.pipelineRunner,
		make(chan struct{}),
		lggr,
		d.cfg.JobPipelineMaxSuccessfulRuns(),
	)

	oracleCtx := job.NewServiceAdapter(oracle)
	return append([]job.ServiceCtx{runResultSaver, ocr2Provider, oracleCtx}, pluginServices...), nil
}
