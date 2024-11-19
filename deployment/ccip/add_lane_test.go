package ccipdeployment

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	commonutils "github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// TestAddLane covers the workflow of adding a lane between two chains and enabling it.
// It also covers the case where the onRamp is disabled on the OffRamp contract initially and then enabled.
func TestAddLane(t *testing.T) {
	t.Parallel()
	// We add more chains to the chainlink nodes than the number of chains where CCIP is deployed.
	e := NewMemoryEnvironmentWithJobs(t, logger.TestLogger(t), 4, 4)
	// Here we have CR + nodes set up, but no CCIP contracts deployed.
	state, err := LoadOnchainState(e.Env)
	require.NoError(t, err)

	selectors := e.Env.AllChainSelectors()
	// deploy CCIP contracts on two chains
	chain1, chain2 := selectors[0], selectors[1]

	feeds := state.Chains[e.FeedChainSel].USDFeeds
	tokenConfig := NewTestTokenConfig(feeds)
	newAddresses := deployment.NewMemoryAddressBook()
	err = DeployPrerequisiteChainContracts(e.Env, newAddresses, e.Env.AllChainSelectors())
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(newAddresses))

	// Set up CCIP contracts and a DON per chain.
	newAddresses = deployment.NewMemoryAddressBook()
	err = DeployCCIPContracts(e.Env, newAddresses, DeployCCIPContractConfig{
		HomeChainSel:   e.HomeChainSel,
		FeedChainSel:   e.FeedChainSel,
		TokenConfig:    tokenConfig,
		MCMSConfig:     NewTestMCMSConfig(t, e.Env),
		ChainsToDeploy: []uint64{chain1, chain2},
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(newAddresses))

	// We expect no lanes available on any chain.
	state, err = LoadOnchainState(e.Env)
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
	require.NoError(t, AddLaneWithDefaultPrices(e.Env, state, chain1, chain2))

	ReplayLogs(t, e.Env.Offchain, replayBlocks)
	time.Sleep(30 * time.Second)
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

	latesthdr, err := e.Env.Chains[chain2].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock := latesthdr.Number.Uint64()
	// Send traffic on the first lane and it should not be processed by the plugin as onRamp is disabled
	// we will check this by confirming that the message is not executed by the end of the test
	msgSentEvent1 := TestSendRequest(t, e.Env, state, chain1, chain2, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[chain2].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	require.Equal(t, uint64(1), msgSentEvent1.SequenceNumber)

	// Add another lane
	require.NoError(t, AddLaneWithDefaultPrices(e.Env, state, chain2, chain1))

	// Send traffic on the second lane and it should succeed
	latesthdr, err = e.Env.Chains[chain1].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	startBlock2 := latesthdr.Number.Uint64()
	msgSentEvent2 := TestSendRequest(t, e.Env, state, chain2, chain1, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[chain2].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	require.Equal(t, uint64(1), msgSentEvent2.SequenceNumber)
	require.NoError(t, commonutils.JustError(ConfirmExecWithSeqNr(t, e.Env.Chains[chain2], e.Env.Chains[chain1], state.Chains[chain1].OffRamp, &startBlock2, msgSentEvent2.SequenceNumber)))

	// now check for the previous message from chain 1 to chain 2 that it has not been executed till now as the onRamp was disabled
	ConfirmNoExecConsistentlyWithSeqNr(t, e.Env.Chains[chain1], e.Env.Chains[chain2], state.Chains[chain2].OffRamp, msgSentEvent1.SequenceNumber, 30*time.Second)

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
	require.NoError(t, commonutils.JustError(ConfirmExecWithSeqNr(t, e.Env.Chains[chain1], e.Env.Chains[chain2], state.Chains[chain2].OffRamp, &startBlock, msgSentEvent1.SequenceNumber)))
}
