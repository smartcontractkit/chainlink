# Mercury Documentation

## Useful Links

[Configuration Builder](https://github.com/smartcontractkit/the-most-amazing-mercury-contract-configuration-tool)

[Contracts](https://github.com/smartcontractkit/chainlink/contracts/src/v0.8/llo-feeds)

[OCR3 Config Documentation](https://github.com/smartcontractkit/libocr/blob/master/offchainreporting2plus/internal/config/ocr3config/public_config.go)




### Example Feed Configuration

```json
{
  "feedId": "0x14e044f932bb959cc2aa8dc1ba110c09224e639aae00264c1ffc2a0830904a3c",
  "chainId": 42161, // source chain id
  "contractAddress": "0x14e044f932bb959cc2aa8dc1ba110c09224e639a", // verifier contract address
  "configCount": 1, // the index of this config
  "signers": [
    "0x000....01",
    "0x000....02",
    "0x000....03",
    "0x000....04"
    ], // NOP signing addresses,
    "transmitters": [
    "0x000....11",
    "0x000....12",
    "0x000....13",
    "0x000....14"
    ], // NOP transmitter addresses
    "offchainConfig": {
      "baseUSDFee": "0.1", // 10c base fee to verify the report
      "deltaCertifiedCommitRequest": "1s",
      "deltaGrace": "0s",
      "deltaInitial": "600ms",
      "deltaProgress": "2s",
      "deltaResend": "10s",
      "deltaRound": "250ms",
      "deltaStage": "0s",
      "expirationWindow": "86400", //window in in which a report can be verified in seconds
      "f": 3,
      "maxDurationObservation": "250ms",
      "maxDurationQuery": "50ms",
      "maxDurationShouldAcceptAttestedReport": "50ms",
      "maxDurationShouldTransmitAcceptedReport": "50ms",
      "rMax": "25",
      "s": [
        4
      ]
    },
  "offchainConfigVersion": 30,
  "onchainConfig": {
    "max": "99999999999999999999999999999",
    "min": "1"
  }
}
```

## Jobs

### Bootstrap

**🚨 Important config**

`relayConfig.chainID` - target chain id. (the chain we pull block numbers from)

`contractID` - the contract address of the verifier contract.

<details><summary>Example bootstrap TOML</summary>

```toml
type = "bootstrap"
relay = "evm"
schemaVersion = 1
name = "$feed_name"
contractID = "$verifier_contract_address"
feedID = "$feed_id" # IMPORTANT - DON'T FORGET THIS OR IT WON'T WORK
contractConfigTrackerPollInterval = "15s"

[relayConfig]
chainID = $evm_chain_id
fromBlock = $from_block
```
</details>

### OCR2

<details><summary>Example OCR2 Mercury TOML</summary>

```toml
type = "offchainreporting2"
schemaVersion = 1
name = "$feed_name"
forwardingAllowed = false
maxTaskDuration = "1s"
contractID = "$verifier_contract_address"
feedID = "$feed_id"
contractConfigTrackerPollInterval = "15s"
ocrKeyBundleID = "$key_bundle_id"
p2pv2Bootstrappers = [
  "$bootstrapper_address>"
]
relay = "evm"
pluginType = "mercury"
transmitterID = "$csa_public_key"

observationSource = """
  // ncfx
	ds1_payload          [type=bridge name="ncfx" timeout="50ms" requestData="{\\"data\\":{\\"endpoint\\":\\"crypto-lwba\\",\\"from\\":\\"ETH\\",\\"to\\":\\"USD\\"}}"];
  ds1_median           [type=jsonparse path="data,mid"];
  ds1_bid              [type=jsonparse path="data,bid"];
  ds1_ask              [type=jsonparse path="data,ask"];
  
  ds1_median_multiply  [type=multiply times=100000000];
  ds1_bid_multiply     [type=multiply times=100000000];
  ds1_ask_multiply     [type=multiply times=100000000];

  // tiingo
  ds2_payload          [type=bridge name="tiingo" timeout="50ms" requestData="{\\"data\\":{\\"endpoint\\":\\"crypto-lwba\\",\\"from\\":\\"ETH\\",\\"to\\":\\"USD\\"}}"];
  ds2_median           [type=jsonparse path="data,mid"];
  ds2_bid              [type=jsonparse path="data,bid"];
  ds2_ask              [type=jsonparse path="data,ask"];

  ds2_median_multiply  [type=multiply times=100000000];
  ds2_bid_multiply     [type=multiply times=100000000];
  ds2_ask_multiply     [type=multiply times=100000000];

  // coinmetrics
  ds3_payload          [type=bridge name="coinmetrics" timeout="50ms" requestData="{\\"data\\":{\\"endpoint\\":\\"crypto-lwba\\",\\"from\\":\\"ETH\\",\\"to\\":\\"USD\\"}}"];
  ds3_median           [type=jsonparse path="data,mid"];
  ds3_bid              [type=jsonparse path="data,bid"];
  ds3_ask              [type=jsonparse path="data,ask"];

  ds3_median_multiply  [type=multiply times=100000000];
  ds3_bid_multiply     [type=multiply times=100000000];
  ds3_ask_multiply     [type=multiply times=100000000];

  ds1_payload -> ds1_median -> ds1_median_multiply -> benchmark_price;
  ds2_payload -> ds2_median -> ds2_median_multiply -> benchmark_price;
  ds3_payload -> ds3_median -> ds3_median_multiply -> benchmark_price;

  benchmark_price [type=median allowedFaults=2 index=0];

  ds1_payload -> ds1_bid -> ds1_bid_multiply -> bid_price;
  ds2_payload -> ds2_bid -> ds2_bid_multiply -> bid_price;
  ds3_payload -> ds3_bid -> ds3_bid_multiply -> bid_price;

  bid_price [type=median allowedFaults=2 index=1];

  ds1_payload -> ds1_ask -> ds1_ask_multiply -> ask_price;
  ds2_payload -> ds2_ask -> ds2_ask_multiply -> ask_price;
  ds3_payload -> ds3_ask -> ds3_ask_multiply -> ask_price;

  ask_price [type=median allowedFaults=2 index=2];
"""

[pluginConfig]
serverURL = "$mercury_server_url"
serverPubKey = "$mercury_server_public_key"

[relayConfig]
chainID = $evm_chain_id
fromBlock = $from_block
```
</details>

## Nodes

**🚨 Important config**

`OCR2.Enabled` - must be `true` - Mercury uses OCR2.

`P2P.V2.Enabled` - required in order for OCR2 to work.

`Feature.LogPoller` - required in order for OCR2 to work. You will get fatal errors if not set.

`JobPipeline.MaxSuccessfulRuns` - set to `0` to disable saving pipeline runs to reduce load on the db. Obviously this means you won’t see anything in the UI.

`TelemetryIngress.SendInterval`  - How frequently to send telemetry batches. Mercury generates a lot of telemetry data due to the throughput. `100ms` has been tested for a single feed with 5 nodes - this will need to be monitored (along with relevant config) as we add more feeds to a node.

`Database` - **must** increase connection limits above the standard defaults

<details><summary>Example node config TOML</summary>

```toml
RootDir = '$ROOT_DIR'

[JobPipeline]
MaxSuccessfulRuns = 0 # you may set to some small value like '10' or similar if you like looking at job runs in the UI

[Feature]
UICSAKeys = true # required
LogPoller = true # required

[Log]
Level = 'info' # this should be 'debug' for chainlink internal deployments, nops may use 'info' to reduce log volume

[Log.File]
< standard values >

[WebServer]
< standard values >

[WebServer.TLS]
< standard values >

[[EVM]]
ChainID = '42161' # change as needed based on target chain

[OCR]
Enabled = false # turn off OCR 1

[P2P]
TraceLogging = false # this should be 'true' for chainlink internal deployments, we may ask nops to set this to true for debugging
PeerID = '$PEERID'

[P2P.V2]
Enabled = true # required
DefaultBootstrappers = < mercury bootstrap nodes > # Note that this should ideally be set in the job spec, this is just a fallback
# Make sure these IPs are properly configured in the firewall. May not be necessary for internal nodes
AnnounceAddresses = ['$EXTERNAL_IP:$EXTERNAL_PORT'] # Use whichever port you like, pls randomize, MAKE SURE ITS CONFIGURED IN THE FIREWALL
ListenAddresses = ['0.0.0.0:$INTERNAL_PORT'] # Use whichever port you like, pls randomize, MAKE SURE ITS CONFIGURED IN THE FIREWALL

[OCR2]
Enabled = true # required
KeyBundleID = '$KEY_BUNDLE_ID' # Note that this should ideally be set in the job spec, this is just a fallback
CaptureEATelemetry = true

[TelemetryIngress]
UniConn = false
SendInterval = '250ms'
BufferSize = 300
MaxBatchSize = 100

[[TelemetryIngress.Endpoints]]
Network = 'EVM'
ChainID = '42161' # change as needed based on target chain
URL = '$TELEMETRY_ENDPOINT_URL' # Provided by Chainlink Labs RSTP team
ServerPubKey = '$TELEMETRY_PUB_KEY' # Provided by Chainlink Labs RSTP team

[Database]
MaxIdleConns = 100 # should equal or greater than total number of mercury jobs
MaxOpenConns = 400 # caution! ensure postgres is configured to support this

[[EVM.Nodes]]
< put RPC nodes here >
```
</details>
