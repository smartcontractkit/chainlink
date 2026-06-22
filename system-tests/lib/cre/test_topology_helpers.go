package cre

// NewDonFamilyGatewayPairingTestTopology returns a two-zone workflow + gateway topology
// with don_family pairing enabled. Intended for unit tests in other packages.
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
