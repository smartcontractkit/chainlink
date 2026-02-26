package environment

import (
	"testing"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
	"github.com/stretchr/testify/require"
)

func TestFilteredRemoteStopConfigKeepsOnlyRemoteComponents(t *testing.T) {
	cfg := &envconfig.Config{
		Blockchains: []*envconfig.Blockchain{
			{Placement: envconfig.PlacementLocal},
			{Placement: envconfig.PlacementRemote},
		},
		NodeSets: []*cre.NodeSet{
			{Placement: "local"},
			{Placement: "remote"},
		},
		JD: &envconfig.JobDistributor{Placement: envconfig.PlacementRemote},
	}

	filtered := filteredRemoteStopConfig(cfg)
	require.Len(t, filtered.Blockchains, 1)
	require.Equal(t, envconfig.PlacementRemote, filtered.Blockchains[0].Placement)
	require.Len(t, filtered.NodeSets, 1)
	require.Equal(t, "remote", filtered.NodeSets[0].Placement)
	require.NotNil(t, filtered.JD)
	require.Equal(t, envconfig.PlacementRemote, filtered.JD.Placement)
}

func TestCaptureRemoteAgentStateReadsExpectedEnvVars(t *testing.T) {
	t.Setenv(envRemoteAgentURL, "http://203.0.113.10:8080")
	t.Setenv(runtimecfg.EnvRemoteAgentEC2InstanceID, "i-abc")
	t.Setenv(envRemoteAgentPort, "18080")
	t.Setenv("AWS_PROFILE", "fallback-profile")

	state := captureRemoteAgentState()
	require.Equal(t, "http://203.0.113.10:8080", state.RemoteAgentURL)
	require.Equal(t, "i-abc", state.RemoteAgentEC2InstanceID)
	require.Equal(t, "18080", state.RemoteAgentPort)
	require.Equal(t, "fallback-profile", state.AWSProfile)
}
