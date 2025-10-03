package ccip

import (
	"fmt"
	"math/big"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"

	"github.com/smartcontractkit/chainlink-ton/pkg/bindings"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/bindings/onramp"
	"github.com/smartcontractkit/chainlink-ton/pkg/ccip/codec"
	"github.com/smartcontractkit/chainlink-ton/pkg/ton/tracetracking"
	"github.com/smartcontractkit/chainlink-ton/pkg/ton/wrappers"

	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIPMessaging_TON2EVM(t *testing.T) {
	t.Skip("Currently skipping TON2EVM, Debugging EVM2TON")
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
	//t.Skip("Test stalls because TON test assertions aren't implemented yet")
	// Setup 2 chains (EVM and Ton) and a single lane.
	// ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithTonChains(1),
	)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	evmChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	slices.Sort(evmChainSelectors)
	allTonChainSelectors := maps.Keys(e.Env.BlockChains.TonChains())
	sourceChain := evmChainSelectors[0]
	destChain := allTonChainSelectors[0]

	t.Logf("=== Test Configuration ===")
	t.Logf("  Source (EVM): %d", sourceChain)
	t.Logf("  Dest (TON):   %d", destChain)
	t.Logf("  OnRamp:       %s", state.Chains[sourceChain].OnRamp.Address())

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

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

	t.Logf("Waiting for event filter registration (~2 mins)...")
	testhelpers.WaitForEventFilterRegistrationOnLane(t, state, e.Env.Offchain, sourceChain, destChain)

	t.Run("message to contract receiver", func(t *testing.T) {
		tonChain := e.Env.BlockChains.TonChains()[destChain]
		offRampAddr := state.TonChains[destChain].OffRamp

		receiver, err := deployReceiverContract(tonChain, &offRampAddr)
		require.NoError(t, err)

		t.Logf("  OffRamp:  %s", offRampAddr.String())
		t.Logf("  Receiver: %s", receiver.String())

		// TODO: should receiver address be saved in state?
		ccipChainState := state.TonChains[destChain]
		ccipChainState.ReceiverAddress = *receiver
		state.TonChains[destChain] = ccipChainState

		ac := codec.NewAddressCodec()
		receiverBytes, err := ac.AddressStringToBytes(receiver.String())
		require.NoError(t, err)
		require.Equal(t, 36, len(receiverBytes), "receiver bytes should be 36 bytes")

		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType: mt.ValidationTypeExec,
				TestSetup:      setup,
				Nonce:          nil, // TON nonce check is skipped
				Receiver:       receiverBytes,
				MsgData:        []byte{}, // TODO: empty data fails?
				// MsgData:                []byte("hello CCIPReceiver"), // TODO: empty data fails?
				ExtraArgs:              testhelpers.MakeEVMExtraArgsV2(100000, false),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
		// TODO: need a test case with wallet receiver(no reply nor received events)
	})

	_ = out
}

// TODO: do we want to have a changeset for receiver? probably for staging validation
func deployReceiverContract(tonChain ton.Chain, offRampAddr *address.Address) (*address.Address, error) {
	// parse compiled contract
	codeCell, err := wrappers.ParseCompiledContract(bindings.GetBuildDir("examples.receiver.compiled.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse Receiver compiled contract: %w", err)
	}

	// Create initial storage - must match TypeScript: beginCell().storeAddress(offRampAddress).endCell()
	// Note: Unlike other contracts that use tlb.ToCell(storage) which creates empty root + ref structure,
	// receiver.tolk expects a simple cell with address stored directly in root cell.
	receiverStorage := cell.BeginCell().
		MustStoreAddr(offRampAddr).
		EndCell()

	conn := tracetracking.NewSignedAPIClient(tonChain.Client, *tonChain.Wallet)
	contract, _, err := wrappers.Deploy(
		&conn,
		codeCell,
		receiverStorage,
		tlb.MustFromTON("5"), // TODO: Configurable
		cell.BeginCell().EndCell(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Receiver contract: %w", err)
	}
	receiver := contract.Address

	return receiver, nil
}
