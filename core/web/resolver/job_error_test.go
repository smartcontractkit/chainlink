package resolver

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/job"
)

// JobErrors are only embedded on the job and are not fetchable by it's own id,
// so we test the job error resolvers by fetching a job by id.

func TestResolver_JobErrors(t *testing.T) {
	var (
		id      = int32(1)
		errorID = int64(200)
	)

	testCases := []GQLTestCase{
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobTx", id).Return(job.Job{
					JobSpecErrors: []job.SpecError{
						{
							ID:          errorID,
							Description: "no contract code at given address",
							Occurrences: 1,
							CreatedAt:   f.Timestamp(),
						},
					},
				}, nil)
			},
			query: `
				query GetJob {
					job(id: "1") {
						... on Job {
							errors {
								id
								description
								occurrences
								createdAt
							}
						}
					}
				}
			`,
			result: `
				{
					"job": {
						"errors": [{
							"id": "200",
							"description": "no contract code at given address",
							"occurrences": 1,
							"createdAt": "2021-01-01T00:00:00Z"
						}]
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}
