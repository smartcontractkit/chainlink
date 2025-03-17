package ccip

import (
	"math/big"
	"testing"

	"golang.org/x/exp/maps"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestTokenTransfer(t *testing.T) {
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
			Receiver: utils.RandomAddress(),
			ExpectedTokenBalances: map[common.Address]*big.Int{
				destToken.Address(): oneE18,
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
			Receiver: state.Chains[destChain].Receiver.Address(),
			ExpectedTokenBalances: map[common.Address]*big.Int{
				destToken.Address(): oneE18,
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
			Receiver:  state.Chains[sourceChain].Receiver.Address(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(300_000, false),
			ExpectedTokenBalances: map[common.Address]*big.Int{
				selfServeSrcToken.Address(): new(big.Int).Add(oneE18, oneE18),
				srcToken.Address():          oneE18,
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
			Receiver:  utils.RandomAddress(),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(1, false),
			ExpectedTokenBalances: map[common.Address]*big.Int{
				selfServeSrcToken.Address(): oneE18,
				srcToken.Address():          new(big.Int).Add(oneE18, oneE18),
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
			Receiver:  state.Chains[sourceChain].Receiver.Address(),
			Data:      []byte("this should be reverted because gasLimit is too low, no tokens are transferred as well"),
			ExtraArgs: testhelpers.MakeEVMExtraArgsV2(1, false),
			ExpectedTokenBalances: map[common.Address]*big.Int{
				selfServeSrcToken.Address(): big.NewInt(0),
				srcToken.Address():          big.NewInt(0),
			},
			ExpectedStatus: testhelpers.EXECUTION_STATE_FAILURE,
		},
	}

	startBlocks, expectedSeqNums, expectedExecutionStates, expectedTokenBalances :=
		testhelpers.TransferMultiple(ctx, t, e, state, tcs)

	err = testhelpers.ConfirmMultipleCommits(
		t,
		e.Chains,
		state.Chains,
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

	testhelpers.WaitForTokenBalances(ctx, t, e.Chains, expectedTokenBalances)
}
