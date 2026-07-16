package cre

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func TestGatewayConfiguration_ExternalHTTPURL(t *testing.T) {
	t.Parallel()

	cfg := &GatewayConfiguration{
		Incoming: Incoming{
			Protocol:     "http",
			Path:         "/",
			ExternalPort: 5002,
		},
	}

	require.Equal(t, "http://localhost:5002/", cfg.ExternalHTTPURL(infra.Provider{Type: infra.Docker}))
	require.Equal(t, "http://cre-ns-gateway.example:5002/", cfg.ExternalHTTPURL(infra.Provider{
		Type: infra.Kubernetes,
		Kubernetes: &infra.KubernetesInput{
			Namespace:      "cre-ns",
			ExternalDomain: "example",
		},
	}))

	cfg.Incoming.Host = "explicit-host"
	require.Equal(t, "http://explicit-host:5002/", cfg.ExternalHTTPURL(infra.Provider{Type: infra.Docker}))
}

func TestGatewayConfiguration_ToConnectorGateway(t *testing.T) {
	t.Parallel()

	t.Run("includes DonID when gateway DON is set", func(t *testing.T) {
		t.Parallel()

		cfg := &GatewayConfiguration{
			AuthGatewayID: "gateway-node-1",
			GatewayDonID:  "gateway_don_eu",
			Outgoing: Outgoing{
				Host: "gateway-eu-node0",
				Port: 5003,
				Path: "/node",
			},
		}

		got := cfg.ToConnectorGateway()
		require.NotNil(t, got.ID)
		require.Equal(t, "gateway-node-1", *got.ID)
		require.NotNil(t, got.DonID)
		require.Equal(t, "gateway_don_eu", *got.DonID)
		require.NotNil(t, got.URL)
		require.Equal(t, "ws://gateway-eu-node0:5003/node", *got.URL)
	})

	t.Run("omits DonID for legacy single-gateway configs", func(t *testing.T) {
		t.Parallel()

		cfg := &GatewayConfiguration{
			AuthGatewayID: "gateway-node-0",
			Outgoing: Outgoing{
				Host: "bootstrap-gateway-us-node0",
				Port: 5003,
				Path: "/node",
			},
		}

		got := cfg.ToConnectorGateway()
		require.Nil(t, got.DonID)
	})
}
