package localcapmgr

import "github.com/smartcontractkit/chainlink/v2/core/config"

// CapabilityConfigProvider supplies node-local capability config overrides, keyed by
// capability ID. Today it is backed by node TOML ([Capabilities.Local]); it is the seam
// through which the offchain capabilities registry (GlobalConfig) will later be layered
// (offchain-wins) without changing the LocalCapabilityManager. See the Offchain
// Capabilities Registry design.
//
// Only the capability config map moves offchain. Node-local infra (binary path override)
// and node role (registry-based launch allowlist) remain sourced from TOML and are read
// directly from config.LocalCapabilities, not through this provider.
type CapabilityConfigProvider interface {
	// LocalConfigOverrides returns config key/values to merge for the given capability,
	// or nil when there is no override.
	LocalConfigOverrides(capID string) map[string]string
}

// tomlCapabilityConfigProvider is the default provider, backed by node TOML config.
type tomlCapabilityConfigProvider struct {
	localCfg config.LocalCapabilities
}

func (p tomlCapabilityConfigProvider) LocalConfigOverrides(capID string) map[string]string {
	if p.localCfg == nil {
		return nil
	}
	capCfg := p.localCfg.GetCapabilityConfig(capID)
	if capCfg == nil {
		return nil
	}
	return capCfg.Config()
}
