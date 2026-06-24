// Tests for workflow DON selection during env workflow deploy.
//
// Covers ResolveWorkflowDONMetadata and resolveWorkflowDONByFamily: how --workflow-don-name,
// --don-family, and --shard-index pick the target workflow DON from saved local CRE topology.
// Also tests the default docker cp container pattern derived from nodesets.name.
package environment

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func testTopologyWithWorkflowDONs(t *testing.T, wfDONs ...*cre.DonMetadata) *cre.Topology {
	t.Helper()

	dons := append([]*cre.DonMetadata{bootstrapDONMetadata()}, wfDONs...)
	dm, err := cre.NewDonsMetadata(dons, infra.Provider{Type: infra.Docker})
	require.NoError(t, err)
	return &cre.Topology{DonsMetadata: dm}
}

func TestWorkflowContainerPatternForDON(t *testing.T) {
	t.Parallel()

	require.Equal(t, "feeds-zone-a-node", workflowContainerPatternForDON(&cre.DonMetadata{Name: "feeds-zone-a"}))
	require.Equal(t, "shard0-node", workflowContainerPatternForDON(&cre.DonMetadata{Name: "shard0"}))
}

// Single-DON CI topologies need no selector flags.
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

// --workflow-don-name wins over --don-family when both are present in the selector.
func TestResolveWorkflowDONMetadata_explicitNameWinsOverDonFamily(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, shardedWorkflowDONTestTopology(t))

	don, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{
		ExplicitName: "shard1",
		DonFamily:    envconfig.DefaultDONFamily,
	})
	require.NoError(t, err)
	require.Equal(t, "shard1", don.Name)
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

func TestResolveWorkflowDONMetadata_shardedByFamily_shardIndex0(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, shardedWorkflowDONTestTopology(t))

	shard0 := uint(0)
	don, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{
		DonFamily:  envconfig.DefaultDONFamily,
		ShardIndex: &shard0,
	})
	require.NoError(t, err)
	require.Equal(t, "shard0", don.Name)
}

// shard0 and shard1 share don_family — family alone is ambiguous without --shard-index.
func TestResolveWorkflowDONMetadata_shardedRequiresShardIndex(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, shardedWorkflowDONTestTopology(t))

	_, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{DonFamily: envconfig.DefaultDONFamily})
	require.Error(t, err)
	require.Contains(t, err.Error(), "--shard-index")
}

func TestResolveWorkflowDONMetadata_unknownDonFamily(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, shardedWorkflowDONTestTopology(t))

	_, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{DonFamily: "unknown-family"})
	require.Error(t, err)
	require.Contains(t, err.Error(), `no workflow DON with don_family "unknown-family"`)
}

func TestResolveWorkflowDONMetadata_unknownShardIndex(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, shardedWorkflowDONTestTopology(t))

	shard99 := uint(99)
	_, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{
		DonFamily:  envconfig.DefaultDONFamily,
		ShardIndex: &shard99,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "shard_index 99")
}

// Multiple non-shard workflow DONs sharing a family require --workflow-don-name.
func TestResolveWorkflowDONMetadata_nonShardDuplicateFamily(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, testTopologyWithWorkflowDONs(t,
		&cre.DonMetadata{Name: "wf-a", ID: 1, DonFamily: "shared-family", Flags: []string{cre.WorkflowDON}},
		&cre.DonMetadata{Name: "wf-b", ID: 2, DonFamily: "shared-family", Flags: []string{cre.WorkflowDON}},
	))

	_, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{DonFamily: "shared-family"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "--workflow-don-name")
}

// Table-driven unit tests for resolveWorkflowDONByFamily (helper used by ResolveWorkflowDONMetadata).
func TestResolveWorkflowDONByFamily(t *testing.T) {
	t.Parallel()

	wfDONs := []*cre.DonMetadata{
		{Name: "shard0", DonFamily: envconfig.DefaultDONFamily, ShardIndex: 0, Flags: []string{cre.WorkflowDON, cre.ShardDON}},
		{Name: "shard1", DonFamily: envconfig.DefaultDONFamily, ShardIndex: 1, Flags: []string{cre.WorkflowDON, cre.ShardDON}},
		{Name: "feeds-zone-a", DonFamily: "feeds-zone-a", Flags: []string{cre.WorkflowDON}},
	}

	tests := []struct {
		name       string
		family     string
		shardIndex *uint
		wantName   string
		wantErr    string
	}{
		{
			name:     "unique family",
			family:   "feeds-zone-a",
			wantName: "feeds-zone-a",
		},
		{
			name:       "sharded family with index",
			family:     envconfig.DefaultDONFamily,
			shardIndex: new(uint(1)),
			wantName:   "shard1",
		},
		{
			name:    "unknown family",
			family:  "missing",
			wantErr: `no workflow DON with don_family "missing"`,
		},
		{
			name:    "sharded family without index",
			family:  envconfig.DefaultDONFamily,
			wantErr: "--shard-index",
		},
		{
			name:       "unknown shard index",
			family:     envconfig.DefaultDONFamily,
			shardIndex: new(uint(99)),
			wantErr:    "shard_index 99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			don, err := resolveWorkflowDONByFamily(wfDONs, tt.family, tt.shardIndex)
			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantName, don.Name)
		})
	}
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
		DonFamily:     envconfig.DefaultDONFamily,
		NodesMetadata: []*cre.NodeMetadata{{Roles: []string{cre.BootstrapNode}}},
	}
}
