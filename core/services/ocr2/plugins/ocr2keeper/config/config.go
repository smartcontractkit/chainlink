package config

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	// OCRInstancesNumber is the default OCR instances to be running
	OCRInstancesNumber = 1
)

// PluginConfig contains custom arguments for the OCR2Keeper plugin.
type PluginConfig struct {
	OCRInstances uint8 `json:"ocrInstances"`

	// TODO
	LinkEthFeedAddress string   `json:"linkEthFeedAddress"`
	ConfirmationDelays []uint32 `json:"confirmationDelays"`
	LookbackBlocks     int64    `json:"lookbackBlocks"`
}

// Validate validates that the given OCR2Keeper plugin configuration is correct.
func (cfg PluginConfig) Validate() error {
	// Number of OCR instances must be equal 100
	if cfg.OCRInstances != OCRInstancesNumber {
		return fmt.Errorf("ocrInstances must be %d", OCRInstancesNumber)
	}

	return nil

	// NOTE: similar to the above.
	if cfg.LinkEthFeedAddress == "" {
		return errors.New("linkEthFieldAddress field must be provided")
	}

	if len(cfg.ConfirmationDelays) != 8 {
		// TODO: this should not be here, should read it from the chain
		// However that's not possible at the moment, since that field is private.
		return errors.New("confirmationDelays field must be an array of 8 elements")
	}

	if cfg.LookbackBlocks <= 0 {
		return errors.New("lookbackBlocks field must be positive")
	}

	return nil
}
