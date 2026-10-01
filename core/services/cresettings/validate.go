package cresettings

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
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

	configType := resolveConfigType(spec)
	switch configType {
	case ConfigTypeSettings, ConfigTypeShardAssignment, ConfigTypeCapRegistry:
	default:
		return jb, fmt.Errorf("unknown config_type %q", configType)
	}

	// Each config_type has exactly one payload field. Rejecting the other one keeps the hash,
	// the delegate's routing, and what is persisted unambiguous.
	if configType == ConfigTypeCapRegistry {
		if spec.Settings != "" {
			return jb, fmt.Errorf("settings must be empty for config_type %q; use offchain_config", ConfigTypeCapRegistry)
		}
	} else if spec.OffchainConfig != "" {
		return jb, fmt.Errorf("offchain_config is only valid for config_type %q, got %q", ConfigTypeCapRegistry, configType)
	}

	// payload is the hashed content and varies by config_type: settings-based types hash
	// the Settings TOML; capabilities_registry hashes the OffchainConfig proto-JSON.
	payload := spec.Settings

	switch configType {
	case ConfigTypeShardAssignment:
		if _, err = ParseShardAssignmentConfig(spec.Settings); err != nil {
			return jb, errors.Wrap(err, "invalid shard_assignment config")
		}
	case ConfigTypeSettings:
		if _, err = settings.NewTOMLGetter([]byte(spec.Settings)); err != nil {
			return jb, errors.Wrap(err, "invalid settings toml")
		}
	case ConfigTypeCapRegistry:
		payload = spec.OffchainConfig
		if err = globalconfig.Validate(spec.OffchainConfig); err != nil {
			return jb, errors.Wrap(err, "invalid capabilities_registry config")
		}
	default:
		return jb, fmt.Errorf("unknown config_type %q", configType)
	}

	shaSum := sha256.Sum256([]byte(payload))
	hash := hex.EncodeToString(shaSum[:])
	if spec.Hash == "" {
		spec.Hash = hash
	} else if spec.Hash != hash {
		return jb, fmt.Errorf("invalid sha256 hash %s: calculated %s from: \n%s", spec.Hash, hash, payload)
	}

	return jb, nil
}

// configTypeEmbeddedCapRegistry is what resolveConfigType returns for a config_type key of
// "capabilities_registry" embedded in Settings. That legacy form is not accepted for
// capabilities_registry: only the top-level config_type field selects it, which is the
// discriminator persisted in cre_settings_specs.config_type (and covered by its single-job
// unique index). Being unknown, it is rejected by validation and by the delegate.
const configTypeEmbeddedCapRegistry = "capabilities_registry (embedded in settings; use the top-level config_type field)"

// resolveConfigType returns the config_type for a spec: the top-level ConfigType field when
// set, else a config_type key embedded in Settings (legacy; not valid for
// capabilities_registry), else the default "settings".
func resolveConfigType(spec job.CRESettingsSpec) string {
	if spec.ConfigType != "" {
		return spec.ConfigType
	}
	if spec.Settings != "" {
		if ct, ok := extractConfigType(spec.Settings); ok {
			if ct == ConfigTypeCapRegistry {
				return configTypeEmbeddedCapRegistry
			}
			return ct
		}
	}
	return ConfigTypeSettings
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
