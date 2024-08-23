# Automation Debugging Script

Use this script to debug and diagnose possible issues with registered upkeeps in Automation v2 registries. The script can debug custom logic upkeeps, log-trigger upkeeps, and upkeeps that use StreamsLookup.

## Setup

Before starting, you will need:

- A registered [upkeep](https://docs.chain.link/chainlink-automation/overview/getting-started)
- A working [Go](https://go.dev/doc/install) installation, please use this Go [version](https://github.com/smartcontractkit/chainlink/blob/develop/go.mod#L3)

1. Clone the chainlink [repo](https://github.com/smartcontractkit/chainlink) and navigate to the `core/scripts/chaincli`
    directory:
    ```
    git clone https://github.com/smartcontractkit/chainlink.git && cd chainlink/core/scripts/chaincli
    ```
1. Create a `.env` file based on the example `.env.debugging.example`:

    ```
    cp .env.debugging.example .env
    ```

## Configuration

Fill in the values for these mandatory fields in your `.env` file:

- `NODE_URL`: Archival node URL for the network to "simulate" the upkeep. Use your own node or get an endpoint from Alchemy or Infura.
- `KEEPER_REGISTRY_ADDRESS`: Address of the registry where your upkeep is registered. Refer to the [Supported Networks](https://docs.chain.link/chainlink-automation/overview/supported-networks#configurations) doc for registry addresses.
 
 For example
 ![Example_ENV_file](/core/scripts/chaincli/images/env_file_example.png "Example .ENV file")

#### StreamsLookup (optional)

If your targeted upkeep involves StreamsLookup, please provide the following details. If you are using Data Streams v0.3 (which is likely), only provide the `DATA_STREAMS_URL`. Ignore `DATA_STREAMS_LEGACY_URL`.

- `DATA_STREAMS_ID`
- `DATA_STREAMS_KEY`
- `DATA_STREAMS_URL`

#### Tenderly integration (optional)

For detailed transaction simulation logs, set up Tenderly credentials. Refer to the [Tenderly documentation](https://docs.tenderly.co/other/platform-access/how-to-generate-api-access-tokens) to learn how to create an API key, account name, and project name on Tenderly.

- `TENDERLY_KEY`
- `TENDERLY_ACCOUNT_NAME`
- `TENDERLY_PROJECT_NAME`

### Usage

Execute the following command based on your upkeep type:

- For custom logic: 

    ```bash
    go run main.go keeper debug UPKEEP_ID [BLOCK_NUMBER]
    ```
    If you don't specify a block number, the debugging script uses the latest block for checkUpkeep and simulatePerformUpkeep. For conditional upkeeps using streams lookup, a BLOCK_NUMBER is required.

- For log trigger upkeep:

    ```bash
    go run main.go keeper debug UPKEEP_ID TX_HASH LOG_INDEX
    ```

### What the debugging script checks

1. The script runs these basic checks on all upkeeps based on the TX_HASH or BLOCK_NUMBER (if provided)
    - Verify upkeep status: active, paused, or canceled
    - Check upkeep balance

2. **For Custom Logic Upkeep:**
    - Check conditional upkeep
    - Simulate `performUpkeep`

3. **For Log Trigger Upkeep:**
    - Check if the upkeep has already run for log-trigger-based upkeep
    - Verify if log matches trigger configuration
    - Check upkeep
    - If check result indicates a streams lookup is required (TargetCheckReverted):
        - Verify if the upkeep is allowed to use Mercury
        - Execute Mercury request
        - Execute check callback

    - Simulate `performUpkeep`

#### Examples
- Eligible and log trigger based and using mercury lookup v0.3:

    ```bash
    go run main.go keeper debug 5591498142036749453487419299781783197030971023186134955311257372668222176389 0xdc6d0e547a5aa85fefa5b0f3a37e3493eafb5aeba8b5f3071ce53c9e9a539e9c 0
    ```

- Ineligible and conditional upkeep:

    ```bash
    go run main.go keeper debug 52635131310730056105456985154251306793887717546629785340977553840883117540096
    ```

- Ineligible and Log does not match trigger config:

    ```bash
    go run main.go keeper debug 5591498142036749453487419299781783197030971023186134955311257372668222176389 0xc0686ae85d2a7a976ef46df6c613517b9fd46f23340ac583be4e44f5c8b7a186 1
    ```
### Common issues with Upkeeps and how to resolve them

#### All upkeeps

- Upkeep is underfunded
  - Underfunded upkeeps will not perform. Fund your upkeep in the Automation [app](https://automation.chain.link/)
- Upkeep is paused
  - Unpause your upkeep in the Automation [app](https://automation.chain.link/)
- Insufficient check gas
  - There is a limit of 10,000,000 (as per Automation v2) on the amount of gas that can be used to "simulate" your `checkUpkeep` function.
  - To diagnose if your upkeep is running out of check gas, you will need to enable the Tenderly options above and then open the simulation link once you run the script.
  - ![Insufficient Check Gas](/core/scripts/chaincli/images/insufficient_check_gas.png "Open the Tenderly simulation and switch to debug mode")
  - ![Out of Gas](/core/scripts/chaincli/images/tenderly_out_of_check_gas.png "Tenderly shows checkUpkeeps has consumed all available gas and is now out of gas")   
  - You will need to adjust your checkUpkeep to consume less gas than the limit
- Insufficient perform gas
  - Your upkeep's perform transaction uses more gas than you specified
  - ![Insufficient Perform Gas](/core/scripts/chaincli/images/insufficient_perform_gas.png "Insufficient perform gas")
  - Use the Automation [app](https://automation.chain.link/) and increase the gas limit of your upkeep
  - The maximum supported perform gas is 5,000,000

#### Log-trigger upkeeps

Log-trigger upkeeps require that you also supply the txn hash containing the log and the index of the log that would have triggered your upkeep. You can find both in the block scanner of the chain in question. For example the txn hash is in the URL and the block number in the green circle on the left.
![Txn Hash and Log Index Number](/core/scripts/chaincli/images/txnHash_and_index.png "Find txn hash and log index in block scanner") 

- Log doesn't match the trigger config
  - Log-trigger upkeeps come with a filter (aka trigger config), if the emitted log doesn't match the filter, the upkeep won't run.
  - ![Log doesn't match](/core/scripts/chaincli/images/log_trigger_log_doesnt_match.png "Log doesn't match trigger config") 
  - Use the Automation [app](https://automation.chain.link/) to update the upkeep's trigger config to match the log.
---
