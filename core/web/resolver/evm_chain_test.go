package resolver

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/ethereum/go-ethereum/common"
	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestResolver_Chains(t *testing.T) {
	var (
		chainID = *utils.NewBigI(1)
		nodeID  = int32(200)

		query = `
			query GetChains {
				chains {
					results {
						id
						enabled
						createdAt
						nodes {
							id
						}
						config {
							blockHistoryEstimatorBlockDelay
							ethTxReaperThreshold
							ethTxResendAfterThreshold
							evmEIP1559DynamicFees
							evmGasLimitMultiplier
							chainType
							gasEstimatorMode
							linkContractAddress
							keySpecificConfigs {
								address
								config {
									blockHistoryEstimatorBlockDelay
									evmEIP1559DynamicFees
								}
							}
						}
					}
					metadata {
						total
					}
				}
			}`
	)
	linkContractAddress := newRandomAddress().String()

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "chains"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				threshold, err := models.MakeDuration(1 * time.Minute)
				require.NoError(t, err)

				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.Mocks.evmORM.PutChains(types.DBChain{
					ID:        chainID,
					Enabled:   true,
					CreatedAt: f.Timestamp(),
					Cfg: &types.ChainCfg{
						BlockHistoryEstimatorBlockDelay: null.IntFrom(1),
						EthTxReaperThreshold:            &threshold,
						EthTxResendAfterThreshold:       &threshold,
						EvmEIP1559DynamicFees:           null.BoolFrom(true),
						EvmGasLimitMultiplier:           null.FloatFrom(1.23),
						GasEstimatorMode:                null.StringFrom("BlockHistory"),
						ChainType:                       null.StringFrom("optimism"),
						LinkContractAddress:             null.StringFrom(linkContractAddress),
						KeySpecific: map[string]types.ChainCfg{
							"test-address": {
								BlockHistoryEstimatorBlockDelay: null.IntFrom(0),
								EvmEIP1559DynamicFees:           null.BoolFrom(false),
							},
						},
					},
				})
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
				f.Mocks.chainSet.On("GetNodesByChainIDs", mock.Anything, []utils.Big{chainID}).
					Return([]types.Node{
						{
							ID:         nodeID,
							EVMChainID: chainID,
						},
					}, nil)
			},
			query: query,
			result: fmt.Sprintf(`
			{
				"chains": {
					"results": [{
						"id": "1",
						"enabled": true,
						"createdAt": "2021-01-01T00:00:00Z",
						"config": {
							"blockHistoryEstimatorBlockDelay": 1,
							"ethTxReaperThreshold": "1m0s",
							"ethTxResendAfterThreshold": "1m0s",
							"evmEIP1559DynamicFees": true,
							"evmGasLimitMultiplier": 1.23,
							"chainType": "OPTIMISM",
							"gasEstimatorMode": "BLOCK_HISTORY",
							"linkContractAddress": "%s",
							"keySpecificConfigs": [{
								"address": "test-address",
								"config": {
									"blockHistoryEstimatorBlockDelay": 0,
									"evmEIP1559DynamicFees": false
								}
							}]
						},
						"nodes": [{
							"id": "200"
						}]
					}],
					"metadata": {
						"total": 1
					}
				}
			}`, linkContractAddress),
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_Chain(t *testing.T) {
	var (
		chainID = *utils.NewBigI(1)
		nodeID  = int32(200)
		query   = `
			query GetChain {
				chain(id: "1") {
					... on Chain {
						id
						enabled
						createdAt
						nodes {
							id
						}
						config {
							blockHistoryEstimatorBlockDelay
							ethTxReaperThreshold
							ethTxResendAfterThreshold
							evmEIP1559DynamicFees
							evmGasLimitMultiplier
							chainType
							gasEstimatorMode
							keySpecificConfigs {
								address
								config {
									blockHistoryEstimatorBlockDelay
									evmEIP1559DynamicFees
								}
							}
						}
					}
					... on NotFoundError {
						code
						message
					}
				}
			}
		`
	)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "chain"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				threshold, err := models.MakeDuration(1 * time.Minute)
				require.NoError(t, err)

				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
				f.Mocks.evmORM.PutChains(types.DBChain{
					ID:        chainID,
					Enabled:   true,
					CreatedAt: f.Timestamp(),
					Cfg: &types.ChainCfg{
						BlockHistoryEstimatorBlockDelay: null.IntFrom(1),
						EthTxReaperThreshold:            &threshold,
						EthTxResendAfterThreshold:       &threshold,
						EvmEIP1559DynamicFees:           null.BoolFrom(true),
						EvmGasLimitMultiplier:           null.FloatFrom(1.23),
						GasEstimatorMode:                null.StringFrom("BlockHistory"),
						ChainType:                       null.StringFrom("optimism"),
						KeySpecific: map[string]types.ChainCfg{
							"test-address": {
								BlockHistoryEstimatorBlockDelay: null.IntFrom(0),
								EvmEIP1559DynamicFees:           null.BoolFrom(false),
							},
						},
					},
				})
				f.Mocks.chainSet.On("GetNodesByChainIDs", mock.Anything, []utils.Big{chainID}).
					Return([]types.Node{
						{
							ID:         nodeID,
							EVMChainID: chainID,
						},
					}, nil)
			},
			query: query,
			result: `
				{
					"chain": {
						"id": "1",
						"enabled": true,
						"createdAt": "2021-01-01T00:00:00Z",
						"config": {
							"blockHistoryEstimatorBlockDelay": 1,
							"ethTxReaperThreshold": "1m0s",
							"ethTxResendAfterThreshold": "1m0s",
							"evmEIP1559DynamicFees": true,
							"evmGasLimitMultiplier": 1.23,
							"chainType": "OPTIMISM",
							"gasEstimatorMode": "BLOCK_HISTORY",
							"keySpecificConfigs": [{
								"address": "test-address",
								"config": {
									"blockHistoryEstimatorBlockDelay": 0,
									"evmEIP1559DynamicFees": false
								}
							}]
						},
						"nodes": [{
							"id": "200"
						}]
					}
				}`,
		},
		{
			name:          "not found error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
			},
			query: query,
			result: `
				{
					"chain": {
						"code": "NOT_FOUND",
						"message": "chain not found"
					}
				}`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_CreateChain(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation CreateChain($input: CreateChainInput!) {
			createChain(input: $input) {
				... on CreateChainSuccess {
					chain {
						id
						enabled
						createdAt
						config {
							blockHistoryEstimatorBlockDelay
							ethTxReaperThreshold
							chainType
							gasEstimatorMode
							linkContractAddress
							keySpecificConfigs {
								address
								config {
									blockHistoryEstimatorBlockDelay
									ethTxReaperThreshold
									chainType
									gasEstimatorMode
								}
							}
						}
					}
				}
				... on InputErrors {
					errors {
						path
						message
						code
					}
				}
			}
		}`

	data, err := json.Marshal(map[string]interface{}{
		"address": "some-address",
		"config": map[string]interface{}{
			"blockHistoryEstimatorBlockDelay": 0,
			"ethTxReaperThreshold":            "1m0s",
			"chainType":                       "XDAI",
			"gasEstimatorMode":                "BLOCK_HISTORY",
		},
	})
	require.NoError(t, err)

	// Ugly hack to avoid type check issues when using slices of maps against the GQL test library...
	// This is because the library internally is trying to assert the slice values against map[string]interface{}
	var keySpecificConfig interface{}
	err = json.Unmarshal(data, &keySpecificConfig)
	require.NoError(t, err)

	linkContractAddress := newRandomAddress().String()

	input := map[string]interface{}{
		"input": map[string]interface{}{
			"id": "1233",
			"config": map[string]interface{}{
				"blockHistoryEstimatorBlockDelay": 1,
				"ethTxReaperThreshold":            "1m0s",
				"chainType":                       "OPTIMISM",
				"gasEstimatorMode":                "BLOCK_HISTORY",
				"linkContractAddress":             linkContractAddress,
			},
			"keySpecificConfigs": []interface{}{
				keySpecificConfig,
			},
		},
	}
	badInput := map[string]interface{}{
		"input": map[string]interface{}{
			"id": "1233",
			"config": map[string]interface{}{
				"ethTxReaperThreshold": "asdadadsa",
				"chainType":            "OPTIMISM",
				"gasEstimatorMode":     "BLOCK_HISTORY",
			},
			"keySpecificConfigs": []interface{}{},
		},
	}

	threshold, err := models.MakeDuration(1 * time.Minute)
	require.NoError(t, err)

	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: input}, "createChain"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				cfg := types.ChainCfg{
					BlockHistoryEstimatorBlockDelay: null.IntFrom(1),
					EthTxReaperThreshold:            &threshold,
					GasEstimatorMode:                null.StringFrom("BlockHistory"),
					ChainType:                       null.StringFrom("optimism"),
					LinkContractAddress:             null.StringFrom(linkContractAddress),
					KeySpecific: map[string]types.ChainCfg{
						"some-address": {
							BlockHistoryEstimatorBlockDelay: null.IntFrom(0),
							EthTxReaperThreshold:            &threshold,
							GasEstimatorMode:                null.StringFrom("BlockHistory"),
							ChainType:                       null.StringFrom("xdai"),
						},
					},
				}

				f.Mocks.chainSet.On("Add", mock.Anything, *utils.NewBigI(1233), &cfg).Return(types.DBChain{
					ID:        *utils.NewBigI(1233),
					Enabled:   true,
					CreatedAt: f.Timestamp(),
					Cfg:       &cfg,
				}, nil)
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
			},
			query:     mutation,
			variables: input,
			result: fmt.Sprintf(`
				{
					"createChain": {
						"chain": {
							"id": "1233",
							"enabled": true,
							"createdAt": "2021-01-01T00:00:00Z",
							"config": {
								"blockHistoryEstimatorBlockDelay": 1,
								"ethTxReaperThreshold": "1m0s",
								"chainType": "OPTIMISM",
								"gasEstimatorMode": "BLOCK_HISTORY",
								"linkContractAddress": "%v",
								"keySpecificConfigs": [
									{
										"address": "some-address",
										"config": {
											"blockHistoryEstimatorBlockDelay": 0,
											"ethTxReaperThreshold": "1m0s",
											"chainType": "XDAI",
											"gasEstimatorMode": "BLOCK_HISTORY"
										}
									}
								]
							}
						}
					}
				}`, linkContractAddress),
		},
		{
			name:          "input errors",
			authenticated: true,
			query:         mutation,
			variables:     badInput,
			result: `
				{
					"createChain": {
						"errors": [{
							"path": "EthTxReaperThreshold",
							"message": "invalid value",
							"code": "INVALID_INPUT"
						}]
					}
				}`,
		},
		{
			name:          "create chain generic error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				cfg := types.ChainCfg{
					BlockHistoryEstimatorBlockDelay: null.IntFrom(1),
					EthTxReaperThreshold:            &threshold,
					GasEstimatorMode:                null.StringFrom("BlockHistory"),
					ChainType:                       null.StringFrom("optimism"),
					LinkContractAddress:             null.StringFrom(linkContractAddress),
					KeySpecific: map[string]types.ChainCfg{
						"some-address": {
							BlockHistoryEstimatorBlockDelay: null.IntFrom(0),
							EthTxReaperThreshold:            &threshold,
							GasEstimatorMode:                null.StringFrom("BlockHistory"),
							ChainType:                       null.StringFrom("xdai"),
						},
					},
				}

				f.Mocks.chainSet.On("Add", mock.Anything, *utils.NewBigI(1233), &cfg).Return(types.DBChain{
					ID:        *utils.NewBigI(1233),
					Enabled:   true,
					CreatedAt: f.Timestamp(),
					Cfg:       &cfg,
				}, gError)
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
			},
			query:     mutation,
			variables: input,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"createChain"},
					Message:       "error",
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_DeleteChain(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation DeleteChain($id: ID!) {
			deleteChain(id: $id) {
				... on DeleteChainSuccess {
					chain {
						id
					}
				}
				... on NotFoundError {
					message
					code
				}
			}
		}`
	variables := map[string]interface{}{
		"id": "123",
	}
	chainID := *utils.NewBigI(123)
	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "deleteChain"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.evmORM.PutChains(types.DBChain{ID: chainID})
				f.Mocks.chainSet.On("Remove", chainID).Return(nil)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"deleteChain": {
						"chain": {
							"id": "123"
						}
					}
				}`,
		},
		{
			name:          "not found error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"deleteChain": {
						"code": "NOT_FOUND",
						"message": "chain not found"
					}
				}`,
		},
		{
			name:          "generic error on delete",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.evmORM.PutChains(types.DBChain{ID: chainID})
				f.Mocks.chainSet.On("Remove", chainID).Return(gError)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
			},
			query:     mutation,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"deleteChain"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_UpdateChain(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation UpdateChain($id: ID!, $input: UpdateChainInput!) {
			updateChain(id: $id, input: $input) {
				... on UpdateChainSuccess {
					chain {
						id
						enabled
						createdAt
						config {
							blockHistoryEstimatorBlockDelay
							ethTxReaperThreshold
							chainType
							gasEstimatorMode
							linkContractAddress
							keySpecificConfigs {
								address
								config {
									blockHistoryEstimatorBlockDelay
									ethTxReaperThreshold
									chainType
									gasEstimatorMode
								}
							}
						}
					}
				}
				... on NotFoundError {
					message
					code
				}
				... on InputErrors {
					errors {
						path
						message
						code
					}
				}
			}
		}`
	chainID := *utils.NewBigI(1233)
	data, err := json.Marshal(map[string]interface{}{
		"address": "some-address",
		"config": map[string]interface{}{
			"blockHistoryEstimatorBlockDelay": 0,
			"ethTxReaperThreshold":            "1m0s",
			"chainType":                       "XDAI",
			"gasEstimatorMode":                "BLOCK_HISTORY",
		},
	})
	require.NoError(t, err)

	var keySpecificConfig interface{}
	err = json.Unmarshal(data, &keySpecificConfig)
	require.NoError(t, err)

	linkContractAddress := newRandomAddress().String()

	input := map[string]interface{}{
		"id": "1233",
		"input": map[string]interface{}{
			"enabled": true,
			"config": map[string]interface{}{
				"blockHistoryEstimatorBlockDelay": 1,
				"ethTxReaperThreshold":            "1m0s",
				"chainType":                       "OPTIMISM",
				"gasEstimatorMode":                "BLOCK_HISTORY",
				"linkContractAddress":             linkContractAddress,
			},
			"keySpecificConfigs": []interface{}{
				keySpecificConfig,
			},
		},
	}
	badInput := map[string]interface{}{
		"id": "1233",
		"input": map[string]interface{}{
			"enabled": true,
			"config": map[string]interface{}{
				"ethTxReaperThreshold": "asdadadsa",
				"chainType":            "OPTIMISM",
				"gasEstimatorMode":     "BLOCK_HISTORY",
			},
			"keySpecificConfigs": []interface{}{},
		},
	}

	threshold, err := models.MakeDuration(1 * time.Minute)
	require.NoError(t, err)

	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: input}, "updateChain"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				cfg := types.ChainCfg{
					BlockHistoryEstimatorBlockDelay: null.IntFrom(1),
					EthTxReaperThreshold:            &threshold,
					GasEstimatorMode:                null.StringFrom("BlockHistory"),
					ChainType:                       null.StringFrom("optimism"),
					LinkContractAddress:             null.StringFrom(linkContractAddress),
					KeySpecific: map[string]types.ChainCfg{
						"some-address": {
							BlockHistoryEstimatorBlockDelay: null.IntFrom(0),
							EthTxReaperThreshold:            &threshold,
							GasEstimatorMode:                null.StringFrom("BlockHistory"),
							ChainType:                       null.StringFrom("xdai"),
						},
					},
				}

				f.Mocks.chainSet.On("Configure", mock.Anything, chainID, true, &cfg).Return(types.DBChain{
					ID:        chainID,
					Enabled:   true,
					CreatedAt: f.Timestamp(),
					Cfg:       &cfg,
				}, nil)
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
			},
			query:     mutation,
			variables: input,
			result: fmt.Sprintf(`
				{
					"updateChain": {
						"chain": {
							"id": "1233",
							"enabled": true,
							"createdAt": "2021-01-01T00:00:00Z",
							"config": {
								"blockHistoryEstimatorBlockDelay": 1,
								"ethTxReaperThreshold": "1m0s",
								"chainType": "OPTIMISM",
								"gasEstimatorMode": "BLOCK_HISTORY",
								"linkContractAddress": "%s",
								"keySpecificConfigs": [
									{
										"address": "some-address",
										"config": {
											"blockHistoryEstimatorBlockDelay": 0,
											"ethTxReaperThreshold": "1m0s",
											"chainType": "XDAI",
											"gasEstimatorMode": "BLOCK_HISTORY"
										}
									}
								]
							}
						}
					}
				}`, linkContractAddress),
		},
		{
			name:          "input errors",
			authenticated: true,
			query:         mutation,
			variables:     badInput,
			result: `
				{
					"updateChain": {
						"errors": [{
							"path": "EthTxReaperThreshold",
							"message": "invalid value",
							"code": "INVALID_INPUT"
						}]
					}
				}`,
		},
		{
			name:          "not found error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				cfg := types.ChainCfg{
					BlockHistoryEstimatorBlockDelay: null.IntFrom(1),
					EthTxReaperThreshold:            &threshold,
					GasEstimatorMode:                null.StringFrom("BlockHistory"),
					ChainType:                       null.StringFrom("optimism"),
					LinkContractAddress:             null.StringFrom(linkContractAddress),
					KeySpecific: map[string]types.ChainCfg{
						"some-address": {
							BlockHistoryEstimatorBlockDelay: null.IntFrom(0),
							EthTxReaperThreshold:            &threshold,
							GasEstimatorMode:                null.StringFrom("BlockHistory"),
							ChainType:                       null.StringFrom("xdai"),
						},
					},
				}

				f.Mocks.chainSet.On("Configure", mock.Anything, chainID, true, &cfg).Return(types.DBChain{}, sql.ErrNoRows)
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
			},
			query:     mutation,
			variables: input,
			result: `
				{
					"updateChain": {
						"code": "NOT_FOUND",
						"message": "chain not found"
					}
				}`,
		},
		{
			name:          "generic error on update",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				cfg := types.ChainCfg{
					BlockHistoryEstimatorBlockDelay: null.IntFrom(1),
					EthTxReaperThreshold:            &threshold,
					GasEstimatorMode:                null.StringFrom("BlockHistory"),
					ChainType:                       null.StringFrom("optimism"),
					LinkContractAddress:             null.StringFrom(linkContractAddress),
					KeySpecific: map[string]types.ChainCfg{
						"some-address": {
							BlockHistoryEstimatorBlockDelay: null.IntFrom(0),
							EthTxReaperThreshold:            &threshold,
							GasEstimatorMode:                null.StringFrom("BlockHistory"),
							ChainType:                       null.StringFrom("xdai"),
						},
					},
				}

				f.Mocks.chainSet.On("Configure", mock.Anything, chainID, true, &cfg).Return(types.DBChain{}, gError)
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
			},
			query:     mutation,
			variables: input,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"updateChain"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}

// Using a local version, since there would be an import cycle if `newRandomAddress()` were to be called in this context.
func newRandomAddress() common.Address {
	b := make([]byte, 20)
	_, _ = rand.Read(b) // Assignment for errcheck. Only used in tests so we can ignore.

	return common.BytesToAddress(b)
}
