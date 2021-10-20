package cmd

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	webpresenters "github.com/smartcontractkit/chainlink/core/web/presenters"
)

// Renderer implements the Render method.
type Renderer interface {
	Render(interface{}, ...string) error
}

// RendererJSON is used to render JSON data.
type RendererJSON struct {
	io.Writer
}

// Render writes the given input as a JSON string.
func (rj RendererJSON) Render(v interface{}, _ ...string) error {
	b, err := utils.FormatJSON(v)
	if err != nil {
		return err
	}

	// Append a new line
	b = append(b, []byte("\n")...)
	if _, err = rj.Write(b); err != nil {
		return err
	}
	return nil
}

// RendererTable is used for data to be rendered as a table.
type RendererTable struct {
	io.Writer
}

type TableRenderer interface {
	RenderTable(rt RendererTable) error
}

// Render returns a formatted table of text for a given Job or presenter
// and relevant information.
func (rt RendererTable) Render(v interface{}, headers ...string) error {
	for _, h := range headers {
		fmt.Println(h)
	}

	switch typed := v.(type) {
	case *[]models.JobSpec:
		return rt.renderJobs(*typed)
	case *presenters.JobSpec:
		return rt.renderJob(*typed)
	case *[]presenters.JobRun:
		return rt.renderJobRuns(*typed)
	case *presenters.JobRun:
		return rt.renderJobRun(*typed)
	case *presenters.ServiceAgreement:
		return rt.renderServiceAgreement(*typed)
	case *webpresenters.ExternalInitiatorAuthentication:
		return rt.renderExternalInitiatorAuthentication(*typed)
	case *web.ConfigPatchResponse:
		return rt.renderConfigPatchResponse(typed)
	case *presenters.ConfigPrinter:
		return rt.renderConfiguration(*typed)
	case *webpresenters.PipelineRunResource:
		return rt.renderPipelineRun(*typed)
	case *webpresenters.ServiceLogConfigResource:
		return rt.renderLogPkgConfig(*typed)
	case *[]VRFKeyPresenter:
		return rt.renderVRFKeys(*typed)
	case TableRenderer:
		return typed.RenderTable(rt)
	default:
		return fmt.Errorf("unable to render object of type %T: %v", typed, typed)
	}
}

func (rt RendererTable) renderLogPkgConfig(serviceLevelLog webpresenters.ServiceLogConfigResource) error {
	table := rt.newTable([]string{"ID", "Service", "LogLevel"})
	for i, svcName := range serviceLevelLog.ServiceName {
		table.Append([]string{
			serviceLevelLog.ID,
			svcName,
			serviceLevelLog.LogLevel[i],
		})
	}

	render("ServiceLogConfig", table)
	return nil
}

func (rt RendererTable) renderVRFKeys(keys []VRFKeyPresenter) error {
	var rows [][]string

	for _, key := range keys {
		rows = append(rows, []string{
			key.Compressed,
			key.Uncompressed,
			key.Hash,
			key.CreatedAt.String(),
			key.UpdatedAt.String(),
			key.FriendlyDeletedAt(),
		})
	}

	renderList([]string{"Compressed", "Uncompressed", "Hash", "Created", "Updated", "Deleted"}, rows, rt.Writer)

	return nil
}

func (rt RendererTable) renderJobs(jobs []models.JobSpec) error {
	table := rt.newTable([]string{"ID", "Name", "Created At", "Initiators", "Tasks"})
	for _, v := range jobs {
		table.Append(jobRowToStrings(v))
	}

	render("Jobs", table)
	return nil
}

func (rt RendererTable) renderConfiguration(cp presenters.ConfigPrinter) error {
	table := rt.newTable([]string{"Key", "Value"})
	schemaT := reflect.TypeOf(config.ConfigSchema{})
	cpT := reflect.TypeOf(cp.EnvPrinter)
	cpV := reflect.ValueOf(cp.EnvPrinter)

	for index := 0; index < cpT.NumField(); index++ {
		item := cpT.FieldByIndex([]int{index})
		schemaItem, ok := schemaT.FieldByName(item.Name)
		if !ok {
			logger.Panicf("Field %s missing from store.Schema", item.Name)
		}
		envName, ok := schemaItem.Tag.Lookup("env")
		if !ok {
			continue
		}
		field := cpV.FieldByIndex(item.Index)

		if stringer, ok := field.Interface().(fmt.Stringer); ok {
			if stringer != reflect.Zero(reflect.TypeOf(stringer)).Interface() {
				table.Append([]string{
					envName,
					stringer.String(),
				})
			}
		} else {
			table.Append([]string{
				envName,
				fmt.Sprintf("%v", field),
			})
		}
	}

	render("Configuration", table)
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

func renderList(fields []string, items [][]string, writer io.Writer) {
	var maxLabelLength int
	for _, field := range fields {
		if len(field) > maxLabelLength {
			maxLabelLength = len(field)
		}
	}
	var itemsRendered []string
	var maxLineLength int
	for _, row := range items {
		var lines []string
		for i, field := range fields {
			diff := maxLabelLength - len(field)
			spaces := strings.Repeat(" ", diff)
			line := fmt.Sprintf("%v: %v%v", field, spaces, row[i])
			if len(line) > maxLineLength {
				maxLineLength = len(line)
			}
			lines = append(lines, line)
		}
		itemsRendered = append(itemsRendered, strings.Join(lines, "\n"))
	}
	divider := strings.Repeat("-", maxLineLength)
	listRendered := divider + "\n" + strings.Join(itemsRendered, "\n"+divider+"\n")
	_, err := writer.Write([]byte(listRendered))
	if err != nil {
		// Handles errcheck
		return
	}
}

func jobRowToStrings(job models.JobSpec) []string {
	p := presenters.JobSpec{JobSpec: job}
	return []string{
		p.ID.String(),
		p.Name,
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

	err := rt.renderJobTasks(job)
	return err
}

func (rt RendererTable) renderJobRun(run presenters.JobRun) error {
	err := rt.renderJobRuns([]presenters.JobRun{run})
	return err
}

func (rt RendererTable) renderJobSingles(j presenters.JobSpec) error {
	table := rt.newTable([]string{"ID", "Name", "Created At", "Start At", "End At", "Min Payment"})
	table.Append([]string{
		j.ID.String(),
		j.Name,
		j.FriendlyCreatedAt(),
		j.FriendlyStartAt(),
		j.FriendlyEndAt(),
		j.FriendlyMinPayment(),
	})
	render("Job", table)
	return nil
}

func (rt RendererTable) renderJobInitiators(j presenters.JobSpec) error {
	table := rt.newTable([]string{"Type", "Schedule", "Run At", "Address"})
	for _, i := range j.Initiators {
		p := presenters.Initiator{Initiator: i}
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
	table := rt.newTable([]string{"Type", "Config", "Value"})
	table.SetAutoWrapText(false)
	for _, t := range j.Tasks {
		p := presenters.TaskSpec{TaskSpec: t}
		keys, values := p.FriendlyParams()
		table.Append([]string{p.Type.String(), keys, values})
	}

	render("Tasks", table)
	return nil
}

func (rt RendererTable) renderJobRuns(runs []presenters.JobRun) error {
	table := rt.newTable([]string{"ID", "Status", "Created", "Completed", "Result", "Error"})
	for _, jr := range runs {
		table.Append([]string{
			jr.ID.String(),
			string(jr.GetStatus()),
			utils.ISO8601UTC(jr.CreatedAt),
			utils.NullISO8601UTC(jr.FinishedAt),
			jr.Result.Data.String(),
			jr.ErrorString(),
		})
	}

	render("Runs", table)
	return nil
}

func (rt RendererTable) renderServiceAgreement(sa presenters.ServiceAgreement) error {
	table := rt.newTable([]string{"ID", "Created At", "Payment", "Expiration", "Aggregator", "AggInit", "AggFulfill"})
	table.Append([]string{
		sa.ID,
		sa.FriendlyCreatedAt(),
		sa.FriendlyPayment(),
		sa.FriendlyExpiration(),
		sa.FriendlyAggregator(),
		sa.FriendlyAggregatorInitMethod(),
		sa.FriendlyAggregatorFulfillMethod(),
	})
	render("Service Agreement", table)
	return nil
}

func (rt RendererTable) renderExternalInitiatorAuthentication(eia webpresenters.ExternalInitiatorAuthentication) error {
	table := rt.newTable([]string{"Name", "URL", "AccessKey", "Secret", "OutgoingToken", "OutgoingSecret"})
	table.Append([]string{
		eia.Name,
		eia.URL.String(),
		eia.AccessKey,
		eia.Secret,
		eia.OutgoingToken,
		eia.OutgoingSecret,
	})
	render("External Initiator Credentials:", table)
	return nil
}

func (rt RendererTable) newTable(headers []string) *tablewriter.Table {
	table := tablewriter.NewWriter(rt)
	table.SetHeader(headers)
	return table
}

func (rt RendererTable) renderConfigPatchResponse(config *web.ConfigPatchResponse) error {
	table := rt.newTable([]string{"Config", "Old Value", "New Value"})
	table.Append([]string{
		"EthGasPriceDefault",
		config.EthGasPriceDefault.From,
		config.EthGasPriceDefault.To,
	})
	render("Configuration Changes", table)
	return nil
}

func (rt RendererTable) renderPipelineRun(run webpresenters.PipelineRunResource) error {
	table := rt.newTable([]string{"ID", "Created At", "Finished At"})

	var finishedAt string
	if !run.FinishedAt.IsZero() {
		finishedAt = run.FinishedAt.String()
	}

	row := []string{
		run.GetID(),
		run.CreatedAt.String(),
		finishedAt,
	}
	table.Append(row)

	render("Pipeline Run", table)
	return nil
}
