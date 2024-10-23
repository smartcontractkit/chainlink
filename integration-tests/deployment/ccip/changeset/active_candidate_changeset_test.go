package changeset

import (
	"fmt"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"testing"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

	"github.com/stretchr/testify/require"

	ccdeploy "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	jobv1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/job/v1"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestActiveCandidate(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := ccdeploy.Context(t)
	tenv := ccdeploy.NewMemoryEnvironment(t, lggr, 3, 5)
	e := tenv.Env

	state, err := ccdeploy.LoadOnchainState(tenv.Env, tenv.Ab)
	require.NoError(t, err)
	require.NotNil(t, state.Chains[tenv.HomeChainSel].LinkToken)

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	tokenConfig := ccdeploy.NewTokenConfig()
	tokenConfig.UpsertTokenInfo(ccdeploy.LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: feeds[ccdeploy.LinkSymbol].Address().String(),
			Decimals:          ccdeploy.LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)

	output, err := InitialDeployChangeSet(tenv.Ab, tenv.Env, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:       tenv.HomeChainSel,
		FeedChainSel:       tenv.FeedChainSel,
		ChainsToDeploy:     tenv.Env.AllChainSelectors(),
		TokenConfig:        tokenConfig,
		MCMSConfig:         ccdeploy.NewTestMCMSConfig(t, e),
		FeeTokenContracts:  tenv.FeeTokenContracts,
		CapabilityRegistry: state.Chains[tenv.HomeChainSel].CapabilityRegistry.Address(),
		OCRSecrets:         deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	// Get new state after migration.
	state, err = ccdeploy.LoadOnchainState(e, tenv.Ab)
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
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	//startBlocks := make(map[uint64]*uint64)
	//// Send a message from each chain to every other chain.
	//expectedSeqNum := make(map[uint64]uint64)
	//for src := range e.Chains {
	//	for dest, destChain := range e.Chains {
	//		if src == dest {
	//			continue
	//		}
	//		latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
	//		require.NoError(t, err)
	//		block := latesthdr.Number.Uint64()
	//		startBlocks[dest] = &block
	//		seqNum := ccdeploy.SendRequest(t, e, state, src, dest, false)
	//		expectedSeqNum[dest] = seqNum
	//	}
	//}
	//
	//// Wait for all commit reports to land.
	//ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.
	//for dest := range e.Chains {
	//	linkAddress := state.Chains[dest].LinkToken.Address()
	//	feeQuoter := state.Chains[dest].FeeQuoter
	//	timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
	//	require.NoError(t, err)
	//	require.Equal(t, ccdeploy.MockLinkPrice, timestampedPrice.Value)
	//}

	// Wait for all exec reports to land
	//ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

	homeCS, destCS := tenv.HomeChainSel, tenv.FeedChainSel

	for _, source := range e.AllChainSelectors() {
		tx, err := state.Chains[source].OnRamp.TransferOwnership(e.Chains[source].DeployerKey, state.Chains[source].Timelock.Address())
		require.NoError(t, err)
		_, err = deployment.ConfirmIfNoError(e.Chains[source], tx, err)
		require.NoError(t, err)
		tx, err = state.Chains[source].FeeQuoter.TransferOwnership(e.Chains[source].DeployerKey, state.Chains[source].Timelock.Address())
		require.NoError(t, err)
		_, err = deployment.ConfirmIfNoError(e.Chains[source], tx, err)
		require.NoError(t, err)
	}
	// Transfer CR contract ownership
	tx, err := state.Chains[homeCS].CapabilityRegistry.TransferOwnership(e.Chains[homeCS].DeployerKey, state.Chains[homeCS].Timelock.Address())
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Chains[homeCS], tx, err)
	require.NoError(t, err)
	tx, err = state.Chains[homeCS].CCIPHome.TransferOwnership(e.Chains[homeCS].DeployerKey, state.Chains[homeCS].Timelock.Address())
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Chains[homeCS], tx, err)
	require.NoError(t, err)

	acceptOwnershipProposal, err := ccdeploy.GenerateAcceptOwnershipProposal(state, homeCS, e.AllChainSelectors())
	require.NoError(t, err)
	acceptOwnershipExec := ccdeploy.SignProposal(t, e, acceptOwnershipProposal)
	// Apply the accept ownership proposal to all the chains.
	for _, sel := range e.AllChainSelectors() {
		ccdeploy.ExecuteProposal(t, e, acceptOwnershipExec, state, sel)
	}

	err = ccdeploy.ConfirmRequestOnSourceAndDest(t, e, state, homeCS, destCS, 1)
	require.NoError(t, err)

	// [ACTIVE, CANDIDATE] setup by setting candidate through cap reg
	capReg, ccipHome := state.Chains[homeCS].CapabilityRegistry, state.Chains[homeCS].CCIPHome
	donID, err := ccdeploy.DonIDForChain(capReg, ccipHome, destCS)
	require.NoError(t, err)
	donInfo, err := state.Chains[homeCS].CapabilityRegistry.GetDON(nil, donID)
	require.NoError(t, err)
	require.Equal(t, 5, len(donInfo.NodeP2PIds))
	require.Equal(t, uint32(4), donInfo.ConfigCount)

	state, err = ccdeploy.LoadOnchainState(e, tenv.Ab)
	require.NoError(t, err)

	// delete a non-bootstrap node
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	newNodeIDs := []string{}
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
		e.Logger,
		deployment.XXXGenerateTestOCRSecrets(),
		state.Chains[destCS].OffRamp,
		e.Chains[destCS],
		destCS,
		tokenConfig.GetTokenInfo(e.Logger, state.Chains[destCS].LinkToken),
		nodes.NonBootstraps(),
		rmnHomeAddress,
	)
	require.NoError(t, err)

	fmt.Println("built ocr3 configurations")

	//var mcmsOps []mcms.Operation
	setCandidateMCMSOps, err := ccdeploy.SetCandidateOnExistingDon(
		ocr3ConfigMap[cctypes.PluginTypeCCIPCommit],
		state.Chains[homeCS].CapabilityRegistry,
		state.Chains[homeCS].CCIPHome,
		destCS,
		nodes.NonBootstraps(),
	)
	require.NoError(t, err)
	//mcmsOps = append(mcmsOps, setCandidateMCMSOps...)

	fmt.Println("creating proposal 1")
	setCandidateProposal, err := ccdeploy.BuildProposalFromBatches(state, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(homeCS),
		Batch:           setCandidateMCMSOps,
	}}, "set new candidates on commit plugin", "0s")
	fmt.Println("set new candidates on commit plugin")

	require.NoError(t, err)
	setCandidateSigned := ccdeploy.SignProposal(t, e, setCandidateProposal)
	fmt.Println("signed proposal 1")
	ccdeploy.ExecuteProposal(t, e, setCandidateSigned, state, homeCS)

	// create the op for the commit plugin as well
	setCandidateMCMSOps, err = ccdeploy.SetCandidateOnExistingDon(
		ocr3ConfigMap[cctypes.PluginTypeCCIPExec],
		state.Chains[homeCS].CapabilityRegistry,
		state.Chains[homeCS].CCIPHome,
		destCS,
		nodes.NonBootstraps(),
	)
	require.NoError(t, err)

	fmt.Println("creating proposal")
	setCandidateProposal, err = ccdeploy.BuildProposalFromBatches(state, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(homeCS),
		Batch:           setCandidateMCMSOps,
	}}, "set new candidates on commit and exec plugins", "0s")
	require.NoError(t, err)

	setCandidateSigned = ccdeploy.SignProposal(t, e, setCandidateProposal)
	fmt.Println("signed proposal 2")
	ccdeploy.ExecuteProposal(t, e, setCandidateSigned, state, homeCS)

	// check setup was successful by confirming number of nodes from cap reg
	donInfo, err = state.Chains[homeCS].CapabilityRegistry.GetDON(nil, donID)
	fmt.Printf("DonID is %d", donID)
	require.NoError(t, err)
	require.Equal(t, 4, len(donInfo.NodeP2PIds))
	require.Equal(t, uint32(6), donInfo.ConfigCount)
	// [ACTIVE, CANDIDATE] done setup

	// [ACTIVE, CANDIDATE] make sure we can still send successful transaction without updating job specs
	// this one fails
	fmt.Println("Sending request number 2")
	err = ccdeploy.ConfirmRequestOnSourceAndDest(t, e, state, homeCS, destCS, 2)
	require.NoError(t, err)
	// [ACTIVE, CANDIDATE] done send successful transaction on active

	// [NEW ACTIVE, NO CANDIDATE] promote to active
	// confirm by getting old candidate digest and making sure new active matches
	oldCandidateDigest, err := state.Chains[homeCS].CCIPHome.GetCandidateDigest(nil, donID, uint8(cctypes.PluginTypeCCIPExec))
	require.NoError(t, err)

	promoteOps, err := ccdeploy.PromoteAllCandidatesForChainOps(state.Chains[homeCS].CapabilityRegistry, state.Chains[homeCS].CCIPHome, destCS, nodes.NonBootstraps())
	require.NoError(t, err)

	promoteCandidateProposal, err := ccdeploy.BuildProposalFromBatches(state, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(homeCS),
		Batch:           promoteOps,
	}}, "promote candidates and revoke actives", "0s")
	require.NoError(t, err)
	promoteCandidateSigned := ccdeploy.SignProposal(t, e, promoteCandidateProposal)
	ccdeploy.ExecuteProposal(t, e, promoteCandidateSigned, state, homeCS)

	newActiveDigest, err := state.Chains[homeCS].CCIPHome.GetActiveDigest(nil, donID, uint8(cctypes.PluginTypeCCIPExec))
	require.NoError(t, err)
	require.Equal(t, oldCandidateDigest, newActiveDigest)

	newCandidateDigest, err := state.Chains[homeCS].CCIPHome.GetCandidateDigest(nil, donID, uint8(cctypes.PluginTypeCCIPCommit))
	require.NoError(t, err)
	require.Equal(t, newCandidateDigest, [32]byte{})
	// [NEW ACTIVE, NO CANDIDATE] done promoting

	// [NEW ACTIVE, NO CANDIDATE] send successful request on new active
	donInfo, err = state.Chains[homeCS].CapabilityRegistry.GetDON(nil, donID)
	require.NoError(t, err)
	require.Equal(t, uint32(8), donInfo.ConfigCount)

	err = ccdeploy.ConfirmRequestOnSourceAndDest(t, e, state, homeCS, destCS, 3)
	require.NoError(t, err)
	// [NEW ACTIVE, NO CANDIDATE] done sending successful request
}
