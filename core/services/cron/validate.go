package cron

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func ValidateCronSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{Pipeline: *pipeline.NewTaskDAG()}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, errors.Wrap(err, "toml error on load")
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on spec")
	}

	var spec job.CronSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on job")
	}

	jb.CronSpec = &spec
	if jb.Type != job.Cron {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}
	if jb.SchemaVersion != uint32(1) {
		return jb, errors.Errorf("the only supported schema version is currently 1, got %v", jb.SchemaVersion)
	}

	if _, err := cron.New().AddFunc(spec.CronSchedule, func() {}); err != nil {
		return jb, errors.Errorf("error parsing cron schedule: %v", err)
	}

	return jb, nil
}
