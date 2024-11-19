package smoke

import (
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"golang.org/x/exp/maps"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestUSDCTokenTransfer(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	tenv, cluster, _ := testsetups.NewLocalDevEnvironmentWithDefaultPrice(t, lggr)

	var endpoint string
	// When inmemory env then spin up in memory mock server
	if cluster == nil {
		server := mockAttestationResponse()
		defer server.Close()
		endpoint = server.URL
	} else {
		err := actions.SetMockServerWithUSDCAttestation(tenv.Env.MockAdapter, nil)
		require.NoError(t, err)
		endpoint = tenv.Env.MockAdapter.InternalEndpoint
	}

	e := tenv.Env
	state, err := ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Chains)
	sourceChain := allChainSelectors[0]
	destChain := allChainSelectors[1]

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	tokenConfig := ccdeploy.NewTokenConfig()
	tokenConfig.UpsertTokenInfo(ccdeploy.LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: cciptypes.UnknownEncodedAddress(feeds[ccdeploy.LinkSymbol].Address().String()),
			Decimals:          ccdeploy.LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)

	output, err := changeset.DeployPrerequisites(e, changeset.DeployPrerequisiteConfig{
		ChainSelectors: e.AllChainSelectors(),
	})
	require.NoError(t, err)
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(output.AddressBook))

	// Apply migration
	output, err = changeset.InitialDeploy(e, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		ChainsToDeploy: e.AllChainSelectors(),
		TokenConfig:    tokenConfig,
		MCMSConfig:     ccdeploy.NewTestMCMSConfig(t, e),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
		USDCConfig: ccdeploy.USDCConfig{
			Enabled: true,
			USDCAttestationConfig: ccdeploy.USDCAttestationConfig{
				API:         endpoint,
				APITimeout:  commonconfig.MustNewDuration(time.Second),
				APIInterval: commonconfig.MustNewDuration(500 * time.Millisecond),
			},
		},
	})
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook))

	state, err = ccdeploy.LoadOnchainState(e)
	require.NoError(t, err)

	srcUSDC, dstUSDC, err := ccdeploy.ConfigureUSDCTokenPools(lggr, e.Chains, sourceChain, destChain, state)
	require.NoError(t, err)

	srcToken, _, dstToken, _, err := ccdeploy.DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		sourceChain,
		destChain,
		state,
		e.ExistingAddresses,
		"MY_TOKEN",
	)
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

	mintAndAllow(t, e, state, map[uint64][]*burn_mint_erc677.BurnMintERC677{
		sourceChain: {srcUSDC, srcToken},
		destChain:   {dstUSDC, dstToken},
	})

	err = ccdeploy.UpdateFeeQuoterForUSDC(lggr, e.Chains[sourceChain], state.Chains[sourceChain], destChain, srcUSDC)
	require.NoError(t, err)

	err = ccdeploy.UpdateFeeQuoterForUSDC(lggr, e.Chains[destChain], state.Chains[destChain], sourceChain, dstUSDC)
	require.NoError(t, err)

	// MockE2EUSDCTransmitter always mint 1, see MockE2EUSDCTransmitter.sol for more details
	tinyOneCoin := new(big.Int).SetUint64(1)

	tcs := []struct {
		name                  string
		receiver              common.Address
		sourceChain           uint64
		destChain             uint64
		tokens                []router.ClientEVMTokenAmount
		data                  []byte
		expectedTokenBalances map[common.Address]*big.Int
	}{
		{
			name:        "single USDC token transfer to EOA",
			receiver:    utils.RandomAddress(),
			sourceChain: destChain,
			destChain:   sourceChain,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  dstUSDC.Address(),
					Amount: tinyOneCoin,
				}},
			expectedTokenBalances: map[common.Address]*big.Int{
				srcUSDC.Address(): tinyOneCoin,
			},
		},
		{
			name:        "multiple USDC tokens within the same message",
			receiver:    utils.RandomAddress(),
			sourceChain: destChain,
			destChain:   sourceChain,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  dstUSDC.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  dstUSDC.Address(),
					Amount: tinyOneCoin,
				},
			},
			expectedTokenBalances: map[common.Address]*big.Int{
				// 2 coins because of the same receiver
				srcUSDC.Address(): new(big.Int).Add(tinyOneCoin, tinyOneCoin),
			},
		},
		{
			name:        "USDC token together with another token transferred to EOA",
			receiver:    utils.RandomAddress(),
			sourceChain: sourceChain,
			destChain:   destChain,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  srcUSDC.Address(),
					Amount: tinyOneCoin,
				},
				{
					Token:  srcToken.Address(),
					Amount: new(big.Int).Mul(tinyOneCoin, big.NewInt(10)),
				},
			},
			expectedTokenBalances: map[common.Address]*big.Int{
				dstUSDC.Address():  tinyOneCoin,
				dstToken.Address(): new(big.Int).Mul(tinyOneCoin, big.NewInt(10)),
			},
		},
		{
			name:        "programmable token transfer to valid contract receiver",
			receiver:    state.Chains[destChain].Receiver.Address(),
			sourceChain: sourceChain,
			destChain:   destChain,
			tokens: []router.ClientEVMTokenAmount{
				{
					Token:  srcUSDC.Address(),
					Amount: tinyOneCoin,
				},
			},
			data: []byte("hello world"),
			expectedTokenBalances: map[common.Address]*big.Int{
				dstUSDC.Address(): tinyOneCoin,
			},
		},
	}

	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			initialBalances := map[common.Address]*big.Int{}
			for token := range tt.expectedTokenBalances {
				initialBalance := getTokenBalance(t, token, tt.receiver, e.Chains[tt.destChain])
				initialBalances[token] = initialBalance
			}

			transferAndWaitForSuccess(
				t,
				e,
				state,
				tt.sourceChain,
				tt.destChain,
				tt.tokens,
				tt.receiver,
				tt.data,
			)

			for token, balance := range tt.expectedTokenBalances {
				expected := new(big.Int).Add(initialBalances[token], balance)
				waitForTheTokenBalance(t, token, tt.receiver, e.Chains[tt.destChain], expected)
			}
		})
	}
}

// mintAndAllow mints tokens for deployers and allow router to spend them
func mintAndAllow(
	t *testing.T,
	e deployment.Environment,
	state ccdeploy.CCIPOnChainState,
	tkMap map[uint64][]*burn_mint_erc677.BurnMintERC677,
) {
	for chain, tokens := range tkMap {
		for _, token := range tokens {
			twoCoins := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(2))

			tx, err := token.Mint(
				e.Chains[chain].DeployerKey,
				e.Chains[chain].DeployerKey.From,
				new(big.Int).Mul(twoCoins, big.NewInt(10)),
			)
			require.NoError(t, err)
			_, err = e.Chains[chain].Confirm(tx)
			require.NoError(t, err)

			tx, err = token.Approve(e.Chains[chain].DeployerKey, state.Chains[chain].Router.Address(), twoCoins)
			require.NoError(t, err)
			_, err = e.Chains[chain].Confirm(tx)
			require.NoError(t, err)
		}
	}
}

// transferAndWaitForSuccess sends a message from sourceChain to destChain and waits for it to be executed
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

	msgSentEvent := ccdeploy.TestSendRequest(t, env, state, sourceChain, destChain, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(receiver.Bytes(), 32),
		Data:         data,
		TokenAmounts: tokens,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	expectedSeqNum[destChain] = msgSentEvent.SequenceNumber

	// Wait for all commit reports to land.
	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, env, state, expectedSeqNum, startBlocks)

	// Wait for all exec reports to land
	ccdeploy.ConfirmExecWithSeqNrForAll(t, env, state, expectedSeqNum, startBlocks)
}

// mockAttestationResponse mocks the USDC attestation server, it returns random Attestation.
// We don't need to return exactly the same attestation, because our Mocked USDC contract doesn't rely on any specific
// value, but instead of that it just checks if the attestation is present. Therefore, it makes the test a bit simpler
// and doesn't require very detailed mocks. Please see tests in chainlink-ccip for detailed tests using real attestations
func mockAttestationResponse() *httptest.Server {
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
	return server
}

func waitForTheTokenBalance(
	t *testing.T,
	token common.Address,
	receiver common.Address,
	chain deployment.Chain,
	expected *big.Int,
) {
	tokenContract, err := burn_mint_erc677.NewBurnMintERC677(token, chain.Client)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		actualBalance, err := tokenContract.BalanceOf(&bind.CallOpts{Context: tests.Context(t)}, receiver)
		require.NoError(t, err)

		t.Log("Waiting for the token balance",
			"expected", expected,
			"actual", actualBalance,
			"token", token,
			"receiver", receiver,
		)

		return actualBalance.Cmp(expected) == 0
	}, tests.WaitTimeout(t), 100*time.Millisecond)
}

func getTokenBalance(
	t *testing.T,
	token common.Address,
	receiver common.Address,
	chain deployment.Chain,
) *big.Int {
	tokenContract, err := burn_mint_erc677.NewBurnMintERC677(token, chain.Client)
	require.NoError(t, err)

	balance, err := tokenContract.BalanceOf(&bind.CallOpts{Context: tests.Context(t)}, receiver)
	require.NoError(t, err)

	t.Log("Getting token balance",
		"actual", balance,
		"token", token,
		"receiver", receiver,
	)

	return balance
}
