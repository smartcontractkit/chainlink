package resolver

import (
	"database/sql"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// This tests the main fields on the job results. Embedded spec testing is done
// in the `spec_test` file
func TestResolver_Jobs(t *testing.T) {
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
						spec {
							__typename
						}
					}
					metadata {
						total
					}
				}
			}`
	)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "jobs"),
		{
			name:          "get jobs success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobs", 0, 50).Return([]job.Job{
					{
						ID:                          1,
						Name:                        null.StringFrom("job1"),
						SchemaVersion:               1,
						MaxTaskDuration:             models.Interval(1 * time.Second),
						ExternalJobID:               externalJobID,
						CreatedAt:                   f.Timestamp(),
						Type:                        job.OffchainReporting,
						OffchainreportingOracleSpec: &job.OffchainReportingOracleSpec{},
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
							"createdAt": "2021-01-01T00:00:00Z",
							"spec": {
								"__typename": "OCRSpec"
							}
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

func TestResolver_Job(t *testing.T) {
	var (
		id            = int32(1)
		externalJobID = uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001"))

		query = `
			query GetJob {
				job(id: "1") {
					... on Job {
						id
						name
						schemaVersion
						maxTaskDuration
						externalJobID
						createdAt
						spec {
							__typename
						}
					}
					... on NotFoundError {
						message
						code
					}
				}
			}
		`
	)

	testCases := []GQLTestCase{
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobTx", id).Return(job.Job{
					ID:                          1,
					Name:                        null.StringFrom("job1"),
					SchemaVersion:               1,
					MaxTaskDuration:             models.Interval(1 * time.Second),
					ExternalJobID:               externalJobID,
					CreatedAt:                   f.Timestamp(),
					Type:                        job.OffchainReporting,
					OffchainreportingOracleSpec: &job.OffchainReportingOracleSpec{},
				}, nil)
			},
			query: query,
			result: `
				{
					"job": {
						"id": "1",
						"name": "job1",
						"schemaVersion": 1,
						"maxTaskDuration": "1s",
						"externalJobID": "00000000-0000-0000-0000-000000000001",
						"createdAt": "2021-01-01T00:00:00Z",
						"spec": {
							"__typename": "OCRSpec"
						}
					}
				}
			`,
		},
		{
			name:          "not found",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobTx", id).Return(job.Job{}, sql.ErrNoRows)
			},
			query: query,
			result: `
				{
					"job": {
						"code": "NOT_FOUND",
						"message": "job not found"
					}
				}
			`,
		},
	}

	RunGQLTests(t, testCases)
}
