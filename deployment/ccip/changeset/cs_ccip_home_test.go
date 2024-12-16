package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/stretchr/testify/require"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func TestActiveCandidate(t *testing.T) {
	t.Skipf("to be enabled after latest cl-ccip is compatible")
	t.Parallel()
	tenv := NewMemoryEnvironment(t,
		WithChains(3),
		WithNodes(5))
	e := tenv.Env
	state, err := LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	allChains := maps.Keys(e.Chains)

	// Add all lanes
	require.NoError(t, AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[SourceDestPair]uint64)
	expectedSeqNumExec := make(map[SourceDestPair][]uint64)
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
			expectedSeqNumExec[SourceDestPair{
				SourceChainSelector: src,
				DestChainSelector:   dest,
			}] = []uint64{msgSentEvent.SequenceNumber}
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
	ConfirmExecWithSeqNrsForAll(t, e, state, expectedSeqNumExec, startBlocks)

	// compose the transfer ownership and accept ownership changesets
	timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts)
	for _, chain := range allChains {
		timelockContracts[chain] = &proposalutils.TimelockExecutionContracts{
			Timelock:  state.Chains[chain].Timelock,
			CallProxy: state.Chains[chain].CallProxy,
		}
	}

	_, err = commonchangeset.ApplyChangesets(t, e, timelockContracts, []commonchangeset.ChangesetApplication{
		// note this doesn't have proposals.
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.TransferToMCMSWithTimelock),
			Config:    genTestTransferOwnershipConfig(tenv, allChains, state),
		},
	})
	require.NoError(t, err)
	// Apply the accept ownership proposal to all the chains.

	err = ConfirmRequestOnSourceAndDest(t, e, state, tenv.HomeChainSel, tenv.FeedChainSel, 2)
	require.NoError(t, err)

	// [ACTIVE, CANDIDATE] setup by setting candidate through cap reg
	capReg, ccipHome := state.Chains[tenv.HomeChainSel].CapabilityRegistry, state.Chains[tenv.HomeChainSel].CCIPHome
	donID, err := internal.DonIDForChain(capReg, ccipHome, tenv.FeedChainSel)
	require.NoError(t, err)
	require.NotEqual(t, uint32(0), donID)
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
	ccipOCRParams := DefaultOCRParams(
		tenv.FeedChainSel,
		tokenConfig.GetTokenInfo(e.Logger, state.Chains[tenv.FeedChainSel].LinkToken, state.Chains[tenv.FeedChainSel].Weth9),
		nil,
	)
	ocr3ConfigMap, err := internal.BuildOCR3ConfigForCCIPHome(
		e.OCRSecrets,
		state.Chains[tenv.FeedChainSel].OffRamp,
		e.Chains[tenv.FeedChainSel],
		nodes.NonBootstraps(),
		rmnHomeAddress,
		ccipOCRParams.OCRParameters,
		ccipOCRParams.CommitOffChainConfig,
		ccipOCRParams.ExecuteOffChainConfig,
	)
	require.NoError(t, err)

	var (
		timelocksPerChain = map[uint64]common.Address{
			tenv.HomeChainSel: state.Chains[tenv.HomeChainSel].Timelock.Address(),
		}
		proposerMCMSes = map[uint64]*gethwrappers.ManyChainMultiSig{
			tenv.HomeChainSel: state.Chains[tenv.HomeChainSel].ProposerMcm,
		}
	)
	setCommitCandidateOp, err := setCandidateOnExistingDon(
		e.Logger,
		deployment.SimTransactOpts(),
		tenv.Env.Chains[tenv.HomeChainSel],
		ocr3ConfigMap[cctypes.PluginTypeCCIPCommit],
		state.Chains[tenv.HomeChainSel].CapabilityRegistry,
		state.Chains[tenv.HomeChainSel].CCIPHome,
		tenv.FeedChainSel,
		nodes.NonBootstraps(),
		true,
	)
	require.NoError(t, err)
	setCommitCandidateProposal, err := proposalutils.BuildProposalFromBatches(timelocksPerChain, proposerMCMSes, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(tenv.HomeChainSel),
		Batch:           setCommitCandidateOp,
	}}, "set new candidates on commit plugin", 0)
	require.NoError(t, err)
	setCommitCandidateSigned := proposalutils.SignProposal(t, e, setCommitCandidateProposal)
	proposalutils.ExecuteProposal(t, e, setCommitCandidateSigned, &proposalutils.TimelockExecutionContracts{
		Timelock:  state.Chains[tenv.HomeChainSel].Timelock,
		CallProxy: state.Chains[tenv.HomeChainSel].CallProxy,
	}, tenv.HomeChainSel)

	// create the op for the commit plugin as well
	setExecCandidateOp, err := setCandidateOnExistingDon(
		e.Logger,
		deployment.SimTransactOpts(),
		tenv.Env.Chains[tenv.HomeChainSel],
		ocr3ConfigMap[cctypes.PluginTypeCCIPExec],
		state.Chains[tenv.HomeChainSel].CapabilityRegistry,
		state.Chains[tenv.HomeChainSel].CCIPHome,
		tenv.FeedChainSel,
		nodes.NonBootstraps(),
		true,
	)
	require.NoError(t, err)

	setExecCandidateProposal, err := proposalutils.BuildProposalFromBatches(timelocksPerChain, proposerMCMSes, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(tenv.HomeChainSel),
		Batch:           setExecCandidateOp,
	}}, "set new candidates on commit and exec plugins", 0)
	require.NoError(t, err)
	setExecCandidateSigned := proposalutils.SignProposal(t, e, setExecCandidateProposal)
	proposalutils.ExecuteProposal(t, e, setExecCandidateSigned, &proposalutils.TimelockExecutionContracts{
		Timelock:  state.Chains[tenv.HomeChainSel].Timelock,
		CallProxy: state.Chains[tenv.HomeChainSel].CallProxy,
	}, tenv.HomeChainSel)

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

	promoteOps, err := promoteAllCandidatesForChainOps(
		tenv.Env.Chains[tenv.HomeChainSel],
		deployment.SimTransactOpts(),
		state.Chains[tenv.HomeChainSel].CapabilityRegistry,
		state.Chains[tenv.HomeChainSel].CCIPHome,
		tenv.FeedChainSel,
		nodes.NonBootstraps(),
		true)
	require.NoError(t, err)
	promoteProposal, err := proposalutils.BuildProposalFromBatches(timelocksPerChain, proposerMCMSes, []timelock.BatchChainOperation{{
		ChainIdentifier: mcms.ChainIdentifier(tenv.HomeChainSel),
		Batch:           promoteOps,
	}}, "promote candidates and revoke actives", 0)
	require.NoError(t, err)
	promoteSigned := proposalutils.SignProposal(t, e, promoteProposal)
	proposalutils.ExecuteProposal(t, e, promoteSigned, &proposalutils.TimelockExecutionContracts{
		Timelock:  state.Chains[tenv.HomeChainSel].Timelock,
		CallProxy: state.Chains[tenv.HomeChainSel].CallProxy,
	}, tenv.HomeChainSel)
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

func Test_PromoteCandidate(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mcmsEnabled bool
	}{
		{
			name:        "MCMS enabled",
			mcmsEnabled: true,
		},
		{
			name:        "MCMS disabled",
			mcmsEnabled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testcontext.Get(t)
			tenv := NewMemoryEnvironment(t,
				WithChains(2),
				WithNodes(4))
			state, err := LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			// Deploy to all chains.
			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, source, dest)
			}

			var (
				capReg   = state.Chains[tenv.HomeChainSel].CapabilityRegistry
				ccipHome = state.Chains[tenv.HomeChainSel].CCIPHome
			)
			donID, err := internal.DonIDForChain(capReg, ccipHome, dest)
			require.NoError(t, err)
			require.NotEqual(t, uint32(0), donID)
			candidateDigestCommitBefore, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPCommit))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, candidateDigestCommitBefore)
			candidateDigestExecBefore, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPExec))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, candidateDigestExecBefore)

			var mcmsConfig *MCMSConfig
			if tc.mcmsEnabled {
				mcmsConfig = &MCMSConfig{
					MinDelay: 0,
				}
			}
			_, err = commonchangeset.ApplyChangesets(t, tenv.Env, map[uint64]*proposalutils.TimelockExecutionContracts{
				tenv.HomeChainSel: {
					Timelock:  state.Chains[tenv.HomeChainSel].Timelock,
					CallProxy: state.Chains[tenv.HomeChainSel].CallProxy,
				},
			}, []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(PromoteAllCandidatesChangeset),
					Config: PromoteAllCandidatesChangesetConfig{
						HomeChainSelector: tenv.HomeChainSel,
						DONChainSelector:  dest,
						MCMS:              mcmsConfig,
					},
				},
			})
			require.NoError(t, err)

			// after promoting the zero digest, active digest should also be zero
			activeDigestCommit, err := ccipHome.GetActiveDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPCommit))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, activeDigestCommit)

			activeDigestExec, err := ccipHome.GetActiveDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPExec))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, activeDigestExec)
		})
	}
}

func Test_SetCandidate(t *testing.T) {
	for _, tc := range []struct {
		name        string
		mcmsEnabled bool
	}{
		{
			name:        "MCMS enabled",
			mcmsEnabled: true,
		},
		{
			name:        "MCMS disabled",
			mcmsEnabled: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testcontext.Get(t)
			tenv := NewMemoryEnvironment(t,
				WithChains(2),
				WithNodes(4))
			state, err := LoadOnchainState(tenv.Env)
			require.NoError(t, err)

			// Deploy to all chains.
			allChains := maps.Keys(tenv.Env.Chains)
			source := allChains[0]
			dest := allChains[1]

			if tc.mcmsEnabled {
				// Transfer ownership to timelock so that we can promote the zero digest later down the line.
				transferToTimelock(t, tenv, state, source, dest)
			}

			var (
				capReg   = state.Chains[tenv.HomeChainSel].CapabilityRegistry
				ccipHome = state.Chains[tenv.HomeChainSel].CCIPHome
			)
			donID, err := internal.DonIDForChain(capReg, ccipHome, dest)
			require.NoError(t, err)
			require.NotEqual(t, uint32(0), donID)
			candidateDigestCommitBefore, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPCommit))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, candidateDigestCommitBefore)
			candidateDigestExecBefore, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPExec))
			require.NoError(t, err)
			require.Equal(t, [32]byte{}, candidateDigestExecBefore)

			var mcmsConfig *MCMSConfig
			if tc.mcmsEnabled {
				mcmsConfig = &MCMSConfig{
					MinDelay: 0,
				}
			}
			tokenConfig := NewTestTokenConfig(state.Chains[tenv.FeedChainSel].USDFeeds)
			_, err = commonchangeset.ApplyChangesets(t, tenv.Env, map[uint64]*proposalutils.TimelockExecutionContracts{
				tenv.HomeChainSel: {
					Timelock:  state.Chains[tenv.HomeChainSel].Timelock,
					CallProxy: state.Chains[tenv.HomeChainSel].CallProxy,
				},
			}, []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(SetCandidateChangeset),
					Config: SetCandidateChangesetConfig{
						HomeChainSelector: tenv.HomeChainSel,
						FeedChainSelector: tenv.FeedChainSel,
						DONChainSelector:  dest,
						PluginType:        types.PluginTypeCCIPCommit,
						CCIPOCRParams: DefaultOCRParams(
							tenv.FeedChainSel,
							tokenConfig.GetTokenInfo(logger.TestLogger(t), state.Chains[dest].LinkToken, state.Chains[dest].Weth9),
							nil,
						),
						MCMS: mcmsConfig,
					},
				},
				{
					Changeset: commonchangeset.WrapChangeSet(SetCandidateChangeset),
					Config: SetCandidateChangesetConfig{
						HomeChainSelector: tenv.HomeChainSel,
						FeedChainSelector: tenv.FeedChainSel,
						DONChainSelector:  dest,
						PluginType:        types.PluginTypeCCIPExec,
						CCIPOCRParams: DefaultOCRParams(
							tenv.FeedChainSel,
							tokenConfig.GetTokenInfo(logger.TestLogger(t), state.Chains[dest].LinkToken, state.Chains[dest].Weth9),
							nil,
						),
						MCMS: mcmsConfig,
					},
				},
			})
			require.NoError(t, err)

			// after setting a new candidate on both plugins, the candidate config digest
			// should be nonzero.
			candidateDigestCommitAfter, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPCommit))
			require.NoError(t, err)
			require.NotEqual(t, [32]byte{}, candidateDigestCommitAfter)
			require.NotEqual(t, candidateDigestCommitBefore, candidateDigestCommitAfter)

			candidateDigestExecAfter, err := ccipHome.GetCandidateDigest(&bind.CallOpts{
				Context: ctx,
			}, donID, uint8(types.PluginTypeCCIPExec))
			require.NoError(t, err)
			require.NotEqual(t, [32]byte{}, candidateDigestExecAfter)
			require.NotEqual(t, candidateDigestExecBefore, candidateDigestExecAfter)
		})
	}
}

func transferToTimelock(
	t *testing.T,
	tenv DeployedEnv,
	state CCIPOnChainState,
	source,
	dest uint64) {
	// Transfer ownership to timelock so that we can promote the zero digest later down the line.
	_, err := commonchangeset.ApplyChangesets(t, tenv.Env, map[uint64]*proposalutils.TimelockExecutionContracts{
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
			Config:    genTestTransferOwnershipConfig(tenv, []uint64{source, dest}, state),
		},
	})
	require.NoError(t, err)
	assertTimelockOwnership(t, tenv, []uint64{source, dest}, state)
}
