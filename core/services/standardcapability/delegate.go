package standardcapability

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
	logger                logger.Logger
	ds                    sqlutil.DataSource
	jobORM                job.ORM
	registry              core.CapabilitiesRegistry
	cfg                   plugins.RegistrarConfig
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	pipelineRunner        pipeline.Runner
	relayers              RelayGetter

	isNewlyCreatedJob bool
}

func NewDelegate(logger logger.Logger, ds sqlutil.DataSource, jobORM job.ORM, registry core.CapabilitiesRegistry,
	cfg plugins.RegistrarConfig, monitoringEndpointGen telemetry.MonitoringEndpointGenerator, pipelineRunner pipeline.Runner,
	relayers RelayGetter) *Delegate {
	return &Delegate{logger: logger, ds: ds, jobORM: jobORM, registry: registry, cfg: cfg, monitoringEndpointGen: monitoringEndpointGen, pipelineRunner: pipelineRunner,
		relayers: relayers, isNewlyCreatedJob: false}
}

func (d *Delegate) JobType() job.Type {
	return job.StandardCapability
}

func (d *Delegate) BeforeJobCreated(job job.Job) {
	// This is only called first time the job is created
	d.isNewlyCreatedJob = true
}

func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) ([]job.ServiceCtx, error) {
	log := d.logger.Named("StandardCapability").Named(spec.StandardCapabilitySpec.GetID())

	kvStore := job.NewKVStore(spec.ID, d.ds, log)
	telemetryService := generic.NewTelemetryAdapter(d.monitoringEndpointGen)
	errorLog := &ErrorLog{jobID: spec.ID, recordError: d.jobORM.RecordError}
	pr := generic.NewPipelineRunnerAdapter(log, spec, d.pipelineRunner)

	relayerSet, err := generic.NewRelayerSet(d.relayers, spec.ExternalJobID, spec.ID, d.isNewlyCreatedJob)
	if err != nil {
		return nil, fmt.Errorf("failed to create relayer set: %w", err)
	}

	standardCapability := NewStandardCapability(log, spec.StandardCapabilitySpec, d.cfg, telemetryService, kvStore, d.registry, errorLog,
		pr, relayerSet)

	return []job.ServiceCtx{standardCapability}, nil
}

func (d *Delegate) AfterJobCreated(job job.Job) {}

func (d *Delegate) BeforeJobDeleted(job job.Job) {}

func (d *Delegate) OnDeleteJob(ctx context.Context, jb job.Job) error { return nil }

func ValidatedStandardCapabilitySpec(tomlString string) (job.Job, error) {
	var jb = job.Job{ExternalJobID: uuid.New()}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, errors.Wrap(err, "toml error on load standard capability")
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on standard capability spec")
	}

	var spec job.StandardCapabilitySpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on standard capability job")
	}

	jb.StandardCapabilitySpec = &spec
	if jb.Type != job.StandardCapability {
		return jb, errors.Errorf("standard capability unsupported job type %s", jb.Type)
	}

	if len(jb.StandardCapabilitySpec.Command) == 0 {
		return jb, errors.Errorf("standard capability command must be set")
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
