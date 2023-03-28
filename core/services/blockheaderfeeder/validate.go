package blockheaderfeeder

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

	err = validateChainID(spec.EVMChainID.Int64())
	if err != nil {
		return jb, err
	}

	// Defaults
	if spec.WaitBlocks == 0 {
		spec.WaitBlocks = 256
	}
	if spec.LookbackBlocks == 0 {
		spec.LookbackBlocks = 1000
	}
	if spec.PollPeriod == 0 {
		spec.PollPeriod = 15 * time.Second
	}
	if spec.RunTimeout == 0 {
		spec.RunTimeout = 30 * time.Second
	}
	if spec.StoreBlockhashesBatchSize == 0 {
		spec.StoreBlockhashesBatchSize = 10
	}
	if spec.GetBlockhashesBatchSize == 0 {
		spec.GetBlockhashesBatchSize = 100
	}

	if spec.WaitBlocks < 256 {
		return jb, errors.New(`"waitBlocks" must be greater than or equal to 256`)
	}
	if spec.LookbackBlocks <= 256 {
		return jb, errors.New(`"lookbackBlocks" must be greater than 256`)
	}
	if spec.WaitBlocks >= spec.LookbackBlocks {
		return jb, errors.New(`"lookbackBlocks" must be greater than "waitBlocks"`)
	}

	jb.BlockHeaderFeederSpec = &spec

	return jb, nil
}

func notSet(field string) error {
	return errors.Errorf("%q must be set", field)
}

// validateChainID validates whether the given chain is supported
// Avax chain is not supported because block header format
// is different from go-ethereum types.Header.
// Special handling for Avax chains is not yet supported
func validateChainID(evmChainID int64) error {
	if evmChainID == 43114 || // C-chain mainnet
		evmChainID == 43113 { // Fuji testnet
		return errors.Errorf("unsupported chain")
	}
	return nil
}
