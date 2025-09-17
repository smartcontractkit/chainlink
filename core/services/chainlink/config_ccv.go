package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

type ccvConfig struct {
	s toml.CCVSecrets
	c toml.CCVConfig
}

func (c *ccvConfig) Enabled() bool {
	if c.c.Enabled == nil {
		return false
	}
	return *c.c.Enabled
}

func (c ccvConfig) ExecutorIndexerAPIKey() string {
	if c.s.Executor == nil || c.s.Executor.IndexerAPIKey == nil {
		return ""
	}
	return string(*c.s.Executor.IndexerAPIKey)
}
