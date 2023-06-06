package resolver

import (
	"database/sql"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/webhook"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
)

func TestQuery_PaginatedJobRuns(t *testing.T) {
	t.Parallel()

	query := `
		query GetJobsRuns {
			jobRuns {
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
		unauthorizedTestCase(GQLTestCase{query: query}, "jobRuns"),
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
					"jobRuns": {
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
					Path:          []interface{}{"jobRuns"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_JobRun(t *testing.T) {
	t.Parallel()

	query := `
		query GetJobRun($id: ID!) {
			jobRun(id: $id) {
				... on JobRun {
					id
					allErrors
					createdAt
					fatalErrors
					finishedAt
					inputs
					job {
						id
						name
					}
					outputs
					status
				}
				... on NotFoundError {
					code
					message
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": "2",
	}
	gError := errors.New("error")
	_, idError := stringutils.ToInt64("asdasads")

	inputs := pipeline.JSONSerializable{}
	err := inputs.UnmarshalJSON([]byte(`{"foo": "bar"}`))
	require.NoError(t, err)

	outputs := pipeline.JSONSerializable{}
	err = outputs.UnmarshalJSON([]byte(`[{"baz": "bar"}]`))
	require.NoError(t, err)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query, variables: variables}, "jobRun"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.jobORM.On("FindPipelineRunByID", int64(2)).Return(pipeline.Run{
					ID:             2,
					PipelineSpecID: 5,
					CreatedAt:      f.Timestamp(),
					FinishedAt:     null.TimeFrom(f.Timestamp()),
					AllErrors:      pipeline.RunErrors{null.StringFrom("fatal error"), null.String{}},
					FatalErrors:    pipeline.RunErrors{null.StringFrom("fatal error"), null.String{}},
					Inputs:         inputs,
					Outputs:        outputs,
					State:          pipeline.RunStatusErrored,
				}, nil)
				f.Mocks.jobORM.On("FindJobsByPipelineSpecIDs", []int32{5}).Return([]job.Job{
					{
						ID:             1,
						PipelineSpecID: 2,
						Name:           null.StringFrom("first-one"),
					},
					{
						ID:             2,
						PipelineSpecID: 5,
						Name:           null.StringFrom("second-one"),
					},
				}, nil)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
			},
			query:     query,
			variables: variables,
			result: `
				{
					"jobRun": {
						"id": "2",
						"allErrors": ["fatal error"],
						"createdAt": "2021-01-01T00:00:00Z",
						"fatalErrors": ["fatal error"],
						"finishedAt": "2021-01-01T00:00:00Z",
						"inputs": "{\"foo\":\"bar\"}",
						"job": {
							"id": "2",
							"name": "second-one"
						},
						"outputs": ["{\"baz\":\"bar\"}"],
						"status": "ERRORED"
					}
				}`,
		},
		{
			name:          "not found error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.jobORM.On("FindPipelineRunByID", int64(2)).Return(pipeline.Run{}, sql.ErrNoRows)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
			},
			query:     query,
			variables: variables,
			result: `
				{
					"jobRun": {
						"code": "NOT_FOUND",
						"message": "job run not found"
					}
				}`,
		},
		{
			name:          "generic error on FindPipelineRunByID()",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.jobORM.On("FindPipelineRunByID", int64(2)).Return(pipeline.Run{}, gError)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
			},
			query:     query,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"jobRun"},
					Message:       gError.Error(),
				},
			},
		},
		{
			name:          "invalid ID error",
			authenticated: true,
			query:         query,
			variables: map[string]interface{}{
				"id": "asdasads",
			},
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: idError,
					Path:          []interface{}{"jobRun"},
					Message:       idError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_RunJob(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation RunJob($id: ID!) {
			runJob(id: $id) {
				... on RunJobSuccess {
					jobRun {
						id
						allErrors
						createdAt
						fatalErrors
						finishedAt
						inputs
						outputs
						status
					}
				}
				... on RunJobCannotRunError {
					code
					message
				}
				... on NotFoundError {
					code
					message
				}
			}
		}`
	id := int32(12)
	idStr := stringutils.FromInt32(id)
	variables := map[string]interface{}{
		"id": idStr,
	}

	inputs := pipeline.JSONSerializable{}
	err := inputs.UnmarshalJSON([]byte(`{"foo": "bar"}`))
	require.NoError(t, err)

	outputs := pipeline.JSONSerializable{}
	err = outputs.UnmarshalJSON([]byte(`[{"baz": "bar"}]`))
	require.NoError(t, err)

	gError := errors.New("error")
	_, idErr := stringutils.ToInt32("some random ID with some specific length that should not work")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "runJob"),
		{
			name:          "success without body",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("RunJobV2", mock.Anything, id, (map[string]interface{})(nil)).Return(int64(25), nil)
				f.Mocks.pipelineORM.On("FindRun", int64(25)).Return(pipeline.Run{
					ID:             2,
					PipelineSpecID: 5,
					CreatedAt:      f.Timestamp(),
					FinishedAt:     null.TimeFrom(f.Timestamp()),
					AllErrors:      pipeline.RunErrors{null.StringFrom("fatal error"), null.String{}},
					FatalErrors:    pipeline.RunErrors{null.StringFrom("fatal error"), null.String{}},
					Inputs:         inputs,
					Outputs:        outputs,
					State:          pipeline.RunStatusErrored,
				}, nil)
				f.App.On("PipelineORM").Return(f.Mocks.pipelineORM)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"runJob": {
						"jobRun": {
							"id": "2",
							"allErrors": ["fatal error"],
							"createdAt": "2021-01-01T00:00:00Z",
							"fatalErrors": ["fatal error"],
							"finishedAt": "2021-01-01T00:00:00Z",
							"inputs": "{\"foo\":\"bar\"}",
							"outputs": ["{\"baz\":\"bar\"}"],
							"status": "ERRORED"
						}
					}
				}`,
		},
		{
			name:          "invalid ID error",
			authenticated: true,
			query:         mutation,
			variables: map[string]interface{}{
				"id": "some random ID with some specific length that should not work",
			},
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: idErr,
					Path:          []interface{}{"runJob"},
					Message:       idErr.Error(),
				},
			},
		},
		{
			name:          "not found job error",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("RunJobV2", mock.Anything, id, (map[string]interface{})(nil)).Return(int64(25), webhook.ErrJobNotExists)
			},
			query: mutation,
			variables: map[string]interface{}{
				"id": idStr,
			},
			result: `
				{
					"runJob": {
						"code": "NOT_FOUND",
						"message": "job does not exist"
					}
				}`,
		},
		{
			name:          "generic error on RunJobV2",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("RunJobV2", mock.Anything, id, (map[string]interface{})(nil)).Return(int64(25), gError)
			},
			query: mutation,
			variables: map[string]interface{}{
				"id": idStr,
			},
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"runJob"},
					Message:       gError.Error(),
				},
			},
		},
		{
			name:          "generic error on FindRun",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("RunJobV2", mock.Anything, id, (map[string]interface{})(nil)).Return(int64(25), nil)
				f.Mocks.pipelineORM.On("FindRun", int64(25)).Return(pipeline.Run{}, gError)
				f.App.On("PipelineORM").Return(f.Mocks.pipelineORM)
			},
			query: mutation,
			variables: map[string]interface{}{
				"id": idStr,
			},
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"runJob"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
