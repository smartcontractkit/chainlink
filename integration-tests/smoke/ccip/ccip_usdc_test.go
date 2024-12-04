package smoke

import (
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

/*
* Chain topology for this test
* 	chainA (USDC, MY_TOKEN)
*			|
*			| ------- chainC (USDC, MY_TOKEN)
*			|
* 	chainB (USDC)
 */
func TestUSDCTokenTransfer(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := tests.Context(t)
	config := &changeset.TestConfigs{
		IsUSDC: true,
	}
	//tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr, config)
	tenv := changeset.NewMemoryEnvironmentWithJobsAndContracts(t, lggr, memory.MemoryEnvironmentConfig{
		Chains:             3,
		NumOfUsersPerChain: 3,
		Nodes:              5,
		Bootstraps:         1,
	}, config)

	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Chains)
	chainA := allChainSelectors[0]
	chainC := allChainSelectors[1]
	chainB := allChainSelectors[2]

	ownerChainA := e.Chains[chainA].DeployerKey
	ownerChainC := e.Chains[chainC].DeployerKey
	ownerChainB := e.Chains[chainB].DeployerKey

	aChainUSDC, cChainUSDC, err := changeset.ConfigureUSDCTokenPools(lggr, e.Chains, chainA, chainC, state)
	require.NoError(t, err)

	bChainUSDC, _, err := changeset.ConfigureUSDCTokenPools(lggr, e.Chains, chainB, chainC, state)
	require.NoError(t, err)

	aChainToken, _, cChainToken, _, err := changeset.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		chainA,
		chainC,
		ownerChainA,
		ownerChainC,
		state,
		e.ExistingAddresses,
		"MY_TOKEN",
	)
	require.NoError(t, err)

	// Add all lanes
	require.NoError(t, changeset.AddLanesForAll(e, state))

	changeset.MintAndAllow(
		t,
		e,
		state,
		map[uint64][]changeset.MintTokenInfo{
			chainA: {
				changeset.NewMintTokenInfo(ownerChainA, aChainUSDC, aChainToken),
			},
			chainB: {
				changeset.NewMintTokenInfo(ownerChainB, bChainUSDC),
			},
			chainC: {
				changeset.NewMintTokenInfo(ownerChainC, cChainUSDC, cChainToken),
			},
		},
	)

	err = updateFeeQuoters(lggr, e, state, chainA, chainB, chainC, aChainUSDC, bChainUSDC, cChainUSDC)
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
			sourceChain: chainC,
			destChain:   chainA,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  cChainUSDC.Address(),
					Amount: tinyOneCoin,
				}},
			expectedTokenBalances: map[common.Address]*big.Int{
				aChainUSDC.Address(): tinyOneCoin,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:        "multiple USDC tokens within the same message",
			receiver:    utils.RandomAddress(),
			sourceChain: chainC,
			destChain:   chainA,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  cChainUSDC.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  cChainUSDC.Address(),
					Amount: tinyOneCoin,
				},
			},
			expectedTokenBalances: map[common.Address]*big.Int{
				// 2 coins because of the same receiver
				aChainUSDC.Address(): new(big.Int).Add(tinyOneCoin, tinyOneCoin),
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:        "USDC token together with another token transferred to EOA",
			receiver:    utils.RandomAddress(),
			sourceChain: chainA,
			destChain:   chainC,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  aChainUSDC.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  aChainToken.Address(),
					Amount: new(big.Int).Mul(tinyOneCoin, big.NewInt(10)),
				},
			},
			expectedTokenBalances: map[common.Address]*big.Int{
				cChainUSDC.Address():  tinyOneCoin,
				cChainToken.Address(): new(big.Int).Mul(tinyOneCoin, big.NewInt(10)),
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:        "programmable token transfer to valid contract receiver",
			receiver:    state.Chains[chainC].Receiver.Address(),
			sourceChain: chainA,
			destChain:   chainC,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  aChainUSDC.Address(),
					Amount: tinyOneCoin,
				},
			},
			data: []byte("hello world"),
			expectedTokenBalances: map[common.Address]*big.Int{
				cChainUSDC.Address(): tinyOneCoin,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
	}

	t.Run("USDC transfers from single source to single dest", func(t *testing.T) {
		startBlocks := make(map[uint64]*uint64)
		expectedTokenBalances := make(map[uint64]map[changeset.TokenReceiverIdentifier]*big.Int)
		expectedSeqNums := make(map[changeset.SourceDestPair]ccipocr3.SeqNumRange)
		expectedExecutionStates := make(map[changeset.SourceDestPair]map[uint64]int)

		for _, tt := range tcs {
			t.Run(tt.name, func(t *testing.T) {
				expectedTokenBalances = changeset.AppendExpectedTokenBalances(
					tt.destChain, tt.receiver, tt.expectedTokenBalances, expectedTokenBalances,
				)

				pairId := changeset.SourceDestPair{
					SourceChainSelector: tt.sourceChain,
					DestChainSelector:   tt.destChain,
				}

				msg, blocks := changeset.Transfer(
					ctx, t, e, state, tt.sourceChain, tt.destChain, tt.tokens, tt.receiver, tt.data, nil)
				if _, ok := expectedExecutionStates[pairId]; !ok {
					expectedExecutionStates[pairId] = make(map[uint64]int)
				}
				expectedExecutionStates[pairId][msg.SequenceNumber] = tt.expectedExecutionState

				if _, ok := startBlocks[tt.destChain]; !ok {
					startBlocks[tt.destChain] = blocks[tt.destChain]
				}

				seqNr, ok := expectedSeqNums[pairId]
				if ok {
					expectedSeqNums[pairId] = ccipocr3.NewSeqNumRange(seqNr.Start(), ccipocr3.SeqNum(msg.SequenceNumber))
				} else {
					expectedSeqNums[pairId] = ccipocr3.NewSeqNumRange(ccipocr3.SeqNum(msg.SequenceNumber), ccipocr3.SeqNum(msg.SequenceNumber))
				}
			})
		}

		err := changeset.ConfirmCommit(
			t,
			e.Chains,
			state.Chains,
			startBlocks,
			false,
			expectedSeqNums,
		)
		require.NoError(t, err)

		execStates := changeset.ConfirmExecWithSeqNrsForAll(
			t,
			e,
			state,
			changeset.SeqNumberRageToSlice(expectedSeqNums),
			startBlocks,
		)
		require.Equal(t, expectedExecutionStates, execStates)

		changeset.WaitForTokenBalances(ctx, t, e.Chains, expectedTokenBalances)
	})

	t.Run("multi-source USDC transfer targeting the same dest receiver", func(t *testing.T) {
		sendSingleTokenTransfer := func(source, dest uint64, token common.Address, receiver common.Address) (*onramp.OnRampCCIPMessageSent, changeset.SourceDestPair) {
			msg := changeset.TestSendRequest(t, e, state, source, dest, false, router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(receiver.Bytes(), 32),
				Data:         []byte{},
				TokenAmounts: []router.ClientEVMTokenAmount{{Token: token, Amount: tinyOneCoin}},
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			})
			return msg, changeset.SourceDestPair{
				SourceChainSelector: source,
				DestChainSelector:   dest,
			}
		}

		receiver := utils.RandomAddress()

		startBlocks := make(map[uint64]*uint64)
		expectedSeqNum := make(map[changeset.SourceDestPair]uint64)
		expectedSeqNumExec := make(map[changeset.SourceDestPair][]uint64)

		latesthdr, err := e.Chains[chainC].Client.HeaderByNumber(testcontext.Get(t), nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[chainC] = &block

		message1, message1ID := sendSingleTokenTransfer(chainA, chainC, aChainUSDC.Address(), receiver)
		expectedSeqNum[message1ID] = message1.SequenceNumber
		expectedSeqNumExec[message1ID] = []uint64{message1.SequenceNumber}

		message2, message2ID := sendSingleTokenTransfer(chainB, chainC, bChainUSDC.Address(), receiver)
		expectedSeqNum[message2ID] = message2.SequenceNumber
		expectedSeqNumExec[message2ID] = []uint64{message2.SequenceNumber}

		changeset.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)
		states := changeset.ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

		require.Equal(t, changeset.EXECUTION_STATE_SUCCESS, states[message1ID][message1.SequenceNumber])
		require.Equal(t, changeset.EXECUTION_STATE_SUCCESS, states[message2ID][message2.SequenceNumber])

		// We sent 1 coin from each source chain, so we should have 2 coins on the destination chain
		// Receiver is randomly generated so we don't need to get the initial balance first
		expectedBalance := new(big.Int).Add(tinyOneCoin, tinyOneCoin)
		changeset.WaitForTheTokenBalance(ctx, t, cChainUSDC.Address(), receiver, e.Chains[chainC], expectedBalance)
	})
}

func updateFeeQuoters(
	lggr logger.Logger,
	e deployment.Environment,
	state changeset.CCIPOnChainState,
	chainA, chainB, chainC uint64,
	aChainUSDC, bChainUSDC, cChainUSDC *burn_mint_erc677.BurnMintERC677,
) error {
	updateFeeQtrGrp := errgroup.Group{}
	updateFeeQtrGrp.Go(func() error {
		return changeset.UpdateFeeQuoterForUSDC(lggr, e.Chains[chainA], state.Chains[chainA], chainC, aChainUSDC)
	})
	updateFeeQtrGrp.Go(func() error {
		return changeset.UpdateFeeQuoterForUSDC(lggr, e.Chains[chainB], state.Chains[chainB], chainC, bChainUSDC)
	})
	updateFeeQtrGrp.Go(func() error {
		return changeset.UpdateFeeQuoterForUSDC(lggr, e.Chains[chainC], state.Chains[chainC], chainA, cChainUSDC)
	})
	return updateFeeQtrGrp.Wait()
}
