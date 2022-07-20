package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// JobRenderer wraps the JSONAPI Job Resource and adds rendering functionality
type JobPresenter struct {
	JAID // This is needed to render the id for a JSONAPI Resource as normal JSON
	presenters.JobResource
}

// ToRows returns the job as a multiple rows per task
func (p JobPresenter) ToRows() [][]string {
	row := [][]string{}

	// Produce a row when there are no tasks
	if len(p.FriendlyTasks()) == 0 {
		row = append(row, p.toRow(""))

		return row
	}

	for _, t := range p.FriendlyTasks() {
		row = append(row, p.toRow(t))
	}

	return row
}

// ToRow generates a row for a task
func (p JobPresenter) toRow(task string) []string {
	return []string{
		p.GetID(),
		p.Name,
		p.Type.String(),
		task,
		p.FriendlyCreatedAt(),
	}
}

// GetTasks extracts the tasks from the dependency graph
func (p JobPresenter) GetTasks() ([]string, error) {
	types := []string{}
	pipeline, err := pipeline.Parse(p.PipelineSpec.DotDAGSource)
	if err != nil {
		return nil, err
	}

	for _, t := range pipeline.Tasks {
		types = append(types, fmt.Sprintf("%s %s", t.DotID(), t.Type()))
	}

	return types, nil
}

// FriendlyTasks returns the tasks
func (p JobPresenter) FriendlyTasks() []string {
	taskTypes, err := p.GetTasks()
	if err != nil {
		return []string{"error parsing DAG"}
	}

	return taskTypes
}

// FriendlyCreatedAt returns the created at timestamp of the spec which matches the
// type in RFC3339 format.
func (p JobPresenter) FriendlyCreatedAt() string {
	switch p.Type {
	case presenters.DirectRequestJobSpec:
		if p.DirectRequestSpec != nil {
			return p.DirectRequestSpec.CreatedAt.Format(time.RFC3339)
		}
	case presenters.FluxMonitorJobSpec:
		if p.FluxMonitorSpec != nil {
			return p.FluxMonitorSpec.CreatedAt.Format(time.RFC3339)
		}
	case presenters.OffChainReportingJobSpec:
		if p.OffChainReportingSpec != nil {
			return p.OffChainReportingSpec.CreatedAt.Format(time.RFC3339)
		}
	case presenters.KeeperJobSpec:
		if p.KeeperSpec != nil {
			return p.KeeperSpec.CreatedAt.Format(time.RFC3339)
		}
	case presenters.CronJobSpec:
		if p.CronSpec != nil {
			return p.CronSpec.CreatedAt.Format(time.RFC3339)
		}
	case presenters.VRFJobSpec:
		if p.VRFSpec != nil {
			return p.VRFSpec.CreatedAt.Format(time.RFC3339)
		}
	case presenters.WebhookJobSpec:
		if p.WebhookSpec != nil {
			return p.WebhookSpec.CreatedAt.Format(time.RFC3339)
		}
	case presenters.BlockhashStoreJobSpec:
		if p.BlockhashStoreSpec != nil {
			return p.BlockhashStoreSpec.CreatedAt.Format(time.RFC3339)
		}
	case presenters.BootstrapJobSpec:
		if p.BootstrapSpec != nil {
			return p.BootstrapSpec.CreatedAt.Format(time.RFC3339)
		}
	default:
		return "unknown"
	}

	// This should never occur since the job should always have a spec matching
	// the type
	return "N/A"
}

// RenderTable implements TableRenderer
func (p *JobPresenter) RenderTable(rt RendererTable) error {
	table := rt.newTable([]string{"ID", "Name", "Type", "Tasks", "Created At"})
	table.SetAutoMergeCells(true)
	for _, r := range p.ToRows() {
		table.Append(r)
	}

	render("Jobs", table)
	return nil
}

type JobPresenters []JobPresenter

// RenderTable implements TableRenderer
func (ps JobPresenters) RenderTable(rt RendererTable) error {
	table := rt.newTable([]string{"ID", "Name", "Type", "Tasks", "Created At"})
	table.SetAutoMergeCells(true)
	for _, p := range ps {
		for _, r := range p.ToRows() {
			table.Append(r)
		}
	}

	render("Jobs (V2)", table)
	return nil
}

// ListJobs lists all jobs
func (cli *Client) ListJobs(c *cli.Context) (err error) {
	return cli.getPage("/v2/jobs", c.Int("page"), &JobPresenters{})
}

// ShowJob displays the details of a job
func (cli *Client) ShowJob(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must provide the id of the job"))
	}
	id := c.Args().First()
	resp, err := cli.HTTP.Get("/v2/jobs/" + id)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &JobPresenter{})
}

// CreateJob creates a job
// Valid input is a TOML string or a path to TOML file
func (cli *Client) CreateJob(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("must pass in TOML or filepath"))
	}

	tomlString, err := getTOMLString(c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}

	request, err := json.Marshal(web.CreateJobRequest{
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
		fmt.Printf("Response: '%v', Status: %d\n", string(body), resp.StatusCode)
		return cli.errorOut(err)
	}

	err = cli.renderAPIResponse(resp, &JobPresenter{}, "Job created")
	return err
}

// DeleteJob deletes a job
func (cli *Client) DeleteJob(c *cli.Context) error {
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

// TriggerPipelineRun triggers a job run based on a job ID
func (cli *Client) TriggerPipelineRun(c *cli.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the job id to trigger a run"))
	}
	resp, err := cli.HTTP.Post("/v2/jobs/"+c.Args().First()+"/runs", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var run presenters.PipelineRunResource
	err = cli.renderAPIResponse(resp, &run, "Pipeline run successfully triggered")
	return err
}
