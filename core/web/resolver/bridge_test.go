package resolver

import (
	"database/sql"
	"net/url"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func Test_Bridges(t *testing.T) {
	t.Parallel()

	var (
		query = `
			query GetBridges {
				bridges {
					results {
						name
						url
						confirmations
						outgoingToken
						minimumContractPayment
						createdAt
					}
					metadata {
						total
					}
				}
			}`
	)

	bridgeURL, err := url.Parse("https://external.adapter")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "bridges"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("BridgeORM").Return(f.Mocks.bridgeORM)
				f.Mocks.bridgeORM.On("BridgeTypes", PageDefaultOffset, PageDefaultLimit).Return([]bridges.BridgeType{
					{
						Name:                   "bridge1",
						URL:                    models.WebURL(*bridgeURL),
						Confirmations:          uint32(1),
						OutgoingToken:          "outgoingToken",
						MinimumContractPayment: assets.NewLinkFromJuels(1),
						CreatedAt:              f.Timestamp(),
					},
				}, 1, nil)
			},
			query: query,
			result: `
			{
				"bridges": {
					"results": [{
						"name": "bridge1",
						"url": "https://external.adapter",
						"confirmations": 1,
						"outgoingToken": "outgoingToken",
						"minimumContractPayment": "1",
						"createdAt": "2021-01-01T00:00:00Z"
					}],
					"metadata": {
						"total": 1
					}
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}

func Test_Bridge(t *testing.T) {
	var (
		query = `
			query GetBridge{
				bridge(name: "bridge1") {
					... on Bridge {
						name
						url
						confirmations
						outgoingToken
						minimumContractPayment
						createdAt
					}
					... on NotFoundError {
						message
						code
					}
				}
			}`

		name = bridges.TaskType("bridge1")
	)
	bridgeURL, err := url.Parse("https://external.adapter")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "bridge"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("BridgeORM").Return(f.Mocks.bridgeORM)
				f.Mocks.bridgeORM.On("FindBridge", name).Return(bridges.BridgeType{
					Name:                   name,
					URL:                    models.WebURL(*bridgeURL),
					Confirmations:          uint32(1),
					OutgoingToken:          "outgoingToken",
					MinimumContractPayment: assets.NewLinkFromJuels(1),
					CreatedAt:              f.Timestamp(),
				}, nil)
			},
			query: query,
			result: `{
				"bridge": {
					"name": "bridge1",
					"url": "https://external.adapter",
					"confirmations": 1,
					"outgoingToken": "outgoingToken",
					"minimumContractPayment": "1",
					"createdAt": "2021-01-01T00:00:00Z"
				}
			}`,
		},
		{
			name:          "not found",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("BridgeORM").Return(f.Mocks.bridgeORM)
				f.Mocks.bridgeORM.On("FindBridge", name).Return(bridges.BridgeType{}, sql.ErrNoRows)
			},
			query: query,
			result: `{
				"bridge": {
					"message": "bridge not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}

func Test_CreateBridge(t *testing.T) {
	t.Parallel()

	var (
		name     = bridges.TaskType("bridge1")
		mutation = `
			mutation createBridge($input: CreateBridgeInput!) {
				createBridge(input: $input) {
					... on CreateBridgeSuccess {
						bridge {
							name
							url
							confirmations
							outgoingToken
							minimumContractPayment
							createdAt
						}
					}
				}
			}`
		variables = map[string]interface{}{
			"input": map[string]interface{}{
				"name":                   "bridge1",
				"url":                    "https://external.adapter",
				"confirmations":          1,
				"minimumContractPayment": "1",
			},
		}
	)
	bridgeURL, err := url.Parse("https://external.adapter")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "createBridge"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("BridgeORM").Return(f.Mocks.bridgeORM)
				f.Mocks.bridgeORM.On("FindBridge", name).Return(bridges.BridgeType{}, sql.ErrNoRows)
				f.Mocks.bridgeORM.On("CreateBridgeType", mock.IsType(&bridges.BridgeType{})).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*bridges.BridgeType)
						*arg = bridges.BridgeType{
							Name:                   name,
							URL:                    models.WebURL(*bridgeURL),
							Confirmations:          uint32(1),
							OutgoingToken:          "outgoingToken",
							MinimumContractPayment: assets.NewLinkFromJuels(1),
							CreatedAt:              f.Timestamp(),
						}
					}).
					Return(nil)
			},
			query:     mutation,
			variables: variables,
			// We should test equality for the generated token but since it is
			// generated by a non mockable object, we can't do this right now.
			result: `
				{
					"createBridge": {
						"bridge": {
							"name": "bridge1",
							"url": "https://external.adapter",
							"confirmations": 1,
							"outgoingToken": "outgoingToken",
							"minimumContractPayment": "1",
							"createdAt": "2021-01-01T00:00:00Z"
						}
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}

func Test_UpdateBridge(t *testing.T) {
	var (
		name     = bridges.TaskType("bridge1")
		mutation = `
			mutation updateBridge($input: UpdateBridgeInput!) {
				updateBridge(name: "bridge1", input: $input) {
					... on UpdateBridgeSuccess {
						bridge {
							name
							url
							confirmations
							outgoingToken
							minimumContractPayment
							createdAt
						}
					}
					... on NotFoundError {
						message
						code
					}
				}
			}`
		variables = map[string]interface{}{
			"input": map[string]interface{}{
				"name":                   "bridge-updated",
				"url":                    "https://external.adapter.new",
				"confirmations":          2,
				"minimumContractPayment": "2",
			},
		}
	)
	bridgeURL, err := url.Parse("https://external.adapter")
	require.NoError(t, err)

	newBridgeURL, err := url.Parse("https://external.adapter.new")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "updateBridge"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				// Initialize the existing bridge
				bridge := bridges.BridgeType{
					Name:                   name,
					URL:                    models.WebURL(*bridgeURL),
					Confirmations:          uint32(1),
					OutgoingToken:          "outgoingToken",
					MinimumContractPayment: assets.NewLinkFromJuels(1),
					CreatedAt:              f.Timestamp(),
				}

				f.App.On("BridgeORM").Return(f.Mocks.bridgeORM)
				f.Mocks.bridgeORM.On("FindBridge", name).Return(bridge, nil)

				btr := &bridges.BridgeTypeRequest{
					Name:                   bridges.TaskType("bridge-updated"),
					URL:                    models.WebURL(*newBridgeURL),
					Confirmations:          2,
					MinimumContractPayment: assets.NewLinkFromJuels(2),
				}

				f.Mocks.bridgeORM.On("UpdateBridgeType", mock.IsType(&bridges.BridgeType{}), btr).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*bridges.BridgeType)
						*arg = bridges.BridgeType{
							Name:                   "bridge-updated",
							URL:                    models.WebURL(*newBridgeURL),
							Confirmations:          2,
							OutgoingToken:          "outgoingToken",
							MinimumContractPayment: assets.NewLinkFromJuels(2),
							CreatedAt:              f.Timestamp(),
						}
					}).
					Return(nil)
			},
			query:     mutation,
			variables: variables,
			result: `{
				"updateBridge": {
					"bridge": {
						"name": "bridge-updated",
						"url": "https://external.adapter.new",
						"confirmations": 2,
						"outgoingToken": "outgoingToken",
						"minimumContractPayment": "2",
						"createdAt": "2021-01-01T00:00:00Z"
					}
				}
			}`,
		},
		{
			name:          "not found",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("BridgeORM").Return(f.Mocks.bridgeORM)
				f.Mocks.bridgeORM.On("FindBridge", name).Return(bridges.BridgeType{}, sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `{
				"updateBridge": {
					"message": "bridge not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}
