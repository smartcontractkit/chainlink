package smoke

import (
	"fmt"
	"testing"

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
	//chains := e.Chains

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

	// now accept the jobs
	for _, n := range nodes {
		jobsToAccept, exists := nodeIdToJobIds[n.NodeId]
		require.True(t, exists, "node %s has no jobs to accept", n.NodeId)
		for _, jobID := range jobsToAccept {
			require.NoError(t, n.AcceptJob(ctx, jobID), "node %s failed to accept job %s", n.NodeId, jobID)
		}
	}
	fmt.Println("Jobs accepted")
}
