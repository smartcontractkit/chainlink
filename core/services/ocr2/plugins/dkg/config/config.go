package config

// PluginConfig contains custom arguments for the DKG plugin.
type PluginConfig struct {
	EncryptionPublicKey string `json:"encryptionPublicKey"`
	SigningPublicKey    string `json:"signingPublicKey"`
	KeyID               string `json:"keyID"`
	DKGAddress          string `json:"dkgAddress"`
}

// ValidatePluginConfig validates that the given DKG plugin configuration is correct.
func ValidatePluginConfig(config PluginConfig) error {
	// TODO
	return nil
}
