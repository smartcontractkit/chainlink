package standardcapabilities

import (
	"context"
	"crypto"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/p2pkey"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/compute"
	gatewayconnector "github.com/smartcontractkit/chainlink/v2/core/capabilities/gateway_connector"
	triggercap "github.com/smartcontractkit/chainlink/v2/core/capabilities/triggers"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	webapitarget "github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/target"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/trigger"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr/capregconfig"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/generic"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities/conversions"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type RelayGetter interface {
	Get(id types.RelayID) (loop.Relayer, error)
	GetIDToRelayerMap() map[types.RelayID]loop.Relayer
}

type Delegate struct {
	logger                  logger.Logger
	ds                      sqlutil.DataSource
	jobORM                  job.ORM
	registry                core.CapabilitiesRegistry
	cfg                     plugins.RegistrarConfig
	monitoringEndpointGen   telemetry.MonitoringEndpointGenerator
	pipelineRunner          pipeline.Runner
	relayers                RelayGetter
	gatewayConnectorWrapper *gatewayconnector.ServiceWrapper
	ks                      keystore.Master
	getPeerID               func() (p2ptypes.PeerID, error)
	ocrPeerWrapper          *ocrcommon.SingletonPeerWrapper
	newOracleFactoryFn      NewOracleFactoryFn
	computeFetcherFactoryFn compute.FetcherFactory
	selectorOpts            []func(*gateway.RoundRobinSelector)
	orgResolver             orgresolver.OrgResolver
	creSettings             core.SettingsBroadcaster
	ocrConfigService        capregconfig.OCRConfigService
	localCfg                coreconfig.LocalCapabilities

	isNewlyCreatedJob bool
}

const (
	commandOverrideForWebAPITrigger       = "__builtin_web-api-trigger"
	commandOverrideForWebAPITarget        = "__builtin_web-api-target"
	commandOverrideForCustomComputeAction = "__builtin_custom-compute-action"
)

type NewOracleFactoryFn func(generic.OracleFactoryParams) (core.OracleFactory, error)

func NewDelegate(
	logger logger.Logger,
	ds sqlutil.DataSource,
	jobORM job.ORM,
	registry core.CapabilitiesRegistry,
	cfg plugins.RegistrarConfig,
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator,
	pipelineRunner pipeline.Runner,
	relayers RelayGetter,
	gatewayConnectorWrapper *gatewayconnector.ServiceWrapper,
	ks keystore.Master,
	getPeerID func() (p2ptypes.PeerID, error),
	ocrPeerWrapper *ocrcommon.SingletonPeerWrapper,
	newOracleFactoryFn NewOracleFactoryFn,
	fetcherFactoryFn compute.FetcherFactory,
	orgResolver orgresolver.OrgResolver,
	creSettings core.SettingsBroadcaster,
	ocrConfigService capregconfig.OCRConfigService,
	localCfg coreconfig.LocalCapabilities,
	opts ...func(*gateway.RoundRobinSelector),
) *Delegate {
	return &Delegate{
		logger:                  logger,
		ds:                      ds,
		jobORM:                  jobORM,
		registry:                registry,
		cfg:                     cfg,
		monitoringEndpointGen:   monitoringEndpointGen,
		pipelineRunner:          pipelineRunner,
		relayers:                relayers,
		isNewlyCreatedJob:       false,
		gatewayConnectorWrapper: gatewayConnectorWrapper,
		ks:                      ks,
		getPeerID:               getPeerID,
		ocrPeerWrapper:          ocrPeerWrapper,
		newOracleFactoryFn:      newOracleFactoryFn,
		computeFetcherFactoryFn: fetcherFactoryFn,
		orgResolver:             orgResolver,
		creSettings:             creSettings,
		ocrConfigService:        ocrConfigService,
		localCfg:                localCfg,
		selectorOpts:            opts,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.StandardCapabilities
}

func (d *Delegate) BeforeJobCreated(job job.Job) {
	// This is only called first time the job is created
	d.isNewlyCreatedJob = true
}

func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) ([]job.ServiceCtx, error) {
	command := spec.StandardCapabilitiesSpec.Command
	configJSON := spec.StandardCapabilitiesSpec.Config

	if d.localCfg != nil {
		capabilityID := conversions.GetCapabilityIDFromCommand(command, configJSON)
		if capabilityID != "" && d.localCfg.IsAllowlisted(capabilityID) {
			return nil, fmt.Errorf(
				"capability %q is in the RegistryBasedLaunchAllowlist and will be started from the on-chain registry; "+
					"remove the job spec and let the LocalCapabilityManager handle it via [Capabilities.Local] TOML config",
				capabilityID,
			)
		}
	}

	// Job-spec boot path: capability DON ID is not carried in the spec, so the
	// host best-effort resolves it from the capability registry inside NewServices.
	// On a node that belongs to multiple DONs running the same capability, the
	// registry lookup cannot disambiguate which DON this plugin serves, so it
	// resolves to 0 and the plugin falls back to the consumer workflow's DON ID
	// for event labeling. Carrying the DON ID in the job spec would close that
	// gap; tracked as a follow-up. See CRE-4409.
	return d.NewServices(ctx, command, configJSON, spec.ID, spec.Name.ValueOrZero(), spec.ExternalJobID, spec.StandardCapabilitiesSpec.OracleFactory, 0)
}

// NewServices builds the per-job services for a Standard Capabilities LOOP.
//
// capabilityDonID is the on-chain DON ID this plugin process is being spawned
// for, when known by the caller. Pass 0 from the job-spec boot path (NewServices
// will resolve it from the capability registry); pass a nonzero value from the
// LocalCapabilityManager boot path where the (capID, donID) pairing is already
// known. The resolved value is plumbed to the plugin via
// StandardCapabilitiesDependencies.CapabilityDonID at Initialise time.
func (d *Delegate) NewServices(
	ctx context.Context,
	command string,
	configJSON string,
	jobID int32,
	jobName string,
	externalJobID uuid.UUID,
	oracleFactoryConfig job.OracleFactoryConfig,
	capabilityDonID uint32,
) ([]job.ServiceCtx, error) {
	log := d.logger.Named("StandardCapabilities").Named(strconv.Itoa(int(jobID))).Named(jobName)

	kvStore := job.NewKVStore(jobID, d.ds)

	// Enable signing and decryption for the capability, if available.
	var ks core.Keystore
	var decrypter core.Decrypter
	var signer crypto.Signer
	if d.ks.Workflow() != nil {
		workflowKeys, err := d.ks.Workflow().GetAll()
		if err != nil {
			return nil, fmt.Errorf("failed to get workflow keys: %w", err)
		}
		if len(workflowKeys) > 0 {
			decrypter = &workflowKeys[0]
		}
	}
	if d.ks.P2P() != nil && d.getPeerID != nil {
		peerID, err := d.getPeerID()
		if err != nil {
			log.Warnw("getPeerID() failed, will extract default peerID from Keystore", "error", err)
		}
		p2pKey, err := d.ks.P2P().GetOrFirst(p2pkey.PeerID(peerID))
		if err != nil {
			return nil, fmt.Errorf("external peer wrapper does not pertain to a valid P2P key %x: %w", peerID, err)
		}
		signer = p2pKey
	}
	ks, err := core.NewSignerDecrypter(core.StandardCapabilityAccount, signer, decrypter)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer decrypter: %w", err)
	}

	relayerSet, err := generic.NewRelayerSet(d.relayers, externalJobID, jobID, d.isNewlyCreatedJob)
	if err != nil {
		return nil, fmt.Errorf("failed to create relayer set: %w", err)
	}

	ocrEvmKeyBundles, err := d.ks.OCR2().GetAllOfType(corekeys.EVM)
	if err != nil {
		return nil, err
	}

	var ocrEvmKeyBundle ocr2key.KeyBundle
	if len(ocrEvmKeyBundles) == 0 {
		ocrEvmKeyBundle, err = d.ks.OCR2().Create(ctx, corekeys.EVM)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create OCR key bundle")
		}
	} else {
		if len(ocrEvmKeyBundles) > 1 {
			log.Infof("found %d EVM OCR key bundles, which may cause unexpected behavior if using the OracleFactory", len(ocrEvmKeyBundles))
		}
		ocrEvmKeyBundle = ocrEvmKeyBundles[0]
	}

	capabilityID := conversions.GetCapabilityIDFromCommand(command, configJSON)
	if d.ocrConfigService != nil && capabilityID == "" {
		log.Warnw("No capability ID mapping for command, using legacy config only",
			"command", command)
	}

	// Best-effort resolve the authoritative capability DON ID for this plugin
	// process when the caller did not provide one (job-spec boot path). The LOOP
	// uses this to label emitted events with the sending DON ID rather than the
	// consumer workflow's DON. Resolution can return 0 (e.g. on a node that
	// belongs to multiple DONs for this capability); in that case the plugin
	// falls back to the workflow DON ID. See CRE-4409.
	if capabilityDonID == 0 && capabilityID != "" {
		capabilityDonID = resolveCapabilityDonID(ctx, log, d.registry, d.getPeerID, capabilityID)
		log.Debugw("Resolved capability DON ID from registry", "capabilityID", capabilityID, "donID", capabilityDonID)
	}

	var oracleFactory core.OracleFactory
	// NOTE: special case for custom Oracle Factory for use in tests
	if d.newOracleFactoryFn != nil {
		oracleFactory, err = d.newOracleFactoryFn(generic.OracleFactoryParams{
			Logger:           log,
			JobORM:           d.jobORM,
			JobID:            jobID,
			JobName:          jobName,
			KB:               ocrEvmKeyBundle,
			Config:           oracleFactoryConfig,
			PeerWrapper:      d.ocrPeerWrapper,
			RelayerSet:       relayerSet,
			OCRConfigService: d.ocrConfigService,
			CapabilityID:     capabilityID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create oracle factory from function: %w", err)
		}
	} else {
		log.Debug("oracleFactoryConfig: ", oracleFactoryConfig)

		if oracleFactoryConfig.Enabled && d.ocrPeerWrapper == nil {
			return nil, errors.New("P2P stack required for Oracle Factory")
		}

		oracleFactory, err = generic.NewOracleFactory(generic.OracleFactoryParams{
			Logger:                 log,
			JobORM:                 d.jobORM,
			JobID:                  jobID,
			JobName:                jobName,
			KB:                     ocrEvmKeyBundle,
			Config:                 oracleFactoryConfig,
			OnchainSigningStrategy: oracleFactoryConfig.OnchainSigning,
			PeerWrapper:            d.ocrPeerWrapper,
			RelayerSet:             relayerSet,
			OcrKeystore:            d.ks.OCR2(),
			EthKeystore:            d.ks.Eth(),
			OCRConfigService:       d.ocrConfigService,
			CapabilityID:           capabilityID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create oracle factory: %w", err)
		}
	}
	var connector connector.GatewayConnector
	if d.gatewayConnectorWrapper != nil {
		connector = d.gatewayConnectorWrapper.GetGatewayConnector()
	}

	// NOTE: special cases for built-in capabilities (to be moved into LOOPPs in the future)
	if command == commandOverrideForWebAPITrigger {
		if d.gatewayConnectorWrapper == nil || connector == nil {
			return nil, errors.New("gateway connector is required for web API Trigger capability")
		}
		triggerSrvc, err := trigger.NewTrigger(configJSON, d.registry, connector, log)
		if err != nil {
			return nil, fmt.Errorf("failed to create a Web API Trigger service: %w", err)
		}
		return []job.ServiceCtx{triggerSrvc}, nil
	}

	if command == commandOverrideForWebAPITarget {
		if d.gatewayConnectorWrapper == nil || connector == nil {
			return nil, errors.New("gateway connector is required for web API Target capability")
		}
		if len(configJSON) == 0 {
			return nil, errors.New("config is empty")
		}
		var targetCfg webapi.ServiceConfig
		err := toml.Unmarshal([]byte(configJSON), &targetCfg)
		if err != nil {
			return nil, err
		}
		lggr := d.logger.Named("WebAPITarget")
		handler, err := webapi.NewOutgoingConnectorHandler(connector, targetCfg, capabilities.MethodWebAPITarget, lggr, d.selectorOpts...)
		if err != nil {
			return nil, err
		}
		capability, err := webapitarget.NewCapability(targetCfg, d.registry, handler, lggr)
		if err != nil {
			return nil, err
		}
		return []job.ServiceCtx{capability, handler}, nil
	}

	if command == commandOverrideForCustomComputeAction {
		var fetcherFactoryFn compute.FetcherFactory
		var services []job.ServiceCtx
		var cfg compute.Config

		tomlErr := toml.Unmarshal([]byte(configJSON), &cfg)
		if tomlErr != nil {
			return nil, tomlErr
		}

		if d.computeFetcherFactoryFn != nil {
			fetcherFactoryFn = d.computeFetcherFactoryFn
		} else {
			if d.gatewayConnectorWrapper == nil || connector == nil {
				return nil, errors.New("gateway connector is required for custom compute capability")
			}

			lggr := d.logger.Named("ComputeAction")

			handler, err := webapi.NewOutgoingConnectorHandler(connector, cfg.ServiceConfig, capabilities.MethodComputeAction, lggr, d.selectorOpts...)
			if err != nil {
				return nil, err
			}
			services = append(services, handler)

			idGeneratorFn := func() string {
				return uuid.New().String()
			}

			fetcherFactoryFn, err = compute.NewOutgoingConnectorFetcherFactory(handler, idGeneratorFn)
			if err != nil {
				return nil, fmt.Errorf("failed to create fetcher factory: %w", err)
			}
		}

		if len(configJSON) == 0 {
			return nil, errors.New("config is empty")
		}

		computeSrvc, err := compute.NewAction(cfg, log, d.registry, fetcherFactoryFn)
		if err != nil {
			return nil, err
		}
		services = append(services, computeSrvc)

		return services, nil
	}

	dependencies := core.StandardCapabilitiesDependencies{
		Config:             configJSON,
		Store:              kvStore,
		CapabilityRegistry: d.registry,
		RelayerSet:         relayerSet,
		OracleFactory:      oracleFactory,
		GatewayConnector:   connector,
		P2PKeystore:        ks,
		OrgResolver:        d.orgResolver,
		CRESettings:        d.creSettings,
		TriggerEventStore:  triggercap.NewTriggerEventStore(d.ds),
		CapabilityDonID:    capabilityDonID,
	}
	standardCapability := NewStandardCapabilities(log, command, configJSON, d.cfg, dependencies)

	return []job.ServiceCtx{standardCapability}, nil
}

// resolveCapabilityDonID best-effort resolves the on-chain DON ID this node is
// running the given capability for, by looking up DONsForCapability and filtering
// by the local node's PeerID.
//
// This is always best-effort: it never returns an error. Any failure — including
// infrastructure issues like getPeerID failing or the registry being unavailable —
// results in returning 0 with a warning logged. The caller then falls back to
// labeling events with the consumer workflow's DON ID. See CRE-4409.
func resolveCapabilityDonID(ctx context.Context, lggr logger.Logger, registry core.CapabilitiesRegistry, getPeerID func() (p2ptypes.PeerID, error), capabilityID string) uint32 {
	if registry == nil {
		lggr.Warnw("Capabilities registry is nil; falling back to workflow DON ID for event labeling", "capabilityID", capabilityID)
		return 0
	}
	if getPeerID == nil {
		lggr.Warnw("getPeerID is nil; falling back to workflow DON ID for event labeling", "capabilityID", capabilityID)
		return 0
	}
	peerID, err := getPeerID()
	if err != nil {
		lggr.Warnw("Failed to get local peer ID; falling back to workflow DON ID for event labeling", "capabilityID", capabilityID, "err", err)
		return 0
	}
	dwns, err := registry.DONsForCapability(ctx, capabilityID)
	if err != nil {
		lggr.Warnw("DONsForCapability failed; falling back to workflow DON ID for event labeling", "capabilityID", capabilityID, "err", err)
		return 0
	}
	var matched []uint32
	for _, d := range dwns {
		for _, n := range d.Nodes {
			if n.PeerID != nil && *n.PeerID == peerID {
				matched = append(matched, d.DON.ID)
				break
			}
		}
	}
	switch len(matched) {
	case 1:
		return matched[0]
	case 0:
		lggr.Warnw("No DON found for local peer on capability; falling back to workflow DON ID for event labeling",
			"peerID", peerID, "capabilityID", capabilityID)
		return 0
	default:
		lggr.Warnw("Local peer belongs to multiple DONs for capability; cannot disambiguate, falling back to workflow DON ID for event labeling",
			"peerID", peerID, "capabilityID", capabilityID, "matched", matched)
		return 0
	}
}

func (d *Delegate) AfterJobCreated(job job.Job) {}

func (d *Delegate) BeforeJobDeleted(job job.Job) {}

func (d *Delegate) OnDeleteJob(ctx context.Context, jb job.Job) error { return nil }

func ValidatedStandardCapabilitiesSpec(tomlString string) (job.Job, error) {
	jb := job.Job{ExternalJobID: uuid.New()}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, errors.Wrap(err, "toml error on load standard capabilities")
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on standard capabilities spec")
	}

	var spec job.StandardCapabilitiesSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on standard capabilities job")
	}

	jb.StandardCapabilitiesSpec = &spec
	if jb.Type != job.StandardCapabilities {
		return jb, errors.Errorf("standard capabilities unsupported job type %s", jb.Type)
	}

	if len(jb.StandardCapabilitiesSpec.Command) == 0 {
		return jb, errors.Errorf("standard capabilities command must be set")
	}

	// Skip validation if Oracle Factory is not enabled
	if !jb.StandardCapabilitiesSpec.OracleFactory.Enabled {
		return jb, nil
	}

	// If Oracle Factory is enabled, it must have at least one bootstrap peer
	if len(jb.StandardCapabilitiesSpec.OracleFactory.BootstrapPeers) == 0 {
		return jb, errors.New("no bootstrap peers found")
	}

	// Validate bootstrap peers
	_, err = ocrcommon.ParseBootstrapPeers(jb.StandardCapabilitiesSpec.OracleFactory.BootstrapPeers)
	if err != nil {
		return jb, errors.Wrap(err, "failed to parse bootstrap peers")
	}

	return jb, nil
}

type ErrorLog struct {
	jobID       int32
	recordError func(ctx context.Context, jobID int32, description string) error
}

func (l *ErrorLog) SaveError(ctx context.Context, msg string) error {
	return l.recordError(ctx, l.jobID, msg)
}
