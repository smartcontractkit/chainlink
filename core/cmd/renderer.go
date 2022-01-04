package cmd

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/config/envvar"
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
	case *webpresenters.ExternalInitiatorAuthentication:
		return rt.renderExternalInitiatorAuthentication(*typed)
	case *web.ConfigPatchResponse:
		return rt.renderConfigPatchResponse(typed)
	case *config.ConfigPrinter:
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
		})
	}

	renderList([]string{"Compressed", "Uncompressed", "Hash"}, rows, rt.Writer)

	return nil
}

func (rt RendererTable) renderConfiguration(cp config.ConfigPrinter) error {
	table := rt.newTable([]string{"Key", "Value"})
	schemaT := reflect.TypeOf(envvar.ConfigSchema{})
	cpT := reflect.TypeOf(cp.EnvPrinter)
	cpV := reflect.ValueOf(cp.EnvPrinter)

	for index := 0; index < cpT.NumField(); index++ {
		item := cpT.FieldByIndex([]int{index})
		schemaItem, ok := schemaT.FieldByName(item.Name)
		if !ok {
			panic(fmt.Sprintf("Field %s missing from store.Schema", item.Name))
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
			for _, l := range strings.Split(line, "\n") {
				if len(l) > maxLineLength {
					maxLineLength = len(l)
				}
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
		"EvmGasPriceDefault",
		config.EvmGasPriceDefault.From,
		config.EvmGasPriceDefault.To,
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
