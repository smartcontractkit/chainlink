package ccip

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	aptosapi "github.com/aptos-labs/aptos-go-sdk/api"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	mlt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagelimitationstest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIP_Messaging_EVM2Aptos(t *testing.T) {
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithAptosChains(1),
	)

	evmChainSelectors := e.Env.AllChainSelectors()
	aptosChainSelectors := maps.Keys(e.Env.BlockChains.AptosChains())

	fmt.Println("EVM: ", evmChainSelectors)
	fmt.Println("Aptos: ", aptosChainSelectors)

	// Deploy dummy receiver contract
	t.Log("Deploying CCIPDummyReceiver...")
	testhelpers.DeployAptosCCIPReceiver(t, e.Env)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := evmChainSelectors[0]
	destChain := aptosChainSelectors[0]

	t.Log("Source chain (EVM): ", sourceChain, "Dest chain (Aptos): ", destChain)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	var (
		replayed bool
		nonce    uint64
		sender   = common.LeftPadBytes(e.Env.Chains[sourceChain].DeployerKey.From.Bytes(), 32)
		setup    = messagingtest.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
		)

		// Tokens
		NATIVE_FEE_TOKEN = "0x0"
		EVM_LINK_TOKEN   = state.Chains[sourceChain].LinkToken
		WETH_TOKEN       = state.Chains[sourceChain].Weth9
	)

	t.Log("Deploying CCIPDummyReceiver...")
	testhelpers.DeployAptosCCIPReceiver(t, e.Env)
	receiver := state.AptosChains[destChain].ReceiverAddress

	ccipChainState := state.AptosChains[destChain]
	ctx := testcontext.Get(t)
	callOpts := &bind.CallOpts{Context: ctx}
	srcFeeQuoterDestChainConfig, err := state.Chains[sourceChain].FeeQuoter.GetDestChainConfig(callOpts, destChain)
	require.NoError(t, err, "Failed to get destination chain config")

	// grant mint role
	tx, err := EVM_LINK_TOKEN.GrantMintRole(e.Env.Chains[sourceChain].DeployerKey, common.BytesToAddress(sender))
	_, err = cldf.ConfirmIfNoError(e.Env.Chains[sourceChain], tx, err)
	require.NoError(t, err)

	// mint token and approve to router
	tx, err = EVM_LINK_TOKEN.Mint(e.Env.Chains[sourceChain].DeployerKey, common.BytesToAddress(sender), deployment.E18Mult(10_000))
	_, err = cldf.ConfirmIfNoError(e.Env.Chains[sourceChain], tx, err)
	require.NoError(t, err)

	tx, err = EVM_LINK_TOKEN.Approve(e.Env.Chains[sourceChain].DeployerKey, state.Chains[sourceChain].Router.Address(), math.MaxBig256)
	_, err = cldf.ConfirmIfNoError(e.Env.Chains[sourceChain], tx, err)
	require.NoError(t, err)

	// Deposit 1 ETH to get WETH
	wethTransactOpts := *e.Env.Chains[sourceChain].DeployerKey
	wethTransactOpts.Value = deployment.E18Mult(1)
	tx, err = WETH_TOKEN.Deposit(&wethTransactOpts)
	_, err = cldf.ConfirmIfNoError(e.Env.Chains[sourceChain], tx, err)
	require.NoError(t, err)

	tx, err = WETH_TOKEN.Approve(e.Env.Chains[sourceChain].DeployerKey, state.Chains[sourceChain].Router.Address(), math.MaxBig256)
	_, err = cldf.ConfirmIfNoError(e.Env.Chains[sourceChain], tx, err)
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

	t.Run("Hello World Message - Should Succeed", func(t *testing.T) {
		message := []byte("Hello Aptos, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				Replayed:       replayed,
				Nonce:          &nonce,
				ValidationType: messagingtest.ValidationTypeExec,
				Receiver:       ccipChainState.ReceiverAddress[:],
				MsgData:        message,
				// true for out of order execution, which is necessary and enforced for Aptos
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(100000, true),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				FeeToken:               NATIVE_FEE_TOKEN,
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
				Replayed:       replayed,
				Nonce:          &nonce,
				ValidationType: messagingtest.ValidationTypeExec,
				Receiver:       ccipChainState.ReceiverAddress[:],
				MsgData:        message,
				// true for out of order execution, which is necessary and enforced for Aptos
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(100000, true),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				FeeToken:               NATIVE_FEE_TOKEN,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertAptosMessageReceivedMatchesSource(t, e, destChain, receiver, message, 1) },
				},
			},
		)
	})

	t.Run("Max Gas Limit - Should Succeed", func(t *testing.T) {
		message := []byte("Hello Aptos, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				Replayed:       replayed,
				Nonce:          &nonce,
				ValidationType: messagingtest.ValidationTypeExec,
				Receiver:       ccipChainState.ReceiverAddress[:],
				MsgData:        message,
				// true for out of order execution, which is necessary and enforced for Aptos
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit), true),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				FeeToken:               NATIVE_FEE_TOKEN,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertAptosMessageReceivedMatchesSource(t, e, destChain, receiver, message, 2) },
				},
			},
		)
	})

	t.Run("Not Enough Gas on Destination - Should Fail (Status = 3)", func(t *testing.T) {
		message := []byte("Hello Aptos, from EVM!")
		gasLimit := uint64(1) // Obvious failure, but we want to test that the status is 3

		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				Replayed:       replayed,
				Nonce:          &nonce,
				ValidationType: messagingtest.ValidationTypeExec,
				Receiver:       ccipChainState.ReceiverAddress[:],
				MsgData:        message,
				// true for out of order execution, which is necessary and enforced for Aptos
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(gasLimit, true),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_FAILURE,
				FeeToken:               NATIVE_FEE_TOKEN,
			},
		)
	})

	t.Run("Fee Token (LINK) - Should Succeed", func(t *testing.T) {
		message := []byte("Hello Aptos, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				Replayed:       replayed,
				Nonce:          &nonce,
				ValidationType: messagingtest.ValidationTypeExec,
				Receiver:       ccipChainState.ReceiverAddress[:],
				MsgData:        message,
				// true for out of order execution, which is necessary and enforced for Aptos
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit), true),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				FeeToken:               EVM_LINK_TOKEN.Address().String(),
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertAptosMessageReceivedMatchesSource(t, e, destChain, receiver, message, 2) },
				},
			},
		)
	})

	t.Run("Fee Token (WETH) - Should Succeed", func(t *testing.T) {
		message := []byte("Hello Aptos, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:      setup,
				Replayed:       replayed,
				Nonce:          &nonce,
				ValidationType: messagingtest.ValidationTypeExec,
				Receiver:       ccipChainState.ReceiverAddress[:],
				MsgData:        message,
				// true for out of order execution, which is necessary and enforced for Aptos
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit), true),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				FeeToken:               WETH_TOKEN.Address().String(),
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
				FeeToken:  common.HexToAddress(NATIVE_FEE_TOKEN),
				ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(mltTestSetup.SrcFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1, true),
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
				FeeToken:  common.HexToAddress(NATIVE_FEE_TOKEN),
				ExtraArgs: []byte{},
			},
			ExpRevert: true,
		})
	})

	t.Run("OutOfOrder Execution False - Should Fail", func(t *testing.T) {
		message := []byte("Hello Aptos, from EVM!")
		mlt.Run(mlt.TestCase{
			TestSetup: mltTestSetup,
			Name:      "Missing ExtraArgs - Should Fail",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  ccipChainState.ReceiverAddress[:],
				Data:      message,
				FeeToken:  common.HexToAddress(NATIVE_FEE_TOKEN),
				ExtraArgs: testhelpers.MakeEVMExtraArgsV2(100000, false),
			},
			ExpRevert: true,
		})
	})
}

func assertAptosMessageReceivedMatchesSource(t *testing.T, e testhelpers.DeployedEnv, destChain uint64, dummyReceiver aptos.AccountAddress, message []byte, sequenceNumber uint64) {
	events, err := getLatestDummyReceiverEvent(t, e.Env.BlockChains.AptosChains()[destChain].Client, dummyReceiver, sequenceNumber)
	require.NoError(t, err)
	require.Equal(t, 1, len(events))

	data, ok := events[0].Data["data"].(string)
	require.True(t, ok)
	bs, err := hex.DecodeString(data[2:])
	require.NoError(t, err)
	require.Equal(t, message, bs)
}

func getLatestDummyReceiverEvent(t *testing.T, rpcClient aptos.AptosRpcClient, dummyReceiver aptos.AccountAddress, sequenceNumber uint64) ([]*aptosapi.Event, error) {
	limit := uint64(1)
	return rpcClient.EventsByHandle(dummyReceiver, fmt.Sprintf("%s::dummy_receiver::CCIPReceiverState", dummyReceiver.String()), "received_message_events", &sequenceNumber, &limit)
}

func Test_CCIP_Messaging_Aptos2EVM(t *testing.T) {
	ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithAptosChains(1),
	)
	evmChainSelectors := e.Env.AllChainSelectors()
	aptosChainSelectors := maps.Keys(e.Env.BlockChains.AptosChains())

	fmt.Println("EVM: ", evmChainSelectors)
	fmt.Println("Aptos: ", aptosChainSelectors)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := aptosChainSelectors[0]
	destChain := evmChainSelectors[1]

	t.Log("Source chain (Aptos): ", sourceChain, "Dest chain (EVM): ", destChain)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	var (
		replayed      bool
		nonce         uint64
		senderAddress = e.Env.BlockChains.AptosChains()[sourceChain].DeployerSigner.AccountAddress()
		sender        = common.LeftPadBytes(senderAddress[:], 32)
		out           messagingtest.TestCaseOutput
		setup         = messagingtest.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
		)
	)

	t.Run("Message to EVM", func(t *testing.T) {
		latestHead, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		message := []byte("Hello EVM, from Aptos!")
		out = messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Replayed:               replayed,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				FeeToken:               "0xa",
				Receiver:               state.Chains[destChain].Receiver.Address().Bytes(),
				MsgData:                message,
				ExtraArgs:              nil,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
						iter, err := state.Chains[destChain].Receiver.FilterMessageReceived(&bind.FilterOpts{
							Context: ctx,
							Start:   latestHead,
						})
						require.NoError(t, err)
						require.True(t, iter.Next())
						// MessageReceived doesn't emit the data unfortunately, so can't check that.
					},
				},
			},
		)
	})

	fmt.Printf("out: %v\n", out)
}
