package ccipdeployment

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// TestAddLane covers the workflow of adding a lane
// between existing supported chains in CCIP.
func TestAddLane(t *testing.T) {
	// TODO: The offchain code doesn't yet support partial lane
	// enablement, need to address then re-enable this test.
	t.Skip()
	e := NewEnvironmentWithCRAndJobs(t, logger.TestLogger(t), 3)
	// Here we have CR + nodes set up, but no CCIP contracts deployed.
	state, err := LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	// Set up CCIP contracts and a DON per chain.
	ab, err := DeployCCIPContracts(e.Env, DeployCCIPContractConfig{
		HomeChainSel:     e.HomeChainSel,
		CCIPOnChainState: state,
	})
	require.NoError(t, err)
	require.NoError(t, e.Ab.Merge(ab))

	// We expect no lanes available on any chain.
	state, err = LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	for _, chain := range state.Chains {
		offRamps, err := chain.Router.GetOffRamps(nil)
		require.NoError(t, err)
		require.Len(t, offRamps, 0)
	}

	// Add one lane and send traffic.
	from, to := e.Env.AllChainSelectors()[0], e.Env.AllChainSelectors()[1]
	require.NoError(t, AddLane(e.Env, state, from, to))

	for sel, chain := range state.Chains {
		offRamps, err := chain.Router.GetOffRamps(nil)
		require.NoError(t, err)
		if sel == to {
			require.Len(t, offRamps, 1)
		} else {
			require.Len(t, offRamps, 0)
		}
	}
	seqNum := SendRequest(t, e.Env, state, from, to)
	require.Equal(t, uint64(1), seqNum)
	ConfirmExecution(t, e.Env.Chains[from], e.Env.Chains[to], state.Chains[to].OffRamp, seqNum)

	// TODO: Add a second lane, then disable the first and
	// ensure we can send on the second but not the first.
}
