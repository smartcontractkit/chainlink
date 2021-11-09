package resolver

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

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
					spec
					status
					externalJobID
					multiAddrs
				}
				... on NotFoundError {
					message
					code
				}
			}
		}`

	jpID := int64(1)
	ejID := uuid.NullUUID{UUID: uuid.NewV4(), Valid: true}
	result := `
		{
			"jobProposal": {
				"id": "1",
				"spec": "some-spec",
				"status": "approved",
				"externalJobID": "%s",
				"multiAddrs": ["1", "2"]
			}
		}`

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "jobProposal"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("GetJobProposal", jpID).Return(&feeds.JobProposal{
					ID:             jpID,
					Spec:           "some-spec",
					Status:         feeds.JobProposalStatusApproved,
					ExternalJobID:  ejID,
					FeedsManagerID: 1,
					Multiaddrs:     []string{"1", "2"},
					ProposedAt:     time.Time{},
				}, nil)
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
			},
			query:  query,
			result: fmt.Sprintf(result, ejID.UUID.String()),
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
