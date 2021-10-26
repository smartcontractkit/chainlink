package resolver

import (
	"database/sql"
	"testing"

	"github.com/graph-gophers/graphql-go/gqltesting"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
	feedsMocks "github.com/smartcontractkit/chainlink/core/services/feeds/mocks"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
)

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
						"jobTypes": ["fluxmonitor"],
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
					"jobTypes": ["fluxmonitor"],
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
