package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
)

func AddLaneTest(t *testing.T, e DeployedEnv) {
	// Here we have CR + nodes set up, but no CCIP contracts deployed.
	state, err := ccipdeployment.LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	// Set up CCIP contracts and a DON per chain.
	ab, err := ccipdeployment.DeployCCIPContracts(e.Env, ccipdeployment.DeployCCIPContractConfig{
		HomeChainSel:     e.HomeChainSel,
		FeedChainSel:     e.FeedChainSel,
		TokenConfig:      e.TokenConfig,
		CCIPOnChainState: state,
	})
	require.NoError(t, err)
	require.NoError(t, e.Ab.Merge(ab))

	// We expect no lanes available on any chain.
	state, err = ccipdeployment.LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	for _, chain := range state.Chains {
		offRamps, err := chain.Router.GetOffRamps(nil)
		require.NoError(t, err)
		require.Len(t, offRamps, 0)
	}

	// Add one lane and send traffic.
	from, to := e.Env.AllChainSelectors()[0], e.Env.AllChainSelectors()[1]
	require.NoError(t, ccipdeployment.AddLane(e.Env, state, from, to))

	for sel, chain := range state.Chains {
		offRamps, err := chain.Router.GetOffRamps(nil)
		require.NoError(t, err)
		if sel == to {
			require.Len(t, offRamps, 1)
		} else {
			require.Len(t, offRamps, 0)
		}
	}
	latesthdr, err := e.Env.Chains[to].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock := latesthdr.Number.Uint64()
	seqNum := SendRequest(t, e.Env, state, from, to, false)
	require.Equal(t, uint64(1), seqNum)
	require.NoError(t, ConfirmExecWithSeqNr(t, e.Env.Chains[from], e.Env.Chains[to], state.Chains[to].OffRamp, &startBlock, seqNum))

	// TODO: Add a second lane, then disable the first and
	// ensure we can send on the second but not the first.
}
