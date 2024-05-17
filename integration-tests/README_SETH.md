# Seth-Specific Instructions

## Table of Contents
1. [Introduction](#introduction)
2. [How to Set Configuration Values](#how-to-set-configuration-values)
   1. [Example](#example)
   2. [Documentation and Further Details](#documentation-and-further-details)
3. [How to Set Seth Logging Level](#how-to-set-seth-logging-level)
   1. [Locally](#locally)
   2. [Remote Runner](#remote-runner)
4. [How to Set Seth Network Configuration](#how-to-set-seth-network-configuration)
   1. [Overview of Configuration Usage](#overview-of-configuration-usage)
   2. [Seth-Specific Network Configuration](#seth-specific-network-configuration)
   3. [Steps for Adding a New Network](#steps-for-adding-a-new-network)
      1. [Network is Already Defined in known_networks.go](#network-is-already-defined-in-known_networksgo)
      2. [It's a New Network](#its-a-new-network)
5. [How to Use Seth CLI](#how-to-use-seth-cli)
6. [How to Get Fallback (Hardcoded) Values](#how-to-get-fallback-hardcoded-values)
   1. [Steps to Use Seth CLI for Fallback Values](#steps-to-use-seth-cli-for-fallback-values)
7. [Ephemeral and Static Keys Explained](#ephemeral-and-static-keys-explained)
   1. [Understanding the Keys](#understanding-the-keys)
   2. [How to Set Ephemeral Keys in the TOML](#how-to-set-ephemeral-keys-in-the-toml)
   3. [How to Set Static Keys in the TOML](#how-to-set-static-keys-in-the-toml)
      1. [As a List of Wallet Keys in Your Network Configuration](#as-a-list-of-wallet-keys-in-your-network-configuration)
      2. [As Base64-Encoded Keyfile Stored as GHA Secret](#as-base64-encoded-keyfile-stored-as-gha-secret)
8. [How to Split Funds Between Static Keys](#how-to-split-funds-between-static-keys)
9. [How to Return Funds From Static Keys to the Root Key](#how-to-return-funds-from-static-keys-to-the-root-key)
   1. [How to Rebalance Static Keys](#how-to-rebalance-static-keys)
10. [How to Deal with "TX Fee Exceeds the Configured Cap" Error](#how-to-deal-with-tx-fee-exceeds-the-configured-cap-error)
11. [How to Use Seth's Synchronous API](#how-to-use-seths-synchronous-api)
12. [How to Read Event Data from Transactions](#how-to-read-event-data-from-transactions)

## Introduction

[Seth](https://github.com/smartcontractkit/seth) is the Ethereum client we use for integration tests. It is designed to be a thin wrapper over `go-ethereum` client that adds a couple of key features:
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
For a comprehensive description of all available configuration options, refer to the `[Seth]` section of configuration documentation in the [default.toml](./testconfig/default.toml) file or consult the Seth [README.md on GitHub](https://github.com/smartcontractkit/seth/blob/master/README.md).

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
You are still expected to manually provide some Seth-specific network configurations related to the network you are using:

- Fallback gas price
- Fallback gas tip/fee cap
- Fallback gas limit (used for contract deployment and interaction)
- Fallback transfer fee (used for transferring funds between accounts)
- Network name
- Chain ID (critical for matching with EVMNetwork)
- Transaction timeout

### Steps for adding a new network

#### Network is already defined in [known_networks.go](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/networks/known_networks.go)
In that case you need add only Seth-specific network configuration to `[[Seth.networks]]` table. Here's an example:
```toml
[[Seth.networks]]
name = "ARBITRUM_SEPOLIA"
chain_id = "421614"
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
evm_chain_id = 421614
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
* you need **both** networks: one for EVM and one for Seth
* websocket URL and private keys from the `EVMNetwork` will be copied over to the `Seth.Network` configuration so you don't need to provide them again
* it's advised to not set the gas limit, unless your test fails without it (might happen when interacting with new networks due bugs or gas estimation quirks); Seth will try to estimate gas for each interaction
* chain ID of `Seth.Network` must match the one from `EVMNetwork` configuration

While this covers the essentials, it is advisable to consult the Seth documentation for detailed settings related to gas estimation, tracing, etc.

## How to use Seth CLI
The most important thing to keep in mind that the CLI requires you to provide a couple of settings via environment variables, in addition to a TOML configuration file. Here's a general breakdown of the required settings:
* `keys` commands requires `SETH_KEYFILE_PATH`, `SETH_CONFIG_PATH` and `ROOT_PRIVATE_KEY` environment variables
* `gas` command requires `SETH_CONFIG_PATH` environment variable

You can find a sample `Seth.toml` file [here](https://github.com/smartcontractkit/seth/blob/master/seth.toml). Currently you cannot use your test TOML file as a Seth configuration file, but we will add ability that in the future.

## How to get Fallback (Hardcoded) Values
There are two primary methods to obtain fallback values for network configuration:
1. **Web-Based Gas Tracker**: Model fallback values based on historical gas prices.
2. **Seth CLI**: This method is more complex, but works for any network. We will focus on it due to its broad applicability.

### Steps to Use Seth CLI for Fallback Values
1. **Clone the Seth Repository:**
   Clone the repository from GitHub using:
```bash
git clone https://github.com/smartcontractkit/seth
```

2. **Configure Network Details:**
   Add your network details in the `seth.toml` file:
```toml
[[Networks]]
name = "my_network"
chain_id = "43113"
urls_secret = ["RPC you want to use"]
```

3. **Run Seth CLI:**
   Execute the command to get fallback gas prices:
```bash
SETH_CONFIG_PATH=seth.toml go run cmd/seth/seth.go -n my_network gas -b 10000 -tp 0.99
```
The network name passed in the CLI must match the one in your TOML file (it is case-sensitive). The `-b` flag specifies the number of blocks to consider for gas estimation, and `-tp` denotes the tip percentage.

4. **Copy Fallback Values:**
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
chain_id = "667"
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
KEYFILE_PATH=keyfile_my_network.toml ROOT_PRIVATE_KEY=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 SETH_CONFIG_PATH=seth.toml go run cmd/seth/seth.go -n my_network keys split -a 10 -b 1
```
The `-a <N>` option specifies the number of keys to distribute funds to, and `-b <N>` denotes the buffer (in ethers) to be left on the root key.

## How to Return Funds From Static Keys to the Root Key
Returning funds from static keys to the root key is a simple process. Execute the following command:
```bash
KEYFILE_PATH=keyfile_my_network.toml ROOT_PRIVATE_KEY=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 SETH_CONFIG_PATH=seth.toml go run cmd/seth/seth.go -n my_network keys return
```
This command will return all funds from static keys read from `keyfile_my_network.toml` to the root key.

## How to Rebalance Static Keys
Rebalancing static keys is a more complex process that involves redistributing funds among keys. Currently, there's no built-in functionality for this in Seth, but you can achieve it by following these steps:
1. **Return Funds**: Use the `keys return` command to return all funds to the root key.
2. **Split Funds**: Use the `keys split` command to redistribute funds among the keys as needed.

Once you've completed these steps, remember to upload new keyfile to the CI (as a base64-ed secret).

**When performing any keyfile-related operations it is advised to keep copies of files in 1password, so you can easily restore them if needed**. That's especially important for rebalancing, because you will not be able to download the keyfile from the CI since it's a secret.

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
