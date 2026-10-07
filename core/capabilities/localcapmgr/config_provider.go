package localcapmgr

import (
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
)

// CapabilityConfigProvider supplies capability config overrides to merge into the config a
// capability is started with, keyed by capability ID and on-chain DON ID. It is the seam
// through which the offchain capabilities registry (GlobalConfig) is layered over on-chain config
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

// offchainCapabilityConfigProvider is backed by a snapshot of the offchain capabilities
// registry taken once per reconcile. It returns the offchain spec_config for a (capID, donID),
// matching the shape the TOML provider yields. The offchain registry is keyed by on-chain DON
// name, so donID is resolved through donNames (see offchainDONNames).
//
// Missing entries are not an error: when the DON has no usable name, the payload has no config
// for the DON, no entry for the capability, or no spec_config, it returns nil and every key
// keeps its on-chain/TOML value.
// A spec_config that cannot be converted is rejected at ingestion (globalconfig.Validate), so it
// cannot reach here from an applied payload; if it ever does, the offchain layer is skipped for
// that capability (on-chain/TOML values are used) and a warning is logged.
type offchainCapabilityConfigProvider struct {
	reg      *capabilitiespb.OffchainCapabilitiesRegistry
	donNames map[uint32]string
	version  uint64
	lggr     logger.Logger
}

func (p offchainCapabilityConfigProvider) LocalConfigOverrides(capID string, donID uint32) map[string]any {
	donName, ok := p.donNames[donID]
	if !ok {
		return nil
	}
	capCfg := p.reg.GetDons()[donName].GetCapabilities()[capID]
	out, err := globalconfig.SpecConfigMap(capCfg.GetSpecConfig())
	if err != nil {
		if p.lggr != nil {
			p.lggr.Warnw("Invalid offchain spec_config, ignoring offchain override",
				"capID", capID, "donID", donID, "donName", donName, "offchainVersion", p.version, "error", err)
		}
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
