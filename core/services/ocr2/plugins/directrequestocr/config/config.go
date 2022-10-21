package config

// The PluginConfig struct contains the custom arguments needed for the Median plugin.
type PluginConfig struct {
	MinIncomingConfirmations uint32 `json:"minIncomingConfirmations"`
	MaxRequestsPerOCRRound   uint32 `json:"maxRequestsPerOCRRound"`
	RequestE2eTimeoutMillis  uint32 `json:"requestE2eTimeoutMillis"`
}

func ValidatePluginConfig(config PluginConfig) error {
	return nil
}
