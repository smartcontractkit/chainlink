package resolver

import (
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestResolver_Nodes(t *testing.T) {
	t.Parallel()

	var (
		chainID = *utils.NewBigI(1)
		nodeID  = int32(200)

		query = `
			query GetNodes {
				nodes {
					results {
						id
						name
						chain {
							id
						}
					}
					metadata {
						total
					}
				}
			}`
	)
	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "nodes"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
				f.Mocks.chainSet.On("GetNodes", mock.Anything, PageDefaultOffset, PageDefaultLimit).Return([]types.Node{
					{
						ID:         nodeID,
						Name:       "node-name",
						EVMChainID: chainID,
					},
				}, 1, nil)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.Mocks.evmORM.PutChains(chains.ChainConfig{ID: chainID.String()})
			},
			query: query,
			result: `
			{
				"nodes": {
					"results": [{
						"id": "node-name",
						"name": "node-name",
						"chain": {
							"id": "1"
						}
					}],
					"metadata": {
						"total": 1
					}
				}
			}`,
		},
		{
			name:          "generic error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.chainSet.On("GetNodes", mock.Anything, PageDefaultOffset, PageDefaultLimit).Return([]types.Node{}, 0, gError)
				f.App.On("GetChains").Return(chainlink.Chains{EVM: f.Mocks.chainSet})
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"nodes"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}

func Test_NodeQuery(t *testing.T) {
	t.Parallel()

	query := `
		query GetNode {
			node(id: "node-name") {
				... on Node {
					name
					wsURL
					httpURL
				}
				... on NotFoundError {
					message
					code
				}
			}
		}`

	nodeID := int32(200)
	const name = "node-name"

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "node"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.Mocks.evmORM.AddNodes(types.Node{
					ID:      nodeID,
					Name:    name,
					WSURL:   null.StringFrom("ws://some-url"),
					HTTPURL: null.StringFrom("http://some-url"),
				})
			},
			query: query,
			result: `
			{
				"node": {
					"name": "node-name",
					"wsURL": "ws://some-url",
					"httpURL": "http://some-url"
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
				"node": {
					"message": "node not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}
