package resolver

import (
	"database/sql"
	"testing"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_NodeQuery(t *testing.T) {
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
			name:          "not found",
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
