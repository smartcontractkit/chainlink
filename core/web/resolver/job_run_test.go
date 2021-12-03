package resolver

import (
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestQuery_PaginatedJobsRuns(t *testing.T) {
	t.Parallel()

	query := `
		query GetJobsRuns {
			jobsRuns {
				results {
					id
				}
				metadata {
					total
				}
			}
		}`

	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "jobsRuns"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.jobORM.On("PipelineRuns", (*int32)(nil), PageDefaultOffset, PageDefaultLimit).Return([]pipeline.Run{
					{
						ID: int64(200),
					},
				}, 1, nil)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
			},
			query: query,
			result: `
				{
					"jobsRuns": {
						"results": [{
							"id": "200"
						}],
						"metadata": {
							"total": 1
						}
					}
				}`,
		},
		{
			name:          "generic error on PipelineRuns()",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.jobORM.On("PipelineRuns", (*int32)(nil), PageDefaultOffset, PageDefaultLimit).Return(nil, 0, gError)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"jobsRuns"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
