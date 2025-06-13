package crib

import (
	"errors"
	"fmt"
	"github.com/AlekSi/pointer"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"math/rand"
	"sort"
	//"github.com/stretchr/testify/require"
)

// LaneConfig represents a unidirectional lane from source to destination
type LaneConfig struct {
	SourceChain      uint64 `toml:",omitempty"`
	DestinationChain uint64 `toml:",omitempty"`
}

// LaneConfiguration defines how lanes should be configured for the load test
type LaneConfiguration struct {
	// Mode determines how lanes are configured
	// "any-to-any" - traditional full mesh (default)
	// "random-lanes" - generate random lanes based on count
	Mode *string `toml:",omitempty"`

	// NumLanes - number of random lanes to generate when Mode is "random-lanes"
	NumLanes *int `toml:",omitempty"`

	// EnsureBidirectional - when true, ensures that for every A->B lane, there's also a B->A lane
	EnsureBidirectional *bool `toml:",omitempty"`

	// Internal fields for caching and deterministic generation
	generatedLanes []LaneConfig `toml:"-"` // Cache for generated lanes
}

const (
	LaneModeAnyToAny    = "any-to-any"
	LaneModeRandomLanes = "random-lanes"
)

func (lc *LaneConfiguration) Validate(e *cldf.Environment) error {
	if lc == nil {
		// Default to any-to-any if not specified
		return nil
	}

	mode := pointer.GetString(lc.Mode)
	if mode == "" {
		mode = LaneModeAnyToAny
	}

	chains := e.BlockChains.ListChainSelectors()
	chainSet := make(map[uint64]bool)
	for _, chain := range chains {
		chainSet[chain] = true
	}

	switch mode {
	case LaneModeAnyToAny:
		// No additional validation needed
		return nil

	case LaneModeRandomLanes:
		if lc.NumLanes == nil || *lc.NumLanes <= 0 {
			return errors.New("NumLanes must be provided and greater than 0 when Mode is 'random-lanes'")
		}

		maxPossibleLanes := len(chains) * (len(chains) - 1)
		if *lc.NumLanes > maxPossibleLanes {
			return fmt.Errorf("NumLanes (%d) cannot exceed maximum possible lanes (%d) for %d chains",
				*lc.NumLanes, maxPossibleLanes, len(chains))
		}

		// Calculate minimum lanes needed for connectivity
		minLanesNeeded := calculateMinimumLanesNeeded(len(chains), pointer.GetBool(lc.EnsureBidirectional))
		if *lc.NumLanes < minLanesNeeded {
			return fmt.Errorf("NumLanes (%d) is too low to ensure each chain is both source and destination. Minimum needed: %d",
				*lc.NumLanes, minLanesNeeded)
		}

	default:
		return fmt.Errorf("invalid Mode: %s. Must be one of: %s, %s, %s",
			mode, LaneModeAnyToAny, LaneModeRandomLanes)
	}

	return nil
}

// GetLanes returns the list of lanes based on the configuration
// This is the main entry point - it caches results for deterministic behavior
func (lc *LaneConfiguration) GetLanes() ([]LaneConfig, error) {
	if lc == nil {
		return nil, errors.New("lane configuration is nil")
	}

	if len(lc.generatedLanes) == 0 {
		return nil, errors.New("lanes have not been generated yet")
	}

	return lc.generatedLanes, nil
}

// GenerateLanes creates the list of lanes based on the configuration
func (lc *LaneConfiguration) GenerateLanes(chains []uint64) []LaneConfig {
	mode := pointer.GetString(lc.Mode)
	if mode == "" {
		mode = LaneModeAnyToAny
	}

	switch mode {
	case LaneModeAnyToAny:
		lc.generatedLanes = generateAnyToAnyLanes(chains)
		return lc.generatedLanes

	case LaneModeRandomLanes:
		if lc.NumLanes == nil {
			return []LaneConfig{}
		}

		lc.generatedLanes = generateRandomLanesWithMinConnectivity(chains, *lc.NumLanes, pointer.GetBool(lc.EnsureBidirectional))

		return lc.generatedLanes

	default:
		// Default to any-to-any if mode is not recognized
		lc.generatedLanes = generateAnyToAnyLanes(chains)
		return lc.generatedLanes
	}
}

// Helper functions for lane generation
func generateAnyToAnyLanes(chains []uint64) []LaneConfig {
	var lanes []LaneConfig

	for _, src := range chains {
		for _, dst := range chains {
			if src != dst {
				lanes = append(lanes, LaneConfig{
					SourceChain:      src,
					DestinationChain: dst,
				})
			}
		}
	}

	return lanes
}

func generateRandomLanesWithMinConnectivity(chains []uint64, numLanes int, bidirectional bool) []LaneConfig {
	rng := rand.New(rand.NewSource(rand.Int63()))

	// Step 1: Ensure minimum connectivity - each chain must be both source and destination
	var guaranteedLanes []LaneConfig

	// Shuffle chains for randomness in connectivity pattern
	shuffledChains := make([]uint64, len(chains))
	copy(shuffledChains, chains)
	rng.Shuffle(len(shuffledChains), func(i, j int) {
		shuffledChains[i], shuffledChains[j] = shuffledChains[j], shuffledChains[i]
	})

	// Create minimum connectivity: each chain as source at least once, each chain as destination at least once
	// We'll create a cycle to ensure connectivity while using minimal lanes
	for i := 0; i < len(shuffledChains); i++ {
		src := shuffledChains[i]
		dst := shuffledChains[(i+1)%len(shuffledChains)] // Create a cycle

		guaranteedLanes = append(guaranteedLanes, LaneConfig{
			SourceChain:      src,
			DestinationChain: dst,
		})
	}

	// If bidirectional, add reverse lanes for guaranteed connectivity
	if bidirectional {
		reverseLanes := make([]LaneConfig, len(guaranteedLanes))
		for i, lane := range guaranteedLanes {
			reverseLanes[i] = LaneConfig{
				SourceChain:      lane.DestinationChain,
				DestinationChain: lane.SourceChain,
			}
		}
		guaranteedLanes = append(guaranteedLanes, reverseLanes...)
	}

	// Step 2: Fill remaining slots with random lanes
	if numLanes <= len(guaranteedLanes) {
		// If requested lanes <= guaranteed lanes, just return a subset of guaranteed lanes
		if numLanes < len(guaranteedLanes) {
			return guaranteedLanes[:numLanes]
		}
		return guaranteedLanes
	}

	// Create set of used lanes to avoid duplicates
	usedLanes := make(map[string]bool)
	for _, lane := range guaranteedLanes {
		laneKey := fmt.Sprintf("%d->%d", lane.SourceChain, lane.DestinationChain)
		usedLanes[laneKey] = true
	}

	// Generate additional random lanes
	allPossibleLanes := generateAnyToAnyLanes(chains)
	var availableLanes []LaneConfig

	// Filter out already used lanes
	for _, lane := range allPossibleLanes {
		laneKey := fmt.Sprintf("%d->%d", lane.SourceChain, lane.DestinationChain)
		if !usedLanes[laneKey] {
			availableLanes = append(availableLanes, lane)
		}
	}

	// Shuffle available lanes
	rng.Shuffle(len(availableLanes), func(i, j int) {
		availableLanes[i], availableLanes[j] = availableLanes[j], availableLanes[i]
	})

	// Add random lanes until we reach numLanes
	remainingSlots := numLanes - len(guaranteedLanes)
	if remainingSlots > len(availableLanes) {
		remainingSlots = len(availableLanes)
	}

	finalLanes := append(guaranteedLanes, availableLanes[:remainingSlots]...)

	return finalLanes
}

// calculateMinimumLanesNeeded calculates minimum lanes needed for connectivity
func calculateMinimumLanesNeeded(numChains int, bidirectional bool) int {
	if numChains <= 1 {
		return 0
	}

	// Minimum is a cycle: each chain -> next chain
	minLanes := numChains

	if bidirectional {
		// If bidirectional, we need reverse lanes too
		minLanes *= 2
	}

	return minLanes
}

// GetConnectedChains returns all chains that are involved in the configured lanes
func (lc *LaneConfiguration) GetConnectedChains() []uint64 {
	lanes, err := lc.GetLanes()
	if err != nil {
		return nil
	}

	chainSet := make(map[uint64]bool)
	for _, lane := range lanes {
		chainSet[lane.SourceChain] = true
		chainSet[lane.DestinationChain] = true
	}

	var connectedChains []uint64
	for chain := range chainSet {
		connectedChains = append(connectedChains, chain)
	}

	// Sort for deterministic order
	sort.Slice(connectedChains, func(i, j int) bool {
		return connectedChains[i] < connectedChains[j]
	})

	return connectedChains
}

// GetSourceChainsForDestination returns all source chains that can send to a specific destination
func (lc *LaneConfiguration) GetSourceChainsForDestination(destination uint64) []uint64 {
	lanes, err := lc.GetLanes()
	if err != nil {
		return nil
	}

	var sources []uint64
	for _, lane := range lanes {
		if lane.DestinationChain == destination {
			sources = append(sources, lane.SourceChain)
		}
	}

	// Sort for deterministic order
	sort.Slice(sources, func(i, j int) bool {
		return sources[i] < sources[j]
	})

	return sources
}

// GetDestinationChainsForSource returns all destination chains that a source can send to
func (lc *LaneConfiguration) GetDestinationChainsForSource(source uint64) []uint64 {
	lanes, err := lc.GetLanes()
	if err != nil {
		return nil
	}

	var destinations []uint64
	for _, lane := range lanes {
		if lane.SourceChain == source {
			destinations = append(destinations, lane.DestinationChain)
		}
	}

	// Sort for deterministic order
	sort.Slice(destinations, func(i, j int) bool {
		return destinations[i] < destinations[j]
	})

	return destinations
}

// Example TOML configurations:

// Any-to-any (traditional full mesh)
/*
Mode = "any-to-any"
*/

// Random lanes - deterministic based on configuration
/*

Mode = "random-lanes"
NumLanes = 350
EnsureBidirectional = true
*/
