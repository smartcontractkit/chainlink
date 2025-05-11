package v1_6_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_ActiveCandidate(t *testing.T) {
	t.Parallel()
	// Setup an environment with 2 chains, a source and a dest.
	// We want to have the active instance execute a few messages
	// and then setup a candidate instance. The candidate instance
	// should not be able to transmit anything until we make it active.
	tenv, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithNumOfNodes(4))
	state, err := changeset.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	// Deploy to all chains.
	allChains := maps.Keys(tenv.Env.Chains)
	source := allChains[0]
	dest := allChains[1]

	// Connect source to dest
	sourceState := state.Chains[source]
	tenv.Env, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateOnRampsDestsChangeset),
			v1_6.UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]v1_6.OnRampDestinationUpdate{
					source: {
						dest: {
							IsEnabled:        true,
							AllowListEnabled: false,
						},
					},
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterPricesChangeset),
			v1_6.UpdateFeeQuoterPricesConfig{
				PricesByChain: map[uint64]v1_6.FeeQuoterPriceUpdatePerSource{
					source: {
						TokenPrices: map[common.Address]*big.Int{
							sourceState.LinkToken.Address(): testhelpers.DefaultLinkPrice,
							sourceState.Weth9.Address():     testhelpers.DefaultWethPrice,
						},
						GasPrices: map[uint64]*big.Int{
							dest: testhelpers.DefaultGasPrice,
						},
					},
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterDestsChangeset),
			v1_6.UpdateFeeQuoterDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
					source: {
						dest: v1_6.DefaultFeeQuoterDestChainConfig(true),
					},
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateOffRampSourcesChangeset),
			v1_6.UpdateOffRampSourcesConfig{
				UpdatesByChain: map[uint64]map[uint64]v1_6.OffRampSourceUpdate{
					dest: {
						source: {
							IsEnabled:                 true,
							IsRMNVerificationDisabled: true,
						},
					},
				},
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateRouterRampsChangeset),
			v1_6.UpdateRouterRampsConfig{
				UpdatesByChain: map[uint64]v1_6.RouterUpdates{
					// onRamp update on source chain
					source: {
						OnRampUpdates: map[uint64]bool{
							dest: true,
						},
					},
					// offramp update on dest chain
					dest: {
						OffRampUpdates: map[uint64]bool{
							source: true,
						},
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// check that source router has dest enabled
	onRamp, err := sourceState.Router.GetOnRamp(&bind.CallOpts{
		Context: testcontext.Get(t),
	}, dest)
	require.NoError(t, err)
	require.NotEqual(t, common.HexToAddress("0x0"), onRamp, "expected onRamp to be set")

	// Transfer ownership so that we can set new candidate configs
	// and set new config digest on the offramp.
	_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelock),
			testhelpers.GenTestTransferOwnershipConfig(tenv, allChains, state),
		),
	)
	require.NoError(t, err)
	testhelpers.AssertTimelockOwnership(t, tenv, allChains, state)

	sendMsg := func() {
		block, err := testhelpers.LatestBlock(testcontext.Get(t), tenv.Env, dest)
		require.NoError(t, err)
		msgSentEvent := testhelpers.TestSendRequest(t, tenv.Env, state, source, dest, false, router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
			Data:         []byte("hello world"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		})

		var (
			startBlocks = map[uint64]*uint64{
				dest: &block,
			}
			expectedSeqNum = map[testhelpers.SourceDestPair]uint64{
				{
					SourceChainSelector: source,
					DestChainSelector:   dest,
				}: msgSentEvent.SequenceNumber,
			}
			expectedSeqNumExec = map[testhelpers.SourceDestPair][]uint64{
				{
					SourceChainSelector: source,
					DestChainSelector:   dest,
				}: {msgSentEvent.SequenceNumber},
			}
		)

		// Confirm execution of the message
		testhelpers.ConfirmCommitForAllWithExpectedSeqNums(t, tenv.Env, state, expectedSeqNum, startBlocks)
		testhelpers.ConfirmExecWithSeqNrsForAll(t, tenv.Env, state, expectedSeqNumExec, startBlocks)
	}

	// send a message from source to dest and ensure that it gets executed
	sendMsg()

	var (
		capReg   = state.Chains[tenv.HomeChainSel].CapabilityRegistry
		ccipHome = state.Chains[tenv.HomeChainSel].CCIPHome
	)
	donID, err := internal.DonIDForChain(capReg, ccipHome, dest)
	require.NoError(t, err)
	candidateDigestCommitBefore, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
		Context: testcontext.Get(t),
	}, donID, uint8(types.PluginTypeCCIPCommit))
	require.NoError(t, err)
	require.Equal(t, [32]byte{}, candidateDigestCommitBefore)
	candidateDigestExecBefore, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
		Context: testcontext.Get(t),
	}, donID, uint8(types.PluginTypeCCIPExec))
	require.NoError(t, err)
	require.Equal(t, [32]byte{}, candidateDigestExecBefore)

	// Now we can add a candidate config, send another request, and observe behavior.
	// The candidate config should not be able to execute messages.
	tokenConfig := changeset.NewTestTokenConfig(state.Chains[tenv.FeedChainSel].USDFeeds)
	_, err = commonchangeset.Apply(t, tenv.Env, tenv.TimelockContracts(t),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.SetCandidateChangeset),
			v1_6.SetCandidateChangesetConfig{
				SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
					HomeChainSelector: tenv.HomeChainSel,
					FeedChainSelector: tenv.FeedChainSel,
					MCMS: &proposalutils.TimelockConfig{
						MinDelay: 0,
					},
				},
				PluginInfo: []v1_6.SetCandidatePluginInfo{
					{
						// NOTE: this is technically not a new chain, but needed for validation.
						OCRConfigPerRemoteChainSelector: map[uint64]v1_6.CCIPOCRParams{
							dest: v1_6.DeriveOCRParamsForCommit(v1_6.SimulationTest, tenv.FeedChainSel, tokenConfig.GetTokenInfo(logger.TestLogger(t),
								state.Chains[dest].LinkToken.Address(),
								state.Chains[dest].Weth9.Address()), nil),
						},
						PluginType: types.PluginTypeCCIPCommit,
					},
					{
						// NOTE: this is technically not a new chain, but needed for validation.
						OCRConfigPerRemoteChainSelector: map[uint64]v1_6.CCIPOCRParams{
							dest: v1_6.DeriveOCRParamsForExec(v1_6.SimulationTest, nil, nil),
						},
						PluginType: types.PluginTypeCCIPExec,
					},
				},
			},
		),
	)
	require.NoError(t, err)

	// check that CCIPHome state is updated with the new candidate configs
	// for the dest chain DON.
	candidateDigestCommit, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
		Context: testcontext.Get(t),
	}, donID, uint8(types.PluginTypeCCIPCommit))
	require.NoError(t, err)
	require.NotEqual(t, candidateDigestCommit, candidateDigestCommitBefore)
	candidateDigestExec, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
		Context: testcontext.Get(t),
	}, donID, uint8(types.PluginTypeCCIPExec))
	require.NoError(t, err)
	require.NotEqual(t, candidateDigestExec, candidateDigestExecBefore)

	// send a message from source to dest and ensure that it gets executed after the candidate config is set
	sendMsg()
}
