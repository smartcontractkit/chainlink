package crib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTier struct {
	NumChains int
}

type mockChainTiers struct {
	Tiers []mockTier
}

func toChainTiers(mt mockChainTiers) *ChainTiers {
	tiers := make([]TierConfigs, len(mt.Tiers))
	for i, t := range mt.Tiers {
		tiers[i] = TierConfigs{NumChains: t.NumChains}
	}
	return &ChainTiers{Tiers: tiers}
}

func Test_getTierChainSelectors(t *testing.T) {
	otherSelectors := []uint64{909606746561742123, 5548718428018410741, 789068866484373046, 5721565186521185178, 964127714438319834}
	defaultSelectors := []uint64{3379446385462418246, 12463857294658392847, 12922642891491394802}
	allSelectors := append(otherSelectors, defaultSelectors...)

	tests := []struct {
		name           string
		inputSelectors []uint64
		chainTiers     mockChainTiers
		wantTiers      [][]uint64
	}{
		{
			name:           "single tier, all selectors",
			inputSelectors: allSelectors,
			chainTiers: mockChainTiers{
				Tiers: []mockTier{{NumChains: len(allSelectors)}},
			},
			wantTiers: [][]uint64{{3379446385462418246, 12463857294658392847, 12922642891491394802, 909606746561742123, 5548718428018410741, 789068866484373046, 5721565186521185178, 964127714438319834}},
		},
		{
			name:           "two tiers, split",
			inputSelectors: allSelectors,
			chainTiers: mockChainTiers{
				Tiers: []mockTier{{NumChains: 4}, {NumChains: 4}},
			},
			wantTiers: [][]uint64{
				{3379446385462418246, 12463857294658392847, 12922642891491394802, 909606746561742123},
				{5548718428018410741, 789068866484373046, 5721565186521185178, 964127714438319834},
			},
		},
		{
			name:           "fewer than priority selectors ",
			inputSelectors: allSelectors,
			chainTiers: mockChainTiers{
				Tiers: []mockTier{{NumChains: 2}, {NumChains: 6}},
			},
			wantTiers: [][]uint64{
				{3379446385462418246, 12463857294658392847},
				{12922642891491394802, 909606746561742123, 5548718428018410741, 789068866484373046, 5721565186521185178, 964127714438319834},
			},
		},
		{
			name:           "fewer than all selectors ",
			inputSelectors: []uint64{12463857294658392847, 3379446385462418246, 12922642891491394802, 909606746561742123, 5548718428018410741},
			chainTiers: mockChainTiers{
				Tiers: []mockTier{{NumChains: 3}, {NumChains: 2}},
			},
			wantTiers: [][]uint64{
				{3379446385462418246, 12463857294658392847, 12922642891491394802},
				{909606746561742123, 5548718428018410741},
			},
		},
		{
			name:           "evm only",
			inputSelectors: []uint64{3379446385462418246, 12922642891491394802, 909606746561742123, 5548718428018410741, 789068866484373046},
			chainTiers: mockChainTiers{
				Tiers: []mockTier{{NumChains: 3}, {NumChains: 2}},
			},
			wantTiers: [][]uint64{
				{3379446385462418246, 12922642891491394802, 909606746561742123},
				{5548718428018410741, 789068866484373046},
			},
		},
		{
			name:           "no tiers",
			inputSelectors: []uint64{3379446385462418246, 12922642891491394802, 909606746561742123, 5548718428018410741, 789068866484373046},
			chainTiers: mockChainTiers{
				Tiers: nil,
			},
			wantTiers: [][]uint64{},
		},
		{
			name:           "three tiers",
			inputSelectors: []uint64{3379446385462418246, 12922642891491394802, 909606746561742123, 5548718428018410741, 789068866484373046},
			chainTiers: mockChainTiers{
				Tiers: []mockTier{{NumChains: 2}, {NumChains: 2}, {NumChains: 1}},
			},
			wantTiers: [][]uint64{
				{3379446385462418246, 12922642891491394802},
				{909606746561742123, 5548718428018410741},
				{789068866484373046},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getTierChainSelectors(tt.inputSelectors, toChainTiers(tt.chainTiers))
			assert.Equal(t, tt.wantTiers, got)
		})
	}
}
