package resolver

import (
	"database/sql"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
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

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "chains"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				threshold, err := models.MakeDuration(1 * time.Minute)
				require.NoError(t, err)

				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.Mocks.evmORM.On("Chains", PageDefaultOffset, PageDefaultLimit).Return([]types.Chain{
					{
						ID:        chainID,
						Enabled:   true,
						CreatedAt: f.Timestamp(),
						Cfg: types.ChainCfg{
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
					},
				}, 1, nil)
				f.Mocks.evmORM.On("GetNodesByChainIDs", []utils.Big{chainID}).
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
			}`,
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
				f.Mocks.evmORM.On("Chain", chainID).Return(types.Chain{
					ID:        chainID,
					Enabled:   true,
					CreatedAt: f.Timestamp(),
					Cfg: types.ChainCfg{
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
				}, nil)
				f.Mocks.evmORM.On("GetNodesByChainIDs", []utils.Big{chainID}).
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
				f.Mocks.evmORM.On("Chain", chainID).Return(types.Chain{}, sql.ErrNoRows)
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
			"chainType":                       "EXCHAIN",
			"gasEstimatorMode":                "BLOCK_HISTORY",
		},
	})
	require.NoError(t, err)

	// Ugly hack to avoid type check issues when using slices of maps against the GQL test library...
	// This is because the library internally is trying to assert the slice values against map[string]interface{}
	var keySpecificConfig interface{}
	err = json.Unmarshal(data, &keySpecificConfig)
	require.NoError(t, err)

	input := map[string]interface{}{
		"input": map[string]interface{}{
			"id": "1233",
			"config": map[string]interface{}{
				"blockHistoryEstimatorBlockDelay": 1,
				"ethTxReaperThreshold":            "1m0s",
				"chainType":                       "OPTIMISM",
				"gasEstimatorMode":                "BLOCK_HISTORY",
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
					KeySpecific: map[string]types.ChainCfg{
						"some-address": {
							BlockHistoryEstimatorBlockDelay: null.IntFrom(0),
							EthTxReaperThreshold:            &threshold,
							GasEstimatorMode:                null.StringFrom("BlockHistory"),
							ChainType:                       null.StringFrom("exchain"),
						},
					},
				}

				f.Mocks.chainSet.On("Add", big.NewInt(1233), cfg).Return(types.Chain{
					ID:        *utils.NewBigI(1),
					Enabled:   true,
					CreatedAt: f.Timestamp(),
					Cfg:       cfg,
				}, nil)
				f.App.On("GetChainSet").Return(f.Mocks.chainSet)
			},
			query:     mutation,
			variables: input,
			result: `
				{
					"createChain": {
						"chain": {
							"id": "1",
							"enabled": true,
							"createdAt": "2021-01-01T00:00:00Z",
							"config": {
								"blockHistoryEstimatorBlockDelay": 1,
								"ethTxReaperThreshold": "1m0s",
								"chainType": "OPTIMISM",
								"gasEstimatorMode": "BLOCK_HISTORY",
								"keySpecificConfigs": [
									{
										"address": "some-address",
										"config": {
											"blockHistoryEstimatorBlockDelay": 0,
											"ethTxReaperThreshold": "1m0s",
											"chainType": "EXCHAIN",
											"gasEstimatorMode": "BLOCK_HISTORY"
										}
									}
								]
							}
						}
					}
				}`,
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
					KeySpecific: map[string]types.ChainCfg{
						"some-address": {
							BlockHistoryEstimatorBlockDelay: null.IntFrom(0),
							EthTxReaperThreshold:            &threshold,
							GasEstimatorMode:                null.StringFrom("BlockHistory"),
							ChainType:                       null.StringFrom("exchain"),
						},
					},
				}

				f.Mocks.chainSet.On("Add", big.NewInt(1233), cfg).Return(types.Chain{
					ID:        *utils.NewBigI(1),
					Enabled:   true,
					CreatedAt: f.Timestamp(),
					Cfg:       cfg,
				}, gError)
				f.App.On("GetChainSet").Return(f.Mocks.chainSet)
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
				f.Mocks.evmORM.On("Chain", chainID).Return(types.Chain{
					ID: chainID,
				}, nil)
				f.Mocks.chainSet.On("Remove", chainID.ToInt()).Return(nil)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.App.On("GetChainSet").Return(f.Mocks.chainSet)
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
				f.Mocks.evmORM.On("Chain", chainID).Return(types.Chain{}, sql.ErrNoRows)
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
				f.Mocks.evmORM.On("Chain", chainID).Return(types.Chain{
					ID: chainID,
				}, nil)
				f.Mocks.chainSet.On("Remove", chainID.ToInt()).Return(gError)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.App.On("GetChainSet").Return(f.Mocks.chainSet)
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
