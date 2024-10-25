package ccipdeployment

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// TestAddLane covers the workflow of adding a lane between two chains and enabling it.
// It also covers the case where the onRamp is disabled on the OffRamp contract initially and then enabled.
func TestAddLane(t *testing.T) {
	// We add more chains to the chainlink nodes than the number of chains where CCIP is deployed.
	e := NewMemoryEnvironmentWithJobs(t, logger.TestLogger(t), 4, 4)
	// Here we have CR + nodes set up, but no CCIP contracts deployed.
	state, err := LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)

	selectors := e.Env.AllChainSelectors()
	// deploy CCIP contracts on two chains
	chain1, chain2 := selectors[0], selectors[1]

	feeds := state.Chains[e.FeedChainSel].USDFeeds
	tokenConfig := NewTestTokenConfig(feeds)

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
		OCRSecrets:         deployment.XXXGenerateTestOCRSecrets(),
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

	// Add one lane from chain1 to chain 2 and send traffic.
	require.NoError(t, AddLane(e.Env, state, chain1, chain2))

	// disable the onRamp initially on OffRamp
	disableRampTx, err := state.Chains[chain2].OffRamp.ApplySourceChainConfigUpdates(e.Env.Chains[chain2].DeployerKey, []offramp.OffRampSourceChainConfigArgs{
		{
			Router:              state.Chains[chain2].Router.Address(),
			SourceChainSelector: chain1,
			IsEnabled:           false,
			OnRamp:              common.LeftPadBytes(state.Chains[chain1].OnRamp.Address().Bytes(), 32),
		},
	})
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[chain2], disableRampTx, err)
	require.NoError(t, err)

	for _, sel := range []uint64{chain1, chain2} {
		chain := state.Chains[sel]
		offRamps, err := chain.Router.GetOffRamps(nil)
		require.NoError(t, err)
		if sel == chain2 {
			require.Len(t, offRamps, 1)
			srcCfg, err := chain.OffRamp.GetSourceChainConfig(nil, chain1)
			require.NoError(t, err)
			require.Equal(t, common.LeftPadBytes(state.Chains[chain1].OnRamp.Address().Bytes(), 32), srcCfg.OnRamp)
			require.False(t, srcCfg.IsEnabled)
		} else {
			require.Len(t, offRamps, 0)
		}
	}

	time.Sleep(30 * time.Second)
	ReplayLogs(t, e.Env.Offchain, replayBlocks)

	latesthdr, err := e.Env.Chains[chain2].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock := latesthdr.Number.Uint64()

	replayBlocks, err = LatestBlocksByChain(testcontext.Get(t), e.Env.Chains)
	require.NoError(t, err)

	time.Sleep(30 * time.Second)
	ReplayLogs(t, e.Env.Offchain, replayBlocks)
	// Send traffic on the first lane and it should not be processed by the plugin as onRamp is disabled
	// we will check this by confirming that the message is not executed by the end of the test
	seqNum := TestSendRequest(t, e.Env, state, chain1, chain2, false)
	require.Equal(t, uint64(1), seqNum)

	// Add another lane
	require.NoError(t, AddLane(e.Env, state, chain2, chain1))

	// Send traffic on the second lane and it should succeed
	latesthdr, err = e.Env.Chains[chain1].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock2 := latesthdr.Number.Uint64()
	seqNum2 := TestSendRequest(t, e.Env, state, chain2, chain1, false)
	require.Equal(t, uint64(1), seqNum)
	require.NoError(t,
		ConfirmCommitWithExpectedSeqNumRange(
			t, e.Env.Chains[chain2], e.Env.Chains[chain1],
			state.Chains[chain1].OffRamp, &startBlock,
			cciptypes.SeqNumRange{
				cciptypes.SeqNum(seqNum),
				cciptypes.SeqNum(seqNum),
			}))
	require.NoError(t, ConfirmExecWithSeqNr(t, e.Env.Chains[chain2], e.Env.Chains[chain1], state.Chains[chain1].OffRamp, &startBlock2, seqNum2))

	// now check for the previous message from chain 1 to chain 2 that it has not been executed till now as the onRamp was disabled
	ConfirmNoExecConsistentlyWithSeqNr(t, e.Env.Chains[chain1], e.Env.Chains[chain2], state.Chains[chain2].OffRamp, seqNum2, 1*time.Minute)

	// enable the onRamp on OffRamp
	enableRampTx, err := state.Chains[chain2].OffRamp.ApplySourceChainConfigUpdates(e.Env.Chains[chain2].DeployerKey, []offramp.OffRampSourceChainConfigArgs{
		{
			Router:              state.Chains[chain2].Router.Address(),
			SourceChainSelector: chain1,
			IsEnabled:           true,
			OnRamp:              common.LeftPadBytes(state.Chains[chain1].OnRamp.Address().Bytes(), 32),
		},
	})
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[chain2], enableRampTx, err)
	require.NoError(t, err)

	srcCfg, err := state.Chains[chain2].OffRamp.GetSourceChainConfig(nil, chain1)
	require.NoError(t, err)
	require.Equal(t, common.LeftPadBytes(state.Chains[chain1].OnRamp.Address().Bytes(), 32), srcCfg.OnRamp)
	require.True(t, srcCfg.IsEnabled)

	// we need the replay here otherwise plugin is not able to locate the message
	ReplayLogs(t, e.Env.Offchain, replayBlocks)
	time.Sleep(30 * time.Second)
	// Now that the onRamp is enabled, the request should be processed
	require.NoError(t, ConfirmExecWithSeqNr(t, e.Env.Chains[chain1], e.Env.Chains[chain2], state.Chains[chain2].OffRamp, &startBlock, seqNum))
}
