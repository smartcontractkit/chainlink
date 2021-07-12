package fluxmonitorv2

import (
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	coreorm "github.com/smartcontractkit/chainlink/core/store/orm"
)

func ValidatedFluxMonitorSpec(config *coreorm.Config, ts string) (job.Job, error) {
	var jb = job.Job{
		ExternalJobID: uuid.NewV4(), // Default to generating a uuid, can be overwritten by the specified one in tomlString.
	}
	var spec job.FluxMonitorSpec
	tree, err := toml.Load(ts)
	if err != nil {
		return jb, err
	}
	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, err
	}
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, err
	}
	jb.FluxMonitorSpec = &spec

	if jb.Type != job.FluxMonitor {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}
	if jb.SchemaVersion != uint32(1) {
		return jb, errors.Errorf("the only supported schema version is currently 1, got %v", jb.SchemaVersion)
	}
	if jb.Pipeline.HasAsync() {
		return jb, errors.Errorf("async=true tasks are not supported for %v", jb.Type)
	}

	// Find the smallest of all the timeouts
	// and ensure the polling period is greater than that.
	minTaskTimeout, aTimeoutSet, err := jb.Pipeline.MinTimeout()
	if err != nil {
		return jb, err
	}
	timeouts := []time.Duration{
		config.DefaultHTTPTimeout().Duration(),
		time.Duration(jb.MaxTaskDuration),
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

	if !validatePollTimer(jb.FluxMonitorSpec.PollTimerDisabled, minTimeout, jb.FluxMonitorSpec.PollTimerPeriod) {
		return jb, errors.Errorf("pollTimer.period must be equal or greater than %v, got %v", minTimeout, jb.FluxMonitorSpec.PollTimerPeriod)
	}

	return jb, nil
}

// validatePollTime validates the period is greater than the min timeout for an
// enabled poll timer.
func validatePollTimer(disabled bool, minTimeout time.Duration, period time.Duration) bool {
	// Disabled timers do not need to validate the period
	if disabled {
		return true
	}

	return period >= minTimeout
}
