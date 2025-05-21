package ccip

import (
	"fmt"
	"strings"
	"testing"

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
		out      messagingtest.TestCaseOutput
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
	ccipChainState := state.AptosChains[destChain]

	t.Run("Sould Succeed Message From Evm to Aptos", func(t *testing.T) {
		message := []byte("Hello Aptos, from EVM!")
		out = messagingtest.Run(t,
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
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
						// TODO: check dummy receiver events
						// dummyReceiver := state.AptosChains[destChain].ReceiverAddress
						// events, err := e.Env.AptosChains[destChain].Client.EventsByHandle(dummyReceiver, fmt.Sprintf("%s::dummy_receiver::ReceivedMessage", dummyReceiver), "received_message_events", nil, nil)
						// require.NoError(t, err)
						// require.Len(t, events, 1)
						// var receivedMessage module_dummy_receiver.ReceivedMessage
						// err = codec.DecodeAptosJsonValue(events[0].Data, &receivedMessage)
						// require.NoError(t, err)
						// require.Equal(t, message, receivedMessage.Data)
					},
				},
			},
		)
	})

	ctx := testcontext.Get(t)
	callOpts := &bind.CallOpts{Context: ctx}

	destChainConfig, err := state.Chains[sourceChain].FeeQuoter.GetDestChainConfig(callOpts, destChain)
	require.NoError(t, err, "Failed to get destination chain config")

	testSetup := mlt.NewTestSetup(
		t,
		state,
		sourceChain,
		destChain,
		common.HexToAddress("0x0"),
		destChainConfig,
		false, // testRouter
		true,  // validateResp
		mlt.WithDeployedEnv(e),
	)

	receiver := state.AptosChains[destChain].ReceiverAddress

	tcs := []mlt.TestCase{
		{
			TestSetup: testSetup,
			Name:      "send message to an EOA",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  common.LeftPadBytes(receiver[:], 32),
				Data:      []byte(strings.Repeat("0", int(testSetup.SrcFeeQuoterDestChainConfig.MaxDataBytes))),
				FeeToken:  common.HexToAddress("0x0"),
				ExtraArgs: testhelpers.MakeEVMExtraArgsV2(100000, true),
			},
		},
		{
			TestSetup: testSetup,
			Name:      "send message with data length equal to maximum data bytes allowed",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  ccipChainState.ReceiverAddress[:],
				Data:      []byte(strings.Repeat("0", int(testSetup.SrcFeeQuoterDestChainConfig.MaxDataBytes))),
				FeeToken:  common.HexToAddress("0x0"),
				ExtraArgs: testhelpers.MakeEVMExtraArgsV2(100000, true),
			},
		},
		{
			TestSetup: testSetup,
			Name:      "send message with gas limit equal to maximum gas limit allowed",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  ccipChainState.ReceiverAddress[:],
				Data:      []byte(strings.Repeat("0", int(testSetup.SrcFeeQuoterDestChainConfig.MaxDataBytes))),
				FeeToken:  common.HexToAddress("0x0"),
				ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(testSetup.SrcFeeQuoterDestChainConfig.MaxPerMsgGasLimit), true),
			},
		},
		{
			TestSetup: testSetup,
			Name:      "send message with data length exceeding maximum data bytes allowed",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  ccipChainState.ReceiverAddress[:],
				Data:      []byte(strings.Repeat("0", int(testSetup.SrcFeeQuoterDestChainConfig.MaxDataBytes)+1)),
				FeeToken:  common.HexToAddress("0x0"),
				ExtraArgs: nil,
			},
			ExpRevert: true,
		},
		{
			TestSetup: testSetup,
			Name:      "send message with gas limit exceeding maximum gas limit allowed",
			Msg: router.ClientEVM2AnyMessage{
				Receiver:  ccipChainState.ReceiverAddress[:],
				Data:      []byte("abc"),
				FeeToken:  common.HexToAddress("0x0"),
				ExtraArgs: testhelpers.MakeEVMExtraArgsV2(uint64(testSetup.SrcFeeQuoterDestChainConfig.MaxPerMsgGasLimit)+1, true),
			},
			ExpRevert: true,
		},
		{
			TestSetup: testSetup,
			Name:      "send message without extra args should fail with invalid args",
			Msg: router.ClientEVM2AnyMessage{
				Receiver: ccipChainState.ReceiverAddress[:],
				Data:     []byte("abc"),
				FeeToken: common.HexToAddress("0x0"),
			},
			ExpRevert: true,
		},
	}

	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[testhelpers.SourceDestPair]uint64)
	expectedSeqNumExec := make(map[testhelpers.SourceDestPair][]uint64)

	for _, tc := range tcs {

		startBlocks[tc.DestChain] = nil

		tco := mlt.Run(tc)

		if tco.MsgSentEvent != nil {
			expectedSeqNum[testhelpers.SourceDestPair{
				SourceChainSelector: tc.SrcChain,
				DestChainSelector:   tc.DestChain,
			}] = tco.MsgSentEvent.SequenceNumber

			expectedSeqNumExec[testhelpers.SourceDestPair{
				SourceChainSelector: tc.SrcChain,
				DestChainSelector:   tc.DestChain,
			}] = []uint64{tco.MsgSentEvent.SequenceNumber}
		}
	}

	fmt.Printf("out: %v\n", out)
}
