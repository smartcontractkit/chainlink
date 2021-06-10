package cron

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

type Delegate struct {
	pipelineRunner pipeline.Runner
}

var _ job.Delegate = (*Delegate)(nil)

func NewDelegate(pipelineRunner pipeline.Runner) *Delegate {
	return &Delegate{
		pipelineRunner: pipelineRunner,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.Cron
}

func (Delegate) OnJobCreated(spec job.Job) {}
func (Delegate) OnJobDeleted(spec job.Job) {}

// ServicesForSpec returns the scheduler to be used for running cron jobs
func (d *Delegate) ServicesForSpec(spec job.Job) (services []job.Service, err error) {
	if spec.CronSpec == nil {
		return nil, errors.Errorf("services.Delegate expects a *jobSpec.CronSpec to be present, got %v", spec)
	}

	cron, err := NewCronFromJobSpec(spec, d.pipelineRunner)
	if err != nil {
		return nil, err
	}

	return []job.Service{cron}, nil
}
