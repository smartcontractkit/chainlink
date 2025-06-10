package config

import (
	"github.com/pkg/errors"

	sm_plugin "github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

// ValidateSecureMintConfig validates the secure mint plugin config.
func ValidateSecureMintConfig(cfg *sm_plugin.PorOffchainConfig) error {
	if cfg == nil {
		return errors.New("secure mint config cannot be nil")
	}

	if cfg.MaxChains <= 0 {
		return errors.New("secure mint config MaxChains must be positive")
	}

	return nil
}
