package cre

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDonMetadata_DonFamilies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		don      DonMetadata
		expected []string
	}{
		{name: "primary only", don: DonMetadata{DonFamily: "zone-a"}, expected: []string{"zone-a"}},
		{
			name:     "primary plus shard family",
			don:      DonMetadata{DonFamily: "zone-a", AdditionalDonFamilies: []string{"zone-a_shard-0"}},
			expected: []string{"zone-a", "zone-a_shard-0"},
		},
		{
			name:     "dedups repeated families",
			don:      DonMetadata{DonFamily: "zone-a", AdditionalDonFamilies: []string{"zone-a", "zone-a_shard-1", "zone-a_shard-1"}},
			expected: []string{"zone-a", "zone-a_shard-1"},
		},
		{
			name:     "trims whitespace and drops empties",
			don:      DonMetadata{DonFamily: " zone-a ", AdditionalDonFamilies: []string{"  ", " zone-a_shard-2 "}},
			expected: []string{"zone-a", "zone-a_shard-2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.expected, tt.don.DonFamilies())
		})
	}
}

func TestShardedCapabilityTopology_FamilyLayout(t *testing.T) {
	t.Parallel()

	topology := shardedCapabilityTestTopology(t)

	byName := map[string]*DonMetadata{}
	for _, d := range topology.DonsMetadata.List() {
		byName[d.Name] = d
	}

	require.Equal(t, []string{"zone-a", "zone-a_shard-0"}, byName["workflow-shard-0"].DonFamilies())
	require.Equal(t, []string{"zone-a", "zone-a_shard-1"}, byName["workflow-shard-1"].DonFamilies())
	require.Equal(t, []string{"zone-a_shard-0"}, byName["cap-shard-0"].DonFamilies())
	require.Equal(t, []string{"zone-a_shard-1"}, byName["cap-shard-1"].DonFamilies())
	require.Equal(t, []string{"zone-a"}, byName["shared-cap"].DonFamilies())

	require.Equal(t, []string{"zone-a"}, topology.WorkflowDONFamilies())
	require.ElementsMatch(t, []DonFamilyGatewayPair{
		{DonFamily: "zone-a", WorkflowDONName: "workflow-shard-0", GatewayDONName: "bootstrap-gateway"},
		{DonFamily: "zone-a", WorkflowDONName: "workflow-shard-1", GatewayDONName: "bootstrap-gateway"},
	}, topology.DonFamilyGatewayPairings())

	require.Len(t, topology.GatewayConnectorsForDonFamily("zone-a").Configurations, 1)
	require.Empty(t, topology.GatewayConnectorsForDonFamily("zone-a_shard-0").Configurations)
	require.Empty(t, topology.GatewayConnectorsForDonFamily("zone-a_shard-1").Configurations)
}

func shardedCapabilityTestTopology(t *testing.T) *Topology {
	t.Helper()

	conn := testGatewayConnector("gateway-node-0", "bootstrap-gateway.local", 5002)

	topology := &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "workflow-shard-0", ID: 1, DonFamily: "zone-a", AdditionalDonFamilies: []string{"zone-a_shard-0"}, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "workflow-shard-1", ID: 2, DonFamily: "zone-a", AdditionalDonFamilies: []string{"zone-a_shard-1"}, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "cap-shard-0", ID: 3, DonFamily: "zone-a_shard-0", Flags: []string{CapabilitiesDON}},
				{Name: "cap-shard-1", ID: 4, DonFamily: "zone-a_shard-1", Flags: []string{CapabilitiesDON}},
				{Name: "shared-cap", ID: 5, DonFamily: "zone-a", Flags: []string{CapabilitiesDON}},
				{Name: "bootstrap-gateway", DonFamily: "zone-a", NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			},
		},
		GatewayConnectors: &GatewayConnectors{
			Configurations: []*DonGatewayConfiguration{conn},
		},
		gatewayConnectorsByDon: map[string]*DonGatewayConfiguration{
			"bootstrap-gateway": conn,
		},
	}
	require.NoError(t, topology.initDonFamilyGatewayPairing())
	return topology
}
