package evm

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	mocklogpoller "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type chainReaderTestHelper struct {
}

func (crTestHelper chainReaderTestHelper) makeChainReaderConfig(abi string, params map[string]any) evmtypes.ChainReaderConfig {
	return evmtypes.ChainReaderConfig{
		ChainContractReaders: map[string]evmtypes.ChainContractReader{
			"MyContract": {
				ContractABI: abi,
				ChainReaderDefinitions: map[string]evmtypes.ChainReaderDefinition{
					"MyGenericMethod": {
						ChainSpecificName: "name",
						Params:            params,
						CacheEnabled:      false,
						ReadType:          evmtypes.Method,
					},
				},
			},
		},
	}
}

func (crTestHelper chainReaderTestHelper) makeChainReaderConfigFromStrings(abi string, chainReadingDefinitions string) (evmtypes.ChainReaderConfig, error) {
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

	abi = strings.Replace(abi, `"`, `\"`, -1)
	formattedCfgJsonString := fmt.Sprintf(chainReaderConfigTemplate, abi, chainReadingDefinitions)
	var chainReaderConfig evmtypes.ChainReaderConfig
	err := json.Unmarshal([]byte(formattedCfgJsonString), &chainReaderConfig)
	return chainReaderConfig, err
}

func TestNewChainReader(t *testing.T) {
	lggr := logger.TestLogger(t)
	lp := mocklogpoller.NewLogPoller(t)
	chain := mocks.NewChain(t)
	contractID := testutils.NewAddress()
	contractABI := `[{"inputs":[{"internalType":"string","name":"param","type":"string"}],"name":"name","stateMutability":"view","type":"function"}]`

	t.Run("happy path", func(t *testing.T) {
		params := make(map[string]any)
		params["param"] = ""
		chainReaderConfig := chainReaderTestHelper{}.makeChainReaderConfig(contractABI, params)
		chain.On("LogPoller").Return(lp)
		_, err := NewChainReaderService(lggr, chain.LogPoller(), contractID, chainReaderConfig)
		assert.NoError(t, err)
	})

	t.Run("invalid config", func(t *testing.T) {
		invalidChainReaderConfig := chainReaderTestHelper{}.makeChainReaderConfig(contractABI, map[string]any{}) // missing param
		_, err := NewChainReaderService(lggr, chain.LogPoller(), contractID, invalidChainReaderConfig)
		assert.ErrorIs(t, err, commontypes.ErrInvalidConfig)
	})

	t.Run("ChainReader config is empty", func(t *testing.T) {
		emptyChainReaderConfig := evmtypes.ChainReaderConfig{}
		_, err := NewChainReaderService(lggr, chain.LogPoller(), contractID, emptyChainReaderConfig)
		assert.ErrorIs(t, err, commontypes.ErrInvalidConfig)
		assert.ErrorContains(t, err, "no contract readers defined")
	})
}

func TestChainReaderStartClose(t *testing.T) {
	lggr := logger.TestLogger(t)
	lp := mocklogpoller.NewLogPoller(t)
	cr := chainReader{
		lggr: lggr,
		lp:   lp,
	}
	err := cr.Start(testutils.Context(t))
	assert.NoError(t, err)
	err = cr.Close()
	assert.NoError(t, err)
}

// TODO Chain Reading Definitions return values are WIP, waiting on codec work and BCF-2789
func TestValidateChainReaderConfig_HappyPath(t *testing.T) {
	type testCase struct {
		name                    string
		abiInput                string
		chainReadingDefinitions string
	}

	var testCases []testCase
	testCases = append(testCases,
		testCase{
			name:     "eventWithMultipleIndexedTopics",
			abiInput: `{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"}],"name":"Swap","type":"event"}`,
			chainReadingDefinitions: `"Swap":{
											"chainSpecificName": "Swap",
											"params":{
												"sender": "0x0",
												"to": "0x0"
											},
											"readType": 1
										}`,
		})

	testCases = append(testCases,
		testCase{
			name:     "methodWithOneParamAndMultipleResponses",
			abiInput: `{"constant":true,"inputs":[{"internalType":"address","name":"_user","type":"address"}],"name":"getUserAccountData","payable":false,"stateMutability":"view","type":"function"}`,
			chainReadingDefinitions: `"getUserAccountData":{
											"chainSpecificName": "getUserAccountData",
											"params":{
												"_user": "0x0"
											},
											"readType": 0
										}`,
		})

	testCases = append(testCases,
		testCase{
			name:     "methodWithMultipleParamsAndOneResult",
			abiInput: `{"inputs":[{"internalType":"address","name":"_input","type":"address"},{"internalType":"address","name":"_output","type":"address"},{"internalType":"uint256","name":"_inputQuantity","type":"uint256"}],"name":"getSwapOutput","stateMutability":"view","type":"function"}`,
			chainReadingDefinitions: `"getSwapOutput":{
											"chainSpecificName": "getSwapOutput",
											"params":{
												"_input":"0x0",
												"_output":"0x0",
												"_inputQuantity":"0x0"
											},
											"readType": 0
										}`,
		})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := chainReaderTestHelper{}.makeChainReaderConfigFromStrings(tc.abiInput, tc.chainReadingDefinitions)
			assert.NoError(t, err)
			assert.NoError(t, validateChainReaderConfig(cfg))
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
		cfg, err := chainReaderTestHelper{}.makeChainReaderConfigFromStrings(largeABI, manyChainReadingDefinitions)
		assert.NoError(t, err)
		assert.NoError(t, validateChainReaderConfig(cfg))
	})
}

// TODO Chain Reading Definitions return values are WIP, waiting on codec work and BCF-2789
func TestValidateChainReaderConfig_BadPath(t *testing.T) {
	type testCase struct {
		name                    string
		abiInput                string
		chainReadingDefinitions string
		expected                error
	}

	var testCases []testCase
	mismatchedEventArgumentsTestABI := `{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"}],"name":"Swap","type":"event"}`
	testCases = append(testCases,
		testCase{
			name:     "mismatched abi and event chain reading param values",
			abiInput: mismatchedEventArgumentsTestABI,
			chainReadingDefinitions: `"Swap":{
													"chainSpecificName": "Swap",
													"params":{
														"malformedParam": "0x0"
													},
													"readType": 1
												}`,
			expected: fmt.Errorf("invalid chainreading definition: \"Swap\" for contract: \"testContract\", err: params: [malformedParam] don't match abi event indexed inputs: [sender]"),
		})

	mismatchedFunctionArgumentsTestABI := `{"constant":true,"inputs":[{"internalType":"address","name":"from","type":"address"}],"name":"Swap","payable":false,"stateMutability":"view","type":"function"}`
	testCases = append(testCases,
		testCase{
			name:     "mismatched abi and method chain reading param values",
			abiInput: mismatchedFunctionArgumentsTestABI,
			chainReadingDefinitions: `"Swap":{
											"chainSpecificName": "Swap",
											"params":{
												"malformedParam": "0x0"
											},
											"readType": 0
										}`,
			expected: fmt.Errorf("invalid chainreading definition: \"Swap\" for contract: \"testContract\", err: params: [malformedParam] don't match abi method inputs: [from]"),
		},
	)

	testCases = append(testCases,
		testCase{
			name:     "event doesn't exist",
			abiInput: `{"constant":true,"inputs":[],"name":"someName","payable":false,"stateMutability":"view","type":"function"}`,
			chainReadingDefinitions: `"TestMethod":{
											"chainSpecificName": "Swap",
											"readType": 1
										}`,
			expected: fmt.Errorf("invalid chainreading definition: \"TestMethod\" for contract: \"testContract\", err: event: Swap doesn't exist"),
		},
	)

	testCases = append(testCases,
		testCase{
			name:     "method doesn't exist",
			abiInput: `{"constant":true,"inputs":[],"name":"someName","payable":false,"stateMutability":"view","type":"function"}`,
			chainReadingDefinitions: `"TestMethod":{
											"chainSpecificName": "Swap",
											"readType": 0
										}`,
			expected: fmt.Errorf("invalid chainreading definition: \"TestMethod\" for contract: \"testContract\", err: method: \"Swap\" doesn't exist"),
		},
	)

	testCases = append(testCases, testCase{
		name:     "invalid abi",
		abiInput: `broken abi`,
		chainReadingDefinitions: `"TestMethod":{
											"chainSpecificName": "Swap",
											"readType": 0
										}`,
		expected: fmt.Errorf("invalid abi"),
	})

	testCases = append(testCases, testCase{
		name:                    "invalid read type",
		abiInput:                `{"constant":true,"inputs":[],"name":"someName","payable":false,"stateMutability":"view","type":"function"}`,
		chainReadingDefinitions: `"TestMethod":{"readType": 59}`,
		expected:                fmt.Errorf("invalid chainreading definition read type: 59"),
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := chainReaderTestHelper{}.makeChainReaderConfigFromStrings(tc.abiInput, tc.chainReadingDefinitions)
			assert.NoError(t, err)
			if tc.expected == nil {
				assert.NoError(t, validateChainReaderConfig(cfg))
			} else {
				assert.ErrorContains(t, validateChainReaderConfig(cfg), tc.expected.Error())
			}
		})
	}
}
