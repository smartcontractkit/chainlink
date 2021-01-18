package fluxmonitorv2

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

type FluxMonitorDelegate struct {
	pipelineRunner pipeline.Runner
	db             *gorm.DB
}

func (d *FluxMonitorDelegate) JobType() job.Type {
	return job.FluxMonitor
}

func (d *FluxMonitorDelegate) ServicesForSpec(spec job.SpecDB) (services []job.Service, err error) {
	if spec.FluxMonitorSpec == nil {
		return nil, errors.Errorf("FluxMonitorDelegate expects a *job.FluxMonitorSpec to be present, got %v", spec)
	}
	// TODO
	return nil, nil
}

func NewFluxMonitorDelegate(pipelineRunner pipeline.Runner, db *gorm.DB) *FluxMonitorDelegate {
	return &FluxMonitorDelegate{
		pipelineRunner,
		db,
	}
}
