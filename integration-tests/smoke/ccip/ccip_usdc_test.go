package ccip

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

/*
* Chain topology for this test
* 	chainA (USDC, MY_TOKEN)
*			|
*			| ------- chainC (USDC, MY_TOKEN)
*			|
* 	chainB (USDC)
 */
func TestUSDCTokenTransfer(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := t.Context()
	tenv, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfUsersPerChain(3),
		testhelpers.WithNumOfChains(3),
		testhelpers.WithUSDC(),
	)

	e := tenv.Env
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	evmChains := e.BlockChains.EVMChains()
	allChainSelectors := maps.Keys(evmChains)
	chainA := allChainSelectors[0]
	chainC := allChainSelectors[1]
	chainB := allChainSelectors[2]

	ownerChainA := evmChains[chainA].DeployerKey
	ownerChainC := evmChains[chainC].DeployerKey
	ownerChainB := evmChains[chainB].DeployerKey

	aChainUSDC, cChainUSDC, err := testhelpers.ConfigureUSDCTokenPools(lggr, evmChains, chainA, chainC, state)
	require.NoError(t, err)

	bChainUSDC, _, err := testhelpers.ConfigureUSDCTokenPools(lggr, evmChains, chainB, chainC, state)
	require.NoError(t, err)

	aChainToken, _, cChainToken, _, err := testhelpers.DeployTransferableToken(
		lggr,
		tenv.Env.BlockChains.EVMChains(),
		chainA,
		chainC,
		ownerChainA,
		ownerChainC,
		state,
		e.ExistingAddresses,
		"MY_TOKEN",
	)
	require.NoError(t, err)

	// Add all lanes
	testhelpers.AddLanesForAll(t, &tenv, state)

	testhelpers.MintAndAllow(
		t,
		e,
		state,
		map[uint64][]testhelpers.MintTokenInfo{
			chainA: {
				testhelpers.NewMintTokenInfo(ownerChainA, aChainUSDC, aChainToken),
			},
			chainB: {
				testhelpers.NewMintTokenInfo(ownerChainB, bChainUSDC),
			},
			chainC: {
				testhelpers.NewMintTokenInfo(ownerChainC, cChainUSDC, cChainToken),
			},
		},
	)

	err = updateFeeQuoters(t, lggr, e, state, chainA, chainB, chainC, aChainUSDC, bChainUSDC, cChainUSDC)
	require.NoError(t, err)

	// MockE2EUSDCTransmitter always mint 1, see MockE2EUSDCTransmitter.sol for more details
	tinyOneCoin := new(big.Int).SetUint64(1)

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:        "single USDC token transfer to EOA",
			Receiver:    utils.RandomAddress().Bytes(),
			SourceChain: chainC,
			DestChain:   chainA,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  cChainUSDC.Address(),
					Amount: tinyOneCoin,
				}},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: aChainUSDC.Address().Bytes(), Amount: tinyOneCoin},
			},
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		},
		{
			Name:        "multiple USDC tokens within the same message",
			Receiver:    utils.RandomAddress().Bytes(),
			SourceChain: chainC,
			DestChain:   chainA,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  cChainUSDC.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  cChainUSDC.Address(),
					Amount: tinyOneCoin,
				},
			},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				// 2 coins because of the same Receiver
				{Token: aChainUSDC.Address().Bytes(), Amount: new(big.Int).Add(tinyOneCoin, tinyOneCoin)},
			},
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		},
		{
			Name:        "USDC token together with another token transferred to EOA",
			Receiver:    utils.RandomAddress().Bytes(),
			SourceChain: chainA,
			DestChain:   chainC,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  aChainUSDC.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  aChainToken.Address(),
					Amount: new(big.Int).Mul(tinyOneCoin, big.NewInt(10)),
				},
			},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: cChainUSDC.Address().Bytes(), Amount: tinyOneCoin},
				{Token: cChainToken.Address().Bytes(), Amount: new(big.Int).Mul(tinyOneCoin, big.NewInt(10))},
			},
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		},
		{
			Name:        "USDC programmable token transfer to valid contract receiver",
			Receiver:    state.MustGetEVMChainState(chainC).Receiver.Address().Bytes(),
			SourceChain: chainA,
			DestChain:   chainC,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  aChainUSDC.Address(),
					Amount: tinyOneCoin,
				},
			},
			Data: []byte("hello world"),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: cChainUSDC.Address().Bytes(), Amount: tinyOneCoin},
			},
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		},
		{
			Name:        "USDC programmable token transfer with too little gas",
			Receiver:    state.MustGetEVMChainState(chainB).Receiver.Address().Bytes(),
			SourceChain: chainC,
			DestChain:   chainB,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  cChainUSDC.Address(),
					Amount: tinyOneCoin,
				},
			},
			Data: []byte("gimme more gas to execute that!"),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: bChainUSDC.Address().Bytes(), Amount: new(big.Int).SetUint64(0)},
			},
			ExtraArgs:      testhelpers.MakeEVMExtraArgsV2(1, false),
			ExpectedStatus: testhelpers.EXECUTION_STATE_FAILURE,
		},
		{
			Name:        "USDC token transfer from a different source chain",
			Receiver:    utils.RandomAddress().Bytes(),
			SourceChain: chainB,
			DestChain:   chainC,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  bChainUSDC.Address(),
					Amount: tinyOneCoin,
				},
			},
			Data: nil,
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: cChainUSDC.Address().Bytes(), Amount: tinyOneCoin},
			},
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
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

func updateFeeQuoters(
	t *testing.T,
	lggr logger.Logger,
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	chainA, chainB, chainC uint64,
	aChainUSDC, bChainUSDC, cChainUSDC *burn_mint_erc677.BurnMintERC677,
) error {
	evmChains := e.BlockChains.EVMChains()
	updateFeeQtrGrp := errgroup.Group{}
	updateFeeQtrGrp.Go(func() error {
		return testhelpers.UpdateFeeQuoterForUSDC(t, e, lggr, evmChains[chainA], chainC)
	})
	updateFeeQtrGrp.Go(func() error {
		return testhelpers.UpdateFeeQuoterForUSDC(t, e, lggr, evmChains[chainB], chainC)
	})
	updateFeeQtrGrp.Go(func() error {
		err1 := testhelpers.UpdateFeeQuoterForUSDC(t, e, lggr, evmChains[chainC], chainA)
		if err1 != nil {
			return err1
		}
		return testhelpers.UpdateFeeQuoterForUSDC(t, e, lggr, evmChains[chainC], chainB)
	})
	return updateFeeQtrGrp.Wait()
}
