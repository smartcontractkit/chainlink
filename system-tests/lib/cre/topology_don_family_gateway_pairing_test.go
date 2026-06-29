// Tests for don_family gateway↔workflow pairing at env start.
package cre

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const testDONFamily = "test-don-family"

func TestGatewayConnectorsForDonFamily_pairing(t *testing.T) {
	t.Parallel()

	topology := donFamilyGatewayPairingTestTopology(t)

	connA := topology.GatewayConnectorsForDonFamily("feeds-zone-a")
	require.Len(t, connA.Configurations, 1)
	require.Equal(t, "gateway-node-0", connA.Configurations[0].AuthGatewayID)

	connB := topology.GatewayConnectorsForDonFamily("feeds-zone-b")
	require.Len(t, connB.Configurations, 1)
	require.Equal(t, "gateway-node-1", connB.Configurations[0].AuthGatewayID)
}

func TestGatewayConnectorsForDonFamily_unknownAndEmptyFamily(t *testing.T) {
	t.Parallel()

	topology := donFamilyGatewayPairingTestTopology(t)

	require.Empty(t, topology.GatewayConnectorsForDonFamily("unknown-family").Configurations)
	require.Empty(t, topology.GatewayConnectorsForDonFamily("").Configurations)
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
				{Name: "workflow", DonFamily: testDONFamily, ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "bootstrap-gateway", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			},
		},
	}

	err := topology.initDonFamilyGatewayPairing()
	require.Error(t, err)
	require.Contains(t, err.Error(), `gateway DON "bootstrap-gateway" has no don_family`)
}

func TestInitDonFamilyGatewayPairing_requiresDonFamilyOnWorkflow(t *testing.T) {
	t.Parallel()

	topology := &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "workflow", ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "bootstrap-gateway", DonFamily: testDONFamily, ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			},
		},
	}

	err := topology.initDonFamilyGatewayPairing()
	require.Error(t, err)
	require.Contains(t, err.Error(), "has no don_family")
}

func TestDonFamilyGatewayPairingSummary_dataFeedsLocalTopology(t *testing.T) {
	t.Parallel()

	topology := donFamilyGatewayPairingTestTopology(t)

	summary := topology.DonFamilyGatewayPairingSummary()
	require.Contains(t, summary, "feeds-zone-a → gateway-zone-a (don_family=feeds-zone-a)")
	require.Contains(t, summary, "feeds-zone-b → gateway-zone-b (don_family=feeds-zone-b)")
}

func TestDonFamilyGatewayPairings_dataFeedsLocalTopology(t *testing.T) {
	t.Parallel()

	topology := donFamilyGatewayPairingTestTopology(t)

	pairs := topology.DonFamilyGatewayPairings()
	require.Equal(t, []DonFamilyGatewayPair{
		{DonFamily: "feeds-zone-a", WorkflowDONName: "feeds-zone-a", GatewayDONName: "gateway-zone-a"},
		{DonFamily: "feeds-zone-b", WorkflowDONName: "feeds-zone-b", GatewayDONName: "gateway-zone-b"},
	}, pairs)
}

func TestDonFamilyGatewayPairings_multiGatewaySameFamily(t *testing.T) {
	t.Parallel()

	topology := multiGatewaySameFamilyTestTopology(t)

	pairs := topology.DonFamilyGatewayPairings()
	require.Len(t, pairs, 2)
	require.Equal(t, DonFamilyGatewayPair{
		DonFamily: testDONFamily, WorkflowDONName: "workflow", GatewayDONName: "bootstrap-gateway-us",
	}, pairs[0])
	require.Equal(t, DonFamilyGatewayPair{
		DonFamily: testDONFamily, WorkflowDONName: "workflow", GatewayDONName: "gateway-eu",
	}, pairs[1])

	conn := topology.GatewayConnectorsForDonFamily(testDONFamily)
	require.Len(t, conn.Configurations, 2)
}

func TestGatewayServiceConfigsForGateway_pairing(t *testing.T) {
	t.Parallel()

	topology := donFamilyGatewayPairingTestTopology(t)
	services := []GatewayServiceConfig{{
		ServiceName: "workflows",
		Handlers:    []string{"web-api-capabilities"},
		DONs:        []string{"feeds-zone-a", "feeds-zone-b"},
	}}

	scopedA := topology.GatewayServiceConfigsForGateway("gateway-zone-a", services)
	require.Equal(t, []string{"feeds-zone-a"}, scopedA[0].DONs)

	scopedB := topology.GatewayServiceConfigsForGateway("gateway-zone-b", services)
	require.Equal(t, []string{"feeds-zone-b"}, scopedB[0].DONs)
}

func TestGatewayServiceConfigsForGateway_preservesCapabilitiesDONForVault(t *testing.T) {
	t.Parallel()

	topology := &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "workflow", DonFamily: testDONFamily, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "capabilities", DonFamily: testDONFamily, Flags: []string{CapabilitiesDON, VaultCapability}},
				{Name: "bootstrap-gateway", DonFamily: testDONFamily, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			},
		},
	}
	require.NoError(t, topology.initDonFamilyGatewayPairing())

	services := []GatewayServiceConfig{{
		ServiceName: "vault",
		Handlers:    []string{"vault"},
		DONs:        []string{"capabilities"},
	}}

	scoped := topology.GatewayServiceConfigsForGateway("bootstrap-gateway", services)
	require.Equal(t, []string{"capabilities"}, scoped[0].DONs)
}

func TestWorkflowDONFamilies(t *testing.T) {
	t.Parallel()

	topology := donFamilyGatewayPairingTestTopology(t)
	require.ElementsMatch(t, []string{"feeds-zone-a", "feeds-zone-b"}, topology.WorkflowDONFamilies())
}

func TestInitDonFamilyGatewayPairing_skipsWithoutGateway(t *testing.T) {
	t.Parallel()

	topology := &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "workflow", DonFamily: testDONFamily, Flags: []string{WorkflowDON}},
				{Name: "capabilities", DonFamily: testDONFamily, Flags: []string{CapabilitiesDON}},
			},
		},
	}

	err := topology.initDonFamilyGatewayPairing()
	require.NoError(t, err)
	require.Nil(t, topology.gatewayDonFamilyPairing)
	require.Equal(t, []string{testDONFamily}, topology.WorkflowDONFamilies())
}

func TestDonFamilyForDON(t *testing.T) {
	t.Parallel()

	topology := donFamilyGatewayPairingTestTopology(t)

	require.Equal(t, "feeds-zone-a", topology.DonFamilyForDON("feeds-zone-a"))
	require.Equal(t, "feeds-zone-b", topology.DonFamilyForDON("gateway-zone-b"))
	require.Empty(t, topology.DonFamilyForDON("unknown"))
}

func donFamilyGatewayPairingTestTopology(t *testing.T) *Topology {
	t.Helper()

	topology := &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "feeds-zone-a", ID: 1, DonFamily: "feeds-zone-a", Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "feeds-zone-b", ID: 2, DonFamily: "feeds-zone-b", Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "gateway-zone-a", DonFamily: "feeds-zone-a", NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
				{Name: "gateway-zone-b", DonFamily: "feeds-zone-b", NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
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
	require.NoError(t, topology.initDonFamilyGatewayPairing())
	return topology
}

func multiGatewaySameFamilyTestTopology(t *testing.T) *Topology {
	t.Helper()

	connUS := testGatewayConnector("gateway-node-0", "bootstrap-gateway-us.local", 5002)
	connEU := testGatewayConnector("gateway-node-1", "gateway-eu.local", 5004)

	topology := &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "workflow", ID: 1, DonFamily: testDONFamily, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "bootstrap-gateway-us", DonFamily: testDONFamily, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
				{Name: "gateway-eu", DonFamily: testDONFamily, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			},
		},
		GatewayConnectors: &GatewayConnectors{
			Configurations: []*DonGatewayConfiguration{connUS, connEU},
		},
		gatewayConnectorsByDon: map[string]*DonGatewayConfiguration{
			"bootstrap-gateway-us": connUS,
			"gateway-eu":           connEU,
		},
	}
	require.NoError(t, topology.initDonFamilyGatewayPairing())
	return topology
}

func testGatewayConnector(authID, host string, externalPort int) *DonGatewayConfiguration {
	return &DonGatewayConfiguration{
		GatewayConfiguration: &GatewayConfiguration{
			AuthGatewayID: authID,
			Incoming: Incoming{
				Protocol:     "http",
				Host:         host,
				Path:         "/",
				ExternalPort: externalPort,
			},
		},
	}
}
