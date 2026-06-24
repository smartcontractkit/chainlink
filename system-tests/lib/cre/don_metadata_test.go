package cre

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDonMetadata_IsWorkflowDON_nilNodeSet(t *testing.T) {
	t.Parallel()

	bootstrap := &DonMetadata{
		Name:          "bootstrap",
		NodesMetadata: []*NodeMetadata{{Roles: []string{BootstrapNode}}},
	}
	require.False(t, bootstrap.IsWorkflowDON())
	require.False(t, bootstrap.IsShardDON())
}

func TestDonsMetadata_WorkflowDONs_skipsBootstrapWithoutNodeSet(t *testing.T) {
	t.Parallel()

	dm := uncheckedDonsMetadata([]*DonMetadata{
		{
			Name:          "bootstrap",
			NodesMetadata: []*NodeMetadata{{Roles: []string{BootstrapNode}}},
		},
		{Name: "feeds-zone-a", ID: 1, Flags: []string{WorkflowDON}},
		{Name: "feeds-zone-b", ID: 2, Flags: []string{WorkflowDON}},
	})

	wfDONs, err := dm.WorkflowDONs()
	require.NoError(t, err)
	require.Len(t, wfDONs, 2)
	require.Equal(t, "feeds-zone-a", wfDONs[0].Name)
	require.Equal(t, "feeds-zone-b", wfDONs[1].Name)
}

func uncheckedDonsMetadata(dons []*DonMetadata) *DonsMetadata {
	return &DonsMetadata{dons: dons}
}
