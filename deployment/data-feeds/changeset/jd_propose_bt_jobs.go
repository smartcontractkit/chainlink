package changeset

import (
	"context"
	"errors"
	"fmt"
	"time"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/offchain"
)

const (
	btTimeout = 120 * time.Second
)

// ProposeBtJobsToJDChangeset is a changeset that reads a boostrap spec from a file and proposes jobs to JD
var ProposeBtJobsToJDChangeset = cldf.CreateChangeSet(proposeBtJobsToJDLogic, proposeBtJobsToJDPrecondition)

func proposeBtJobsToJDLogic(env deployment.Environment, c types.ProposeBtJobsConfig) (deployment.ChangesetOutput, error) {
	ctx, cancel := context.WithTimeout(env.GetContext(), btTimeout)
	defer cancel()

	bootstrapJobSpec, err := offchain.JobSpecFromBootstrap(c.NodeFilter.DONID, c.ChainSelector, c.BootstrapJobName, c.Contract)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create job spec from bootstrap: %w", err)
	}

	return offchain.ProposeJobs(ctx, env, bootstrapJobSpec, nil, c.NodeFilter)
}

func proposeBtJobsToJDPrecondition(env deployment.Environment, c types.ProposeBtJobsConfig) error {
	if c.NodeFilter == nil {
		return errors.New("node labels are required")
	}

	_, ok := env.Chains[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	return nil
}
