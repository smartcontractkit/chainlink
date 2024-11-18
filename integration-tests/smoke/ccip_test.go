package smoke

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestInitialDeployOnLocal(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr)
	e := tenv.Env

	state, err := ccdeploy.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	output, err := changeset.DeployPrerequisites(tenv.Env, changeset.DeployPrerequisiteConfig{
		ChainSelectors: tenv.Env.AllChainSelectors(),
	})
	require.NoError(t, err)
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(output.AddressBook))

	// Apply migration
	output, err = changeset.InitialDeploy(tenv.Env, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		ChainsToDeploy: tenv.Env.AllChainSelectors(),
		TokenConfig:    ccdeploy.NewTestTokenConfig(feeds),
		MCMSConfig:     ccdeploy.NewTestMCMSConfig(t, e),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(output.AddressBook))
	// Get new state after migration.
	state, err = ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	ccdeploy.ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)

	// Apply the jobs.
	for nodeID, jobs := range output.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}

	// Add all lanes
	require.NoError(t, ccdeploy.AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[uint64]uint64)
	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block
			msgSentEvent := ccdeploy.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
				Data:         []byte("hello world"),
				TokenAmounts: nil,
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			})
			expectedSeqNum[dest] = msgSentEvent.SequenceNumber
		}
	}

	// Wait for all commit reports to land.
	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.
	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, ccdeploy.MockLinkPrice, timestampedPrice.Value)
	}

	// Wait for all exec reports to land
	ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

	// TODO: Apply the proposal.
}

func TestTokenTransfer(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr)

	e := tenv.Env
	state, err := ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	output, err := changeset.DeployPrerequisites(e, changeset.DeployPrerequisiteConfig{
		ChainSelectors: e.AllChainSelectors(),
	})
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook))

	// Apply migration
	output, err = changeset.InitialDeploy(e, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		ChainsToDeploy: e.AllChainSelectors(),
		TokenConfig:    ccdeploy.NewTestTokenConfig(state.Chains[tenv.FeedChainSel].USDFeeds),
		MCMSConfig:     ccdeploy.NewTestMCMSConfig(t, e),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook))
	// Get new state after migration and mock USDC token deployment.
	state, err = ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	srcToken, _, dstToken, _, err := ccdeploy.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		tenv.HomeChainSel,
		tenv.FeedChainSel,
		state,
		e.ExistingAddresses,
		"MY_TOKEN",
	)
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	ccdeploy.ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)

	// Apply the jobs.
	for nodeID, jobs := range output.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}

	// Add all lanes
	require.NoError(t, ccdeploy.AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[uint64]uint64)

	twoCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2))
	tx, err := srcToken.Mint(
		e.Chains[tenv.HomeChainSel].DeployerKey,
		e.Chains[tenv.HomeChainSel].DeployerKey.From,
		new(big.Int).Mul(twoCoins, big.NewInt(10)),
	)
	require.NoError(t, err)
	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err)

	tx, err = dstToken.Mint(
		e.Chains[tenv.FeedChainSel].DeployerKey,
		e.Chains[tenv.FeedChainSel].DeployerKey.From,
		new(big.Int).Mul(twoCoins, big.NewInt(10)),
	)
	require.NoError(t, err)
	_, err = e.Chains[tenv.FeedChainSel].Confirm(tx)
	require.NoError(t, err)

	tx, err = srcToken.Approve(e.Chains[tenv.HomeChainSel].DeployerKey, state.Chains[tenv.HomeChainSel].Router.Address(), twoCoins)
	require.NoError(t, err)
	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err)
	tx, err = dstToken.Approve(e.Chains[tenv.FeedChainSel].DeployerKey, state.Chains[tenv.FeedChainSel].Router.Address(), twoCoins)
	require.NoError(t, err)
	_, err = e.Chains[tenv.FeedChainSel].Confirm(tx)
	require.NoError(t, err)

	tokens := map[uint64][]router.ClientEVMTokenAmount{
		tenv.HomeChainSel: {{
			Token:  srcToken.Address(),
			Amount: twoCoins,
		}},
		tenv.FeedChainSel: {{
			Token:  dstToken.Address(),
			Amount: twoCoins,
		}},
	}

	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block

			var (
				receiver = common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32)
				data     = []byte("hello world")
				feeToken = common.HexToAddress("0x0")
			)
			if src == tenv.HomeChainSel && dest == tenv.FeedChainSel {
				msgSentEvent := ccdeploy.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
					Receiver:     receiver,
					Data:         data,
					TokenAmounts: tokens[src],
					FeeToken:     feeToken,
					ExtraArgs:    nil,
				})
				expectedSeqNum[dest] = msgSentEvent.SequenceNumber
			} else {
				msgSentEvent := ccdeploy.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
					Receiver:     receiver,
					Data:         data,
					TokenAmounts: nil,
					FeeToken:     feeToken,
					ExtraArgs:    nil,
				})
				expectedSeqNum[dest] = msgSentEvent.SequenceNumber
			}
		}
	}

	// Wait for all commit reports to land.
	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.
	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, ccdeploy.MockLinkPrice, timestampedPrice.Value)
	}

	// Wait for all exec reports to land
	ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

	balance, err := dstToken.BalanceOf(nil, state.Chains[tenv.FeedChainSel].Receiver.Address())
	require.NoError(t, err)
	require.Equal(t, twoCoins, balance)
}
