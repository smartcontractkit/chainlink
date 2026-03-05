package ccip

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	sui_deployment "github.com/smartcontractkit/chainlink-sui/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIP_Fees_Sui2EVM(t *testing.T) {
	ctx := testhelpers.Context(t)

	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	sourceChain := suiChainSelectors[0]
	destChain := evmChainSelectors[0]

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	t.Log("Source chain (SUI):", sourceChain, "Dest chain (EVM):", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	suiChain := e.Env.BlockChains.SuiChains()[sourceChain]

	suiCoinMetadataId, err := testhelpers.GetSuiNativeCoinMetadataID(ctx, suiChain.Client)
	require.NoError(t, err, "failed to get SUI CoinMetadata ID")
	t.Log("SUI CoinMetadata ID:", suiCoinMetadataId)

	testhelpers.RegisterSuiNativeFeeToken(t, e.Env, sourceChain, suiCoinMetadataId)

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	receiver := state.Chains[destChain].Receiver.Address()

	linkFeeTokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 1000000000000)
	linkFeeToken := linkFeeTokenOutput.Objects.MintedLinkTokenObjectId

	t.Run("Send message with LINK fee token", func(t *testing.T) {
		msg := testhelpers.SuiSendRequest{
			Receiver:  common.LeftPadBytes(receiver.Bytes(), 32),
			Data:      []byte("Hello EVM, from SUI with LINK fee!"),
			FeeToken:  linkFeeToken,
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
		}

		msgSentEvent := testhelpers.TestSendRequest(t, e.Env, state, sourceChain, destChain, false, msg)
		require.NotNil(t, msgSentEvent)
		t.Log("LINK fee message sent, seqNum:", msgSentEvent.SequenceNumber)

		seqNum := ccipocr3.SeqNum(msgSentEvent.SequenceNumber)
		expectedSeqNums := map[testhelpers.SourceDestPair]ccipocr3.SeqNumRange{
			{SourceChainSelector: sourceChain, DestChainSelector: destChain}: ccipocr3.NewSeqNumRange(seqNum, seqNum),
		}
		startBlocks := map[uint64]*uint64{}
		block, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		startBlocks[destChain] = &block

		err = testhelpers.ConfirmMultipleCommits(t, e.Env, state, startBlocks, false, expectedSeqNums)
		require.NoError(t, err)

		execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
			t, e.Env, state,
			testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
			startBlocks,
		)
		for _, states := range execStates {
			for _, s := range states {
				require.Equal(t, testhelpers.EXECUTION_STATE_SUCCESS, s)
			}
		}
	})

	t.Run("Send message with native SUI fee token", func(t *testing.T) {
		suiFeeCoinId := testhelpers.SplitSuiCoinForFee(ctx, t, suiChain, 10_000_000_000)

		msg := testhelpers.SuiSendRequest{
			Receiver:           common.LeftPadBytes(receiver.Bytes(), 32),
			Data:               []byte("Hello EVM, from SUI with native SUI fee!"),
			FeeToken:           suiFeeCoinId,
			FeeTokenCoinType:   testhelpers.SuiNativeCoinType,
			FeeTokenMetadataID: suiCoinMetadataId,
			ExtraArgs:          testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
		}

		msgSentEvent := testhelpers.TestSendRequest(t, e.Env, state, sourceChain, destChain, false, msg)
		require.NotNil(t, msgSentEvent)
		t.Log("Native SUI fee message sent, seqNum:", msgSentEvent.SequenceNumber)

		seqNum := ccipocr3.SeqNum(msgSentEvent.SequenceNumber)
		expectedSeqNums := map[testhelpers.SourceDestPair]ccipocr3.SeqNumRange{
			{SourceChainSelector: sourceChain, DestChainSelector: destChain}: ccipocr3.NewSeqNumRange(seqNum, seqNum),
		}
		startBlocks := map[uint64]*uint64{}
		block, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		startBlocks[destChain] = &block

		err = testhelpers.ConfirmMultipleCommits(t, e.Env, state, startBlocks, false, expectedSeqNums)
		require.NoError(t, err)

		execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
			t, e.Env, state,
			testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
			startBlocks,
		)
		for _, states := range execStates {
			for _, s := range states {
				require.Equal(t, testhelpers.EXECUTION_STATE_SUCCESS, s)
			}
		}
	})

	t.Run("Token transfer with LINK fee token", func(t *testing.T) {
		updatedEnv, evmToken, _, err := testhelpers.HandleTokenAndBurnMintTokenPoolDeploymentForSUI(
			e.Env, sourceChain, destChain, []testhelpers.TokenPoolRateLimiterConfig{
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

		state, err = stateview.LoadOnchainState(e.Env)
		require.NoError(t, err)

		tokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 1_000_000_000)
		feeOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 100_000_000_000)

		tcs := []testhelpers.TestTransferRequest{
			{
				Name:           "Token transfer with LINK fee",
				SourceChain:    sourceChain,
				DestChain:      destChain,
				Receiver:       updatedEnv.BlockChains.EVMChains()[destChain].DeployerKey.From.Bytes(),
				ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
				FeeToken:       feeOutput.Objects.MintedLinkTokenObjectId,
				SuiTokens: []testhelpers.SuiTokenAmount{
					{
						TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
						Token:         tokenOutput.Objects.MintedLinkTokenObjectId,
						Amount:        1_000_000_000,
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

		startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances := testhelpers.TransferMultiple(ctx, t, updatedEnv, state, tcs)

		err = testhelpers.ConfirmMultipleCommits(t, updatedEnv, state, startBlocks, false, expectedSeqNums)
		require.NoError(t, err)

		execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
			t, updatedEnv, state,
			testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
			startBlocks,
		)
		require.Equal(t, expectedExecutionStates, execStates)
		testhelpers.WaitForTokenBalances(ctx, t, updatedEnv, expectedTokenBalances)
	})

	t.Run("Token transfer with native SUI fee token", func(t *testing.T) {
		state, err = stateview.LoadOnchainState(e.Env)
		require.NoError(t, err)

		tokenOutput := mintLinkTokenOnSui(t, e.Env, sourceChain, 1_000_000_000)
		suiFeeCoinId := testhelpers.SplitSuiCoinForFee(ctx, t, suiChain, 10_000_000_000)

		msg := testhelpers.SuiSendRequest{
			Receiver:           common.LeftPadBytes(e.Env.BlockChains.EVMChains()[destChain].DeployerKey.From.Bytes(), 32),
			Data:               []byte{},
			FeeToken:           suiFeeCoinId,
			FeeTokenCoinType:   testhelpers.SuiNativeCoinType,
			FeeTokenMetadataID: suiCoinMetadataId,
			ExtraArgs:          testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
			TokenAmounts: []testhelpers.SuiTokenAmount{
				{
					TokenPoolType: sui_deployment.TokenPoolTypeBurnMint,
					Token:         tokenOutput.Objects.MintedLinkTokenObjectId,
					Amount:        1_000_000_000,
				},
			},
		}

		msgSentEvent := testhelpers.TestSendRequest(t, e.Env, state, sourceChain, destChain, false, msg)
		require.NotNil(t, msgSentEvent)
		t.Log("Native SUI fee token transfer sent, seqNum:", msgSentEvent.SequenceNumber)

		seqNum := ccipocr3.SeqNum(msgSentEvent.SequenceNumber)
		expectedSeqNums := map[testhelpers.SourceDestPair]ccipocr3.SeqNumRange{
			{SourceChainSelector: sourceChain, DestChainSelector: destChain}: ccipocr3.NewSeqNumRange(seqNum, seqNum),
		}
		startBlocks := map[uint64]*uint64{}
		block, err := testhelpers.LatestBlock(ctx, e.Env, destChain)
		require.NoError(t, err)
		startBlocks[destChain] = &block

		err = testhelpers.ConfirmMultipleCommits(t, e.Env, state, startBlocks, false, expectedSeqNums)
		require.NoError(t, err)

		execStates := testhelpers.ConfirmExecWithSeqNrsForAll(
			t, e.Env, state,
			testhelpers.SeqNumberRangeToSlice(expectedSeqNums),
			startBlocks,
		)
		for _, states := range execStates {
			for _, s := range states {
				require.Equal(t, testhelpers.EXECUTION_STATE_SUCCESS, s)
			}
		}
	})

	t.Run("Send with invalid fee token should fail", func(t *testing.T) {
		msg := testhelpers.SuiSendRequest{
			Receiver:  common.LeftPadBytes(receiver.Bytes(), 32),
			Data:      []byte("should fail"),
			FeeToken:  "0x0000000000000000000000000000000000000000000000000000000000000bad",
			ExtraArgs: testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
		}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}
		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err, "expected error for invalid fee token")
		t.Log("Invalid fee token correctly rejected:", err)
	})

	t.Run("Send with insufficient SUI fee token balance should fail", func(t *testing.T) {
		tinyCoinId := testhelpers.SplitSuiCoinForFee(ctx, t, suiChain, 1)

		msg := testhelpers.SuiSendRequest{
			Receiver:           common.LeftPadBytes(receiver.Bytes(), 32),
			Data:               []byte("should fail - insufficient balance"),
			FeeToken:           tinyCoinId,
			FeeTokenCoinType:   testhelpers.SuiNativeCoinType,
			FeeTokenMetadataID: suiCoinMetadataId,
			ExtraArgs:          testhelpers.MakeBCSEVMExtraArgsV2(big.NewInt(300000), false),
		}

		baseOpts := []ccipclient.SendReqOpts{
			ccipclient.WithSourceChain(sourceChain),
			ccipclient.WithDestChain(destChain),
			ccipclient.WithTestRouter(false),
			ccipclient.WithMessage(msg),
		}
		_, err := testhelpers.SendRequest(e.Env, state, baseOpts...)
		require.Error(t, err, "expected error for insufficient fee token balance")
		t.Log("Insufficient balance correctly rejected:", err)
	})
}
