package config

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
)

func TestCommitConfig(t *testing.T) {
	tests := []struct {
		name                    string
		cfg                     CommitPluginJobSpecConfig
		expectedValidationError error
	}{
		{
			name: "valid config",
			cfg: CommitPluginJobSpecConfig{
				SourceStartBlock:       222,
				DestStartBlock:         333,
				OffRamp:                ccipcalc.HexToAddress("0x123"),
				TokenPricesUSDPipeline: `merge [type=merge left="{}" right="{\"0xC79b96044906550A5652BCf20a6EA02f139B9Ae5\":\"1000000000000000000\"}"];`,
				PriceGetterConfig: &DynamicPriceGetterConfig{
					AggregatorPrices: map[common.Address]AggregatorPriceConfig{
						common.HexToAddress("0x0820c05e1fba1244763a494a52272170c321cad3"): {
							ChainID:                   1000,
							AggregatorContractAddress: common.HexToAddress("0xb8dabd288955d302d05ca6b011bb46dfa3ea7acf"),
						},
						common.HexToAddress("0x4a98bb4d65347016a7ab6f85bea24b129c9a1272"): {
							ChainID:                   1337,
							AggregatorContractAddress: common.HexToAddress("0xb80244cc8b0bb18db071c150b36e9bcb8310b236"),
						},
					},
					StaticPrices: map[common.Address]StaticPriceConfig{
						common.HexToAddress("0xec8c353470ccaa4f43067fcde40558e084a12927"): {
							ChainID: 1057,
							Price:   big.NewInt(1000000000000000000),
						},
					},
				},
			},
			expectedValidationError: nil,
		},
		{
			name: "missing dynamic aggregator contract address",
			cfg: CommitPluginJobSpecConfig{
				SourceStartBlock:       222,
				DestStartBlock:         333,
				OffRamp:                ccipcalc.HexToAddress("0x123"),
				TokenPricesUSDPipeline: `merge [type=merge left="{}" right="{\"0xC79b96044906550A5652BCf20a6EA02f139B9Ae5\":\"1000000000000000000\"}"];`,
				PriceGetterConfig: &DynamicPriceGetterConfig{
					AggregatorPrices: map[common.Address]AggregatorPriceConfig{
						common.HexToAddress("0x0820c05e1fba1244763a494a52272170c321cad3"): {
							ChainID:                   1000,
							AggregatorContractAddress: common.HexToAddress("0xb8dabd288955d302d05ca6b011bb46dfa3ea7acf"),
						},
						common.HexToAddress("0x4a98bb4d65347016a7ab6f85bea24b129c9a1272"): {
							ChainID:                   1337,
							AggregatorContractAddress: common.HexToAddress(""),
						},
					},
					StaticPrices: map[common.Address]StaticPriceConfig{
						common.HexToAddress("0xec8c353470ccaa4f43067fcde40558e084a12927"): {
							ChainID: 1057,
							Price:   big.NewInt(1000000000000000000),
						},
					},
				},
			},
			expectedValidationError: fmt.Errorf("aggregator contract address is zero"),
		},
		{
			name: "missing chain ID",
			cfg: CommitPluginJobSpecConfig{
				SourceStartBlock:       222,
				DestStartBlock:         333,
				OffRamp:                ccipcalc.HexToAddress("0x123"),
				TokenPricesUSDPipeline: `merge [type=merge left="{}" right="{\"0xC79b96044906550A5652BCf20a6EA02f139B9Ae5\":\"1000000000000000000\"}"];`,
				PriceGetterConfig: &DynamicPriceGetterConfig{
					AggregatorPrices: map[common.Address]AggregatorPriceConfig{
						common.HexToAddress("0x0820c05e1fba1244763a494a52272170c321cad3"): {
							ChainID:                   1000,
							AggregatorContractAddress: common.HexToAddress("0xb8dabd288955d302d05ca6b011bb46dfa3ea7acf"),
						},
						common.HexToAddress("0x4a98bb4d65347016a7ab6f85bea24b129c9a1272"): {
							ChainID:                   1337,
							AggregatorContractAddress: common.HexToAddress("0xb80244cc8b0bb18db071c150b36e9bcb8310b236"),
						},
					},
					StaticPrices: map[common.Address]StaticPriceConfig{
						common.HexToAddress("0xec8c353470ccaa4f43067fcde40558e084a12927"): {
							ChainID: 0,
							Price:   big.NewInt(1000000000000000000),
						},
					},
				},
			},
			expectedValidationError: fmt.Errorf("chain id is zero"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Verify proper marshall/unmarshalling of the config.
			bts, err := json.Marshal(test.cfg)
			require.NoError(t, err)
			parsedConfig := CommitPluginJobSpecConfig{}
			require.NoError(t, json.Unmarshal(bts, &parsedConfig))
			require.Equal(t, test.cfg, parsedConfig)

			// Ensure correctness of price getter configuration.
			pgc := test.cfg.PriceGetterConfig
			err = pgc.Validate()
			if test.expectedValidationError != nil {
				require.ErrorContains(t, err, test.expectedValidationError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, uint64(1000), pgc.AggregatorPrices[common.HexToAddress("0x0820c05e1fba1244763a494a52272170c321cad3")].ChainID)
				require.Equal(t, uint64(1337), pgc.AggregatorPrices[common.HexToAddress("0x4a98bb4d65347016a7ab6f85bea24b129c9a1272")].ChainID)
				require.Equal(t, uint64(1057), pgc.StaticPrices[common.HexToAddress("0xec8c353470ccaa4f43067fcde40558e084a12927")].ChainID)
			}
		})
	}
}

func TestExecutionConfig(t *testing.T) {
	exampleConfig := ExecPluginJobSpecConfig{
		SourceStartBlock: 222,
		DestStartBlock:   333,
	}

	bts, err := json.Marshal(exampleConfig)
	require.NoError(t, err)

	parsedConfig := ExecPluginJobSpecConfig{}
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

func TestUnmarshallDynamicPriceConfig(t *testing.T) {
	jsonCfg := `
{
	"aggregatorPrices": {
		"0x0820c05e1fba1244763a494a52272170c321cad3": {
			"chainID": "1000",
			"contractAddress": "0xb8dabd288955d302d05ca6b011bb46dfa3ea7acf"
		},
		"0x4a98bb4d65347016a7ab6f85bea24b129c9a1272": {
			"chainID": "1337",
			"contractAddress": "0xb80244cc8b0bb18db071c150b36e9bcb8310b236"
		}
	},
	"staticPrices": {
		"0xec8c353470ccaa4f43067fcde40558e084a12927": {
			"chainID": "1057",
			"price": 1000000000000000000
		}
	}
}
`
	var cfg DynamicPriceGetterConfig
	err := json.Unmarshal([]byte(jsonCfg), &cfg)
	require.NoError(t, err)
	err = cfg.Validate()
	require.NoError(t, err)
}
