package blockhashstore

import (
	"time"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// ValidatedSpec validates and converts the given toml string to a job.Job.
func ValidatedSpec(tomlString string) (job.Job, error) {
	jb := job.Job{
		// Default to generating a UUID, can be overwritten by the specified one in tomlString.
		ExternalJobID: uuid.NewV4(),
	}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, errors.Wrap(err, "loading toml")
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, errors.Wrap(err, "unmarshalling toml spec")
	}

	if jb.Type != job.BlockhashStore {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}

	var spec job.BlockhashStoreSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, errors.Wrap(err, "unmarshalling toml job")
	}

	// Required fields
	if spec.CoordinatorV1Address == nil && spec.CoordinatorV2Address == nil {
		return jb, errors.New(
			`at least one of "coordinatorV1Address" and "coordinatorV2Address" must be set`)
	}
	if spec.BlockhashStoreAddress == "" {
		return jb, notSet("blockhashStoreAddress")
	}
	if spec.EVMChainID == nil {
		return jb, notSet("evmChainID")
	}

	// Defaults
	if spec.WaitBlocks == 0 {
		spec.WaitBlocks = 100
	}
	if spec.LookbackBlocks == 0 {
		spec.LookbackBlocks = 200
	}
	if spec.PollPeriod == 0 {
		spec.PollPeriod = 30 * time.Second
	}
	if spec.RunTimeout == 0 {
		spec.RunTimeout = 30 * time.Second
	}

	// Validation
	if spec.WaitBlocks >= spec.LookbackBlocks {
		return jb, errors.New(`"waitBlocks" must be less than "lookbackBlocks"`)
	}
	if spec.WaitBlocks >= 256 {
		return jb, errors.New(`"waitBlocks" must be less than 256`)
	}
	if spec.LookbackBlocks >= 256 {
		return jb, errors.New(`"lookbackBlocks" must be less than 256`)
	}

	jb.BlockhashStoreSpec = &spec

	return jb, nil
}

func notSet(field string) error {
	return errors.Errorf("%q must be set", field)
}
