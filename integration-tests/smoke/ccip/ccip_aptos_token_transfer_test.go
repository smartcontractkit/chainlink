package ccip

import (
	"fmt"
	"math/big"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_CCIP_TokenTransfer_EVM2Aptos(t *testing.T) {
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithAptosChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	aptosChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyAptos))

	fmt.Println("EVM: ", evmChainSelectors)
	fmt.Println("Aptos: ", aptosChainSelectors)

	// Deploy the dummy receiver contract
	t.Log("Deploying CCIPDummyReceiver...")
	testhelpers.DeployAptosCCIPReceiver(t, e.Env)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := evmChainSelectors[0]
	destChain := aptosChainSelectors[0]
	deployerSourceChain := e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey
	deployerDestChain := e.Env.BlockChains.AptosChains()[destChain].DeployerSigner.AccountAddress()

	t.Log("Source chain (EVM): ", sourceChain, "Dest chain (Aptos): ", destChain)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	evmToken, _, aptosToken, _, err := testhelpers.DeployTransferableTokenAptos(t, lggr, e.Env, sourceChain, destChain, "TOKEN", nil)
	require.NoError(t, err)

	testhelpers.MintAndAllow(
		t,
		e.Env,
		state,
		map[uint64][]testhelpers.MintTokenInfo{
			sourceChain: {
				testhelpers.NewMintTokenInfo(deployerSourceChain, evmToken),
			},
		},
	)

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       deployerDestChain[:],
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e18),
				},
			},
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(100000, true),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  aptosToken[:],
					Amount: big.NewInt(1e8),
				},
			},
		},
	}

	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		e.Env,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		e.Env,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, e.Env, expectedTokenBalances)
}
func Test_CCIP_TokenTransfer_Aptos2EVM(t *testing.T) {
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithAptosChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	aptosChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyAptos))

	fmt.Println("EVM: ", evmChainSelectors)
	fmt.Println("Aptos: ", aptosChainSelectors)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := aptosChainSelectors[0]
	destChain := evmChainSelectors[0]

	deployerSourceChain := e.Env.BlockChains.AptosChains()[sourceChain].DeployerSigner.AccountAddress()
	deployerDestChain := e.Env.BlockChains.EVMChains()[destChain].DeployerKey

	// Chain State
	destChainState := state.Chains[destChain]

	// Receiver Address
	ccipReceiverAddress := destChainState.Receiver.Address()

	t.Log("Source chain (EVM): ", sourceChain, "Dest chain (Aptos): ", destChain)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	evmToken, _, aptosToken, _, err := testhelpers.DeployTransferableTokenAptos(t, lggr, e.Env, destChain, sourceChain, "TOKEN", &config.TokenMint{
		To:     deployerSourceChain,
		Amount: 10e8,
	})
	require.NoError(t, err)

	// Fee Tokens
	var NativeFeeToken = "0xa" // coin

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       deployerDestChain.From.Bytes(),
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			AptosTokens: []testhelpers.AptosTokenAmount{
				{
					Token:  aptosToken,
					Amount: 1e8,
				},
			},
			FeeToken:  "0xa",
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(100000), true),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(1e18),
				},
			},
		},
		{
			Name:           "Send token and message to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       deployerDestChain.From.Bytes(),
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			Data:           []byte("Hello, World!"),
			AptosTokens: []testhelpers.AptosTokenAmount{
				{
					Token:  aptosToken,
					Amount: 1e8,
				},
			},
			FeeToken:  NativeFeeToken,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(0), true),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(1e18),
				},
			},
		},
		{
			Name:           "Send token to Receiver",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       ccipReceiverAddress.Bytes(),
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			AptosTokens: []testhelpers.AptosTokenAmount{
				{
					Token:  aptosToken,
					Amount: 1e8,
				},
			},
			FeeToken:  NativeFeeToken,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(0), true),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(1e18),
				},
			},
		},
	}

	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		e.Env,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		e.Env,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, e.Env, expectedTokenBalances)
}
