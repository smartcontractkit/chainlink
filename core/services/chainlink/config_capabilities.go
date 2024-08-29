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

func (c *capabilitiesConfig) ExternalRegistry() config.CapabilitiesExternalRegistry {
	return &capabilitiesExternalRegistry{
		c: c.c.ExternalRegistry,
	}
}

func (c *capabilitiesConfig) Dispatcher() config.Dispatcher {
	return &dispatcher{d: c.c.Dispatcher}
}

type dispatcher struct {
	d toml.Dispatcher
}

func (d *dispatcher) SupportedVersion() int {
	return *d.d.SupportedVersion
}

func (d *dispatcher) ReceiverBufferSize() int {
	return *d.d.ReceiverBufferSize
}

func (d *dispatcher) RateLimit() config.DispatcherRateLimit {
	return &dispatcherRateLimit{r: d.d.RateLimit}
}

type dispatcherRateLimit struct {
	r toml.DispatcherRateLimit
}

func (r *dispatcherRateLimit) GlobalRPS() float64 {
	return *r.r.GlobalRPS
}

func (r *dispatcherRateLimit) GlobalBurst() int {
	return *r.r.GlobalBurst
}

func (r *dispatcherRateLimit) PerSenderRPS() float64 {
	return *r.r.PerSenderRPS
}

func (r *dispatcherRateLimit) PerSenderBurst() int {
	return *r.r.PerSenderBurst
}

type capabilitiesExternalRegistry struct {
	c toml.ExternalRegistry
}

func (c *capabilitiesExternalRegistry) RelayID() types.RelayID {
	return types.NewRelayID(c.NetworkID(), c.ChainID())
}

func (c *capabilitiesExternalRegistry) NetworkID() string {
	return *c.c.NetworkID
}

func (c *capabilitiesExternalRegistry) ChainID() string {
	return *c.c.ChainID
}

func (c *capabilitiesExternalRegistry) Address() string {
	return *c.c.Address
}
