package smoke

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestInitialDeployOnLocal(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr, nil)
	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	// Add all lanes
	require.NoError(t, changeset.AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[changeset.SourceDestPair]uint64)
	expectedSeqNumExec := make(map[changeset.SourceDestPair][]uint64)
	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block
			msgSentEvent := changeset.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
				Data:         []byte("hello world"),
				TokenAmounts: nil,
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			})
			expectedSeqNum[changeset.SourceDestPair{
				SourceChainSelector: src,
				DestChainSelector:   dest,
			}] = msgSentEvent.SequenceNumber
			expectedSeqNumExec[changeset.SourceDestPair{
				SourceChainSelector: src,
				DestChainSelector:   dest,
			}] = []uint64{msgSentEvent.SequenceNumber}
		}
	}

	// Wait for all commit reports to land.
	changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.
	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, changeset.MockLinkPrice, timestampedPrice.Value)
	}

	// Wait for all exec reports to land
	changeset.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

	// TODO: Apply the proposal.
}

func TestTokenTransfer(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr, nil)
	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	srcToken, _, dstToken, _, err := changeset.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		tenv.HomeChainSel,
		tenv.FeedChainSel,
		state,
		e.ExistingAddresses,
		"MY_TOKEN",
	)
	require.NoError(t, err)

	// Add all lanes
	require.NoError(t, changeset.AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[changeset.SourceDestPair]uint64)
	expectedSeqNumExec := make(map[changeset.SourceDestPair][]uint64)

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
				receiver     = common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32)
				data         = []byte("hello world")
				feeToken     = common.HexToAddress("0x0")
				msgSentEvent *onramp.OnRampCCIPMessageSent
			)
			if src == tenv.HomeChainSel && dest == tenv.FeedChainSel {
				msgSentEvent = changeset.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
					Receiver:     receiver,
					Data:         data,
					TokenAmounts: tokens[src],
					FeeToken:     feeToken,
					ExtraArgs:    nil,
				})
			} else {
				msgSentEvent = changeset.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
					Receiver:     receiver,
					Data:         data,
					TokenAmounts: nil,
					FeeToken:     feeToken,
					ExtraArgs:    nil,
				})
			}

			expectedSeqNum[changeset.SourceDestPair{
				SourceChainSelector: src,
				DestChainSelector:   dest,
			}] = msgSentEvent.SequenceNumber
			expectedSeqNumExec[changeset.SourceDestPair{
				SourceChainSelector: src,
				DestChainSelector:   dest,
			}] = []uint64{msgSentEvent.SequenceNumber}
		}
	}

	// Wait for all commit reports to land.
	changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.
	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, changeset.MockLinkPrice, timestampedPrice.Value)
	}

	// Wait for all exec reports to land
	changeset.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

	balance, err := dstToken.BalanceOf(nil, state.Chains[tenv.FeedChainSel].Receiver.Address())
	require.NoError(t, err)
	require.Equal(t, twoCoins, balance)
}

// setupTokens deploys transferable tokens on the source and dest, mints tokens for the source and dest, and
// approves the router to spend the tokens
func setupTokens(t *testing.T, tenv changeset.DeployedEnv) (srcToken *burn_mint_erc677.BurnMintERC677, dstToken *burn_mint_erc677.BurnMintERC677) {
	lggr := logger.TestLogger(t)

	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	tenCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(10))

	// Deploy the token to test transferring
	srcToken, _, dstToken, _, err = changeset.DeployTransferableToken(
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
	// ctx := changeset.Context(t)
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr, nil)
	e := tenv.Env

	// feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	output, err := changeset.DeployPrerequisites(tenv.Env, changeset.DeployPrerequisiteConfig{
		ChainSelectors: tenv.Env.AllChainSelectors(),
	})
	require.NoError(t, err)
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(output.AddressBook))

	expectedSeqNum := make(map[changeset.SourceDestPair]uint64)
	expectedSeqNumExec := make(map[changeset.SourceDestPair][]uint64)

	require.NoError(t, err)
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(output.AddressBook))
	// Get new state after migration.
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	srcToken, dstToken := setupTokens(t, tenv)

	// Ensure capreg logs are up to date.
	changeset.ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)

	// Add all lanes
	require.NoError(t, changeset.AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.

	// Mint 2 tokens to be transferred
	twoCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2))

	linkToken := state.Chains[tenv.HomeChainSel].LinkToken
	maxUint256 := math.MaxBig256

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

	t.Run("Send Token Pay with Link token home chain -> remote", func(t *testing.T) {
		src := tenv.HomeChainSel
		dest, destChain := tenv.FeedChainSel, e.Chains[tenv.FeedChainSel]

		sourceDestPair := changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}

		// Approve to spend link token
		tx, err := linkToken.Approve(e.Chains[src].DeployerKey, state.Chains[src].Router.Address(), maxUint256)
		require.NoError(t, err)
		_, err = e.Chains[src].Confirm(tx)
		require.NoError(t, err)

		// Get the fee Token Balance Before
		srctokenBalance, err := srcToken.BalanceOf(nil, e.Chains[src].DeployerKey.From)
		require.NoError(t, err)
		require.GreaterOrEqual(t, srctokenBalance.Int64(), twoCoins.Int64())

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

		msgSentEvent := changeset.TestSendRequest(t, e, state, src, dest, false, ccipMessage)

		expectedSeqNum[changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}] = msgSentEvent.SequenceNumber
		expectedSeqNumExec[changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}] = []uint64{msgSentEvent.SequenceNumber}

		// Check the fee token balance after the request and ensure fee tokens were spent
		feeTokenBalanceAfter, err := linkToken.BalanceOf(nil, e.Chains[tenv.HomeChainSel].DeployerKey.From)
		require.NoError(t, err)
		require.Equal(t, feeTokenBalanceAfter, new(big.Int).Sub(feeTokenBalanceBefore, srcFee))

		// Wait for all commit reports to land.
		changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

		// After commit is reported on all chains, token prices should be updated in FeeQuoter.
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, changeset.MockLinkPrice, timestampedPrice.Value)

		// Wait for all exec reports to land
		changeset.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

		balance, err := dstToken.BalanceOf(nil, state.Chains[tenv.FeedChainSel].Receiver.Address())
		require.NoError(t, err)
		require.Equal(t, twoCoins, balance)

		// Delete from the expcted seq num the chain that was just tested so that we don't pass it to the
		// commit report in the next test
		delete(expectedSeqNum, sourceDestPair)
	})

	t.Run("Send Token Pay with native remote chain -> home", func(t *testing.T) {
		// Assign Src to the Home Chain Selector and destChain to the Feed Chain Selector
		src := tenv.FeedChainSel
		dest, destChain := tenv.HomeChainSel, e.Chains[tenv.HomeChainSel]

		sourceDestPair := changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
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
			TokenAmounts: tokens[src],
			FeeToken:     feeToken,
			ExtraArgs:    nil,
		}

		msgSentEvent := changeset.TestSendRequest(t, e, state, src, dest, false, ccipMessage)
		expectedSeqNum[changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}] = msgSentEvent.SequenceNumber
		expectedSeqNumExec[changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}] = []uint64{msgSentEvent.SequenceNumber}

		// Wait for all commit reports to land.
		changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

		// After commit is reported on all chains, token prices should be updated in FeeQuoter.
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, changeset.MockLinkPrice, timestampedPrice.Value)

		// Wait for all exec reports to land
		changeset.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

		balance, err := srcToken.BalanceOf(nil, state.Chains[tenv.HomeChainSel].Receiver.Address())
		require.NoError(t, err)
		require.Equal(t, twoCoins, balance)

		delete(expectedSeqNum, sourceDestPair)
	})

	t.Run("Send Token pay with wrapped native home chain -> remote", func(t *testing.T) {
		// Assign Src to the Home Chain Selector and destChain to the Feed Chain Selector
		src := tenv.HomeChainSel
		dest, destChain := tenv.FeedChainSel, e.Chains[tenv.FeedChainSel]

		WETH := state.Chains[src].Weth9

		// Approve to spend WETH token as feeToken
		tx, err := WETH.Approve(e.Chains[tenv.HomeChainSel].DeployerKey, state.Chains[tenv.HomeChainSel].Router.Address(), maxUint256)
		require.NoError(t, err)
		_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
		require.NoError(t, err)

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

		srcFee, err := state.Chains[src].Router.GetFee(nil, dest, ccipMessage)
		require.NoError(t, err)

		// We need to acquire some WETH to send with the tx so we deposit some ETH into the WETH contract
		// but deployerKey is a reference so we need to dereference it and then reassign it
		// Also the fees have been unusually high so we need to deposit more than usual to ensure the test doesn't fail
		depositOps := *e.Chains[tenv.HomeChainSel].DeployerKey
		depositOps.Value = srcFee
		tx, err = WETH.Deposit(&depositOps)
		require.NoError(t, err)
		_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
		require.NoError(t, err)

		// Get the fee Token Balance Before
		feeTokenBalanceBefore, err := WETH.BalanceOf(nil, e.Chains[src].DeployerKey.From)
		require.NoError(t, err)

		msgSentEvent := changeset.TestSendRequest(t, e, state, src, dest, false, ccipMessage)
		expectedSeqNum[changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}] = msgSentEvent.SequenceNumber
		expectedSeqNumExec[changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}] = []uint64{msgSentEvent.SequenceNumber}

		// Check the fee token balance after the request and ensure fee tokens were spent
		feeTokenBalanceAfter, err := WETH.BalanceOf(nil, e.Chains[tenv.HomeChainSel].DeployerKey.From)
		require.NoError(t, err)
		require.Equal(t, feeTokenBalanceAfter, new(big.Int).Sub(feeTokenBalanceBefore, srcFee))

		// Wait for all commit reports to land.
		changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

		// After commit is reported on all chains, token prices should be updated in FeeQuoter.
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, changeset.MockLinkPrice, timestampedPrice.Value)

		// Wait for all exec reports to land
		changeset.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

		balance, err := dstToken.BalanceOf(nil, state.Chains[dest].Receiver.Address())
		require.NoError(t, err)

		// The balance should be 4 since we've already sent 2 tokens over this lane direction in the
		// first part of the test, so balance should already be two
		require.Equal(t, new(big.Int).Mul(twoCoins, big.NewInt(2)), balance)
	})
	t.Run("Send Token but revert not enough tokens", func(t *testing.T) {
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
			data     = []byte("")
			feeToken = linkToken.Address()
		)

		// Increase the token send amount to more than available to intentionally cause a revert
		tokens[src][0].Amount = new(big.Int).Mul(twoCoins, big.NewInt(100))

		ccipMessage := router.ClientEVM2AnyMessage{
			Receiver:     receiver,
			Data:         data,
			TokenAmounts: tokens[src],
			FeeToken:     feeToken,
			ExtraArgs:    nil,
		}

		_, _, err = changeset.CCIPSendRequest(
			e,
			state,
			src, dest,
			true,
			ccipMessage,
		)

		require.Error(t, err)

	})
}

func Test_PricingForMessages(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	// ctx := changeset.Context(t)
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr, nil)
	e := tenv.Env

	// feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	output, err := changeset.DeployPrerequisites(tenv.Env, changeset.DeployPrerequisiteConfig{
		ChainSelectors: tenv.Env.AllChainSelectors(),
	})
	require.NoError(t, err)
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(output.AddressBook))

	expectedSeqNum := make(map[changeset.SourceDestPair]uint64)
	expectedSeqNumExec := make(map[changeset.SourceDestPair][]uint64)

	require.NoError(t, err)
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(output.AddressBook))
	// Get new state after migration.
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	srcToken, dstToken := setupTokens(t, tenv)

	// Ensure capreg logs are up to date.
	changeset.ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)

	// Add all lanes
	require.NoError(t, changeset.AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.

	// Mint 2 tokens to be transferred
	twoCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2))

	linkToken := state.Chains[tenv.HomeChainSel].LinkToken
	maxUint256 := math.MaxBig256

	emptyTokenArray := make([]router.ClientEVMTokenAmount, 0)

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

	t.Run("Send message Pay with Link token home chain -> remote", func(t *testing.T) {
		src := tenv.HomeChainSel
		dest, destChain := tenv.FeedChainSel, e.Chains[tenv.FeedChainSel]

		sourceDestPair := changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}

		// Approve to spend link token
		tx, err := linkToken.Approve(e.Chains[src].DeployerKey, state.Chains[src].Router.Address(), maxUint256)
		require.NoError(t, err)
		_, err = e.Chains[src].Confirm(tx)
		require.NoError(t, err)

		// Get the fee Token Balance Before
		srctokenBalance, err := srcToken.BalanceOf(nil, e.Chains[src].DeployerKey.From)
		require.NoError(t, err)
		require.GreaterOrEqual(t, srctokenBalance.Int64(), twoCoins.Int64())

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

		msgSentEvent := changeset.TestSendRequest(t, e, state, src, dest, false, ccipMessage)

		expectedSeqNum[changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}] = msgSentEvent.SequenceNumber
		expectedSeqNumExec[changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}] = []uint64{msgSentEvent.SequenceNumber}

		// Check the fee token balance after the request and ensure fee tokens were spent
		feeTokenBalanceAfter, err := linkToken.BalanceOf(nil, e.Chains[tenv.HomeChainSel].DeployerKey.From)
		require.NoError(t, err)
		require.Equal(t, feeTokenBalanceAfter, new(big.Int).Sub(feeTokenBalanceBefore, srcFee))

		// Wait for all commit reports to land.
		changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

		// After commit is reported on all chains, token prices should be updated in FeeQuoter.
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, changeset.MockLinkPrice, timestampedPrice.Value)

		// Wait for all exec reports to land
		changeset.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

		balance, err := dstToken.BalanceOf(nil, state.Chains[tenv.FeedChainSel].Receiver.Address())
		require.NoError(t, err)
		require.Equal(t, twoCoins, balance)

		// Delete from the expcted seq num the chain that was just tested so that we don't pass it to the
		// commit report in the next test
		delete(expectedSeqNum, sourceDestPair)
	})

	t.Run("Send message Pay with native remote chain -> home", func(t *testing.T) {
		// Assign Src to the Home Chain Selector and destChain to the Feed Chain Selector
		src := tenv.FeedChainSel
		dest, destChain := tenv.HomeChainSel, e.Chains[tenv.HomeChainSel]

		sourceDestPair := changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}

		// get the header for the destination chain and the relevant block number
		latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[dest] = &block

		// Send to the receiver on the destination chain paying with native token
		var (
			receiver = common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32)
			data     = []byte("hello world")
			feeToken = common.HexToAddress("0x0")
		)

		ccipMessage := router.ClientEVM2AnyMessage{
			Receiver:     receiver,
			Data:         data,
			TokenAmounts: emptyTokenArray,
			FeeToken:     feeToken,
			ExtraArgs:    nil,
		}

		msgSentEvent := changeset.TestSendRequest(t, e, state, src, dest, false, ccipMessage)
		expectedSeqNum[changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}] = msgSentEvent.SequenceNumber
		expectedSeqNumExec[changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}] = []uint64{msgSentEvent.SequenceNumber}

		// Wait for all commit reports to land.
		changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

		// After commit is reported on all chains, token prices should be updated in FeeQuoter.
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, changeset.MockLinkPrice, timestampedPrice.Value)

		// Wait for all exec reports to land
		changeset.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

		delete(expectedSeqNum, sourceDestPair)
	})

	t.Run("Send message pay with wrapped native home chain -> remote", func(t *testing.T) {
		// Assign Src to the Home Chain Selector and destChain to the Feed Chain Selector
		src := tenv.HomeChainSel
		dest, destChain := tenv.FeedChainSel, e.Chains[tenv.FeedChainSel]

		WETH := state.Chains[src].Weth9

		// Approve to spend WETH token as feeToken
		tx, err := WETH.Approve(e.Chains[tenv.HomeChainSel].DeployerKey, state.Chains[tenv.HomeChainSel].Router.Address(), maxUint256)
		require.NoError(t, err)
		_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
		require.NoError(t, err)

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
			data     = []byte("hello world")
			feeToken = state.Chains[src].Weth9.Address()
		)

		ccipMessage := router.ClientEVM2AnyMessage{
			Receiver:     receiver,
			Data:         data,
			TokenAmounts: tokens[src],
			FeeToken:     feeToken,
			ExtraArgs:    nil,
		}

		srcFee, err := state.Chains[src].Router.GetFee(nil, dest, ccipMessage)
		require.NoError(t, err)

		// We need to acquire some WETH to send with the tx so we deposit some ETH into the WETH contract
		// but deployerKey is a reference so we need to dereference it and then reassign it
		// Also the fees have been unusually high so we need to deposit more than usual to ensure the test doesn't fail
		depositOps := *e.Chains[tenv.HomeChainSel].DeployerKey
		depositOps.Value = srcFee
		tx, err = WETH.Deposit(&depositOps)
		require.NoError(t, err)
		_, err = e.Chains[tenv.HomeChainSel].Confirm(tx)
		require.NoError(t, err)

		// Get the fee Token Balance Before
		feeTokenBalanceBefore, err := WETH.BalanceOf(nil, e.Chains[src].DeployerKey.From)
		require.NoError(t, err)

		msgSentEvent := changeset.TestSendRequest(t, e, state, src, dest, false, ccipMessage)
		expectedSeqNum[changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}] = msgSentEvent.SequenceNumber
		expectedSeqNumExec[changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}] = []uint64{msgSentEvent.SequenceNumber}

		// Check the fee token balance after the request and ensure fee tokens were spent
		feeTokenBalanceAfter, err := WETH.BalanceOf(nil, e.Chains[tenv.HomeChainSel].DeployerKey.From)
		require.NoError(t, err)
		require.Equal(t, feeTokenBalanceAfter, new(big.Int).Sub(feeTokenBalanceBefore, srcFee))

		// Wait for all commit reports to land.
		changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

		// After commit is reported on all chains, token prices should be updated in FeeQuoter.
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, changeset.MockLinkPrice, timestampedPrice.Value)

		// Wait for all exec reports to land
		changeset.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

		balance, err := dstToken.BalanceOf(nil, state.Chains[dest].Receiver.Address())
		require.NoError(t, err)

		// The balance should be 4 since we've already sent 2 tokens over this lane direction in the
		// first part of the test, so balance should already be two
		require.Equal(t, new(big.Int).Mul(twoCoins, big.NewInt(2)), balance)
	})
	t.Run("Send Programmable Token Transfers paying in link", func(t *testing.T) {
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

		msgSentEvent := changeset.TestSendRequest(t, e, state, src, dest, false, ccipMessage)

		expectedSeqNum[changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}] = msgSentEvent.SequenceNumber
		expectedSeqNumExec[changeset.SourceDestPair{
			SourceChainSelector: src,
			DestChainSelector:   dest,
		}] = []uint64{msgSentEvent.SequenceNumber}

		// Check the fee token balance after the request and ensure fee tokens were spent
		feeTokenBalanceAfter, err := linkToken.BalanceOf(nil, e.Chains[tenv.HomeChainSel].DeployerKey.From)
		require.NoError(t, err)
		require.Equal(t, feeTokenBalanceAfter, new(big.Int).Sub(feeTokenBalanceBefore, srcFee))

		balanceBefore, err := dstToken.BalanceOf(nil, state.Chains[tenv.FeedChainSel].Receiver.Address())
		require.NoError(t, err)

		// Wait for all commit reports to land.
		changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

		// After commit is reported on all chains, token prices should be updated in FeeQuoter.
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, changeset.MockLinkPrice, timestampedPrice.Value)

		// Wait for all exec reports to land
		changeset.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

		balanceAfter, err := dstToken.BalanceOf(nil, state.Chains[tenv.FeedChainSel].Receiver.Address())
		require.NoError(t, err)
		require.Equal(t, new(big.Int).Add(balanceBefore, twoCoins), balanceAfter)

		// Delete from the expcted seq num the chain that was just tested so that we don't pass it to the
	})

}
