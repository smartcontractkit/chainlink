package vrf

import (
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

type Delegate struct {
	vorm ORM
	pr   pipeline.Runner
	porm pipeline.ORM
}

func NewDelegate(vorm ORM, pr pipeline.Runner, porm pipeline.ORM) *Delegate {
	return &Delegate{
		vorm: vorm,
		pr:   pr,
		porm: porm,
	}
}

// JobType implements the job.Delegate interface
func (d *Delegate) JobType() job.Type {
	return job.VRF
}

// ServicesForSpec returns the flux monitor service for the job spec
func (d *Delegate) ServicesForSpec(spec job.Job) ([]job.Service, error) {
	// TODO
	return []job.Service{}, nil
}
