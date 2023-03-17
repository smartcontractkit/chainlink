package resolver

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/ethereum/go-ethereum/common"
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

				f.Mocks.evmORM.PutChains(types.ChainConfig{
					ID:      chainID,
					Enabled: true,
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
				f.Mocks.evmORM.AddNodes(types.Node{
					ID:         nodeID,
					Name:       "node-name",
					EVMChainID: chainID,
				})
				f.Mocks.chainSet.On("GetNodesByChainIDs", mock.Anything, []utils.Big{chainID}).
					Return(f.Mocks.evmORM.GetNodesByChainIDs([]utils.Big{chainID}))
			},
			query: query,
			result: fmt.Sprintf(`
			{
				"chains": {
					"results": [{
						"id": "1",
						"enabled": true,
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
							"id": "node-name"
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
				f.Mocks.evmORM.PutChains(types.ChainConfig{
					ID:      chainID,
					Enabled: true,
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
				f.Mocks.evmORM.AddNodes(types.Node{
					ID:         nodeID,
					Name:       "node-name",
					EVMChainID: chainID,
				})
				f.Mocks.chainSet.On("GetNodesByChainIDs", mock.Anything, []utils.Big{chainID}).
					Return(f.Mocks.evmORM.GetNodesByChainIDs([]utils.Big{chainID}))
			},
			query: query,
			result: `
				{
					"chain": {
						"id": "1",
						"enabled": true,
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
							"id": "node-name"
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

// Using a local version, since there would be an import cycle if `newRandomAddress()` were to be called in this context.
func newRandomAddress() common.Address {
	b := make([]byte, 20)
	_, _ = rand.Read(b) // Assignment for errcheck. Only used in tests so we can ignore.

	return common.BytesToAddress(b)
}
