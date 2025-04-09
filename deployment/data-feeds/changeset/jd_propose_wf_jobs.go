package changeset

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/offchain"
)

const (
	timeout = 120 * time.Second
)

// ProposeWfJobsToJDChangeset is a changeset that reads a workflow yaml spec from a file and proposes jobs to JD
var ProposeWfJobsToJDChangeset = deployment.CreateChangeSet(proposeWfJobsToJDLogic, proposeWfJobsToJDPrecondition)

func proposeWfJobsToJDLogic(env deployment.Environment, c types.ProposeWfJobsConfig) (deployment.ChangesetOutput, error) {
	ctx, cancel := context.WithTimeout(env.GetContext(), timeout)
	defer cancel()

	workflowJobSpec, workflowName, err := offchain.JobSpecFromWorkflow(c.InputFS, c.InputFileName, c.WorkflowJobName)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to create job spec from workflow: %w", err)
	}

	return offchain.ProposeJobs(ctx, env, workflowJobSpec, &workflowName, c.NodeFilter)
}

func proposeWfJobsToJDPrecondition(_ deployment.Environment, c types.ProposeWfJobsConfig) error {
	if c.NodeFilter == nil {
		return errors.New("node filters are required")
	}
	if c.InputFileName == "" {
		return errors.New("input file name is required")
	}
	wfYaml, err := c.InputFS.ReadFile(c.InputFileName)
	if err != nil {
		return fmt.Errorf("failed to load workflow spec from %s: %w", c.InputFileName, err)
	}
	wfStr := string(wfYaml)
	_, err = workflows.ParseWorkflowSpecYaml(wfStr)
	if err != nil {
		return fmt.Errorf("failed to parse workflow spec from %s: %w", c.InputFileName, err)
	}

	return nil
}
