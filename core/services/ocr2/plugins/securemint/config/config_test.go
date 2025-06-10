package config

import (
	"testing"

	sm_plugin "github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

func TestValidatePluginConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *sm_plugin.PorOffchainConfig
		wantErr bool
	}{
		{
			name: "valid config with MaxChains=2",
			cfg: &sm_plugin.PorOffchainConfig{
				MaxChains: 2,
			},
			wantErr: false,
		},
		{
			name: "valid config with MaxChains=0",
			cfg: &sm_plugin.PorOffchainConfig{
				MaxChains: 0,
			},
			wantErr: true,
		},
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSecureMintConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSecureMintConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
