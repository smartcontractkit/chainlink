package ccip

import (
	"math/big"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
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

func Test_CCIP_RegulatedTokenTransfer_EVM2Aptos(t *testing.T) {
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

	evmToken, _, aptosToken, _, err := testhelpers.DeployRegulatedTransferableTokenAptos(t, lggr, e.Env, sourceChain, destChain, "Regulated Token", nil)
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
			Name:           "Send regulated token to EOA",
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
			Name:           "Send regulated token to EOA with gas limit set to 0",
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
			Name:           "Send regulated token to contract",
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

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(t, e.Env, state, testhelpers.SeqNumberRangeToSlice(expectedSeqNums), startBlocks)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, e.Env, expectedTokenBalances)
}

func Test_CCIP_RegulatedTokenTransfer_Aptos2EVM(t *testing.T) {
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

	lggr.Debug("Source chain (Aptos): ", sourceChain, "Dest chain (EVM): ", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	evmToken, _, aptosToken, _, err := testhelpers.DeployRegulatedTransferableTokenAptos(t, lggr, e.Env, destChain, sourceChain, "Regulated Token", &config.TokenMint{
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
			Name:           "Send regulated token to EOA",
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
			Name:           "Send regulated token to EOA with gas limit set to 0",
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
			Name:           "Send regulated token and message to EOA",
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
			Name:           "Send regulated token and message to EOA without setting ExtraArgs",
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
			Name:           "Send regulated token to Receiver",
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

	t.Run("Send regulated token to CCIP Receiver setting gas above max gas allowed - should fail", func(t *testing.T) {
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

	t.Run("Send regulated token to CCIP Receiver with token amount set to 0 - should fail", func(t *testing.T) {
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

	t.Run("Send invalid regulated token to CCIP Receiver - should fail", func(t *testing.T) {
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
