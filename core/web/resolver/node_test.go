package resolver

import (
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func TestResolver_Nodes(t *testing.T) {
	t.Parallel()

	var (
		chainID = *big.NewI(1)

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
				f.App.On("GetRelayers").Return(chainlink.RelayerChainInteroperators(f.Mocks.relayerChainInterops))
				f.Mocks.relayerChainInterops.Nodes = []types.NodeStatus{
					{
						Name:    "node-name",
						ChainID: chainID.String(),
						Config:  `Name = 'node-name'`,
					},
				}
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.Mocks.evmORM.PutChains(toml.EVMConfig{ChainID: &chainID})
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
				f.Mocks.relayerChainInterops.NodesErr = gError
				f.App.On("GetRelayers").Return(f.Mocks.relayerChainInterops)
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
					order
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
				f.Mocks.evmORM.PutChains(toml.EVMConfig{Nodes: []*toml.Node{{
					Name:    &name,
					WSURL:   commonconfig.MustParseURL("ws://some-url"),
					HTTPURL: commonconfig.MustParseURL("http://some-url"),
					Order:   ptr(int32(11)),
				}}})
			},
			query: query,
			result: `
			{
				"node": {
					"name": "node-name",
					"wsURL": "ws://some-url",
					"httpURL": "http://some-url",
					"order": 11
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

func ptr[T any](t T) *T { return &t }
