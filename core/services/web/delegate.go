package web

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

type Delegate struct {
	pipelineRunner pipeline.Runner
}

func NewDelegate(pipelineRunner pipeline.Runner) *Delegate {
	return &Delegate{pipelineRunner}
}

func (d *Delegate) JobType() job.Type {
	return job.Web
}

// ServicesForSpec returns the delegate to be used for running web jobs
func (d *Delegate) ServicesForSpec(spec job.Job) (services []job.Service, err error) {
	if spec.WebSpec == nil {
		return nil, errors.Errorf("services.Delegate expects a *jobSpec.WebSpec to be present, got %v", spec)
	}

	web, err := NewFromJobSpec(spec, d.pipelineRunner)
	if err != nil {
		return nil, err
	}

	return []job.Service{web}, nil
}
