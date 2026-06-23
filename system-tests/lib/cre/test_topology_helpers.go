package cre

const testDefaultDONFamily = "test-don-family"

// testGatewayConnector builds a DonGatewayConfiguration with Incoming populated for URL tests.
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

// NewDonFamilyGatewayPairingTestTopologyWithIncoming is like NewDonFamilyGatewayPairingTestTopology
// but gateway connectors include Incoming host/port for gateway URL resolver tests.
func NewDonFamilyGatewayPairingTestTopologyWithIncoming() *Topology {
	connA := testGatewayConnector("gateway-node-0", "gateway-zone-a.local", 5002)
	connB := testGatewayConnector("gateway-node-1", "gateway-zone-b.local", 5004)

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
			Configurations: []*DonGatewayConfiguration{connA, connB},
		},
		gatewayConnectorsByDon: map[string]*DonGatewayConfiguration{
			"gateway-zone-a": connA,
			"gateway-zone-b": connB,
		},
	}
	if err := topology.initDonFamilyGatewayPairing(); err != nil {
		panic(err)
	}
	return topology
}

// NewMultiGatewaySameFamilyTestTopology builds one workflow DON and two gateway nodesets sharing
// don_family (mirrors workflow-gateway-capabilities-multi-gateway-don.toml shape).
func NewMultiGatewaySameFamilyTestTopology() *Topology {
	family := testDefaultDONFamily
	connUS := testGatewayConnector("gateway-node-0", "bootstrap-gateway-us.local", 5002)
	connEU := testGatewayConnector("gateway-node-1", "gateway-eu.local", 5004)

	topology := &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "workflow", ID: 1, DonFamily: family, Flags: []string{WorkflowDON, HTTPActionCapability}},
				{Name: "bootstrap-gateway-us", DonFamily: family, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
				{Name: "gateway-eu", DonFamily: family, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
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
	if err := topology.initDonFamilyGatewayPairing(); err != nil {
		panic(err)
	}
	return topology
}

// NewShardedWorkflowTestTopology builds shard0/shard1 workflow DONs with a shared gateway family
// and pairing initialized (for deploy resolver shard tests).
func NewShardedWorkflowTestTopology() *Topology {
	family := testDefaultDONFamily
	conn := testGatewayConnector("gateway-node-0", "bootstrap-gateway.local", 5002)

	topology := &Topology{
		DonsMetadata: &DonsMetadata{
			dons: []*DonMetadata{
				{Name: "shard0", ID: 1, DonFamily: family, ShardIndex: 0, Flags: []string{WorkflowDON, ShardDON, HTTPActionCapability}},
				{Name: "shard1", ID: 2, DonFamily: family, ShardIndex: 1, Flags: []string{WorkflowDON, ShardDON, HTTPActionCapability}},
				{Name: "bootstrap-gateway", DonFamily: family, NodesMetadata: []*NodeMetadata{{Roles: []string{GatewayNode}}}},
			},
		},
		GatewayConnectors: &GatewayConnectors{
			Configurations: []*DonGatewayConfiguration{conn},
		},
		gatewayConnectorsByDon: map[string]*DonGatewayConfiguration{
			"bootstrap-gateway": conn,
		},
	}
	if err := topology.initDonFamilyGatewayPairing(); err != nil {
		panic(err)
	}
	return topology
}
