package vrf

import (
	"bytes"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
)

var (
	ErrKeyNotSet = errors.New("key not set")
)

func ValidateVRFSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{Pipeline: *pipeline.NewTaskDAG()}

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
