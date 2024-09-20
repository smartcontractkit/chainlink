package smoke

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/changeset"
	jobv1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/job/v1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test0002_InitialDeployOnLocal(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := ccipdeployment.Context(t)
	tenv := ccipdeployment.NewDeployedLocalDevEnvironment(t, lggr)
	e := tenv.Env
	nodes := tenv.Nodes

	state, err := ccipdeployment.LoadOnchainState(tenv.Env, tenv.Ab)
	require.NoError(t, err)

	// Apply migration
	output, err := changeset.Apply0002(tenv.Env, ccipdeployment.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
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
	for _, n := range nodes {
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

	// Wait for all exec reports to land
	ccipdeployment.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)
}
