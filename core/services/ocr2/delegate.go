package ocr2

import (
	"encoding/json"
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
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/blockhashes"
	ocr2vrfconfig "github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/config"
	ocr2coordinator "github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/coordinator"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/juelsfeecoin"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/reportserializer"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/promwrapper"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/relay"
	evmrelay "github.com/smartcontractkit/chainlink/core/services/relay/evm"
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
func (d *Delegate) ServicesForSpec(jobSpec job.Job) ([]job.ServiceCtx, error) {
	spec := jobSpec.OCR2OracleSpec
	if spec == nil {
		return nil, errors.Errorf("offchainreporting2.Delegate expects an *job.Offchainreporting2OracleSpec to be present, got %v", jobSpec)
	}
	if !spec.TransmitterID.Valid {
		return nil, errors.Errorf("expected a transmitterID to be specified")
	}
	relayer, exists := d.relayers[spec.Relay]
	if !exists {
		return nil, errors.Errorf("%s relay does not exist is it enabled?", spec.Relay)
	}

	lggr := d.lggr.Named("OCR").With(
		"contractID", spec.ContractID,
		"jobName", jobSpec.Name.ValueOrZero(),
		"jobID", jobSpec.ID,
	)

	if spec.Relay == relay.EVM {
		chainIDInterface, ok := spec.RelayConfig["chainID"]
		if !ok {
			return nil, errors.New("chainID must be provided in relay config")
		}
		chainID, ok := chainIDInterface.(float64)
		if !ok {
			return nil, errors.Errorf("invalid chain type got %T want float64", chainIDInterface)
		}
		chain, err2 := d.chainSet.Get(big.NewInt(int64(chainID)))
		if err2 != nil {
			return nil, errors.Wrap(err2, "get chainset")
		}

		var sendingKeys []string
		ethSendingKeys, err2 := d.ethKs.GetAll()
		if err2 != nil {
			return nil, errors.Wrap(err2, "get eth sending keys")
		}

		// Automatically provide the node's local sending keys to the job spec.
		for _, s := range ethSendingKeys {
			sendingKeys = append(sendingKeys, s.Address.String())
		}
		spec.RelayConfig["sendingKeys"] = sendingKeys

		// effectiveTransmitterAddress is the transmitter address registered on the ocr contract. This is by default the EOA account on the node.
		// In the case of forwarding, the transmitter address is the forwarder contract deployed onchain between EOA and OCR contract.
		effectiveTransmitterAddress := spec.TransmitterID
		if jobSpec.ForwardingAllowed {
			fwdrAddress, fwderr := chain.TxManager().GetForwarderForEOA(common.HexToAddress(spec.TransmitterID.String))
			if fwderr == nil {
				effectiveTransmitterAddress = null.StringFrom(fwdrAddress.String())
			} else {
				lggr.Warnw("Skipping forwarding for job, will fallback to default behavior", "job", jobSpec.Name, "err", fwderr)
			}
		}
		spec.RelayConfig["effectiveTransmitterAddress"] = effectiveTransmitterAddress
	}

	ocrDB := NewDB(d.db, spec.ID, d.lggr, d.cfg)
	peerWrapper := d.peerWrapper
	if peerWrapper == nil {
		return nil, errors.New("cannot setup OCR2 job service, libp2p peer was missing")
	} else if !peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR2 jobs require a started and running p2p v2 peer")
	}

	ocrLogger := logger.NewOCRWrapper(lggr, true, func(msg string) {
		d.lggr.ErrorIf(d.jobORM.RecordError(jobSpec.ID, msg), "unable to record error")
	})

	lc := validate.ToLocalConfig(d.cfg, *spec)
	if err := libocr2.SanityCheckLocalConfig(lc); err != nil {
		return nil, err
	}
	d.lggr.Infow("OCR2 job using local config",
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

	var pluginOracle plugins.OraclePlugin
	var ocr2Provider types.Plugin
	switch spec.PluginType {
	case job.Median:
		medianProvider, err2 := relayer.NewMedianProvider(
			types.RelayArgs{
				ExternalJobID: jobSpec.ExternalJobID,
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
			MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID),
			OffchainConfigDigester:       medianProvider.OffchainConfigDigester(),
			OffchainKeyring:              kb,
			OnchainKeyring:               kb,
		}
		return median.NewMedianServices(jobSpec, medianProvider, d.pipelineRunner, runResults, lggr, ocrLogger, oracleArgsNoPlugin)
	case job.DKG:
		chainIDInterface, ok := jobSpec.OCR2OracleSpec.RelayConfig["chainID"]
		if !ok {
			return nil, errors.New("chainID must be provided in relay config")
		}
		chainID := int64(chainIDInterface.(float64))
		chain, err2 := d.chainSet.Get(big.NewInt(chainID))
		if err2 != nil {
			return nil, errors.Wrap(err2, "get chainset")
		}
		ocr2vrfRelayer := evmrelay.NewOCR2VRFRelayer(d.db, chain, lggr.Named("OCR2VRFRelayer"))
		dkgProvider, err2 := ocr2vrfRelayer.NewDKGProvider(
			types.RelayArgs{
				ExternalJobID: jobSpec.ExternalJobID,
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
			ContractTransmitter:          dkgProvider.ContractTransmitter(),
			ContractConfigTracker:        dkgProvider.ContractConfigTracker(),
			Database:                     ocrDB,
			LocalConfig:                  lc,
			Logger:                       ocrLogger,
			MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID),
			OffchainConfigDigester:       dkgProvider.OffchainConfigDigester(),
			OffchainKeyring:              kb,
			OnchainKeyring:               kb,
		}
		return dkg.NewDKGServices(
			jobSpec,
			dkgProvider,
			ocrLogger,
			d.dkgSignKs,
			d.dkgEncryptKs,
			chain.Client(),
			oracleArgsNoPlugin)
	case job.OCR2VRF:
		chainIDInterface, ok := jobSpec.OCR2OracleSpec.RelayConfig["chainID"]
		if !ok {
			return nil, errors.New("chainID must be provided in relay config")
		}
		chainID := int64(chainIDInterface.(float64))
		chain, err2 := d.chainSet.Get(big.NewInt(chainID))
		if err2 != nil {
			return nil, errors.Wrap(err2, "get chainset")
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

		ocr2vrfRelayer := evmrelay.NewOCR2VRFRelayer(d.db, chain, lggr.Named("OCR2VRFRelayer"))

		vrfProvider, err2 := ocr2vrfRelayer.NewOCR2VRFProvider(
			types.RelayArgs{
				ExternalJobID: jobSpec.ExternalJobID,
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
				ExternalJobID: jobSpec.ExternalJobID,
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

		juelsPerFeeCoin, err2 := juelsfeecoin.NewLinkEthPriceProvider(
			common.HexToAddress(cfg.LinkEthFeedAddress), chain.Client(), 1*time.Second)
		if err2 != nil {
			return nil, errors.Wrap(err2, "new link eth price provider")
		}

		// No need to error check here, we check these keys exist when validating
		// the configuration.
		encryptionSecretKey, _ := d.dkgEncryptKs.Get(cfg.DKGEncryptionPublicKey)
		signingSecretKey, _ := d.dkgSignKs.Get(cfg.DKGSigningPublicKey)
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
			cfg.LookbackBlocks,
			chain.LogPoller(),
			chain.Config().EvmFinalityDepth(),
		)
		if err2 != nil {
			return nil, errors.Wrap(err2, "create ocr2vrf coordinator")
		}
		l := d.lggr.Named("OCR2VRF").With(
			"jobName", jobSpec.Name.ValueOrZero(),
			"jobID", jobSpec.ID,
		)
		vrfLogger := logger.NewOCRWrapper(l.With(
			"vrfContractID", spec.ContractID), true, func(msg string) {
			d.lggr.ErrorIf(d.jobORM.RecordError(jobSpec.ID, msg), "unable to record error")
		})
		dkgLogger := logger.NewOCRWrapper(l.With(
			"dkgContractID", cfg.DKGContractAddress), true, func(msg string) {
			d.lggr.ErrorIf(d.jobORM.RecordError(jobSpec.ID, msg), "unable to record error")
		})
		dkgReportingPluginFactoryDecorator := func(wrapped ocr2types.ReportingPluginFactory) ocr2types.ReportingPluginFactory {
			return promwrapper.NewPromFactory(wrapped, "DKG", string(relay.EVM), chain.ID())
		}
		vrfReportingPluginFactoryDecorator := func(wrapped ocr2types.ReportingPluginFactory) ocr2types.ReportingPluginFactory {
			return promwrapper.NewPromFactory(wrapped, "OCR2VRF", string(relay.EVM), chain.ID())
		}
		oracles, err2 := ocr2vrf.NewOCR2VRF(ocr2vrf.DKGVRFArgs{
			VRFLogger:                          vrfLogger,
			DKGLogger:                          dkgLogger,
			BinaryNetworkEndpointFactory:       peerWrapper.Peer2,
			V2Bootstrappers:                    bootstrapPeers,
			OffchainKeyring:                    kb,
			OnchainKeyring:                     kb,
			VRFOffchainConfigDigester:          vrfProvider.OffchainConfigDigester(),
			VRFContractConfigTracker:           vrfProvider.ContractConfigTracker(),
			VRFContractTransmitter:             vrfProvider.ContractTransmitter(),
			VRFDatabase:                        ocrDB,
			VRFLocalConfig:                     lc,
			VRFMonitoringEndpoint:              d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID),
			DKGContractConfigTracker:           dkgProvider.ContractConfigTracker(),
			DKGOffchainConfigDigester:          dkgProvider.OffchainConfigDigester(),
			DKGContract:                        dkgpkg.NewOnchainContract(dkgContract, &altbn_128.G2{}),
			DKGContractTransmitter:             dkgProvider.ContractTransmitter(),
			DKGDatabase:                        ocrDB,
			DKGLocalConfig:                     lc,
			DKGMonitoringEndpoint:              d.monitoringEndpointGen.GenMonitoringEndpoint(cfg.DKGContractAddress),
			Blockhashes:                        blockhashes.NewFixedBlockhashProvider(chain.LogPoller(), d.lggr, 256),
			Serializer:                         reportserializer.NewReportSerializer(&altbn_128.G1{}),
			JulesPerFeeCoin:                    juelsPerFeeCoin,
			Coordinator:                        coordinator,
			Esk:                                encryptionSecretKey.KyberScalar(),
			Ssk:                                signingSecretKey.KyberScalar(),
			KeyID:                              keyID,
			DKGReportingPluginFactoryDecorator: dkgReportingPluginFactoryDecorator,
			VRFReportingPluginFactoryDecorator: vrfReportingPluginFactoryDecorator,
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
		)

		// NOTE: we return from here with the services because the OCR2VRF oracles are defined
		// and exported from the ocr2vrf library. It takes care of running the DKG and OCR2VRF
		// oracles under the hood together.
		oracleCtx := job.NewServiceAdapter(oracles)
		return []job.ServiceCtx{runResultSaver, vrfProvider, dkgProvider, oracleCtx}, nil
	case job.OCR2Keeper:
		keeperProvider, rgstry, encoder, logProvider, err2 := ocr2keeper.EVMDependencies(jobSpec, d.db, lggr, d.chainSet, d.pipelineRunner)
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
			MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID),
			OffchainConfigDigester:       keeperProvider.OffchainConfigDigester(),
			OffchainKeyring:              kb,
			OnchainKeyring:               kb,
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
		)

		return []job.ServiceCtx{
			runResultSaver,
			keeperProvider,
			pluginService,
		}, nil
	case job.OCR2DirectRequest:
		// TODO: relayer for DR-OCR plugin: https://app.shortcut.com/chainlinklabs/story/54051/relayer-for-the-ocr-plugin
		drProvider, err2 := relayer.NewMedianProvider(
			types.RelayArgs{
				ExternalJobID: jobSpec.ExternalJobID,
				JobID:         spec.ID,
				ContractID:    spec.ContractID,
				RelayConfig:   spec.RelayConfig.Bytes(),
			}, types.PluginArgs{
				TransmitterID: spec.TransmitterID.String,
				PluginConfig:  spec.PluginConfig.Bytes(),
			})
		if err2 != nil {
			return nil, err2
		}
		ocr2Provider = drProvider

		var relayConfig evmrelay.RelayConfig
		err2 = json.Unmarshal(spec.RelayConfig.Bytes(), &relayConfig)
		if err2 != nil {
			return nil, err2
		}
		chain, err2 := d.chainSet.Get(relayConfig.ChainID.ToInt())
		if err2 != nil {
			return nil, err2
		}
		// TODO replace with a DB: https://app.shortcut.com/chainlinklabs/story/54049/database-table-in-core-node
		pluginORM := drocr_service.NewInMemoryORM()
		pluginOracle, _ = directrequestocr.NewDROracle(jobSpec, d.pipelineRunner, d.jobORM, ocr2Provider, pluginORM, chain, lggr, ocrLogger)
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
		MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID),
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
		lggr)

	oracleCtx := job.NewServiceAdapter(oracle)
	return append([]job.ServiceCtx{runResultSaver, ocr2Provider, oracleCtx}, pluginServices...), nil
}
