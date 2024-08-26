package chainlink

import (
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
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

func (c *capabilitiesConfig) WorkflowConnectorConfig() config.WorkflowConnectorConfig {
	return &workflowConnectorConfig{
		c: c.c.WorkflowConnectorConfig,
	}
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

type workflowConnectorConfig struct {
	c toml.WorkflowConnectorConfig
}

func (c *workflowConnectorConfig) ChainIDForNodeKey() big.Int {
	return *c.c.ChainIDForNodeKey
}

func (c *workflowConnectorConfig) GatewayConnectorConfig() config.GatewayConnectorConfig {
	return &gatewayConnectorConfig{
		c: *c.c.GatewayConnectorConfig,
	}
}

type gatewayConnectorConfig struct {
	c toml.GatewayConnectorConfig
}

func (c *gatewayConnectorConfig) NodeAddress() string {
	return *c.c.NodeAddress
}

func (c *gatewayConnectorConfig) DonId() string {
	return *c.c.DonId
}

func (c *gatewayConnectorConfig) Gateways() []config.ConnectorGatewayConfig {
	t := []config.ConnectorGatewayConfig{}
	for index, element := range c.c.Gateways {
		t[index] = &connectorGatewayConfig{element}
	}
	return t
}

func (c *gatewayConnectorConfig) WsClientConfig() network.WebSocketClientConfig {
	return *c.c.WsClientConfig
}

func (c *gatewayConnectorConfig) AuthMinChallengeLen() int {
	return *c.c.AuthMinChallengeLen
}

func (c *gatewayConnectorConfig) AuthTimestampToleranceSec() uint32 {
	return *c.c.AuthTimestampToleranceSec
}

type connectorGatewayConfig struct {
	c toml.ConnectorGatewayConfig
}

func (c *connectorGatewayConfig) Id() string {
	return *c.c.Id
}

func (c *connectorGatewayConfig) URL() string {
	return *c.c.URL
}
