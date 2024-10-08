package ccipdeployment

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// TestAddLane covers the workflow of adding a lane
// between existing supported chains in CCIP.
func TestAddLane(t *testing.T) {
	e := NewMemoryEnvironmentWithJobs(t, logger.TestLogger(t), 3)
	// Here we have CR + nodes set up, but no CCIP contracts deployed.
	state, err := LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	from, to := e.Env.AllChainSelectors()[0], e.Env.AllChainSelectors()[1]

	// Set up CCIP contracts and a DON per chain.
	err = DeployCCIPContracts(e.Env, e.Ab, DeployCCIPContractConfig{
		HomeChainSel:       e.HomeChainSel,
		FeedChainSel:       e.FeedChainSel,
		TokenConfig:        NewTokenConfig(),
		MCMSConfig:         NewTestMCMSConfig(t, e.Env),
		FeeTokenContracts:  e.FeeTokenContracts,
		ChainsToDeploy:     []uint64{from, to},
		CapabilityRegistry: state.Chains[e.HomeChainSel].CapabilityRegistry.Address(),
	})
	require.NoError(t, err)

	// We expect no lanes available on any chain.
	state, err = LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	for sel, chain := range state.Chains {
		if sel != from && sel != to {
			continue
		}
		offRamps, err := chain.Router.GetOffRamps(nil)
		require.NoError(t, err)
		require.Len(t, offRamps, 0)
	}

	replayBlocks, err := LatestBlocksByChain(testcontext.Get(t), e.Env.Chains)
	require.NoError(t, err)
	// Add one lane and send traffic.
	require.NoError(t, AddLane(e.Env, state, from, to))
	require.NoError(t, AddLane(e.Env, state, to, from))
	for sel, chain := range state.Chains {
		if sel != from && sel != to {
			continue
		}
		offRamps, err := chain.Router.GetOffRamps(nil)
		require.NoError(t, err)
		if sel == to {
			require.Len(t, offRamps, 1)
		} else {
			//			require.Len(t, offRamps, 0)
		}
	}
	time.Sleep(30 * time.Second)
	ReplayLogs(t, e.Env.Offchain, replayBlocks)

	latesthdr, err := e.Env.Chains[to].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock := latesthdr.Number.Uint64()
	seqNum := SendRequest(t, e.Env, state, from, to, false)
	require.Equal(t, uint64(1), seqNum)
	require.NoError(t,
		ConfirmCommitWithExpectedSeqNumRange(
			t, e.Env.Chains[from], e.Env.Chains[to],
			state.Chains[to].OffRamp, &startBlock,
			ccipocr3.SeqNumRange{
				ccipocr3.SeqNum(seqNum),
				ccipocr3.SeqNum(seqNum),
			}))
	require.NoError(t, ConfirmExecWithSeqNr(t, e.Env.Chains[from], e.Env.Chains[to], state.Chains[to].OffRamp, &startBlock, seqNum))

	// TODO: Add a second lane, then disable the first and
	// ensure we can send on the second but not the first.
}
