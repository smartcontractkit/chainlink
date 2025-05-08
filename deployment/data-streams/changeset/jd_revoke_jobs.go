package changeset

import (
	"errors"
	"fmt"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink/deployment"
)

var _ deployment.ChangeSetV2[CsRevokeJobSpecsConfig] = CsRevokeJobSpecs{}

type CsRevokeJobSpecsConfig struct {
	JobIDs []string
}

type CsRevokeJobSpecs struct{}

func (CsRevokeJobSpecs) Apply(e deployment.Environment, cfg CsRevokeJobSpecsConfig) (deployment.ChangesetOutput, error) {
	revokedJobs := make([]deployment.ProposedJob, 0, len(cfg.JobIDs))
	for _, jobID := range cfg.JobIDs {
		resp, err := e.Offchain.RevokeJob(e.GetContext(), &jobv1.RevokeJobRequest{
			IdOneof: &jobv1.RevokeJobRequest_Id{
				Id: jobID,
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
	if len(config.JobIDs) == 0 {
		return errors.New("job ids are required")
	}

	return nil
}
