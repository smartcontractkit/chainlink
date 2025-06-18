package resolver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-evm/pkg/chains"
	clnull "github.com/smartcontractkit/chainlink/v2/core/null"
	"github.com/smartcontractkit/chainlink/v2/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
)

// This tests the main fields on the job results. Embedded spec testing is done
// in the `spec_test` file
func TestResolver_Jobs(t *testing.T) {
	var (
		externalJobID = uuid.MustParse(("00000000-0000-0000-0000-000000000001"))

		query = `
			query GetJobs {
				jobs {
					results {
						id
						createdAt
						externalJobID
						gasLimit
						forwardingAllowed
						maxTaskDuration
						name
						schemaVersion
						spec {
							__typename
						}
						runs {
							__typename
							results {
								id
							}
							metadata {
								total
							}
						}
						observationSource
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
			before: func(ctx context.Context, f *gqlTestFramework) {
				plnSpecID := int32(12)

				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobs", mock.Anything, 0, 50).Return([]job.Job{
					{
						ID:              1,
						Name:            null.StringFrom("job1"),
						SchemaVersion:   1,
						MaxTaskDuration: models.Interval(1 * time.Second),
						ExternalJobID:   externalJobID,
						CreatedAt:       f.Timestamp(),
						Type:            job.OffchainReporting,
						PipelineSpecID:  plnSpecID,
						OCROracleSpec:   &job.OCROracleSpec{},
						PipelineSpec: &pipeline.Spec{
							DotDagSource: "ds1 [type=bridge name=voter_turnout];",
						},
					},
				}, 1, nil)
				f.Mocks.jobORM.
					On("FindPipelineRunIDsByJobID", mock.Anything, int32(1), 0, 50).
					Return([]int64{200}, nil)
				f.Mocks.jobORM.
					On("FindPipelineRunsByIDs", mock.Anything, []int64{200}).
					Return([]pipeline.Run{{ID: 200}}, nil)
				f.Mocks.jobORM.
					On("CountPipelineRunsByJobID", mock.Anything, int32(1)).
					Return(int32(1), nil)
			},
			query: query,
			result: `
				{
					"jobs": {
						"results": [{
							"id": "1",
							"createdAt": "2021-01-01T00:00:00Z",
							"externalJobID": "00000000-0000-0000-0000-000000000001",
							"gasLimit": null,
							"forwardingAllowed": false,
							"maxTaskDuration": "1s",
							"name": "job1",
							"schemaVersion": 1,
							"spec": {
								"__typename": "OCRSpec"
							},
							"runs": {
								"__typename": "JobRunsPayload",
								"results": [{
									"id": "200"
								}],
								"metadata": {
									"total": 1
								}
							},
							"observationSource": "ds1 [type=bridge name=voter_turnout];"
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
		externalJobID = uuid.MustParse(("00000000-0000-0000-0000-000000000001"))

		query = `
			query GetJob {
				job(id: "1") {
					... on Job {
						id
						createdAt
						externalJobID
						gasLimit
						maxTaskDuration
						name
						schemaVersion
						spec {
							__typename
						}
						runs {
							__typename
							results {
								id
							}
							metadata {
								total
							}
						}
						observationSource
					}
					... on NotFoundError {
						code
						message
					}
				}
			}
		`
		exampleJobResult = `
				{
					"job": {
						"id": "1",
						"createdAt": "2021-01-01T00:00:00Z",
						"externalJobID": "00000000-0000-0000-0000-000000000001",
						"gasLimit": 123,
						"maxTaskDuration": "1s",
						"name": "job1",
						"schemaVersion": 1,
						"spec": {
							"__typename": "OCRSpec"
						},
						"runs": {
							"__typename": "JobRunsPayload",
							"results": [{
								"id": "200"
							}],
							"metadata": {
								"total": 1
							}
						},
						"observationSource": "ds1 [type=bridge name=voter_turnout];"
					}
				}
			`
	)

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "job"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					ID:              1,
					Name:            null.StringFrom("job1"),
					SchemaVersion:   1,
					GasLimit:        clnull.Uint32From(123),
					MaxTaskDuration: models.Interval(1 * time.Second),
					ExternalJobID:   externalJobID,
					CreatedAt:       f.Timestamp(),
					Type:            job.OffchainReporting,
					OCROracleSpec:   &job.OCROracleSpec{},
					PipelineSpec: &pipeline.Spec{
						DotDagSource: "ds1 [type=bridge name=voter_turnout];",
					},
				}, nil)
				f.Mocks.jobORM.
					On("FindPipelineRunIDsByJobID", mock.Anything, int32(1), 0, 50).
					Return([]int64{200}, nil)
				f.Mocks.jobORM.
					On("FindPipelineRunsByIDs", mock.Anything, []int64{200}).
					Return([]pipeline.Run{{ID: 200}}, nil)
				f.Mocks.jobORM.
					On("CountPipelineRunsByJobID", mock.Anything, int32(1)).
					Return(int32(1), nil)
			},
			query:  query,
			result: exampleJobResult,
		},
		{
			name:          "not found",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{}, sql.ErrNoRows)
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
		{
			name:          "show job when chainID is disabled",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					ID:              1,
					Name:            null.StringFrom("job1"),
					SchemaVersion:   1,
					GasLimit:        clnull.Uint32From(123),
					MaxTaskDuration: models.Interval(1 * time.Second),
					ExternalJobID:   externalJobID,
					CreatedAt:       f.Timestamp(),
					Type:            job.OffchainReporting,
					OCROracleSpec:   &job.OCROracleSpec{},
					PipelineSpec: &pipeline.Spec{
						DotDagSource: "ds1 [type=bridge name=voter_turnout];",
					},
				}, chains.ErrNoSuchChainID)
				f.Mocks.jobORM.
					On("FindPipelineRunIDsByJobID", mock.Anything, int32(1), 0, 50).
					Return([]int64{200}, nil)
				f.Mocks.jobORM.
					On("FindPipelineRunsByIDs", mock.Anything, []int64{200}).
					Return([]pipeline.Run{{ID: 200}}, nil)
				f.Mocks.jobORM.
					On("CountPipelineRunsByJobID", mock.Anything, int32(1)).
					Return(int32(1), nil)
			},
			query:  query,
			result: exampleJobResult,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_CreateJob(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation CreateJob($input: CreateJobInput!) {
			createJob(input: $input) {
				... on CreateJobSuccess {
					job {
						id
						createdAt
						externalJobID
						maxTaskDuration
						name
						schemaVersion
					}
				}
				... on InputErrors {
					errors {
						path
						message
						code
					}
				}
			}
		}`
	uuid := uuid.New()
	spec := fmt.Sprintf(testspecs.DirectRequestSpecTemplate, uuid, uuid)
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"TOML": spec,
		},
	}
	invalid := map[string]interface{}{
		"input": map[string]interface{}{
			"TOML": "some wrong value",
		},
	}
	jb, err := directrequest.ValidatedDirectRequestSpec(spec)
	assert.NoError(t, err)

	d, err := json.Marshal(map[string]interface{}{
		"createJob": map[string]interface{}{
			"job": map[string]interface{}{
				"id":              "0",
				"maxTaskDuration": "0s",
				"name":            jb.Name,
				"schemaVersion":   1,
				"createdAt":       "0001-01-01T00:00:00Z",
				"externalJobID":   jb.ExternalJobID.String(),
			},
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "createJob"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetConfig").Return(f.Mocks.cfg)
				f.App.On("AddJobV2", mock.Anything, &jb).Return(nil)
			},
			query:     mutation,
			variables: variables,
			result:    expected,
		},
		{
			name:          "invalid TOML error",
			authenticated: true,
			query:         mutation,
			variables:     invalid,
			result: `
				{
					"createJob": {
						"errors": [{
							"code": "INVALID_INPUT",
							"message": "failed to parse TOML: (1, 6): was expecting token =, but got \"wrong\" instead",
							"path": "TOML spec"
						}]
					}
				}`,
		},
		{
			name:          "generic error when adding the job",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetConfig").Return(f.Mocks.cfg)
				f.App.On("AddJobV2", mock.Anything, &jb).Return(gError)
			},
			query:     mutation,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"createJob"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_DeleteJob(t *testing.T) {
	t.Parallel()

	id := int32(123)
	extJID := uuid.New()
	mutation := `
		mutation DeleteJob($id: ID!) {
			deleteJob(id: $id) {
				... on DeleteJobSuccess {
					job {
						id
						createdAt
						externalJobID
						maxTaskDuration
						name
						schemaVersion
					}
				}
				... on NotFoundError {
						code
						message
					}
				}
		}`
	variables := map[string]interface{}{
		"id": "123",
	}
	invalidVariables := map[string]interface{}{
		"id": "asdadada",
	}
	d, err := json.Marshal(map[string]interface{}{
		"deleteJob": map[string]interface{}{
			"job": map[string]interface{}{
				"id":              "123",
				"maxTaskDuration": "2s",
				"name":            "test-job",
				"schemaVersion":   0,
				"createdAt":       "2021-01-01T00:00:00Z",
				"externalJobID":   extJID.String(),
			},
		},
	})
	assert.NoError(t, err)
	expected := string(d)

	gError := errors.New("error")
	_, idError := stringutils.ToInt64("asdadada")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "deleteJob"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{
					ID:              id,
					Name:            null.StringFrom("test-job"),
					ExternalJobID:   extJID,
					MaxTaskDuration: models.Interval(2 * time.Second),
					CreatedAt:       f.Timestamp(),
				}, nil)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.App.On("DeleteJob", mock.Anything, id).Return(nil)
			},
			query:     mutation,
			variables: variables,
			result:    expected,
		},
		{
			name:          "not found on FindJob()",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{}, sql.ErrNoRows)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"deleteJob": {
						"code": "NOT_FOUND",
						"message": "job not found"
					}
				}
			`,
		},
		{
			name:          "not found on DeleteJob()",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{}, nil)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.App.On("DeleteJob", mock.Anything, id).Return(sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"deleteJob": {
						"code": "NOT_FOUND",
						"message": "job not found"
					}
				}
			`,
		},
		{
			name:          "generic error on FindJob()",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{}, gError)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
			},
			query:     mutation,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"deleteJob"},
					Message:       gError.Error(),
				},
			},
		},
		{
			name:          "generic error on DeleteJob()",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.jobORM.On("FindJobWithoutSpecErrors", mock.Anything, id).Return(job.Job{}, nil)
				f.App.On("JobORM").Return(f.Mocks.jobORM)
				f.App.On("DeleteJob", mock.Anything, id).Return(gError)
			},
			query:     mutation,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"deleteJob"},
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
					Path:          []interface{}{"deleteJob"},
					Message:       idError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
