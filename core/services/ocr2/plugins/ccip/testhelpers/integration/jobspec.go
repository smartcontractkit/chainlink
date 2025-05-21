package integrationtesthelpers

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// OCR2TaskJobSpec represents an OCR2 job that is given to other nodes, meant to communicate with the bootstrap node,
// and provide their answers
type OCR2TaskJobSpec struct {
	Name              string `toml:"name"`
	JobType           string `toml:"type"`
	MaxTaskDuration   string `toml:"maxTaskDuration"` // Optional
	ForwardingAllowed bool   `toml:"forwardingAllowed"`
	OCR2OracleSpec    job.OCR2OracleSpec
	ObservationSource string `toml:"observationSource"` // List of commands for the Chainlink node
	ExternalJobID     string `toml:"externalJobID"`
}

// Type returns the type of the job
func (o *OCR2TaskJobSpec) Type() string { return o.JobType }

// String representation of the job
func (o *OCR2TaskJobSpec) String() (string, error) {
	var feedID string
	if o.OCR2OracleSpec.FeedID != nil {
		feedID = o.OCR2OracleSpec.FeedID.Hex()
	}
	externalID, err := ExternalJobID(o.Name)
	if err != nil {
		return "", err
	}
	specWrap := struct {
		Name                     string
		JobType                  string
		ExternalJobID            string
		MaxTaskDuration          string
		ForwardingAllowed        bool
		ContractID               string
		FeedID                   string
		Relay                    string
		PluginType               string
		RelayConfig              map[string]interface{}
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
		ExternalJobID:         externalID,
		ForwardingAllowed:     o.ForwardingAllowed,
		MaxTaskDuration:       o.MaxTaskDuration,
		ContractID:            o.OCR2OracleSpec.ContractID,
		FeedID:                feedID,
		Relay:                 o.OCR2OracleSpec.Relay,
		PluginType:            string(o.OCR2OracleSpec.PluginType),
		RelayConfig:           o.OCR2OracleSpec.RelayConfig,
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
externalJobID                          = "{{.ExternalJobID}}"
forwardingAllowed                      = {{.ForwardingAllowed}}
{{if .MaxTaskDuration}}
maxTaskDuration                        = "{{ .MaxTaskDuration }}" {{end}}
{{if .PluginType}}
pluginType                             = "{{ .PluginType }}" {{end}}
relay                                  = "{{.Relay}}"
schemaVersion                          = 1
contractID                             = "{{.ContractID}}"
{{if .FeedID}}
feedID                                 = "{{.FeedID}}"
{{end}}
{{if eq .JobType "offchainreporting2" }}
ocrKeyBundleID                         = "{{.OCRKeyBundleID}}" {{end}}
{{if eq .JobType "offchainreporting2" }}
transmitterID                          = "{{.TransmitterID}}" {{end}}
{{if .BlockchainTimeout}}
blockchainTimeout                      = "{{.BlockchainTimeout}}"
{{end}}
{{if .ContractConfirmations}}
contractConfigConfirmations            = {{.ContractConfirmations}}
{{end}}
{{if .TrackerPollInterval}}
contractConfigTrackerPollInterval      = "{{.TrackerPollInterval}}"
{{end}}
{{if .TrackerSubscribeInterval}}
contractConfigTrackerSubscribeInterval = "{{.TrackerSubscribeInterval}}"
{{end}}
{{if .P2PV2Bootstrappers}}
p2pv2Bootstrappers                     = [{{range .P2PV2Bootstrappers}}"{{.}}",{{end}}]{{end}}
{{if .MonitoringEndpoint}}
monitoringEndpoint                     = "{{.MonitoringEndpoint}}" {{end}}
{{if .ObservationSource}}
observationSource                      = """
{{.ObservationSource}}
"""{{end}}
{{if eq .JobType "offchainreporting2" }}
[pluginConfig]{{range $key, $value := .PluginConfig}}
{{$key}} = {{$value}}{{end}}
{{end}}
[relayConfig]{{range $key, $value := .RelayConfig}}
{{$key}} = {{$value}}{{end}}
`
	return MarshallTemplate(specWrap, "OCR2 Job", ocr2TemplateString)
}

// MarshallTemplate Helper to marshall templates
func MarshallTemplate(jobSpec interface{}, name, templateString string) (string, error) {
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

type JobType string

const (
	Commit    JobType = "commit"
	Execution JobType = "exec"
	Bootstrap JobType = "bootstrap"
)

func JobName(jobType JobType, source string, destination, version string) string {
	if version != "" {
		return fmt.Sprintf("ccip-%s-%s-%s-%s", jobType, source, destination, version)
	}
	return fmt.Sprintf("ccip-%s-%s-%s", jobType, source, destination)
}

type CCIPJobSpecParams struct {
	Name                   string
	Version                string
	OffRamp                common.Address
	CommitStore            common.Address
	SourceChainName        string
	DestChainName          string
	DestEvmChainId         uint64
	TokenPricesUSDPipeline string
	PriceGetterConfig      string
	SourceStartBlock       uint64
	DestStartBlock         uint64
	USDCAttestationAPI     string
	USDCConfig             *config.USDCConfig
	LBTCConfig             *config.LBTCConfig
	P2PV2Bootstrappers     pq.StringArray
}

func (params CCIPJobSpecParams) Validate() error {
	if params.CommitStore == common.HexToAddress("0x0") {
		return errors.New("must set commit store address")
	}
	return nil
}

func (params CCIPJobSpecParams) ValidateCommitJobSpec() error {
	commonErr := params.Validate()
	if commonErr != nil {
		return commonErr
	}
	if params.OffRamp == common.HexToAddress("0x0") {
		return errors.New("OffRamp cannot be empty for execution job")
	}
	// Validate token prices config
	// NB: only validate the dynamic price getter config if present since we could also be using the pipeline instead.
	// NB: make this test mandatory once we switch to dynamic price getter only.
	if params.PriceGetterConfig != "" {
		if _, err := pricegetter.NewDynamicPriceGetterConfig(params.PriceGetterConfig); err != nil {
			return fmt.Errorf("invalid price getter config: %w", err)
		}
	}
	return nil
}

func (params CCIPJobSpecParams) ValidateExecJobSpec() error {
	commonErr := params.Validate()
	if commonErr != nil {
		return commonErr
	}
	if params.OffRamp == common.HexToAddress("0x0") {
		return errors.New("OffRamp cannot be empty for execution job")
	}
	return nil
}

// CommitJobSpec generates template for CCIP-relay job spec.
// OCRKeyBundleID,TransmitterID need to be set from the calling function
func (params CCIPJobSpecParams) CommitJobSpec() (*OCR2TaskJobSpec, error) {
	err := params.ValidateCommitJobSpec()
	if err != nil {
		return nil, fmt.Errorf("invalid job spec params: %w", err)
	}

	pluginConfig := map[string]interface{}{
		"offRamp": fmt.Sprintf(`"%s"`, params.OffRamp.Hex()),
	}
	if params.TokenPricesUSDPipeline != "" {
		pluginConfig["tokenPricesUSDPipeline"] = fmt.Sprintf(`"""
%s
"""`, params.TokenPricesUSDPipeline)
	}
	if params.PriceGetterConfig != "" {
		pluginConfig["priceGetterConfig"] = fmt.Sprintf(`"""
%s
"""`, params.PriceGetterConfig)
	}

	ocrSpec := job.OCR2OracleSpec{
		Relay:                             relay.NetworkEVM,
		PluginType:                        types.CCIPCommit,
		ContractID:                        params.CommitStore.Hex(),
		ContractConfigConfirmations:       1,
		ContractConfigTrackerPollInterval: models.Interval(20 * time.Second),
		P2PV2Bootstrappers:                params.P2PV2Bootstrappers,
		PluginConfig:                      pluginConfig,
		RelayConfig: map[string]interface{}{
			"chainID": params.DestEvmChainId,
		},
	}
	if params.DestStartBlock > 0 {
		ocrSpec.PluginConfig["destStartBlock"] = params.DestStartBlock
	}
	if params.SourceStartBlock > 0 {
		ocrSpec.PluginConfig["sourceStartBlock"] = params.SourceStartBlock
	}
	return &OCR2TaskJobSpec{
		OCR2OracleSpec: ocrSpec,
		JobType:        "offchainreporting2",
		Name:           JobName(Commit, params.SourceChainName, params.DestChainName, params.Version),
	}, nil
}

// ExecutionJobSpec generates template for CCIP-execution job spec.
// OCRKeyBundleID,TransmitterID need to be set from the calling function
func (params CCIPJobSpecParams) ExecutionJobSpec() (*OCR2TaskJobSpec, error) {
	err := params.ValidateExecJobSpec()
	if err != nil {
		return nil, err
	}
	ocrSpec := job.OCR2OracleSpec{
		Relay:                             relay.NetworkEVM,
		PluginType:                        types.CCIPExecution,
		ContractID:                        params.OffRamp.Hex(),
		ContractConfigConfirmations:       1,
		ContractConfigTrackerPollInterval: models.Interval(20 * time.Second),

		P2PV2Bootstrappers: params.P2PV2Bootstrappers,
		PluginConfig:       map[string]interface{}{},
		RelayConfig: map[string]interface{}{
			"chainID": params.DestEvmChainId,
		},
	}
	if params.DestStartBlock > 0 {
		ocrSpec.PluginConfig["destStartBlock"] = params.DestStartBlock
	}
	if params.SourceStartBlock > 0 {
		ocrSpec.PluginConfig["sourceStartBlock"] = params.SourceStartBlock
	}
	if params.USDCAttestationAPI != "" {
		ocrSpec.PluginConfig["USDCConfig.AttestationAPI"] = fmt.Sprintf("\"%s\"", params.USDCAttestationAPI)
		ocrSpec.PluginConfig["USDCConfig.SourceTokenAddress"] = fmt.Sprintf("\"%s\"", utils.RandomAddress().String())
		ocrSpec.PluginConfig["USDCConfig.SourceMessageTransmitterAddress"] = fmt.Sprintf("\"%s\"", utils.RandomAddress().String())
		ocrSpec.PluginConfig["USDCConfig.AttestationAPITimeoutSeconds"] = 5
	}
	if params.USDCConfig != nil {
		ocrSpec.PluginConfig["USDCConfig.AttestationAPI"] = fmt.Sprintf(`"%s"`, params.USDCConfig.AttestationAPI)
		ocrSpec.PluginConfig["USDCConfig.SourceTokenAddress"] = fmt.Sprintf(`"%s"`, params.USDCConfig.SourceTokenAddress)
		ocrSpec.PluginConfig["USDCConfig.SourceMessageTransmitterAddress"] = fmt.Sprintf(`"%s"`, params.USDCConfig.SourceMessageTransmitterAddress)
		ocrSpec.PluginConfig["USDCConfig.AttestationAPITimeoutSeconds"] = params.USDCConfig.AttestationAPITimeoutSeconds
	}
	if params.LBTCConfig != nil {
		ocrSpec.PluginConfig["LBTCConfig.AttestationAPI"] = fmt.Sprintf(`"%s"`, params.LBTCConfig.AttestationAPI)
		ocrSpec.PluginConfig["LBTCConfig.SourceTokenAddress"] = fmt.Sprintf(`"%s"`, params.LBTCConfig.SourceTokenAddress)
		ocrSpec.PluginConfig["LBTCConfig.AttestationAPITimeoutSeconds"] = params.LBTCConfig.AttestationAPITimeoutSeconds
	}
	return &OCR2TaskJobSpec{
		OCR2OracleSpec: ocrSpec,
		JobType:        "offchainreporting2",
		Name:           JobName(Execution, params.SourceChainName, params.DestChainName, params.Version),
	}, err
}

func (params CCIPJobSpecParams) BootstrapJob(contractID string) *OCR2TaskJobSpec {
	bootstrapSpec := job.OCR2OracleSpec{
		ContractID:                        contractID,
		Relay:                             relay.NetworkEVM,
		ContractConfigConfirmations:       1,
		ContractConfigTrackerPollInterval: models.Interval(20 * time.Second),
		RelayConfig: map[string]interface{}{
			"chainID": params.DestEvmChainId,
		},
	}
	return &OCR2TaskJobSpec{
		Name:           fmt.Sprintf("%s-%s-%s", Bootstrap, params.SourceChainName, params.DestChainName),
		JobType:        "bootstrap",
		OCR2OracleSpec: bootstrapSpec,
	}
}

func (c *CCIPIntegrationTestHarness) NewCCIPJobSpecParams(tokenPricesUSDPipeline string, priceGetterConfig string, configBlock uint64, usdcAttestationAPI string) CCIPJobSpecParams {
	return CCIPJobSpecParams{
		CommitStore:            c.Dest.CommitStore.Address(),
		OffRamp:                c.Dest.OffRamp.Address(),
		DestEvmChainId:         c.Dest.ChainID,
		SourceChainName:        "SimulatedSource",
		DestChainName:          "SimulatedDest",
		TokenPricesUSDPipeline: tokenPricesUSDPipeline,
		PriceGetterConfig:      priceGetterConfig,
		DestStartBlock:         configBlock,
		USDCAttestationAPI:     usdcAttestationAPI,
	}
}

func ExternalJobID(jobName string) (string, error) {
	in := []byte(jobName)
	sha256Hash := sha256.New()
	sha256Hash.Write(in)
	in = sha256Hash.Sum(nil)[:16]
	// tag as valid UUID v4 https://github.com/google/uuid/blob/0f11ee6918f41a04c201eceeadf612a377bc7fbc/version4.go#L53-L54
	in[6] = (in[6] & 0x0f) | 0x40 // Version 4
	in[8] = (in[8] & 0x3f) | 0x80 // Variant is 10
	id, err := uuid.FromBytes(in)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
