package resolver

import (
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestResolver_Nodes(t *testing.T) {
	t.Parallel()

	var (
		chainID = *utils.NewBigI(1)

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
				f.Mocks.chainSet.On("NodeStatuses", mock.Anything, PageDefaultOffset, PageDefaultLimit).Return([]types.NodeStatus{
					{
						Name:    "node-name",
						ChainID: chainID.String(),
						Config:  `Name = 'node-name'`,
					},
				}, 1, nil)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.Mocks.evmORM.PutChains(v2.EVMConfig{ChainID: &chainID})
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
				f.Mocks.chainSet.On("NodeStatuses", mock.Anything, PageDefaultOffset, PageDefaultLimit).Return([]types.NodeStatus{}, 0, gError)
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

	var name = "node-name"

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "node"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.Mocks.evmORM.PutChains(v2.EVMConfig{Nodes: []*v2.Node{{
					Name:    &name,
					WSURL:   models.MustParseURL("ws://some-url"),
					HTTPURL: models.MustParseURL("http://some-url"),
				}}})
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
