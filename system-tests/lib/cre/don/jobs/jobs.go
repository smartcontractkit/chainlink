package jobs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/ratelimit"
	"golang.org/x/sync/errgroup"

	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

// defined as variables to allow for easy testing
var loadNodeProposalIDs = func(ctx context.Context, node *cre.Node) (map[string]string, error) {
	jd, err := node.Clients.GQLClient.GetJobDistributor(ctx, node.JobDistributorDetails.JDID)
	if err != nil {
		return nil, err
	}
	if jd.GetJobProposals() == nil {
		return nil, fmt.Errorf("no job proposals found for node %s", node.Name)
	}

	proposalIDsBySpec := make(map[string]string, len(jd.JobProposals))
	for _, proposal := range jd.JobProposals {
		proposalIDsBySpec[proposal.LatestSpec.Definition] = proposal.Id
	}

	return proposalIDsBySpec, nil
}

// defined as variables to allow for easy testing
var approveJobProposalSpec = func(ctx context.Context, node *cre.Node, proposalID string) error {
	approvedSpec, err := node.Clients.GQLClient.ApproveJobProposalSpec(ctx, proposalID, false)
	if err != nil {
		return err
	}
	if approvedSpec == nil {
		return fmt.Errorf("no job proposal spec found for job id %s", proposalID)
	}

	return nil
}

func Approve(ctx context.Context, _ cldf_offchain.Client, dons *cre.Dons, nodeToSpecs map[string][]string) error {
	nodeByID := make(map[string]*cre.Node)
	for _, don := range dons.List() {
		for _, node := range don.Nodes {
			nodeByID[node.JobDistributorDetails.NodeID] = node
		}
	}

	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(4)

	for nodeID, jobSpecs := range nodeToSpecs {
		node, ok := nodeByID[nodeID]
		if !ok {
			return fmt.Errorf("node with id %s not found", nodeID)
		}

		eg.Go(func() error {
			proposalIDsBySpec, err := loadNodeProposalIDs(egCtx, node)
			if err != nil {
				return err
			}

			for _, jobSpec := range jobSpecs {
				proposalID, ok := proposalIDsBySpec[jobSpec]
				if !ok {
					return fmt.Errorf("no job proposal found for job spec %s", jobSpec)
				}
				if err := accept(egCtx, node, proposalID, jobSpec); err != nil {
					return err
				}
			}

			return nil
		})
	}

	return eg.Wait()
}

func Create(ctx context.Context, offChainClient cldf_offchain.Client, dons *cre.Dons, jobSpecs cre.DonJobs) error {
	if len(jobSpecs) == 0 {
		return nil
	}

	eg := &errgroup.Group{}
	jobRateLimit := ratelimit.New(5)

	for _, jobReq := range jobSpecs {
		eg.Go(func() error {
			jobRateLimit.Take()
			timeout := time.Second * 60
			ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			_, pErr := offChainClient.ProposeJob(ctxWithTimeout, jobReq)
			if pErr != nil {
				fmt.Println("Failed jobspec proposal:")
				fmt.Println(jobReq)
				return fmt.Errorf("failed to propose job for node %s: %w", jobReq.NodeId, pErr)
			}

			for _, don := range dons.List() {
				for _, node := range don.Nodes {
					if node.JobDistributorDetails.NodeID != jobReq.NodeId {
						continue
					}

					if err := accept(ctx, node, "", jobReq.Spec); err != nil {
						return err
					}
				}
			}

			if ctx.Err() != nil {
				return errors.Wrapf(pErr, "timed out after %s proposing job for node %s", timeout.String(), jobReq.NodeId)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, "failed to create at least one job for DON")
	}

	return nil
}

func accept(ctx context.Context, node *cre.Node, proposalID, jobSpec string) error {
	var err error
	if proposalID == "" {
		err = node.AcceptJob(ctx, jobSpec)
	} else {
		err = approveJobProposalSpec(ctx, node, proposalID)
	}
	if err != nil {
		// Workflow specs get auto approved
		if strings.Contains(err.Error(), "cannot approve an approved spec") && strings.Contains(jobSpec, `type = "workflow"`) {
			return nil
		}
		fmt.Println("Failed jobspec proposal for node ", node.Name)
		fmt.Println(jobSpec)

		return fmt.Errorf("failed to accept job for node %s. err: %w", node.Name, err)
	}

	return nil
}
