package resolver

import (
	"encoding/json"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
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
						createdAt
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
						CreatedAt:  f.Timestamp(),
					},
				}, 1, nil)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.Mocks.evmORM.PutChains(types.DBChain{ID: chainID})
			},
			query: query,
			result: `
			{
				"nodes": {
					"results": [{
						"id": "node-name",
						"name": "node-name",
						"createdAt": "2021-01-01T00:00:00Z",
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

func Test_CreateNodeMutation(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation CreateNode($input: CreateNodeInput!) {
			createNode(input: $input) {
				... on CreateNodeSuccess {
					node {
						name
						wsURL
						httpURL
						chain {
							id
							enabled
						}
					}
				}
			}
		}`
	createNodeInput := types.NewNode{
		Name:       "node-name",
		EVMChainID: *utils.NewBigI(1),
		WSURL:      null.StringFrom("ws://some-url"),
		HTTPURL:    null.StringFrom("http://some-url"),
		SendOnly:   true,
	}
	input := map[string]interface{}{
		"input": map[string]interface{}{
			"name":       createNodeInput.Name,
			"evmChainID": createNodeInput.EVMChainID,
			"wsURL":      createNodeInput.WSURL,
			"httpURL":    createNodeInput.HTTPURL,
			"sendOnly":   createNodeInput.SendOnly,
		},
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: input}, "createNode"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.Mocks.evmORM.PutChains(types.DBChain{ID: *utils.NewBigI(1), Enabled: true})
			},
			query:     mutation,
			variables: input,
			result: `
				{
					"createNode": {
						"node": {
							"name":       "node-name",
							"wsURL":      "ws://some-url",
							"httpURL":    "http://some-url",
							"chain": {
								"id": "1",
								"enabled": true
							}
						}
					}
				}`,
		},
	}

	RunGQLTests(t, testCases)
}

func Test_DeleteNodeMutation(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation DeleteNode($id: ID!) {
			deleteNode(id: $id) {
				... on DeleteNodeSuccess {
					node {
						name
						wsURL
						httpURL
					}
				}
				... on NotFoundError {
					message
					code
				}
			}
		}`

	fakeID := int32(2)
	const name = "node-name"
	fakeNode := types.Node{
		ID:         fakeID,
		Name:       "node-name",
		EVMChainID: *utils.NewBigI(1),
		WSURL:      null.StringFrom("ws://some-url"),
		HTTPURL:    null.StringFrom("http://some-url"),
		SendOnly:   true,
	}

	variables := map[string]interface{}{
		"id": name,
	}

	d, err := json.Marshal(map[string]interface{}{
		"deleteNode": map[string]interface{}{
			"node": map[string]interface{}{
				"name":    fakeNode.Name,
				"wsURL":   fakeNode.WSURL,
				"httpURL": fakeNode.HTTPURL,
			},
		},
	})
	assert.NoError(t, err)

	expected := string(d)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "deleteNode"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
				f.Mocks.evmORM.AddNodes(fakeNode)
			},
			query:     mutation,
			variables: variables,
			result:    expected,
		},
		{
			name:          "not found error on fetch",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"deleteNode": {
						"code": "NOT_FOUND",
						"message": "node not found"
					}
				}`,
		},
		{
			name:          "not found error on delete",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"deleteNode": {
						"code": "NOT_FOUND",
						"message": "node not found"
					}
				}`,
		},
	}

	RunGQLTests(t, testCases)
}
