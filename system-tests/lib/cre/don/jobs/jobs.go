package jobs

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/ratelimit"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func CreateWithRetry(offChainClient deployment.OffchainClient, don *devenv.DON, flags []string, jobSpecs types.DonJobs) error {
	if len(jobSpecs) == 0 {
		return nil
	}

	if unknownErr := checkForUnknownJobs(jobSpecs); unknownErr != nil {
		return errors.Wrap(unknownErr, "failed to create jobs")
	}

	//errCh := make(chan error, calculateJobCount(jobSpecs))
	//var wg sync.WaitGroup

	const (
		maxRetries = 5
	)

	for i, jobDesc := range SupportedJobs {
		if keystoneflags.HasFlag(flags, jobDesc.Flag) {
			if jobReqs, ok := jobSpecs[jobDesc]; ok {
				for i2, jobReq := range jobReqs {
					//wg.Add(1)
					//go func(jobReq *jobv1.ProposeJobRequest, jobDesc types.JobDescription) {
					//	defer wg.Done()
					time.Sleep(time.Second * 2)
					framework.L.Info().Msgf("%d/%d %d/%d Creating job on %s", i, len(SupportedJobs), i2, len(jobReqs), jobReq.NodeId)
					var lastErr error
					for attempt := 1; attempt <= maxRetries; attempt++ {
						framework.L.Info().Msgf("attempt %d", attempt)
						ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
						_, err := offChainClient.ProposeJob(ctx, jobReq)
						cancel()

						if err == nil {
							break
						}

						if strings.Contains(err.Error(), "cannot approve an approved spec") {
							break
						}

						lastErr = errors2.Join(lastErr, err)
						framework.L.Error().Msgf("Got error %s", err.Error())
					}

					//errCh <- errors.Wrapf(lastErr, "failed to propose job %s for node %s after %d attempts", jobDesc.Flag, jobReq.NodeId, maxRetries)
					//}(jobReq, jobDesc)

					if lastErr != nil {

						//return lastErr
					}
				}
			}
		}
	}

	//wg.Wait()
	//close(errCh)

	//var finalErr error
	//for err := range errCh {
	//	if finalErr == nil {
	//		finalErr = err
	//	} else {
	//		finalErr = errors.Wrap(finalErr, err.Error())
	//	}
	//}
	//
	//if finalErr != nil {
	//	return errors.Wrap(finalErr, "failed to create at least one job for DON")
	//}

	return nil
}

func Create(offChainClient deployment.OffchainClient, don *devenv.DON, flags []string, jobSpecs types.DonJobs) error {
	if len(jobSpecs) == 0 {
		return nil
	}

	eg := &errgroup.Group{}
	jobRateLimit := ratelimit.New(5)

	for jobDesc, jobReqs := range jobSpecs {
		for _, jobReq := range jobReqs {
			eg.Go(func() error {
				jobRateLimit.Take()
				timeout := time.Second * 60
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				_, err := offChainClient.ProposeJob(ctx, jobReq)
				if err != nil {
					return errors.Wrapf(err, "failed to propose job %s for node %s", jobDesc.Flag, jobReq.NodeId)
				}
				if ctx.Err() != nil {
					return errors.Wrapf(err, "timed out after %s proposing job %s for node %s", timeout.String(), jobDesc.Flag, jobReq.NodeId)
				}

				return nil
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, "failed to create at least one job for DON")
	}

	return nil
}
