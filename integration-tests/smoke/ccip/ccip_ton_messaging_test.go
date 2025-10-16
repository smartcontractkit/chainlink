package ccip

import (
	"math/big"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/xssnick/tonutils-go/tlb"

	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings/onramp"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/codec"
	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIPMessaging_TON2EVM(t *testing.T) {
	// setup environment with 1 ton chain
	e, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithTonChains(1),
	)

	// load state
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// get chain selectors
	allTonChainSelectors := maps.Keys(e.Env.BlockChains.TonChains())
	sourceChain := allTonChainSelectors[0]
	evmChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	slices.Sort(evmChainSelectors) // make evm chains sorted for deterministic test results
	destChain := evmChainSelectors[0]
	t.Log("Chain selectors",
		"TON", allTonChainSelectors,
		"EVM", evmChainSelectors,
		"home", e.HomeChainSel,
		"feed", e.FeedChainSel,
		"source", sourceChain,
		"dest", destChain,
	)

	// setup lane
	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	// encode sender address(deployer address)
	ac := codec.NewAddressCodec()
	tonChain := e.Env.BlockChains.TonChains()[sourceChain]
	addrBytes, err := ac.AddressStringToBytes(tonChain.WalletAddress.String())
	require.NoError(t, err)

	// wait for event filter registration
	t.Logf("Waiting for event filter registration (~2 mins)...")
	testhelpers.WaitForEventFilterRegistrationOnLane(t, state, e.Env.Offchain, sourceChain, destChain)
	// ready to test
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
		extraArgs, err := tlb.ToCell(onramp.GenericExtraArgsV2{
			GasLimit:                 big.NewInt(1000000),
			AllowOutOfOrderExecution: true,
		})
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
				ExtraArgs:              extraArgs.ToBOC(),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})

	_ = out
}

func Test_CCIPMessaging_EVM2TON(t *testing.T) {
	// setup environment with 1 ton chain
	e, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithTonChains(1),
	)

	// load state
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// get chain selectors
	evmChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	slices.Sort(evmChainSelectors)
	allTonChainSelectors := maps.Keys(e.Env.BlockChains.TonChains())
	sourceChain := evmChainSelectors[0]
	destChain := allTonChainSelectors[0]

	t.Log("Chain selectors",
		"TON", allTonChainSelectors,
		"EVM", evmChainSelectors,
		"home", e.HomeChainSel,
		"feed", e.FeedChainSel,
		"source", sourceChain,
		"dest", destChain,
	)
	t.Logf("  OnRamp:       %s", state.Chains[sourceChain].OnRamp.Address())

	// setup lane
	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	// wait for event filter registration
	t.Logf("Waiting for event filter registration (~2 mins)...")
	testhelpers.WaitForEventFilterRegistrationOnLane(t, state, e.Env.Offchain, sourceChain, destChain)

	// ready to test
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

	t.Run("message to contract receiver", func(t *testing.T) {
		offRampAddr := state.TonChains[destChain].OffRamp
		receiverAddr := state.TonChains[destChain].ReceiverAddress

		t.Logf("  TON OffRamp:  %s", offRampAddr.String())
		t.Logf("  TON Receiver: %s", receiverAddr.String())

		ac := codec.NewAddressCodec()
		receiverBytes, err := ac.AddressStringToBytes(receiverAddr.String())
		require.NoError(t, err)
		require.Len(t, receiverBytes, 36, "receiver bytes should be 36 bytes")

		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Nonce:                  nil, // TON nonce check is skipped
				Receiver:               receiverBytes,
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(100000, false),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})
	_ = out
}
