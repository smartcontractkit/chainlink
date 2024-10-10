package ccipdeployment

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// TestAddLane covers the workflow of adding a lane
// between existing supported chains in CCIP.
func TestAddLane(t *testing.T) {
	// The test does not work if the environment has more than 2 chains and the contracts are deployed on only 2 chains.
	e := NewMemoryEnvironmentWithJobs(t, logger.TestLogger(t), 2)
	// Here we have CR + nodes set up, but no CCIP contracts deployed.
	state, err := LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)

	selectors := e.Env.AllChainSelectors()
	chain1, chain2 := selectors[0], selectors[1]

	feeds := state.Chains[e.FeedChainSel].USDFeeds
	tokenConfig := NewTokenConfig()
	tokenConfig.UpsertTokenInfo(LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: feeds[LinkSymbol].Address().String(),
			Decimals:          LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)
	feeTokenContracts := make(map[uint64]FeeTokenContracts)
	for _, sel := range []uint64{chain1, chain2} {
		feeTokenContracts[sel] = e.FeeTokenContracts[sel]
	}
	// Set up CCIP contracts and a DON per chain.
	err = DeployCCIPContracts(e.Env, e.Ab, DeployCCIPContractConfig{
		HomeChainSel:       e.HomeChainSel,
		FeedChainSel:       e.FeedChainSel,
		TokenConfig:        tokenConfig,
		MCMSConfig:         NewTestMCMSConfig(t, e.Env),
		FeeTokenContracts:  feeTokenContracts,
		ChainsToDeploy:     []uint64{chain1, chain2},
		CapabilityRegistry: state.Chains[e.HomeChainSel].CapabilityRegistry.Address(),
	})
	require.NoError(t, err)

	// We expect no lanes available on any chain.
	state, err = LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)
	for _, sel := range []uint64{chain1, chain2} {
		chain := state.Chains[sel]
		offRamps, err := chain.Router.GetOffRamps(nil)
		require.NoError(t, err)
		require.Len(t, offRamps, 0)
	}

	replayBlocks, err := LatestBlocksByChain(testcontext.Get(t), e.Env.Chains)
	require.NoError(t, err)
	// Add one lane and send traffic.

	require.NoError(t, AddLane(e.Env, state, chain1, chain2))

	for _, sel := range []uint64{chain1, chain2} {
		chain := state.Chains[sel]
		offRamps, err := chain.Router.GetOffRamps(nil)
		require.NoError(t, err)
		if sel == chain2 {
			require.Len(t, offRamps, 1)
			srcCfg, err := chain.OffRamp.GetSourceChainConfig(nil, chain1)
			require.NoError(t, err)
			require.Equal(t, common.LeftPadBytes(state.Chains[chain1].OnRamp.Address().Bytes(), 32), srcCfg.OnRamp)
		} else {
			require.Len(t, offRamps, 0)
		}
	}

	time.Sleep(30 * time.Second)
	ReplayLogs(t, e.Env.Offchain, replayBlocks)

	latesthdr, err := e.Env.Chains[chain2].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock := latesthdr.Number.Uint64()
	seqNum := SendRequest(t, e.Env, state, chain1, chain2, false)
	require.Equal(t, uint64(1), seqNum)
	require.NoError(t,
		ConfirmCommitWithExpectedSeqNumRange(
			t, e.Env.Chains[chain1], e.Env.Chains[chain2],
			state.Chains[chain2].OffRamp, &startBlock,
			cciptypes.SeqNumRange{
				cciptypes.SeqNum(seqNum),
				cciptypes.SeqNum(seqNum),
			}))
	require.NoError(t, ConfirmExecWithSeqNr(t, e.Env.Chains[chain1], e.Env.Chains[chain2], state.Chains[chain2].OffRamp, &startBlock, seqNum))

	// Add another lane
	replayBlocks, err = LatestBlocksByChain(testcontext.Get(t), e.Env.Chains)
	require.NoError(t, err)

	// disable onRamp for previous lane chain1 -> chain2
	tx, err := state.Chains[chain2].OffRamp.ApplySourceChainConfigUpdates(
		e.Env.Chains[chain2].DeployerKey, []offramp.OffRampSourceChainConfigArgs{
			{
				Router:              state.Chains[chain2].Router.Address(),
				SourceChainSelector: chain1,
				IsEnabled:           false,
				OnRamp:              common.LeftPadBytes(state.Chains[chain1].OnRamp.Address().Bytes(), 32),
			},
		})
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[chain2], tx, err)
	require.NoError(t, err)

	require.NoError(t, AddLane(e.Env, state, chain2, chain1))
	srcCfg, err := state.Chains[chain1].OffRamp.GetSourceChainConfig(nil, chain2)
	require.NoError(t, err)
	require.Equal(t, common.LeftPadBytes(state.Chains[chain2].OnRamp.Address().Bytes(), 32), srcCfg.OnRamp)
	require.Equal(t, true, srcCfg.IsEnabled)
	require.Equal(t, state.Chains[chain1].Router.Address(), srcCfg.Router)

	time.Sleep(30 * time.Second)
	ReplayLogs(t, e.Env.Offchain, replayBlocks)

	// Send traffic on the first lane and it should fail
	/*latesthdr, err = e.Env.Chains[chain2].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock = latesthdr.Number.Uint64()
	seqNum = SendRequest(t, e.Env, state, chain1, chain2, false)*/

	// Send traffic on the second lane and it should succeed
	latesthdr, err = e.Env.Chains[chain1].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock = latesthdr.Number.Uint64()
	seqNum = SendRequest(t, e.Env, state, chain2, chain1, false)
	require.Equal(t, uint64(1), seqNum)
	require.NoError(t,
		ConfirmCommitWithExpectedSeqNumRange(
			t, e.Env.Chains[chain2], e.Env.Chains[chain1],
			state.Chains[chain2].OffRamp, &startBlock,
			cciptypes.SeqNumRange{
				cciptypes.SeqNum(seqNum),
				cciptypes.SeqNum(seqNum),
			}))
	require.NoError(t, ConfirmExecWithSeqNr(t, e.Env.Chains[chain2], e.Env.Chains[chain1], state.Chains[chain1].OffRamp, &startBlock, seqNum))

	/*require.Equal(t, uint64(2), seqNum)
	require.Error(t, ConfirmExecWithSeqNr(t, e.Env.Chains[chain1], e.Env.Chains[chain2], state.Chains[chain2].OffRamp, &startBlock, seqNum))*/
}
