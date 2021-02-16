package fluxmonitorv2

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"gorm.io/gorm"
)

type Delegate struct {
	pipelineRunner pipeline.Runner
	db             *gorm.DB
}

func NewDelegate(pipelineRunner pipeline.Runner, db *gorm.DB) *Delegate {
	return &Delegate{
		pipelineRunner,
		db,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.FluxMonitor
}

func (d *Delegate) ServicesForSpec(spec job.SpecDB) (services []job.Service, err error) {
	if spec.FluxMonitorSpec == nil {
		return nil, errors.Errorf("Delegate expects a *job.FluxMonitorSpec to be present, got %v", spec)
	}
	// TODO
	return nil, nil
}
