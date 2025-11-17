package crib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getTierChainSelectors(t *testing.T) {
	otherSelectors := []uint64{909606746561742123, 5548718428018410741, 789068866484373046, 5721565186521185178, 964127714438319834}
	defaultSelectors := []uint64{3379446385462418246, 12463857294658392847, 12922642891491394802}
	//nolint:gocritic // append is used here to combine the two slices for testing purposes
	allSelectors := append(otherSelectors, defaultSelectors...)

	tests := []struct {
		name                 string
		inputSelectors       []uint64
		highTierCount        int
		lowTierCount         int
		expectedHighTierSels []uint64
		expectedLowTierSels  []uint64
	}{
		{
			name:                 "single tier, all selectors",
			inputSelectors:       allSelectors,
			highTierCount:        len(allSelectors),
			expectedHighTierSels: []uint64{3379446385462418246, 12463857294658392847, 12922642891491394802, 909606746561742123, 5548718428018410741, 789068866484373046, 5721565186521185178, 964127714438319834},
			expectedLowTierSels:  []uint64{},
		},
		{
			name:                 "two tiers, split",
			inputSelectors:       allSelectors,
			highTierCount:        4,
			lowTierCount:         4,
			expectedHighTierSels: []uint64{3379446385462418246, 12463857294658392847, 12922642891491394802, 909606746561742123},
			expectedLowTierSels:  []uint64{5548718428018410741, 789068866484373046, 5721565186521185178, 964127714438319834},
		},
		{
			name:                 "fewer than priority selectors ",
			inputSelectors:       allSelectors,
			highTierCount:        2,
			lowTierCount:         6,
			expectedHighTierSels: []uint64{3379446385462418246, 12463857294658392847},
			expectedLowTierSels:  []uint64{12922642891491394802, 909606746561742123, 5548718428018410741, 789068866484373046, 5721565186521185178, 964127714438319834},
		},
		{
			name:                 "fewer than all selectors ",
			inputSelectors:       []uint64{12463857294658392847, 3379446385462418246, 12922642891491394802, 909606746561742123, 5548718428018410741},
			highTierCount:        3,
			lowTierCount:         2,
			expectedHighTierSels: []uint64{3379446385462418246, 12463857294658392847, 12922642891491394802},
			expectedLowTierSels:  []uint64{909606746561742123, 5548718428018410741},
		},
		{
			name:                 "evm only",
			inputSelectors:       []uint64{3379446385462418246, 12922642891491394802, 909606746561742123, 5548718428018410741, 789068866484373046},
			highTierCount:        3,
			lowTierCount:         2,
			expectedHighTierSels: []uint64{3379446385462418246, 12922642891491394802, 909606746561742123},
			expectedLowTierSels:  []uint64{5548718428018410741, 789068866484373046},
		},
		{
			name:                 "no tiers",
			inputSelectors:       []uint64{3379446385462418246, 12922642891491394802, 909606746561742123, 5548718428018410741, 789068866484373046},
			expectedHighTierSels: []uint64{},
			expectedLowTierSels:  []uint64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			highTierSels, lowTierSels := getTierChainSelectors(tt.inputSelectors, tt.highTierCount, tt.lowTierCount)
			assert.Equal(t, tt.expectedHighTierSels, highTierSels)
			assert.Equal(t, tt.expectedLowTierSels, lowTierSels)
		})
	}
}
