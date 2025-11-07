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

func Approve(ctx context.Context, offChainClient cldf_offchain.Client, dons *cre.Dons, nodeToSpecs map[string][]string) error {
	for nodeID, jobSpecs := range nodeToSpecs {
		for _, don := range dons.List() {
			for _, node := range don.Nodes {
				if node.JobDistributorDetails.NodeID != nodeID {
					continue
				}

				for _, jobSpec := range jobSpecs {
					if err := accept(ctx, node, jobSpec); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
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

					if err := accept(ctx, node, jobReq.Spec); err != nil {
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

func accept(ctx context.Context, node *cre.Node, jobSpec string) error {
	// TODO: is there a way to accept the job with proposal id?
	if err := node.AcceptJob(ctx, jobSpec); err != nil {
		// Workflow specs get auto approved
		// TODO: Narrow down scope by checking type == workflow
		if strings.Contains(err.Error(), "cannot approve an approved spec") {
			return nil
		}
		fmt.Println("Failed jobspec proposal:")
		fmt.Println(jobSpec)

		return fmt.Errorf("failed to accept job. err: %w", err)
	}

	return nil
}
