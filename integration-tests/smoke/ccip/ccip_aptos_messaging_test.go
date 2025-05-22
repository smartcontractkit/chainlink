package ccip

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	aptosapi "github.com/aptos-labs/aptos-go-sdk/api"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
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

	evmChainSelectors := maps.Keys(e.Env.Chains)
	aptosChainSelectors := maps.Keys(e.Env.AptosChains)

	fmt.Println("EVM: ", evmChainSelectors)
	fmt.Println("Aptos: ", aptosChainSelectors)

	// Deploy dummy receiver
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
	)

	t.Log("Deploying CCIPDummyReceiver...")
	testhelpers.DeployAptosCCIPReceiver(t, e.Env)
	receiver := state.AptosChains[destChain].ReceiverAddress

	ccipChainState := state.AptosChains[destChain]
	ctx := testcontext.Get(t)
	callOpts := &bind.CallOpts{Context: ctx}
	srcFeeQuoterDestChainConfig, err := state.Chains[sourceChain].FeeQuoter.GetDestChainConfig(callOpts, destChain)
	require.NoError(t, err, "Failed to get destination chain config")

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
				FeeToken:               "0x0",
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
				FeeToken:               "0x0",
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
				FeeToken:               "0x0",
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
				FeeToken:  common.HexToAddress("0x0"),
				ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(mltTestSetup.SrcFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1, true),
			},
			ExpRevert: true,
		})
	})

	// mltTestSetup := mlt.NewTestSetup(
	// 	t,
	// 	state,
	// 	sourceChain,
	// 	destChain,
	// 	common.HexToAddress("0x0"),
	// 	srcFeeQuoterDestChainConfig,
	// 	false, // testRouter
	// 	true,  // validateResp
	// 	mlt.WithDeployedEnv(e),
	// )

	// tcs := []mlt.TestCase{
	// 	{
	// 		TestSetup: mltTestSetup,
	// 		Name:      "send message with gas limit exceeding maximum gas limit allowed",
	// 		Msg: router.ClientEVM2AnyMessage{
	// 			Receiver:  ccipChainState.ReceiverAddress[:],
	// 			Data:      []byte("abc"),
	// 			FeeToken:  common.HexToAddress("0x0"),
	// 			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(mltTestSetup.SrcFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1, true),
	// 		},
	// 		ExpRevert: true,
	// 	},
	// 	{
	// 		TestSetup: mltTestSetup,
	// 		Name:      "send message without extra args should fail with invalid args",
	// 		Msg: router.ClientEVM2AnyMessage{
	// 			Receiver: ccipChainState.ReceiverAddress[:],
	// 			Data:     []byte("abc"),
	// 			FeeToken: common.HexToAddress("0x0"),
	// 		},
	// 		ExpRevert: true,
	// 	},
	// }
}

func assertAptosMessageReceivedMatchesSource(t *testing.T, e testhelpers.DeployedEnv, destChain uint64, dummyReceiver aptos.AccountAddress, message []byte, sequenceNumber uint64) {
	events, err := getLatestDummyReceiverEvent(t, e.Env.AptosChains[destChain].Client, dummyReceiver, sequenceNumber)
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
