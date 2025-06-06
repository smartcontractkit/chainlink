package changeset

import (
	"errors"
	"fmt"
	"slices"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

var (
	// RevokeJobsChangeset revokes job proposals with the given jobIDs through JD. It can only be used on
	// the proposals that are in the pending and cancelled state.
	RevokeJobsChangeset = cldf.CreateChangeSet(revokeJobsLogic, revokeJobsPrecondition)

	// DeleteJobChangeset sends a delete request to the node where the job is running and marks it as deleted in Job Distributor.
	// If the node is not connected or the delete request fails, the deletion process is halted.
	// Nodes are expected to cancel the job once the request is sent by JD.
	// Refer to integration-tests/smoke/ccip/ccip_jobspec_test.go for node operations example after DeleteJobChangeset.
	DeleteJobChangeset = cldf.CreateChangeSet(deleteJobsLogic, deleteJobsPrecondition)
)

func revokeJobsPrecondition(env cldf.Environment, jobIDs []string) error {
	proposals, err := env.Offchain.ListProposals(env.GetContext(), &jobv1.ListProposalsRequest{
		Filter: &jobv1.ListProposalsRequest_Filter{
			JobIds: jobIDs,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to list proposals for jobIDs %s: %w", jobIDs, err)
	}
	for _, proposal := range proposals.Proposals {
		if proposal.Status != jobv1.ProposalStatus_PROPOSAL_STATUS_PROPOSED && proposal.Status != jobv1.ProposalStatus_PROPOSAL_STATUS_CANCELLED {
			return fmt.Errorf("proposal %s is not in PROPOSED or CANCELLED state", proposal.Id)
		}
	}
	return nil
}

func revokeJobsLogic(env cldf.Environment, jobIDs []string) (cldf.ChangesetOutput, error) {
	var successfullyRevoked []string
	for _, jobID := range jobIDs {
		res, err := env.Offchain.RevokeJob(env.GetContext(), &jobv1.RevokeJobRequest{
			IdOneof: &jobv1.RevokeJobRequest_Id{Id: jobID},
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to revoke job %s: %w", jobID, err)
		}
		if res == nil {
			return cldf.ChangesetOutput{}, errors.New("revoke job response is nil")
		}
		if res.Proposal == nil || res.Proposal.Status != jobv1.ProposalStatus_PROPOSAL_STATUS_REVOKED {
			return cldf.ChangesetOutput{}, errors.New("revoke job response is not in cancelled state")
		}
		successfullyRevoked = append(successfullyRevoked, jobID)
	}
	if len(successfullyRevoked) == 0 {
		return cldf.ChangesetOutput{}, errors.New("no jobs were revoked")
	}
	if len(successfullyRevoked) != len(jobIDs) {
		return cldf.ChangesetOutput{}, fmt.Errorf("not all jobs were revoked, successfully revoked %s, expected %s", successfullyRevoked, jobIDs)
	}
	env.Logger.Infof("successfully revoked jobs %s", successfullyRevoked)
	return cldf.ChangesetOutput{}, nil
}

func deleteJobsPrecondition(env cldf.Environment, jobIDs []string) error {
	jobs, err := env.Offchain.ListJobs(env.GetContext(), &jobv1.ListJobsRequest{
		Filter: &jobv1.ListJobsRequest_Filter{
			Ids: jobIDs,
		},
	})
	if err != nil {
		return err
	}
	if len(jobs.Jobs) != len(jobIDs) {
		var found []string
		for _, job := range jobs.Jobs {
			if job.DeletedAt != nil {
				return fmt.Errorf("job %s is already deleted", job.Id)
			}
			found = append(found, job.Id)
		}
		return fmt.Errorf("not all jobs found in listJobs response, returned jobs with ids %s, expected %s", found, jobIDs)
	}
	return nil
}

// DeleteJobChangeset sends the delete job request to nodes for the given jobID.
// nops needs to cancel the job once the request is sent by JD.
func deleteJobsLogic(env cldf.Environment, jobIDs []string) (cldf.ChangesetOutput, error) {
	jobIDsToDelete, err := jobsToDelete(env, jobIDs)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get jobIDs to delete: %w", err)
	}
	for _, jobID := range jobIDsToDelete {
		res, err := env.Offchain.DeleteJob(env.GetContext(), &jobv1.DeleteJobRequest{
			IdOneof: &jobv1.DeleteJobRequest_Id{Id: jobID},
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to delete job %s: %w", jobID, err)
		}
		if res == nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("delete job response is nil for job %s", jobID)
		}
		if res.Job == nil || res.Job.DeletedAt == nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("delete job response is not in deleted state for job %s", jobID)
		}
	}
	env.Logger.Infof("successfully deleted jobs %s", jobIDs)
	return cldf.ChangesetOutput{}, nil
}

func jobsToDelete(env cldf.Environment, jobIDs []string) ([]string, error) {
	jobs, err := env.Offchain.ListProposals(env.GetContext(), &jobv1.ListProposalsRequest{
		Filter: &jobv1.ListProposalsRequest_Filter{
			JobIds: jobIDs,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list proposals for jobIDs %v: %w", jobIDs, err)
	}
	if len(jobs.Proposals) == 0 {
		return nil, fmt.Errorf("no proposals found for jobIDs %s", jobIDs)
	}
	jobIDsToDelete := make([]string, 0)
	for _, proposal := range jobs.Proposals {
		if proposal.Status == jobv1.ProposalStatus_PROPOSAL_STATUS_PROPOSED ||
			proposal.Status == jobv1.ProposalStatus_PROPOSAL_STATUS_APPROVED ||
			proposal.Status == jobv1.ProposalStatus_PROPOSAL_STATUS_PENDING {
			jobIDsToDelete = append(jobIDsToDelete, proposal.JobId)
		}
	}
	// remove duplicates
	jobIDsToDelete = slices.Compact(jobIDsToDelete)
	return jobIDsToDelete, nil
}
