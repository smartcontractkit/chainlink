package ccvexecutor

import (
	"fmt"

	"github.com/pelletier/go-toml"

	"github.com/smartcontractkit/chainlink-ccv/executor"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func ValidatedCCVExecutorSpec(tomlString string) (jb job.Job, err error) {
	var spec job.CCVExecutorSpec
	tree, err := toml.Load(tomlString)
	if err != nil {
		return job.Job{}, fmt.Errorf("toml error on load: %w", err)
	}

	err = tree.Unmarshal(&spec)
	if err != nil {
		return job.Job{}, fmt.Errorf("toml unmarshal error on spec: %w", err)
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return job.Job{}, err
	}
	jb.CCVExecutorSpec = &spec

	if jb.Type != job.CCVExecutor {
		return job.Job{}, fmt.Errorf("the only supported type is currently 'ccvexecutor', got %s", jb.Type)
	}

	var cfg executor.Configuration
	err = toml.Unmarshal([]byte(jb.CCVExecutorSpec.ExecutorConfig), &cfg)
	if err != nil {
		return job.Job{}, fmt.Errorf("failed to unmarshal executorConfig into the executor config struct: %w", err)
	}

	// validation functions are called in ServicesForSpec
	// so that this method can be used in tests w/out crafting
	// a valid configuration.

	return jb, nil
}
