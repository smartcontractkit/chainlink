package ccip

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/maps"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	soltestutils "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/testutils"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solstate "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	soltokens "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/message_hasher"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// duplicated from messagingtest
func sleepAndReplay(t *testing.T, e testhelpers.DeployedEnv, chainSelectors ...uint64) {
	time.Sleep(30 * time.Second)
	replayBlocks := make(map[uint64]uint64)
	for _, selector := range chainSelectors {
		replayBlocks[selector] = 1
	}
	testhelpers.ReplayLogs(t, e.Env.Offchain, replayBlocks)
}

func TestTokenTransfer_EVM2EVM(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := tests.Context(t)

	tenv, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfUsersPerChain(3))

	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(e.Chains), 2)

	allChainSelectors := maps.Keys(e.Chains)
	sourceChain, destChain := allChainSelectors[0], allChainSelectors[1]
	ownerSourceChain := e.Chains[sourceChain].DeployerKey
	ownerDestChain := e.Chains[destChain].DeployerKey

	require.GreaterOrEqual(t, len(tenv.Users[sourceChain]), 2)
	require.GreaterOrEqual(t, len(tenv.Users[destChain]), 2)
	selfServeSrcTokenPoolDeployer := tenv.Users[sourceChain][1]
	selfServeDestTokenPoolDeployer := tenv.Users[destChain][1]

	oneE18 := new(big.Int).SetUint64(1e18)

	// Deploy tokens and pool by CCIP Owner
	srcToken, _, destToken, _, err := testhelpers.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		sourceChain,
		destChain,
		ownerSourceChain,
		ownerDestChain,
		state,
		e.ExistingAddresses,
		"OWNER_TOKEN",
	)
	require.NoError(t, err)

	// Deploy Self Serve tokens and pool
	selfServeSrcToken, _, selfServeDestToken, _, err := testhelpers.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		sourceChain,
		destChain,
		selfServeSrcTokenPoolDeployer,
		selfServeDestTokenPoolDeployer,
		state,
		e.ExistingAddresses,
		"SELF_SERVE_TOKEN",
	)
	require.NoError(t, err)
	testhelpers.AddLanesForAll(t, &tenv, state)

	testhelpers.MintAndAllow(
		t,
		e,
		state,
		map[uint64][]testhelpers.MintTokenInfo{
			sourceChain: {
				testhelpers.NewMintTokenInfo(selfServeSrcTokenPoolDeployer, selfServeSrcToken),
				testhelpers.NewMintTokenInfo(ownerSourceChain, srcToken),
			},
			destChain: {
				testhelpers.NewMintTokenInfo(selfServeDestTokenPoolDeployer, selfServeDestToken),
				testhelpers.NewMintTokenInfo(ownerDestChain, destToken),
			},
		},
	)

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:        "Send token to EOA",
			SourceChain: sourceChain,
			DestChain:   destChain,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken.Address(),
					Amount: oneE18,
				},
			},
			Receiver: utils.RandomAddress().Bytes(),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: destToken.Address().Bytes(), Amount: oneE18},
			},
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		},
		{
			Name:        "Send token to contract",
			SourceChain: sourceChain,
			DestChain:   destChain,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken.Address(),
					Amount: oneE18,
				},
			},
			Receiver: state.Chains[destChain].Receiver.Address().Bytes(),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: destToken.Address().Bytes(), Amount: oneE18},
			},
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		},
		{
			Name:        "Send N tokens to contract",
			SourceChain: destChain,
			DestChain:   sourceChain,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  selfServeDestToken.Address(),
					Amount: oneE18,
				},
				{
					Token:  destToken.Address(),
					Amount: oneE18,
				},
				{
					Token:  selfServeDestToken.Address(),
					Amount: oneE18,
				},
			},
			Receiver:  state.Chains[sourceChain].Receiver.Address().Bytes(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(300_000, false),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: selfServeSrcToken.Address().Bytes(), Amount: new(big.Int).Add(oneE18, oneE18)},
				{Token: srcToken.Address().Bytes(), Amount: oneE18},
			},
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		},
		{
			Name:        "Sending token transfer with custom gasLimits to the EOA is successful",
			SourceChain: destChain,
			DestChain:   sourceChain,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  selfServeDestToken.Address(),
					Amount: oneE18,
				},
				{
					Token:  destToken.Address(),
					Amount: new(big.Int).Add(oneE18, oneE18),
				},
			},
			Receiver:  utils.RandomAddress().Bytes(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(1, false),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: selfServeSrcToken.Address().Bytes(), Amount: oneE18},
				{Token: srcToken.Address().Bytes(), Amount: new(big.Int).Add(oneE18, oneE18)},
			},
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		},
		{
			Name:        "Sending PTT with too low gas limit leads to the revert when receiver is a contract",
			SourceChain: destChain,
			DestChain:   sourceChain,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  selfServeDestToken.Address(),
					Amount: oneE18,
				},
				{
					Token:  destToken.Address(),
					Amount: oneE18,
				},
			},
			Receiver:  state.Chains[sourceChain].Receiver.Address().Bytes(),
			Data:      []byte("this should be reverted because gasLimit is too low, no tokens are transferred as well"),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(1, false),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: selfServeSrcToken.Address().Bytes(), Amount: big.NewInt(0)},
				{Token: srcToken.Address().Bytes(), Amount: big.NewInt(0)},
			},
			ExpectedStatus: testhelpers.EXECUTION_STATE_FAILURE,
		},
	}

	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances :=
		testhelpers.TransferMultiple(ctx, t, e, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		e,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		e,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, e, expectedTokenBalances)
}

func TestTokenTransfer_EVM2Solana(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	tenv, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfUsersPerChain(3),
		testhelpers.WithSolChains(1))

	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(e.Chains), 2)

	allChainSelectors := maps.Keys(e.Chains)
	allSolChainSelectors := maps.Keys(e.SolChains)
	sourceChain, destChain := allChainSelectors[0], allSolChainSelectors[0]
	ownerSourceChain := e.Chains[sourceChain].DeployerKey
	// ownerDestChain := e.SolChains[destChain].DeployerKey

	require.GreaterOrEqual(t, len(tenv.Users[sourceChain]), 2) // TODO: ???

	oneE9 := new(big.Int).SetUint64(1e9)

	// Deploy tokens and pool by CCIP Owner
	srcToken, _, destToken, err := testhelpers.DeployTransferableTokenSolana(
		lggr,
		e,
		sourceChain,
		destChain,
		ownerSourceChain,
		"OWNER_TOKEN",
	)
	require.NoError(t, err)

	// testhelpers.AddLanesForAll(t, &tenv, state) TODO:, fixed for Solana now
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &tenv, state, sourceChain, destChain, false)

	testhelpers.MintAndAllow(
		t,
		e,
		state,
		map[uint64][]testhelpers.MintTokenInfo{
			sourceChain: {
				testhelpers.NewMintTokenInfo(ownerSourceChain, srcToken),
			},
			// destChain: {
			// 	testhelpers.NewMintTokenInfo(ownerDestChain, destToken),
			// },
		},
	)
	// TODO: how to do MintAndAllow on Solana?
	tokenReceiver, _, ferr := soltokens.FindAssociatedTokenAddress(solana.Token2022ProgramID, destToken, state.SolChains[destChain].Receiver)
	require.NoError(t, ferr)

	extraArgs, err := ccipevm.SerializeClientSVMExtraArgsV1(message_hasher.ClientSVMExtraArgsV1{
		TokenReceiver: tokenReceiver,
		// Accounts: accounts,
	})
	require.NoError(t, err)

	// TODO: test both with ATA pre-initialized and not

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:        "Send token to contract",
			SourceChain: sourceChain,
			DestChain:   destChain,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  srcToken.Address(),
					Amount: oneE9,
				},
			},
			TokenReceiver: tokenReceiver.Bytes(),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				// due to the differences in decimals, 1e9 on EVM results to 1 on SVM
				{Token: destToken.Bytes(), Amount: big.NewInt(1)},
			},
			ExtraArgs:      extraArgs,
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		},
		// {
		// 	Name:        "Send N tokens to contract",
		// 	SourceChain: destChain,
		// 	DestChain:   sourceChain,
		// 	Tokens: []router.ClientEVMTokenAmount{
		// 		{
		// 			Token:  selfServeDestToken.Address(),
		// 			Amount: oneE9,
		// 		},
		// 		{
		// 			Token:  destToken.Address(),
		// 			Amount: oneE9,
		// 		},
		// 		{
		// 			Token:  selfServeDestToken.Address(),
		// 			Amount: oneE9,
		// 		},
		// 	},
		// 	Receiver:  state.Chains[sourceChain].Receiver.Address().Bytes(),
		// 	ExtraArgs: testhelpers.MakeEVMExtraArgsV2(300_000, false),
		// 	ExpectedTokenBalances: []testhelpers.ExpectedBalance{
		// 		{selfServeSrcToken.Address().Bytes(), new(big.Int).Add(oneE18, oneE18)},
		// 		{srcToken.Address().Bytes(), oneE18},
		// 	},
		// 	ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		// },
	}

	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances :=
		testhelpers.TransferMultiple(ctx, t, e, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		e,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		e,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, e, expectedTokenBalances)
}

func TestTokenTransfer_Solana2EVM(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()

	tenv, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfUsersPerChain(3),
		testhelpers.WithSolChains(1))

	e := tenv.Env
	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(e.Chains), 2)

	allChainSelectors := maps.Keys(e.Chains)
	allSolChainSelectors := maps.Keys(e.SolChains)
	sourceChain, destChain := allSolChainSelectors[0], allChainSelectors[0]
	sender := e.SolChains[sourceChain].DeployerKey
	ownerSourceChain := sender.PublicKey()
	ownerDestChain := e.Chains[destChain].DeployerKey

	require.GreaterOrEqual(t, len(tenv.Users[destChain]), 2) // TODO: ???

	const oneE9 uint64 = 1e9

	// Deploy tokens and pool by CCIP Owner
	destToken, _, srcToken, err := testhelpers.DeployTransferableTokenSolana(
		lggr,
		e,
		destChain,
		sourceChain,
		ownerDestChain,
		"OWNER_TOKEN",
	)
	require.NoError(t, err)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &tenv, state, sourceChain, destChain, false)

	// TODO: handle in setup
	deployer := e.SolChains[sourceChain].DeployerKey
	rpcClient := e.SolChains[sourceChain].Client

	// create ATA for user
	tokenProgram := solana.TokenProgramID
	wSOL := solana.SolMint
	ixAtaUser, deployerWSOL, uerr := soltokens.CreateAssociatedTokenAccount(tokenProgram, wSOL, deployer.PublicKey(), deployer.PublicKey())
	require.NoError(t, uerr)

	billingSignerPDA, _, err := solstate.FindFeeBillingSignerPDA(state.SolChains[sourceChain].Router)
	require.NoError(t, err)

	// Approve CCIP to transfer the user's token for billing
	ixApprove, err := soltokens.TokenApproveChecked(1e9, 9, tokenProgram, deployerWSOL, wSOL, billingSignerPDA, deployer.PublicKey(), []solana.PublicKey{})
	require.NoError(t, err)

	soltestutils.SendAndConfirm(ctx, t, rpcClient, []solana.Instruction{ixAtaUser, ixApprove}, *deployer, solconfig.DefaultCommitment)

	// fund user WSOL (transfer SOL + syncNative)
	transferAmount := 1.0 * solana.LAMPORTS_PER_SOL
	ixTransfer, err := soltokens.NativeTransfer(tokenProgram, transferAmount, deployer.PublicKey(), deployerWSOL)
	require.NoError(t, err)
	ixSync, err := soltokens.SyncNative(tokenProgram, deployerWSOL)
	require.NoError(t, err)
	soltestutils.SendAndConfirm(ctx, t, rpcClient, []solana.Instruction{ixTransfer, ixSync}, *deployer, solconfig.DefaultCommitment)
	// END: handle in setup

	testhelpers.MintAndAllow(
		t,
		e,
		state,
		map[uint64][]testhelpers.MintTokenInfo{
			// sourceChain: {
			// 	testhelpers.NewMintTokenInfo(ownerSourceChain, srcToken),
			// },
			destChain: {
				testhelpers.NewMintTokenInfo(ownerDestChain, destToken),
			},
		},
	)

	// TODO: extract as MintAndAllow on Solana? mint already previously happened
	userTokenAccount, _, err := soltokens.FindAssociatedTokenAddress(solana.Token2022ProgramID, srcToken, ownerSourceChain)
	require.NoError(t, err)

	ixApprove2, err := soltokens.TokenApproveChecked(1000, 9, solana.Token2022ProgramID, userTokenAccount, srcToken, billingSignerPDA, ownerSourceChain, nil)
	require.NoError(t, err)

	ixs := []solana.Instruction{ixApprove2}
	result := soltestutils.SendAndConfirm(ctx, t, rpcClient, ixs, *sender, solconfig.DefaultCommitment)
	require.NotNil(t, result)
	// END: extract as MintAndAllow on Solana

	// ---
	emptyEVMExtraArgsV2 := []byte{}
	extraArgs := emptyEVMExtraArgsV2

	// extraArgs := soltestutils.MustSerializeExtraArgs(t, fee_quoter.EVMExtraArgsV2{
	// 	GasLimit: bin.Uint128{Lo: 500_000, Hi: 0}, // TODO: why is default not enough
	// }, solccip.EVMExtraArgsV2Tag)

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:        "Send token to contract",
			SourceChain: sourceChain,
			DestChain:   destChain,
			SolTokens: []ccip_router.SVMTokenAmount{
				{
					Token:  srcToken,
					Amount: 1,
				},
			},
			Receiver: state.Chains[destChain].Receiver.Address().Bytes(),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				// due to the differences in decimals, 1 on SVM results to 1e9 on EVM
				{Token: common.LeftPadBytes(destToken.Address().Bytes(), 32), Amount: new(big.Int).SetUint64(oneE9)},
			},
			ExtraArgs:      extraArgs,
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		},
		// {
		// 	Name:        "Send N tokens to contract",
		// 	SourceChain: destChain,
		// 	DestChain:   sourceChain,
		// 	Tokens: []router.ClientEVMTokenAmount{
		// 		{
		// 			Token:  selfServeDestToken.Address(),
		// 			Amount: oneE9,
		// 		},
		// 		{
		// 			Token:  destToken.Address(),
		// 			Amount: oneE9,
		// 		},
		// 		{
		// 			Token:  selfServeDestToken.Address(),
		// 			Amount: oneE9,
		// 		},
		// 	},
		// 	Receiver:  state.Chains[sourceChain].Receiver.Address().Bytes(),
		// 	ExtraArgs: testhelpers.MakeEVMExtraArgsV2(300_000, false),
		// 	ExpectedTokenBalances: []testhelpers.ExpectedBalance{
		// 		{selfServeSrcToken.Address().Bytes(), new(big.Int).Add(oneE18, oneE18)},
		// 		{srcToken.Address().Bytes(), oneE18},
		// 	},
		// 	ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		// },
	}

	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances :=
		testhelpers.TransferMultiple(ctx, t, e, state, tcs)

	// HACK: we need to replay blocks only after the CCIP plugin has already properly booted
	sleepAndReplay(t, tenv, sourceChain, destChain)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		e,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		e,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, e, expectedTokenBalances)
}
