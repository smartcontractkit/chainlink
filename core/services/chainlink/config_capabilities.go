package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var _ config.Capabilities = (*capabilitiesConfig)(nil)

type capabilitiesConfig struct {
	c toml.Capabilities
}

func (c *capabilitiesConfig) Peering() config.P2P {
	return &p2p{c: c.c.Peering}
}

func (c *capabilitiesConfig) Registry() config.CapabilitiesRegistry {
	return &capabilitiesRegistry{
		c: c.c.Registry,
	}
}

type capabilitiesRegistry struct {
	c toml.Registry
}

func (c *capabilitiesRegistry) RelayID() types.RelayID {
	return types.NewRelayID(c.NetworkID(), c.ChainID())
}

func (c *capabilitiesRegistry) NetworkID() string {
	return c.c.NetworkID
}

func (c *capabilitiesRegistry) ChainID() string {
	return c.c.ChainID
}

func (c *capabilitiesRegistry) RemoteAddress() string {
	return c.c.RemoteAddress
}
