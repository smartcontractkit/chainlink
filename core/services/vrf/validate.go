package vrf

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
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

	var spec job.VRFSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on job")
	}
	var empty secp256k1.PublicKey
	if bytes.Equal(spec.PublicKey[:], empty[:]) {
		return jb, errors.Wrap(ErrKeyNotSet, "publicKey")
	}
	if spec.MinIncomingConfirmations == 0 {
		return jb, errors.Wrap(ErrKeyNotSet, "minIncomingConfirmations")
	}
	if spec.CoordinatorAddress.String() == "" {
		return jb, errors.Wrap(ErrKeyNotSet, "coordinatorAddress")
	}
	if spec.RequestedConfsDelay < 0 {
		return jb, errors.Wrap(ErrKeyNotSet, "requestedConfsDelay must be >= 0")
	}
	// If a request timeout is not provided set it to a reasonable default.
	if spec.RequestTimeout == 0 {
		spec.RequestTimeout = 24 * time.Hour
	}

	if spec.BatchFulfillmentEnabled && spec.BatchCoordinatorAddress == nil {
		return jb, errors.Wrap(ErrKeyNotSet, "batch coordinator address must be provided if batchFulfillmentEnabled = true")
	}

	if spec.BatchFulfillmentGasMultiplier <= 0 {
		spec.BatchFulfillmentGasMultiplier = 1.15
	}

	if spec.ChunkSize == 0 {
		spec.ChunkSize = 20
	}

	if spec.BackoffMaxDelay < spec.BackoffInitialDelay {
		return jb, fmt.Errorf("backoff max delay (%s) cannot be less than backoff initial delay (%s)",
			spec.BackoffMaxDelay.String(), spec.BackoffInitialDelay.String())
	}

	if spec.GasLanePrice != nil && spec.GasLanePrice.Cmp(assets.GWei(0)) <= 0 {
		return jb, fmt.Errorf("gasLanePrice must be positive, given: %s", spec.GasLanePrice.String())
	}

	var foundVRFTask bool
	for _, t := range jb.Pipeline.Tasks {
		if t.Type() == pipeline.TaskTypeVRF || t.Type() == pipeline.TaskTypeVRFV2 {
			foundVRFTask = true
		}

		if t.Type() == pipeline.TaskTypeVRFV2 {
			if len(spec.FromAddresses) == 0 {
				return jb, errors.Wrap(ErrKeyNotSet, "fromAddreses needs to have a non-zero length")
			}
		}
	}
	if !foundVRFTask {
		return jb, errors.Wrapf(ErrKeyNotSet, "invalid pipeline, expected a vrf task")
	}

	jb.VRFSpec = &spec

	return jb, nil
}
