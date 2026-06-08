package cre

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
)

func TestResolveNodesetZone_explicit(t *testing.T) {
	t.Parallel()
	require.Equal(t, "zone-a", ResolveNodesetZone("feeds-zone-a", "zone-a"))
	require.Equal(t, "custom", ResolveNodesetZone("feeds-zone-a", "custom"))
}

func TestResolveNodesetZone_nameSuffix(t *testing.T) {
	t.Parallel()
	require.Equal(t, "zone-b", ResolveNodesetZone("gateway-zone-b", ""))
	require.Empty(t, ResolveNodesetZone("gateway", ""))
}

func TestGatewayConnectorsForWorkflow_zonePairing(t *testing.T) {
	t.Parallel()

	topology := zoneGatewayPairingTestTopology(true)

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
			{Name: "feeds-zone-a", Zone: "zone-a", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "feeds-zone-b", Zone: "zone-b", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
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

func TestGatewayConnectorsForWorkflow_zonePairing_unknownWorkflow(t *testing.T) {
	t.Parallel()

	topology := zoneGatewayPairingTestTopology(true)

	conn := topology.GatewayConnectorsForWorkflow("unknown-don")
	require.Empty(t, conn.Configurations)
}

func TestGatewayConnectorsForWorkflow_withoutZonesUsesLegacy(t *testing.T) {
	t.Parallel()

	topology := zoneGatewayPairingTestTopology(false)

	require.False(t, topology.zoneGatewayPairingEnabled())

	conn := topology.GatewayConnectorsForWorkflow("feeds-zone-b")
	require.Len(t, conn.Configurations, 2)
}

func TestTopology_validateZoneGatewayPairing_missingGateway(t *testing.T) {
	t.Parallel()

	topology := &Topology{
		DonsMetadata: &DonsMetadata{dons: []*DonMetadata{
			{Name: "feeds-zone-a", Zone: "zone-a", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "feeds-zone-b", Zone: "zone-b", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "feeds-zone-c", Zone: "zone-c", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "gateway-zone-a", Zone: "zone-a", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			{Name: "gateway-zone-b", Zone: "zone-b", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
		}},
	}

	err := topology.validateZoneGatewayPairing()
	require.Error(t, err)
	require.Contains(t, err.Error(), "no gateway DON is defined for that zone")
}

func TestGatewayConnectorsForWorkflow_dataFeedsLocalTopology(t *testing.T) {
	t.Parallel()

	topology := zoneGatewayPairingTestTopology(true)

	require.True(t, topology.zoneGatewayPairingEnabled())

	pairs := topology.GatewayZonePairings()
	require.Equal(t, [][2]string{
		{"feeds-zone-a", "gateway-zone-a"},
		{"feeds-zone-b", "gateway-zone-b"},
	}, pairs)
}

func TestGatewayServiceConfigsForGateway_zonePairing(t *testing.T) {
	t.Parallel()

	topology := zoneGatewayPairingTestTopology(true)
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

func zoneGatewayPairingTestTopology(zoned bool) *Topology {
	zoneA, zoneB := "", ""
	if zoned {
		zoneA, zoneB = "zone-a", "zone-b"
	}

	return &Topology{
		DonsMetadata: &DonsMetadata{dons: []*DonMetadata{
			{Name: "feeds-zone-a", Zone: zoneA, ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "feeds-zone-b", Zone: zoneB, ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON}},
			{Name: "gateway-zone-a", Zone: zoneA, ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			{Name: "gateway-zone-b", Zone: zoneB, ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
		}},
		GatewayConnectors: &GatewayConnectors{Configurations: []*DonGatewayConfiguration{
			{GatewayConfiguration: &GatewayConfiguration{AuthGatewayID: "gateway-node-0"}, DONName: "gateway-zone-a"},
			{GatewayConfiguration: &GatewayConfiguration{AuthGatewayID: "gateway-node-1"}, DONName: "gateway-zone-b"},
		}},
		gatewayConnectorsByDon: map[string]*DonGatewayConfiguration{
			"gateway-zone-a": {GatewayConfiguration: &GatewayConfiguration{AuthGatewayID: "gateway-node-0"}, DONName: "gateway-zone-a"},
			"gateway-zone-b": {GatewayConfiguration: &GatewayConfiguration{AuthGatewayID: "gateway-node-1"}, DONName: "gateway-zone-b"},
		},
	}
}
