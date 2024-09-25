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
	gatewayconnector "github.com/smartcontractkit/chainlink/v2/core/capabilities/gateway_connector"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/generic"
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

	isNewlyCreatedJob bool
}

const (
	commandOverrideForWebAPITrigger = "__builtin_web-api-trigger"
)

func NewDelegate(logger logger.Logger, ds sqlutil.DataSource, jobORM job.ORM, registry core.CapabilitiesRegistry,
	cfg plugins.RegistrarConfig, monitoringEndpointGen telemetry.MonitoringEndpointGenerator, pipelineRunner pipeline.Runner,
	relayers RelayGetter, gatewayConnectorWrapper *gatewayconnector.ServiceWrapper) *Delegate {
	return &Delegate{logger: logger, ds: ds, jobORM: jobORM, registry: registry, cfg: cfg, monitoringEndpointGen: monitoringEndpointGen, pipelineRunner: pipelineRunner,
		relayers: relayers, isNewlyCreatedJob: false, gatewayConnectorWrapper: gatewayConnectorWrapper}
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

	// NOTE: special cases for built-in capabilities (to be moved into LOOPPs in the future)
	if spec.StandardCapabilitiesSpec.Command == commandOverrideForWebAPITrigger {
		if d.gatewayConnectorWrapper == nil {
			return nil, errors.New("gateway connector is required for web API Trigger capability")
		}
		connector := d.gatewayConnectorWrapper.GetGatewayConnector()
		triggerSrvc, err := webapi.NewTrigger(spec.StandardCapabilitiesSpec.Config, d.registry, connector, log)
		if err != nil {
			return nil, fmt.Errorf("failed to create a Web API Trigger service: %w", err)
		}
		return []job.ServiceCtx{triggerSrvc}, nil
	}

	standardCapability := newStandardCapabilities(log, spec.StandardCapabilitiesSpec, d.cfg, telemetryService, kvStore, d.registry, errorLog,
		pr, relayerSet)

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

	return jb, nil
}

type ErrorLog struct {
	jobID       int32
	recordError func(ctx context.Context, jobID int32, description string) error
}

func (l *ErrorLog) SaveError(ctx context.Context, msg string) error {
	return l.recordError(ctx, l.jobID, msg)
}
