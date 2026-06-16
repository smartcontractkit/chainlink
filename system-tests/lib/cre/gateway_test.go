package cre

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
