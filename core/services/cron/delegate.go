package cron

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
)

type Delegate struct{}

func NewDelegate() *Delegate {
	return &Delegate{}
}

func (d *Delegate) JobType() job.Type {
	return job.CronJob
}

// ServicesForSpec returns the scheduler to be used for running cron jobs
func (d *Delegate) ServicesForSpec(spec job.Job) (services []job.Service, err error) {
	if spec.CronRequestSpec == nil {
		return nil, errors.Errorf("services.Delegate expects a *job.CronRequestSpec to be present, got %v", spec)
	}

	// TODO: do we need the scheduler service?
	// services = append(services, scheduler)
	cron, err := NewJobFromSpec(spec)
	if err != nil {
		return nil, err
	}

	return []job.Service{cron}, nil
}

// Start implements the job.Service interface.
func (cron *CronJob) Start() error {
	cron.logger.Debug("Starting cron job")

	// TODO: do stuff..., run cron schedule? run schedule? (ref ../scheduler.go)

	return nil
}

// Close implements the job.Service interface. It stops this instance from
// polling, cleaning up resources.
func (cron *CronJob) Close() error {

	// TODO: cleanup any cron stuff? maybe not needed? CLeanup scheduler? (ref ../scheduler.go)

	return nil
}
