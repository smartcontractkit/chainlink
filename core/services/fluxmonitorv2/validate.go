package fluxmonitorv2

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func ValidatedFluxMonitorSpec(ts string) (job.SpecDB, error) {
	var specDB = job.SpecDB{
		Pipeline: *pipeline.NewTaskDAG(),
	}
	var spec job.FluxMonitorSpec
	tree, err := toml.Load(ts)
	if err != nil {
		return specDB, err
	}
	err = tree.Unmarshal(&specDB)
	if err != nil {
		return specDB, err
	}
	err = tree.Unmarshal(&spec)
	if err != nil {
		return specDB, err
	}
	specDB.FluxMonitorSpec = &spec

	if specDB.Type != job.FluxMonitor {
		return specDB, errors.Errorf("unsupported type %s", specDB.Type)
	}
	if specDB.SchemaVersion != uint32(1) {
		return specDB, errors.Errorf("the only supported schema version is currently 1, got %v", specDB.SchemaVersion)
	}
	return specDB, nil
}
