package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.Capabilities = (*capabilitiesConfig)(nil)

type capabilitiesConfig struct {
	c toml.Capabilities
}

func (c *capabilitiesConfig) Peering() config.P2P {
	return &p2p{c: c.c.Peering}
}
