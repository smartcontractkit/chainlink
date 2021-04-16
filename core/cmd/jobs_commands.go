package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/urfave/cli"
	"go.uber.org/multierr"
)

// JobRenderer wraps the JSONAPI Job Resource and adds rendering functionality
type JobPresenter struct {
	JAID // This is needed to render the id for a JSONAPI Resource as normal JSON
	presenters.JobResource
}

// ToRows returns the job as a multiple rows per task
func (j JobPresenter) ToRows() [][]string {
	row := [][]string{}

	// Produce a row when there are no tasks
	if len(j.FriendlyTasks()) == 0 {
		row = append(row, j.toRow(""))

		return row
	}

	for _, t := range j.FriendlyTasks() {
		row = append(row, j.toRow(t))
	}

	return row
}

// ToRow generates a row for a task
func (j JobPresenter) toRow(task string) []string {
	return []string{
		j.GetID(),
		j.Name,
		j.Type.String(),
		task,
		j.FriendlyCreatedAt(),
	}
}

// GetTasks extracts the tasks from the dependency graph
func (j JobPresenter) GetTasks() ([]string, error) {
	types := []string{}
	dag := pipeline.NewTaskDAG()
	err := dag.UnmarshalText([]byte(j.PipelineSpec.DotDAGSource))
	if err != nil {
		return nil, err
	}

	tasks, err := dag.TasksInDependencyOrder()
	if err != nil {
		return nil, err
	}

	// Reverse the order as dependency tasks start from output and end at the
	// inputs.
	for i := len(tasks) - 1; i >= 0; i-- {
		t := tasks[i]
		types = append(types, fmt.Sprintf("%s %s", t.DotID(), t.Type()))
	}

	return types, nil
}

// FriendlyTasks returns the tasks
func (j JobPresenter) FriendlyTasks() []string {
	taskTypes, err := j.GetTasks()
	if err != nil {
		return []string{"error parsing DAG"}
	}

	return taskTypes
}

// FriendlyCreatedAt returns the created at timestamp of the spec which matches the
// type in RFC3339 format.
func (j JobPresenter) FriendlyCreatedAt() string {
	switch j.Type {
	case presenters.DirectRequestJobSpec:
		if j.DirectRequestSpec != nil {
			return j.DirectRequestSpec.CreatedAt.Format(time.RFC3339)
		}
	case presenters.FluxMonitorJobSpec:
		if j.FluxMonitorSpec != nil {
			return j.FluxMonitorSpec.CreatedAt.Format(time.RFC3339)
		}
	case presenters.OffChainReportingJobSpec:
		if j.OffChainReportingSpec != nil {
			return j.OffChainReportingSpec.CreatedAt.Format(time.RFC3339)
		}
	case presenters.KeeperJobSpec:
		if j.KeeperSpec != nil {
			return j.KeeperSpec.CreatedAt.Format(time.RFC3339)
		}
	default:
		return "unknown"
	}

	// This should never occur since the job should always have a spec matching
	// the type
	return "N/A"
}

// ListJobsV2 lists all v2 jobs
func (cli *Client) ListJobsV2(c *cli.Context) (err error) {
	return cli.getPage("/v2/jobs", c.Int("page"), &[]JobPresenter{})
}

// CreateJobV2 creates a V2 job
// Valid input is a TOML string or a path to TOML file
func (cli *Client) CreateJobV2(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass in TOML or filepath"))
	}

	tomlString, err := getTOMLString(c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}

	request, err := json.Marshal(models.CreateJobSpecRequest{
		TOML: tomlString,
	})
	if err != nil {
		return cli.errorOut(err)
	}

	resp, err := cli.HTTP.Post("/v2/jobs", bytes.NewReader(request))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode >= 400 {
		body, rerr := ioutil.ReadAll(resp.Body)
		if err != nil {
			err = multierr.Append(err, rerr)
			return cli.errorOut(err)
		}
		fmt.Printf("Error : %v\n", string(body))
		return cli.errorOut(err)
	}

	var presenter JobPresenter
	err = cli.renderAPIResponse(resp, &presenter, "Job created")
	return err
}

// DeleteJobV2 deletes a V2 job
func (cli *Client) DeleteJobV2(c *cli.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass the job id to be archived"))
	}
	resp, err := cli.HTTP.Delete("/v2/jobs/" + c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}
	_, err = cli.parseResponse(resp)
	if err != nil {
		return cli.errorOut(err)
	}

	fmt.Printf("Job %v Deleted\n", c.Args().First())
	return nil
}
