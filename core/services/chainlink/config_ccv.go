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
