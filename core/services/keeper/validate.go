package keeper

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func ValidatedKeeperSpec(tomlString string) (job.Job, error) {
	var j = job.Job{
		Pipeline: *pipeline.NewTaskDAG(),
	}
	var spec job.KeeperSpec
	tree, err := toml.Load(tomlString)
	if err != nil {
		return j, err
	}
	err = tree.Unmarshal(&j)
	if err != nil {
		return j, err
	}
	err = tree.Unmarshal(&spec)
	if err != nil {
		return j, err
	}
	j.KeeperSpec = &spec

	if j.Type != job.Keeper {
		return j, errors.Errorf("unsupported type %s", j.Type)
	}
	if j.SchemaVersion != uint32(1) {
		return j, errors.Errorf("the only supported schema version is currently 1, got %d", j.SchemaVersion)
	}
	return j, nil
}
