package workflows

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
)

type Delegate struct {
	registry core.CapabilitiesRegistry
	logger   logger.Logger
	store    store.Store
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
func (d *Delegate) ServicesForSpec(_ context.Context, spec job.Job) ([]job.ServiceCtx, error) {
	cfg := Config{
		Lggr:          d.logger,
		Spec:          spec.WorkflowSpec.Workflow,
		WorkflowID:    spec.WorkflowSpec.WorkflowID,
		WorkflowOwner: spec.WorkflowSpec.WorkflowOwner,
		WorkflowName:  spec.WorkflowSpec.WorkflowName,
		Registry:      d.registry,
		Store:         d.store,
	}
	engine, err := NewEngine(cfg)
	if err != nil {
		return nil, err
	}
	return []job.ServiceCtx{engine}, nil
}

func NewDelegate(
	logger logger.Logger,
	registry core.CapabilitiesRegistry,
	store store.Store,
) *Delegate {
	return &Delegate{logger: logger, registry: registry, store: store}
}

func ValidatedWorkflowJobSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{ExternalJobID: uuid.New()}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, fmt.Errorf("toml error on load: %w", err)
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, fmt.Errorf("toml unmarshal error on spec: %w", err)
	}
	if jb.Type != job.Workflow {
		return jb, fmt.Errorf("unsupported type %s, expected %s", jb.Type, job.Workflow)
	}

	var spec job.WorkflowSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, fmt.Errorf("toml unmarshal error on workflow spec: %w", err)
	}

	err = spec.Validate()
	if err != nil {
		return jb, fmt.Errorf("invalid WorkflowSpec: %w", err)
	}

	// ensure the embedded workflow graph is valid
	_, err = Parse(spec.Workflow)
	if err != nil {
		return jb, fmt.Errorf("failed to parse workflow graph: %w", err)
	}
	jb.WorkflowSpec = &spec
	jb.WorkflowSpecID = &spec.ID

	return jb, nil
}
