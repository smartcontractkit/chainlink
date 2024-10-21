package standardcapabilities

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/compute"
	gatewayconnector "github.com/smartcontractkit/chainlink/v2/core/capabilities/gateway_connector"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	webapitarget "github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/target"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/trigger"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/generic"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type RelayGetter interface {
	Get(id types.RelayID) (loop.Relayer, error)
	GetIDToRelayerMap() (map[types.RelayID]loop.Relayer, error)
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
	peerWrapper             *ocrcommon.SingletonPeerWrapper

	isNewlyCreatedJob bool
}

const (
	commandOverrideForWebAPITrigger       = "__builtin_web-api-trigger"
	commandOverrideForWebAPITarget        = "__builtin_web-api-target"
	commandOverrideForCustomComputeAction = "__builtin_custom-compute-action"
)

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
	peerWrapper *ocrcommon.SingletonPeerWrapper,
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
		peerWrapper:             peerWrapper,
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
	log := d.logger.Named("StandardCapabilities").Named(spec.StandardCapabilitiesSpec.GetID())

	kvStore := job.NewKVStore(spec.ID, d.ds, log)
	telemetryService := generic.NewTelemetryAdapter(d.monitoringEndpointGen)
	errorLog := &ErrorLog{jobID: spec.ID, recordError: d.jobORM.RecordError}
	pr := generic.NewPipelineRunnerAdapter(log, spec, d.pipelineRunner)

	relayerSet, err := generic.NewRelayerSet(d.relayers, spec.ExternalJobID, spec.ID, d.isNewlyCreatedJob)
	if err != nil {
		return nil, fmt.Errorf("failed to create relayer set: %w", err)
	}

	ocrKeyBundles, err := d.ks.OCR2().GetAll()
	if err != nil {
		return nil, err
	}

	if len(ocrKeyBundles) > 1 {
		return nil, fmt.Errorf("expected exactly one OCR key bundle, but found: %d", len(ocrKeyBundles))
	}

	var ocrKeyBundle ocr2key.KeyBundle
	if len(ocrKeyBundles) == 0 {
		ocrKeyBundle, err = d.ks.OCR2().Create(ctx, chaintype.EVM)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create OCR key bundle")
		}
	} else {
		ocrKeyBundle = ocrKeyBundles[0]
	}

	ethKeyBundles, err := d.ks.Eth().GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if len(ethKeyBundles) > 1 {
		return nil, fmt.Errorf("expected exactly one ETH key bundle, but found: %d", len(ethKeyBundles))
	}

	var ethKeyBundle ethkey.KeyV2
	if len(ethKeyBundles) == 0 {
		ethKeyBundle, err = d.ks.Eth().Create(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create ETH key bundle")
		}
	} else {
		ethKeyBundle = ethKeyBundles[0]
	}

	log.Debug("oracleFactoryConfig: ", spec.StandardCapabilitiesSpec.OracleFactory)

	if spec.StandardCapabilitiesSpec.OracleFactory.Enabled && d.peerWrapper == nil {
		return nil, errors.New("P2P stack required for Oracle Factory")
	}

	oracleFactory, err := generic.NewOracleFactory(generic.OracleFactoryParams{
		Logger:        log,
		JobORM:        d.jobORM,
		JobID:         spec.ID,
		JobName:       spec.Name.ValueOrZero(),
		KB:            ocrKeyBundle,
		Config:        spec.StandardCapabilitiesSpec.OracleFactory,
		PeerWrapper:   d.peerWrapper,
		RelayerSet:    relayerSet,
		TransmitterID: ethKeyBundle.Address.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create oracle factory: %w", err)
	}
	// NOTE: special cases for built-in capabilities (to be moved into LOOPPs in the future)
	if spec.StandardCapabilitiesSpec.Command == commandOverrideForWebAPITrigger {
		if d.gatewayConnectorWrapper == nil {
			return nil, errors.New("gateway connector is required for web API Trigger capability")
		}
		connector := d.gatewayConnectorWrapper.GetGatewayConnector()
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
		connector := d.gatewayConnectorWrapper.GetGatewayConnector()
		if len(spec.StandardCapabilitiesSpec.Config) == 0 {
			return nil, errors.New("config is empty")
		}
		var targetCfg webapi.Config
		err := toml.Unmarshal([]byte(spec.StandardCapabilitiesSpec.Config), &targetCfg)
		if err != nil {
			return nil, err
		}
		lggr := d.logger.Named("WebAPITarget")
		handler, err := webapi.NewOutgoingConnectorHandler(connector, targetCfg, capabilities.MethodWebAPITarget, lggr)
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
		if d.gatewayConnectorWrapper == nil {
			return nil, errors.New("gateway connector is required for custom compute capability")
		}

		if len(spec.StandardCapabilitiesSpec.Config) == 0 {
			return nil, errors.New("config is empty")
		}

		var fetchCfg webapi.Config
		err := toml.Unmarshal([]byte(spec.StandardCapabilitiesSpec.Config), &fetchCfg)
		if err != nil {
			return nil, err
		}
		lggr := d.logger.Named("ComputeAction")

		handler, err := webapi.NewOutgoingConnectorHandler(d.gatewayConnectorWrapper.GetGatewayConnector(), fetchCfg, capabilities.MethodComputeAction, lggr)
		if err != nil {
			return nil, err
		}

		idGeneratorFn := func() string {
			return uuid.New().String()
		}

		computeSrvc := compute.NewAction(fetchCfg, log, d.registry, handler, idGeneratorFn)
		return []job.ServiceCtx{computeSrvc}, nil
	}

	standardCapability := newStandardCapabilities(log, spec.StandardCapabilitiesSpec, d.cfg, telemetryService, kvStore, d.registry, errorLog,
		pr, relayerSet, oracleFactory)

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
