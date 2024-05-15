# Seth-Specific Instructions

## Configuring with `seth.toml`

When operating Seth either as a library or through the Command Line Interface (CLI), it typically requires a path to the `seth.toml` file for configuration settings. In our workflow, however, you can bypass this requirement by utilizing the `TestConfig` TOML configuration
and configuring Seth settings directly in the TOML file(s).

### How to Set Configuration Values

Place all Seth-specific configuration entries under the `[Seth]` section in your TOML file. This can be done in files such as `default.toml` or `overrides.toml` (or any product-specific TOML) .

#### Example:
If your `Seth.toml` file contains a setting like this:
```toml
tracing_level = "reverted"
```

You can replicate this setting in `default.toml` or `overrides.toml` (or any product-specific TOML) by adding:
```toml
[Seth]
tracing_level = "reverted"
```

Documentation and Further Details
For a comprehensive description of all available configuration options, refer to the configuration documentation in the [Seth.toml](https://github.com/smartcontractkit/seth/blob/master/seth.toml) file or consult the [README.md on GitHub](https://github.com/smartcontractkit/seth/blob/master/README.md).

# Setting Log Level
## Locally

To adjust the logging level for Seth when running locally, use the environment variable `SETH_LOG_LEVEL`. For basic tracing and decoding information, set this variable to `debug`. For more detailed tracing, use the `trace` level.

## Remote Runner
To set the Seth log level in the Remote Runner, use the `TEST_SETH_LOG_LEVEL` environment variable. In the future, we plan to implement automatic forwarding of the `SETH_LOG_LEVEL` environment variable. Until then, you must set it explicitly.

# Seth Network Configuration
Seth's network configuration is entirely separate from the traditional `EVMNetwork`, and the two cannot be used interchangeably. Currently, both configurations must be provided for tests to function properly, as different parts of the test utilize each configuration.

## Overview of Configuration Usage
While most of the test logic relies on the `EVMNetwork` struct, Seth employs its own network configuration. To facilitate ease of use, we have introduced convenience methods that duplicate certain fields from `EVMNetwork` to `seth.Network`, eliminating the need to specify the same values twice. The following fields are automatically copied:

- Private keys
- RPC endpoints
- EIP-1559 support (only for simulated networks)
- Default gas limit (only for simulated networks)

## Example Code

Here is a raw example demonstrating how to merge configurations:

```go
import "github.com/smartcontractkit/chainlink-testing-framework/utils/seth"

readSethCfg := config.GetSethConfig()
if readSethCfg == nil {
    return nil, fmt.Errorf("Seth config not found")
}

sethCfg, err := utils.MergeSethAndEvmNetworkConfigs(network, *readSethCfg)
if err != nil {
    return nil, errors.Wrapf(err, "Error merging seth and evm network configs")
}
```

In practice, you rarely need to invoke these methods directly. It is recommended to use the following higher-level convenience methods:
```go
config, err := tc.GetConfig("Smoke", tc.OCR)
require.NoError(t, err, "Error getting config")

// validate Seth config before anything else
network := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0]

seth, err := actions_seth.GetChainClient(config, network)
require.NoError(t, err, "Error getting Seth client")
// or
sethWithConfigValidated, err := actions_seth.GetChainClientWithConfigFunction(config, network, actions_seth.OneEphemeralKeysLiveTestnetCheckFn)
require.NoError(t, err, "Error getting Seth client")
```

`actions_seth.OneEphemeralKeysLiveTestnetCheckFn` is a validation function that checks if the network is a live network and if ephemeral keys are being used. It is recommended to use this function when working with live networks. It makes
sure that you are not trying to use ephemeral keys on live networks, which would result in loss of funds associated with those keys.

## Seth-Specific Network Configuration
You are still expected to manually provide some Seth-specific network configurations related to the network you are using:
- Fallback gas price (for non-EIP-1559 networks), gas tip/fee cap (for EIP-1559 networks)
- Fallback gas limit (used for contract deployment and interaction)
- Fallback transfer fee (used for transferring funds between accounts)
- Network name
- Chain ID (critical for matching with EVMNetwork)
- Transaction timeout
- While this covers the essentials, it is advisable to consult the Seth documentation for detailed settings related to gas estimation, tracing, etc.

## Obtaining Fallback (Hardcoded) Values

There are two primary methods to obtain fallback values for network configurations:

1. **Web-Based Gas Tracker**: Model fallback values based on historical gas prices.
2. **Seth CLI**: This method is more complex but works for any network. We will focus on this due to its broad applicability.

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

4. **Copy Fallback Values:**
From the output, copy the relevant fallback prices into your network configuration. Here's an example of what you might see:
```bash
 5:08PM INF Fallback prices for TOML config:
gas_price = 121487901046
gas_tip_cap = 1000000000
gas_fee_cap = 122487901046 
```

This method ensures you get the most accurate and network-specific fallback values, optimizing your setup for current network conditions.

# Ephemeral and Static Keys

Understanding the difference between ephemeral and static keys is essential for effective and safe use of Seth. Here are a couple of use cases where you might need them:

- **Parallel Operations**: If you need to run multiple operations simultaneously without waiting for the previous one to finish, remember that Seth is synchronous and requires different keys for each goroutine.
- **Load Generation**: To generate a large volume of transactions in a short time.

## Understanding the Keys

- **Ephemeral Keys**: These are generated on the fly, are not stored, and should not be used on live networks, because any funds associated will be lost. Use these keys only when it's acceptable to lose the funds.
- **Static Keys**: These are specified in the Seth configuration and are suitable for use on live networks. Funds associated with these keys are not lost post-test since you retain copies of the private keys.

Most tests, especially on live networks, will restrict the use of ephemeral keys.

## How to Set Ephemeral Keys in the TOML

Setting ephemeral keys is straightforward:
```toml
[Seth]
ephemeral_addresses_number = 10
```

## How to Set Static Keys in the TOML

There are several methods to set static keys, but here are two:

### As a List of Wallet Keys in Your Network Configuration

Add it directly to your test TOML:
```toml
[Network.WalletKeys]
arbitrum_sepolia=["first_key","second_key"]
```
This method is ideal for local tests, but should be avoided in continuous integration (CI) environments.

### As Base64-Encoded Keyfile Stored as GHA Secret

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
gh secret set SETH_KEYFILE_BASE64 -b $(cat keyfile.toml | base64)
```

## How to Split Funds Between Static Keys

Managing funds across multiple static keys can be complex, especially if your tests require a substantial number of keys to generate adequate load. To simplify this process, follow these steps:

1. **Fund a Root Key**: Start by funding a key (referred to as the root key) with the total amount of funds you intend to distribute among other keys.
2. **Use Seth to Distribute Funds**: Execute the command below to split the funds from the root key to other keys:
```
KEYFILE_PATH=keyfile_my_network.toml ROOT_PRIVATE_KEY=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 SETH_CONFIG_PATH=seth.toml go run cmd/seth/seth.go -n my_network keys split -a 10 -b 1
```
The `-a <N>` option specifies the number of keys to distribute funds to, and `-b <N>` denotes the buffer (in ethers) to be left on the root key.

## Synchronous API

Seth is designed with a synchronous API to enhance usability and predictability. This feature is implemented through the `seth.Decode()` function, which waits for each transaction to be mined before proceeding. Depending on the Seth configuration, the function will:

- **Decode transactions only if they are reverted**: This is the default setting.
- **Always decode transactions**: This occurs if the `tracing_level` is set to `all`.

This approach simplifies the way transactions are handled, making them more predictable and easier to debug. Therefore, it is highly recommended that you wrap all contract interactions in that method.

# Getting Event Data from Transactions

Retrieving event data from transactions in Seth involves a few steps but is not overly complicated. Below is a Go function example that illustrates how to capture event data from a specific transaction:

```go
func (v *EthereumVRFCoordinatorV2_5) CancelSubscription(subID *big.Int, to common.Address) (*seth.DecodedTransaction, *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled, error) {
  tx, err := v.client.Decode(v.coordinator.CancelSubscription(v.client.NewTXOpts(), subID, to))
  if err != nil {
      return nil, nil, err
  }
  var cancelEvent *vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled
  for _, log := range tx.Receipt.Logs {
    for _, topic := range log.Topics {
        if topic.Cmp(vrf_coordinator_v2_5.VRFCoordinatorV25SubscriptionCanceled{}.Topic()) == 0 {
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
