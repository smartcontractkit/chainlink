package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type creConfig struct {
	s toml.CreSecrets
	c toml.CreConfig
}

func (c *creConfig) StreamsApiKey() string {
	if c.s.Streams == nil || c.s.Streams.ApiKey == nil {
		return ""
	}
	return string(*c.s.Streams.ApiKey)
}

func (c *creConfig) StreamsApiSecret() string {
	if c.s.Streams == nil || c.s.Streams.ApiSecret == nil {
		return ""
	}
	return string(*c.s.Streams.ApiSecret)
}

func (c *creConfig) WsURL() string {
	if c.c.Streams == nil || c.c.Streams.WsURL == nil {
		return ""
	}
	return *c.c.Streams.WsURL
}

func (c *creConfig) RestURL() string {
	if c.c.Streams == nil || c.c.Streams.RestURL == nil {
		return ""
	}
	return *c.c.Streams.RestURL
}
