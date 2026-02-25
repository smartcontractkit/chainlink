package environment

import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/stretchr/testify/require"
)

func TestBuildRemoteNodeSetInputRequiresImageOrBuildFields(t *testing.T) {
	nodeSet := &cre.NodeSet{
		Input: &simple_node_set.Input{
			Name: "remote-don",
		},
		NodeSpecs: []*cre.NodeSpecWithRole{
			{
				Input: &clnode.Input{
					Node: &clnode.NodeInput{
						Image: "",
					},
				},
			},
		},
	}

	_, err := buildRemoteNodeSetInput(nodeSet)
	require.Error(t, err, "expected missing image/build validation error")
	require.Contains(t, err.Error(), "must set node.image or docker build fields", "expected image validation error")
}

func TestBuildRemoteNodeSetInputRejectsImageAndBuildFieldsTogether(t *testing.T) {
	nodeSet := &cre.NodeSet{
		Input: &simple_node_set.Input{
			Name: "remote-don",
		},
		NodeSpecs: []*cre.NodeSpecWithRole{
			{
				Input: &clnode.Input{
					Node: &clnode.NodeInput{
						Image:          "repo/chainlink:tag",
						DockerContext:  "../../../..",
						DockerFilePath: "core/chainlink.Dockerfile",
					},
				},
			},
		},
	}

	_, err := buildRemoteNodeSetInput(nodeSet)
	require.Error(t, err, "expected image+build conflict validation error")
	require.Contains(t, err.Error(), "either node.image or docker build fields", "expected image/build conflict error")
}

func TestRewriteRemoteNodeSetOutputForLocalAccess_LocalOnlyNoop(t *testing.T) {
	err := rewriteRemoteNodeSetOutputForLocalAccess(nil, "203.0.113.10")
	require.NoError(t, err, "expected local-only no-op rewrite to succeed")
}

func TestNormalizeForExecution_RemoteRewritesGatewayIncomingHost(t *testing.T) {
	topology, nodeSet := mustBuildRemoteGatewayTopology(t)
	normalizeForExecution(topology, []*cre.NodeSet{nodeSet}, "203.0.113.10")

	require.NotNil(t, topology.GatewayConnectors)
	require.Len(t, topology.GatewayConnectors.Configurations, 1)
	require.Equal(
		t,
		"203.0.113.10",
		topology.GatewayConnectors.Configurations[0].Incoming.Host,
		"expected remote nodeset rewrite to expose gateway incoming via EC2 host",
	)
}

func TestRewriteRemoteNodeSetOutputForLocalAccess_RemoteRewritesNodeExternalURL(t *testing.T) {
	output := &simple_node_set.Output{
		CLNodes: []*clnode.Output{
			{
				Node: &clnode.NodeOut{
					ExternalURL: "http://127.0.0.1:6688",
				},
			},
		},
	}

	err := rewriteRemoteNodeSetOutputForLocalAccess(output, "203.0.113.10")
	require.NoError(t, err, "expected remote rewrite to succeed")
	require.Equal(t, "http://203.0.113.10:6688", output.CLNodes[0].Node.ExternalURL)
}

func TestRewriteRemoteNodeSetOutputForLocalAccess_InvalidNodeExternalURLFails(t *testing.T) {
	output := &simple_node_set.Output{
		CLNodes: []*clnode.Output{
			{
				Node: &clnode.NodeOut{
					ExternalURL: "://bad-url",
				},
			},
		},
	}

	err := rewriteRemoteNodeSetOutputForLocalAccess(output, "203.0.113.10")
	require.Error(t, err, "expected invalid node external URL to fail rewrite")
	require.Contains(t, err.Error(), "failed to parse address", "expected parse failure context")
}

func mustBuildRemoteGatewayTopology(t *testing.T) (*cre.Topology, *cre.NodeSet) {
	t.Helper()

	provider := infra.Provider{Type: infra.Docker}
	nodeSet := &cre.NodeSet{
		Input: &simple_node_set.Input{Name: "workflow"},
		NodeSpecs: []*cre.NodeSpecWithRole{
			{
				Input: &clnode.Input{Node: &clnode.NodeInput{}},
				Roles: []cre.NodeType{cre.BootstrapNode, cre.GatewayNode},
			},
		},
		Placement: "remote",
	}

	donMetadata, err := cre.NewDonMetadata(nodeSet, 1, provider, nil)
	require.NoError(t, err, "failed to build DonMetadata")
	donsMetadata, err := cre.NewDonsMetadata([]*cre.DonMetadata{donMetadata}, provider)
	require.NoError(t, err, "failed to build DonsMetadata")

	gatewayNode, hasGateway := donMetadata.Gateway()
	require.True(t, hasGateway, "expected gateway node in metadata")

	topology := &cre.Topology{
		DonsMetadata: donsMetadata,
		GatewayConnectors: &cre.GatewayConnectors{
			Configurations: []*cre.DonGatewayConfiguration{
				{
					GatewayConfiguration: &cre.GatewayConfiguration{
						NodeUUID: gatewayNode.UUID,
						Incoming: cre.Incoming{
							Host:         "bootstrap-gateway-node0",
							ExternalPort: 5002,
						},
					},
				},
			},
		},
	}
	return topology, nodeSet
}
