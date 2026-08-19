package ccip

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	suiBind "github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_dummy_receiver "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_dummy_receiver/ccip_dummy_receiver"
	codec "github.com/smartcontractkit/chainlink-sui/codec"
	sui_cs "github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

// solana2SuiMessagingFixtures is shared setup for Solana->Sui messaging tests.
type solana2SuiMessagingFixtures struct {
	e                      testhelpers.DeployedEnv
	sourceChain, destChain uint64
	state                  stateview.CCIPOnChainState
	setup                  messagingtest.TestSetup
	receiverByte           []byte
	receiverObjectIDs      [][32]byte
	receiverPkgID          string
	receiverStateObjID     string
}

// prepareSolana2SuiMessagingTest brings up a Solana+Sui env, wires the Solana->Sui lane, and
// deploys + registers a Sui dummy receiver. It mirrors prepareEVM2SuiMessagingTest with a
// Solana source instead of an EVM source.
func prepareSolana2SuiMessagingTest(t *testing.T) solana2SuiMessagingFixtures {
	t.Helper()
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithSolChains(1),
		testhelpers.WithSuiChains(1),
	)

	solChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySolana))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := solChainSelectors[0]
	destChain := suiChainSelectors[0]

	t.Log("Source chain (Solana): ", sourceChain, "Dest chain (Sui): ", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	var setup messagingtest.TestSetup

	_, output, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.DeployDummyReceiver{}, sui_cs.DeployDummyReceiverConfig{
			SuiChainSelector: destChain,
			McmsOwner:        "0x1",
		}),
	})
	require.NoError(t, err)

	rawOutput := output[0].Reports[0]

	outputMap, ok := rawOutput.Output.(sui_ops.OpTxResult[ccipops.DeployDummyReceiverObjects])
	require.True(t, ok)

	id := strings.TrimPrefix(outputMap.PackageId, "0x")
	receiverByteDecoded, err := hex.DecodeString(id)
	require.NoError(t, err)

	updatedEnv, _, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.RegisterDummyReceiver{}, sui_cs.RegisterDummyReceiverConfig{
			SuiChainSelector:       destChain,
			OwnerCapObjectId:       outputMap.Objects.OwnerCapObjectId,
			CCIPObjectRefObjectId:  state.SuiChains[destChain].CCIPObjectRef,
			DummyReceiverPackageId: outputMap.PackageId,
		}),
	})
	require.NoError(t, err)
	e.Env = updatedEnv

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	e.RefreshAdapters()

	sender := common.LeftPadBytes(e.Env.BlockChains.SolanaChains()[sourceChain].DeployerKey.PublicKey().Bytes(), 32)
	setup = messagingtest.NewTestSetupWithDeployedEnv(
		t,
		e,
		state,
		sourceChain,
		destChain,
		sender,
		false, // test router
	)

	var clockObj [32]byte
	copy(clockObj[:], hexutil.MustDecode(
		"0x0000000000000000000000000000000000000000000000000000000000000006",
	))

	var stateObj [32]byte
	copy(stateObj[:], hexutil.MustDecode(
		outputMap.Objects.CCIPReceiverStateObjectId,
	))

	receiverObjectIDs := [][32]byte{clockObj, stateObj}

	return solana2SuiMessagingFixtures{
		e: e, sourceChain: sourceChain, destChain: destChain,
		state: state, setup: setup,
		receiverByte: receiverByteDecoded, receiverObjectIDs: receiverObjectIDs,
		receiverPkgID: outputMap.PackageId, receiverStateObjID: outputMap.Objects.CCIPReceiverStateObjectId,
	}
}

// Test_CCIP_Messaging_Solana2Sui_Success sends a message-only CCIP message from Solana to Sui
// using the SuiExtraArgsV1 extra args introduced by chainlink-ccip PR #2239, then asserts the
// Sui dummy receiver ran ccip_receive with no token transfer. This is the on-chain counterpart
// to the relayer-only Solana->Sui workaround described in
// core/capabilities/ccip/ccipsui/SOLANA_TO_SUI.md.
//
// Skipped until the Sui side of the Solana<->Sui lane is wired by the test helpers. AddLane
// wires the Solana source for a Sui remote via the FamilySui case in
// AddLaneSolanaChangesetsV0_1_0, but the Sui OffRamp must also be configured to accept the
// Solana source chain. For EVM<->Sui that is done by the chainlink-sui ConnectSuiToEVM
// changeset at env bring-up; no ConnectSuiToSolana analog exists yet. Remove this skip once
// that Sui offramp source-chain wiring lands.
func Test_CCIP_Messaging_Solana2Sui_Success(t *testing.T) {
	t.Skip("Solana->Sui lane needs Sui OffRamp source-chain config wiring (ConnectSuiToSolana analog); see ccipsui/SOLANA_TO_SUI.md")

	fx := prepareSolana2SuiMessagingTest(t)
	var nonce uint64

	waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.destChain])

	t.Run("Message to Sui", func(t *testing.T) {
		testhelpers.WaitForEventFilterRegistrationOnLane(t, fx.state, fx.e.Env.Offchain, fx.sourceChain, fx.destChain)

		message := []byte("Hello Sui, from Solana!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              fx.setup,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				Receiver:               fx.receiverByte,
				MsgData:                message,
				FeeToken:               "", // native SOL, converted to wSOL via Sync Native
				ExtraArgs:              testhelpers.MakeSolanaSuiExtraArgsV1(1_000_000, true, fx.receiverObjectIDs, [32]byte{}),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})

	waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.destChain])

	// Message-only: ccip_receive ran but carries no token transfer, so the receiver must
	// have stored zero dest token amounts.
	ctx := testcontext.Get(t)
	suiChain := fx.e.Env.BlockChains.SuiChains()[fx.destChain]
	receiverContract, err := module_dummy_receiver.NewDummyReceiver(fx.receiverPkgID, suiChain.Client)
	require.NoError(t, err)
	receiverStateObj := codec.Object{Id: fx.receiverStateObjID}
	devInspectOpts := &suiBind.CallOpts{
		Signer:           suiChain.Signer,
		WaitForExecution: true,
	}
	counter, err := receiverContract.DevInspect().GetCounter(ctx, devInspectOpts, receiverStateObj)
	require.NoError(t, err)
	require.Positive(t, counter, "dummy receiver ccip_receive did not run for the message-only case")
	destTokenAmounts, err := receiverContract.DevInspect().GetDestTokenAmounts(ctx, devInspectOpts, receiverStateObj)
	require.NoError(t, err)
	require.Empty(t, destTokenAmounts, "message-only path must store no dest token amounts")
}

// Test_CCIP_Messaging_Sui2Solana_Success is the reverse-lane regression for Solana<->Sui.
// Sui->Solana already worked before PR #2239; this guards against regressions from the shared
// ccipsui parseExtraDataMap / codec changes.
//
// Skipped: the Sui source side of the lane is not wired by AddLane, whose fromFamily switch
// has no FamilySui case, so no Sui OnRamp dest config or fee-quoter pricing for Solana is
// applied. Sui-source SVM-dest extra args also need a Sui/BCS builder that does not exist yet.
// Add this once the Sui->Solana lane helpers and Sui-source SVM extra-args builder land.
func Test_CCIP_Messaging_Sui2Solana_Success(t *testing.T) {
	t.Skip("Sui->Solana lane needs AddLane FamilySui source case + Sui-source SVM extra-args builder")
}
