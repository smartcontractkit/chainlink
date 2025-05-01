package changeset

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
)

const (
	defaultJobSpecsTimeout = 120 * time.Second
)

func chainAndAddresses(e deployment.Environment, chainSel uint64) (chainID string, addresses map[string]deployment.TypeAndVersion, err error) {
	chainID, err = chainsel.GetChainIDFromSelector(chainSel)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get chain ID from selector: %w", err)
	}

	addresses, err = e.ExistingAddresses.AddressesForChain(chainSel)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	return chainID, addresses, nil
}

// proposeAllOrNothing proposes all jobs in the list and if any of them fail, it will revoke all already made proposals.
func proposeAllOrNothing(ctx context.Context, oc cldf.OffchainClient, prs []*job.ProposeJobRequest) (proposedJobs []deployment.ProposedJob, err error) {
	var proposals []*job.ProposeJobResponse
	var p *job.ProposeJobResponse
	for _, pr := range prs {
		p, err = oc.ProposeJob(ctx, pr)
		if err != nil {
			break
		}
		proposedJobs = append(proposedJobs, deployment.ProposedJob{
			JobID: p.Proposal.JobId,
			Node:  pr.NodeId,
			Spec:  pr.Spec,
		})
		proposals = append(proposals, p)
	}

	if err != nil {
		// There's an error, so we need to revoke all proposals we just made.
		var errs []error
		for _, pr := range proposals {
			if _, errRevoke := oc.RevokeJob(ctx, &job.RevokeJobRequest{
				IdOneof: &job.RevokeJobRequest_Id{Id: pr.Proposal.JobId},
			}); errRevoke != nil {
				errs = append(errs, fmt.Errorf("failed to revoke job %s: %w", pr.Proposal.JobId, errRevoke))
			}
		}
		// If we got any errors while trying to cancel, we need to return them, so we know we sent some jobs twice.
		if len(errs) > 0 {
			err = errors.Join(err, errors.Join(errs...))
		}
	}

	return proposedJobs, err
}

// fetchExternalJobID looks for an existing job that matches the given labels and returns its ID.
// If no job is found, it returns a nil UUID.
func fetchExternalJobID(e deployment.Environment, nodeID string, selectors []*ptypes.Selector) (externalJobID uuid.UUID, err error) {
	var nodeIDs []string
	if nodeID != "" {
		nodeIDs = []string{nodeID}
	}
	jobsResp, err := e.Offchain.ListJobs(e.GetContext(), &job.ListJobsRequest{
		Filter: &job.ListJobsRequest_Filter{
			NodeIds:   nodeIDs,
			Selectors: selectors,
		},
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to list jobs: %w", err)
	}

	switch len(jobsResp.Jobs) {
	case 0:
		// No job found, return nil UUID
	case 1:
		// One job found, return its ID
		externalJobID, err = uuid.Parse(jobsResp.Jobs[0].Uuid)
		if err != nil {
			err = fmt.Errorf("failed to parse external job ID: %w", err)
		}
	default:
		// More than one job found, return error
		err = fmt.Errorf("multiple jobs found: %d", len(jobsResp.Jobs))
	}
	return externalJobID, err
}
