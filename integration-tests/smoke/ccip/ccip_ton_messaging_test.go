package ccip

import (
	"encoding/base64"
	"math/big"
	"slices"
	"testing"

	"github.com/xssnick/tonutils-go/tlb"

	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings/onramp"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/codec"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIPMessaging_TON2EVM(t *testing.T) {
	e, _, _ := testsetups.NewIntegrationEnvironment(t, testhelpers.WithTonChains(1))

	t.Logf("Environment: %+v", e.Env)
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Logf("Loaded state: %v", state)
	_ = state

	// make evm chains sorted for deterministic test results
	evmChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	slices.Sort(evmChainSelectors)

	allTonChainSelectors := maps.Keys(e.Env.BlockChains.TonChains())
	sourceChain := allTonChainSelectors[0]
	destChain := evmChainSelectors[0]
	t.Log("TON chain selectors:", allTonChainSelectors,
		", EVM chain selectors:", evmChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector:", sourceChain,
		", dest chain selector:", destChain,
	)

	tonChain := e.Env.BlockChains.TonChains()[sourceChain]
	ac := codec.NewAddressCodec()
	addrBytes, err := ac.AddressStringToBytes(tonChain.WalletAddress.String())
	require.NoError(t, err)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	var (
		sender = addrBytes
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
		receiver := common.LeftPadBytes(e.Env.BlockChains.EVMChains()[destChain].DeployerKey.From.Bytes(), 32)
		require.NoError(t, err)

		ea := onramp.GenericExtraArgsV2{
			GasLimit:                 big.NewInt(1000000),
			AllowOutOfOrderExecution: true,
		}
		c, err := tlb.ToCell(ea)
		require.NoError(t, err)
		out = mt.Run(
			t,
			mt.TestCase{
				Replayed:               true,
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Nonce:                  nil, // TON nonce check is skipped
				Receiver:               receiver,
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              c.ToBOC(),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})

	_ = out
}

func Test_CCIPMessaging_EVM2TON(t *testing.T) {
	t.Skip("Test stalls because TON test assertions aren't implemented yet")
	// Setup 2 chains (EVM and Ton) and a single lane.
	// ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(t, testhelpers.WithTonChains(1))

	t.Logf("Environment: %+v", e.Env)
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Logf("Loaded state: %v", state)
	_ = state

	evmChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	slices.Sort(evmChainSelectors)
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

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	var (
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
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Nonce:                  nil, // TON nonce check is skipped
				Receiver:               receiverBase64Bytes,
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(100000, false),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS, // state would be failed
			},
		)
	})

	_ = out
}
