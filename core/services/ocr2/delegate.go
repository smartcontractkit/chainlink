package ocr2

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	dkgpkg "github.com/smartcontractkit/ocr2vrf/dkg"
	"github.com/smartcontractkit/ocr2vrf/ocr2vrf"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/dkg"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/dkg/persistence"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper"
	ocr2vrfconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2vrf/config"
	ocr2coordinator "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2vrf/coordinator"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2vrf/juelsfeecoin"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2vrf/reasonablegasprice"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2vrf/reportserializer"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/promwrapper"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
func (d *Delegate) OnDeleteJob(jb job.Job, q pg.Queryer) error {
	// If the job spec is malformed in any way, we report the error but return nil so that
	//  the job deletion itself isn't blocked.  However, if UnregisterFilter returns an
	//  error, that means it failed to remove a valid active filter from the db.  We do abort the job deletion
	//  in that case, since it should be easy for the user to retry and will avoid leaving the db in
	//  an inconsistent state.  This assumes UnregisterFilter will return nil if the filter wasn't found
	//  at all (no rows deleted).

	spec := jb.OCR2OracleSpec
	if spec == nil {
		d.lggr.Errorf("offchainreporting2.Delegate.OnDeleteJob called with wrong job type, ignoring non-OCR2 spec %v", jb)
		return nil
	}
	if spec.Relay != relay.EVM {
		return nil
	}

	chainID, err := spec.RelayConfig.EVMChainID()
	if err != nil {
		d.lggr.Errorf("OCR2 jobs spec missing chainID")
		return nil
	}
	chain, err := d.chainSet.Get(big.NewInt(chainID))
	if err != nil {
		d.lggr.Error(err)
		return nil
	}
	lp := chain.LogPoller()

	var filters []string
	switch spec.PluginType {
	case job.OCR2VRF:
		filters, err = ocr2coordinator.FilterNamesFromSpec(spec)
		if err != nil {
			d.lggr.Errorw("failed to derive ocr2vrf filter names from spec", "err", err, "spec", spec)
		}
	case job.OCR2Keeper:
		filters, err = ocr2keeper.FilterNamesFromSpec(spec)
		if err != nil {
			d.lggr.Errorw("failed to derive ocr2keeper filter names from spec", "err", err, "spec", spec)
		}
	default:
		return nil
	}

	rargs := types.RelayArgs{
		ExternalJobID: jb.ExternalJobID,
		JobID:         spec.ID,
		ContractID:    spec.ContractID,
		New:           false,
		RelayConfig:   spec.RelayConfig.Bytes(),
	}

	relayFilters, err := evmrelay.FilterNamesFromRelayArgs(rargs)
	if err != nil {
		d.lggr.Errorw("Failed to derive evm relay filter names from relay args", "err", err, "rargs", rargs)
		return nil
	}

	filters = append(filters, relayFilters...)

	for _, filter := range filters {
		d.lggr.Debugf("Unregistering %s filter", filter)
		err = lp.UnregisterFilter(filter, q)
		if err != nil {
			return errors.Wrapf(err, "Failed to unregister filter %s", filter)
		}
	}
	return nil
}

// ServicesForSpec returns the OCR2 services that need to run for this job
func (d *Delegate) ServicesForSpec(jb job.Job) ([]job.ServiceCtx, error) {
	spec := jb.OCR2OracleSpec
	if spec == nil {
		return nil, errors.Errorf("offchainreporting2.Delegate expects an *job.Offchainreporting2OracleSpec to be present, got %v", jb)
	}
	if !spec.TransmitterID.Valid {
		return nil, errors.Errorf("expected a transmitterID to be specified")
	}
	transmitterID := spec.TransmitterID.String
	relayer, exists := d.relayers[spec.Relay]
	if !exists {
		return nil, errors.Errorf("%s relay does not exist is it enabled?", spec.Relay)
	}
	effectiveTransmitterID := transmitterID

	lggr := logger.Sugared(d.lggr.Named("OCR").With(
		"contractID", spec.ContractID,
		"jobName", jb.Name.ValueOrZero(),
		"jobID", jb.ID,
	))
	feedID := spec.FeedID
	if feedID != (common.Hash{}) {
		lggr = logger.Sugared(lggr.With("feedID", spec.FeedID))
		spec.RelayConfig["feedID"] = feedID
	}

	if spec.PluginType == job.Mercury {
		if feedID == (common.Hash{}) {
			return nil, errors.Errorf("ServicesForSpec: mercury job type requires feedID")
		}
		if len(transmitterID) != 64 {
			return nil, errors.Errorf("ServicesForSpec: mercury job type requires transmitter ID to be a 32-byte hex string, got: %q", transmitterID)
		}
		if _, err := hex.DecodeString(transmitterID); err != nil {
			return nil, errors.Wrapf(err, "ServicesForSpec: mercury job type requires transmitter ID to be a 32-byte hex string, got: %q", transmitterID)
		}
	}

	if spec.Relay == relay.EVM {
		chainID, err2 := spec.RelayConfig.EVMChainID()
		if err2 != nil {
			return nil, errors.Wrap(err2, "ServicesForSpec failed to get chainID")
		}
		chain, err2 := d.chainSet.Get(big.NewInt(chainID))
		if err2 != nil {
			return nil, errors.Wrap(err2, "ServicesForSpec failed to get chainset")
		}

		if spec.PluginType != job.Mercury {
			if !common.IsHexAddress(transmitterID) {
				return nil, errors.Errorf("transmitterID is not valid EVM hex address, got: %v", transmitterID)
			}
			if spec.RelayConfig["sendingKeys"] == nil {
				spec.RelayConfig["sendingKeys"] = []string{transmitterID}
			}

			// effectiveTransmitterID is the transmitter address registered on the ocr contract. This is by default the EOA account on the node.
			// In the case of forwarding, the transmitter address is the forwarder contract deployed onchain between EOA and OCR contract.
			if jb.ForwardingAllowed { // FIXME: ForwardingAllowed cannot be set with Mercury, validate this
				fwdrAddress, fwderr := chain.TxManager().GetForwarderForEOA(common.HexToAddress(transmitterID))
				if fwderr == nil {
					effectiveTransmitterID = fwdrAddress.String()
				} else {
					lggr.Warnw("Skipping forwarding for job, will fallback to default behavior", "job", jb.Name, "err", fwderr)
				}
			}
		}
	}
	spec.RelayConfig["effectiveTransmitterID"] = effectiveTransmitterID

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

	switch spec.PluginType {
	case job.Mercury:
		mercuryProvider, err2 := relayer.NewMercuryProvider(
			types.RelayArgs{
				ExternalJobID: jb.ExternalJobID,
				JobID:         spec.ID,
				ContractID:    spec.ContractID,
				New:           d.isNewlyCreatedJob,
				RelayConfig:   spec.RelayConfig.Bytes(),
			}, types.PluginArgs{
				TransmitterID: transmitterID,
				PluginConfig:  spec.PluginConfig.Bytes(),
			})
		if err2 != nil {
			return nil, err2
		}
		oracleArgsNoPlugin := libocr2.OracleArgs{
			BinaryNetworkEndpointFactory: peerWrapper.Peer2,
			V2Bootstrappers:              bootstrapPeers,
			ContractTransmitter:          mercuryProvider.ContractTransmitter(),
			ContractConfigTracker:        mercuryProvider.ContractConfigTracker(),
			Database:                     ocrDB,
			LocalConfig:                  lc,
			Logger:                       ocrLogger,
			// FIXME: It looks like telemetry is uniquely keyed by contractID
			// but mercury runs multiple feeds per contract.
			// How can we scope this to a more granular level?
			// https://smartcontract-it.atlassian.net/browse/MERC-227
			MonitoringEndpoint:     d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID, synchronization.OCR2Mercury),
			OffchainConfigDigester: mercuryProvider.OffchainConfigDigester(),
			OffchainKeyring:        kb,
			OnchainKeyring:         kb,
		}
		return mercury.NewServices(jb, mercuryProvider, d.pipelineRunner, runResults, lggr, oracleArgsNoPlugin, d.cfg)
	case job.Median:
		medianProvider, err2 := relayer.NewMedianProvider(
			types.RelayArgs{
				ExternalJobID: jb.ExternalJobID,
				JobID:         spec.ID,
				ContractID:    spec.ContractID,
				New:           d.isNewlyCreatedJob,
				RelayConfig:   spec.RelayConfig.Bytes(),
			}, types.PluginArgs{
				TransmitterID: transmitterID,
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
		eaMonitoringEndpoint := d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID, synchronization.EnhancedEA)
		return median.NewMedianServices(jb, medianProvider, d.pipelineRunner, runResults, lggr, ocrLogger, oracleArgsNoPlugin, d.cfg, eaMonitoringEndpoint)
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
				TransmitterID: transmitterID,
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
				TransmitterID: transmitterID,
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
				TransmitterID: transmitterID,
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
			chain.GasEstimator(),
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
		functionsProvider, err2 := evmrelay.NewFunctionsProvider(
			d.chainSet,
			types.RelayArgs{
				ExternalJobID: jb.ExternalJobID,
				JobID:         spec.ID,
				ContractID:    spec.ContractID,
				RelayConfig:   spec.RelayConfig.Bytes(),
				New:           d.isNewlyCreatedJob,
			},
			types.PluginArgs{
				TransmitterID: transmitterID,
				PluginConfig:  spec.PluginConfig.Bytes(),
			},
			lggr.Named("FunctionsRelayer"),
			d.ethKs,
		)
		if err2 != nil {
			return nil, err2
		}

		var relayConfig evmrelaytypes.RelayConfig
		err2 = json.Unmarshal(spec.RelayConfig.Bytes(), &relayConfig)
		if err2 != nil {
			return nil, err2
		}
		chain, err2 := d.chainSet.Get(relayConfig.ChainID.ToInt())
		if err2 != nil {
			return nil, err2
		}

		sharedOracleArgs := libocr2.OracleArgs{
			BinaryNetworkEndpointFactory: peerWrapper.Peer2,
			V2Bootstrappers:              bootstrapPeers,
			ContractTransmitter:          functionsProvider.ContractTransmitter(),
			ContractConfigTracker:        functionsProvider.ContractConfigTracker(),
			Database:                     ocrDB,
			LocalConfig:                  lc,
			Logger:                       ocrLogger,
			MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID, synchronization.OCR2Functions),
			OffchainConfigDigester:       functionsProvider.OffchainConfigDigester(),
			OffchainKeyring:              kb,
			OnchainKeyring:               kb,
			ReportingPluginFactory:       nil, // To be set by NewFunctionsServices
		}

		functionsServicesConfig := functions.FunctionsServicesConfig{
			Job:            jb,
			PipelineRunner: d.pipelineRunner,
			JobORM:         d.jobORM,
			OCR2JobConfig:  d.cfg,
			DB:             d.db,
			Chain:          chain,
			ContractID:     spec.ContractID,
			Lggr:           lggr,
			MailMon:        d.mailMon,
		}

		functionsServices, err := functions.NewFunctionsServices(&sharedOracleArgs, &functionsServicesConfig)
		if err != nil {
			return nil, errors.Wrap(err, "error calling NewFunctionsServices")
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

		return append([]job.ServiceCtx{runResultSaver, functionsProvider}, functionsServices...), nil
	default:
		return nil, errors.Errorf("plugin type %s not supported", spec.PluginType)
	}
}
