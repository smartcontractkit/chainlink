package environment

import (
	"testing"

	"github.com/rs/zerolog"
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

func TestResolveRemoteRuntimeForSetupSkipsResolutionWhenNoRemoteComponents(t *testing.T) {
	execPlan, planErr := buildPlacementPlan(
		[]*config.Blockchain{{Placement: config.PlacementLocal}},
		&config.JobDistributor{Placement: config.PlacementLocal},
		[]*cre.NodeSet{{Placement: "local"}},
	)
	require.NoError(t, planErr)

	runtime, err := resolveRemoteRuntimeForSetup(
		zerolog.Nop(),
		execPlan,
	)
	require.NoError(t, err)
	require.Nil(t, runtime, "expected nil runtime when no remote components are configured")
}

func TestBuildExecutionPlanIncludesPlacementAndRemoteFlags(t *testing.T) {
	execPlan, err := buildPlacementPlan(
		[]*config.Blockchain{{Placement: config.PlacementRemote}},
		&config.JobDistributor{Placement: config.PlacementLocal},
		[]*cre.NodeSet{{Placement: "local"}, {Placement: "remote"}},
	)
	require.NoError(t, err, "expected execution plan build to succeed")
	require.NotNil(t, execPlan, "expected non-nil execution plan")
	require.NotNil(t, execPlan.NodeSetPlacement, "expected nodeset placement summary")
	require.True(t, execPlan.NodeSetPlacement.HasLocalTargets, "expected local nodeset placement")
	require.True(t, execPlan.NodeSetPlacement.HasRemoteTargets, "expected remote nodeset placement")
	require.True(t, execPlan.HasRemoteComponents, "expected remote components flag")
}
