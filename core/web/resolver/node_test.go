package resolver

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_NodeQuery(t *testing.T) {
	t.Parallel()

	query := `
		query GetNode {
			node(id: "200") {
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
	notFoundQuery := `
		query GetNode {
			node(id: "1") {
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

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "node"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.evmORM.On("Node", nodeID).Return(types.Node{
					ID:      nodeID,
					Name:    "node-name",
					WSURL:   null.StringFrom("ws://some-url"),
					HTTPURL: null.StringFrom("http://some-url"),
				}, nil)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
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
				f.Mocks.evmORM.On("Node", int32(1)).Return(types.Node{}, sql.ErrNoRows)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
			},
			query: notFoundQuery,
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
				f.Mocks.evmORM.On("CreateNode", createNodeInput).Return(types.Node{
					ID:         int32(1),
					Name:       createNodeInput.Name,
					EVMChainID: createNodeInput.EVMChainID,
					WSURL:      createNodeInput.WSURL,
					HTTPURL:    createNodeInput.HTTPURL,
					SendOnly:   createNodeInput.SendOnly,
				}, nil)
				f.Mocks.evmORM.On("GetChainsByIDs", []utils.Big{createNodeInput.EVMChainID}).Return([]types.Chain{
					{ID: *utils.NewBigI(1), Enabled: true},
				}, nil)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
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
	fakeNode := types.Node{
		ID:         fakeID,
		Name:       "node-name",
		EVMChainID: *utils.NewBigI(1),
		WSURL:      null.StringFrom("ws://some-url"),
		HTTPURL:    null.StringFrom("http://some-url"),
		SendOnly:   true,
	}

	variables := map[string]interface{}{
		"id": fakeID,
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
				f.Mocks.evmORM.On("Node", fakeID).Return(fakeNode, nil)
				f.Mocks.evmORM.On("DeleteNode", int64(2)).Return(nil)
				f.App.On("EVMORM").Return(f.Mocks.evmORM)
			},
			query:     mutation,
			variables: variables,
			result:    expected,
		},
		{
			name:          "not found error on fetch",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.evmORM.On("Node", fakeID).Return(types.Node{}, sql.ErrNoRows)
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
				f.Mocks.evmORM.On("Node", fakeID).Return(fakeNode, nil)
				f.Mocks.evmORM.On("DeleteNode", int64(2)).Return(evm.ErrNoRowsAffected)
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
