package resolver

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func Test_Jobs(t *testing.T) {
	var (
		externalJobID = uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001"))

		query = `
			query GetJobs {
				jobs {
					results {
						id
						name
						schemaVersion
						maxTaskDuration
						externalJobID
						createdAt
					}
				}
			}`
	)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "jobs"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobs", 0, 50).Return([]job.Job{
					{
						ID:              1,
						Name:            null.StringFrom("job1"),
						SchemaVersion:   1,
						MaxTaskDuration: models.Interval(1 * time.Second),
						ExternalJobID:   externalJobID,
						CreatedAt:       f.Timestamp(),
					},
				}, 1, nil)
			},
			query: query,
			result: `
				{
					"jobs": {
						"results": [{
							"id": "1",
							"name": "job1",
							"schemaVersion": 1,
							"maxTaskDuration": "1s",
							"externalJobID": "00000000-0000-0000-0000-000000000001",
							"createdAt": "2021-01-01T00:00:00Z"
						}]
					}
				}`,
		},
	}

	RunGQLTests(t, testCases)
}
