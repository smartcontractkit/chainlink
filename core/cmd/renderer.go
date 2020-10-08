package cmd

import (
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/olekukonko/tablewriter"
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
	case *[]presenters.Tx:
		return rt.renderTxs(*typed)
	case *presenters.Tx:
		return rt.renderTx(*typed)
	case *presenters.ExternalInitiatorAuthentication:
		return rt.renderExternalInitiatorAuthentication(*typed)
	case *web.ConfigPatchResponse:
		return rt.renderConfigPatchResponse(typed)
	case *presenters.ConfigPrinter:
		return rt.renderConfiguration(*typed)
	default:
		return fmt.Errorf("unable to render object of type %T: %v", typed, typed)
	}
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

	table.Append([]string{
		"ACCOUNT_ADDRESS",
		cp.AccountAddress,
	})

	schemaT := reflect.TypeOf(orm.ConfigSchema{})
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

func (rt RendererTable) renderExternalInitiatorAuthentication(eia presenters.ExternalInitiatorAuthentication) error {
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

func (rt RendererTable) renderTxAttempts(attempts []models.TxAttempt) error {
	table := rt.newTable([]string{"TxID", "Hash", "GasPrice", "SentAt", "Confirmed"})
	for _, a := range attempts {
		table.Append([]string{
			fmt.Sprint(a.TxID),
			a.Hash.Hex(),
			fmt.Sprint(a.GasPrice),
			fmt.Sprint(a.SentAt),
			fmt.Sprint(a.Confirmed),
		})
	}

	render("Ethereum Transaction Attempts", table)
	return nil
}

func (rt RendererTable) renderTx(tx presenters.Tx) error {
	table := rt.newTable([]string{"From", "Nonce", "To", "Confirmed"})
	table.Append([]string{
		tx.From.Hex(),
		tx.Nonce,
		tx.To.Hex(),
		fmt.Sprint(tx.Confirmed),
	})

	render(fmt.Sprintf("Ethereum Transaction %v", tx.Hash.Hex()), table)
	return nil
}

func (rt RendererTable) renderTxs(txs []presenters.Tx) error {
	table := rt.newTable([]string{"Hash", "Nonce", "From", "GasPrice", "SentAt", "Confirmed"})
	for _, tx := range txs {
		table.Append([]string{
			tx.Hash.Hex(),
			tx.Nonce,
			tx.From.Hex(),
			tx.GasPrice,
			tx.SentAt,
			fmt.Sprint(tx.Confirmed),
		})
	}

	render("Ethereum Transactions", table)
	return nil
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
