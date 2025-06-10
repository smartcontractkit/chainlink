package ccip

import (
	"math/big"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/common"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

// Test_AddChainE2E tests the end-to-end functionality of adding a new chain to the existing
// chains using two consolidated changesets.
// 1. AddCandidatesForNewChainChangeset
// 2. PromoteNewChainForConfigChangeset
// Scenario:
// - Deploy 4 chains: home, feed, third, and fourth.
// - Home and feed chains are already deployed and connected to each other.
// - Third chain is connected to home and feed chains using consolidated changesets.
// - Fourth chain is connected to the third chain using consolidated changesets.
// - Messages are sent between all valid lanes to verify the setup.
func Test_AddChainE2E(t *testing.T) {
	e, _, tEnv := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(4),
		testhelpers.WithPrerequisiteDeploymentOnly(nil),
	)
	initialSetToDeploy := []uint64{e.HomeChainSel, e.FeedChainSel}
	remainingChains := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM), chain.WithChainSelectorsExclusion(initialSetToDeploy))
	thirdChain := remainingChains[0]
	fourthChain := remainingChains[1]

	e = testhelpers.AddCCIPContractsToEnvironment(t, initialSetToDeploy, tEnv, false)
	// Need to update what the RMNProxy is pointing to, otherwise plugin will not work.
	var err error
	e.Env, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.SetRMNRemoteOnRMNProxyChangeset),
			v1_6.SetRMNRemoteOnRMNProxyConfig{
				ChainSelectors: initialSetToDeploy,
			},
		),
	)
	require.NoError(t, err)
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	tEnv.UpdateDeployedEnvironment(e)

	// wire up initial set of chains all - all
	for _, src := range initialSetToDeploy {
		for _, dest := range initialSetToDeploy {
			if src != dest {
				testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(
					t,
					&e,
					state,
					src,
					dest,
					false,
				)
			}
		}
	}

	testhelpers.SleepAndReplay(t, e.Env, 10*time.Second, initialSetToDeploy...)

	// need donIDClaimer contract to be deployed before we can deploy the new chain
	e.Env, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(
			v1_6.DeployDonIDClaimerChangeset,
			v1_6.DeployDonIDClaimerConfig{},
		))
	require.NoError(t, err, "must deploy donIDClaimer contract")

	// Transfer ownership of the home and feed chain contracts to the timelock
	e.Env, err = TransferOwnership(t, e.Env, e.HomeChainSel, initialSetToDeploy, state)
	require.NoError(t, err, "must transfer ownership of home %d and feed "+
		"chain %d contracts to the timelock", e.HomeChainSel, e.FeedChainSel)

	// setup the third chain with home and feed chain
	e.Env = SetupNewChain(t, e.HomeChainSel, e.FeedChainSel, thirdChain, initialSetToDeploy, e.Env, state)

	// setup the fourth chain with third chain alone
	e.Env = SetupNewChain(t, e.HomeChainSel, e.FeedChainSel, fourthChain, []uint64{thirdChain}, e.Env, state)

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	testhelpers.SleepAndReplay(t, e.Env, 30*time.Second, e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))...)
	SendMsgs(t,
		e.Env,
		[]testhelpers.SourceDestPair{
			{
				SourceChainSelector: e.HomeChainSel,
				DestChainSelector:   thirdChain,
			},
			{
				SourceChainSelector: e.FeedChainSel,
				DestChainSelector:   thirdChain,
			},
			{
				SourceChainSelector: thirdChain,
				DestChainSelector:   e.FeedChainSel,
			},
			{
				SourceChainSelector: thirdChain,
				DestChainSelector:   e.HomeChainSel,
			},
			{
				SourceChainSelector: thirdChain,
				DestChainSelector:   fourthChain,
			},
			{
				SourceChainSelector: fourthChain,
				DestChainSelector:   thirdChain,
			},
			{
				SourceChainSelector: e.HomeChainSel,
				DestChainSelector:   e.FeedChainSel,
			},
			{
				SourceChainSelector: e.FeedChainSel,
				DestChainSelector:   e.HomeChainSel,
			},
		},
		state,
		false,
	)
}

func SetupNewChain(
	t *testing.T,
	homeChain, feedChain, chainToDeploy uint64,
	remoteChains []uint64,
	env cldf.Environment,
	state stateview.CCIPOnChainState,
) cldf.Environment {
	nodeInfo, err := deployment.NodeInfo(env.NodeIDs, env.Offchain)
	require.NoError(t, err, "must get node info")
	mcmsDeploymentCfg := proposalutils.SingleGroupTimelockConfigV2(t)
	tokenConfig := shared.NewTestTokenConfig(state.MustGetEVMChainState(feedChain).USDFeeds)

	// Build remote chain configurations
	remoteChainsDefinition := make([]v1_6.ChainDefinition, len(remoteChains))
	for i, selector := range remoteChains {
		remoteChainsDefinition[i] = v1_6.ChainDefinition{
			ConnectionConfig: v1_6.ConnectionConfig{
				RMNVerificationDisabled: true,
				AllowListEnabled:        false,
			},
			Selector: selector,
			GasPrice: testhelpers.DefaultGasPrice,
			TokenPrices: map[common.Address]*big.Int{
				state.Chains[selector].LinkToken.Address(): testhelpers.DefaultLinkPrice,
				state.Chains[selector].Weth9.Address():     testhelpers.DefaultWethPrice,
			},
			FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
		}
	}

	newChainDefinition := v1_6.NewChainDefinition{
		ChainDefinition: v1_6.ChainDefinition{
			ConnectionConfig: v1_6.ConnectionConfig{
				RMNVerificationDisabled: true,
				AllowListEnabled:        false,
			},
			Selector: chainToDeploy,
			GasPrice: testhelpers.DefaultGasPrice,
			TokenPrices: map[common.Address]*big.Int{
				state.Chains[chainToDeploy].LinkToken.Address(): testhelpers.DefaultLinkPrice,
				state.Chains[chainToDeploy].Weth9.Address():     testhelpers.DefaultWethPrice,
			},
			FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
		},
		ChainContractParams: ccipseq.ChainContractParams{
			FeeQuoterParams: ccipops.DefaultFeeQuoterParams(),
			OffRampParams:   ccipops.DefaultOffRampParams(),
		},
		ConfigOnHome: v1_6.ChainConfig{
			Readers: nodeInfo.NonBootstraps().PeerIDs(),
			FChain:  uint8(len(nodeInfo.NonBootstraps().PeerIDs()) / 3), // #nosec G115 - Overflow is not a concern in this test scenario
			EncodableChainConfig: chainconfig.ChainConfig{
				GasPriceDeviationPPB:      cciptypes.BigInt{Int: big.NewInt(testhelpers.DefaultGasPriceDeviationPPB)},
				DAGasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(testhelpers.DefaultDAGasPriceDeviationPPB)},
				OptimisticConfirmations:   globals.OptimisticConfirmations,
				ChainFeeDeviationDisabled: false,
			},
		},
		CommitOCRParams: v1_6.DeriveOCRParamsForCommit(
			v1_6.SimulationTest,
			feedChain,
			tokenConfig.GetTokenInfo(env.Logger,
				state.MustGetEVMChainState(chainToDeploy).LinkToken.Address(),
				state.MustGetEVMChainState(chainToDeploy).Weth9.Address(),
			),
			func(ocrParams v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
				ocrParams.CommitOffChainConfig.RMNEnabled = false
				return ocrParams
			},
		),
		ExecOCRParams: v1_6.DeriveOCRParamsForExec(v1_6.SimulationTest, nil, nil),
	}

	// Add candidate for new chain using AddCandidatesForNewChainChangeset
	env, err = commonchangeset.Apply(t, env,
		commonchangeset.Configure(
			v1_6.AddCandidatesForNewChainChangeset,
			v1_6.AddCandidatesForNewChainConfig{
				HomeChainSelector:    homeChain,
				FeedChainSelector:    feedChain,
				NewChain:             newChainDefinition,
				RemoteChains:         remoteChainsDefinition,
				MCMSDeploymentConfig: &mcmsDeploymentCfg,
				MCMSConfig: &proposalutils.TimelockConfig{
					MinDelay:   0 * time.Second,
					MCMSAction: mcmstypes.TimelockActionSchedule,
				},
				DonIDOffSet: nil,
			},
		),
	)
	require.NoError(t, err, "must apply AddCandidatesForNewChainChangeset, "+
		"new chain: %d, remote chains: %v", chainToDeploy, remoteChains)

	// Apply PromoteNewChainForConfigChangeset
	env, err = commonchangeset.Apply(t, env,
		commonchangeset.Configure(
			v1_6.PromoteNewChainForConfigChangeset,
			v1_6.PromoteNewChainForConfig{
				HomeChainSelector: homeChain,
				NewChain:          newChainDefinition,
				RemoteChains:      remoteChainsDefinition,
				TestRouter:        pointer.ToBool(false),
				MCMSConfig: &proposalutils.TimelockConfig{
					MinDelay:   0 * time.Second,
					MCMSAction: mcmstypes.TimelockActionSchedule,
				},
			},
		),
	)
	require.NoError(t, err, "must apply PromoteNewChainForConfigChangeset, "+
		"new chain: %d, remote chains: %v", chainToDeploy, remoteChains)

	return env
}

func SendMsgs(
	t *testing.T,
	env cldf.Environment,
	sourceDestPairs []testhelpers.SourceDestPair,
	state stateview.CCIPOnChainState,
	testRouter bool,
) {
	t.Helper()

	var (
		startBlocks        = make(map[uint64]*uint64)
		expectedSeqNum     = make(map[testhelpers.SourceDestPair]uint64)
		expectedSeqNumExec = make(map[testhelpers.SourceDestPair][]uint64)
	)
	evmChains := env.BlockChains.EVMChains()
	for _, pair := range sourceDestPairs {
		latesthdr, err := evmChains[pair.DestChainSelector].Client.HeaderByNumber(testcontext.Get(t), nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		msgSentEvent := testhelpers.TestSendRequest(
			t,
			env,
			state,
			pair.SourceChainSelector,
			pair.DestChainSelector,
			testRouter,
			router.ClientEVM2AnyMessage{
				Receiver: common.LeftPadBytes(
					state.Chains[pair.DestChainSelector].Receiver.Address().Bytes(),
					32,
				),
				Data:         []byte("hello world"),
				TokenAmounts: nil,
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			})
		startBlocks[pair.DestChainSelector] = &block
		expectedSeqNum[pair] = msgSentEvent.SequenceNumber
		expectedSeqNumExec[pair] = append(expectedSeqNumExec[pair], msgSentEvent.SequenceNumber)
	}
	testhelpers.SleepAndReplay(t, env, 10*time.Second, env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))...)
	testhelpers.ConfirmCommitForAllWithExpectedSeqNums(t, env, state,
		testhelpers.ToSeqRangeMap(expectedSeqNum), startBlocks)
	testhelpers.ConfirmExecWithSeqNrsForAll(t, env, state, expectedSeqNumExec, startBlocks)
}

func TransferOwnership(
	t *testing.T,
	env cldf.Environment,
	homeChainSelector uint64,
	initialSetToDeploy []uint64,
	state stateview.CCIPOnChainState,
) (cldf.Environment, error) {
	t.Helper()

	// Transfer home and feed chain contracts ownership to the timelock
	contractsToTransfer := make(map[uint64][]common.Address, len(initialSetToDeploy))
	for _, selector := range initialSetToDeploy {
		contractsToTransfer[selector] = []common.Address{
			state.Chains[selector].OnRamp.Address(),
			state.Chains[selector].OffRamp.Address(),
			state.Chains[selector].Router.Address(),
			state.Chains[selector].FeeQuoter.Address(),
			state.Chains[selector].RMNProxy.Address(),
			state.Chains[selector].NonceManager.Address(),
			state.Chains[selector].TokenAdminRegistry.Address(),
			state.Chains[selector].RMNRemote.Address(),
		}
	}
	contractsToTransfer[homeChainSelector] = append(
		contractsToTransfer[homeChainSelector],
		state.Chains[homeChainSelector].CCIPHome.Address(),
		state.Chains[homeChainSelector].CapabilityRegistry.Address(),
	)

	return commonchangeset.Apply(t, env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
			commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: contractsToTransfer,
				MCMSConfig: proposalutils.TimelockConfig{
					MinDelay: 0 * time.Second,
				},
			},
		),
	)
}
