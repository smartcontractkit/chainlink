package environment

import (
	"net/url"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/stretchr/testify/require"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/tunnel"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
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
	err := rewriteRemoteNodeSetOutputForLocalAccess(nil, 0, nil, nil, "203.0.113.10")
	require.NoError(t, err, "expected local-only no-op rewrite to succeed")
}

func TestRewriteRemoteNodeSetOutputForLocalAccess_RemoteRewritesGatewayIncomingHost(t *testing.T) {
	topology, nodeSet := mustBuildRemoteGatewayTopology(t)
	output := &simple_node_set.Output{}

	err := rewriteRemoteNodeSetOutputForLocalAccess(topology, 0, nodeSet, output, "203.0.113.10")
	require.NoError(t, err, "expected remote rewrite to succeed")

	require.NotNil(t, topology.GatewayConnectors)
	require.Len(t, topology.GatewayConnectors.Configurations, 1)
	require.Equal(
		t,
		"203.0.113.10",
		topology.GatewayConnectors.Configurations[0].Incoming.Host,
		"expected remote nodeset rewrite to expose gateway incoming via EC2 host",
	)
}

func TestRewriteRemoteNodeSetOutputForLocalAccess_InvalidNodeExternalURLFails(t *testing.T) {
	topology, nodeSet := mustBuildRemoteGatewayTopology(t)
	output := &simple_node_set.Output{
		CLNodes: []*clnode.Output{
			{
				Node: &clnode.NodeOut{
					ExternalURL: "://bad-url",
				},
			},
		},
	}

	err := rewriteRemoteNodeSetOutputForLocalAccess(topology, 0, nodeSet, output, "203.0.113.10")
	require.Error(t, err, "expected invalid node external URL to fail rewrite")
	require.Contains(t, err.Error(), "failed to parse url", "expected parse failure context")
}

func TestParseCustomPortMapping(t *testing.T) {
	t.Run("valid mapping", func(t *testing.T) {
		hostPort, containerPort, err := parseCustomPortMapping("127.0.0.1:18080:8080")
		require.NoError(t, err, "expected valid mapping to parse")
		require.Equal(t, 18080, hostPort)
		require.Equal(t, 8080, containerPort)
	})

	t.Run("missing separator", func(t *testing.T) {
		_, _, err := parseCustomPortMapping("8080")
		require.Error(t, err, "expected malformed mapping to fail")
		require.Contains(t, err.Error(), "expected hostPort:containerPort")
	})

	t.Run("invalid host port", func(t *testing.T) {
		_, _, err := parseCustomPortMapping("bad:8080")
		require.Error(t, err, "expected invalid host port to fail")
		require.Contains(t, err.Error(), "invalid host port")
	})

	t.Run("invalid container port", func(t *testing.T) {
		_, _, err := parseCustomPortMapping("18080:bad")
		require.Error(t, err, "expected invalid container port to fail")
		require.Contains(t, err.Error(), "invalid container port")
	})
}

func TestNodeSetResolveURLPort(t *testing.T) {
	tests := []struct {
		name      string
		rawURL    string
		wantPort  int
		wantError string
	}{
		{name: "explicit port", rawURL: "http://node:1234", wantPort: 1234},
		{name: "http default", rawURL: "http://node", wantPort: 80},
		{name: "https default", rawURL: "https://node", wantPort: 443},
		{name: "unsupported scheme without port", rawURL: "tcp://node", wantError: "unsupported scheme"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := url.Parse(tt.rawURL)
			require.NoError(t, err, "test setup should parse URL")
			port, err := nodeSetResolveURLPort(parsed)
			if tt.wantError != "" {
				require.Error(t, err, "expected port resolution failure")
				require.Contains(t, err.Error(), tt.wantError)
				return
			}
			require.NoError(t, err, "expected port resolution success")
			require.Equal(t, tt.wantPort, port)
		})
	}
}

func TestNodeSetEndpointFromURL(t *testing.T) {
	ref, err := nodeSetEndpointFromURL("nodeset:workflow", "node-0-api", "http://node-0:8081")
	require.NoError(t, err, "expected endpoint ref to parse")
	require.Equal(t, "nodeset:workflow", ref.ComponentID)
	require.Equal(t, "node-0-api", ref.EndpointName)
	require.Equal(t, "http", ref.Scheme)
	require.Equal(t, "node-0", ref.Host)
	require.Equal(t, 8081, ref.Port)

	ref, err = nodeSetEndpointFromURL("nodeset:workflow", "node-0-api", "   ")
	require.NoError(t, err, "expected blank URL to be ignored")
	require.Nil(t, ref, "expected nil endpoint for blank URL")

	_, err = nodeSetEndpointFromURL("nodeset:workflow", "node-0-api", "http://")
	require.Error(t, err, "expected empty hostname to fail")
	require.Contains(t, err.Error(), "empty hostname")
}

func TestGatewayLocalPortFromBindings(t *testing.T) {
	bindings := []tunnel.TunnelBinding{
		{EndpointRef: tunnel.EndpointRef{EndpointName: "node-0-custom-0-5002"}, LocalPort: 22002},
		{EndpointRef: tunnel.EndpointRef{EndpointName: "node-1-custom-0-5002"}, LocalPort: 22012},
	}

	port, ok := gatewayLocalPortFromBindings(0, 5002, bindings)
	require.True(t, ok, "expected matching binding")
	require.Equal(t, 22002, port)

	_, ok = gatewayLocalPortFromBindings(0, 6000, bindings)
	require.False(t, ok, "expected non-matching container port to return false")
}

func TestRewriteGatewayIncomingForNodeSetBindings(t *testing.T) {
	topology, nodeSet := mustBuildRemoteGatewayTopology(t)
	nodeSet.NodeSpecs[0].Input.Node.CustomPorts = []string{"18080:5002"}

	bindings := []tunnel.TunnelBinding{
		{EndpointRef: tunnel.EndpointRef{EndpointName: "node-0-custom-0-5002"}, LocalPort: 22002},
	}

	rewriteGatewayIncomingForNodeSetBindings(topology, 0, nodeSet, bindings)
	cfg := topology.GatewayConnectors.Configurations[0].GatewayConfiguration
	require.Equal(t, "127.0.0.1", cfg.Incoming.Host, "incoming host should be local during binding mode")
	require.Equal(t, 22002, cfg.Incoming.ExternalPort, "incoming external port should be rewritten from binding")
}

func TestRewriteCustomPortMappingHostPort(t *testing.T) {
	rewritten := rewriteCustomPortMappingHostPort("127.0.0.1:18080:8080", 22080)
	require.Equal(t, "127.0.0.1:22080:8080", rewritten)

	unchanged := rewriteCustomPortMappingHostPort("bad", 22080)
	require.True(t, strings.EqualFold("bad", unchanged), "malformed mapping should remain unchanged")
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
