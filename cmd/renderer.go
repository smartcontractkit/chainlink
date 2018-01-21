package cmd

import (
	"fmt"
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
)

type Renderer interface {
	Render(interface{}) error
}

type RendererJSON struct {
	io.Writer
}

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

type RendererNoOp struct{}

func (rj RendererNoOp) Render(v interface{}) error { return nil }

type RendererTable struct {
	io.Writer
}

func (rt RendererTable) Render(v interface{}) error {
	switch typed := v.(type) {
	case *[]models.Job:
		rt.renderJobs(*typed)
	case *presenters.Job:
		rt.renderJob(*typed)
	default:
		return fmt.Errorf("Unable to render object: %v", typed)
	}

	return nil
}

func (rt RendererTable) renderJobs(jobs []models.Job) error {
	table := tablewriter.NewWriter(rt)
	table.SetHeader([]string{"ID", "Created", "Initiators", "Tasks", "End"})
	for _, v := range jobs {
		table.Append(jobRowToStrings(v))
	}

	table.Render()
	return nil
}

func jobRowToStrings(job models.Job) []string {
	p := presenters.Job{job, nil}
	return []string{
		p.ID,
		p.FriendlyCreatedAt(),
		p.FriendlyInitiators(),
		p.FriendlyTasks(),
		p.FriendlyEndAt(),
	}
}

func (rt RendererTable) renderJob(job presenters.Job) error {
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

func (rt RendererTable) renderJobSingles(j presenters.Job) error {
	table := tablewriter.NewWriter(rt)
	table.SetHeader([]string{"ID", "Created", "End"})
	table.Append([]string{j.ID, j.FriendlyCreatedAt(), j.FriendlyEndAt()})
	table.Render()
	return nil
}

func (rt RendererTable) renderJobInitiators(j presenters.Job) error {
	table := tablewriter.NewWriter(rt)
	table.SetHeader([]string{"Initiator", "Schedule", "Run At", "Address"})
	for _, i := range j.Initiators {
		p := presenters.Initiator{i}
		table.Append([]string{
			p.Type,
			p.Schedule.String(),
			p.FriendlyRunAt(),
			p.FriendlyAddress(),
		})
	}

	table.Render()
	return nil
}

func (rt RendererTable) renderJobTasks(j presenters.Job) error {
	table := tablewriter.NewWriter(rt)
	table.SetHeader([]string{"Order", "Task", "Params"})
	for o, t := range j.Tasks {
		p := presenters.Task{t}
		params, err := p.FriendlyParams()
		if err != nil {
			return err
		}

		table.Append([]string{strconv.Itoa(o), p.Type, params})
	}

	table.Render()
	return nil
}

func (rt RendererTable) renderJobRuns(j presenters.Job) error {
	table := tablewriter.NewWriter(rt)
	table.SetHeader([]string{"Job Run", "Status", "Created", "Result", "Error"})
	for _, jr := range j.Runs {
		output, err := jr.Result.Output.String()
		if err != nil {
			return err
		}
		table.Append([]string{
			jr.ID,
			jr.Status,
			utils.ISO8601UTC(jr.CreatedAt),
			output,
			jr.Result.ErrorMessage.String,
		})
	}

	table.Render()
	return nil
}
