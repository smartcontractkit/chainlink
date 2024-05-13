# Seth-specific instructions

# Using `seth.toml`
When used a standalone or via CLI Seth is expecting to be given a path to `seth.toml` file with configuration. In our case it's not necessary as Seth-specific config
can be passed using `TestConfig` TOML configuration. All Seth-specific values have to placed under `[Seth]` table in the TOML file.

For example, if `Seth.toml` has a configuration key:
```toml
tracing_level = "reverted"
```

You can set that, for example in `default.toml` or `overrides.toml` by adding:
```toml
[Seth]
tracing_level = "reverted"
```

Detailed description of all the available options can be found in the [Seth.toml](https://github.com/smartcontractkit/seth/blob/master/seth.toml) or its [README.md](https://github.com/smartcontractkit/seth/blob/master/README.md). 

# Setting log level
## Locally
Use `SETH_LOG_LEVEL` env var. If you want to see basic tracing/decoding information be sure to set it to `debug`. Even more detailed traces
are available with `trace` level.

## Remote Runner
Use `TEST_SETH_LOG_LEVEL` to set Seth log level in the Remote Runner. In the future we will add auto-forwarding for `SETH_LOG_LEVEL` env var, but for now you need to do it explicitly.

# Seth network configuration
Seth's network config is completely separate from good old `EVMNetwork` and you can't use it interchangeably. For now you need to provide both for tests to work, since each
of them is used in different parts of the test.

Most of test logic is based on the `EVMNetwork` struct, but Seth is using its own network configuration. We have added convenience methods that copy some of the `EVMNetwork` fields
to `seth.Network` so that you don't have to provide the same values twice. This is true for following fields:
* private keys
* rpc endpoints
* EIP-1559 support (simulated networks only)
* default gas limit (simulated networks only)

Raw example:
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

In reality there's rarely a need to even call these methods directly. It's recommended to use following higher-level convenience methods:
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

Nonetheless you are still expected to provide some Seth-specific network configuration related to the network you are using, like:
* fallback gas price (for non-EIP-1559 networks), gas tip/fee cap (for EIP-1559 networks)
* fallback gas limit (used for contract deployment and interaction)
* fallback transfer fee (used for transferring funds between accounts)
* network name 
* chain id (really important as it's used for matching it with `EVMNetwork`)
* transaction timeout

That's the bare minimum, but it's advised to read Seth documentation to learn about settings related to gas estimation, tracing, etc.

## Where to get fallback (hardcoded) values from
There are two ways to get these values:
* first one is using a web-based gas tracker and modeling fallback values on historical gas prices
* second one is using Seth CLI to get the values from the network you are using

Since the second one is much more complex, but will work for any network, we will focus on it.

1. Clone [Seth repository](https://github.com/smartcontractkit/seth)
2. Add required network details for your network to `seth.toml` file:
```toml
[[Networks]]
name = "my_network"
chain_id = "43113"
urls_secret = ["RPC you want to use"]
```
3. Run `SETH_CONFIG_PATH=seth.toml go run cmd/seth/seth.go -n my_network gas -b 10000 -tp 0.99` where `-b <N>` is the number of last blocks that will be used to calculate the fallback gas price and `-tp <N>` is the percentile of the gas prices that will be used as the fallback gas prices.
4. Copy the following part of the output to your network configuration:
```toml
5:08PM INF Fallback prices for TOML config:
gas_price = 121487901046
gas_tip_cap = 1000000000
gas_fee_cap = 122487901046 
```

# Ephemeral and static keys
In order to use Seth effectively you need to understand the difference between ephemeral and static keys, but first, why would you need any of them? Two use cases come to mind:
* parallel operations - you want to run multiple operations in parallel, but you don't want to wait for the previous one to finish (remember, Seth is **synchronous** and requires you to use different keys for each goroutine)
* load generation - you want to generate a lot of transactions in a short period of time

Now, let's get to the keys:
* Ephemeral keys are generated on the fly. They are not stored anywhere and are not meant to be used on live networks. You should use them only, when you are fine with losing all funds associated with them. 
* Static keys are ones that were passed to Seth in the configuration. They should be used on live networks, since funds associated with them are not lost after the test ends, because you have copies of the private keys.

In fact most (but not necessarily all) tests will not let you use ephemeral keys on live networks.

## How to set ephemeral keys in the TOML
It's really simple:
```toml
[Seth]
ephemeral_addresses_number = 10
```

## How to set static keys in the TOML
There are various ways to pass them, but I will show you two:

### As a list of wallet keys in your network configuration
Simply add it to your test TOML:
```toml
[Network.WalletKeys]
arbitrum_sepolia=["first_key","second_key"]
```
This is best to use, when running tests on local, but should be avoided in CI. There another method should be used:

### As Base64-encoded keyfile stored as GHA secret
This is a safer and preferred way, but it involves a bit more work. 

First the prerequisites:
1. Your Seth must be configured to read the keyfile from env var, by setting the following in your TOML:
```toml
[Seth]
keyfile_source = "base64_env"
```
2. Your pipeline needs to set the secret with Base64-encoded keyfile as an env var called `SETH_KEYFILE_BASE64`. Seth will read
and decode it automatically.

Once your pipeline is set up, you can add the keyfile to your GHA secrets. Here's how you can do it:
1. First you need to create a keyfile (check below)
2. Then you need to Base64-encode it and add it to your GHA secrets, you can use GH CLI for that:
```bash
gh secret set SETH_KEYFILE_BASE64 -b $(cat keyfile.toml | base64)
```

## How to split funds between static keys
Now you might ask how to split funds between multiple static keys. Let's say your test requires 50 different static keys to generate
sufficient load. Splitting funds manually would be a nightmare. That's why we have a simple way to do it:
1. Fund some key, let's call it root key, with the desired amount that will be used as source of funds.
2. Use seth to split the funds between all the keys:
```bash
KEYFILE_PATH=keyfile_my_network.toml ROOT_PRIVATE_KEY=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 SETH_CONFIG_PATH=seth.toml go run cmd/seth/seth.go -n my_network  keys split -a 10 -b 1
```
Where `-a <N>` is the number of keys you want to split the funds between and `-b <N>` is the buffer in ether that will be left on the root key.

# Synchronous API
One of the main drivers behind Seth was making its API synchronous. This can be achieved by wrapping all transactions in `seth.Decode()` function that will wait for the transaction to be mined. Then depending on Seth settings
it will decode it either only if it was reverted (default behavior) or always (if `tracing_level` is set to `all`).

# Getting event data from transactions
Getting event data from transaction is not overly complicated, but it requires a couple of steps:
```go

    // first execute the tranasaction and wrap it in Decode function
	tx, err := v.client.Decode(v.coordinator.CancelSubscription(v.client.NewTXOpts(), subID, to)
	if err != nil {
        return nil, errors.Wrap(err, "Error executing transaction")
    }
	
	// now look for the event in the transaction logs
	var event vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled
	for i, e := range tx.Events {
	    if len(e.Topics) == 0 {
	        return nil, fmt.Errorf("no topics in event %d", i)
	    }
	switch e.Topics[0] {
	    case vrf_coordinator_v2.VRFCoordinatorV2SubscriptionCanceled{}.Topic().String():
	        if to, ok := e.EventData["to"].(common.Address); ok {
	            event.To = to
	        } else {
	            return nil, fmt.Errorf("'to' not found in the event")
	        }
	        if amount, ok := e.EventData["amount"].(*big.Int); ok {
	            event.Amount = amount
	        } else {
	            return nil, fmt.Errorf("'amount' not found in the event")
	        }
	        if subId, ok := e.EventData["subId"].(uint64); ok {
	            event.SubId = subId
	        } else {
	            return nil, fmt.Errorf("'subId' not found in the event")
	        }
	    }
	}
    
```