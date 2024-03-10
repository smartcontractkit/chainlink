package keeper

import (
	"strings"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// ValidatedKeeperSpec analyses the tomlString passed as parameter and
// returns a newly-created Job if there are no validation errors inside the toml.
func ValidatedKeeperSpec(tomlString string) (job.Job, error) {
	// Create a new job with a randomly generated uuid, which can be replaced with the one from tomlString.
	var j = job.Job{
		ExternalJobID: uuid.New(),
	}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return j, err
	}

	if err := tree.Unmarshal(&j); err != nil {
		return j, err
	}

	var spec job.KeeperSpec
	if err := tree.Unmarshal(&spec); err != nil {
		return j, err
	}

	j.KeeperSpec = &spec

	if j.Type != job.Keeper {
		return j, errors.Errorf("unsupported type %s", j.Type)
	}

	if strings.Contains(tomlString, "observationSource") ||
		strings.Contains(tomlString, "ObservationSource") {
		return j, errors.New("There should be no 'observationSource' parameter included in the toml")
	}

	return j, nil
}
