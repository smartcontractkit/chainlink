package jobs

var (
	VRFV2JobFormatted = `type = "vrf"
name = "vrf_v2"
schemaVersion = 1
coordinatorAddress = "%s"
batchCoordinatorAddress = "%s"
batchFulfillmentEnabled = %t
batchFulfillmentGasMultiplier = 1.1
publicKey = "%s"
minIncomingConfirmations = %d
evmChainID = "%d"
fromAddresses = ["%s"]
pollPeriod = "300ms"
requestTimeout = "30m0s"
observationSource = """decode_log   [type=ethabidecodelog
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
              data="$(vrf.output)"]
simulate     [type=ethcall
              from="%s"
              to="%s"
              gas="$(estimate_gas)"
              gasPrice="$(jobSpec.maxGasPrice)"
              extractRevertReason=true
              contract="%s"
              data="$(vrf.output)"]
decode_log->vrf->estimate_gas->simulate
"""`

	VRFV2PlusJobFormatted = `
type = "vrf"
name = "vrf_v2_plus"
schemaVersion = 1
coordinatorAddress = "%s"
batchCoordinatorAddress = "%s"
batchFulfillmentEnabled = %t
batchFulfillmentGasMultiplier = 1.1
publicKey = "%s"
minIncomingConfirmations = %d
evmChainID = "%d"
fromAddresses = ["%s"]
pollPeriod = "300ms"
requestTimeout = "30m0s"
observationSource = """
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
						 data="$(generate_proof.output)"]
simulate_fulfillment    [type=ethcall
						 from="%s"
                         to="%s"
		                 gas="$(estimate_gas)"
		                 gasPrice="$(jobSpec.maxGasPrice)"
		                 extractRevertReason=true
		                 contract="%s"
		                 data="$(generate_proof.output)"]
decode_log->generate_proof->estimate_gas->simulate_fulfillment
"""
`

	BHSJobFormatted = `type = "blockhashstore"
schemaVersion = 1
name = "blockhashstore"
forwardingAllowed = false
coordinatorV2Address = "%s"
waitBlocks = %d
lookbackBlocks = %d
blockhashStoreAddress = "%s"
pollPeriod = "30s"
runTimeout = "1m0s"
evmChainID = "%d"
fromAddresses = ["%s"]
`
	BHFJobFormatted = `type = "blockheaderfeeder"
schemaVersion = 1
name = "blockheaderfeeder"
forwardingAllowed = false
coordinatorV2Address = "%s"
waitBlocks = 256
lookbackBlocks = 1_000
blockhashStoreAddress = "%s"
batchBlockhashStoreAddress = "%s"
pollPeriod = "10s"
runTimeout = "30s"
evmChainID = "%d"
fromAddresses = ["%s"]
getBlockhashesBatchSize = 50
storeBlockhashesBatchSize = 10
`
)
