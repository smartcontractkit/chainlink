package config

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

// PluginConfig contains configuration for the Ring OCR plugin.
type PluginConfig struct {
	ShardConfigAddr string `json:"shardConfigAddr" toml:"shardConfigAddr"`
}

// Unmarshal parses the plugin config from JSON bytes.
func (p *PluginConfig) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

// Validate validates the plugin configuration.
func (p *PluginConfig) Validate() error {
	if p.ShardConfigAddr == "" {
		return errors.New("shardConfigAddr is required")
	}
	if !common.IsHexAddress(p.ShardConfigAddr) {
		return errors.New("shardConfigAddr is not a valid hex address")
	}
	return nil
}
