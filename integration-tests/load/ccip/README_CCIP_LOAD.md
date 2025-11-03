# CCIP Load Test Runbook Staging Deployment

## Prerequisites

- CRIB and Chainlink repos in the same parent directory
- CRIB repo on `main` branch
- Local PostgreSQL database running (Required for Sui CR)

## CRIB Setup

### Directory Structure

Create the deployment data directory:
```bash
cd <CRIB_REPO>
mkdir -p deployments/ccip-v2/.tmp/data
```

### chain-details.json

Create `deployments/ccip-v2/.tmp/data/chain-details.json` with RPC endpoints:

```json
[
  {
    "ChainID": "11155420",
    "ChainName": "evm-opt-sepolia",
    "ChainType": "EVM",
    "WSRPCs": [{"External": "wss://...", "Internal": "wss://..."}],
    "HTTPRPCs": [{"External": "https://...", "Internal": "https://..."}]
  },
  {
    "ChainID": "2",
    "ChainName": "aptos-testnet",
    "ChainType": "APTOS",
    "WSRPCs": [{"External": "", "Internal": ""}],
    "HTTPRPCs": [{"External": "https://...", "Internal": "https://..."}]
  }
]
```

**Note:** Number of RPCs in this file determines the number of destination chains for the load test.

### address-book.json

Create `deployments/ccip-v2/.tmp/data/address-book.json` with deployment addresses for your target environment.

## Chainlink Repo Setup

### Environment Variables

Export required database and observability configuration:

```bash
export CL_DATABASE_URL=postgresql://postgres:postgres@localhost:5432/chainlink_test?sslmode=disable
export LOKI_URL=
export LOKI_TENANT_ID=
```

**Note:** Private keys are configured in `overrides.toml` under `[Load.CCIP.Load.TestnetConfig]`, not as environment variables.


### Test Configuration

Create `~/.testsecrets`:

```bash
E2E_TEST_COMMON_PRIVATE_KEY=" "
E2E_TEST_SOLANA_SECRET="<solana_secret>"
E2E_TEST_CHAINLINK_IMAGE="<ecr_registry>.dkr.ecr.us-west-2.amazonaws.com/chainlink"
E2E_TEST_CHAINLINK_VERSION=develop
CHAINLINK_VERSION=develop
```

### Load Test Configuration

Edit `integration-tests/testconfig/overrides.toml` - focus on the `[Load.CCIP.Load]` section:

```toml
[Load.CCIP.Load]
Testnet = true
ChaosMode = 0  # 0=none, 1=rpc latency, 2=full chaos
RPCLatency = "400ms"
RPCJitter = "0ms"

MessageDetails = [
    {Ratio = 50,  MsgType = "TokenTransfer"},
    {Ratio = 50,  MsgType = "Messaging", DataLengthBytes = 10, DestGasLimit=2_500_000},
]

SolanaDataSize = 50
RequestFrequency = "5s"  # Per destination chain
LoadDuration = "4h"
TestLabel = "your-test-label"
NumDestinationChains = 4  # Must match chain-details.json RPCs count
CribEnvDirectory = "../../../../crib/deployments/ccip-v2/.tmp"
TimeoutDuration = "6h"  # Set 1h+ longer than LoadDuration
OOOExecution = false

[Load.CCIP.Load.TestnetConfig]
Testnet = false
EVMPrivateKey = ""
AptosPrivateKey = ""
SolanaPrivateKey = ""
FundingAmountEth = 20000000000000000 # 0.02 ETH
FundingAmountSol = 0
FundingAmountApt = 0

[Load.CCIP.Load.TestnetConfig.SuiConfig]
SuiPrivateKey = ""
SuiFeeTokenObjectId = ""
SuiTestReceiverAddress = ""
SuiStateReceiverStateObjectId = ""
```

## Configuration Parameters

### MessageDetails
- **MsgType**: `TokenTransfer`, `Messaging`, or `ProgrammableTokenTransfer`
- **Ratio**: Must sum to 100 across all message types
- **DataLengthBytes**: Required for `Messaging` and `ProgrammableTokenTransfer`
- **DestGasLimit**: Optional, defaults per chain type

### RequestFrequency
Rate per destination chain:
- High: `5s`
- Medium: `10s`
- Soft: `30s`

Example: `5s` with 4 chains = 1 msg per chain every 5s = 4 msg/5s total

### SolanaDataSize
Overrides `DataLengthBytes` for Solana due to block size limitations.

### TestnetConfig
Configure chain-specific private keys and funding amounts:
- **EVMPrivateKey**: Private key for EVM chain funding account
- **AptosPrivateKey**: Private key for Aptos chain funding account
- **SolanaPrivateKey**: Private key for Solana chain funding account
- **FundingAmountEth**: Amount in wei to fund each EVM test account
- **FundingAmountSol**: Amount in lamports to fund each Solana test account
- **FundingAmountApt**: Amount in octas to fund each Aptos test account

### SuiConfig
Sui-specific configuration (nested under TestnetConfig):
- **SuiPrivateKey**: Private key for Sui account
- **SuiFeeTokenObjectId**: Object ID of the fee token on Sui
- **SuiTestReceiverAddress**: Address of the test receiver contract on Sui
- **SuiStateReceiverStateObjectId**: State object ID for the receiver on Sui

### NumDestinationChains
Must match the number of RPCs in `chain-details.json`.

### ChaosMode
- `0`: No chaos
- `1`: RPC latency (uses `RPCLatency` and `RPCJitter`)
- `2`: Full chaos suite (requires `[CCIP.Chaos]` configuration)

## Running the Test

```bash
cd integration-tests/load/ccip
go test -run ^TestCCIPLoad_RPS$ -v -timeout 0
```

## Test Flow

1. **Account Creation**: Creates N load accounts (where N = `NumDestinationChains`)
2. **Funding Phase**:
   - EVM accounts: Native tokens (per `FundingAmountEth`) + BnM tokens
   - Aptos accounts: Native tokens (per `FundingAmountApt`) + BnM tokens
   - Solana accounts: Native tokens (per `FundingAmountSol`)
   - Sui: Uses default specified account (no additional accounts)
3. **Load Execution**: Sends messages per `RequestFrequency` for `LoadDuration`

**Note:** Private keys for created accounts are logged for fund recovery.

## Monitoring

### Logs

Load statistics printed per iteration:
```
12:27PM INF Load stats CallTimeout=0 Component=ccipLoad Failed=0 Success=3
```

Account creation logs:
```
# Aptos
Created new Aptos sender on Chain 743186221051783445 | Address: 0x494c... | PrivateKey: ed25519-priv-0xfacc...

# EVM
New address created Addr=0x5c6ecd44e73A2C9Ec6fBb7be7ceae693410AA110
create account load testing account {"private key": "5e7f86ee09...", "selector": 13264668187771770619}
```

### Dashboards

Monitor test execution using available observability dashboards. Filter by `TestLabel` from config to view test-specific metrics.

## Fund Recovery

### Aptos

Use the fund recovery script to return APT from test accounts to main account. Script requires:
- Private keys from test logs
- Target address (main funding account)
- Updates to `privateKeys` slice in script

### EVM

Manually transfer funds using logged private keys.

### Sui

No fund recovery needed - uses single account.

## Funding Amounts

Current defaults (adjustable in config):
- **Aptos**: Configurable via `FundingAmountApt` (in octas)
- **EVM**: Configurable via `FundingAmountEth` (in wei)
- **Solana**: Configurable via `FundingAmountSol` (in lamports)
- **Sui**: Uses pre-funded account

## References

- Load test configuration: `integration-tests/testconfig/ccip/load.go`
- Test implementation: `integration-tests/load/ccip/ccip_test.go`
- Sui helpers: `integration-tests/load/ccip/sui_helpers.go`
- General helpers: `integration-tests/load/ccip/helpers.go`