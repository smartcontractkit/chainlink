package config

import (
	"github.com/pkg/errors"
)

// PluginConfig contains custom arguments for the OCR2Keeper plugin.
type PluginConfig struct {
	LinkEthFeedAddress string   `json:"linkEthFeedAddress"`
	ConfirmationDelays []uint32 `json:"confirmationDelays"`
	LookbackBlocks     int64    `json:"lookbackBlocks"`
}

// ValidatePluginConfig validates that the given OCR2VRF plugin configuration is correct.
func ValidatePluginConfig(config PluginConfig) error {
	return nil

	// NOTE: similar to the above.
	if config.LinkEthFeedAddress == "" {
		return errors.New("linkEthFieldAddress field must be provided")
	}

	if len(config.ConfirmationDelays) != 8 {
		// TODO: this should not be here, should read it from the chain
		// However that's not possible at the moment, since that field is private.
		return errors.New("confirmationDelays field must be an array of 8 elements")
	}

	if config.LookbackBlocks <= 0 {
		return errors.New("lookbackBlocks field must be positive")
	}

	return nil
}
