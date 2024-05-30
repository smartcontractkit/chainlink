package resolver

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func Test_Bridges(t *testing.T) {
	t.Parallel()

	var (
		query = `
			query GetBridges {
				bridges {
					results {
						id
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
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("BridgeORM").Return(f.Mocks.bridgeORM)
				f.Mocks.bridgeORM.On("BridgeTypes", mock.Anything, PageDefaultOffset, PageDefaultLimit).Return([]bridges.BridgeType{
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
						"id": "bridge1",
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
	t.Parallel()

	var (
		query = `
			query GetBridge{
				bridge(id: "bridge1") {
					... on Bridge {
						id
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

		name = bridges.BridgeName("bridge1")
	)
	bridgeURL, err := url.Parse("https://external.adapter")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "bridge"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("BridgeORM").Return(f.Mocks.bridgeORM)
				f.Mocks.bridgeORM.On("FindBridge", mock.Anything, name).Return(bridges.BridgeType{
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
					"id": "bridge1",
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
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("BridgeORM").Return(f.Mocks.bridgeORM)
				f.Mocks.bridgeORM.On("FindBridge", mock.Anything, name).Return(bridges.BridgeType{}, sql.ErrNoRows)
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
		name     = bridges.BridgeName("bridge1")
		mutation = `
			mutation createBridge($input: CreateBridgeInput!) {
				createBridge(input: $input) {
					... on CreateBridgeSuccess {
						bridge {
							id
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
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("BridgeORM").Return(f.Mocks.bridgeORM)
				f.Mocks.bridgeORM.On("FindBridge", mock.Anything, name).Return(bridges.BridgeType{}, sql.ErrNoRows)
				f.Mocks.bridgeORM.On("CreateBridgeType", mock.Anything, mock.IsType(&bridges.BridgeType{})).
					Run(func(args mock.Arguments) {
						arg := args.Get(1).(*bridges.BridgeType)
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
							"id": "bridge1",
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
	t.Parallel()

	var (
		name     = bridges.BridgeName("bridge1")
		mutation = `
			mutation updateBridge($id: ID!, $input: UpdateBridgeInput!) {
				updateBridge(id: $id, input: $input) {
					... on UpdateBridgeSuccess {
						bridge {
							id
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
			"id": "bridge1",
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
			before: func(ctx context.Context, f *gqlTestFramework) {
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
				f.Mocks.bridgeORM.On("FindBridge", mock.Anything, name).Return(bridge, nil)

				btr := &bridges.BridgeTypeRequest{
					Name:                   bridges.BridgeName("bridge-updated"),
					URL:                    models.WebURL(*newBridgeURL),
					Confirmations:          2,
					MinimumContractPayment: assets.NewLinkFromJuels(2),
				}

				f.Mocks.bridgeORM.On("UpdateBridgeType", mock.Anything, mock.IsType(&bridges.BridgeType{}), btr).
					Run(func(args mock.Arguments) {
						arg := args.Get(1).(*bridges.BridgeType)
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
						"id": "bridge-updated",
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
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("BridgeORM").Return(f.Mocks.bridgeORM)
				f.Mocks.bridgeORM.On("FindBridge", mock.Anything, name).Return(bridges.BridgeType{}, sql.ErrNoRows)
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

func Test_DeleteBridgeMutation(t *testing.T) {
	t.Parallel()

	name := bridges.BridgeName("bridge1")

	bridgeURL, err := url.Parse("https://test-url.com")
	require.NoError(t, err)

	link := assets.Link{}
	err = json.Unmarshal([]byte(`"1"`), &link)
	assert.NoError(t, err)

	mutation := `
		mutation DeleteBridge($id: ID!) {
			deleteBridge(id: $id) {
				... on DeleteBridgeSuccess {
					bridge {
						id
						name
						url
						confirmations
						outgoingToken
						minimumContractPayment
					}
				}
				... on NotFoundError {
					message
					code
				}
				... on DeleteBridgeInvalidNameError {
					message
					code
				}
				... on DeleteBridgeConflictError {
					message
					code
				}
			}
		}`

	variables := map[string]interface{}{
		"id": name.String(),
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "deleteBridge"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				bridge := bridges.BridgeType{
					Name:                   name,
					URL:                    models.WebURL(*bridgeURL),
					Confirmations:          1,
					OutgoingToken:          "some-token",
					MinimumContractPayment: &link,
				}

				f.Mocks.bridgeORM.On("FindBridge", mock.Anything, name).Return(bridge, nil)
				f.Mocks.bridgeORM.On("DeleteBridgeType", mock.Anything, &bridge).Return(nil)
				f.Mocks.jobORM.On("FindJobIDsWithBridge", mock.Anything, name.String()).Return([]int32{}, nil)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.App.On("BridgeORM").Return(f.Mocks.bridgeORM)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"deleteBridge": {
						"bridge": {
							"id": "bridge1",
							"name": "bridge1",
							"url": "https://test-url.com",
							"confirmations": 1,
							"outgoingToken": "some-token",
							"minimumContractPayment": "1"
						}
					}
				}`,
		},
		{
			name:          "invalid bridge type name",
			authenticated: true,
			query:         mutation,
			variables: map[string]interface{}{
				"id": "][]$$$$324adfas",
			},
			result: `
				{
					"deleteBridge": {
						"message": "task type validation: name ][]$$$$324adfas contains invalid characters",
						"code": "UNPROCESSABLE"
					}
				}`,
		},
		{
			name:          "invalid bridge type name",
			authenticated: true,
			query:         mutation,
			variables: map[string]interface{}{
				"id": "bridge1",
			},
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.bridgeORM.On("FindBridge", mock.Anything, name).Return(bridges.BridgeType{}, sql.ErrNoRows)
				f.App.On("BridgeORM").Return(f.Mocks.bridgeORM)
			},
			result: `
				{
					"deleteBridge": {
						"message": "bridge not found",
						"code": "NOT_FOUND"
					}
				}`,
		},
		{
			name:          "bridge with jobs associated",
			authenticated: true,
			query:         mutation,
			variables: map[string]interface{}{
				"id": "bridge1",
			},
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.bridgeORM.On("FindBridge", mock.Anything, name).Return(bridges.BridgeType{}, nil)
				f.Mocks.jobORM.On("FindJobIDsWithBridge", mock.Anything, name.String()).Return([]int32{1}, nil)
				f.App.On("BridgeORM").Return(f.Mocks.bridgeORM)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
			},
			result: `
				{
					"deleteBridge": {
						"message": "bridge has jobs associated with it",
						"code": "UNPROCESSABLE"
					}
				}`,
		},
	}

	RunGQLTests(t, testCases)
}
