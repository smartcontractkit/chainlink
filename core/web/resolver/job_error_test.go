package resolver

import (
	"database/sql"
	"encoding/json"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
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
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", id).Return(job.Job{
					ID: int32(1),
				}, nil)
				f.Mocks.jobORM.On("FindSpecErrorsByJobIDs", []int32{1}, mock.Anything).Return([]job.SpecError{
					{
						ID:          errorID,
						Description: "no contract code at given address",
						Occurrences: 1,
						CreatedAt:   f.Timestamp(),
						JobID:       int32(1),
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

func TestResolver_DismissJobError(t *testing.T) {
	t.Parallel()

	id := int64(1)
	mutation := `
		mutation DismissJobError($id: ID!) {
			dismissJobError(id: $id) {
				... on DismissJobErrorSuccess {
					jobError {
						id
						description
						occurrences
						createdAt
					}
				}
				... on NotFoundError {
						code
						message
					}
				}
		}`
	variables := map[string]interface{}{
		"id": "1",
	}
	invalidVariables := map[string]interface{}{
		"id": "asdadada",
	}
	d, err := json.Marshal(map[string]interface{}{
		"dismissJobError": map[string]interface{}{
			"jobError": map[string]interface{}{
				"id":          "1",
				"occurrences": 5,
				"description": "test-description",
				"createdAt":   "2021-01-01T00:00:00Z",
			},
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	gError := errors.New("error")

	_, idError := stringutils.ToInt64("asdadada")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "dismissJobError"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.jobORM.On("FindSpecError", id).Return(job.SpecError{
					ID:          id,
					Occurrences: 5,
					Description: "test-description",
					CreatedAt:   f.Timestamp(),
				}, nil)
				f.Mocks.jobORM.On("DismissError", mock.Anything, id).Return(nil)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
			},
			query:     mutation,
			variables: variables,
			result:    expected,
		},
		{
			name:          "not found on FindSpecError()",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.jobORM.On("FindSpecError", id).Return(job.SpecError{}, sql.ErrNoRows)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"dismissJobError": {
						"code": "NOT_FOUND",
						"message": "JobSpecError not found"
					}
				}
			`,
		},
		{
			name:          "not found on DismissError()",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.jobORM.On("FindSpecError", id).Return(job.SpecError{}, nil)
				f.Mocks.jobORM.On("DismissError", mock.Anything, id).Return(sql.ErrNoRows)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"dismissJobError": {
						"code": "NOT_FOUND",
						"message": "JobSpecError not found"
					}
				}
			`,
		},
		{
			name:          "generic error on FindSpecError()",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.jobORM.On("FindSpecError", id).Return(job.SpecError{}, gError)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
			},
			query:     mutation,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"dismissJobError"},
					Message:       gError.Error(),
				},
			},
		},
		{
			name:          "generic error on DismissError()",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.jobORM.On("FindSpecError", id).Return(job.SpecError{}, nil)
				f.Mocks.jobORM.On("DismissError", mock.Anything, id).Return(gError)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
			},
			query:     mutation,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"dismissJobError"},
					Message:       gError.Error(),
				},
			},
		},
		{
			name:          "error on ID parsing",
			authenticated: true,
			query:         mutation,
			variables:     invalidVariables,
			result:        `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: idError,
					Path:          []interface{}{"dismissJobError"},
					Message:       idError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
