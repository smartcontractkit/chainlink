package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/stretchr/testify/require"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestActiveCandidate(t *testing.T) {
	t.Skipf("to be enabled after latest cl-ccip is compatible")

	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	tenv := ccdeploy.NewMemoryEnvironment(t, lggr, 3, 5, ccdeploy.MockLinkPrice, ccdeploy.MockWethPrice)
	e := tenv.Env

	state, err := ccdeploy.LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	require.NotNil(t, state.Chains[tenv.HomeChainSel].LinkToken)

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	tokenConfig := ccdeploy.NewTestTokenConfig(feeds)

	output, err := InitialDeploy(tenv.Env, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		ChainsToDeploy: tenv.Env.AllChainSelectors(),
		TokenConfig:    tokenConfig,
		MCMSConfig:     ccdeploy.NewTestMCMSConfig(t, e),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	// Get new state after migration.
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(output.AddressBook))
	state, err = ccdeploy.LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	homeCS, destCS := tenv.HomeChainSel, tenv.FeedChainSel

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
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[uint64]uint64)
	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block
			msgSentEvent := ccdeploy.TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
				Data:         []byte("hello world"),
				TokenAmounts: nil,
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			})
			expectedSeqNum[dest] = msgSentEvent.SequenceNumber
		}
	}

	// Wait for all commit reports to land.
	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	//After commit is reported on all chains, token prices should be updated in FeeQuoter.
	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, ccdeploy.MockLinkPrice, timestampedPrice.Value)
	}

	//Wait for all exec reports to land
	ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

	// transfer ownership
	ccdeploy.TransferAllOwnership(t, state, homeCS, e)
	acceptOwnershipProposal, err := ccdeploy.GenerateAcceptOwnershipProposal(state, homeCS, e.AllChainSelectors())
	require.NoError(t, err)
	acceptOwnershipExec := ccdeploy.SignProposal(t, e, acceptOwnershipProposal)
	for _, sel := range e.AllChainSelectors() {
		ccdeploy.ExecuteProposal(t, e, acceptOwnershipExec, state, sel)
	}
	// Apply the accept ownership proposal to all the chains.

	err = ccdeploy.ConfirmRequestOnSourceAndDest(t, e, state, homeCS, destCS, 2)
	require.NoError(t, err)

	// [ACTIVE, CANDIDATE] setup by setting candidate through cap reg
	capReg, ccipHome := state.Chains[homeCS].CapabilityRegistry, state.Chains[homeCS].CCIPHome
	donID, err := ccdeploy.DonIDForChain(capReg, ccipHome, destCS)
	require.NoError(t, err)
	donInfo, err := state.Chains[homeCS].CapabilityRegistry.GetDON(nil, donID)
	require.NoError(t, err)
	require.Equal(t, 5, len(donInfo.NodeP2PIds))
	require.Equal(t, uint32(4), donInfo.ConfigCount)

	state, err = ccdeploy.LoadOnchainState(e)
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
	rmnHomeAddress := state.Chains[homeCS].RMNHome.Address()
	ocr3ConfigMap, err := ccdeploy.BuildOCR3ConfigForCCIPHome(
		deployment.XXXGenerateTestOCRSecrets(),
		state.Chains[destCS].OffRamp,
		e.Chains[destCS],
		destCS,
		tokenConfig.GetTokenInfo(e.Logger, state.Chains[destCS].LinkToken, state.Chains[destCS].Weth9),
		nodes.NonBootstraps(),
		rmnHomeAddress,
		nil,
	)
	require.NoError(t, err)

	setCommitCandidateOp, err := ccdeploy.SetCandidateOnExistingDon(
		ocr3ConfigMap[cctypes.PluginTypeCCIPCommit],
		state.Chains[homeCS].CapabilityRegistry,
		state.Chains[homeCS].CCIPHome,
		destCS,
		nodes.NonBootstraps(),
	)
	require.NoError(t, err)
	setCommitCandidateProposal, err := ccdeploy.BuildProposalFromBatches(state, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(homeCS),
		Batch:           setCommitCandidateOp,
	}}, "set new candidates on commit plugin", 0)
	require.NoError(t, err)
	setCommitCandidateSigned := ccdeploy.SignProposal(t, e, setCommitCandidateProposal)
	ccdeploy.ExecuteProposal(t, e, setCommitCandidateSigned, state, homeCS)

	// create the op for the commit plugin as well
	setExecCandidateOp, err := ccdeploy.SetCandidateOnExistingDon(
		ocr3ConfigMap[cctypes.PluginTypeCCIPExec],
		state.Chains[homeCS].CapabilityRegistry,
		state.Chains[homeCS].CCIPHome,
		destCS,
		nodes.NonBootstraps(),
	)
	require.NoError(t, err)

	setExecCandidateProposal, err := ccdeploy.BuildProposalFromBatches(state, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(homeCS),
		Batch:           setExecCandidateOp,
	}}, "set new candidates on commit and exec plugins", 0)
	require.NoError(t, err)
	setExecCandidateSigned := ccdeploy.SignProposal(t, e, setExecCandidateProposal)
	ccdeploy.ExecuteProposal(t, e, setExecCandidateSigned, state, homeCS)

	// check setup was successful by confirming number of nodes from cap reg
	donInfo, err = state.Chains[homeCS].CapabilityRegistry.GetDON(nil, donID)
	require.NoError(t, err)
	require.Equal(t, 4, len(donInfo.NodeP2PIds))
	require.Equal(t, uint32(6), donInfo.ConfigCount)
	// [ACTIVE, CANDIDATE] done setup

	// [ACTIVE, CANDIDATE] make sure we can still send successful transaction without updating job specs
	err = ccdeploy.ConfirmRequestOnSourceAndDest(t, e, state, homeCS, destCS, 3)
	require.NoError(t, err)
	// [ACTIVE, CANDIDATE] done send successful transaction on active

	// [NEW ACTIVE, NO CANDIDATE] promote to active
	// confirm by getting old candidate digest and making sure new active matches
	oldCandidateDigest, err := state.Chains[homeCS].CCIPHome.GetCandidateDigest(nil, donID, uint8(cctypes.PluginTypeCCIPExec))
	require.NoError(t, err)

	promoteOps, err := ccdeploy.PromoteAllCandidatesForChainOps(state.Chains[homeCS].CapabilityRegistry, state.Chains[homeCS].CCIPHome, destCS, nodes.NonBootstraps())
	require.NoError(t, err)
	promoteProposal, err := ccdeploy.BuildProposalFromBatches(state, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(homeCS),
		Batch:           promoteOps,
	}}, "promote candidates and revoke actives", 0)
	require.NoError(t, err)
	promoteSigned := ccdeploy.SignProposal(t, e, promoteProposal)
	ccdeploy.ExecuteProposal(t, e, promoteSigned, state, homeCS)
	// [NEW ACTIVE, NO CANDIDATE] done promoting

	// [NEW ACTIVE, NO CANDIDATE] check onchain state
	newActiveDigest, err := state.Chains[homeCS].CCIPHome.GetActiveDigest(nil, donID, uint8(cctypes.PluginTypeCCIPExec))
	require.NoError(t, err)
	require.Equal(t, oldCandidateDigest, newActiveDigest)

	newCandidateDigest, err := state.Chains[homeCS].CCIPHome.GetCandidateDigest(nil, donID, uint8(cctypes.PluginTypeCCIPCommit))
	require.NoError(t, err)
	require.Equal(t, newCandidateDigest, [32]byte{})
	// [NEW ACTIVE, NO CANDIDATE] done checking on chain state

	// [NEW ACTIVE, NO CANDIDATE] send successful request on new active
	donInfo, err = state.Chains[homeCS].CapabilityRegistry.GetDON(nil, donID)
	require.NoError(t, err)
	require.Equal(t, uint32(8), donInfo.ConfigCount)

	err = ccdeploy.ConfirmRequestOnSourceAndDest(t, e, state, homeCS, destCS, 4)
	require.NoError(t, err)
	// [NEW ACTIVE, NO CANDIDATE] done sending successful request
}
