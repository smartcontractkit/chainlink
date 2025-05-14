package changeset

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
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
	var filter *jobv1.ListJobsRequest_Filter
	switch {
	case len(cfg.UUIDs) > 0 && len(cfg.StreamIDs) == 0:
		filter = &jobv1.ListJobsRequest_Filter{
			Uuids: cfg.UUIDs,
		}
	case len(cfg.StreamIDs) > 0 && len(cfg.UUIDs) == 0:
		ids := make([]string, len(cfg.StreamIDs))
		for i, id := range cfg.StreamIDs {
			ids[i] = strconv.FormatUint(uint64(id), 10)
		}
		filter = &jobv1.ListJobsRequest_Filter{
			Selectors: []*ptypes.Selector{
				{
					Key:   devenv.LabelStreamIDKey,
					Op:    ptypes.SelectorOp_IN,
					Value: pointer.To(strings.Join(ids, ",")),
				},
			},
		}
	default:
		return cldf.ChangesetOutput{}, errors.New("either job ids or stream ids are required")
	}

	// Fetch the internal job IDs from the job distributor:
	jobsResp, err := e.Offchain.ListJobs(e.GetContext(), &jobv1.ListJobsRequest{
		Filter: filter,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to list jobs: %w", err)
	}
	if len(cfg.UUIDs) > 0 && len(jobsResp.Jobs) != len(cfg.UUIDs) {
		return cldf.ChangesetOutput{}, errors.New("failed to find jobs for all provided UUIDs")
	}

	revokedJobs := make([]cldf.ProposedJob, 0, len(jobsResp.Jobs))
	for _, job := range jobsResp.Jobs {
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

func (f CsRevokeJobSpecs) VerifyPreconditions(_ cldf.Environment, config CsRevokeJobSpecsConfig) error {
	if (len(config.UUIDs) == 0 && len(config.StreamIDs) == 0) || (len(config.UUIDs) > 0 && len(config.StreamIDs) > 0) {
		return errors.New("either job ids or stream ids are required")
	}

	return nil
}
