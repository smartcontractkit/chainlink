package ccip

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	module_lock_release_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/lock_release_token_pool"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	suiBind "github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	sui_deployment "github.com/smartcontractkit/chainlink-sui/deployment"
	sui_cs "github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIPTokenTransfer_Sui2EVM_LockReleaseTokenPool(t *testing.T) {
	e, sourceChain, destChain := testSetupTokenTransferSui2Evm(t)

	feeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 100000000000)
	linkTokenOutput1 := mintLinkTokenOnSui(t, e.Env, sourceChain, 5000000000)
	linkTokenOutput2 := mintLinkTokenOnSui(t, e.Env, sourceChain, 2000000000)
	linkTokenOutput3 := mintLinkTokenOnSui(t, e.Env, sourceChain, 1500000000)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// Receiver Address
	ccipReceiverAddress := state.Chains[destChain].Receiver.Address()

	// Token Pool setup on both SUI and EVM
	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndLockReleaseTokenPoolDeploymentForSUI(e.Env, sourceChain, destChain, []testhelpers.TokenPoolRateLimiterConfig{
		{
			RemoteChainSelector: destChain,
			OutboundIsEnabled:   false,
			OutboundCapacity:    100000,
			OutboundRate:        100,
			InboundIsEnabled:    false,
			InboundCapacity:     100000,
			InboundRate:         100,
		},
	}) // SourceChain = SUI, destChain = EVM
	require.NoError(t, err)
	e.Env = updatedEnv

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       updatedEnv.BlockChains.EVMChains()[destChain].DeployerKey.From.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			FeeToken:       feeTokenOutput.Objects.MintedLinkTokenObjectId,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeLockRelease,
					Token:         linkTokenOutput1.Objects.MintedLinkTokenObjectId,
					Amount:        5000000000, // Send 5 LINK to EVM
				},
			},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(5e18),
				},
			},
		},
	}

	ctx := testhelpers.Context(t)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, updatedEnv, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		updatedEnv,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		updatedEnv,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, updatedEnv, expectedTokenBalances)

	suiState, err := sui_deployment.LoadOnchainStatesui(e.Env)
	require.NoError(t, err)

	lockReleaseTokenPool, err := module_lock_release_token_pool.NewLockReleaseTokenPool(suiState[sourceChain].LnRTokenPools[testhelpers.TokenSymbolLINK].PackageID, e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)
	balance, err := lockReleaseTokenPool.DevInspect().GetBalance(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
	}, []string{state.SuiChains[sourceChain].LinkTokenAddress + "::link::LINK"}, suiBind.Object{Id: suiState[sourceChain].LnRTokenPools[testhelpers.TokenSymbolLINK].StateObjectId})
	require.NoError(t, err)
	t.Log("Current Balance: ", balance)
	require.Equal(t, uint64(5000000000), balance)

	deps := getOpTxDeps(e.Env.BlockChains.SuiChains()[sourceChain])

	// enable allowlist but not adding the current sender to the allowlist
	_, err = operations.ExecuteOperation(e.Env.OperationsBundle, lockreleasetokenpoolops.LockReleaseTokenPoolProvideLiquidityOp, deps, lockreleasetokenpoolops.LockReleaseTokenPoolProvideLiquidityInput{
		LockReleaseTokenPoolPackageId: suiState[sourceChain].LnRTokenPools[testhelpers.TokenSymbolLINK].PackageID,
		StateObjectId:                 suiState[sourceChain].LnRTokenPools[testhelpers.TokenSymbolLINK].StateObjectId,
		RebalancerCapObjectId:         suiState[sourceChain].LnRTokenPools[testhelpers.TokenSymbolLINK].RebalancerCapIds[0],
		CoinObjectTypeArg:             state.SuiChains[sourceChain].LinkTokenAddress + "::link::LINK",
		Coin:                          linkTokenOutput2.Objects.MintedLinkTokenObjectId,
	})
	require.NoError(t, err)

	balance, err = lockReleaseTokenPool.DevInspect().GetBalance(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
	}, []string{state.SuiChains[sourceChain].LinkTokenAddress + "::link::LINK"}, suiBind.Object{Id: suiState[sourceChain].LnRTokenPools[testhelpers.TokenSymbolLINK].StateObjectId})
	require.NoError(t, err)
	t.Log("New Balance: ", balance)
	require.Equal(t, uint64(7000000000), balance)

	suifeeQuoter, err := module_fee_quoter.NewFeeQuoter(suiState[sourceChain].CCIPAddress, e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)

	suiFeeQuoterDestChainConfig, err := suifeeQuoter.DevInspect().GetDestChainConfig(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
	}, suiBind.Object{Id: suiState[sourceChain].CCIPObjectRef}, destChain)
	require.NoError(t, err, "Failed to get destination chain config")

	t.Run("Send invalid token to CCIP Receiver - should fail", func(t *testing.T) {
		msg := testhelpers.SuiSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  feeTokenOutput.Objects.MintedLinkTokenObjectId,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(suiFeeQuoterDestChainConfig.MaxPerMsgGasLimit)), false),
			TokenAmounts: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeLockRelease,
					Token:         "0x0",
					Amount:        1e9,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertSuiSourceRevertExpectedError(t, err, "failed to resolve CallArg at index 2", "failed to resolve UnresolvedObject 0x0000000000000000000000000000000000000000000000000000000000000000")
		t.Log("Expected error: ", err)
	})

	t.Run("Send token to CCIP Receiver setting gas above max gas allowed - should fail", func(t *testing.T) {
		msg := testhelpers.SuiSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  feeTokenOutput.Objects.MintedLinkTokenObjectId,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(suiFeeQuoterDestChainConfig.MaxPerMsgGasLimit+10)), false),
			TokenAmounts: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeLockRelease,
					Token:         linkTokenOutput3.Objects.MintedLinkTokenObjectId,
					Amount:        1500000000,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertSuiSourceRevertExpectedError(t, err, "transaction failed with error", "function_name: Some(\"resolve_generic_gas_limit\") }, 18)")
		t.Log("Expected error: ", err)
	})
}

func Test_CCIPTokenTransfer_Sui2EVM_BurnMintTokenPool(t *testing.T) {
	e, sourceChain, destChain := testSetupTokenTransferSui2Evm(t)

	feeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000000)
	linkTokenOutput1 := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000)
	linkTokenOutput2 := mintLinkTokenOnSui(t, e.Env, sourceChain, 2000000000)
	linkTokenOutput3 := mintLinkTokenOnSui(t, e.Env, sourceChain, 1500000000)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// Receiver Address
	ccipReceiverAddress := state.Chains[destChain].Receiver.Address()

	// Token Pool setup on both SUI and EVM
	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndBurnMintTokenPoolDeploymentForSUI(e.Env, sourceChain, destChain, []testhelpers.TokenPoolRateLimiterConfig{
		{
			RemoteChainSelector: destChain,
			OutboundIsEnabled:   false,
			OutboundCapacity:    100000,
			OutboundRate:        100,
			InboundIsEnabled:    false,
			InboundCapacity:     100000,
			InboundRate:         100,
		},
	}) // SourceChain = SUI, destChain = EVM
	require.NoError(t, err)
	e.Env = updatedEnv

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       updatedEnv.BlockChains.EVMChains()[destChain].DeployerKey.From.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			FeeToken:       feeTokenOutput.Objects.MintedLinkTokenObjectId,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         linkTokenOutput1.Objects.MintedLinkTokenObjectId,
					Amount:        1000000000, // Send 1 LINK to EVM
				},
			},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(1e18),
				},
			},
		},
		{
			Name:           "Send token to Receiver",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       ccipReceiverAddress.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         linkTokenOutput2.Objects.MintedLinkTokenObjectId,
					Amount:        2000000000, // Send 2 LINK to EVM
				},
			},
			FeeToken: feeTokenOutput.Objects.MintedLinkTokenObjectId,
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(2e18),
				},
			},
		},
	}

	ctx := testhelpers.Context(t)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, updatedEnv, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		updatedEnv,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		updatedEnv,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, updatedEnv, expectedTokenBalances)

	suiState, err := sui_deployment.LoadOnchainStatesui(e.Env)
	require.NoError(t, err)

	suifeeQuoter, err := module_fee_quoter.NewFeeQuoter(suiState[sourceChain].CCIPAddress, e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)

	suiFeeQuoterDestChainConfig, err := suifeeQuoter.DevInspect().GetDestChainConfig(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
	}, suiBind.Object{Id: suiState[sourceChain].CCIPObjectRef}, destChain)
	require.NoError(t, err, "Failed to get destination chain config")

	t.Run("Send invalid token to CCIP Receiver - should fail", func(t *testing.T) {
		msg := testhelpers.SuiSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  feeTokenOutput.Objects.MintedLinkTokenObjectId,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(suiFeeQuoterDestChainConfig.MaxPerMsgGasLimit)), false),
			TokenAmounts: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         "0x0",
					Amount:        1e9,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertSuiSourceRevertExpectedError(t, err, "failed to resolve CallArg at index 2", "failed to resolve UnresolvedObject 0x0000000000000000000000000000000000000000000000000000000000000000")
		t.Log("Expected error: ", err)
	})

	t.Run("Send token to CCIP Receiver setting gas above max gas allowed - should fail", func(t *testing.T) {
		msg := testhelpers.SuiSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:      []byte("Hello, World!"),
			FeeToken:  feeTokenOutput.Objects.MintedLinkTokenObjectId,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(suiFeeQuoterDestChainConfig.MaxPerMsgGasLimit+10)), false),
			TokenAmounts: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         linkTokenOutput3.Objects.MintedLinkTokenObjectId,
					Amount:        1500000000,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertSuiSourceRevertExpectedError(t, err, "transaction failed with error", "function_name: Some(\"resolve_generic_gas_limit\") }, 18)")
		t.Log("Expected error: ", err)
	})
}

func Test_CCIPTokenTransfer_Sui2EVM_BurnMintTokenPool_ThenGloballyCursedUncursed(t *testing.T) {
	e, sourceChain, destChain := testSetupTokenTransferSui2Evm(t)

	feeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000000)
	linkTokenOutput1 := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000)
	linkTokenOutput2 := mintLinkTokenOnSui(t, e.Env, sourceChain, 2000000000)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// Receiver Address
	ccipReceiverAddress := state.Chains[destChain].Receiver.Address()

	// Token Pool setup on both SUI and EVM
	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndBurnMintTokenPoolDeploymentForSUI(e.Env, sourceChain, destChain, []testhelpers.TokenPoolRateLimiterConfig{
		{
			RemoteChainSelector: destChain,
			OutboundIsEnabled:   false,
			OutboundCapacity:    100000,
			OutboundRate:        100,
			InboundIsEnabled:    false,
			InboundCapacity:     100000,
			InboundRate:         100,
		},
	}) // SourceChain = SUI, destChain = EVM
	require.NoError(t, err)
	e.Env = updatedEnv

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       updatedEnv.BlockChains.EVMChains()[destChain].DeployerKey.From.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			FeeToken:       feeTokenOutput.Objects.MintedLinkTokenObjectId,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         linkTokenOutput1.Objects.MintedLinkTokenObjectId,
					Amount:        1000000000, // Send 1 LINK to EVM
				},
			},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(1e18),
				},
			},
		},
	}

	ctx := testhelpers.Context(t)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, updatedEnv, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		updatedEnv,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		updatedEnv,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, updatedEnv, expectedTokenBalances)

	suiState, err := sui_deployment.LoadOnchainStatesui(e.Env)
	require.NoError(t, err)

	suiChain := e.Env.BlockChains.SuiChains()[sourceChain]
	require.NotNil(t, suiChain)

	deps := getOpTxDeps(suiChain)

	// curse globally
	_, err = operations.ExecuteOperation(e.Env.OperationsBundle, ccipops.RMNRemoteCurseOp, deps, ccipops.RMNRemoteCurseInput{
		CCIPPackageId:    suiState[sourceChain].CCIPAddress,
		StateObjectId:    suiState[sourceChain].CCIPObjectRef,
		OwnerCapObjectId: suiState[sourceChain].CCIPOwnerCapObjectId,
		Subject: []byte{
			0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		},
	})
	require.NoError(t, err)

	t.Run("Destination chain is cursed - should fail", func(t *testing.T) {
		msg := testhelpers.SuiSendRequest{
			Receiver: common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:     []byte("Hello, World!"),
			FeeToken: feeTokenOutput.Objects.MintedLinkTokenObjectId,
			TokenAmounts: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         linkTokenOutput2.Objects.MintedLinkTokenObjectId,
					Amount:        2000000000,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertSuiSourceRevertExpectedError(t, err, "failed to execute ccip_send with err: transaction failed with error: MoveAbort", "function_name: Some(\"validate_lock_or_burn\") }, 3)")
		t.Log("Expected error: ", err)
	})

	// uncurse globally
	_, err = operations.ExecuteOperation(e.Env.OperationsBundle, ccipops.RMNRemoteUncurseOp, deps, ccipops.RMNRemoteUncurseInput{
		CCIPPackageId:    suiState[sourceChain].CCIPAddress,
		StateObjectId:    suiState[sourceChain].CCIPObjectRef,
		OwnerCapObjectId: suiState[sourceChain].CCIPOwnerCapObjectId,
		Subject: []byte{
			0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		},
	})
	require.NoError(t, err)

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	tcs = []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       updatedEnv.BlockChains.EVMChains()[destChain].DeployerKey.From.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			FeeToken:       feeTokenOutput.Objects.MintedLinkTokenObjectId,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         linkTokenOutput2.Objects.MintedLinkTokenObjectId,
					Amount:        2000000000, // Send 2 LINK to EVM
				},
			},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(3e18), // 1e18 from first transfer + 2e18 from second transfer
				},
			},
		},
	}

	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances = testhelpers.TransferMultiple(ctx, t, updatedEnv, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		updatedEnv,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates = testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		updatedEnv,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, updatedEnv, expectedTokenBalances)
}

func Test_CCIPTokenTransfer_Sui2EVM_BurnMintTokenPool_WithAllowlist(t *testing.T) {
	e, sourceChain, destChain := testSetupTokenTransferSui2Evm(t)
	feeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000000)
	linkTokenOutput1 := mintLinkTokenOnSui(t, e.Env, sourceChain, 2000000000)
	linkTokenOutput2 := mintLinkTokenOnSui(t, e.Env, sourceChain, 1500000000)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// Receiver Address
	ccipReceiverAddress := state.Chains[destChain].Receiver.Address()

	// Token Pool setup on both SUI and EVM
	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndBurnMintTokenPoolDeploymentForSUI(e.Env, sourceChain, destChain, []testhelpers.TokenPoolRateLimiterConfig{
		{
			RemoteChainSelector: destChain,
			OutboundIsEnabled:   false,
			OutboundCapacity:    100000,
			OutboundRate:        100,
			InboundIsEnabled:    false,
			InboundCapacity:     100000,
			InboundRate:         100,
		},
	}) // SourceChain = SUI, destChain = EVM
	require.NoError(t, err)
	e.Env = updatedEnv

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to Receiver",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       ccipReceiverAddress.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         linkTokenOutput1.Objects.MintedLinkTokenObjectId,
					Amount:        2000000000, // send 2 LINK to EVM
				},
			},
			FeeToken: feeTokenOutput.Objects.MintedLinkTokenObjectId,
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(2e18),
				},
			},
		},
	}

	ctx := testhelpers.Context(t)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, updatedEnv, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		updatedEnv,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		updatedEnv,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, updatedEnv, expectedTokenBalances)

	suiChain := e.Env.BlockChains.SuiChains()[sourceChain]

	deps := getOpTxDeps(suiChain)

	// reload state to include the BnM token pool
	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// enable allowlist but not adding the current sender to the allowlist
	_, err = operations.ExecuteOperation(e.Env.OperationsBundle, burnminttokenpoolops.BurnMintTokenPoolSetAllowlistEnabledOp, deps, burnminttokenpoolops.BurnMintTokenPoolSetAllowlistEnabledInput{
		BurnMintPackageId: state.SuiChains[sourceChain].BnMTokenPools[testhelpers.TokenSymbolLINK].PackageID,
		StateObjectId:     state.SuiChains[sourceChain].BnMTokenPools[testhelpers.TokenSymbolLINK].StateObjectId,
		OwnerCap:          state.SuiChains[sourceChain].BnMTokenPools[testhelpers.TokenSymbolLINK].OwnerCapObjectId,
		CoinObjectTypeArg: state.SuiChains[sourceChain].LinkTokenAddress + "::link::LINK",
		Enabled:           true,
	})
	require.NoError(t, err)

	t.Run("Sender not in allowlist - should fail", func(t *testing.T) {
		msg := testhelpers.SuiSendRequest{
			Receiver: common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32),
			Data:     []byte("Hello, World!"),
			FeeToken: feeTokenOutput.Objects.MintedLinkTokenObjectId,
			TokenAmounts: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         linkTokenOutput2.Objects.MintedLinkTokenObjectId,
					Amount:        1500000000,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertSuiSourceRevertExpectedError(t, err, "failed to execute ccip_send with err: transaction failed with error: MoveAbort", "function_name: Some(\"validate_lock_or_burn\") }, 1)")
		t.Log("Expected error: ", err)
	})

	signerAddress, err := deps.Signer.GetAddress()
	require.NoError(t, err)
	_, err = operations.ExecuteOperation(e.Env.OperationsBundle, burnminttokenpoolops.BurnMintTokenPoolApplyAllowlistUpdatesOp, deps, burnminttokenpoolops.BurnMintTokenPoolApplyAllowlistUpdatesInput{
		BurnMintPackageId: state.SuiChains[sourceChain].BnMTokenPools[testhelpers.TokenSymbolLINK].PackageID,
		StateObjectId:     state.SuiChains[sourceChain].BnMTokenPools[testhelpers.TokenSymbolLINK].StateObjectId,
		OwnerCap:          state.SuiChains[sourceChain].BnMTokenPools[testhelpers.TokenSymbolLINK].OwnerCapObjectId,
		CoinObjectTypeArg: state.SuiChains[sourceChain].LinkTokenAddress + "::link::LINK",
		Removes:           []string{},
		Adds:              []string{signerAddress},
	})
	require.NoError(t, err)

	tcs = []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to Receiver",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       ccipReceiverAddress.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         linkTokenOutput2.Objects.MintedLinkTokenObjectId,
					Amount:        1500000000, // send 1.5 LINK to EVM
				},
			},
			FeeToken: feeTokenOutput.Objects.MintedLinkTokenObjectId,
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(3.5e18),
				},
			},
		},
	}

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances = testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		e.Env,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates = testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		e.Env,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, e.Env, expectedTokenBalances)
}

func Test_CCIPTokenTransfer_Sui2EVM_BurnMintTokenPool_WithRateLimit(t *testing.T) {
	e, sourceChain, destChain := testSetupTokenTransferSui2Evm(t)
	feeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000000)
	linkTokenOutput1 := mintLinkTokenOnSui(t, e.Env, sourceChain, 10)
	linkTokenOutput2 := mintLinkTokenOnSui(t, e.Env, sourceChain, 60)
	linkTokenOutput3 := mintLinkTokenOnSui(t, e.Env, sourceChain, 999999999999)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// Receiver Address
	ccipReceiverAddress := state.Chains[destChain].Receiver.Address()

	// Token Pool setup on both SUI and EVM
	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndBurnMintTokenPoolDeploymentForSUI(e.Env, sourceChain, destChain, []testhelpers.TokenPoolRateLimiterConfig{
		{
			RemoteChainSelector: destChain,
			OutboundIsEnabled:   true,
			OutboundCapacity:    10000000,
			OutboundRate:        100,
			InboundIsEnabled:    true,
			InboundCapacity:     10000000,
			InboundRate:         100,
		},
	}) // SourceChain = SUI, destChain = EVM
	require.NoError(t, err)
	e.Env = updatedEnv

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       updatedEnv.BlockChains.EVMChains()[destChain].DeployerKey.From.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			FeeToken:       feeTokenOutput.Objects.MintedLinkTokenObjectId,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         linkTokenOutput1.Objects.MintedLinkTokenObjectId,
					Amount:        10,
				},
			},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(1e10),
				},
			},
		},
		{
			Name:           "Send token to Receiver",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       ccipReceiverAddress.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         linkTokenOutput2.Objects.MintedLinkTokenObjectId,
					Amount:        60,
				},
			},
			FeeToken: feeTokenOutput.Objects.MintedLinkTokenObjectId,
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(6e10),
				},
			},
		},
	}

	ctx := testhelpers.Context(t)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, updatedEnv, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		updatedEnv,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		updatedEnv,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, updatedEnv, expectedTokenBalances)

	t.Run("Send token above Sui's outbound rate limit - should fail", func(t *testing.T) {
		msg := testhelpers.SuiSendRequest{
			Receiver: common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32),
			Data:     []byte("Hello, World!"),
			FeeToken: feeTokenOutput.Objects.MintedLinkTokenObjectId,
			TokenAmounts: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         linkTokenOutput3.Objects.MintedLinkTokenObjectId,
					Amount:        999999999999,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertSuiSourceRevertExpectedError(t, err, "failed to execute ccip_send with err: transaction failed with error: MoveAbort", "function_name: Some(\"consume\") }, 1)")
		t.Log("Expected error: ", err)
	})
}

func mintLinkTokenOnSui(t *testing.T, e cldf.Environment, sourceChain uint64, amount uint64) sui_ops.OpTxResult[linkops.MintLinkTokenOutput] {
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	_, transferTokenOutput, err := commoncs.ApplyChangesets(t, e, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.MintLinkToken{}, sui_cs.MintLinkTokenConfig{
			ChainSelector:  sourceChain,
			TokenPackageId: state.SuiChains[sourceChain].LinkTokenAddress,
			TreasuryCapId:  state.SuiChains[sourceChain].LinkTokenTreasuryCapId,
			Amount:         amount,
		}),
	})
	require.NoError(t, err)
	rawOutputTransferToken := transferTokenOutput[0].Reports[0]
	outputMapTransferToken, ok := rawOutputTransferToken.Output.(sui_ops.OpTxResult[linkops.MintLinkTokenOutput])
	require.True(t, ok)
	return outputMapTransferToken
}

func Test_CCIPTokenTransfer_Sui2EVM_ManagedTokenPool_ThenCurseUncurse(t *testing.T) {
	e, sourceChain, destChain := testSetupTokenTransferSui2Evm(t)

	feeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000000)
	linkTokenOutput1 := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000)
	linkTokenOutput2 := mintLinkTokenOnSui(t, e.Env, sourceChain, 2000000000)
	linkTokenOutput3 := mintLinkTokenOnSui(t, e.Env, sourceChain, 1500000000)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// Receiver Address
	ccipReceiverAddress := state.Chains[destChain].Receiver.Address()

	// Token Pool setup on both SUI and EVM
	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndManagedTokenPoolDeploymentForSUI(e.Env, sourceChain, destChain, []testhelpers.TokenPoolRateLimiterConfig{
		{
			RemoteChainSelector: destChain,
			OutboundIsEnabled:   false,
			OutboundCapacity:    100000,
			OutboundRate:        100,
			InboundIsEnabled:    false,
			InboundCapacity:     100000,
			InboundRate:         100,
		},
	}) // SourceChain = SUI, destChain = EVM
	require.NoError(t, err)
	e.Env = updatedEnv

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       updatedEnv.BlockChains.EVMChains()[destChain].DeployerKey.From.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			FeeToken:       feeTokenOutput.Objects.MintedLinkTokenObjectId,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeManaged,
					Token:         linkTokenOutput1.Objects.MintedLinkTokenObjectId,
					Amount:        1000000000, // Send 1 LINK to EVM
				},
			},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(1e18),
				},
			},
		},
		{
			Name:           "Send token to Receiver",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       ccipReceiverAddress.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeManaged,
					Token:         linkTokenOutput2.Objects.MintedLinkTokenObjectId,
					Amount:        2000000000, // Send 2 LINK to EVM
				},
			},
			FeeToken: feeTokenOutput.Objects.MintedLinkTokenObjectId,
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(2e18),
				},
			},
		},
	}

	ctx := testhelpers.Context(t)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, updatedEnv, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		updatedEnv,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		updatedEnv,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, updatedEnv, expectedTokenBalances)

	suiState, err := sui_deployment.LoadOnchainStatesui(e.Env)
	require.NoError(t, err)

	suiChain := e.Env.BlockChains.SuiChains()[sourceChain]
	require.NotNil(t, suiChain)

	deps := getOpTxDeps(suiChain)

	// Convert evmChain selector to []byte
	selectorBytes := make([]byte, 16)
	binary.BigEndian.PutUint64(selectorBytes[8:], destChain)

	// curse destination chain
	_, err = operations.ExecuteOperation(e.Env.OperationsBundle, ccipops.RMNRemoteCurseOp, deps, ccipops.RMNRemoteCurseInput{
		CCIPPackageId:    suiState[sourceChain].CCIPAddress,
		StateObjectId:    suiState[sourceChain].CCIPObjectRef,
		OwnerCapObjectId: suiState[sourceChain].CCIPOwnerCapObjectId,
		Subject:          selectorBytes,
	})
	require.NoError(t, err)

	t.Run("Destination chain is cursed - should fail", func(t *testing.T) {
		msg := testhelpers.SuiSendRequest{
			Receiver: common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32),
			Data:     []byte("Hello, World!"),
			FeeToken: feeTokenOutput.Objects.MintedLinkTokenObjectId,
			TokenAmounts: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeManaged,
					Token:         linkTokenOutput3.Objects.MintedLinkTokenObjectId,
					Amount:        1500000000,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertSuiSourceRevertExpectedError(t, err, "failed to execute ccip_send with err: transaction failed with error: MoveAbort", "function_name: Some(\"validate_lock_or_burn\") }, 3)")
		t.Log("Expected error: ", err)
	})

	// uncurse destination chain
	_, err = operations.ExecuteOperation(e.Env.OperationsBundle, ccipops.RMNRemoteUncurseOp, deps, ccipops.RMNRemoteUncurseInput{
		CCIPPackageId:    suiState[sourceChain].CCIPAddress,
		StateObjectId:    suiState[sourceChain].CCIPObjectRef,
		OwnerCapObjectId: suiState[sourceChain].CCIPOwnerCapObjectId,
		Subject:          selectorBytes,
	})
	require.NoError(t, err)

	tcs = []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to EOA after uncursing",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       updatedEnv.BlockChains.EVMChains()[destChain].DeployerKey.From.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			FeeToken:       feeTokenOutput.Objects.MintedLinkTokenObjectId,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeManaged,
					Token:         linkTokenOutput3.Objects.MintedLinkTokenObjectId,
					Amount:        1500000000, // Send 1.5 LINK to EVM
				},
			},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(2.5e18), // 1e18 from first transfer + 1.5e18 from second transfer
				},
			},
		},
	}

	// reload state
	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	ctx = testhelpers.Context(t)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances = testhelpers.TransferMultiple(ctx, t, updatedEnv, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		updatedEnv,
		state,
		startBlocks,
		false,
		expectedSeqNums,
	)
	require.NoError(t, err)

	execStates = testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		updatedEnv,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, updatedEnv, expectedTokenBalances)
}

// func Test_CCIPTokenTransfer_EVM2SUI_ManagedTokenPool_NoRateLimit(t *testing.T) {
// 	e, sourceChain, destChain, deployerSourceChain, suiTokenBytes, suiAddr := testSetupHelperEvm2Sui(t)

// 	state, err := stateview.LoadOnchainState(e.Env)
// 	require.NoError(t, err)

// 	// Token Pool setup on both SUI and EVM
// 	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndManagedTokenPoolDeploymentForSUI(e.Env, destChain, sourceChain, []testhelpers.TokenPoolRateLimiterConfig{
// 		{
// 			RemoteChainSelector: sourceChain,
// 			OutboundIsEnabled:   false,
// 			OutboundCapacity:    100000,
// 			OutboundRate:        100,
// 			InboundIsEnabled:    false,
// 			InboundCapacity:     100000,
// 			InboundRate:         100,
// 		},
// 	}) // sourceChain=EVM, destChain=SUI
// 	require.NoError(t, err)

// 	// update env to include deployed contracts
// 	e.Env = updatedEnv

// 	testhelpers.MintAndAllow(
// 		t,
// 		e.Env,
// 		state,
// 		map[uint64][]testhelpers.MintTokenInfo{
// 			sourceChain: {
// 				testhelpers.NewMintTokenInfo(deployerSourceChain, evmToken),
// 			},
// 		},
// 	)

// 	emptyReceiver := hexutil.MustDecode(
// 		"0x0000000000000000000000000000000000000000000000000000000000000000", // receiver packageID
// 	)

// 	tcs := []testhelpers.TestTransferRequest{
// 		{
// 			Name:             "Send token to EOA",
// 			SourceChain:      sourceChain,
// 			DestChain:        destChain,
// 			Receiver:         emptyReceiver,
// 			TokenReceiverATA: suiAddr[:], // tokenReceiver extracted from extraArgs (the address that actually gets the token)
// 			ExpectedStatus:   testhelpers.EXECUTION_STATE_SUCCESS,
// 			Tokens: []router.ClientEVMTokenAmount{
// 				{
// 					Token:  evmToken.Address(),
// 					Amount: big.NewInt(1e18),
// 				},
// 			},
// 			ExtraArgs: testhelpers.MakeSuiExtraArgs(0, true, [][32]byte{}, suiAddr),
// 			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
// 				{
// 					Token:  suiTokenBytes,
// 					Amount: big.NewInt(1e9),
// 				},
// 			},
// 		},
// 	}

// 	ctx := testhelpers.Context(t)
// 	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, updatedEnv, state, tcs)

// 	err = testhelpers.ConfirmMultipleCommits(
// 		t,
// 		updatedEnv,
// 		state,
// 		startBlocks,
// 		false,
// 		expectedSeqNums,
// 	)
// 	require.NoError(t, err)

// 	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
// 		t,
// 		updatedEnv,
// 		state,
// 		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
// 		startBlocks,
// 	)
// 	require.Equal(t, expectedExecutionStates, execStates)

// 	testhelpers.WaitForTokenBalances(ctx, t, updatedEnv, expectedTokenBalances)
// }

func Test_CCIPTokenTransfer_SUI2EVM_ManagedTokenPool_WithRateLimit(t *testing.T) {
	e, sourceChain, destChain := testSetupTokenTransferSui2Evm(t)

	feeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000000)
	linkTokenOutput1 := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000)
	linkTokenOutput2 := mintLinkTokenOnSui(t, e.Env, sourceChain, 20000000000)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	ccipReceiverAddress := state.Chains[destChain].Receiver.Address()

	// Token Pool setup on both SUI and EVM
	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndManagedTokenPoolDeploymentForSUI(e.Env, sourceChain, destChain, []testhelpers.TokenPoolRateLimiterConfig{
		{
			RemoteChainSelector: destChain,
			OutboundIsEnabled:   true,
			OutboundCapacity:    100000,
			OutboundRate:        100,
			InboundIsEnabled:    true,
			InboundCapacity:     2000000000,
			InboundRate:         100000,
		},
	}) // sourceChain=SUI, destChain=EVM
	require.NoError(t, err)

	// update env to include deployed contracts
	e.Env = updatedEnv

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to EOA, within rate limit",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       updatedEnv.BlockChains.EVMChains()[destChain].DeployerKey.From.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			FeeToken:       feeTokenOutput.Objects.MintedLinkTokenObjectId,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeManaged,
					Token:         linkTokenOutput1.Objects.MintedLinkTokenObjectId,
					Amount:        1000000000, // Send 1 LINK to EVM
				},
			},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(1e18),
				},
			},
		},
	}

	ctx := testhelpers.Context(t)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)

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

	t.Run("Send tokens exceeding Sui's outbound rate limit - should fail", func(t *testing.T) {
		msg := testhelpers.SuiSendRequest{
			Receiver: common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32),
			Data:     []byte("Hello, World!"),
			FeeToken: feeTokenOutput.Objects.MintedLinkTokenObjectId,
			TokenAmounts: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeManaged,
					Token:         linkTokenOutput2.Objects.MintedLinkTokenObjectId,
					Amount:        20000000000,
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertSuiSourceRevertExpectedError(t, err, "failed to execute ccip_send with err: transaction failed with error: MoveAbort", "function_name: Some(\"validate_lock_or_burn\") }, 3)")
		t.Log("Expected error: ", err)
	})
}

// func Test_CCIPTokenTransfer_EVM2SUI_ManagedTokenPool_WithRateLimit(t *testing.T) {
// 	e, sourceChain, destChain, deployerSourceChain, suiTokenBytes, suiAddr := testSetupHelperEvm2Sui(t)

// 	// Token Pool setup on both SUI and EVM
// 	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndManagedTokenPoolDeploymentForSUI(e.Env, destChain, sourceChain, []testhelpers.TokenPoolRateLimiterConfig{
// 		{
// 			RemoteChainSelector: sourceChain,
// 			OutboundIsEnabled:   true,
// 			OutboundCapacity:    100000,
// 			OutboundRate:        100,
// 			InboundIsEnabled:    true,
// 			InboundCapacity:     2000000000,
// 			InboundRate:         100000,
// 		},
// 	}) // sourceChain=EVM, destChain=SUI
// 	require.NoError(t, err)

// 	// update env to include deployed contracts
// 	e.Env = updatedEnv

// 	state, err := stateview.LoadOnchainState(e.Env)
// 	require.NoError(t, err)

// 	testhelpers.MintAndAllow(
// 		t,
// 		e.Env,
// 		state,
// 		map[uint64][]testhelpers.MintTokenInfo{
// 			sourceChain: {
// 				testhelpers.NewMintTokenInfo(deployerSourceChain, evmToken),
// 			},
// 		},
// 	)

// 	emptyReceiver := hexutil.MustDecode(
// 		"0x0000000000000000000000000000000000000000000000000000000000000000", // receiver packageID
// 	)

// 	tcs := []testhelpers.TestTransferRequest{
// 		{
// 			Name:             "Send token to EOA",
// 			SourceChain:      sourceChain,
// 			DestChain:        destChain,
// 			Receiver:         emptyReceiver,
// 			TokenReceiverATA: suiAddr[:], // tokenReceiver extracted from extraArgs (the address that actually gets the token)
// 			ExpectedStatus:   testhelpers.EXECUTION_STATE_SUCCESS,
// 			Tokens: []router.ClientEVMTokenAmount{
// 				{
// 					Token:  evmToken.Address(),
// 					Amount: big.NewInt(1e18),
// 				},
// 			},
// 			ExtraArgs: testhelpers.MakeSuiExtraArgs(0, true, [][32]byte{}, suiAddr),
// 			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
// 				{
// 					Token:  suiTokenBytes,
// 					Amount: big.NewInt(1e9),
// 				},
// 			},
// 		},
// 	}

// 	ctx := testhelpers.Context(t)
// 	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)

// 	err = testhelpers.ConfirmMultipleCommits(
// 		t,
// 		e.Env,
// 		state,
// 		startBlocks,
// 		false,
// 		expectedSeqNums,
// 	)
// 	require.NoError(t, err)

// 	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
// 		t,
// 		e.Env,
// 		state,
// 		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
// 		startBlocks,
// 	)
// 	require.Equal(t, expectedExecutionStates, execStates)

// 	testhelpers.WaitForTokenBalances(ctx, t, e.Env, expectedTokenBalances)

// 	t.Run("Send tokens exceeding Sui's inbound rate limit - should fail", func(t *testing.T) {
// 		msg := router.ClientEVM2AnyMessage{
// 			FeeToken:  evmToken.Address(),
// 			Receiver:  emptyReceiver,
// 			ExtraArgs: testhelpers.MakeSuiExtraArgs(0, true, [][32]byte{}, suiAddr),
// 			TokenAmounts: []router.ClientEVMTokenAmount{
// 				{
// 					Token:  evmToken.Address(),
// 					Amount: big.NewInt(5e18), // send 5 LINK
// 				},
// 			}}

// 		baseOpts := []ccipclient.SendReqOpts{
// 			ccipclient.WithSourceChain(sourceChain),
// 			ccipclient.WithDestChain(destChain),
// 			ccipclient.WithTestRouter(false),
// 			ccipclient.WithMessage(msg),
// 		}

// 		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
// 		require.Error(t, err)
// 		require.Contains(t, err.Error(), "execution reverted")
// 		t.Log("Expected error: ", err)
// 	})
// }

func Test_CCIPTokenTransfer_EVM2SUI_BurnMintTokenPool(t *testing.T) {
	e, sourceChain, destChain, deployerSourceChain, suiTokenBytes, suiAddr := testSetupHelperEvm2Sui(t)

	// Token Pool setup on both SUI and EVM
	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndBurnMintTokenPoolDeploymentForSUI(e.Env, destChain, sourceChain, []testhelpers.TokenPoolRateLimiterConfig{
		{
			RemoteChainSelector: sourceChain,
			OutboundIsEnabled:   false,
			OutboundCapacity:    100000,
			OutboundRate:        100,
			InboundIsEnabled:    false,
			InboundCapacity:     100000,
			InboundRate:         100,
		},
	}) // sourceChain=EVM, destChain=SUI
	require.NoError(t, err)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	// update env to include deployed contracts
	e.Env = updatedEnv

	testhelpers.MintAndAllow(
		t,
		e.Env,
		state,
		map[uint64][]testhelpers.MintTokenInfo{
			sourceChain: {
				testhelpers.NewMintTokenInfo(deployerSourceChain, evmToken),
			},
		},
	)

	// Deploy SUI Receiver
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

	// register the receiver
	_, _, err = commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.RegisterDummyReceiver{}, sui_cs.RegisterDummyReceiverConfig{
			SuiChainSelector:       destChain,
			OwnerCapObjectId:       outputMap.Objects.OwnerCapObjectId,
			CCIPObjectRefObjectId:  state.SuiChains[destChain].CCIPObjectRef,
			DummyReceiverPackageId: outputMap.PackageId,
		}),
	})
	require.NoError(t, err)

	receiverByte := receiverByteDecoded

	var clockObj [32]byte
	copy(clockObj[:], hexutil.MustDecode(
		"0x0000000000000000000000000000000000000000000000000000000000000006",
	))

	var stateObj [32]byte
	copy(stateObj[:], hexutil.MustDecode(
		outputMap.Objects.CCIPReceiverStateObjectId,
	))

	receiverObjectIDs := [][32]byte{clockObj, stateObj}

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:             "Send token to EOA",
			SourceChain:      sourceChain,
			DestChain:        destChain,
			Receiver:         receiverByte, // receiver contract pkgId
			TokenReceiverATA: suiAddr[:],   // tokenReceiver extracted from extraArgs (the address that actually gets the token)
			ExpectedStatus:   testhelpers.EXECUTION_STATE_SUCCESS,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e18),
				},
			},
			ExtraArgs: testhelpers.MakeSuiExtraArgs(10000000, true, receiverObjectIDs, suiAddr),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  suiTokenBytes,
					Amount: big.NewInt(1e9),
				},
			},
		},
	}

	ctx := testhelpers.Context(t)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)

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

	callOpts := &bind.CallOpts{Context: ctx}
	srcFeeQuoterDestChainConfig, err := state.Chains[sourceChain].FeeQuoter.GetDestChainConfig(callOpts, destChain)
	require.NoError(t, err, "Failed to get destination chain fee quoter config")

	t.Run("Send token to CCIP Receiver setting gas above max gas allowed - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			Receiver:  receiverByte,
			Data:      []byte("Hello, World!"),
			FeeToken:  evmToken.Address(),
			ExtraArgs: testhelpers.MakeSuiExtraArgs(uint64(srcFeeQuoterDestChainConfig.MaxPerMsgGasLimit+1), true, receiverObjectIDs, stateObj),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e8),
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
		t.Log("Expected error: ", err)
	})

	t.Run("Send multiple token - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			Receiver:  receiverByte,
			Data:      []byte("Hello, World!"),
			FeeToken:  evmToken.Address(),
			ExtraArgs: testhelpers.MakeSuiExtraArgs(10000000, true, receiverObjectIDs, stateObj),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1),
				},
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1),
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
		t.Log("Expected error: ", err)
	})

	t.Run("Send invalid token to CCIP Receiver - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			Receiver:  receiverByte,
			Data:      []byte("Hello, World!"),
			FeeToken:  evmToken.Address(),
			ExtraArgs: testhelpers.MakeSuiExtraArgs(10000000, true, receiverObjectIDs, stateObj),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  common.HexToAddress("0x0000000000000000000000000000000000000000"), // Invalid token
					Amount: big.NewInt(1e8),
				},
			}}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")
		t.Log("Expected error: ", err)
	})
}

// func Test_CCIPPureTokenTransfer_EVM2SUI_BurnMintTokenPool(t *testing.T) {
// 	e, sourceChain, destChain, deployerSourceChain, suiTokenBytes, suiAddr := testSetupHelperEvm2Sui(t)

// 	// Token Pool setup on both SUI and EVM
// 	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndBurnMintTokenPoolDeploymentForSUI(e.Env, destChain, sourceChain, []testhelpers.TokenPoolRateLimiterConfig{
// 		{
// 			RemoteChainSelector: sourceChain,
// 			OutboundIsEnabled:   false,
// 			OutboundCapacity:    100000,
// 			OutboundRate:        100,
// 			InboundIsEnabled:    false,
// 			InboundCapacity:     100000,
// 			InboundRate:         100,
// 		},
// 	}) // sourceChain=EVM, destChain=SUI
// 	require.NoError(t, err)

// 	state, err := stateview.LoadOnchainState(e.Env)
// 	require.NoError(t, err)

// 	// update env to include deployed contracts
// 	e.Env = updatedEnv

// 	testhelpers.MintAndAllow(
// 		t,
// 		e.Env,
// 		state,
// 		map[uint64][]testhelpers.MintTokenInfo{
// 			sourceChain: {
// 				testhelpers.NewMintTokenInfo(deployerSourceChain, evmToken),
// 			},
// 		},
// 	)

// 	emptyReceiver := hexutil.MustDecode(
// 		"0x0000000000000000000000000000000000000000000000000000000000000000", // receiver packageID
// 	)

// 	tcs := []testhelpers.TestTransferRequest{
// 		// Pure token transfer
// 		// ReceiverObjectIds = empty
// 		// token.Receiver = non empty (maybe EOA or object)
// 		// message.Receiver = empty
// 		// don't need extraArgs gasLimit, can be set to 0
// 		{
// 			Name:             "Send token to EOA with - Pure Token Transfer",
// 			SourceChain:      sourceChain,
// 			DestChain:        destChain,
// 			Data:             []byte{},
// 			Receiver:         emptyReceiver, // empty Receiver
// 			TokenReceiverATA: suiAddr[:],    // tokenReceiver extracted from extraArgs (the address that actually gets the token)
// 			ExpectedStatus:   testhelpers.EXECUTION_STATE_SUCCESS,
// 			Tokens: []router.ClientEVMTokenAmount{
// 				{
// 					Token:  evmToken.Address(),
// 					Amount: big.NewInt(1e18),
// 				},
// 			},
// 			ExtraArgs: testhelpers.MakeSuiExtraArgs(0, true, [][32]byte{}, suiAddr),
// 			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
// 				{
// 					Token:  suiTokenBytes,
// 					Amount: big.NewInt(1e9),
// 				},
// 			},
// 		},
// 	}

// 	ctx := testhelpers.Context(t)
// 	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)

// 	err = testhelpers.ConfirmMultipleCommits(
// 		t,
// 		e.Env,
// 		state,
// 		startBlocks,
// 		false,
// 		expectedSeqNums,
// 	)
// 	require.NoError(t, err)

// 	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
// 		t,
// 		e.Env,
// 		state,
// 		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
// 		startBlocks,
// 	)
// 	require.Equal(t, expectedExecutionStates, execStates)

// 	testhelpers.WaitForTokenBalances(ctx, t, e.Env, expectedTokenBalances)
// }

// func Test_CCIPProgrammableTokenTransfer_EVM2SUI_BurnMintTokenPool(t *testing.T) {
// 	e, sourceChain, destChain, deployerSourceChain, _, _ := testSetupHelperEvm2Sui(t)

// 	// Token Pool setup on both SUI and EVM
// 	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndBurnMintTokenPoolDeploymentForSUI(e.Env, destChain, sourceChain, []testhelpers.TokenPoolRateLimiterConfig{
// 		{
// 			RemoteChainSelector: sourceChain,
// 			OutboundIsEnabled:   false,
// 			OutboundCapacity:    100000,
// 			OutboundRate:        100,
// 			InboundIsEnabled:    false,
// 			InboundCapacity:     100000,
// 			InboundRate:         100,
// 		},
// 	}) // sourceChain=EVM, destChain=SUI
// 	require.NoError(t, err)

// 	state, err := stateview.LoadOnchainState(e.Env)
// 	require.NoError(t, err)

// 	// update env to include deployed contracts
// 	e.Env = updatedEnv

// 	testhelpers.MintAndAllow(
// 		t,
// 		e.Env,
// 		state,
// 		map[uint64][]testhelpers.MintTokenInfo{
// 			sourceChain: {
// 				testhelpers.NewMintTokenInfo(deployerSourceChain, evmToken),
// 			},
// 		},
// 	)

// 	// Deploy SUI Receiver
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

// 	// register the receiver
// 	_, _, err = commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
// 		commoncs.Configure(sui_cs.RegisterDummyReceiver{}, sui_cs.RegisterDummyReceiverConfig{
// 			SuiChainSelector:       destChain,
// 			OwnerCapObjectId:       outputMap.Objects.OwnerCapObjectId,
// 			CCIPObjectRefObjectId:  state.SuiChains[destChain].CCIPObjectRef,
// 			DummyReceiverPackageId: outputMap.PackageId,
// 		}),
// 	})
// 	require.NoError(t, err)

// 	receiverByte := receiverByteDecoded

// 	var clockObj [32]byte
// 	copy(clockObj[:], hexutil.MustDecode(
// 		"0x0000000000000000000000000000000000000000000000000000000000000006",
// 	))

// 	var stateObj [32]byte
// 	copy(stateObj[:], hexutil.MustDecode(
// 		outputMap.Objects.CCIPReceiverStateObjectId,
// 	))

// 	receiverObjectIDs := [][32]byte{clockObj, stateObj}

// 	tcs := []testhelpers.TestTransferRequest{
// 		// Programmable token transfer
// 		// can be thought of as two separate paths tokenPool release/mint + message ccip_receive
// 		// receiverObjectIds = non empty (with clock & receiverStateValue)
// 		// token.Receiver = non empty(maybe EOA or object)
// 		// message.Receiver = receiverPackageId
// 		// extraArgs gasLimit > 0
// 		{
// 			Name:             "Send token to an Object",
// 			SourceChain:      sourceChain,
// 			DestChain:        destChain,
// 			Data:             []byte("Hello Sui From EVM"),
// 			Receiver:         receiverByte, // receiver contract pkgId
// 			TokenReceiverATA: stateObj[:],  // tokenReceiver extracted from extraArgs (the object that actually gets the token)
// 			ExpectedStatus:   testhelpers.EXECUTION_STATE_SUCCESS,
// 			Tokens: []router.ClientEVMTokenAmount{
// 				{
// 					Token:  evmToken.Address(),
// 					Amount: big.NewInt(1e18),
// 				},
// 			},
// 			ExtraArgs:             testhelpers.MakeSuiExtraArgs(1000000, true, receiverObjectIDs, stateObj), // receiver is objectId this time
// 			ExpectedTokenBalances: []testhelpers.ExpectedBalance{},
// 		},
// 	}

// 	ctx := testhelpers.Context(t)
// 	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)

// 	err = testhelpers.ConfirmMultipleCommits(
// 		t,
// 		e.Env,
// 		state,
// 		startBlocks,
// 		false,
// 		expectedSeqNums,
// 	)
// 	require.NoError(t, err)

// 	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
// 		t,
// 		e.Env,
// 		state,
// 		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
// 		startBlocks,
// 	)
// 	require.Equal(t, expectedExecutionStates, execStates)

// 	testhelpers.WaitForTokenBalances(ctx, t, e.Env, expectedTokenBalances)
// }

// func Test_CCIPZeroGasLimitTokenTransfer_EVM2SUI_BurnMintTokenPool(t *testing.T) {
// 	e, sourceChain, destChain, deployerSourceChain, suiTokenBytes, suiAddr := testSetupHelperEvm2Sui(t)

// 	// Token Pool setup on both SUI and EVM
// 	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndBurnMintTokenPoolDeploymentForSUI(e.Env, destChain, sourceChain, []testhelpers.TokenPoolRateLimiterConfig{
// 		{
// 			RemoteChainSelector: sourceChain,
// 			OutboundIsEnabled:   false,
// 			OutboundCapacity:    100000,
// 			OutboundRate:        100,
// 			InboundIsEnabled:    false,
// 			InboundCapacity:     100000,
// 			InboundRate:         100,
// 		},
// 	}) // sourceChain=EVM, destChain=SUI
// 	require.NoError(t, err)

// 	state, err := stateview.LoadOnchainState(e.Env)
// 	require.NoError(t, err)

// 	// update env to include deployed contracts
// 	e.Env = updatedEnv

// 	testhelpers.MintAndAllow(
// 		t,
// 		e.Env,
// 		state,
// 		map[uint64][]testhelpers.MintTokenInfo{
// 			sourceChain: {
// 				testhelpers.NewMintTokenInfo(deployerSourceChain, evmToken),
// 			},
// 		},
// 	)

// 	// Deploy SUI Receiver
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

// 	// register the receiver
// 	_, _, err = commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
// 		commoncs.Configure(sui_cs.RegisterDummyReceiver{}, sui_cs.RegisterDummyReceiverConfig{
// 			SuiChainSelector:       destChain,
// 			OwnerCapObjectId:       outputMap.Objects.OwnerCapObjectId,
// 			CCIPObjectRefObjectId:  state.SuiChains[destChain].CCIPObjectRef,
// 			DummyReceiverPackageId: outputMap.PackageId,
// 		}),
// 	})
// 	require.NoError(t, err)

// 	receiverByte := receiverByteDecoded

// 	var clockObj [32]byte
// 	copy(clockObj[:], hexutil.MustDecode(
// 		"0x0000000000000000000000000000000000000000000000000000000000000006",
// 	))

// 	var stateObj [32]byte
// 	copy(stateObj[:], hexutil.MustDecode(
// 		outputMap.Objects.CCIPReceiverStateObjectId,
// 	))

// 	receiverObjectIDs := [][32]byte{clockObj, stateObj}

// 	tcs := []testhelpers.TestTransferRequest{
// 		{
// 			Name:             "Send token To EOA + include a receiver but keep gasLimit to 0",
// 			SourceChain:      sourceChain,
// 			DestChain:        destChain,
// 			Data:             []byte("Hello Sui From EVM"),
// 			Receiver:         receiverByte, // non empty Receiver
// 			TokenReceiverATA: suiAddr[:],   // tokenReceiver extracted from extraArgs (the address that actually gets the token)
// 			ExpectedStatus:   testhelpers.EXECUTION_STATE_SUCCESS,
// 			Tokens: []router.ClientEVMTokenAmount{
// 				{
// 					Token:  evmToken.Address(),
// 					Amount: big.NewInt(1e18),
// 				},
// 			},
// 			ExtraArgs: testhelpers.MakeSuiExtraArgs(0, true, receiverObjectIDs, suiAddr), // keep gasLimit to 0
// 			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
// 				{
// 					Token:  suiTokenBytes,
// 					Amount: big.NewInt(1e9),
// 				},
// 			},
// 		},
// 	}

// 	ctx := testhelpers.Context(t)
// 	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)

// 	err = testhelpers.ConfirmMultipleCommits(
// 		t,
// 		e.Env,
// 		state,
// 		startBlocks,
// 		false,
// 		expectedSeqNums,
// 	)
// 	require.NoError(t, err)

// 	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
// 		t,
// 		e.Env,
// 		state,
// 		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
// 		startBlocks,
// 	)
// 	require.Equal(t, expectedExecutionStates, execStates)

// 	testhelpers.WaitForTokenBalances(ctx, t, e.Env, expectedTokenBalances)
// }

func testSetupTokenTransferSui2Evm(t *testing.T) (e testhelpers.DeployedEnv, sourceChain uint64, destChain uint64) {
	e, _, _ = testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	sourceChain = suiChainSelectors[0]
	destChain = evmChainSelectors[0]

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Log("Source chain (Sui): ", sourceChain, "Dest chain (EVM): ", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	return e, sourceChain, destChain
}

func testSetupHelperEvm2Sui(t *testing.T) (e testhelpers.DeployedEnv, sourceChain uint64, destChain uint64, deployerSourceChain *bind.TransactOpts, suiTokenBytes []byte, suiAddr [32]byte) {
	e, _, _ = testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	sourceChain = evmChainSelectors[0]
	destChain = suiChainSelectors[0]

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Log("Source chain (EVM): ", sourceChain, "Dest chain (SUI): ", destChain)

	deployerSourceChain = e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey
	suiTokenHex := state.SuiChains[destChain].LinkTokenAddress
	suiTokenHex = strings.TrimPrefix(suiTokenHex, "0x")

	suiTokenBytes, err = hex.DecodeString(suiTokenHex)
	require.NoError(t, err)

	require.Len(t, suiTokenBytes, 32, "expected 32-byte sui address")

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	// get sui address in [32]bytes for extraArgs.TokenReceiver
	suiAddrStr, err := e.Env.BlockChains.SuiChains()[destChain].Signer.GetAddress()
	require.NoError(t, err)

	suiAddrStr = strings.TrimPrefix(suiAddrStr, "0x")

	addrBytes, err := hex.DecodeString(suiAddrStr)
	require.NoError(t, err)

	require.Len(t, addrBytes, 32, "expected 32-byte sui address")
	copy(suiAddr[:], addrBytes)

	return e, sourceChain, destChain, deployerSourceChain, suiTokenBytes, suiAddr
}

func getOpTxDeps(suiChain sui.Chain) sui_ops.OpTxDeps {
	return sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: suiChain.Signer,
		GetCallOpts: func() *suiBind.CallOpts {
			b := uint64(400_000_000)
			return &suiBind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
		SuiRPC: suiChain.URL,
	}
}

func assertSuiSourceRevertExpectedError(t *testing.T, err error, execRevertErrorMsg string, execRevertCauseErrorMsg string) {
	require.Error(t, err)
	fmt.Println("Error: ", err.Error())
	require.Contains(t, err.Error(), execRevertErrorMsg)
	require.Contains(t, err.Error(), execRevertCauseErrorMsg)
}
