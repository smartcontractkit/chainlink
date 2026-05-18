package blockhashstore

import (
	"time"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

var EmptyAddress = utils.ZeroAddress.Hex()

// ValidatedSpec validates and converts the given toml string to a job.Job.
func ValidatedSpec(tomlString string) (job.Job, error) {
	jb := job.Job{
		// Default to generating a UUID, can be overwritten by the specified one in tomlString.
		ExternalJobID: uuid.New(),
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
	if spec.CoordinatorV1Address == nil && spec.CoordinatorV2Address == nil && spec.CoordinatorV2PlusAddress == nil {
		return jb, errors.New(
			`at least one of "coordinatorV1Address", "coordinatorV2Address" and "coordinatorV2PlusAddress" must be set`)
	}
	if spec.BlockhashStoreAddress == "" {
		return jb, notSet("blockhashStoreAddress")
	}
	if spec.EVMChainID == nil {
		return jb, notSet("evmChainID")
	}
	if spec.TrustedBlockhashStoreAddress != nil && spec.TrustedBlockhashStoreAddress.Hex() != EmptyAddress && spec.TrustedBlockhashStoreBatchSize == 0 {
		return jb, notSet("trustedBlockhashStoreBatchSize")
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
	if spec.HeartbeatPeriod < 0 {
		return jb, errors.New(`"heartbeatPeriod" must be greater than 0`)
	}
	// spec.HeartbeatPeriodTime == 0, default is heartbeat disabled

	// Validation
	if spec.WaitBlocks >= spec.LookbackBlocks {
		return jb, errors.New(`"waitBlocks" must be less than "lookbackBlocks"`)
	}
	if (spec.TrustedBlockhashStoreAddress == nil || spec.TrustedBlockhashStoreAddress.Hex() == EmptyAddress) && spec.WaitBlocks >= 256 {
		return jb, errors.New(`"waitBlocks" must be less than 256`)
	}
	if (spec.TrustedBlockhashStoreAddress == nil || spec.TrustedBlockhashStoreAddress.Hex() == EmptyAddress) && spec.LookbackBlocks >= 256 {
		return jb, errors.New(`"lookbackBlocks" must be less than 256`)
	}

	jb.BlockhashStoreSpec = &spec

	return jb, nil
}

func notSet(field string) error {
	return errors.Errorf("%q must be set", field)
}
