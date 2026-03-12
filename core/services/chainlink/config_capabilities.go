package chainlink

import (
	"regexp"
	"sync"

	"github.com/smartcontractkit/libocr/commontypes"

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

func (c *capabilitiesConfig) SharedPeering() config.SharedPeering {
	return &sharedPeering{s: c.c.SharedPeering}
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

func (d *dispatcher) SendToSharedPeer() bool {
	return *d.d.SendToSharedPeer
}

type sharedPeering struct {
	s toml.SharedPeering
}

func (s *sharedPeering) Enabled() bool {
	return *s.s.Enabled
}

func (s *sharedPeering) Bootstrappers() (locators []commontypes.BootstrapperLocator) {
	if d := s.s.Bootstrappers; d != nil {
		return *d
	}
	return nil
}

func (s *sharedPeering) StreamConfig() config.StreamConfig {
	return &streamConfig{c: s.s.StreamConfig}
}

type streamConfig struct {
	c toml.StreamConfig
}

func (c *streamConfig) IncomingMessageBufferSize() int {
	return *c.c.IncomingMessageBufferSize
}

func (c *streamConfig) OutgoingMessageBufferSize() int {
	return *c.c.OutgoingMessageBufferSize
}

func (c *streamConfig) MaxMessageLenBytes() int {
	return *c.c.MaxMessageLenBytes
}

func (c *streamConfig) MessageRateLimiterRate() float64 {
	return *c.c.MessageRateLimiterRate
}

func (c *streamConfig) MessageRateLimiterCapacity() uint32 {
	return *c.c.MessageRateLimiterCapacity
}

func (c *streamConfig) BytesRateLimiterRate() float64 {
	return *c.c.BytesRateLimiterRate
}

func (c *streamConfig) BytesRateLimiterCapacity() uint32 {
	return *c.c.BytesRateLimiterCapacity
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

func (c *capabilitiesExternalRegistry) ContractVersion() string {
	return *c.c.ContractVersion
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

func (c *capabilitiesWorkflowRegistry) ContractVersion() string {
	return *c.c.ContractVersion
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

func (c *capabilitiesWorkflowRegistry) SyncStrategy() string {
	return *c.c.SyncStrategy
}

func (c *capabilitiesWorkflowRegistry) MaxConcurrency() int {
	return *c.c.MaxConcurrency
}

func (c *capabilitiesWorkflowRegistry) WorkflowStorage() config.WorkflowStorage {
	return &workflowStorage{
		c: c.c.WorkflowStorage,
	}
}

func (c *capabilitiesWorkflowRegistry) AdditionalSources() []config.AdditionalWorkflowSource {
	sources := make([]config.AdditionalWorkflowSource, len(c.c.AdditionalSourcesConfig))
	for i, src := range c.c.AdditionalSourcesConfig {
		sources[i] = &additionalWorkflowSource{c: src}
	}
	return sources
}

type workflowStorage struct {
	c toml.WorkflowStorage
}

func (c *workflowStorage) URL() string {
	return *c.c.URL
}

func (c *workflowStorage) TLSEnabled() bool {
	return *c.c.TLSEnabled
}

func (c *workflowStorage) ArtifactStorageHost() string {
	return *c.c.ArtifactStorageHost
}

type additionalWorkflowSource struct {
	c toml.AdditionalWorkflowSource
}

func (a *additionalWorkflowSource) GetURL() string {
	if a.c.URL == nil {
		return ""
	}
	return *a.c.URL
}

func (a *additionalWorkflowSource) GetTLSEnabled() bool {
	if a.c.TLSEnabled == nil {
		return true // Default to true
	}
	return *a.c.TLSEnabled
}

func (a *additionalWorkflowSource) GetName() string {
	if a.c.Name == nil {
		return ""
	}
	return *a.c.Name
}

type gatewayConnector struct {
	c toml.GatewayConnector
}

func (c *gatewayConnector) ChainIDForNodeKey() string {
	return *c.c.ChainIDForNodeKey
}

// NodeAddress is the address used to sign gateway handshakes.
// Empty string signals auto-discovery.
func (c *gatewayConnector) NodeAddress() string {
	if c.c.NodeAddress == nil || *c.c.NodeAddress == "" {
		return ""
	}
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

func (c *capabilitiesConfig) Local() config.LocalCapabilities {
	return &localCapabilities{c: c.c.Local}
}

var _ config.LocalCapabilities = (*localCapabilities)(nil)

type localCapabilities struct {
	c toml.LocalCapabilities

	// compiledRegexes caches compiled regex patterns from RegistryBasedLaunchAllowlist.
	// Lazily initialized on first call to IsAllowlisted.
	compiledRegexes []*regexp.Regexp
	regexOnce       sync.Once
}

func (l *localCapabilities) RegistryBasedLaunchAllowlist() []string {
	return l.c.RegistryBasedLaunchAllowlist
}

func (l *localCapabilities) Capabilities() map[string]config.CapabilityNodeConfig {
	if l.c.Capabilities == nil {
		return nil
	}
	result := make(map[string]config.CapabilityNodeConfig, len(l.c.Capabilities))
	for k, v := range l.c.Capabilities {
		result[k] = &capabilityNodeConfig{c: v}
	}
	return result
}

func (l *localCapabilities) compileRegexes() {
	l.compiledRegexes = make([]*regexp.Regexp, 0, len(l.c.RegistryBasedLaunchAllowlist))
	for _, pattern := range l.c.RegistryBasedLaunchAllowlist {
		// Patterns are validated at config load time, so compilation should not fail.
		// If it does fail (shouldn't happen), skip the invalid pattern.
		if re, err := regexp.Compile(pattern); err == nil {
			l.compiledRegexes = append(l.compiledRegexes, re)
		}
	}
}

func (l *localCapabilities) IsAllowlisted(capabilityID string) bool {
	l.regexOnce.Do(l.compileRegexes)
	for _, re := range l.compiledRegexes {
		if re.MatchString(capabilityID) {
			return true
		}
	}
	return false
}

func (l *localCapabilities) GetCapabilityConfig(capabilityID string) config.CapabilityNodeConfig {
	if l.c.Capabilities == nil {
		return nil
	}
	if c, ok := l.c.Capabilities[capabilityID]; ok {
		return &capabilityNodeConfig{c: c}
	}
	return nil
}

var _ config.CapabilityNodeConfig = (*capabilityNodeConfig)(nil)

type capabilityNodeConfig struct {
	c toml.CapabilityNodeConfig
}

func (c *capabilityNodeConfig) BinaryPathOverride() string {
	if c.c.BinaryPathOverride == nil {
		return ""
	}
	return *c.c.BinaryPathOverride
}

func (c *capabilityNodeConfig) Config() map[string]string {
	return c.c.Config
}
