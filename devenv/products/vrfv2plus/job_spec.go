package vrfv2plus

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

// TxPipelineSpec defines the observation source pipeline for a VRFv2Plus job.
type TxPipelineSpec struct {
	Address               string
	EstimateGasMultiplier float64
	FromAddress           string
	SimulationBlock       *string // nil, "latest", or "pending"
}

func (d *TxPipelineSpec) Type() string { return "vrf_pipeline_v2plus" }

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
             abi="RandomWordsRequested(bytes32 indexed keyHash,uint256 requestId,uint256 preSeed,uint256 indexed subId,uint16 minimumRequestConfirmations,uint32 callbackGasLimit,uint32 numWords,bytes extraArgs,address indexed sender)"
             data="$(jobRun.logData)"
             topics="$(jobRun.logTopics)"]
generate_proof [type=vrfv2plus
                publicKey="$(jobSpec.publicKey)"
                requestBlockHash="$(jobRun.logBlockHash)"
                requestBlockNumber="$(jobRun.logBlockNumber)"
                topics="$(jobRun.logTopics)"]
estimate_gas [type=estimategaslimit
             to="{{ .Address }}"
             multiplier="{{ .EstimateGasMultiplier }}"
             data="$(generate_proof.output)"
			 %s]
simulate_fulfillment [type=ethcall
					  from="{{ .FromAddress }}"
                      to="{{ .Address }}"
                      gas="$(estimate_gas)"
                      gasPrice="$(jobSpec.maxGasPrice)"
                      extractRevertReason=true
                      contract="{{ .Address }}"
                      data="$(generate_proof.output)"
					  %s]
decode_log->generate_proof->estimate_gas->simulate_fulfillment`

	sourceString := fmt.Sprintf(sourceTemplate, optionalSimBlock, optionalSimBlock)
	return marshallTemplate(d, "VRFV2 Plus pipeline template", sourceString)
}

// JobSpec defines the full VRFv2Plus job for a Chainlink node.
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
observationSource = """
{{.ObservationSource}}
"""
`
	return marshallTemplate(v, "VRFV2 PLUS Job", vrfTemplateString)
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
