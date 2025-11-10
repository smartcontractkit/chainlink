package cresettings

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func ValidatedCRESettingsSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{
		ExternalJobID: uuid.New(), // Default to generating a uuid, can be overwritten by the specified one in tomlString.
	}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, errors.Wrap(err, "toml error on load")
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on spec")
	}

	var spec job.CRESettingsSpec
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, errors.Wrap(err, "toml unmarshal error on job")
	}

	jb.CRESettingsSpec = &spec
	if jb.Type != job.CRESettings {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}

	_, err = settings.NewTOMLGetter([]byte(jb.CRESettingsSpec.Settings))
	if err != nil {
		return jb, errors.Wrap(err, "invalid settings toml")
	}
	shaSum := sha256.Sum256([]byte(spec.Settings))
	hash := hex.EncodeToString(shaSum[:])
	if spec.Hash == "" {
		spec.Hash = hash
	} else if spec.Hash != hash {
		return jb, fmt.Errorf("invalid sha256 hash %s: calculated %s", spec.Hash, hash)
	}

	return jb, nil
}
