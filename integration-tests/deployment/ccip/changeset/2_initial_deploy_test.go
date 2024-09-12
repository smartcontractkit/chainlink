package changeset

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	jobv1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/job/v1"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test0002_InitialDeploy(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := ccipdeployment.Context(t)
	tenv := ccipdeployment.NewDeployedTestEnvironment(t, lggr)
	e := tenv.Env
	nodes := tenv.Nodes
	chains := e.Chains

	state, err := ccipdeployment.LoadOnchainState(tenv.Env, tenv.Ab)
	require.NoError(t, err)

	// Apply migration
	output, err := Apply0002(tenv.Env, ccipdeployment.DeployCCIPContractConfig{
		HomeChainSel: tenv.HomeChainSel,
		// Capreg/config already exist.
		CCIPOnChainState: state,
	})
	require.NoError(t, err)
	// Get new state after migration.
	state, err = ccipdeployment.LoadOnchainState(e, output.AddressBook)
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	require.NoError(t, ReplayAllLogs(nodes, chains))

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
	// Wait for plugins to register filters?
	// TODO: Investigate how to avoid.
	time.Sleep(30 * time.Second)

	// Ensure job related logs are up to date.
	require.NoError(t, ReplayAllLogs(nodes, chains))

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
			// record the block number for each chain
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

	// TODO: Apply the proposal.
}

func ReplayAllLogs(nodes map[string]memory.Node, chains map[uint64]deployment.Chain) error {
	for _, node := range nodes {
		for sel := range chains {
			if err := node.ReplayLogs(map[uint64]uint64{sel: 1}); err != nil {
				return err
			}
		}
	}
	return nil
}
