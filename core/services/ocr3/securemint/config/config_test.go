package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Validate(t *testing.T) {
	tests := []struct {
		name string
		cfg  *SecureMintConfig
		err  bool
	}{
		{
			name: "valid config",
			cfg: &SecureMintConfig{
				Token:    "eth",
				Reserves: "platform",
			},
			err: false,
		},
		{
			name: "nil config",
			cfg:  nil,
			err:  true,
		},
		{
			name: "invalid token",
			cfg: &SecureMintConfig{
				Token:    "",
				Reserves: "platform",
			},
			err: true,
		},
		{
			name: "invalid reserves",
			cfg: &SecureMintConfig{
				Token:    "eth",
				Reserves: "",
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParseSecureMintConfig(t *testing.T) {
	tests := []struct {
		name                   string
		configJSON             string
		expectedToken          string
		expectedReserves       string
		expectedChainSelectors []string
		expectError            bool
	}{
		{
			name:        "empty config is invalid",
			configJSON:  "",
			expectError: true,
		},
		{
			name:                   "custom values",
			configJSON:             `{"token": "btc", "reserves": "custom", "chainSelectors": ["8953668971247136127", "729797994450396300"]}`,
			expectedToken:          "btc",
			expectedReserves:       "custom",
			expectedChainSelectors: []string{"8953668971247136127", "729797994450396300"},
			expectError:            false,
		},
		{
			name:                   "partial config uses empty string",
			configJSON:             `{"token": "link", "chainSelectors": ["8953668971247136127", "729797994450396300"]}`,
			expectedToken:          "link",
			expectedReserves:       "",
			expectedChainSelectors: []string{"8953668971247136127", "729797994450396300"},
			expectError:            false,
		},
		{
			name:                   "partial config uses empty string 2",
			configJSON:             `{"reserves": "custom", "chainSelectors": ["8953668971247136127", "729797994450396300"]}`,
			expectedToken:          "",
			expectedReserves:       "custom",
			expectedChainSelectors: []string{"8953668971247136127", "729797994450396300"},
			expectError:            false,
		},
		{
			name:                   "partial config uses empty slice",
			configJSON:             `{"token": "btc", "reserves": "custom"}`,
			expectedToken:          "btc",
			expectedReserves:       "custom",
			expectedChainSelectors: nil,
			expectError:            false,
		},
		{
			name:             "invalid JSON",
			configJSON:       `{"token": "btc", "reserves":}`,
			expectedToken:    "",
			expectedReserves: "",
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := Parse([]byte(tt.configJSON))

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, config)
			require.Equal(t, tt.expectedToken, config.Token)
			require.Equal(t, tt.expectedReserves, config.Reserves)
			require.Equal(t, tt.expectedChainSelectors, config.ChainSelectors)
		})
	}
}
