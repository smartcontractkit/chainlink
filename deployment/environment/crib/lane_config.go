package crib

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	solState "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
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

// calculateMinimumLanesNeeded calculates minimum lanes needed for connectivity where each chain
// must be both a source and destination.
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

// DiscoverLanesFromDeployedState reverse engineers the lane configuration from deployed state
func (lc *LaneConfiguration) DiscoverLanesFromDeployedState(env cldf.Environment, state *stateview.CCIPOnChainState) error {
	var discoveredLanes []LaneConfig

	evmChains := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(selectors.FamilyEVM))
	solChains := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(selectors.FamilySolana))
	allChains := append(evmChains, solChains...)

	// Discover EVM to EVM lanes
	for _, srcChain := range evmChains {
		srcChainState, exists := state.Chains[srcChain]
		if !exists {
			continue
		}

		// Check which destination chains are configured on the OnRamp
		destinations, err := lc.getEnabledDestinationsFromOnRamp(srcChainState, allChains)
		if err != nil {
			return fmt.Errorf("failed to get enabled destinations for EVM chain %d: %w", srcChain, err)
		}

		for _, dstChain := range destinations {
			discoveredLanes = append(discoveredLanes, LaneConfig{
				SourceChain:      srcChain,
				DestinationChain: dstChain,
			})
		}
	}

	// Discover Solana to EVM lanes
	for _, srcChain := range solChains {
		srcChainState, exists := state.SolChains[srcChain]
		if !exists {
			continue
		}

		// Check which EVM destination chains are configured on the Solana Router
		destinations, err := lc.getEnabledDestinationsFromSolanaRouter(srcChainState, allChains)
		if err != nil {
			return fmt.Errorf("failed to get enabled EVM destinations for Solana chain %d: %w", srcChain, err)
		}

		for _, dstChain := range destinations {
			discoveredLanes = append(discoveredLanes, LaneConfig{
				SourceChain:      srcChain,
				DestinationChain: dstChain,
			})
		}
	}

	// Sort lanes for deterministic behavior
	sort.Slice(discoveredLanes, func(i, j int) bool {
		if discoveredLanes[i].SourceChain != discoveredLanes[j].SourceChain {
			return discoveredLanes[i].SourceChain < discoveredLanes[j].SourceChain
		}
		return discoveredLanes[i].DestinationChain < discoveredLanes[j].DestinationChain
	})

	// Store discovered lanes in the same field used by deployment configuration
	lc.generatedLanes = discoveredLanes
	return nil
}

// getEnabledDestinationsFromOnRamp checks which destinations are enabled on the OnRamp
func (lc *LaneConfiguration) getEnabledDestinationsFromOnRamp(chainState evm.CCIPChainState, candidateDestinations []uint64) ([]uint64, error) {
	var enabledDestinations []uint64

	// For each candidate destination, check if it's enabled on the OnRamp
	for _, dstChain := range candidateDestinations {
		isEnabled, err := lc.isDestinationEnabledOnOnRamp(chainState, dstChain)
		if err != nil {
			// Log but continue - some destinations might not be configured
			continue
		}

		if isEnabled {
			enabledDestinations = append(enabledDestinations, dstChain)
		}
	}

	return enabledDestinations, nil
}

// getEnabledDestinationsFromSolanaRouter checks which destinations are enabled on the Solana Router
func (lc *LaneConfiguration) getEnabledDestinationsFromSolanaRouter(chainState solState.CCIPChainState, candidateDestinations []uint64) ([]uint64, error) {
	var enabledDestinations []uint64

	// For each candidate destination, check if it's enabled on the Solana Router
	for _, dstChain := range candidateDestinations {
		isEnabled, err := lc.isDestinationEnabledOnSolanaRouter(chainState, dstChain)
		if err != nil {
			// Log but continue - some destinations might not be configured
			continue
		}

		if isEnabled {
			enabledDestinations = append(enabledDestinations, dstChain)
		}
	}

	return enabledDestinations, nil
}

// isDestinationEnabledOnOnRamp checks if a destination is enabled on the EVM OnRamp
func (lc *LaneConfiguration) isDestinationEnabledOnOnRamp(chainState evm.CCIPChainState, destinationChain uint64) (bool, error) {
	destConfig, err := chainState.OnRamp.GetDestChainConfig(&bind.CallOpts{}, destinationChain)
	if err != nil {
		// If we can't get the config, assume it's not enabled
		return false, err
	}

	// Check if the destination is enabled (router address should not be zero)
	return destConfig.Router != common.HexToAddress("0x0"), nil
}

// isDestinationEnabledOnSolanaRouter checks if a destination is enabled on the Solana Router
func (lc *LaneConfiguration) isDestinationEnabledOnSolanaRouter(chainState solState.CCIPChainState, destinationChain uint64) (bool, error) {
	panic("isDestinationEnabledOnSolanaRouter not implemented yet") // TODO: Implement this function
}

// GetSourceChainsForDestination returns all source chains that can send to a specific destination
func (lc *LaneConfiguration) GetSourceChainsForDestination(destination uint64) []uint64 {
	if lc == nil {
		return nil
	}

	var sources []uint64
	for _, lane := range lc.generatedLanes {
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
	if lc == nil {
		return nil
	}

	var destinations []uint64
	for _, lane := range lc.generatedLanes {
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

// LaneStats provides statistics about the discovered lane configuration
type LaneStats struct {
	TotalLanes        int
	UniqueChains      int
	AvgLanesPerChain  float64
	MaxLanesPerChain  int
	MinLanesPerChain  int
	SourceChains      int
	DestinationChains int
}

// GetLaneStats For metrics and reporting on the lane configuration
func (lc *LaneConfiguration) GetLaneStats() LaneStats {
	if lc == nil {
		return LaneStats{}
	}

	chainLaneCount := make(map[uint64]int)
	sourceChains := make(map[uint64]bool)
	destChains := make(map[uint64]bool)

	for _, lane := range lc.generatedLanes {
		chainLaneCount[lane.SourceChain]++
		chainLaneCount[lane.DestinationChain]++
		sourceChains[lane.SourceChain] = true
		destChains[lane.DestinationChain] = true
	}

	stats := LaneStats{
		TotalLanes:        len(lc.generatedLanes),
		UniqueChains:      len(chainLaneCount),
		SourceChains:      len(sourceChains),
		DestinationChains: len(destChains),
	}

	if len(chainLaneCount) > 0 {
		total := 0
		max := 0
		min := int(^uint(0) >> 1) // max int

		for _, count := range chainLaneCount {
			total += count
			if count > max {
				max = count
			}
			if count < min {
				min = count
			}
		}

		stats.AvgLanesPerChain = float64(total) / float64(len(chainLaneCount))
		stats.MaxLanesPerChain = max
		stats.MinLanesPerChain = min
	}

	return stats
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
