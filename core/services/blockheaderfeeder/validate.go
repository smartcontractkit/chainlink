package blockheaderfeeder

import (
	"time"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/services/job"
)

// ValidatedBlockHeaderFeederSpec validates and converts the given toml string to a job.Job.
func ValidatedBlockHeaderFeederSpec(tomlString string) (job.Job, error) {
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

	if jb.Type != job.BlockHeaderFeeder {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}

	var spec job.BlockHeaderFeederSpec
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
	if spec.BatchBlockhashStoreAddress == "" {
		return jb, notSet("batchBlockhashStoreAddress")
	}
	if spec.EVMChainID == nil {
		return jb, notSet("evmChainID")
	}

	// Defaults
	// TODO: revisit defaults
	if spec.LookbackBlocks == 0 {
		spec.LookbackBlocks = 1000
	}
	if spec.PollPeriod == 0 {
		spec.PollPeriod = 30 * time.Second
	}
	if spec.RunTimeout == 0 {
		spec.RunTimeout = 30 * time.Second
	}
	if spec.StoreBlockhashesBatchSize == 0 {
		spec.StoreBlockhashesBatchSize = 10
	}
	if spec.GetBlockhashesBatchSize == 0 {
		spec.GetBlockhashesBatchSize = 10
	}

	if spec.LookbackBlocks <= 256 {
		return jb, errors.New(`"lookback" must be greater than 256`)
	}

	jb.BlockHeaderFeederSpec = &spec

	return jb, nil
}

func notSet(field string) error {
	return errors.Errorf("%q must be set", field)
}
