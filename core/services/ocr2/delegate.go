package ocr2

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"
	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	dkgpkg "github.com/smartcontractkit/ocr2vrf/dkg"
	"github.com/smartcontractkit/ocr2vrf/ocr2vrf"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/dkg"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/median"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/blockhashes"
	ocr2vrfconfig "github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/config"
	ocr2coordinator "github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/coordinator"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/juelsfeecoin"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/reportserializer"
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
	relayers              map[relay.Network]types.Relayer
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
		relayers,
	}
}

func (d Delegate) JobType() job.Type {
	return job.OffchainReporting2
}

func (Delegate) OnJobCreated(spec job.Job) {}
func (Delegate) OnJobDeleted(spec job.Job) {}

func (Delegate) AfterJobCreated(spec job.Job)  {}
func (Delegate) BeforeJobDeleted(spec job.Job) {}

// ServicesForSpec returns the OCR2 services that need to run for this job
func (d Delegate) ServicesForSpec(jobSpec job.Job) ([]job.ServiceCtx, error) {
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

	ocrDB := NewDB(d.db, spec.ID, d.lggr, d.cfg)
	peerWrapper := d.peerWrapper
	if peerWrapper == nil {
		return nil, errors.New("cannot setup OCR2 job service, libp2p peer was missing")
	} else if !peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR2 jobs require a started and running peer. Did you forget to specify P2P_LISTEN_PORT?")
	}

	lggr := d.lggr.Named("OCR").With(
		"contractID", spec.ContractID,
		"jobName", jobSpec.Name.ValueOrZero(),
		"jobID", jobSpec.ID,
	)
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
				RelayConfig:   spec.RelayConfig.Bytes(),
			}, types.PluginArgs{
				TransmitterID: spec.TransmitterID.String,
				PluginConfig:  spec.PluginConfig.Bytes(),
			})
		if err2 != nil {
			return nil, err2
		}
		ocr2Provider = medianProvider
		pluginOracle, err = median.NewMedian(jobSpec, medianProvider, d.pipelineRunner, runResults, lggr, ocrLogger)
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
				RelayConfig:   spec.RelayConfig.Bytes(),
			}, types.PluginArgs{
				TransmitterID: spec.TransmitterID.String,
				PluginConfig:  spec.PluginConfig.Bytes(),
			})
		if err2 != nil {
			return nil, err2
		}
		ocr2Provider = dkgProvider
		pluginOracle, err = dkg.NewDKG(
			jobSpec,
			dkgProvider,
			lggr.Named("DKG"),
			ocrLogger,
			d.dkgSignKs,
			d.dkgEncryptKs,
			chain.Client())
		if err != nil {
			return nil, errors.Wrap(err, "error while instantiating DKG")
		}
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
			common.HexToAddress(cfg.DKGContractAddress),
			chain.Client(),
			cfg.LookbackBlocks,
			chain.LogPoller(),
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
			VRFMonitoringEndpoint:        d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID),
			DKGContractConfigTracker:     dkgProvider.ContractConfigTracker(),
			DKGOffchainConfigDigester:    dkgProvider.OffchainConfigDigester(),
			DKGContract:                  dkgpkg.NewOnchainContract(dkgContract, &altbn_128.G2{}),
			DKGContractTransmitter:       dkgProvider.ContractTransmitter(),
			DKGDatabase:                  ocrDB,
			DKGLocalConfig:               lc,
			DKGMonitoringEndpoint:        d.monitoringEndpointGen.GenMonitoringEndpoint(cfg.DKGContractAddress),
			Blockhashes:                  blockhashes.NewFixedBlockhashProvider(chain.Client(), 256, 256),
			Serializer:                   reportserializer.NewReportSerializer(&altbn_128.G1{}),
			JulesPerFeeCoin:              juelsPerFeeCoin,
			Coordinator:                  coordinator,
			Esk:                          encryptionSecretKey.KyberScalar(),
			Ssk:                          signingSecretKey.KyberScalar(),
			KeyID:                        keyID,
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
			lggr)

		// NOTE: we return from here with the services because the OCR2VRF oracles are defined
		// and exported from the ocr2vrf library. It takes care of running the DKG and OCR2VRF
		// oracles under the hood together.
		oracleCtx := job.NewServiceAdapter(oracles)
		return []job.ServiceCtx{runResultSaver, vrfProvider, oracleCtx}, nil
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
