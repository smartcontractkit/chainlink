// Tests for don_family gateway↔workflow pairing at env start.
//
// Covers initDonFamilyGatewayPairing validation, per-family gateway connector lookup,
// gateway worker DON scoping, and pairing summaries. Uses in-memory topologies from
// test_topology_helpers.go — not Docker E2E.
package cre

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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

// Unknown and empty don_family values return empty connector lists without error.
func TestGatewayConnectorsForDonFamily_unknownAndEmptyFamily(t *testing.T) {
	t.Parallel()

	topology := NewDonFamilyGatewayPairingTestTopology()

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
				{Name: "workflow", DonFamily: testDefaultDONFamily, ns: &NodeSet{DONTypes: []string{WorkflowDON}}, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "bootstrap-gateway", ns: &NodeSet{}, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			},
		},
	}

	err := topology.initDonFamilyGatewayPairing()
	require.Error(t, err)
	require.Contains(t, err.Error(), `gateway DON "bootstrap-gateway" has no don_family`)
}

// Env start fails when a gateway topology has workflow or gateway nodesets without don_family.
func TestInitDonFamilyGatewayPairing_requiresDonFamilyOnWorkflow(t *testing.T) {
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

func TestDonFamilyGatewayPairingSummary_dataFeedsLocalTopology(t *testing.T) {
	t.Parallel()

	topology := NewDonFamilyGatewayPairingTestTopology()

	summary := topology.DonFamilyGatewayPairingSummary()
	require.Contains(t, summary, "feeds-zone-a → gateway-zone-a (don_family=feeds-zone-a)")
	require.Contains(t, summary, "feeds-zone-b → gateway-zone-b (don_family=feeds-zone-b)")
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

// One workflow DON may pair with multiple gateway nodesets that share the same don_family.
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

// Non-gateway topologies skip pairing; WorkflowDONFamilies stays nil for registry setup fallback.
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

// Topologies without http-action skip pairing validation entirely.
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
