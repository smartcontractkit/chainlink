package cron

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func ValidatedCronSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{
		ExternalJobID: uuid.NewV4(), // Default to generating a uuid, can be overwritten by the specified one in tomlString.
	}

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
	if err := utils.ValidateCronSchedule(spec.CronSchedule); err != nil {
		return jb, errors.Wrapf(err, "while validating cron schedule '%v'", spec.CronSchedule)
	}

	return jb, nil
}
