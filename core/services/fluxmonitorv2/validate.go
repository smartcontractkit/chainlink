package fluxmonitorv2

import (
	"time"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

func ValidatedFluxMonitorSpec(config *orm.Config, ts string) (job.SpecDB, error) {
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

	// Find the smallest of all the timeouts
	// and ensure the polling period is greater than that.
	minTaskTimeout, aTimeoutSet, err := specDB.Pipeline.MinTimeout()
	if err != nil {
		return specDB, err
	}
	timeouts := []time.Duration{
		config.DefaultHTTPTimeout().Duration(),
		time.Duration(specDB.MaxTaskDuration),
	}
	if aTimeoutSet {
		timeouts = append(timeouts, minTaskTimeout)
	}
	var minTimeout time.Duration = 1<<63 - 1
	for _, timeout := range timeouts {
		if timeout < minTimeout {
			minTimeout = timeout
		}
	}
	if !spec.PollTimerDisabled && spec.PollTimerPeriod < minTimeout {
		return specDB, errors.Errorf("pollTimer.period must be equal or greater than %v, got %v", minTimeout, spec.PollTimerPeriod)
	}
	return specDB, nil
}
