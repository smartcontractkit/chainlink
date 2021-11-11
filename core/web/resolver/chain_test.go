package resolver

import (
	"database/sql"
	"testing"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_Chains(t *testing.T) {
	var (
		chainID = *utils.NewBigI(1)
		nodeID  = int32(200)

		query = `
			query GetChains {
				chains {
					id
					enabled
					createdAt
					nodes {
						id
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
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.Mocks.evmORM.On("Chains", PageDefaultOffset, PageDefaultLimit).Return([]types.Chain{
					{
						ID:        chainID,
						Enabled:   true,
						CreatedAt: f.Timestamp(),
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
				"chains": [{
					"id": "1",
					"enabled": true,
					"createdAt": "2021-01-01T00:00:00Z",
					"nodes": [{
						"id": "200"
					}]
				}]
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
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.Mocks.evmORM.On("Chain", chainID).Return(types.Chain{
					ID:        chainID,
					Enabled:   true,
					CreatedAt: f.Timestamp(),
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
