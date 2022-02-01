package resolver

import (
	"database/sql"
	"fmt"
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/services/feeds"
)

func TestResolver_GetJobProposal(t *testing.T) {
	t.Parallel()

	query := `
		query GetJobProposal {
			jobProposal(id: "1") {
				... on JobProposal {
					id
					status
					externalJobID
					remoteUUID
					multiAddrs
					pendingUpdate
					feedsManager {
						id
						name
					}
				}
				... on NotFoundError {
					message
					code
				}
			}
		}`

	jpID := int64(1)
	ejID := uuid.NullUUID{UUID: uuid.NewV4(), Valid: true}
	rUUID := uuid.NewV4()
	result := `
		{
			"jobProposal": {
				"id": "1",
				"status": "APPROVED",
				"externalJobID": "%s",
				"remoteUUID": "%s",
				"multiAddrs": ["1", "2"],
				"pendingUpdate": false,
				"feedsManager": {
					"id": "1",
					"name": "manager"
				}
			}
		}`

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "jobProposal"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("ListManagersByIDs", []int64{1}).Return([]feeds.FeedsManager{
					{
						ID:   1,
						Name: "manager",
					},
				}, nil)
				f.Mocks.feedsSvc.On("GetJobProposal", jpID).Return(&feeds.JobProposal{
					ID:             jpID,
					Status:         feeds.JobProposalStatusApproved,
					ExternalJobID:  ejID,
					RemoteUUID:     rUUID,
					FeedsManagerID: 1,
					Multiaddrs:     []string{"1", "2"},
					PendingUpdate:  false,
				}, nil)
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
			},
			query:  query,
			result: fmt.Sprintf(result, ejID.UUID.String(), rUUID.String()),
		},
		{
			name:          "not found error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("GetJobProposal", jpID).Return(nil, sql.ErrNoRows)
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
			},
			query: query,
			result: `
			{
				"jobProposal": {
					"message": "job proposal not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}
