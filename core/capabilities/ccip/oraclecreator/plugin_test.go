package oraclecreator

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

func TestPluginTypeToTelemetryType(t *testing.T) {
	tests := []struct {
		name       string
		pluginType types.PluginType
		expected   synchronization.TelemetryType
		expectErr  bool
	}{
		{
			name:       "CCIPCommit plugin type",
			pluginType: types.PluginTypeCCIPCommit,
			expected:   synchronization.OCR3CCIPCommit,
			expectErr:  false,
		},
		{
			name:       "CCIPExec plugin type",
			pluginType: types.PluginTypeCCIPExec,
			expected:   synchronization.OCR3CCIPExec,
			expectErr:  false,
		},
		{
			name:       "Unknown plugin type",
			pluginType: types.PluginType(99),
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pluginTypeToTelemetryType(tt.pluginType)
			if tt.expectErr {
				require.Error(t, err, "Expected an error for plugin type %d", tt.pluginType)
				require.Equal(t, synchronization.TelemetryType(""), result, "Expected empty result for plugin type %d", tt.pluginType)
			} else {
				require.NoError(t, err, "Unexpected error for plugin type %d", tt.pluginType)
				require.Equal(t, tt.expected, result, "Unexpected telemetry type for plugin type %d", tt.pluginType)
			}
		})
	}
}
