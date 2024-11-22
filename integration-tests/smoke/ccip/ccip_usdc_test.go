package smoke

import (
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
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
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

/*
* Chain topology for this test
* 	sourceChain1 (USDC, MY_TOKEN)
*			|
*			| ------- destChain (USDC, MY_TOKEN)
*			|
* 	sourceChain2 (USDC)
 */
func TestUSDCTokenTransfer(t *testing.T) {
	lggr := logger.TestLogger(t)
	config := &changeset.TestConfigs{
		IsUSDC: true,
	}
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr, config)
	//tenv := changeset.NewMemoryEnvironmentWithJobsAndContracts(t, lggr, 3, 4, config)

	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Chains)
	sourceChain1 := allChainSelectors[0]
	destChain := allChainSelectors[1]
	sourceChain2 := allChainSelectors[2]

	sourceUSDC1, destUSDC, err := changeset.ConfigureUSDCTokenPools(lggr, e.Chains, sourceChain1, destChain, state)
	require.NoError(t, err)

	sourceUSDC2, _, err := changeset.ConfigureUSDCTokenPools(lggr, e.Chains, sourceChain2, destChain, state)
	require.NoError(t, err)

	sourceToken1, _, destToken, _, err := changeset.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		sourceChain1,
		destChain,
		state,
		e.ExistingAddresses,
		"MY_TOKEN",
	)
	require.NoError(t, err)

	// Add all lanes
	require.NoError(t, changeset.AddLanesForAll(e, state))

	mintAndAllow(t, e, state, map[uint64][]*burn_mint_erc677.BurnMintERC677{
		sourceChain1: {sourceUSDC1, sourceToken1},
		sourceChain2: {sourceUSDC2},
		destChain:    {destUSDC, destToken},
	})

	err = changeset.UpdateFeeQuoterForUSDC(lggr, e.Chains[sourceChain1], state.Chains[sourceChain1], destChain, sourceUSDC1)
	require.NoError(t, err)

	err = changeset.UpdateFeeQuoterForUSDC(lggr, e.Chains[sourceChain2], state.Chains[sourceChain2], destChain, sourceUSDC2)
	require.NoError(t, err)

	err = changeset.UpdateFeeQuoterForUSDC(lggr, e.Chains[destChain], state.Chains[destChain], sourceChain1, destUSDC)
	require.NoError(t, err)

	// MockE2EUSDCTransmitter always mint 1, see MockE2EUSDCTransmitter.sol for more details
	tinyOneCoin := new(big.Int).SetUint64(1)

	tcs := []struct {
		name                   string
		receiver               common.Address
		sourceChain            uint64
		destChain              uint64
		tokens                 []router.ClientEVMTokenAmount
		data                   []byte
		expectedTokenBalances  map[common.Address]*big.Int
		expectedExecutionState int
	}{
		{
			name:        "single USDC token transfer to EOA",
			receiver:    utils.RandomAddress(),
			sourceChain: destChain,
			destChain:   sourceChain1,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  destUSDC.Address(),
					Amount: tinyOneCoin,
				}},
			expectedTokenBalances: map[common.Address]*big.Int{
				sourceUSDC1.Address(): tinyOneCoin,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:        "multiple USDC tokens within the same message",
			receiver:    utils.RandomAddress(),
			sourceChain: destChain,
			destChain:   sourceChain1,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  destUSDC.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  destUSDC.Address(),
					Amount: tinyOneCoin,
				},
			},
			expectedTokenBalances: map[common.Address]*big.Int{
				// 2 coins because of the same receiver
				sourceUSDC1.Address(): new(big.Int).Add(tinyOneCoin, tinyOneCoin),
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:        "USDC token together with another token transferred to EOA",
			receiver:    utils.RandomAddress(),
			sourceChain: sourceChain1,
			destChain:   destChain,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  sourceUSDC1.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  sourceToken1.Address(),
					Amount: new(big.Int).Mul(tinyOneCoin, big.NewInt(10)),
				},
			},
			expectedTokenBalances: map[common.Address]*big.Int{
				destUSDC.Address():  tinyOneCoin,
				destToken.Address(): new(big.Int).Mul(tinyOneCoin, big.NewInt(10)),
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:        "programmable token transfer to valid contract receiver",
			receiver:    state.Chains[destChain].Receiver.Address(),
			sourceChain: sourceChain1,
			destChain:   destChain,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  sourceUSDC1.Address(),
					Amount: tinyOneCoin,
				},
			},
			data: []byte("hello world"),
			expectedTokenBalances: map[common.Address]*big.Int{
				destUSDC.Address(): tinyOneCoin,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
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
				tt.expectedExecutionState,
			)

			for token, balance := range tt.expectedTokenBalances {
				expected := new(big.Int).Add(initialBalances[token], balance)
				waitForTheTokenBalance(t, token, tt.receiver, e.Chains[tt.destChain], expected)
			}
		})
	}

	t.Run("multi-source USDC transfer targeting the same dest receiver", func(t *testing.T) {
		receiver := utils.RandomAddress()

		startBlocks := make(map[uint64]*uint64)
		expectedSeqNum := make(map[changeset.SourceDestPair]uint64)
		expectedSeqNumExec := make(map[changeset.SourceDestPair][]uint64)

		latesthdr, err := e.Chains[destChain].Client.HeaderByNumber(testcontext.Get(t), nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[destChain] = &block

		message1ID := changeset.SourceDestPair{
			SourceChainSelector: sourceChain1,
			DestChainSelector:   destChain,
		}
		message1 := changeset.TestSendRequest(t, e, state, sourceChain1, destChain, false, router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(receiver.Bytes(), 32),
			Data:         []byte{},
			TokenAmounts: []router.ClientEVMTokenAmount{{Token: sourceUSDC1.Address(), Amount: tinyOneCoin}},
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		})
		expectedSeqNum[message1ID] = message1.SequenceNumber
		expectedSeqNumExec[message1ID] = []uint64{message1.SequenceNumber}

		message2 := changeset.TestSendRequest(t, e, state, sourceChain2, destChain, false, router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(receiver.Bytes(), 32),
			Data:         []byte{},
			TokenAmounts: []router.ClientEVMTokenAmount{{Token: sourceUSDC2.Address(), Amount: tinyOneCoin}},
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		})
		message2ID := changeset.SourceDestPair{
			SourceChainSelector: sourceChain2,
			DestChainSelector:   destChain,
		}
		expectedSeqNum[message2ID] = message2.SequenceNumber
		expectedSeqNumExec[message2ID] = []uint64{message2.SequenceNumber}

		changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)
		states := changeset.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

		require.Equal(t, changeset.EXECUTION_STATE_SUCCESS, states[message1ID][message1.SequenceNumber])
		require.Equal(t, changeset.EXECUTION_STATE_SUCCESS, states[message2ID][message2.SequenceNumber])

		// We sent 1 coin from each source chain, so we should have 2 coins on the destination chain
		expectedBalance := new(big.Int).Add(tinyOneCoin, tinyOneCoin)
		waitForTheTokenBalance(t, destUSDC.Address(), receiver, e.Chains[destChain], expectedBalance)
	})
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
	expectedStatus int,
) {
	identifier := changeset.SourceDestPair{
		SourceChainSelector: sourceChain,
		DestChainSelector:   destChain,
	}

	startBlocks := make(map[uint64]*uint64)
	expectedSeqNum := make(map[changeset.SourceDestPair]uint64)
	expectedSeqNumExec := make(map[changeset.SourceDestPair][]uint64)

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
	expectedSeqNum[identifier] = msgSentEvent.SequenceNumber
	expectedSeqNumExec[identifier] = []uint64{msgSentEvent.SequenceNumber}

	// Wait for all commit reports to land.
	changeset.ConfirmCommitForAllWithExpectedSeqNums(t, env, state, expectedSeqNum, startBlocks)

	// Wait for all exec reports to land
	states := changeset.ConfirmExecWithSeqNrsForAll(t, env, state, expectedSeqNumExec, startBlocks)
	require.Equal(t, expectedStatus, states[identifier][msgSentEvent.SequenceNumber])
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
