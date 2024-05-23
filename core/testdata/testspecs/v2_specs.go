package testspecs

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/webhook"
)

var (
	CronSpecTemplate = `
type                = "cron"
schemaVersion       = 1
schedule            = "CRON_TZ=UTC * 0 0 1 1 *"
externalJobID       =  "%s"
observationSource   = """
ds          [type=http method=GET url="https://chain.link/ETH-USD"];
ds_parse    [type=jsonparse path="data,price"];
ds_multiply [type=multiply times=100];
ds -> ds_parse -> ds_multiply;
"""
`
	CronSpecDotSepTemplate = `
type                = "cron"
schemaVersion       = 1
schedule            = "CRON_TZ=UTC * 0 0 1 1 *"
externalJobID       =  "%s"
observationSource   = """
ds          [type=http method=GET url="https://chain.link/ETH-USD"];
ds_parse    [type=jsonparse path="data.price" separator="."];
ds_multiply [type=multiply times=100];
ds -> ds_parse -> ds_multiply;
"""
`
	DirectRequestSpecNoExternalJobID = `
type                = "directrequest"
schemaVersion       = 1
name                = "%s"
contractAddress     = "0x613a38AC1659769640aaE063C651F48E0250454C"
evmChainID 			= "0"
observationSource   = """
    ds1          [type=http method=GET url="http://example.com" allowunrestrictednetworkaccess="true"];
    ds1_parse    [type=jsonparse path="USD"];
    ds1_multiply [type=multiply times=100];
    ds1 -> ds1_parse -> ds1_multiply;
"""
`
	DirectRequestSpecTemplate = `
type                = "directrequest"
schemaVersion       = 1
name                = "%s"
contractAddress     = "0x613a38AC1659769640aaE063C651F48E0250454C"
externalJobID       =  "%s"
evmChainID 			= "0"
observationSource   = """
    ds1          [type=http method=GET url="http://example.com" allowunrestrictednetworkaccess="true"];
    ds1_parse    [type=jsonparse path="USD"];
    ds1_multiply [type=multiply times=100];
    ds1 -> ds1_parse -> ds1_multiply;
"""
`
	DirectRequestSpecWithRequestersAndMinContractPaymentTemplate = `
type                         = "directrequest"
schemaVersion                = 1
requesters                   = ["0xaaaa1F8ee20f5565510B84f9353F1E333E753B7a", "0xbbbb70F0e81C6F3430dfdC9fa02fB22BdD818C4e"]
minContractPaymentLinkJuels  = "1000000000000000000000"
name                         = "%s"
contractAddress              = "0x613a38AC1659769640aaE063C651F48E0250454C"
externalJobID                = "%s"
evmChainID                   = 0
observationSource            = """
    ds1          [type=http method=GET url="http://example.com" allowunrestrictednetworkaccess="true"];
    ds1_parse    [type=jsonparse path="USD"];
    ds1_multiply [type=multiply times=100];
    ds1 -> ds1_parse -> ds1_multiply;
"""
`
	FluxMonitorSpecTemplate = `
type                = "fluxmonitor"
schemaVersion       = 1
name                = "%s"
contractAddress     = "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
externalJobID       =  "%s"
evmChainID          = 0
threshold = 0.5
absoluteThreshold = 0.0 # optional

idleTimerPeriod = "1s"
idleTimerDisabled = false

pollTimerPeriod = "1m"
pollTimerDisabled = false

observationSource = """
// data source 1
ds1 [type=http method=GET url="https://pricesource1.com" requestData="{\\"coin\\": \\"ETH\\", \\"market\\": \\"USD\\"}"];
ds1_parse [type=jsonparse path="latest"];

// data source 2
ds2 [type=http method=GET url="https://pricesource1.com" requestData="{\\"coin\\": \\"ETH\\", \\"market\\": \\"USD\\"}"];
ds2_parse [type=jsonparse path="latest"];

ds1 -> ds1_parse -> answer1;
ds2 -> ds2_parse -> answer1;

answer1 [type=median index=0];
"""
`

	OCR2EVMSpecMinimalTemplate = `type = "offchainreporting2"
schemaVersion = 1
name = "%s"
relay = "evm"
contractID = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pv2Bootstrappers = []
transmitterID = "0xF67D0290337bca0847005C7ffD1BC75BA9AAE6e4"
pluginType         = "median"
observationSource = """
	ds          [type=http method=GET url="https://chain.link/ETH-USD"];
	ds_parse    [type=jsonparse path="data.price" separator="."];
	ds_multiply [type=multiply times=100];
	ds -> ds_parse -> ds_multiply;
"""
[relayConfig]
chainID = 0
[pluginConfig]
`
	WebhookSpecNoBodyTemplate = `
type            = "webhook"
schemaVersion   = 1
externalJobID   = "%s"
observationSource   = """
    fetch          [type=bridge name="%s"]
    parse_request  [type=jsonparse path="data,result"];
    multiply       [type=multiply times="100"];
    submit         [type=bridge name="%s" includeInputAtKey="result"];

    fetch -> parse_request -> multiply -> submit;
"""
`

	WebhookSpecWithBodyTemplate = `
type            = "webhook"
schemaVersion   = 1
externalJobID   = "%s"
observationSource   = """
    parse_request  [type=jsonparse path="data,result" data="$(jobRun.requestBody)"];
    multiply       [type=multiply times="100"];
    send_to_bridge [type=bridge name="%s" includeInputAtKey="result" ];

    parse_request -> multiply -> send_to_bridge;
"""
`

	OCRBootstrapSpec = `
type			= "bootstrap"
name			= "%s"
relay			= "evm"
schemaVersion	= 1
contractID		= "0x613a38AC1659769640aaE063C651F48E0250454C"
[relayConfig]
chainID			= 1337
`
)

func GetOCRBootstrapSpec() string {
	return fmt.Sprintf(OCRBootstrapSpec, uuid.New())
}

func GetDirectRequestSpec() string {
	uuid := uuid.New()
	return GetDirectRequestSpecWithUUID(uuid)
}

func GetDirectRequestSpecWithUUID(u uuid.UUID) string {
	return fmt.Sprintf(DirectRequestSpecTemplate, u, u)
}

func GetOCR2EVMSpecMinimal() string {
	return fmt.Sprintf(OCR2EVMSpecMinimalTemplate, uuid.New())
}

func GetWebhookSpecNoBody(u uuid.UUID, fetchBridge, submitBridge string) string {
	return fmt.Sprintf(WebhookSpecNoBodyTemplate, u, fetchBridge, submitBridge)
}

type KeeperSpecParams struct {
	Name              string
	ContractAddress   string
	FromAddress       string
	EvmChainID        int
	ObservationSource string
}

type KeeperSpec struct {
	KeeperSpecParams
	toml string
}

func (os KeeperSpec) Toml() string {
	return os.toml
}

func GenerateKeeperSpec(params KeeperSpecParams) KeeperSpec {
	template := `
type            		 	= "keeper"
schemaVersion   		 	= 1
name            		 	= "%s"
contractAddress 		 	= "%s"
fromAddress     		 	= "%s"
evmChainID      		 	= %d
externalJobID   		 	=  "123e4567-e89b-12d3-a456-426655440002"
observationSource = """%s"""
`
	escapedObvSource := strings.ReplaceAll(params.ObservationSource, `\`, `\\`)
	return KeeperSpec{
		KeeperSpecParams: params,
		toml:             fmt.Sprintf(template, params.Name, params.ContractAddress, params.FromAddress, params.EvmChainID, escapedObvSource),
	}
}

type VRFSpecParams struct {
	JobID                         string
	Name                          string
	CoordinatorAddress            string
	VRFVersion                    vrfcommon.Version
	BatchCoordinatorAddress       string
	VRFOwnerAddress               string
	BatchFulfillmentEnabled       bool
	CustomRevertsPipelineEnabled  bool
	BatchFulfillmentGasMultiplier float64
	MinIncomingConfirmations      int
	FromAddresses                 []string
	PublicKey                     string
	ObservationSource             string
	EVMChainID                    string
	RequestedConfsDelay           int
	RequestTimeout                time.Duration
	V2                            bool
	ChunkSize                     int
	BackoffInitialDelay           time.Duration
	BackoffMaxDelay               time.Duration
	GasLanePrice                  *assets.Wei
	PollPeriod                    time.Duration
}

type VRFSpec struct {
	VRFSpecParams
	toml string
}

func (vs VRFSpec) Toml() string {
	return vs.toml
}

func GenerateVRFSpec(params VRFSpecParams) VRFSpec {
	jobID := "123e4567-e89b-12d3-a456-426655440000"
	if params.JobID != "" {
		jobID = params.JobID
	}
	name := "vrf-primary"
	if params.Name != "" {
		name = params.Name
	}
	vrfVersion := vrfcommon.V2
	if params.VRFVersion != "" {
		vrfVersion = params.VRFVersion
	}
	coordinatorAddress := "0xABA5eDc1a551E55b1A570c0e1f1055e5BE11eca7"
	if params.CoordinatorAddress != "" {
		coordinatorAddress = params.CoordinatorAddress
	}
	batchCoordinatorAddress := "0x5C7B1d96CA3132576A84423f624C2c492f668Fea"
	if params.BatchCoordinatorAddress != "" {
		batchCoordinatorAddress = params.BatchCoordinatorAddress
	}
	vrfOwnerAddress := "0x5383C25DA15b1253463626243215495a3718beE4"
	if params.VRFOwnerAddress != "" && vrfVersion == vrfcommon.V2 {
		vrfOwnerAddress = params.VRFOwnerAddress
	}
	pollPeriod := 5 * time.Second
	if params.PollPeriod > 0 && (vrfVersion == vrfcommon.V2 || vrfVersion == vrfcommon.V2Plus) {
		pollPeriod = params.PollPeriod
	}
	batchFulfillmentGasMultiplier := 1.0
	if params.BatchFulfillmentGasMultiplier >= 1.0 {
		batchFulfillmentGasMultiplier = params.BatchFulfillmentGasMultiplier
	}
	confirmations := 6
	if params.MinIncomingConfirmations != 0 {
		confirmations = params.MinIncomingConfirmations
	}
	gasLanePrice := assets.GWei(100)
	if params.GasLanePrice != nil {
		gasLanePrice = params.GasLanePrice
	}
	requestTimeout := 24 * time.Hour
	if params.RequestTimeout != 0 {
		requestTimeout = params.RequestTimeout
	}
	publicKey := "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
	if params.PublicKey != "" {
		publicKey = params.PublicKey
	}
	chunkSize := 20
	if params.ChunkSize != 0 {
		chunkSize = params.ChunkSize
	}
	observationSource := fmt.Sprintf(`
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
submit_tx  [type=ethtx to="%s"
            data="$(encode_tx)"
            minConfirmations="0"
            from="$(jobSpec.from)"
            txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"
            transmitChecker="{\\"CheckerType\\": \\"vrf_v1\\", \\"VRFCoordinatorAddress\\": \\"%s\\"}"]
decode_log->vrf->encode_tx->submit_tx
`, coordinatorAddress, coordinatorAddress)
	if params.V2 {
		observationSource = fmt.Sprintf(`
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
              to="%s"
              multiplier="1.1"
              data="$(vrf.output)"
]
simulate [type=ethcall
          to="%s"
		  gas="$(estimate_gas)"
		  gasPrice="$(jobSpec.maxGasPrice)"
		  extractRevertReason=true
		  contract="%s"
		  data="$(vrf.output)"
]
decode_log->vrf->estimate_gas->simulate
`, coordinatorAddress, coordinatorAddress, coordinatorAddress)
	}
	if vrfVersion == vrfcommon.V2Plus {
		observationSource = fmt.Sprintf(`
decode_log              [type=ethabidecodelog
                         abi="RandomWordsRequested(bytes32 indexed keyHash,uint256 requestId,uint256 preSeed,uint256 indexed subId,uint16 minimumRequestConfirmations,uint32 callbackGasLimit,uint32 numWords,bytes extraArgs,address indexed sender)"
                         data="$(jobRun.logData)"
                         topics="$(jobRun.logTopics)"]
generate_proof          [type=vrfv2plus
                         publicKey="$(jobSpec.publicKey)"
                         requestBlockHash="$(jobRun.logBlockHash)"
                         requestBlockNumber="$(jobRun.logBlockNumber)"
                         topics="$(jobRun.logTopics)"]
estimate_gas            [type=estimategaslimit
                         to="%s"
                         multiplier="1.1"
                         data="$(generate_proof.output)"
						 block="latest"
]
simulate_fulfillment    [type=ethcall
                         to="%s"
		                 gas="$(estimate_gas)"
		                 gasPrice="$(jobSpec.maxGasPrice)"
		                 extractRevertReason=true
		                 contract="%s"
		                 data="$(generate_proof.output)"
						 block="latest"
]
decode_log->generate_proof->estimate_gas->simulate_fulfillment
`, coordinatorAddress, coordinatorAddress, coordinatorAddress)
	}
	if params.ObservationSource != "" {
		observationSource = params.ObservationSource
	}
	if params.EVMChainID == "" {
		params.EVMChainID = "0"
	}
	template := `
externalJobID = "%s"
type = "vrf"
schemaVersion = 1
name = "%s"
coordinatorAddress = "%s"
evmChainID         =  "%s"
batchCoordinatorAddress = "%s"
batchFulfillmentEnabled = %v
batchFulfillmentGasMultiplier = %s
customRevertsPipelineEnabled = %v
minIncomingConfirmations = %d
requestedConfsDelay = %d
requestTimeout = "%s"
publicKey = "%s"
chunkSize = %d
backoffInitialDelay = "%s"
backoffMaxDelay = "%s"
gasLanePrice = "%s"
pollPeriod = "%s"
observationSource = """
%s
"""
`
	toml := fmt.Sprintf(template,
		jobID, name, coordinatorAddress, params.EVMChainID, batchCoordinatorAddress,
		params.BatchFulfillmentEnabled, strconv.FormatFloat(batchFulfillmentGasMultiplier, 'f', 2, 64),
		params.CustomRevertsPipelineEnabled,
		confirmations, params.RequestedConfsDelay, requestTimeout.String(), publicKey, chunkSize,
		params.BackoffInitialDelay.String(), params.BackoffMaxDelay.String(), gasLanePrice.String(),
		pollPeriod.String(), observationSource)
	if len(params.FromAddresses) != 0 {
		var addresses []string
		for _, address := range params.FromAddresses {
			addresses = append(addresses, fmt.Sprintf("%q", address))
		}
		toml = toml + "\n" + fmt.Sprintf(`fromAddresses = [%s]`, strings.Join(addresses, ", "))
	}
	if vrfVersion == vrfcommon.V2 {
		toml = toml + "\n" + fmt.Sprintf(`vrfOwnerAddress = "%s"`, vrfOwnerAddress)
	}

	return VRFSpec{VRFSpecParams: VRFSpecParams{
		JobID:                    jobID,
		Name:                     name,
		CoordinatorAddress:       coordinatorAddress,
		BatchCoordinatorAddress:  batchCoordinatorAddress,
		BatchFulfillmentEnabled:  params.BatchFulfillmentEnabled,
		MinIncomingConfirmations: confirmations,
		PublicKey:                publicKey,
		ObservationSource:        observationSource,
		EVMChainID:               params.EVMChainID,
		RequestedConfsDelay:      params.RequestedConfsDelay,
		RequestTimeout:           requestTimeout,
		ChunkSize:                chunkSize,
		BackoffInitialDelay:      params.BackoffInitialDelay,
		BackoffMaxDelay:          params.BackoffMaxDelay,
		VRFOwnerAddress:          vrfOwnerAddress,
		VRFVersion:               vrfVersion,
		PollPeriod:               pollPeriod,
	}, toml: toml}
}

type OCRSpecParams struct {
	JobID              string
	Name               string
	TransmitterAddress string
	ContractAddress    string
	DS1BridgeName      string
	DS2BridgeName      string
	EVMChainID         string
}

type OCRSpec struct {
	OCRSpecParams
	toml string
}

func (os OCRSpec) Toml() string {
	return os.toml
}

func GenerateOCRSpec(params OCRSpecParams) OCRSpec {
	jobID := params.JobID
	if jobID == "" {
		jobID = uuid.New().String()
	}
	transmitterAddress := "0xF67D0290337bca0847005C7ffD1BC75BA9AAE6e4"
	if params.TransmitterAddress != "" {
		transmitterAddress = params.TransmitterAddress
	}
	contractAddress := "0x613a38AC1659769640aaE063C651F48E0250454C"
	if params.ContractAddress != "" {
		contractAddress = params.ContractAddress
	}
	name := params.Name
	if params.Name == "" {
		name = jobID
	}
	ds1BridgeName := fmt.Sprintf("automatically_generated_bridge_%s", uuid.New().String())
	if params.DS1BridgeName != "" {
		ds1BridgeName = params.DS1BridgeName
	}
	ds2BridgeName := fmt.Sprintf("automatically_generated_bridge_%s", uuid.New().String())
	if params.DS2BridgeName != "" {
		ds2BridgeName = params.DS2BridgeName
	}

	evmChainID := "0"
	if params.EVMChainID != "" {
		evmChainID = params.EVMChainID
	}
	template := `
type               = "offchainreporting"
schemaVersion      = 1
name               = "%s"
contractAddress    = "%s"
evmChainID         = %s
p2pPeerID          = "12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6X"
externalJobID      =  "%s"
p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
isBootstrapPeer    = false
keyBundleID        = "f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5"
monitoringEndpoint = "chain.link:4321"
transmitterAddress = "%s"
observationTimeout = "10s"
blockchainTimeout  = "20s"
contractConfigTrackerSubscribeInterval = "2m"
contractConfigTrackerPollInterval = "1m"
contractConfigConfirmations = 3
observationSource = """
    // data source 1
    ds1          [type=bridge name="%s"];
    ds1_parse    [type=jsonparse path="one,two"];
    ds1_multiply [type=multiply times=1.23];

    // data source 2
    ds2          [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData="{\\"hi\\": \\"hello\\"}"];
    ds2_parse    [type=jsonparse path="three,four"];
    ds2_multiply [type=multiply times=4.56];

    ds1 -> ds1_parse -> ds1_multiply -> answer1;
    ds2 -> ds2_parse -> ds2_multiply -> answer1;

    answer1 [type=median                      index=0];
    answer2 [type=bridge name="%s" index=1];
"""
`
	return OCRSpec{OCRSpecParams: OCRSpecParams{
		JobID:              jobID,
		Name:               name,
		TransmitterAddress: transmitterAddress,
		DS1BridgeName:      ds1BridgeName,
		DS2BridgeName:      ds2BridgeName,
	}, toml: fmt.Sprintf(template, name, contractAddress, evmChainID, jobID, transmitterAddress, ds1BridgeName, ds2BridgeName)}
}

type WebhookSpecParams struct {
	ExternalInitiators []webhook.TOMLWebhookSpecExternalInitiator
}

type WebhookSpec struct {
	WebhookSpecParams
	toml string
}

func (ws WebhookSpec) Toml() string {
	return ws.toml
}

func GenerateWebhookSpec(params WebhookSpecParams) (ws WebhookSpec) {
	var externalInitiatorsTOMLs []string
	for _, wsEI := range params.ExternalInitiators {
		s := fmt.Sprintf(`{ name = "%s", spec = '%s' }`, wsEI.Name, wsEI.Spec)
		externalInitiatorsTOMLs = append(externalInitiatorsTOMLs, s)
	}
	externalInitiatorsTOML := strings.Join(externalInitiatorsTOMLs, ",\n")
	template := `
type            = "webhook"
schemaVersion   = 1
externalInitiators = [
    %s
]
observationSource   = """
ds          [type=http method=GET url="https://chain.link/ETH-USD"];
ds_parse    [type=jsonparse path="data,price"];
ds_multiply [type=multiply times=100];
ds -> ds_parse -> ds_multiply;
"""
`
	ws.toml = fmt.Sprintf(template, externalInitiatorsTOML)
	ws.WebhookSpecParams = params

	return ws
}

// BlockhashStoreSpecParams defines params for building a blockhash store job spec.
type BlockhashStoreSpecParams struct {
	JobID                          string
	Name                           string
	CoordinatorV1Address           string
	CoordinatorV2Address           string
	CoordinatorV2PlusAddress       string
	WaitBlocks                     int
	HeartbeatPeriod                time.Duration
	LookbackBlocks                 int
	BlockhashStoreAddress          string
	TrustedBlockhashStoreAddress   string
	TrustedBlockhashStoreBatchSize int32
	PollPeriod                     time.Duration
	RunTimeout                     time.Duration
	EVMChainID                     int64
	FromAddresses                  []string
}

// BlockhashStoreSpec defines a blockhash store job spec.
type BlockhashStoreSpec struct {
	BlockhashStoreSpecParams
	toml string
}

// Toml returns the BlockhashStoreSpec in TOML string form.
func (bhs BlockhashStoreSpec) Toml() string {
	return bhs.toml
}

// GenerateBlockhashStoreSpec creates a BlockhashStoreSpec from the given params.
func GenerateBlockhashStoreSpec(params BlockhashStoreSpecParams) BlockhashStoreSpec {
	if params.JobID == "" {
		params.JobID = "123e4567-e89b-12d3-a456-426655442222"
	}

	if params.Name == "" {
		params.Name = "blockhash-store"
	}

	if params.CoordinatorV1Address == "" {
		params.CoordinatorV1Address = "0x19D20b4Ec0424A530C3C1cDe874445E37747eb18"
	}

	if params.CoordinatorV2Address == "" {
		params.CoordinatorV2Address = "0x2498e651Ae17C2d98417C4826F0816Ac6366A95E"
	}

	if params.CoordinatorV2PlusAddress == "" {
		params.CoordinatorV2PlusAddress = "0x92B5e28Ac583812874e4271380c7d070C5FB6E6b"
	}

	if params.TrustedBlockhashStoreAddress == "" {
		params.TrustedBlockhashStoreAddress = utils.ZeroAddress.Hex()
	}

	if params.TrustedBlockhashStoreBatchSize == 0 {
		params.TrustedBlockhashStoreBatchSize = 20
	}

	if params.WaitBlocks == 0 {
		params.WaitBlocks = 100
	}

	if params.LookbackBlocks == 0 {
		params.LookbackBlocks = 200
	}

	if params.BlockhashStoreAddress == "" {
		params.BlockhashStoreAddress = "0x31Ca8bf590360B3198749f852D5c516c642846F6"
	}

	if params.PollPeriod == 0 {
		params.PollPeriod = 30 * time.Second
	}

	if params.RunTimeout == 0 {
		params.RunTimeout = 15 * time.Second
	}

	var formattedFromAddresses string
	if params.FromAddresses == nil {
		formattedFromAddresses = `["0x4bd43cb108Bc3742e484f47E69EBfa378cb6278B"]`
	} else {
		var addresses []string
		for _, address := range params.FromAddresses {
			addresses = append(addresses, fmt.Sprintf("%q", address))
		}
		formattedFromAddresses = fmt.Sprintf("[%s]", strings.Join(addresses, ", "))
	}

	template := `
type = "blockhashstore"
schemaVersion = 1
name = "%s"
coordinatorV1Address = "%s"
coordinatorV2Address = "%s"
coordinatorV2PlusAddress = "%s"
waitBlocks = %d
lookbackBlocks = %d
blockhashStoreAddress = "%s"
trustedBlockhashStoreAddress = "%s"
trustedBlockhashStoreBatchSize = %d
pollPeriod = "%s"
runTimeout = "%s"
evmChainID = "%d"
fromAddresses = %s
heartbeatPeriod = "%s"
`
	toml := fmt.Sprintf(template, params.Name, params.CoordinatorV1Address,
		params.CoordinatorV2Address, params.CoordinatorV2PlusAddress, params.WaitBlocks, params.LookbackBlocks,
		params.BlockhashStoreAddress, params.TrustedBlockhashStoreAddress, params.TrustedBlockhashStoreBatchSize, params.PollPeriod.String(), params.RunTimeout.String(),
		params.EVMChainID, formattedFromAddresses, params.HeartbeatPeriod.String())

	return BlockhashStoreSpec{BlockhashStoreSpecParams: params, toml: toml}
}

// BlockHeaderFeederSpecParams defines params for building a block header feeder job spec.
type BlockHeaderFeederSpecParams struct {
	JobID                      string
	Name                       string
	CoordinatorV1Address       string
	CoordinatorV2Address       string
	CoordinatorV2PlusAddress   string
	WaitBlocks                 int
	LookbackBlocks             int
	BlockhashStoreAddress      string
	BatchBlockhashStoreAddress string
	PollPeriod                 time.Duration
	RunTimeout                 time.Duration
	EVMChainID                 int64
	FromAddresses              []string
	GetBlockhashesBatchSize    uint16
	StoreBlockhashesBatchSize  uint16
}

// BlockHeaderFeederSpec defines a block header feeder job spec.
type BlockHeaderFeederSpec struct {
	BlockHeaderFeederSpecParams
	toml string
}

// Toml returns the BlockhashStoreSpec in TOML string form.
func (b BlockHeaderFeederSpec) Toml() string {
	return b.toml
}

// GenerateBlockHeaderFeederSpec creates a BlockHeaderFeederSpec from the given params.
func GenerateBlockHeaderFeederSpec(params BlockHeaderFeederSpecParams) BlockHeaderFeederSpec {
	if params.JobID == "" {
		params.JobID = "123e4567-e89b-12d3-a456-426655442211"
	}

	if params.Name == "" {
		params.Name = "blockheaderfeeder"
	}

	if params.CoordinatorV1Address == "" {
		params.CoordinatorV1Address = "0x2d7F888fE0dD469bd81A12f77e6291508f714d4B"
	}

	if params.CoordinatorV2Address == "" {
		params.CoordinatorV2Address = "0x2d7F888fE0dD469bd81A12f77e6291508f714d4B"
	}

	if params.CoordinatorV2PlusAddress == "" {
		params.CoordinatorV2PlusAddress = "0x2d7F888fE0dD469bd81A12f77e6291508f714d4B"
	}

	if params.WaitBlocks == 0 {
		params.WaitBlocks = 256
	}

	if params.LookbackBlocks == 0 {
		params.LookbackBlocks = 500
	}

	if params.BlockhashStoreAddress == "" {
		params.BlockhashStoreAddress = "0x016D54091ee83D42aF46e4F2d7177D0A232D2bDa"
	}

	if params.BatchBlockhashStoreAddress == "" {
		params.BatchBlockhashStoreAddress = "0xde08B57586839BfF5DB58Bdd7FdeB7142Bff3795"
	}

	if params.PollPeriod == 0 {
		params.PollPeriod = 60 * time.Second
	}

	if params.RunTimeout == 0 {
		params.RunTimeout = 30 * time.Second
	}

	if params.GetBlockhashesBatchSize == 0 {
		params.GetBlockhashesBatchSize = 10
	}

	if params.StoreBlockhashesBatchSize == 0 {
		params.StoreBlockhashesBatchSize = 5
	}

	var formattedFromAddresses string
	if params.FromAddresses == nil {
		formattedFromAddresses = `["0xBe0b739f841bC113D4F4e4CdD16086ffAbB5f39f"]`
	} else {
		var addresses []string
		for _, address := range params.FromAddresses {
			addresses = append(addresses, fmt.Sprintf("%q", address))
		}
		formattedFromAddresses = fmt.Sprintf("[%s]", strings.Join(addresses, ", "))
	}

	template := `
type = "blockheaderfeeder"
schemaVersion = 1
name = "%s"
coordinatorV1Address = "%s"
coordinatorV2Address = "%s"
coordinatorV2PlusAddress = "%s"
waitBlocks = %d
lookbackBlocks = %d
blockhashStoreAddress = "%s"
batchBlockhashStoreAddress = "%s"
pollPeriod = "%s"
runTimeout = "%s"
evmChainID = "%d"
fromAddresses = %s
getBlockhashesBatchSize = %d
storeBlockhashesBatchSize = %d
`
	toml := fmt.Sprintf(template, params.Name, params.CoordinatorV1Address,
		params.CoordinatorV2Address, params.CoordinatorV2PlusAddress, params.WaitBlocks, params.LookbackBlocks,
		params.BlockhashStoreAddress, params.BatchBlockhashStoreAddress, params.PollPeriod.String(),
		params.RunTimeout.String(), params.EVMChainID, formattedFromAddresses, params.GetBlockhashesBatchSize,
		params.StoreBlockhashesBatchSize)

	return BlockHeaderFeederSpec{BlockHeaderFeederSpecParams: params, toml: toml}
}

type StreamSpecParams struct {
	Name     string
	StreamID uint64
}

type StreamSpec struct {
	StreamSpecParams
	toml string
}

// Toml returns the BlockhashStoreSpec in TOML string form.
func (b StreamSpec) Toml() string {
	return b.toml
}

func GenerateStreamSpec(params StreamSpecParams) StreamSpec {
	template := `
type = "stream"
schemaVersion = 1
name = "%s"
streamID = %d
observationSource = """
ds          [type=http method=GET url="https://chain.link/ETH-USD"];
ds_parse    [type=jsonparse path="data,price"];
ds_multiply [type=multiply times=100];
ds -> ds_parse -> ds_multiply;
"""
`

	toml := fmt.Sprintf(template, params.Name, params.StreamID)
	return StreamSpec{StreamSpecParams: params, toml: toml}
}

type WorkflowSpec struct {
	toml string
}

func (w WorkflowSpec) Toml() string {
	return w.toml
}

func GenerateWorkflowSpec(id, owner, name, spec string) WorkflowSpec {
	template := `
type = "workflow"
schemaVersion = 1
name = "test-spec"
workflowId = "%s"
workflowOwner = "%s"
workflowName = "%s"
workflow = """
%s
"""
`
	toml := fmt.Sprintf(template, id, owner, name, spec)
	return WorkflowSpec{toml: toml}
}
