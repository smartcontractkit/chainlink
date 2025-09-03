package ccip

import (
	"fmt"
	"math/big"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func assertAptosRegulatedTokenRevertExpectedError(t *testing.T, err error, execRevertErrorMsg string, execRevertCauseErrorMsg string) {
	require.Error(t, err)
	fmt.Println("Error: ", err.Error())
	require.Contains(t, err.Error(), execRevertErrorMsg)
	require.Contains(t, err.Error(), execRevertCauseErrorMsg)
}

func Test_CCIP_RegulatedTokenTransfer_EVM2Aptos(t *testing.T) {
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithAptosChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	aptosChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyAptos))

	// Deploy the dummy receiver contract
	testhelpers.DeployAptosCCIPReceiver(t, e.Env)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := evmChainSelectors[0]
	destChain := aptosChainSelectors[0]
	deployerSourceChain := e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey
	deployerDestChain := e.Env.BlockChains.AptosChains()[destChain].DeployerSigner.AccountAddress()
	ccipChainState := state.AptosChains[destChain]

	lggr.Debug("Source chain (EVM): ", sourceChain, "Dest chain (Aptos): ", destChain)

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)

	evmToken, _, aptosToken, _, err := testhelpers.DeployRegulatedTransferableTokenAptos(t, lggr, e.Env, sourceChain, destChain, "Regulated Token", nil)
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

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:           "Send regulated token to EOA",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       deployerDestChain[:],
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e18),
				},
			},
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(100000, true),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  aptosToken[:],
					Amount: big.NewInt(1e8),
				},
			},
		},
		{
			Name:           "Send regulated token to EOA with gas limit set to 0",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       deployerDestChain[:],
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e18),
				},
			},
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(0, true),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{
					Token:  aptosToken[:],
					Amount: big.NewInt(1e8),
				},
			},
		},
		{
			Name:           "Send regulated token to contract",
			SourceChain:    sourceChain,
			DestChain:      destChain,
			Receiver:       ccipChainState.ReceiverAddress[:],
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  evmToken.Address(),
					Amount: big.NewInt(1e18),
				},
			},
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(100000, true),
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

	execStates := testhelpers.ConfirmExecWithSeqNrsForAll(t, e.Env, state, testhelpers.SeqNumberRangeToSlice(expectedSeqNums), startBlocks)
	require.Equal(t, expectedExecutionStates, execStates)

	testhelpers.WaitForTokenBalances(ctx, t, e.Env, expectedTokenBalances)
}
