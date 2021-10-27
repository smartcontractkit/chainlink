package resolver

import (
	"database/sql"
	"testing"

	"github.com/graph-gophers/graphql-go/gqltesting"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
	feedsMocks "github.com/smartcontractkit/chainlink/core/services/feeds/mocks"
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
		f        = setupFramework(t)
		feedsSvc = &feedsMocks.Service{}
	)

	pubKey, err := crypto.PublicKeyFromHex("3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808")
	require.NoError(t, err)

	t.Cleanup(func() {
		mock.AssertExpectationsForObjects(t,
			feedsSvc,
		)
	})

	f.App.On("GetFeedsService").Return(feedsSvc)
	feedsSvc.On("ListManagers").Return([]feeds.FeedsManager{
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

	gqltesting.RunTest(t, &gqltesting.Test{
		Context: f.Ctx,
		Schema:  f.RootSchema,
		Query: `
			{
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
					}
				}
			}
		`,
		ExpectedResult: `
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
						"createdAt": "2021-01-01T00:00:00Z"
					}]
				}
			}
		`,
	})
}

func Test_FeedsManager(t *testing.T) {
	var (
		mgrID = int64(1)
		query = `
		{
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

	testCases := []struct {
		name   string
		before func(*gqlTestFramework, *feedsMocks.Service)
		query  string
		result string
	}{
		{
			name: "success",
			before: func(f *gqlTestFramework, feedsSvc *feedsMocks.Service) {
				f.App.On("GetFeedsService").Return(feedsSvc)
				feedsSvc.On("GetManager", mgrID).Return(&feeds.FeedsManager{
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
			name: "not found",
			before: func(f *gqlTestFramework, feedsSvc *feedsMocks.Service) {
				f.App.On("GetFeedsService").Return(feedsSvc)
				feedsSvc.On("GetManager", mgrID).Return(nil, sql.ErrNoRows)
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

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var (
				f        = setupFramework(t)
				feedsSvc = &feedsMocks.Service{}
			)

			t.Cleanup(func() {
				mock.AssertExpectationsForObjects(t,
					feedsSvc,
				)
			})

			if tc.before != nil {
				tc.before(f, feedsSvc)
			}

			gqltesting.RunTest(t, &gqltesting.Test{
				Context:        f.Ctx,
				Schema:         f.RootSchema,
				Query:          tc.query,
				ExpectedResult: tc.result,
			})
		})
	}
}

func Test_CreateFeedsManager(t *testing.T) {
	var (
		mgrID     = int64(1)
		name      = "manager1"
		uri       = "localhost:2000"
		pubKeyHex = "3b0f149627adb7b6fafe1497a9dfc357f22295a5440786c3bc566dfdb0176808"

		query = `
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

	testCases := []struct {
		name      string
		before    func(*gqlTestFramework, *feedsMocks.Service)
		query     string
		variables map[string]interface{}
		result    string
	}{
		{
			name: "success",
			before: func(f *gqlTestFramework, feedsSvc *feedsMocks.Service) {
				f.App.On("GetFeedsService").Return(feedsSvc)
				feedsSvc.On("RegisterManager", &feeds.FeedsManager{
					Name:                      name,
					URI:                       uri,
					PublicKey:                 *pubKey,
					JobTypes:                  pq.StringArray([]string{"fluxmonitor"}),
					IsOCRBootstrapPeer:        false,
					OCRBootstrapPeerMultiaddr: null.StringFromPtr(nil),
				}).Return(mgrID, nil)
				feedsSvc.On("GetManager", mgrID).Return(&feeds.FeedsManager{
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
			query:     query,
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
			name: "single feeds manager error",
			before: func(f *gqlTestFramework, feedsSvc *feedsMocks.Service) {
				f.App.On("GetFeedsService").Return(feedsSvc)
				feedsSvc.
					On("RegisterManager", mock.IsType(&feeds.FeedsManager{})).
					Return(int64(0), feeds.ErrSingleFeedsManager)
			},
			query:     query,
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
			name: "not found",
			before: func(f *gqlTestFramework, feedsSvc *feedsMocks.Service) {
				f.App.On("GetFeedsService").Return(feedsSvc)
				feedsSvc.On("RegisterManager", mock.IsType(&feeds.FeedsManager{})).Return(mgrID, nil)
				feedsSvc.On("GetManager", mgrID).Return(nil, sql.ErrNoRows)
			},
			query:     query,
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
			name:  "invalid input public key",
			query: query,
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

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var (
				f        = setupFramework(t)
				feedsSvc = &feedsMocks.Service{}
			)

			t.Cleanup(func() {
				mock.AssertExpectationsForObjects(t,
					feedsSvc,
				)
			})

			if tc.before != nil {
				tc.before(f, feedsSvc)
			}

			gqltesting.RunTest(t, &gqltesting.Test{
				Context:        f.Ctx,
				Schema:         f.RootSchema,
				Query:          tc.query,
				Variables:      tc.variables,
				ExpectedResult: tc.result,
			})
		})
	}
}
