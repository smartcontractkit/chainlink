package environment

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

func testTopologyWithWorkflowDONs(t *testing.T, wfDONs ...*cre.DonMetadata) *cre.Topology {
	t.Helper()

	dons := append([]*cre.DonMetadata{bootstrapDONMetadata()}, wfDONs...)
	return &cre.Topology{DonsMetadata: cre.NewUncheckedDonsMetadata(dons)}
}

func TestWorkflowContainerPatternForDON(t *testing.T) {
	t.Parallel()

	require.Equal(t, "feeds-zone-a-node", workflowContainerPatternForDON(&cre.DonMetadata{Name: "feeds-zone-a"}))
	require.Equal(t, "shard0-node", workflowContainerPatternForDON(&cre.DonMetadata{Name: "shard0"}))
}

func TestResolveWorkflowDONMetadata_singleDON(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, testTopologyWithWorkflowDONs(t,
		&cre.DonMetadata{Name: "workflow", ID: 1, DonFamily: envconfig.DefaultDONFamily, Flags: []string{cre.WorkflowDON}},
	))

	don, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{})
	require.NoError(t, err)
	require.Equal(t, "workflow", don.Name)
}

func TestResolveWorkflowDONMetadata_singleDONByFamily(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, testTopologyWithWorkflowDONs(t,
		&cre.DonMetadata{Name: "workflow", ID: 1, DonFamily: envconfig.DefaultDONFamily, Flags: []string{cre.WorkflowDON}},
	))

	don, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{DonFamily: envconfig.DefaultDONFamily})
	require.NoError(t, err)
	require.Equal(t, "workflow", don.Name)
}

func TestResolveWorkflowDONMetadata_multiDONExplicitName(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, multiWorkflowDONTestTopology(t))

	don, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{ExplicitName: "feeds-zone-b"})
	require.NoError(t, err)
	require.Equal(t, "feeds-zone-b", don.Name)
}

func TestResolveWorkflowDONMetadata_multiDONByFamily(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, multiWorkflowDONTestTopology(t))

	don, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{DonFamily: "feeds-zone-a"})
	require.NoError(t, err)
	require.Equal(t, "feeds-zone-a", don.Name)
}

func TestResolveWorkflowDONMetadata_multiDONRequiresSelector(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, multiWorkflowDONTestTopology(t))

	_, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "--don-family")
}

func TestResolveWorkflowDONMetadata_unknownExplicitName(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, multiWorkflowDONTestTopology(t))

	_, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{ExplicitName: "unknown"})
	require.Error(t, err)
	require.Contains(t, err.Error(), `workflow DON "unknown" not found`)
}

func TestResolveWorkflowDONMetadata_shardedByFamilyAndShardIndex(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, shardedWorkflowDONTestTopology(t))

	shard1 := uint(1)
	don, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{
		DonFamily:  envconfig.DefaultDONFamily,
		ShardIndex: &shard1,
	})
	require.NoError(t, err)
	require.Equal(t, "shard1", don.Name)
}

func TestResolveWorkflowDONMetadata_shardedRequiresShardIndex(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, shardedWorkflowDONTestTopology(t))

	_, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{DonFamily: envconfig.DefaultDONFamily})
	require.Error(t, err)
	require.Contains(t, err.Error(), "--shard-index")
}

func testStateResolver(t *testing.T, topology *cre.Topology) *LocalCREStateResolver {
	t.Helper()

	return &LocalCREStateResolver{topology: topology}
}

func multiWorkflowDONTestTopology(t *testing.T) *cre.Topology {
	t.Helper()

	return testTopologyWithWorkflowDONs(t,
		&cre.DonMetadata{Name: "feeds-zone-a", ID: 1, DonFamily: "feeds-zone-a", Flags: []string{cre.WorkflowDON}},
		&cre.DonMetadata{Name: "feeds-zone-b", ID: 2, DonFamily: "feeds-zone-b", Flags: []string{cre.WorkflowDON}},
	)
}

func shardedWorkflowDONTestTopology(t *testing.T) *cre.Topology {
	t.Helper()

	return testTopologyWithWorkflowDONs(t,
		&cre.DonMetadata{Name: "shard0", ID: 1, DonFamily: envconfig.DefaultDONFamily, ShardIndex: 0, Flags: []string{cre.WorkflowDON, cre.ShardDON}},
		&cre.DonMetadata{Name: "shard1", ID: 2, DonFamily: envconfig.DefaultDONFamily, ShardIndex: 1, Flags: []string{cre.WorkflowDON, cre.ShardDON}},
	)
}

func bootstrapDONMetadata() *cre.DonMetadata {
	return &cre.DonMetadata{
		Name:          "bootstrap",
		NodesMetadata: []*cre.NodeMetadata{{Roles: []string{cre.BootstrapNode}}},
	}
}
