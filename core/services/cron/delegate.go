package cron

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

type Delegate struct {
	pipelineRunner pipeline.Runner
	lggr           logger.Logger
}

var _ job.Delegate = (*Delegate)(nil)

func NewDelegate(pipelineRunner pipeline.Runner, lggr logger.Logger) *Delegate {
	return &Delegate{
		pipelineRunner: pipelineRunner,
		lggr:           lggr,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.Cron
}

func (Delegate) AfterJobCreated(spec job.Job)  {}
func (Delegate) BeforeJobDeleted(spec job.Job) {}

// ServicesForSpec returns the scheduler to be used for running cron jobs
func (d *Delegate) ServicesForSpec(spec job.Job) (services []job.ServiceCtx, err error) {
	if spec.CronSpec == nil {
		return nil, errors.Errorf("services.Delegate expects a *jobSpec.CronSpec to be present, got %v", spec)
	}

	cron, err := NewCronFromJobSpec(spec, d.pipelineRunner, d.lggr)
	if err != nil {
		return nil, err
	}

	return []job.ServiceCtx{cron}, nil
}
