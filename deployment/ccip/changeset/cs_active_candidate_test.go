package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"golang.org/x/exp/maps"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_ActiveCandidate(t *testing.T) {
	// Setup an environment with 2 chains, a source and a dest.
	// We want to have the active instance execute a few messages
	// and then setup a candidate instance. The candidate instance
	// should not be able to transmit anything until we make it active.
	lggr := logger.TestLogger(t)
	tenv := NewMemoryEnvironmentWithJobsAndContracts(t, lggr, memory.MemoryEnvironmentConfig{
		Chains:             2,
		NumOfUsersPerChain: 1,
		Nodes:              4,
		Bootstraps:         1,
	}, nil)
	state, err := LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	// Deploy to all chains.
	allChains := maps.Keys(tenv.Env.Chains)
	source := allChains[0]
	dest := allChains[1]
	t.Logf("source: %d, dest: %d, home chain: %d, feed chain: %d", source, dest, tenv.HomeChainSel, tenv.FeedChainSel)
	newAddresses := deployment.NewMemoryAddressBook()
	err = deployPrerequisiteChainContracts(tenv.Env, newAddresses, allChains)
	require.NoError(t, err)
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(newAddresses))

	// Connect source to dest
	require.NoError(t, AddLaneWithDefaultPricesAndFeeQuoterConfig(tenv.Env, state, source, dest, false))

	// Transfer ownership so that we can set new candidate configs
	// and set new config digest on the offramp.
	_, err = commonchangeset.ApplyChangesets(t, tenv.Env, map[uint64]*commonchangeset.TimelockExecutionContracts{
		source: {
			Timelock:  state.Chains[source].Timelock,
			CallProxy: state.Chains[source].CallProxy,
		},
		dest: {
			Timelock:  state.Chains[dest].Timelock,
			CallProxy: state.Chains[dest].CallProxy,
		},
		tenv.HomeChainSel: {
			Timelock:  state.Chains[tenv.HomeChainSel].Timelock,
			CallProxy: state.Chains[tenv.HomeChainSel].CallProxy,
		},
	}, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.TransferToMCMSWithTimelock),
			Config:    genTestTransferOwnershipConfig(tenv, allChains, state),
		},
	})
	require.NoError(t, err)
	assertTimelockOwnership(t, tenv, allChains, state)

	// send a message from source to dest and ensure that it gets executed
	latesthdr, err := tenv.Env.Chains[dest].Client.HeaderByNumber(testcontext.Get(t), nil)
	require.NoError(t, err)
	block := latesthdr.Number.Uint64()
	msgSentEvent := TestSendRequest(t, tenv.Env, state, source, dest, false, router.ClientEVM2AnyMessage{
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
		expectedSeqNum = map[SourceDestPair]uint64{
			{
				SourceChainSelector: source,
				DestChainSelector:   dest,
			}: msgSentEvent.SequenceNumber,
		}
		expectedSeqNumExec = map[SourceDestPair][]uint64{
			{
				SourceChainSelector: source,
				DestChainSelector:   dest,
			}: {msgSentEvent.SequenceNumber},
		}
	)

	// Confirm execution of the message
	ConfirmCommitForAllWithExpectedSeqNums(t, tenv.Env, state, expectedSeqNum, startBlocks)
	ConfirmExecWithSeqNrsForAll(t, tenv.Env, state, expectedSeqNumExec, startBlocks)

	nodes, err := deployment.NodeInfo(tenv.Env.NodeIDs, tenv.Env.Offchain)
	require.NoError(t, err)

	var nodeIDs []string
	for _, node := range nodes {
		nodeIDs = append(nodeIDs, node.NodeID)
	}

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
	tokenConfig := NewTestTokenConfig(state.Chains[tenv.FeedChainSel].USDFeeds)
	_, err = commonchangeset.ApplyChangesets(t, tenv.Env, map[uint64]*commonchangeset.TimelockExecutionContracts{
		tenv.HomeChainSel: {
			Timelock:  state.Chains[tenv.HomeChainSel].Timelock,
			CallProxy: state.Chains[tenv.HomeChainSel].CallProxy,
		},
	}, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(SetCandidatePluginChangeset),
			Config: AddDonAndSetCandidateChangesetConfig{
				HomeChainSelector: tenv.HomeChainSel,
				FeedChainSelector: tenv.FeedChainSel,
				// NOTE: this is technically not a new chain, but needed for validation.
				NewChainSelector: dest,
				PluginType:       types.PluginTypeCCIPCommit,
				NodeIDs:          nodeIDs,
				CCIPOCRParams: DefaultOCRParams(
					tenv.FeedChainSel,
					tokenConfig.GetTokenInfo(logger.TestLogger(t), state.Chains[dest].LinkToken, state.Chains[dest].Weth9),
					nil,
				),
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(SetCandidatePluginChangeset),
			Config: AddDonAndSetCandidateChangesetConfig{
				HomeChainSelector: tenv.HomeChainSel,
				FeedChainSelector: tenv.FeedChainSel,
				// NOTE: this is technically not a new chain, but needed for validation.
				NewChainSelector: dest,
				PluginType:       types.PluginTypeCCIPExec,
				NodeIDs:          nodeIDs,
				CCIPOCRParams: DefaultOCRParams(
					tenv.FeedChainSel,
					tokenConfig.GetTokenInfo(logger.TestLogger(t), state.Chains[dest].LinkToken, state.Chains[dest].Weth9),
					nil,
				),
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
}
