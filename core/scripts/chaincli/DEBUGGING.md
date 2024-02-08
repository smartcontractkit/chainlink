##  Automation Debugging Script

### Context

Use this script to debug and diagnose possible issues with registered upkeeps in Automation v2 registries. The script will debug both custom logic and log-trigger upkeeps as well as upkeeps using StreamsLookup.

### Setup

Before starting, you will need:
1. A registered [upkeep](https://docs.chain.link/chainlink-automation/overview/getting-started)
1. Git clone the chainlink [repo](https://github.com/smartcontractkit/chainlink)
1. A working [Go](https://go.dev/doc/install) installation
1. Change directory to `core/scripts/chaincli` and create a `.env` file based on the example `.env.debugging.example`

### Configuration in `.env` File

#### Mandatory Fields

Ensure the following fields are provided in your `.env` file:

- `NODE_URL`: Archival node URL for the network to "simulate" the upkeep. User your own node or get an endpoint from Alchemy or Infura
- `KEEPER_REGISTRY_ADDRESS`: Address of the Registry where your upkeep is registered. Refer to the [Supported Networks](https://docs.chain.link/chainlink-automation/overview/supported-networks#configurations) doc for addresses.
 
 For example
 ![Example_ENV_file](/images/env_file_example.png "Example .ENV file")

#### Optional Fields (StreamsLookup)

If your targeted upkeep involves StreamsLookup, please provide the following details. If you are using Data Streams v0.3 (which is likely), only provide the DATA_STREAMS_URL. Ignore DATA_STREAMS_LEGACY_URL.

- `DATA_STREAMS_ID`
- `DATA_STREAMS_KEY`
- `DATA_STREAMS_LEGACY_URL`
- `DATA_STREAMS_URL`

#### Optional Fields (Tenderly Integration)

For detailed transaction simulation logs, set up Tenderly credentials. Refer to the [Tenderly Documentation](https://docs.tenderly.co/other/platform-access/how-to-generate-api-access-tokens) for creating an API key, account name, and project name.

- `TENDERLY_KEY`
- `TENDERLY_ACCOUNT_NAME`
- `TENDERLY_PROJECT_NAME`

### Usage

Execute the following command based on your upkeep type:

- For custom logic, if a block number is given we use that block, otherwise we use the latest block:

    ```bash
    go run main.go keeper debug UPKEEP_ID [OPTIONAL BLOCK_NUMBER]
    ```

- For log trigger upkeep:

    ```bash
    go run main.go keeper debug UPKEEP_ID TX_HASH LOG_INDEX
    ```

### Checks Performed by the Debugging Script

1. **Fetch and Sanity Check Upkeep:**
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

### Examples
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
- Upkeep is underfunded
  - Underfunded upkeeps will not perform. Fund your upkeep in the Automation [app](https://automation.chain.link/)
- Upkeep is paused
  - Unpause your upkeep in the Automation [app](https://automation.chain.link/)
- Insufficient check gas
  - There is a limit of 10_000_000 (as per Automation v2) on the amount of gas that can be used to "simulate" your checkUpkeep function.
  - To diagnose if your upkeep is running out of check gas you will need to enable the Tenderly options above and then open the simulation link once you run the script
  - ![Insufficient Check Gas](/images/insufficient_check_gas.png "Open the Tenderly simulation and switch to debug mode")
  - ![Out of Gas](/images/tenderly_out_of_check_gas.png "Tenderly shows checkUpkeeps has consumed all available gas and is now out of gas")   
  - You will need to adjust your checkUpkeep to consume less gas than the limit
- Insufficient perform gas
  - Your upkeep's perform transaction uses more gas than you specified
  - ![Insufficient Perform Gas](/images/insufficient_perform_gas.png "Insufficient perform gas")
  - Use the Automation [app](https://automation.chain.link/) and increase the gas limit of your upkeep
  - The maximum supported perform gas is 5_000_000



---