package smoke

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestUSDCTokenTransfer(t *testing.T) {
	lggr := logger.TestLogger(t)
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr)

	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Chains)
	sourceChain := allChainSelectors[0]
	destChain := allChainSelectors[1]

	srcUSDC, dstUSDC, err := changeset.ConfigureUSDCTokenPools(lggr, e.Chains, sourceChain, destChain, state)
	require.NoError(t, err)

	srcToken, _, dstToken, _, err := changeset.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		sourceChain,
		destChain,
		state,
		e.ExistingAddresses,
		"MY_TOKEN",
	)
	require.NoError(t, err)

	// Add all lanes
	require.NoError(t, changeset.AddLanesForAll(e, state))

	mintAndAllow(t, e, state, map[uint64][]*burn_mint_erc677.BurnMintERC677{
		sourceChain: {srcUSDC, srcToken},
		destChain:   {dstUSDC, dstToken},
	})

	err = changeset.UpdateFeeQuoterForUSDC(lggr, e.Chains[sourceChain], state.Chains[sourceChain], destChain, srcUSDC)
	require.NoError(t, err)

	err = changeset.UpdateFeeQuoterForUSDC(lggr, e.Chains[destChain], state.Chains[destChain], sourceChain, dstUSDC)
	require.NoError(t, err)

	// MockE2EUSDCTransmitter always mint 1, see MockE2EUSDCTransmitter.sol for more details
	tinyOneCoin := new(big.Int).SetUint64(1)

	tcs := []struct {
		name                  string
		receiver              common.Address
		sourceChain           uint64
		destChain             uint64
		tokens                []router.ClientEVMTokenAmount
		data                  []byte
		expectedTokenBalances map[common.Address]*big.Int
	}{
		{
			name:        "single USDC token transfer to EOA",
			receiver:    utils.RandomAddress(),
			sourceChain: destChain,
			destChain:   sourceChain,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  dstUSDC.Address(),
					Amount: tinyOneCoin,
				}},
			expectedTokenBalances: map[common.Address]*big.Int{
				srcUSDC.Address(): tinyOneCoin,
			},
		},
		{
			name:        "multiple USDC tokens within the same message",
			receiver:    utils.RandomAddress(),
			sourceChain: destChain,
			destChain:   sourceChain,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  dstUSDC.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  dstUSDC.Address(),
					Amount: tinyOneCoin,
				},
			},
			expectedTokenBalances: map[common.Address]*big.Int{
				// 2 coins because of the same receiver
				srcUSDC.Address(): new(big.Int).Add(tinyOneCoin, tinyOneCoin),
			},
		},
		{
			name:        "USDC token together with another token transferred to EOA",
			receiver:    utils.RandomAddress(),
			sourceChain: sourceChain,
			destChain:   destChain,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  srcUSDC.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  srcToken.Address(),
					Amount: new(big.Int).Mul(tinyOneCoin, big.NewInt(10)),
				},
			},
			expectedTokenBalances: map[common.Address]*big.Int{
				dstUSDC.Address():  tinyOneCoin,
				dstToken.Address(): new(big.Int).Mul(tinyOneCoin, big.NewInt(10)),
			},
		},
		{
			name:        "programmable token transfer to valid contract receiver",
			receiver:    state.Chains[destChain].Receiver.Address(),
			sourceChain: sourceChain,
			destChain:   destChain,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  srcUSDC.Address(),
					Amount: tinyOneCoin,
				},
			},
			data: []byte("hello world"),
			expectedTokenBalances: map[common.Address]*big.Int{
				dstUSDC.Address(): tinyOneCoin,
			},
		},
	}

	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			initialBalances := map[common.Address]*big.Int{}
			for token := range tt.expectedTokenBalances {
				initialBalance := getTokenBalance(t, token, tt.receiver, e.Chains[tt.destChain])
				initialBalances[token] = initialBalance
			}

			transferAndWaitForSuccess(
				t,
				e,
				state,
				tt.sourceChain,
				tt.destChain,
				tt.tokens,
				tt.receiver,
				tt.data,
			)

			for token, balance := range tt.expectedTokenBalances {
				expected := new(big.Int).Add(initialBalances[token], balance)
				waitForTheTokenBalance(t, token, tt.receiver, e.Chains[tt.destChain], expected)
			}
		})
	}
}

// mintAndAllow mints tokens for deployers and allow router to spend them
func mintAndAllow(
	t *testing.T,
	e deployment.Environment,
	state changeset.CCIPOnChainState,
	tkMap map[uint64][]*burn_mint_erc677.BurnMintERC677,
) {
	for chain, tokens := range tkMap {
		for _, token := range tokens {
			twoCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2))

			tx, err := token.Mint(
				e.Chains[chain].DeployerKey,
				e.Chains[chain].DeployerKey.From,
				new(big.Int).Mul(twoCoins, big.NewInt(10)),
			)
			require.NoError(t, err)
			_, err = e.Chains[chain].Confirm(tx)
			require.NoError(t, err)

			tx, err = token.Approve(e.Chains[chain].DeployerKey, state.Chains[chain].Router.Address(), twoCoins)
			require.NoError(t, err)
			_, err = e.Chains[chain].Confirm(tx)
			require.NoError(t, err)
		}
	}
}

// transferAndWaitForSuccess sends a message from sourceChain to destChain and waits for it to be executed
func transferAndWaitForSuccess(
	t *testing.T,
	env deployment.Environment,
	state changeset.CCIPOnChainState,
	sourceChain, destChain uint64,
	tokens []router.ClientEVMTokenAmount,
	receiver common.Address,
	data []byte,
) {
	startBlocks := make(map[uint64]*uint64)
	expectedSeqNum := make(map[changeset.SourceDestPair]uint64)

	latesthdr, err := env.Chains[destChain].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	block := latesthdr.Number.Uint64()
	startBlocks[destChain] = &block

	msgSentEvent := changeset.TestSendRequest(t, env, state, sourceChain, destChain, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(receiver.Bytes(), 32),
		Data:         data,
		TokenAmounts: tokens,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	expectedSeqNum[changeset.SourceDestPair{
		SourceChainSelector: sourceChain,
		DestChainSelector:   destChain,
	}] = msgSentEvent.SequenceNumber

	// Wait for all commit reports to land.
	changeset.ConfirmCommitForAllWithExpectedSeqNums(t, env, state, expectedSeqNum, startBlocks)

	// Wait for all exec reports to land
	changeset.ConfirmExecWithSeqNrForAll(t, env, state, expectedSeqNum, startBlocks)
}

func waitForTheTokenBalance(
	t *testing.T,
	token common.Address,
	receiver common.Address,
	chain deployment.Chain,
	expected *big.Int,
) {
	tokenContract, err := burn_mint_erc677.NewBurnMintERC677(token, chain.Client)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		actualBalance, err := tokenContract.BalanceOf(&bind.CallOpts{Context: tests.Context(t)}, receiver)
		require.NoError(t, err)

		t.Log("Waiting for the token balance",
			"expected", expected,
			"actual", actualBalance,
			"token", token,
			"receiver", receiver,
		)

		return actualBalance.Cmp(expected) == 0
	}, tests.WaitTimeout(t), 100*time.Millisecond)
}

func getTokenBalance(
	t *testing.T,
	token common.Address,
	receiver common.Address,
	chain deployment.Chain,
) *big.Int {
	tokenContract, err := burn_mint_erc677.NewBurnMintERC677(token, chain.Client)
	require.NoError(t, err)

	balance, err := tokenContract.BalanceOf(&bind.CallOpts{Context: tests.Context(t)}, receiver)
	require.NoError(t, err)

	t.Log("Getting token balance",
		"actual", balance,
		"token", token,
		"receiver", receiver,
	)

	return balance
}
