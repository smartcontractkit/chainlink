package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type creConfig struct {
	c toml.CreSecrets
}

func (c *creConfig) StreamsApiKey() string {
	if c.c.Streams == nil || c.c.Streams.ApiKey == nil {
		return ""
	}
	return string(*c.c.Streams.ApiKey)
}

func (c *creConfig) StreamsApiSecret() string {
	if c.c.Streams == nil || c.c.Streams.ApiSecret == nil {
		return ""
	}
	return string(*c.c.Streams.ApiSecret)
} 