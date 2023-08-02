package jobs

var (
	VRFJobFormatted = `type = "vrf"
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
              to="%s"
              gas="$(estimate_gas)"
              gasPrice="$(jobSpec.maxGasPrice)"
              extractRevertReason=true
              contract="%s"
              data="$(vrf.output)"]
decode_log->vrf->estimate_gas->simulate
"""`
)
