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
	// invalid operation: cannot indirect
	// c.c.GatewayConnectorConfig (variable of type "github.com/smartcontractkit/chainlink/v2/core/config/toml".ConnectorConfig)
	// compilerInvalidIndirection
	return &gatewayConnectorConfig{
		c: *c.c.GatewayConnectorConfig,
	}
}

type gatewayConnectorConfig struct {
	c toml.GatewayConnectorConfig
}

func (c *gatewayConnectorConfig) NodeAddress() string {
	// why not *c.c.NodeAddress like the above?
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
