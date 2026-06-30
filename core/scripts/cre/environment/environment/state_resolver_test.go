// Tests for LocalCREStateResolver gateway URL helpers used during workflow deploy.
package environment

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func TestGatewayURLForDonFamily_returnsFamilyConnector(t *testing.T) {
	t.Parallel()

	topology := stateResolverTestTopologyWithIncoming(t)
	resolver := &LocalCREStateResolver{topology: topology}

	url, err := resolver.GatewayURLForDonFamily("feeds-zone-a")
	require.NoError(t, err)
	require.Equal(t, "http://gateway-zone-a.local:5002/", url)
}

func TestGatewayURLForDonFamily_unknownFamily(t *testing.T) {
	t.Parallel()

	topology := stateResolverTestTopologyWithIncoming(t)
	resolver := &LocalCREStateResolver{topology: topology}

	_, err := resolver.GatewayURLForDonFamily("unknown-family")
	require.Error(t, err)
	require.Contains(t, err.Error(), `no gateway connector found for don_family "unknown-family"`)
}

func TestFormatGatewayURL_usesInfraHostWhenIncomingHostEmpty(t *testing.T) {
	t.Parallel()

	cfg := &cre.DonGatewayConfiguration{
		GatewayConfiguration: &cre.GatewayConfiguration{
			Incoming: cre.Incoming{
				Protocol:     "http",
				Path:         "/vault",
				ExternalPort: 5002,
			},
		},
	}
	resolver := &LocalCREStateResolver{
		cfg: &config.Config{
			Infra: &infra.Provider{Type: infra.Docker},
		},
	}

	url, err := resolver.formatGatewayURL(cfg)
	require.NoError(t, err)
	require.Equal(t, "http://localhost:5002/vault", url)
}

func stateResolverTestTopologyWithIncoming(t *testing.T) *cre.Topology {
	t.Helper()

	connA := deployTestGatewayConnector("gateway-node-0", "gateway-zone-a.local", 5002)
	connB := deployTestGatewayConnector("gateway-node-1", "gateway-zone-b.local", 5004)

	dm, err := cre.NewDonsMetadata([]*cre.DonMetadata{
		deployTestBootstrapDON(),
		{Name: "feeds-zone-a", ID: 1, DonFamily: "feeds-zone-a", Flags: []string{cre.WorkflowDON, cre.HTTPActionCapability}},
		{Name: "feeds-zone-b", ID: 2, DonFamily: "feeds-zone-b", Flags: []string{cre.WorkflowDON, cre.HTTPActionCapability}},
		{Name: "gateway-zone-a", DonFamily: "feeds-zone-a", NodesMetadata: []*cre.NodeMetadata{{Roles: []string{cre.GatewayNode}}}},
		{Name: "gateway-zone-b", DonFamily: "feeds-zone-b", NodesMetadata: []*cre.NodeMetadata{{Roles: []string{cre.GatewayNode}}}},
	}, infra.Provider{Type: infra.Docker})
	require.NoError(t, err)

	topology := &cre.Topology{
		DonsMetadata: dm,
		GatewayConnectors: &cre.GatewayConnectors{
			Configurations: []*cre.DonGatewayConfiguration{connA, connB},
		},
	}
	require.NoError(t, topology.EnsureGatewayDonFamilyPairing())
	return topology
}
