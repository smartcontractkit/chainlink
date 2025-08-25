package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type creConfig struct {
	s toml.CreSecrets
	c toml.CreConfig
}

func (c *creConfig) StreamsAPIKey() string {
	if c.s.Streams == nil || c.s.Streams.APIKey == nil {
		return ""
	}
	return string(*c.s.Streams.APIKey)
}

func (c *creConfig) StreamsAPISecret() string {
	if c.s.Streams == nil || c.s.Streams.APISecret == nil {
		return ""
	}
	return string(*c.s.Streams.APISecret)
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

type workflowFetcherConfig struct {
	url string
}

func (w *workflowFetcherConfig) URL() string {
	return w.url
}

func (c *creConfig) WorkflowFetcher() config.WorkflowFetcher {
	if c.c.WorkflowFetcher == nil || c.c.WorkflowFetcher.URL == nil {
		return &workflowFetcherConfig{url: ""}
	}
	return &workflowFetcherConfig{url: *c.c.WorkflowFetcher.URL}
}

func (c *creConfig) UseLocalTimeProvider() bool {
	if c.c.UseLocalTimeProvider == nil {
		return true // default to local time provider since DON Time plugin may not be running
	}
	return *c.c.UseLocalTimeProvider
}
