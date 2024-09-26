package smoke

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/changeset"
	jobv1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/job/v1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test0002_InitialDeployOnLocal(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := ccipdeployment.Context(t)
	tenv := ccipdeployment.NewLocalDevEnvironmentWithCRAndFeeds(t, lggr)
	e := tenv.Env
	don := tenv.DON

	state, err := ccipdeployment.LoadOnchainState(tenv.Env, tenv.Ab)
	require.NoError(t, err)

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	tokenConfig := ccipdeployment.NewTokenConfig()
	tokenConfig.UpsertTokenInfo(ccipdeployment.LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: feeds[ccipdeployment.LinkSymbol].Address().String(),
			Decimals:          ccipdeployment.LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)
	// Apply migration
	output, err := changeset.Apply0002(tenv.Env, ccipdeployment.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		TokenConfig:    tokenConfig,
		ChainsToDeploy: tenv.Env.AllChainSelectors(),
		// Capreg/config already exist.
		CCIPOnChainState: state,
	})
	require.NoError(t, err)
	// Get new state after migration.
	state, err = ccipdeployment.LoadOnchainState(e, output.AddressBook)
	require.NoError(t, err)

	// Apply the jobs.
	nodeIdToJobIds := make(map[string][]string)
	for nodeID, jobs := range output.JobSpecs {
		nodeIdToJobIds[nodeID] = make([]string, 0, len(jobs))
		for _, job := range jobs {
			res, err := e.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
			require.NotNil(t, res.Proposal)
			nodeIdToJobIds[nodeID] = append(nodeIdToJobIds[nodeID], res.Proposal.JobId)
		}
	}

	// Accept all the jobs for this node.
	for _, n := range don.Nodes {
		jobsToAccept, exists := nodeIdToJobIds[n.NodeId]
		require.True(t, exists, "node %s has no jobs to accept", n.NodeId)
		for i, jobID := range jobsToAccept {
			require.NoError(t, n.AcceptJob(ctx, strconv.Itoa(i+1)), "node -%s failed to accept job %s", n.Name, jobID)
		}
	}
	t.Log("Jobs accepted")

	// Add all lanes
	require.NoError(t, ccipdeployment.AddLanesForAll(e, state))
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
			seqNum := ccipdeployment.SendRequest(t, e, state, src, dest, false)
			expectedSeqNum[dest] = seqNum
		}
	}

	// Wait for all commit reports to land.
	ccipdeployment.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.
	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, ccipdeployment.MockLinkPrice, timestampedPrice.Value)
	}
	// Wait for all exec reports to land
	ccipdeployment.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)
}
