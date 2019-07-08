package cmd

import (
	"fmt"
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
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
		return rt.renderJobs(*typed)
	case *presenters.JobSpec:
		return rt.renderJob(*typed)
	case *[]presenters.JobRun:
		return rt.renderJobRuns(*typed)
	case *presenters.JobRun:
		return rt.renderJobRun(*typed)
	case *models.BridgeType:
		return rt.renderBridge(*typed)
	case *models.BridgeTypeAuthentication:
		return rt.renderBridgeAuthentication(*typed)
	case *[]models.BridgeType:
		return rt.renderBridges(*typed)
	case *[]presenters.AccountBalance:
		return rt.renderAccountBalances(*typed)
	case *presenters.ServiceAgreement:
		return rt.renderServiceAgreement(*typed)
	case *[]models.TxAttempt:
		return rt.renderTxAttempts(*typed)
	case *models.ExternalInitiatorAuthentication:
		return rt.renderExternalInitiatorAuthentication(*typed)
	default:
		return fmt.Errorf("Unable to render object of type %T: %v", typed, typed)
	}
}

func (rt RendererTable) renderJobs(jobs []models.JobSpec) error {
	table := rt.newTable([]string{"ID", "Created At", "Initiators", "Tasks"})
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
	p := presenters.JobSpec{JobSpec: job}
	return []string{
		p.ID,
		p.FriendlyCreatedAt(),
		p.FriendlyInitiators(),
		p.FriendlyTasks(),
	}
}

func bridgeRowToStrings(bridge models.BridgeType) []string {
	return []string{
		bridge.Name.String(),
		bridge.URL.String(),
		strconv.FormatUint(uint64(bridge.Confirmations), 10),
	}
}

func (rt RendererTable) renderBridges(bridges []models.BridgeType) error {
	table := rt.newTable([]string{"Name", "URL", "Confirmations"})
	for _, v := range bridges {
		table.Append(bridgeRowToStrings(v))
	}

	render("Bridges", table)
	return nil
}

func (rt RendererTable) renderBridge(bridge models.BridgeType) error {
	table := rt.newTable([]string{"Name", "URL", "Default Confirmations", "Outgoing Token"})
	table.Append([]string{
		bridge.Name.String(),
		bridge.URL.String(),
		strconv.FormatUint(uint64(bridge.Confirmations), 10),
		bridge.OutgoingToken,
	})
	render("Bridge", table)
	return nil
}

func (rt RendererTable) renderBridgeAuthentication(bridge models.BridgeTypeAuthentication) error {
	table := rt.newTable([]string{"Name", "URL", "Default Confirmations", "Incoming Token", "Outgoing Token"})
	table.Append([]string{
		bridge.Name.String(),
		bridge.URL.String(),
		strconv.FormatUint(uint64(bridge.Confirmations), 10),
		bridge.IncomingToken,
		bridge.OutgoingToken,
	})
	render("Bridge", table)
	return nil
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
	table := rt.newTable([]string{"ID", "Created At", "Start At", "End At"})
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
			jr.ID,
			string(jr.Status),
			utils.ISO8601UTC(jr.CreatedAt),
			utils.NullISO8601UTC(jr.FinishedAt),
			jr.Result.Data.String(),
			jr.Result.ErrorMessage.String,
		})
	}

	render("Runs", table)
	return nil
}

func (rt RendererTable) renderAccountBalances(balances []presenters.AccountBalance) error {
	table := rt.newTable([]string{"Address", "ETH", "LINK"})
	for _, ab := range balances {
		table.Append([]string{
			ab.Address,
			ab.EthBalance.String(),
			ab.LinkBalance.String(),
		})
	}
	render("Account Balance", table)
	return nil
}

func (rt RendererTable) renderServiceAgreement(sa presenters.ServiceAgreement) error {
	table := rt.newTable([]string{"ID", "Created At", "Payment", "Expiration"})
	table.Append([]string{
		sa.ID,
		sa.FriendlyCreatedAt(),
		sa.FriendlyPayment(),
		sa.FriendlyExpiration(),
	})
	render("Service Agreement", table)
	return nil
}

func (rt RendererTable) renderExternalInitiatorAuthentication(eia models.ExternalInitiatorAuthentication) error {
	table := rt.newTable([]string{"ID", "Value"})
	table.Append([]string{
		"AccessKey",
		eia.AccessKey,
		"Secret",
		eia.Secret,
	})
	render("External Initiator Credentials", table)
	return nil
}

func (rt RendererTable) newTable(headers []string) *tablewriter.Table {
	table := tablewriter.NewWriter(rt)
	table.SetHeader(headers)
	return table
}

func (rt RendererTable) renderTxAttempts(attempts []models.TxAttempt) error {
	table := rt.newTable([]string{"TxID", "Hash", "GasPrice", "SentAt", "Confirmed"})
	for _, a := range attempts {
		table.Append([]string{
			fmt.Sprint(a.TxID),
			a.Hash.Hex(),
			fmt.Sprint(a.GasPrice),
			fmt.Sprint(a.SentAt),
			fmt.Sprint(a.Confirmed()),
		})
	}

	render("Tx Attempts", table)
	return nil
}
