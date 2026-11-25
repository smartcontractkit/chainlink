package contracts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
)

func TestComputeConfigs_Removals(t *testing.T) {
	t.Parallel()

	cfg := func(id string) CapabilityConfig {
		return CapabilityConfig{
			Capability: Capability{CapabilityID: id},
			Config:     map[string]any{"defaultConfig": map[string]any{}},
		}
	}
	onChain := func(ids ...string) []capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration {
		var out []capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration
		for _, id := range ids {
			out = append(out, capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration{CapabilityId: id, Config: []byte{0x1}})
		}
		return out
	}

	tests := []struct {
		name      string
		cfgs      []CapabilityConfig
		existing  []capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration
		merge     bool
		removeIDs []string
		wantIDs   []string
		wantErr   []string
	}{
		{
			name:     "merge=false with full list is ok",
			cfgs:     []CapabilityConfig{cfg("a@1"), cfg("b@1")},
			existing: onChain("a@1", "b@1"),
			wantIDs:  []string{"a@1", "b@1"},
		},
		{
			name:     "merge=false adding a new capability while keeping all is ok",
			cfgs:     []CapabilityConfig{cfg("a@1"), cfg("b@1"), cfg("c@1")},
			existing: onChain("a@1", "b@1"),
			wantIDs:  []string{"a@1", "b@1", "c@1"},
		},
		{
			name:     "merge=false missing one on-chain capability fails and names it",
			cfgs:     []CapabilityConfig{cfg("a@1")},
			existing: onChain("a@1", "b@1"),
			wantErr:  []string{"b@1", "removeCapabilityIDs", "mergeCapabilityConfigsWithOnChain"},
		},
		{
			name:      "merge=false missing one but named for removal removes it",
			cfgs:      []CapabilityConfig{cfg("a@1")},
			existing:  onChain("a@1", "b@1"),
			removeIDs: []string{"b@1"},
			wantIDs:   []string{"a@1"},
		},
		{
			name:     "empty list with merge=false on a DON with capabilities fails",
			existing: onChain("a@1", "b@1"),
			wantErr:  []string{"a@1", "b@1"},
		},
		{
			name:      "empty list with every capability named for removal is allowed",
			existing:  onChain("a@1", "b@1"),
			removeIDs: []string{"a@1", "b@1"},
		},
		{
			name:     "empty list with merge=false on a DON without capabilities is ok",
			existing: nil,
		},
		{
			name:     "merge=true never removes",
			cfgs:     []CapabilityConfig{cfg("a@1")},
			existing: onChain("a@1", "b@1"),
			merge:    true,
			wantIDs:  []string{"a@1", "b@1"},
		},
		{
			name:      "merge=true still honours an explicit removal",
			cfgs:      []CapabilityConfig{cfg("a@1")},
			existing:  onChain("a@1", "b@1"),
			merge:     true,
			removeIDs: []string{"b@1"},
			wantIDs:   []string{"a@1"},
		},
		{
			name:      "removal naming a capability not on the DON fails",
			cfgs:      []CapabilityConfig{cfg("a@1")},
			existing:  onChain("a@1"),
			removeIDs: []string{"zzz@1"},
			wantErr:   []string{"zzz@1", "not on the DON"},
		},
		{
			name:      "capability both in configs and removal list fails",
			cfgs:      []CapabilityConfig{cfg("a@1"), cfg("b@1")},
			existing:  onChain("a@1", "b@1"),
			removeIDs: []string{"b@1"},
			wantErr:   []string{"b@1", "both"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := computeConfigs(tc.cfgs, tc.existing, tc.merge, tc.removeIDs)
			if len(tc.wantErr) > 0 {
				require.Error(t, err)
				for _, s := range tc.wantErr {
					assert.Contains(t, err.Error(), s)
				}
				return
			}
			require.NoError(t, err)
			var ids []string
			for _, c := range got {
				ids = append(ids, c.CapabilityId)
			}
			assert.Equal(t, tc.wantIDs, ids)
		})
	}
}
