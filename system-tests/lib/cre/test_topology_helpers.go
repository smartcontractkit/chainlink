package cre

// NewDonFamilyGatewayPairingTestTopology builds an in-memory two-family workflow + gateway topology
// (feeds-zone-a/b, gateway-zone-a/b) with initDonFamilyGatewayPairing already applied.
//
// Used by topology_don_family_gateway_pairing_test.go and cross-package deploy resolver tests
// without spinning up Docker. Not a substitute for chainlink-data-feeds local CRE E2E.
func NewDonFamilyGatewayPairingTestTopology() *Topology {
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
	if err := topology.initDonFamilyGatewayPairing(); err != nil {
		panic(err)
	}
	return topology
}
