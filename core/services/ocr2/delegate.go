package ocr2

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/commontypes"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocr2keepers20 "github.com/smartcontractkit/ocr2keepers/pkg/v2"
	ocr2keepers20config "github.com/smartcontractkit/ocr2keepers/pkg/v2/config"
	ocr2keepers20coordinator "github.com/smartcontractkit/ocr2keepers/pkg/v2/coordinator"
	ocr2keepers20polling "github.com/smartcontractkit/ocr2keepers/pkg/v2/observer/polling"
	ocr2keepers20runner "github.com/smartcontractkit/ocr2keepers/pkg/v2/runner"
	ocr2keepers21config "github.com/smartcontractkit/ocr2keepers/pkg/v3/config"
	ocr2keepers21 "github.com/smartcontractkit/ocr2keepers/pkg/v3/plugin"
	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	dkgpkg "github.com/smartcontractkit/ocr2vrf/dkg"
	"github.com/smartcontractkit/ocr2vrf/ocr2vrf"
	"github.com/smartcontractkit/sqlx"

	relaylogger "github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/dkg"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/dkg/persistence"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper"
	ocr2keeper21core "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
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
	functionsRelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/functions"
	evmmercury "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"

	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

var ErrJobSpecNoRelayer = errors.New("OCR2 job spec could not get relayer id")

type RelayGetter interface {
	Get(id relay.ID) (loop.Relayer, error)
}
type Delegate struct {
	db                    *sqlx.DB
	jobORM                job.ORM
	bridgeORM             bridges.ORM
	mercuryORM            evmmercury.ORM
	pipelineRunner        pipeline.Runner
	peerWrapper           *ocrcommon.SingletonPeerWrapper
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	cfg                   DelegateConfig
	lggr                  logger.Logger
	ks                    keystore.OCR2
	dkgSignKs             keystore.DKGSign
	dkgEncryptKs          keystore.DKGEncrypt
	ethKs                 keystore.Eth
	RelayGetter
	isNewlyCreatedJob bool // Set to true if this is a new job freshly added, false if job was present already on node boot.
	mailMon           *utils.MailboxMonitor

	legacyChains evm.LegacyChainContainer // legacy: use relayers instead
}

type DelegateConfig interface {
	plugins.RegistrarConfig
	OCR2() ocr2Config
	JobPipeline() jobPipelineConfig
	Database() pg.QConfig
	Insecure() insecureConfig
	Mercury() coreconfig.Mercury
	Threshold() coreconfig.Threshold
}

// concrete implementation of DelegateConfig so it can be explicitly composed
type delegateConfig struct {
	plugins.RegistrarConfig
	ocr2        ocr2Config
	jobPipeline jobPipelineConfig
	database    pg.QConfig
	insecure    insecureConfig
	mercury     mercuryConfig
	threshold   thresholdConfig
}

func (d *delegateConfig) JobPipeline() jobPipelineConfig {
	return d.jobPipeline
}

func (d *delegateConfig) Database() pg.QConfig {
	return d.database
}

func (d *delegateConfig) Insecure() insecureConfig {
	return d.insecure
}

func (d *delegateConfig) Threshold() coreconfig.Threshold {
	return d.threshold
}

func (d *delegateConfig) Mercury() coreconfig.Mercury {
	return d.mercury
}

func (d *delegateConfig) OCR2() ocr2Config {
	return d.ocr2
}

type ocr2Config interface {
	BlockchainTimeout() time.Duration
	CaptureEATelemetry() bool
	ContractConfirmations() uint16
	ContractPollInterval() time.Duration
	ContractTransmitterTransmitTimeout() time.Duration
	DatabaseTimeout() time.Duration
	KeyBundleID() (string, error)
	TraceLogging() bool
}

type insecureConfig interface {
	OCRDevelopmentMode() bool
}

type jobPipelineConfig interface {
	MaxSuccessfulRuns() uint64
	ResultWriteQueueDepth() uint64
}

type mercuryConfig interface {
	Credentials(credName string) *models.MercuryCredentials
}

type thresholdConfig interface {
	ThresholdKeyShare() string
}

func NewDelegateConfig(ocr2Cfg ocr2Config, m coreconfig.Mercury, t coreconfig.Threshold, i insecureConfig, jp jobPipelineConfig, qconf pg.QConfig, pluginProcessCfg plugins.RegistrarConfig) DelegateConfig {
	return &delegateConfig{
		ocr2:            ocr2Cfg,
		RegistrarConfig: pluginProcessCfg,
		jobPipeline:     jp,
		database:        qconf,
		insecure:        i,
		mercury:         m,
		threshold:       t,
	}
}

var _ job.Delegate = (*Delegate)(nil)

func NewDelegate(
	db *sqlx.DB,
	jobORM job.ORM,
	bridgeORM bridges.ORM,
	mercuryORM evmmercury.ORM,
	pipelineRunner pipeline.Runner,
	peerWrapper *ocrcommon.SingletonPeerWrapper,
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator,
	legacyChains evm.LegacyChainContainer,
	lggr logger.Logger,
	cfg DelegateConfig,
	ks keystore.OCR2,
	dkgSignKs keystore.DKGSign,
	dkgEncryptKs keystore.DKGEncrypt,
	ethKs keystore.Eth,
	relayers RelayGetter,
	mailMon *utils.MailboxMonitor,
	eventBroadcaster pg.EventBroadcaster,
) *Delegate {
	return &Delegate{
		db:                    db,
		jobORM:                jobORM,
		bridgeORM:             bridgeORM,
		mercuryORM:            mercuryORM,
		pipelineRunner:        pipelineRunner,
		peerWrapper:           peerWrapper,
		monitoringEndpointGen: monitoringEndpointGen,
		legacyChains:          legacyChains,
		cfg:                   cfg,
		lggr:                  lggr.Named("OCR2"),
		ks:                    ks,
		dkgSignKs:             dkgSignKs,
		dkgEncryptKs:          dkgEncryptKs,
		ethKs:                 ethKs,
		RelayGetter:           relayers,
		isNewlyCreatedJob:     false,
		mailMon:               mailMon,
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
	//  the job deletion itself isn't blocked.

	spec := jb.OCR2OracleSpec
	if spec == nil {
		d.lggr.Errorf("offchainreporting2.Delegate.OnDeleteJob called with wrong job type, ignoring non-OCR2 spec %v", jb)
		return nil
	}

	rid, err := spec.RelayID()
	if err != nil {
		d.lggr.Errorw("DeleteJob: "+ErrJobSpecNoRelayer.Error(), "err", err)
		return nil
	}
	// we only have clean to do for the EVM
	if rid.Network == relay.EVM {
		return d.cleanupEVM(jb, q, rid)
	}
	return nil
}

// cleanupEVM is a helper for clean up EVM specific state when a job is deleted
func (d *Delegate) cleanupEVM(jb job.Job, q pg.Queryer, relayID relay.ID) error {
	//  If UnregisterFilter returns an
	//  error, that means it failed to remove a valid active filter from the db.  We do abort the job deletion
	//  in that case, since it should be easy for the user to retry and will avoid leaving the db in
	//  an inconsistent state.  This assumes UnregisterFilter will return nil if the filter wasn't found
	//  at all (no rows deleted).
	spec := jb.OCR2OracleSpec
	chain, err := d.legacyChains.Get(relayID.ChainID)
	if err != nil {
		d.lggr.Error("cleanupEVM: failed to chain get chain %s", "err", relayID.ChainID, err)
		return nil
	}
	lp := chain.LogPoller()

	var filters []string
	switch spec.PluginType {
	case types.OCR2VRF:
		filters, err = ocr2coordinator.FilterNamesFromSpec(spec)
		if err != nil {
			d.lggr.Errorw("failed to derive ocr2vrf filter names from spec", "err", err, "spec", spec)
		}
	case types.OCR2Keeper:
		filters, err = ocr2keeper.FilterNamesFromSpec20(spec)
		if err != nil {
			d.lggr.Errorw("failed to derive ocr2keeper filter names from spec", "err", err, "spec", spec)
		}
	default:
		return nil
	}

	rargs := types.RelayArgs{
		ExternalJobID: jb.ExternalJobID,
		JobID:         jb.ID,
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
		err = lp.UnregisterFilter(filter, pg.WithQueryer(q))
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
		return nil, errors.Errorf("offchainreporting2.Delegate expects an *job.OCR2OracleSpec to be present, got %v", jb)
	}

	transmitterID := spec.TransmitterID.String
	effectiveTransmitterID := transmitterID

	lggrCtx := loop.ContextValues{
		JobID:   jb.ID,
		JobName: jb.Name.ValueOrZero(),

		ContractID:    spec.ContractID,
		TransmitterID: transmitterID,
	}
	if spec.FeedID != nil && (*spec.FeedID != (common.Hash{})) {
		lggrCtx.FeedID = *spec.FeedID
		spec.RelayConfig["feedID"] = spec.FeedID
	}
	lggr := logger.Sugared(d.lggr.Named(jb.ExternalJobID.String()).With(lggrCtx.Args()...))

	rid, err := spec.RelayID()
	if err != nil {
		return nil, fmt.Errorf("ServicesForSpec: %w: %w", ErrJobSpecNoRelayer, err)
	}

	if rid.Network == relay.EVM {
		lggr = logger.Sugared(lggr.With("evmChainID", rid.ChainID))

		chain, err2 := d.legacyChains.Get(rid.ChainID)
		if err2 != nil {
			return nil, fmt.Errorf("ServicesForSpec: could not get EVM chain %s: %w", rid.ChainID, err2)
		}
		effectiveTransmitterID, err2 = GetEVMEffectiveTransmitterID(&jb, chain, lggr)
		if err2 != nil {
			return nil, fmt.Errorf("ServicesForSpec failed to get evm transmitterID: %w", err2)
		}
	}
	spec.RelayConfig["effectiveTransmitterID"] = effectiveTransmitterID

	ocrDB := NewDB(d.db, spec.ID, 0, lggr, d.cfg.Database())
	if d.peerWrapper == nil {
		return nil, errors.New("cannot setup OCR2 job service, libp2p peer was missing")
	} else if !d.peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR2 jobs require a started and running p2p v2 peer")
	}

	ocrLogger := relaylogger.NewOCRWrapper(lggr, d.cfg.OCR2().TraceLogging(), func(msg string) {
		lggr.ErrorIf(d.jobORM.RecordError(jb.ID, msg), "unable to record error")
	})

	lc, err := validate.ToLocalConfig(d.cfg.OCR2(), d.cfg.Insecure(), *spec)
	if err != nil {
		return nil, err
	}
	if err = libocr2.SanityCheckLocalConfig(lc); err != nil {
		return nil, err
	}
	lggr.Infow("OCR2 job using local config",
		"BlockchainTimeout", lc.BlockchainTimeout,
		"ContractConfigConfirmations", lc.ContractConfigConfirmations,
		"ContractConfigTrackerPollInterval", lc.ContractConfigTrackerPollInterval,
		"ContractTransmitterTransmitTimeout", lc.ContractTransmitterTransmitTimeout,
		"DatabaseTimeout", lc.DatabaseTimeout,
	)

	bootstrapPeers, err := ocrcommon.GetValidatedBootstrapPeers(spec.P2PV2Bootstrappers, d.peerWrapper.P2PConfig().V2().DefaultBootstrappers())
	if err != nil {
		return nil, err
	}
	lggr.Debugw("Using bootstrap peers", "peers", bootstrapPeers)
	// Fetch the specified OCR2 key bundle
	var kbID string
	if spec.OCRKeyBundleID.Valid {
		kbID = spec.OCRKeyBundleID.String
	} else if kbID, err = d.cfg.OCR2().KeyBundleID(); err != nil {
		return nil, err
	}
	kb, err := d.ks.Get(kbID)
	if err != nil {
		return nil, err
	}

	spec.CaptureEATelemetry = d.cfg.OCR2().CaptureEATelemetry()

	runResults := make(chan *pipeline.Run, d.cfg.JobPipeline().ResultWriteQueueDepth())

	ctx := lggrCtx.ContextWithValues(context.Background())
	switch spec.PluginType {
	case types.Mercury:
		return d.newServicesMercury(ctx, lggr, jb, runResults, bootstrapPeers, kb, ocrDB, lc, ocrLogger)

	case types.Median:
		return d.newServicesMedian(ctx, lggr, jb, runResults, bootstrapPeers, kb, ocrDB, lc, ocrLogger)

	case types.DKG:
		return d.newServicesDKG(lggr, jb, bootstrapPeers, kb, ocrDB, lc, ocrLogger)

	case types.OCR2VRF:
		return d.newServicesOCR2VRF(lggr, jb, runResults, bootstrapPeers, kb, ocrDB, lc)

	case types.OCR2Keeper:
		return d.newServicesOCR2Keepers(lggr, jb, runResults, bootstrapPeers, kb, ocrDB, lc, ocrLogger)

	case types.Functions:
		const (
			_ int32 = iota
			thresholdPluginId
			s4PluginId
		)
		thresholdPluginDB := NewDB(d.db, spec.ID, thresholdPluginId, lggr, d.cfg.Database())
		s4PluginDB := NewDB(d.db, spec.ID, s4PluginId, lggr, d.cfg.Database())
		return d.newServicesOCR2Functions(lggr, jb, runResults, bootstrapPeers, kb, ocrDB, thresholdPluginDB, s4PluginDB, lc, ocrLogger)

	default:
		return nil, errors.Errorf("plugin type %s not supported", spec.PluginType)
	}
}

func GetEVMEffectiveTransmitterID(jb *job.Job, chain evm.Chain, lggr logger.SugaredLogger) (string, error) {
	spec := jb.OCR2OracleSpec
	if spec.PluginType == types.Mercury {
		return spec.TransmitterID.String, nil
	}

	if spec.RelayConfig["sendingKeys"] == nil {
		spec.RelayConfig["sendingKeys"] = []string{spec.TransmitterID.String}
	} else if !spec.TransmitterID.Valid {
		sendingKeys, err := job.SendingKeysForJob(jb)
		if err != nil {
			return "", err
		}
		if len(sendingKeys) > 1 && spec.PluginType != types.OCR2VRF {
			return "", errors.New("only ocr2 vrf should have more than 1 sending key")
		}
		spec.TransmitterID = null.StringFrom(sendingKeys[0])
	}

	// effectiveTransmitterID is the transmitter address registered on the ocr contract. This is by default the EOA account on the node.
	// In the case of forwarding, the transmitter address is the forwarder contract deployed onchain between EOA and OCR contract.
	// ForwardingAllowed cannot be set with Mercury, so this should always be false for mercury jobs
	if jb.ForwardingAllowed {
		if chain == nil {
			return "", fmt.Errorf("job forwarding requires non-nil chain")
		}
		effectiveTransmitterID, err := chain.TxManager().GetForwarderForEOA(common.HexToAddress(spec.TransmitterID.String))
		if err == nil {
			return effectiveTransmitterID.String(), nil
		} else if spec.TransmitterID.Valid {
			lggr.Warnw("Skipping forwarding for job, will fallback to default behavior", "job", jb.Name, "err", err)
			// this shouldn't happen unless behaviour above was changed
		} else {
			return "", errors.New("failed to get forwarder address and transmitterID is not set")
		}
	}

	return spec.TransmitterID.String, nil
}

func (d *Delegate) newServicesMercury(
	ctx context.Context,
	lggr logger.SugaredLogger,
	jb job.Job,
	runResults chan *pipeline.Run,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	ocrDB *db,
	lc ocrtypes.LocalConfig,
	ocrLogger commontypes.Logger,
) ([]job.ServiceCtx, error) {
	if jb.OCR2OracleSpec.FeedID == nil || (*jb.OCR2OracleSpec.FeedID == (common.Hash{})) {
		return nil, errors.Errorf("ServicesForSpec: mercury job type requires feedID")
	}
	spec := jb.OCR2OracleSpec
	transmitterID := spec.TransmitterID.String
	if len(transmitterID) != 64 {
		return nil, errors.Errorf("ServicesForSpec: mercury job type requires transmitter ID to be a 32-byte hex string, got: %q", transmitterID)
	}
	if _, err := hex.DecodeString(transmitterID); err != nil {
		return nil, errors.Wrapf(err, "ServicesForSpec: mercury job type requires transmitter ID to be a 32-byte hex string, got: %q", transmitterID)
	}

	rid, err := spec.RelayID()
	if err != nil {
		return nil, fmt.Errorf("mercury services: %w: %w", ErrJobSpecNoRelayer, err)
	}
	if rid.Network != relay.EVM {
		return nil, fmt.Errorf("mercury services: expected EVM relayer got %s", rid.Network)
	}
	relayer, err := d.RelayGetter.Get(rid)
	if err != nil {
		return nil, fmt.Errorf("failed to get relay %s is it enabled?: %w", spec.Relay, err)
	}
	chain, err := d.legacyChains.Get(rid.ChainID)
	if err != nil {
		return nil, fmt.Errorf("mercury services: failed to get chain %s: %w", rid.ChainID, err)
	}

	provider, err2 := relayer.NewPluginProvider(ctx,
		types.RelayArgs{
			ExternalJobID: jb.ExternalJobID,
			JobID:         jb.ID,
			ContractID:    spec.ContractID,
			New:           d.isNewlyCreatedJob,
			RelayConfig:   spec.RelayConfig.Bytes(),
			ProviderType:  string(spec.PluginType),
		}, types.PluginArgs{
			TransmitterID: transmitterID,
			PluginConfig:  spec.PluginConfig.Bytes(),
		})
	if err2 != nil {
		return nil, err2
	}

	mercuryProvider, ok := provider.(types.MercuryProvider)
	if !ok {
		return nil, errors.New("could not coerce PluginProvider to MercuryProvider")
	}

	oracleArgsNoPlugin := libocr2.MercuryOracleArgs{
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		ContractTransmitter:          mercuryProvider.ContractTransmitter(),
		ContractConfigTracker:        mercuryProvider.ContractConfigTracker(),
		Database:                     ocrDB,
		LocalConfig:                  lc,
		Logger:                       ocrLogger,
		MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.FeedID.String(), synchronization.OCR3Mercury, rid.Network, rid.ChainID),
		OffchainConfigDigester:       mercuryProvider.OffchainConfigDigester(),
		OffchainKeyring:              kb,
		OnchainKeyring:               kb,
	}

	chEnhancedTelem := make(chan ocrcommon.EnhancedTelemetryMercuryData, 100)

	mercuryServices, err2 := mercury.NewServices(jb, mercuryProvider, d.pipelineRunner, runResults, lggr, oracleArgsNoPlugin, d.cfg.JobPipeline(), chEnhancedTelem, chain, d.mercuryORM, (mercuryutils.FeedID)(*spec.FeedID))

	if ocrcommon.ShouldCollectEnhancedTelemetryMercury(&jb) {
		enhancedTelemService := ocrcommon.NewEnhancedTelemetryService(&jb, chEnhancedTelem, make(chan struct{}), d.monitoringEndpointGen.GenMonitoringEndpoint(spec.FeedID.String(), synchronization.EnhancedEAMercury, rid.Network, rid.ChainID), lggr.Named("EnhancedTelemetryMercury"))
		mercuryServices = append(mercuryServices, enhancedTelemService)
	}

	return mercuryServices, err2
}

func (d *Delegate) newServicesMedian(
	ctx context.Context,
	lggr logger.SugaredLogger,
	jb job.Job,
	runResults chan *pipeline.Run,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	ocrDB *db,
	lc ocrtypes.LocalConfig,
	ocrLogger commontypes.Logger,
) ([]job.ServiceCtx, error) {
	spec := jb.OCR2OracleSpec

	rid, err := spec.RelayID()
	if err != nil {
		return nil, fmt.Errorf("median services: %w: %w", ErrJobSpecNoRelayer, err)
	}

	oracleArgsNoPlugin := libocr2.OCR2OracleArgs{
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		Database:                     ocrDB,
		LocalConfig:                  lc,
		Logger:                       ocrLogger,
		MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID, synchronization.OCR2Median, rid.Network, rid.ChainID),
		OffchainKeyring:              kb,
		OnchainKeyring:               kb,
	}
	errorLog := &errorLog{jobID: jb.ID, recordError: d.jobORM.RecordError}
	enhancedTelemChan := make(chan ocrcommon.EnhancedTelemetryData, 100)
	mConfig := median.NewMedianConfig(d.cfg.JobPipeline().MaxSuccessfulRuns(), d.cfg)

	relayer, err := d.RelayGetter.Get(rid)
	if err != nil {
		return nil, fmt.Errorf("median services; failed to get relay %s is it enabled?: %w", spec.Relay, err)
	}

	medianServices, err2 := median.NewMedianServices(ctx, jb, d.isNewlyCreatedJob, relayer, d.pipelineRunner, runResults, lggr, oracleArgsNoPlugin, mConfig, enhancedTelemChan, errorLog)

	if ocrcommon.ShouldCollectEnhancedTelemetry(&jb) {
		enhancedTelemService := ocrcommon.NewEnhancedTelemetryService(&jb, enhancedTelemChan, make(chan struct{}), d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID, synchronization.EnhancedEA, rid.Network, rid.ChainID), lggr.Named("EnhancedTelemetry"))
		medianServices = append(medianServices, enhancedTelemService)
	}

	return medianServices, err2
}

func (d *Delegate) newServicesDKG(
	lggr logger.SugaredLogger,
	jb job.Job,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	ocrDB *db,
	lc ocrtypes.LocalConfig,
	ocrLogger commontypes.Logger,
) ([]job.ServiceCtx, error) {
	spec := jb.OCR2OracleSpec
	rid, err := spec.RelayID()
	if err != nil {
		return nil, fmt.Errorf("DKG services: %w: %w", ErrJobSpecNoRelayer, err)
	}
	if rid.Network != relay.EVM {
		return nil, fmt.Errorf("DKG services: expected EVM relayer got %s", rid.Network)
	}

	chain, err2 := d.legacyChains.Get(rid.ChainID)
	if err2 != nil {
		return nil, fmt.Errorf("DKG services: failed to get chain %s: %w", rid.ChainID, err2)
	}
	ocr2vrfRelayer := evmrelay.NewOCR2VRFRelayer(d.db, chain, lggr.Named("OCR2VRFRelayer"), d.ethKs)
	dkgProvider, err2 := ocr2vrfRelayer.NewDKGProvider(
		types.RelayArgs{
			ExternalJobID: jb.ExternalJobID,
			JobID:         jb.ID,
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
	oracleArgsNoPlugin := libocr2.OCR2OracleArgs{
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
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
		d.cfg.Database(),
		chain.ID(),
		spec.Relay,
	)
}

func (d *Delegate) newServicesOCR2VRF(
	lggr logger.SugaredLogger,
	jb job.Job,
	runResults chan *pipeline.Run,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	ocrDB *db,
	lc ocrtypes.LocalConfig,
) ([]job.ServiceCtx, error) {
	spec := jb.OCR2OracleSpec

	rid, err := spec.RelayID()
	if err != nil {
		return nil, fmt.Errorf("VRF services: %w: %w", ErrJobSpecNoRelayer, err)
	}
	if rid.Network != relay.EVM {
		return nil, fmt.Errorf("VRF services: expected EVM relayer got %s", rid.Network)
	}
	chain, err2 := d.legacyChains.Get(rid.ChainID)
	if err2 != nil {
		return nil, fmt.Errorf("VRF services: failed to get chain (%s): %w", rid.ChainID, err2)
	}
	if jb.ForwardingAllowed != chain.Config().EVM().Transactions().ForwardersEnabled() {
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
	transmitterID := spec.TransmitterID.String

	vrfProvider, err2 := ocr2vrfRelayer.NewOCR2VRFProvider(
		types.RelayArgs{
			ExternalJobID: jb.ExternalJobID,
			JobID:         jb.ID,
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
			JobID:         jb.ID,
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

	timeout := 5 * time.Second
	interval := 60 * time.Second
	juelsLogger := lggr.Named("JuelsFeeCoin").With("contract", cfg.LinkEthFeedAddress, "timeout", timeout, "interval", interval)
	juelsPerFeeCoin, err2 := juelsfeecoin.NewLinkEthPriceProvider(
		common.HexToAddress(cfg.LinkEthFeedAddress), chain.Client(), timeout, interval, juelsLogger)
	if err2 != nil {
		return nil, errors.Wrap(err2, "new link eth price provider")
	}

	reasonableGasPrice := reasonablegasprice.NewReasonableGasPriceProvider(
		chain.GasEstimator(),
		timeout,
		chain.Config().EVM().GasEstimator().PriceMax(),
		chain.Config().EVM().GasEstimator().EIP1559DynamicFees(),
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
		chain.Config().EVM().FinalityDepth(),
	)
	if err2 != nil {
		return nil, errors.Wrap(err2, "create ocr2vrf coordinator")
	}
	l := lggr.Named("OCR2VRF").With(
		"jobName", jb.Name.ValueOrZero(),
		"jobID", jb.ID,
	)
	vrfLogger := relaylogger.NewOCRWrapper(l.With(
		"vrfContractID", spec.ContractID), d.cfg.OCR2().TraceLogging(), func(msg string) {
		lggr.ErrorIf(d.jobORM.RecordError(jb.ID, msg), "unable to record error")
	})
	dkgLogger := relaylogger.NewOCRWrapper(l.With(
		"dkgContractID", cfg.DKGContractAddress), d.cfg.OCR2().TraceLogging(), func(msg string) {
		lggr.ErrorIf(d.jobORM.RecordError(jb.ID, msg), "unable to record error")
	})
	dkgReportingPluginFactoryDecorator := func(wrapped ocrtypes.ReportingPluginFactory) ocrtypes.ReportingPluginFactory {
		return promwrapper.NewPromFactory(wrapped, "DKG", string(relay.EVM), chain.ID())
	}
	vrfReportingPluginFactoryDecorator := func(wrapped ocrtypes.ReportingPluginFactory) ocrtypes.ReportingPluginFactory {
		return promwrapper.NewPromFactory(wrapped, "OCR2VRF", string(relay.EVM), chain.ID())
	}
	noopMonitoringEndpoint := telemetry.NoopAgent{}
	oracles, err2 := ocr2vrf.NewOCR2VRF(ocr2vrf.DKGVRFArgs{
		VRFLogger:                    vrfLogger,
		DKGLogger:                    dkgLogger,
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		OffchainKeyring:              kb,
		OnchainKeyring:               kb,
		VRFOffchainConfigDigester:    vrfProvider.OffchainConfigDigester(),
		VRFContractConfigTracker:     vrfProvider.ContractConfigTracker(),
		VRFContractTransmitter:       vrfProvider.ContractTransmitter(),
		VRFDatabase:                  ocrDB,
		VRFLocalConfig:               lc,
		VRFMonitoringEndpoint:        d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID, synchronization.OCR2VRF, rid.Network, rid.ChainID),
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
		DKGSharePersistence:                persistence.NewShareDB(d.db, lggr.Named("DKGShareDB"), d.cfg.Database(), chain.ID(), spec.Relay),
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
		d.cfg.JobPipeline().MaxSuccessfulRuns(),
	)

	// NOTE: we return from here with the services because the OCR2VRF oracles are defined
	// and exported from the ocr2vrf library. It takes care of running the DKG and OCR2VRF
	// oracles under the hood together.
	oracleCtx := job.NewServiceAdapter(oracles)
	return []job.ServiceCtx{runResultSaver, vrfProvider, dkgProvider, oracleCtx}, nil
}

func (d *Delegate) newServicesOCR2Keepers(
	lggr logger.SugaredLogger,
	jb job.Job,
	runResults chan *pipeline.Run,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	ocrDB *db,
	lc ocrtypes.LocalConfig,
	ocrLogger commontypes.Logger,
) ([]job.ServiceCtx, error) {
	spec := jb.OCR2OracleSpec
	var cfg ocr2keeper.PluginConfig
	if err := json.Unmarshal(spec.PluginConfig.Bytes(), &cfg); err != nil {
		return nil, errors.Wrap(err, "unmarshal ocr2keepers plugin config")
	}

	if err := ocr2keeper.ValidatePluginConfig(cfg); err != nil {
		return nil, errors.Wrap(err, "ocr2keepers plugin config validation failure")
	}

	switch cfg.ContractVersion {
	case "v2.1":
		return d.newServicesOCR2Keepers21(lggr, jb, runResults, bootstrapPeers, kb, ocrDB, lc, ocrLogger, cfg, spec)
	case "v2.0":
		return d.newServicesOCR2Keepers20(lggr, jb, runResults, bootstrapPeers, kb, ocrDB, lc, ocrLogger, cfg, spec)
	default:
		return d.newServicesOCR2Keepers20(lggr, jb, runResults, bootstrapPeers, kb, ocrDB, lc, ocrLogger, cfg, spec)
	}
}

func (d *Delegate) newServicesOCR2Keepers21(
	lggr logger.SugaredLogger,
	jb job.Job,
	runResults chan *pipeline.Run,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	ocrDB *db,
	lc ocrtypes.LocalConfig,
	ocrLogger commontypes.Logger,
	cfg ocr2keeper.PluginConfig,
	spec *job.OCR2OracleSpec,
) ([]job.ServiceCtx, error) {
	credName, err2 := jb.OCR2OracleSpec.PluginConfig.MercuryCredentialName()
	if err2 != nil {
		return nil, errors.Wrap(err2, "failed to get mercury credential name")
	}

	mc := d.cfg.Mercury().Credentials(credName)
	rid, err := spec.RelayID()
	if err != nil {
		return nil, fmt.Errorf("keeper2 services: %w: %w", ErrJobSpecNoRelayer, err)
	}
	if rid.Network != relay.EVM {
		return nil, fmt.Errorf("keeper2 services: expected EVM relayer got %s", rid.Network)
	}

	chain, err2 := d.legacyChains.Get(rid.ChainID)
	if err2 != nil {
		return nil, fmt.Errorf("keeper2 services: failed to get chain %s: %w", rid.ChainID, err2)
	}

	keeperProvider, services, err2 := ocr2keeper.EVMDependencies21(jb, d.db, lggr, chain, d.pipelineRunner, mc, kb, d.cfg.Database())
	if err2 != nil {
		return nil, errors.Wrap(err2, "could not build dependencies for ocr2 keepers")
	}
	// set some defaults
	conf := ocr2keepers21config.ReportingFactoryConfig{
		CacheExpiration:       ocr2keepers21config.DefaultCacheExpiration,
		CacheEvictionInterval: ocr2keepers21config.DefaultCacheClearInterval,
		MaxServiceWorkers:     ocr2keepers21config.DefaultMaxServiceWorkers,
		ServiceQueueLength:    ocr2keepers21config.DefaultServiceQueueLength,
	}

	// override if set in config
	if cfg.CacheExpiration.Value() != 0 {
		conf.CacheExpiration = cfg.CacheExpiration.Value()
	}

	if cfg.CacheEvictionInterval.Value() != 0 {
		conf.CacheEvictionInterval = cfg.CacheEvictionInterval.Value()
	}

	if cfg.MaxServiceWorkers != 0 {
		conf.MaxServiceWorkers = cfg.MaxServiceWorkers
	}

	if cfg.ServiceQueueLength != 0 {
		conf.ServiceQueueLength = cfg.ServiceQueueLength
	}

	dConf := ocr2keepers21.DelegateConfig{
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		ContractTransmitter:          evmrelay.NewKeepersOCR3ContractTransmitter(keeperProvider.ContractTransmitter()),
		ContractConfigTracker:        keeperProvider.ContractConfigTracker(),
		KeepersDatabase:              ocrDB,
		Logger:                       ocrLogger,
		MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID, synchronization.OCR2Automation, rid.Network, rid.ChainID),
		OffchainConfigDigester:       keeperProvider.OffchainConfigDigester(),
		OffchainKeyring:              kb,
		OnchainKeyring:               services.Keyring(),
		LocalConfig:                  lc,
		LogProvider:                  services.LogEventProvider(),
		EventProvider:                services.TransmitEventProvider(),
		Runnable:                     services.Registry(),
		Encoder:                      services.Encoder(),
		BlockSubscriber:              services.BlockSubscriber(),
		RecoverableProvider:          services.LogRecoverer(),
		PayloadBuilder:               services.PayloadBuilder(),
		UpkeepProvider:               services.UpkeepProvider(),
		UpkeepStateUpdater:           services.UpkeepStateStore(),
		UpkeepTypeGetter:             ocr2keeper21core.GetUpkeepType,
		WorkIDGenerator:              ocr2keeper21core.UpkeepWorkID,
		// TODO: Clean up the config
		CacheExpiration:       cfg.CacheExpiration.Value(),
		CacheEvictionInterval: cfg.CacheEvictionInterval.Value(),
		MaxServiceWorkers:     cfg.MaxServiceWorkers,
		ServiceQueueLength:    cfg.ServiceQueueLength,
	}

	pluginService, err := ocr2keepers21.NewDelegate(dConf)
	if err != nil {
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
		d.cfg.JobPipeline().MaxSuccessfulRuns(),
	)

	return []job.ServiceCtx{
		runResultSaver,
		keeperProvider,
		services.Registry(),
		services.BlockSubscriber(),
		services.LogEventProvider(),
		services.LogRecoverer(),
		services.UpkeepStateStore(),
		services.TransmitEventProvider(),
		pluginService,
	}, nil
}

func (d *Delegate) newServicesOCR2Keepers20(
	lggr logger.SugaredLogger,
	jb job.Job,
	runResults chan *pipeline.Run,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	ocrDB *db,
	lc ocrtypes.LocalConfig,
	ocrLogger commontypes.Logger,
	cfg ocr2keeper.PluginConfig,
	spec *job.OCR2OracleSpec,
) ([]job.ServiceCtx, error) {

	rid, err := spec.RelayID()
	if err != nil {
		return nil, fmt.Errorf("keepers2.0 services: %w: %w", ErrJobSpecNoRelayer, err)
	}
	if rid.Network != relay.EVM {
		return nil, fmt.Errorf("keepers2.0 services: expected EVM relayer got %s", rid.Network)
	}
	chain, err2 := d.legacyChains.Get(rid.ChainID)
	if err2 != nil {
		return nil, fmt.Errorf("keepers2.0 services: failed to get chain (%s): %w", rid.ChainID, err2)
	}

	keeperProvider, rgstry, encoder, logProvider, err2 := ocr2keeper.EVMDependencies20(jb, d.db, lggr, chain, d.pipelineRunner)
	if err2 != nil {
		return nil, errors.Wrap(err2, "could not build dependencies for ocr2 keepers")
	}

	w := &logWriter{log: lggr.Named("Automation Dependencies")}

	// set some defaults
	conf := ocr2keepers20config.ReportingFactoryConfig{
		CacheExpiration:       ocr2keepers20config.DefaultCacheExpiration,
		CacheEvictionInterval: ocr2keepers20config.DefaultCacheClearInterval,
		MaxServiceWorkers:     ocr2keepers20config.DefaultMaxServiceWorkers,
		ServiceQueueLength:    ocr2keepers20config.DefaultServiceQueueLength,
	}

	// override if set in config
	if cfg.CacheExpiration.Value() != 0 {
		conf.CacheExpiration = cfg.CacheExpiration.Value()
	}

	if cfg.CacheEvictionInterval.Value() != 0 {
		conf.CacheEvictionInterval = cfg.CacheEvictionInterval.Value()
	}

	if cfg.MaxServiceWorkers != 0 {
		conf.MaxServiceWorkers = cfg.MaxServiceWorkers
	}

	if cfg.ServiceQueueLength != 0 {
		conf.ServiceQueueLength = cfg.ServiceQueueLength
	}

	runr, err2 := ocr2keepers20runner.NewRunner(
		log.New(w, "[automation-plugin-runner] ", log.Lshortfile),
		rgstry,
		encoder,
		conf.MaxServiceWorkers,
		conf.ServiceQueueLength,
		conf.CacheExpiration,
		conf.CacheEvictionInterval,
	)
	if err2 != nil {
		return nil, errors.Wrap(err2, "failed to create automation pipeline runner")
	}

	condObs := &ocr2keepers20polling.PollingObserverFactory{
		Logger:  log.New(w, "[automation-plugin-conditional-observer] ", log.Lshortfile),
		Source:  rgstry,
		Heads:   rgstry,
		Runner:  runr,
		Encoder: encoder,
	}

	coord := &ocr2keepers20coordinator.CoordinatorFactory{
		Logger:     log.New(w, "[automation-plugin-coordinator] ", log.Lshortfile),
		Encoder:    encoder,
		Logs:       logProvider,
		CacheClean: conf.CacheEvictionInterval,
	}

	dConf := ocr2keepers20.DelegateConfig{
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		ContractTransmitter:          keeperProvider.ContractTransmitter(),
		ContractConfigTracker:        keeperProvider.ContractConfigTracker(),
		KeepersDatabase:              ocrDB,
		LocalConfig:                  lc,
		Logger:                       ocrLogger,
		MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID, synchronization.OCR2Automation, rid.Network, rid.ChainID),
		OffchainConfigDigester:       keeperProvider.OffchainConfigDigester(),
		OffchainKeyring:              kb,
		OnchainKeyring:               kb,
		ConditionalObserverFactory:   condObs,
		CoordinatorFactory:           coord,
		Encoder:                      encoder,
		Runner:                       runr,
		// the following values are not needed in the delegate config anymore
		CacheExpiration:       cfg.CacheExpiration.Value(),
		CacheEvictionInterval: cfg.CacheEvictionInterval.Value(),
		MaxServiceWorkers:     cfg.MaxServiceWorkers,
		ServiceQueueLength:    cfg.ServiceQueueLength,
	}

	pluginService, err := ocr2keepers20.NewDelegate(dConf)
	if err != nil {
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
		d.cfg.JobPipeline().MaxSuccessfulRuns(),
	)

	return []job.ServiceCtx{
		job.NewServiceAdapter(runr),
		runResultSaver,
		keeperProvider,
		rgstry,
		logProvider,
		pluginService,
	}, nil
}

func (d *Delegate) newServicesOCR2Functions(
	lggr logger.SugaredLogger,
	jb job.Job,
	runResults chan *pipeline.Run,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	functionsOcrDB *db,
	thresholdOcrDB *db,
	s4OcrDB *db,
	lc ocrtypes.LocalConfig,
	ocrLogger commontypes.Logger,
) ([]job.ServiceCtx, error) {
	spec := jb.OCR2OracleSpec

	rid, err := spec.RelayID()
	if err != nil {
		return nil, fmt.Errorf("functions services: %w: %w", ErrJobSpecNoRelayer, err)
	}
	if rid.Network != relay.EVM {
		return nil, fmt.Errorf("functions services: expected EVM relayer got %s", rid.Network)
	}
	chain, err := d.legacyChains.Get(rid.ChainID)
	if err != nil {
		return nil, fmt.Errorf("functions services: failed to get chain %s: %w", rid.ChainID, err)
	}
	createPluginProvider := func(pluginType functionsRelay.FunctionsPluginType, relayerName string) (evmrelaytypes.FunctionsProvider, error) {
		return evmrelay.NewFunctionsProvider(
			chain,
			types.RelayArgs{
				ExternalJobID: jb.ExternalJobID,
				JobID:         jb.ID,
				ContractID:    spec.ContractID,
				RelayConfig:   spec.RelayConfig.Bytes(),
				New:           d.isNewlyCreatedJob,
			},
			types.PluginArgs{
				TransmitterID: spec.TransmitterID.String,
				PluginConfig:  spec.PluginConfig.Bytes(),
			},
			lggr.Named(relayerName),
			d.ethKs,
			pluginType,
		)
	}

	functionsProvider, err := createPluginProvider(functionsRelay.FunctionsPlugin, "FunctionsRelayer")
	if err != nil {
		return nil, err
	}

	thresholdProvider, err := createPluginProvider(functionsRelay.ThresholdPlugin, "FunctionsThresholdRelayer")
	if err != nil {
		return nil, err
	}

	s4Provider, err := createPluginProvider(functionsRelay.S4Plugin, "FunctionsS4Relayer")
	if err != nil {
		return nil, err
	}

	functionsOracleArgs := libocr2.OCR2OracleArgs{
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		ContractTransmitter:          functionsProvider.ContractTransmitter(),
		ContractConfigTracker:        functionsProvider.ContractConfigTracker(),
		Database:                     functionsOcrDB,
		LocalConfig:                  lc,
		Logger:                       ocrLogger,
		MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID, synchronization.OCR2Functions, rid.Network, rid.ChainID),
		OffchainConfigDigester:       functionsProvider.OffchainConfigDigester(),
		OffchainKeyring:              kb,
		OnchainKeyring:               kb,
		ReportingPluginFactory:       nil, // To be set by NewFunctionsServices
	}

	noopMonitoringEndpoint := telemetry.NoopAgent{}

	thresholdOracleArgs := libocr2.OCR2OracleArgs{
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		ContractTransmitter:          thresholdProvider.ContractTransmitter(),
		ContractConfigTracker:        thresholdProvider.ContractConfigTracker(),
		Database:                     thresholdOcrDB,
		LocalConfig:                  lc,
		Logger:                       ocrLogger,
		// Telemetry ingress for OCR2Threshold is currently not supported so a noop monitoring endpoint is being used
		MonitoringEndpoint:     &noopMonitoringEndpoint,
		OffchainConfigDigester: thresholdProvider.OffchainConfigDigester(),
		OffchainKeyring:        kb,
		OnchainKeyring:         kb,
		ReportingPluginFactory: nil, // To be set by NewFunctionsServices
	}

	s4OracleArgs := libocr2.OCR2OracleArgs{
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		ContractTransmitter:          s4Provider.ContractTransmitter(),
		ContractConfigTracker:        s4Provider.ContractConfigTracker(),
		Database:                     s4OcrDB,
		LocalConfig:                  lc,
		Logger:                       ocrLogger,
		// Telemetry ingress for OCR2S4 is currently not supported so a noop monitoring endpoint is being used
		MonitoringEndpoint:     &noopMonitoringEndpoint,
		OffchainConfigDigester: s4Provider.OffchainConfigDigester(),
		OffchainKeyring:        kb,
		OnchainKeyring:         kb,
		ReportingPluginFactory: nil, // To be set by NewFunctionsServices
	}

	encryptedThresholdKeyShare := d.cfg.Threshold().ThresholdKeyShare()
	var thresholdKeyShare []byte
	if len(encryptedThresholdKeyShare) > 0 {
		encryptedThresholdKeyShareBytes, err2 := hex.DecodeString(encryptedThresholdKeyShare)
		if err2 != nil {
			return nil, errors.Wrap(err2, "failed to decode ThresholdKeyShare hex string")
		}
		thresholdKeyShare, err2 = kb.NaclBoxOpenAnonymous(encryptedThresholdKeyShareBytes)
		if err2 != nil {
			return nil, errors.Wrap(err2, "failed to decrypt ThresholdKeyShare")
		}
	}

	functionsServicesConfig := functions.FunctionsServicesConfig{
		Job:               jb,
		JobORM:            d.jobORM,
		BridgeORM:         d.bridgeORM,
		QConfig:           d.cfg.Database(),
		DB:                d.db,
		Chain:             chain,
		ContractID:        spec.ContractID,
		Logger:            lggr,
		MailMon:           d.mailMon,
		URLsMonEndpoint:   d.monitoringEndpointGen.GenMonitoringEndpoint(spec.ContractID, synchronization.FunctionsRequests, rid.Network, rid.ChainID),
		EthKeystore:       d.ethKs,
		ThresholdKeyShare: thresholdKeyShare,
		LogPollerWrapper:  functionsProvider.LogPollerWrapper(),
	}

	functionsServices, err := functions.NewFunctionsServices(&functionsOracleArgs, &thresholdOracleArgs, &s4OracleArgs, &functionsServicesConfig)
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
		d.cfg.JobPipeline().MaxSuccessfulRuns(),
	)

	return append([]job.ServiceCtx{runResultSaver, functionsProvider, thresholdProvider, s4Provider}, functionsServices...), nil
}

// errorLog implements [loop.ErrorLog]
type errorLog struct {
	jobID       int32
	recordError func(jobID int32, description string, qopts ...pg.QOpt) error
}

func (l *errorLog) SaveError(ctx context.Context, msg string) error {
	return l.recordError(l.jobID, msg)
}

type logWriter struct {
	log logger.Logger
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	l.log.Debug(string(p), nil)
	n = len(p)
	return
}
