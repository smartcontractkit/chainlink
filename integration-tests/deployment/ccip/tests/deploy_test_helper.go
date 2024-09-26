package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func DeployCCIPContractsTest(t *testing.T, e deployment.Environment) {
	lggr := logger.TestLogger(t)
	// Deploy all the CCIP contracts.
	homeChain := e.AllChainSelectors()[HomeChainIndex]
	addressBook, _, err := ccipdeployment.DeployCapReg(lggr, e.Chains[homeChain])
	require.NoError(t, err)
	s, err := ccipdeployment.LoadOnchainState(e, addressBook)
	require.NoError(t, err)

	feedChain := e.AllChainSelectors()[FeedChainIndex]
	feedAddresses, _, err := ccipdeployment.DeployMockFeeds(lggr, e.Chains[feedChain])
	require.NoError(t, err)

	// Merge the feed addresses into the address book.
	require.NoError(t, addressBook.Merge(feedAddresses))

	// Load the state after deploying the cap reg and feeds.
	homeAndFeedStates, err := ccipdeployment.LoadOnchainState(e, addressBook)
	require.NoError(t, err)
	require.NotNil(t, s.Chains[homeChain].CapabilityRegistry)
	require.NotNil(t, s.Chains[homeChain].CCIPConfig)
	require.NotNil(t, homeAndFeedStates.Chains[feedChain].USDFeeds)

	ab, err := ccipdeployment.DeployCCIPContracts(e, ccipdeployment.DeployCCIPContractConfig{
		HomeChainSel:     homeChain,
		FeedChainSel:     feedChain,
		ChainsToDeploy:   e.AllChainSelectors(),
		TokenConfig:      ccipdeployment.NewTokenConfig(),
		CCIPOnChainState: homeAndFeedStates,
	})
	require.NoError(t, err)
	state, err := ccipdeployment.LoadOnchainState(e, ab)
	require.NoError(t, err)
	snap, err := state.View(e.AllChainSelectors())
	require.NoError(t, err)

	// Assert expect every deployed address to be in the address book.
	// TODO (CCIP-3047): Add the rest of CCIPv2 representation
	b, err := json.MarshalIndent(snap, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
}

func JobSpecGenerationTest(t *testing.T, e deployment.Environment) {
	js, err := ccipdeployment.NewCCIPJobSpecs(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	for node, jb := range js {
		fmt.Println(node, jb)
	}
	// TODO: Add job assertions
}
