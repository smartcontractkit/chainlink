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

	topology := cre.NewDonFamilyGatewayPairingTestTopologyWithIncoming()
	resolver := &LocalCREStateResolver{topology: topology}

	url, err := resolver.GatewayURLForDonFamily("feeds-zone-a")
	require.NoError(t, err)
	require.Equal(t, "http://gateway-zone-a.local:5002/", url)
}

func TestGatewayURLForDonFamily_unknownFamily(t *testing.T) {
	t.Parallel()

	topology := cre.NewDonFamilyGatewayPairingTestTopologyWithIncoming()
	resolver := &LocalCREStateResolver{topology: topology}

	_, err := resolver.GatewayURLForDonFamily("unknown-family")
	require.Error(t, err)
	require.Contains(t, err.Error(), `no gateway connector found for don_family "unknown-family"`)
}

func TestResolveGatewayURL_fallsBackToFirstGateway(t *testing.T) {
	t.Parallel()

	conn := &cre.DonGatewayConfiguration{
		GatewayConfiguration: &cre.GatewayConfiguration{
			AuthGatewayID: "fallback-gateway",
			Incoming: cre.Incoming{
				Protocol:     "http",
				Host:         "fallback.local",
				Path:         "/",
				ExternalPort: 5002,
			},
		},
	}
	topology := &cre.Topology{
		GatewayConnectors: &cre.GatewayConnectors{
			Configurations: []*cre.DonGatewayConfiguration{conn},
		},
	}
	resolver := &LocalCREStateResolver{topology: topology}

	url, err := resolver.resolveGatewayURL("unknown-family")
	require.NoError(t, err)
	require.Equal(t, "http://fallback.local:5002/", url)
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
