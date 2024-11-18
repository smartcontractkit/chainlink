package smoke

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/multicall3"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_CCIPBatching(t *testing.T) {
	// Setup 3 chains, with 2 lanes going to the dest.
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	// Will load 3 chains when specified by the overrides.toml or env vars (E2E_TEST_SELECTED_NETWORK).
	// See e2e-tests.yml.
	e, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr)

	state, err := ccdeploy.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Env.Chains)
	require.Len(t, allChainSelectors, 3, "this test expects 3 chains")
	sourceChain1 := allChainSelectors[0]
	sourceChain2 := allChainSelectors[1]
	destChain := allChainSelectors[2]
	t.Log("All chain selectors:", allChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector 1:", sourceChain1,
		", source chain selector 2:", sourceChain2,
		", dest chain selector:", destChain,
	)
	output, err := changeset.DeployPrerequisites(e.Env, changeset.DeployPrerequisiteConfig{
		ChainSelectors: e.Env.AllChainSelectors(),
	})
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(output.AddressBook))

	tokenConfig := ccdeploy.NewTestTokenConfig(state.Chains[e.FeedChainSel].USDFeeds)
	// Apply migration
	output, err = changeset.InitialDeploy(e.Env, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:   e.HomeChainSel,
		FeedChainSel:   e.FeedChainSel,
		ChainsToDeploy: allChainSelectors,
		TokenConfig:    tokenConfig,
		MCMSConfig:     ccdeploy.NewTestMCMSConfig(t, e.Env),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	require.NoError(t, e.Env.ExistingAddresses.Merge(output.AddressBook))
	// Get new state after migration.
	state, err = ccdeploy.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	ccdeploy.ReplayLogs(t, e.Env.Offchain, e.ReplayBlocks)

	// Apply the jobs.
	for nodeID, jobs := range output.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Env.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}

	// connect sourceChain1 and sourceChain2 to destChain
	require.NoError(t, ccdeploy.AddLaneWithDefaultPrices(e.Env, state, sourceChain1, destChain))
	require.NoError(t, ccdeploy.AddLaneWithDefaultPrices(e.Env, state, sourceChain2, destChain))

	// var (
	// 	replayed      bool
	// 	nonce         map[uint64]uint64
	// 	senderSource1 = common.LeftPadBytes(e.Env.Chains[sourceChain1].DeployerKey.From.Bytes(), 32)
	// 	senderSource2 = common.LeftPadBytes(e.Env.Chains[sourceChain2].DeployerKey.From.Bytes(), 32)
	// 	outSource1    messagingTestCaseOutput
	// 	outSource2    messagingTestCaseOutput
	// 	setupSource1  = testCaseSetup{
	// 		t:            t,
	// 		sender:       senderSource1,
	// 		deployedEnv:  e,
	// 		onchainState: state,
	// 		sourceChain:  sourceChain1,
	// 		destChain:    destChain,
	// 	}
	// 	setupSource2 = testCaseSetup{
	// 		t:            t,
	// 		sender:       senderSource2,
	// 		deployedEnv:  e,
	// 		onchainState: state,
	// 		sourceChain:  sourceChain2,
	// 		destChain:    destChain,
	// 	}
	// )

	t.Run("batch data only messages from multiple sources", func(t *testing.T) {
		// Generate some messages, on each source.
		// Send them from a multicall contract, i.e multiple ccip messages in the same tx.
		// assert they are committed in the same batch.
		msg := router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello world"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		}
		fee, err := state.Chains[sourceChain1].Router.GetFee(&bind.CallOpts{
			Context: ctx,
		}, destChain, msg)
		require.NoError(t, err)

		// Send the tx with the message through the multicall
		calldata := ccdeploy.CCIPSendCalldata(t, destChain, msg)
		tx, err := state.Chains[sourceChain1].Multicall3.Aggregate3Value(
			&bind.TransactOpts{
				From:   e.Env.Chains[sourceChain1].DeployerKey.From,
				Signer: e.Env.Chains[sourceChain1].DeployerKey.Signer,
				Value:  fee,
			},
			[]multicall3.Multicall3Call3Value{
				{
					Target:       state.Chains[sourceChain1].Router.Address(),
					AllowFailure: false,
					CallData:     calldata,
					Value:        fee,
				},
			},
		)
		require.NoError(t, err)
		_, err = e.Env.Chains[sourceChain1].Confirm(tx)
		require.NoError(t, err)

		// check that the message was emitted
		iter, err := state.Chains[sourceChain1].OnRamp.FilterCCIPMessageSent(
			nil, []uint64{destChain}, nil,
		)
		require.NoError(t, err)
		require.True(t, iter.Next())
		require.Equal(t, msg.Receiver, iter.Event.Message.Receiver)
	})

	t.Run("batch mix of data only messages and token messages from multiple sources", func(t *testing.T) {
		// Generate some messages, on each source.
		// Send them from a multicall contract, i.e multiple ccip messages in the same tx.
		// assert they are committed in the same batch.
	})
}
