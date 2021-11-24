package resolver

import (
	"database/sql"
	"testing"
	"time"

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
