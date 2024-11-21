package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/stretchr/testify/require"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestActiveCandidate(t *testing.T) {
	t.Skipf("to be enabled after latest cl-ccip is compatible")

	lggr := logger.TestLogger(t)
	tenv := NewMemoryEnvironmentWithJobsAndContracts(t, lggr, 3, 5)
	e := tenv.Env
	state, err := LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	// Add all lanes
	require.NoError(t, AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[SourceDestPair]uint64)
	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block
			msgSentEvent := TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
				Data:         []byte("hello world"),
				TokenAmounts: nil,
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			})
			expectedSeqNum[SourceDestPair{
				SourceChainSelector: src,
				DestChainSelector:   dest,
			}] = msgSentEvent.SequenceNumber
		}
	}

	// Wait for all commit reports to land.
	ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	//After commit is reported on all chains, token prices should be updated in FeeQuoter.
	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, MockLinkPrice, timestampedPrice.Value)
	}

	//Wait for all exec reports to land
	ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

	// transfer ownership
	TransferAllOwnership(t, state, tenv.HomeChainSel, e)
	acceptOwnershipProposal, err := GenerateAcceptOwnershipProposal(state, tenv.HomeChainSel, e.AllChainSelectors())
	require.NoError(t, err)
	acceptOwnershipExec := commonchangeset.SignProposal(t, e, acceptOwnershipProposal)
	for _, sel := range e.AllChainSelectors() {
		commonchangeset.ExecuteProposal(t, e, acceptOwnershipExec, state.Chains[sel].Timelock, sel)
	}
	// Apply the accept ownership proposal to all the chains.

	err = ConfirmRequestOnSourceAndDest(t, e, state, tenv.HomeChainSel, tenv.FeedChainSel, 2)
	require.NoError(t, err)

	// [ACTIVE, CANDIDATE] setup by setting candidate through cap reg
	capReg, ccipHome := state.Chains[tenv.HomeChainSel].CapabilityRegistry, state.Chains[tenv.HomeChainSel].CCIPHome
	donID, err := internal.DonIDForChain(capReg, ccipHome, tenv.FeedChainSel)
	require.NoError(t, err)
	donInfo, err := state.Chains[tenv.HomeChainSel].CapabilityRegistry.GetDON(nil, donID)
	require.NoError(t, err)
	require.Equal(t, 5, len(donInfo.NodeP2PIds))
	require.Equal(t, uint32(4), donInfo.ConfigCount)

	state, err = LoadOnchainState(e)
	require.NoError(t, err)

	// delete a non-bootstrap node
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	var newNodeIDs []string
	// make sure we delete a node that is NOT bootstrap.
	// we will remove bootstrap later by calling nodes.NonBootstrap()
	if nodes[0].IsBootstrap {
		newNodeIDs = e.NodeIDs[:len(e.NodeIDs)-1]
	} else {
		newNodeIDs = e.NodeIDs[1:]
	}
	nodes, err = deployment.NodeInfo(newNodeIDs, e.Offchain)
	require.NoError(t, err)

	// this will construct ocr3 configurations for the
	// commit and exec plugin we will be using
	rmnHomeAddress := state.Chains[tenv.HomeChainSel].RMNHome.Address()
	tokenConfig := NewTestTokenConfig(state.Chains[tenv.FeedChainSel].USDFeeds)
	ocr3ConfigMap, err := internal.BuildOCR3ConfigForCCIPHome(
		deployment.XXXGenerateTestOCRSecrets(),
		state.Chains[tenv.FeedChainSel].OffRamp,
		e.Chains[tenv.FeedChainSel],
		tenv.FeedChainSel,
		tokenConfig.GetTokenInfo(e.Logger, state.Chains[tenv.FeedChainSel].LinkToken, state.Chains[tenv.FeedChainSel].Weth9),
		nodes.NonBootstraps(),
		rmnHomeAddress,
		nil,
	)
	require.NoError(t, err)

	setCommitCandidateOp, err := SetCandidateOnExistingDon(
		ocr3ConfigMap[cctypes.PluginTypeCCIPCommit],
		state.Chains[tenv.HomeChainSel].CapabilityRegistry,
		state.Chains[tenv.HomeChainSel].CCIPHome,
		tenv.FeedChainSel,
		nodes.NonBootstraps(),
	)
	require.NoError(t, err)
	setCommitCandidateProposal, err := BuildProposalFromBatches(state, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(tenv.HomeChainSel),
		Batch:           setCommitCandidateOp,
	}}, "set new candidates on commit plugin", 0)
	require.NoError(t, err)
	setCommitCandidateSigned := commonchangeset.SignProposal(t, e, setCommitCandidateProposal)
	commonchangeset.ExecuteProposal(t, e, setCommitCandidateSigned, state.Chains[tenv.HomeChainSel].Timelock, tenv.HomeChainSel)

	// create the op for the commit plugin as well
	setExecCandidateOp, err := SetCandidateOnExistingDon(
		ocr3ConfigMap[cctypes.PluginTypeCCIPExec],
		state.Chains[tenv.HomeChainSel].CapabilityRegistry,
		state.Chains[tenv.HomeChainSel].CCIPHome,
		tenv.FeedChainSel,
		nodes.NonBootstraps(),
	)
	require.NoError(t, err)

	setExecCandidateProposal, err := BuildProposalFromBatches(state, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(tenv.HomeChainSel),
		Batch:           setExecCandidateOp,
	}}, "set new candidates on commit and exec plugins", 0)
	require.NoError(t, err)
	setExecCandidateSigned := commonchangeset.SignProposal(t, e, setExecCandidateProposal)
	commonchangeset.ExecuteProposal(t, e, setExecCandidateSigned, state.Chains[tenv.HomeChainSel].Timelock, tenv.HomeChainSel)

	// check setup was successful by confirming number of nodes from cap reg
	donInfo, err = state.Chains[tenv.HomeChainSel].CapabilityRegistry.GetDON(nil, donID)
	require.NoError(t, err)
	require.Equal(t, 4, len(donInfo.NodeP2PIds))
	require.Equal(t, uint32(6), donInfo.ConfigCount)
	// [ACTIVE, CANDIDATE] done setup

	// [ACTIVE, CANDIDATE] make sure we can still send successful transaction without updating job specs
	err = ConfirmRequestOnSourceAndDest(t, e, state, tenv.HomeChainSel, tenv.FeedChainSel, 3)
	require.NoError(t, err)
	// [ACTIVE, CANDIDATE] done send successful transaction on active

	// [NEW ACTIVE, NO CANDIDATE] promote to active
	// confirm by getting old candidate digest and making sure new active matches
	oldCandidateDigest, err := state.Chains[tenv.HomeChainSel].CCIPHome.GetCandidateDigest(nil, donID, uint8(cctypes.PluginTypeCCIPExec))
	require.NoError(t, err)

	promoteOps, err := PromoteAllCandidatesForChainOps(state.Chains[tenv.HomeChainSel].CapabilityRegistry, state.Chains[tenv.HomeChainSel].CCIPHome, tenv.FeedChainSel, nodes.NonBootstraps())
	require.NoError(t, err)
	promoteProposal, err := BuildProposalFromBatches(state, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(tenv.HomeChainSel),
		Batch:           promoteOps,
	}}, "promote candidates and revoke actives", 0)
	require.NoError(t, err)
	promoteSigned := commonchangeset.SignProposal(t, e, promoteProposal)
	commonchangeset.ExecuteProposal(t, e, promoteSigned, state.Chains[tenv.HomeChainSel].Timelock, tenv.HomeChainSel)
	// [NEW ACTIVE, NO CANDIDATE] done promoting

	// [NEW ACTIVE, NO CANDIDATE] check onchain state
	newActiveDigest, err := state.Chains[tenv.HomeChainSel].CCIPHome.GetActiveDigest(nil, donID, uint8(cctypes.PluginTypeCCIPExec))
	require.NoError(t, err)
	require.Equal(t, oldCandidateDigest, newActiveDigest)

	newCandidateDigest, err := state.Chains[tenv.HomeChainSel].CCIPHome.GetCandidateDigest(nil, donID, uint8(cctypes.PluginTypeCCIPCommit))
	require.NoError(t, err)
	require.Equal(t, newCandidateDigest, [32]byte{})
	// [NEW ACTIVE, NO CANDIDATE] done checking on chain state

	// [NEW ACTIVE, NO CANDIDATE] send successful request on new active
	donInfo, err = state.Chains[tenv.HomeChainSel].CapabilityRegistry.GetDON(nil, donID)
	require.NoError(t, err)
	require.Equal(t, uint32(8), donInfo.ConfigCount)

	err = ConfirmRequestOnSourceAndDest(t, e, state, tenv.HomeChainSel, tenv.FeedChainSel, 4)
	require.NoError(t, err)
	// [NEW ACTIVE, NO CANDIDATE] done sending successful request
}
