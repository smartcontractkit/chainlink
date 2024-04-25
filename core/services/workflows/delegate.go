package workflows

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type Delegate struct {
	registry        core.CapabilitiesRegistry
	logger          logger.Logger
	legacyEVMChains legacyevm.LegacyChainContainer
}

var _ job.Delegate = (*Delegate)(nil)

func (d *Delegate) JobType() job.Type {
	return job.Workflow
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {}

func (d *Delegate) AfterJobCreated(jb job.Job) {}

func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) OnDeleteJob(context.Context, job.Job) error { return nil }

// ServicesForSpec satisfies the job.Delegate interface.
func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) ([]job.ServiceCtx, error) {
	// NOTE: we temporarily do registration inside ServicesForSpec, this will be moved out of job specs in the future
	err := targets.InitializeWrite(d.registry, d.legacyEVMChains, d.logger)
	if err != nil {
		d.logger.Errorw("could not initialize writes", err)
	}

	cfg := Config{
		Lggr:       d.logger,
		Spec:       spec.WorkflowSpec.Workflow,
		WorkflowID: spec.WorkflowSpec.WorkflowID,
		Registry:   d.registry,
	}
	engine, err := NewEngine(cfg)
	if err != nil {
		return nil, err
	}
	return []job.ServiceCtx{engine}, nil
}

func NewDelegate(logger logger.Logger, registry core.CapabilitiesRegistry, legacyEVMChains legacyevm.LegacyChainContainer) *Delegate {
	return &Delegate{logger: logger, registry: registry, legacyEVMChains: legacyEVMChains}
}

func ValidatedWorkflowSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{ExternalJobID: uuid.New()}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, fmt.Errorf("toml error on load: %w", err)
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, fmt.Errorf("toml unmarshal error on spec: %w", err)
	}

	var spec job.WorkflowSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, fmt.Errorf("toml unmarshal error on job: %w", err)
	}

	if err := spec.Validate(); err != nil {
		return jb, err
	}

	jb.WorkflowSpec = &spec
	if jb.Type != job.Workflow {
		return jb, fmt.Errorf("unsupported type %s", jb.Type)
	}

	return jb, nil
}
