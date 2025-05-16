package jobs

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
)

var _ cldf.ChangeSetV2[CsRevokeJobSpecsConfig] = CsRevokeJobSpecs{}

// CsRevokeJobSpecsConfig is the configuration for the revoking a job.
// In order to revoke a job, we need to know one of two things:
// 1. The external job ID (UUID) of the job.
// 2. The stream ID to which the job belongs.
//
// Note that only one set of IDs (UUIDs or stream IDs) is allowed.
type CsRevokeJobSpecsConfig struct {
	// UUIDs is a list of external job IDs to revoke.
	UUIDs []string

	StreamIDs []uint32
}

type CsRevokeJobSpecs struct{}

func (CsRevokeJobSpecs) Apply(e cldf.Environment, cfg CsRevokeJobSpecsConfig) (cldf.ChangesetOutput, error) {
	jobs, err := findJobsForIDs(e, cfg.UUIDs, cfg.StreamIDs)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to find jobs: %w", err)
	}

	revokedJobs := make([]cldf.ProposedJob, 0, len(jobs))
	for _, job := range jobs {
		resp, err := e.Offchain.RevokeJob(e.GetContext(), &jobv1.RevokeJobRequest{
			IdOneof: &jobv1.RevokeJobRequest_Id{
				Id: job.GetId(),
			},
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to revoke job: %w", err)
		}
		revokedJobs = append(revokedJobs, cldf.ProposedJob{
			JobID: resp.GetProposal().GetJobId(),
			Spec:  resp.GetProposal().GetSpec(),
		})
	}

	return cldf.ChangesetOutput{
		Jobs: revokedJobs,
	}, nil
}

// findJobsForIDs finds jobs for either the given UUIDs or stream IDs.
func findJobsForIDs(e cldf.Environment, uuids []string, streamIDs []uint32) ([]*jobv1.Job, error) {
	if (len(uuids) == 0 && len(streamIDs) == 0) || (len(uuids) > 0 && len(streamIDs) > 0) {
		return nil, errors.New("either job ids or stream ids are required")
	}
	if len(uuids) > 0 {
		return findJobsForUUIDs(e, uuids)
	}
	return findJobsForStreamIDs(e, streamIDs)
}

func findJobsForUUIDs(e cldf.Environment, uuids []string) ([]*jobv1.Job, error) {
	// Fetch the internal job IDs from the job distributor:
	jobsResp, err := e.Offchain.ListJobs(e.GetContext(), &jobv1.ListJobsRequest{
		Filter: &jobv1.ListJobsRequest_Filter{
			Uuids: uuids,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}
	if len(jobsResp.Jobs) != len(uuids) {
		return nil, errors.New("failed to find jobs for all provided UUIDs")
	}
	return jobsResp.Jobs, nil
}

func findJobsForStreamIDs(e cldf.Environment, streamIDs []uint32) ([]*jobv1.Job, error) {
	var jobs []*jobv1.Job
	// We need to collect the jobs for each stream ID separately because the label we use is a flag and we cannot
	// select by "OR" logic.
	for _, sid := range streamIDs {
		jobsResp, err := e.Offchain.ListJobs(e.GetContext(), &jobv1.ListJobsRequest{
			Filter: &jobv1.ListJobsRequest_Filter{
				Selectors: []*ptypes.Selector{
					{
						Key: utils.StreamIDLabel(sid),
						Op:  ptypes.SelectorOp_EXIST,
					},
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list jobs: %w", err)
		}
		jobs = append(jobs, jobsResp.Jobs...)
	}
	return jobs, nil
}

func (f CsRevokeJobSpecs) VerifyPreconditions(_ cldf.Environment, config CsRevokeJobSpecsConfig) error {
	if (len(config.UUIDs) == 0 && len(config.StreamIDs) == 0) || (len(config.UUIDs) > 0 && len(config.StreamIDs) > 0) {
		return errors.New("either job ids or stream ids are required")
	}

	return nil
}
