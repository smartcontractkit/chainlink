# Seth-Specific Instructions

## Table of Contents
- [Seth-Specific Instructions](#seth-specific-instructions)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [How to Set Configuration Values](#how-to-set-configuration-values)
    - [Example:](#example)
    - [Documentation and Further Details](#documentation-and-further-details)
  - [How to set Seth logging level](#how-to-set-seth-logging-level)
    - [Locally](#locally)
    - [Remote Runner](#remote-runner)
  - [How to set Seth Network Configuration](#how-to-set-seth-network-configuration)
    - [Overview of Configuration Usage](#overview-of-configuration-usage)
    - [Seth-Specific Network Configuration](#seth-specific-network-configuration)
    - [Steps for adding a new network](#steps-for-adding-a-new-network)
      - [Network is already defined in known\_networks.go](#network-is-already-defined-in-known_networksgo)
      - [It's a new network](#its-a-new-network)
    - [Things to remember:](#things-to-remember)
  - [How to use Seth CLI](#how-to-use-seth-cli)
  - [How to get Fallback (Hardcoded) Values](#how-to-get-fallback-hardcoded-values)
    - [Steps to Use Seth CLI for Fallback Values](#steps-to-use-seth-cli-for-fallback-values)
  - [Ephemeral and Static Keys explained](#ephemeral-and-static-keys-explained)
    - [Understanding the Keys](#understanding-the-keys)
    - [How to Set Ephemeral Keys in the TOML](#how-to-set-ephemeral-keys-in-the-toml)
    - [How to Set Static Keys in the TOML](#how-to-set-static-keys-in-the-toml)
      - [As a List of Wallet Keys in Your Network Configuration](#as-a-list-of-wallet-keys-in-your-network-configuration)
      - [As Base64-Encoded Keyfile Stored as GHA Secret](#as-base64-encoded-keyfile-stored-as-gha-secret)
    - [Setting Up Your Pipeline](#setting-up-your-pipeline)
  - [How to Split Funds Between Static Keys](#how-to-split-funds-between-static-keys)
  - [How to Return Funds From Static Keys to the Root Key](#how-to-return-funds-from-static-keys-to-the-root-key)
  - [How to Rebalance or Top Up Static Keys](#how-to-rebalance-or-top-up-static-keys)
    - [Rebalancing Static Keys](#rebalancing-static-keys)
    - [Topping Up Static Keys](#topping-up-static-keys)
  - [How to Deal with "TX Fee Exceeds the Configured Cap" Error](#how-to-deal-with-tx-fee-exceeds-the-configured-cap-error)
  - [How to use Seth's synchronous API](#how-to-use-seths-synchronous-api)
  - [How to read Event Data from transactions](#how-to-read-event-data-from-transactions)
  - [How to deal with Failed Transactions](#how-to-deal-with-failed-transactions)
  - [How Automated Gas Estimation works](#how-automated-gas-estimation-works)
        - [Legacy Transactions](#legacy-transactions)
        - [EIP-1559 Transactions](#eip-1559-transactions)
        - [Adjustment factor](#adjustment-factor)
        - [Buffer percents](#buffer-percents)
  - [How to tweak Automated Gas Estimation](#how-to-tweak-automated-gas-estimation)
  - [How to debug with 'execution reverted' error](#how-to-debug-with-execution-reverted-error)
  - [How to split non-native tokens between keys](#how-to-split-non-native-tokens-between-keys)

## Introduction

[Seth](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/seth) is the Ethereum client we use for integration tests. It is designed to be a thin wrapper over `go-ethereum` client that adds a couple of key features:
* key management
* transaction decoding and tracing
* gas estimation

To use it you don't need to add any specific configuration to your TOML files. Reasonable defaults have been added to `default.toml` file under `[Seth]` section. For some of the products
we have added a couple of product-scoped overrides. For example for Automation's Load tests we have increased ephemeral addresses count from `10` to `100`:
```toml
[Load.Seth]
ephemeral_addresses_number = 100
```

Feel free to modify the configuration to suit your needs, but remember to scope it correctly, so that it doesn't impact other products. You can find more information about TOML configuration and override precedences [here](./testconfig/README.md).

## How to Set Configuration Values
Place all Seth-specific configuration entries under the `[Seth]` section in your TOML file. This can be done in files such as `default.toml` or `overrides.toml` or any product-specific TOML located in the [testconfig](./testconfig) folder.

### Example:
```toml
[Seth]
tracing_level = "all" # trace all transactions regardless of whether they are reverted or not
```

### Documentation and Further Details
For a comprehensive description of all available configuration options, refer to the `[Seth]` section of configuration documentation in the [default.toml](./testconfig/default.toml) file or consult the Seth [README.md on GitHub](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/seth/blob/master/README.md).

## How to set Seth logging level
### Locally
To adjust the logging level for Seth when running locally, use the environment variable `SETH_LOG_LEVEL`. For basic tracing and decoding information, set this variable to `debug`. For more detailed tracing, use the `trace` level.
### Remote Runner
To set the Seth log level in the Remote Runner, use the `TEST_SETH_LOG_LEVEL` environment variable. In the future, we plan to implement automatic forwarding of the `SETH_LOG_LEVEL` environment variable. Until then, you must set it explicitly.

## How to set Seth Network Configuration
Seth's network configuration is entirely separate from the traditional `EVMNetwork`, and the two cannot be used interchangeably. Currently, both configurations must be provided for tests to function properly, as different parts of the test utilize each configuration.

### Overview of Configuration Usage
While most of the test logic relies on the `EVMNetwork` struct, Seth employs its own network configuration. To facilitate ease of use, we have introduced convenience methods that duplicate certain fields from `EVMNetwork` to `seth.Network`, eliminating the need to specify the same values twice. The following fields are automatically copied:

- Private keys
- RPC endpoints
- EIP-1559 support (only for simulated networks)
- Default gas limit (only for simulated networks)

### Seth-Specific Network Configuration
Since version v1.0.11 Seth will use a default network configuration if you don't provide one in your TOML file. It is stored in `default.toml` and named `Default`. It has both EIP-1559 and gas estimation enabled, because Seth will disable both if they are not supported by the network.
If for whatever reason you want to use different default settings for you product, please add them to your product-specific TOML file.

You can still provide your own network configuration in your TOML file, if you need to override the default settings. When you do that, you need to provide all the following fields:
- Fallback gas price
- Fallback gas tip/fee cap
- Fallback gas limit (used for contract deployment and interaction)
- Fallback transfer fee (used for transferring funds between accounts)
- Network name
- Transaction timeout

### Steps for adding a new network

#### Network is already defined in [known_networks.go](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/networks/known_networks.go)
If you are fine with the default network configuration, you don't need to do anything. Otherwise, you need add only Seth-specific network configuration to `[[Seth.networks]]` table. Here's an example:
```toml
[[Seth.networks]]
name = "ARBITRUM_SEPOLIA"
transaction_timeout = "10m"
transfer_gas_fee = 50_000
# gas_limit = 15_000_000
# legacy transactions fallback gas price
gas_price = 200_000_000
# EIP-1559 transactions fallback gas tip cap and fee cap
eip_1559_dynamic_fees = true
gas_fee_cap = 400_000_000
gas_tip_cap = 200_000_000
# if set to true we will estimate gas for every transaction
gas_price_estimation_enabled = true
# how many last blocks to use, when estimating gas for a transaction
gas_price_estimation_blocks = 100
# priority of the transaction, can be "fast", "standard" or "slow" (the higher the priority, the higher adjustment factor will be used for gas estimation) [default: "standard"]
gas_price_estimation_tx_priority = "standard"
```

Name of the network doesn't really matter and is used only for logging purposes. Chain ID must match the one from `known_networks.go` file.

**Warning!** Please do not use the values from above-mentioned example. They should be replaced with the actual values obtained from gas tracker or Seth CLI (more on that later). 

#### It's a new network
Apart from above-mentioned fields you either need to add the network to `known_networks.go` file in the [CTF](https://github.com/smartcontractkit/chainlink-testing-framework) or define it in your test TOML file. 
Here's an example of how to define a new `EVMNetwork` network in your test TOML file:
```toml
[Network.EVMNetworks.ARBITRUM_SEPOLIA]
evm_name = "ARBITRUM_SEPOLIA"
evm_urls = ["rpc ws endpoint"]
evm_http_urls = ["rpc http endpoint"]
client_implementation = "Ethereum"
evm_keys = ["private keys you want to use"]
evm_simulated = false
evm_chainlink_transaction_limit = 5000
evm_minimum_confirmations = 1
evm_gas_estimation_buffer = 10000
evm_supports_eip1559 = false
evm_default_gas_limit = 6000000
```

### Things to remember:
* you need **both** networks: one for EVM and one for Seth (unless you are fine with the default settings)
* websocket URL and private keys from the `EVMNetwork` will be copied over to the `Seth.Network` configuration, so you don't need to provide them again
* it's advised to not set the gas limit, unless your test fails without it (might happen when interacting with new networks due bugs or gas estimation quirks); Seth will try to estimate gas for each interaction
* **name** of `Seth.Network` must match the one from `EVMNetwork` configuration

While this covers the essentials, it is advisable to consult the Seth documentation for detailed settings related to gas estimation, tracing, etc.

## How to use Seth CLI
The most important thing to keep in mind that the CLI requires you to provide a couple of settings via environment variables, in addition to a TOML configuration file. Here's a general breakdown of the required settings:
* `keys` commands requires `SETH_KEYFILE_PATH`, `SETH_CONFIG_PATH` and `SETH_ROOT_PRIVATE_KEY` environment variables
* `gas` and `stats` command requires `SETH_CONFIG_PATH` environment variable

You can find a sample `Seth.toml` file [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/seth/blob/master/seth.toml). Currently, you cannot use your test TOML file as a Seth configuration file, but we will add ability that in the future.

## How to get Fallback (Hardcoded) Values
There are two primary methods to obtain fallback values for network configuration:
1. **Web-Based Gas Tracker**: Model fallback values based on historical gas prices.
2. **Seth CLI**: This method is more complex, but works for any network. We will focus on it due to its broad applicability.

### Steps to Use Seth CLI for Fallback Values
1. **Clone the Seth Repository:**
   Clone the repository from GitHub using:
```bash
git clone https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/seth
```

2. **Run Seth CLI:**
   Execute the command to get fallback gas prices:
```bash
SETH_CONFIG_PATH=seth.toml go run cmd/seth/seth.go -u https://RPC_TO_USE -b 10000 -tp 0.99
```
The network name passed in the CLI must match the one in your TOML file (it is case-sensitive). The `-b` flag specifies the number of blocks to consider for gas estimation, and `-tp` denotes the tip percentage.

3**Copy Fallback Values:**
   From the output, copy the relevant fallback prices into your network configuration in test TOML. Here's an example of what you might see:
```bash
 5:08PM INF Fallback prices for TOML config:
gas_price = 121487901046
gas_tip_cap = 1000000000
gas_fee_cap = 122487901046 
```

5. **Update TOML Configuration:**
   Update your network configuration with the copied values:
```toml
[[Seth.networks]]
name = "my_network"
transaction_timeout = "10m"
transfer_gas_fee = 21_000
eip_1559_dynamic_fees = true
gas_fee_cap = 122487901046
gas_tip_cap = 1000000000
gas_price_estimation_enabled = true
gas_price_estimation_blocks = 100
gas_price_estimation_tx_priority = "standard"
```

This method ensures you get the most accurate and network-specific fallback values, optimizing your setup for current network conditions.

## Ephemeral and Static Keys explained
Understanding the difference between ephemeral and static keys is essential for effective and safe use of Seth. 

### Understanding the Keys
- **Ephemeral Keys**: These are generated on the fly, are not stored, and should not be used on live networks, because any funds associated will be lost. Use these keys only when it's acceptable to lose the funds.
- **Static Keys**: These are specified in the Seth configuration and are suitable for use on live networks. Funds associated with these keys are not lost post-test since you retain copies of the private keys.

Here are a couple of use cases where you might need to use ephemeral keys or more than one static key:

- **Parallel Operations**: If you need to run multiple operations simultaneously without waiting for the previous one to finish, remember that Seth is synchronous and requires different keys for each goroutine.
- **Load Generation**: To generate a large volume of transactions in a short time.

Most tests, especially on live networks, will restrict the use of ephemeral keys.

### How to Set Ephemeral Keys in the TOML
Setting ephemeral keys is straightforward:
```toml
[Seth]
ephemeral_addresses_number = 10
```

### How to Set Static Keys in the TOML
There are several methods to set static keys, but here are two:

#### As a List of Wallet Keys in Your Network Configuration
Add it directly to your test TOML:
```toml
[Network.WalletKeys]
arbitrum_sepolia=["first_key","second_key"]
```
This method is ideal for local tests, but should be avoided in continuous integration (CI) environments.

#### As Base64-Encoded Keyfile Stored as GHA Secret
This safer, preferred method involves more setup:

1. **Configuration**: Your Seth must be configured to read the keyfile in Base64-encoded version from an environment variable, by setting in your TOML:
```
[Seth]
keyfile_source = "base64_env"
```
2. **Pipeline Setup**: Your pipeline must have the secret with the Base64-encoded keyfile exposed as an environment variable named `SETH_KEYFILE_BASE64`. Seth will automatically read and decode it given the above-mentioned configuration.

### Setting Up Your Pipeline
Here's how to add the keyfile to your GitHub Actions secrets:
1. Create a keyfile (instructions provided below).
2. Base64-encode the keyfile and add it to your GitHub Actions secrets using the GitHub CLI:
```
gh secret set SETH_MY_NETWORK_KEYFILE_BASE64 -b $(cat keyfile_my_network.toml | base64)
```

It is advised to use a separate keyfile for each network to avoid confusion and potential errors. If you need to run your test on multiple networks you should add logic to your pipeline that will set the correct keyfile based on the network you are testing.

## How to Split Funds Between Static Keys
Managing funds across multiple static keys can be complex, especially if your tests require a substantial number of keys to generate adequate load. To simplify this process, follow these steps:

1. **Fund a Root Key**: Start by funding a key (referred to as the root key) with the total amount of funds you intend to distribute among other keys.
2. **Use Seth to Distribute Funds**: Execute the command below to split the funds from the root key to other keys:
```
SETH_KEYFILE_PATH=keyfile_my_network.toml SETH_ROOT_PRIVATE_KEY=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 SETH_CONFIG_PATH=seth.toml go run cmd/seth/seth.go -n my_network keys fund -a 10 -b 1
```
The `-a <N>` option specifies the number of keys to distribute funds to, and `-b <N>` denotes the buffer (in ethers) to be left on the root key.

## How to Return Funds From Static Keys to the Root Key
Returning funds from static keys to the root key is a simple process. Execute the following command:
```bash
KEYFILE_PATH=keyfile_my_network.toml ROOT_PRIVATE_KEY=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 SETH_CONFIG_PATH=seth.toml go run cmd/seth/seth.go -n my_network keys return
```
This command will return all funds from static keys read from `keyfile_my_network.toml` to the root key.

## How to Rebalance or Top Up Static Keys
Currently, Seth doesn't have a single step function for this, but you can follow these steps to rebalance or top up static keys:

### Rebalancing Static Keys
1. **Return Funds to Root Key**: Use the command `keys return` to transfer all funds back to the root key.
2. **Redistribute Funds**: Use the command `keys funds` to allocate the funds among the keys as needed.

After rebalancing, upload the new keyfile to the CI as a base64-encoded secret.

### Topping Up Static Keys
1. **Fund the Root Key**: Add funds to the root key.
2. **Distribute Funds**: Use the `keys fund` command to distribute these funds among the keys. 

The key is to understand that `keys fund` command will use existing keys if a keyfile is found, or generate new private keys and save them if the file doesn't exist.

For both rebalancing and topping up, you don't need to upload the keyfile to the CI again, as the private keys remain the same, only their on-chain balances change.

**Tip**: Always keep copies of keyfiles in 1Password to easily restore them if needed. This is crucial for rebalancing because you cannot download the keyfile from the CI since it's a secret.

## How to Deal with "TX Fee Exceeds the Configured Cap" Error
If the gas prices set for a transaction and the gas limit result in the transaction fee exceeding the maximum fee set for a given RPC node, you can try the following solutions:
1. **Try a Different RPC Node**: This setting is usually specific to the RPC node. Try using a different node, as it might have a higher limit.
2. **Decrease Gas Price**: If you are using legacy transactions and not using automated gas estimation, try decreasing the gas price in your TOML configuration. This will lower the overall transaction fee.
3. **Decrease Fee Cap**: If you are using EIP-1559 transactions and not using automated gas estimation, try decreasing the fee cap in your TOML configuration. You should also decrease the tip cap accordingly, as the fee cap includes both the base fee and the tip. This will lower the overall transaction fee.
4. **Decrease Gas Limit**: If you are using a hardcoded gas limit, try decreasing it. This will lower the overall transaction fee regardless of the transaction type.
5. **Use Gas Estimation**: If you are not using automated gas estimation, enable it. This will make Seth estimate gas for each transaction and adjust the gas price accordingly, which could prevent the error if your hardcoded gas-related values were too high.
6. **Use Different Gas Estimation Settings**: If you are using automated gas estimation, try lowering the gas estimation settings in your TOML configuration. Adjust the number of blocks used for estimation or the priority of the transaction.
7. **Disable Gas Estimations**: If you are using automated gas estimation, you can try disabling it. This will make Seth use the hardcoded gas-related values from your TOML configuration. This could prevent the error if you set the values low enough, but be aware it might lead to other issues, such as long waits for transaction inclusion in a block.

## How to use Seth's synchronous API
Seth is designed with a synchronous API to enhance usability and predictability. This feature is implemented through the `seth.Decode()` function, which waits for each transaction to be mined before proceeding. Depending on the Seth configuration, the function will:

- **Decode transactions only if they are reverted**: This is the default setting.
- **Always decode transactions**: This occurs if the `tracing_level` is set to `all`.
- **Always try to get revert reason**: if the transaction is reverted, Seth will try to get the revert reason, regardless of the `tracing_level` setting.

This approach simplifies the way transactions are handled, making them more predictable and easier to debug. Therefore, it is highly recommended that you wrap all contract interactions in that method.

## How to read Event Data from transactions
Retrieving event data from transactions in Seth involves a few steps but is not overly complicated. Below is a Go function example that illustrates how to capture event data from a specific transaction:

```go
func (v *EthereumVRFCoordinatorV2_5) CancelSubscription(subID *big.Int, to common.Address) (*seth.DecodedTransaction, *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled, error) {
 // execute a transaction
  tx, err := v.client.Decode(v.coordinator.CancelSubscription(v.client.NewTXOpts(), subID, to))
  if err != nil {
      return nil, nil, err
  }
  
  // define the event you are looking for
  var cancelEvent *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled
  // iterate over receipt logs until you find a topic that matches the event you are looking for
  for _, log := range tx.Receipt.Logs {
    for _, topic := range log.Topics {
        if topic.Cmp(vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled{}.Topic()) == 0 {
			// use geth wrapper to parse the log to the event struct
            cancelEvent, err = v.coordinator.ParseSubscriptionCanceled(*log)
            if err != nil {
                return nil, nil, fmt.Errorf("parsing SubscriptionCanceled log failed, err: %w", err)
            }
		}
    }
  }
  return tx, cancelEvent, err
}
```

This function demonstrates how to decode a transaction, check for specific event topics in the transaction logs, and parse those events. It's a structured approach to handling event data that is crucial for applications that need to respond to specific changes in state reflected by those events.

## How to deal with Failed Transactions
When a state-changing interaction with a contract fails it's often difficult to know why it failed without using dedicated tools like Tenderly. That's why we have included transaction tracing and decoding in Seth. By default, it only applies to reverted transactions,
but you can enable it for all transactions by setting `tracing_level` to `all` in your TOML configuration. The information is printed to the console (with a mix of `debug` and `trace` levels), but you can also automatically save tracing information to a json file by setting `trace_to_json=true`.

If you suspect your tests might run into failed transactions in hard to reproduce circumstances it's advised to:
* set `SETH_LOG_LEVEL` and `TEST_SETH_LOG_LEVEL` to `trace`
* for tests running in the CI/Docker enable `trace_to_json` in your TOML configuration and add upload `./traces` directory to the artifacts in your CI pipeline (it will be located in the same directory as `./logs`)
* make sure that in your code you upload contracts only using Seth's `ContractLoader` helper struct as it will automatically add contract's ABI to Seth's ABI cache -- and that's a prerequisite for decoding transactions (more info below)
* make sure that in your code when you load contracts manually you always add its ABI to Seth's ABI cache using `seth.ContractStore.AddABI()` function
* make sure that contracts you are using are not very heavily optimised (as the optimiser might remove custom revert reasons) and that your EVM target version is >= Constantinople hard fork (as custom revert reasons are not supported in earlier versions)

Example of loading contract using `ContractLoader`:
```go
func LoadOffchainAggregator(client *seth.Client, contractAddress common.Address) (offchainaggregator.OffchainAggregator, error) {
	// initialize contract loader with the generic type of the Geth contract wrapper
	loader := seth.NewContractLoader[offchainaggregator.OffchainAggregator](client)
	// call load function with contract name, address, ABI getter function and Geth wrapper constructor function
	ocr, err := loader.LoadContract("OffChainAggregator", contractAddress, offchainaggregator.OffchainAggregatorMetaData.GetAbi, offchainaggregator.NewOffchainAggregator)
	if err != nil {
		return EthereumOffchainAggregator{}, fmt.Errorf("failed to instantiate OCR instance: %w", err)
	}
 }
```

## How Automated Gas Estimation works
Gas estimation varies based on whether the network is a private Ethereum Network or a live network.

- **Private Ethereum Networks**: no estimation is ever done. We always use hardcoded values.

For real networks, the estimation process differs for legacy transactions and those compliant with EIP-1559:

##### Legacy Transactions
1. **Initial Price**: Query the network node for the current suggested gas price.
2. **Priority Adjustment**: Modify the initial price based on `gas_price_estimation_tx_priority`. Higher priority increases the price to ensure faster inclusion in a block.
3. **Congestion Analysis**: Examine the last X blocks (as specified by `gas_price_estimation_blocks`) to determine network congestion, calculating the usage rate of gas in each block and giving recent blocks more weight.
4. **Buffering**: Add a buffer to the adjusted gas price to increase transaction reliability during high congestion.

##### EIP-1559 Transactions
1. **Tip Fee Query**: Ask the node for the current recommended tip fee.
2. **Fee History Analysis**: Gather the base fee and tip history from recent blocks to establish a fee baseline.
3. **Fee Selection**: Use the greater of the node's suggested tip or the historical average tip for upcoming calculations.
4. **Priority and Adjustment**: Increase the base and tip fees based on transaction priority (`gas_price_estimation_tx_priority`), which influences how much you are willing to spend to expedite your transaction.
5. **Final Fee Calculation**: Sum the base fee and adjusted tip to set the `gas_fee_cap`.
6. **Congestion Buffer**: Similar to legacy transactions, analyze congestion and apply a buffer to both the fee cap and the tip to secure transaction inclusion.

Understanding and setting these parameters correctly ensures that your transactions are processed efficiently and cost-effectively on the network.

Finally, `gas_price_estimation_tx_priority` is also used, when deciding, which percentile to use for base fee and tip for historical fee data. Here's how that looks:
```go
		case Priority_Fast:
			baseFee = stats.GasPrice.Perc99
			historicalGasTipCap = stats.TipCap.Perc99
		case Priority_Standard:
			baseFee = stats.GasPrice.Perc50
			historicalGasTipCap = stats.TipCap.Perc50
		case Priority_Slow:
			baseFee = stats.GasPrice.Perc25
			historicalGasTipCap = stats.TipCap.Perc25
```

##### Adjustment factor
All values are multiplied by the adjustment factor, which is calculated based on `gas_price_estimation_tx_priority`:
```go
	case Priority_Fast:
		return 1.2
	case Priority_Standard:
		return 1.0
	case Priority_Slow:
		return 0.8
```
For fast transactions we will increase gas price by 20%, for standard we will use the value as is and for slow we will decrease it by 20%.

##### Buffer percents
We further adjust the gas price by adding a buffer to it, based on congestion rate:
```go
	case Congestion_Low:
		return 1.10, nil
	case Congestion_Medium:
		return 1.20, nil
	case Congestion_High:
		return 1.30, nil
	case Congestion_VeryHigh:
		return 1.40, nil
```

For low congestion rate we will increase gas price by 10%, for medium by 20%, for high by 30% and for very high by 40%.  
We cache block header data in an in-memory cache, so we don't have to fetch it every time we estimate gas. The cache has capacity equal to `gas_price_estimation_blocks` and every time we add a new element, we remove one that is least frequently used and oldest (with block number being a constant and chain always moving forward it makes no sense to keep old blocks).

It's important to know that in order to use congestion metrics we need to fetch at least 80% of the requested blocks. If that fails, we will skip this part of the estimation and only adjust the gas price based on priority.

For both transaction types if any of the steps fails, we fallback to hardcoded values.

## How to tweak Automated Gas Estimation
Now that you know how automated gas estimation works, you can tweak it to better suit your needs. Here are some tips to help you optimize your gas estimation process:
* **Adjust the Gas Price Estimation Blocks**: Increase or decrease the number of blocks used for gas estimation based on network conditions. More blocks provide a more accurate picture of network congestion and can smooth out short spikes of high gas prices. On the other hand, fewer blocks can speed up the estimation process or even be a prerequisite for the estimation to work (if the RPC is slow). Remember also that longer block range can lead to less accurate estimation, if network conditions changed very recently (although the algorithm tries to counter that by assigning higher weights to more recent blocks.
* **Set the Gas Price Estimation Tx Priority**: Choose the transaction priority that best suits your needs. A higher priority will increase the gas price, ensuring faster inclusion in a block. Conversely, a lower priority will reduce the gas price, potentially saving you money.
* **Disable Gas Estimation**: If you prefer to use hardcoded values, you can disable automated gas estimation. This will make Seth use the hardcoded values from your TOML configuration, which can be useful if you want to avoid the overhead of estimation or if you have specific gas prices you want to use or in a rare case that estimations are inaccurate.

## How to debug with 'execution reverted' error
When you encounter the 'execution reverted' error without any additional details, it can be challenging to determine the cause. This error indicates that the Ethereum Virtual Machine (EVM) couldn't return a specific revert reason. Here are some tips to help you debug this issue:
1. **Set SETH_LOG_LEVEL**: Increase the log level to `trace` to get more detailed information about the transaction. This will help you identify the cause of the revert.
2. **Use Custom Revert Reasons**: Ensure your contract uses custom error/custom messages with `revert`/`require` statements. Without them, you won't get specific reasons for the error.
3. **Check Solidity Version**: Make sure your contract uses a Solidity version that supports custom revert reasons (>= 0.8.4). The version is specified at the beginning of your contract file, like this:
   ```solidity
   // SPDX-License-Identifier: MIT
   pragma solidity ^0.8.4;
   ```
   If the version is lower, you won't get custom revert reasons. Upgrading the version might introduce breaking changes, so do not commit these changes without thorough testing.
4. **Inspect `delegatecall` Usage**: If the method you are calling uses `delegatecall`, it might be the cause of the issue. `delegatecall` is a low-level call that doesn't return revert reasons, so you'll always get the 'execution reverted' error unless you manually handle it in assembly.
5. **Debug with Events**: If you're still unable to find the cause, you can modify the contract to emit an event with the revert reason. For example, define an event like `event DebugEvent(string reason)` and add it before the points where you suspect a revert might occur. You will be able to find events in the transaction receipt's logs. Remember to remove these changes after debugging.
6. **Check EVM Node Logs**: If all else fails, the issue might be that Seth cannot decode the revert reason. Look at the logs of the EVM node for 'execution reverted' errors and any additional information (e.g., `errdata=0x46f08154`). If you find more details there but not in your test output, it might be a bug in Seth, and you should contact the Test Tooling team.

By following these steps, you can systematically identify and address the cause of 'execution reverted' errors in your smart contracts.

## How to split non-native tokens between keys
Currently, Seth doesn't support splitting non-native tokens between keys. If you need to split non-native tokens between keys, you can use the following approach:
1. **Fund root key** with the total amount of tokens you want to distribute.
2. **Deploy v contract** that allows to execute multiple operations in a single transaction.
3. **Prepare payload for Multicall contract** that will call `transfer` function on the token contract for each key.
4. **Execute Multicall contract** with the payload prepared in the previous step.


You can find sample code for doing that for LINK token in [actions](./actions/actions.go) file as `SendLinkFundsToDeploymentAddresses()` method.
