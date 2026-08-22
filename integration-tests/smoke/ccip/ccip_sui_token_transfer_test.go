package ccip

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"
	suiBind "github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	module_dummy_receiver "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_dummy_receiver/ccip_dummy_receiver"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	codec "github.com/smartcontractkit/chainlink-sui/codec"
	sui_deployment "github.com/smartcontractkit/chainlink-sui/deployment"
	sui_cs "github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIPTokenTransfer_Sui2EVM_LockReleaseTokenPool_Plain(t *testing.T) {
	e, sourceChain, destChain := testSetupTokenTransferSui2Evm(t)

	feeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 100000000000)
	linkTokenOutput1 := mintLinkTokenOnSui(t, e.Env, sourceChain, 5000000000)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

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
	})
	require.NoError(t, err)
	e.Env = updatedEnv

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       updatedEnv.BlockChains.EVMChains()[destChain].DeployerKey.From.Bytes(),
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			FeeToken:       feeTokenOutput.Objects.MintedLinkTokenObjectId,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeLockRelease,
					Token:         linkTokenOutput1.Objects.MintedLinkTokenObjectId,
					Amount:        5000000000,
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
}

func Test_CCIPTokenTransfer_Sui2EVM_LockReleaseTokenPool_Revert(t *testing.T) {
	e, sourceChain, destChain := testSetupTokenTransferSui2Evm(t)

	feeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 100000000000)
	linkTokenOutput1 := mintLinkTokenOnSui(t, e.Env, sourceChain, 5000000000)
	linkTokenOutput3 := mintLinkTokenOnSui(t, e.Env, sourceChain, 1500000000)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	ccipReceiverAddress := state.Chains[destChain].Receiver.Address()

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
	})
	require.NoError(t, err)
	e.Env = updatedEnv

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       updatedEnv.BlockChains.EVMChains()[destChain].DeployerKey.From.Bytes(),
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			FeeToken:       feeTokenOutput.Objects.MintedLinkTokenObjectId,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeLockRelease,
					Token:         linkTokenOutput1.Objects.MintedLinkTokenObjectId,
					Amount:        5000000000,
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

	waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])

	suiState, err := sui_deployment.LoadOnchainStatesui(e.Env)
	require.NoError(t, err)

	suifeeQuoter, err := module_fee_quoter.NewFeeQuoter(suiState[sourceChain].CCIPAddress, e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)

	suiFeeQuoterDestChainConfig, err := suifeeQuoter.DevInspect().GetDestChainConfig(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
	}, suiBind.Object{Id: suiState[sourceChain].CCIPObjectRef}, destChain)
	require.NoError(t, err, "Failed to get destination chain config")

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Run("Send invalid token to CCIP Receiver - should fail", func(t *testing.T) {
		waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])
		msg := testhelpers.SuiSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32),
			Data:      []byte("Hello, World!"),
			FeeToken:  feeTokenOutput.Objects.MintedLinkTokenObjectId,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(suiFeeQuoterDestChainConfig.MaxPerMsgGasLimit)), false),
			TokenAmounts: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeLockRelease,
					Token:         "0x0",
					Amount:        1e9,
				},
			},
		}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertSuiSourceRevertExpectedError(t, err, "failed to resolve token coin object", "failed to resolve UnresolvedObject 0x0000000000000000000000000000000000000000000000000000000000000000")
		t.Log("Expected error: ", err)
	})

	t.Run("Send token to CCIP Receiver setting gas above max gas allowed - should fail", func(t *testing.T) {
		waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])
		msg := testhelpers.SuiSendRequest{
			Receiver:  common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32),
			Data:      []byte("Hello, World!"),
			FeeToken:  feeTokenOutput.Objects.MintedLinkTokenObjectId,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(int64(suiFeeQuoterDestChainConfig.MaxPerMsgGasLimit+10)), false),
			TokenAmounts: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeLockRelease,
					Token:         linkTokenOutput3.Objects.MintedLinkTokenObjectId,
					Amount:        1500000000,
				},
			},
		}

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

func Test_CCIPTokenTransfer_Sui2EVM_BurnMintTokenPool_Plain(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CCIP-11054")
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

	waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])

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
		waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])
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
			},
		}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}

		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		assertSuiSourceRevertExpectedError(t, err, "failed to resolve token coin object", "failed to resolve UnresolvedObject 0x0000000000000000000000000000000000000000000000000000000000000000")
		t.Log("Expected error: ", err)
	})

	t.Run("Send token to CCIP Receiver setting gas above max gas allowed - should fail", func(t *testing.T) {
		waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])
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
			},
		}

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
	linkTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 5000000000)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	ccipReceiverAddress := state.Chains[destChain].Receiver.Address()

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

	waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	suiState, err := sui_deployment.LoadOnchainStatesui(e.Env)
	require.NoError(t, err)

	suiChain := e.Env.BlockChains.SuiChains()[sourceChain]
	require.NotNil(t, suiChain)

	deps := getOpTxDeps(suiChain)

	// Ensure Sui state is fully consistent before curse operation
	waitForSuiRPCSyncSlow(t, e.Env.BlockChains.SuiChains()[sourceChain])

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
		waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])
		msg := testhelpers.SuiSendRequest{
			Receiver: common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32), // left-pad 20-byte address up to 32 bytes to make it compatible with evm
			Data:     []byte("Hello, World!"),
			FeeToken: feeTokenOutput.Objects.MintedLinkTokenObjectId,
			TokenAmounts: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         linkTokenOutput.Objects.MintedLinkTokenObjectId,
					Amount:        2000000000,
				},
			},
		}

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

	// Ensure Sui state is fully consistent before uncurse operation
	waitForSuiRPCSyncSlow(t, e.Env.BlockChains.SuiChains()[sourceChain])

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

	waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	evmDeployer := updatedEnv.BlockChains.EVMChains()[destChain].DeployerKey.From
	preBal := evmBurnMint677BalanceOf(t, updatedEnv, destChain, evmToken, evmDeployer)
	// Sui LINK 9 decimals → EVM 18 decimals: multiply Sui amount by 1e9 for minted wei on dest.
	transferWei := new(big.Int).Mul(big.NewInt(2000000000), big.NewInt(1_000_000_000))
	expectedEVMBal := new(big.Int).Add(preBal, transferWei)

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to EOA after uncursing",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       evmDeployer.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			FeeToken:       feeTokenOutput.Objects.MintedLinkTokenObjectId,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         linkTokenOutput.Objects.MintedLinkTokenObjectId,
					Amount:        2000000000, // Send 2 LINK to EVM
				},
			},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: expectedEVMBal,
				},
			},
		},
	}

	waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])

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
}

func Test_CCIPTokenTransfer_Sui2EVM_BurnMintTokenPool_WithAllowlist_DenylistedSender(t *testing.T) {
	e, sourceChain, destChain := testSetupTokenTransferSui2Evm(t)
	feeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000000)
	linkTokenOutput2 := mintLinkTokenOnSui(t, e.Env, sourceChain, 4000000000)

	updatedEnv, _, _, err := testhelpers.HandleTokenAndBurnMintTokenPoolDeploymentForSUI(e.Env, sourceChain, destChain, []testhelpers.TokenPoolRateLimiterConfig{
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

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	ccipReceiverAddress := state.Chains[destChain].Receiver.Address()

	suiChain := e.Env.BlockChains.SuiChains()[sourceChain]
	deps := getOpTxDeps(suiChain)

	// enable allowlist but not adding the current sender to the allowlist
	_, err = operations.ExecuteOperation(e.Env.OperationsBundle, burnminttokenpoolops.BurnMintTokenPoolSetAllowlistEnabledOp, deps, burnminttokenpoolops.BurnMintTokenPoolSetAllowlistEnabledInput{
		BurnMintPackageId: state.SuiChains[sourceChain].BnMTokenPools[testhelpers.TokenSymbolLINK].PackageID,
		StateObjectId:     state.SuiChains[sourceChain].BnMTokenPools[testhelpers.TokenSymbolLINK].StateObjectId,
		OwnerCap:          state.SuiChains[sourceChain].BnMTokenPools[testhelpers.TokenSymbolLINK].OwnerCapObjectId,
		CoinObjectTypeArg: state.SuiChains[sourceChain].LinkTokenAddress + "::link::LINK",
		Enabled:           true,
	})
	require.NoError(t, err)

	waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])
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
		},
	}

	baseOpts := []ccipclient.SendReqOpts{
		ccipclient.WithSourceChain(sourceChain),
		ccipclient.WithDestChain(destChain),
		ccipclient.WithTestRouter(false),
		ccipclient.WithMessage(msg),
	}

	_, err = testhelpers.SendRequest(e.Env, state, baseOpts...)
	assertSuiSourceRevertExpectedError(t, err, "failed to execute ccip_send with err: transaction failed with error: MoveAbort", "function_name: Some(\"validate_lock_or_burn\") }, 1)")
	t.Log("Expected error: ", err)
}

func Test_CCIPTokenTransfer_Sui2EVM_BurnMintTokenPool_WithAllowlist_AfterSignerAdded(t *testing.T) {
	e, sourceChain, destChain := testSetupTokenTransferSui2Evm(t)
	feeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000000)
	linkTokenOutput2 := mintLinkTokenOnSui(t, e.Env, sourceChain, 4000000000)

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

	waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])

	e.Env = updatedEnv

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	ccipReceiverAddress := state.Chains[destChain].Receiver.Address()

	suiChain := e.Env.BlockChains.SuiChains()[sourceChain]
	deps := getOpTxDeps(suiChain)

	_, err = operations.ExecuteOperation(e.Env.OperationsBundle, burnminttokenpoolops.BurnMintTokenPoolSetAllowlistEnabledOp, deps, burnminttokenpoolops.BurnMintTokenPoolSetAllowlistEnabledInput{
		BurnMintPackageId: state.SuiChains[sourceChain].BnMTokenPools[testhelpers.TokenSymbolLINK].PackageID,
		StateObjectId:     state.SuiChains[sourceChain].BnMTokenPools[testhelpers.TokenSymbolLINK].StateObjectId,
		OwnerCap:          state.SuiChains[sourceChain].BnMTokenPools[testhelpers.TokenSymbolLINK].OwnerCapObjectId,
		CoinObjectTypeArg: state.SuiChains[sourceChain].LinkTokenAddress + "::link::LINK",
		Enabled:           true,
	})
	require.NoError(t, err)

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

	waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	preRecvBal := evmBurnMint677BalanceOf(t, updatedEnv, destChain, evmToken, ccipReceiverAddress)
	transferWei := new(big.Int).Mul(big.NewInt(1500000000), big.NewInt(1_000_000_000))
	expectedRecvBal := new(big.Int).Add(preRecvBal, transferWei)

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to Receiver after signer allowlisted",
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
					Amount: expectedRecvBal,
				},
			},
		},
	}

	waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])

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
}

func Test_CCIPTokenTransfer_Sui2EVM_BurnMintTokenPool_WithRateLimit(t *testing.T) {
	e, sourceChain, destChain := testSetupTokenTransferSui2Evm(t)
	feeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000000)
	linkTokenOutput3 := mintLinkTokenOnSui(t, e.Env, sourceChain, 999999999999)

	updatedEnv, _, _, err := testhelpers.HandleTokenAndBurnMintTokenPoolDeploymentForSUI(e.Env, sourceChain, destChain, []testhelpers.TokenPoolRateLimiterConfig{
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

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	ccipReceiverAddress := state.Chains[destChain].Receiver.Address()

	t.Run("Send token above Sui's outbound rate limit - should fail", func(t *testing.T) {
		waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])
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
			},
		}

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

func mintLinkTokenOnSui(t *testing.T, e cldf.Environment, sourceChain, amount uint64) sui_ops.OpTxResult[linkops.MintLinkTokenOutput] {
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

// evmBurnMint677BalanceOf reads an ERC-677 balance for WaitForTokenBalances expectations.
// CCIP mints on top of any balance already held by the account (e.g. from earlier deliveries in the same test).
func evmBurnMint677BalanceOf(t *testing.T, env cldf.Environment, destChain uint64, token *burn_mint_erc677.BurnMintERC677, account common.Address) *big.Int {
	t.Helper()
	ctx := testhelpers.Context(t)
	bal, err := token.BalanceOf(&bind.CallOpts{Context: ctx}, account)
	require.NoError(t, err)
	return new(big.Int).Set(bal)
}

func Test_CCIPTokenTransfer_Sui2EVM_ManagedTokenPool_ThenCurseUncurse(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CCIP-11054")
	e, sourceChain, destChain := testSetupTokenTransferSui2Evm(t)

	feeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000000)
	linkTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 5000000000)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	ccipReceiverAddress := state.Chains[destChain].Receiver.Address()

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

	waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	suiState, err := sui_deployment.LoadOnchainStatesui(e.Env)
	require.NoError(t, err)

	suiChain := e.Env.BlockChains.SuiChains()[sourceChain]
	require.NotNil(t, suiChain)

	deps := getOpTxDeps(suiChain)

	// Convert evmChain selector to []byte
	selectorBytes := make([]byte, 16)
	binary.BigEndian.PutUint64(selectorBytes[8:], destChain)

	// Ensure Sui state is fully consistent before curse operation
	waitForSuiRPCSyncSlow(t, e.Env.BlockChains.SuiChains()[sourceChain])

	// curse destination chain
	_, err = operations.ExecuteOperation(e.Env.OperationsBundle, ccipops.RMNRemoteCurseOp, deps, ccipops.RMNRemoteCurseInput{
		CCIPPackageId:    suiState[sourceChain].CCIPAddress,
		StateObjectId:    suiState[sourceChain].CCIPObjectRef,
		OwnerCapObjectId: suiState[sourceChain].CCIPOwnerCapObjectId,
		Subject:          selectorBytes,
	})
	require.NoError(t, err)

	t.Run("Destination chain is cursed - should fail", func(t *testing.T) {
		waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])
		msg := testhelpers.SuiSendRequest{
			Receiver: common.LeftPadBytes(ccipReceiverAddress.Bytes(), 32),
			Data:     []byte("Hello, World!"),
			FeeToken: feeTokenOutput.Objects.MintedLinkTokenObjectId,
			TokenAmounts: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeManaged,
					Token:         linkTokenOutput.Objects.MintedLinkTokenObjectId,
					Amount:        1500000000,
				},
			},
		}

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

	// Ensure Sui state is fully consistent before uncurse operation
	waitForSuiRPCSyncSlow(t, e.Env.BlockChains.SuiChains()[sourceChain])

	// uncurse destination chain
	_, err = operations.ExecuteOperation(e.Env.OperationsBundle, ccipops.RMNRemoteUncurseOp, deps, ccipops.RMNRemoteUncurseInput{
		CCIPPackageId:    suiState[sourceChain].CCIPAddress,
		StateObjectId:    suiState[sourceChain].CCIPObjectRef,
		OwnerCapObjectId: suiState[sourceChain].CCIPOwnerCapObjectId,
		Subject:          selectorBytes,
	})
	require.NoError(t, err)

	waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	evmDeployer := updatedEnv.BlockChains.EVMChains()[destChain].DeployerKey.From
	preBal := evmBurnMint677BalanceOf(t, updatedEnv, destChain, evmToken, evmDeployer)
	transferWei := new(big.Int).Mul(big.NewInt(1500000000), big.NewInt(1_000_000_000))
	expectedEVMBal := new(big.Int).Add(preBal, transferWei)

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to EOA after uncursing",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       evmDeployer.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			FeeToken:       feeTokenOutput.Objects.MintedLinkTokenObjectId,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeManaged,
					Token:         linkTokenOutput.Objects.MintedLinkTokenObjectId,
					Amount:        1500000000, // Send 1.5 LINK to EVM
				},
			},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: expectedEVMBal,
				},
			},
		},
	}

	waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])

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

	waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])

	testhelpers.WaitForTokenBalances(ctx, t, updatedEnv, expectedTokenBalances)
}

func Test_CCIPTokenTransfer_EVM2Sui_ManagedTokenPool_NoRateLimit(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CCIP-11054")

	e, sourceChain, destChain, deployerSourceChain, suiTokenBytes, suiAddr := testSetupHelperEvm2Sui(t)

	// Token Pool setup on both SUI and EVM
	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndManagedTokenPoolDeploymentForSUI(e.Env, destChain, sourceChain, []testhelpers.TokenPoolRateLimiterConfig{
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

	// update env to include deployed contracts
	e.Env = updatedEnv

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

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

	emptyReceiver := hexutil.MustDecode(
		"0x0000000000000000000000000000000000000000000000000000000000000000", // receiver packageID
	)

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:             "Send token to EOA",
			SourceChain:      sourceChain,
			DestChain:        destChain,
			Receiver:         emptyReceiver,
			TokenReceiverATA: suiAddr[:], // tokenReceiver extracted from extraArgs (the address that actually gets the token)
			ExpectedStatus:   testhelpers.EXECUTION_STATE_SUCCESS,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e18),
				},
			},
			ExtraArgs: testhelpers.MakeSuiExtraArgs(0, true, [][32]byte{}, suiAddr),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  suiTokenBytes,
					Amount: big.NewInt(1e9),
				},
			},
		},
	}

	ctx := testhelpers.Context(t)
	prepareEvm2SuiTransferLane(t, e, state, sourceChain, destChain)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)
	replayEvm2SuiTransferLane(t, e, sourceChain, destChain)

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
}

func Test_CCIPTokenTransfer_Sui2EVM_ManagedTokenPool_WithRateLimit(t *testing.T) {
	e, sourceChain, destChain := testSetupTokenTransferSui2Evm(t)

	feeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000000)
	linkTokenOutput2 := mintLinkTokenOnSui(t, e.Env, sourceChain, 20000000000)

	updatedEnv, _, _, err := testhelpers.HandleTokenAndManagedTokenPoolDeploymentForSUI(e.Env, sourceChain, destChain, []testhelpers.TokenPoolRateLimiterConfig{
		{
			RemoteChainSelector: destChain,
			OutboundIsEnabled:   true,
			OutboundCapacity:    10000000000, // 10 LINK; a single 20 LINK send exceeds outbound capacity
			OutboundRate:        1000000000,
			InboundIsEnabled:    true,
			InboundCapacity:     2000000000,
			InboundRate:         100000,
		},
	}) // sourceChain=SUI, destChain=EVM
	require.NoError(t, err)

	e.Env = updatedEnv

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	ccipReceiverAddress := state.Chains[destChain].Receiver.Address()

	t.Run("Send tokens exceeding Sui's outbound rate limit - should fail", func(t *testing.T) {
		waitForSuiRPCSync(t, e.Env.BlockChains.SuiChains()[sourceChain])
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
			},
		}

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

func Test_CCIPTokenTransfer_EVM2Sui_ManagedTokenPool_WithRateLimit(t *testing.T) {
	e, sourceChain, destChain, deployerSourceChain, _, suiAddr := testSetupHelperEvm2Sui(t)

	// Token Pool setup on both SUI and EVM
	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndManagedTokenPoolDeploymentForSUI(e.Env, destChain, sourceChain, []testhelpers.TokenPoolRateLimiterConfig{
		{
			RemoteChainSelector: sourceChain,
			OutboundIsEnabled:   true,
			OutboundCapacity:    100000,
			OutboundRate:        100,
			InboundIsEnabled:    true,
			InboundCapacity:     2000000000,
			InboundRate:         100000,
		},
	}) // sourceChain=EVM, destChain=SUI
	require.NoError(t, err)

	// update env to include deployed contracts
	e.Env = updatedEnv

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

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

	emptyReceiver := hexutil.MustDecode(
		"0x0000000000000000000000000000000000000000000000000000000000000000", // receiver packageID
	)

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Run("Send tokens exceeding Sui's inbound rate limit - should fail", func(t *testing.T) {
		msg := router.ClientEVM2AnyMessage{
			FeeToken:  evmToken.Address(),
			Receiver:  emptyReceiver,
			ExtraArgs: testhelpers.MakeSuiExtraArgs(0, true, [][32]byte{}, suiAddr),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(5e18), // send 5 LINK
				},
			},
		}

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

func Test_CCIPTokenTransfer_EVM2Sui_BurnMintTokenPool(t *testing.T) {
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

	// update env to include deployed contracts
	e.Env = updatedEnv

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

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
			ExtraArgs: testhelpers.MakeSuiExtraArgs(1000000, true, receiverObjectIDs, suiAddr),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  suiTokenBytes,
					Amount: big.NewInt(1e9),
				},
			},
		},
	}

	ctx := testhelpers.Context(t)
	prepareEvm2SuiTransferLane(t, e, state, sourceChain, destChain)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)
	replayEvm2SuiTransferLane(t, e, sourceChain, destChain)

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
			},
		}

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
			ExtraArgs: testhelpers.MakeSuiExtraArgs(1000000, true, receiverObjectIDs, stateObj),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1),
				},
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1),
				},
			},
		}

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
			ExtraArgs: testhelpers.MakeSuiExtraArgs(1000000, true, receiverObjectIDs, stateObj),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  common.HexToAddress("0x0000000000000000000000000000000000000000"), // Invalid token
					Amount: big.NewInt(1e8),
				},
			},
		}

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

// Sui offramp must surface the destination (local) token amount to CCIP receivers, not the unconverted source-chain amount.
//
// Flow: send 1e18 of an 18-dec EVM token + a message to a registered Sui dummy receiver (the 9-dec
// Sui counterpart). The pool converts 18->9 and mints 1e9. The dummy receiver stores the amount it
// reads via client::get_token_and_amount into CCIPReceiverState.dest_token_amounts.
func Test_CCIP_EVM2Sui_DestTokenAmount_ReportsLocalAmount(t *testing.T) {
	e, sourceChain, destChain, deployerSourceChain, suiTokenBytes, suiAddr := testSetupHelperEvm2Sui(t)

	// Token pool setup on both SUI (9-dec LINK, burn-mint) and EVM (18-dec token).
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
	e.Env = updatedEnv

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

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

	// Deploy + register the Sui dummy receiver (it stores the receiver-visible amount).
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

	// Token + message to the dummy receiver: 1e18 source (18 dec) -> 1e9 local (9 dec).
	// Token lands in the signer wallet (suiAddr) so the  existing balance check stays meaningful;
	// the message leg runs (gas > 0) and ccip_receive stores the receiver-visible amount.
	tcs := []testhelpers.TestTransferRequest{
		{
			Name:             "token+message, receiver sees local amount",
			SourceChain:      sourceChain,
			DestChain:        destChain,
			Data:             []byte("Hello Sui from EVM"),
			Receiver:         receiverByte, // receiver contract pkgId -> ccip_receive runs
			TokenReceiverATA: suiAddr[:],   // wallet that actually receives the minted token
			ExpectedStatus:   testhelpers.EXECUTION_STATE_SUCCESS,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e18),
				},
			},
			ExtraArgs: testhelpers.MakeSuiExtraArgs(1000000, true, receiverObjectIDs, suiAddr),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  suiTokenBytes,
					Amount: big.NewInt(1e9),
				},
			},
		},
	}

	ctx := testhelpers.Context(t)
	prepareEvm2SuiTransferLane(t, e, state, sourceChain, destChain)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)
	replayEvm2SuiTransferLane(t, e, sourceChain, destChain)

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

	// Unchanged-path guard: the pool minted the correct local amount to the wallet.
	testhelpers.WaitForTokenBalances(ctx, t, e.Env, expectedTokenBalances)

	// dest_token_amounts[0].amount must be the local amount (1e9), NOT the source amount (1e18).
	suiChain := e.Env.BlockChains.SuiChains()[destChain]
	receiverContract, err := module_dummy_receiver.NewDummyReceiver(outputMap.PackageId, suiChain.Client)
	require.NoError(t, err)

	receiverStateObj := codec.Object{Id: outputMap.Objects.CCIPReceiverStateObjectId}
	devInspectOpts := &suiBind.CallOpts{
		Signer:           suiChain.Signer,
		WaitForExecution: true,
	}

	// ccip_receive ran (counter incremented) and stored exactly one token amount.
	counter, err := receiverContract.DevInspect().GetCounter(ctx, devInspectOpts, receiverStateObj)
	require.NoError(t, err)
	require.Positive(t, counter, "dummy receiver ccip_receive did not run")

	destTokenAmounts, err := receiverContract.DevInspect().GetDestTokenAmounts(ctx, devInspectOpts, receiverStateObj)
	require.NoError(t, err)
	require.Len(t, destTokenAmounts, 1, "expected exactly one dest token amount stored by the receiver")

	// GetDestTokenAmounts already decodes each TokenAmount (token + amount); read the amount
	// directly. Calling get_token_amount_amount separately is not viable via DevInspect because its
	// &TokenAmount argument is a plain struct (not an object/UID) and cannot be encoded as a MoveCall
	// arg.
	receiverVisibleAmount := destTokenAmounts[0].Amount
	require.NotNil(t, receiverVisibleAmount)

	// The fix: receiver sees the local (minted) amount, not the source-chain amount.
	require.Equalf(t, 0, receiverVisibleAmount.Cmp(big.NewInt(1e9)),
		"receiver-visible amount must be the local amount (1e9), got %s", receiverVisibleAmount.String())
	require.NotEqualf(t, 0, receiverVisibleAmount.Cmp(big.NewInt(1e18)),
		"receiver is seeing the unconverted source amount (1e18)")
}

// Test_CCIP_ReceiverNotRegistered_EVM2Sui covers destination-side receiver-failure handling for
// the two cases that behave identically on the Sui offramp:
//  1. Receiver package id NOT registered in the receiver registry (a real, deployed package that
//     was never registered).
//  2. Receiver address that does not exist on Sui (Invalid 32-byte address).
//
// On-chain (offramp.move pre_execute_single_report): has_valid_message_receiver =
// (!data.is_empty() || gas_limit != 0) && is_registered_receiver(...). When the receiver is not
// registered, NO Any2SuiMessage is created, the receiver callback is skipped, execution is marked
// SUCCESS, and the message is permanently dropped (terminal, no retry). A token leg in the same
// message completes independently -- assets are released without the receiver-side app logic ever
// running. This test locks in that behavior: SUCCESS, token delivered (token+message case), and the
// unregistered receiver's state unchanged (ccip_receive never ran).
func Test_CCIP_ReceiverNotRegistered_EVM2Sui(t *testing.T) {
	e, sourceChain, destChain, deployerSourceChain, suiTokenBytes, suiAddr := testSetupHelperEvm2Sui(t)

	// Token pool setup (needed for the token+message case); harmless for message-only cases.
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
	})
	require.NoError(t, err)
	e.Env = updatedEnv

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

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

	// Deploy a dummy receiver but DO NOT register it -> a valid, deployed package id that is not in
	// the receiver registry. We can still DevInspect its CCIPReceiverState to prove ccip_receive
	// never ran.
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

	unregisteredReceiverByte, err := hex.DecodeString(strings.TrimPrefix(outputMap.PackageId, "0x"))
	require.NoError(t, err)

	// Invalid 32-byte address: valid Sui address format, but no package exists and it is not
	// registered. (Distinct from the source-side "invalid receiver" revert test, which uses a
	// malformed non-32-byte address that the onramp rejects at send time.)
	nonexistentReceiver := common.HexToHash("0xdeadbeef00000000000000000000000000000000000000000000000000deadbeef")

	// Receiver object IDs are irrelevant when the receiver is not registered (no ccip_receive call),
	// so pass empty -- mirrors Test_CCIP_EVM2Sui_ZeroReceiver.
	emptyReceiverObjectIDs := [][32]byte{}

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "message-only to unregistered (deployed) receiver -> success, message dropped",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Data:           []byte("msg to unregistered receiver"),
			Receiver:       unregisteredReceiverByte,
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			ExtraArgs:      testhelpers.MakeSuiExtraArgs(1000000, true, emptyReceiverObjectIDs, [32]byte{}),
		},
		{
			Name:             "token+message to unregistered (deployed) receiver -> success, token delivered, message dropped",
			SourceChain:      sourceChain,
			DestChain:        destChain,
			Data:             []byte("token+msg to unregistered receiver"),
			Receiver:         unregisteredReceiverByte,
			TokenReceiverATA: suiAddr[:], // wallet receives the minted token
			ExpectedStatus:   testhelpers.EXECUTION_STATE_SUCCESS,
			Tokens: []router.ClientEVMTokenAmount{
				{Token: evmToken.Address(), Amount: big.NewInt(1e18)},
			},
			ExtraArgs: testhelpers.MakeSuiExtraArgs(1000000, true, emptyReceiverObjectIDs, suiAddr),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: suiTokenBytes, Amount: big.NewInt(1e9)},
			},
		},
		{
			Name:           "message-only to nonexistent (garbage) receiver -> success, message dropped",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Data:           []byte("msg to nonexistent receiver"),
			Receiver:       nonexistentReceiver[:],
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			ExtraArgs:      testhelpers.MakeSuiExtraArgs(1000000, true, emptyReceiverObjectIDs, [32]byte{}),
		},
	}

	ctx := testhelpers.Context(t)
	prepareEvm2SuiTransferLane(t, e, state, sourceChain, destChain)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)
	replayEvm2SuiTransferLane(t, e, sourceChain, destChain)

	err = testhelpers.ConfirmMultipleCommits(t, e.Env, state, startBlocks, false, expectedSeqNums)
	require.NoError(t, err)

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
		t,
		e.Env,
		state,
		testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
		startBlocks,
	)
	// All three cases: offramp skips the unregistered/nonexistent receiver and marks SUCCESS.
	require.Equal(t, expectedExecutionStates, execStates)

	// Token leg of the token+message case was still delivered (assets released without app logic).
	testhelpers.WaitForTokenBalances(ctx, t, e.Env, expectedTokenBalances)

	// Prove the message was dropped: the unregistered (deployed) receiver's ccip_receive never ran.
	suiChain := e.Env.BlockChains.SuiChains()[destChain]
	receiverContract, err := module_dummy_receiver.NewDummyReceiver(outputMap.PackageId, suiChain.Client)
	require.NoError(t, err)
	receiverStateObj := codec.Object{Id: outputMap.Objects.CCIPReceiverStateObjectId}
	devInspectOpts := &suiBind.CallOpts{Signer: suiChain.Signer, WaitForExecution: true}

	counter, err := receiverContract.DevInspect().GetCounter(ctx, devInspectOpts, receiverStateObj)
	require.NoError(t, err)
	require.Zero(t, counter, "unregistered receiver's ccip_receive must not have run (message dropped)")

	destTokenAmounts, err := receiverContract.DevInspect().GetDestTokenAmounts(ctx, devInspectOpts, receiverStateObj)
	require.NoError(t, err)
	require.Empty(t, destTokenAmounts, "unregistered receiver must have stored no dest token amounts")
}

func Test_CCIPPureTokenTransfer_EVM2Sui_BurnMintTokenPool(t *testing.T) {
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

	// update env to include deployed contracts
	e.Env = updatedEnv

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

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

	emptyReceiver := hexutil.MustDecode(
		"0x0000000000000000000000000000000000000000000000000000000000000000", // receiver packageID
	)

	tcs := []testhelpers.TestTransferRequest{
		// Pure token transfer
		// ReceiverObjectIds = empty
		// token.Receiver = non empty (maybe EOA or object)
		// message.Receiver = empty
		// don't need extraArgs gasLimit, can be set to 0
		{
			Name:             "Send token to EOA with - Pure Token Transfer",
			SourceChain:      sourceChain,
			DestChain:        destChain,
			Data:             []byte{},
			Receiver:         emptyReceiver, // empty Receiver
			TokenReceiverATA: suiAddr[:],    // tokenReceiver extracted from extraArgs (the address that actually gets the token)
			ExpectedStatus:   testhelpers.EXECUTION_STATE_SUCCESS,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e18),
				},
			},
			ExtraArgs: testhelpers.MakeSuiExtraArgs(0, true, [][32]byte{}, suiAddr),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  suiTokenBytes,
					Amount: big.NewInt(1e9),
				},
			},
		},
	}

	ctx := testhelpers.Context(t)
	prepareEvm2SuiTransferLane(t, e, state, sourceChain, destChain)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)
	replayEvm2SuiTransferLane(t, e, sourceChain, destChain)

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

func Test_CCIPProgrammableTokenTransfer_EVM2Sui_BurnMintTokenPool(t *testing.T) {
	e, sourceChain, destChain, deployerSourceChain, _, _ := testSetupHelperEvm2Sui(t)

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

	// update env to include deployed contracts
	e.Env = updatedEnv

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

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
		// Programmable token transfer
		// can be thought of as two separate paths tokenPool release/mint + message ccip_receive
		// receiverObjectIds = non empty (with clock & receiverStateValue)
		// token.Receiver = non empty(maybe EOA or object)
		// message.Receiver = receiverPackageId
		// extraArgs gasLimit > 0
		{
			Name:             "Send token to an Object",
			SourceChain:      sourceChain,
			DestChain:        destChain,
			Data:             []byte("Hello Sui From EVM"),
			Receiver:         receiverByte, // receiver contract pkgId
			TokenReceiverATA: stateObj[:],  // tokenReceiver extracted from extraArgs (the object that actually gets the token)
			ExpectedStatus:   testhelpers.EXECUTION_STATE_SUCCESS,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e18),
				},
			},
			ExtraArgs:             testhelpers.MakeSuiExtraArgs(1000000, true, receiverObjectIDs, stateObj), // receiver is objectId this time
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{},
		},
	}

	ctx := testhelpers.Context(t)
	prepareEvm2SuiTransferLane(t, e, state, sourceChain, destChain)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)
	replayEvm2SuiTransferLane(t, e, sourceChain, destChain)

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

func Test_CCIPZeroGasLimitTokenTransfer_EVM2Sui_BurnMintTokenPool(t *testing.T) {
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

	// update env to include deployed contracts
	e.Env = updatedEnv

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

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
			Name:             "Send token To EOA + include a receiver but keep gasLimit to 0",
			SourceChain:      sourceChain,
			DestChain:        destChain,
			Data:             []byte("Hello Sui From EVM"),
			Receiver:         receiverByte, // non empty Receiver
			TokenReceiverATA: suiAddr[:],   // tokenReceiver extracted from extraArgs (the address that actually gets the token)
			ExpectedStatus:   testhelpers.EXECUTION_STATE_SUCCESS,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e18),
				},
			},
			ExtraArgs: testhelpers.MakeSuiExtraArgs(0, true, receiverObjectIDs, suiAddr), // keep gasLimit to 0
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  suiTokenBytes,
					Amount: big.NewInt(1e9),
				},
			},
		},
	}

	ctx := testhelpers.Context(t)
	prepareEvm2SuiTransferLane(t, e, state, sourceChain, destChain)
	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)
	replayEvm2SuiTransferLane(t, e, sourceChain, destChain)

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

func testSetupTokenTransferSui2Evm(t *testing.T) (e testhelpers.DeployedEnv, sourceChain, destChain uint64) {
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

// Test_CCIP_TokenTransfer_EVM2Sui_PoolReleaseOrMintTransmitterOwned_Rejected
// verifies the offchain ReleaseOrMintParams ownership guard end-to-end. A
// TEST-only malicious Sui burn-mint LINK pool is registered with an
// executor-transmitter-owned SUI coin id in its release_or_mint_params, and its
// release_or_mint declares an extra &mut Coin<SUI> drain tail. During a normal
// EVM->Sui token transfer the production execute PTB builder appends that coin
// to the pool callback under the transmitter's signature; the guard rejects the
// transmitter-owned entry during PTB build, so the execute PTB is never
// submitted, EXECUTION_STATE_SUCCESS is never persisted, and the transmitter's
// coin is not drained. Mirrors Test_CCIP_Messaging_EVM2Sui_TransmitterOwnedTail_Rejected.
func Test_CCIP_TokenTransfer_EVM2Sui_PoolReleaseOrMintTransmitterOwned_Rejected(t *testing.T) {
	e, sourceChain, destChain, deployerSourceChain, _, suiAddr := testSetupHelperEvm2Sui(t)

	suiChain := e.Env.BlockChains.SuiChains()[destChain]
	waitForSuiRPCSync(t, suiChain)

	// release_or_mint_params is fixed at pool initialize, so resolve the exec
	// transmitter's SUI gas coin before deploying the malicious pool.
	transmitterAddr := suiExecTransmitterAddress(t, e, destChain)
	require.NotEmpty(t, transmitterAddr)
	coins, err := suiChain.Client.QueryCoinsByAddress(testhelpers.Context(t), transmitterAddr, "0x2::coin::Coin<0x2::sui::SUI>")
	require.NoError(t, err)
	require.NotEmpty(t, coins, "exec transmitter must own SUI gas coins")
	transmitterCoinID := coins[0].GetObjectId()
	preCoinBalance := new(big.Int).SetUint64(coins[0].GetBalance())

	// Deploy the malicious LINK burn-mint pool as the Sui dest pool, registered
	// with the transmitter coin id in release_or_mint_params. The EVM-side
	// token/pool + cross-registration are handled by the sibling helper.
	updatedEnv, evmToken, _, _, _, err := testhelpers.HandleMaliciousBurnMintTokenPoolDeploymentForSUI(
		e.Env, destChain, sourceChain, "0x1", []string{transmitterCoinID},
	)
	require.NoError(t, err)
	e.Env = updatedEnv

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

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

	// Deploy a benign dummy receiver for the message leg. The exploit vector is
	// the pool, not the receiver, so the receiverObjectIDs carry no drain coin.
	_, output, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.DeployDummyReceiver{}, sui_cs.DeployDummyReceiverConfig{
			SuiChainSelector: destChain,
			McmsOwner:        "0x1",
		}),
	})
	require.NoError(t, err)

	rawOutput := output[0].Reports[0]
	receiverOutput, ok := rawOutput.Output.(sui_ops.OpTxResult[ccipops.DeployDummyReceiverObjects])
	require.True(t, ok)

	id := strings.TrimPrefix(receiverOutput.PackageId, "0x")
	receiverByte, err := hex.DecodeString(id)
	require.NoError(t, err)

	_, _, err = commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.RegisterDummyReceiver{}, sui_cs.RegisterDummyReceiverConfig{
			SuiChainSelector:       destChain,
			OwnerCapObjectId:       receiverOutput.Objects.OwnerCapObjectId,
			CCIPObjectRefObjectId:  state.SuiChains[destChain].CCIPObjectRef,
			DummyReceiverPackageId: receiverOutput.PackageId,
		}),
	})
	require.NoError(t, err)

	var clockObj [32]byte
	copy(clockObj[:], hexutil.MustDecode(
		"0x0000000000000000000000000000000000000000000000000000000000000006",
	))
	var receiverStateObj [32]byte
	copy(receiverStateObj[:], hexutil.MustDecode(receiverOutput.Objects.CCIPReceiverStateObjectId))
	receiverObjectIDs := [][32]byte{clockObj, receiverStateObj}

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:             "Pool release_or_mint transmitter-owned coin drain attempt",
			SourceChain:      sourceChain,
			DestChain:        destChain,
			Receiver:         receiverByte,
			TokenReceiverATA: suiAddr[:],
			// Exploit message is never finalized: the guard aborts the execute PTB
			// before submission, so no ExecutionStateChanged event fires. The real
			// assertion is the explicit GetExecutionState DevInspect below (require.Error).
			// UNTOUCHED is the honest expected state; TransferMultiple's expected-states
			// map is discarded here, so this field is documentary only.
			ExpectedStatus: testhelpers.EXECUTION_STATE_UNTOUCHED,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e18),
				},
			},
			ExtraArgs: testhelpers.MakeSuiExtraArgs(1_000_000, true, receiverObjectIDs, suiAddr),
			// No ExpectedTokenBalances: the guard aborts the execute PTB before
			// submission, so no token is released/minted on Sui. The real balance
			// check is the explicit suiCoinByID assertion on the transmitter coin
			// below; WaitForTokenBalances' accumulator is discarded here anyway.
		},
	}

	ctx := testhelpers.Context(t)
	prepareEvm2SuiTransferLane(t, e, state, sourceChain, destChain)
	startBlocks, expectedSeqNums, _, _ := testhelpers.TransferMultiple(ctx, t, e.Env, state, tcs)

	pair := testhelpers.SourceDestPair{SourceChainSelector: sourceChain, DestChainSelector: destChain}
	seqRange, ok := expectedSeqNums[pair]
	require.True(t, ok, "expected a sequence number for the exploit transfer")
	seqExploit := uint64(seqRange.Start())

	// Let the DON replay + settle so the source event is OCR-committed and the
	// dest exec is attempted. The guard rejects the transmitter-owned
	// release_or_mint_params entry during PTB build, so no execute transaction is
	// submitted and EXECUTION_STATE_SUCCESS never commits.
	replayEvm2SuiTransferLane(t, e, sourceChain, destChain)

	// Confirm the report was committed, so the exec rejection below is meaningful.
	err = testhelpers.ConfirmMultipleCommits(t, e.Env, state, startBlocks, false, expectedSeqNums)
	require.NoError(t, err)

	// With no executed transaction, the sequence is absent and get_execution_state
	// aborts EUnknownSequenceNumber (atomic-rollback / never-finalized guarantee).
	suiState, err := sui_deployment.LoadOnchainStatesui(e.Env)
	require.NoError(t, err)
	offrampContract, err := module_offramp.NewOfframp(suiState[destChain].EffectiveOffRampPackageID(), suiChain.Client)
	require.NoError(t, err)
	offRampStateObj := codec.Object{Id: suiState[destChain].OffRampStateObjectId}
	devInspectOpts := &suiBind.CallOpts{
		Signer:           suiChain.Signer,
		WaitForExecution: true,
	}
	_, err = offrampContract.DevInspect().GetExecutionState(ctx, devInspectOpts, offRampStateObj, sourceChain, seqExploit)
	require.Error(t, err, "exploit execute must not finalize (get_execution_state must NOT return SUCCESS)")

	// No drain: the guard rejects the transmitter-owned release_or_mint_params
	// entry, so the pool's release_or_mint never receives the coin. The coin may
	// still be selected as the gas coin for unrelated transactions, so its balance
	// can drop by a tiny fraction, but a drain would transfer its value off the
	// transmitter or split it down to ~0. Assert the coin is still owned by the
	// transmitter AND retained most of its value.
	postCoinBalance, stillOwned := suiCoinByID(t, suiChain, transmitterAddr, transmitterCoinID)
	require.Truef(t, stillOwned,
		"transmitter coin %s is no longer owned by the transmitter after execute (guard should prevent the pool from taking it)",
		transmitterCoinID)
	decrease := new(big.Int).Sub(preCoinBalance, postCoinBalance)
	half := new(big.Int).Quo(preCoinBalance, big.NewInt(2))
	require.Negativef(t, decrease.Cmp(half),
		"transmitter coin drained: pre=%s post=%s decrease=%s (guard should prevent the pool from taking the coin's value; only gas should be charged)",
		preCoinBalance.String(), postCoinBalance.String(), decrease.String())

	// Lane-not-stuck control: a subsequent honest EVM->Sui arbitrary-data message
	// must still finalize SUCCESS. The guard skips only the exploit's token-pool
	// command (non-retryable, off-chain); the offramp must keep processing later
	// messages, proving the lane is not head-of-line blocked by the rejected seq.
	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	honestTcs := []testhelpers.TestTransferRequest{
		{
			Name:             "EVM2Sui arbitrary message after rejected exploit (lane not stuck)",
			SourceChain:      sourceChain,
			DestChain:        destChain,
			Receiver:         receiverByte,
			TokenReceiverATA: []byte{},
			ExpectedStatus:   testhelpers.EXECUTION_STATE_SUCCESS,
			Tokens:           []router.ClientEVMTokenAmount{}, // arbitrary data, no token-pool command
			Data:             []byte("lane not stuck"),
			// Message-only: token receiver must be zero. The Sui offramp asserts
			// token_receiver == @0x0 iff the message carries no token amounts, else
			// it aborts EInvalidTokenReceiver at init_execute and the whole
			// execute PTB reverts, rolling back ExecutionStateChanged.
			ExtraArgs: testhelpers.MakeSuiExtraArgs(1_000_000, true, receiverObjectIDs, [32]byte{}),
		},
	}
	honestStartBlocks, honestExpectedSeqNums, honestExpectedExecStates, _ := testhelpers.TransferMultiple(ctx, t, e.Env, state, honestTcs)
	replayEvm2SuiTransferLane(t, e, sourceChain, destChain)
	err = testhelpers.ConfirmMultipleCommits(t, e.Env, state, honestStartBlocks, false, honestExpectedSeqNums)
	require.NoError(t, err)
	honestExecStates := testhelpers.ConfirmExecWithSeqNrsForAll(t, e.Env, state, testhelpers.SeqNumberRangeToSlice(honestExpectedSeqNums), honestStartBlocks)
	require.Equal(t, honestExpectedExecStates, honestExecStates, "honest EVM2Sui message must execute after the rejected exploit (lane not stuck)")
}

func testSetupHelperEvm2Sui(t *testing.T) (e testhelpers.DeployedEnv, sourceChain, destChain uint64, deployerSourceChain *bind.TransactOpts, suiTokenBytes []byte, suiAddr [32]byte) {
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

func assertSuiSourceRevertExpectedError(t *testing.T, err error, execRevertErrorMsg, execRevertCauseErrorMsg string) {
	require.Error(t, err)
	fmt.Println("Error: ", err.Error())
	require.Contains(t, err.Error(), execRevertErrorMsg)
	require.Contains(t, err.Error(), execRevertCauseErrorMsg)
}

// prepareEvm2SuiTransferLane waits for CCIPMessageSent log-poller registration on the EVM
// source chain before sending, matching ccip_token_transfer_test.go.
func prepareEvm2SuiTransferLane(t *testing.T, e testhelpers.DeployedEnv, state stateview.CCIPOnChainState, sourceChain, destChain uint64) {
	t.Helper()
	testhelpers.WaitForEventFilterRegistrationOnLane(t, state, e.Env.Offchain, sourceChain, destChain)
}

// replayEvm2SuiTransferLane replays EVM logs and waits for async replay to settle so the DON
// can OCR-commit before we poll the Sui offramp (Sui replay is a no-op in nodetestutils).
func replayEvm2SuiTransferLane(t *testing.T, e testhelpers.DeployedEnv, sourceChain, destChain uint64) {
	t.Helper()
	_ = destChain
	messagingtest.SleepReplayAndSettle(t, e.Env, 30*time.Second, sourceChain)
}

// waitForSuiRPCSync blocks until the Sui fullnode JSON-RPC view has had a chance to index
// recent transactions, by waiting for the latest checkpoint sequence to advance.
//
// Background: since Sui v1.33 the JSON-RPC node silently ignores `requestType:
// "WaitForLocalExecution"` and returns as soon as effects are certified
// (https://forums.sui.io/t/deprecating-waitforlocalexecution/45988). The Typescript SDK
// works around this by polling `client.waitForTransaction({ digest })`, but the Go SDK
// (block-vision/sui-go-sdk) does not. As a result, a tight sequence like
// "mutating admin tx -> ccip_send" can fetch stale owned-object versions (e.g. the gas
// coin) and the validators reject the second tx with "Object ... Version ... is not
// available for consumption" — masking the Move abort we are trying to assert.
//
// Call this helper at the top of any "should fail" subtest that submits a Sui tx
// immediately after a previous Sui tx in the same test. This is a test-side band-aid;
// the proper fix belongs in chainlink-sui/bindings/bind (poll sui_getTransactionBlock on
// the returned digest, matching the Typescript SDK behavior).
// waitForSuiRPCSyncWithOptions provides configurable Sui RPC synchronization
func waitForSuiRPCSyncWithOptions(t *testing.T, suiChain sui.Chain, timeout time.Duration, minAdvance uint64) {
	t.Helper()

	const pollInterval = 1 * time.Second

	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	defer cancel()

	beforeCheckpoint, err := suiChain.Client.GetLatestCheckpoint(ctx)
	if err != nil {
		t.Logf("waitForSuiRPCSync: failed to read initial checkpoint seq (%v); falling back to fixed sleep", err)
		// Use context-aware sleep instead of fixed sleep
		fallbackSleep := min(timeout/2, 5*time.Second)
		select {
		case <-time.After(fallbackSleep):
			return
		case <-ctx.Done():
			t.Logf("waitForSuiRPCSync: context cancelled during fallback sleep")
			return
		}
	}

	before := beforeCheckpoint.GetSequenceNumber()

	t.Logf("waitForSuiRPCSync: starting sync wait from checkpoint %d (timeout: %v, minAdvance: %d)", before, timeout, minAdvance)

	// Use context deadline instead of separate time tracking
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.Logf("waitForSuiRPCSync: timeout waiting for checkpoint to advance from %d (context: %v)", before, ctx.Err())
			return
		case <-ticker.C:
			afterCheckpoint, cerr := suiChain.Client.GetLatestCheckpoint(ctx)
			if cerr != nil {
				// Log but continue - network might be temporarily unavailable
				t.Logf("waitForSuiRPCSync: temporary error reading checkpoint: %v", cerr)
				continue
			}
			after := afterCheckpoint.GetSequenceNumber()

			// Check if we've advanced enough checkpoints
			if after >= before+minAdvance {
				t.Logf("waitForSuiRPCSync: sync complete - checkpoint advanced from %d to %d", before, after)
				return
			}

			// Log progress for debugging
			if after > before {
				t.Logf("waitForSuiRPCSync: checkpoint progressing: %d -> %d (need %d total advance)", before, after, minAdvance)
			}
		}
	}
}

// waitForSuiRPCSync provides default Sui RPC synchronization with improved reliability
func waitForSuiRPCSync(t *testing.T, suiChain sui.Chain) {
	// Use more conservative defaults for better reliability
	waitForSuiRPCSyncWithOptions(t, suiChain, 30*time.Second, 1)
}

// waitForSuiRPCSyncFast provides faster synchronization for less critical operations
func waitForSuiRPCSyncFast(t *testing.T, suiChain sui.Chain) {
	waitForSuiRPCSyncWithOptions(t, suiChain, 15*time.Second, 1)
}

// waitForSuiRPCSyncSlow provides extended synchronization for critical operations like upgrades
func waitForSuiRPCSyncSlow(t *testing.T, suiChain sui.Chain) {
	waitForSuiRPCSyncWithOptions(t, suiChain, 60*time.Second, 2)
}

// waitForSuiRPCSyncUpgrade provides maximum synchronization for post-upgrade event indexing
func waitForSuiRPCSyncUpgrade(t *testing.T, suiChain sui.Chain) {
	waitForSuiRPCSyncWithOptions(t, suiChain, 180*time.Second, 5)
}

// waitForSuiRPCSyncCritical provides extended synchronization for critical CCIP system reconfigurations
func waitForSuiRPCSyncCritical(t *testing.T, suiChain sui.Chain) {
	waitForSuiRPCSyncWithOptions(t, suiChain, 300*time.Second, 10)
}
