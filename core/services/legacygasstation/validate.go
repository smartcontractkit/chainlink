package legacygasstation

import (
	"time"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// ValidatedServerSpec validates and converts the given toml string to a job.Job.
func ValidatedServerSpec(tomlString string) (job.Job, error) {
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

	if jb.Type != job.LegacyGasStationServer {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}

	var spec job.LegacyGasStationServerSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, errors.Wrap(err, "unmarshalling toml job")
	}

	// Required fields
	if spec.ForwarderAddress == "" {
		return jb, notSet("forwarderAddress")
	}
	if spec.EVMChainID == nil {
		return jb, notSet("evmChainID")
	}
	if spec.CCIPChainSelector == nil {
		return jb, notSet("ccipChainSelector")
	}
	if spec.FromAddresses == nil {
		return jb, notSet("fromAddresses")
	}

	jb.LegacyGasStationServerSpec = &spec

	return jb, nil
}

// ValidatedSidecarSpec validates and converts the given toml string to a job.Job.
func ValidatedSidecarSpec(tomlString string) (job.Job, error) {
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

	if jb.Type != job.LegacyGasStationSidecar {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}

	var spec job.LegacyGasStationSidecarSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, errors.Wrap(err, "unmarshalling toml job")
	}

	// Required fields
	if spec.ForwarderAddress == "" {
		return jb, notSet("forwarderAddress")
	}
	if spec.OffRampAddress == "" {
		return jb, notSet("offRampAddress")
	}
	if spec.EVMChainID == nil {
		return jb, notSet("evmChainID")
	}
	if spec.CCIPChainSelector == nil {
		return jb, notSet("ccipChainSelector")
	}

	if spec.StatusUpdateURL == "" {
		return jb, notSet("statusUpdateURL")
	}

	// Defaults
	if spec.LookbackBlocks == 0 {
		spec.LookbackBlocks = 10000
	}
	if spec.PollPeriod == 0 {
		spec.PollPeriod = 15 * time.Second
	}
	if spec.RunTimeout == 0 {
		spec.RunTimeout = 30 * time.Second
	}

	jb.LegacyGasStationSidecarSpec = &spec

	return jb, nil
}

func notSet(field string) error {
	return errors.Errorf("%q must be set", field)
}
