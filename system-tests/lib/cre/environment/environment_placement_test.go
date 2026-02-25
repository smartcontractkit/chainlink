package environment

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func TestSummarizeNodeSetPlacement_AllowsMixedPlacements(t *testing.T) {
	nodeSets := []*cre.NodeSet{
		{Placement: "local"},
		{Placement: "remote"},
	}

	summary, err := summarizeNodeSetPlacement(nodeSets)
	require.NoError(t, err, "summarizeNodeSetPlacement should succeed")
	require.True(t, summary.HasLocalTargets, "expected local placement to be detected")
	require.True(t, summary.HasRemoteTargets, "expected remote placement to be detected")
}

func TestHasRemoteComponents(t *testing.T) {
	tests := []struct {
		name       string
		blockchains []*config.Blockchain
		jd         *config.JobDistributor
		nodeSets   []*cre.NodeSet
		want       bool
	}{
		{
			name: "none remote",
			blockchains: []*config.Blockchain{
				{Placement: config.PlacementLocal},
			},
			jd:       &config.JobDistributor{Placement: config.PlacementLocal},
			nodeSets: []*cre.NodeSet{{Placement: "local"}},
			want:     false,
		},
		{
			name: "remote blockchain",
			blockchains: []*config.Blockchain{
				{Placement: config.PlacementRemote},
			},
			want: true,
		},
		{
			name: "remote jd",
			jd:   &config.JobDistributor{Placement: config.PlacementRemote},
			want: true,
		},
		{
			name:     "remote nodeset",
			nodeSets: []*cre.NodeSet{{Placement: "remote"}},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasRemoteComponents(tt.blockchains, tt.jd, tt.nodeSets)
			require.Equalf(t, tt.want, got, "expected hasRemoteComponents() to return %v", tt.want)
		})
	}
}
