package ccip

import (
	"testing"

	"github.com/smartcontractkit/chainlink/deployment/environment/crib"

	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

// TestLaneDiscovery_AnyToAny tests lane discovery when all chains are connected to each other
func TestLaneDiscovery_AnyToAny(t *testing.T) {
	//TODO: ADD Solana lanes
	tenv, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfChains(3),
	)

	e := tenv.Env
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	// Add all lanes (any-to-any setup)
	testhelpers.AddLanesForAll(t, &tenv, state)

	state, err = stateview.LoadOnchainState(e)
	require.NoError(t, err)

	// Discover lanes from deployed state
	laneConfig := &crib.LaneConfiguration{}
	err = laneConfig.DiscoverLanesFromDeployedState(e, &state)
	require.NoError(t, err)

	// Verify discovered lanes
	discoveredLanes, err := laneConfig.GetLanes()
	require.NoError(t, err)
	chains := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))

	// Should have n*(n-1) lanes for n chains (any-to-any)
	expectedLaneCount := len(chains) * (len(chains) - 1)
	require.Equal(t, expectedLaneCount, len(discoveredLanes),
		"Should discover %d lanes for %d chains in any-to-any setup", expectedLaneCount, len(chains))

	// Verify all chains are connected
	connectedChains := laneConfig.GetConnectedChains()
	require.Equal(t, len(chains), len(connectedChains),
		"All chains should be connected")

	// Verify each chain can reach every other chain
	for _, src := range chains {
		destinations := laneConfig.GetDestinationChainsForSource(src)
		require.Equal(t, len(chains)-1, len(destinations),
			"Each chain should have %d destinations", len(chains)-1)

		sources := laneConfig.GetSourceChainsForDestination(src)
		require.Equal(t, len(chains)-1, len(sources),
			"Each chain should have %d sources", len(chains)-1)
	}

	// Verify statistics
	stats := laneConfig.GetLaneStats()
	require.Equal(t, expectedLaneCount, stats.TotalLanes)
	require.Equal(t, len(chains), stats.UniqueChains)
	require.Equal(t, len(chains), stats.SourceChains)
	require.Equal(t, len(chains), stats.DestinationChains)
}

// TODO: ADD Solana lanes
// TestLaneDiscovery_PartialConnectivity tests lane discovery with limited connectivity
func TestLaneDiscovery_PartialConnectivity(t *testing.T) {
	tenv, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfChains(4),
	)

	e := tenv.Env
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	chains := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
	require.Len(t, chains, 4, "Should have 4 chains")

	chainA, chainB, chainC, chainD := chains[0], chains[1], chains[2], chains[3]

	// Setup partial connectivity: A->B, A->C,  B->C, C->D, D->A (cycle)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &tenv, state, chainA, chainB, false)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &tenv, state, chainA, chainC, false)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &tenv, state, chainB, chainC, false)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &tenv, state, chainC, chainD, false)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &tenv, state, chainD, chainA, false)

	// Reload state after adding lanes
	state, err = stateview.LoadOnchainState(e)
	require.NoError(t, err)

	// Discover lanes from deployed state
	laneConfig := &crib.LaneConfiguration{}
	err = laneConfig.DiscoverLanesFromDeployedState(e, &state)
	require.NoError(t, err)

	// Verify discovered lanes
	discoveredLanes, err := laneConfig.GetLanes()
	require.NoError(t, err)
	require.Equal(t, 5, len(discoveredLanes), "Should discover exactly 5 lanes")

	// Verify specific lanes exist
	expectedLanes := []crib.LaneConfig{
		{SourceChain: chainA, DestinationChain: chainB},
		{SourceChain: chainA, DestinationChain: chainC},
		{SourceChain: chainB, DestinationChain: chainC},
		{SourceChain: chainC, DestinationChain: chainD},
		{SourceChain: chainD, DestinationChain: chainA},
	}

	for _, expectedLane := range expectedLanes {
		found := false
		for _, discoveredLane := range discoveredLanes {
			if discoveredLane.SourceChain == expectedLane.SourceChain &&
				discoveredLane.DestinationChain == expectedLane.DestinationChain {
				found = true
				break
			}
		}
		require.True(t, found, "Expected lane %d->%d not found",
			expectedLane.SourceChain, expectedLane.DestinationChain)
	}

	// Verify connectivity patterns
	require.Equal(t, []uint64{chainB, chainC}, laneConfig.GetDestinationChainsForSource(chainA))
	require.Equal(t, []uint64{chainC}, laneConfig.GetDestinationChainsForSource(chainB))
	require.Equal(t, []uint64{chainD}, laneConfig.GetDestinationChainsForSource(chainC))
	require.Equal(t, []uint64{chainA}, laneConfig.GetDestinationChainsForSource(chainD))

	require.Equal(t, []uint64{chainD}, laneConfig.GetSourceChainsForDestination(chainA))
	require.Equal(t, []uint64{chainA}, laneConfig.GetSourceChainsForDestination(chainB))
	require.Equal(t, []uint64{chainA, chainB}, laneConfig.GetSourceChainsForDestination(chainC))
	require.Equal(t, []uint64{chainC}, laneConfig.GetSourceChainsForDestination(chainD))
}

// TestLaneDiscovery_EmptyState tests lane discovery with no lanes configured
func TestLaneDiscovery_EmptyState(t *testing.T) {
	tenv, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithNumOfUsersPerChain(1),
	)

	e := tenv.Env
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	// Don't add any lanes - test with empty state

	// Discover lanes from deployed state
	laneConfig := &crib.LaneConfiguration{}
	err = laneConfig.DiscoverLanesFromDeployedState(e, &state)
	require.NoError(t, err)

	// Verify no lanes discovered
	discoveredLanes, err := laneConfig.GetLanes()
	require.Error(t, err, "lanes have not been generated yet")
	require.Empty(t, discoveredLanes, "Should discover no lanes in empty state")

	// Verify empty connectivity
	connectedChains := laneConfig.GetConnectedChains()
	require.Empty(t, connectedChains, "Should have no connected chains")

	// Verify empty statistics
	stats := laneConfig.GetLaneStats()
	require.Equal(t, 0, stats.TotalLanes)
	require.Equal(t, 0, stats.UniqueChains)
	require.Equal(t, 0, stats.SourceChains)
	require.Equal(t, 0, stats.DestinationChains)
}

// TestLaneDiscovery_NilConfiguration tests behavior with nil configuration
func TestLaneDiscovery_NilConfiguration(t *testing.T) {
	var laneConfig *crib.LaneConfiguration = nil

	// Test GetLanes with nil config
	lanes, err := laneConfig.GetLanes()
	require.Error(t, err, "lane configuration is nil")
	require.Nil(t, lanes, "GetLanes should return nil slice for nil config")

	// Test GetConnectedChains with nil config
	chains := laneConfig.GetConnectedChains()
	require.Nil(t, chains, "GetConnectedChains should return nil slice for nil config")

	// Test GetSourceChainsForDestination with nil config
	sources := laneConfig.GetSourceChainsForDestination(12345)
	require.Nil(t, sources, "GetSourceChainsForDestination should return nil slice for nil config")

	// Test GetDestinationChainsForSource with nil config
	destinations := laneConfig.GetDestinationChainsForSource(12345)
	require.Nil(t, destinations, "GetDestinationChainsForSource should return nil slice for nil config")

	// Test GetLaneStats with nil config
	stats := laneConfig.GetLaneStats()
	require.Equal(t, crib.LaneStats{}, stats, "GetLaneStats should return zero stats for nil config")
}
