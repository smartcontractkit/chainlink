package smoke

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"golang.org/x/exp/maps"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestUSDCTokenTransfer(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	tenv := ccdeploy.NewMemoryEnvironment(t, lggr, 2, 4)

	e := tenv.Env
	state, err := ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Chains)
	sourceChain := allChainSelectors[0]
	destChain := allChainSelectors[1]

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"status": "complete",
			"attestation": "0x9049623e91719ef2aa63c55f357be2529b0e7122ae552c18aff8db58b4633c4d3920ff03d3a6d1ddf11f06bf64d7fd60d45447ac81f527ba628877dc5ca759651b08ffae25a6d3b1411749765244f0a1c131cbfe04430d687a2e12fd9d2e6dc08e118ad95d94ad832332cf3c4f7a4f3da0baa803b7be024b02db81951c0f0714de1b"
		}`

		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}))
	defer server.Close()

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	tokenConfig := ccdeploy.NewTokenConfig()
	tokenConfig.UpsertTokenInfo(ccdeploy.LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: cciptypes.UnknownEncodedAddress(feeds[ccdeploy.LinkSymbol].Address().String()),
			Decimals:          ccdeploy.LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)

	// Apply migration
	output, err := changeset.InitialDeploy(e, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		ChainsToDeploy: e.AllChainSelectors(),
		TokenConfig:    tokenConfig,
		MCMSConfig:     ccdeploy.NewTestMCMSConfig(t, e),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
		USDCConfig: ccdeploy.USDCConfig{
			Enabled: true,
			USDCAttestationConfig: ccdeploy.USDCAttestationConfig{
				API:         server.URL,
				APITimeout:  commonconfig.MustNewDuration(time.Second),
				APIInterval: commonconfig.MustNewDuration(500 * time.Millisecond),
			},
		},
	})
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook))
	// Get new state after migration and mock USDC token deployment.
	state, err = ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	srcUSDC, dstUSDC, err := ccdeploy.ConfigureUSDCTokenPools(lggr, e.Chains, sourceChain, destChain, state)
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	ccdeploy.ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)

	// Apply the jobs.
	for nodeID, jobs := range output.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}

	// Add all lanes
	require.NoError(t, ccdeploy.AddLanesForAll(e, state))

	twoCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2))
	tx, err := srcUSDC.Mint(
		e.Chains[sourceChain].DeployerKey,
		e.Chains[sourceChain].DeployerKey.From,
		new(big.Int).Mul(twoCoins, big.NewInt(10)),
	)
	require.NoError(t, err)
	_, err = e.Chains[sourceChain].Confirm(tx)
	require.NoError(t, err)

	tx, err = dstUSDC.Mint(
		e.Chains[destChain].DeployerKey,
		e.Chains[destChain].DeployerKey.From,
		new(big.Int).Mul(twoCoins, big.NewInt(10)),
	)
	require.NoError(t, err)
	_, err = e.Chains[destChain].Confirm(tx)
	require.NoError(t, err)

	tx, err = srcUSDC.Approve(e.Chains[sourceChain].DeployerKey, state.Chains[sourceChain].Router.Address(), twoCoins)
	require.NoError(t, err)
	_, err = e.Chains[sourceChain].Confirm(tx)
	require.NoError(t, err)
	tx, err = dstUSDC.Approve(e.Chains[destChain].DeployerKey, state.Chains[destChain].Router.Address(), twoCoins)
	require.NoError(t, err)
	_, err = e.Chains[destChain].Confirm(tx)
	require.NoError(t, err)

	err = ccdeploy.UpdateFeeQuoterForUSDC(e.Chains[sourceChain], state.Chains[sourceChain], destChain, srcUSDC)
	require.NoError(t, err)

	err = ccdeploy.UpdateFeeQuoterForUSDC(e.Chains[destChain], state.Chains[destChain], sourceChain, dstUSDC)
	require.NoError(t, err)

	tinyOneCoin := new(big.Int).SetUint64(1)

	//tokens := map[uint64][]router.ClientEVMTokenAmount{
	//	sourceChain: {{
	//		Token:  srcUSDC.Address(),
	//		Amount: tinyOneCoin,
	//	}},
	//	destChain: {{
	//		Token:  dstUSDC.Address(),
	//		Amount: tinyOneCoin,
	//	}},
	//}

	//t.Run("single USDC token transfer to EOA", func(t *testing.T) {
	//	var (
	//		receiver = common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32)
	//		data     = []byte("hello world")
	//		feeToken = common.HexToAddress("0x0")
	//	)
	//
	//	latesthdr, err := dChain.Client.HeaderByNumber(testcontext.Get(t), nil)
	//	require.NoError(t, err)
	//	block := latesthdr.Number.Uint64()
	//	startBlocks[dest] = &block
	//
	//	if src == sourceChain && dest == destChain {
	//		seqNum := ccdeploy.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
	//			Receiver:     receiver,
	//			Data:         data,
	//			TokenAmounts: tokens[src],
	//			FeeToken:     feeToken,
	//			ExtraArgs:    nil,
	//		})
	//		expectedSeqNum[dest] = seqNum
	//	}
	//
	//	// Wait for all commit reports to land.
	//	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)
	//
	//	// Wait for all exec reports to land
	//	ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)
	//
	//	balance, err := dstUSDC.BalanceOf(nil, state.Chains[destChain].Receiver.Address())
	//	require.NoError(t, err)
	//	require.Equal(t, tinyOneCoin, balance)
	//})
	//
	//t.Run("USDC token transfer from different sources to the same destination", func(t *testing.T) {
	//
	//})
	//
	t.Run("programmable token transfer with USDC to contract", func(t *testing.T) {
		transferAndWaitForSuccess(
			t,
			e,
			state,
			sourceChain,
			destChain,
			[]router.ClientEVMTokenAmount{
				{
					Token:  srcUSDC.Address(),
					Amount: tinyOneCoin,
				}},
			state.Chains[destChain].Receiver.Address(),
			[]byte("hello world"),
		)

		balance, err := dstUSDC.BalanceOf(&bind.CallOpts{Context: tests.Context(t)}, state.Chains[destChain].Receiver.Address())
		require.NoError(t, err)
		require.Equal(t, tinyOneCoin, balance)
	})

	//
	//t.Run("USDC token together with other token transfer", func(t *testing.T) {
	//
	//})
}

func transferAndWaitForSuccess(
	t *testing.T,
	env deployment.Environment,
	state ccdeploy.CCIPOnChainState,
	sourceChain, destChain uint64,
	tokens []router.ClientEVMTokenAmount,
	receiver common.Address,
	data []byte,
) {
	startBlocks := make(map[uint64]*uint64)
	expectedSeqNum := make(map[uint64]uint64)

	latesthdr, err := env.Chains[destChain].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	block := latesthdr.Number.Uint64()
	startBlocks[destChain] = &block

	seqNum := ccdeploy.TestSendRequest(t, env, state, sourceChain, destChain, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(receiver.Bytes(), 32),
		Data:         data,
		TokenAmounts: tokens,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	expectedSeqNum[destChain] = seqNum

	// Wait for all commit reports to land.
	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, env, state, expectedSeqNum, startBlocks)

	// Wait for all exec reports to land
	ccdeploy.ConfirmExecWithSeqNrForAll(t, env, state, expectedSeqNum, startBlocks)
}
