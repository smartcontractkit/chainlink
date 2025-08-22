package config

import (
	"errors"
	"maps"
	"slices"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type Config struct {
	Blockchains       []blockchain.Input              `toml:"blockchains" validate:"required"`
	NodeSets          []*cre.CapabilitiesAwareNodeSet `toml:"nodesets" validate:"required"`
	JD                *jd.Input                       `toml:"jd" validate:"required"`
	Infra             *infra.Input                    `toml:"infra" validate:"required"`
	Fake              *fake.Input                     `toml:"fake" validate:"required"`
	S3ProviderInput   *s3provider.Input               `toml:"s3provider"`
	CapabilityConfigs map[string]cre.CapabilityConfig `toml:"capability_configs"` // capability flag -> capability config
}

// Validate performs validation checks on the configuration, ensuring all required fields
// are present and all referenced capabilities are known to the system.
func (c Config) Validate(capabilityFlagsProvider cre.CapabilityFlagsProvider) error {
	if c.JD.CSAEncryptionKey == "" {
		return errors.New("jd.csa_encryption_key must be provided")
	}

	for _, nodeSet := range c.NodeSets {
		for _, capability := range nodeSet.Capabilities {
			if !slices.Contains(capabilityFlagsProvider.SupportedCapabilityFlags(), capability) {
				return errors.New("unknown capability: " + capability + ". Make sure you have added it to the capabilityFlagsProvider")
			}
		}

		for capability := range nodeSet.ChainCapabilities {
			if !slices.Contains(capabilityFlagsProvider.SupportedCapabilityFlags(), capability) {
				return errors.New("unknown capability: " + capability + ". Make sure you have added it to the capabilityFlagsProvider")
			}
		}
	}

	return nil
}

// ResolveCapabilityForChain merges defaults with chain override for a capability on a given chain.
// Returns (enabled, mergedConfig).
func ResolveCapabilityForChain(
	capName string,
	caps map[string]*cre.ChainCapabilityConfig,
	defaults map[string]any,
	chainID uint64,
) (bool, map[string]any, error) {
	if caps == nil {
		return false, nil, nil
	}
	cfg, ok := caps[capName]
	if !ok {
		return false, nil, nil
	}
	enabled := slices.Contains(cfg.EnabledChains, chainID)
	if !enabled {
		return false, nil, nil
	}
	merged := map[string]any{}
	if defaults != nil {
		// copy defaults
		maps.Copy(merged, defaults)
	}
	if co, ok := cfg.ChainOverrides[chainID]; ok {
		// override with chain-specific values
		maps.Copy(merged, co)
	}
	return true, merged, nil
}

// ResolveCapabilityConfigForDON merges global defaults with DON-specific overrides for capabilities
// that don't have chain-specific configuration (like cron, web-api-target, web-api-trigger).
// Returns the merged configuration.
func ResolveCapabilityConfigForDON(
	capabilityName string,
	globalDefaults map[string]any,
	donOverrides map[string]map[string]any,
) map[string]any {
	merged := map[string]any{}

	// Start with global defaults
	if globalDefaults != nil {
		maps.Copy(merged, globalDefaults)
	}

	// Apply DON-specific overrides
	if donOverrides != nil {
		if overrides, ok := donOverrides[capabilityName]; ok {
			maps.Copy(merged, overrides)
		}
	}

	return merged
}
