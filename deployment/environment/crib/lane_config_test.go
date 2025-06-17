package crib

import (
	"k8s.io/utils/pointer"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRandomLanesWithMinConnectivity(t *testing.T) {
	rand.Seed(12345)

	tests := []struct {
		name         string
		chains       []uint64
		numLanes     int
		shouldError  bool
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
		{
			name:         "Request fewer lanes than minimum",
			chains:       []uint64{1, 2, 3},
			numLanes:     4, // Less than minimum 6
			validateFunc: validatePartialBidirectionalConnectivity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldError {
				require.Panics(t, func() {
					generateRandomLanesWithMinConnectivity(tt.chains, tt.numLanes)
				})
				return
			}

			lanes := generateRandomLanesWithMinConnectivity(tt.chains, tt.numLanes)
			tt.validateFunc(t, lanes, tt.chains, tt.numLanes)
		})
	}
}

func TestBidirectionalPairGeneration(t *testing.T) {
	chains := []uint64{1, 2, 3, 4}
	numLanes := 12 // full connectivity

	lanes := generateRandomLanesWithMinConnectivity(chains, numLanes)

	// Validate that we have complete bidirectional pairs
	bidirectionalPairs := findBidirectionalPairs(lanes)

	// With 4 chains and 12 lanes, we should have 6 bidirectional pairs
	require.Equal(t, 6, len(bidirectionalPairs), "Should have exactly 3 bidirectional pairs")

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
		name     string
		lc       *LaneConfiguration
		chains   []uint64
		expected int
	}{
		{
			name: "Random lanes with bidirectional",
			lc: &LaneConfiguration{
				Mode:     pointer.String(LaneModeRandomLanes),
				NumLanes: pointer.Int(6),
			},
			chains:   []uint64{1, 2, 3},
			expected: 6,
		},
		{
			name: "Any-to-any mode",
			lc: &LaneConfiguration{
				Mode: pointer.String(LaneModeAnyToAny),
			},
			chains:   []uint64{1, 2, 3},
			expected: 6, // 3*2 = 6 total lanes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lanes := tt.lc.GenerateLanes(tt.chains)

			// Validate basic properties
			require.Equal(t, len(lanes), tt.expected)

			if tt.lc.Mode != nil && *tt.lc.Mode == LaneModeAnyToAny {
				validateFullBidirectionalConnectivity(t, lanes, tt.chains, tt.expected)
			}
		})
	}
}
