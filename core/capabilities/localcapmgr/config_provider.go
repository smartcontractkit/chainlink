package localcapmgr

import (
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

// CapabilityConfigProvider supplies capability config overrides to merge into the config a
// capability is started with, keyed by capability ID and on-chain DON ID. It is the seam
// through which the offchain capabilities registry (GlobalConfig) is layered over node TOML
// ([Capabilities.Local]) without changing the LocalCapabilityManager. See the Offchain
// Capabilities Registry design.
//
// Only the capability config map flows through this seam. Node-local infra (binary path
// override) and node role (registry-based launch allowlist) remain sourced from TOML and are
// read directly from config.LocalCapabilities, not through this provider.
type CapabilityConfigProvider interface {
	// LocalConfigOverrides returns config key/values to merge for the given capability on the
	// given DON, or nil when there is no override. The DON ID is honored by the offchain
	// provider (whose config is DON-scoped); the TOML provider ignores it.
	LocalConfigOverrides(capID string, donID uint32) map[string]any
}

// tomlCapabilityConfigProvider is the default provider, backed by node TOML config. TOML
// config is not DON-scoped, so donID is ignored.
type tomlCapabilityConfigProvider struct {
	localCfg config.LocalCapabilities
}

func (p tomlCapabilityConfigProvider) LocalConfigOverrides(capID string, _ uint32) map[string]any {
	if p.localCfg == nil {
		return nil
	}
	capCfg := p.localCfg.GetCapabilityConfig(capID)
	if capCfg == nil {
		return nil
	}
	return toAnyMap(capCfg.Config())
}

// offchainCapabilityConfigProvider is backed by the offchain capabilities registry. It returns
// the offchain spec_config for a (capID, donID), matching the shape the TOML provider yields.
type offchainCapabilityConfigProvider struct {
	registry *globalconfig.GlobalConfig
	lggr     logger.Logger
}

func (p offchainCapabilityConfigProvider) LocalConfigOverrides(capID string, donID uint32) map[string]any {
	if p.registry == nil {
		return nil
	}
	reg, _ := p.registry.LoadParsed()
	if reg == nil {
		return nil
	}
	don := reg.GetDons()[donID]
	if don == nil {
		return nil
	}
	capCfg := don.GetCapabilityConfigs()[capID]
	if capCfg == nil {
		return nil
	}
	sc := capCfg.GetSpecConfig()
	if sc == nil {
		return nil
	}
	m, err := values.FromMapValueProto(sc)
	if err != nil || m == nil {
		if p.lggr != nil {
			p.lggr.Warnw("Failed to convert offchain spec_config, ignoring offchain override",
				"capID", capID, "donID", donID, "error", err)
		}
		return nil
	}
	unwrapped, err := m.Unwrap()
	if err != nil {
		if p.lggr != nil {
			p.lggr.Warnw("Failed to unwrap offchain spec_config, ignoring offchain override",
				"capID", capID, "donID", donID, "error", err)
		}
		return nil
	}
	out, ok := unwrapped.(map[string]any)
	if !ok {
		return nil
	}
	return out
}

func toAnyMap(m map[string]string) map[string]any {
	if m == nil {
		return nil
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
