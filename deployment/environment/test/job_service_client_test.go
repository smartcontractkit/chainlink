package test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink/deployment/helpers/pointer"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
)

func TestNewJobServiceClient(t *testing.T) {
	t.Parallel()

	// Create a mock JobApprover getter
	mockGetter := &mockJobApproverGetter{
		jobApprovers: make(map[string]*mockJobApprover),
	}
	mockGetter.jobApprovers["node-1"] = &mockJobApprover{}

	// Create a new JobServiceClient
	client := NewJobServiceClient(mockGetter)

	// Assert that the client was initialized correctly
	require.NotNil(t, client)
	require.NotNil(t, client.jobStore)
	require.NotNil(t, client.proposalStore)
	require.NotNil(t, client.jobApproverStore)

	// Assert that it's the same getter we passed in
	require.Equal(t, mockGetter, client.jobApproverStore)
}

func TestBatchProposeJob(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	// Create mock job approvers
	mockGetter := &mockJobApproverGetter{
		jobApprovers: make(map[string]*mockJobApprover),
	}
	mockGetter.jobApprovers["node-1"] = &mockJobApprover{}
	mockGetter.jobApprovers["node-2"] = &mockJobApprover{}

	client := NewJobServiceClient(mockGetter)

	// Valid job spec for testing
	externalJobID := uuid.NewString()
	jobSpec := createValidJobSpec(externalJobID)

	t.Run("successful batch propose to multiple nodes", func(t *testing.T) {
		req := &jobv1.BatchProposeJobRequest{
			NodeIds: []string{"node-1", "node-2"},
			Spec:    jobSpec,
		}

		resp, err := client.BatchProposeJob(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.SuccessResponses, 2)
		require.Empty(t, resp.FailedResponses)

		// Verify that both nodes have the job
		require.Contains(t, resp.SuccessResponses, "node-1")
		require.Contains(t, resp.SuccessResponses, "node-2")
	})

	t.Run("one node fails", func(t *testing.T) {
		mockGetter.jobApprovers["node-2"].shouldFail = true

		req := &jobv1.BatchProposeJobRequest{
			NodeIds: []string{"node-1", "node-2"},
			Spec:    jobSpec,
		}

		resp, err := client.BatchProposeJob(ctx, req)
		require.Error(t, err) // Total error should be non-nil
		require.NotNil(t, resp)
		require.Len(t, resp.SuccessResponses, 1)
		require.Len(t, resp.FailedResponses, 1)

		require.Contains(t, resp.SuccessResponses, "node-1")
		require.Contains(t, resp.FailedResponses, "node-2")

		// Reset for next tests
		mockGetter.jobApprovers["node-2"].shouldFail = false
	})

	t.Run("node not found", func(t *testing.T) {
		req := &jobv1.BatchProposeJobRequest{
			NodeIds: []string{"node-1", "non-existent-node"},
			Spec:    jobSpec,
		}

		resp, err := client.BatchProposeJob(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "node not found")
	})

	t.Run("no nodes provided", func(t *testing.T) {
		req := &jobv1.BatchProposeJobRequest{
			NodeIds: []string{},
			Spec:    jobSpec,
		}

		resp, err := client.BatchProposeJob(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "no nodes found")
	})
}

func TestProposeJob(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Create mock job approvers
	mockGetter := &mockJobApproverGetter{
		jobApprovers: make(map[string]*mockJobApprover),
	}
	mockGetter.jobApprovers["node-1"] = &mockJobApprover{}

	client := NewJobServiceClient(mockGetter)

	t.Run("successful job proposal", func(t *testing.T) {
		externalJobID := uuid.NewString()
		jobSpec := createValidJobSpec(externalJobID)

		req := &jobv1.ProposeJobRequest{
			NodeId: "node-1",
			Spec:   jobSpec,
		}

		resp, err := client.ProposeJob(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Proposal)
		require.Equal(t, externalJobID, resp.Proposal.JobId)
		require.Equal(t, jobSpec, resp.Proposal.Spec)
		require.Equal(t, int64(1), resp.Proposal.Revision)
		require.Equal(t, jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED, resp.Proposal.Status)
	})

	t.Run("successful job proposal without approval", func(t *testing.T) {
		externalJobID := uuid.NewString()
		jobSpec := createValidJobSpec(externalJobID)

		req := &jobv1.ProposeJobRequest{
			NodeId: "node-1",
			Spec:   jobSpec,
			Labels: []*ptypes.Label{
				{
					Key: LabelDoNotAutoApprove,
				},
			},
		}

		resp, err := client.ProposeJob(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Proposal)
		require.Equal(t, externalJobID, resp.Proposal.JobId)
		require.Equal(t, jobSpec, resp.Proposal.Spec)
		require.Equal(t, int64(1), resp.Proposal.Revision)
		require.Equal(t, jobv1.ProposalStatus_PROPOSAL_STATUS_PENDING, resp.Proposal.Status)
	})

	t.Run("node not found", func(t *testing.T) {
		externalJobID := uuid.NewString()
		jobSpec := createValidJobSpec(externalJobID)

		req := &jobv1.ProposeJobRequest{
			NodeId: "non-existent-node",
			Spec:   jobSpec,
		}

		resp, err := client.ProposeJob(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "node not found")
	})

	t.Run("invalid job spec", func(t *testing.T) {
		req := &jobv1.ProposeJobRequest{
			NodeId: "node-1",
			Spec:   "invalid job spec",
		}

		resp, err := client.ProposeJob(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("invalid job spec", func(t *testing.T) {
		jobSpec := `
type = "test"
name = "Test Job"
`

		req := &jobv1.ProposeJobRequest{
			NodeId: "node-1",
			Spec:   jobSpec,
		}

		resp, err := client.ProposeJob(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("approval failure", func(t *testing.T) {
		mockGetter.jobApprovers["node-1"].shouldFail = true

		externalJobID := uuid.NewString()
		jobSpec := createValidJobSpec(externalJobID)

		req := &jobv1.ProposeJobRequest{
			NodeId: "node-1",
			Spec:   jobSpec,
		}

		resp, err := client.ProposeJob(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "failed to auto approve job")

		// Reset for next tests
		mockGetter.jobApprovers["node-1"].shouldFail = false
	})
}

func TestRevokeJob(t *testing.T) {
	t.Parallel()

	mockGetter := &mockJobApproverGetter{
		jobApprovers: make(map[string]*mockJobApprover),
	}
	mockGetter.jobApprovers["node-1"] = &mockJobApprover{}
	client := NewJobServiceClient(mockGetter)

	proposeJob := func(autoApprove bool) *jobv1.Proposal {
		externalJobID := uuid.NewString()
		jobSpec := createValidJobSpec(externalJobID)

		labels := []*ptypes.Label{}
		if !autoApprove {
			labels = append(labels, &ptypes.Label{
				Key: LabelDoNotAutoApprove,
			})
		}
		req := &jobv1.ProposeJobRequest{
			NodeId: "node-1",
			Spec:   jobSpec,
			Labels: labels,
		}

		resp, err := client.ProposeJob(t.Context(), req)
		require.NoError(t, err)

		return resp.GetProposal()
	}

	tests := []struct {
		name          string
		autoApprove   bool
		overrideJobID string
		expectedError string
	}{
		{
			name:        "successful revoke job",
			autoApprove: false,
		},
		{
			name:          "cannot revoke approved job",
			autoApprove:   true,
			expectedError: "is not revokable: status PROPOSAL_STATUS_APPROVED",
		},
		{
			name:          "job not found",
			overrideJobID: "non-existent-job",
			expectedError: "no proposals found for job",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			job := proposeJob(tc.autoApprove)

			if tc.overrideJobID != "" {
				job.JobId = tc.overrideJobID
			}

			resp, err := client.RevokeJob(t.Context(), &jobv1.RevokeJobRequest{
				IdOneof: &jobv1.RevokeJobRequest_Id{Id: job.JobId},
			})
			if tc.expectedError != "" {
				require.Error(t, err)
				require.Nil(t, resp)
				require.Contains(t, err.Error(), tc.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, job.JobId, resp.Proposal.JobId)
			}
		})
	}
}

func TestGetJob(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Create mock job approvers
	mockGetter := &mockJobApproverGetter{
		jobApprovers: make(map[string]*mockJobApprover),
	}
	mockGetter.jobApprovers["node-1"] = &mockJobApprover{}

	client := NewJobServiceClient(mockGetter)

	// First propose a job to have something to get
	externalJobID := uuid.NewString()
	jobSpec := createValidJobSpec(externalJobID)

	proposeReq := &jobv1.ProposeJobRequest{
		NodeId: "node-1",
		Spec:   jobSpec,
	}

	_, err := client.ProposeJob(ctx, proposeReq)
	require.NoError(t, err)

	t.Run("get existing job", func(t *testing.T) {
		req := &jobv1.GetJobRequest{
			IdOneof: &jobv1.GetJobRequest_Id{
				Id: externalJobID,
			},
		}

		resp, err := client.GetJob(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Job)
		require.Equal(t, externalJobID, resp.Job.Id)
		require.Equal(t, "node-1", resp.Job.NodeId)
	})

	t.Run("get non-existent job", func(t *testing.T) {
		req := &jobv1.GetJobRequest{
			IdOneof: &jobv1.GetJobRequest_Id{
				Id: "non-existent-job",
			},
		}

		resp, err := client.GetJob(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "failed to get job")
	})
}

func TestGetProposal(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Create mock job approvers
	mockGetter := &mockJobApproverGetter{
		jobApprovers: make(map[string]*mockJobApprover),
	}
	mockGetter.jobApprovers["node-1"] = &mockJobApprover{}

	client := NewJobServiceClient(mockGetter)

	// First propose a job to have a proposal to get
	externalJobID := uuid.NewString()
	jobSpec := createValidJobSpec(externalJobID)

	proposeReq := &jobv1.ProposeJobRequest{
		NodeId: "node-1",
		Spec:   jobSpec,
	}

	proposeResp, err := client.ProposeJob(ctx, proposeReq)
	require.NoError(t, err)
	require.NotNil(t, proposeResp)

	proposalID := proposeResp.Proposal.Id

	t.Run("get existing proposal", func(t *testing.T) {
		req := &jobv1.GetProposalRequest{
			Id: proposalID,
		}

		resp, err := client.GetProposal(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Proposal)
		require.Equal(t, proposalID, resp.Proposal.Id)
		require.Equal(t, externalJobID, resp.Proposal.JobId)
	})

	t.Run("get non-existent proposal", func(t *testing.T) {
		req := &jobv1.GetProposalRequest{
			Id: "non-existent-proposal",
		}

		resp, err := client.GetProposal(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "failed to get proposal")
	})
}

func TestListJobs(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Create mock job approvers
	mockGetter := &mockJobApproverGetter{
		jobApprovers: make(map[string]*mockJobApprover),
	}
	mockGetter.jobApprovers["node-1"] = &mockJobApprover{}
	mockGetter.jobApprovers["node-2"] = &mockJobApprover{}

	client := NewJobServiceClient(mockGetter)

	// Propose multiple jobs for different nodes
	externalJobID1 := uuid.NewString()
	jobSpec1 := createValidJobSpec(externalJobID1)

	externalJobID2 := uuid.NewString()
	jobSpec2 := createValidJobSpec(externalJobID2)

	_, err := client.ProposeJob(ctx, &jobv1.ProposeJobRequest{
		NodeId: "node-1",
		Spec:   jobSpec1,
	})
	require.NoError(t, err)

	_, err = client.ProposeJob(ctx, &jobv1.ProposeJobRequest{
		NodeId: "node-2",
		Spec:   jobSpec2,
	})
	require.NoError(t, err)

	t.Run("list all jobs", func(t *testing.T) {
		req := &jobv1.ListJobsRequest{}

		resp, err := client.ListJobs(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Jobs, 2)
	})

	t.Run("filter by node ID", func(t *testing.T) {
		req := &jobv1.ListJobsRequest{
			Filter: &jobv1.ListJobsRequest_Filter{
				NodeIds: []string{"node-1"},
			},
		}

		resp, err := client.ListJobs(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Jobs, 1)
		require.Equal(t, "node-1", resp.Jobs[0].NodeId)
	})

	t.Run("filter by job ID", func(t *testing.T) {
		req := &jobv1.ListJobsRequest{
			Filter: &jobv1.ListJobsRequest_Filter{
				Ids: []string{externalJobID1},
			},
		}

		resp, err := client.ListJobs(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Jobs, 1)
		require.Equal(t, externalJobID1, resp.Jobs[0].Id)
	})

	t.Run("filter by uuid", func(t *testing.T) {
		req := &jobv1.ListJobsRequest{
			Filter: &jobv1.ListJobsRequest_Filter{
				Uuids: []string{externalJobID2},
			},
		}

		resp, err := client.ListJobs(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Jobs, 1)
		require.Equal(t, externalJobID2, resp.Jobs[0].Uuid)
	})

	t.Run("filter by non-existent ID", func(t *testing.T) {
		req := &jobv1.ListJobsRequest{
			Filter: &jobv1.ListJobsRequest_Filter{
				Ids: []string{"non-existent-id"},
			},
		}

		resp, err := client.ListJobs(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Empty(t, resp.Jobs)
	})

	t.Run("invalid filter combination", func(t *testing.T) {
		req := &jobv1.ListJobsRequest{
			Filter: &jobv1.ListJobsRequest_Filter{
				NodeIds: []string{"node-1"},
				Ids:     []string{externalJobID1},
				Uuids:   []string{externalJobID1},
			},
		}

		resp, err := client.ListJobs(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "only one of NodeIds, Uuids or Ids can be set")
	})
}

func TestListProposals(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Create mock job approvers
	mockGetter := &mockJobApproverGetter{
		jobApprovers: make(map[string]*mockJobApprover),
	}
	mockGetter.jobApprovers["node-1"] = &mockJobApprover{}

	client := NewJobServiceClient(mockGetter)

	// Propose multiple jobs to create proposals
	externalJobID1 := uuid.NewString()
	jobSpec1 := createValidJobSpec(externalJobID1)

	externalJobID2 := uuid.NewString()
	jobSpec2 := createValidJobSpec(externalJobID2)

	proposeResp1, err := client.ProposeJob(ctx, &jobv1.ProposeJobRequest{
		NodeId: "node-1",
		Spec:   jobSpec1,
	})
	require.NoError(t, err)
	proposalID1 := proposeResp1.Proposal.Id

	_, err = client.ProposeJob(ctx, &jobv1.ProposeJobRequest{
		NodeId: "node-1",
		Spec:   jobSpec2,
	})
	require.NoError(t, err)
	//	proposalID2 := proposeResp2.Proposal.Id

	t.Run("list all proposals", func(t *testing.T) {
		req := &jobv1.ListProposalsRequest{}

		resp, err := client.ListProposals(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Proposals, 2)
	})

	t.Run("filter by proposal ID", func(t *testing.T) {
		req := &jobv1.ListProposalsRequest{
			Filter: &jobv1.ListProposalsRequest_Filter{
				Ids: []string{proposalID1},
			},
		}

		resp, err := client.ListProposals(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Proposals, 1)
		require.Equal(t, proposalID1, resp.Proposals[0].Id)
	})

	t.Run("filter by job ID", func(t *testing.T) {
		req := &jobv1.ListProposalsRequest{
			Filter: &jobv1.ListProposalsRequest_Filter{
				JobIds: []string{externalJobID2},
			},
		}

		resp, err := client.ListProposals(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Proposals, 1)
		require.Equal(t, externalJobID2, resp.Proposals[0].JobId)
	})

	t.Run("filter by non-existent ID", func(t *testing.T) {
		req := &jobv1.ListProposalsRequest{
			Filter: &jobv1.ListProposalsRequest_Filter{
				Ids: []string{"non-existent-id"},
			},
		}

		resp, err := client.ListProposals(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Empty(t, resp.Proposals)
	})

	t.Run("invalid filter combination", func(t *testing.T) {
		req := &jobv1.ListProposalsRequest{
			Filter: &jobv1.ListProposalsRequest_Filter{
				Ids:    []string{proposalID1},
				JobIds: []string{externalJobID1},
			},
		}

		resp, err := client.ListProposals(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "only one of Ids or JobIds can be set")
	})
}

func TestMapJobStore(t *testing.T) {
	t.Parallel()

	store := newMapJobStore()

	job1 := &jobv1.Job{
		Id:     "job-1",
		Uuid:   "uuid-1",
		NodeId: "node-1",
	}

	job2 := &jobv1.Job{
		Id:     "job-2",
		Uuid:   "uuid-2",
		NodeId: "node-2",
	}

	t.Run("put and get", func(t *testing.T) {
		err := store.put(job1.Id, job1)
		require.NoError(t, err)

		retrieved, err := store.get(job1.Id)
		require.NoError(t, err)
		require.Equal(t, job1, retrieved)

		// Test non-existent job
		_, err = store.get("non-existent")
		require.ErrorIs(t, err, errNoExist)
	})

	t.Run("list", func(t *testing.T) {
		// Add job2
		err := store.put(job2.Id, job2)
		require.NoError(t, err)

		// List all jobs
		jobs, err := store.list(nil)
		require.NoError(t, err)
		require.Len(t, jobs, 2)

		// Filter by node ID
		jobs, err = store.list(&jobv1.ListJobsRequest_Filter{
			NodeIds: []string{"node-1"},
		})
		require.NoError(t, err)
		require.Len(t, jobs, 1)
		require.Equal(t, "node-1", jobs[0].NodeId)

		// Filter by job ID
		jobs, err = store.list(&jobv1.ListJobsRequest_Filter{
			Ids: []string{"job-2"},
		})
		require.NoError(t, err)
		require.Len(t, jobs, 1)
		require.Equal(t, "job-2", jobs[0].Id)

		// Filter by UUID
		jobs, err = store.list(&jobv1.ListJobsRequest_Filter{
			Uuids: []string{"uuid-1"},
		})
		require.NoError(t, err)
		require.Len(t, jobs, 1)
		require.Equal(t, "uuid-1", jobs[0].Uuid)
	})

	t.Run("delete", func(t *testing.T) {
		err := store.delete(job1.Id)
		require.NoError(t, err)

		_, err = store.get(job1.Id)
		require.ErrorIs(t, err, errNoExist)

		// Deleting non-existent job should not error
		err = store.delete("non-existent")
		require.NoError(t, err)
	})
}

func TestMapProposalStore(t *testing.T) {
	t.Parallel()

	store := newMapProposalStore()

	proposal1 := &jobv1.Proposal{
		Id:    "proposal-1",
		JobId: "job-1",
	}

	proposal2 := &jobv1.Proposal{
		Id:    "proposal-2",
		JobId: "job-2",
	}

	t.Run("put and get", func(t *testing.T) {
		err := store.put(proposal1.Id, proposal1)
		require.NoError(t, err)

		retrieved, err := store.get(proposal1.Id)
		require.NoError(t, err)
		require.Equal(t, proposal1, retrieved)

		// Test non-existent proposal
		_, err = store.get("non-existent")
		require.Error(t, err)
		require.Contains(t, err.Error(), "proposal not found")
	})

	t.Run("list", func(t *testing.T) {
		// Add proposal2
		err := store.put(proposal2.Id, proposal2)
		require.NoError(t, err)

		// List all proposals
		proposals, err := store.list(nil)
		require.NoError(t, err)
		require.Len(t, proposals, 2)

		// Filter by proposal ID
		proposals, err = store.list(&jobv1.ListProposalsRequest_Filter{
			Ids: []string{"proposal-1"},
		})
		require.NoError(t, err)
		require.Len(t, proposals, 1)
		require.Equal(t, "proposal-1", proposals[0].Id)

		// Filter by job ID
		proposals, err = store.list(&jobv1.ListProposalsRequest_Filter{
			JobIds: []string{"job-2"},
		})
		require.NoError(t, err)
		require.Len(t, proposals, 1)
		require.Equal(t, "job-2", proposals[0].JobId)
	})

	t.Run("delete", func(t *testing.T) {
		err := store.delete(proposal1.Id)
		require.NoError(t, err)

		_, err = store.get(proposal1.Id)
		require.Error(t, err)
		require.Contains(t, err.Error(), "proposal not found")
	})
}

func TestMatchesSelectors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		selectors []*ptypes.Selector
		job       *jobv1.Job
		expected  bool
	}{
		// SelectorOp_EXIST cases
		{
			name: "SelectorOp_EXIST match",
			selectors: []*ptypes.Selector{
				{
					Key: "label1",
					Op:  ptypes.SelectorOp_EXIST,
				},
			},
			job: &jobv1.Job{
				Labels: []*ptypes.Label{
					{
						Key: "label1",
					},
				},
			},
			expected: true,
		},
		{
			name: "SelectorOp_EXIST does not exist",
			selectors: []*ptypes.Selector{
				{
					Key: "label1",
					Op:  ptypes.SelectorOp_EXIST,
				},
			},
			job:      &jobv1.Job{},
			expected: false,
		},
		// SelectorOp_NOT_EXIST cases
		{
			name: "SelectorOp_NOT_EXIST match",
			selectors: []*ptypes.Selector{
				{
					Key: "label1",
					Op:  ptypes.SelectorOp_NOT_EXIST,
				},
			},
			job: &jobv1.Job{
				Labels: []*ptypes.Label{
					{
						Key: "label1",
					},
				},
			},
			expected: false,
		},
		{
			name: "SelectorOp_NOT_EXIST does not exist",
			selectors: []*ptypes.Selector{
				{
					Key: "label1",
					Op:  ptypes.SelectorOp_NOT_EXIST,
				},
			},
			job:      &jobv1.Job{},
			expected: true,
		},
		// SelectorOp_EQ cases
		{
			name: "SelectorOp_EQ match",
			selectors: []*ptypes.Selector{
				{
					Key:   "label1",
					Op:    ptypes.SelectorOp_EQ,
					Value: pointer.To("value1"),
				},
			},
			job: &jobv1.Job{
				Labels: []*ptypes.Label{
					{
						Key:   "label1",
						Value: pointer.To("value1"),
					},
				},
			},
			expected: true,
		},
		{
			name: "SelectorOp_EQ mismatched label value",
			selectors: []*ptypes.Selector{
				{
					Key:   "label1",
					Op:    ptypes.SelectorOp_EQ,
					Value: pointer.To("value1"),
				},
			},
			job: &jobv1.Job{
				Labels: []*ptypes.Label{
					{
						Key:   "label1",
						Value: pointer.To("NOT THE VALUE WE NEED"),
					},
				},
			},
			expected: false,
		},
		{
			name: "SelectorOp_EQ with missing label",
			selectors: []*ptypes.Selector{
				{
					Key:   "label1",
					Op:    ptypes.SelectorOp_EQ,
					Value: pointer.To("value1"),
				},
			},
			job:      &jobv1.Job{},
			expected: false,
		},
		// SelectorOp_NOT_EQ cases
		{
			name: "SelectorOp_NOT_EQ match",
			selectors: []*ptypes.Selector{
				{
					Key:   "label1",
					Op:    ptypes.SelectorOp_NOT_EQ,
					Value: pointer.To("value1"),
				},
			},
			job: &jobv1.Job{
				Labels: []*ptypes.Label{
					{
						Key:   "label1",
						Value: pointer.To("value1"),
					},
				},
			},
			expected: false,
		},
		{
			name: "SelectorOp_NOT_EQ mismatched label value",
			selectors: []*ptypes.Selector{
				{
					Key:   "label1",
					Op:    ptypes.SelectorOp_NOT_EQ,
					Value: pointer.To("value1"),
				},
			},
			job: &jobv1.Job{
				Labels: []*ptypes.Label{
					{
						Key:   "label1",
						Value: pointer.To("NOT THE VALUE WE NEED"),
					},
				},
			},
			expected: true,
		},
		{
			name: "SelectorOp_NOT_EQ with missing label",
			selectors: []*ptypes.Selector{
				{
					Key:   "label1",
					Op:    ptypes.SelectorOp_NOT_EQ,
					Value: pointer.To("value1"),
				},
			},
			job:      &jobv1.Job{},
			expected: true,
		},
		// SelectorOp_IN cases
		{
			name: "SelectorOp_IN match",
			selectors: []*ptypes.Selector{
				{
					Key:   "label1",
					Op:    ptypes.SelectorOp_IN,
					Value: pointer.To("value1,value2"),
				},
			},
			job: &jobv1.Job{
				Labels: []*ptypes.Label{
					{
						Key:   "label1",
						Value: pointer.To("value1"),
					},
				},
			},
			expected: true,
		},
		{
			name: "SelectorOp_IN mismatched label value",
			selectors: []*ptypes.Selector{
				{
					Key:   "label1",
					Op:    ptypes.SelectorOp_IN,
					Value: pointer.To("value1,value2"),
				},
			},
			job: &jobv1.Job{
				Labels: []*ptypes.Label{
					{
						Key:   "label1",
						Value: pointer.To("NOT THE VALUE WE NEED"),
					},
				},
			},
			expected: false,
		},
		{
			name: "SelectorOp_IN with missing label",
			selectors: []*ptypes.Selector{
				{
					Key:   "label1",
					Op:    ptypes.SelectorOp_IN,
					Value: pointer.To("value1"),
				},
			},
			job:      &jobv1.Job{},
			expected: false,
		},
		// SelectorOp_NOT_IN cases
		{
			name: "SelectorOp_NOT_IN match",
			selectors: []*ptypes.Selector{
				{
					Key:   "label1",
					Op:    ptypes.SelectorOp_NOT_IN,
					Value: pointer.To("value1,value2"),
				},
			},
			job: &jobv1.Job{
				Labels: []*ptypes.Label{
					{
						Key:   "label1",
						Value: pointer.To("value1"),
					},
				},
			},
			expected: false,
		},
		{
			name: "SelectorOp_NOT_IN mismatched label value",
			selectors: []*ptypes.Selector{
				{
					Key:   "label1",
					Op:    ptypes.SelectorOp_NOT_IN,
					Value: pointer.To("value1,value2"),
				},
			},
			job: &jobv1.Job{
				Labels: []*ptypes.Label{
					{
						Key:   "label1",
						Value: pointer.To("NOT THE VALUE WE NEED"),
					},
				},
			},
			expected: true,
		},
		{
			name: "SelectorOp_NOT_IN with missing label",
			selectors: []*ptypes.Selector{
				{
					Key:   "label1",
					Op:    ptypes.SelectorOp_NOT_IN,
					Value: pointer.To("value1"),
				},
			},
			job:      &jobv1.Job{},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := matchesSelectors(tc.selectors, tc.job)
			require.Equal(t, tc.expected, result)
		})
	}
}

// need some non-ocr job type to avoid the ocr validation and the p2pwrapper check
func createValidJobSpec(externalJobID string) string {
	tomlString := `
type = "standardcapabilities"
schemaVersion = 1
externalJobID = "%s"
name = "hacking-%s"
forwardingAllowed = false
command = "/home/capabilities/nowhere"
config = ""
`
	return fmt.Sprintf(tomlString, externalJobID, externalJobID)
}

// mockJobApprover is a mock implementation of the JobApprover interface
type mockJobApprover struct {
	shouldFail bool
}

func (m *mockJobApprover) AutoApproveJob(ctx context.Context, p *feeds.ProposeJobArgs) error {
	if m.shouldFail {
		return errors.New("mock approval failure")
	}
	return nil
}

// mockJobApproverGetter is a mock implementation of the getter[JobApprover] interface
type mockJobApproverGetter struct {
	jobApprovers map[string]*mockJobApprover
}

func (m *mockJobApproverGetter) Get(id string) (JobApprover, error) {
	approver, ok := m.jobApprovers[id]
	if !ok {
		return nil, errors.New("job approver not found")
	}
	return approver, nil
}
