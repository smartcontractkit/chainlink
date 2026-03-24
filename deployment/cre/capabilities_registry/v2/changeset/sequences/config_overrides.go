package sequences

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
)

// CapabilityConfigOverride specifies a per-DON override for capability configs.
// If CapabilityID is empty, the override is applied to all capabilities for that DON.
// Config is deep-merged into the base CapabilityConfig.Config.
type CapabilityConfigOverride struct {
	CapabilityID string         `json:"capabilityID" yaml:"capabilityID"`
	Config       map[string]any `json:"config" yaml:"config"`
}

// resolveCapabilityConfigsForDON returns a copy of baseConfigs with per-DON overrides deep-merged in.
// If an override has an empty CapabilityID, it applies to all capabilities.
// If an override has a CapabilityID, it applies only to the matching capability.
func resolveCapabilityConfigsForDON(baseConfigs []contracts.CapabilityConfig, overrides []CapabilityConfigOverride) ([]contracts.CapabilityConfig, error) {
	if len(overrides) == 0 {
		return baseConfigs, nil
	}

	result := make([]contracts.CapabilityConfig, len(baseConfigs))
	for i, base := range baseConfigs {
		result[i] = contracts.CapabilityConfig{
			Capability: base.Capability,
			Config:     deepCopyMap(base.Config),
		}
	}

	for _, override := range overrides {
		if override.Config == nil {
			continue
		}
		applied := false
		for i := range result {
			if override.CapabilityID == "" || override.CapabilityID == result[i].Capability.CapabilityID {
				result[i].Config = deepMergeMaps(result[i].Config, override.Config)
				applied = true
			}
		}
		if override.CapabilityID != "" && !applied {
			return nil, fmt.Errorf("override references capability ID %q which does not exist in the base capability configs", override.CapabilityID)
		}
	}

	return result, nil
}

// deepMergeMaps recursively merges override into base, returning a new map.
// For nested map[string]any values, it recurses. For all other types the override value wins.
// Neither input is mutated.
func deepMergeMaps(base, override map[string]any) map[string]any {
	if base == nil && override == nil {
		return nil
	}
	result := deepCopyMap(base)
	if result == nil {
		result = make(map[string]any)
	}
	for k, overrideVal := range override {
		baseVal, exists := result[k]
		if exists {
			baseMap, baseIsMap := baseVal.(map[string]any)
			overrideMap, overrideIsMap := overrideVal.(map[string]any)
			if baseIsMap && overrideIsMap {
				result[k] = deepMergeMaps(baseMap, overrideMap)
				continue
			}
		}
		result[k] = deepCopyValue(overrideVal)
	}
	return result
}

func deepCopyMap(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	result := make(map[string]any, len(m))
	for k, v := range m {
		result[k] = deepCopyValue(v)
	}
	return result
}

func deepCopyValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		return deepCopyMap(val)
	case []any:
		cp := make([]any, len(val))
		for i, item := range val {
			cp[i] = deepCopyValue(item)
		}
		return cp
	default:
		return v
	}
}
