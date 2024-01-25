##  Automation Debugging Script

### Context

The debugging script is a tool within ChainCLI designed to facilitate the debugging of upkeeps in Automation v21, covering both conditional and log-based scenarios.

### Setup

Before starting, you will need:
1. Git clone this chainlink [repo](https://github.com/smartcontractkit/chainlink)
2. A working [Go](https://go.dev/doc/install) installation
2. Change directory to `core/scripts/chaincli` and create a `.env` file based on the example `.env.debugging.example`

### Configuration in `.env` File

#### Mandatory Fields

Ensure the following fields are provided in your `.env` file:

- `NODE_URL`: Archival node URL
- `KEEPER_REGISTRY_ADDRESS`: Address of the Keeper Registry contract. Refer to the [Supported Networks](https://docs.chain.link/chainlink-automation/overview/supported-networks#configurations) doc for addresses.

#### Optional Fields (Streams Lookup)

If your targeted upkeep involves streams lookup, please provide the following details. If you are using Data Streams v0.3 (which is likely), only provide the DATA_STREAMS_URL. The DATA_STREAMS_LEGACY_URL is specifically for Data Streams v0.2.

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

- For conditional upkeep, if a block number is given we use that block, otherwise we use the latest block:

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

2. **For Conditional Upkeep:**
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
---