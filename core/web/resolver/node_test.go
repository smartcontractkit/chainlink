package resolver

import (
	"context"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	chainlinkmocks "github.com/smartcontractkit/chainlink/v2/core/services/chainlink/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/web/testutils"
)

func TestResolver_Nodes(t *testing.T) {
	t.Parallel()

	query := `
			query GetNodes {
				nodes {
					results {
						id
						name
						chain {
							id
							network
						}
					}
					metadata {
						total
					}
				}
			}`
	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "nodes"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetRelayers").Return(&chainlinkmocks.FakeRelayerChainInteroperators{
					Nodes: []chainlink.NetworkNodeStatus{
						{
							Network: relay.NetworkEVM,
							NodeStatus: types.NodeStatus{
								ChainID: "1",
								Name:    "node-name",
								Config:  "Name='node-name'\nOrder=11\nHTTPURL='http://some-url'\nWSURL='ws://some-url'",
								State:   "alive",
							},
						},
						{
							Network: relay.NetworkAptos,
							NodeStatus: types.NodeStatus{
								ChainID: "2",
								Name:    "aptos-node",
								Config:  "Name='aptos-node'\nURL='http://aptos-url'",
								State:   "alive",
							},
						},
						{
							Network: relay.NetworkStellar,
							NodeStatus: types.NodeStatus{
								ChainID: "stellar-testnet",
								Name:    "stellar-node",
								Config:  "Name='stellar-node'\nURL='http://stellar-url'",
								State:   "alive",
							},
						},
					},
					Relayers: map[types.RelayID]loop.Relayer{
						{
							Network: relay.NetworkEVM,
							ChainID: "1",
						}: &testutils.MockRelayer{ChainStatus: types.ChainStatus{
							ID:      "1",
							Enabled: true,
							Config:  "",
						}},
						{
							Network: relay.NetworkAptos,
							ChainID: "2",
						}: &testutils.MockRelayer{ChainStatus: types.ChainStatus{
							ID:      "2",
							Enabled: true,
						}},
						{
							Network: relay.NetworkStellar,
							ChainID: "stellar-testnet",
						}: &testutils.MockRelayer{ChainStatus: types.ChainStatus{
							ID:      "stellar-testnet",
							Enabled: true,
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
							"id": "1",
							"network": "evm"
						}
					}, {
						"id": "aptos-node",
						"name": "aptos-node",
						"chain": {
							"id": "2",
							"network": "aptos"
						}
					}, {
						"id": "stellar-node",
						"name": "stellar-node",
						"chain": {
							"id": "stellar-testnet",
							"network": "stellar"
						}
					}],
					"metadata": {
						"total": 3
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
					Path:          []any{"nodes"},
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
				f.App.On("GetRelayers").Return(&chainlinkmocks.FakeRelayerChainInteroperators{Relayers: map[types.RelayID]loop.Relayer{
					{
						Network: relay.NetworkEVM,
						ChainID: "1",
					}: &testutils.MockRelayer{NodeStatuses: []types.NodeStatus{
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
			name:          "success non-evm node",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetRelayers").Return(&chainlinkmocks.FakeRelayerChainInteroperators{Relayers: map[types.RelayID]loop.Relayer{
					{
						Network: relay.NetworkStellar,
						ChainID: "stellar-testnet",
					}: &testutils.MockRelayer{
						ChainStatus: types.ChainStatus{ID: "stellar-testnet", Enabled: true},
						NodeStatuses: []types.NodeStatus{
							{
								ChainID: "stellar-testnet",
								Name:    "node-name",
								Config:  "Name='node-name'\nURL='http://stellar-url'",
							},
						},
					},
				}})
			},
			query: `
				query GetNode {
					node(id: "node-name") {
						... on Node {
							name
							chain {
								id
								network
							}
						}
					}
				}`,
			result: `
			{
				"node": {
					"name": "node-name",
					"chain": {
						"id": "stellar-testnet",
						"network": "stellar"
					}
				}
			}`,
		},
		{
			name:          "not found error",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetRelayers").Return(&chainlinkmocks.FakeRelayerChainInteroperators{Relayers: map[types.RelayID]loop.Relayer{}})
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
