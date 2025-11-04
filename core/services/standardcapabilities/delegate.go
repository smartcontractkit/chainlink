package standardcapabilities

import (
	"context"
	"crypto"
	"fmt"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/compute"
	gatewayconnector "github.com/smartcontractkit/chainlink/v2/core/capabilities/gateway_connector"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	webapitarget "github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/target"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/trigger"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/generic"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
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
	log := d.logger.Named("StandardCapabilities").Named(spec.StandardCapabilitiesSpec.GetID()).Named(spec.Name.ValueOrZero())

	kvStore := job.NewKVStore(spec.ID, d.ds)

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

	telemetryService := generic.NewTelemetryAdapter(d.monitoringEndpointGen)
	errorLog := &ErrorLog{jobID: spec.ID, recordError: d.jobORM.RecordError}
	pr := generic.NewPipelineRunnerAdapter(log, spec, d.pipelineRunner)

	relayerSet, err := generic.NewRelayerSet(d.relayers, spec.ExternalJobID, spec.ID, d.isNewlyCreatedJob)
	if err != nil {
		return nil, fmt.Errorf("failed to create relayer set: %w", err)
	}

	ocrEvmKeyBundles, err := d.ks.OCR2().GetAllOfType(chaintype.EVM)
	if err != nil {
		return nil, err
	}

	var ocrEvmKeyBundle ocr2key.KeyBundle
	if len(ocrEvmKeyBundles) == 0 {
		ocrEvmKeyBundle, err = d.ks.OCR2().Create(ctx, chaintype.EVM)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create OCR key bundle")
		}
	} else {
		if len(ocrEvmKeyBundles) > 1 {
			log.Infof("found %d EVM OCR key bundles, which may cause unexpected behavior if using the OracleFactory", len(ocrEvmKeyBundles))
		}
		ocrEvmKeyBundle = ocrEvmKeyBundles[0]
	}

	var oracleFactory core.OracleFactory
	// NOTE: special case for custom Oracle Factory for use in tests
	if d.newOracleFactoryFn != nil {
		oracleFactory, err = d.newOracleFactoryFn(generic.OracleFactoryParams{
			Logger:      log,
			JobORM:      d.jobORM,
			JobID:       spec.ID,
			JobName:     spec.Name.ValueOrZero(),
			KB:          ocrEvmKeyBundle,
			Config:      spec.StandardCapabilitiesSpec.OracleFactory,
			PeerWrapper: d.ocrPeerWrapper,
			RelayerSet:  relayerSet,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create oracle factory from function: %w", err)
		}
	} else {
		log.Debug("oracleFactoryConfig: ", spec.StandardCapabilitiesSpec.OracleFactory)

		if spec.StandardCapabilitiesSpec.OracleFactory.Enabled && d.ocrPeerWrapper == nil {
			return nil, errors.New("P2P stack required for Oracle Factory")
		}

		oracleFactory, err = generic.NewOracleFactory(generic.OracleFactoryParams{
			Logger:                 log,
			JobORM:                 d.jobORM,
			JobID:                  spec.ID,
			JobName:                spec.Name.ValueOrZero(),
			KB:                     ocrEvmKeyBundle,
			Config:                 spec.StandardCapabilitiesSpec.OracleFactory,
			OnchainSigningStrategy: spec.StandardCapabilitiesSpec.OracleFactory.OnchainSigning,
			PeerWrapper:            d.ocrPeerWrapper,
			RelayerSet:             relayerSet,
			OcrKeystore:            d.ks.OCR2(),
			EthKeystore:            d.ks.Eth(),
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
	if spec.StandardCapabilitiesSpec.Command == commandOverrideForWebAPITrigger {
		if d.gatewayConnectorWrapper == nil {
			return nil, errors.New("gateway connector is required for web API Trigger capability")
		}
		triggerSrvc, err := trigger.NewTrigger(spec.StandardCapabilitiesSpec.Config, d.registry, connector, log)
		if err != nil {
			return nil, fmt.Errorf("failed to create a Web API Trigger service: %w", err)
		}
		return []job.ServiceCtx{triggerSrvc}, nil
	}

	if spec.StandardCapabilitiesSpec.Command == commandOverrideForWebAPITarget {
		if d.gatewayConnectorWrapper == nil {
			return nil, errors.New("gateway connector is required for web API Target capability")
		}
		if len(spec.StandardCapabilitiesSpec.Config) == 0 {
			return nil, errors.New("config is empty")
		}
		var targetCfg webapi.ServiceConfig
		err := toml.Unmarshal([]byte(spec.StandardCapabilitiesSpec.Config), &targetCfg)
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

	if spec.StandardCapabilitiesSpec.Command == commandOverrideForCustomComputeAction {
		var fetcherFactoryFn compute.FetcherFactory
		var services []job.ServiceCtx
		var cfg compute.Config

		tomlErr := toml.Unmarshal([]byte(spec.StandardCapabilitiesSpec.Config), &cfg)
		if tomlErr != nil {
			return nil, tomlErr
		}

		if d.computeFetcherFactoryFn != nil {
			fetcherFactoryFn = d.computeFetcherFactoryFn
		} else {
			if d.gatewayConnectorWrapper == nil {
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

		if len(spec.StandardCapabilitiesSpec.Config) == 0 {
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
		Config:             spec.StandardCapabilitiesSpec.Config,
		TelemetryService:   telemetryService,
		Store:              kvStore,
		CapabilityRegistry: d.registry,
		ErrorLog:           errorLog,
		PipelineRunner:     pr,
		RelayerSet:         relayerSet,
		OracleFactory:      oracleFactory,
		GatewayConnector:   connector,
		P2PKeystore:        ks,
		OrgResolver:        d.orgResolver,
	}
	standardCapability := NewStandardCapabilities(log, spec.StandardCapabilitiesSpec, d.cfg, dependencies)

	return []job.ServiceCtx{standardCapability}, nil
}

func (d *Delegate) AfterJobCreated(job job.Job) {}

func (d *Delegate) BeforeJobDeleted(job job.Job) {}

func (d *Delegate) OnDeleteJob(ctx context.Context, jb job.Job) error { return nil }

func ValidatedStandardCapabilitiesSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{ExternalJobID: uuid.New()}

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
