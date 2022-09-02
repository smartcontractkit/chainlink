package config

import "fmt"

// PluginConfig contains custom arguments for the OCR2Keeper plugin.
type PluginConfig struct {
	// Reporting plugin limits
	MaxQueryLength       int `json:"maxQueryLength"`
	MaxObservationLength int `json:"maxObservationLength"`
	MaxReportLength      int `json:"maxReportLength"`
}

// Validate validates that the given OCR2Keeper plugin configuration is correct.
func (cfg PluginConfig) Validate() error {
	if cfg.MaxQueryLength == 0 {
		return fmt.Errorf("pluginConfig.maxQueryLength must be positive value")
	}

	if cfg.MaxObservationLength == 0 {
		return fmt.Errorf("pluginConfig.maxQueryLength must be positive value")
	}

	if cfg.MaxReportLength == 0 {
		return fmt.Errorf("pluginConfig.maxQueryLength must be positive value")
	}

	return nil
}
