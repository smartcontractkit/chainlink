package cron

import (
	"crypto/sha256"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

// TODO: use for validation test
func ValidateCronJobSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{Pipeline: *pipeline.NewTaskDAG()}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, err
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, err
	}

	var spec job.CronJobSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, err
	}

	spec.OnChainJobSpecID = sha256.Sum256([]byte(tomlString))
	jb.CronRequestSpec = &spec

	if jb.Type != job.CronJob {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}
	if jb.SchemaVersion != uint32(1) {
		return jb, errors.Errorf("the only supported schema version is currently 1, got %v", jb.SchemaVersion)
	}
	return jb, nil
}
