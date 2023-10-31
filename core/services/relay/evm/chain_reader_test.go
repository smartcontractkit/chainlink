package evm_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mocklogpoller "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func TestChainReaderStartClose(t *testing.T) {
	lggr := logger.TestLogger(t)
	lp := mocklogpoller.NewLogPoller(t)
	chainReader, err := evm.NewChainReaderService(lggr, lp)
	require.NoError(t, err)
	require.NotNil(t, chainReader)
	err = chainReader.Start(testutils.Context(t))
	assert.NoError(t, err)
	err = chainReader.Close()
	assert.NoError(t, err)
}

func TestValidateChainReaderConfig(t *testing.T) {
	chainReaderConfigTemplate := `{
	   "chainContractReaders": {
	     "testContract": {
			   "contractName": "testContract",
			   "contractABI":  "[%s]",
			   "chainReaderDefinitions": {
					%s
				}
	     }
	   }
	}`

	type testCase struct {
		name                    string
		abiInput                string
		chainReadingDefinitions string
	}

	var testCases []testCase
	testCases = append(testCases,
		testCase{
			name:     "eventWithNoIndexedTopics",
			abiInput: `{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint112","name":"reserve0","type":"uint112"},{"indexed":false,"internalType":"uint112","name":"reserve1","type":"uint112"}],"name":"Sync","type":"event"}`,
			chainReadingDefinitions: ` "Sync":{
											"chainSpecificName": "Sync",
											"returnValues": [
												"reserve0",
												"reserve1"
											],
											"readType": 1
										}`,
		})

	testCases = append(testCases,
		testCase{
			name:     "eventWithMultipleIndexedTopics",
			abiInput: `{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount0In","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount1In","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount0Out","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount1Out","type":"uint256"},{"indexed":true,"internalType":"address","name":"to","type":"address"}],"name":"Swap","type":"event"}`,
			chainReadingDefinitions: `"Swap":{
											"chainSpecificName": "Swap",
											"params":{
												"sender": "0x0",
												"to": "0x0"
											},
											"returnValues": [
												"sender",
												"amount0In",
												"amount1In",
												"amount0Out",
												"amount1Out",
												"to"
											],
											"readType": 1
										}`,
		})

	testCases = append(testCases,
		testCase{
			name:     "functionWithOneParamAndMultipleResponses",
			abiInput: `{"constant":true,"inputs":[{"internalType":"address","name":"_user","type":"address"}],"name":"getUserAccountData","outputs":[{"internalType":"uint256","name":"totalLiquidityETH","type":"uint256"},{"internalType":"uint256","name":"totalCollateralETH","type":"uint256"},{"internalType":"uint256","name":"totalBorrowsETH","type":"uint256"},{"internalType":"uint256","name":"totalFeesETH","type":"uint256"},{"internalType":"uint256","name":"availableBorrowsETH","type":"uint256"},{"internalType":"uint256","name":"currentLiquidationThreshold","type":"uint256"},{"internalType":"uint256","name":"ltv","type":"uint256"},{"internalType":"uint256","name":"healthFactor","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}`,
			chainReadingDefinitions: `"getUserAccountData":{
											"chainSpecificName": "getUserAccountData",
											"params":{
												"_user": "0x0"
											},
											"returnValues": [
												"totalLiquidityETH",
												"totalCollateralETH",
												"totalBorrowsETH",
												"totalFeesETH",
												"availableBorrowsETH",
												"currentLiquidationThreshold",
												"healthFactor"
											],
											"readType": 0
										}`,
		})

	testCases = append(testCases,
		testCase{
			name:     "functionWithNoParamsAndOneResponseWithNoName",
			abiInput: `{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}`,
			chainReadingDefinitions: `"name":{
											"chainSpecificName": "name",
											"returnValues": [
												""
											],
											"readType": 0
										}`,
		})

	testCases = append(testCases,
		testCase{
			name:     "functionWithMultipleParamsAndOneResult",
			abiInput: `{"inputs":[{"internalType":"address","name":"_input","type":"address"},{"internalType":"address","name":"_output","type":"address"},{"internalType":"uint256","name":"_inputQuantity","type":"uint256"}],"name":"getSwapOutput","outputs":[{"internalType":"uint256","name":"swapOutput","type":"uint256"}],"stateMutability":"view","type":"function"}`,
			chainReadingDefinitions: `"getSwapOutput":{
											"chainSpecificName": "getSwapOutput",
											"params":{
												"_input":"0x0",
												"_output":"0x0",
												"_inputQuantity":"0x0"
											},
											"returnValues": [
												"swapOutput"
											],
											"readType": 0
										}`,
		})

	// TODO how to handle return values for tuples
	/*testCases = append(testCases,
	testCase{
		name: "functionWithOneParamAndMultipleTupleResponse",
		 struct BassetPersonal {
		    // Address of the bAsset
		    address addr;
		    // Address of the bAsset
		    address integrator;
		    // An ERC20 can charge transfer fee, for example USDT, DGX tokens.
		    bool hasTxFee; // takes a byte in storage
		    // Status of the bAsset
		    BassetStatus status;
		}

		// Status of the Basset - has it broken its peg?
		enum BassetStatus {
		    Default,
		    Normal,
		    BrokenBelowPeg,
		    BrokenAbovePeg,
		    Blacklisted,
		    Liquidating,
		    Liquidated,
		    Failed
		}

		struct BassetData {
		    // 1 Basset * ratio / ratioScale == x Masset (relative value)
		    // If ratio == 10e8 then 1 bAsset = 10 mAssets
		    // A ratio is divised as 10^(18-tokenDecimals) * measurementMultiple(relative value of 1 base unit)
		    uint128 ratio;
		    // Amount of the Basset that is held in Collateral
		    uint128 vaultBalance;
		}
		abiInput: `{"inputs":[{"internalType":"address","name":"_bAsset","type":"address"}],"name":"getBasset","outputs":[{"components":[{"internalType":"address","name":"addr","type":"address"},{"internalType":"address","name":"integrator","type":"address"},{"internalType":"bool","name":"hasTxFee","type":"bool"},{"internalType":"enum BassetStatus","name":"status","type":"uint8"}],"internalType":"struct BassetPersonal","name":"personal","type":"tuple"},{"components":[{"internalType":"uint128","name":"ratio","type":"uint128"},{"internalType":"uint128","name":"vaultBalance","type":"uint128"}],"internalType":"struct BassetData","name":"vaultData","type":"tuple"}],"stateMutability":"view","type":"function"}`,
		chainReadingDefinitions: `getBasset:{
										chainSpecificName: getBasset,
										params:{
											_bAsset:"0x0",
										},
										returnValues: [
											TODO,
										]
										readType: 0,
									}`,
	})
	*/

	// TODO how to handle return values for tuples
	/*
		testCases = append(testCases,
			testCase{
				name: "functionWithNoParamsAndTupleResponse",
				 struct FeederConfig {
					uint256 supply;
					uint256 a;
					WeightLimits limits;
				}

				struct WeightLimits {
					uint128 min;
					uint128 max;
				}
				abiInput: `{"inputs":[],"name":"getConfig","outputs":[{"components":[{"internalType":"uint256","name":"supply","type":"uint256"},{"internalType":"uint256","name":"a","type":"uint256"},{"components":[{"internalType":"uint128","name":"min","type":"uint128"},{"internalType":"uint128","name":"max","type":"uint128"}],"internalType":"struct WeightLimits","name":"limits","type":"tuple"}],"internalType":"struct FeederConfig","name":"config","type":"tuple"}],"stateMutability":"view","type":"function"}`,
				chainReadingDefinitions: `getConfig:{
												chainSpecificName: getConfig,
												params:{},
												returnValues: [
													TODO,
												]
												readType: 0,
											}`,
			})*/

	var cfg types.ChainReaderConfig
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			abiString := strings.Replace(tc.abiInput, `"`, `\"`, -1)
			formattedCfgJsonString := fmt.Sprintf(chainReaderConfigTemplate, abiString, tc.chainReadingDefinitions)
			assert.NoError(t, json.Unmarshal([]byte(formattedCfgJsonString), &cfg))
			assert.NoError(t, evm.ValidateChainReaderConfig(cfg))
		})
	}

	t.Run("large config with all test cases", func(t *testing.T) {
		var largeABI string
		var manyChainReadingDefinitions string
		for _, tc := range testCases {
			largeABI += tc.abiInput + ","
			manyChainReadingDefinitions += tc.chainReadingDefinitions + ","
		}

		largeABI = largeABI[:len(largeABI)-1]
		manyChainReadingDefinitions = manyChainReadingDefinitions[:len(manyChainReadingDefinitions)-1]
		formattedCfgJsonString := fmt.Sprintf(chainReaderConfigTemplate, strings.Replace(largeABI, `"`, `\"`, -1), manyChainReadingDefinitions)
		fmt.Println(formattedCfgJsonString)
		assert.NoError(t, json.Unmarshal([]byte(formattedCfgJsonString), &cfg))
		assert.NoError(t, evm.ValidateChainReaderConfig(cfg))
	})
}
