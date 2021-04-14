package cron

import (
	"github.com/smartcontractkit/chainlink/core/store"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"gorm.io/gorm"
)

type Delegate struct {
	pipelineRunner pipeline.Runner
	db             *gorm.DB
	store          store.Store
	config         Config
}

func NewDelegate(pipelineRunner pipeline.Runner, store store.Store, db *gorm.DB, config Config) *Delegate {
	return &Delegate{
		pipelineRunner: pipelineRunner,
		config:         config,
		db:             db,
		store:          store,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.CronJob
}

// ServicesForSpec returns the scheduler to be used for running cron jobs
func (d *Delegate) ServicesForSpec(spec job.Job) (services []job.Service, err error) {
	if spec.CronRequestSpec == nil {
		return nil, errors.Errorf("services.Delegate expects a *jobSpec.CronRequestSpec to be present, got %v", spec)
	}

	cron, err := NewFromJobSpec(spec, d.store, d.config, d.pipelineRunner, NewORM(d.db))
	if err != nil {
		return nil, err
	}

	return []job.Service{cron}, nil
}
