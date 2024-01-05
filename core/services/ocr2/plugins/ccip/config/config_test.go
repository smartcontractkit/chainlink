package config

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

func TestCommitConfig(t *testing.T) {
	exampleConfig := CommitPluginJobSpecConfig{
		SourceStartBlock:       222,
		DestStartBlock:         333,
		OffRamp:                common.HexToAddress("0x123"),
		TokenPricesUSDPipeline: `merge [type=merge left="{}" right="{\"0xC79b96044906550A5652BCf20a6EA02f139B9Ae5\":\"1000000000000000000\"}"];`,
	}

	bts, err := json.Marshal(exampleConfig)
	require.NoError(t, err)

	parsedConfig := CommitPluginJobSpecConfig{}
	require.NoError(t, json.Unmarshal(bts, &parsedConfig))

	require.Equal(t, exampleConfig, parsedConfig)
}

func TestExecutionConfig(t *testing.T) {
	exampleConfig := ExecutionPluginJobSpecConfig{
		SourceStartBlock: 222,
		DestStartBlock:   333,
	}

	bts, err := json.Marshal(exampleConfig)
	require.NoError(t, err)

	parsedConfig := ExecutionPluginJobSpecConfig{}
	require.NoError(t, json.Unmarshal(bts, &parsedConfig))

	require.Equal(t, exampleConfig, parsedConfig)
}

func TestUSDCValidate(t *testing.T) {
	testcases := []struct {
		config USDCConfig
		err    string
	}{
		{
			config: USDCConfig{},
			err:    "AttestationAPI is required",
		},
		{
			config: USDCConfig{
				AttestationAPI: "api",
			},
			err: "SourceTokenAddress is required",
		},
		{
			config: USDCConfig{
				AttestationAPI:     "api",
				SourceTokenAddress: utils.ZeroAddress,
			},
			err: "SourceTokenAddress is required",
		},
		{
			config: USDCConfig{
				AttestationAPI:     "api",
				SourceTokenAddress: utils.RandomAddress(),
			},
			err: "SourceMessageTransmitterAddress is required",
		},
		{
			config: USDCConfig{
				AttestationAPI:                  "api",
				SourceTokenAddress:              utils.RandomAddress(),
				SourceMessageTransmitterAddress: utils.ZeroAddress,
			},
			err: "SourceMessageTransmitterAddress is required",
		},
		{
			config: USDCConfig{
				AttestationAPI:                  "api",
				SourceTokenAddress:              utils.RandomAddress(),
				SourceMessageTransmitterAddress: utils.RandomAddress(),
			},
			err: "",
		},
		{
			config: USDCConfig{
				AttestationAPI:                  "api",
				SourceTokenAddress:              utils.RandomAddress(),
				SourceMessageTransmitterAddress: utils.RandomAddress(),
				AttestationAPITimeoutSeconds:    -1,
			},
			err: "AttestationAPITimeoutSeconds must be non-negative",
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(fmt.Sprintf("error = %s", tc.err), func(t *testing.T) {
			t.Parallel()
			err := tc.config.ValidateUSDCConfig()
			if tc.err != "" {
				require.ErrorContains(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
