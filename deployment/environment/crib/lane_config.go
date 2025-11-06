package crib

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"sort"

	"github.com/AlekSi/pointer"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_onramp"

	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_router"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	ccipSolState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	aptosState "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/aptos"
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
	// "tiered-lanes" - generate lanes with priority to higher tiered chains
	Mode *string `toml:",omitempty"`

	// NumLanes - number of random lanes to generate when Mode is "random-lanes"
	NumLanes *int `toml:",omitempty"`

	// Internal fields for caching
	generatedLanes []LaneConfig

	HighTierChainCount *int `toml:",omitempty"` // Optional, used for tiered lanes to specify how many chains are in the high tier
	LowTierChainCount  *int `toml:",omitempty"` // Optional, used for tiered lanes to specify how many chains are in the low tier
}

const (
	LaneModeAnyToAny    = "any-to-any"
	LaneModeRandomLanes = "random-lanes"
	LaneModeChainTiers  = "tiered-lanes"
)

// Validate checks the lane configuration for correctness, ensuring that
// the mode is set and that the number of lanes is valid for the given mode based on the expected number of chains.
func (lc *LaneConfiguration) Validate(chainCount int) error {
	if lc == nil {
		return errors.New("lane configuration is nil")
	}

	mode := pointer.GetString(lc.Mode)
	if mode == "" {
		return errors.New("mode must be set in LaneConfiguration")
	}

	switch mode {
	case LaneModeAnyToAny:
		// No additional validation needed
		return nil
	case LaneModeChainTiers:
		if lc.HighTierChainCount == nil || lc.LowTierChainCount == nil {
			return errors.New("HighTierChainCount and LowTierChainCount must be provided when Mode is 'tiered-lanes'")
		}

		if *lc.HighTierChainCount+*lc.LowTierChainCount != chainCount {
			return fmt.Errorf("HighTierChainCount (%d) + LowTierChainCount (%d) must equal total chain count (%d)",
				*lc.HighTierChainCount, *lc.LowTierChainCount, chainCount)
		}
	case LaneModeRandomLanes:
		if lc.NumLanes == nil || *lc.NumLanes <= 0 {
			return errors.New("NumLanes must be provided and greater than 0 when Mode is 'random-lanes'")
		}

		maxPossibleLanes := chainCount * (chainCount - 1)
		if *lc.NumLanes > maxPossibleLanes {
			return fmt.Errorf("NumLanes (%d) cannot exceed maximum possible lanes (%d) for %d chains",
				*lc.NumLanes, maxPossibleLanes, chainCount)
		}

		// Calculate minimum lanes needed for connectivity
		minLanesNeeded := calculateMinimumLanesNeeded(chainCount)
		if *lc.NumLanes < minLanesNeeded {
			return fmt.Errorf("NumLanes (%d) is too low to ensure each chain is both source and destination"+
				"bidirecionally. Minimum needed: %d",
				*lc.NumLanes, minLanesNeeded)
		}

	default:
		return fmt.Errorf("invalid Mode: %s. Must be one of: %s, %s, %s",
			mode, LaneModeAnyToAny, LaneModeRandomLanes, LaneModeChainTiers)
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
		panic("LaneConfiguration mode is not set, cannot generate lanes")
	}

	if lc.generatedLanes != nil {
		// If lanes are already generated, return cached result
		return lc.generatedLanes
	}

	switch mode {
	case LaneModeAnyToAny:
		lc.generatedLanes = generateAnyToAnyLanes(chains)
		return lc.generatedLanes
	case LaneModeChainTiers:
		lc.generatedLanes = generateChainTierLanes(chains, *lc.HighTierChainCount, *lc.LowTierChainCount)
		return lc.generatedLanes
	case LaneModeRandomLanes:
		if lc.NumLanes == nil {
			return []LaneConfig{}
		}

		lc.generatedLanes = generateBidirectionalRandomLanesWithMinConnectivity(chains, *lc.NumLanes)

		return lc.generatedLanes

	default:
		// Default to any-to-any if mode is not recognized
		lc.generatedLanes = generateAnyToAnyLanes(chains)
		return lc.generatedLanes
	}
}

// generateChainTierLanes generates lanes where chains of a 'high' tier are connected to all chains
// chains of a 'low' tier are only connected to chains of a 'high' tier.
func generateChainTierLanes(chains []uint64, highTierCount, lowtierCount int) []LaneConfig {
	uniqueLanes := mapset.NewSet[LaneConfig]()
	highTierSels, _ := getTierChainSelectors(chains, highTierCount, lowtierCount)
	for _, src := range highTierSels {
		for _, dst := range chains {
			if src != dst {
				// make lanes bidirectional
				uniqueLanes.Add(LaneConfig{
					SourceChain:      src,
					DestinationChain: dst,
				})
				uniqueLanes.Add(LaneConfig{
					SourceChain:      dst,
					DestinationChain: src,
				})
			}
		}
	}
	return uniqueLanes.ToSlice()
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

func generateBidirectionalRandomLanesWithMinConnectivity(chains []uint64, numLanes int) []LaneConfig {
	if len(chains) <= 1 {
		// If there's only one chain or none, no lanes can be generated
		return []LaneConfig{}
	}
	rng := rand.New(rand.NewSource(rand.Int63()))

	// Ensure minimum connectivity - each chain must be both source and destination
	var generatedLanes []LaneConfig

	// Shuffle chains for randomness in connectivity pattern
	shuffledChains := make([]uint64, len(chains))
	copy(shuffledChains, chains)
	rng.Shuffle(len(shuffledChains), func(i, j int) {
		shuffledChains[i], shuffledChains[j] = shuffledChains[j], shuffledChains[i]
	})

	// Create minimum connectivity: each chain as source and destination bidirectionally
	for i := range shuffledChains {
		// First cycle - connect to next chain
		src := shuffledChains[i]
		dst := shuffledChains[(i+1)%len(shuffledChains)]
		generatedLanes = append(generatedLanes, LaneConfig{
			SourceChain:      src,
			DestinationChain: dst,
		})
		// bidirectional connection
		generatedLanes = append(generatedLanes, LaneConfig{
			SourceChain:      dst,
			DestinationChain: src,
		})
	}

	// Fill remaining slots with random lanes
	if numLanes <= len(generatedLanes) {
		return generatedLanes
	}

	// Create set of used lanes to avoid duplicates
	usedLanes := make(map[LaneConfig]bool)
	for _, lane := range generatedLanes {
		usedLanes[lane] = true
	}

	// Generate additional random lanes
	allPossibleLanes := generateAnyToAnyLanes(chains)
	var availableLanes []LaneConfig

	// Filter out already used lanes
	for _, lane := range allPossibleLanes {
		if !usedLanes[lane] {
			availableLanes = append(availableLanes, lane)
		}
	}

	// Shuffle available lanes
	rng.Shuffle(len(availableLanes), func(i, j int) {
		availableLanes[i], availableLanes[j] = availableLanes[j], availableLanes[i]
	})

	for _, availableLane := range availableLanes {
		if len(generatedLanes) >= numLanes {
			break
		}
		// Add only if it doesn't already exist in guaranteed lanes
		if !usedLanes[availableLane] {
			// Add the available lane and its reverse to ensure bidirectionality
			reverseLane := LaneConfig{
				SourceChain:      availableLane.DestinationChain,
				DestinationChain: availableLane.SourceChain,
			}
			generatedLanes = append(generatedLanes, availableLane)
			generatedLanes = append(generatedLanes, reverseLane)
			usedLanes[availableLane] = true
			usedLanes[reverseLane] = true
		}
	}

	return generatedLanes
}

// calculateMinimumLanesNeeded calculates minimum lanes needed for connectivity where each chain
// must be both a source and destination.
func calculateMinimumLanesNeeded(numChains int) int {
	if numChains <= 1 {
		return 0
	}

	// Minimum is:
	// bidirectional lanes for each chain
	// each chain[i] <-> [chain[i+1]]
	minLanes := numChains * 2

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
	slices.Sort(connectedChains)

	return connectedChains
}

// DiscoverLanesFromDeployedState reverse engineers the lane configuration from deployed state
func (lc *LaneConfiguration) DiscoverLanesFromDeployedState(env cldf.Environment, state *stateview.CCIPOnChainState) error {
	var discoveredLanes []LaneConfig

	evmChains := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(selectors.FamilyEVM))
	solChains := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(selectors.FamilySolana))
	aptosChains := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(selectors.FamilyAptos))
	//nolint: gocritic // append is fine here
	allChains := append(evmChains, solChains...)
	allChains = append(allChains, aptosChains...)

	// Discover EVM to EVM lanes
	for _, srcChain := range evmChains {
		srcChainState, exists := state.Chains[srcChain]
		if !exists {
			continue
		}

		// Check which destination chains are configured on the OnRamp
		destinations, err := lc.getEnabledDestinationsFromOnRamp(srcChainState, srcChain, allChains)
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
		destinations, err := lc.getEnabledDestinationsFromSolanaRouter(env, srcChain, srcChainState, allChains)
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

	// Discover Aptos to EVM lanes
	for _, srcChain := range aptosChains {
		srcChainState, exists := state.AptosChains[srcChain]
		if !exists {
			continue
		}

		// Check which EVM destination chains are configured on the Aptos Router
		destinations, err := lc.getEnabledDestinationsFromAptosRouter(env, srcChain, srcChainState, allChains)
		if err != nil {
			return fmt.Errorf("failed to get enabled EVM destinations for Aptos chain %d: %w", srcChain, err)
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
func (lc *LaneConfiguration) getEnabledDestinationsFromOnRamp(
	chainState evm.CCIPChainState,
	srcSelector uint64,
	candidateDestinations []uint64) ([]uint64, error) {
	var enabledDestinations []uint64

	// For each candidate destination, check if it's enabled on the OnRamp
	for _, dstChain := range candidateDestinations {
		if dstChain == srcSelector {
			continue
		}
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
func (lc *LaneConfiguration) getEnabledDestinationsFromSolanaRouter(env cldf.Environment, selector uint64, chainState solState.CCIPChainState, candidateDestinations []uint64) ([]uint64, error) {
	var enabledDestinations []uint64

	// For each candidate destination, check if it's enabled on the Solana Router
	for _, dstChain := range candidateDestinations {
		if dstChain == selector {
			continue
		}
		// we don't verify against error because if the destination is not configured, it will return an error
		isEnabled, _ := lc.isDestinationEnabledOnSolanaRouter(env.GetContext(), chainState, dstChain, env.BlockChains.SolanaChains()[selector].Client)
		if isEnabled {
			enabledDestinations = append(enabledDestinations, dstChain)
		}
	}

	return enabledDestinations, nil
}

func (lc *LaneConfiguration) getEnabledDestinationsFromAptosRouter(env cldf.Environment, selector uint64, chainState aptosState.CCIPChainState, candidateDestinations []uint64) ([]uint64, error) {
	var enabledDestinations []uint64

	// For each candidate destination, check if it's enabled on the Aptos Router
	for _, dstChain := range candidateDestinations {
		if dstChain == selector {
			continue
		}
		// we don't verify against error because if the destination is not configured, it will return an error
		isEnabled, _ := lc.isDestinationEnabledOnAptosRouter(env, selector, chainState, dstChain)
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
func (lc *LaneConfiguration) isDestinationEnabledOnSolanaRouter(ctx context.Context, chainState solState.CCIPChainState, destinationChain uint64, client *solrpc.Client) (bool, error) {
	routerRemoteStatePDA, _ := ccipSolState.FindDestChainStatePDA(destinationChain, chainState.Router)
	var destChainStateAccount solRouter.DestChain
	err := solCommonUtil.GetAccountDataBorshInto(ctx, client, routerRemoteStatePDA, cldf_solana.SolDefaultCommitment, &destChainStateAccount)
	if err != nil {
		// If we can't get the config, assume it's not enabled
		return false, fmt.Errorf("failed to get dest chain state for %d: %w", destinationChain, err)
	}
	return true, nil
}

// isDestinationEnabledOnAptosRouter checks if a destination is enabled on the Aptos Router
func (lc *LaneConfiguration) isDestinationEnabledOnAptosRouter(env cldf.Environment, aptosChainSelector uint64, chainState aptosState.CCIPChainState, destinationChain uint64) (bool, error) {
	// Get the client from the environment
	client := env.BlockChains.AptosChains()[aptosChainSelector].Client

	// Bind to the OnRamp contract
	boundOnRamp := ccip_onramp.Bind(chainState.CCIPAddress, client)

	// Use IsChainSupported directly for the specific destination
	isSupported, err := boundOnRamp.Onramp().IsChainSupported(nil, destinationChain)
	if err != nil {
		// If we can't check support, assume it's not enabled
		return false, fmt.Errorf("failed to check if destination chain is supported on Aptos onRamp: %w", err)
	}

	return isSupported, nil
}

// GetSourceChainsForDestination returns all source chains that can send to a specific destination
func (lc *LaneConfiguration) GetSourceChainsForDestination(destination uint64) []uint64 {
	if lc == nil {
		panic("LaneConfiguration is nil, cannot get source chains for destination")
	}

	var sources []uint64
	for _, lane := range lc.generatedLanes {
		if lane.DestinationChain == destination {
			sources = append(sources, lane.SourceChain)
		}
	}

	// Sort for deterministic order
	slices.Sort(sources)

	return sources
}

// GetDestinationChainsForSource returns all destination chains that a source can send to
func (lc *LaneConfiguration) GetDestinationChainsForSource(source uint64) []uint64 {
	if lc == nil {
		panic("LaneConfiguration is nil, cannot get destination chains for source")
	}

	var destinations []uint64
	for _, lane := range lc.generatedLanes {
		if lane.SourceChain == source {
			destinations = append(destinations, lane.DestinationChain)
		}
	}

	// Sort for deterministic order
	slices.Sort(destinations)

	return destinations
}

// LaneStats provides statistics about the discovered lane configuration
type LaneStats struct {
	TotalLanes        int
	UniqueChains      int
	SourceChains      int
	DestinationChains int
}

// GetLaneStats For metrics and reporting on the lane configuration
func (lc *LaneConfiguration) GetLaneStats() LaneStats {
	if lc == nil {
		panic("LaneConfiguration is nil")
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

	return stats
}

func (lc *LaneConfiguration) LogLaneConfigInfo(lggr logger.Logger) {
	if lc == nil {
		lggr.Warn("LaneConfiguration is nil, cannot log stats")
		return
	}

	stats := lc.GetLaneStats()
	lggr.Infow("Lane Configuration Stats",
		"TotalLanes", stats.TotalLanes,
		"UniqueChains", stats.UniqueChains,
		"SourceChains", stats.SourceChains,
		"DestinationChains", stats.DestinationChains,
		"GeneratedLanes", lc.generatedLanes,
	)
}

// Example TOML configurations:

// Any-to-any (traditional full mesh)
/*
Mode = "any-to-any"
*/

// Random lanes
/*

Mode = "random-lanes"
NumLanes = 350
*/
