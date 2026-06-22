package environment

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func testTopologyWithWorkflowDONs(t *testing.T, wfDONs ...*cre.DonMetadata) *cre.Topology {
	t.Helper()

	dons := append([]*cre.DonMetadata{bootstrapDONMetadata()}, wfDONs...)
	return &cre.Topology{DonsMetadata: cre.NewUncheckedDonsMetadata(dons)}
}

func TestContainerPatternMatchesDON(t *testing.T) {
	t.Parallel()

	don := &cre.DonMetadata{Name: "feeds-zone-a"}

	require.True(t, containerPatternMatchesDON("feeds-zone-a", don))
	require.True(t, containerPatternMatchesDON("feeds-zone-a-node", don))
	require.True(t, containerPatternMatchesDON("feeds-zone-a-node-0", don))
	require.False(t, containerPatternMatchesDON("workflow-node", don))
	require.False(t, containerPatternMatchesDON("feeds-zone-b-node", don))
}

func TestContainerPatternMatchesDON_prefixAmbiguity(t *testing.T) {
	t.Parallel()

	parent := &cre.DonMetadata{Name: "feeds"}
	child := &cre.DonMetadata{Name: "feeds-zone-a"}

	require.False(t, containerPatternMatchesDON("feeds-zone-a-node", parent))
	require.True(t, containerPatternMatchesDON("feeds-zone-a-node", child))
}

func TestResolveWorkflowDONMetadata_singleDONLegacy(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, testTopologyWithWorkflowDONs(t,
		&cre.DonMetadata{Name: "workflow", ID: 1, Flags: []string{cre.WorkflowDON}},
	))

	don, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{ContainerPattern: "workflow-node"})
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

func TestResolveWorkflowDONMetadata_multiDONContainerPattern(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, multiWorkflowDONTestTopology(t))

	don, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{ContainerPattern: "feeds-zone-a-node"})
	require.NoError(t, err)
	require.Equal(t, "feeds-zone-a", don.Name)
}

func TestResolveWorkflowDONMetadata_multiDONGenericPatternFails(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, multiWorkflowDONTestTopology(t))

	_, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{ContainerPattern: "workflow-node"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "--workflow-don-name")
}

func TestResolveWorkflowDONMetadata_multiDONRequiresSelector(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, multiWorkflowDONTestTopology(t))

	_, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "--workflow-don-name")
}

func TestResolveWorkflowDONMetadata_unknownExplicitName(t *testing.T) {
	t.Parallel()

	resolver := testStateResolver(t, multiWorkflowDONTestTopology(t))

	_, err := resolver.ResolveWorkflowDONMetadata(workflowDONSelector{ExplicitName: "unknown"})
	require.Error(t, err)
	require.Contains(t, err.Error(), `workflow DON "unknown" not found`)
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

func bootstrapDONMetadata() *cre.DonMetadata {
	return &cre.DonMetadata{
		Name:          "bootstrap",
		NodesMetadata: []*cre.NodeMetadata{{Roles: []string{cre.BootstrapNode}}},
	}
}
