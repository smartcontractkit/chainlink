package ccip

import (
	"encoding/base64"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIPMessaging_EVM2Ton(t *testing.T) {
	// Setup 2 chains (EVM and Ton) and a single lane.
	// TODO: Looks like 2 EVM chains are required for test config
	// ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(t, testhelpers.WithTonChains(1))

	t.Logf("Environment: %+v", e.Env)
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// t.Logf("Loaded state: %v", state)
	_ = state

	evmChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	allTonChainSelectors := maps.Keys(e.Env.BlockChains.TonChains())
	sourceChain := evmChainSelectors[0]
	destChain := allTonChainSelectors[0]
	t.Log("EVM chain selectors:", evmChainSelectors,
		", TON chain selectors:", allTonChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector:", sourceChain,
		", dest chain selector:", destChain,
	)
	// connect a single lane, source to dest
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	var (
		nonce  uint64
		sender = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)
		out    mt.TestCaseOutput
		setup  = mt.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
		)
	)

	t.Run("message to contract implementing CCIPReceiver", func(t *testing.T) {
		ccipChainState := state.TonChains[destChain]
		receiver := ccipChainState.ReceiverAddress
		receiverBase64Bytes, err := base64.RawURLEncoding.DecodeString(receiver.String())
		require.NoError(t, err)
		// Prepare 36-byte raw address
		receiver.FlagsToByte()
		out = mt.Run(
			t,
			mt.TestCase{
				TestSetup:              setup,
				Replayed:               true,
				Nonce:                  &nonce,
				Receiver:               receiverBase64Bytes,
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(100000, false),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS, // state would be failed
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
					},
				},
			},
		)
	})

	_ = out
}
