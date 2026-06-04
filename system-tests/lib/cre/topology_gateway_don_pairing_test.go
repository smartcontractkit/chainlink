package cre

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
)

func TestGatewayConnectorsForWorkflow_gatewayDonPairing(t *testing.T) {
	t.Parallel()

	topology := gatewayDonPairingTestTopology(true)

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

func TestGatewayConnectorsForWorkflow_gatewayDonPairing_unknownWorkflow(t *testing.T) {
	t.Parallel()

	topology := gatewayDonPairingTestTopology(true)

	conn := topology.GatewayConnectorsForWorkflow("unknown-don")
	require.Empty(t, conn.Configurations)
}

func TestGatewayConnectorsForWorkflow_equalCountsWithoutFlagUsesLegacy(t *testing.T) {
	t.Parallel()

	topology := gatewayDonPairingTestTopology(false)

	require.False(t, topology.gatewayDonPairingEnabled())

	conn := topology.GatewayConnectorsForWorkflow("feeds-zone-b")
	require.Len(t, conn.Configurations, 2)
}

func TestTopology_validateGatewayDONPairing_countMismatch(t *testing.T) {
	t.Parallel()

	topology := &Topology{
		GatewayDONPairing: true,
		DonsMetadata: &DonsMetadata{dons: []*DonMetadata{
			{Name: "feeds-zone-a", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "feeds-zone-b", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "feeds-zone-c", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "gateway-zone-a", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			{Name: "gateway-zone-b", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
		}},
	}

	err := topology.validateGatewayDONPairing()
	require.Error(t, err)
	require.Contains(t, err.Error(), "cre_topology.gateway_don_pairing requires equal workflow and gateway DON counts")
}

func TestGatewayConnectorsForWorkflow_dataFeedsLocalTopology(t *testing.T) {
	t.Parallel()

	topology := gatewayDonPairingTestTopology(true)

	require.True(t, topology.gatewayDonPairingEnabled())
	gatewayNames, workflowNames := topology.gatewayAndWorkflowDONNames()
	require.Equal(t, []string{"gateway-zone-a", "gateway-zone-b"}, gatewayNames)
	require.Equal(t, []string{"feeds-zone-a", "feeds-zone-b"}, workflowNames)

	pairs := topology.GatewayDONPairings()
	require.Equal(t, [][2]string{
		{"feeds-zone-a", "gateway-zone-a"},
		{"feeds-zone-b", "gateway-zone-b"},
	}, pairs)
}

func TestGatewayServiceConfigsForGateway_gatewayDonPairing(t *testing.T) {
	t.Parallel()

	topology := gatewayDonPairingTestTopology(true)
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

func gatewayDonPairingTestTopology(enabled bool) *Topology {
	return &Topology{
		GatewayDONPairing: enabled,
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
