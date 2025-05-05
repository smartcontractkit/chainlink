package v1_6_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	ccip "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/validate"
)

func TestJobSpecChangeset(t *testing.T) {
	t.Parallel()
	var err error
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNoJobsAndContracts())
	e := tenv.Env
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	e, err = commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
			v1_6.DeployHomeChainConfig{
				HomeChainSel:     tenv.HomeChainSel,
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Chains[tenv.HomeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					testhelpers.TestNodeOperator: nodes.NonBootstraps().PeerIDs(),
				},
			},
		),
	)
	require.NoError(t, err)
	output, err := v1_6.CCIPCapabilityJobspecChangeset(e, nil)
	require.NoError(t, err)
	require.NotEmpty(t, output.Jobs)
	nodeIDs := make(map[string]struct{})
	for _, job := range output.Jobs {
		require.NotEmpty(t, job.JobID)
		require.NotEmpty(t, job.Spec)
		require.NotEmpty(t, job.Node)
		nodeIDs[job.Node] = struct{}{}
		_, err = ccip.ValidatedCCIPSpec(job.Spec)
		require.NoError(t, err)
	}
	for _, node := range nodes {
		_, ok := nodeIDs[node.NodeID]
		require.True(t, ok)
	}
}

func TestJobSpecChangesetIdempotent(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(t)
	// we call the changeset again to ensure that it doesn't return any new job specs
	// as the job specs are already created in the first call
	output, err := v1_6.CCIPCapabilityJobspecChangeset(e.Env, nil)
	require.NoError(t, err)
	require.Empty(t, output.Jobs)
}
