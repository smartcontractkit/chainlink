package resolver

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestResolver_Chains(t *testing.T) {
	var (
		chainID = *utils.NewBigI(1)

		query = `
			query GetChains {
				chains {
					results {
						id
						enabled
						config
					}
					metadata {
						total
					}
				}
			}`
		configTOML = `BlockHistoryEstimatorBlockDelay = 1
EthTxReaperThreshold = '1m'
EthTxResendAfterThreshold = '1m'
EvmEIP1559DynamicFees = true
EvmGasLimitMultiplier = 1.23
GasEstimatorMode = "BlockHistory"
ChainType = "optimism"
[[KeySpecific]]
Address = "test-address"
BlockHistoryEstimatorBlockDelay = 0
EvmEIP1559DynamicFees = false
`
	)
	configTOMLEscaped, err := json.Marshal(configTOML)
	require.NoError(t, err)
	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "chains"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("EVMORM").Return(f.Mocks.evmORM)

				f.Mocks.evmORM.PutChains(chains.ChainConfig{
					ID:      chainID.String(),
					Enabled: true,
					Cfg:     configTOML,
				})
			},
			query: query,
			result: fmt.Sprintf(`
			{
				"chains": {
					"results": [{
						"id": "1",
						"enabled": true,
						"config": %s
					}],
					"metadata": {
						"total": 1
					}
				}
			}`, configTOMLEscaped),
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_Chain(t *testing.T) {
	var (
		chainID = "1"
		query   = `
			query GetChain {
				chain(id: "1") {
					... on Chain {
						id
						enabled
						config
					}
					... on NotFoundError {
						code
						message
					}
				}
			}
		`
		configTOML = `BlockHistoryEstimatorBlockDelay = 1
EthTxReaperThreshold = '1m'
EthTxResendAfterThreshold = '1m'
EvmEIP1559DynamicFees = true
EvmGasLimitMultiplier = 1.23
GasEstimatorMode = "BlockHistory"
ChainType = "optimism"
[[KeySpecific]]
Address = "test-address"
BlockHistoryEstimatorBlockDelay = 0
EvmEIP1559DynamicFees = false
`
	)
	configTOMLEscaped, err := json.Marshal(configTOML)
	require.NoError(t, err)
	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "chain"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.Mocks.evmORM.PutChains(chains.ChainConfig{
					ID:      chainID,
					Enabled: true,
					Cfg:     configTOML,
				})
			},
			query: query,
			result: fmt.Sprintf(`
				{
					"chain": {
						"id": "1",
						"enabled": true,
						"config": %s
					}
				}`, configTOMLEscaped),
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
