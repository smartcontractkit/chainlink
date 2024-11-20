package smoke

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
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

// setupTokens deploys transferable tokens on the source and dest, mints tokens for the source and dest, and
// approves the router to spend the tokens
func setupTokens(t *testing.T, tenv ccdeploy.DeployedEnv) (srcToken *burn_mint_erc677.BurnMintERC677, dstToken *burn_mint_erc677.BurnMintERC677) {
	lggr := logger.TestLogger(t)

	e := tenv.Env
	state, err := ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	tenCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(10))

	// Deploy the token to test transferring
	srcToken, _, dstToken, _, err = ccdeploy.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		tenv.HomeChainSel,
		tenv.FeedChainSel,
		state,
		e.ExistingAddresses,
		"MY_TOKEN",
	)

	linkToken := state.Chains[tenv.HomeChainSel].LinkToken

	require.NoError(t, err)

	tx, err := srcToken.Mint(
		e.Chains[tenv.HomeChainSel].DeployerKey,
		e.Chains[tenv.HomeChainSel].DeployerKey.From,
		new(big.Int).Mul(tenCoins, big.NewInt(10)),
	)
	require.NoError(t, err)

	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err)

	// Mint a destination token
	tx, err = dstToken.Mint(
		e.Chains[tenv.FeedChainSel].DeployerKey,
		e.Chains[tenv.FeedChainSel].DeployerKey.From,
		new(big.Int).Mul(tenCoins, big.NewInt(10)),
	)

	// Confirm the mint tx
	require.NoError(t, err)
	_, err = e.Chains[tenv.FeedChainSel].Confirm(tx)
	require.NoError(t, err)

	maxUint256 := math.MaxBig256

	// Approve the router to spend the tokens and confirm the tx's
	// To prevent having to approve the router for every transfer, we approve a sufficiently large amount
	tx, err = srcToken.Approve(e.Chains[tenv.HomeChainSel].DeployerKey, state.Chains[tenv.HomeChainSel].Router.Address(), maxUint256)
	require.NoError(t, err)
	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err)

	tx, err = dstToken.Approve(e.Chains[tenv.FeedChainSel].DeployerKey, state.Chains[tenv.FeedChainSel].Router.Address(), maxUint256)
	require.NoError(t, err)
	_, err = e.Chains[tenv.FeedChainSel].Confirm(tx)
	require.NoError(t, err)

	// Grant mint and burn roles to the deployer key for the newly deployed linkToken
	// Since those roles are not granted automatically
	tx, err = linkToken.GrantMintAndBurnRoles(e.Chains[tenv.HomeChainSel].DeployerKey, e.Chains[tenv.HomeChainSel].DeployerKey.From)
	require.NoError(t, err)
	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err)

	// Mint link token and confirm the tx
	tx, err = linkToken.Mint(
		e.Chains[tenv.HomeChainSel].DeployerKey,
		e.Chains[tenv.HomeChainSel].DeployerKey.From,
		new(big.Int).Mul(tenCoins, big.NewInt(10)),
	)
	require.NoError(t, err)

	_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
	require.NoError(t, err)

	return srcToken, dstToken
}

func Test_PricingForTokenTransfers(t *testing.T) {
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

	srcToken, dstToken := setupTokens(t, tenv)

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

	// Mint 2 tokens to be transferred
	twoCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2))

	linkToken := state.Chains[tenv.HomeChainSel].LinkToken
	maxUint256 := math.MaxBig256

	t.Run("Send Token Pay with Link token home chain -> remote", func(t *testing.T) {
		src := tenv.HomeChainSel
		dest, destChain := tenv.FeedChainSel, e.Chains[tenv.FeedChainSel]

		// Approve to spend link token
		tx, err := linkToken.Approve(e.Chains[src].DeployerKey, state.Chains[src].Router.Address(), maxUint256)
		require.NoError(t, err)
		_, err = e.Chains[src].Confirm(tx)
		require.NoError(t, err)

		// Get the fee Token Balance Before
		srctokenBalance, err := srcToken.BalanceOf(nil, e.Chains[src].DeployerKey.From)
		require.NoError(t, err)
		require.GreaterOrEqual(t, srctokenBalance.Int64(), twoCoins.Int64())

		// Create two ClientEVMTokenAmount structs to be passed to the router
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

		// Assign Src to the Home Chain Selector and destChain to the Feed Chain Selector

		// get the header for the destination chain and the relevant block number
		latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[dest] = &block

		// Send to the receiver on the destination chain paying with LINK token
		var (
			receiver = common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32)
			data     = []byte("hello world")
			feeToken = linkToken.Address()
		)

		ccipMessage := router.ClientEVM2AnyMessage{
			Receiver:     receiver,
			Data:         data,
			TokenAmounts: tokens[src],
			FeeToken:     feeToken,
			ExtraArgs:    nil,
		}

		// Get the fee Token Balance Before
		feeTokenBalanceBefore, err := linkToken.BalanceOf(nil, e.Chains[src].DeployerKey.From)
		require.NoError(t, err)

		// Check the fee Amount
		srcFee, err := state.Chains[src].Router.GetFee(nil, dest, ccipMessage)
		require.NoError(t, err)

		seqNum := ccdeploy.TestSendRequest(t, e, state, src, dest, false, ccipMessage).SequenceNumber
		expectedSeqNum[dest] = seqNum

		// Check the fee token balance after the request and ensure fee tokens were spent
		feeTokenBalanceAfter, err := linkToken.BalanceOf(nil, e.Chains[tenv.HomeChainSel].DeployerKey.From)
		require.NoError(t, err)
		require.Equal(t, feeTokenBalanceAfter, new(big.Int).Sub(feeTokenBalanceBefore, srcFee))

		// Wait for all commit reports to land.
		ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

		// After commit is reported on all chains, token prices should be updated in FeeQuoter.
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, ccdeploy.MockLinkPrice, timestampedPrice.Value)

		// Wait for all exec reports to land
		ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

		balance, err := dstToken.BalanceOf(nil, state.Chains[tenv.FeedChainSel].Receiver.Address())
		require.NoError(t, err)
		require.Equal(t, twoCoins, balance)

		// Delete from the expcted seq num the chain that was just tested so that we don't pass it to the
		// commit report in the next test
		delete(expectedSeqNum, dest)
	})

	t.Run("Send Token Pay with native remote chain -> home", func(t *testing.T) {
		// Create two ClientEVMTokenAmount structs to be passed to the router
		tokens := map[uint64][]router.ClientEVMTokenAmount{
			tenv.HomeChainSel: {{
				Token:  dstToken.Address(),
				Amount: twoCoins,
			}},
			tenv.FeedChainSel: {{
				Token:  srcToken.Address(),
				Amount: twoCoins,
			}},
		}

		// Assign Src to the Home Chain Selector and destChain to the Feed Chain Selector
		src := tenv.FeedChainSel
		dest, destChain := tenv.HomeChainSel, e.Chains[tenv.HomeChainSel]

		// get the header for the destination chain and the relevant block number
		latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[dest] = &block

		// Send to the receiver on the destination chain paying with native token
		var (
			receiver = common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32)
			data     = []byte("")
			feeToken = common.HexToAddress("0x0")
		)

		ccipMessage := router.ClientEVM2AnyMessage{
			Receiver:     receiver,
			Data:         data,
			TokenAmounts: tokens[dest],
			FeeToken:     feeToken,
			ExtraArgs:    nil,
		}

		seqNum := ccdeploy.TestSendRequest(t, e, state, src, dest, false, ccipMessage).SequenceNumber
		expectedSeqNum[dest] = seqNum

		// Wait for all commit reports to land.
		ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

		// After commit is reported on all chains, token prices should be updated in FeeQuoter.
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, ccdeploy.MockLinkPrice, timestampedPrice.Value)

		// Wait for all exec reports to land
		ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

		balance, err := srcToken.BalanceOf(nil, state.Chains[tenv.HomeChainSel].Receiver.Address())
		require.NoError(t, err)
		require.Equal(t, twoCoins, balance)

		delete(expectedSeqNum, dest)
	})

	t.Run("Send Token pay with wrapped native home chain -> remote", func(t *testing.T) {
		// Assign Src to the Home Chain Selector and destChain to the Feed Chain Selector
		src := tenv.HomeChainSel
		dest, destChain := tenv.FeedChainSel, e.Chains[tenv.FeedChainSel]

		WETH := state.Chains[src].Weth9

		// We need to acquire some WETH to send with the tx so we deposit some ETH into the WETH contract
		// but deployerKey is a reference so we need to dereference it and then reassign it
		depositOps := *e.Chains[tenv.HomeChainSel].DeployerKey
		depositOps.Value = new(big.Int).Mul(twoCoins, big.NewInt(2))
		tx, err := WETH.Deposit(&depositOps)
		require.NoError(t, err)
		_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
		require.NoError(t, err)

		// Approve to spend WETH token as feeToken
		fmt.Printf("--- APPROVING ROUTER TO SPEND---")
		tx, err = WETH.Approve(e.Chains[tenv.HomeChainSel].DeployerKey, state.Chains[tenv.HomeChainSel].Router.Address(), twoCoins)
		require.NoError(t, err)
		fmt.Printf("tx: %v\n", tx)
		_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
		require.NoError(t, err)
		fmt.Printf("--- APPROVED ---")

		allowance, err := WETH.Allowance(nil, e.Chains[src].DeployerKey.From, state.Chains[src].Router.Address())
		require.NoError(t, err)
		require.GreaterOrEqual(t, allowance.Int64(), twoCoins.Int64())

		// Create two ClientEVMTokenAmount structs to be passed to the router
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

		srcTokenBal, err := srcToken.BalanceOf(nil, e.Chains[src].DeployerKey.From)
		require.NoError(t, err)
		require.GreaterOrEqual(t, srcTokenBal.Int64(), twoCoins.Int64())

		// get the header for the destination chain and the relevant block number
		latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[dest] = &block

		// Send to the receiver on the destination chain paying with LINK token
		var (
			receiver = common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32)
			data     = []byte("")
			feeToken = state.Chains[src].Weth9.Address()
		)

		ccipMessage := router.ClientEVM2AnyMessage{
			Receiver:     receiver,
			Data:         data,
			TokenAmounts: tokens[src],
			FeeToken:     feeToken,
			ExtraArgs:    nil,
		}

		// Get the fee Token Balance Before
		feeTokenBalanceBefore, err := WETH.BalanceOf(nil, e.Chains[src].DeployerKey.From)
		require.NoError(t, err)

		// Check the fee Amount
		srcFee, err := state.Chains[src].Router.GetFee(nil, dest, ccipMessage)
		require.NoError(t, err)

		// Check that we have enough WETH to pay the fee
		require.GreaterOrEqual(t, feeTokenBalanceBefore.Int64(), srcFee.Int64())
		fmt.Printf("srcFee: %v\n", srcFee)
		fmt.Printf("feeTokenBalanceBefore: %v\n", feeTokenBalanceBefore)

		seqNum := ccdeploy.TestSendRequest(t, e, state, src, dest, false, ccipMessage).SequenceNumber
		expectedSeqNum[dest] = seqNum

		// Check the fee token balance after the request and ensure fee tokens were spent
		feeTokenBalanceAfter, err := WETH.BalanceOf(nil, e.Chains[tenv.HomeChainSel].DeployerKey.From)
		require.NoError(t, err)
		require.Equal(t, feeTokenBalanceAfter, new(big.Int).Sub(feeTokenBalanceBefore, srcFee))

		// Wait for all commit reports to land.
		ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

		// After commit is reported on all chains, token prices should be updated in FeeQuoter.
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, ccdeploy.MockLinkPrice, timestampedPrice.Value)

		// Wait for all exec reports to land
		ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

		balance, err := dstToken.BalanceOf(nil, state.Chains[dest].Receiver.Address())
		require.NoError(t, err)

		// The balance should be 4 since we've already sent 2 tokens over this lane direction in the
		// first part of the test, so balance should already be two
		require.Equal(t, new(big.Int).Mul(twoCoins, big.NewInt(2)), balance)
	})

	t.Run("Send Token but revert not enough tokens", func(t *testing.T) {
		// Approve the router to spend the tokens and confirm the tx's
		tx, err := srcToken.Approve(e.Chains[tenv.HomeChainSel].DeployerKey, state.Chains[tenv.HomeChainSel].Router.Address(), twoCoins)
		require.NoError(t, err)
		_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
		require.NoError(t, err)

		tx, err = dstToken.Approve(e.Chains[tenv.FeedChainSel].DeployerKey, state.Chains[tenv.FeedChainSel].Router.Address(), twoCoins)
		require.NoError(t, err)
		_, err = e.Chains[tenv.FeedChainSel].Confirm(tx)
		require.NoError(t, err)

		existingBalanceSrc, err := dstToken.BalanceOf(nil, e.Chains[tenv.HomeChainSel].DeployerKey.From)
		require.NoError(t, err)

		existingBalanceDst, err := dstToken.BalanceOf(nil, e.Chains[tenv.FeedChainSel].DeployerKey.From)
		require.NoError(t, err)

		// Assign Src to the Home Chain Selector and destChain to the Feed Chain Selector
		src := tenv.FeedChainSel
		dest, destChain := tenv.HomeChainSel, e.Chains[tenv.HomeChainSel]

		// Create two ClientEVMTokenAmount structs to be passed to the router
		tokens := map[uint64][]router.ClientEVMTokenAmount{
			tenv.HomeChainSel: {{
				Token: dstToken.Address(),
				// Send twice as many tokens as available on purpose
				Amount: new(big.Int).Mul(existingBalanceSrc, big.NewInt(2)),
			}},
			tenv.FeedChainSel: {{
				Token:  srcToken.Address(),
				Amount: new(big.Int).Mul(existingBalanceDst, big.NewInt(10)),
			}},
		}

		// get the header for the destination chain and the relevant block number
		latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[dest] = &block

		// Send to the receiver on the destination chain paying with native token
		var (
			receiver = common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32)
			data     = []byte("")
			feeToken = common.HexToAddress("0x0")
		)

		ccipMessage := router.ClientEVM2AnyMessage{
			Receiver:     receiver,
			Data:         data,
			TokenAmounts: tokens[dest],
			FeeToken:     feeToken,
			ExtraArgs:    nil,
		}

		// Send the CCP Request
		tx, _, err = ccdeploy.CCIPSendRequest(
			e,
			state,
			src, dest,
			false,
			ccipMessage,
		)

		// Check if the transaction reverted
		require.Error(t, err)
		require.Nil(t, tx)
	})
}

func Test_PricingForMessages(t *testing.T) {
	// Deploy the environment
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr)

	// Load the test state
	e := tenv.Env
	state, err := ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	tokenConfig := ccdeploy.NewTokenConfig()
	tokenConfig.UpsertTokenInfo(ccdeploy.LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: cciptypes.UnknownEncodedAddress(feeds[ccdeploy.LinkSymbol].Address().String()),
			Decimals:          ccdeploy.LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)

	// Apply migration
	output, err := changeset.InitialDeploy(e, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		ChainsToDeploy: e.AllChainSelectors(),
		TokenConfig:    tokenConfig,
		MCMSConfig:     ccdeploy.NewTestMCMSConfig(t, e),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook))
	// Get new state after migration and mock USDC token deployment.
	state, err = ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	ccdeploy.ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)

	// Create an empty tokenAmounts map since we are sending a message without a token
	tokens := []router.ClientEVMTokenAmount{}

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

	// Mint 2 tokens to be transferred
	twoCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2))

	linkToken := state.Chains[tenv.HomeChainSel].LinkToken

	t.Run("Send Token Pay with Link token home chain -> remote", func(t *testing.T) {
		// Approve to spend link token
		tx, err := linkToken.Approve(e.Chains[tenv.HomeChainSel].DeployerKey, state.Chains[tenv.HomeChainSel].Router.Address(), twoCoins)
		require.NoError(t, err)
		_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
		require.NoError(t, err)

		// Assign Src to the Home Chain Selector and destChain to the Feed Chain Selector
		src := tenv.HomeChainSel
		dest, destChain := tenv.FeedChainSel, e.Chains[tenv.FeedChainSel]

		// get the header for the destination chain and the relevant block number
		latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[dest] = &block

		// Send to the receiver on the destination chain paying with LINK token
		var (
			receiver = common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32)
			data     = []byte("hello world")
			feeToken = linkToken.Address()
		)

		ccipMessage := router.ClientEVM2AnyMessage{
			Receiver:     receiver,
			Data:         data,
			TokenAmounts: tokens,
			FeeToken:     feeToken,
			ExtraArgs:    nil,
		}

		// Get the fee Token Balance Before
		feeTokenBalanceBefore, err := linkToken.BalanceOf(nil, e.Chains[tenv.HomeChainSel].DeployerKey.From)
		require.NoError(t, err)

		// Check the fee Amount
		srcFee, err := state.Chains[src].Router.GetFee(nil, dest, ccipMessage)
		require.NoError(t, err)

		seqNum := ccdeploy.TestSendRequest(t, e, state, src, dest, false, ccipMessage).SequenceNumber
		expectedSeqNum[dest] = seqNum

		// Check the fee token balance after the request and ensure fee tokens were spent
		feeTokenBalanceAfter, err := linkToken.BalanceOf(nil, e.Chains[tenv.HomeChainSel].DeployerKey.From)
		require.NoError(t, err)
		require.Equal(t, feeTokenBalanceAfter, new(big.Int).Sub(feeTokenBalanceBefore, srcFee))

		// Wait for all commit reports to land.
		ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

		// After commit is reported on all chains, token prices should be updated in FeeQuoter.
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, ccdeploy.MockLinkPrice, timestampedPrice.Value)

		// Wait for all exec reports to land
		ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

		// Delete from the expcted seq num the chain that was just tested so that we don't pass it to the
		// commit report in the next test
		delete(expectedSeqNum, dest)
	})

	t.Run("Send Token Pay with native remote chain -> home", func(t *testing.T) {
		// Assign Src to the Home Chain Selector and destChain to the Feed Chain Selector
		src := tenv.FeedChainSel
		dest, destChain := tenv.HomeChainSel, e.Chains[tenv.HomeChainSel]

		// get the header for the destination chain and the relevant block number
		latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[dest] = &block

		// Send to the receiver on the destination chain paying with native token
		var (
			receiver = common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32)
			data     = []byte("")
			feeToken = common.HexToAddress("0x0")
		)

		ccipMessage := router.ClientEVM2AnyMessage{
			Receiver:     receiver,
			Data:         data,
			TokenAmounts: tokens,
			FeeToken:     feeToken,
			ExtraArgs:    nil,
		}

		seqNum := ccdeploy.TestSendRequest(t, e, state, src, dest, false, ccipMessage).SequenceNumber
		expectedSeqNum[dest] = seqNum

		// Wait for all commit reports to land.
		ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

		// After commit is reported on all chains, token prices should be updated in FeeQuoter.
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, ccdeploy.MockLinkPrice, timestampedPrice.Value)

		// Wait for all exec reports to land
		ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

		delete(expectedSeqNum, dest)
	})

	t.Run("Send Token pay with wrapped native home chain -> remote", func(t *testing.T) {
		WETH := state.Chains[tenv.HomeChainSel].Weth9

		// We need to acquire some WETH to send with the tx so we deposit some ETH into the WETH contract
		// but deployerKey is a reference so we need to dereference it and then reassign it
		depositOps := *e.Chains[tenv.HomeChainSel].DeployerKey
		depositOps.Value = new(big.Int).Mul(twoCoins, big.NewInt(2))
		tx, err := WETH.Deposit(&depositOps)
		require.NoError(t, err)
		_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
		require.NoError(t, err)

		// Approve to spend WETH token as feeToken
		tx, err = WETH.Approve(e.Chains[tenv.HomeChainSel].DeployerKey, state.Chains[tenv.HomeChainSel].Router.Address(), new(big.Int).Mul(twoCoins, big.NewInt(2)))
		require.NoError(t, err)
		_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
		require.NoError(t, err)

		// Assign Src to the Home Chain Selector and destChain to the Feed Chain Selector
		src := tenv.HomeChainSel
		dest, destChain := tenv.FeedChainSel, e.Chains[tenv.FeedChainSel]

		// get the header for the destination chain and the relevant block number
		latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[dest] = &block

		// Send to the receiver on the destination chain paying with LINK token
		var (
			receiver = common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32)
			data     = []byte("hello world")
			feeToken = state.Chains[src].Weth9.Address()
		)

		ccipMessage := router.ClientEVM2AnyMessage{
			Receiver:     receiver,
			Data:         data,
			TokenAmounts: tokens,
			FeeToken:     feeToken,
			ExtraArgs:    nil,
		}

		// Get the fee Token Balance Before
		feeTokenBalanceBefore, err := WETH.BalanceOf(nil, e.Chains[src].DeployerKey.From)
		require.NoError(t, err)

		// Check the fee Amount
		srcFee, err := state.Chains[src].Router.GetFee(nil, dest, ccipMessage)
		require.NoError(t, err)

		seqNum := ccdeploy.TestSendRequest(t, e, state, src, dest, false, ccipMessage).SequenceNumber
		expectedSeqNum[dest] = seqNum

		// Check the fee token balance after the request and ensure fee tokens were spent
		feeTokenBalanceAfter, err := WETH.BalanceOf(nil, e.Chains[tenv.HomeChainSel].DeployerKey.From)
		require.NoError(t, err)
		require.Equal(t, feeTokenBalanceAfter, new(big.Int).Sub(feeTokenBalanceBefore, srcFee))

		// Wait for all commit reports to land.
		ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

		// After commit is reported on all chains, token prices should be updated in FeeQuoter.
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, ccdeploy.MockLinkPrice, timestampedPrice.Value)

		// Wait for all exec reports to land
		ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)
	})
}
