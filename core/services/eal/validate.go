package eal

import (
	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func ValidatedEALSpec(tomlString string) (jb job.Job, err error) {
	jb = job.Job{
		// Default to generating a UUID, can be overwritten by the specified one in tomlString.
		ExternalJobID: uuid.New(),
	}
	var spec job.EALSpec
	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, errors.Wrap(err, "toml error on load")
	}
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on spec")
	}
	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on job")
	}

	if jb.Type != job.EAL {
		return jb, errors.Errorf("the only supported type is 'EAL', got %s", jb.Type)
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

	jb.EALSpec = &spec

	return jb, nil
}

func notSet(field string) error {
	return errors.Errorf("%q must be set", field)
}
