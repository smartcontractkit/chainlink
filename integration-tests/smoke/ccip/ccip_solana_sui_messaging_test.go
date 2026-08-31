package ccip

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	solccip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"

	suiutil "github.com/smartcontractkit/chainlink-sui/bindings/utils"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	soltesthelpers "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

// // solana2SuiMessagingFixtures is shared setup for Solana->Sui messaging tests.
// type solana2SuiMessagingFixtures struct {
// 	e                      testhelpers.DeployedEnv
// 	sourceChain, destChain uint64
// 	state                  stateview.CCIPOnChainState
// 	setup                  messagingtest.TestSetup
// 	receiverByte           []byte
// 	receiverObjectIDs      [][32]byte
// 	receiverPkgID          string
// 	receiverStateObjID     string
// }

// // prepareSolana2SuiMessagingTest brings up a Solana+Sui env, wires the Solana->Sui lane, and
// // deploys + registers a Sui dummy receiver. It mirrors prepareEVM2SuiMessagingTest with a
// // Solana source instead of an EVM source.
// func prepareSolana2SuiMessagingTest(t *testing.T) solana2SuiMessagingFixtures {
// 	t.Helper()
// 	e, _, _ := testsetups.NewIntegrationEnvironment(
// 		t,
// 		testhelpers.WithSolChains(1),
// 		testhelpers.WithSuiChains(1),
// 	)

// 	solChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySolana))
// 	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

// 	state, err := stateview.LoadOnchainState(e.Env)
// 	require.NoError(t, err)

// 	sourceChain := solChainSelectors[0]
// 	destChain := suiChainSelectors[0]

// 	t.Log("Source chain (Solana): ", sourceChain, "Dest chain (Sui): ", destChain)

// 	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
// 	require.NoError(t, err)

// 	var setup messagingtest.TestSetup

// 	_, output, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
// 		commoncs.Configure(sui_cs.DeployDummyReceiver{}, sui_cs.DeployDummyReceiverConfig{
// 			SuiChainSelector: destChain,
// 			McmsOwner:        "0x1",
// 		}),
// 	})
// 	require.NoError(t, err)

// 	rawOutput := output[0].Reports[0]

// 	outputMap, ok := rawOutput.Output.(sui_ops.OpTxResult[ccipops.DeployDummyReceiverObjects])
// 	require.True(t, ok)

// 	id := strings.TrimPrefix(outputMap.PackageId, "0x")
// 	receiverByteDecoded, err := hex.DecodeString(id)
// 	require.NoError(t, err)

// 	updatedEnv, _, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
// 		commoncs.Configure(sui_cs.RegisterDummyReceiver{}, sui_cs.RegisterDummyReceiverConfig{
// 			SuiChainSelector:       destChain,
// 			OwnerCapObjectId:       outputMap.Objects.OwnerCapObjectId,
// 			CCIPObjectRefObjectId:  state.SuiChains[destChain].CCIPObjectRef,
// 			DummyReceiverPackageId: outputMap.PackageId,
// 		}),
// 	})
// 	require.NoError(t, err)
// 	e.Env = updatedEnv

// 	state, err = stateview.LoadOnchainState(e.Env)
// 	require.NoError(t, err)
// 	e.RefreshAdapters()

// 	sender := common.LeftPadBytes(e.Env.BlockChains.SolanaChains()[sourceChain].DeployerKey.PublicKey().Bytes(), 32)
// 	setup = messagingtest.NewTestSetupWithDeployedEnv(
// 		t,
// 		e,
// 		state,
// 		sourceChain,
// 		destChain,
// 		sender,
// 		false, // test router
// 	)

// 	var clockObj [32]byte
// 	copy(clockObj[:], hexutil.MustDecode(
// 		"0x0000000000000000000000000000000000000000000000000000000000000006",
// 	))

// 	var stateObj [32]byte
// 	copy(stateObj[:], hexutil.MustDecode(
// 		outputMap.Objects.CCIPReceiverStateObjectId,
// 	))

// 	receiverObjectIDs := [][32]byte{clockObj, stateObj}

// 	return solana2SuiMessagingFixtures{
// 		e: e, sourceChain: sourceChain, destChain: destChain,
// 		state: state, setup: setup,
// 		receiverByte: receiverByteDecoded, receiverObjectIDs: receiverObjectIDs,
// 		receiverPkgID: outputMap.PackageId, receiverStateObjID: outputMap.Objects.CCIPReceiverStateObjectId,
// 	}
// }

// // Test_CCIP_Messaging_Solana2Sui_Success sends a message-only CCIP message from Solana to Sui
// // using the SuiExtraArgsV1 extra args introduced by chainlink-ccip PR #2239, then asserts the
// // Sui dummy receiver ran ccip_receive with no token transfer. This is the on-chain counterpart
// // to the relayer-only Solana->Sui workaround described in
// // core/capabilities/ccip/ccipsui/SOLANA_TO_SUI.md.
// //
// // Skipped until the Sui side of the Solana<->Sui lane is wired by the test helpers. AddLane
// // wires the Solana source for a Sui remote via the FamilySui case in
// // AddLaneSolanaChangesetsV0_1_0, but the Sui OffRamp must also be configured to accept the
// // Solana source chain. For EVM<->Sui that is done by the chainlink-sui ConnectSuiToEVM
// // changeset at env bring-up; no ConnectSuiToSolana analog exists yet. Remove this skip once
// // that Sui offramp source-chain wiring lands.
// func Test_CCIP_Messaging_Solana2Sui_Success(t *testing.T) {
// 	t.Skip("Solana->Sui lane needs Sui OffRamp source-chain config wiring (ConnectSuiToSolana analog); see ccipsui/SOLANA_TO_SUI.md")

// 	fx := prepareSolana2SuiMessagingTest(t)
// 	var nonce uint64

// 	waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.destChain])

// 	t.Run("Message to Sui", func(t *testing.T) {
// 		testhelpers.WaitForEventFilterRegistrationOnLane(t, fx.state, fx.e.Env.Offchain, fx.sourceChain, fx.destChain)

// 		message := []byte("Hello Sui, from Solana!")
// 		messagingtest.Run(t,
// 			messagingtest.TestCase{
// 				TestSetup:              fx.setup,
// 				Nonce:                  &nonce,
// 				ValidationType:         messagingtest.ValidationTypeExec,
// 				Receiver:               fx.receiverByte,
// 				MsgData:                message,
// 				FeeToken:               "", // native SOL, converted to wSOL via Sync Native
// 				ExtraArgs:              testhelpers.MakeSolanaSuiExtraArgsV1(1_000_000, true, fx.receiverObjectIDs, [32]byte{}),
// 				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
// 			},
// 		)
// 	})

// 	waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.destChain])

// 	// Message-only: ccip_receive ran but carries no token transfer, so the receiver must
// 	// have stored zero dest token amounts.
// 	ctx := testcontext.Get(t)
// 	suiChain := fx.e.Env.BlockChains.SuiChains()[fx.destChain]
// 	receiverContract, err := module_dummy_receiver.NewDummyReceiver(fx.receiverPkgID, suiChain.Client)
// 	require.NoError(t, err)
// 	receiverStateObj := codec.Object{Id: fx.receiverStateObjID}
// 	devInspectOpts := &suiBind.CallOpts{
// 		Signer:           suiChain.Signer,
// 		WaitForExecution: true,
// 	}
// 	counter, err := receiverContract.DevInspect().GetCounter(ctx, devInspectOpts, receiverStateObj)
// 	require.NoError(t, err)
// 	require.Positive(t, counter, "dummy receiver ccip_receive did not run for the message-only case")
// 	destTokenAmounts, err := receiverContract.DevInspect().GetDestTokenAmounts(ctx, devInspectOpts, receiverStateObj)
// 	require.NoError(t, err)
// 	require.Empty(t, destTokenAmounts, "message-only path must store no dest token amounts")
// }

// // Test_CCIP_Messaging_Sui2Solana_Success is the reverse-lane regression for Solana<->Sui.
// // Sui->Solana already worked before PR #2239; this guards against regressions from the shared
// // ccipsui parseExtraDataMap / codec changes.
// //
// // Skipped: the Sui source side of the lane is not wired by AddLane, whose fromFamily switch
// // has no FamilySui case, so no Sui OnRamp dest config or fee-quoter pricing for Solana is
// // applied. Sui-source SVM-dest extra args also need a Sui/BCS builder that does not exist yet.
// // Add this once the Sui->Solana lane helpers and Sui-source SVM extra-args builder land.
// func Test_CCIP_Messaging_Sui2Solana_Success(t *testing.T) {
// 	t.Skip("Sui->Solana lane needs AddLane FamilySui source case + Sui-source SVM extra-args builder")
// }

// sui2SolanaMessagingFixtures is shared setup for Sui->Solana messaging tests.
type sui2SolanaMessagingFixtures struct {
	e                                  testhelpers.DeployedEnv
	sourceChain, destChain             uint64
	state                              stateview.CCIPOnChainState
	setup                              mt.TestSetup
	suiLinkFeeToken                    string
	nonce                              *uint64
	receiverProgram                    solana.PublicKey
	receiverTargetAccountPDA           solana.PublicKey
	receiverExternalExecutionConfigPDA solana.PublicKey
}

// prepareSui2SolanaMessagingTest brings up an EVM+Sui+Solana env, wires the Sui->Solana lane
// (via the lanes.ConnectChains Sui<->Solana case in AddLane), deploys the Solana CCIP receiver,
// mints a Sui LINK fee token, and builds the messagingtest setup. It mirrors
// prepareSui2EvmMessagingTest (Sui2EVM test) for the Sui source side and the EVM2Solana
// messaging test for the Solana dest side (receiver PDA + accounts).
func prepareSui2SolanaMessagingTest(t *testing.T) sui2SolanaMessagingFixtures {
	t.Helper()
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2), // Validate requires >=2 EVM chains
		testhelpers.WithSuiChains(1),
		testhelpers.WithSolChains(1),
		testhelpers.WithOCRConfigOverride(func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
			params.ExecuteOffChainConfig.InflightCacheExpiry = *config.MustNewDuration(1 * time.Minute)
			params.ExecuteOffChainConfig.MessageVisibilityInterval = *config.MustNewDuration(1 * time.Hour)
			params.ExecuteOffChainConfig.MultipleReportsEnabled = true
			params.ExecuteOffChainConfig.MaxReportMessages = 1
			params.ExecuteOffChainConfig.MaxSingleChainReports = 1
			return params
		}),
	)

	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chainsel.FamilySui))
	solChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chainsel.FamilySolana))
	sourceChain := suiChainSelectors[0]
	destChain := solChainSelectors[0]

	// Solana dest requires the CCIP receiver program to be deployed.
	testhelpers.DeploySolanaCcipReceiver(t, e.Env)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Log("Source chain (Sui): ", sourceChain, "Dest chain (Solana): ", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	// Sui source pays fees in LINK; mint a LINK coin object to fund ccip_send.
	feeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000000)

	// Sui sender address -> 32-byte left-pad for the Solana message sender field.
	suiSenderAddr, err := e.Env.BlockChains.SuiChains()[sourceChain].Signer.GetAddress()
	require.NoError(t, err)
	normalizedAddr, err := suiutil.ConvertStringToAddressBytes(suiSenderAddr)
	require.NoError(t, err)
	sender := common.LeftPadBytes(normalizedAddr[:], 32)

	setup := mt.NewTestSetupWithDeployedEnv(
		t,
		e,
		state,
		sourceChain,
		destChain,
		sender,
		false, // testRouter
	)

	receiverProgram := state.SolChains[destChain].Receiver
	receiverTargetAccountPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("counter")}, receiverProgram)
	receiverExternalExecutionConfigPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("external_execution_config")}, receiverProgram)

	n := uint64(0)
	return sui2SolanaMessagingFixtures{
		e:                                  e,
		sourceChain:                        sourceChain,
		destChain:                          destChain,
		state:                              state,
		setup:                              setup,
		suiLinkFeeToken:                    feeTokenOutput.Objects.MintedLinkTokenObjectId,
		nonce:                              &n,
		receiverProgram:                    receiverProgram,
		receiverTargetAccountPDA:           receiverTargetAccountPDA,
		receiverExternalExecutionConfigPDA: receiverExternalExecutionConfigPDA,
	}
}

// Test_CCIP_Messaging_Sui2Solana_Success sends a message-only CCIP message from Sui to Solana
// using SVMExtraArgsV1 extra args encoded in BCS (the Sui source family's convention) via
// MakeSuiSourceSVMExtraArgsV1, then asserts the Solana CCIPReceiver ran ccip_receive and
// incremented its counter. The Sui->Solana off-chain codec is already complete: the Solana dest
// MessageHasher decodes the source extraArgs via the source-family-dispatching codec bundle,
// which routes FamilySui to the ccipaptos BCS decoder that decodes SVMExtraArgsV1.
func Test_CCIP_Messaging_Sui2Solana_Success(t *testing.T) {
	fx := prepareSui2SolanaMessagingTest(t)
	ctx := testhelpers.Context(t)
	solChains := fx.e.Env.BlockChains.SolanaChains()

	waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.sourceChain])

	t.Run("Message to Solana CCIPReceiver", func(t *testing.T) {
		testhelpers.WaitForEventFilterRegistrationOnLane(t, fx.state, fx.e.Env.Offchain, fx.sourceChain, fx.destChain)

		accounts := [][32]byte{
			fx.receiverExternalExecutionConfigPDA,
			fx.receiverTargetAccountPDA,
			solana.SystemProgramID,
		}

		// Sui source emits SVMExtraArgsV1 in BCS (tag 0x1f3b3aba + BCS fields) for a Solana dest.
		extraArgs := testhelpers.MakeSuiSourceSVMExtraArgsV1(
			80_000, // computeUnits
			solccip.GenerateBitMapForIndexes([]int{0, 1}), // accountIsWritableBitmap
			true,       // allowOutOfOrderExecution
			[32]byte{}, // tokenReceiver: zero (message-only)
			accounts,
		)

		// counter starts at 0
		var receiverCounterAccount soltesthelpers.ReceiverCounter
		err := solcommon.GetAccountDataBorshInto(ctx, solChains[fx.destChain].Client, fx.receiverTargetAccountPDA, solconfig.DefaultCommitment, &receiverCounterAccount)
		require.NoError(t, err, "failed to get receiver counter account before send")
		require.Equal(t, uint8(0), receiverCounterAccount.Value)

		message := []byte("Hello Solana, from Sui!")
		mt.Run(t,
			mt.TestCase{
				TestSetup:              fx.setup,
				Nonce:                  fx.nonce,
				ValidationType:         mt.ValidationTypeExec,
				Receiver:               fx.receiverProgram.Bytes(),
				MsgData:                message,
				FeeToken:               fx.suiLinkFeeToken,
				ExtraArgs:              extraArgs,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
						var after soltesthelpers.ReceiverCounter
						err := solcommon.GetAccountDataBorshInto(ctx, solChains[fx.destChain].Client, fx.receiverTargetAccountPDA, solconfig.DefaultCommitment, &after)
						require.NoError(t, err, "failed to get receiver counter account after exec")
						require.Equal(t, uint8(1), after.Value, "Solana CCIP receiver counter should increment after ccip_receive")
					},
				},
			},
		)
	})

	waitForSuiRPCSync(t, fx.e.Env.BlockChains.SuiChains()[fx.sourceChain])
}

// ----------------------------------------------------------------------------
// DEFERRED: Sui->Solana burn-mint token transfer smoke test
// ----------------------------------------------------------------------------
//
// Test_CCIPTokenTransfer_Sui2Solana_BurnMintTokenPool is intentionally NOT
// implemented yet. Unlike the messaging test above (which reuses the existing
// family-agnostic lanes.ConnectChains wiring + a BCS SVMExtraArgsV1 builder),
// the token-transfer case has no reusable testhelper and hits a real gap:
//
//   1. No combined Sui<->Solana burn-mint pool helper exists.
//        - HandleTokenAndBurnMintTokenPoolDeploymentForSUI (test_sui_helpers.go:734)
//          is EVM-dest-hardcoded (deploys an EVM burn-mint pool, cross-regs with
//          EVM 20-byte addresses).
//        - DeployTransferableTokenSolanaV0_1_1 (test_helpers_solana_v0_1_1.go:115)
//          requires an EVM leg (explicit FamilyEVM check).
//
//   2. The legacy Solana cross-registration changeset is a BLOCKER for a Sui
//      remote: SetupTokenPoolForRemoteChain (cs_token_pool.go:1174) reads the
//      remote pool/token from EVM on-chain state via EVMChains()[selector] and
//      hard-rejects non-EVM remotes at cs_token_pool.go:480-483.
//
// The USABLE path (mapped, not yet wired) is the family-agnostic lanes
// token-adapter framework:
//   - Sui source:  SuiTokenAdapter.DeployTokenPoolForToken +
//                  ConfigureTokenForTransfersSequence
//     (chainlink-sui/deployment/adapters/token_adapter.go:637, :370). The Sui
//     ApplyChainUpdates op accepts a 32-byte Solana remote: pool via
//     StrToBytes (variable length), token via StrTo32 (left-pad to 32, rejects
//     only >32) -- op_burn_mint_token_pool.go:100, :111. PASS HEX, not base58.
//   - Solana dest: SolanaAdapter.DeployTokenPoolForToken +
//                  ConfigureTokenForTransfersSequence
//     (chainlink-ccip/chains/solana/deployment/v1_6_0/sequences/{adapter.go:806,
//      tokens.go:31}). The cross-register op UpsertRemoteChainConfigBurnMint
//     (burnmint.go:186) takes RemoteTokenAddress/RemotePoolAddress as raw []byte
//     with NO family check; the on-chain RemoteAddress.Address is [32]byte, so a
//     32-byte Sui object ID fits. ConfigureTokenForTransfersSequence also folds
//     in RegisterTokenAdminRegistry + Accept + SetPool + rate limits.
//   - Lane:        addSuiSolanaMixedLane (test_helpers_sui_solana_lanes.go) --
//                  already wired by the messaging test.
//   - Send:        TransferMultiple with a TestTransferRequest{SourceChain: sui,
//                  DestChain: sol, SuiTokens: [...], Receiver: <sol ATA 32B>,
//                  FeeToken: ...} dispatches to SendRequestSui
//                  (test_helpers_solana_v0_1_0.go:1998, :538).
//
// Two caveats to handle when implementing (from the path mapping):
//   - Transfer's Sui branch does NOT populate SuiSendRequest.TokenReceiverATA
//     (test_helpers_solana_v0_1_0.go:1887-1894); verify SendRequestSui needs it
//     for token messages to Solana.
//   - TransferMultiple's Sui-source branch books expected balances against
//     tt.Receiver (:2000), not tt.TokenReceiverATA; set tt.Receiver = the Solana
//     ATA bytes (32) so WaitForTokenBalances checks the right account.
//
// This was deferred because the chainlink smoke suite does not yet use the lanes
// token-adapter framework anywhere (no reference test for a non-EVM-Sui token
// pair), and Sui+Solana in-memory containers cannot be run in this sandbox to
// iterate. Add this test once a lanes-token reference exists or nix/CI iteration
// is available. When added, also re-add the in-memory-tests.json entry:
//   {"name":"Test_CCIPTokenTransfer_Sui2Solana_BurnMintTokenPool",
//    "test":"Test_CCIPTokenTransfer_Sui2Solana_BurnMintTokenPool",
//    "timeout":"25m","parallel":1,"plugins":true,"runs_on":"cpu=16/ram=64",
//    "free_disk":true,"aptos":"","sui":"mainnet-v1.75.2"}
