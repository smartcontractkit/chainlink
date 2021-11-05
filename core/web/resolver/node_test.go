package resolver

import (
	"testing"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
)

func Test_NodeQuery(t *testing.T) {
	query := `
		query GetNode {
			node(id: "200") {
				name
				wsURL
				httpURL
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
	}

	RunGQLTests(t, testCases)
}
