package resolver

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
)

func TestResolver_ApproveJobProposalSpec(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation ApproveJobProposalSpec($id: ID!) {
			approveJobProposalSpec(id: $id) {
				... on ApproveJobProposalSpecSuccess {
					spec {
						id
					}
				}
				... on NotFoundError {
					message
					code
				}
				... on JobAlreadyExistsError {
					message
					code
				}
			}
		}`

	specID := int64(1)
	result := `
		{
			"approveJobProposalSpec": {
				"spec": {
					"id": "1"
				}
			}
		}`
	variables := map[string]interface{}{
		"id": "1",
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "approveJobProposalSpec"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("ApproveSpec", mock.Anything, specID, false).Return(nil)
				f.Mocks.feedsSvc.On("GetSpec", mock.Anything, specID).Return(&feeds.JobProposalSpec{
					ID: specID,
				}, nil)
			},
			query:     mutation,
			variables: variables,
			result:    result,
		},
		{
			name:          "not found error on approval",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("ApproveSpec", mock.Anything, specID, false).Return(sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"approveJobProposalSpec": {
					"message": "spec not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
		{
			name:          "not found error on fetch",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("ApproveSpec", mock.Anything, specID, false).Return(nil)
				f.Mocks.feedsSvc.On("GetSpec", mock.Anything, specID).Return(nil, sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"approveJobProposalSpec": {
					"message": "spec not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
		{
			name:          "unprocessable error on approval if job already exists",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("ApproveSpec", mock.Anything, specID, false).Return(feeds.ErrJobAlreadyExists)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"approveJobProposalSpec": {
					"message": "a job for this contract address already exists - please use the 'force' option to replace it",
					"code": "UNPROCESSABLE"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_CancelJobProposalSpec(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation CancelJobProposalSpec($id: ID!) {
			cancelJobProposalSpec(id: $id) {
				... on CancelJobProposalSpecSuccess {
					spec {
						id
					}
				}
				... on NotFoundError {
					message
					code
				}
			}
		}`

	specID := int64(1)
	result := `
		{
			"cancelJobProposalSpec": {
				"spec": {
					"id": "1"
				}
			}
		}`
	variables := map[string]interface{}{
		"id": "1",
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "cancelJobProposalSpec"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("CancelSpec", mock.Anything, specID).Return(nil)
				f.Mocks.feedsSvc.On("GetSpec", mock.Anything, specID).Return(&feeds.JobProposalSpec{
					ID: specID,
				}, nil)
			},
			query:     mutation,
			variables: variables,
			result:    result,
		},
		{
			name:          "not found error on cancel",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("CancelSpec", mock.Anything, specID).Return(sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"cancelJobProposalSpec": {
					"message": "spec not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
		{
			name:          "not found error on fetch",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("CancelSpec", mock.Anything, specID).Return(nil)
				f.Mocks.feedsSvc.On("GetSpec", mock.Anything, specID).Return(nil, sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"cancelJobProposalSpec": {
					"message": "spec not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_RejectJobProposalSpec(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation RejectJobProposalSpec($id: ID!) {
			rejectJobProposalSpec(id: $id) {
				... on RejectJobProposalSpecSuccess {
					spec {
						id
					}
				}
				... on NotFoundError {
					message
					code
				}
			}
		}`

	specID := int64(1)
	result := `
		{
			"rejectJobProposalSpec": {
				"spec": {
					"id": "1"
				}
			}
		}`
	variables := map[string]interface{}{
		"id": "1",
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "rejectJobProposalSpec"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("RejectSpec", mock.Anything, specID).Return(nil)
				f.Mocks.feedsSvc.On("GetSpec", mock.Anything, specID).Return(&feeds.JobProposalSpec{
					ID: specID,
				}, nil)
			},
			query:     mutation,
			variables: variables,
			result:    result,
		},
		{
			name:          "not found error on reject",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("RejectSpec", mock.Anything, specID).Return(sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"rejectJobProposalSpec": {
					"message": "spec not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
		{
			name:          "not found error on fetch",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("RejectSpec", mock.Anything, specID).Return(nil)
				f.Mocks.feedsSvc.On("GetSpec", mock.Anything, specID).Return(nil, sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"rejectJobProposalSpec": {
					"message": "spec not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_UpdateJobProposalSpecDefinition(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation UpdateJobProposalSpecDefinition($id: ID!, $input: UpdateJobProposalSpecDefinitionInput!) {
			updateJobProposalSpecDefinition(id: $id, input: $input) {
				... on UpdateJobProposalSpecDefinitionSuccess {
					spec {
						id
					}
				}
				... on NotFoundError {
					message
					code
				}
			}
		}`

	specID := int64(1)
	result := `
		{
			"updateJobProposalSpecDefinition": {
				"spec": {
					"id": "1"
				}
			}
		}`
	variables := map[string]interface{}{
		"id": "1",
		"input": map[string]interface{}{
			"definition": "",
		},
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "updateJobProposalSpecDefinition"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("UpdateSpecDefinition", mock.Anything, specID, "").Return(nil)
				f.Mocks.feedsSvc.On("GetSpec", mock.Anything, specID).Return(&feeds.JobProposalSpec{
					ID: specID,
				}, nil)
			},
			query:     mutation,
			variables: variables,
			result:    result,
		},
		{
			name:          "not found error on update",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("UpdateSpecDefinition", mock.Anything, specID, "").Return(sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"updateJobProposalSpecDefinition": {
					"message": "spec not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
		{
			name:          "not found error on fetch",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
				f.Mocks.feedsSvc.On("UpdateSpecDefinition", mock.Anything, specID, "").Return(nil)
				f.Mocks.feedsSvc.On("GetSpec", mock.Anything, specID).Return(nil, sql.ErrNoRows)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"updateJobProposalSpecDefinition": {
					"message": "spec not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}

// Tests resolving a job proposal spec. Since there is not GetJobProposalSpec
// query, we rely on the GetJobProposal query to fetch the nested specs
func TestResolver_GetJobProposal_Spec(t *testing.T) {
	t.Parallel()

	timestamp := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	query := `
		query GetJobProposal {
			jobProposal(id: "1") {
				... on JobProposal {
					id
					specs {
						id
						definition
						status
						version
						statusUpdatedAt
						createdAt
						updatedAt
					}
				}
				... on NotFoundError {
					message
					code
				}
			}
		}`

	jpID := int64(1)
	spec := feeds.JobProposalSpec{
		ID:              100,
		Definition:      "name='spec'",
		Status:          feeds.SpecStatusPending,
		JobProposalID:   jpID,
		Version:         1,
		StatusUpdatedAt: timestamp,
		CreatedAt:       timestamp,
		UpdatedAt:       timestamp,
	}
	specs := []feeds.JobProposalSpec{spec}
	result := `
		{
			"jobProposal": {
				"id": "1",
				"specs": [{
					"id": "100",
					"definition": "name='spec'",
					"status": "PENDING",
					"version": 1,
					"statusUpdatedAt": "2021-01-01T00:00:00Z",
					"createdAt": "2021-01-01T00:00:00Z",
					"updatedAt": "2021-01-01T00:00:00Z"
				}]
			}
		}`

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "jobProposal"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("GetJobProposal", mock.Anything, jpID).Return(&feeds.JobProposal{
					ID:             jpID,
					Status:         feeds.JobProposalStatusApproved,
					FeedsManagerID: 1,
					Multiaddrs:     []string{"1", "2"},
					PendingUpdate:  false,
				}, nil)
				f.Mocks.feedsSvc.
					On("ListSpecsByJobProposalIDs", mock.Anything, []int64{jpID}).
					Return(specs, nil)
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
			},
			query:  query,
			result: result,
		},
	}

	RunGQLTests(t, testCases)
}
