package ccip

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	aptos_call_opts "github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	aptos_feequoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func assertAptosSourceRevertExpectedError(t *testing.T, err error, execRevertErrorMsg string, execRevertCauseErrorMsg string) {
	require.Error(t, err)
	fmt.Println("Error: ", err.Error())
	require.Contains(t, err.Error(), execRevertErrorMsg)
	require.Contains(t, err.Error(), execRevertCauseErrorMsg)
}

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

	// Deploy the dummy receiver contract
	testhelpers.DeployAptosCCIPReceiver(t, e.Env)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := evmChainSelectors[0]
	destChain := aptosChainSelectors[0]
	deployerSourceChain := e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey
	deployerDestChain := e.Env.BlockChains.AptosChains()[destChain].DeployerSigner.AccountAddress()
	ccipChainState := state.AptosChains[destChain]

	lggr.Debug("Source chain (EVM): ", sourceChain, "Dest chain (Aptos): ", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

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
		{
			Name:           "Send token to EOA with gas limit set to 0",
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
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(0, true),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  aptosToken[:],
					Amount: big.NewInt(1e8),
				},
			},
		},
		{
			Name:           "Send token to Receiver",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       ccipChainState.ReceiverAddress[:],
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
		{
			Name:           "Send token and message to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       deployerDestChain[:],
			Data:           []byte("Hello, World!"),
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

	callOpts := &bind.CallOpts{Context: ctx}
	srcFeeQuoterDestChainConfig, err := state.Chains[sourceChain].FeeQuoter.GetDestChainConfig(callOpts, destChain)
	require.NoError(t, err, "Failed to get destination chain fee quoter config")

	t.Run("Send token to CCIP Receiver setting gas above max gas allowed - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			Receiver:  ccipChainState.ReceiverAddress[:],
			Data:      []byte("Hello, World!"),
			FeeToken:  evmToken.Address(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit), true),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e8),
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
		t.Log("Expected error: ", err)
	})

	t.Run("Send token to CCIP Receiver with token amount set to 0 - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			Receiver:  ccipChainState.ReceiverAddress[:],
			Data:      []byte("Hello, World!"),
			FeeToken:  evmToken.Address(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(100, true),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(0),
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
		t.Log("Expected error: ", err)
	})

	t.Run("Send invalid token to CCIP Receiver - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			Receiver:  ccipChainState.ReceiverAddress[:],
			Data:      []byte("Hello, World!"),
			FeeToken:  evmToken.Address(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit), true),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  common.HexToAddress("0x0000000000000000000000000000000000000000"), // Invalid token
					Amount: big.NewInt(1e8),
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
		t.Log("Expected error: ", err)
	})
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

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	evmToken, _, aptosToken, _, err := testhelpers.DeployTransferableTokenAptos(t, lggr, e.Env, destChain, sourceChain, "TOKEN", &config.TokenMint{
		To:     deployerSourceChain,
		Amount: 10e8,
	})
	require.NoError(t, err)

	// Fee Tokens
	var NativeFeeToken = "0xa" // coin

	// Invalid Fee Token
	var aptosInvalidToken aptos.AccountAddress

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
			FeeToken:  shared.AptosAPTAddress,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(100000), true),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(1e18),
				},
			},
		},
		{
			Name:           "Send token to EOA with gas limit set to 0",
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
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(0), true),
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
			Name:           "Send token and message to EOA without setting ExtraArgs",
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
			FeeToken: NativeFeeToken,
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

	// parse the aptos native fee token hex string into an Aptos AccountAddress

	var aptosFeeToken aptos.AccountAddress
	require.NoError(t, aptosFeeToken.ParseStringRelaxed(NativeFeeToken))

	aptosCallOpts := &aptos_call_opts.CallOpts{}

	aptosFeeQuoter := aptos_feequoter.NewFeeQuoter(
		state.AptosChains[sourceChain].CCIPAddress,
		e.Env.BlockChains.AptosChains()[sourceChain].Client)

	aptosFeeQuoterDestChainConfig, err := aptosFeeQuoter.GetDestChainConfig(aptosCallOpts, destChain)
	require.NoError(t, err, "Failed to get destination chain fee quoter config")

	t.Run("Send token to CCIP Receiver setting gas above max gas allowed - should fail", func(t *testing.T) {
		msg := testhelpers.AptosSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  aptosFeeToken,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(aptosFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1), false),
			TokenAmounts: []testhelpers.AptosTokenAmount{
				{
					Token:  aptosToken,
					Amount: 1e8,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertAptosSourceRevertExpectedError(t, err, "transaction reverted", "E_MESSAGE_GAS_LIMIT_TOO_HIGH")
		t.Log("Expected error: ", err)
	})

	t.Run("Send token to CCIP Receiver with token amount set to 0 - should fail", func(t *testing.T) {
		msg := testhelpers.AptosSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  aptosFeeToken,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(0), true),
			TokenAmounts: []testhelpers.AptosTokenAmount{
				{
					Token:  aptosToken,
					Amount: 0,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertAptosSourceRevertExpectedError(t, err, "transaction reverted", "E_CANNOT_SEND_ZERO_TOKENS")
		t.Log("Expected error: ", err)
	})

	t.Run("Send invalid token to CCIP Receiver - should fail", func(t *testing.T) {
		msg := testhelpers.AptosSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  aptosFeeToken,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(aptosFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1), false),
			TokenAmounts: []testhelpers.AptosTokenAmount{
				{
					Token:  aptosInvalidToken,
					Amount: 1e8,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertAptosSourceRevertExpectedError(t, err, "ABORTED", "invalid_input")
		t.Log("Expected error: ", err)
	})
}

// ########################
// # Burn Mint Token Pool #
// ########################

func Test_CCIP_TokenTransfer_BnM_EVM2Aptos(t *testing.T) {
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithAptosChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	aptosChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyAptos))

	// Deploy the dummy receiver contract
	testhelpers.DeployAptosCCIPReceiver(t, e.Env)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := evmChainSelectors[0]
	destChain := aptosChainSelectors[0]
	deployerSourceChain := e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey
	deployerDestChain := e.Env.BlockChains.AptosChains()[destChain].DeployerSigner.AccountAddress()
	ccipChainState := state.AptosChains[destChain]

	lggr.Debug("Source chain (EVM): ", sourceChain, "Dest chain (Aptos): ", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	evmToken, _, aptosToken, _, err := testhelpers.DeployBnMTokenAptos(t, lggr, e.Env, sourceChain, destChain, "TOKEN", nil)
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
		{
			Name:           "Send token to EOA with gas limit set to 0",
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
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(0, true),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  aptosToken[:],
					Amount: big.NewInt(1e8),
				},
			},
		},
		{
			Name:           "Send token to Receiver",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       ccipChainState.ReceiverAddress[:],
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
		{
			Name:           "Send token and message to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       deployerDestChain[:],
			Data:           []byte("Hello, World!"),
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

	callOpts := &bind.CallOpts{Context: ctx}
	srcFeeQuoterDestChainConfig, err := state.Chains[sourceChain].FeeQuoter.GetDestChainConfig(callOpts, destChain)
	require.NoError(t, err, "Failed to get destination chain fee quoter config")

	t.Run("Send token to CCIP Receiver setting gas above max gas allowed - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			Receiver:  ccipChainState.ReceiverAddress[:],
			Data:      []byte("Hello, World!"),
			FeeToken:  evmToken.Address(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit), true),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e8),
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
		t.Log("Expected error: ", err)
	})

	t.Run("Send token to CCIP Receiver with token amount set to 0 - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			Receiver:  ccipChainState.ReceiverAddress[:],
			Data:      []byte("Hello, World!"),
			FeeToken:  evmToken.Address(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(100, true),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(0),
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
		t.Log("Expected error: ", err)
	})

	t.Run("Send invalid token to CCIP Receiver - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			Receiver:  ccipChainState.ReceiverAddress[:],
			Data:      []byte("Hello, World!"),
			FeeToken:  evmToken.Address(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit), true),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  common.HexToAddress("0x0000000000000000000000000000000000000000"), // Invalid token
					Amount: big.NewInt(1e8),
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
		t.Log("Expected error: ", err)
	})
}

func Test_CCIP_TokenTransfer_BnM_Aptos2EVM(t *testing.T) {
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithAptosChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	aptosChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyAptos))

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

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	evmToken, _, aptosToken, _, err := testhelpers.DeployBnMTokenAptos(t, lggr, e.Env, destChain, sourceChain, "TOKEN", &config.TokenMint{
		To:     deployerSourceChain,
		Amount: 10e8,
	})
	require.NoError(t, err)

	// Fee Tokens
	var NativeFeeToken = "0xa" // coin

	// Invalid Fee Token
	var aptosInvalidToken aptos.AccountAddress

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
			FeeToken:  shared.AptosAPTAddress,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(100000), true),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(1e18),
				},
			},
		},
		{
			Name:           "Send token to EOA with gas limit set to 0",
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
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(0), true),
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
			Name:           "Send token and message to EOA without setting ExtraArgs",
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
			FeeToken: NativeFeeToken,
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

	// parse the aptos native fee token hex string into an Aptos AccountAddress

	var aptosFeeToken aptos.AccountAddress
	require.NoError(t, aptosFeeToken.ParseStringRelaxed(NativeFeeToken))

	aptosCallOpts := &aptos_call_opts.CallOpts{}

	aptosFeeQuoter := aptos_feequoter.NewFeeQuoter(
		state.AptosChains[sourceChain].CCIPAddress,
		e.Env.BlockChains.AptosChains()[sourceChain].Client)

	aptosFeeQuoterDestChainConfig, err := aptosFeeQuoter.GetDestChainConfig(aptosCallOpts, destChain)
	require.NoError(t, err, "Failed to get destination chain fee quoter config")

	t.Run("Send token to CCIP Receiver setting gas above max gas allowed - should fail", func(t *testing.T) {
		msg := testhelpers.AptosSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  aptosFeeToken,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(aptosFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1), false),
			TokenAmounts: []testhelpers.AptosTokenAmount{
				{
					Token:  aptosToken,
					Amount: 1e8,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertAptosSourceRevertExpectedError(t, err, "transaction reverted", "E_MESSAGE_GAS_LIMIT_TOO_HIGH")
		t.Log("Expected error: ", err)
	})

	t.Run("Send token to CCIP Receiver with token amount set to 0 - should fail", func(t *testing.T) {
		msg := testhelpers.AptosSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  aptosFeeToken,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(0), true),
			TokenAmounts: []testhelpers.AptosTokenAmount{
				{
					Token:  aptosToken,
					Amount: 0,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertAptosSourceRevertExpectedError(t, err, "transaction reverted", "E_CANNOT_SEND_ZERO_TOKENS")
		t.Log("Expected error: ", err)
	})

	t.Run("Send invalid token to CCIP Receiver - should fail", func(t *testing.T) {
		msg := testhelpers.AptosSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  aptosFeeToken,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(aptosFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1), false),
			TokenAmounts: []testhelpers.AptosTokenAmount{
				{
					Token:  aptosInvalidToken,
					Amount: 1e8,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertAptosSourceRevertExpectedError(t, err, "ABORTED", "invalid_input")
		t.Log("Expected error: ", err)
	})
}

// ##############################################
// # Lock Release Token Pool - with TransferRef #
// ##############################################

func Test_CCIP_TokenTransfer_LnR_EVM2Aptos(t *testing.T) {
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithAptosChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	aptosChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyAptos))

	// Deploy the dummy receiver contract
	testhelpers.DeployAptosCCIPReceiver(t, e.Env)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := evmChainSelectors[0]
	destChain := aptosChainSelectors[0]
	deployerSourceChain := e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey
	deployerDestChain := e.Env.BlockChains.AptosChains()[destChain].DeployerSigner
	deployerAddressDestChain := deployerDestChain.AccountAddress()
	destChainClient := e.Env.BlockChains.AptosChains()[destChain].Client
	ccipChainState := state.AptosChains[destChain]

	lggr.Debug("Source chain (EVM): ", sourceChain, "Dest chain (Aptos): ", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	evmToken, _, aptosToken, aptosTokenPool, err := testhelpers.DeployLnRTokenAptos(t, lggr, e.Env, sourceChain, destChain, "TOKEN", &config.TokenMint{
		To:     deployerAddressDestChain,
		Amount: 100e8,
	},
		// Enable dispatch hooks
		true,
	)
	require.NoError(t, err)

	// Provide liquidity to the Aptos LnR token pool
	poolStore, err := aptosTokenPool.LockReleaseTokenPool().GetStoreAddress(nil)
	require.NoError(t, err)
	// Send all minted fund to the TP's primary store - the test expects there to be a balance of 0 else the assertions will fail
	payload, err := aptos.FungibleAssetPrimaryStoreTransferPayload(&aptosToken, poolStore, 100e8)
	require.NoError(t, err)
	rawTx, err := destChainClient.BuildTransaction(deployerAddressDestChain, aptos.TransactionPayload{Payload: payload})
	require.NoError(t, err)
	signedTx, err := rawTx.SignedTransaction(deployerDestChain)
	require.NoError(t, err)
	tx, err := destChainClient.SubmitTransaction(signedTx)
	require.NoError(t, err)
	data, err := destChainClient.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to send liquidity to lock release token pool: %v", data.VmStatus)

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
			Receiver:       deployerAddressDestChain[:],
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
		{
			Name:           "Send token to EOA with gas limit set to 0",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       deployerAddressDestChain[:],
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e18),
				},
			},
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(0, true),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  aptosToken[:],
					Amount: big.NewInt(1e8),
				},
			},
		},
		{
			Name:           "Send token to Receiver",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       ccipChainState.ReceiverAddress[:],
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
		{
			Name:           "Send token and message to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       deployerAddressDestChain[:],
			Data:           []byte("Hello, World!"),
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

	callOpts := &bind.CallOpts{Context: ctx}
	srcFeeQuoterDestChainConfig, err := state.Chains[sourceChain].FeeQuoter.GetDestChainConfig(callOpts, destChain)
	require.NoError(t, err, "Failed to get destination chain fee quoter config")

	t.Run("Send token to CCIP Receiver setting gas above max gas allowed - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			Receiver:  ccipChainState.ReceiverAddress[:],
			Data:      []byte("Hello, World!"),
			FeeToken:  evmToken.Address(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit), true),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e8),
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
		t.Log("Expected error: ", err)
	})

	t.Run("Send token to CCIP Receiver with token amount set to 0 - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			Receiver:  ccipChainState.ReceiverAddress[:],
			Data:      []byte("Hello, World!"),
			FeeToken:  evmToken.Address(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(100, true),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(0),
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
		t.Log("Expected error: ", err)
	})

	t.Run("Send invalid token to CCIP Receiver - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			Receiver:  ccipChainState.ReceiverAddress[:],
			Data:      []byte("Hello, World!"),
			FeeToken:  evmToken.Address(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit), true),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  common.HexToAddress("0x0000000000000000000000000000000000000000"), // Invalid token
					Amount: big.NewInt(1e8),
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
		t.Log("Expected error: ", err)
	})
}

func Test_CCIP_TokenTransfer_LnR_Aptos2EVM(t *testing.T) {
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithAptosChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	aptosChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyAptos))

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

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	evmToken, _, aptosToken, _, err := testhelpers.DeployLnRTokenAptos(t, lggr, e.Env, destChain, sourceChain, "TOKEN", &config.TokenMint{
		To:     deployerSourceChain,
		Amount: 110e8,
	},
		// Enable dispatch hooks
		true,
	)
	require.NoError(t, err)

	// Fee Tokens
	var NativeFeeToken = "0xa" // coin

	// Invalid Fee Token
	var aptosInvalidToken aptos.AccountAddress

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
			FeeToken:  shared.AptosAPTAddress,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(100000), true),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(1e18),
				},
			},
		},
		{
			Name:           "Send token to EOA with gas limit set to 0",
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
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(0), true),
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
			Name:           "Send token and message to EOA without setting ExtraArgs",
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
			FeeToken: NativeFeeToken,
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

	// parse the aptos native fee token hex string into an Aptos AccountAddress

	var aptosFeeToken aptos.AccountAddress
	require.NoError(t, aptosFeeToken.ParseStringRelaxed(NativeFeeToken))

	aptosCallOpts := &aptos_call_opts.CallOpts{}

	aptosFeeQuoter := aptos_feequoter.NewFeeQuoter(
		state.AptosChains[sourceChain].CCIPAddress,
		e.Env.BlockChains.AptosChains()[sourceChain].Client)

	aptosFeeQuoterDestChainConfig, err := aptosFeeQuoter.GetDestChainConfig(aptosCallOpts, destChain)
	require.NoError(t, err, "Failed to get destination chain fee quoter config")

	t.Run("Send token to CCIP Receiver setting gas above max gas allowed - should fail", func(t *testing.T) {
		msg := testhelpers.AptosSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  aptosFeeToken,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(aptosFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1), false),
			TokenAmounts: []testhelpers.AptosTokenAmount{
				{
					Token:  aptosToken,
					Amount: 1e8,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertAptosSourceRevertExpectedError(t, err, "transaction reverted", "E_MESSAGE_GAS_LIMIT_TOO_HIGH")
		t.Log("Expected error: ", err)
	})

	t.Run("Send token to CCIP Receiver with token amount set to 0 - should fail", func(t *testing.T) {
		msg := testhelpers.AptosSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  aptosFeeToken,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(0), true),
			TokenAmounts: []testhelpers.AptosTokenAmount{
				{
					Token:  aptosToken,
					Amount: 0,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertAptosSourceRevertExpectedError(t, err, "transaction reverted", "E_CANNOT_SEND_ZERO_TOKENS")
		t.Log("Expected error: ", err)
	})

	t.Run("Send invalid token to CCIP Receiver - should fail", func(t *testing.T) {
		msg := testhelpers.AptosSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  aptosFeeToken,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(aptosFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1), false),
			TokenAmounts: []testhelpers.AptosTokenAmount{
				{
					Token:  aptosInvalidToken,
					Amount: 1e8,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertAptosSourceRevertExpectedError(t, err, "ABORTED", "invalid_input")
		t.Log("Expected error: ", err)
	})
}

// #################################################
// # Lock Release Token Pool - without TransferRef #
// #################################################

func Test_CCIP_TokenTransfer_LnR_without_TransferRef_EVM2Aptos(t *testing.T) {
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithAptosChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	aptosChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyAptos))

	// Deploy the dummy receiver contract
	testhelpers.DeployAptosCCIPReceiver(t, e.Env)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := evmChainSelectors[0]
	destChain := aptosChainSelectors[0]
	deployerSourceChain := e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey
	deployerDestChain := e.Env.BlockChains.AptosChains()[destChain].DeployerSigner
	deployerAddressDestChain := deployerDestChain.AccountAddress()
	destChainClient := e.Env.BlockChains.AptosChains()[destChain].Client
	ccipChainState := state.AptosChains[destChain]

	lggr.Debug("Source chain (EVM): ", sourceChain, "Dest chain (Aptos): ", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	evmToken, _, aptosToken, aptosTokenPool, err := testhelpers.DeployLnRTokenAptos(t, lggr, e.Env, sourceChain, destChain, "TOKEN", &config.TokenMint{
		To:     deployerAddressDestChain,
		Amount: 100e8,
	},
		// Disable dispatch hooks
		false,
	)
	require.NoError(t, err)

	// Provide liquidity to the Aptos LnR token pool
	poolStore, err := aptosTokenPool.LockReleaseTokenPool().GetStoreAddress(nil)
	require.NoError(t, err)
	// Send all minted fund to the TP's primary store - the test expects there to be a balance of 0 else the assertions will fail
	payload, err := aptos.FungibleAssetPrimaryStoreTransferPayload(&aptosToken, poolStore, 100e8)
	require.NoError(t, err)
	rawTx, err := destChainClient.BuildTransaction(deployerAddressDestChain, aptos.TransactionPayload{Payload: payload})
	require.NoError(t, err)
	signedTx, err := rawTx.SignedTransaction(deployerDestChain)
	require.NoError(t, err)
	tx, err := destChainClient.SubmitTransaction(signedTx)
	require.NoError(t, err)
	data, err := destChainClient.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, data.Success, "failed to send liquidity to lock release token pool: %v", data.VmStatus)

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
			Receiver:       deployerAddressDestChain[:],
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
		{
			Name:           "Send token to EOA with gas limit set to 0",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       deployerAddressDestChain[:],
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e18),
				},
			},
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(0, true),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  aptosToken[:],
					Amount: big.NewInt(1e8),
				},
			},
		},
		{
			Name:           "Send token to Receiver",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       ccipChainState.ReceiverAddress[:],
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
		{
			Name:           "Send token and message to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       deployerAddressDestChain[:],
			Data:           []byte("Hello, World!"),
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

	callOpts := &bind.CallOpts{Context: ctx}
	srcFeeQuoterDestChainConfig, err := state.Chains[sourceChain].FeeQuoter.GetDestChainConfig(callOpts, destChain)
	require.NoError(t, err, "Failed to get destination chain fee quoter config")

	t.Run("Send token to CCIP Receiver setting gas above max gas allowed - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			Receiver:  ccipChainState.ReceiverAddress[:],
			Data:      []byte("Hello, World!"),
			FeeToken:  evmToken.Address(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit), true),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e8),
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
		t.Log("Expected error: ", err)
	})

	t.Run("Send token to CCIP Receiver with token amount set to 0 - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			Receiver:  ccipChainState.ReceiverAddress[:],
			Data:      []byte("Hello, World!"),
			FeeToken:  evmToken.Address(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(100, true),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(0),
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
		t.Log("Expected error: ", err)
	})

	t.Run("Send invalid token to CCIP Receiver - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			Receiver:  ccipChainState.ReceiverAddress[:],
			Data:      []byte("Hello, World!"),
			FeeToken:  evmToken.Address(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit), true),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  common.HexToAddress("0x0000000000000000000000000000000000000000"), // Invalid token
					Amount: big.NewInt(1e8),
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
		t.Log("Expected error: ", err)
	})
}

func Test_CCIP_TokenTransfer_LnR_without_TransferRef_Aptos2EVM(t *testing.T) {
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithAptosChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	aptosChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyAptos))

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

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	evmToken, _, aptosToken, _, err := testhelpers.DeployLnRTokenAptos(t, lggr, e.Env, destChain, sourceChain, "TOKEN", &config.TokenMint{
		To:     deployerSourceChain,
		Amount: 110e8,
	},
		// Disable dispatch hooks
		false,
	)
	require.NoError(t, err)

	// Fee Tokens
	var NativeFeeToken = "0xa" // coin

	// Invalid Fee Token
	var aptosInvalidToken aptos.AccountAddress

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
			FeeToken:  shared.AptosAPTAddress,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(100000), true),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(1e18),
				},
			},
		},
		{
			Name:           "Send token to EOA with gas limit set to 0",
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
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(0), true),
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
			Name:           "Send token and message to EOA without setting ExtraArgs",
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
			FeeToken: NativeFeeToken,
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

	// parse the aptos native fee token hex string into an Aptos AccountAddress

	var aptosFeeToken aptos.AccountAddress
	require.NoError(t, aptosFeeToken.ParseStringRelaxed(NativeFeeToken))

	aptosCallOpts := &aptos_call_opts.CallOpts{}

	aptosFeeQuoter := aptos_feequoter.NewFeeQuoter(
		state.AptosChains[sourceChain].CCIPAddress,
		e.Env.BlockChains.AptosChains()[sourceChain].Client)

	aptosFeeQuoterDestChainConfig, err := aptosFeeQuoter.GetDestChainConfig(aptosCallOpts, destChain)
	require.NoError(t, err, "Failed to get destination chain fee quoter config")

	t.Run("Send token to CCIP Receiver setting gas above max gas allowed - should fail", func(t *testing.T) {
		msg := testhelpers.AptosSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  aptosFeeToken,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(aptosFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1), false),
			TokenAmounts: []testhelpers.AptosTokenAmount{
				{
					Token:  aptosToken,
					Amount: 1e8,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertAptosSourceRevertExpectedError(t, err, "transaction reverted", "E_MESSAGE_GAS_LIMIT_TOO_HIGH")
		t.Log("Expected error: ", err)
	})

	t.Run("Send token to CCIP Receiver with token amount set to 0 - should fail", func(t *testing.T) {
		msg := testhelpers.AptosSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  aptosFeeToken,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(0), true),
			TokenAmounts: []testhelpers.AptosTokenAmount{
				{
					Token:  aptosToken,
					Amount: 0,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertAptosSourceRevertExpectedError(t, err, "transaction reverted", "E_CANNOT_SEND_ZERO_TOKENS")
		t.Log("Expected error: ", err)
	})

	t.Run("Send invalid token to CCIP Receiver - should fail", func(t *testing.T) {
		msg := testhelpers.AptosSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  aptosFeeToken,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(aptosFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1), false),
			TokenAmounts: []testhelpers.AptosTokenAmount{
				{
					Token:  aptosInvalidToken,
					Amount: 1e8,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertAptosSourceRevertExpectedError(t, err, "ABORTED", "invalid_input")
		t.Log("Expected error: ", err)
	})
}
