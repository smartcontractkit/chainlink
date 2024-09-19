package ccipdeployment

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployCCIPContracts(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     2,
		Nodes:      4,
	})
	// Deploy all the CCIP contracts.
	homeChain := e.AllChainSelectors()[HomeChainIndex]
	addressBook, _, err := DeployCapReg(lggr, e.Chains[homeChain])
	require.NoError(t, err)
	s, err := LoadOnchainState(e, addressBook)
	require.NoError(t, err)

	feedChain := e.AllChainSelectors()[FeedChainIndex]
	feedAddresses, _, err := DeployFeeds(lggr, e.Chains[feedChain])
	require.NoError(t, err)

	// Merge the feed addresses into the address book.
	require.NoError(t, addressBook.Merge(feedAddresses))

	// Load the state after deploying the cap reg and feeds.
	homeAndFeedStates, err := LoadOnchainState(e, addressBook)
	require.NoError(t, err)
	require.NotNil(t, s.Chains[homeChain].CapabilityRegistry)
	require.NotNil(t, s.Chains[homeChain].CCIPConfig)
	require.NotNil(t, homeAndFeedStates.Chains[feedChain].USDFeeds)

	ab, err := DeployCCIPContracts(e, DeployCCIPContractConfig{
		HomeChainSel:     homeChain,
		FeedChainSel:     feedChain,
		ChainsToDeploy:   e.AllChainSelectors(),
		TokenConfig:      NewTokenConfig(),
		CCIPOnChainState: homeAndFeedStates,
	})
	require.NoError(t, err)
	state, err := LoadOnchainState(e, ab)
	require.NoError(t, err)
	snap, err := state.Snapshot(e.AllChainSelectors())
	require.NoError(t, err)

	// Assert expect every deployed address to be in the address book.
	// TODO (CCIP-3047): Add the rest of CCIPv2 representation
	b, err := json.MarshalIndent(snap, "", "	")
	require.NoError(t, err)
	fmt.Println(string(b))
}

func TestJobSpecGeneration(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
		Nodes:  1,
	})
	js, err := NewCCIPJobSpecs(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	for node, jb := range js {
		fmt.Println(node, jb)
	}
	// TODO: Add job assertions
}
