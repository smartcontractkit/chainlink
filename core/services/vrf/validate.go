package vrf

import (
	"bytes"

	uuid "github.com/satori/go.uuid"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
)

var (
	ErrKeyNotSet = errors.New("key not set")
)

func ValidatedVRFSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{
		ExternalJobID: uuid.NewV4(), // Default to generating a uuid, can be overwritten by the specified one in tomlString.
	}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, errors.Wrap(err, "toml error on load")
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on spec")
	}
	if jb.Type != job.VRF {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}
	if jb.SchemaVersion != uint32(1) {
		return jb, errors.Errorf("the only supported schema version is currently 1, got %v", jb.SchemaVersion)
	}
	if jb.Pipeline.HasAsync() {
		return jb, errors.Errorf("async=true tasks are not supported for %v", jb.Type)
	}

	var spec job.VRFSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on job")
	}
	var empty secp256k1.PublicKey
	if bytes.Equal(spec.PublicKey[:], empty[:]) {
		return jb, errors.Wrap(ErrKeyNotSet, "publicKey")
	}
	if spec.Confirmations == 0 {
		return jb, errors.Wrap(ErrKeyNotSet, "confirmations")
	}
	if spec.CoordinatorAddress.String() == "" {
		return jb, errors.Wrap(ErrKeyNotSet, "coordinatorAddress")
	}

	jb.VRFSpec = &spec

	return jb, nil
}
