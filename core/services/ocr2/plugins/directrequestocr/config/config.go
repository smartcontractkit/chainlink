package config

// This Config is part of the job spec.
// TODO: Define which values should be in the config stored in the contract vs here.
// https://app.shortcut.com/chainlinklabs/story/56615/config-for-reporting-plugin
type PluginConfig struct {
	MinIncomingConfirmations uint32 `json:"minIncomingConfirmations"`
	MaxRequestsPerOCRRound   uint32 `json:"maxRequestsPerOCRRound"`
	RequestE2eTimeoutMillis  uint32 `json:"requestE2eTimeoutMillis"`
}

func ValidatePluginConfig(config PluginConfig) error {
	return nil
}
