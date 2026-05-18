package integrationtesthelpers

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func NewJobSpec(params *LMJobSpecParams) (*OCR3TaskJobSpec, error) {
	params.Defaults()
	if err := params.Validate(); err != nil {
		return nil, err
	}
	return params.JobSpec(), nil
}

func NewBootsrapJobSpec(params *LMJobSpecParams) (*OCR3TaskJobSpec, error) {
	params.Defaults()
	if err := params.Validate(); err != nil {
		return nil, err
	}
	return params.BootstrapJobSpec(), nil
}

type LMJobSpecParams struct {
	Name                    string
	Type                    string
	ChainID                 uint64
	ContractID              string
	OCRKeyBundleID          string
	TransmitterID           string
	RelayFromBlock          int64
	FollowerChains          string
	LiquidityManagerAddress common.Address
	NetworkSelector         uint64
	CfgTrackerInterval      time.Duration
	P2PV2Bootstrappers      pq.StringArray
}

func (p *LMJobSpecParams) Defaults() {
	if len(p.Name) == 0 {
		p.Name = "liquiditymanager"
	}
	if len(p.Type) == 0 {
		p.Type = "ping-pong"
	}
	if p.CfgTrackerInterval == 0 {
		p.CfgTrackerInterval = 15 * time.Second
	}
}

func (p *LMJobSpecParams) Validate() error {
	if len(p.ContractID) == 0 {
		return fmt.Errorf("ContractID is required")
	}

	// TODO: Add more validations

	return nil
}

func (p *LMJobSpecParams) JobSpec() *OCR3TaskJobSpec {
	// NOTE: doing a workaround to specify rebalancerConfig in the pluginConfig
	pluginConfig := job.JSONConfig{
		"closePluginTimeoutSec":   fmt.Sprintf(`%d`, 10),
		"liquidityManagerAddress": fmt.Sprintf(`"%s"`, p.LiquidityManagerAddress.Hex()),
		"liquidityManagerNetwork": fmt.Sprintf(`"%d"`, p.NetworkSelector) + `
[pluginConfig.rebalancerConfig]
type = ` + fmt.Sprintf("\"%s\"\n", p.Type),
	}

	relayCfg := map[string]interface{}{
		"chainID":   p.ChainID,
		"fromBlock": p.RelayFromBlock,
	}

	ocrSpec := job.OCR2OracleSpec{
		Relay:                             relay.NetworkEVM,
		PluginType:                        "liquiditymanager",
		ContractID:                        p.ContractID,
		OCRKeyBundleID:                    null.StringFrom(p.OCRKeyBundleID),
		TransmitterID:                     null.StringFrom(p.TransmitterID),
		ContractConfigConfirmations:       1,
		ContractConfigTrackerPollInterval: models.Interval(p.CfgTrackerInterval),
		P2PV2Bootstrappers:                p.P2PV2Bootstrappers,
		PluginConfig:                      pluginConfig,
		RelayConfig:                       relayCfg,
	}

	return &OCR3TaskJobSpec{
		Name:              p.Name,
		JobType:           "offchainreporting2",
		MaxTaskDuration:   "30s",
		ForwardingAllowed: false,
		OCR2OracleSpec:    ocrSpec,
	}
}

func (p *LMJobSpecParams) BootstrapJobSpec() *OCR3TaskJobSpec {
	bootstrapSpec := job.OCR2OracleSpec{
		ContractID:                        p.ContractID,
		Relay:                             relay.NetworkEVM,
		ContractConfigConfirmations:       1,
		ContractConfigTrackerPollInterval: models.Interval(p.CfgTrackerInterval),
		RelayConfig: map[string]interface{}{
			"chainID":   p.ChainID,
			"fromBlock": p.RelayFromBlock,
		},
	}
	return &OCR3TaskJobSpec{
		Name:           fmt.Sprintf("bootstrap-%d-%s", p.ChainID, p.ContractID),
		JobType:        "bootstrap",
		OCR2OracleSpec: bootstrapSpec,
	}
}

// OCR3TaskJobSpec represents an OCR2 job for chainlink node
type OCR3TaskJobSpec struct {
	Name              string `toml:"name"`
	JobType           string `toml:"type"`
	MaxTaskDuration   string `toml:"maxTaskDuration"` // Optional
	ForwardingAllowed bool   `toml:"forwardingAllowed"`
	OCR2OracleSpec    job.OCR2OracleSpec
	ObservationSource string `toml:"observationSource"` // List of commands for the Chainlink node
}

// Type returns the type of the job
func (o *OCR3TaskJobSpec) Type() string {
	return o.JobType
}

// String representation of the job
func (o *OCR3TaskJobSpec) String() (string, error) {
	var feedID string
	if o.OCR2OracleSpec.FeedID != nil {
		feedID = o.OCR2OracleSpec.FeedID.Hex()
	}
	relayConfig, err := toml.Marshal(struct {
		RelayConfig job.JSONConfig `toml:"relayConfig"`
	}{RelayConfig: o.OCR2OracleSpec.RelayConfig})
	if err != nil {
		return "", fmt.Errorf("failed to marshal relay config: %w", err)
	}
	specWrap := struct {
		Name                     string
		JobType                  string
		MaxTaskDuration          string
		ForwardingAllowed        bool
		ContractID               string
		FeedID                   string
		Relay                    string
		PluginType               string
		RelayConfig              string
		PluginConfig             map[string]interface{}
		P2PV2Bootstrappers       []string
		OCRKeyBundleID           string
		MonitoringEndpoint       string
		TransmitterID            string
		BlockchainTimeout        time.Duration
		TrackerSubscribeInterval time.Duration
		TrackerPollInterval      time.Duration
		ContractConfirmations    uint16
		ObservationSource        string
	}{
		Name:                  o.Name,
		JobType:               o.JobType,
		ForwardingAllowed:     o.ForwardingAllowed,
		MaxTaskDuration:       o.MaxTaskDuration,
		ContractID:            o.OCR2OracleSpec.ContractID,
		FeedID:                feedID,
		Relay:                 o.OCR2OracleSpec.Relay,
		PluginType:            string(o.OCR2OracleSpec.PluginType),
		RelayConfig:           string(relayConfig),
		PluginConfig:          o.OCR2OracleSpec.PluginConfig,
		P2PV2Bootstrappers:    o.OCR2OracleSpec.P2PV2Bootstrappers,
		OCRKeyBundleID:        o.OCR2OracleSpec.OCRKeyBundleID.String,
		MonitoringEndpoint:    o.OCR2OracleSpec.MonitoringEndpoint.String,
		TransmitterID:         o.OCR2OracleSpec.TransmitterID.String,
		BlockchainTimeout:     o.OCR2OracleSpec.BlockchainTimeout.Duration(),
		ContractConfirmations: o.OCR2OracleSpec.ContractConfigConfirmations,
		TrackerPollInterval:   o.OCR2OracleSpec.ContractConfigTrackerPollInterval.Duration(),
		ObservationSource:     o.ObservationSource,
	}
	ocr2TemplateString := `
type                                   = "{{ .JobType }}"
name                                   = "{{.Name}}"
forwardingAllowed                      = {{.ForwardingAllowed}}
{{- if .MaxTaskDuration}}
maxTaskDuration                        = "{{ .MaxTaskDuration }}" {{end}}
{{- if .PluginType}}
pluginType                             = "{{ .PluginType }}" {{end}}
relay                                  = "{{.Relay}}"
schemaVersion                          = 1
contractID                             = "{{.ContractID}}"
{{- if .FeedID}}
feedID                                 = "{{.FeedID}}"
{{end}}
{{- if eq .JobType "offchainreporting2" }}
ocrKeyBundleID                         = "{{.OCRKeyBundleID}}" {{end}}
{{- if eq .JobType "offchainreporting2" }}
transmitterID                          = "{{.TransmitterID}}" {{end}}
{{- if .BlockchainTimeout}}
blockchainTimeout                      = "{{.BlockchainTimeout}}"
{{end}}
{{- if .ContractConfirmations}}
contractConfigConfirmations            = {{.ContractConfirmations}}
{{end}}
{{- if .TrackerPollInterval}}
contractConfigTrackerPollInterval      = "{{.TrackerPollInterval}}"
{{end}}
{{- if .TrackerSubscribeInterval}}
contractConfigTrackerSubscribeInterval = "{{.TrackerSubscribeInterval}}"
{{end}}
{{- if .P2PV2Bootstrappers}}
p2pv2Bootstrappers                     = [{{range .P2PV2Bootstrappers}}"{{.}}",{{end}}]{{end}}
{{- if .MonitoringEndpoint}}
monitoringEndpoint                     = "{{.MonitoringEndpoint}}" {{end}}
{{- if .ObservationSource}}
observationSource                      = """
{{.ObservationSource}}
"""{{end}}
{{if eq .JobType "offchainreporting2" }}
[pluginConfig]{{range $key, $value := .PluginConfig}}
{{$key}} = {{$value}}{{end}}
{{end}}
{{.RelayConfig}}
`
	return marshallTemplate(specWrap, "OCR2 Job", ocr2TemplateString)
}

// marshallTemplate Helper to marshall templates
func marshallTemplate(jobSpec interface{}, name, templateString string) (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.New(name).Parse(templateString)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buf, jobSpec)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}
