package generic

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/stretchr/testify/assert"
)

func TestNewOracleFactoryConfig(t *testing.T) {
	tests := []struct {
		name          string
		config        string
		expectedError error
		expected      *oracleFactoryConfig
	}{
		{
			name:   "valid config",
			config: `{"enabled": true, "traceLogging": true, "bootstrapPeers": [ "12D3KooWDSxjWrKDscvUx7xEebZ6Dez1W1nfcS2jjCRuFVoujpMh@localhost:8001" ]}`,
			expected: &oracleFactoryConfig{
				Enabled:      true,
				TraceLogging: true,
				BootstrapPeers: []commontypes.BootstrapperLocator{
					{
						PeerID: "12D3KooWDSxjWrKDscvUx7xEebZ6Dez1W1nfcS2jjCRuFVoujpMh",
						Addrs:  []string{"localhost:8001"},
					},
				},
			},
		},
		{
			name:          "invalid JSON",
			config:        `{"enabled": true, "traceLogging": true, "bootstrapPeers": [ "12D3KooWDSxjWrKDscvUx7xEebZ6Dez1W1nfcS2jjCRuFVoujpMh@localhost:8001" `,
			expectedError: errors.New("failed to unmarshal oracle factory config"),
		},
		{
			name:   "disabled config",
			config: `{"enabled": false, "traceLogging": true, "bootstrapPeers": [ "12D3KooWDSxjWrKDscvUx7xEebZ6Dez1W1nfcS2jjCRuFVoujpMh@localhost:8001" ]}`,
			expected: &oracleFactoryConfig{
				Enabled: false,
			},
		},
		{
			name:          "enabled with no bootstrap peers",
			config:        `{"enabled": true, "traceLogging": true, "bootstrapPeers": []}`,
			expectedError: errors.New("no bootstrap peers found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := NewOracleFactoryConfig(tt.config)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, config)
			}
		})
	}
}
