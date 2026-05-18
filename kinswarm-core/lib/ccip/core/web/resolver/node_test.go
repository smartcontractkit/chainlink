package resolver

import (
	"context"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	chainlinkmocks "github.com/smartcontractkit/chainlink/v2/core/services/chainlink/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/web/testutils"
)

func TestResolver_Nodes(t *testing.T) {
	t.Parallel()

	var (
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
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetRelayers").Return(&chainlinkmocks.FakeRelayerChainInteroperators{
					Nodes: []types.NodeStatus{
						{
							ChainID: "1",
							Name:    "node-name",
							Config:  "Name='node-name'\nOrder=11\nHTTPURL='http://some-url'\nWSURL='ws://some-url'",
							State:   "alive",
						},
					},
					Relayers: []loop.Relayer{
						testutils.MockRelayer{ChainStatus: types.ChainStatus{
							ID:      "1",
							Enabled: true,
							Config:  "",
						}},
					},
				})
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
			before: func(ctx context.Context, f *gqlTestFramework) {
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

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "node"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetRelayers").Return(&chainlinkmocks.FakeRelayerChainInteroperators{Relayers: []loop.Relayer{
					testutils.MockRelayer{NodeStatuses: []types.NodeStatus{
						{
							Name:   "node-name",
							Config: "Name='node-name'\nOrder=11\nHTTPURL='http://some-url'\nWSURL='ws://some-url'",
						},
					}},
				}})
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
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetRelayers").Return(&chainlinkmocks.FakeRelayerChainInteroperators{Relayers: []loop.Relayer{}})
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
