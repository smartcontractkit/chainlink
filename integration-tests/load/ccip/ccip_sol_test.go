// Will be deleted only tmp file for testing message passing with CRIB
//
//nolint:golang-ci-lint // Will be deleted only tmp
package ccip

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/message_hasher"
	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	soltestutils "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/testutils"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solccip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solstate "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	soltokens "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	soltesthelpers "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/solana"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/crib"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestCCIPSol2EvmCRIB(t *testing.T) {
	// comment out when executing the test
	// t.Skip("Skipping test as this test should not be auto triggered")
	lggr := logger.Test(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	// get user defined configurations
	config, err := tc.GetConfig([]string{"Load"}, tc.CCIP)
	require.NoError(t, err)
	userOverrides := config.CCIP.Load

	// generate environment from crib-produced files
	cribEnv := crib.NewDevspaceEnvFromStateDir(lggr, *userOverrides.CribEnvDirectory)
	cribDeployOutput, err := cribEnv.GetConfig(simChainTestKey, solTestKey)
	require.NoError(t, err)
	cribEnvironment, err := crib.NewDeployEnvironmentFromCribOutput(lggr, cribDeployOutput)
	require.NoError(t, err)
	require.NotNil(t, cribEnvironment)
	userOverrides.Validate(t, cribEnvironment)

	allChainSelectors := cribEnvironment.AllChainSelectors()
	allSolChainSelectors := cribEnvironment.AllChainSelectorsSolana()

	e := testhelpers.DeployedEnv{
		Env:          *cribEnvironment,
		HomeChainSel: allChainSelectors[0],
		FeedChainSel: allChainSelectors[0],
	}
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := allSolChainSelectors[0]
	destChain := allChainSelectors[0]
	t.Log("All chain selectors:", allChainSelectors,
		", sol chain selectors:", allSolChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector:", sourceChain,
		", dest chain selector:", destChain,
	)

	var (
		// nonce  uint64
		sender = common.LeftPadBytes(e.Env.SolChains[sourceChain].DeployerKey.PublicKey().Bytes(), 32)
		out    mt.TestCaseOutput
		setup  = mt.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
			true,  // validateResp
		)
	)

	if false {
		// TODO: handle in setup
		deployer := e.Env.SolChains[sourceChain].DeployerKey
		rpcClient := e.Env.SolChains[sourceChain].Client

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
	}

	emptyEVMExtraArgsV2 := []byte{}
	SetProgramIDsSafe(state.SolChains[sourceChain])

	t.Logf("Dest chain %v", state.Chains[destChain].Receiver.Address())
	t.Run("message to contract implementing CCIPReceiver", func(t *testing.T) {
		extraArgs := emptyEVMExtraArgsV2
		latestHead, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		out = mt.Run(
			mt.TestCase{
				TestSetup:              setup,
				Replayed:               true,
				Nonce:                  1,
				Receiver:               state.Chains[destChain].Receiver.Address().Bytes(),
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              extraArgs, // default extraArgs
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
						iter, err := state.Chains[destChain].Receiver.FilterMessageReceived(&bind.FilterOpts{
							Context: ctx,
							Start:   latestHead,
						})
						require.NoError(t, err)
						require.True(t, iter.Next())
						// MessageReceived doesn't emit the data unfortunately, so can't check that.
					},
				},
			},
		)

		_ = out // avoid unused error
	})
}

func TestCCIPEvm2SolCRIB(t *testing.T) {
	// comment out when executing the test
	// t.Skip("Skipping test as this test should not be auto triggered")
	lggr := logger.Test(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	// get user defined configurations
	config, err := tc.GetConfig([]string{"Load"}, tc.CCIP)
	require.NoError(t, err)
	userOverrides := config.CCIP.Load

	// generate environment from crib-produced files
	cribEnv := crib.NewDevspaceEnvFromStateDir(lggr, *userOverrides.CribEnvDirectory)
	cribDeployOutput, err := cribEnv.GetConfig(simChainTestKey, solTestKey)
	require.NoError(t, err)
	cribEnvironment, err := crib.NewDeployEnvironmentFromCribOutput(lggr, cribDeployOutput)
	require.NoError(t, err)
	require.NotNil(t, cribEnvironment)
	userOverrides.Validate(t, cribEnvironment)

	allChainSelectors := cribEnvironment.AllChainSelectors()
	allSolChainSelectors := cribEnvironment.AllChainSelectorsSolana()

	e := testhelpers.DeployedEnv{
		Env:          *cribEnvironment,
		HomeChainSel: allChainSelectors[0],
		FeedChainSel: allChainSelectors[0],
	}
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := allChainSelectors[0]
	destChain := allSolChainSelectors[0]
	t.Log("All chain selectors:", allChainSelectors,
		", sol chain selectors:", allSolChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector:", sourceChain,
		", dest chain selector:", destChain,
	)

	var (
		nonce  uint64
		sender = common.LeftPadBytes(e.Env.Chains[sourceChain].DeployerKey.From.Bytes(), 32)
		out    mt.TestCaseOutput
		setup  = mt.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
			true,  // validateResp
		)
	)

	t.Run("message to contract implementing CCIPReceiver", func(t *testing.T) {
		receiverProgram := state.SolChains[destChain].Receiver
		receiver := receiverProgram.Bytes()
		receiverTargetAccountPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("counter")}, receiverProgram)
		receiverExternalExecutionConfigPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("external_execution_config")}, receiverProgram)

		accounts := [][32]byte{
			receiverExternalExecutionConfigPDA,
			receiverTargetAccountPDA,
			solana.SystemProgramID,
		}

		extraArgs, err := ccipevm.SerializeClientSVMExtraArgsV1(message_hasher.ClientSVMExtraArgsV1{
			AccountIsWritableBitmap: solccip.GenerateBitMapForIndexes([]int{0, 1}),
			Accounts:                accounts,
			ComputeUnits:            80_000,
		})
		require.NoError(t, err)

		// check that counter is 0
		var receiverCounterAccount soltesthelpers.ReceiverCounter
		err = solcommon.GetAccountDataBorshInto(ctx, e.Env.SolChains[destChain].Client, receiverTargetAccountPDA, solconfig.DefaultCommitment, &receiverCounterAccount)
		require.NoError(t, err, "failed to get account info")
		// require.Equal(t, uint8(0), receiverCounterAccount.Value)

		out = mt.Run(
			mt.TestCase{
				TestSetup:              setup,
				Replayed:               true,
				Nonce:                  nonce,
				Receiver:               receiver,
				MsgData:                []byte("hello CCIPReceiver"),
				ExtraArgs:              extraArgs,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
						var receiverCounterAccount soltesthelpers.ReceiverCounter
						err = solcommon.GetAccountDataBorshInto(ctx, e.Env.SolChains[destChain].Client, receiverTargetAccountPDA, solconfig.DefaultCommitment, &receiverCounterAccount)
						require.NoError(t, err, "failed to get account info")
						require.Equal(t, uint8(1), receiverCounterAccount.Value)
					},
				},
			},
		)
	})

	_ = out
}

func TestTokenTransfer_EVM2SolanaCRIB(t *testing.T) {
	lggr := logger.Test(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	// get user defined configurations
	config, err := tc.GetConfig([]string{"Load"}, tc.CCIP)
	require.NoError(t, err)
	userOverrides := config.CCIP.Load

	// generate environment from crib-produced files
	cribEnv := crib.NewDevspaceEnvFromStateDir(lggr, *userOverrides.CribEnvDirectory)
	cribDeployOutput, err := cribEnv.GetConfig(simChainTestKey, solTestKey)
	require.NoError(t, err)
	cribEnvironment, err := crib.NewDeployEnvironmentFromCribOutput(lggr, cribDeployOutput)
	require.NoError(t, err)
	require.NotNil(t, cribEnvironment)
	userOverrides.Validate(t, cribEnvironment)

	allChainSelectors := cribEnvironment.AllChainSelectors()
	allSolChainSelectors := cribEnvironment.AllChainSelectorsSolana()
	e := testhelpers.DeployedEnv{
		Env:          *cribEnvironment,
		HomeChainSel: allChainSelectors[0],
		FeedChainSel: allChainSelectors[0],
	}
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain, destChain := allChainSelectors[0], allSolChainSelectors[0]
	ownerSourceChain := e.Env.Chains[sourceChain].DeployerKey
	// ownerDestChain := e.SolChains[destChain].DeployerKey

	// require.GreaterOrEqual(t, len(e.Users[sourceChain]), 2) // TODO: ???

	oneE9 := new(big.Int).SetUint64(1e9)

	SetProgramIDsSafe(state.SolChains[destChain])

	// Deploy tokens and pool by CCIP Owner
	srcToken, _, destToken, err := testhelpers.DeployTransferableTokenSolana(
		t,
		lggr,
		e.Env,
		sourceChain,
		destChain,
		ownerSourceChain,
		e.Env.ExistingAddresses, //nolint:staticcheck // SA1019
		"OWNER_TOKEN",
	)
	require.NoError(t, err)

	testhelpers.MintAndAllow(
		t,
		e.Env,
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
		testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		e.Env,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		e.Env,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, e.Env, expectedTokenBalances)
}

func TestTokenTransfer_Solana2EVMCRIB(t *testing.T) {
	lggr := logger.Test(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	// get user defined configurations
	config, err := tc.GetConfig([]string{"Load"}, tc.CCIP)
	require.NoError(t, err)
	userOverrides := config.CCIP.Load

	// generate environment from crib-produced files
	cribEnv := crib.NewDevspaceEnvFromStateDir(lggr, *userOverrides.CribEnvDirectory)
	cribDeployOutput, err := cribEnv.GetConfig(simChainTestKey, solTestKey)
	require.NoError(t, err)
	cribEnvironment, err := crib.NewDeployEnvironmentFromCribOutput(lggr, cribDeployOutput)
	require.NoError(t, err)
	require.NotNil(t, cribEnvironment)
	userOverrides.Validate(t, cribEnvironment)

	allChainSelectors := cribEnvironment.AllChainSelectors()
	allSolChainSelectors := cribEnvironment.AllChainSelectorsSolana()
	e := testhelpers.DeployedEnv{
		Env:          *cribEnvironment,
		HomeChainSel: allChainSelectors[0],
		FeedChainSel: allChainSelectors[0],
	}
	state, err := changeset.LoadOnchainState(e.Env)
	require.NoError(t, err)
	sourceChain, destChain := allSolChainSelectors[0], allChainSelectors[0]
	sender := e.Env.SolChains[sourceChain].DeployerKey
	ownerSourceChain := sender.PublicKey()
	ownerDestChain := e.Env.Chains[destChain].DeployerKey

	// require.GreaterOrEqual(t, len(e.Users[destChain]), 2) // TODO: ???

	const oneE9 uint64 = 1e9

	SetProgramIDsSafe(state.SolChains[sourceChain])

	// Deploy tokens and pool by CCIP Owner
	destToken, _, srcToken, err := testhelpers.DeployTransferableTokenSolana(
		t,
		lggr,
		e.Env,
		destChain,
		sourceChain,
		ownerDestChain,
		e.Env.ExistingAddresses, //nolint:staticcheck // SA1019
		"OWNER_TOKEN",
	)
	require.NoError(t, err)

	// TODO: handle in setup
	deployer := e.Env.SolChains[sourceChain].DeployerKey
	rpcClient := e.Env.SolChains[sourceChain].Client

	if false {
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
	}

	testhelpers.MintAndAllow(
		t,
		e.Env,
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

	billingSignerPDA, _, err := solstate.FindFeeBillingSignerPDA(state.SolChains[sourceChain].Router)
	require.NoError(t, err)

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
		testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		e.Env,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		e.Env,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, e.Env, expectedTokenBalances)
}
