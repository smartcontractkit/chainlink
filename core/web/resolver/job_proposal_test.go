package resolver

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"

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
	ejID := uuid.NullUUID{UUID: uuid.NewV4(), Valid: true}
	result := `
		{
			"jobProposal": {
				"id": "1",
				"spec": "some-spec",
				"status": "APPROVED",
				"externalJobID": "%s",
				"multiAddrs": ["1", "2"],
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
				f.Mocks.feedsSvc.On("GetManagers", []int64{1}).Return([]feeds.FeedsManager{
					{
						ID:   1,
						Name: "manager",
					},
				}, nil)
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

func TestResolver_ApproveJobProposal(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation ApproveJobProposal($id: ID!) {
			approveJobProposal(id: $id) {
				... on ApproveJobProposalSuccess {
					jobProposal {
						id
						spec
						status
						externalJobID
						multiAddrs
						feedsManager {
							id
							name
						}
					}
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
			"approveJobProposal": {
				"jobProposal": {
					"id": "1",
					"spec": "some-spec",
					"status": "APPROVED",
					"externalJobID": "%s",
					"multiAddrs": ["1", "2"],
					"feedsManager": {
						"id": "1",
						"name": "manager"
					}
				}
			}
		}`
	variables := map[string]interface{}{
		"id": "1",
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "approveJobProposal"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("GetManagers", []int64{1}).Return([]feeds.FeedsManager{
					{
						ID:   1,
						Name: "manager",
					},
				}, nil)
				f.Mocks.feedsSvc.On("ApproveJobProposal", mock.Anything, jpID).Return(nil)
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
			query:     mutation,
			variables: variables,
			result:    fmt.Sprintf(result, ejID.UUID.String()),
		},
		{
			name:          "not found error on approval",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("ApproveJobProposal", mock.Anything, jpID).Return(sql.ErrNoRows)
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"approveJobProposal": {
					"message": "job proposal not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
		{
			name:          "not found error on fetch",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("ApproveJobProposal", mock.Anything, jpID).Return(nil)
				f.Mocks.feedsSvc.On("GetJobProposal", jpID).Return(nil, sql.ErrNoRows)
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"approveJobProposal": {
					"message": "job proposal not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_CancelJobProposal(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation CancelJobProposal($id: ID!) {
			cancelJobProposal(id: $id) {
				... on CancelJobProposalSuccess {
					jobProposal {
						id
						spec
						status
						externalJobID
						multiAddrs
						feedsManager {
							id
							name
						}
					}
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
			"cancelJobProposal": {
				"jobProposal": {
					"id": "1",
					"spec": "some-spec",
					"status": "CANCELLED",
					"externalJobID": "%s",
					"multiAddrs": ["1", "2"],
					"feedsManager": {
						"id": "1",
						"name": "manager"
					}
				}
			}
		}`
	variables := map[string]interface{}{
		"id": "1",
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "cancelJobProposal"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("GetManagers", []int64{1}).Return([]feeds.FeedsManager{
					{
						ID:   1,
						Name: "manager",
					},
				}, nil)
				f.Mocks.feedsSvc.On("CancelJobProposal", mock.Anything, jpID).Return(nil)
				f.Mocks.feedsSvc.On("GetJobProposal", jpID).Return(&feeds.JobProposal{
					ID:             jpID,
					Spec:           "some-spec",
					Status:         feeds.JobProposalStatusCancelled,
					ExternalJobID:  ejID,
					FeedsManagerID: 1,
					Multiaddrs:     []string{"1", "2"},
					ProposedAt:     time.Time{},
				}, nil)
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
			},
			query:     mutation,
			variables: variables,
			result:    fmt.Sprintf(result, ejID.UUID.String()),
		},
		{
			name:          "not found error on cancel",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("CancelJobProposal", mock.Anything, jpID).Return(sql.ErrNoRows)
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"cancelJobProposal": {
					"message": "job proposal not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
		{
			name:          "not found error on fetch",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("CancelJobProposal", mock.Anything, jpID).Return(nil)
				f.Mocks.feedsSvc.On("GetJobProposal", jpID).Return(nil, sql.ErrNoRows)
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"cancelJobProposal": {
					"message": "job proposal not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_RejectJobProposal(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation RejectJobProposal($id: ID!) {
			rejectJobProposal(id: $id) {
				... on RejectJobProposalSuccess {
					jobProposal {
						id
						spec
						status
						externalJobID
						multiAddrs
						feedsManager {
							id
							name
						}
					}
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
			"rejectJobProposal": {
				"jobProposal": {
					"id": "1",
					"spec": "some-spec",
					"status": "REJECTED",
					"externalJobID": "%s",
					"multiAddrs": ["1", "2"],
					"feedsManager": {
						"id": "1",
						"name": "manager"
					}
				}
			}
		}`
	variables := map[string]interface{}{
		"id": "1",
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "rejectJobProposal"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("GetManagers", []int64{1}).Return([]feeds.FeedsManager{
					{
						ID:   1,
						Name: "manager",
					},
				}, nil)
				f.Mocks.feedsSvc.On("RejectJobProposal", mock.Anything, jpID).Return(nil)
				f.Mocks.feedsSvc.On("GetJobProposal", jpID).Return(&feeds.JobProposal{
					ID:             jpID,
					Spec:           "some-spec",
					Status:         feeds.JobProposalStatusRejected,
					ExternalJobID:  ejID,
					FeedsManagerID: 1,
					Multiaddrs:     []string{"1", "2"},
					ProposedAt:     time.Time{},
				}, nil)
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
			},
			query:     mutation,
			variables: variables,
			result:    fmt.Sprintf(result, ejID.UUID.String()),
		},
		{
			name:          "not found error on reject",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("RejectJobProposal", mock.Anything, jpID).Return(sql.ErrNoRows)
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"rejectJobProposal": {
					"message": "job proposal not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
		{
			name:          "not found error on fetch",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("RejectJobProposal", mock.Anything, jpID).Return(nil)
				f.Mocks.feedsSvc.On("GetJobProposal", jpID).Return(nil, sql.ErrNoRows)
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"rejectJobProposal": {
					"message": "job proposal not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}

func TestResolver_UpdateJobSpecProposal(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation UpdateJobProposalSpec($id: ID!, $input: UpdateJobProposalSpecInput!) {
			updateJobProposalSpec(id: $id, input: $input) {
				... on UpdateJobProposalSpecSuccess {
					jobProposal {
						id
						spec
						status
						externalJobID
						multiAddrs
						feedsManager {
							id
							name
						}
					}
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
			"updateJobProposalSpec": {
				"jobProposal": {
					"id": "1",
					"spec": "some-spec",
					"status": "APPROVED",
					"externalJobID": "%s",
					"multiAddrs": ["1", "2"],
					"feedsManager": {
						"id": "1",
						"name": "manager"
					}
				}
			}
		}`
	variables := map[string]interface{}{
		"id": "1",
		"input": map[string]interface{}{
			"spec": "",
		},
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "updateJobProposalSpec"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("GetManagers", []int64{1}).Return([]feeds.FeedsManager{
					{
						ID:   1,
						Name: "manager",
					},
				}, nil)
				f.Mocks.feedsSvc.On("UpdateJobProposalSpec", mock.Anything, jpID, "").Return(nil)
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
			query:     mutation,
			variables: variables,
			result:    fmt.Sprintf(result, ejID.UUID.String()),
		},
		{
			name:          "not found error on update",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("UpdateJobProposalSpec", mock.Anything, jpID, "").Return(sql.ErrNoRows)
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"updateJobProposalSpec": {
					"message": "job proposal not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
		{
			name:          "not found error on fetch",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.Mocks.feedsSvc.On("UpdateJobProposalSpec", mock.Anything, jpID, "").Return(nil)
				f.Mocks.feedsSvc.On("GetJobProposal", jpID).Return(nil, sql.ErrNoRows)
				f.App.On("GetFeedsService").Return(f.Mocks.feedsSvc)
			},
			query:     mutation,
			variables: variables,
			result: `
			{
				"updateJobProposalSpec": {
					"message": "job proposal not found",
					"code": "NOT_FOUND"
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}
