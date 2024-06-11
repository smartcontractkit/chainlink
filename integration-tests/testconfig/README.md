# TOML is the Ultimate Choice!

## Introduction

Final implementation has undergone minor adjustments in comparison to the approach by Adam Hamric, Anindita Ghosh, and Sergey Kudasov stated in the ADR. The primary changes are as follows:

* `TEST_LOG_LEVEL` remains an environment variable, pending the release of version 2.
* `TEST_TYPE` is also kept as an environment variable to facilitate dynamic configuration selection by some tests.
* TOML configuration of Chainlink nodes themselves has not been added, awaiting version 2.
* The hierarchy of configuration overrides has been streamlined for simplicity.

By design, all test configurations are intended to reside within the `testconfig` package, organized into application-specific folders. However, the system can locate these configurations in any folder within the `integration-tests` directory, selecting the first one found. To identify the configurations in use, execute tests with the `debug` log level.

The `testconfig` package serves as a centralized resource for accessing configurations across all products, including shared settings like logging and network preferences, as well as initial funding for Chainlink nodes. Product configurations, if present, are subjected to validation based on logical assumptions and observed code values. The `TestConfig` structure includes a `Save()` method, allowing for the preservation of test configurations after all adjustments have been applied.

## Configuration and Overrides

The order of precedence for overrides is as follows:

* [File `default.toml`](#defaulttoml)
* [Product-specific file, e.g., `[product_name].toml`](#product-specific-configurations)
* [File `overrides.toml`](#overridestoml)
* [Environment variable `BASE64_CONFIG_OVERRIDE`](#base64_config_override)

### default.toml
That file is envisioned to contain fundamental and universally applicable settings, such as logging configurations, private Ethereum network settings or Seth networks settings for known networks.

### Product-specific configurations
Product-specific configurations, such as those in `[product_name].toml`, house the bulk of default and variant settings, supporting default configurations like the following in `log_poller.toml`, which should be used by all Log Poller tests:

```toml
# product defaults
[LogPoller]
[LogPoller.General]
generator = "looped"
contracts = 2
events_per_tx = 4
use_finality_tag = true
log_poll_interval = "500ms"
# 0 disables backup poller
backup_log_poller_block_delay = 0

[LogPoller.Looped]
execution_count = 100
min_emit_wait_time_ms = 200
max_emit_wait_time_ms = 500
```

### overrides.toml
This file is recommended for local use to adjust dynamic variables or modify predefined settings. At the very minimum it should contain the Chainlink image and version, as shown in the example below:

```toml
[ChainlinkImage]
image = "your image name"
version = "your tag"
```

### `BASE64_CONFIG_OVERRIDE`
This environment variable is primarily intended for use in continuous integration environments, enabling the substitution of default settings with confidential or user-specific parameters. For instance:

```bash
cat << EOF > config.toml
[Network]
selected_networks=["$SELECTED_NETWORKS"]

[ChainlinkImage]
image="$CHAINLINK_IMAGE"
version="$CHAINLINK_VERSION"
postgres_version="$CHAINLINK_POSTGRES_VERSION"

[Pyroscope]
enabled=$pyroscope_enabled
server_url="$PYROSCOPE_SERVER"
environment="$PYROSCOPE_ENVIRONMENT"
key_secret="$PYROSCOPE_KEY"

[Logging]
test_log_collect=false
run_id="$RUN_ID"

[Logging.LogStream]
log_targets=["$LOG_TARGETS"]

[Logging.Loki]
tenant_id="$LOKI_TENANT_ID"
endpoint="$LOKI_ENDPOINT"
basic_auth_secret="$LOKI_BASIC_AUTH"

[Logging.Grafana]
base_url="$GRAFANA_URL"
dashboard_url="$GRAFANA_DASHBOARD_URL"
EOF

BASE64_CONFIG_OVERRIDE=$(cat config.toml | base64 -w 0)
echo ::add-mask::$BASE64_CONFIG_OVERRIDE
echo "BASE64_CONFIG_OVERRIDE=$BASE64_CONFIG_OVERRIDE" >> $GITHUB_ENV
```
**It is highly recommended to use reusable GHA actions present in [.actions](../../../.github/.actions) to generate and apply the base64-encoded configuration.** Own implementation of `BASE64_CONFIG_OVERRIDE` generation is discouraged and should be used only if existing actions do not cover the use case. But even in that case it might be a better idea to extend existing actions.
This variable is automatically relayed to Kubernetes-based tests, eliminating the need for manual intervention in test scripts.


## Named Configurations

Named configurations allow for the customization of settings through unique identifiers, such as a test name or type, acting as specific overrides. Here's how you can define and use these configurations:

For instance, to tailor configurations for a particular test, you might define it as follows:

```toml
# Here the configuration name is "TestLogManyFiltersPollerFinalityTag"
[TestLogManyFiltersPollerFinalityTag.LogPoller.General]
contracts = 300
```

Alternatively, for a configuration that applies to a certain type of test, as seen in `vrfv2.toml`, you could specify:

```toml
# Here the configuration name is "Soak"
[Soak.VRFv2.Common]
cancel_subs_after_test_run = true
```

When processing TOML files, the system initially searches for a general (unnamed) configuration. If a named configuration is found, it can specifically override the general (unnamed) settings, providing a targeted approach to configuration management based on distinct identifiers like test names or types.

### Chainlink Node TOML config

Find default node config in `testconfig/default.toml`

To set custom config for Chainlink Node use `NodeConfig.BaseConfigTOML` in TOML. Example:
```toml
[NodeConfig]
BaseConfigTOML = """
[Feature]
FeedsManager = true
LogPoller = true
UICSAKeys = true

[Log]
Level = 'debug'
JSONConsole = true

[Log.File]
MaxSize = '0b'

[OCR]
Enabled = true
DefaultTransactionQueueDepth = 0
"""
```
Note that you cannot override individual values in BaseConfigTOML. You must provide the entire configuration.


To set base config for EVM chains use `NodeConfig.CommonChainConfigTOML`. Example:
```toml
CommonChainConfigTOML = """
AutoCreateKey = true
FinalityDepth = 1
MinContractPayment = 0

[GasEstimator]
PriceMax = '200 gwei'
LimitDefault = 6000000
FeeCapDefault = '200 gwei'
"""
```

This is the default configuration used for all EVM chains unless ChainConfigTOMLByChainID is specified.

To set custom per-chain config use `[NodeConfig.ChainConfigTOMLByChainID]`. Example:
```toml
[NodeConfig.ChainConfigTOMLByChainID]
# applicable for arbitrum-goerli chain
421613 = """
[GasEstimator]
PriceMax = '400 gwei'
LimitDefault = 100000000
FeeCapDefault = '200 gwei'
BumpThreshold = 60
BumpPercent = 20
BumpMin = '100 gwei'
"""
```

For more examples see `example.toml` in product TOML configs like `testconfig/automation/example.toml`.

### Setting env vars for Chainlink Node

To set env vars for Chainlink Node use `WithCLNodeOptions()` and `WithNodeEnvVars()` when building a test environment. Example:

```go
envs := map[string]string{
    "CL_LOOPP_HOSTNAME": "hostname",
}
testEnv, err := test_env.NewCLTestEnvBuilder().
    WithTestInstance(t).
    WithTestConfig(&config).
    WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
    WithMockAdapter().
    WithCLNodes(clNodeCount).
    WithCLNodeOptions(test_env.WithNodeEnvVars(envs)).
    WithFunding(big.NewFloat(.1)).
    WithStandardCleanup().
    WithSeth().
    Build()
```

## Local/Kubernetes Usage

GitHub workflows in this repository have been updated to dynamically generate and utilize base64-encoded TOML configurations derived from user inputs or environment variables. For local execution or remote Kubernetes runners, users must manually supply certain variables, which cannot be embedded in configuration files due to their sensitive or dynamic nature.

Essential variables might include:

* Chainlink image and version
* Test duration for specific tests (e.g., load, soak)
* Configuration specific to Loki (mandatory for certain tests)
* Grafana dashboard URLs

For local testing, it is advisable to place these variables in the `overrides.toml` file. For Kubernetes or remote runners, the process involves creating a TOML file with the necessary values, encoding it in base64, and setting the result as the `BASE64_CONFIG_OVERRIDE` environment variable.

## Embedded config

Because Go automatically excludes TOML files during the compilation of binaries, we must take deliberate steps to include our configuration files in the compiled binary. This can be accomplished by using a custom build tag `-o embed`. Implementing this tag will incorporate all the default configurations located in the `./testconfig` folder directly into the binary. Therefore, when executing tests from the binary, you'll only need to supply the `overrides.toml` file. This file should list only the settings you wish to modify; all other configurations will be sourced from the embedded configurations. You can access these embedded configurations [here](.integration-tests/testconfig/configs_embed.go).

## To bear in mind

### Validation failures

When the system encounters even a single setting related to a specific product or configuration within the configurations, it triggers a comprehensive validation of the entire configuration for that product. This approach is based on the assumption that if any configuration for a given product is specified, the entire set of configurations for that product must be complete and valid. This is particularly crucial when dealing with the `overrides.toml` file, where it's easy to overlook the need to comment out or adjust values when switching between configurations for different products. Essentially, the presence of any specific configuration detail necessitates that all relevant configurations for that product be fully defined and correct to prevent validation errors.

## Possible nil pointers

If no configuration values are set for a product or its logging parameters, the system won't perform validation checks. This can lead to a 'nil pointer exception' error if you attempt to access a configuration property later on. This situation arises because we use pointers to facilitate optional overrides; accessing an unset (nil) pointer will cause an error. To avoid such issues, especially when general validations might not cover every scenario, it's crucial for users to ensure that all necessary configuration options are explicitly set. Additionally, it's highly recommended to implement test-specific validations to confirm that all required values for a particular test are indeed established. This proactive approach helps prevent runtime errors and ensures smooth test execution.

## Contributing

It's crucial to incorporate all new test configuration settings directly into the TOML configuration files, steering clear of using environment variables for this purpose. Our goal is to centralize all configuration details, including examples, within the same package. This approach simplifies the process of understanding the available configuration options and identifying the appropriate values to use for each setting.

## Reusing TestConfig in other projects

To ensure the cleanliness and simplicity of your project's configuration, it's advised against using the `testconfig` code as a direct library in other projects. The reason is that much of this code is tailored specifically to its current application, which might not align with the requirements of your project. Your project might not necessitate any overrides or could perhaps benefit from a simpler configuration approach.

However, if you find a need to utilize some methods from this project, the recommended practice is to implement the required interfaces within your project's configuration package, rather than directly copying and pasting code. For instance, if you aim to incorporate a setup action similar to the `SetupVRFV2Environment` for VRFv2, like the one shown below:

```go
func SetupVRFV2Environment(
    env *test_env.CLClusterTestEnv,
    nodesToCreate []vrfcommon.VRFNodeType,
    vrfv2TestConfig types.VRFv2TestConfig,
    useVRFOwner bool,
    useTestCoordinator bool,
    linkToken contracts.LinkToken,
    mockNativeLINKFeed contracts.MockETHLINKFeed,
    registerProvingKeyAgainstAddress string,
    numberOfTxKeysToCreate int,
    numberOfConsumers int,
    numberOfSubToCreate int,
    l zerolog.Logger,
) (*vrfcommon.VRFContracts, []uint64, *vrfcommon.VRFKeyData, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, error) {
```

You should not replicate the entire `TestConfig` structure. Instead, create an implementation of the `types.VRFv2TestConfig` interface in your project and use that as the parameter. This approach allows you to maintain a streamlined and focused configuration package in your project.

## Known Issues/Limitations

* Duplicate file names in different locations may lead to unpredictable configurations being selected.
* The use of pointer fields for optional configuration elements necessitates careful handling, especially for programmatic modifications, to avoid unintended consequences. The `MustCopy()` function is recommended for creating deep copies of configurations for isolated modifications. Unfortunately some of the custom types are not copied at all, you need to set them manually. It's true for example for `blockchain.StrDuration` type.
