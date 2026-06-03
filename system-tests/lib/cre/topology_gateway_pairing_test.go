package cre

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
)

func TestGatewayConnectorsForWorkflow_shardedPairing(t *testing.T) {
	t.Parallel()

	topology := shardedGatewayTestTopology()

	connA := topology.GatewayConnectorsForWorkflow("feeds-zone-a")
	require.Len(t, connA.Configurations, 1)
	require.Equal(t, "gateway-node-0", connA.Configurations[0].AuthGatewayID)

	connB := topology.GatewayConnectorsForWorkflow("feeds-zone-b")
	require.Len(t, connB.Configurations, 1)
	require.Equal(t, "gateway-node-1", connB.Configurations[0].AuthGatewayID)
}

func TestGatewayConnectorsForWorkflow_singleGatewayLegacy(t *testing.T) {
	t.Parallel()

	topology := &Topology{
		DonsMetadata: &DonsMetadata{dons: []*DonMetadata{
			{Name: "feeds-zone-a", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "feeds-zone-b", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "gateway", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
		}},
		GatewayConnectors: &GatewayConnectors{Configurations: []*DonGatewayConfiguration{
			{GatewayConfiguration: &GatewayConfiguration{AuthGatewayID: "gateway-node-0"}},
		}},
	}

	conn := topology.GatewayConnectorsForWorkflow("feeds-zone-b")
	require.Len(t, conn.Configurations, 1)
	require.Equal(t, "gateway-node-0", conn.Configurations[0].AuthGatewayID)
}

func TestGatewayConnectorsForWorkflow_shardedPairing_unknownWorkflow(t *testing.T) {
	t.Parallel()

	topology := shardedGatewayTestTopology()

	conn := topology.GatewayConnectorsForWorkflow("unknown-don")
	require.Empty(t, conn.Configurations)
}

func TestGatewayConnectorsForWorkflow_mismatchedCountsUsesLegacy(t *testing.T) {
	t.Parallel()

	topology := &Topology{
		DonsMetadata: &DonsMetadata{dons: []*DonMetadata{
			{Name: "feeds-zone-a", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "feeds-zone-b", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "feeds-zone-c", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "gateway-zone-a", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			{Name: "gateway-zone-b", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
		}},
		GatewayConnectors: &GatewayConnectors{Configurations: []*DonGatewayConfiguration{
			{GatewayConfiguration: &GatewayConfiguration{AuthGatewayID: "gateway-node-0"}},
			{GatewayConfiguration: &GatewayConfiguration{AuthGatewayID: "gateway-node-1"}},
		}},
	}

	require.False(t, topology.shardedGatewayPairingEnabled())

	conn := topology.GatewayConnectorsForWorkflow("feeds-zone-c")
	require.Len(t, conn.Configurations, 2)
}

func TestGatewayConnectorsForWorkflow_dataFeedsLocalTopology(t *testing.T) {
	t.Parallel()

	topology := shardedGatewayTestTopology()

	require.True(t, topology.shardedGatewayPairingEnabled())
	gatewayNames, workflowNames := topology.gatewayAndWorkflowDONNames()
	require.Equal(t, []string{"gateway-zone-a", "gateway-zone-b"}, gatewayNames)
	require.Equal(t, []string{"feeds-zone-a", "feeds-zone-b"}, workflowNames)
}

func TestGatewayServiceConfigsForGateway_shardedPairing(t *testing.T) {
	t.Parallel()

	topology := shardedGatewayTestTopology()
	services := []GatewayServiceConfig{{
		ServiceName: pkg.ServiceNameWorkflows,
		Handlers:    []string{pkg.GatewayHandlerTypeWebAPICapabilities},
		DONs:        []string{"feeds-zone-a", "feeds-zone-b"},
	}}

	scopedA := topology.GatewayServiceConfigsForGateway("gateway-zone-a", services)
	require.Equal(t, []string{"feeds-zone-a"}, scopedA[0].DONs)

	scopedB := topology.GatewayServiceConfigsForGateway("gateway-zone-b", services)
	require.Equal(t, []string{"feeds-zone-b"}, scopedB[0].DONs)
}

func shardedGatewayTestTopology() *Topology {
	return &Topology{
		DonsMetadata: &DonsMetadata{dons: []*DonMetadata{
			{Name: "feeds-zone-a", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "feeds-zone-b", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "gateway-zone-a", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			{Name: "gateway-zone-b", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
		}},
		GatewayConnectors: &GatewayConnectors{Configurations: []*DonGatewayConfiguration{
			{GatewayConfiguration: &GatewayConfiguration{AuthGatewayID: "gateway-node-0"}},
			{GatewayConfiguration: &GatewayConfiguration{AuthGatewayID: "gateway-node-1"}},
		}},
	}
}
