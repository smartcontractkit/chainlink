package config

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// SecureMintConfig holds secure mint specific configuration
type SecureMintConfig struct {
	Token    string `json:"token"`
	Reserves string `json:"reserves"`
}

// Parse parses the secure mint configuration from JSON bytes
func Parse(configBytes []byte) (*SecureMintConfig, error) {
	if len(configBytes) == 0 {
		return nil, errors.New("secure mint config cannot be empty")
	}

	var config SecureMintConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal SecureMintConfig")
	}

	return &config, nil
}

// Validate validates the secure mint plugin-specific config.
func (cfg *SecureMintConfig) Validate() error {
	if cfg == nil {
		return errors.New("secure mint plugin config cannot be nil")
	}

	if cfg.Token == "" {
		return errors.New("token cannot be empty")
	}

	if cfg.Reserves == "" {
		return errors.New("reserves cannot be empty")
	}

	return nil
}
