package resolver

import (
	"database/sql"
	"testing"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
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
