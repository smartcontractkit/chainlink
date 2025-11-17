package crib

import (
	"testing"

	"k8s.io/utils/ptr"

	"github.com/stretchr/testify/require"
)

func TestGenerateBidirectionalRandomLanesWithMinConnectivity(t *testing.T) {
	tests := []struct {
		name         string
		chains       []uint64
		numLanes     int
		validateFunc func(t *testing.T, lanes []LaneConfig, chains []uint64, numLanes int)
	}{
		{
			name:         "Empty chains",
			chains:       []uint64{},
			numLanes:     0,
			validateFunc: validateEmptyResult,
		},
		{
			name:         "Single chain",
			chains:       []uint64{1},
			numLanes:     0,
			validateFunc: validateEmptyResult,
		},
		{
			name:         "Four chains - all possible lanes",
			chains:       []uint64{1, 2, 3, 4},
			numLanes:     12, // 4*3 = 12 total possible
			validateFunc: validateFullBidirectionalConnectivity,
		},
		{
			name:         "Large chain set with random selection",
			chains:       []uint64{1, 2, 3, 4, 5, 6},
			numLanes:     20,
			validateFunc: validatePartialBidirectionalConnectivity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lanes := generateBidirectionalRandomLanesWithMinConnectivity(tt.chains, tt.numLanes)
			tt.validateFunc(t, lanes, tt.chains, tt.numLanes)
		})
	}
}

func TestBidirectionalPairGeneration(t *testing.T) {
	chains := []uint64{1, 2, 3, 4}
	numLanes := 12 // full connectivity

	lanes := generateBidirectionalRandomLanesWithMinConnectivity(chains, numLanes)

	// Validate that we have complete bidirectional pairs
	bidirectionalPairs := findBidirectionalPairs(lanes)

	// With 4 chains and 12 lanes, we should have 6 bidirectional pairs
	require.Len(t, bidirectionalPairs, 6, "Should have exactly 3 bidirectional pairs")

	// Each chain should be reachable from every other chain
	require.True(t, isFullyConnected(lanes, chains))
}

// Validation helper functions

func validateEmptyResult(t *testing.T, lanes []LaneConfig, chains []uint64, numLanes int) {
	require.Empty(t, lanes, "Should return empty lanes for empty chains or zero numLanes")
}

func validateFullBidirectionalConnectivity(t *testing.T, lanes []LaneConfig, chains []uint64, numLanes int) {
	if len(chains) <= 1 {
		validateEmptyResult(t, lanes, chains, numLanes)
		return
	}

	// Basic validations
	require.LessOrEqual(t, len(lanes), numLanes, "Should not exceed requested number of lanes")

	// Validate minimum connectivity - each chain should be both source and destination
	validateMinimumConnectivity(t, lanes, chains)

	validateLaneBidirecionality(t, lanes)

	require.True(t, isFullyConnected(lanes, chains))
}

// validateLaneBidirecionality checks that each lane has a reverse pair and no self-loops
func validateLaneBidirecionality(t *testing.T, lanes []LaneConfig) {
	laneSet := make(map[LaneConfig]bool)
	// Validate no self-loops
	for _, lane := range lanes {
		laneSet[lane] = true
		require.NotEqual(t, lane.SourceChain, lane.DestinationChain, "No self-loops allowed")
	}

	// Validate that each lane has a reverse pair
	for _, lane := range lanes {
		reverseLane := LaneConfig{
			SourceChain:      lane.DestinationChain,
			DestinationChain: lane.SourceChain,
		}
		require.Contains(t, laneSet, reverseLane, "Each lane should have a reverse pair")
	}
}

func validatePartialBidirectionalConnectivity(t *testing.T, lanes []LaneConfig, chains []uint64, numLanes int) {
	if len(chains) <= 1 {
		validateEmptyResult(t, lanes, chains, numLanes)
		return
	}

	// Should not exceed requested lanes
	require.LessOrEqual(t, len(lanes), numLanes, "Should not exceed requested number of lanes")
	validateMinimumConnectivity(t, lanes, chains)
	validateLaneBidirecionality(t, lanes)
	require.False(t, isFullyConnected(lanes, chains), "Should not be fully connected in partial mode")
}

func validateMinimumConnectivity(t *testing.T, lanes []LaneConfig, chains []uint64) {
	if len(chains) <= 1 {
		return
	}

	sourceChains := make(map[uint64]bool)
	destChains := make(map[uint64]bool)

	for _, lane := range lanes {
		sourceChains[lane.SourceChain] = true
		destChains[lane.DestinationChain] = true
	}

	// Each chain should appear as both source and destination
	for _, chain := range chains {
		require.True(t, sourceChains[chain], "Chain %d should be a source", chain)
		require.True(t, destChains[chain], "Chain %d should be a destination", chain)
	}
}

func isFullyConnected(lanes []LaneConfig, chains []uint64) bool {
	// Build adjacency map
	adjacency := make(map[uint64]map[uint64]bool)
	for _, chain := range chains {
		adjacency[chain] = make(map[uint64]bool)
	}

	for _, lane := range lanes {
		adjacency[lane.SourceChain][lane.DestinationChain] = true
	}

	// Verify each chain can reach every other chain (directly)
	for _, src := range chains {
		for _, dst := range chains {
			if src != dst {
				_, exists := adjacency[src][dst]
				if !exists {
					return false
				}
			}
		}
	}

	return true
}

func findBidirectionalPairs(lanes []LaneConfig) [][]LaneConfig {
	// Create a map to find reverse lanes
	laneMap := make(map[LaneConfig]bool)
	for _, lane := range lanes {
		laneMap[lane] = true
	}

	var pairs [][]LaneConfig
	processed := make(map[LaneConfig]bool)

	for _, lane := range lanes {
		if processed[lane] {
			continue
		}

		reverseLane := LaneConfig{
			SourceChain:      lane.DestinationChain,
			DestinationChain: lane.SourceChain,
		}

		if laneMap[reverseLane] && !processed[reverseLane] {
			pairs = append(pairs, []LaneConfig{lane, reverseLane})
			processed[lane] = true
			processed[reverseLane] = true
		}
	}

	return pairs
}

func TestLaneConfiguration_GenerateLanes_BidirectionalMode(t *testing.T) {
	tests := []struct {
		name            string
		lc              *LaneConfiguration
		chains          []uint64
		expected        int
		validationError bool
	}{
		{
			name: "Random lanes with bidirectional",
			lc: &LaneConfiguration{
				Mode:     ptr.To(LaneModeRandomLanes),
				NumLanes: ptr.To(6),
			},
			chains:   []uint64{1, 2, 3},
			expected: 6,
		},
		{
			name: "Nil mode",
			lc: &LaneConfiguration{
				NumLanes: ptr.To(5),
			},
			chains:          []uint64{1, 2, 3, 4},
			validationError: true,
		},
		{
			name: "Random lanes with bidirectional - wrong lane count",
			lc: &LaneConfiguration{
				Mode:     ptr.To(LaneModeRandomLanes),
				NumLanes: ptr.To(5),
			},
			chains:          []uint64{1, 2, 3, 4},
			validationError: true,
		},
		{
			name: "Random lanes with bidirectional - odd lane count",
			lc: &LaneConfiguration{
				Mode:     ptr.To(LaneModeRandomLanes),
				NumLanes: ptr.To(9),
			},
			chains:   []uint64{1, 2, 3, 4},
			expected: 10, // requested 9, but should generate 10 to ensure all lanes are bidirectional
		},
		{
			name: "Any-to-any mode",
			lc: &LaneConfiguration{
				Mode: ptr.To(LaneModeAnyToAny),
			},
			chains:   []uint64{1, 2, 3},
			expected: 6, // 3*2 = 6 total lanes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.lc.Validate(len(tt.chains))
			if tt.validationError {
				require.Error(t, err,
					"Lane configuration validation should fail for invalid cases")
				return
			}

			require.NoError(t, err,
				"Lane configuration validation should not fail")

			lanes := tt.lc.GenerateLanes(tt.chains)

			require.Len(t, lanes, tt.expected)

			if tt.lc.Mode != nil && *tt.lc.Mode == LaneModeAnyToAny {
				validateFullBidirectionalConnectivity(t, lanes, tt.chains, tt.expected)
			}

			if tt.lc.Mode != nil && *tt.lc.Mode == LaneModeRandomLanes {
				if tt.expected == len(tt.chains)*(len(tt.chains)-1) {
					validateFullBidirectionalConnectivity(t, lanes, tt.chains, tt.expected)
				} else {
					validatePartialBidirectionalConnectivity(t, lanes, tt.chains, tt.expected)
				}
			}
		})
	}
}

func Test_generateChainTierLanes(t *testing.T) {
	chains := []uint64{3379446385462418246, 12463857294658392847, 12922642891491394802, 909606746561742123, 5548718428018410741}

	t.Run("happy path", func(t *testing.T) {
		lanes := generateChainTierLanes(chains, 2, 3)
		// Only 3379446385462418246 and 2 are sources, all are destinations except self
		expected := []LaneConfig{
			// chain 3379446385462418246 should have bidirectional lanes to all other chains
			{SourceChain: 3379446385462418246, DestinationChain: 12463857294658392847},
			{SourceChain: 3379446385462418246, DestinationChain: 12922642891491394802},
			{SourceChain: 3379446385462418246, DestinationChain: 909606746561742123},
			{SourceChain: 3379446385462418246, DestinationChain: 5548718428018410741},
			{SourceChain: 12922642891491394802, DestinationChain: 3379446385462418246},
			{SourceChain: 909606746561742123, DestinationChain: 3379446385462418246},
			{SourceChain: 5548718428018410741, DestinationChain: 3379446385462418246},
			// chain 2 should have bidirectional lanes to all other chains
			{SourceChain: 12463857294658392847, DestinationChain: 3379446385462418246},
			{SourceChain: 12463857294658392847, DestinationChain: 12922642891491394802},
			{SourceChain: 12463857294658392847, DestinationChain: 909606746561742123},
			{SourceChain: 12463857294658392847, DestinationChain: 5548718428018410741},
			{SourceChain: 12922642891491394802, DestinationChain: 12463857294658392847},
			{SourceChain: 909606746561742123, DestinationChain: 12463857294658392847},
			{SourceChain: 5548718428018410741, DestinationChain: 12463857294658392847},
		}
		require.ElementsMatch(t, expected, lanes)
		for _, lane := range lanes {
			require.NotEqual(t, lane.SourceChain, lane.DestinationChain, "no self-loops")
		}
	})

	t.Run("empty chains returns empty", func(t *testing.T) {
		lanes := generateChainTierLanes([]uint64{}, 0, 0)
		require.Empty(t, lanes)
	})
}
