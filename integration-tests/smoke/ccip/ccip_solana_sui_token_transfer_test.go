package ccip

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	soltestutils "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/testutils"
	ccip_router "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_0/ccip_router"
	solstate "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	soltokens "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"

	testcontext "github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

// solana2SuiTokenFixtures is shared setup for the Solana->Sui burn-mint token-transfer tests.
type solana2SuiTokenFixtures struct {
	e            testhelpers.DeployedEnv
	sourceChain  uint64 // Solana
	destChain    uint64 // Sui
	state        stateview.CCIPOnChainState
	solTokenMint solana.PublicKey
	suiLinkPkgID string
	suiAddr      [32]byte // Sui deployer address (the token_receiver) as 32 bytes
	suiAddrStr   string   // Sui deployer address as a 0x-hex string (for coin-balance queries)
	wSOL         solana.PublicKey
}

// prepareSolana2SuiTokenTransferTest brings up a 1-Solana + 1-Sui env, wires the Solana->Sui lane,
// deploys + cross-registers a burn-mint LINK pool pair, and sets up the Solana source fee (wSOL)
// and SPL-2022 token ATA + CCIP spend approvals. It mirrors the Solana source fee/ATA setup in
// TestTokenTransfer_Solana2EVM (ccip_token_transfer_test.go) and the Sui dest token-receiver
// derivation in testSetupHelperEvm2Sui (ccip_sui_token_transfer_test.go).
func prepareSolana2SuiTokenTransferTest(t *testing.T) solana2SuiTokenFixtures {
	t.Helper()
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithSolChains(1),
		testhelpers.WithSuiChains(1),
	)

	solChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chainsel.FamilySolana))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chainsel.FamilySui))
	sourceChain := solChainSelectors[0]
	destChain := suiChainSelectors[0]

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Log("Source chain (Solana): ", sourceChain, "Dest chain (Sui): ", destChain)

	require.NoError(t, testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false))

	// Deploy + cross-register the Solana<->Sui burn-mint LINK pool pair (helper is hybrid:
	// Sui DeployTPAndConfigure + Solana E2ETokenPool + one v1_6_0 UpsertRemoteChainConfigBurnMint
	// op for the Solana->Sui remote, since the legacy Solana cross-reg is EVM-only).
	solTokenMint, _, suiLinkPkgID, err := testhelpers.DeployAndCrossRegisterBurnMintPoolSolanaSui(t, &e, sourceChain, destChain)
	require.NoError(t, err)

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// --- Solana source fee (wSOL) + token ATA + CCIP spend approvals -----------------
	// Mirrors TestTokenTransfer_Solana2EVM lines 392-444.
	solChains := e.Env.BlockChains.SolanaChains()
	deployer := solChains[sourceChain].DeployerKey
	rpcClient := solChains[sourceChain].Client
	ctx := testcontext.Get(t)

	wSOL := solana.SolMint
	ixAtaUser, deployerWSOL, err := soltokens.CreateAssociatedTokenAccount(solana.TokenProgramID, wSOL, deployer.PublicKey(), deployer.PublicKey())
	require.NoError(t, err)

	billingSignerPDA, _, err := solstate.FindFeeBillingSignerPDA(state.SolChains[sourceChain].Router)
	require.NoError(t, err)

	// Approve CCIP to spend wSOL for fee billing.
	ixApprove, err := soltokens.TokenApproveChecked(1e9, 9, solana.TokenProgramID, deployerWSOL, wSOL, billingSignerPDA, deployer.PublicKey(), []solana.PublicKey{})
	require.NoError(t, err)
	soltestutils.SendAndConfirm(ctx, t, rpcClient, []solana.Instruction{ixAtaUser, ixApprove}, *deployer, solconfig.DefaultCommitment)

	// Fund the wSOL ATA (transfer SOL + syncNative).
	transferAmount := 1.0 * solana.LAMPORTS_PER_SOL
	ixTransfer, err := soltokens.NativeTransfer(transferAmount, deployer.PublicKey(), deployerWSOL)
	require.NoError(t, err)
	ixSync, err := soltokens.SyncNative(solana.TokenProgramID, deployerWSOL)
	require.NoError(t, err)
	soltestutils.SendAndConfirm(ctx, t, rpcClient, []solana.Instruction{ixTransfer, ixSync}, *deployer, solconfig.DefaultCommitment)

	// SPL-2022 token ATA + approve CCIP to take custody of the token for the burn-mint pool.
	ownerSource := deployer.PublicKey()
	userTokenAccount, _, err := soltokens.FindAssociatedTokenAddress(solana.Token2022ProgramID, solTokenMint, ownerSource)
	require.NoError(t, err)
	ixApproveToken, err := soltokens.TokenApproveChecked(1000, 9, solana.Token2022ProgramID, userTokenAccount, solTokenMint, billingSignerPDA, ownerSource, nil)
	require.NoError(t, err)
	soltestutils.SendAndConfirm(ctx, t, rpcClient, []solana.Instruction{ixApproveToken}, *deployer, solconfig.DefaultCommitment)

	// Sui receiver wallet (token_receiver) = Sui deployer address as 32 bytes. Mirrors
	// testSetupHelperEvm2Sui's suiAddr derivation.
	suiAddrStr, err := e.Env.BlockChains.SuiChains()[destChain].Signer.GetAddress()
	require.NoError(t, err)
	addrBytes, err := hex.DecodeString(strings.TrimPrefix(suiAddrStr, "0x"))
	require.NoError(t, err)
	require.Len(t, addrBytes, 32, "expected 32-byte sui address")
	var suiAddr [32]byte
	copy(suiAddr[:], addrBytes)

	return solana2SuiTokenFixtures{
		e:            e,
		sourceChain:  sourceChain,
		destChain:    destChain,
		state:        state,
		solTokenMint: solTokenMint,
		suiLinkPkgID: suiLinkPkgID,
		suiAddr:      suiAddr,
		suiAddrStr:   suiAddrStr,
		wSOL:         wSOL,
	}
}

// suiLinkBalance returns the total Sui LINK coin balance held by account, keyed on the full coin
// type 0x2::coin::Coin<<linkPkgID>::link::LINK> (mirrors WaitForTokenBalanceSui's query). Used for
// a before/after delta assertion since the exact minted amount depends on source/dest decimals.
func suiLinkBalance(t *testing.T, ctx context.Context, chain cldf_sui.Chain, account, linkPkgID string) *big.Int {
	t.Helper()
	coins, err := chain.Client.QueryCoinsByAddress(ctx, account, "0x2::coin::Coin<"+linkPkgID+"::link::LINK>")
	require.NoError(t, err)
	balance := new(big.Int)
	for _, coin := range coins {
		balance.Add(balance, new(big.Int).SetUint64(coin.GetBalance()))
	}
	return balance
}

// Test_CCIPTokenTransfer_Solana2Sui_BurnMintTokenPool sends a single burn-mint token from Solana
// to Sui with message.receiver = 0 (pure token transfer, no ccip_receive execution) and
// token_receiver = the Sui deployer wallet, then asserts the Sui wallet's LINK balance increased.
//
// This is the token-transfer counterpart to Test_CCIP_Messaging_Solana2Sui_Success and exercises
// the SuiExtraArgsV1 (tag 0x21ea4ca9) Solana-source path end-to-end through the 1.6.4 fee-quoter
// Sui dest dispatch + the Sui OffRamp burn-mint minting.
//
// message.receiver = 0 is permitted by the fee-quoter when gas_limit == 0 (its 2-arg
// validate_sui_address short-circuits), and the Sui OffRamp skips an unregistered/zero receiver
// while still delivering the minted token to token_receiver. Because TransferMultiple's
// Solana-source branch books ExpectedTokenBalances against tt.Receiver (which is zero here), the
// balance is asserted manually via suiLinkBalance rather than via ExpectedTokenBalances.
func Test_CCIPTokenTransfer_Solana2Sui_BurnMintTokenPool(t *testing.T) {
	t.Parallel()
	fx := prepareSolana2SuiTokenTransferTest(t)
	ctx := testcontext.Get(t)
	e := fx.e.Env
	suiChain := e.BlockChains.SuiChains()[fx.destChain]

	waitForSuiRPCSync(t, suiChain)
	testhelpers.WaitForEventFilterRegistrationOnLane(t, fx.state, e.Offchain, fx.sourceChain, fx.destChain)

	balanceBefore := suiLinkBalance(t, ctx, suiChain, fx.suiAddrStr, fx.suiLinkPkgID)

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:        "Send token to Sui EOA",
			SourceChain: fx.sourceChain,
			DestChain:   fx.destChain,
			// message.receiver = 0: pure token transfer, no ccip_receive. The minted token is
			// delivered to token_receiver (the Sui wallet) set in ExtraArgs below.
			Receiver: make([]byte, 32),
			FeeToken: fx.wSOL.String(),
			SolTokens: []ccip_router.SVMTokenAmount{
				{Token: fx.solTokenMint, Amount: 1},
			},
			// gas_limit = 0 (no receiver execution); token_receiver = Sui wallet (non-zero, required
			// when tokens are present). receiverObjectIDs = nil (no ccip_receive, no receiver objects).
			ExtraArgs:      testhelpers.MakeSolanaSuiExtraArgsV1(0, true, nil, fx.suiAddr),
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		},
	}

	startBlocks, expectedSeqNums, expectedExecutionStates, _ := testhelpers.TransferMultiple(ctx, t, e, fx.state, tcs)

	require.NoError(t, testhelpers.ConfirmMultipleCommits(t, e, fx.state, startBlocks, false, expectedSeqNums))

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t, e, fx.state, testhelpers.SeqNumberRangeToSlice(expectedSeqNums), startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	waitForSuiRPCSync(t, suiChain)

	// Assert the Sui wallet received the minted LINK (balance strictly increased). Tolerant of
	// decimals/amount exactness per the user's "ignore fee nits" constraint.
	require.Eventually(t, func() bool {
		return suiLinkBalance(t, ctx, suiChain, fx.suiAddrStr, fx.suiLinkPkgID).Cmp(balanceBefore) > 0
	}, 10*time.Minute, 2*time.Second, "Sui wallet LINK balance did not increase after Solana->Sui token transfer")
}

// Test_CCIPTokenTransfer_Solana2Sui_BurnMintTokenPool_ZeroTokenReceiver_Revert asserts that
// ccip_send rejects a Solana->Sui token transfer whose SuiExtraArgsV1 token_receiver is zero while
// a token is present. This is the novel 1.6.4 fee-quoter Sui gating (InvalidTokenReceiver,
// messages.rs), and it fires at send time so no relayer round-trip is needed.
//
// message.receiver is set to a valid non-zero Sui address with gas_limit > 0 so the fee-quoter's
// receiver validation passes and the only failing check is the token_receiver one, isolating it.
func Test_CCIPTokenTransfer_Solana2Sui_BurnMintTokenPool_ZeroTokenReceiver_Revert(t *testing.T) {
	t.Parallel()
	fx := prepareSolana2SuiTokenTransferTest(t)
	e := fx.e.Env

	waitForSuiRPCSync(t, e.BlockChains.SuiChains()[fx.destChain])
	testhelpers.WaitForEventFilterRegistrationOnLane(t, fx.state, e.Offchain, fx.sourceChain, fx.destChain)

	msg := ccip_router.SVM2AnyMessage{
		Receiver:     common.LeftPadBytes(fx.suiAddr[:], 32), // valid non-zero receiver
		Data:         []byte{},
		TokenAmounts: []ccip_router.SVMTokenAmount{{Token: fx.solTokenMint, Amount: 1}},
		FeeToken:     fx.wSOL,
		// token_receiver = 0 with a token present -> fee-quoter InvalidTokenReceiver at ccip_send.
		ExtraArgs: testhelpers.MakeSolanaSuiExtraArgsV1(1_000_000, true, nil, [32]byte{}),
	}

	_, err := testhelpers.SendRequest(e, fx.state,
		ccipclient.WithSourceChain(fx.sourceChain),
		ccipclient.WithDestChain(fx.destChain),
		ccipclient.WithTestRouter(false),
		ccipclient.WithMessage(msg),
	)
	require.Error(t, err)
	t.Log("Expected ccip_send rejection (zero token_receiver with token present): ", err)
}
