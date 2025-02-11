package ccip

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"

	"github.com/smartcontractkit/chainlink-integrations/evm/assets"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// setupTokens deploys transferable tokens on the source and dest, mints tokens for the source and dest, and
// approves the router to spend the tokens
func setupTokens(
	t *testing.T,
	state changeset.CCIPOnChainState,
	tenv testhelpers.DeployedEnv,
	src, dest uint64,
	transferTokenMintAmount,
	feeTokenMintAmount *big.Int,
) (
	srcToken *burn_mint_erc677.BurnMintERC677,
	dstToken *burn_mint_erc677.BurnMintERC677,
) {
	lggr := logger.TestLogger(t)
	e := tenv.Env

	// Deploy the token to test transferring
	srcToken, _, dstToken, _, err := testhelpers.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		src,
		dest,
		tenv.Env.Chains[src].DeployerKey,
		tenv.Env.Chains[dest].DeployerKey,
		state,
		tenv.Env.ExistingAddresses,
		"MY_TOKEN",
	)
	require.NoError(t, err)

	linkToken := state.Chains[src].LinkToken

	tx, err := srcToken.Mint(
		e.Chains[src].DeployerKey,
		e.Chains[src].DeployerKey.From,
		transferTokenMintAmount,
	)
	_, err = deployment.ConfirmIfNoError(e.Chains[src], tx, err)
	require.NoError(t, err)

	// Mint a destination token
	tx, err = dstToken.Mint(
		e.Chains[dest].DeployerKey,
		e.Chains[dest].DeployerKey.From,
		transferTokenMintAmount,
	)
	_, err = deployment.ConfirmIfNoError(e.Chains[dest], tx, err)
	require.NoError(t, err)

	// Approve the router to spend the tokens and confirm the tx's
	// To prevent having to approve the router for every transfer, we approve a sufficiently large amount
	tx, err = srcToken.Approve(e.Chains[src].DeployerKey, state.Chains[src].Router.Address(), math.MaxBig256)
	_, err = deployment.ConfirmIfNoError(e.Chains[src], tx, err)
	require.NoError(t, err)

	tx, err = dstToken.Approve(e.Chains[dest].DeployerKey, state.Chains[dest].Router.Address(), math.MaxBig256)
	_, err = deployment.ConfirmIfNoError(e.Chains[dest], tx, err)
	require.NoError(t, err)

	// Grant mint and burn roles to the deployer key for the newly deployed linkToken
	// Since those roles are not granted automatically
	tx, err = linkToken.GrantMintAndBurnRoles(e.Chains[src].DeployerKey, e.Chains[src].DeployerKey.From)
	_, err = deployment.ConfirmIfNoError(e.Chains[src], tx, err)
	require.NoError(t, err)

	// Mint link token and confirm the tx
	tx, err = linkToken.Mint(
		e.Chains[src].DeployerKey,
		e.Chains[src].DeployerKey.From,
		feeTokenMintAmount,
	)
	_, err = deployment.ConfirmIfNoError(e.Chains[src], tx, err)
	require.NoError(t, err)

	return srcToken, dstToken
}

func Test_CCIPFees(t *testing.T) {
	t.Parallel()
	tenv, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithMultiCall3(),
	)
	e := tenv.Env

	allChains := tenv.Env.AllChainSelectors()
	require.Len(t, allChains, 2, "need two chains for this test")
	sourceChain := allChains[0]
	destChain := allChains[1]

	// Get new state after migration.
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	srcToken, dstToken := setupTokens(
		t,
		state,
		tenv,
		sourceChain,
		destChain,
		deployment.E18Mult(10_000),
		deployment.E18Mult(10_000),
	)

	// Ensure capreg logs are up to date.
	testhelpers.ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)

	// Add all lanes
	testhelpers.AddLanesForAll(t, &tenv, state)

	t.Run("Send programmable token transfer pay with Link token", func(t *testing.T) {
		runFeeTokenTestCase(feeTokenTestCase{
			t:   t,
			dst: destChain,
			src: sourceChain,
			env: tenv,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken.Address(),
					Amount: deployment.E18Mult(2),
				},
			},
			feeToken:           state.Chains[sourceChain].LinkToken.Address(),
			data:               []byte("hello ptt world"),
			receiver:           common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32),
			srcToken:           srcToken,
			dstToken:           dstToken,
			assertTokenBalance: true,
		})
	})

	t.Run("Send programmable token transfer pay with native", func(t *testing.T) {
		runFeeTokenTestCase(feeTokenTestCase{
			t: t,

			// note the order of src and dest is reversed here
			src: destChain,
			dst: sourceChain,

			env: tenv,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  dstToken.Address(),
					Amount: deployment.E18Mult(2),
				},
			},
			feeToken: common.HexToAddress("0x0"),
			data:     []byte("hello ptt world"),
			receiver: common.LeftPadBytes(state.Chains[sourceChain].Receiver.Address().Bytes(), 32),

			// note the order of src and dest is reversed here
			srcToken:           dstToken,
			dstToken:           srcToken,
			assertTokenBalance: true,
		})
	})

	t.Run("Send programmable token transfer pay with wrapped native", func(t *testing.T) {
		runFeeTokenTestCase(feeTokenTestCase{
			t:   t,
			src: sourceChain,
			dst: destChain,
			env: tenv,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken.Address(),
					Amount: deployment.E18Mult(2),
				},
			},
			feeToken:           state.Chains[sourceChain].Weth9.Address(),
			data:               []byte("hello ptt world"),
			receiver:           common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32),
			srcToken:           srcToken,
			dstToken:           dstToken,
			assertTokenBalance: true,
		})
	})

	t.Run("Send programmable token transfer but revert not enough tokens", func(t *testing.T) {
		// Send to the receiver on the destination chain paying with LINK token
		var (
			receiver = common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32)
			data     = []byte("")
			feeToken = state.Chains[sourceChain].LinkToken.Address()
		)

		// Increase the token send amount to more than available to intentionally cause a revert
		ccipMessage := router.ClientEVM2AnyMessage{
			Receiver: receiver,
			Data:     data,
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken.Address(),
					Amount: deployment.E18Mult(100_000_000),
				},
			},
			FeeToken:  feeToken,
			ExtraArgs: nil,
		}

		_, _, err = testhelpers.CCIPSendRequest(
			e,
			state,
			&testhelpers.CCIPSendReqConfig{
				Sender:         e.Chains[sourceChain].DeployerKey,
				IsTestRouter:   true,
				SourceChain:    sourceChain,
				DestChain:      destChain,
				Evm2AnyMessage: ccipMessage,
			},
		)
		require.Error(t, err)
	})

	t.Run("Send data-only message pay with link token", func(t *testing.T) {
		runFeeTokenTestCase(feeTokenTestCase{
			t:   t,
			src: sourceChain,
			dst: destChain,
			env: tenv,
			// no tokens, only data
			tokenAmounts:       nil,
			feeToken:           state.Chains[sourceChain].LinkToken.Address(),
			data:               []byte("hello link world"),
			receiver:           common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32),
			srcToken:           srcToken,
			dstToken:           dstToken,
			assertTokenBalance: false,
		})
	})

	t.Run("Send message pay with native", func(t *testing.T) {
		runFeeTokenTestCase(feeTokenTestCase{
			t:   t,
			src: sourceChain,
			dst: destChain,
			env: tenv,
			// no tokens, only data
			tokenAmounts:       nil,
			feeToken:           common.HexToAddress("0x0"),
			data:               []byte("hello native world"),
			receiver:           common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32),
			srcToken:           srcToken,
			dstToken:           dstToken,
			assertTokenBalance: false,
		})
	})

	t.Run("Send message pay with wrapped native", func(t *testing.T) {
		runFeeTokenTestCase(feeTokenTestCase{
			t:   t,
			src: sourceChain,
			dst: destChain,
			env: tenv,
			// no tokens, only data
			tokenAmounts:       nil,
			feeToken:           state.Chains[sourceChain].Weth9.Address(),
			data:               []byte("hello wrapped native world"),
			receiver:           common.LeftPadBytes(state.Chains[destChain].Receiver.Address().Bytes(), 32),
			srcToken:           srcToken,
			dstToken:           dstToken,
			assertTokenBalance: false,
		})
	})
}

type feeTokenTestCase struct {
	t                  *testing.T
	src, dst           uint64
	env                testhelpers.DeployedEnv
	srcToken, dstToken *burn_mint_erc677.BurnMintERC677
	tokenAmounts       []router.ClientEVMTokenAmount
	feeToken           common.Address
	receiver           []byte
	data               []byte
	assertTokenBalance bool
}

func runFeeTokenTestCase(tc feeTokenTestCase) {
	ctx := tests.Context(tc.t)
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	expectedSeqNum := make(map[testhelpers.SourceDestPair]uint64)
	expectedSeqNumExec := make(map[testhelpers.SourceDestPair][]uint64)

	srcChain := tc.env.Env.Chains[tc.src]
	dstChain := tc.env.Env.Chains[tc.dst]

	state, err := changeset.LoadOnchainState(tc.env.Env)
	require.NoError(tc.t, err)

	var dstTokBalanceBefore *big.Int
	if tc.assertTokenBalance {
		var err error
		dstTokBalanceBefore, err = tc.dstToken.BalanceOf(nil, state.Chains[tc.dst].Receiver.Address())
		require.NoError(tc.t, err)
		tc.t.Logf("destination token balance before of receiver %s: %s",
			state.Chains[tc.dst].Receiver.Address(),
			dstTokBalanceBefore.String())
	}

	// if fee token is not native then approve the router to spend the fee token from the sender.
	var feeTokenWrapper *burn_mint_erc677.BurnMintERC677
	if tc.feeToken != common.HexToAddress("0x0") {
		if tc.feeToken == state.Chains[tc.src].Weth9.Address() {
			// Deposit some ETH into the WETH contract
			weth9, err := weth9.NewWETH9(state.Chains[tc.src].Weth9.Address(), srcChain.Client)
			require.NoError(tc.t, err)

			balance, err := srcChain.Client.BalanceAt(ctx, srcChain.DeployerKey.From, nil)
			require.NoError(tc.t, err)

			tc.t.Logf("balance before deposit: %s", balance.String())

			srcChain.DeployerKey.Value = assets.Ether(100).ToInt()
			tx, err := weth9.Deposit(srcChain.DeployerKey)
			_, err = deployment.ConfirmIfNoError(srcChain, tx, err)
			require.NoError(tc.t, err)
			srcChain.DeployerKey.Value = big.NewInt(0)
		}

		var err error
		feeTokenWrapper, err = burn_mint_erc677.NewBurnMintERC677(tc.feeToken, srcChain.Client)
		require.NoError(tc.t, err)

		// Approve the router to spend fee token
		tx, err := feeTokenWrapper.Approve(srcChain.DeployerKey, state.Chains[tc.src].Router.Address(), math.MaxBig256)

		_, err = deployment.ConfirmIfNoError(srcChain, tx, err)
		require.NoError(tc.t, err)
	}

	// get the header for the destination chain and the relevant block number
	latesthdr, err := dstChain.Client.HeaderByNumber(testcontext.Get(tc.t), nil)
	require.NoError(tc.t, err)
	block := latesthdr.Number.Uint64()
	startBlocks[tc.dst] = &block

	// Get the fee Token Balance Before, if not fee token set get native balance.
	var feeTokenBalanceBefore *big.Int
	if feeTokenWrapper != nil {
		feeTokenBalanceBefore, err = feeTokenWrapper.BalanceOf(&bind.CallOpts{
			Context: ctx,
		}, srcChain.DeployerKey.From)
		require.NoError(tc.t, err)
	} else {
		feeTokenBalanceBefore, err = srcChain.Client.BalanceAt(ctx, srcChain.DeployerKey.From, nil)
		require.NoError(tc.t, err)
	}
	tc.t.Logf("fee token balance before: %s, fee token enabled: %s",
		feeTokenBalanceBefore.String(), tc.feeToken.String())

	msgSentEvent := testhelpers.TestSendRequest(
		tc.t,
		tc.env.Env,
		state,
		tc.src,
		tc.dst,
		false,
		router.ClientEVM2AnyMessage{
			Receiver:     tc.receiver,
			Data:         tc.data,
			TokenAmounts: tc.tokenAmounts,
			FeeToken:     tc.feeToken,
			ExtraArgs:    nil,
		},
	)

	expectedSeqNum[testhelpers.SourceDestPair{
		SourceChainSelector: tc.src,
		DestChainSelector:   tc.dst,
	}] = msgSentEvent.SequenceNumber
	expectedSeqNumExec[testhelpers.SourceDestPair{
		SourceChainSelector: tc.src,
		DestChainSelector:   tc.dst,
	}] = []uint64{msgSentEvent.SequenceNumber}

	// Check the fee token balance after the request and ensure fee tokens were spent
	var feeTokenBalanceAfter *big.Int
	if feeTokenWrapper != nil {
		feeTokenBalanceAfter, err = feeTokenWrapper.BalanceOf(&bind.CallOpts{
			Context: ctx,
		}, srcChain.DeployerKey.From)
		require.NoError(tc.t, err)
	} else {
		feeTokenBalanceAfter, err = srcChain.Client.BalanceAt(ctx, srcChain.DeployerKey.From, nil)
		require.NoError(tc.t, err)
	}
	tc.t.Logf("fee token balance after: %s, fee token: %s, fee paid: %s",
		feeTokenBalanceAfter.String(), tc.feeToken.String(), msgSentEvent.Message.FeeTokenAmount)
	// in the case we have no fee token, native is also used to pay for the tx,
	// so we have to subtract that as well
	if feeTokenWrapper == nil {
		receipt, err := srcChain.Client.TransactionReceipt(ctx, msgSentEvent.Raw.TxHash)
		require.NoError(tc.t, err)
		txCostWei := new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), receipt.EffectiveGasPrice)
		feeTokenBalanceBefore.Sub(feeTokenBalanceBefore, txCostWei)
	}
	require.Equal(
		tc.t,
		feeTokenBalanceAfter,
		new(big.Int).Sub(feeTokenBalanceBefore, msgSentEvent.Message.FeeTokenAmount),
	)

	// Wait for all commit reports to land.
	testhelpers.ConfirmCommitForAllWithExpectedSeqNums(tc.t, tc.env.Env, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.
	linkAddress := state.Chains[tc.dst].LinkToken.Address()
	feeQuoter := state.Chains[tc.dst].FeeQuoter
	timestampedPrice, err := feeQuoter.GetTokenPrice(&bind.CallOpts{
		Context: ctx,
	}, linkAddress)
	require.NoError(tc.t, err)
	require.Equal(tc.t, changeset.MockLinkPrice, timestampedPrice.Value)

	// Wait for all exec reports to land
	testhelpers.ConfirmExecWithSeqNrsForAll(tc.t, tc.env.Env, state, expectedSeqNumExec, startBlocks)

	if tc.assertTokenBalance {
		require.Len(tc.t, tc.tokenAmounts, 1)
		expectedTransferAmount := tc.tokenAmounts[0].Amount

		balanceAfter, err := tc.dstToken.BalanceOf(&bind.CallOpts{
			Context: ctx,
		}, state.Chains[tc.dst].Receiver.Address())
		require.NoError(tc.t, err)
		require.Equal(
			tc.t,
			new(big.Int).Add(dstTokBalanceBefore, expectedTransferAmount),
			balanceAfter,
		)
	}
}
