package ccip

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common/hexutil"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"

	sui_cs "github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIPTokenTransfer_Sui2EVM(t *testing.T) {
	ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	fmt.Println("EVM: ", evmChainSelectors[0])
	fmt.Println("Sui: ", suiChainSelectors[0])

	sourceChain := suiChainSelectors[0]
	destChain := evmChainSelectors[0]

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Log("Source chain (Sui): ", sourceChain, "Dest chain (EVM): ", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	_, err = e.Env.BlockChains.SuiChains()[sourceChain].Signer.GetAddress()
	require.NoError(t, err)

	// SUI FeeToken
	// mint link token to use as feeToken
	_, feeTokenOutput, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.MintLinkToken{}, sui_cs.MintLinkTokenConfig{
			ChainSelector:  sourceChain,
			TokenPackageId: state.SuiChains[sourceChain].LinkTokenAddress,
			TreasuryCapId:  state.SuiChains[sourceChain].LinkTokenTreasuryCapId,
			Amount:         1000000000000, // 1000Link with 1e9,
		}),
	})
	require.NoError(t, err)

	rawOutput := feeTokenOutput[0].Reports[0]
	outputMap, ok := rawOutput.Output.(sui_ops.OpTxResult[linkops.MintLinkTokenOutput])
	require.True(t, ok)

	// SUI TransferToken
	// mint link token to use as Transfer Token
	_, transferTokenOutput, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.MintLinkToken{}, sui_cs.MintLinkTokenConfig{
			ChainSelector:  sourceChain,
			TokenPackageId: state.SuiChains[sourceChain].LinkTokenAddress,
			TreasuryCapId:  state.SuiChains[sourceChain].LinkTokenTreasuryCapId,
			Amount:         1000000000, // 1Link with 1e9
		}),
	})
	require.NoError(t, err)

	rawOutputTransferToken := transferTokenOutput[0].Reports[0]
	outputMapTransferToken, ok := rawOutputTransferToken.Output.(sui_ops.OpTxResult[linkops.MintLinkTokenOutput])
	require.True(t, ok)

	// mint more token
	_, transferTokenOutput1, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.MintLinkToken{}, sui_cs.MintLinkTokenConfig{
			ChainSelector:  sourceChain,
			TokenPackageId: state.SuiChains[sourceChain].LinkTokenAddress,
			TreasuryCapId:  state.SuiChains[sourceChain].LinkTokenTreasuryCapId,
			Amount:         2000000000, // 1Link with 1e9
		}),
	})
	require.NoError(t, err)

	rawOutputTransferToken1 := transferTokenOutput1[0].Reports[0]
	outputMapTransferToken1, ok := rawOutputTransferToken1.Output.(sui_ops.OpTxResult[linkops.MintLinkTokenOutput])
	require.True(t, ok)

	// Receiver Address
	ccipReceiverAddress := state.Chains[destChain].Receiver.Address()

	// Token Pool setup on both SUI and EVM
	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndPoolDeploymentForSUI(e.Env, sourceChain, destChain) // SourceChain = SUI, destChain = EVM
	require.NoError(t, err)
	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send token to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       updatedEnv.BlockChains.EVMChains()[destChain].DeployerKey.From.Bytes(), // internally left padded to 32byte
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			FeeToken:       outputMap.Objects.MintedLinkTokenObjectId,
			SuiTokens: []testhelpers.SuiTokenAmount{
				{
					Token:  outputMapTransferToken.Objects.MintedLinkTokenObjectId,
					Amount: 1000000000, // Send 1Link to EVM
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
					Token:  outputMapTransferToken1.Objects.MintedLinkTokenObjectId,
					Amount: 2000000000, // Send 1Link to EVM
				},
			},
			FeeToken: outputMap.Objects.MintedLinkTokenObjectId,
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  evmToken.Address().Bytes(),
					Amount: big.NewInt(2e18),
				},
			},
		},
	}

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

func Test_CCIPTokenTransfer_EVM2SUI(t *testing.T) {
	ctx := testhelpers.Context(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	fmt.Println("EVM: ", evmChainSelectors[0])
	fmt.Println("Sui: ", suiChainSelectors[0])

	sourceChain := evmChainSelectors[0]
	destChain := suiChainSelectors[0]

	t.Log("Source chain (Sui): ", sourceChain, "Dest chain (EVM): ", destChain)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	deployerSourceChain := e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey
	var suiTokenAddr [32]byte
	suiTokenHex := state.SuiChains[destChain].LinkTokenAddress
	suiTokenHex = strings.TrimPrefix(suiTokenHex, "0x")

	suiTokenBytes, err := hex.DecodeString(suiTokenHex)
	require.NoError(t, err)

	require.NoError(t, err)

	require.Len(t, suiTokenBytes, 32, "expected 32-byte sui address")
	copy(suiTokenAddr[:], suiTokenBytes)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	// get sui address in [32]bytes for extraArgs.TokenReceiver
	var suiAddr [32]byte
	suiAddrStr, err := e.Env.BlockChains.SuiChains()[destChain].Signer.GetAddress()
	require.NoError(t, err)

	suiAddrStr = strings.TrimPrefix(suiAddrStr, "0x")

	addrBytes, err := hex.DecodeString(suiAddrStr)
	require.NoError(t, err)

	require.Len(t, addrBytes, 32, "expected 32-byte sui address")
	copy(suiAddr[:], addrBytes)

	// Token Pool setup on both SUI and EVM
	updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndPoolDeploymentForSUI(e.Env, destChain, sourceChain) // sourceChain=EVM, destChain=SUI
	require.NoError(t, err)

	state, err = stateview.LoadOnchainState(e.Env)
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

	// TODO: might be needed for validation
	// getPoolBySourceToken
	// onRamp, err := onramp.NewOnRamp(state.Chains[sourceChain].OnRamp.Address(), e.Env.BlockChains.EVMChains()[sourceChain].Client)
	// require.NoError(t, err)

	// poolAddr, err := onRamp.GetPoolBySourceToken(&bind.CallOpts{}, destChain, evmToken.Address())
	// require.NoError(t, err)

	// fmt.Println("POOL ADDR: ", poolAddr)

	// getRemoteToken
	// tp, err := burn_from_mint_token_pool.NewBurnFromMintTokenPool(evmTokenPool.Address(), e.Env.BlockChains.EVMChains()[sourceChain].Client)
	// require.NoError(t, err)

	// remoteToken, err := tp.GetRemoteToken(&bind.CallOpts{}, destChain)
	// require.NoError(t, err)

	// remotePool, err := tp.GetRemotePools(&bind.CallOpts{}, destChain)
	// require.NoError(t, err)

	// fmt.Println("REMOTETOKEN: ", remoteToken)
	// fmt.Println("REMOTEPOOL: ", remotePool)

	// fmt.Println("TOKENBALANCE TEST: RECEIVER: ", suiAddrStr, " TOKENN: ", suiTokenHex)

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
		{
			Name:             "Send token to an Object",
			SourceChain:      sourceChain,
			DestChain:        destChain,
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
