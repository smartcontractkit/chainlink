package vrfv2

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

// TxPipelineSpec is the classic VRF Coordinator V2 observation source (vrf_pipeline_v2).
type TxPipelineSpec struct {
	Address               string
	EstimateGasMultiplier float64
	FromAddress           string
	SimulationBlock       *string
}

func (d *TxPipelineSpec) Type() string { return "vrf_pipeline_v2" }

func (d *TxPipelineSpec) String() (string, error) {
	optionalSimBlock := ""
	if d.SimulationBlock != nil {
		sb := *d.SimulationBlock
		if sb != "latest" && sb != "pending" {
			return "", fmt.Errorf("invalid SimulationBlock value: %s", sb)
		}
		optionalSimBlock = fmt.Sprintf("block=\"%s\"", sb)
	}
	sourceTemplate := `
decode_log   [type=ethabidecodelog
             abi="RandomWordsRequested(bytes32 indexed keyHash,uint256 requestId,uint256 preSeed,uint64 indexed subId,uint16 minimumRequestConfirmations,uint32 callbackGasLimit,uint32 numWords,address indexed sender)"
             data="$(jobRun.logData)"
             topics="$(jobRun.logTopics)"]
vrf          [type=vrfv2
             publicKey="$(jobSpec.publicKey)"
             requestBlockHash="$(jobRun.logBlockHash)"
             requestBlockNumber="$(jobRun.logBlockNumber)"
             topics="$(jobRun.logTopics)"]
estimate_gas [type=estimategaslimit
             to="{{ .Address }}"
             multiplier="{{ .EstimateGasMultiplier }}"
             data="$(vrf.output)"
			 %s]
simulate [type=ethcall
          from="{{ .FromAddress }}"
          to="{{ .Address }}"
          gas="$(estimate_gas)"
          gasPrice="$(jobSpec.maxGasPrice)"
          extractRevertReason=true
          contract="{{ .Address }}"
          data="$(vrf.output)"
		  %s]
decode_log->vrf->estimate_gas->simulate`

	sourceString := fmt.Sprintf(sourceTemplate, optionalSimBlock, optionalSimBlock)
	return marshallTemplate(d, "VRFv2 pipeline template", sourceString)
}

// JobSpec is the full VRF v2 job TOML (type vrf).
type JobSpec struct {
	Name                          string        `toml:"name"`
	CoordinatorAddress            string        `toml:"coordinatorAddress"`
	BatchCoordinatorAddress       string        `toml:"batchCoordinatorAddress"`
	PublicKey                     string        `toml:"publicKey"`
	ExternalJobID                 string        `toml:"externalJobID"`
	ObservationSource             string        `toml:"observationSource"`
	MinIncomingConfirmations      int           `toml:"minIncomingConfirmations"`
	FromAddresses                 []string      `toml:"fromAddresses"`
	EVMChainID                    string        `toml:"evmChainID"`
	ForwardingAllowed             bool          `toml:"forwardingAllowed"`
	BatchFulfillmentEnabled       bool          `toml:"batchFulfillmentEnabled"`
	BatchFulfillmentGasMultiplier float64       `toml:"batchFulfillmentGasMultiplier"`
	BackOffInitialDelay           time.Duration `toml:"backOffInitialDelay"`
	BackOffMaxDelay               time.Duration `toml:"backOffMaxDelay"`
	PollPeriod                    time.Duration `toml:"pollPeriod"`
	RequestTimeout                time.Duration `toml:"requestTimeout"`
}

func (v *JobSpec) Type() string { return "vrf" }

func (v *JobSpec) String() (string, error) {
	vrfTemplateString := `
type                     = "vrf"
schemaVersion            = 1
name                     = "{{.Name}}"
forwardingAllowed        = {{.ForwardingAllowed}}
coordinatorAddress       = "{{.CoordinatorAddress}}"
{{ if .BatchFulfillmentEnabled }}batchCoordinatorAddress                = "{{.BatchCoordinatorAddress}}"{{ else }}{{ end }}
fromAddresses            = [{{range .FromAddresses}}"{{.}}",{{end}}]
evmChainID               = "{{.EVMChainID}}"
minIncomingConfirmations = {{.MinIncomingConfirmations}}
publicKey                = "{{.PublicKey}}"
externalJobID            = "{{.ExternalJobID}}"
batchFulfillmentEnabled = {{.BatchFulfillmentEnabled}}
batchFulfillmentGasMultiplier = {{.BatchFulfillmentGasMultiplier}}
backoffInitialDelay     = "{{.BackOffInitialDelay}}"
backoffMaxDelay         = "{{.BackOffMaxDelay}}"
pollPeriod              = "{{.PollPeriod}}"
requestTimeout          = "{{.RequestTimeout}}"
customRevertsPipelineEnabled = true
observationSource = """
{{.ObservationSource}}
"""
`
	return marshallTemplate(v, "VRFv2 Job", vrfTemplateString)
}

func marshallTemplate(jobSpec any, name, templateString string) (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.New(name).Parse(templateString)
	if err != nil {
		return "", err
	}
	if err := tmpl.Execute(&buf, jobSpec); err != nil {
		return "", err
	}
	return buf.String(), nil
}
