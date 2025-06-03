package ccip

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
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

	t.Log("Source chain (EVM): ", sourceChain, "Dest chain (Aptos): ", destChain)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	var (
		replayed bool
		nonce    uint64
		sender   = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)
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

	t.Run("Message to Aptos", func(t *testing.T) {
		ccipChainState := state.AptosChains[destChain]
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
						// events, err := e.Env.AptosChains[destChain].Client.EventsByHandle(dummyReceiver, fmt.Sprintf("%s::dummy_receiver::CCIPReceiverState", dummyReceiver), "received_message_events", nil, nil)
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

	fmt.Printf("out: %v\n", out)
}

func Test_CCIP_Messaging_Aptos2EVM(t *testing.T) {
	ctx := testhelpers.Context(t)
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
