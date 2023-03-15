package client

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/services/job"
)

// EIServiceConfig represents External Initiator service config
type EIServiceConfig struct {
	URL string
}

// ChainlinkConfig represents the variables needed to connect to a Chainlink node
type ChainlinkConfig struct {
	URL      string
	Email    string
	Password string
	RemoteIP string
}

// ResponseSlice is the generic model that can be used for all Chainlink API responses that are an slice
type ResponseSlice struct {
	Data []map[string]interface{}
}

// Response is the generic model that can be used for all Chainlink API responses
type Response struct {
	Data map[string]interface{}
}

// JobRunsResponse job runs
type JobRunsResponse struct {
	Data []RunsResponseData `json:"data"`
	Meta RunsMetaResponse   `json:"meta"`
}

// RunsResponseData runs response data
type RunsResponseData struct {
	Type       string                 `json:"type"`
	ID         string                 `json:"id"`
	Attributes RunsAttributesResponse `json:"attributes"`
}

// RunsAttributesResponse runs attributes
type RunsAttributesResponse struct {
	Meta       interface{}   `json:"meta"`
	Errors     []interface{} `json:"errors"`
	Inputs     RunInputs     `json:"inputs"`
	TaskRuns   []TaskRun     `json:"taskRuns"`
	CreatedAt  time.Time     `json:"createdAt"`
	FinishedAt time.Time     `json:"finishedAt"`
}

// DecodeLogTaskRun is "ethabidecodelog" task run info,
// also used for "RequestID" tracing in perf tests
type DecodeLogTaskRun struct {
	Fee       int    `json:"fee"`
	JobID     []int  `json:"jobID"`
	KeyHash   []int  `json:"keyHash"`
	RequestID []byte `json:"requestID"`
	Sender    string `json:"sender"`
}

// TaskRun is pipeline task run info
type TaskRun struct {
	Type       string      `json:"type"`
	CreatedAt  time.Time   `json:"createdAt"`
	FinishedAt time.Time   `json:"finishedAt"`
	Output     string      `json:"output"`
	Error      interface{} `json:"error"`
	DotID      string      `json:"dotId"`
}

type NodeKeysBundle struct {
	OCR2Key    OCR2Key
	PeerID     string
	TXKey      TxKey
	P2PKeys    P2PKeys
	EthAddress string
}

// RunInputs run inputs (value)
type RunInputs struct {
	Parse int `json:"parse"`
}

// RunsMetaResponse runs meta
type RunsMetaResponse struct {
	Count int `json:"count"`
}

// BridgeType is the model that represents the bridge when read or created on a Chainlink node
type BridgeType struct {
	Data BridgeTypeData `json:"data"`
}

// BridgeTypeData is the model that represents the bridge when read or created on a Chainlink node
type BridgeTypeData struct {
	Attributes BridgeTypeAttributes `json:"attributes"`
}

// BridgeTypeAttributes is the model that represents the bridge when read or created on a Chainlink node
type BridgeTypeAttributes struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	RequestData string `json:"requestData,omitempty"`
}

// Session is the form structure used for authenticating
type Session struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ExportedEVMKey holds all details needed to recreate a private key of the Chainlink node
type ExportedEVMKey struct {
	Address string `json:"address"`
	Crypto  struct {
		Cipher       string `json:"cipher"`
		CipherText   string `json:"ciphertext"`
		CipherParams struct {
			Iv string `json:"iv"`
		} `json:"cipherparams"`
		Kdf       string `json:"kdf"`
		KDFParams struct {
			DkLen int    `json:"dklen"`
			N     int    `json:"n"`
			P     int    `json:"p"`
			R     int    `json:"r"`
			Salt  string `json:"salt"`
		} `json:"kdfparams"`
		Mac string `json:"mac"`
	} `json:"crypto"`
	ID      string `json:"id"`
	Version int    `json:"version"`
}

// VRFExportKey is the model that represents the exported VRF key
type VRFExportKey struct {
	PublicKey string `json:"PublicKey"`
	VrfKey    struct {
		Address string `json:"address"`
		Crypto  struct {
			Cipher       string `json:"cipher"`
			Ciphertext   string `json:"ciphertext"`
			Cipherparams struct {
				Iv string `json:"iv"`
			} `json:"cipherparams"`
			Kdf       string `json:"kdf"`
			Kdfparams struct {
				Dklen int    `json:"dklen"`
				N     int    `json:"n"`
				P     int    `json:"p"`
				R     int    `json:"r"`
				Salt  string `json:"salt"`
			} `json:"kdfparams"`
			Mac string `json:"mac"`
		} `json:"crypto"`
		Version int `json:"version"`
	} `json:"vrf_key"`
}

// VRFKeyAttributes is the model that represents the created VRF key attributes when read
type VRFKeyAttributes struct {
	Compressed   string      `json:"compressed"`
	Uncompressed string      `json:"uncompressed"`
	Hash         string      `json:"hash"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	DeletedAt    interface{} `json:"deletedAt"`
}

// VRFKeyData is the model that represents the created VRF key's data when read
type VRFKeyData struct {
	Type       string           `json:"type"`
	ID         string           `json:"id"`
	Attributes VRFKeyAttributes `json:"attributes"`
}

// VRFKey is the model that represents the created VRF key when read
type VRFKey struct {
	Data VRFKeyData `json:"data"`
}

// VRFKeys is the model that represents the created VRF keys when read
type VRFKeys struct {
	Data []VRFKey `json:"data"`
}

// DKGSignKeyAttributes is the model that represents the created DKG Sign key attributes when read
type DKGSignKeyAttributes struct {
	PublicKey string `json:"publicKey"`
}

// DKGSignKeyData is the model that represents the created DKG Sign key's data when read
type DKGSignKeyData struct {
	Type       string               `json:"type"`
	ID         string               `json:"id"`
	Attributes DKGSignKeyAttributes `json:"attributes"`
}

// DKGSignKey is the model that represents the created DKG Sign key when read
type DKGSignKey struct {
	Data DKGSignKeyData `json:"data"`
}

// DKGSignKeys is the model that represents the created DKGSignData key when read
type DKGSignKeys struct {
	Data []DKGSignKey `json:"data"`
}

// DKGEncryptKeyAttributes is the model that represents the created DKG Encrypt key attributes when read
type DKGEncryptKeyAttributes struct {
	PublicKey string `json:"publicKey"`
}

// DKGEncryptKeyData is the model that represents the created DKG Encrypt key's data when read
type DKGEncryptKeyData struct {
	Type       string                  `json:"type"`
	ID         string                  `json:"id"`
	Attributes DKGEncryptKeyAttributes `json:"attributes"`
}

// DKGEncryptKey is the model that represents the created DKG Encrypt key when read
type DKGEncryptKey struct {
	Data DKGEncryptKeyData `json:"data"`
}

// DKGEncryptKeys is the model that represents the created DKGEncryptKeys key when read
type DKGEncryptKeys struct {
	Data []DKGEncryptKey `json:"data"`
}

// OCRKeys is the model that represents the created OCR keys when read
type OCRKeys struct {
	Data []OCRKeyData `json:"data"`
}

// OCRKey is the model that represents the created OCR keys when read
type OCRKey struct {
	Data OCRKeyData `json:"data"`
}

// OCRKeyData is the model that represents the created OCR keys when read
type OCRKeyData struct {
	Attributes OCRKeyAttributes `json:"attributes"`
	ID         string           `json:"id"`
}

// OCRKeyAttributes is the model that represents the created OCR keys when read
type OCRKeyAttributes struct {
	ConfigPublicKey       string `json:"configPublicKey"`
	OffChainPublicKey     string `json:"offChainPublicKey"`
	OnChainSigningAddress string `json:"onChainSigningAddress"`
}

// OCR2Keys is the model that represents the created OCR2 keys when read
type OCR2Keys struct {
	Data []OCR2KeyData `json:"data"`
}

// OCR2Key is the model that represents the created OCR2 keys when read
type OCR2Key struct {
	Data OCR2KeyData `json:"data"`
}

// OCR2KeyData is the model that represents the created OCR2 keys when read
type OCR2KeyData struct {
	Type       string            `json:"type"`
	Attributes OCR2KeyAttributes `json:"attributes"`
	ID         string            `json:"id"`
}

// OCR2KeyAttributes is the model that represents the created OCR2 keys when read
type OCR2KeyAttributes struct {
	ChainType         string `json:"chainType"`
	ConfigPublicKey   string `json:"configPublicKey"`
	OffChainPublicKey string `json:"offchainPublicKey"`
	OnChainPublicKey  string `json:"onchainPublicKey"`
}

// P2PKeys is the model that represents the created P2P keys when read
type P2PKeys struct {
	Data []P2PKeyData `json:"data"`
}

// P2PKey is the model that represents the created P2P keys when read
type P2PKey struct {
	Data P2PKeyData `json:"data"`
}

// P2PKeyData is the model that represents the created P2P keys when read
type P2PKeyData struct {
	Attributes P2PKeyAttributes `json:"attributes"`
}

// P2PKeyAttributes is the model that represents the created P2P keys when read
type P2PKeyAttributes struct {
	ID        int    `json:"id"`
	PeerID    string `json:"peerId"`
	PublicKey string `json:"publicKey"`
}

// CSAKeys is the model that represents the created CSA keys when read
type CSAKeys struct {
	Data []CSAKeyData `json:"data"`
}

// CSAKey is the model that represents the created CSA key when created
type CSAKey struct {
	Data CSAKeyData `json:"data"`
}

// CSAKeyData is the model that represents the created CSA key when read
type CSAKeyData struct {
	Type       string           `json:"type"`
	ID         string           `json:"id"`
	Attributes CSAKeyAttributes `json:"attributes"`
}

// CSAKeyAttributes is the model that represents the attributes of a CSA Key
type CSAKeyAttributes struct {
	PublicKey string `json:"publicKey"`
	Version   int    `json:"version"`
}

// ETHKeys is the model that represents the created ETH keys when read
type ETHKeys struct {
	Data []ETHKeyData `json:"data"`
}

// ETHKey is the model that represents the created ETH keys when read
type ETHKey struct {
	Data ETHKeyData `json:"data"`
}

// ETHKeyData is the model that represents the created ETH keys when read
type ETHKeyData struct {
	Attributes ETHKeyAttributes `json:"attributes"`
}

// ETHKeyAttributes is the model that represents the created ETH keys when read
type ETHKeyAttributes struct {
	Address    string `json:"address"`
	ETHBalance string `json:"ethBalance"`
	ChainID    string `json:"evmChainID"`
}

// TxKeys is the model that represents the created keys when read
type TxKeys struct {
	Data []TxKeyData `json:"data"`
}

// TxKey is the model that represents the created keys when read
type TxKey struct {
	Data TxKeyData `json:"data"`
}

// TxKeyData is the model that represents the created keys when read
type TxKeyData struct {
	Type       string          `json:"type"`
	ID         string          `json:"id"`
	Attributes TxKeyAttributes `json:"attributes"`
}

// TxKeyAttributes is the model that represents the created keys when read
type TxKeyAttributes struct {
	PublicKey string `json:"publicKey"`

	// starknet specific (uses contract model instead of EOA)
	AccountAddr string `json:"accountAddr,omitempty"`
	StarkKey    string `json:"starkPubKey,omitempty"`
}

type SingleTransactionDataWrapper struct {
	Data TransactionData `json:"data"`
}

type SendEtherRequest struct {
	DestinationAddress string `json:"address"`
	FromAddress        string `json:"from"`
	Amount             string `json:"amount"`
	EVMChainID         int    `json:"evmChainID,omitempty"`
	AllowHigherAmounts bool   `json:"allowHigherAmounts"`
}

// EIAttributes is the model that represents the EI keys when created and read
type EIAttributes struct {
	Name              string `json:"name,omitempty"`
	URL               string `json:"url,omitempty"`
	IncomingAccessKey string `json:"incomingAccessKey,omitempty"`
	AccessKey         string `json:"accessKey,omitempty"`
	Secret            string `json:"incomingSecret,omitempty"`
	OutgoingToken     string `json:"outgoingToken,omitempty"`
	OutgoingSecret    string `json:"outgoingSecret,omitempty"`
}

// EIKeys is the model that represents the EI configs when read
type EIKeys struct {
	Data []EIKey `json:"data"`
}

// EIKeyCreate is the model that represents the EI config when created
type EIKeyCreate struct {
	Data EIKey `json:"data"`
}

// EIKey is the model that represents the EI configs when read
type EIKey struct {
	Attributes EIAttributes `json:"attributes"`
}

type SolanaChainConfig struct {
	BlockRate           null.String
	ConfirmPollPeriod   null.String
	OCR2CachePollPeriod null.String
	OCR2CacheTTL        null.String
	TxTimeout           null.String
	SkipPreflight       null.Bool
	Commitment          null.String
}

// SolanaChainAttributes is the model that represents the solana chain
type SolanaChainAttributes struct {
	ChainID string            `json:"chainID"`
	Config  SolanaChainConfig `json:"config"`
}

// SolanaChain is the model that represents the solana chain when read
type SolanaChain struct {
	Attributes SolanaChainAttributes `json:"attributes"`
}

// SolanaChainCreate is the model that represents the solana chain when created
type SolanaChainCreate struct {
	Data SolanaChain `json:"data"`
}

// SolanaNodeAttributes is the model that represents the solana noded
type SolanaNodeAttributes struct {
	Name          string `json:"name"`
	SolanaChainID string `json:"solanaChainId" db:"solana_chain_id"`
	SolanaURL     string `json:"solanaURL" db:"solana_url"`
}

// SolanaNode is the model that represents the solana node when read
type SolanaNode struct {
	Attributes SolanaNodeAttributes `json:"attributes"`
}

// SolanaNodeCreate is the model that represents the solana node when created
type SolanaNodeCreate struct {
	Data SolanaNode `json:"data"`
}

type StarkNetChainConfig struct {
	OCR2CachePollPeriod null.String
	OCR2CacheTTL        null.String
	RequestTimeout      null.String
	TxTimeout           null.Bool
	TxSendFrequency     null.String
	TxMaxBatchSize      null.String
}

// StarkNetChainAttributes is the model that represents the starknet chain
type StarkNetChainAttributes struct {
	Type    string              `json:"type"`
	ChainID string              `json:"chainID"`
	Config  StarkNetChainConfig `json:"config"`
}

// StarkNetChain is the model that represents the starknet chain when read
type StarkNetChain struct {
	Attributes StarkNetChainAttributes `json:"attributes"`
}

// StarkNetChainCreate is the model that represents the starknet chain when created
type StarkNetChainCreate struct {
	Data StarkNetChain `json:"data"`
}

// StarkNetNodeAttributes is the model that represents the starknet node
type StarkNetNodeAttributes struct {
	Name    string `json:"name"`
	ChainID string `json:"chainId"`
	Url     string `json:"url"`
}

// StarkNetNode is the model that represents the starknet node when read
type StarkNetNode struct {
	Attributes StarkNetNodeAttributes `json:"attributes"`
}

// StarkNetNodeCreate is the model that represents the starknet node when created
type StarkNetNodeCreate struct {
	Data StarkNetNode `json:"data"`
}

// SpecForm is the form used when creating a v2 job spec, containing the TOML of the v2 job
type SpecForm struct {
	TOML string `json:"toml"`
}

// Spec represents a job specification that contains information about the job spec
type Spec struct {
	Data SpecData `json:"data"`
}

// SpecData contains the ID of the job spec
type SpecData struct {
	ID string `json:"id"`
}

// JobForm is the form used when creating a v2 job spec, containing the TOML of the v2 job
type JobForm struct {
	TOML string `json:"toml"`
}

// Job contains the job data for a given job
type Job struct {
	Data JobData `json:"data"`
}

// JobData contains the ID for a given job
type JobData struct {
	ID string `json:"id"`
}

// JobSpec represents the different possible job types that Chainlink nodes can handle
type JobSpec interface {
	Type() string
	// String Returns TOML representation of the job
	String() (string, error)
}

// CronJobSpec represents a cron job spec
type CronJobSpec struct {
	Schedule          string `toml:"schedule"`          // CRON job style schedule string
	ObservationSource string `toml:"observationSource"` // List of commands for the Chainlink node
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

// PipelineSpec common API call pipeline
type PipelineSpec struct {
	BridgeTypeAttributes BridgeTypeAttributes
	DataPath             string
}

// Type is common_pipeline
func (d *PipelineSpec) Type() string {
	return "common_pipeline"
}

// String representation of the pipeline
func (d *PipelineSpec) String() (string, error) {
	sourceString := `
		fetch [type=bridge name="{{.BridgeTypeAttributes.Name}}" requestData="{{.BridgeTypeAttributes.RequestData}}"];
		parse [type=jsonparse path="{{.DataPath}}"];
		fetch -> parse;`
	return marshallTemplate(d, "API call pipeline template", sourceString)
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

// DirectRequestTxPipelineSpec oracle request with tx callback
type DirectRequestTxPipelineSpec struct {
	BridgeTypeAttributes BridgeTypeAttributes
	DataPath             string
}

// Type returns the type of the pipeline
func (d *DirectRequestTxPipelineSpec) Type() string {
	return "directrequest_pipeline"
}

// String representation of the pipeline
func (d *DirectRequestTxPipelineSpec) String() (string, error) {
	sourceString := `
            decode_log   [type=ethabidecodelog
                         abi="OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes data)"
                         data="$(jobRun.logData)"
                         topics="$(jobRun.logTopics)"]
			encode_tx  [type=ethabiencode
                        abi="fulfill(bytes32 _requestId, uint256 _data)"
                        data=<{
                          "_requestId": $(decode_log.requestId),
                          "_data": $(parse)
                         }>
                       ]
			fetch  [type=bridge name="{{.BridgeTypeAttributes.Name}}" requestData="{{.BridgeTypeAttributes.RequestData}}"];
			parse  [type=jsonparse path="{{.DataPath}}"]
            submit [type=ethtx to="$(decode_log.requester)" data="$(encode_tx)" failOnRevert=true]
			decode_log -> fetch -> parse -> encode_tx -> submit`
	return marshallTemplate(d, "Direct request pipeline template", sourceString)
}

// DirectRequestJobSpec represents a direct request spec
type DirectRequestJobSpec struct {
	Name                     string `toml:"name"`
	ContractAddress          string `toml:"contractAddress"`
	ExternalJobID            string `toml:"externalJobID"`
	MinIncomingConfirmations string `toml:"minIncomingConfirmations"`
	ObservationSource        string `toml:"observationSource"` // List of commands for the Chainlink node
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
	ObservationSource string        `toml:"observationSource"` // List of commands for the Chainlink node
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

// OCRTaskJobSpec represents an OCR job that is given to other nodes, meant to communicate with the bootstrap node,
// and provide their answers
type OCRTaskJobSpec struct {
	Name                     string        `toml:"name"`
	BlockChainTimeout        time.Duration `toml:"blockchainTimeout"`                      // Optional
	ContractConfirmations    int           `toml:"contractConfigConfirmations"`            // Optional
	TrackerPollInterval      time.Duration `toml:"contractConfigTrackerPollInterval"`      // Optional
	TrackerSubscribeInterval time.Duration `toml:"contractConfigTrackerSubscribeInterval"` // Optional
	ForwardingAllowed        bool          `toml:"forwardingAllowed"`                      // Optional, by default false
	ContractAddress          string        `toml:"contractAddress"`                        // Address of the OCR contract
	P2PBootstrapPeers        []*Chainlink  `toml:"p2pBootstrapPeers"`                      // P2P ID of the bootstrap node
	IsBootstrapPeer          bool          `toml:"isBootstrapPeer"`                        // Typically false
	P2PPeerID                string        `toml:"p2pPeerID"`                              // This node's P2P ID
	KeyBundleID              string        `toml:"keyBundleID"`                            // ID of this node's OCR key bundle
	MonitoringEndpoint       string        `toml:"monitoringEndpoint"`                     // Typically "chain.link:4321"
	TransmitterAddress       string        `toml:"transmitterAddress"`                     // ETH address this node will use to transmit its answer
	ObservationSource        string        `toml:"observationSource"`                      // List of commands for the Chainlink node
}

// P2PData holds the remote ip and the peer id and port
type P2PData struct {
	RemoteIP   string
	RemotePort string
	PeerID     string
}

func (p *P2PData) P2PV2Bootstrapper() string {
	if p.RemotePort == "" {
		p.RemotePort = "6690"
	}
	return fmt.Sprintf("%s@%s:%s", p.PeerID, p.RemoteIP, p.RemotePort)
}

// Type returns the type of the job
func (o *OCRTaskJobSpec) Type() string { return "offchainreporting" }

// String representation of the job
func (o *OCRTaskJobSpec) String() (string, error) {
	// Pre-process P2P data for easier templating
	peers := []P2PData{}
	for _, peer := range o.P2PBootstrapPeers {
		p2pKeys, err := peer.MustReadP2PKeys()
		if err != nil {
			return "", err
		}
		peers = append(peers, P2PData{
			RemoteIP: peer.RemoteIP(),
			PeerID:   p2pKeys.Data[0].Attributes.PeerID,
		})
	}
	specWrap := struct {
		Name                     string
		BlockChainTimeout        time.Duration
		ContractConfirmations    int
		TrackerPollInterval      time.Duration
		TrackerSubscribeInterval time.Duration
		ContractAddress          string
		P2PBootstrapPeers        []P2PData
		IsBootstrapPeer          bool
		P2PPeerID                string
		KeyBundleID              string
		MonitoringEndpoint       string
		TransmitterAddress       string
		ObservationSource        string
		ForwardingAllowed        bool
	}{
		Name:                     o.Name,
		BlockChainTimeout:        o.BlockChainTimeout,
		ContractConfirmations:    o.ContractConfirmations,
		TrackerPollInterval:      o.TrackerPollInterval,
		TrackerSubscribeInterval: o.TrackerSubscribeInterval,
		ContractAddress:          o.ContractAddress,
		P2PBootstrapPeers:        peers,
		IsBootstrapPeer:          o.IsBootstrapPeer,
		P2PPeerID:                o.P2PPeerID,
		KeyBundleID:              o.KeyBundleID,
		MonitoringEndpoint:       o.MonitoringEndpoint,
		TransmitterAddress:       o.TransmitterAddress,
		ObservationSource:        o.ObservationSource,
		ForwardingAllowed:        o.ForwardingAllowed,
	}
	// Results in /dns4//tcp/6690/p2p/12D3KooWAuC9xXBnadsYJpqzZZoB4rMRWqRGpxCrr2mjS7zCoAdN\
	ocrTemplateString := `type = "offchainreporting"
schemaVersion                          = 1
blockchainTimeout                      ={{if not .BlockChainTimeout}} "20s" {{else}} {{.BlockChainTimeout}} {{end}}
contractConfigConfirmations            ={{if not .ContractConfirmations}} 3 {{else}} {{.ContractConfirmations}} {{end}}
contractConfigTrackerPollInterval      ={{if not .TrackerPollInterval}} "1m" {{else}} {{.TrackerPollInterval}} {{end}}
contractConfigTrackerSubscribeInterval ={{if not .TrackerSubscribeInterval}} "2m" {{else}} {{.TrackerSubscribeInterval}} {{end}}
contractAddress                        = "{{.ContractAddress}}"
{{if .P2PBootstrapPeers}}
p2pBootstrapPeers                      = [
  {{range $peer := .P2PBootstrapPeers}}
  "/dns4/{{$peer.RemoteIP}}/tcp/6690/p2p/{{$peer.PeerID}}",
  {{end}}
]
{{else}}
p2pBootstrapPeers                      = []
{{end}}
isBootstrapPeer                        = {{.IsBootstrapPeer}}
p2pPeerID                              = "{{.P2PPeerID}}"
keyBundleID                            = "{{.KeyBundleID}}"
monitoringEndpoint                     ={{if not .MonitoringEndpoint}} "chain.link:4321" {{else}} "{{.MonitoringEndpoint}}" {{end}}
transmitterAddress                     = "{{.TransmitterAddress}}"
forwardingAllowed					   = {{.ForwardingAllowed}}
observationSource                      = """
{{.ObservationSource}}
"""`

	return marshallTemplate(specWrap, "OCR Job", ocrTemplateString)
}

// OCR2TaskJobSpec represents an OCR2 job that is given to other nodes, meant to communicate with the bootstrap node,
// and provide their answers
type OCR2TaskJobSpec struct {
	Name              string `toml:"name"`
	JobType           string `toml:"type"`
	MaxTaskDuration   string `toml:"maxTaskDuration"` // Optional
	OCR2OracleSpec    job.OCR2OracleSpec
	ObservationSource string `toml:"observationSource"` // List of commands for the Chainlink node
}

// Type returns the type of the job
func (o *OCR2TaskJobSpec) Type() string { return o.JobType }

// String representation of the job
func (o *OCR2TaskJobSpec) String() (string, error) {
	specWrap := struct {
		Name                     string
		JobType                  string
		MaxTaskDuration          string
		ContractID               string
		Relay                    string
		PluginType               string
		RelayConfig              map[string]interface{}
		RelayConfigMercuryConfig map[string]interface{}
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
		Name:                     o.Name,
		JobType:                  o.JobType,
		MaxTaskDuration:          o.MaxTaskDuration,
		ContractID:               o.OCR2OracleSpec.ContractID,
		Relay:                    string(o.OCR2OracleSpec.Relay),
		PluginType:               string(o.OCR2OracleSpec.PluginType),
		RelayConfig:              o.OCR2OracleSpec.RelayConfig,
		RelayConfigMercuryConfig: o.OCR2OracleSpec.RelayConfigMercuryConfig,
		PluginConfig:             o.OCR2OracleSpec.PluginConfig,
		P2PV2Bootstrappers:       o.OCR2OracleSpec.P2PV2Bootstrappers,
		OCRKeyBundleID:           o.OCR2OracleSpec.OCRKeyBundleID.String,
		MonitoringEndpoint:       o.OCR2OracleSpec.MonitoringEndpoint.String,
		TransmitterID:            o.OCR2OracleSpec.TransmitterID.String,
		BlockchainTimeout:        o.OCR2OracleSpec.BlockchainTimeout.Duration(),
		ContractConfirmations:    o.OCR2OracleSpec.ContractConfigConfirmations,
		TrackerPollInterval:      o.OCR2OracleSpec.ContractConfigTrackerPollInterval.Duration(),
		ObservationSource:        o.ObservationSource,
	}
	ocr2TemplateString := `
type                                   = "{{ .JobType }}"
name                                   = "{{.Name}}"
{{if .MaxTaskDuration}}
maxTaskDuration                        = "{{ .MaxTaskDuration }}" {{end}}
{{if .PluginType}}
pluginType                             = "{{ .PluginType }}" {{end}}
relay                                  = "{{.Relay}}"
schemaVersion                          = 1
contractID                             = "{{.ContractID}}"
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
{{if .RelayConfigMercuryConfig}}
[relayConfig.MercuryConfig]{{range $key, $value := .RelayConfigMercuryConfig}}
{{$key}} = "{{$value}}"{{end}}
{{end}}
`
	return marshallTemplate(specWrap, "OCR2 Job", ocr2TemplateString)
}

// VRFV2JobSpec represents a VRFV2 job
type VRFV2JobSpec struct {
	Name                     string        `toml:"name"`
	CoordinatorAddress       string        `toml:"coordinatorAddress"` // Address of the VRF Coordinator contract
	PublicKey                string        `toml:"publicKey"`          // Public key of the proving key
	ExternalJobID            string        `toml:"externalJobID"`
	ObservationSource        string        `toml:"observationSource"` // List of commands for the Chainlink node
	MinIncomingConfirmations int           `toml:"minIncomingConfirmations"`
	FromAddresses            []string      `toml:"fromAddresses"`
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
fromAddresses            = [{{range .FromAddresses}}"{{.}}",{{end}}]
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

// VRFJobSpec represents a VRF job
type VRFJobSpec struct {
	Name                     string `toml:"name"`
	CoordinatorAddress       string `toml:"coordinatorAddress"` // Address of the VRF Coordinator contract
	PublicKey                string `toml:"publicKey"`          // Public key of the proving key
	ExternalJobID            string `toml:"externalJobID"`
	ObservationSource        string `toml:"observationSource"` // List of commands for the Chainlink node
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

// WebhookJobSpec reprsents a webhook job
type WebhookJobSpec struct {
	Name              string `toml:"name"`
	Initiator         string `toml:"initiator"`         // External initiator name
	InitiatorSpec     string `toml:"initiatorSpec"`     // External initiator spec object in stringified form
	ObservationSource string `toml:"observationSource"` // List of commands for the Chainlink node
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

// ObservationSourceSpecBridge creates a bridge task spec for json data
func ObservationSourceSpecBridge(bta BridgeTypeAttributes) string {
	return fmt.Sprintf(`
		fetch [type=bridge name="%s" requestData="%s"];
		parse [type=jsonparse path="data,result"];
		fetch -> parse;`, bta.Name, bta.RequestData)
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

type TransactionsData struct {
	Data []TransactionData    `json:"data"`
	Meta TransactionsMetaData `json:"meta"`
}

type TransactionData struct {
	Type       string                `json:"type"`
	ID         string                `json:"id"`
	Attributes TransactionAttributes `json:"attributes"`
}

type TransactionAttributes struct {
	State    string `json:"state"`
	Data     string `json:"data"`
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	ChainID  string `json:"evmChainID"`
	GasLimit string `json:"gasLimit"`
	GasPrice string `json:"gasPrice"`
	Hash     string `json:"hash"`
	RawHex   string `json:"rawHex"`
	Nonce    string `json:"nonce"`
	SentAt   string `json:"sentAt"`
}

type TransactionsMetaData struct {
	Count int `json:"count"`
}

// ChainlinkProfileResults holds the results of asking the Chainlink node to run a PPROF session
type ChainlinkProfileResults struct {
	Reports                 []*ChainlinkProfileResult
	ScheduledProfileSeconds int // How long the profile was scheduled to last
	ActualRunSeconds        int // How long the target function to profile actually took to execute
	NodeIndex               int
}

// ChainlinkProfileResult contains the result of a single PPROF run
type ChainlinkProfileResult struct {
	Type string
	Data []byte
}

// NewBlankChainlinkProfileResults returns all the standard types of profile results with blank data
func NewBlankChainlinkProfileResults() *ChainlinkProfileResults {
	results := &ChainlinkProfileResults{
		Reports: make([]*ChainlinkProfileResult, 0),
	}
	profileStrings := []string{
		"allocs", // A sampling of all past memory allocations
		"block",  // Stack traces that led to blocking on synchronization primitives
		// "cmdline",      // The command line invocation of the current program
		"goroutine",    // Stack traces of all current goroutines
		"heap",         // A sampling of memory allocations of live objects.
		"mutex",        // Stack traces of holders of contended mutexes
		"profile",      // CPU profile.
		"threadcreate", // Stack traces that led to the creation of new OS threads
		"trace",        // A trace of execution of the current program.
	}
	for _, profile := range profileStrings {
		results.Reports = append(results.Reports, &ChainlinkProfileResult{Type: profile})
	}
	return results
}

type CLNodesWithKeys struct {
	Node       *Chainlink
	KeysBundle NodeKeysBundle
}

// Forwarder the model that represents the created Forwarder when created
type Forwarder struct {
	Data ForwarderData `json:"data"`
}

// Forwarders is the model that represents the created Forwarders when read
type Forwarders struct {
	Data []Forwarder `json:"data"`
}

// ForwarderData is the model that represents the created Forwarder when read
type ForwarderData struct {
	ID        string    `json:"id"`
	Address   string    `json:"address"`
	ChainID   string    `json:"chainId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ForwarderAttributes is the model that represents attributes of a Forwarder
type ForwarderAttributes struct {
	Address string `json:"address"`
	ChainID string `json:"chainID"`
}
