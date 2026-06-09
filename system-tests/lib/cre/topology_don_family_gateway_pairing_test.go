package cre

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGatewayConnectorsForDonFamily_legacyWithoutFamilies(t *testing.T) {
	t.Parallel()

	topology := &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "workflow", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "bootstrap-gateway", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			},
		},
		GatewayConnectors: &GatewayConnectors{
			Configurations: []*DonGatewayConfiguration{
				{GatewayConfiguration: &GatewayConfiguration{AuthGatewayID: "gateway-node-0"}},
			},
		},
		gatewayConnectorsByDon: map[string]*DonGatewayConfiguration{
			"bootstrap-gateway": {GatewayConfiguration: &GatewayConfiguration{AuthGatewayID: "gateway-node-0"}},
		},
	}

	require.False(t, topology.DonFamilyGatewayPairingEnabled())
	require.NoError(t, topology.validateDonFamilyGatewayPairing())

	conn := topology.GatewayConnectorsForDonFamily("")
	require.Len(t, conn.Configurations, 1)
	require.Equal(t, "gateway-node-0", conn.Configurations[0].AuthGatewayID)
}

func TestGatewayConnectorsForDonFamily_pairing(t *testing.T) {
	t.Parallel()

	topology := donFamilyGatewayPairingTestTopology()

	connA := topology.GatewayConnectorsForDonFamily("feeds-zone-a")
	require.Len(t, connA.Configurations, 1)
	require.Equal(t, "gateway-node-0", connA.Configurations[0].AuthGatewayID)

	connB := topology.GatewayConnectorsForDonFamily("feeds-zone-b")
	require.Len(t, connB.Configurations, 1)
	require.Equal(t, "gateway-node-1", connB.Configurations[0].AuthGatewayID)
}

func TestGatewayConnectorsForDonFamily_unknownFamily(t *testing.T) {
	t.Parallel()

	topology := donFamilyGatewayPairingTestTopology()

	conn := topology.GatewayConnectorsForDonFamily("unknown-family")
	require.Empty(t, conn.Configurations)
}

func TestTopology_validateDonFamilyGatewayPairing_missingGateway(t *testing.T) {
	t.Parallel()

	topology := &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "feeds-zone-a", DonFamily: "feeds-zone-a", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "feeds-zone-b", DonFamily: "feeds-zone-b", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "feeds-zone-c", DonFamily: "feeds-zone-c", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "gateway-zone-a", DonFamily: "feeds-zone-a", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
				{Name: "gateway-zone-b", DonFamily: "feeds-zone-b", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			},
		},
	}

	err := topology.validateDonFamilyGatewayPairing()
	require.Error(t, err)
	require.Contains(t, err.Error(), "no gateway DON is defined for that family")
}

func TestGatewayDonFamilyPairings_dataFeedsLocalTopology(t *testing.T) {
	t.Parallel()

	topology := donFamilyGatewayPairingTestTopology()

	pairs := topology.GatewayDonFamilyPairings()
	require.Equal(t, []DonFamilyGatewayPair{
		{DonFamily: "feeds-zone-a", WorkflowDONName: "feeds-zone-a", GatewayDONName: "gateway-zone-a"},
		{DonFamily: "feeds-zone-b", WorkflowDONName: "feeds-zone-b", GatewayDONName: "gateway-zone-b"},
	}, pairs)
}

func TestGatewayServiceConfigsForDonFamily_pairing(t *testing.T) {
	t.Parallel()

	topology := donFamilyGatewayPairingTestTopology()
	services := []GatewayServiceConfig{{
		ServiceName: "workflows",
		Handlers:    []string{"web-api-capabilities"},
		DONs:        []string{"feeds-zone-a", "feeds-zone-b"},
	}}

	scopedA := topology.GatewayServiceConfigsForDonFamily("feeds-zone-a", services)
	require.Equal(t, []string{"feeds-zone-a"}, scopedA[0].DONs)

	scopedB := topology.GatewayServiceConfigsForDonFamily("feeds-zone-b", services)
	require.Equal(t, []string{"feeds-zone-b"}, scopedB[0].DONs)
}

func donFamilyGatewayPairingTestTopology() *Topology {
	return &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "feeds-zone-a", DonFamily: "feeds-zone-a", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "feeds-zone-b", DonFamily: "feeds-zone-b", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "gateway-zone-a", DonFamily: "feeds-zone-a", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
				{Name: "gateway-zone-b", DonFamily: "feeds-zone-b", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			},
		},
		GatewayConnectors: &GatewayConnectors{
			Configurations: []*DonGatewayConfiguration{
				{GatewayConfiguration: &GatewayConfiguration{AuthGatewayID: "gateway-node-0"}},
				{GatewayConfiguration: &GatewayConfiguration{AuthGatewayID: "gateway-node-1"}},
			},
		},
		gatewayConnectorsByDon: map[string]*DonGatewayConfiguration{
			"gateway-zone-a": {GatewayConfiguration: &GatewayConfiguration{AuthGatewayID: "gateway-node-0"}},
			"gateway-zone-b": {GatewayConfiguration: &GatewayConfiguration{AuthGatewayID: "gateway-node-1"}},
		},
	}
}
