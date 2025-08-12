package ccip

import (
	"encoding/base64"
	"testing"

	"github.com/xssnick/tonutils-go/tlb"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIPMessaging_EVM2Ton(t *testing.T) {
	t.Skip("Skipping the test temporarily - fix required")
	// Setup 2 chains (EVM and Ton) and a single lane.
	// ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(t, testhelpers.WithTonChains(1))

	t.Logf("Environment: %+v", e.Env)
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Logf("Loaded state: %v", state)
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

	tonChain := e.Env.BlockChains.TonChains()[destChain]
	tonClient := tonChain.Client
	deployerWallet := tonChain.Wallet

	masterInfo, err := tonClient.GetMasterchainInfo(t.Context())
	require.NoError(t, err, "Failed to get masterchain info")
	acc, err := tonClient.GetAccount(t.Context(), masterInfo, deployerWallet.Address())
	require.NoError(t, err, "Failed to get deployer account")
	require.NotNil(t, acc, "Deployer account should not be nil")
	require.NotNil(t, acc.State, "Deployer account state should not be nil")
	require.True(t, acc.IsActive, "Deployer account should be active")

	// Check deployer wallet balance
	expected := tlb.MustFromTON("1000")
	require.GreaterOrEqual(t, acc.State.Balance.Compare(&expected), 0)

	// Should fail, we don't have Fee Quoter support yet for TON chain
	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.Error(t, err, "Expected failure when configuring EVM->TON lane")

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
		t.Skip("Skipping test for now, as it requires a deployed contracts on TON chain")
		ccipChainState := state.TonChains[destChain]
		receiver := ccipChainState.ReceiverAddress
		receiverBase64Bytes, err := base64.RawURLEncoding.DecodeString(receiver.String())
		require.NoError(t, err)
		// Prepare 36-byte raw address
		receiver.FlagsToByte()
		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Nonce:                  &nonce,
				Receiver:               receiverBase64Bytes,
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(100000, false),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS, // state would be failed
			},
		)
	})

	_ = out
}
