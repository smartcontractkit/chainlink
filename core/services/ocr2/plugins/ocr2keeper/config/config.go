package config

// PluginConfig contains custom arguments for the OCR2Keeper plugin.
type PluginConfig struct {
	// MinIncomingConfirmations is the minimum number of confirmations required ty sync updates from keeper registry
	MinIncomingConfirmations uint32 `json:"minIncomingConfirmations"`
}

// Validate validates that the given OCR2Keeper plugin configuration is correct.
func (cfg PluginConfig) Validate() error {
	return nil
}
