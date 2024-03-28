package resolver

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
)

func TestResolver_GetJobProposal(t *testing.T) {
	t.Parallel()

	query := `
		query GetJobProposal {
			jobProposal(id: "1") {
				... on JobProposal {
					id
					name
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
	ejID := uuid.NullUUID{UUID: uuid.New(), Valid: true}
	rUUID := uuid.New()
	name := "job_proposal_1"
	result := `
		{
			"jobProposal": {
				"id": "1",
				"name": "%s",
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
				f.Mocks.feedsSvc.On("ListManagersByIDs", mock.Anything, []int64{1}).Return([]feeds.FeedsManager{
					{
						ID:   1,
						Name: "manager",
					},
				}, nil)
				f.Mocks.feedsSvc.On("GetJobProposal", mock.Anything, jpID).Return(&feeds.JobProposal{
					ID:             jpID,
					Name:           null.StringFrom(name),
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
			result: fmt.Sprintf(result, name, ejID.UUID.String(), rUUID.String()),
		},
		{
			name:          "not found error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("GetJobProposal", mock.Anything, jpID).Return(nil, sql.ErrNoRows)
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
