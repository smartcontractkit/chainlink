package jobs

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	keystoneflags "github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	types "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

var SupportedJobs = []types.JobDescription{
	{Flag: types.OCR3Capability, NodeType: types.BootstrapNode},
	{Flag: types.WorkflowDON, NodeType: types.BootstrapNode},
	{Flag: types.CustomComputeCapability, NodeType: types.BootstrapNode},
	{Flag: types.CronCapability, NodeType: types.WorkerNode},
	{Flag: types.CustomComputeCapability, NodeType: types.WorkerNode},
	{Flag: types.OCR3Capability, NodeType: types.WorkerNode},
	{Flag: types.GatewayDON, NodeType: types.GatewayDON},

	// add more jobs as needed
}

func checkForUnknownJobs(jobSpecs types.DonJobs) error {
	for jobDesc := range jobSpecs {
		found := false
		for _, supportedJob := range SupportedJobs {
			if jobDesc.Flag == supportedJob.Flag {
				found = true
				break
			}
		}

		if !found {
			return errors.Errorf("unknown job type %s", jobDesc.Flag)
		}
	}

	return nil
}

func Create(offChainClient deployment.OffchainClient, don *devenv.DON, flags []string, jobSpecs types.DonJobs) error {
	if len(jobSpecs) == 0 {
		return nil
	}

	if unknownErr := checkForUnknownJobs(jobSpecs); unknownErr != nil {
		return errors.Wrap(unknownErr, "failed to create jobs")
	}

	errCh := make(chan error, calculateJobCount(jobSpecs))
	var wg sync.WaitGroup

	for _, jobDesc := range SupportedJobs {
		if keystoneflags.HasFlag(flags, jobDesc.Flag) {
			if jobReqs, ok := jobSpecs[jobDesc]; ok {
				for _, jobReq := range jobReqs {
					wg.Add(1)
					go func(jobReq *jobv1.ProposeJobRequest) {
						defer wg.Done()
						timeout := time.Second * 60
						ctx, cancel := context.WithTimeout(context.Background(), timeout)
						defer cancel()
						_, err := offChainClient.ProposeJob(ctx, jobReq)
						if err != nil {
							errCh <- errors.Wrapf(err, "failed to propose job %s for node %s", jobDesc.Flag, jobReq.NodeId)
						}
						err = ctx.Err()
						if err != nil {
							errCh <- errors.Wrapf(err, "timed out after %s proposing job %s for node %s", timeout.String(), jobDesc.Flag, jobReq.NodeId)
						}
					}(jobReq)
				}
			}
		}
	}

	wg.Wait()
	close(errCh)

	var finalErr error
	for err := range errCh {
		if finalErr == nil {
			finalErr = err
		} else {
			finalErr = errors.Wrap(finalErr, err.Error())
		}
	}

	if finalErr != nil {
		return errors.Wrap(finalErr, "failed to create at least one job for DON")
	}

	return nil
}

func calculateJobCount(jobSpecs types.DonJobs) int {
	count := 0
	for _, jobSpec := range jobSpecs {
		count += len(jobSpec)
	}

	return count
}
