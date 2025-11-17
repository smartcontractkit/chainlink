package changeset

import (
	"context"
	"errors"
	"time"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/offchain"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

const (
	deleteJobTimeout = 5 * time.Minute
)

// DeleteJobsJDChangeset is a changeset that deletes jobs from JD either using job ids or workflow name
var DeleteJobsJDChangeset = cldf.CreateChangeSet(deleteJobsJDLogic, deleteJobsJDPrecondition)

func deleteJobsJDLogic(env cldf.Environment, c types.DeleteJobsConfig) (cldf.ChangesetOutput, error) {
	ctx, cancel := context.WithTimeout(env.GetContext(), deleteJobTimeout)
	defer cancel()

	offchain.DeleteJobs(ctx, env, c.JobIDs, c.WorkflowName, c.Environment, c.Zone)

	ds := datastore.NewMemoryDataStore()
	// Delete the workflow spec from the datastore if workflow name is provided
	if c.WorkflowName != "" {
		err := UpdateWorkflowMetadataDS(env, ds, c.WorkflowName, "")
		if err == nil {
			return cldf.ChangesetOutput{DataStore: ds}, nil
		}
		env.Logger.Errorf("failed to update workflow spec: %s", err)
	}
	return cldf.ChangesetOutput{}, nil
}

func deleteJobsJDPrecondition(_ cldf.Environment, c types.DeleteJobsConfig) error {
	if len(c.JobIDs) == 0 && c.WorkflowName == "" {
		return errors.New("job ids or workflow name are required")
	}
	return nil
}
