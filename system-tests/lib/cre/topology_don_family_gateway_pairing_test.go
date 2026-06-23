package cre

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGatewayConnectorsForDonFamily_requiresDonFamily(t *testing.T) {
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

	err := topology.initDonFamilyGatewayPairing()
	require.Error(t, err)
	require.Contains(t, err.Error(), "has no don_family")
}

func TestGatewayConnectorsForDonFamily_pairing(t *testing.T) {
	t.Parallel()

	topology := NewDonFamilyGatewayPairingTestTopology()

	connA := topology.GatewayConnectorsForDonFamily("feeds-zone-a")
	require.Len(t, connA.Configurations, 1)
	require.Equal(t, "gateway-node-0", connA.Configurations[0].AuthGatewayID)

	connB := topology.GatewayConnectorsForDonFamily("feeds-zone-b")
	require.Len(t, connB.Configurations, 1)
	require.Equal(t, "gateway-node-1", connB.Configurations[0].AuthGatewayID)
}

func TestGatewayConnectorsForDonFamily_unknownFamily(t *testing.T) {
	t.Parallel()

	topology := NewDonFamilyGatewayPairingTestTopology()

	conn := topology.GatewayConnectorsForDonFamily("unknown-family")
	require.Empty(t, conn.Configurations)
}

func TestGatewayConnectorsForDonFamily_emptyFamily(t *testing.T) {
	t.Parallel()

	topology := NewDonFamilyGatewayPairingTestTopology()

	conn := topology.GatewayConnectorsForDonFamily("")
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

	err := topology.initDonFamilyGatewayPairing()
	require.Error(t, err)
	require.Contains(t, err.Error(), "no gateway DON is defined for that family")
}

func TestTopology_validateDonFamilyGatewayPairing_gatewayMissingDonFamily(t *testing.T) {
	t.Parallel()

	topology := &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "workflow", DonFamily: testDefaultDONFamily, ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "bootstrap-gateway", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			},
		},
	}

	err := topology.initDonFamilyGatewayPairing()
	require.Error(t, err)
	require.Contains(t, err.Error(), `gateway DON "bootstrap-gateway" has no don_family`)
}

func TestDonFamilyGatewayPairingSummary_dataFeedsLocalTopology(t *testing.T) {
	t.Parallel()

	topology := NewDonFamilyGatewayPairingTestTopology()

	summary := topology.DonFamilyGatewayPairingSummary()
	require.Contains(t, summary, "feeds-zone-a → gateway-zone-a (don_family=feeds-zone-a)")
	require.Contains(t, summary, "feeds-zone-b → gateway-zone-b (don_family=feeds-zone-b)")
}

func TestDonFamilyGatewayPairingSummary_requiresDonFamily(t *testing.T) {
	t.Parallel()

	topology := &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "workflow", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "bootstrap-gateway", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			},
		},
	}

	err := topology.initDonFamilyGatewayPairing()
	require.Error(t, err)
	require.Contains(t, err.Error(), "has no don_family")
}

func TestTopology_validateDonFamilyGatewayPairing_partialFamily(t *testing.T) {
	t.Parallel()

	topology := &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "feeds-zone-a", DonFamily: "feeds-zone-a", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "feeds-zone-b", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "gateway-zone-a", DonFamily: "feeds-zone-a", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			},
		},
	}

	err := topology.initDonFamilyGatewayPairing()
	require.Error(t, err)
	require.Contains(t, err.Error(), `workflow DON "feeds-zone-b" has no don_family`)
}

func TestDonFamilyGatewayPairings_dataFeedsLocalTopology(t *testing.T) {
	t.Parallel()

	topology := NewDonFamilyGatewayPairingTestTopology()

	pairs := topology.DonFamilyGatewayPairings()
	require.Equal(t, []DonFamilyGatewayPair{
		{DonFamily: "feeds-zone-a", WorkflowDONName: "feeds-zone-a", GatewayDONName: "gateway-zone-a"},
		{DonFamily: "feeds-zone-b", WorkflowDONName: "feeds-zone-b", GatewayDONName: "gateway-zone-b"},
	}, pairs)
}

func TestDonFamilyGatewayPairings_multiGatewaySameFamily(t *testing.T) {
	t.Parallel()

	topology := NewMultiGatewaySameFamilyTestTopology()

	pairs := topology.DonFamilyGatewayPairings()
	require.Len(t, pairs, 2)
	require.Equal(t, DonFamilyGatewayPair{
		DonFamily: testDefaultDONFamily, WorkflowDONName: "workflow", GatewayDONName: "bootstrap-gateway-us",
	}, pairs[0])
	require.Equal(t, DonFamilyGatewayPair{
		DonFamily: testDefaultDONFamily, WorkflowDONName: "workflow", GatewayDONName: "gateway-eu",
	}, pairs[1])

	conn := topology.GatewayConnectorsForDonFamily(testDefaultDONFamily)
	require.Len(t, conn.Configurations, 2)
}

func TestGatewayServiceConfigsForDonFamily_pairing(t *testing.T) {
	t.Parallel()

	topology := NewDonFamilyGatewayPairingTestTopology()
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

func TestWorkflowDONFamilies(t *testing.T) {
	t.Parallel()

	topology := NewDonFamilyGatewayPairingTestTopology()
	require.ElementsMatch(t, []string{"feeds-zone-a", "feeds-zone-b"}, topology.WorkflowDONFamilies())
}

func TestWorkflowDONFamilies_nilWithoutPairing(t *testing.T) {
	t.Parallel()

	topology := &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "workflow", Flags: []string{WorkflowDON}},
			},
		},
	}

	require.Nil(t, topology.WorkflowDONFamilies())
}

func TestDonFamilyForDON(t *testing.T) {
	t.Parallel()

	topology := NewDonFamilyGatewayPairingTestTopology()

	require.Equal(t, "feeds-zone-a", topology.DonFamilyForDON("feeds-zone-a"))
	require.Equal(t, "feeds-zone-b", topology.DonFamilyForDON("gateway-zone-b"))
	require.Equal(t, "", topology.DonFamilyForDON("unknown"))
}

func TestInitDonFamilyGatewayPairing_skipsWithoutGateway(t *testing.T) {
	t.Parallel()

	topology := &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "workflow", Flags: []string{WorkflowDON}},
				{Name: "capabilities", Flags: []string{CapabilitiesDON}},
			},
		},
	}

	err := topology.initDonFamilyGatewayPairing()
	require.NoError(t, err)
	require.Nil(t, topology.donFamilyPairing)
	require.Nil(t, topology.WorkflowDONFamilies())
}
