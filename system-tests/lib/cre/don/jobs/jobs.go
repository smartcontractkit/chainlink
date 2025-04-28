package jobs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/ratelimit"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func Create(offChainClient deployment.OffchainClient, don *devenv.DON, flags []string, jobSpecs types.DonJobs) error {
	if len(jobSpecs) == 0 {
		return nil
	}

	eg := &errgroup.Group{}
	jobRateLimit := ratelimit.New(5)

	for _, jobReq := range jobSpecs {
		eg.Go(func() error {
			jobRateLimit.Take()
			timeout := time.Second * 60
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			_, err := offChainClient.ProposeJob(ctx, jobReq)
			if err != nil {
				// Workflow specs get auto approved
				// TODO: Narrow down scope by checking type == workflow
				if strings.Contains(err.Error(), "cannot approve an approved spec") {
					return nil
				}
				fmt.Println("Failed jobspec proposal:")
				fmt.Println(jobReq)
				return errors.Wrapf(err, "failed to propose job for node %s", jobReq.NodeId)
			}
			if ctx.Err() != nil {
				return errors.Wrapf(err, "timed out after %s proposing job for node %s", timeout.String(), jobReq.NodeId)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, "failed to create at least one job for DON")
	}

	return nil
}
