package changeset

import (
	"errors"
	"fmt"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink/deployment"
)

var _ deployment.ChangeSetV2[CsRevokeJobSpecsConfig] = CsRevokeJobSpecs{}

type CsRevokeJobSpecsConfig struct {
	// UUIDs is a list of external job IDs to revoke.
	UUIDs []string
}

type CsRevokeJobSpecs struct{}

func (CsRevokeJobSpecs) Apply(e deployment.Environment, cfg CsRevokeJobSpecsConfig) (deployment.ChangesetOutput, error) {
	// Fetch the internal job IDs from the job distributor:
	jobsResp, err := e.Offchain.ListJobs(e.GetContext(), &jobv1.ListJobsRequest{
		Filter: &jobv1.ListJobsRequest_Filter{
			Uuids: cfg.UUIDs,
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to list jobs: %w", err)
	}
	if len(jobsResp.Jobs) != len(cfg.UUIDs) {
		return deployment.ChangesetOutput{}, errors.New("failed to find jobs for all provided UUIDs")
	}

	revokedJobs := make([]deployment.ProposedJob, 0, len(jobsResp.Jobs))
	for _, job := range jobsResp.Jobs {
		resp, err := e.Offchain.RevokeJob(e.GetContext(), &jobv1.RevokeJobRequest{
			IdOneof: &jobv1.RevokeJobRequest_Id{
				Id: job.GetId(),
			},
		})
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to revoke job: %w", err)
		}
		revokedJobs = append(revokedJobs, deployment.ProposedJob{
			JobID: resp.GetProposal().GetJobId(),
			Spec:  resp.GetProposal().GetSpec(),
		})
	}

	return deployment.ChangesetOutput{
		Jobs: revokedJobs,
	}, nil
}

func (f CsRevokeJobSpecs) VerifyPreconditions(_ deployment.Environment, config CsRevokeJobSpecsConfig) error {
	if len(config.UUIDs) == 0 {
		return errors.New("job ids are required")
	}

	return nil
}
