package smoke

import (
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
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
	lggr := logger.TestLogger(t)
	config := &changeset.TestConfigs{}
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr, config)

	/*
	* use this if you are testing locally in memory
	* tenv := changeset.NewMemoryEnvironmentWithJobsAndContracts(t, lggr, 3, 4, config)
	 */

	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Chains)
	sourceChain := allChainSelectors[0]
	destChain := allChainSelectors[1]

	srcToken1, _, destToken1, _, err := changeset.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		sourceChain,
		destChain,
		state,
		e.ExistingAddresses,
		"MY_TOKEN_1",
	)
	require.NoError(t, err)

	srcToken2, _, destToken2, _, err := changeset.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		sourceChain,
		destChain,
		state,
		e.ExistingAddresses,
		"MY_TOKEN_2",
	)
	require.NoError(t, err)

	// Add all lanes.
	require.NoError(t, changeset.AddLanesForAll(e, state))

	mintAndAllow(t, e, state, map[uint64][]*burn_mint_erc677.BurnMintERC677{
		sourceChain: {srcToken1, srcToken2},
		destChain:   {destToken1, destToken2},
	})

	tinyOneCoin := new(big.Int).SetUint64(1)

	// Test scenarios are defined here
	scenarios := []struct {
		name                   string
		srcChain               uint64
		dstChain               uint64
		tokenAmounts           []router.ClientEVMTokenAmount
		receiver               common.Address
		data                   []byte
		expectedTokenBalances  map[common.Address]*big.Int
		expectedExecutionState int
	}{
		{
			name:     "Send token to EOA",
			srcChain: sourceChain,
			dstChain: destChain,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken1.Address(),
					Amount: tinyOneCoin,
				},
			},
			receiver: utils.RandomAddress(),
			expectedTokenBalances: map[common.Address]*big.Int{
				destToken1.Address(): tinyOneCoin,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:     "Send token to contract",
			srcChain: sourceChain,
			dstChain: destChain,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken1.Address(),
					Amount: tinyOneCoin,
				},
			},
			receiver: state.Chains[destChain].Receiver.Address(),
			expectedTokenBalances: map[common.Address]*big.Int{
				destToken1.Address(): tinyOneCoin,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:     "Send 2 tokens to receiver",
			srcChain: destChain,
			dstChain: sourceChain,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  destToken1.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  destToken2.Address(),
					Amount: tinyOneCoin,
				},
			},
			receiver: e.Chains[sourceChain].DeployerKey.From,
			expectedTokenBalances: map[common.Address]*big.Int{
				srcToken1.Address(): tinyOneCoin,
				srcToken2.Address(): tinyOneCoin,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:     "Send N tokens to contract",
			srcChain: destChain,
			dstChain: sourceChain,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  destToken1.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  destToken2.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  destToken1.Address(),
					Amount: tinyOneCoin,
				},
			},
			receiver: state.Chains[sourceChain].Receiver.Address(),
			expectedTokenBalances: map[common.Address]*big.Int{
				srcToken1.Address(): new(big.Int).SetUint64(2),
				srcToken2.Address(): tinyOneCoin,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			t.Run(scenario.name, func(t *testing.T) {
				initialBalances := map[common.Address]*big.Int{}
				for token := range scenario.expectedTokenBalances {
					initialBalance := getTokenBalance(t, token, scenario.receiver, e.Chains[scenario.dstChain])
					initialBalances[token] = initialBalance
				}

				transferAndWaitForSuccess(
					t,
					e,
					state,
					scenario.srcChain,
					scenario.dstChain,
					scenario.tokenAmounts,
					scenario.receiver,
					scenario.data,
					scenario.expectedExecutionState,
				)

				for token, balance := range scenario.expectedTokenBalances {
					expected := new(big.Int).Add(initialBalances[token], balance)
					waitForTheTokenBalance(t, token, scenario.receiver, e.Chains[scenario.dstChain], expected)
				}
			})
		})
	}
}
