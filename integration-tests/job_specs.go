package networks

import (
	"bytes"
	"fmt"
	"math/big"
	"text/template"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/client"
)

// CronJobSpec represents a cron job spec
type CronJobSpec struct {
	Schedule          string `toml:"schedule"`          // CRON job style schedule string
	ObservationSource string `toml:"observationSource"` // List of commands for the chainlink node
}

// Type is cron
func (c *CronJobSpec) Type() string { return "cron" }

// String representation of the job
func (c *CronJobSpec) String() (string, error) {
	cronJobTemplateString := `type     = "cron"
schemaVersion     = 1
schedule          = "{{.Schedule}}"
observationSource = """
{{.ObservationSource}}
"""`
	return marshallTemplate(c, "CRON Job", cronJobTemplateString)
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

// VRFTxPipelineSpec VRF request with tx callback
type VRFTxPipelineSpec struct {
	Address string
}

// Type returns the type of the pipeline
func (d *VRFTxPipelineSpec) Type() string {
	return "vrf_pipeline"
}

// String representation of the pipeline
func (d *VRFTxPipelineSpec) String() (string, error) {
	sourceString := `
decode_log   [type=ethabidecodelog
              abi="RandomnessRequest(bytes32 keyHash,uint256 seed,bytes32 indexed jobID,address sender,uint256 fee,bytes32 requestID)"
              data="$(jobRun.logData)"
              topics="$(jobRun.logTopics)"]
vrf          [type=vrf
              publicKey="$(jobSpec.publicKey)"
              requestBlockHash="$(jobRun.logBlockHash)"
              requestBlockNumber="$(jobRun.logBlockNumber)"
              topics="$(jobRun.logTopics)"]
encode_tx    [type=ethabiencode
              abi="fulfillRandomnessRequest(bytes proof)"
              data="{\\"proof\\": $(vrf)}"]
submit_tx  [type=ethtx to="{{.Address}}"
            data="$(encode_tx)"
            txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
decode_log->vrf->encode_tx->submit_tx`
	return marshallTemplate(d, "VRF pipeline template", sourceString)
}

// BlockhashStoreJobSpec represents a blockhashstore job
type BlockhashStoreJobSpec struct {
	Name                  string `toml:"name"`
	CoordinatorV2Address  string `toml:"coordinatorV2Address"` // Address of the VRF Coordinator contract
	WaitBlocks            int    `toml:"waitBlocks"`
	LookbackBlocks        int    `toml:"lookbackBlocks"`
	BlockhashStoreAddress string `toml:"blockhashStoreAddress"`
	PollPeriod            string `toml:"pollPeriod"`
	RunTimeout            string `toml:"runTimeout"`
	EVMChainID            string `toml:"evmChainID"`
}

// Type returns the type of the job
func (b *BlockhashStoreJobSpec) Type() string { return "blockhashstore" }

// String representation of the job
func (b *BlockhashStoreJobSpec) String() (string, error) {
	vrfTemplateString := `
type                     = "blockhashstore"
schemaVersion            = 1
name                     = "{{.Name}}"
coordinatorV2Address     = "{{.CoordinatorV2Address}}"
waitBlocks               = {{.WaitBlocks}}
lookbackBlocks           = {{.LookbackBlocks}}
blockhashStoreAddress    = "{{.BlockhashStoreAddress}}"
pollPeriod               = "{{.PollPeriod}}"
runTimeout               = "{{.RunTimeout}}"
evmChainID               = "{{.EVMChainID}}"
`
	return marshallTemplate(b, "BlockhashStore Job", vrfTemplateString)
}

// VRFJobSpec represents a VRF job
type VRFJobSpec struct {
	Name                     string `toml:"name"`
	CoordinatorAddress       string `toml:"coordinatorAddress"` // Address of the VRF Coordinator contract
	PublicKey                string `toml:"publicKey"`          // Public key of the proving key
	ExternalJobID            string `toml:"externalJobID"`
	ObservationSource        string `toml:"observationSource"` // List of commands for the chainlink node
	MinIncomingConfirmations int    `toml:"minIncomingConfirmations"`
}

// Type returns the type of the job
func (v *VRFJobSpec) Type() string { return "vrf" }

// String representation of the job
func (v *VRFJobSpec) String() (string, error) {
	vrfTemplateString := `
type                     = "vrf"
schemaVersion            = 1
name                     = "{{.Name}}"
coordinatorAddress       = "{{.CoordinatorAddress}}"
minIncomingConfirmations = {{.MinIncomingConfirmations}}
publicKey                = "{{.PublicKey}}"
externalJobID            = "{{.ExternalJobID}}"
observationSource = """
{{.ObservationSource}}
"""
`
	return marshallTemplate(v, "VRF Job", vrfTemplateString)
}

// VRFV2TxPipelineSpec VRFv2 request with tx callback
type VRFV2TxPipelineSpec struct {
	Address string
}

// Type returns the type of the pipeline
func (d *VRFV2TxPipelineSpec) Type() string {
	return "vrf_pipeline_v2"
}

// String representation of the pipeline
func (d *VRFV2TxPipelineSpec) String() (string, error) {
	sourceString := `
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
             multiplier="1.1"
             data="$(vrf.output)"]
simulate [type=ethcall
          to="{{ .Address }}"
          gas="$(estimate_gas)"
          gasPrice="$(jobSpec.maxGasPrice)"
          extractRevertReason=true
          contract="{{ .Address }}"
          data="$(vrf.output)"]
decode_log->vrf->estimate_gas->simulate`
	return marshallTemplate(d, "VRFV2 pipeline template", sourceString)
}

// DirectRequestJobSpec represents a direct request spec
type DirectRequestJobSpec struct {
	Name                     string `toml:"name"`
	ContractAddress          string `toml:"contractAddress"`
	ExternalJobID            string `toml:"externalJobID"`
	MinIncomingConfirmations string `toml:"minIncomingConfirmations"`
	ObservationSource        string `toml:"observationSource"` // List of commands for the chainlink node
}

// Type returns the type of the pipeline
func (d *DirectRequestJobSpec) Type() string { return "directrequest" }

// String representation of the pipeline
func (d *DirectRequestJobSpec) String() (string, error) {
	directRequestTemplateString := `type     = "directrequest"
schemaVersion     = 1
name              = "{{.Name}}"
maxTaskDuration   = "99999s"
contractAddress   = "{{.ContractAddress}}"
externalJobID     = "{{.ExternalJobID}}"
minIncomingConfirmations = {{.MinIncomingConfirmations}}
observationSource = """
{{.ObservationSource}}
"""`
	return marshallTemplate(d, "Direct Request Job", directRequestTemplateString)
}

// FluxMonitorJobSpec represents a flux monitor spec
type FluxMonitorJobSpec struct {
	Name              string        `toml:"name"`
	ContractAddress   string        `toml:"contractAddress"`   // Address of the Flux Monitor script
	Precision         int           `toml:"precision"`         // Optional
	Threshold         float32       `toml:"threshold"`         // Optional
	AbsoluteThreshold float32       `toml:"absoluteThreshold"` // Optional
	IdleTimerPeriod   time.Duration `toml:"idleTimerPeriod"`   // Optional
	IdleTimerDisabled bool          `toml:"idleTimerDisabled"` // Optional
	PollTimerPeriod   time.Duration `toml:"pollTimerPeriod"`   // Optional
	PollTimerDisabled bool          `toml:"pollTimerDisabled"` // Optional
	MaxTaskDuration   time.Duration `toml:"maxTaskDuration"`   // Optional
	ObservationSource string        `toml:"observationSource"` // List of commands for the chainlink node
}

// Type returns the type of the job
func (f *FluxMonitorJobSpec) Type() string { return "fluxmonitor" }

// String representation of the job
func (f *FluxMonitorJobSpec) String() (string, error) {
	fluxMonitorTemplateString := `type              = "fluxmonitor"
schemaVersion     = 1
name              = "{{.Name}}"
contractAddress   = "{{.ContractAddress}}"
precision         ={{if not .Precision}} 0 {{else}} {{.Precision}} {{end}}
threshold         ={{if not .Threshold}} 0.5 {{else}} {{.Threshold}} {{end}}
absoluteThreshold ={{if not .AbsoluteThreshold}} 0.1 {{else}} {{.AbsoluteThreshold}} {{end}}

idleTimerPeriod   ={{if not .IdleTimerPeriod}} "1ms" {{else}} "{{.IdleTimerPeriod}}" {{end}}
idleTimerDisabled ={{if not .IdleTimerDisabled}} false {{else}} {{.IdleTimerDisabled}} {{end}}

pollTimerPeriod   ={{if not .PollTimerPeriod}} "1m" {{else}} "{{.PollTimerPeriod}}" {{end}}
pollTimerDisabled ={{if not .PollTimerDisabled}} false {{else}} {{.PollTimerDisabled}} {{end}}

maxTaskDuration = {{if not .Precision}} "180s" {{else}} {{.Precision}} {{end}}

observationSource = """
{{.ObservationSource}}
"""`
	return marshallTemplate(f, "Flux Monitor Job", fluxMonitorTemplateString)
}

// TODO - OTHER ARREAS AFFECTED
// KeeperJobSpec represents a V2 keeper spec
type KeeperJobSpec struct {
	Name                     string `toml:"name"`
	ContractAddress          string `toml:"contractAddress"`
	FromAddress              string `toml:"fromAddress"` // Hex representation of the from address
	MinIncomingConfirmations int    `toml:"minIncomingConfirmations"`
}

// Type returns the type of the job
func (k *KeeperJobSpec) Type() string { return "keeper" }

// String representation of the job
func (k *KeeperJobSpec) String() (string, error) {
	keeperTemplateString := `
type                     = "keeper"
schemaVersion            = 1
name                     = "{{.Name}}"
contractAddress          = "{{.ContractAddress}}"
fromAddress              = "{{.FromAddress}}"
minIncomingConfirmations = {{.MinIncomingConfirmations}}
`
	return marshallTemplate(k, "Keeper Job", keeperTemplateString)
}

// TODO - OTHER AREAS AFFECTED
// OCRBootstrapJobSpec represents the spec for bootstrapping an OCR job, given to one node that then must be linked
// back to by others by OCRTaskJobSpecs
type OCRBootstrapJobSpec struct {
	Name                     string        `toml:"name"`
	BlockChainTimeout        time.Duration `toml:"blockchainTimeout"`                      // Optional
	ContractConfirmations    int           `toml:"contractConfigConfirmations"`            // Optional
	TrackerPollInterval      time.Duration `toml:"contractConfigTrackerPollInterval"`      // Optional
	TrackerSubscribeInterval time.Duration `toml:"contractConfigTrackerSubscribeInterval"` // Optional
	ContractAddress          string        `toml:"contractAddress"`                        // Address of the OCR contract
	IsBootstrapPeer          bool          `toml:"isBootstrapPeer"`                        // Typically true
	P2PPeerID                string        `toml:"p2pPeerID"`                              // This node's P2P ID
}

// Type returns the type of the job
func (o *OCRBootstrapJobSpec) Type() string { return "offchainreporting" }

// String representation of the job
func (o *OCRBootstrapJobSpec) String() (string, error) {
	ocrTemplateString := `type = "offchainreporting"
schemaVersion                          = 1
blockchainTimeout                      ={{if not .BlockChainTimeout}} "20s" {{else}} {{.BlockChainTimeout}} {{end}}
contractConfigConfirmations            ={{if not .ContractConfirmations}} 3 {{else}} {{.ContractConfirmations}} {{end}}
contractConfigTrackerPollInterval      ={{if not .TrackerPollInterval}} "1m" {{else}} {{.TrackerPollInterval}} {{end}}
contractConfigTrackerSubscribeInterval ={{if not .TrackerSubscribeInterval}} "2m" {{else}} {{.TrackerSubscribeInterval}} {{end}}
contractAddress                        = "{{.ContractAddress}}"
p2pBootstrapPeers                      = []
isBootstrapPeer                        = {{.IsBootstrapPeer}}
p2pPeerID                              = "{{.P2PPeerID}}"`
	return marshallTemplate(o, "OCR Bootstrap Job", ocrTemplateString)
}

// WebhookJobSpec reprsents a webhook job
type WebhookJobSpec struct {
	Name              string `toml:"name"`
	Initiator         string `toml:"initiator"`         // External initiator name
	InitiatorSpec     string `toml:"initiatorSpec"`     // External initiator spec object in stringified form
	ObservationSource string `toml:"observationSource"` // List of commands for the chainlink node
}

// Type returns the type of the job
func (w *WebhookJobSpec) Type() string { return "webhook" }

// String representation of the job
func (w *WebhookJobSpec) String() (string, error) {
	webHookTemplateString := `type = "webhook"
schemaVersion      = 1
name               = "{{.Name}}"
externalInitiators = [
	{ name = "{{.Initiator}}", spec = "{{.InitiatorSpec}}"}
]
observationSource = """
{{.ObservationSource}}
"""`
	return marshallTemplate(w, "Webhook Job", webHookTemplateString)
}

// ObservationSourceSpecHTTP creates a http GET task spec for json data
func ObservationSourceSpecHTTP(url string) string {
	return fmt.Sprintf(`
		fetch [type=http method=GET url="%s"];
		parse [type=jsonparse path="data,result"];
		fetch -> parse;`, url)
}

// OCR2TaskJobSpec represents an OCR2 job that is given to other nodes, meant to communicate with the bootstrap node,
// and provide their answers
type OCR2TaskJobSpec struct {
	Name                     string            `toml:"name"`
	JobType                  string            `toml:"type"`
	ContractID               string            `toml:"contractID"`                             // Address of the OCR contract/account(s)
	Relay                    string            `toml:"relay"`                                  // Name of blockchain relay to use
	PluginType               string            `toml:"pluginType"`                             // Type of report plugin to use
	RelayConfig              map[string]string `toml:"relayConfig"`                            // Relay spec object in stringified form
	P2PV2Bootstrappers       []P2PData         `toml:"p2pv2Bootstrappers"`                     // P2P ID of the bootstrap node
	OCRKeyBundleID           string            `toml:"ocrKeyBundleID"`                         // ID of this node's OCR key bundle
	MonitoringEndpoint       string            `toml:"monitoringEndpoint"`                     // Typically "chain.link:4321"
	TransmitterID            string            `toml:"transmitterID"`                          // ID of address this node will use to transmit
	BlockChainTimeout        time.Duration     `toml:"blockchainTimeout"`                      // Optional
	TrackerSubscribeInterval time.Duration     `toml:"contractConfigTrackerSubscribeInterval"` // Optional
	TrackerPollInterval      time.Duration     `toml:"contractConfigTrackerPollInterval"`      // Optional
	ContractConfirmations    int               `toml:"contractConfigConfirmations"`            // Optional
	ObservationSource        string            `toml:"observationSource"`                      // List of commands for the chainlink node
	JuelsPerFeeCoinSource    string            `toml:"juelsPerFeeCoinSource"`                  // List of commands to fetch JuelsPerFeeCoin value (used to calculate ocr payments)
}

// Type returns the type of the job
func (o *OCR2TaskJobSpec) Type() string { return o.JobType }

// String representation of the job
func (o *OCR2TaskJobSpec) String() (string, error) {
	ocr2TemplateString := `type = "{{ .JobType }}"
schemaVersion                          = 1
blockchainTimeout                      ={{if not .BlockChainTimeout}} "20s" {{else}} "{{.BlockChainTimeout}}" {{end}}
contractConfigConfirmations            ={{if not .ContractConfirmations}} 3 {{else}} {{.ContractConfirmations}} {{end}}
contractConfigTrackerPollInterval      ={{if not .TrackerPollInterval}} "1m" {{else}} "{{.TrackerPollInterval}}" {{end}}
contractConfigTrackerSubscribeInterval ={{if not .TrackerSubscribeInterval}} "2m" {{else}} "{{.TrackerSubscribeInterval}}" {{end}}
name 																	 = "{{.Name}}"
relay																	 = "{{.Relay}}"
contractID		                         = "{{.ContractID}}"
{{if .P2PV2Bootstrappers}}
p2pv2Bootstrappers                      = [
  {{range $peer := .P2PV2Bootstrappers}}
  "{{$peer.PeerID}}@{{$peer.RemoteIP}}:{{if $peer.RemotePort}}{{$peer.RemotePort}}{{else}}6690{{end}}",
  {{end}}
]
{{else}}
p2pv2Bootstrappers                      = []
{{end}}
monitoringEndpoint                     ={{if not .MonitoringEndpoint}} "chain.link:4321" {{else}} "{{.MonitoringEndpoint}}" {{end}}
{{if eq .JobType "offchainreporting2" }}
pluginType                             = "{{ .PluginType }}"
ocrKeyBundleID                         = "{{.OCRKeyBundleID}}"
transmitterID                     		 = "{{.TransmitterID}}"
observationSource                      = """
{{.ObservationSource}}
"""
[pluginConfig]
juelsPerFeeCoinSource                  = """
{{.JuelsPerFeeCoinSource}}
"""
{{end}}

[relayConfig]
{{range $key, $value := .RelayConfig}}
{{$key}} = "{{$value}}"
{{end}}`

	return marshallTemplate(o, "OCR2 Job", ocr2TemplateString)
}

// P2PData holds the remote ip and the peer id and port
type P2PData struct {
	RemoteIP   string
	RemotePort string
	PeerID     string
}

// VRFV2JobSpec represents a VRFV2 job
type VRFV2JobSpec struct {
	Name                     string        `toml:"name"`
	CoordinatorAddress       string        `toml:"coordinatorAddress"` // Address of the VRF Coordinator contract
	PublicKey                string        `toml:"publicKey"`          // Public key of the proving key
	ExternalJobID            string        `toml:"externalJobID"`
	ObservationSource        string        `toml:"observationSource"` // List of commands for the chainlink node
	MinIncomingConfirmations int           `toml:"minIncomingConfirmations"`
	FromAddress              string        `toml:"fromAddress"`
	EVMChainID               string        `toml:"evmChainID"`
	BatchFulfillmentEnabled  bool          `toml:"batchFulfillmentEnabled"`
	BackOffInitialDelay      time.Duration `toml:"backOffInitialDelay"`
	BackOffMaxDelay          time.Duration `toml:"backOffMaxDelay"`
}

// Type returns the type of the job
func (v *VRFV2JobSpec) Type() string { return "vrf" }

// String representation of the job
func (v *VRFV2JobSpec) String() (string, error) {
	vrfTemplateString := `
type                     = "vrf"
schemaVersion            = 1
name                     = "{{.Name}}"
coordinatorAddress       = "{{.CoordinatorAddress}}"
fromAddress              = "{{.FromAddress}}"
evmChainID               = "{{.EVMChainID}}"
minIncomingConfirmations = {{.MinIncomingConfirmations}}
publicKey                = "{{.PublicKey}}"
externalJobID            = "{{.ExternalJobID}}"
batchFulfillmentEnabled = {{.BatchFulfillmentEnabled}}
backoffInitialDelay     = "{{.BackOffInitialDelay}}"
backoffMaxDelay         = "{{.BackOffMaxDelay}}"
observationSource = """
{{.ObservationSource}}
"""
`
	return marshallTemplate(v, "VRFV2 Job", vrfTemplateString)
}

// EncodeOnChainVRFProvingKey encodes uncompressed public VRF key to on-chain representation
func EncodeOnChainVRFProvingKey(vrfKey client.VRFKey) ([2]*big.Int, error) {
	uncompressed := vrfKey.Data.Attributes.Uncompressed
	provingKey := [2]*big.Int{}
	var set1 bool
	var set2 bool
	// strip 0x to convert to int
	provingKey[0], set1 = new(big.Int).SetString(uncompressed[2:66], 16)
	if !set1 {
		return [2]*big.Int{}, fmt.Errorf("can not convert VRF key to *big.Int")
	}
	provingKey[1], set2 = new(big.Int).SetString(uncompressed[66:], 16)
	if !set2 {
		return [2]*big.Int{}, fmt.Errorf("can not convert VRF key to *big.Int")
	}
	return provingKey, nil
}
