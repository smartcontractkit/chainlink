package smoke

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestInitialDeployOnLocal(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	config := &changeset.TestConfigs{}
	tenv, _, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr, config)
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

	// use this if you are testing locally in memory
	// tenv := changeset.NewMemoryEnvironmentWithJobsAndContracts(t, lggr, 2, 4, config)

	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	// Chain and account setup
	allChainSelectors := maps.Keys(e.Chains)
	sourceChain, destChain := allChainSelectors[0], allChainSelectors[1]
	ownerSourceChain := e.Chains[sourceChain].DeployerKey
	ownerDestChain := e.Chains[destChain].DeployerKey

	// Deploy and fund self-serve actors (when using memory)
	//selfServeSrcTokenPoolDeployer := createAndFundSelfServeActor(t, ownerSourceChain, e.Chains[sourceChain], big.NewInt(1e18))
	//selfServeSrcTokenPoolDeployer := createAndFundSelfServeActor(t, ownerDestChain, e.Chains[destChain], big.NewInt(1e18))
	selfServeSrcTokenPoolDeployer := createDeployerKeyFromSimulatedPrivateKey(t, sourceChain)
	selfServeDestTokenPoolDeployer := createDeployerKeyFromSimulatedPrivateKey(t, destChain)

	// Deploy tokens and pool by CCIP Owner
	srcToken, _, destToken, _, err := changeset.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		sourceChain,
		destChain,
		ownerSourceChain,
		ownerDestChain,
		state,
		e.ExistingAddresses,
		"OWNER_TOKEN",
	)
	require.NoError(t, err)

	// Deploy Self Serve tokens and pool
	selfServeSrcToken, _, selfServeDestToken, _, err := changeset.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		sourceChain,
		destChain,
		selfServeSrcTokenPoolDeployer,
		selfServeDestTokenPoolDeployer,
		state,
		e.ExistingAddresses,
		"SELF_SERVE_TOKEN",
	)
	require.NoError(t, err)

	// Add all lanes.
	require.NoError(t, changeset.AddLanesForAll(e, state))

	// Mint and allow tokens for the router
	changeset.MintAndAllow(t, e, state, map[uint64]*bind.TransactOpts{
		sourceChain: ownerSourceChain,
		destChain:   ownerDestChain,
	}, map[uint64][]*burn_mint_erc677.BurnMintERC677{
		sourceChain: {srcToken},
		destChain:   {destToken},
	})
	changeset.MintAndAllow(t, e, state, map[uint64]*bind.TransactOpts{
		sourceChain: selfServeSrcTokenPoolDeployer,
		destChain:   selfServeDestTokenPoolDeployer,
	}, map[uint64][]*burn_mint_erc677.BurnMintERC677{
		sourceChain: {selfServeSrcToken},
		destChain:   {selfServeDestToken},
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
					Token:  srcToken.Address(),
					Amount: tinyOneCoin,
				},
			},
			receiver: utils.RandomAddress(),
			expectedTokenBalances: map[common.Address]*big.Int{
				destToken.Address(): tinyOneCoin,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:     "Send token to contract",
			srcChain: sourceChain,
			dstChain: destChain,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken.Address(),
					Amount: tinyOneCoin,
				},
			},
			receiver: state.Chains[destChain].Receiver.Address(),
			expectedTokenBalances: map[common.Address]*big.Int{
				destToken.Address(): tinyOneCoin,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:     "Send 2 tokens to receiver",
			srcChain: destChain,
			dstChain: sourceChain,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  destToken.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  selfServeDestToken.Address(),
					Amount: tinyOneCoin,
				},
			},
			receiver: e.Chains[sourceChain].DeployerKey.From,
			expectedTokenBalances: map[common.Address]*big.Int{
				srcToken.Address():          tinyOneCoin,
				selfServeSrcToken.Address(): tinyOneCoin,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
		{
			name:     "Send N tokens to contract",
			srcChain: destChain,
			dstChain: sourceChain,
			tokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  selfServeDestToken.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  destToken.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  selfServeDestToken.Address(),
					Amount: tinyOneCoin,
				},
			},
			receiver: state.Chains[sourceChain].Receiver.Address(),
			expectedTokenBalances: map[common.Address]*big.Int{
				selfServeSrcToken.Address(): new(big.Int).SetUint64(2),
				srcToken.Address():          tinyOneCoin,
			},
			expectedExecutionState: changeset.EXECUTION_STATE_SUCCESS,
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			initialBalances := map[common.Address]*big.Int{}
			for token := range scenario.expectedTokenBalances {
				initialBalance := changeset.GetTokenBalance(t, token, scenario.receiver, e.Chains[scenario.dstChain])
				initialBalances[token] = initialBalance
			}

			changeset.TransferAndWaitForSuccess(
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
				changeset.WaitForTheTokenBalance(t, token, scenario.receiver, e.Chains[scenario.dstChain], expected)
			}
		})
	}
}

// _createAndFundSelfServeActor is currently unused but retained for potential local testing.
// nolint:unused
func _createAndFundSelfServeActor(
	t *testing.T,
	deployer *bind.TransactOpts,
	chain deployment.Chain,
	amountToFund *big.Int,
) *bind.TransactOpts {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	actor, err := bind.NewKeyedTransactorWithChainID(key, getChainIdFromSelector(t, chain.Selector))
	require.NoError(t, err)

	nonce, err := chain.Client.PendingNonceAt(tests.Context(t), deployer.From)
	require.NoError(t, err)

	gasPrice, err := chain.Client.SuggestGasPrice(tests.Context(t))
	require.NoError(t, err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &actor.From,
		Value:    amountToFund,
		Gas:      uint64(21000),
		GasPrice: gasPrice,
		Data:     nil,
	})

	signedTx, err := deployer.Signer(deployer.From, tx)
	require.NoError(t, err)

	err = chain.Client.SendTransaction(tests.Context(t), signedTx)
	require.NoError(t, err)

	_, err = chain.Confirm(signedTx)
	require.NoError(t, err)

	return actor
}

func getChainIdFromSelector(t *testing.T, selector uint64) *big.Int {
	chainId, err := chainselectors.GetChainIDFromSelector(selector)
	require.NoError(t, err)
	//convert chainId from string to big.Int
	chainIdBigInt := new(big.Int)
	chainIdBigInt.SetString(chainId, 10)
	return chainIdBigInt
}

func createDeployerKeyFromSimulatedPrivateKey(t *testing.T, selector uint64) *bind.TransactOpts {
	key, err := crypto.HexToECDSA(networks.AdditionalSimulatedPvtKeys[0])
	require.NoError(t, err)
	chainId := getChainIdFromSelector(t, selector)
	id, err := bind.NewKeyedTransactorWithChainID(key, chainId)
	require.NoError(t, err)
	return id
}
