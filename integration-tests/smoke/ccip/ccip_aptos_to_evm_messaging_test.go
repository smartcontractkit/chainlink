package ccip

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/stretchr/testify/require"

	v1_6 "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

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

	defaultDestinationChainConfig := v1_6.DefaultFeeQuoterDestChainConfig(true, destChain)

	var (
		replayed      bool
		nonce         uint64
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

		STANDARD_MESSAGE = []byte("Hello EVM, from Aptos!")

		// Tokens
		NATIVE_FEE_TOKEN = "0xa"
	)

	require.NoError(t, err)

	t.Run("Message from Aptos to EVM", func(t *testing.T) {
		latestHead, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		message := STANDARD_MESSAGE
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Replayed:               replayed,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				FeeToken:               NATIVE_FEE_TOKEN,
				Receiver:               ccipReceiverAddress,
				MsgData:                message,
				ExtraArgs:              nil,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertEvmMessageReceived(t, ctx, state, destChain, latestHead, message) },
				},
			},
		)
	})

	t.Run("Max Data Bytes - Should Succeed", func(t *testing.T) {
		latestHead, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		message := []byte(strings.Repeat("0", int(defaultDestinationChainConfig.MaxDataBytes)))
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Replayed:               replayed,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				FeeToken:               NATIVE_FEE_TOKEN,
				Receiver:               ccipReceiverAddress,
				MsgData:                message,
				ExtraArgs:              nil,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertEvmMessageReceived(t, ctx, state, destChain, latestHead, message) },
				},
			},
		)
	})

	t.Run("Max Gas Limit - Should Succeed", func(t *testing.T) {
		latestHead, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		message := STANDARD_MESSAGE
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Replayed:               replayed,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				FeeToken:               NATIVE_FEE_TOKEN,
				Receiver:               ccipReceiverAddress,
				MsgData:                message,
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(uint64(defaultDestinationChainConfig.MaxPerMsgGasLimit), false),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) { assertEvmMessageReceived(t, ctx, state, destChain, latestHead, message) },
				},
			},
		)
	})
}

func assertEvmMessageReceived(t *testing.T, ctx context.Context, state stateview.CCIPOnChainState, destChain uint64, latestHead uint64, message []byte) {
	receivedMessage, err := state.Chains[destChain].Receiver.FilterMessageReceived(&bind.FilterOpts{
		Context: ctx,
		Start:   latestHead,
	})
	require.NoError(t, err)
	require.Equal(t, message, receivedMessage.Event.Data, "Message data should match the sent message")
}
