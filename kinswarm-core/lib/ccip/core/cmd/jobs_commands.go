package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func initJobsSubCmds(s *Shell) []cli.Command {
	return []cli.Command{
		{
			Name:   "list",
			Usage:  "List all jobs",
			Action: s.ListJobs,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "page",
					Usage: "page of results to display",
				},
			},
		},
		{
			Name:   "show",
			Usage:  "Show a job",
			Action: s.ShowJob,
		},
		{
			Name:   "create",
			Usage:  "Create a job",
			Action: s.CreateJob,
		},
		{
			Name:   "delete",
			Usage:  "Delete a job",
			Action: s.DeleteJob,
		},
		{
			Name:   "run",
			Usage:  "Trigger a job run",
			Action: s.TriggerPipelineRun,
		},
	}
}

// JobPresenter wraps the JSONAPI Job Resource and adds rendering functionality
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
	if strings.TrimSpace(p.PipelineSpec.DotDAGSource) == "" {
		return nil, nil
	}
	var types []string
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
	case presenters.BlockHeaderFeederJobSpec:
		if p.BlockHeaderFeederSpec != nil {
			return p.BlockHeaderFeederSpec.CreatedAt.Format(time.RFC3339)
		}
	case presenters.BootstrapJobSpec:
		if p.BootstrapSpec != nil {
			return p.BootstrapSpec.CreatedAt.Format(time.RFC3339)
		}
	case presenters.GatewayJobSpec:
		if p.GatewaySpec != nil {
			return p.GatewaySpec.CreatedAt.Format(time.RFC3339)
		}
	case presenters.WorkflowJobSpec:
		if p.WorkflowSpec != nil {
			return p.WorkflowSpec.CreatedAt.Format(time.RFC3339)
		}
	case presenters.StandardCapabilitiesJobSpec:
		if p.StandardCapabilitiesSpec != nil {
			return p.StandardCapabilitiesSpec.CreatedAt.Format(time.RFC3339)
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
func (s *Shell) ListJobs(c *cli.Context) (err error) {
	return s.getPage("/v2/jobs", c.Int("page"), &JobPresenters{})
}

// ShowJob displays the details of a job
func (s *Shell) ShowJob(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return s.errorOut(errors.New("must provide the id of the job"))
	}
	id := c.Args().First()
	resp, err := s.HTTP.Get(s.ctx(), "/v2/jobs/"+id)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return s.renderAPIResponse(resp, &JobPresenter{})
}

// CreateJob creates a job
// Valid input is a TOML string or a path to TOML file
func (s *Shell) CreateJob(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return s.errorOut(errors.New("must pass in TOML or filepath"))
	}

	tomlString, err := getTOMLString(c.Args().First())
	if err != nil {
		return s.errorOut(err)
	}

	request, err := json.Marshal(web.CreateJobRequest{
		TOML: tomlString,
	})
	if err != nil {
		return s.errorOut(err)
	}

	resp, err := s.HTTP.Post(s.ctx(), "/v2/jobs", bytes.NewReader(request))
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode >= 400 {
		body, rerr := io.ReadAll(resp.Body)
		if err != nil {
			err = multierr.Append(err, rerr)
			return s.errorOut(err)
		}
		fmt.Printf("Response: '%v', Status: %d\n", string(body), resp.StatusCode)
		return s.errorOut(err)
	}

	err = s.renderAPIResponse(resp, &JobPresenter{}, "Job created")
	return err
}

// DeleteJob deletes a job
func (s *Shell) DeleteJob(c *cli.Context) error {
	if !c.Args().Present() {
		return s.errorOut(errors.New("must pass the job id to be archived"))
	}
	resp, err := s.HTTP.Delete(s.ctx(), "/v2/jobs/"+c.Args().First())
	if err != nil {
		return s.errorOut(err)
	}
	_, err = s.parseResponse(resp)
	if err != nil {
		return s.errorOut(err)
	}

	fmt.Printf("Job %v Deleted\n", c.Args().First())
	return nil
}

// TriggerPipelineRun triggers a job run based on a job ID
func (s *Shell) TriggerPipelineRun(c *cli.Context) error {
	if !c.Args().Present() {
		return s.errorOut(errors.New("Must pass the job id to trigger a run"))
	}
	resp, err := s.HTTP.Post(s.ctx(), "/v2/jobs/"+c.Args().First()+"/runs", nil)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var run presenters.PipelineRunResource
	err = s.renderAPIResponse(resp, &run, "Pipeline run successfully triggered")
	return err
}
