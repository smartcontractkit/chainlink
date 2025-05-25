package ccip

import (
	"math/big"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/common"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

// Test_AddChainE2E tests the end-to-end functionality of adding a new chain to the existing
// chains using two consolidated changesets.
// 1. AddCandidatesForNewChainChangeset
// 2. PromoteNewChainForConfigChangeset
// This is a docker test
func Test_AddChainE2E(t *testing.T) {
	e, _, tEnv := testsetups.NewIntegrationEnvironment(
		t,
		// testhelpers.WithExtraConfigTomls([]string{"Test_AddChainE2E.toml"}),
		testhelpers.WithNumOfChains(3),
		testhelpers.WithPrerequisiteDeploymentOnly(nil),
	)

	toDeploy := e.Env.AllChainSelectorsExcluding([]uint64{e.HomeChainSel, e.FeedChainSel})
	initialSetToDeploy := e.Env.AllChainSelectorsExcluding(toDeploy)

	e = testhelpers.AddCCIPContractsToEnvironment(t, initialSetToDeploy, tEnv, false)
	// Need to update what the RMNProxy is pointing to, otherwise plugin will not work.
	var err error
	e.Env, err = commonchangeset.Apply(t, e.Env, e.TimelockContracts(t),
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

	// Build remote chain configurations
	remoteChains := make([]v1_6.ChainDefinition, len(initialSetToDeploy))
	for i, selector := range initialSetToDeploy {
		remoteChains[i] = v1_6.ChainDefinition{
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
	// Fetch the timelock and call proxy contracts for the initial set of chains
	timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts, len(initialSetToDeploy))
	for _, selector := range initialSetToDeploy {
		timelockContracts[selector] = &proposalutils.TimelockExecutionContracts{
			Timelock:  state.Chains[selector].Timelock,
			CallProxy: state.Chains[selector].CallProxy,
		}
	}

	// Transfer ownership of the home and feed chain contracts to the timelock
	e.Env, err = TransferOwnership(t, e.Env, timelockContracts, e.HomeChainSel, initialSetToDeploy, state)
	require.NoError(t, err, "must transfer ownership of home and feed chain contracts to the timelock")

	// setup the third chain with other two
	nodeInfo, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err, "must get node info")
	mcmsDeploymentCfg := proposalutils.SingleGroupTimelockConfigV2(t)
	chainToDeploy := toDeploy[0]
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
				GasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(testhelpers.DefaultGasPriceDeviationPPB)},
				DAGasPriceDeviationPPB:  cciptypes.BigInt{Int: big.NewInt(testhelpers.DefaultDAGasPriceDeviationPPB)},
				OptimisticConfirmations: globals.OptimisticConfirmations,
			},
		},
		CommitOCRParams: v1_6.DeriveOCRParamsForCommit(v1_6.Default, e.FeedChainSel, nil, nil),
		ExecOCRParams:   v1_6.DeriveOCRParamsForExec(v1_6.Default, nil, nil),
	}

	// need donIDClaimer contract to be deployed before we can deploy the new chain
	e.Env, err = commonchangeset.Apply(t, e.Env, nil,
		commonchangeset.Configure(
			v1_6.DeployDonIDClaimerChangeset,
			v1_6.DeployDonIDClaimerConfig{},
		))
	require.NoError(t, err, "must deploy donIDClaimer contract")

	// Add candidate for new chain using AddCandidatesForNewChainChangeset
	e.Env, err = commonchangeset.Apply(t, e.Env, timelockContracts,
		commonchangeset.Configure(
			v1_6.AddCandidatesForNewChainChangeset,
			v1_6.AddCandidatesForNewChainConfig{
				HomeChainSelector:    e.HomeChainSel,
				FeedChainSelector:    e.FeedChainSel,
				NewChain:             newChainDefinition,
				RemoteChains:         remoteChains,
				MCMSDeploymentConfig: &mcmsDeploymentCfg,
				MCMSConfig: &proposalutils.TimelockConfig{
					MinDelay:   0 * time.Second,
					MCMSAction: mcmstypes.TimelockActionSchedule,
				},
				DonIDOffSet: nil,
			},
		),
	)
	require.NoError(t, err, "must apply AddCandidatesForNewChainChangeset")

	// Apply PromoteNewChainForConfigChangeset
	e.Env, err = commonchangeset.Apply(t, e.Env, timelockContracts,
		commonchangeset.Configure(
			v1_6.PromoteNewChainForConfigChangeset,
			v1_6.PromoteNewChainForConfig{
				HomeChainSelector: e.HomeChainSel,
				NewChain:          newChainDefinition,
				RemoteChains:      remoteChains,
				TestRouter:        pointer.ToBool(false),
				MCMSConfig: &proposalutils.TimelockConfig{
					MinDelay:   0 * time.Second,
					MCMSAction: mcmstypes.TimelockActionSchedule,
				},
			},
		),
	)
	require.NoError(t, err, "must apply PromoteNewChainForConfigChangeset")
	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)
	time.Sleep(120 * time.Second)
	SendMessages(t,
		e.Env,
		[]testhelpers.SourceDestPair{
			{
				SourceChainSelector: e.HomeChainSel,
				DestChainSelector:   chainToDeploy,
			},
			// {
			// 	SourceChainSelector: chainToDeploy,
			// 	DestChainSelector:   e.HomeChainSel,
			// },
			// {
			// 	SourceChainSelector: e.HomeChainSel,
			// 	DestChainSelector:   e.FeedChainSel,
			// },
			// {
			// 	SourceChainSelector: e.FeedChainSel,
			// 	DestChainSelector:   e.HomeChainSel,
			// },
			// {
			// 	SourceChainSelector: e.FeedChainSel,
			// 	DestChainSelector:   chainToDeploy,
			// },
			// {
			// 	SourceChainSelector: chainToDeploy,
			// 	DestChainSelector:   e.FeedChainSel,
			// },
		},
		state,
		false,
	)
}

func SendMessages(
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
	for _, pair := range sourceDestPairs {
		latesthdr, err := env.Chains[pair.DestChainSelector].Client.HeaderByNumber(testcontext.Get(t), nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		// time.Sleep(10 * time.Second)
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
	testhelpers.ConfirmCommitForAllWithExpectedSeqNums(t, env, state, expectedSeqNum, startBlocks)
	testhelpers.ConfirmExecWithSeqNrsForAll(t, env, state, expectedSeqNumExec, startBlocks)
}

func TransferOwnership(
	t *testing.T,
	env cldf.Environment,
	timelockContracts map[uint64]*proposalutils.TimelockExecutionContracts,
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
	)
	contractsToTransfer[homeChainSelector] = append(
		contractsToTransfer[homeChainSelector],
		state.Chains[homeChainSelector].CapabilityRegistry.Address(),
	)

	return commonchangeset.Apply(t, env, timelockContracts,
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
