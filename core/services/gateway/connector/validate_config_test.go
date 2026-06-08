package connector

import (
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	gatewaymocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
)

func validConnectorConfig() *ConnectorConfig {
	return &ConnectorConfig{
		NodeAddress: "0x68902d681c28119f9b2531473a417088bf008e59",
		DonId:       "example_don",
		Gateways: []ConnectorGatewayConfig{
			{ID: "gateway_a", URL: "ws://localhost:8081/a"},
			{ID: "gateway_b", URL: "ws://localhost:8081/b"},
		},
	}
}

func TestValidateConnectorConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		modify     func(*ConnectorConfig)
		wantErr    bool
		errMessage string
	}{
		{
			name:    "legacy single gateway without per-gateway DonID",
			modify:  func(cfg *ConnectorConfig) { cfg.Gateways = cfg.Gateways[:1] },
			wantErr: false,
		},
		{
			name:    "legacy multiple gateways without per-gateway DonID",
			wantErr: false,
		},
		{
			name: "legacy ignores unset per-gateway DonID on all gateways",
			modify: func(cfg *ConnectorConfig) {
				for i := range cfg.Gateways {
					cfg.Gateways[i].DonID = ""
				}
			},
			wantErr: false,
		},
		{
			name: "multi-DON single gateway",
			modify: func(cfg *ConnectorConfig) {
				cfg.Gateways = []ConnectorGatewayConfig{
					{ID: "gateway_us", DonID: "gateway_don_us", URL: "ws://localhost:8081/us"},
				}
			},
			wantErr: false,
		},
		{
			name: "multi-DON all gateways set DonID",
			modify: func(cfg *ConnectorConfig) {
				cfg.Gateways = []ConnectorGatewayConfig{
					{ID: "gateway_us_1", DonID: "gateway_don_us", URL: "ws://localhost:8081/us-1"},
					{ID: "gateway_us_2", DonID: "gateway_don_us", URL: "ws://localhost:8081/us-2"},
					{ID: "gateway_eu_1", DonID: "gateway_don_eu", URL: "ws://localhost:8081/eu-1"},
				}
			},
			wantErr: false,
		},
		{
			name: "empty connector DonID",
			modify: func(cfg *ConnectorConfig) {
				cfg.DonId = ""
			},
			wantErr:    true,
			errMessage: "invalid DON ID",
		},
		{
			name: "connector DonID exceeds max length",
			modify: func(cfg *ConnectorConfig) {
				cfg.DonId = "012345678901234567890123456789012345678901234567890123456789012345"
			},
			wantErr:    true,
			errMessage: "invalid DON ID",
		},
		{
			name: "partial per-gateway DonID",
			modify: func(cfg *ConnectorConfig) {
				cfg.Gateways = []ConnectorGatewayConfig{
					{ID: "gateway_us", DonID: "gateway_don_us", URL: "ws://localhost:8081/us"},
					{ID: "gateway_eu", URL: "ws://localhost:8081/eu"},
				}
			},
			wantErr:    true,
			errMessage: "all gateways must set DonID when multi-DON mode is enabled",
		},
		{
			name: "partial per-gateway DonID with legacy gateway first",
			modify: func(cfg *ConnectorConfig) {
				cfg.Gateways = []ConnectorGatewayConfig{
					{ID: "gateway_legacy", URL: "ws://localhost:8081/legacy"},
					{ID: "gateway_us", DonID: "gateway_don_us", URL: "ws://localhost:8081/us"},
				}
			},
			wantErr:    true,
			errMessage: "all gateways must set DonID when multi-DON mode is enabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := validConnectorConfig()
			if tt.modify != nil {
				tt.modify(cfg)
			}

			err := validateConnectorConfig(cfg)
			if tt.wantErr {
				require.Error(t, err)
				require.EqualError(t, err, tt.errMessage)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestNewGatewayConnector_ConfigValidation(t *testing.T) {
	t.Parallel()

	signer := gatewaymocks.NewSigner(t)
	clock := clockwork.NewFakeClock()
	lggr := logger.Test(t)

	t.Run("legacy config passes", func(t *testing.T) {
		t.Parallel()

		_, err := NewGatewayConnector(validConnectorConfig(), signer, clock, lggr)
		require.NoError(t, err)
	})

	t.Run("multi-DON config passes", func(t *testing.T) {
		t.Parallel()

		cfg := validConnectorConfig()
		cfg.Gateways = []ConnectorGatewayConfig{
			{ID: "gateway_us_1", DonID: "gateway_don_us", URL: "ws://localhost:8081/us-1"},
			{ID: "gateway_eu_1", DonID: "gateway_don_eu", URL: "ws://localhost:8081/eu-1"},
		}

		_, err := NewGatewayConnector(cfg, signer, clock, lggr)
		require.NoError(t, err)
	})

	t.Run("invalid config fails before connector construction", func(t *testing.T) {
		t.Parallel()

		cfg := validConnectorConfig()
		cfg.Gateways = []ConnectorGatewayConfig{
			{ID: "gateway_us", DonID: "gateway_don_us", URL: "ws://localhost:8081/us"},
			{ID: "gateway_eu", URL: "ws://localhost:8081/eu"},
		}

		_, err := NewGatewayConnector(cfg, signer, clock, lggr)
		require.EqualError(t, err, "all gateways must set DonID when multi-DON mode is enabled")
	})
}
