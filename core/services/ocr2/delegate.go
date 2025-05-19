package ocr2

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"gopkg.in/guregu/null.v4"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/commontypes"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ocr2keepers20 "github.com/smartcontractkit/chainlink-automation/pkg/v2"
	ocr2keepers20config "github.com/smartcontractkit/chainlink-automation/pkg/v2/config"
	ocr2keepers20coordinator "github.com/smartcontractkit/chainlink-automation/pkg/v2/coordinator"
	ocr2keepers20polling "github.com/smartcontractkit/chainlink-automation/pkg/v2/observer/polling"
	ocr2keepers20runner "github.com/smartcontractkit/chainlink-automation/pkg/v2/runner"
	ocr2keepers21config "github.com/smartcontractkit/chainlink-automation/pkg/v3/config"
	ocr2keepers21 "github.com/smartcontractkit/chainlink-automation/pkg/v3/plugin"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins/ocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/retirement"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/ccipcommit"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/ccipexec"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/generic"
	lloconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/llo/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/autotelemetry21"
	ocr2keeper21core "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	functionsRelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/functions"
	evmmercury "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type ErrJobSpecNoRelayer struct {
	PluginName string
	Err        error
}

func (e ErrJobSpecNoRelayer) Unwrap() error { return e.Err }

func (e ErrJobSpecNoRelayer) Error() string {
	return fmt.Sprintf("%s services: OCR2 job spec could not get relayer ID: %s", e.PluginName, e.Err)
}

type ErrRelayNotEnabled struct {
	PluginName string
	Relay      string
	Err        error
}

func (e ErrRelayNotEnabled) Unwrap() error { return e.Err }

func (e ErrRelayNotEnabled) Error() string {
	return fmt.Sprintf("%s services: failed to get relay %s, is it enabled? %s", e.PluginName, e.Relay, e.Err)
}

type RelayGetter interface {
	Get(id types.RelayID) (loop.Relayer, error)
	GetIDToRelayerMap() map[types.RelayID]loop.Relayer
}
type Delegate struct {
	ds                    sqlutil.DataSource
	jobORM                job.ORM
	bridgeORM             bridges.ORM
	mercuryORM            evmmercury.ORM
	pipelineRunner        pipeline.Runner
	streamRegistry        streams.Getter
	peerWrapper           *ocrcommon.SingletonPeerWrapper
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	cfg                   DelegateConfig
	lggr                  logger.Logger
	ks                    keystore.OCR2
	ethKs                 keystore.Eth
	RelayGetter
	isNewlyCreatedJob     bool // Set to true if this is a new job freshly added, false if job was present already on node boot.
	mailMon               *mailbox.Monitor
	retirementReportCache retirement.RetirementReportCache

	legacyChains         legacyevm.LegacyChainContainer // legacy: use relayers instead
	capabilitiesRegistry core.CapabilitiesRegistry
}

type DelegateConfig interface {
	plugins.RegistrarConfig
	OCR2() ocr2Config
	JobPipeline() jobPipelineConfig
	Insecure() insecureConfig
	Mercury() coreconfig.Mercury
	Threshold() coreconfig.Threshold
}

// concrete implementation of DelegateConfig so it can be explicitly composed
type delegateConfig struct {
	plugins.RegistrarConfig
	ocr2        ocr2Config
	jobPipeline jobPipelineConfig
	insecure    insecureConfig
	mercury     mercuryConfig
	threshold   thresholdConfig
}

func (d *delegateConfig) JobPipeline() jobPipelineConfig {
	return d.jobPipeline
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
	DefaultTransactionQueueDepth() uint32
	KeyBundleID() (string, error)
	SimulateTransactions() bool
	TraceLogging() bool
	CaptureAutomationCustomTelemetry() bool
	AllowNoBootstrappers() bool
}

type insecureConfig interface {
	OCRDevelopmentMode() bool
}

type jobPipelineConfig interface {
	MaxSuccessfulRuns() uint64
	ResultWriteQueueDepth() uint64
}

type mercuryConfig interface {
	Credentials(credName string) *types.MercuryCredentials
	Cache() coreconfig.MercuryCache
	TLS() coreconfig.MercuryTLS
	Transmitter() coreconfig.MercuryTransmitter
	VerboseLogging() bool
}

type thresholdConfig interface {
	ThresholdKeyShare() string
}

func NewDelegateConfig(ocr2Cfg ocr2Config, m coreconfig.Mercury, t coreconfig.Threshold, i insecureConfig, jp jobPipelineConfig, pluginProcessCfg plugins.RegistrarConfig) DelegateConfig {
	return &delegateConfig{
		ocr2:            ocr2Cfg,
		RegistrarConfig: pluginProcessCfg,
		jobPipeline:     jp,
		insecure:        i,
		mercury:         m,
		threshold:       t,
	}
}

var _ job.Delegate = (*Delegate)(nil)

type DelegateOpts struct {
	Ds                    sqlutil.DataSource
	JobORM                job.ORM
	BridgeORM             bridges.ORM
	MercuryORM            evmmercury.ORM
	PipelineRunner        pipeline.Runner
	StreamRegistry        streams.Getter
	PeerWrapper           *ocrcommon.SingletonPeerWrapper
	MonitoringEndpointGen telemetry.MonitoringEndpointGenerator
	LegacyChains          legacyevm.LegacyChainContainer
	Lggr                  logger.Logger
	Ks                    keystore.OCR2
	EthKs                 keystore.Eth
	Relayers              RelayGetter
	MailMon               *mailbox.Monitor
	CapabilitiesRegistry  core.CapabilitiesRegistry
	RetirementReportCache retirement.RetirementReportCache
}

func NewDelegate(
	opts DelegateOpts,
	cfg DelegateConfig,
) *Delegate {
	return &Delegate{
		ds:                    opts.Ds,
		jobORM:                opts.JobORM,
		bridgeORM:             opts.BridgeORM,
		mercuryORM:            opts.MercuryORM,
		pipelineRunner:        opts.PipelineRunner,
		streamRegistry:        opts.StreamRegistry,
		peerWrapper:           opts.PeerWrapper,
		monitoringEndpointGen: opts.MonitoringEndpointGen,
		legacyChains:          opts.LegacyChains,
		cfg:                   cfg,
		lggr:                  opts.Lggr.Named("OCR2"),
		ks:                    opts.Ks,
		ethKs:                 opts.EthKs,
		RelayGetter:           opts.Relayers,
		isNewlyCreatedJob:     false,
		mailMon:               opts.MailMon,
		capabilitiesRegistry:  opts.CapabilitiesRegistry,
		retirementReportCache: opts.RetirementReportCache,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.OffchainReporting2
}

func (d *Delegate) BeforeJobCreated(_ job.Job) {
	// This is only called first time the job is created
	d.isNewlyCreatedJob = true
}
func (d *Delegate) AfterJobCreated(_ job.Job)  {}
func (d *Delegate) BeforeJobDeleted(_ job.Job) {}
func (d *Delegate) OnDeleteJob(ctx context.Context, jb job.Job) error {
	// If the job spec is malformed in any way, we report the error but return nil so that
	//  the job deletion itself isn't blocked.

	spec := jb.OCR2OracleSpec
	if spec == nil {
		d.lggr.Errorf("offchainreporting2.Delegate.OnDeleteJob called with wrong job type, ignoring non-OCR2 spec %v", jb)
		return nil
	}

	rid, err := spec.RelayID()
	if err != nil {
		d.lggr.Errorw("DeleteJob", "err", ErrJobSpecNoRelayer{Err: err, PluginName: string(spec.PluginType)})
		return nil
	}
	// we only have clean to do for the EVM
	if rid.Network == relay.NetworkEVM {
		return d.cleanupEVM(ctx, jb, rid)
	}
	return nil
}

// cleanupEVM is a helper for clean up EVM specific state when a job is deleted
func (d *Delegate) cleanupEVM(ctx context.Context, jb job.Job, relayID types.RelayID) error {
	//  If UnregisterFilter returns an
	//  error, that means it failed to remove a valid active filter from the db.  We do abort the job deletion
	//  in that case, since it should be easy for the user to retry and will avoid leaving the db in
	//  an inconsistent state.  This assumes UnregisterFilter will return nil if the filter wasn't found
	//  at all (no rows deleted).
	spec := jb.OCR2OracleSpec
	transmitterID := spec.TransmitterID.String
	chain, err := d.legacyChains.Get(relayID.ChainID)
	if err != nil {
		d.lggr.Errorw("cleanupEVM: failed to get chain id", "chainId", relayID.ChainID, "err", err)
		return nil
	}
	lp := chain.LogPoller()

	var filters []string
	switch spec.PluginType {
	case types.OCR2Keeper:
		// Not worth the effort to validate and parse the job spec config to figure out whether this is v2.0 or v2.1,
		// simpler and faster to just Unregister them both
		filters, err = ocr2keeper.FilterNamesFromSpec20(spec)
		if err != nil {
			d.lggr.Errorw("failed to derive ocr2keeper filter names from spec", "err", err, "spec", spec)
		}
		filters21, err2 := ocr2keeper.FilterNamesFromSpec21(spec)
		if err2 != nil {
			d.lggr.Errorw("failed to derive ocr2keeper filter names from spec", "err", err, "spec", spec)
		}
		filters = append(filters, filters21...)
	case types.CCIPCommit:
		// Write PluginConfig bytes to send source/dest relayer provider + info outside of top level rargs/pargs over the wire
		var pluginJobSpecConfig ccipconfig.CommitPluginJobSpecConfig
		err = json.Unmarshal(spec.PluginConfig.Bytes(), &pluginJobSpecConfig)
		if err != nil {
			return err
		}

		dstProvider, err2 := d.ccipCommitGetDstProvider(ctx, jb, pluginJobSpecConfig, transmitterID)
		if err2 != nil {
			return err
		}

		srcProvider, _, err2 := d.ccipCommitGetSrcProvider(ctx, jb, pluginJobSpecConfig, transmitterID, dstProvider)
		if err2 != nil {
			return err
		}
		err2 = ccipcommit.UnregisterCommitPluginLpFilters(srcProvider, dstProvider)
		if err2 != nil {
			d.lggr.Errorw("failed to unregister ccip commit plugin filters", "err", err2, "spec", spec)
		}
		return nil
	case types.CCIPExecution:
		// PROVIDER BASED ARG CONSTRUCTION
		// Write PluginConfig bytes to send source/dest relayer provider + info outside of top level rargs/pargs over the wire
		var pluginJobSpecConfig ccipconfig.ExecPluginJobSpecConfig
		err = json.Unmarshal(spec.PluginConfig.Bytes(), &pluginJobSpecConfig)
		if err != nil {
			return err
		}

		dstProvider, err2 := d.ccipExecGetDstProvider(ctx, jb, pluginJobSpecConfig, transmitterID)
		if err2 != nil {
			return err
		}

		srcProvider, _, err2 := d.ccipExecGetSrcProvider(ctx, jb, pluginJobSpecConfig, transmitterID, dstProvider)
		if err2 != nil {
			return err
		}
		err2 = ccipexec.UnregisterExecPluginLpFilters(srcProvider, dstProvider)
		if err2 != nil {
			d.lggr.Errorw("failed to unregister ccip exec plugin filters", "err", err2, "spec", spec)
		}
		return nil
	case types.LLO:
		var pluginCfg lloconfig.PluginConfig
		err = json.Unmarshal(spec.PluginConfig.Bytes(), &pluginCfg)
		if err != nil {
			return err
		}
		var chainSelector uint64
		chainSelector, err = chainselectors.SelectorFromChainId(chain.ID().Uint64())
		if err != nil {
			return err
		}
		if err = llo.Cleanup(ctx, lp, pluginCfg.ChannelDefinitionsContractAddress, pluginCfg.DonID, d.ds, chainSelector); err != nil {
			// Cleanup is optimistic. Don't return error here, as we don't want
			// to block job deletion
			d.lggr.Errorw("failed to cleanup llo", "err", err, "spec", spec)
		}
		return nil
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
		err = lp.UnregisterFilter(ctx, filter)
		if err != nil {
			return errors.Wrapf(err, "Failed to unregister filter %s", filter)
		}
	}
	return nil
}

// ServicesForSpec returns the OCR2 services that need to run for this job
func (d *Delegate) ServicesForSpec(ctx context.Context, jb job.Job) ([]job.ServiceCtx, error) {
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
	lggr := logger.Sugared(d.lggr.Named(string(jb.Type)).Named(jb.ExternalJobID.String()).With(lggrCtx.Args()...))

	kvStore := job.NewKVStore(jb.ID, d.ds)

	rid, err := spec.RelayID()
	if err != nil {
		return nil, ErrJobSpecNoRelayer{Err: err, PluginName: string(spec.PluginType)}
	}

	if rid.Network == relay.NetworkEVM {
		lggr = logger.Sugared(lggr.With("evmChainID", rid.ChainID))

		chain, err2 := d.legacyChains.Get(rid.ChainID)
		if err2 != nil {
			return nil, fmt.Errorf("ServicesForSpec: could not get EVM chain %s: %w", rid.ChainID, err2)
		}
		effectiveTransmitterID, err2 = GetEVMEffectiveTransmitterID(ctx, &jb, chain, lggr)
		if err2 != nil {
			return nil, fmt.Errorf("ServicesForSpec failed to get evm transmitterID: %w", err2)
		}
	}
	spec.RelayConfig["effectiveTransmitterID"] = effectiveTransmitterID
	spec.RelayConfig.ApplyDefaultsOCR2(d.cfg.OCR2())

	ocrDB := NewDB(d.ds, spec.ID, 0, lggr)
	if d.peerWrapper == nil {
		return nil, errors.New("cannot setup OCR2 job service, libp2p peer was missing")
	} else if !d.peerWrapper.IsStarted() {
		return nil, errors.New("peerWrapper is not started. OCR2 jobs require a started and running p2p v2 peer")
	}

	lc, err := validate.ToLocalConfig(d.cfg.OCR2(), d.cfg.Insecure(), *spec)
	if err != nil {
		return nil, err
	}
	if err = libocr2.SanityCheckLocalConfig(lc); err != nil {
		return nil, err
	}
	lggr.Infow("OCR2 job using local config",
		"BlockchainTimeout", lc.BlockchainTimeout,
		"ContractConfigLoadTimeout", lc.ContractConfigLoadTimeout,
		"ContractConfigConfirmations", lc.ContractConfigConfirmations,
		"ContractConfigTrackerPollInterval", lc.ContractConfigTrackerPollInterval,
		"ContractTransmitterTransmitTimeout", lc.ContractTransmitterTransmitTimeout,
		"DatabaseTimeout", lc.DatabaseTimeout,
		"DefaultMaxDurationInitialization", lc.DefaultMaxDurationInitialization,
	)

	bootstrapPeers, err := ocrcommon.GetValidatedBootstrapPeers(spec.P2PV2Bootstrappers, d.peerWrapper.P2PConfig().V2().DefaultBootstrappers(), spec.AllowNoBootstrappers)
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

	ctx = lggrCtx.ContextWithValues(ctx)
	switch spec.PluginType {
	case types.Mercury:
		return d.newServicesMercury(ctx, lggr, jb, bootstrapPeers, kb, ocrDB, lc)

	case types.LLO:
		return d.newServicesLLO(ctx, lggr, jb, bootstrapPeers, kb, ocrDB, lc)

	case types.Median:
		return d.newServicesMedian(ctx, lggr, jb, bootstrapPeers, kb, kvStore, ocrDB, lc)

	case types.OCR2Keeper:
		return d.newServicesOCR2Keepers(ctx, lggr, jb, bootstrapPeers, kb, ocrDB, lc)

	case types.Functions:
		const (
			_ int32 = iota
			thresholdPluginId
			s4PluginId
		)
		thresholdPluginDB := NewDB(d.ds, spec.ID, thresholdPluginId, lggr)
		s4PluginDB := NewDB(d.ds, spec.ID, s4PluginId, lggr)
		return d.newServicesOCR2Functions(ctx, lggr, jb, bootstrapPeers, kb, ocrDB, thresholdPluginDB, s4PluginDB, lc)

	case types.GenericPlugin:
		return d.newServicesGenericPlugin(ctx, lggr, jb, bootstrapPeers, kb, ocrDB, lc, d.capabilitiesRegistry,
			kvStore)

	case types.CCIPCommit:
		return d.newServicesCCIPCommit(ctx, lggr, jb, bootstrapPeers, kb, ocrDB, lc, transmitterID)
	case types.CCIPExecution:
		return d.newServicesCCIPExecution(ctx, lggr, jb, bootstrapPeers, kb, ocrDB, lc, transmitterID)
	default:
		return nil, errors.Errorf("plugin type %s not supported", spec.PluginType)
	}
}

func GetEVMEffectiveTransmitterID(ctx context.Context, jb *job.Job, chain legacyevm.Chain, lggr logger.SugaredLogger) (string, error) {
	spec := jb.OCR2OracleSpec
	if spec.PluginType == types.Mercury || spec.PluginType == types.LLO {
		return spec.TransmitterID.String, nil
	}

	if spec.RelayConfig["sendingKeys"] == nil {
		spec.RelayConfig["sendingKeys"] = []string{spec.TransmitterID.String}
	} else if !spec.TransmitterID.Valid {
		sendingKeys, err := job.SendingKeysForJob(jb)
		if err != nil {
			return "", err
		}
		if len(sendingKeys) > 1 {
			return "", errors.New("no plugin should have more than 1 sending key")
		}
		spec.TransmitterID = null.StringFrom(sendingKeys[0])
	}

	// effectiveTransmitterID is the transmitter address registered on the ocr contract. This is by default the EOA account on the node.
	// In the case of forwarding, the transmitter address is the forwarder contract deployed onchain between EOA and OCR contract.
	// ForwardingAllowed cannot be set with Mercury, so this should always be false for mercury jobs
	if jb.ForwardingAllowed {
		if chain == nil {
			return "", errors.New("job forwarding requires non-nil chain")
		}

		var err error
		var effectiveTransmitterID common.Address
		// Median forwarders need special handling because of OCR2Aggregator transmitters whitelist.
		if spec.PluginType == types.Median {
			effectiveTransmitterID, err = chain.TxManager().GetForwarderForEOAOCR2Feeds(ctx, common.HexToAddress(spec.TransmitterID.String), common.HexToAddress(spec.ContractID))
		} else {
			effectiveTransmitterID, err = chain.TxManager().GetForwarderForEOA(ctx, common.HexToAddress(spec.TransmitterID.String))
		}

		if err == nil {
			return effectiveTransmitterID.String(), nil
		} else if !spec.TransmitterID.Valid {
			return "", errors.New("failed to get forwarder address and transmitterID is not set")
		}
		lggr.Warnw("Skipping forwarding for job, will fallback to default behavior", "job", jb.Name, "err", err)
	}

	return spec.TransmitterID.String, nil
}

type connProvider interface {
	ClientConn() grpc.ClientConnInterface
}

func (d *Delegate) newServicesGenericPlugin(
	ctx context.Context,
	lggr logger.SugaredLogger,
	jb job.Job,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	ocrDB *db,
	lc ocrtypes.LocalConfig,
	capabilitiesRegistry core.CapabilitiesRegistry,
	keyValueStore core.KeyValueStore,
) (srvs []job.ServiceCtx, err error) {
	spec := jb.OCR2OracleSpec
	// NOTE: we don't need to validate this config, since that happens as part of creating the job.
	// See: validate/validate.go's `validateSpec`.
	pCfg := validate.OCR2GenericPluginConfig{}
	err = json.Unmarshal(spec.PluginConfig.Bytes(), &pCfg)
	if err != nil {
		return nil, err
	}
	// NOTE: we don't need to validate the strategy, since that happens as part of creating the job.
	// See: validate/validate.go's `validateSpec`.
	onchainSigningStrategy := validate.OCR2OnchainSigningStrategy{}
	err = json.Unmarshal(spec.OnchainSigningStrategy.Bytes(), &onchainSigningStrategy)
	if err != nil {
		return nil, err
	}

	plugEnv := env.NewPlugin(pCfg.PluginName)

	command := pCfg.Command
	if command == "" {
		command = plugEnv.Cmd.Get()
	}

	// Add the default pipeline to the pluginConfig
	pCfg.Pipelines = append(
		pCfg.Pipelines,
		validate.PipelineSpec{Name: "__DEFAULT_PIPELINE__", Spec: jb.Pipeline.Source},
	)

	rid, err := spec.RelayID()
	if err != nil {
		return nil, ErrJobSpecNoRelayer{PluginName: pCfg.PluginName, Err: err}
	}

	relayerSet, err := generic.NewRelayerSet(d.RelayGetter, jb.ExternalJobID, jb.ID, d.isNewlyCreatedJob)
	if err != nil {
		return nil, fmt.Errorf("failed to create relayer set: %w", err)
	}

	relayer, err := d.RelayGetter.Get(rid)
	if err != nil {
		return nil, ErrRelayNotEnabled{Err: err, Relay: spec.Relay, PluginName: pCfg.PluginName}
	}

	provider, err := relayer.NewPluginProvider(ctx, types.RelayArgs{
		ExternalJobID: jb.ExternalJobID,
		JobID:         jb.ID,
		OracleSpecID:  spec.ID,
		ContractID:    spec.ContractID,
		New:           d.isNewlyCreatedJob,
		RelayConfig:   spec.RelayConfig.Bytes(),
		ProviderType:  pCfg.ProviderType,
	}, types.PluginArgs{
		TransmitterID: spec.TransmitterID.String,
		PluginConfig:  spec.PluginConfig.Bytes(),
	})
	if err != nil {
		return nil, err
	}
	srvs = append(srvs, provider)

	envVars, err := plugins.ParseEnvFile(plugEnv.Env.Get())
	if err != nil {
		return nil, fmt.Errorf("failed to parse median env file: %w", err)
	}
	if len(pCfg.EnvVars) > 0 {
		for k, v := range pCfg.EnvVars {
			envVars = append(envVars, k+"="+v)
		}
	}

	pluginLggr := lggr.Named(pCfg.PluginName).Named(spec.ContractID).Named(spec.GetID())
	cmdFn, grpcOpts, err := d.cfg.RegisterLOOP(plugins.CmdConfig{
		ID:  fmt.Sprintf("%s-%s-%s", pCfg.PluginName, spec.ContractID, spec.GetID()),
		Cmd: command,
		Env: envVars,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register loop: %w", err)
	}

	errorLog := &errorLog{jobID: jb.ID, recordError: d.jobORM.RecordError}
	var providerClientConn grpc.ClientConnInterface
	providerConn, ok := provider.(connProvider)
	if ok {
		providerClientConn = providerConn.ClientConn()
	} else {
		// We chose to deal with the difference between a LOOP provider and an embedded provider here rather than
		// in NewServerAdapter because this has a smaller blast radius, as the scope of this workaround is to
		// enable the medianpoc for EVM and not touch the other providers.
		// TODO: remove this workaround when the EVM relayer is running inside of an LOOPP
		d.lggr.Info("provider is not a LOOPP provider, switching to provider server")

		ps, err2 := loop.NewProviderServer(provider, types.OCR2PluginType(pCfg.ProviderType), d.lggr)
		if err2 != nil {
			return nil, fmt.Errorf("cannot start EVM provider server: %w", err2)
		}
		providerClientConn, err2 = ps.GetConn()
		if err2 != nil {
			return nil, fmt.Errorf("cannot connect to EVM provider server: %w", err)
		}
		srvs = append(srvs, ps)
	}

	pc, err := json.Marshal(pCfg.Config)
	if err != nil {
		return nil, fmt.Errorf("cannot dump plugin config to string before sending to plugin: %w", err)
	}

	pluginConfig := core.ReportingPluginServiceConfig{
		PluginName:    pCfg.PluginName,
		Command:       command,
		ProviderType:  pCfg.ProviderType,
		TelemetryType: pCfg.TelemetryType,
		PluginConfig:  string(pc),
	}

	pr := generic.NewPipelineRunnerAdapter(pluginLggr, jb, d.pipelineRunner)
	ta := generic.NewTelemetryAdapter(d.monitoringEndpointGen)

	oracleEndpoint := d.monitoringEndpointGen.GenMonitoringEndpoint(
		rid.Network,
		rid.ChainID,
		spec.ContractID,
		synchronization.TelemetryType(pCfg.TelemetryType),
	)

	ocrLogger := ocrcommon.NewOCRWrapper(lggr, d.cfg.OCR2().TraceLogging(), func(ctx context.Context, msg string) {
		lggr.ErrorIf(d.jobORM.RecordError(ctx, jb.ID, msg), "unable to record error")
	})
	srvs = append(srvs, ocrLogger)

	switch pCfg.OCRVersion {
	case 2:
		plugin := reportingplugins.NewLOOPPService(pluginLggr, grpcOpts, cmdFn, pluginConfig, providerClientConn, pr, ta,
			errorLog, keyValueStore, relayerSet)
		oracleArgs := libocr2.OCR2OracleArgs{
			BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
			V2Bootstrappers:              bootstrapPeers,
			Database:                     ocrDB,
			LocalConfig:                  lc,
			Logger:                       ocrLogger,
			MonitoringEndpoint:           oracleEndpoint,
			OffchainKeyring:              kb,
			OnchainKeyring:               kb,
			ContractTransmitter:          provider.ContractTransmitter(),
			ContractConfigTracker:        provider.ContractConfigTracker(),
			OffchainConfigDigester:       provider.OffchainConfigDigester(),
			MetricsRegisterer:            prometheus.WrapRegistererWith(map[string]string{"job_name": jb.Name.ValueOrZero()}, prometheus.DefaultRegisterer),
		}
		oracleArgs.ReportingPluginFactory = plugin
		srvs = append(srvs, plugin)
		oracle, oracleErr := libocr2.NewOracle(oracleArgs)
		if oracleErr != nil {
			return nil, oracleErr
		}
		srvs = append(srvs, job.NewServiceAdapter(oracle))

	case 3:
		// OCR3 with OCR2 OnchainKeyring and ContractTransmitter
		plugin := ocr3.NewLOOPPService(
			pluginLggr,
			grpcOpts,
			cmdFn,
			pluginConfig,
			providerClientConn,
			pr,
			ta,
			errorLog,
			capabilitiesRegistry,
			keyValueStore,
			relayerSet,
		)

		// Adapt the provider's contract transmitter for OCR3, unless
		// the provider exposes an OCR3ContractTransmitter interface, in which case
		// we'll use that instead.
		contractTransmitter := ocr3types.ContractTransmitter[[]byte](
			ocrcommon.NewOCR3ContractTransmitterAdapter(provider.ContractTransmitter()),
		)
		if ocr3Provider, ok := provider.(types.OCR3ContractTransmitter); ok {
			contractTransmitter = ocr3Provider.OCR3ContractTransmitter()
		}
		var onchainKeyringAdapter ocr3types.OnchainKeyring[[]byte]
		if onchainSigningStrategy.IsMultiChain() {
			// We are extracting the config beforehand
			keyBundles := map[string]ocr2key.KeyBundle{}
			for name := range onchainSigningStrategy.ConfigCopy() {
				kbID, ostErr := onchainSigningStrategy.KeyBundleID(name)
				if ostErr != nil {
					return nil, ostErr
				}
				os, ostErr := d.ks.Get(kbID)
				if ostErr != nil {
					return nil, ostErr
				}
				keyBundles[name] = os
			}
			onchainKeyringAdapter, err = ocrcommon.NewOCR3OnchainKeyringMultiChainAdapter(keyBundles, lggr)
			if err != nil {
				return nil, err
			}
		} else {
			onchainKeyringAdapter = ocrcommon.NewOCR3OnchainKeyringAdapter(kb)
		}
		oracleArgs := libocr2.OCR3OracleArgs[[]byte]{
			BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
			V2Bootstrappers:              bootstrapPeers,
			ContractConfigTracker:        provider.ContractConfigTracker(),
			ContractTransmitter:          contractTransmitter,
			Database:                     ocrDB,
			LocalConfig:                  lc,
			Logger:                       ocrLogger,
			MonitoringEndpoint:           oracleEndpoint,
			OffchainConfigDigester:       provider.OffchainConfigDigester(),
			OffchainKeyring:              kb,
			OnchainKeyring:               onchainKeyringAdapter,
			MetricsRegisterer:            prometheus.WrapRegistererWith(map[string]string{"job_name": jb.Name.ValueOrZero()}, prometheus.DefaultRegisterer),
		}
		oracleArgs.ReportingPluginFactory = plugin
		srvs = append(srvs, plugin)
		oracle, err := libocr2.NewOracle(oracleArgs)
		if err != nil {
			return nil, err
		}
		srvs = append(srvs, job.NewServiceAdapter(oracle))

	default:
		return nil, fmt.Errorf("unknown OCR version: %d", pCfg.OCRVersion)
	}

	return srvs, nil
}

func (d *Delegate) newServicesMercury(
	ctx context.Context,
	lggr logger.SugaredLogger,
	jb job.Job,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	ocrDB *db,
	lc ocrtypes.LocalConfig,
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
		return nil, ErrJobSpecNoRelayer{Err: err, PluginName: "mercury"}
	}
	if rid.Network != relay.NetworkEVM {
		return nil, fmt.Errorf("mercury services: expected EVM relayer got %q", rid.Network)
	}
	relayer, err := d.RelayGetter.Get(rid)
	if err != nil {
		return nil, ErrRelayNotEnabled{Err: err, Relay: spec.Relay, PluginName: "mercury"}
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

	lc.ContractConfigTrackerPollInterval = 1 * time.Second // This is the fastest that libocr supports. See: https://github.com/smartcontractkit/offchain-reporting/pull/520

	// Disable OCR debug+info logging for legacy mercury jobs unless tracelogging is enabled, because its simply too verbose (150 jobs => ~50k logs per second)
	ocrLogger := ocrcommon.NewOCRWrapper(llo.NewSuppressedLogger(lggr, d.cfg.OCR2().TraceLogging(), d.cfg.OCR2().TraceLogging()), d.cfg.OCR2().TraceLogging(), func(ctx context.Context, msg string) {
		lggr.ErrorIf(d.jobORM.RecordError(ctx, jb.ID, msg), "unable to record error")
	})

	var relayConfig evmrelaytypes.RelayConfig
	err = json.Unmarshal(jb.OCR2OracleSpec.RelayConfig.Bytes(), &relayConfig)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling relay config: %w", err)
	}

	var telemetryType synchronization.TelemetryType
	if relayConfig.EnableTriggerCapability && len(jb.OCR2OracleSpec.PluginConfig) == 0 {
		telemetryType = synchronization.OCR3DataFeeds
		// First use case for TriggerCapability transmission is Data Feeds, so telemetry should be routed accordingly.
		// This is only true if TriggerCapability is the *only* transmission method (PluginConfig is empty).
	} else {
		telemetryType = synchronization.OCR3Mercury
	}

	oracleArgsNoPlugin := libocr2.MercuryOracleArgs{
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		ContractTransmitter:          mercuryProvider.ContractTransmitter(),
		ContractConfigTracker:        mercuryProvider.ContractConfigTracker(),
		Database:                     ocrDB,
		LocalConfig:                  lc,
		Logger:                       ocrLogger,
		MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(rid.Network, rid.ChainID, spec.FeedID.String(), telemetryType),
		OffchainConfigDigester:       mercuryProvider.OffchainConfigDigester(),
		OffchainKeyring:              kb,
		OnchainKeyring:               kb,
		MetricsRegisterer:            prometheus.WrapRegistererWith(map[string]string{"job_name": jb.Name.ValueOrZero()}, prometheus.DefaultRegisterer),
	}

	chEnhancedTelem := make(chan ocrcommon.EnhancedTelemetryMercuryData, 100)

	mCfg := mercury.NewMercuryConfig(d.cfg.JobPipeline().MaxSuccessfulRuns(), d.cfg.JobPipeline().ResultWriteQueueDepth(), d.cfg)

	mercuryServices, err2 := mercury.NewServices(jb, mercuryProvider, d.pipelineRunner, lggr, oracleArgsNoPlugin, mCfg, chEnhancedTelem, d.mercuryORM, (mercuryutils.FeedID)(*spec.FeedID), relayConfig.EnableTriggerCapability)

	if ocrcommon.ShouldCollectEnhancedTelemetryMercury(jb) {
		enhancedTelemService := ocrcommon.NewEnhancedTelemetryService(&jb, chEnhancedTelem, make(chan struct{}), d.monitoringEndpointGen.GenMonitoringEndpoint(rid.Network, rid.ChainID, spec.FeedID.String(), synchronization.EnhancedEAMercury), lggr.Named("EnhancedTelemetryMercury"))
		mercuryServices = append(mercuryServices, enhancedTelemService)
	} else {
		lggr.Infow("Enhanced telemetry is disabled for mercury job", "job", jb.Name)
	}

	mercuryServices = append(mercuryServices, ocrLogger)

	return mercuryServices, err2
}

func (d *Delegate) newServicesLLO(
	ctx context.Context,
	lggr logger.SugaredLogger,
	jb job.Job,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	ocrDB *db,
	lc ocrtypes.LocalConfig,
) ([]job.ServiceCtx, error) {
	spec := jb.OCR2OracleSpec
	transmitterID := spec.TransmitterID.String
	if len(transmitterID) != 64 {
		return nil, errors.Errorf("ServicesForSpec: streams job type requires transmitter ID to be a 32-byte hex string, got: %q", transmitterID)
	}
	if _, err := hex.DecodeString(transmitterID); err != nil {
		return nil, errors.Wrapf(err, "ServicesForSpec: streams job type requires transmitter ID to be a 32-byte hex string, got: %q", transmitterID)
	}

	rid, err := spec.RelayID()
	if err != nil {
		return nil, ErrJobSpecNoRelayer{Err: err, PluginName: "streams"}
	}
	relayer, err := d.RelayGetter.Get(rid)
	if err != nil {
		return nil, ErrRelayNotEnabled{Err: err, Relay: spec.Relay, PluginName: "streams"}
	}

	provider, err2 := relayer.NewLLOProvider(ctx,
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

	var pluginCfg lloconfig.PluginConfig
	if err = json.Unmarshal(spec.PluginConfig.Bytes(), &pluginCfg); err != nil {
		return nil, err
	}
	lggr = logger.Sugared(lggr.Named("LLO").With("donID", pluginCfg.DonID, "channelDefinitionsContractAddress", pluginCfg.ChannelDefinitionsContractAddress))

	// Handle key bundle IDs explicitly specified in job spec
	kbm := make(map[llotypes.ReportFormat]llo.Key)
	for rfStr, kbid := range pluginCfg.KeyBundleIDs {
		k, err3 := d.ks.Get(kbid)
		if err3 != nil {
			return nil, fmt.Errorf("job %d (%s) specified key bundle ID %q for report format %s, but got error trying to load it: %w", jb.ID, jb.Name.ValueOrZero(), kbid, rfStr, err3)
		}
		rf, err4 := llotypes.ReportFormatFromString(rfStr)
		if err4 != nil {
			return nil, fmt.Errorf("job %d (%s) specified key bundle ID %q for report format %s, but it is not a recognized report format: %w", jb.ID, jb.Name.ValueOrZero(), kbid, rfStr, err4)
		}
		kbm[rf] = k
	}

	// Use the default key bundle if not specified
	// TODO: MERC-3594
	evmKeySignedFormats := []llotypes.ReportFormat{
		// NOTE: Just use EVM key for signing everything. This isn't
		// necessarily required, just seems easiest since it's the only key
		// type available for now.
		llotypes.ReportFormatJSON,
		llotypes.ReportFormatEVMPremiumLegacy,
		llotypes.ReportFormatRetirement,
		llotypes.ReportFormatEVMABIEncodeUnpacked,
		llotypes.ReportFormatCapabilityTrigger,
		llotypes.ReportFormatEVMStreamlined,
	}
	for _, rf := range evmKeySignedFormats {
		if _, exists := kbm[rf]; !exists {
			// Use the first if unspecified
			kbs, err3 := d.ks.GetAllOfType("evm")
			if err3 != nil {
				return nil, err3
			}
			if len(kbs) == 0 {
				return nil, fmt.Errorf("no on-chain signing keys found for report format %s", "evm")
			} else if len(kbs) > 1 {
				lggr.Debugf("Multiple on-chain signing keys found for report format %s, using the first", rf.String())
			}
			kbm[rf] = kbs[0]
		}
	}

	// FIXME: This is a bit confusing because the OCR2 key bundle actually
	// includes an EVM on-chain key... but LLO only uses the key bundle for the
	// offchain keys and the suppoprted onchain keys are defined in the plugin
	// config on the job spec instead.
	// https://smartcontract-it.atlassian.net/browse/MERC-3594
	lggr.Infof("Using on-chain signing keys for LLO job %d (%s): %v", jb.ID, jb.Name.ValueOrZero(), kbm)
	kr := llo.NewOnchainKeyring(lggr, kbm, pluginCfg.DonID)

	telemetryContractID := fmt.Sprintf("%s/%d", spec.ContractID, pluginCfg.DonID)

	cfg := llo.DelegateConfig{
		Logger:     lggr,
		DataSource: d.ds,
		Runner:     d.pipelineRunner,
		Registry:   d.streamRegistry,

		JobName:            jb.Name,
		CaptureEATelemetry: jb.OCR2OracleSpec.CaptureEATelemetry,
		// NOTE: These can be turned off/on, or made configurable in future if
		// necessary
		CaptureObservationTelemetry: jb.OCR2OracleSpec.CaptureEATelemetry,
		CaptureOutcomeTelemetry:     jb.OCR2OracleSpec.CaptureEATelemetry,
		CaptureReportTelemetry:      false,

		ChannelDefinitionCache:   provider.ChannelDefinitionCache(),
		RetirementReportCache:    d.retirementReportCache,
		ShouldRetireCache:        provider.ShouldRetireCache(),
		RetirementReportCodec:    datastreamsllo.StandardRetirementReportCodec{},
		PluginMonitoringEndpoint: d.monitoringEndpointGen.GenMultitypeMonitoringEndpoint(rid.Network, rid.ChainID, telemetryContractID),
		DonID:                    pluginCfg.DonID,
		ChainID:                  rid.ChainID,

		TraceLogging:                 d.cfg.OCR2().TraceLogging(),
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		ContractTransmitter:          provider.ContractTransmitter(),
		ContractConfigTrackers:       provider.ContractConfigTrackers(),
		LocalConfig:                  lc,
		OCR3MonitoringEndpoint:       d.monitoringEndpointGen.GenMonitoringEndpoint(rid.Network, rid.ChainID, telemetryContractID, synchronization.OCR3Mercury),
		OffchainConfigDigester:       provider.OffchainConfigDigester(),
		OffchainKeyring:              kb,
		OnchainKeyring:               kr,

		// Enable verbose logging if either Mercury.VerboseLogging is on or OCR2.TraceLogging is on
		ReportingPluginConfig: datastreamsllo.Config{VerboseLogging: d.cfg.Mercury().VerboseLogging() || d.cfg.OCR2().TraceLogging()},
		NewOCR3DB: func(pluginID int32) ocr3types.Database {
			return NewDB(d.ds, spec.ID, pluginID, lggr)
		},
	}
	oracle, err := llo.NewDelegate(cfg)
	if err != nil {
		return nil, err
	}
	return []job.ServiceCtx{provider, oracle}, nil
}

func (d *Delegate) newServicesMedian(
	ctx context.Context,
	lggr logger.SugaredLogger,
	jb job.Job,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	kvStore job.KVStore,
	ocrDB *db,
	lc ocrtypes.LocalConfig,
) ([]job.ServiceCtx, error) {
	spec := jb.OCR2OracleSpec

	rid, err := spec.RelayID()
	if err != nil {
		return nil, ErrJobSpecNoRelayer{Err: err, PluginName: "median"}
	}

	ocrLogger := ocrcommon.NewOCRWrapper(lggr, d.cfg.OCR2().TraceLogging(), func(ctx context.Context, msg string) {
		lggr.ErrorIf(d.jobORM.RecordError(ctx, jb.ID, msg), "unable to record error")
	})

	oracleArgsNoPlugin := libocr2.OCR2OracleArgs{
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		Database:                     ocrDB,
		LocalConfig:                  lc,
		Logger:                       ocrLogger,
		MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(rid.Network, rid.ChainID, spec.ContractID, synchronization.OCR2Median),
		OffchainKeyring:              kb,
		OnchainKeyring:               kb,
		MetricsRegisterer:            prometheus.WrapRegistererWith(map[string]string{"job_name": jb.Name.ValueOrZero()}, prometheus.DefaultRegisterer),
	}
	errorLog := &errorLog{jobID: jb.ID, recordError: d.jobORM.RecordError}
	enhancedTelemChan := make(chan ocrcommon.EnhancedTelemetryData, 100)
	mConfig := median.NewMedianConfig(
		d.cfg.JobPipeline().MaxSuccessfulRuns(),
		d.cfg.JobPipeline().ResultWriteQueueDepth(),
		d.cfg,
	)

	relayer, err := d.RelayGetter.Get(rid)
	if err != nil {
		return nil, ErrRelayNotEnabled{Err: err, PluginName: "median", Relay: spec.Relay}
	}

	medianServices, err2 := median.NewMedianServices(ctx, jb, d.isNewlyCreatedJob, relayer, kvStore, d.pipelineRunner, lggr, oracleArgsNoPlugin, mConfig, enhancedTelemChan, errorLog)

	if ocrcommon.ShouldCollectEnhancedTelemetry(&jb) {
		enhancedTelemService := ocrcommon.NewEnhancedTelemetryService(&jb, enhancedTelemChan, make(chan struct{}), d.monitoringEndpointGen.GenMonitoringEndpoint(rid.Network, rid.ChainID, spec.ContractID, synchronization.EnhancedEA), lggr.Named("EnhancedTelemetry"))
		medianServices = append(medianServices, enhancedTelemService)
	} else {
		lggr.Infow("Enhanced telemetry is disabled for job", "job", jb.Name)
	}

	medianServices = append(medianServices, ocrLogger)

	return medianServices, err2
}

func (d *Delegate) newServicesOCR2Keepers(
	ctx context.Context,
	lggr logger.SugaredLogger,
	jb job.Job,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	ocrDB *db,
	lc ocrtypes.LocalConfig,
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
		return d.newServicesOCR2Keepers21(ctx, lggr, jb, bootstrapPeers, kb, ocrDB, lc, cfg, spec)
	case "v2.1+":
		// Future contracts of v2.1 (v2.x) will use the same job spec as v2.1
		return d.newServicesOCR2Keepers21(ctx, lggr, jb, bootstrapPeers, kb, ocrDB, lc, cfg, spec)
	case "v2.0":
		return d.newServicesOCR2Keepers20(ctx, lggr, jb, bootstrapPeers, kb, ocrDB, lc, cfg, spec)
	default:
		return d.newServicesOCR2Keepers20(ctx, lggr, jb, bootstrapPeers, kb, ocrDB, lc, cfg, spec)
	}
}

func (d *Delegate) newServicesOCR2Keepers21(
	ctx context.Context,
	lggr logger.SugaredLogger,
	jb job.Job,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	ocrDB *db,
	lc ocrtypes.LocalConfig,
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
		return nil, ErrJobSpecNoRelayer{Err: err, PluginName: "keeper2"}
	}
	if rid.Network != relay.NetworkEVM {
		return nil, fmt.Errorf("keeper2 services: expected EVM relayer got %q", rid.Network)
	}

	transmitterID := spec.TransmitterID.String
	relayer, err := d.RelayGetter.Get(rid)
	if err != nil {
		return nil, ErrRelayNotEnabled{Err: err, Relay: spec.Relay, PluginName: "ocr2keepers"}
	}

	provider, err := relayer.NewPluginProvider(ctx,
		types.RelayArgs{
			ExternalJobID:      jb.ExternalJobID,
			JobID:              jb.ID,
			ContractID:         spec.ContractID,
			New:                d.isNewlyCreatedJob,
			RelayConfig:        spec.RelayConfig.Bytes(),
			ProviderType:       string(spec.PluginType),
			MercuryCredentials: mc,
		}, types.PluginArgs{
			TransmitterID: transmitterID,
			PluginConfig:  spec.PluginConfig.Bytes(),
		})
	if err != nil {
		return nil, err
	}

	keeperProvider, ok := provider.(types.AutomationProvider)
	if !ok {
		return nil, errors.New("could not coerce PluginProvider to AutomationProvider")
	}

	services, err := ocr2keeper.EVMDependencies21(kb)
	if err != nil {
		return nil, errors.Wrap(err, "could not build dependencies for ocr2 keepers")
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
	ocrLogger := ocrcommon.NewOCRWrapper(lggr, d.cfg.OCR2().TraceLogging(), func(ctx context.Context, msg string) {
		lggr.ErrorIf(d.jobORM.RecordError(ctx, jb.ID, msg), "unable to record error")
	})

	dConf := ocr2keepers21.DelegateConfig{
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		ContractTransmitter:          evmrelay.NewKeepersOCR3ContractTransmitter(keeperProvider.ContractTransmitter()),
		ContractConfigTracker:        keeperProvider.ContractConfigTracker(),
		MetricsRegisterer:            prometheus.WrapRegistererWith(map[string]string{"job_name": jb.Name.ValueOrZero()}, prometheus.DefaultRegisterer),
		KeepersDatabase:              ocrDB,
		Logger:                       ocrLogger,
		MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(rid.Network, rid.ChainID, spec.ContractID, synchronization.OCR3Automation),
		OffchainConfigDigester:       keeperProvider.OffchainConfigDigester(),
		OffchainKeyring:              kb,
		OnchainKeyring:               services.Keyring(),
		LocalConfig:                  lc,
		LogProvider:                  keeperProvider.LogEventProvider(),
		EventProvider:                keeperProvider.TransmitEventProvider(),
		Runnable:                     keeperProvider.Registry(),
		Encoder:                      keeperProvider.Encoder(),
		BlockSubscriber:              keeperProvider.BlockSubscriber(),
		RecoverableProvider:          keeperProvider.LogRecoverer(),
		PayloadBuilder:               keeperProvider.PayloadBuilder(),
		UpkeepProvider:               keeperProvider.UpkeepProvider(),
		UpkeepStateUpdater:           keeperProvider.UpkeepStateStore(),
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

	automationServices := []job.ServiceCtx{
		keeperProvider,
		keeperProvider.Registry(),
		keeperProvider.BlockSubscriber(),
		keeperProvider.LogEventProvider(),
		keeperProvider.LogRecoverer(),
		keeperProvider.UpkeepStateStore(),
		keeperProvider.TransmitEventProvider(),
		pluginService,
		ocrLogger,
	}

	if cfg.CaptureAutomationCustomTelemetry != nil && *cfg.CaptureAutomationCustomTelemetry ||
		cfg.CaptureAutomationCustomTelemetry == nil && d.cfg.OCR2().CaptureAutomationCustomTelemetry() {
		endpoint := d.monitoringEndpointGen.GenMonitoringEndpoint(rid.Network, rid.ChainID, spec.ContractID, synchronization.AutomationCustom)
		customTelemService, custErr := autotelemetry21.NewAutomationCustomTelemetryService(
			endpoint,
			lggr,
			keeperProvider.BlockSubscriber(),
			keeperProvider.ContractConfigTracker(),
		)
		if custErr != nil {
			return nil, errors.Wrap(custErr, "Error when creating AutomationCustomTelemetryService")
		}
		automationServices = append(automationServices, customTelemService)
	}

	return automationServices, nil
}

func (d *Delegate) newServicesOCR2Keepers20(
	ctx context.Context,
	lggr logger.SugaredLogger,
	jb job.Job,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	ocrDB *db,
	lc ocrtypes.LocalConfig,
	cfg ocr2keeper.PluginConfig,
	spec *job.OCR2OracleSpec,
) ([]job.ServiceCtx, error) {
	rid, err := spec.RelayID()
	if err != nil {
		return nil, ErrJobSpecNoRelayer{Err: err, PluginName: "keepers2.0"}
	}
	if rid.Network != relay.NetworkEVM {
		return nil, fmt.Errorf("keepers2.0 services: expected EVM relayer got %q", rid.Network)
	}
	chain, err2 := d.legacyChains.Get(rid.ChainID)
	if err2 != nil {
		return nil, fmt.Errorf("keepers2.0 services: failed to get chain (%s): %w", rid.ChainID, err2)
	}

	cid := chain.ID()
	ks := keys.NewChainStore(keystore.NewEthSigner(d.ethKs, cid), cid)
	keeperProvider, rgstry, encoder, logProvider, err2 := ocr2keeper.EVMDependencies20(ctx, jb, d.ds, lggr, chain, ks)
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

	ocrLogger := ocrcommon.NewOCRWrapper(lggr, d.cfg.OCR2().TraceLogging(), func(ctx context.Context, msg string) {
		lggr.ErrorIf(d.jobORM.RecordError(ctx, jb.ID, msg), "unable to record error")
	})

	dConf := ocr2keepers20.DelegateConfig{
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		ContractTransmitter:          keeperProvider.ContractTransmitter(),
		ContractConfigTracker:        keeperProvider.ContractConfigTracker(),
		MetricsRegisterer:            prometheus.WrapRegistererWith(map[string]string{"job_name": jb.Name.ValueOrZero()}, prometheus.DefaultRegisterer),
		KeepersDatabase:              ocrDB,
		LocalConfig:                  lc,
		Logger:                       ocrLogger,
		MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(rid.Network, rid.ChainID, spec.ContractID, synchronization.OCR2Automation),
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

	return []job.ServiceCtx{
		job.NewServiceAdapter(runr),
		keeperProvider,
		rgstry,
		logProvider,
		ocrLogger,
		pluginService,
	}, nil
}

func (d *Delegate) newServicesOCR2Functions(
	ctx context.Context,
	lggr logger.SugaredLogger,
	jb job.Job,
	bootstrapPeers []commontypes.BootstrapperLocator,
	kb ocr2key.KeyBundle,
	functionsOcrDB *db,
	thresholdOcrDB *db,
	s4OcrDB *db,
	lc ocrtypes.LocalConfig,
) ([]job.ServiceCtx, error) {
	spec := jb.OCR2OracleSpec

	rid, err := spec.RelayID()
	if err != nil {
		return nil, ErrJobSpecNoRelayer{Err: err, PluginName: "functions"}
	}
	if rid.Network != relay.NetworkEVM {
		return nil, fmt.Errorf("functions services: expected EVM relayer got %q", rid.Network)
	}
	chain, err := d.legacyChains.Get(rid.ChainID)
	if err != nil {
		return nil, fmt.Errorf("functions services: failed to get chain %s: %w", rid.ChainID, err)
	}
	cid := chain.ID()
	ks := keys.NewChainStore(keystore.NewEthSigner(d.ethKs, cid), cid)
	createPluginProvider := func(pluginType functionsRelay.FunctionsPluginType, relayerName string) (evmrelaytypes.FunctionsProvider, error) {
		return evmrelay.NewFunctionsProvider(
			ctx,
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
			ks,
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

	ocrLogger := ocrcommon.NewOCRWrapper(lggr, d.cfg.OCR2().TraceLogging(), func(ctx context.Context, msg string) {
		lggr.ErrorIf(d.jobORM.RecordError(ctx, jb.ID, msg), "unable to record error")
	})

	functionsOracleArgs := libocr2.OCR2OracleArgs{
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		ContractTransmitter:          functionsProvider.ContractTransmitter(),
		ContractConfigTracker:        functionsProvider.ContractConfigTracker(),
		Database:                     functionsOcrDB,
		LocalConfig:                  lc,
		Logger:                       ocrLogger,
		MonitoringEndpoint:           d.monitoringEndpointGen.GenMonitoringEndpoint(rid.Network, rid.ChainID, spec.ContractID, synchronization.OCR2Functions),
		OffchainConfigDigester:       functionsProvider.OffchainConfigDigester(),
		OffchainKeyring:              kb,
		OnchainKeyring:               kb,
		ReportingPluginFactory:       nil, // To be set by NewFunctionsServices
		MetricsRegisterer:            prometheus.WrapRegistererWith(map[string]string{"job_name": jb.Name.ValueOrZero()}, prometheus.DefaultRegisterer),
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
		MetricsRegisterer:      prometheus.WrapRegistererWith(map[string]string{"job_name": jb.Name.ValueOrZero()}, prometheus.DefaultRegisterer),
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
		MetricsRegisterer:      prometheus.WrapRegistererWith(map[string]string{"job_name": jb.Name.ValueOrZero()}, prometheus.DefaultRegisterer),
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
		DS:                d.ds,
		Chain:             chain,
		ContractID:        spec.ContractID,
		Logger:            lggr,
		MailMon:           d.mailMon,
		URLsMonEndpoint:   d.monitoringEndpointGen.GenMonitoringEndpoint(rid.Network, rid.ChainID, spec.ContractID, synchronization.FunctionsRequests),
		EthKeystore:       ks,
		ThresholdKeyShare: thresholdKeyShare,
		LogPollerWrapper:  functionsProvider.LogPollerWrapper(),
	}

	functionsServices, err := functions.NewFunctionsServices(ctx, &functionsOracleArgs, &thresholdOracleArgs, &s4OracleArgs, &functionsServicesConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error calling NewFunctionsServices")
	}

	return append([]job.ServiceCtx{functionsProvider, thresholdProvider, s4Provider, ocrLogger}, functionsServices...), nil
}

func (d *Delegate) newServicesCCIPCommit(ctx context.Context, lggr logger.SugaredLogger, jb job.Job, bootstrapPeers []commontypes.BootstrapperLocator, kb ocr2key.KeyBundle, ocrDB *db, lc ocrtypes.LocalConfig, transmitterID string) ([]job.ServiceCtx, error) {
	spec := jb.OCR2OracleSpec
	if !isCCIPSupportedNetwork(spec.Relay) {
		return nil, fmt.Errorf("chain not supported for CCIP commit: %s", spec.Relay)
	}
	dstRid, err := spec.RelayID()
	if err != nil {
		return nil, ErrJobSpecNoRelayer{Err: err, PluginName: string(spec.PluginType)}
	}

	logError := func(msg string) {
		lggr.ErrorIf(d.jobORM.RecordError(context.Background(), jb.ID, msg), "unable to record error")
	}

	// Write PluginConfig bytes to send source/dest relayer provider + info outside of top level rargs/pargs over the wire
	var pluginJobSpecConfig ccipconfig.CommitPluginJobSpecConfig
	err = json.Unmarshal(spec.PluginConfig.Bytes(), &pluginJobSpecConfig)
	if err != nil {
		return nil, err
	}

	dstChainID, err := strconv.ParseInt(dstRid.ChainID, 10, 64)
	if err != nil {
		return nil, err
	}

	dstProvider, err := d.ccipCommitGetDstProvider(ctx, jb, pluginJobSpecConfig, transmitterID)
	if err != nil {
		return nil, err
	}

	srcProvider, srcChainID, err := d.ccipCommitGetSrcProvider(ctx, jb, pluginJobSpecConfig, transmitterID, dstProvider)
	if err != nil {
		return nil, err
	}

	oracleArgsNoPlugin := libocr2.OCR2OracleArgs{
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		ContractTransmitter:          dstProvider.ContractTransmitter(),
		ContractConfigTracker:        dstProvider.ContractConfigTracker(),
		Database:                     ocrDB,
		LocalConfig:                  lc,
		MonitoringEndpoint: d.monitoringEndpointGen.GenMonitoringEndpoint(
			dstRid.Network,
			dstRid.ChainID,
			spec.ContractID,
			synchronization.OCR2CCIPCommit,
		),
		OffchainConfigDigester: dstProvider.OffchainConfigDigester(),
		OffchainKeyring:        kb,
		OnchainKeyring:         kb,
		MetricsRegisterer:      prometheus.WrapRegistererWith(map[string]string{"job_name": jb.Name.ValueOrZero()}, prometheus.DefaultRegisterer),
	}

	//nolint:gosec // safe to cast
	return ccipcommit.NewCommitServices(
		ctx,
		d.ds,
		srcProvider,
		dstProvider,
		jb,
		lggr,
		d.pipelineRunner,
		oracleArgsNoPlugin,
		d.isNewlyCreatedJob,
		int64(srcChainID),
		dstChainID,
		logError,
		pluginJobSpecConfig,
		d.RelayGetter,
	)
}

func newCCIPCommitPluginBytes(isSourceProvider bool, sourceStartBlock uint64, destStartBlock uint64) config.CommitPluginConfig {
	return config.CommitPluginConfig{
		IsSourceProvider: isSourceProvider,
		SourceStartBlock: sourceStartBlock,
		DestStartBlock:   destStartBlock,
	}
}

func (d *Delegate) ccipCommitGetDstProvider(ctx context.Context, jb job.Job, pluginJobSpecConfig ccipconfig.CommitPluginJobSpecConfig, transmitterID string) (types.CCIPCommitProvider, error) {
	spec := jb.OCR2OracleSpec
	if !isCCIPSupportedNetwork(spec.Relay) {
		return nil, fmt.Errorf("chain not supported for CCIP commit: %s", spec.Relay)
	}

	dstRid, err := spec.RelayID()
	if err != nil {
		return nil, ErrJobSpecNoRelayer{Err: err, PluginName: string(spec.PluginType)}
	}

	// Write PluginConfig bytes to send source/dest relayer provider + info outside of top level rargs/pargs over the wire
	dstConfigBytes, err := newCCIPCommitPluginBytes(false, pluginJobSpecConfig.SourceStartBlock, pluginJobSpecConfig.DestStartBlock).Encode()
	if err != nil {
		return nil, err
	}

	// Get provider from dest chain
	dstRelayer, err := d.RelayGetter.Get(dstRid)
	if err != nil {
		return nil, err
	}

	provider, err := dstRelayer.NewPluginProvider(ctx,
		types.RelayArgs{
			ContractID:   spec.ContractID,
			RelayConfig:  spec.RelayConfig.Bytes(),
			ProviderType: string(types.CCIPCommit),
		},
		types.PluginArgs{
			TransmitterID: transmitterID,
			PluginConfig:  dstConfigBytes,
		})
	if err != nil {
		return nil, fmt.Errorf("unable to create ccip commit provider: %w", err)
	}
	dstProvider, ok := provider.(types.CCIPCommitProvider)
	if !ok {
		return nil, errors.New("could not coerce PluginProvider to CCIPCommitProvider")
	}

	return dstProvider, nil
}

func (d *Delegate) ccipCommitGetSrcProvider(ctx context.Context, jb job.Job, pluginJobSpecConfig ccipconfig.CommitPluginJobSpecConfig, transmitterID string, dstProvider types.CCIPCommitProvider) (srcProvider types.CCIPCommitProvider, srcChainID uint64, err error) {
	spec := jb.OCR2OracleSpec
	srcConfigBytes, err := newCCIPCommitPluginBytes(true, pluginJobSpecConfig.SourceStartBlock, pluginJobSpecConfig.DestStartBlock).Encode()
	if err != nil {
		return nil, 0, err
	}
	// Use OffRampReader to get src chain ID and fetch the src relayer

	var pluginConfig ccipconfig.CommitPluginJobSpecConfig
	err = json.Unmarshal(spec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, 0, err
	}
	offRampAddress := pluginConfig.OffRamp
	offRampReader, err := dstProvider.NewOffRampReader(ctx, offRampAddress)
	if err != nil {
		return nil, 0, fmt.Errorf("create offRampReader: %w", err)
	}

	offRampConfig, err := offRampReader.GetStaticConfig(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("get offRamp static config: %w", err)
	}

	srcChainID, err = chainselectors.ChainIdFromSelector(offRampConfig.SourceChainSelector)
	if err != nil {
		return nil, 0, err
	}
	srcChainIDstr := strconv.FormatUint(srcChainID, 10)

	// Get provider from source chain
	srcRelayer, err := d.RelayGetter.Get(types.RelayID{Network: spec.Relay, ChainID: srcChainIDstr})
	if err != nil {
		return nil, 0, err
	}
	provider, err := srcRelayer.NewPluginProvider(ctx,
		types.RelayArgs{
			ContractID:   "", // Contract address only valid for dst chain
			RelayConfig:  spec.RelayConfig.Bytes(),
			ProviderType: string(types.CCIPCommit),
		},
		types.PluginArgs{
			TransmitterID: transmitterID,
			PluginConfig:  srcConfigBytes,
		})
	if err != nil {
		return nil, 0, fmt.Errorf("srcRelayer.NewPluginProvider: %w", err)
	}
	srcProvider, ok := provider.(types.CCIPCommitProvider)
	if !ok {
		return nil, 0, errors.New("could not coerce PluginProvider to CCIPCommitProvider")
	}

	return
}

func (d *Delegate) newServicesCCIPExecution(ctx context.Context, lggr logger.SugaredLogger, jb job.Job, bootstrapPeers []commontypes.BootstrapperLocator, kb ocr2key.KeyBundle, ocrDB *db, lc ocrtypes.LocalConfig, transmitterID string) ([]job.ServiceCtx, error) {
	spec := jb.OCR2OracleSpec
	if !isCCIPSupportedNetwork(spec.Relay) {
		return nil, fmt.Errorf("chain not supported for CCIP execution: %s", spec.Relay)
	}
	dstRid, err := spec.RelayID()

	if err != nil {
		return nil, ErrJobSpecNoRelayer{Err: err, PluginName: string(spec.PluginType)}
	}

	logError := func(msg string) {
		lggr.ErrorIf(d.jobORM.RecordError(context.Background(), jb.ID, msg), "unable to record error")
	}

	// PROVIDER BASED ARG CONSTRUCTION
	// Write PluginConfig bytes to send source/dest relayer provider + info outside of top level rargs/pargs over the wire
	var pluginJobSpecConfig ccipconfig.ExecPluginJobSpecConfig
	err = json.Unmarshal(spec.PluginConfig.Bytes(), &pluginJobSpecConfig)
	if err != nil {
		return nil, err
	}

	dstChainID, err := strconv.ParseInt(dstRid.ChainID, 10, 64)
	if err != nil {
		return nil, err
	}

	dstProvider, err := d.ccipExecGetDstProvider(ctx, jb, pluginJobSpecConfig, transmitterID)
	if err != nil {
		return nil, err
	}

	srcProvider, srcChainID, err := d.ccipExecGetSrcProvider(ctx, jb, pluginJobSpecConfig, transmitterID, dstProvider)
	if err != nil {
		return nil, err
	}

	oracleArgsNoPlugin2 := libocr2.OCR2OracleArgs{
		BinaryNetworkEndpointFactory: d.peerWrapper.Peer2,
		V2Bootstrappers:              bootstrapPeers,
		ContractTransmitter:          dstProvider.ContractTransmitter(),
		ContractConfigTracker:        dstProvider.ContractConfigTracker(),
		Database:                     ocrDB,
		LocalConfig:                  lc,
		MonitoringEndpoint: d.monitoringEndpointGen.GenMonitoringEndpoint(
			dstRid.Network,
			dstRid.ChainID,
			spec.ContractID,
			synchronization.OCR2CCIPExec,
		),
		OffchainConfigDigester: dstProvider.OffchainConfigDigester(),
		OffchainKeyring:        kb,
		OnchainKeyring:         kb,
		MetricsRegisterer:      prometheus.WrapRegistererWith(map[string]string{"job_name": jb.Name.ValueOrZero()}, prometheus.DefaultRegisterer),
	}

	return ccipexec.NewExecServices(ctx, lggr, jb, srcProvider, dstProvider, int64(srcChainID), dstChainID, d.isNewlyCreatedJob, oracleArgsNoPlugin2, logError)
}

func (d *Delegate) ccipExecGetDstProvider(ctx context.Context, jb job.Job, pluginJobSpecConfig ccipconfig.ExecPluginJobSpecConfig, transmitterID string) (types.CCIPExecProvider, error) {
	spec := jb.OCR2OracleSpec
	if !isCCIPSupportedNetwork(spec.Relay) {
		return nil, fmt.Errorf("chain not supported for CCIP execution: %s", spec.Relay)
	}
	dstRid, err := spec.RelayID()

	if err != nil {
		return nil, ErrJobSpecNoRelayer{Err: err, PluginName: string(spec.PluginType)}
	}

	// PROVIDER BASED ARG CONSTRUCTION
	// Write PluginConfig bytes to send source/dest relayer provider + info outside of top level rargs/pargs over the wire
	dstConfigBytes, err := newExecPluginConfig(false, pluginJobSpecConfig.SourceStartBlock, pluginJobSpecConfig.DestStartBlock, pluginJobSpecConfig.USDCConfig, pluginJobSpecConfig.LBTCConfig, string(jb.ID)).Encode()
	if err != nil {
		return nil, err
	}

	// Get provider from dest chain
	dstRelayer, err := d.RelayGetter.Get(dstRid)
	if err != nil {
		return nil, err
	}
	provider, err := dstRelayer.NewPluginProvider(ctx,
		types.RelayArgs{
			ContractID:   spec.ContractID,
			RelayConfig:  spec.RelayConfig.Bytes(),
			ProviderType: string(types.CCIPExecution),
		},
		types.PluginArgs{
			TransmitterID: transmitterID,
			PluginConfig:  dstConfigBytes,
		})
	if err != nil {
		return nil, fmt.Errorf("NewPluginProvider failed on dstRelayer: %w", err)
	}
	dstProvider, ok := provider.(types.CCIPExecProvider)
	if !ok {
		return nil, errors.New("could not coerce PluginProvider to CCIPExecProvider")
	}

	return dstProvider, nil
}

func (d *Delegate) ccipExecGetSrcProvider(ctx context.Context, jb job.Job, pluginJobSpecConfig ccipconfig.ExecPluginJobSpecConfig, transmitterID string, dstProvider types.CCIPExecProvider) (srcProvider types.CCIPExecProvider, srcChainID uint64, err error) {
	spec := jb.OCR2OracleSpec
	srcConfigBytes, err := newExecPluginConfig(true, pluginJobSpecConfig.SourceStartBlock, pluginJobSpecConfig.DestStartBlock, pluginJobSpecConfig.USDCConfig, pluginJobSpecConfig.LBTCConfig, string(jb.ID)).Encode()
	if err != nil {
		return nil, 0, err
	}

	// Use OffRampReader to get src chain ID and fetch the src relayer
	offRampAddress := cciptypes.Address(common.HexToAddress(spec.ContractID).String())
	offRampReader, err := dstProvider.NewOffRampReader(ctx, offRampAddress)
	if err != nil {
		return nil, 0, fmt.Errorf("create offRampReader: %w", err)
	}

	offRampConfig, err := offRampReader.GetStaticConfig(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("get offRamp static config: %w", err)
	}

	srcChainID, err = chainselectors.ChainIdFromSelector(offRampConfig.SourceChainSelector)
	if err != nil {
		return nil, 0, err
	}
	srcChainIDstr := strconv.FormatUint(srcChainID, 10)

	// Get provider from source chain
	srcRelayer, err := d.RelayGetter.Get(types.RelayID{Network: spec.Relay, ChainID: srcChainIDstr})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get relayer: %w", err)
	}
	provider, err := srcRelayer.NewPluginProvider(ctx,
		types.RelayArgs{
			ContractID:   "",
			RelayConfig:  spec.RelayConfig.Bytes(),
			ProviderType: string(types.CCIPExecution),
		},
		types.PluginArgs{
			TransmitterID: transmitterID,
			PluginConfig:  srcConfigBytes,
		})
	if err != nil {
		return nil, 0, err
	}
	srcProvider, ok := provider.(types.CCIPExecProvider)
	if !ok {
		return nil, 0, fmt.Errorf("could not coerce PluginProvider to CCIPExecProvider: %w", err)
	}

	return
}

func newExecPluginConfig(isSourceProvider bool, srcStartBlock uint64, dstStartBlock uint64, usdcConfig ccipconfig.USDCConfig, lbtcConfig ccipconfig.LBTCConfig, jobID string) config.ExecPluginConfig {
	return config.ExecPluginConfig{
		IsSourceProvider: isSourceProvider,
		SourceStartBlock: srcStartBlock,
		DestStartBlock:   dstStartBlock,
		USDCConfig:       usdcConfig,
		LBTCConfig:       lbtcConfig,
		JobID:            jobID,
	}
}

// errorLog implements [loop.ErrorLog]
type errorLog struct {
	jobID       int32
	recordError func(ctx context.Context, jobID int32, description string) error
}

func (l *errorLog) SaveError(ctx context.Context, msg string) error {
	return l.recordError(ctx, l.jobID, msg)
}

type logWriter struct {
	log logger.Logger
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	l.log.Debug(string(p), nil)
	n = len(p)
	return
}

func isCCIPSupportedNetwork(network string) bool {
	switch network {
	case relay.NetworkEVM, relay.NetworkTron:
		return true
	default:
		return false
	}
}
