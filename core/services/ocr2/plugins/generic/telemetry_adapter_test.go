package generic_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/generic"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

type mockEndpoint struct {
	network       string
	chainID       string
	contractID    string
	telemetryType string
	payload       []byte
}

func (m *mockEndpoint) SendLog(payload []byte) { m.payload = payload }

type mockGenerator struct{}

func (m *mockGenerator) GenMonitoringEndpoint(network string, chainID string, contractID string, telemetryType synchronization.TelemetryType) commontypes.MonitoringEndpoint {
	return &mockEndpoint{
		network:       network,
		chainID:       chainID,
		contractID:    contractID,
		telemetryType: string(telemetryType),
	}
}

func (m *mockGenerator) GenMultitypeMonitoringEndpoint(network string, chainID string, contractID string) telemetry.MultitypeMonitoringEndpoint {
	return nil
}

func TestTelemetryAdapter(t *testing.T) {
	ta := generic.NewTelemetryAdapter(&mockGenerator{})

	tests := []struct {
		name          string
		contractID    string
		telemetryType string
		networkID     string
		chainID       string
		payload       []byte
		errorMsg      string
	}{
		{
			name:          "valid request",
			contractID:    "contract",
			telemetryType: "mercury",
			networkID:     "solana",
			chainID:       "1337",
			payload:       []byte("uh oh"),
		},
		{
			name:          "no valid contractID",
			telemetryType: "mercury",
			networkID:     "solana",
			chainID:       "1337",
			payload:       []byte("uh oh"),
			errorMsg:      "contractID cannot be empty",
		},
		{
			name:          "no valid chainID",
			contractID:    "contract",
			telemetryType: "mercury",
			networkID:     "solana",
			payload:       []byte("uh oh"),
			errorMsg:      "chainID cannot be empty",
		},
		{
			name:       "no valid telemetryType",
			contractID: "contract",
			networkID:  "solana",
			chainID:    "1337",
			payload:    []byte("uh oh"),
			errorMsg:   "telemetryType cannot be empty",
		},
		{
			name:          "no valid network",
			contractID:    "contract",
			telemetryType: "mercury",
			chainID:       "1337",
			payload:       []byte("uh oh"),
			errorMsg:      "network cannot be empty",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ta.Send(testutils.Context(t), test.networkID, test.chainID, test.contractID, test.telemetryType, test.payload)
			if test.errorMsg != "" {
				assert.ErrorContains(t, err, test.errorMsg)
			} else {
				require.NoError(t, err)
				key := [4]string{test.networkID, test.chainID, test.contractID, test.telemetryType}
				endpoint, ok := ta.Endpoints()[key]
				require.True(t, ok)

				me := endpoint.(*mockEndpoint)
				assert.Equal(t, test.networkID, me.network)
				assert.Equal(t, test.chainID, me.chainID)
				assert.Equal(t, test.contractID, me.contractID)
				assert.Equal(t, test.telemetryType, me.telemetryType)
				assert.Equal(t, test.payload, me.payload)
			}
		})
	}
}
