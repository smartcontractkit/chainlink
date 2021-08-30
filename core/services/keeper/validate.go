package keeper

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/services/job"
)

func ValidatedKeeperSpec(tomlString string) (job.Job, error) {
	var j = job.Job{
		ExternalJobID: uuid.NewV4(), // Default to generating a uuid, can be overwritten by the specified one in tomlString.
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
	return j, nil
}
