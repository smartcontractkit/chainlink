package cmd

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/models/ocrkey"
	"github.com/smartcontractkit/chainlink/core/store/models/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
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
	case *presenters.ServiceAgreement:
		return rt.renderServiceAgreement(*typed)
	case *[]presenters.EthTx:
		return rt.renderEthTxs(*typed)
	case *presenters.EthTx:
		return rt.renderEthTx(*typed)
	case *presenters.ExternalInitiatorAuthentication:
		return rt.renderExternalInitiatorAuthentication(*typed)
	case *web.ConfigPatchResponse:
		return rt.renderConfigPatchResponse(typed)
	case *presenters.ConfigPrinter:
		return rt.renderConfiguration(*typed)
	case *presenters.ETHKey:
		return rt.renderETHKeys([]presenters.ETHKey{*typed})
	case *[]presenters.ETHKey:
		return rt.renderETHKeys(*typed)
	case *p2pkey.EncryptedP2PKey:
		return rt.renderP2PKeys([]p2pkey.EncryptedP2PKey{*typed})
	case *[]p2pkey.EncryptedP2PKey:
		return rt.renderP2PKeys(*typed)
	case *ocrkey.EncryptedKeyBundle:
		return rt.renderOCRKeys([]ocrkey.EncryptedKeyBundle{*typed})
	case *[]ocrkey.EncryptedKeyBundle:
		return rt.renderOCRKeys(*typed)
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

func renderList(fields []string, items [][]string) {
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
	fmt.Println(listRendered)
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

func (rt RendererTable) renderEthTx(tx presenters.EthTx) error {
	table := rt.newTable([]string{"From", "Nonce", "To", "State"})
	table.Append([]string{
		tx.From.Hex(),
		tx.Nonce,
		tx.To.Hex(),
		fmt.Sprint(tx.State),
	})

	render(fmt.Sprintf("Ethereum Transaction %v", tx.Hash.Hex()), table)
	return nil
}

func (rt RendererTable) renderEthTxs(txs []presenters.EthTx) error {
	table := rt.newTable([]string{"Hash", "Nonce", "From", "GasPrice", "SentAt", "State"})
	for _, tx := range txs {
		table.Append([]string{
			tx.Hash.Hex(),
			tx.Nonce,
			tx.From.Hex(),
			tx.GasPrice,
			tx.SentAt,
			fmt.Sprint(tx.State),
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

func (rt RendererTable) renderETHKeys(keys []presenters.ETHKey) error {
	var rows [][]string
	for _, key := range keys {
		var nextNonce string
		if key.NextNonce == nil {
			nextNonce = "0"
		} else {
			nextNonce = fmt.Sprintf("%d", *key.NextNonce)
		}
		var lastUsed string
		if key.LastUsed != nil {
			lastUsed = key.LastUsed.String()
		}
		var deletedAt string
		if !key.DeletedAt.IsZero() {
			deletedAt = key.DeletedAt.Time.String()
		}
		rows = append(rows, []string{
			key.Address,
			key.EthBalance.String(),
			key.LinkBalance.String(),
			nextNonce,
			lastUsed,
			fmt.Sprintf("%v", key.IsFunding),
			key.CreatedAt.String(),
			key.UpdatedAt.String(),
			deletedAt,
		})
	}
	fmt.Println("\n🔑 ETH Keys")
	renderList([]string{"Address", "ETH", "LINK", "Next nonce", "Last used", "Is funding", "Created", "Updated", "Deleted"}, rows)
	return nil
}

func (rt RendererTable) renderP2PKeys(p2pKeys []p2pkey.EncryptedP2PKey) error {
	var rows [][]string
	for _, key := range p2pKeys {
		var deletedAt string
		if !key.DeletedAt.IsZero() {
			deletedAt = key.DeletedAt.Time.String()
		}
		rows = append(rows, []string{
			fmt.Sprintf("%v", key.ID),
			fmt.Sprintf("%v", key.PeerID),
			fmt.Sprintf("%v", key.PubKey),
			fmt.Sprintf("%v", key.CreatedAt),
			fmt.Sprintf("%v", key.UpdatedAt),
			fmt.Sprintf("%v", deletedAt),
		})
	}
	fmt.Println("\n🔑 P2P Keys")
	renderList([]string{"ID", "Peer ID", "Public key", "Created", "Updated", "Deleted"}, rows)
	return nil
}

func (rt RendererTable) renderOCRKeys(ocrKeys []ocrkey.EncryptedKeyBundle) error {
	var rows [][]string
	for _, key := range ocrKeys {
		var deletedAt string
		if !key.DeletedAt.IsZero() {
			deletedAt = key.DeletedAt.Time.String()
		}
		rows = append(rows, []string{
			key.ID.String(),
			key.OnChainSigningAddress.String(),
			key.OffChainPublicKey.String(),
			key.ConfigPublicKey.String(),
			key.CreatedAt.String(),
			key.UpdatedAt.String(),
			deletedAt,
		})
	}
	fmt.Println("\n🔑 OCR Keys")
	renderList([]string{"ID", "On-chain signing addr", "Off-chain pubkey", "Config pubkey", "Created", "Updated", "Deleted"}, rows)
	return nil
}
