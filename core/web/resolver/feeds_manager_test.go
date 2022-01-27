package resolver

import (
	"database/sql"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
)

func Test_ToJobType(t *testing.T) {
	t.Parallel()

	jt, err := ToJobType("fluxmonitor")
	require.NoError(t, err)
	assert.Equal(t, jt, JobTypeFluxMonitor)

	jt, err = ToJobType("ocr")
	require.NoError(t, err)
	assert.Equal(t, jt, JobTypeOCR)

	_, err = ToJobType("xxx")
	require.Error(t, err)
	assert.EqualError(t, err, "invalid job type")
}

func Test_FromJobType(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "fluxmonitor", FromJobTypeInput(JobTypeFluxMonitor))
	assert.Equal(t, "ocr", FromJobTypeInput(JobTypeOCR))
}

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
						jobTypes
						isBootstrapPeer
						bootstrapPeerMultiaddr
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
						ID:                        1,
						Name:                      "manager1",
						URI:                       "localhost:2000",
						PublicKey:                 *pubKey,
						JobTypes:                  []string{"fluxmonitor"},
						IsOCRBootstrapPeer:        true,
						OCRBootstrapPeerMultiaddr: null.StringFrom("/dns4/ocr-bootstrap.chain.link/tcp/0000/p2p/7777777"),
						IsConnectionActive:        true,
						CreatedAt:                 f.Timestamp(),
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
							"jobTypes": ["FLUX_MONITOR"],
							"isBootstrapPeer": true,
							"bootstrapPeerMultiaddr": "/dns4/ocr-bootstrap.chain.link/tcp/0000/p2p/7777777",
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
						jobTypes
						isBootstrapPeer
						isConnectionActive
						bootstrapPeerMultiaddr
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
					ID:                        mgrID,
					Name:                      "manager1",
					URI:                       "localhost:2000",
					PublicKey:                 *pubKey,
					JobTypes:                  []string{"fluxmonitor"},
					IsOCRBootstrapPeer:        true,
					OCRBootstrapPeerMultiaddr: null.StringFrom("/dns4/ocr-bootstrap.chain.link/tcp/0000/p2p/7777777"),
					IsConnectionActive:        true,
					CreatedAt:                 f.Timestamp(),
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
					"jobTypes": ["FLUX_MONITOR"],
					"isBootstrapPeer": true,
					"bootstrapPeerMultiaddr": "/dns4/ocr-bootstrap.chain.link/tcp/0000/p2p/7777777",
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
							jobTypes
							isBootstrapPeer
							isConnectionActive
							bootstrapPeerMultiaddr
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
				"jobTypes":        []interface{}{"FLUX_MONITOR"},
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
				f.Mocks.feedsSvc.On("RegisterManager", &feeds.FeedsManager{
					Name:                      name,
					URI:                       uri,
					PublicKey:                 *pubKey,
					JobTypes:                  pq.StringArray([]string{"fluxmonitor"}),
					IsOCRBootstrapPeer:        false,
					OCRBootstrapPeerMultiaddr: null.StringFromPtr(nil),
				}).Return(mgrID, nil)
				f.Mocks.feedsSvc.On("GetManager", mgrID).Return(&feeds.FeedsManager{
					ID:                        mgrID,
					Name:                      name,
					URI:                       uri,
					PublicKey:                 *pubKey,
					JobTypes:                  []string{"fluxmonitor"},
					IsOCRBootstrapPeer:        false,
					OCRBootstrapPeerMultiaddr: null.StringFromPtr(nil),
					IsConnectionActive:        false,
					CreatedAt:                 f.Timestamp(),
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
						"jobTypes": ["FLUX_MONITOR"],
						"isBootstrapPeer": false,
						"bootstrapPeerMultiaddr": null,
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
					On("RegisterManager", mock.IsType(&feeds.FeedsManager{})).
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
				f.Mocks.feedsSvc.On("RegisterManager", mock.IsType(&feeds.FeedsManager{})).Return(mgrID, nil)
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
					"name":            name,
					"uri":             uri,
					"jobTypes":        []interface{}{"FLUX_MONITOR"},
					"publicKey":       "zzzzz",
					"isBootstrapPeer": false,
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
							jobTypes
							isBootstrapPeer
							isConnectionActive
							bootstrapPeerMultiaddr
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
				"name":            name,
				"uri":             uri,
				"jobTypes":        []interface{}{"FLUX_MONITOR"},
				"publicKey":       pubKeyHex,
				"isBootstrapPeer": false,
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
					ID:                        mgrID,
					Name:                      name,
					URI:                       uri,
					PublicKey:                 *pubKey,
					JobTypes:                  pq.StringArray([]string{"fluxmonitor"}),
					IsOCRBootstrapPeer:        false,
					OCRBootstrapPeerMultiaddr: null.StringFromPtr(nil),
				}).Return(nil)
				f.Mocks.feedsSvc.On("GetManager", mgrID).Return(&feeds.FeedsManager{
					ID:                        mgrID,
					Name:                      name,
					URI:                       uri,
					PublicKey:                 *pubKey,
					JobTypes:                  []string{"fluxmonitor"},
					IsOCRBootstrapPeer:        false,
					OCRBootstrapPeerMultiaddr: null.StringFromPtr(nil),
					IsConnectionActive:        false,
					CreatedAt:                 f.Timestamp(),
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
						"jobTypes": ["FLUX_MONITOR"],
						"isBootstrapPeer": false,
						"bootstrapPeerMultiaddr": null,
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
					"name":            name,
					"uri":             uri,
					"jobTypes":        []interface{}{"FLUX_MONITOR"},
					"publicKey":       "zzzzz",
					"isBootstrapPeer": false,
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
