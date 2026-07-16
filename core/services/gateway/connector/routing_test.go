package connector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGatewayConnector_GatewayIDsForDon(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		gateways []ConnectorGatewayConfig
		donID    string
		want     []string
		wantErr  string
	}{
		{
			name: "legacy empty donID returns all gateways",
			gateways: []ConnectorGatewayConfig{
				{ID: "gateway_a", URL: "ws://localhost:8081/a"},
				{ID: "gateway_b", URL: "ws://localhost:8081/b"},
			},
			donID: "",
			want:  []string{"gateway_a", "gateway_b"},
		},
		{
			name: "legacy non-empty donID returns no matches",
			gateways: []ConnectorGatewayConfig{
				{ID: "gateway_a", URL: "ws://localhost:8081/a"},
			},
			donID: "gateway_don_us",
			want:  nil,
		},
		{
			name: "multi-DON filters by gateway DON",
			gateways: []ConnectorGatewayConfig{
				{ID: "gateway_us_1", DonID: "gateway_don_us", URL: "ws://localhost:8081/us-1"},
				{ID: "gateway_us_2", DonID: "gateway_don_us", URL: "ws://localhost:8081/us-2"},
				{ID: "gateway_eu_1", DonID: "gateway_don_eu", URL: "ws://localhost:8081/eu-1"},
			},
			donID: "gateway_don_us",
			want:  []string{"gateway_us_1", "gateway_us_2"},
		},
		{
			name: "multi-DON empty donID returns all gateways",
			gateways: []ConnectorGatewayConfig{
				{ID: "gateway_us_1", DonID: "gateway_don_us", URL: "ws://localhost:8081/us-1"},
				{ID: "gateway_eu_1", DonID: "gateway_don_eu", URL: "ws://localhost:8081/eu-1"},
			},
			donID: "",
			want:  []string{"gateway_us_1", "gateway_eu_1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			connector, _, _ := newTestConnector(t, &ConnectorConfig{
				NodeAddress: "0x68902d681c28119f9b2531473a417088bf008e59",
				DonId:       "example_don",
				Gateways:    tt.gateways,
			})

			got, err := connector.GatewayIDsForDon(t.Context(), tt.donID)
			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestGatewayConnector_DonIDForGateway(t *testing.T) {
	t.Parallel()

	t.Run("legacy gateway", func(t *testing.T) {
		t.Parallel()

		connector, _, _ := newTestConnector(t, &ConnectorConfig{
			NodeAddress: "0x68902d681c28119f9b2531473a417088bf008e59",
			DonId:       "example_don",
			Gateways: []ConnectorGatewayConfig{
				{ID: "gateway_legacy", URL: "ws://localhost:8081/legacy"},
			},
		})

		donID, err := connector.DonIDForGateway(t.Context(), "gateway_legacy")
		require.NoError(t, err)
		require.Empty(t, donID)
	})

	t.Run("multi-DON gateway", func(t *testing.T) {
		t.Parallel()

		connector, _, _ := newTestConnector(t, &ConnectorConfig{
			NodeAddress: "0x68902d681c28119f9b2531473a417088bf008e59",
			DonId:       "example_don",
			Gateways: []ConnectorGatewayConfig{
				{ID: "gateway_us", DonID: "gateway_don_us", URL: "ws://localhost:8081/us"},
			},
		})

		donID, err := connector.DonIDForGateway(t.Context(), "gateway_us")
		require.NoError(t, err)
		require.Equal(t, "gateway_don_us", donID)
	})

	t.Run("invalid gateway ID", func(t *testing.T) {
		t.Parallel()

		connector, _, _ := newTestConnector(t, parseTOMLConfig(t, defaultConfig))

		_, err := connector.DonIDForGateway(t.Context(), "missing")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid Gateway ID")
	})
}

func TestGatewayConnector_PrimaryDonIDNotImplemented(t *testing.T) {
	t.Parallel()

	connector, _, _ := newTestConnector(t, &ConnectorConfig{
		NodeAddress: "0x68902d681c28119f9b2531473a417088bf008e59",
		DonId:       "example_don",
		Gateways: []ConnectorGatewayConfig{
			{ID: "gateway_us", DonID: "gateway_don_us", URL: "ws://localhost:8081/us"},
		},
	})

	_, err := connector.PrimaryDonID(t.Context())
	require.Error(t, err)
	require.Contains(t, err.Error(), "not implemented")
}

func TestGatewayConnector_DonIDReturnsSourceDON(t *testing.T) {
	t.Parallel()

	connector, _, _ := newTestConnector(t, &ConnectorConfig{
		NodeAddress: "0x68902d681c28119f9b2531473a417088bf008e59",
		DonId:       "workflow_don",
		Gateways: []ConnectorGatewayConfig{
			{ID: "gateway_us", DonID: "gateway_don_us", URL: "ws://localhost:8081/us"},
		},
	})

	donID, err := connector.DonID(t.Context())
	require.NoError(t, err)
	require.Equal(t, "workflow_don", donID)
}
