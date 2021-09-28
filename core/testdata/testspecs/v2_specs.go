package testspecs

import (
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink/core/services/webhook"
)

var (
	CronSpec = `
type                = "cron"
schemaVersion       = 1
schedule            = "CRON_TZ=UTC * 0 0 1 1 *"
externalJobID       =  "123e4567-e89b-12d3-a456-426655440003"
observationSource   = """
ds          [type=http method=GET url="https://chain.link/ETH-USD"];
ds_parse    [type=jsonparse path="data,price"];
ds_multiply [type=multiply times=100];
ds -> ds_parse -> ds_multiply;
"""
`
	DirectRequestSpec = `
type                = "directrequest"
schemaVersion       = 1
name                = "example eth request event spec"
contractAddress     = "0x613a38AC1659769640aaE063C651F48E0250454C"
externalJobID       =  "123e4567-e89b-12d3-a456-426655440004"
observationSource   = """
    ds1          [type=http method=GET url="http://example.com" allowunrestrictednetworkaccess="true"];
    ds1_parse    [type=jsonparse path="USD"];
    ds1_multiply [type=multiply times=100];
    ds1 -> ds1_parse -> ds1_multiply;
"""
`
	DirectRequestSpecWithRequestersAndMinContractPayment = `
type                         = "directrequest"
schemaVersion                = 1
requesters                   = ["0xaaaa1F8ee20f5565510B84f9353F1E333E753B7a", "0xbbbb70F0e81C6F3430dfdC9fa02fB22BdD818C4e"]
minContractPaymentLinkJuels  = "1000000000000000000000"
name                         = "example eth request event spec with requesters and min contract payment"
contractAddress              = "0x613a38AC1659769640aaE063C651F48E0250454C"
externalJobID                = "123e4567-e89b-12d3-a456-426655440014"
observationSource            = """
    ds1          [type=http method=GET url="http://example.com" allowunrestrictednetworkaccess="true"];
    ds1_parse    [type=jsonparse path="USD"];
    ds1_multiply [type=multiply times=100];
    ds1 -> ds1_parse -> ds1_multiply;
"""
`
	FluxMonitorSpec = `
type                = "fluxmonitor"
schemaVersion       = 1
name                = "example flux monitor spec"
contractAddress     = "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
externalJobID       =  "123e4567-e89b-12d3-a456-426655440005"
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
)

type KeeperSpecParams struct {
	ContractAddress string
	FromAddress     string
	EvmChainID      int
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
type            = "keeper"
schemaVersion   = 2
name            = "example keeper spec"
contractAddress = "%s"
fromAddress     = "%s"
evmChainID      = %d
externalJobID   =  "123e4567-e89b-12d3-a456-426655440002"


observationSource = """
encode_check_upkeep_tx   [type=ethabiencode
                          abi="checkUpkeep(uint256 id, address from)"
                          data="{\\"id\\":$(jobSpec.upkeepID),\\"from\\":$(jobSpec.fromAddress)}"]
check_upkeep_tx          [type=ethcall
                          failEarly=true
                          extractRevertReason=true
                          contract="$(jobSpec.contractAddress)"
                          gas="$(jobSpec.checkUpkeepGasLimit)"
                          gasPrice="$(jobSpec.gasPrice)"
                          data="$(encode_check_upkeep_tx)"]
decode_check_upkeep_tx   [type=ethabidecode
                          abi="bytes memory performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth"]
encode_perform_upkeep_tx [type=ethabiencode
                          abi="performUpkeep(uint256 id, bytes calldata performData)"
                          data="{\\"id\\": $(jobSpec.upkeepID),\\"performData\\":$(decode_check_upkeep_tx.performData)}"]
perform_upkeep_tx        [type=ethtx
                          minConfirmations=0
                          to="$(jobSpec.contractAddress)"
                          data="$(encode_perform_upkeep_tx)"
                          gasLimit="$(jobSpec.performUpkeepGasLimit)"
                          txMeta="{\\"jobID\\":$(jobSpec.jobID)}"]
encode_check_upkeep_tx -> check_upkeep_tx -> decode_check_upkeep_tx -> encode_perform_upkeep_tx -> perform_upkeep_tx
"""
`
	return KeeperSpec{
		KeeperSpecParams: params,
		toml:             fmt.Sprintf(template, params.ContractAddress, params.FromAddress, params.EvmChainID),
	}
}

type VRFSpecParams struct {
	JobID              string
	Name               string
	CoordinatorAddress string
	Confirmations      int
	PublicKey          string
	ObservationSource  string
	V2                 bool
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
	coordinatorAddress := "0xABA5eDc1a551E55b1A570c0e1f1055e5BE11eca7"
	if params.CoordinatorAddress != "" {
		coordinatorAddress = params.CoordinatorAddress
	}
	confirmations := 6
	if params.Confirmations != 0 {
		confirmations = params.Confirmations
	}
	publicKey := "0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F8179800"
	if params.PublicKey != "" {
		publicKey = params.PublicKey
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
            txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(decode_log.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
decode_log->vrf->encode_tx->submit_tx
`, coordinatorAddress)
	if params.V2 {
		//encode_tx    [type=ethabiencode
		//abi="fulfillRandomWords(bytes proof, bytes requestCommitment)"
		//data=<{"proof": $(vrf.proof), "requestCommitment": $(vrf.requestCommitment)}>]
		observationSource = fmt.Sprintf(`
decode_log   [type=ethabidecodelog
              abi="RandomWordsRequested(bytes32 indexed keyHash,uint256 requestId,uint256 preSeed,uint64 subId,uint16 minimumRequestConfirmations,uint32 callbackGasLimit,uint32 numWords,address indexed sender)"
              data="$(jobRun.logData)"
              topics="$(jobRun.logTopics)"]
vrf          [type=vrfv2
              publicKey="$(jobSpec.publicKey)"
              requestBlockHash="$(jobRun.logBlockHash)"
              requestBlockNumber="$(jobRun.logBlockNumber)"
              topics="$(jobRun.logTopics)"]
estimate_gas [type=estimategaslimit
              to="%s"
              multiplier="1"
              data="$(vrf.output)"]
submit_tx  [type=ethtx to="%s"
            data="$(vrf.output)"
            gasLimit="$(estimate_gas)"
            minConfirmations="0"
            txMeta="{\\"requestTxHash\\": $(jobRun.logTxHash),\\"requestID\\": $(vrf.requestID),\\"jobID\\": $(jobSpec.databaseID)}"]
decode_log->vrf->estimate_gas->submit_tx
`, coordinatorAddress, coordinatorAddress)
	}
	if params.ObservationSource != "" {
		publicKey = params.ObservationSource
	}
	template := `
externalJobID = "%s"
type = "vrf"
schemaVersion = 1
name = "%s"
coordinatorAddress = "%s"
confirmations = %d
publicKey = "%s"
observationSource = """
%s
"""
`
	return VRFSpec{VRFSpecParams: VRFSpecParams{
		JobID:              jobID,
		Name:               name,
		CoordinatorAddress: coordinatorAddress,
		Confirmations:      confirmations,
		PublicKey:          publicKey,
		ObservationSource:  observationSource,
	}, toml: fmt.Sprintf(template, jobID, name, coordinatorAddress, confirmations, publicKey, observationSource)}
}

type OCRSpecParams struct {
	JobID              string
	Name               string
	TransmitterAddress string
}

type OCRSpec struct {
	OCRSpecParams
	toml string
}

func (os OCRSpec) Toml() string {
	return os.toml
}

func GenerateOCRSpec(params OCRSpecParams) OCRSpec {
	jobID := "123e4567-e89b-12d3-a456-426655440001"
	if params.JobID != "" {
		jobID = params.JobID
	}
	transmitterAddress := "0xF67D0290337bca0847005C7ffD1BC75BA9AAE6e4"
	if params.TransmitterAddress != "" {
		transmitterAddress = params.TransmitterAddress
	}
	name := "web oracle spec"
	if params.Name != "" {
		name = params.Name
	}
	template := `
type               = "offchainreporting"
schemaVersion      = 1
name               = "%s"
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6X"
externalJobID     =  "%s"
p2pBootstrapPeers  = [
    "/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
]
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
    ds1          [type=bridge name=voter_turnout];
    ds1_parse    [type=jsonparse path="one,two"];
    ds1_multiply [type=multiply times=1.23];

    // data source 2
    ds2          [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData="{\\"hi\\": \\"hello\\"}"];
    ds2_parse    [type=jsonparse path="three,four"];
    ds2_multiply [type=multiply times=4.56];

    ds1 -> ds1_parse -> ds1_multiply -> answer1;
    ds2 -> ds2_parse -> ds2_multiply -> answer1;

    answer1 [type=median                      index=0];
    answer2 [type=bridge name=election_winner index=1];
"""
`
	return OCRSpec{OCRSpecParams: OCRSpecParams{
		JobID:              jobID,
		Name:               name,
		TransmitterAddress: transmitterAddress,
	}, toml: fmt.Sprintf(template, name, jobID, transmitterAddress)}
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
