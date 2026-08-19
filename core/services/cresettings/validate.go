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
		ExternalJobID: uuid.New(),
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

	configType := ConfigTypeSettings
	if spec.Settings != "" {
		ct, ok := extractConfigType(spec.Settings)
		if ok {
			configType = ct
		}
	}

	switch configType {
	case ConfigTypeShardAssignment:
		if _, err = ParseShardAssignmentConfig(spec.Settings); err != nil {
			return jb, errors.Wrap(err, "invalid shard_assignment config")
		}
	case ConfigTypeSettings:
		if _, err = settings.NewTOMLGetter([]byte(spec.Settings)); err != nil {
			return jb, errors.Wrap(err, "invalid settings toml")
		}
	default:
		return jb, fmt.Errorf("unknown config_type %q", configType)
	}

	shaSum := sha256.Sum256([]byte(spec.Settings))
	hash := hex.EncodeToString(shaSum[:])
	if spec.Hash == "" {
		spec.Hash = hash
	} else if spec.Hash != hash {
		return jb, fmt.Errorf("invalid sha256 hash %s: calculated %s from: \n%s", spec.Hash, hash, spec.Settings)
	}

	return jb, nil
}

func extractConfigType(settings string) (string, bool) {
	tree, err := toml.Load(settings)
	if err != nil {
		return "", false
	}
	v := tree.Get("config_type")
	if v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}
