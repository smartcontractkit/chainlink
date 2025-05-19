package config

import (
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
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
			expectedValidationError: errors.New("aggregator contract address is zero"),
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
			expectedValidationError: errors.New("chain id is zero"),
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
			require.NoError(t, pgc.MoveDeprecatedFields(123, 345, utils.RandomAddress()))
			err = pgc.Validate()
			if test.expectedValidationError != nil {
				require.ErrorContains(t, err, test.expectedValidationError.Error())
			} else {
				require.NoError(t, err)
				require.Empty(t, pgc.AggregatorPrices, 0)
				require.Len(t, pgc.TokenPrices, 3)
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
		t.Run("error = "+tc.err, func(t *testing.T) {
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

func TestDynamicPriceGetterConfig(t *testing.T) {
	// this test goes through unmarshal -> move deprecated -> validate -> assert equal to expected
	// for verifying e2e config loading and validation
	destChain := chainsel.TEST_1000   // 11787463284727550157
	sourceChain := chainsel.TEST_1338 // 2181150070347029680

	sourceNativeTokenAddress := common.HexToAddress("0x0000c05e1fba1244763a494a52272170c0000000")
	expTokenPrices := []TokenPriceConfig{
		{
			TokenAddress:  sourceNativeTokenAddress,
			ChainSelector: sourceChain.Selector,
			AggregatorConfig: &AggregatorPriceConfig{
				ChainID:                   5678,
				AggregatorContractAddress: common.HexToAddress("0xb8dabd288955d302d05ca6b011bb46dfa3e11111"),
			},
		},
		{
			TokenAddress:  common.HexToAddress("0x0820c05e1fba1244763a494a52272170c321cad3"),
			ChainSelector: destChain.Selector,
			AggregatorConfig: &AggregatorPriceConfig{
				ChainID:                   1000,
				AggregatorContractAddress: common.HexToAddress("0xb8dabd288955d302d05ca6b011bb46dfa3ea7acf"),
			},
		},
		{
			TokenAddress:  common.HexToAddress("0x4a98bb4d65347016a7ab6f85bea24b129c9a1272"),
			ChainSelector: destChain.Selector,
			AggregatorConfig: &AggregatorPriceConfig{
				ChainID:                   1337,
				AggregatorContractAddress: common.HexToAddress("0xb80244cc8b0bb18db071c150b36e9bcb8310b236"),
			},
		},
		{
			TokenAddress:  common.HexToAddress("0xec8c353470ccaa4f43067fcde40558e084a12927"),
			ChainSelector: destChain.Selector,
			StaticConfig: &StaticPriceConfig{
				ChainID: 1057,
				Price:   big.NewInt(1000000000000000000),
			},
		},
	}

	testCases := []struct {
		name     string
		jsonCfg  string
		expCfg   DynamicPriceGetterConfig
		expError bool
	}{
		{
			name:     "empty json cfg",
			jsonCfg:  "{}",
			expCfg:   DynamicPriceGetterConfig{},
			expError: false,
		},
		{
			name: "valid json cfg",
			jsonCfg: `
				{
					"randomFieldsAreIgnored": { "foo": "bar" },
  					"tokenPrices": [
  					  {
  					    "tokenAddress": "0x0000c05e1fba1244763a494a52272170c0000000",
  					    "chainSelector": "2181150070347029680",
  					    "aggregatorConfig": {
  					      "chainID": "5678",
  					      "contractAddress": "0xb8dabd288955d302d05ca6b011bb46dfa3e11111"
  					    }
  					  },
  					  {
  					    "tokenAddress": "0x0820c05e1fba1244763a494a52272170c321cad3",
  					    "chainSelector": "11787463284727550157",
  					    "aggregatorConfig": {
  					      "chainID": "1000",
  					      "contractAddress": "0xb8dabd288955d302d05ca6b011bb46dfa3ea7acf"
  					    }
  					  },
  					  {
  					    "tokenAddress": "0x4a98bb4d65347016a7ab6f85bea24b129c9a1272",
  					    "chainSelector": "11787463284727550157",
  					    "aggregatorConfig": {
  					      "chainID": "1337",
  					      "contractAddress": "0xb80244cc8b0bb18db071c150b36e9bcb8310b236"
  					    }
  					  },
  					  {
  					    "tokenAddress": "0xec8c353470ccaa4f43067fcde40558e084a12927",
  					    "chainSelector": "11787463284727550157",
  					    "staticConfig": {
  					      "chainID": "1057",
  					      "price": 1000000000000000000
  					    }
  					  }
  					]
				}
			`,
			expCfg:   DynamicPriceGetterConfig{TokenPrices: expTokenPrices},
			expError: false,
		},
		{
			name: "deprecated json config should be migrated",
			jsonCfg: `
				{
				  "aggregatorPrices": {
				    "0x0000c05e1fba1244763a494a52272170c0000000": {
				      "chainID": "5678",
				      "contractAddress": "0xb8dabd288955d302d05ca6b011bb46dfa3e11111"
				    },
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
			`,
			expCfg:   DynamicPriceGetterConfig{TokenPrices: expTokenPrices},
			expError: false,
		},
		{
			name: "duplicate token in config first and third have same token address and chain",
			jsonCfg: `
				{
				  "tokenPrices": [
				    {
				      "tokenAddress": "0x0820c05e1fba1244763a494a52272170c321cad3",
				      "chainID": "1010",
				      "aggregatorConfig": {
				        "chainID": "1000",
				        "contractAddress": "0xb8dabd288955d302d05ca6b011bb46dfa3ea7acf"
				      }
				    },
				    {
				      "tokenAddress": "0x4a98bb4d65347016a7ab6f85bea24b129c9a1272",
				      "chainID": "1010",
				      "aggregatorConfig": {
				        "chainID": "1337",
				        "contractAddress": "0xb80244cc8b0bb18db071c150b36e9bcb8310b236"
				      }
				    },
				    {
				      "tokenAddress": "0x0820c05e1fba1244763a494a52272170c321cad3",
				      "chainID": "1010",
				      "staticConfig": {
				        "chainID": "1057",
				        "price": 1000000000000000000
				      }
				    }
				  ]
				}
			`,
			expCfg:   DynamicPriceGetterConfig{},
			expError: true,
		},
		{
			name: "both static and aggregator configs are provided for the second token",
			jsonCfg: `
				{
				  "tokenPrices": [
				    {
				      "tokenAddress": "0x0820c05e1fba1244763a494a52272170c321cad3",
				      "chainSelector": "11787463284727550157",
				      "aggregatorConfig": {
				        "chainID": "1000",
				        "contractAddress": "0xb8dabd288955d302d05ca6b011bb46dfa3ea7acf"
				      }
				    },
				    {
				      "tokenAddress": "0x4a98bb4d65347016a7ab6f85bea24b129c9a1272",
				      "chainSelector": "11787463284727550157",
				      "aggregatorConfig": {
				        "chainID": "1337",
				        "contractAddress": "0xb80244cc8b0bb18db071c150b36e9bcb8310b236"
				      },
				      "staticConfig": {
				        "chainID": "1057",
				        "price": 1000000000000000000
				      }
				    },
				    {
				      "tokenAddress": "0xec8c353470ccaa4f43067fcde40558e084a12927",
				      "chainSelector": "11787463284727550157",
				      "staticConfig": {
				        "chainID": "1057",
				        "price": 1000000000000000000
				      }
				    }
				  ]
				}
			`,
			expCfg:   DynamicPriceGetterConfig{},
			expError: true,
		},
		{
			name: "both aggregator and static price config are empty for the first token",
			jsonCfg: `
				{
				  "tokenPrices": [
				    {
				      "tokenAddress": "0x0820c05e1fba1244763a494a52272170c321cad3",
				      "chainSelector": "11787463284727550157"
				    },
				    {
				      "tokenAddress": "0x4a98bb4d65347016a7ab6f85bea24b129c9a1272",
				      "chainSelector": "11787463284727550157",
				      "aggregatorConfig": {
				        "chainID": "1337",
				        "contractAddress": "0xb80244cc8b0bb18db071c150b36e9bcb8310b236"
				      }
				    },
				    {
				      "tokenAddress": "0xec8c353470ccaa4f43067fcde40558e084a12927",
				      "chainSelector": "11787463284727550157",
				      "staticConfig": {
				        "chainID": "1057",
				        "price": 1000000000000000000
				      }
				    }
				  ]
				}
			`,
			expCfg:   DynamicPriceGetterConfig{},
			expError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cfg DynamicPriceGetterConfig
			err := json.Unmarshal([]byte(tc.jsonCfg), &cfg)
			require.NoError(t, err)

			require.NoError(t, err)
			require.NoError(t, cfg.MoveDeprecatedFields(
				sourceChain.Selector,
				destChain.Selector,
				sourceNativeTokenAddress,
			))

			err = cfg.Validate()
			if tc.expError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expCfg, cfg)
		})
	}
}
