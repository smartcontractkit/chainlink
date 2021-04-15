package cron

import (
	"crypto/sha256"
	"math/big"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/robfig/cron"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func ValidateCronSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{Pipeline: *pipeline.NewTaskDAG()}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, err
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, err
	}

	var spec job.CronSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, err
	}

	spec.OnChainJobSpecID = sha256.Sum256([]byte(tomlString))
	jb.CronSpec = &spec

	if jb.Type != job.CronJob {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}
	if jb.SchemaVersion != uint32(1) {
		return jb, errors.Errorf("the only supported schema version is currently 1, got %v", jb.SchemaVersion)
	}

	_, err = cron.Parse(spec.CronSchedule)
	if err != nil {
		return jb, errors.Errorf("error parsing cron schedule: %v", err)
	}

	if spec.OraclePayment.IsZero() || spec.OraclePayment.ToInt().Cmp(big.NewInt(0)) == -1 {
		return jb, errors.Errorf("invalid oracle payment input")
	}

	return jb, nil
}
