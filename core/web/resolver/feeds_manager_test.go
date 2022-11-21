package resolver

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
)

func Test_FeedsManagers(t *testing.T) {
	var (
		query = `
			query GetFeedsManagers {
				feedsManagers {
					results {
						id
						name
						uri
						publicKey
						isConnectionActive
						createdAt
						jobProposals {
							id
							status
						}
					}
				}
			}`
	)

	pubKey, err := crypto.PublicKeyFromHex("3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "feedsManagers"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("ListJobProposalsByManagersIDs", []int64{1}).Return([]feeds.JobProposal{
					{
						ID:             int64(100),
						FeedsManagerID: int64(1),
						Status:         feeds.JobProposalStatusApproved,
					},
				}, nil)
				f.Mocks.feedsSvc.On("ListManagers").Return([]feeds.FeedsManager{
					{
						ID:                 1,
						Name:               "manager1",
						URI:                "localhost:2000",
						PublicKey:          *pubKey,
						IsConnectionActive: true,
						CreatedAt:          f.Timestamp(),
					},
				}, nil)
			},
			query: query,
			result: `
				{
					"feedsManagers": {
						"results": [{
							"id": "1",
							"name": "manager1",
							"uri": "localhost:2000",
							"publicKey": "3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808",
							"isConnectionActive": true,
							"createdAt": "2021-01-01T00:00:00Z",
							"jobProposals": [{
								"id": "100",
								"status": "APPROVED"
							}]
						}]
					}
				}`,
		},
	}

	RunGQLTests(t, testCases)
}

func Test_FeedsManager(t *testing.T) {
	var (
		mgrID = int64(1)
		query = `
			query GetFeedsManager {
				feedsManager(id: 1) {
					... on FeedsManager {
						id
						name
						uri
						publicKey
						isConnectionActive
						createdAt
					}
					... on NotFoundError {
						message
						code
					}
				}
			}`
	)
	pubKey, err := crypto.PublicKeyFromHex("3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808")
	require.NoError(t, err)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "feedsManager"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("GetManager", mgrID).Return(&feeds.FeedsManager{
					ID:                 mgrID,
					Name:               "manager1",
					URI:                "localhost:2000",
					PublicKey:          *pubKey,
					IsConnectionActive: true,
					CreatedAt:          f.Timestamp(),
				}, nil)
			},
			query: query,
			result: `
			{
				"feedsManager": {
					"id": "1",
					"name": "manager1",
					"uri": "localhost:2000",
					"publicKey": "3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808",
					"isConnectionActive": true,
					"createdAt": "2021-01-01T00:00:00Z"
				}
			}`,
		},
		{
			name:          "not found",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("GetManager", mgrID).Return(nil, sql.ErrNoRows)
			},
			query: query,
			result: `
			{
				"feedsManager": {
					"message": "feeds manager not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}

func Test_CreateFeedsManager(t *testing.T) {
	var (
		mgrID     = int64(1)
		name      = "manager1"
		uri       = "localhost:2000"
		pubKeyHex = "3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808"

		mutation = `
			mutation CreateFeedsManager($input: CreateFeedsManagerInput!) {
				createFeedsManager(input: $input) {
					... on CreateFeedsManagerSuccess {
						feedsManager {
							id
							name
							uri
							publicKey
							isConnectionActive
							createdAt
						}
					}
					... on SingleFeedsManagerError {
						message
						code
					}
					... on NotFoundError {
						message
						code
					}
					... on InputErrors {
						errors {
							path
							message
							code
						}
					}
				}
			}`
		variables = map[string]interface{}{
			"input": map[string]interface{}{
				"name":            name,
				"uri":             uri,
				"jobTypes":        []interface{}{"FLUX_MONITOR", "OCR2"},
				"publicKey":       pubKeyHex,
				"isBootstrapPeer": false,
			},
		}
	)
	pubKey, err := crypto.PublicKeyFromHex(pubKeyHex)
	require.NoError(t, err)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "createFeedsManager"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("RegisterManager", mock.Anything, feeds.RegisterManagerParams{
					Name:      name,
					URI:       uri,
					PublicKey: *pubKey,
				}).Return(mgrID, nil)
				f.Mocks.feedsSvc.On("GetManager", mgrID).Return(&feeds.FeedsManager{
					ID:                 mgrID,
					Name:               name,
					URI:                uri,
					PublicKey:          *pubKey,
					IsConnectionActive: false,
					CreatedAt:          f.Timestamp(),
				}, nil)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"createFeedsManager": {
					"feedsManager": {
						"id": "1",
						"name": "manager1",
						"uri": "localhost:2000",
						"publicKey": "3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808",
						"isConnectionActive": false,
						"createdAt": "2021-01-01T00:00:00Z"
					}
				}
			}`,
		},
		{
			name:          "single feeds manager error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.
					On("RegisterManager", mock.Anything, mock.IsType(feeds.RegisterManagerParams{})).
					Return(int64(0), feeds.ErrSingleFeedsManager)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"createFeedsManager": {
					"message": "only a single feeds manager is supported",
					"code": "UNPROCESSABLE"
				}
			}`,
		},
		{
			name:          "not found",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("RegisterManager", mock.Anything, mock.IsType(feeds.RegisterManagerParams{})).Return(mgrID, nil)
				f.Mocks.feedsSvc.On("GetManager", mgrID).Return(nil, sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"createFeedsManager": {
					"message": "feeds manager not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
		{
			name:          "invalid input public key",
			authenticated: true,
			query:         mutation,
			variables: map[string]interface{}{
				"input": map[string]interface{}{
					"name":      name,
					"uri":       uri,
					"publicKey": "zzzzz",
				},
			},
			result: `
			{
				"createFeedsManager": {
					"errors": [{
						"path": "input/publicKey",
						"message": "invalid hex value",
						"code": "INVALID_INPUT"
					}]
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}

func Test_UpdateFeedsManager(t *testing.T) {
	var (
		mgrID     = int64(1)
		name      = "manager1"
		uri       = "localhost:2000"
		pubKeyHex = "3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808"

		mutation = `
			mutation UpdateFeedsManager($id: ID!, $input: UpdateFeedsManagerInput!) {
				updateFeedsManager(id: $id, input: $input) {
					... on UpdateFeedsManagerSuccess {
						feedsManager {
							id
							name
							uri
							publicKey
							isConnectionActive
							createdAt
						}
					}
					... on NotFoundError {
						message
						code
					}
					... on InputErrors {
						errors {
							path
							message
							code
						}
					}
				}
			}`
		variables = map[string]interface{}{
			"id": "1",
			"input": map[string]interface{}{
				"name":      name,
				"uri":       uri,
				"publicKey": pubKeyHex,
			},
		}
	)
	pubKey, err := crypto.PublicKeyFromHex(pubKeyHex)
	require.NoError(t, err)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "updateFeedsManager"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("UpdateManager", mock.Anything, feeds.FeedsManager{
					ID:        mgrID,
					Name:      name,
					URI:       uri,
					PublicKey: *pubKey,
				}).Return(nil)
				f.Mocks.feedsSvc.On("GetManager", mgrID).Return(&feeds.FeedsManager{
					ID:                 mgrID,
					Name:               name,
					URI:                uri,
					PublicKey:          *pubKey,
					IsConnectionActive: false,
					CreatedAt:          f.Timestamp(),
				}, nil)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"updateFeedsManager": {
					"feedsManager": {
						"id": "1",
						"name": "manager1",
						"uri": "localhost:2000",
						"publicKey": "3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808",
						"isConnectionActive": false,
						"createdAt": "2021-01-01T00:00:00Z"
					}
				}
			}`,
		},
		{
			name:          "not found",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("UpdateManager", mock.Anything, mock.IsType(feeds.FeedsManager{})).Return(nil)
				f.Mocks.feedsSvc.On("GetManager", mgrID).Return(nil, sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"updateFeedsManager": {
					"message": "feeds manager not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
		{
			name:          "invalid input public key",
			authenticated: true,
			query:         mutation,
			variables: map[string]interface{}{
				"id": "1",
				"input": map[string]interface{}{
					"name":      name,
					"uri":       uri,
					"publicKey": "zzzzz",
				},
			},
			result: `
			{
				"updateFeedsManager": {
					"errors": [{
						"path": "input/publicKey",
						"message": "invalid hex value",
						"code": "INVALID_INPUT"
					}]
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}
