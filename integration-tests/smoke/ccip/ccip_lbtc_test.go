package ccip

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

/*
* Chain topology for this test
* 	chainA (LBTC, MY_TOKEN)
*			|
*			| ------- chainC (LBTC, MY_TOKEN)
*			|
* 	chainB (LBTC)
 */
func TestLBTCTokenTransfer(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()
	tenv, _, _ := testsetups.NewIntegrationEnvironment(t,
		testhelpers.WithNumOfUsersPerChain(3),
		testhelpers.WithNumOfChains(3),
		testhelpers.WithLBTC(),
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

	aChainLBTC, cChainLBTC, err := testhelpers.ConfigureLBTCTokenPools(lggr, evmChains, chainA, chainC, state)
	require.NoError(t, err)

	bChainLBTC, _, err := testhelpers.ConfigureLBTCTokenPools(lggr, evmChains, chainB, chainC, state)
	require.NoError(t, err)

	aChainToken, _, cChainToken, _, err := testhelpers.DeployTransferableToken(
		lggr,
		tenv.Env.BlockChains.EVMChains(),
		chainA,
		chainC,
		ownerChainA,
		ownerChainC,
		state,
		e.ExistingAddresses, //nolint:staticcheck // Addressbook is deprecated, but we still use it for the time being
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
				testhelpers.NewMintTokenInfo(ownerChainA, aChainLBTC, aChainToken),
			},
			chainB: {
				testhelpers.NewMintTokenInfo(ownerChainB, bChainLBTC),
			},
			chainC: {
				testhelpers.NewMintTokenInfo(ownerChainC, cChainLBTC, cChainToken),
			},
		},
	)

	updateFeeQtrGrp := errgroup.Group{}
	for _, chainSel1 := range allChainSelectors {
		updateFeeQtrGrp.Go(func() error {
			for _, chainSel2 := range allChainSelectors {
				if chainSel1 == chainSel2 {
					continue
				}
				if err := testhelpers.UpdateFeeQuoterForToken(t, e, lggr, evmChains[chainSel1], chainSel2, shared.LBTCSymbol); err != nil {
					return err
				}
			}
			return nil
		})
	}
	err = updateFeeQtrGrp.Wait()
	require.NoError(t, err)

	tinyOneCoin := new(big.Int).SetUint64(1)

	tcs := []testhelpers.TestTransferRequest{
		{
			Name:        "single LBTC token transfer to EOA",
			Receiver:    utils.RandomAddress().Bytes(),
			SourceChain: chainC,
			DestChain:   chainA,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  cChainLBTC.Address(),
					Amount: tinyOneCoin,
				}},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: aChainLBTC.Address().Bytes(), Amount: tinyOneCoin},
			},
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		},
		{
			Name:        "multiple LBTC tokens within the same message",
			Receiver:    utils.RandomAddress().Bytes(),
			SourceChain: chainC,
			DestChain:   chainA,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  cChainLBTC.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  cChainLBTC.Address(),
					Amount: tinyOneCoin,
				},
			},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				// 2 coins because of the same Receiver
				{Token: aChainLBTC.Address().Bytes(), Amount: new(big.Int).Add(tinyOneCoin, tinyOneCoin)},
			},
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		},
		{
			Name:        "LBTC token together with another token transferred to EOA",
			Receiver:    utils.RandomAddress().Bytes(),
			SourceChain: chainA,
			DestChain:   chainC,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  aChainLBTC.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  aChainToken.Address(),
					Amount: new(big.Int).Mul(tinyOneCoin, big.NewInt(10)),
				},
			},
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: cChainLBTC.Address().Bytes(), Amount: tinyOneCoin},
				{Token: cChainToken.Address().Bytes(), Amount: new(big.Int).Mul(tinyOneCoin, big.NewInt(10))},
			},
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		},
		{
			Name:        "LBTC programmable token transfer to valid contract receiver",
			Receiver:    state.MustGetEVMChainState(chainC).Receiver.Address().Bytes(),
			SourceChain: chainA,
			DestChain:   chainC,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  aChainLBTC.Address(),
					Amount: tinyOneCoin,
				},
			},
			Data: []byte("hello world"),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: cChainLBTC.Address().Bytes(), Amount: tinyOneCoin},
			},
			ExpectedStatus: testhelpers.EXECUTION_STATE_SUCCESS,
		},
		{
			Name:        "LBTC programmable token transfer with too little gas",
			Receiver:    state.MustGetEVMChainState(chainB).Receiver.Address().Bytes(),
			SourceChain: chainC,
			DestChain:   chainB,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  cChainLBTC.Address(),
					Amount: tinyOneCoin,
				},
			},
			Data: []byte("gimme more gas to execute that!"),
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: bChainLBTC.Address().Bytes(), Amount: new(big.Int).SetUint64(0)},
			},
			ExtraArgs:      testhelpers.MakeEVMExtraArgsV2(1, false),
			ExpectedStatus: testhelpers.EXECUTION_STATE_FAILURE,
		},
		{
			Name:        "LBTC token transfer from a different source chain",
			Receiver:    utils.RandomAddress().Bytes(),
			SourceChain: chainB,
			DestChain:   chainC,
			Tokens: []router.ClientEVMTokenAmount{
				{
					Token:  bChainLBTC.Address(),
					Amount: tinyOneCoin,
				},
			},
			Data: nil,
			ExpectedTokenBalances: []testhelpers.ExpectedBalance{
				{Token: cChainLBTC.Address().Bytes(), Amount: tinyOneCoin},
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
