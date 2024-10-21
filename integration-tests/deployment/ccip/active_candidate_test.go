package ccipdeployment

import (
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
)

func Test_ActiveCandidateMigration(t *testing.T) {
	// [SETUP]
	// 2 chains with a lane connecting them.
	// We set up 5 nodes initially. Our candidate configuration will have 4 nodes
	e := NewMemoryEnvironmentWithJobs(t, logger.TestLogger(t), 2, 5)
	state, err := LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)

	// We deploy to all chain selectors
	initialDeploy := e.Env.AllChainSelectors()

	feeds := state.Chains[e.FeedChainSel].USDFeeds
	tokenConfig := NewTokenConfig()
	tokenConfig.UpsertTokenInfo(LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: feeds[LinkSymbol].Address().String(),
			Decimals:          LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)
	err = DeployCCIPContracts(e.Env, e.Ab, DeployCCIPContractConfig{
		HomeChainSel:       e.HomeChainSel,
		FeedChainSel:       e.FeedChainSel,
		ChainsToDeploy:     initialDeploy,
		TokenConfig:        tokenConfig,
		MCMSConfig:         NewTestMCMSConfig(t, e.Env),
		FeeTokenContracts:  e.FeeTokenContracts,
		CapabilityRegistry: state.Chains[e.HomeChainSel].CapabilityRegistry.Address(),
		OCRSecrets:         deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	state, err = LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)

	// Connect all the existing lanes.
	for _, source := range initialDeploy {
		for _, dest := range initialDeploy {
			if source != dest {
				require.NoError(t, AddLane(e.Env, state, source, dest))
			}
		}
	}

	homeCS, destCS := e.HomeChainSel, e.FeedChainSel
	rmnHomeAddress := state.Chains[homeCS].RMNHome.Address()

	// Transfer onramp/fq ownership to timelock.
	// Enable the new dest on the test router.
	for _, source := range initialDeploy {
		tx, err := state.Chains[source].OnRamp.TransferOwnership(e.Env.Chains[source].DeployerKey, state.Chains[source].Timelock.Address())
		require.NoError(t, err)
		_, err = deployment.ConfirmIfNoError(e.Env.Chains[source], tx, err)
		require.NoError(t, err)
		tx, err = state.Chains[source].FeeQuoter.TransferOwnership(e.Env.Chains[source].DeployerKey, state.Chains[source].Timelock.Address())
		require.NoError(t, err)
		_, err = deployment.ConfirmIfNoError(e.Env.Chains[source], tx, err)
		require.NoError(t, err)
	}
	// Transfer CR contract ownership
	tx, err := state.Chains[homeCS].CapabilityRegistry.TransferOwnership(e.Env.Chains[homeCS].DeployerKey, state.Chains[homeCS].Timelock.Address())
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[homeCS], tx, err)
	require.NoError(t, err)
	tx, err = state.Chains[homeCS].CCIPHome.TransferOwnership(e.Env.Chains[homeCS].DeployerKey, state.Chains[homeCS].Timelock.Address())
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(e.Env.Chains[homeCS], tx, err)
	require.NoError(t, err)

	acceptOwnershipProposal, err := GenerateAcceptOwnershipProposal(state, homeCS, initialDeploy)
	require.NoError(t, err)
	acceptOwnershipExec := SignProposal(t, e.Env, acceptOwnershipProposal)
	// Apply the accept ownership proposal to all the chains.
	for _, sel := range initialDeploy {
		ExecuteProposal(t, e.Env, acceptOwnershipExec, state, sel)
	}
	for _, chain := range initialDeploy {
		owner, err2 := state.Chains[chain].OnRamp.Owner(nil)
		require.NoError(t, err2)
		require.Equal(t, state.Chains[chain].Timelock.Address(), owner)
	}
	cfgOwner, err := state.Chains[homeCS].CCIPHome.Owner(nil)
	require.NoError(t, err)
	crOwner, err := state.Chains[homeCS].CapabilityRegistry.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, state.Chains[homeCS].Timelock.Address(), cfgOwner)
	require.Equal(t, state.Chains[homeCS].Timelock.Address(), crOwner)

	// check setup was successful
	donID, err := DonIDForChain(state.Chains[homeCS].CapabilityRegistry, state.Chains[homeCS].CCIPHome, destCS)
	require.NoError(t, err)
	donInfo, err := state.Chains[homeCS].CapabilityRegistry.GetDON(nil, donID)
	require.NoError(t, err)
	require.Equal(t, 5, len(donInfo.NodeP2PIds))
	require.Equal(t, uint32(4), donInfo.ConfigCount)
	// [SETUP] done

	// [ACTIVE ONLY, NO CANDIDATE] send successful request on active
	//latesthdr, err := e.Env.Chains[destCS].Client.HeaderByNumber(testcontext.Get(t), nil)
	//require.NoError(t, err)
	//startBlock := latesthdr.Number.Uint64()
	seqNum := SendRequest(t, e.Env, state, homeCS, destCS, false)
	require.Equal(t, uint64(1), seqNum)

	//ConfirmCommitForAllWithExpectedSeqNums(t, e.Env, state,
	//	map[uint64]uint64{destCS: seqNum, homeCS: 0},
	//	map[uint64]*uint64{destCS: &startBlock, homeCS: &startBlock})
	//
	//require.NoError(t, ConfirmExecWithSeqNr(t, e.Env.Chains[homeCS], e.Env.Chains[destCS], state.Chains[destCS].OffRamp, &startBlock, seqNum))
	// [ACTIVE ONLY, NO CANDIDATE] done

	// [ACTIVE, CANDIDATE] setup by setting candidate through cap reg
	state, err = LoadOnchainState(e.Env, e.Ab)
	require.NoError(t, err)

	// delete the last node from the list of nodes.
	// bootstrap node is first and will be ignored so we delete from the end
	e.Env.NodeIDs = e.Env.NodeIDs[:len(e.Env.NodeIDs)-1]
	nodes, err := deployment.NodeInfo(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)

	// this will construct ocr3 configurations for the
	// commit and exec plugin we will be using
	ocr3ConfigMap, err := BuildOCR3ConfigForCCIPHome(
		e.Env.Logger,
		deployment.XXXGenerateTestOCRSecrets(),
		state.Chains[destCS].OffRamp,
		e.Env.Chains[destCS],
		e.FeedChainSel,
		tokenConfig.GetTokenInfo(e.Env.Logger, state.Chains[destCS].LinkToken),
		nodes.NonBootstraps(),
		rmnHomeAddress,
	)

	require.NoError(t, err)

	var mcmsOps []mcms.Operation

	setCandidateMCMSOps, err := SetCandidateOnExistingDon(
		ocr3ConfigMap[cctypes.PluginTypeCCIPExec],
		state.Chains[homeCS].CapabilityRegistry,
		state.Chains[homeCS].CCIPHome,
		destCS,
		nodes.NonBootstraps(),
	)
	require.NoError(t, err)
	mcmsOps = append(mcmsOps, setCandidateMCMSOps...)

	// create the op for the commit plugin as well
	setCandidateMCMSOps, err = SetCandidateOnExistingDon(
		ocr3ConfigMap[cctypes.PluginTypeCCIPCommit],
		state.Chains[homeCS].CapabilityRegistry,
		state.Chains[homeCS].CCIPHome,
		destCS,
		nodes.NonBootstraps(),
	)
	require.NoError(t, err)
	mcmsOps = append(mcmsOps, setCandidateMCMSOps...)

	tl, mcmMds, err := BuildProposalMetadata(state, []uint64{homeCS})
	setCandidateProposal, err := timelock.NewMCMSWithTimelockProposal(
		"1",
		2004259681, // TODO
		[]mcms.Signature{},
		false,
		mcmMds,
		tl,
		"set new candidates on commit and exec plugins", // TODO
		[]timelock.BatchChainOperation{{
			ChainIdentifier: mcms.ChainIdentifier(homeCS),
			Batch:           mcmsOps,
		}},
		timelock.Schedule, "0s")

	require.NoError(t, err)
	setCandidateSigned := SignProposal(t, e.Env, setCandidateProposal)
	ExecuteProposal(t, e.Env, setCandidateSigned, state, e.HomeChainSel)

	// check setup was successful by confirming number of nodes from cap reg
	donInfo, err = state.Chains[homeCS].CapabilityRegistry.GetDON(nil, donID)
	require.NoError(t, err)
	require.Equal(t, 4, len(donInfo.NodeP2PIds))
	require.Equal(t, uint32(6), donInfo.ConfigCount)
	// [ACTIVE, CANDIDATE] done setup

	// [ACTIVE, CANDIDATE] make sure we can still send successful transaction
	//seqNum = SendRequest(t, e.Env, state, homeCS, destCS, false)
	//require.Equal(t, uint64(2), seqNum)
	//latesthdr, err = e.Env.Chains[destCS].Client.HeaderByNumber(testcontext.Get(t), nil)
	//require.NoError(t, err)
	//startBlock = latesthdr.Number.Uint64()
	seqNum = SendRequest(t, e.Env, state, homeCS, destCS, false)
	require.Equal(t, uint64(2), seqNum)
	// Wait for all commit reports to land.
	//ConfirmCommitForAllWithExpectedSeqNums(t, e.Env, state,
	//	map[uint64]uint64{destCS: seqNum},
	//	map[uint64]*uint64{destCS: &startBlock})
	//require.NoError(t, ConfirmExecWithSeqNr(t, e.Env.Chains[homeCS], e.Env.Chains[destCS], state.Chains[destCS].OffRamp, &startBlock, seqNum))
	// [ACTIVE, CANDIDATE] done send successful transaction on active

	// [NEW ACTIVE, NO CANDIDATE] promote to active
	// confirm by getting old candidate digest and making sure new active matches
	oldCandidateDigest, err := state.Chains[homeCS].CCIPHome.GetCandidateDigest(nil, donID, uint8(cctypes.PluginTypeCCIPExec))

	mcmsOps, err = PromoteAllCandidatesForChainOps(state.Chains[homeCS].CapabilityRegistry, state.Chains[homeCS].CCIPHome, destCS, nodes.NonBootstraps())
	require.NoError(t, err)

	tl, mcmMds, err = BuildProposalMetadata(state, []uint64{homeCS})
	promoteCandidateProposal, err := timelock.NewMCMSWithTimelockProposal(
		"1",
		2004259681, // TODO
		[]mcms.Signature{},
		false,
		mcmMds,
		tl,
		"promote candidates and revoke actives", // TODO
		[]timelock.BatchChainOperation{{
			ChainIdentifier: mcms.ChainIdentifier(homeCS),
			Batch:           mcmsOps,
		}},
		timelock.Schedule, "0s")

	promoteCandidateSigned := SignProposal(t, e.Env, promoteCandidateProposal)
	ExecuteProposal(t, e.Env, promoteCandidateSigned, state, e.HomeChainSel)

	newActiveDigest, err := state.Chains[homeCS].CCIPHome.GetActiveDigest(nil, donID, uint8(cctypes.PluginTypeCCIPExec))
	require.NoError(t, err)
	require.Equal(t, oldCandidateDigest, newActiveDigest)

	newCandidateDigest, err := state.Chains[homeCS].CCIPHome.GetCandidateDigest(nil, donID, uint8(cctypes.PluginTypeCCIPCommit))
	require.NoError(t, err)
	require.Equal(t, newCandidateDigest, emptyUint32)
	// [NEW ACTIVE, NO CANDIDATE] done promoting

	// [NEW ACTIVE, NO CANDIDATE] send successful request on new active
	seqNum = SendRequest(t, e.Env, state, homeCS, destCS, false)
	require.Equal(t, uint64(3), seqNum)
	donInfo, err = state.Chains[homeCS].CapabilityRegistry.GetDON(nil, donID)
	require.NoError(t, err)
	require.Equal(t, uint32(8), donInfo.ConfigCount)
	// [NEW ACTIVE, NO CANDIDATE] done sending successful request
}
