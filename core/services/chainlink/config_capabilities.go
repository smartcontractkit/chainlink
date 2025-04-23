package chainlink

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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

func (c *capabilitiesConfig) WorkflowRegistry() config.CapabilitiesWorkflowRegistry {
	return &capabilitiesWorkflowRegistry{
		c: c.c.WorkflowRegistry,
	}
}

func (c *capabilitiesConfig) RateLimit() config.EngineExecutionRateLimit {
	return &engineExecutionRateLimit{
		rl: c.c.RateLimit,
	}
}

type engineExecutionRateLimit struct {
	rl toml.EngineExecutionRateLimit
}

func (rl *engineExecutionRateLimit) GlobalRPS() float64 {
	return *rl.rl.GlobalRPS
}

func (rl *engineExecutionRateLimit) GlobalBurst() int {
	return *rl.rl.GlobalBurst
}

func (rl *engineExecutionRateLimit) PerSenderRPS() float64 {
	return *rl.rl.PerSenderRPS
}

func (rl *engineExecutionRateLimit) PerSenderBurst() int {
	return *rl.rl.PerSenderBurst
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

func (c *capabilitiesConfig) GatewayConnector() config.GatewayConnector {
	return &gatewayConnector{
		c: c.c.GatewayConnector,
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

type capabilitiesWorkflowRegistry struct {
	c toml.WorkflowRegistry
}

func (c *capabilitiesWorkflowRegistry) RelayID() types.RelayID {
	return types.NewRelayID(c.NetworkID(), c.ChainID())
}

func (c *capabilitiesWorkflowRegistry) NetworkID() string {
	return *c.c.NetworkID
}

func (c *capabilitiesWorkflowRegistry) ChainID() string {
	return *c.c.ChainID
}

func (c *capabilitiesWorkflowRegistry) Address() string {
	return *c.c.Address
}

func (c *capabilitiesWorkflowRegistry) MaxEncryptedSecretsSize() utils.FileSize {
	return *c.c.MaxEncryptedSecretsSize
}

func (c *capabilitiesWorkflowRegistry) MaxBinarySize() utils.FileSize {
	return *c.c.MaxBinarySize
}

func (c *capabilitiesWorkflowRegistry) MaxConfigSize() utils.FileSize {
	return *c.c.MaxConfigSize
}

type gatewayConnector struct {
	c toml.GatewayConnector
}

func (c *gatewayConnector) ChainIDForNodeKey() string {
	return *c.c.ChainIDForNodeKey
}
func (c *gatewayConnector) NodeAddress() string {
	return *c.c.NodeAddress
}

func (c *gatewayConnector) DonID() string {
	return *c.c.DonID
}

func (c *gatewayConnector) Gateways() []config.ConnectorGateway {
	t := make([]config.ConnectorGateway, len(c.c.Gateways))
	for index, element := range c.c.Gateways {
		t[index] = &connectorGateway{element}
	}
	return t
}

func (c *gatewayConnector) WSHandshakeTimeoutMillis() uint32 {
	return *c.c.WSHandshakeTimeoutMillis
}

func (c *gatewayConnector) AuthMinChallengeLen() int {
	return *c.c.AuthMinChallengeLen
}

func (c *gatewayConnector) AuthTimestampToleranceSec() uint32 {
	return *c.c.AuthTimestampToleranceSec
}

type connectorGateway struct {
	c toml.ConnectorGateway
}

func (c *connectorGateway) ID() string {
	return *c.c.ID
}

func (c *connectorGateway) URL() string {
	return *c.c.URL
}
