package changeset_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_ActiveCandidate(t *testing.T) {
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
	tenv.Env, err = commonchangeset.ApplyChangesets(t, tenv.Env, tenv.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.UpdateOnRampsDestsChangeset),
			Config: changeset.UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]changeset.OnRampDestinationUpdate{
					source: {
						dest: {
							IsEnabled:        true,
							AllowListEnabled: false,
						},
					},
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.UpdateFeeQuoterPricesChangeset),
			Config: changeset.UpdateFeeQuoterPricesConfig{
				PricesByChain: map[uint64]changeset.FeeQuoterPriceUpdatePerSource{
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
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.UpdateFeeQuoterDestsChangeset),
			Config: changeset.UpdateFeeQuoterDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]fee_quoter.FeeQuoterDestChainConfig{
					source: {
						dest: changeset.DefaultFeeQuoterDestChainConfig(),
					},
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.UpdateOffRampSourcesChangeset),
			Config: changeset.UpdateOffRampSourcesConfig{
				UpdatesByChain: map[uint64]map[uint64]changeset.OffRampSourceUpdate{
					dest: {
						source: {
							IsEnabled: true,
						},
					},
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.UpdateRouterRampsChangeset),
			Config: changeset.UpdateRouterRampsConfig{
				UpdatesByChain: map[uint64]changeset.RouterUpdates{
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
		},
	})
	require.NoError(t, err)

	// check that source router has dest enabled
	onRamp, err := sourceState.Router.GetOnRamp(&bind.CallOpts{
		Context: testcontext.Get(t),
	}, dest)
	require.NoError(t, err)
	require.NotEqual(t, common.HexToAddress("0x0"), onRamp, "expected onRamp to be set")

	// Transfer ownership so that we can set new candidate configs
	// and set new config digest on the offramp.
	_, err = commonchangeset.ApplyChangesets(t, tenv.Env, tenv.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.TransferToMCMSWithTimelock),
			Config:    testhelpers.GenTestTransferOwnershipConfig(tenv, allChains, state),
		},
	})
	require.NoError(t, err)
	testhelpers.AssertTimelockOwnership(t, tenv, allChains, state)

	sendMsg := func() {
		latesthdr, err := tenv.Env.Chains[dest].Client.HeaderByNumber(testcontext.Get(t), nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
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
	_, err = commonchangeset.ApplyChangesets(t, tenv.Env, tenv.TimelockContracts(t), []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.SetCandidateChangeset),
			Config: changeset.SetCandidateChangesetConfig{
				SetCandidateConfigBase: changeset.SetCandidateConfigBase{
					HomeChainSelector: tenv.HomeChainSel,
					FeedChainSelector: tenv.FeedChainSel,
					MCMS: &changeset.MCMSConfig{
						MinDelay: 0,
					},
				},
				PluginInfo: []changeset.SetCandidatePluginInfo{
					{
						// NOTE: this is technically not a new chain, but needed for validation.
						OCRConfigPerRemoteChainSelector: map[uint64]changeset.CCIPOCRParams{
							dest: changeset.DeriveCCIPOCRParams(
								changeset.WithDefaultCommitOffChainConfig(tenv.FeedChainSel, tokenConfig.GetTokenInfo(logger.TestLogger(t), state.Chains[dest].LinkToken, state.Chains[dest].Weth9)),
							),
						},
						PluginType: types.PluginTypeCCIPCommit,
					},
					{
						// NOTE: this is technically not a new chain, but needed for validation.
						OCRConfigPerRemoteChainSelector: map[uint64]changeset.CCIPOCRParams{
							dest: changeset.DeriveCCIPOCRParams(
								changeset.WithDefaultExecuteOffChainConfig(nil),
							),
						},
						PluginType: types.PluginTypeCCIPExec,
					},
				},
			},
		},
	})
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
