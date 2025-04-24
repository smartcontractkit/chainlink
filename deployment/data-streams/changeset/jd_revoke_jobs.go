package changeset

import (
	"errors"
	"fmt"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink/deployment"
)

var _ deployment.ChangeSetV2[CsRevokeJobSpecsConfig] = CsRevokeJobSpecs{}

type CsRevokeJobSpecsConfig struct {
	JobID string
}

type CsRevokeJobSpecs struct{}

func (CsRevokeJobSpecs) Apply(e deployment.Environment, cfg CsRevokeJobSpecsConfig) (deployment.ChangesetOutput, error) {
	resp, err := e.Offchain.RevokeJob(e.GetContext(), &jobv1.RevokeJobRequest{
		IdOneof: &jobv1.RevokeJobRequest_Id{
			Id: cfg.JobID,
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to revoke job: %w", err)
	}

	return deployment.ChangesetOutput{
		Jobs: []deployment.ProposedJob{
			{
				JobID: resp.GetProposal().GetJobId(),
				Spec:  resp.GetProposal().GetSpec(),
			},
		},
	}, nil
}

func (f CsRevokeJobSpecs) VerifyPreconditions(_ deployment.Environment, config CsRevokeJobSpecsConfig) error {
	if config.JobID == "" {
		return errors.New("job id is required")
	}

	return nil
}
