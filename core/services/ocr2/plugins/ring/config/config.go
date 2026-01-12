package config

import (
	"encoding/json"
)

// PluginConfig contains configuration for the Ring OCR plugin.
type PluginConfig struct {
	ShardConfigAddr string `json:"shardConfigAddr" toml:"shardConfigAddr"`
}

// Unmarshal parses the plugin config from JSON bytes.
func (p *PluginConfig) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}
