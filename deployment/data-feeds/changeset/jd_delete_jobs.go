package changeset

import (
	"context"
	"errors"
	"time"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/offchain"
)

const (
	deleteJobTimeout = 120 * time.Second
)

// DeleteJobsJDChangeset is a changeset that deletes jobs from JD either using job ids or workflow name
var DeleteJobsJDChangeset = cldf.CreateChangeSet(deleteJobsJDLogic, deleteJobsJDPrecondition)

func deleteJobsJDLogic(env cldf.Environment, c types.DeleteJobsConfig) (cldf.ChangesetOutput, error) {
	ctx, cancel := context.WithTimeout(env.GetContext(), deleteJobTimeout)
	defer cancel()

	offchain.DeleteJobs(ctx, env, c.JobIDs, c.WorkflowName)
	return cldf.ChangesetOutput{}, nil
}

func deleteJobsJDPrecondition(_ cldf.Environment, c types.DeleteJobsConfig) error {
	if len(c.JobIDs) == 0 && c.WorkflowName == "" {
		return errors.New("job ids or workflow name are required")
	}
	return nil
}
