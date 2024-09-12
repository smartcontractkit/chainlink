package smoke

import (
	"strconv"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/stretchr/testify/require"

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
		HomeChainSel: tenv.HomeChainSel,
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

	// Wait for plugins to register filters?
	// TODO: Investigate how to avoid.
	time.Sleep(30 * time.Second)

	// Send a request from every router
	// Add all lanes
	require.NoError(t, ccipdeployment.AddLanesForAll(e, state))

	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	for src, srcChain := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			num, err := destChain.LatestBlockNum(ctx)
			require.NoError(t, err)
			startBlocks[dest] = &num
			t.Logf("Sending CCIP request from chain selector %d to chain selector %d",
				src, dest)
			_, err = ccipdeployment.SendMessage(src, dest, e.Chains[src].DeployerKey, srcChain.Confirm, state)
			require.NoError(t, err)
		}
	}

	// Wait for all commit reports to land.
	ccipdeployment.WaitForCommitForAllWithInterval(t, e, state, ccipocr3.SeqNumRange{1, 1}, startBlocks)

	// Wait for all exec reports to land
	ccipdeployment.WaitForExecWithSeqNrForAll(t, e, state, 1, startBlocks)
}
