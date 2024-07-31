package integrationtesthelpers

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestJobSpec(t *testing.T) {
	tests := []struct {
		name     string
		params   *LMJobSpecParams
		err      bool
		expected string
	}{
		{
			name:   "empty",
			params: &LMJobSpecParams{},
			err:    true,
		},
		{
			name: "lm job spec",
			params: &LMJobSpecParams{
				Name:                    "liquiditymanager",
				Type:                    "ping-pong",
				ChainID:                 1337,
				ContractID:              "0x1234567890abcdef",
				TransmitterID:           "0xabcdef1234567890",
				OCRKeyBundleID:          "0xabcdef1234567890",
				CfgTrackerInterval:      20 * time.Second,
				LiquidityManagerAddress: common.HexToAddress("0x444444"),
				NetworkSelector:         1,
			},
			expected: `
type                                   = "offchainreporting2"
name                                   = "liquiditymanager"
forwardingAllowed                      = false
maxTaskDuration                        = "30s" 
pluginType                             = "liquiditymanager" 
relay                                  = "evm"
schemaVersion                          = 1
contractID                             = "0x1234567890abcdef"
ocrKeyBundleID                         = "0xabcdef1234567890" 
transmitterID                          = "0xabcdef1234567890" 
contractConfigConfirmations            = 1

contractConfigTrackerPollInterval      = "20s"


[pluginConfig]
closePluginTimeoutSec = 10
liquidityManagerAddress = "0x0000000000000000000000000000000000444444"
liquidityManagerNetwork = "1"
[pluginConfig.rebalancerConfig]
type = "ping-pong"


[relayConfig]
chainID = 1337
fromBlock = 0

`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jobSpec, err := NewJobSpec(tc.params)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, jobSpec)
			jobSpecToml, err := jobSpec.String()
			require.NoError(t, err)
			t.Logf("job spec toml: \n%s\n", jobSpecToml)
			require.Equal(t, tc.expected, jobSpecToml)
		})
	}
}

func TestBootstrapJobSpec(t *testing.T) {
	tests := []struct {
		name     string
		params   *LMJobSpecParams
		err      bool
		expected string
	}{
		{
			name:   "empty",
			params: &LMJobSpecParams{},
			err:    true,
		},
		{
			name: "lm bootstrap job spec",
			params: &LMJobSpecParams{
				ChainID:            1337,
				ContractID:         "0x1234567890abcdef",
				CfgTrackerInterval: 10 * time.Second,
				RelayFromBlock:     1234,
			},
			expected: `
type                                   = "bootstrap"
name                                   = "bootstrap-1337-0x1234567890abcdef"
forwardingAllowed                      = false
relay                                  = "evm"
schemaVersion                          = 1
contractID                             = "0x1234567890abcdef"
contractConfigConfirmations            = 1

contractConfigTrackerPollInterval      = "10s"


[relayConfig]
chainID = 1337
fromBlock = 1234

`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jobSpec, err := NewBootsrapJobSpec(tc.params)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, jobSpec)
			jobSpecToml, err := jobSpec.String()
			require.NoError(t, err)
			t.Logf("bootstrap job spec toml: \n%s\n", jobSpecToml)
			require.Equal(t, tc.expected, jobSpecToml)
		})
	}
}
