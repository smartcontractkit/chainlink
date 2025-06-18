package feestest

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/weth9"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

// NewFeeTokenTestCase creates a new FeeTokenTestCase to test fee token usage scenarios.
func NewFeeTokenTestCase(
	t *testing.T,
	env cldf.Environment,
	src, dst uint64,
	feeToken common.Address,
	tokenAmounts []router.ClientEVMTokenAmount,
	srcToken, dstToken *burn_mint_erc677.BurnMintERC677,
	receiver []byte,
	data []byte,
	assertTokenBalance, assertExecution bool,
) FeeTokenTestCase {
	return FeeTokenTestCase{
		t:                  t,
		src:                src,
		dst:                dst,
		env:                env,
		srcToken:           srcToken,
		dstToken:           dstToken,
		tokenAmounts:       tokenAmounts,
		feeToken:           feeToken,
		receiver:           receiver,
		data:               data,
		assertTokenBalance: assertTokenBalance,
		assertExecution:    assertExecution,
	}
}

type FeeTokenTestCase struct {
	t                  *testing.T
	src, dst           uint64
	env                cldf.Environment
	srcToken, dstToken *burn_mint_erc677.BurnMintERC677
	tokenAmounts       []router.ClientEVMTokenAmount
	feeToken           common.Address
	receiver           []byte
	data               []byte
	assertTokenBalance bool
	assertExecution    bool
}

func RunFeeTokenTestCase(tc FeeTokenTestCase) {
	ctx := tc.t.Context()
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	expectedSeqNum := make(map[testhelpers.SourceDestPair]uint64)
	expectedSeqNumExec := make(map[testhelpers.SourceDestPair][]uint64)
	evmChains := tc.env.BlockChains.EVMChains()

	srcChain := evmChains[tc.src]
	dstChain := evmChains[tc.dst]

	state, err := stateview.LoadOnchainState(tc.env)
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
			_, err = cldf.ConfirmIfNoError(srcChain, tx, err)
			require.NoError(tc.t, err)
			srcChain.DeployerKey.Value = big.NewInt(0)
		}

		var err error
		feeTokenWrapper, err = burn_mint_erc677.NewBurnMintERC677(tc.feeToken, srcChain.Client)
		require.NoError(tc.t, err)

		bal, err := feeTokenWrapper.BalanceOf(&bind.CallOpts{
			Context: ctx,
		}, srcChain.DeployerKey.From)
		require.NoError(tc.t, err)

		tc.t.Logf("fee token balance before approval: %s", bal.String())

		// Approve the router to spend fee token
		tx, err := feeTokenWrapper.Approve(srcChain.DeployerKey, state.Chains[tc.src].Router.Address(), math.MaxBig256)

		_, err = cldf.ConfirmIfNoError(srcChain, tx, err)
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
		tc.env,
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

	if tc.assertExecution {
		// Wait for all commit reports to land.
		testhelpers.ConfirmCommitForAllWithExpectedSeqNums(tc.t, tc.env, state,
			testhelpers.ToSeqRangeMap(expectedSeqNum), startBlocks)

		// After commit is reported on all chains, token prices should be updated in FeeQuoter.
		linkAddress := state.Chains[tc.dst].LinkToken.Address()
		feeQuoter := state.Chains[tc.dst].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(&bind.CallOpts{
			Context: ctx,
		}, linkAddress)
		require.NoError(tc.t, err)
		require.Equal(tc.t, shared.MockLinkPrice, timestampedPrice.Value)

		// Wait for all exec reports to land
		testhelpers.ConfirmExecWithSeqNrsForAll(tc.t, tc.env, state, expectedSeqNumExec, startBlocks)
	}

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
