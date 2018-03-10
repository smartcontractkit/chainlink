package cmd

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
)

// Renderer implements the Render method.
type Renderer interface {
	Render(interface{}) error
}

// RendererJSON is used to render JSON data.
type RendererJSON struct {
	io.Writer
}

// Render writes the given input as a JSON string.
func (rj RendererJSON) Render(v interface{}) error {
	b, err := utils.FormatJSON(v)
	if err != nil {
		return err
	}
	if _, err = rj.Write(b); err != nil {
		return err
	}
	return nil
}

// RendererTable is used for data to be rendered as a table.
type RendererTable struct {
	io.Writer
}

// Render returns a formatted table of text for a given Job or presenter
// and relevant information.
func (rt RendererTable) Render(v interface{}) error {
	switch typed := v.(type) {
	case *[]models.JobSpec:
		rt.renderJobs(*typed)
	case *presenters.JobSpec:
		rt.renderJob(*typed)
	default:
		return fmt.Errorf("Unable to render object: %v", typed)
	}

	return nil
}

func (rt RendererTable) renderJobs(jobs []models.JobSpec) error {
	table := tablewriter.NewWriter(rt)
	table.SetHeader([]string{"ID", "Created At", "Initiators", "Tasks"})
	for _, v := range jobs {
		table.Append(jobRowToStrings(v))
	}

	render("Jobs", table)
	return nil
}

func render(name string, table *tablewriter.Table) {
	table.SetRowLine(true)
	table.SetColumnSeparator("║")
	table.SetRowSeparator("═")
	table.SetCenterSeparator("╬")

	fmt.Println("╔ " + name)
	table.Render()
}

func jobRowToStrings(job models.JobSpec) []string {
	p := presenters.JobSpec{job, nil}
	return []string{
		p.ID,
		p.FriendlyCreatedAt(),
		p.FriendlyInitiators(),
		p.FriendlyTasks(),
	}
}

func (rt RendererTable) renderJob(job presenters.JobSpec) error {
	if err := rt.renderJobSingles(job); err != nil {
		return err
	}

	if err := rt.renderJobInitiators(job); err != nil {
		return err
	}

	if err := rt.renderJobTasks(job); err != nil {
		return err
	}

	if err := rt.renderJobRuns(job); err != nil {
		return err
	}

	return nil
}

func (rt RendererTable) renderJobSingles(j presenters.JobSpec) error {
	table := tablewriter.NewWriter(rt)
	table.SetHeader([]string{"ID", "Created At", "Start At", "End At"})
	table.Append([]string{
		j.ID,
		j.FriendlyCreatedAt(),
		j.FriendlyStartAt(),
		j.FriendlyEndAt(),
	})
	render("Job", table)
	return nil
}

func (rt RendererTable) renderJobInitiators(j presenters.JobSpec) error {
	table := tablewriter.NewWriter(rt)
	table.SetHeader([]string{"Type", "Schedule", "Run At", "Address"})
	for _, i := range j.Initiators {
		p := presenters.Initiator{i}
		table.Append([]string{
			p.Type,
			p.Schedule.String(),
			p.FriendlyRunAt(),
			p.FriendlyAddress(),
		})
	}

	render("Initiators", table)
	return nil
}

func (rt RendererTable) renderJobTasks(j presenters.JobSpec) error {
	table := tablewriter.NewWriter(rt)
	table.SetHeader([]string{"Type", "Config", "Value"})
	for _, t := range j.Tasks {
		p := presenters.TaskSpec{t}
		keys, values := p.FriendlyParams()
		table.Append([]string{p.Type, keys, values})
	}

	render("Tasks", table)
	return nil
}

func (rt RendererTable) renderJobRuns(j presenters.JobSpec) error {
	table := tablewriter.NewWriter(rt)
	table.SetHeader([]string{"ID", "Status", "Created", "Completed", "Result", "Error"})
	for _, jr := range j.Runs {
		table.Append([]string{
			jr.ID,
			jr.Status,
			utils.ISO8601UTC(jr.CreatedAt),
			utils.NullISO8601UTC(jr.CompletedAt),
			jr.Result.Data.String(),
			jr.Result.ErrorMessage.String,
		})
	}

	render("Runs", table)
	return nil
}
