package environment

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
)

func TestRewriteReconstructedGatewayIncomingHosts_RemoteGatewayUsesEC2IP(t *testing.T) {
	topology, nodeSet := mustBuildRemoteGatewayTopology(t)
	cfg := &config.Config{
		NodeSets: []*cre.NodeSet{nodeSet},
	}
	t.Setenv(runtimecfg.EnvRemoteHostIP, "203.0.113.10")

	err := rewriteReconstructedGatewayIncomingHosts(cfg, topology)
	require.NoError(t, err, "expected remote gateway incoming rewrite to succeed")
	require.Equal(
		t,
		"203.0.113.10",
		topology.GatewayConnectors.Configurations[0].Incoming.Host,
		"expected reconstructed remote gateway incoming host to use EC2 IP",
	)
}

func TestRewriteReconstructedGatewayIncomingHosts_LocalGatewayNoop(t *testing.T) {
	topology, nodeSet := mustBuildRemoteGatewayTopology(t)
	nodeSet.Placement = string(config.PlacementLocal)
	cfg := &config.Config{
		NodeSets: []*cre.NodeSet{nodeSet},
	}
	t.Setenv(runtimecfg.EnvRemoteHostIP, "")
	t.Setenv(runtimecfg.EnvRemoteAgentEC2InstanceID, "")

	err := rewriteReconstructedGatewayIncomingHosts(cfg, topology)
	require.NoError(t, err, "expected local gateway reconstruction rewrite to be a no-op")
	require.Equal(
		t,
		"bootstrap-gateway-node0",
		topology.GatewayConnectors.Configurations[0].Incoming.Host,
		"expected local gateway incoming host to remain unchanged",
	)
}

func TestRewriteReconstructedGatewayIncomingHosts_RewritesOnlyRemoteNodeSets(t *testing.T) {
	remoteTopology, remoteNodeSet := mustBuildRemoteGatewayTopology(t)
	localTopology, localNodeSet := mustBuildRemoteGatewayTopology(t)
	localNodeSet.Placement = string(config.PlacementLocal)
	// Preserve the remote topology gateway config and append a local-only gateway config.
	remoteTopology.GatewayConnectors.Configurations = append(
		remoteTopology.GatewayConnectors.Configurations,
		localTopology.GatewayConnectors.Configurations[0],
	)

	cfg := &config.Config{
		NodeSets: []*cre.NodeSet{remoteNodeSet, localNodeSet},
	}
	t.Setenv(runtimecfg.EnvRemoteHostIP, "203.0.113.77")

	err := rewriteReconstructedGatewayIncomingHosts(cfg, remoteTopology)
	require.NoError(t, err, "expected mixed reconstruction rewrite to succeed")
	require.Equal(t, "203.0.113.77", remoteTopology.GatewayConnectors.Configurations[0].Incoming.Host, "expected remote gateway incoming host rewrite")
	require.Equal(t, "bootstrap-gateway-node0", remoteTopology.GatewayConnectors.Configurations[1].Incoming.Host, "expected local gateway incoming host to remain unchanged")
}
