package changeset

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/offchain"
)

const (
	deleteJobTimeout = 120 * time.Second
)

// DeleteJobsJDChangeset is a changeset that deletes jobs from JD
var DeleteJobsJDChangeset = deployment.CreateChangeSet(deleteJobsJDLogic, deleteJobsJDPrecondition)

func deleteJobsJDLogic(env deployment.Environment, c types.DeleteJobsConfig) (deployment.ChangesetOutput, error) {
	ctx, cancel := context.WithTimeout(env.GetContext(), deleteJobTimeout)
	defer cancel()

	offchain.DeleteJobs(ctx, env, c.JobIDs)
	return deployment.ChangesetOutput{}, nil
}

func deleteJobsJDPrecondition(_ deployment.Environment, c types.DeleteJobsConfig) error {
	return nil
}
