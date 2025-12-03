package ccip

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	aptosapi "github.com/aptos-labs/aptos-go-sdk/api"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	aptos_call_opts "github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	aptos_feequoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	mlt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagelimitationstest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_CCIP_Messaging_EVM2Aptos(t *testing.T) {
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

	lggr.Debug("Source chain (EVM): ", sourceChain, "Dest chain (Aptos): ", destChain)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	var (
		nonce  uint64
		sender = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)

		setup = messagingtest.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // test router
		)

		// Tokens
		nativeFeeToken = "0x0"
		evmLinkToken   = state.Chains[sourceChain].LinkToken
		wethToken      = state.Chains[sourceChain].Weth9
	)
	receiver := state.AptosChains[destChain].ReceiverAddress

	ccipChainState := state.AptosChains[destChain]
	ctx := testcontext.Get(t)
	callOpts := &bind.CallOpts{Context: ctx}
	srcFeeQuoterDestChainConfig, err := state.Chains[sourceChain].FeeQuoter.GetDestChainConfig(callOpts, destChain)
	require.NoError(t, err, "Failed to get destination chain config")

	// grant mint role
	tx, err := evmLinkToken.GrantMintRole(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey, common.BytesToAddress(sender))
	_, err = cldf.ConfirmIfNoError(e.Env.BlockChains.EVMChains()[sourceChain], tx, err)
	require.NoError(t, err)

	// mint token and approve to router
	tx, err = evmLinkToken.Mint(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey, common.BytesToAddress(sender), deployment.E18Mult(10_000))
	_, err = cldf.ConfirmIfNoError(e.Env.BlockChains.EVMChains()[sourceChain], tx, err)
	require.NoError(t, err)

	tx, err = evmLinkToken.Approve(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey, state.Chains[sourceChain].Router.Address(), math.MaxBig256)
	_, err = cldf.ConfirmIfNoError(e.Env.BlockChains.EVMChains()[sourceChain], tx, err)
	require.NoError(t, err)

	// Deposit 1 ETH to get WETH
	wethTransactOpts := *e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey
	wethTransactOpts.Value = deployment.E18Mult(1)
	tx, err = wethToken.Deposit(&wethTransactOpts)
	_, err = cldf.ConfirmIfNoError(e.Env.BlockChains.EVMChains()[sourceChain], tx, err)
	require.NoError(t, err)

	tx, err = wethToken.Approve(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey, state.Chains[sourceChain].Router.Address(), math.MaxBig256)
	_, err = cldf.ConfirmIfNoError(e.Env.BlockChains.EVMChains()[sourceChain], tx, err)
	require.NoError(t, err)

	// For testing messages that revert on source
	mltTestSetup := mlt.NewTestSetup(
		t,
		state,
		sourceChain,
		destChain,
		common.HexToAddress("0x0"),
		srcFeeQuoterDestChainConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(e),
	)

	invalidDestChainSelectorTestSetup := mlt.NewTestSetup(
		t,
		state,
		sourceChain,
		destChain,
		common.HexToAddress("0x0"),
		srcFeeQuoterDestChainConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(e),
	)

	t.Run("Hello World Message - Should Succeed", func(t *testing.T) {
		message := []byte("Hello Aptos, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				Nonce:          &nonce,
				ValidationType: messagingtest.ValidationTypeExec,
				Receiver:       ccipChainState.ReceiverAddress[:],
				MsgData:        message,
				// true for out of order execution, which is necessary and enforced for Aptos
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(100000, true),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				FeeToken:               nativeFeeToken,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertAptosMessageReceivedMatchesSource(t, e, destChain, receiver, message, 0) },
				},
			},
		)
	})

	t.Run("Max Data Bytes - Should Succeed", func(t *testing.T) {
		message := []byte(strings.Repeat("0", int(srcFeeQuoterDestChainConfig.MaxDataBytes)))
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				Nonce:          &nonce,
				ValidationType: messagingtest.ValidationTypeExec,
				Receiver:       ccipChainState.ReceiverAddress[:],
				MsgData:        message,
				// true for out of order execution, which is necessary and enforced for Aptos
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(100000, true),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				FeeToken:               nativeFeeToken,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertAptosMessageReceivedMatchesSource(t, e, destChain, receiver, message, 1) },
				},
			},
		)
	})

	t.Run("Max Gas Limit - Should Succeed", func(t *testing.T) {
		t.Skip("TODO: Test fails with gas limit too high, but should succeed. We add a buffer on top so current MaxPerMsgGasLimit seems to high. Unskip once its fixed")
		message := []byte("Hello Aptos, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				Nonce:          &nonce,
				ValidationType: messagingtest.ValidationTypeExec,
				Receiver:       ccipChainState.ReceiverAddress[:],
				MsgData:        message,
				// true for out of order execution, which is necessary and enforced for Aptos
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit), true),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				FeeToken:               nativeFeeToken,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertAptosMessageReceivedMatchesSource(t, e, destChain, receiver, message, 2) },
				},
			},
		)
	})

	t.Run("Not Enough Gas on Destination - Should Fail (Status = 3)", func(t *testing.T) {
		t.Skip("TODO: Unskip this test when we have a fix for this bug")
		message := []byte("Hello Aptos, from EVM!")
		gasLimit := uint64(1) // Obvious failure, but we want to test that the status is 3

		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				Nonce:          &nonce,
				ValidationType: messagingtest.ValidationTypeExec,
				Receiver:       ccipChainState.ReceiverAddress[:],
				MsgData:        message,
				// true for out of order execution, which is necessary and enforced for Aptos
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(gasLimit, true),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_FAILURE,
				FeeToken:               nativeFeeToken,
			},
		)
	})

	t.Run("Fee Token (LINK) - Should Succeed", func(t *testing.T) {
		message := []byte("Hello Aptos, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				Nonce:          &nonce,
				ValidationType: messagingtest.ValidationTypeExec,
				Receiver:       ccipChainState.ReceiverAddress[:],
				MsgData:        message,
				// true for out of order execution, which is necessary and enforced for Aptos
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(100000, true),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				FeeToken:               evmLinkToken.Address().String(),
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertAptosMessageReceivedMatchesSource(t, e, destChain, receiver, message, 2) },
				},
			},
		)
	})

	t.Run("Fee Token (WETH) - Should Succeed", func(t *testing.T) {
		t.Skip("TODO: Unskip this test when fixed, it fails with low level call ERC20 revert")
		message := []byte("Hello Aptos, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				Nonce:          &nonce,
				ValidationType: messagingtest.ValidationTypeExec,
				Receiver:       ccipChainState.ReceiverAddress[:],
				MsgData:        message,
				// true for out of order execution, which is necessary and enforced for Aptos
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(100000, true),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				FeeToken:               wethToken.Address().String(),
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertAptosMessageReceivedMatchesSource(t, e, destChain, receiver, message, 2) },
				},
			},
		)
	})

	t.Run("Max Data Bytes + 1 - Should Fail", func(t *testing.T) {
		message := []byte(strings.Repeat("0", int(srcFeeQuoterDestChainConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Max Data Bytes + 1 - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  ccipChainState.ReceiverAddress[:],
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1, true),
			},
			ExpRevert: true,
		})
	})

	t.Run("Max Data Bytes + 1 to EOA - Should Fail", func(t *testing.T) {
		atposEOAAddress := e.Env.BlockChains.AptosChains()[destChain].DeployerSigner.AccountAddress()
		message := []byte(strings.Repeat("0", int(srcFeeQuoterDestChainConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Max Data Bytes + 1 to EOA - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  atposEOAAddress[:], // Sending to EOA
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1, true),
			},
			ExpRevert: true,
		})
	})

	t.Run("Missing ExtraArgs - Should Fail", func(t *testing.T) {
		message := []byte("Hello Aptos, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Missing ExtraArgs - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  ccipChainState.ReceiverAddress[:],
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: []byte{},
			},
			ExpRevert: true,
		})
	})

	t.Run("OutOfOrder Execution False - Should Fail", func(t *testing.T) {
		message := []byte("Hello Aptos, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "OutOfOrder Execution False - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  ccipChainState.ReceiverAddress[:],
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: testhelpers.MakeEVMExtraArgsV2(100000, false),
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid receiver - Should Fail", func(t *testing.T) {
		message := []byte("Hello Aptos, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Send message to invalid receiver - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  []byte("0x000"),
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: testhelpers.MakeEVMExtraArgsV2(100000, false),
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid chain selector - Should Fail", func(t *testing.T) {
		message := []byte("Hello Aptos, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: invalidDestChainSelectorTestSetup,
			Name:      "Send message to invalid chain selector - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  ccipChainState.ReceiverAddress[:],
				Data:      message,
				FeeToken:  common.HexToAddress(nativeFeeToken),
				ExtraArgs: testhelpers.MakeEVMExtraArgsV2(100000, false),
			},
			ExpRevert: true,
		})
	})
}

func assertAptosMessageReceivedMatchesSource(t *testing.T, e testhelpers.DeployedEnv, destChain uint64, dummyReceiver aptos.AccountAddress, message []byte, sequenceNumber uint64) {
	events, err := getLatestDummyReceiverEvent(t, e.Env.BlockChains.AptosChains()[destChain].Client, dummyReceiver, sequenceNumber)
	require.NoError(t, err)
	require.Len(t, events, 1)

	data, ok := events[0].Data["data"].(string)
	require.True(t, ok)
	bs, err := hex.DecodeString(data[2:])
	require.NoError(t, err)
	require.Equal(t, message, bs)
}

func getLatestDummyReceiverEvent(t *testing.T, rpcClient aptos.AptosRpcClient, dummyReceiver aptos.AccountAddress, sequenceNumber uint64) ([]*aptosapi.Event, error) {
	limit := uint64(1)
	return rpcClient.EventsByHandle(dummyReceiver, dummyReceiver.String()+"::dummy_receiver::CCIPReceiverState", "received_message_events", &sequenceNumber, &limit)
}

func Test_CCIP_Messaging_Aptos2EVM(t *testing.T) {
	ctx := testhelpers.Context(t)
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
	destChain := evmChainSelectors[1]

	lggr.Debug("Source chain (Aptos): ", sourceChain, "Dest chain (EVM): ", destChain)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	aptosCallOpts := &aptos_call_opts.CallOpts{}

	aptosFeeQuoter := aptos_feequoter.NewFeeQuoter(
		state.AptosChains[sourceChain].CCIPAddress,
		e.Env.BlockChains.AptosChains()[sourceChain].Client)

	aptosFeeQuoterDestChainConfig, err := aptosFeeQuoter.GetDestChainConfig(aptosCallOpts, destChain)
	require.NoError(t, err, "Failed to get destination chain config")

	var (
		senderAddress = e.Env.BlockChains.AptosChains()[sourceChain].DeployerSigner.AccountAddress()
		sender        = common.LeftPadBytes(senderAddress[:], 32)
		setup         = messagingtest.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
		)

		ccipReceiverAddress = state.Chains[destChain].Receiver.Address().Bytes()

		standardMessage = []byte("Hello EVM, from Aptos!")

		// Tokens
		nativeFeeToken = "0xa"
	)

	var addr aptos.AccountAddress
	err = addr.ParseStringRelaxed(nativeFeeToken)
	require.NoError(t, err, "Failed to parse address from string")
	aptosNativeFeeTokenAddress := addr

	require.NoError(t, err)

	// For testing messages that revert on source
	mltTestSetup := mlt.NewTestSetup(
		t,
		state,
		sourceChain,
		destChain,
		common.HexToAddress(nativeFeeToken),
		aptosFeeQuoterDestChainConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(e),
	)

	invalidDestChainSelectorTestSetup := mlt.NewTestSetup(
		t,
		state,
		sourceChain,
		destChain,
		common.HexToAddress("0x0"),
		aptosFeeQuoterDestChainConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(e),
	)

	t.Run("Message from Aptos to EVM", func(t *testing.T) {
		latestHead, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		message := standardMessage
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				ValidationType:         messagingtest.ValidationTypeExec,
				FeeToken:               nativeFeeToken,
				Receiver:               ccipReceiverAddress,
				MsgData:                message,
				ExtraArgs:              nil,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertEvmMessageReceived(ctx, t, state, destChain, latestHead, message) },
				},
			},
		)
	})

	t.Run("Max Data Bytes - Should Succeed", func(t *testing.T) {
		latestHead, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		message := []byte(strings.Repeat("0", int(aptosFeeQuoterDestChainConfig.MaxDataBytes)))
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				ValidationType: messagingtest.ValidationTypeExec,
				FeeToken:       nativeFeeToken,
				Receiver:       ccipReceiverAddress,
				MsgData:        message,
				// Just ensuring enough gas is provided to execute the message, doesn't matter if it's way too much
				ExtraArgs:              testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertEvmMessageReceived(ctx, t, state, destChain, latestHead, message) },
				},
			},
		)
	})

	t.Run("Max Gas Limit - Should Succeed", func(t *testing.T) {
		latestHead, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		message := standardMessage
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				ValidationType:         messagingtest.ValidationTypeExec,
				FeeToken:               nativeFeeToken,
				Receiver:               ccipReceiverAddress,
				MsgData:                message,
				ExtraArgs:              testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(aptosFeeQuoterDestChainConfig.MaxPerMsgGasLimit)), false),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertEvmMessageReceived(ctx, t, state, destChain, latestHead, message) },
				},
			},
		)
	})

	t.Run("Max Data Bytes + 1 - Should Fail", func(t *testing.T) {
		message := []byte(strings.Repeat("0", int(aptosFeeQuoterDestChainConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Max Data Bytes + 1 - Should Fail",
			Msg: testhelpers.AptosSendRequest{
				Receiver:  ccipReceiverAddress,
				Data:      message,
				FeeToken:  aptosNativeFeeTokenAddress,
				ExtraArgs: nil,
			},
			ExpRevert: true,
		})
	})

	t.Run("Max Data Bytes + 1 to EOA - Should Fail", func(t *testing.T) {
		message := []byte(strings.Repeat("0", int(aptosFeeQuoterDestChainConfig.MaxDataBytes)+1))
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Max Data Bytes + 1 to EOA - Should Fail",
			Msg: testhelpers.AptosSendRequest{
				Receiver:  e.Env.BlockChains.EVMChains()[destChain].DeployerKey.From[:], // Sending to EOA
				Data:      message,
				FeeToken:  aptosNativeFeeTokenAddress,
				ExtraArgs: nil,
			},
			ExpRevert: true,
		})
	})

	t.Run("Max Gas Limit + 1 - Should Fail", func(t *testing.T) {
		message := standardMessage
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Max Gas Limit + 1 - Should Fail",
			Msg: testhelpers.AptosSendRequest{
				Receiver:  ccipReceiverAddress,
				Data:      message,
				FeeToken:  aptosNativeFeeTokenAddress,
				ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(aptosFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1), false),
			},
			ExpRevert: true,
		})
	})

	t.Run("Missing ExtraArgs - Should Fail", func(t *testing.T) {
		message := standardMessage
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Missing ExtraArgs - Should Fail",
			Msg: testhelpers.AptosSendRequest{
				Receiver:  ccipReceiverAddress,
				Data:      message,
				FeeToken:  aptosNativeFeeTokenAddress,
				ExtraArgs: []byte{},
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid receiver - Should Fail", func(t *testing.T) {
		message := standardMessage
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Send message to invalid receiver - Should Fail",
			Msg: testhelpers.AptosSendRequest{
				Receiver:  []byte("0x0000"),
				Data:      message,
				FeeToken:  aptosNativeFeeTokenAddress,
				ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
			},
			ExpRevert: true,
		})
	})

	t.Run("Send message to invalid chain selector - Should Fail", func(t *testing.T) {
		message := []byte("Hello Aptos, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: invalidDestChainSelectorTestSetup,
			Name:      "Send message to invalid chain selector - Should Fail",
			Msg: testhelpers.AptosSendRequest{
				Receiver:  ccipReceiverAddress,
				Data:      message,
				FeeToken:  aptosNativeFeeTokenAddress,
				ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
			},
			ExpRevert: true,
		})
	})
}

func assertEvmMessageReceived(ctx context.Context, t *testing.T, state stateview.CCIPOnChainState, destChain uint64, latestHead uint64, message []byte) {
	iter, err := state.Chains[destChain].Receiver.FilterMessageReceived(&bind.FilterOpts{
		Context: ctx,
		Start:   latestHead + 1,
	})
	require.NoError(t, err)
	require.True(t, iter.Next())
	require.Equal(t, message, iter.Event.Data, "Message data should match the sent message")
}
