# Test Configurations

- [Test Configurations](#test-configurations)
  - [Summary](#summary)
  - [Configurations Files and Overrides Precedence](#configurations-files-and-overrides-precedence)
    - [Test Secrets (optional)](#test-secrets-optional)
    - [default.toml](#defaulttoml)
    - [\<product\>.toml](#producttoml)
    - [overrides.toml (optional)](#overridestoml-optional)
    - [`BASE64_CONFIG_OVERRIDE`](#base64_config_override)
  - [Node configurations](#node-configurations)
    - [Spec Properties](#spec-properties)
      - [BaseConfigTOML](#baseconfigtoml)
    - [Network configurations](#network-configurations)
      - [Spec Properties](#spec-properties-1)
        - [CommonChainConfigTOML](#commonchainconfigtoml)
        - [ChainConfigTOMLByChainID](#chainconfigtomlbychainid)
    - [Programmatic configuration](#programmatic-configuration)
    - [Embedded configurations](#embedded-configurations)
  - [Test type and case specific configurations](#test-type-and-case-specific-configurations)
  - [Product-specific configurations](#product-specific-configurations)
    - [Migration tests](#migration-tests)
    - [Automation](#automation)
      - [Specific test secrets](#specific-test-secrets)
    - [OCR](#ocr)
      - [Common OCR configurations](#common-ocr-configurations)
      - [Reuse OCR contracts](#reuse-ocr-contracts)
  - [Worthy to note](#worthy-to-note)
  - [Reusing `testconfig` in other projects](#reusing-testconfig-in-other-projects)

> [!NOTE]
> **Current state and v2**
>
> - There are still several configuration files required to run tests: `.env`, `default.toml, <per-product>.toml`, `testsecrets`.
> - `TEST_LOG_LEVEL` is still kept as an environment variable to facilitate dynamic configuration selection by some tests.
> - Configuration of Chainlink nodes with TOML files separated from test configuration is aimed at v2.

> [!IMPORTANT]
> **Contributing**
>
> It's crucial to incorporate all new test configuration settings directly into the TOML configuration files, steering clear of using environment variables for this purpose. Our goal is to centralize all configuration details, including examples, within the same package. This approach simplifies the process of understanding the available configuration options and identifying the appropriate values to use for each setting.

## Summary

The `testconfig` package represents a storage of test configurations per product/service.

> [!TIP]
> 1. The `testconfig.go` may `Save()` your specification to a configuration file.
> 2. To identify the configurations in use, run tests with the `TEST_LOG_LEVEL=debug`.

## Configurations Files and Overrides Precedence

Overrides work in the order of the following precedence (1 overrides 2 and so on):

1. [Environment variable `BASE64_CONFIG_OVERRIDE`](#base64_config_override)
2. [overrides.toml](#overridestoml-optional)
3. [\<product\>.toml](#producttoml)
4. [default.toml](#defaulttoml)
5. [testsecrets](#test-secrets-optional)

### Test Secrets (optional)

Test secrets are necessary for remote environments (CI, Kubernetes). For more details, visit [Test Secrets in CTF](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/config/README.md#test-secrets).

### default.toml

[default.toml](default.toml) represents common default test settings for logging, network, node and chain client.

> [!TIP]
> It is recommended to provide [product-specific configurations](#producttoml) with explicitly defined values to override default configs.

### \<product\>.toml

Product-specific default configurations stored in `./testconfig/<product_directory>/<product>.toml` (e.g. [ocr2.toml](./ocr2/ocr2.toml)). They explicitly define a [node configuration](#node-configurations), [test type-specific configurations](#test-type-and-case-specific-configurations), and override any [default settings](#defaulttoml) per product/service.

### overrides.toml (optional)

> [!CAUTION]
> Even though `overrides.toml` is git-ignored, pay attention to avoid storing any sensitive data in this file, especially in override-files for remote environments (which are not git-ignored).

The `overrides.toml` enables tests parametrization. When provided, it overrides default and per-product configurations.

1. **For local runs**, store overrides as follows `./testconfig/overrides.toml`.
2. **For remote environments** (CI, Kubernetes), commit overrides to repository under `../integration-tests/testconfig/<product>/overrides/<override-name>.toml` (e.g. [OCR2 overrides](../integration-tests/testconfig/ocr2/overrides/base_sepolia.toml)). See more in [E2E Tests on GitHub CI with overrides](../../.github/E2E_TESTS_ON_GITHUB_CI.md#test-workflows-setup-in-ci).

Alternatively, set `E2E_TEST_CHAINLINK_IMAGE` and `E2E_TEST_CHAINLINK_VERSION` in `~/.testsecrets`

### `BASE64_CONFIG_OVERRIDE`

This environment variable is used for overriding defaults in remote testing ([CI - On Demand Workflows](../../.github/E2E_TESTS_ON_GITHUB_CI.md#on-demand-workflows) and Kubernetes) when triggered from local machine.

Example:

```bash
# Go test
BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) go test <test args>

- - - - - - - - - - - - -

# Make command
BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) make test_<test>
```

## Node configurations

A node configuration consists of two main blocks:

- Network/chain configuration
- Node-specific configuration

### Spec Properties

#### BaseConfigTOML

A node's configuration unrelated to network settings is defined in `[NodeConfig].BaseConfigTOML="""your_config_here"""`, if none - defaults are used.

Example:

```toml
[NodeConfig]
BaseConfigTOML = """
[Feature]
FeedsManager = true
LogPoller = true
UICSAKeys = true

[OCR]
Enabled = true
"""
```

### Network configurations

Chain-specific configuration is composed of the following blocks:

- `[Network].selected_networks` - a list of networks to run tests on.
- `ChainConfigTOMLByChainID` (if an entry with matching chain id is defined) OR `CommonChainConfigTOML` (if no entry with matching chain id is defined).

> [!NOTE]
>
> 1. If a `ChainConfigTOMLByChainID` or `CommonChainConfigTOML` is specified, they will override any defaults a Chainlink Node might have for the given network.
> 2. To override default [BaseConfigTOML](#baseconfigtoml) and/or [CommonChainConfigTOML](#commonchainconfigtoml) provide the entire blocks as if it would be completely new configuration.

#### Spec Properties

##### CommonChainConfigTOML

A network-specific node config for EVM chains: `[NodeConfig].CommonChainConfigTOML="""your_config_here"""`.

Example:

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

##### ChainConfigTOMLByChainID

The `[NodeConfig.ChainConfigTOMLByChainID]` is a custom per-chain config that overrides the [`CommonChainConfigTOML`](#commonchainconfigtoml). See examples in product directories, e.g. [ocr2/example.toml](./ocr2/example.toml).

Example:

```toml
[NodeConfig.ChainConfigTOMLByChainID]
# applicable only to arbitrum-goerli chain
421613 = """
[GasEstimator]
PriceMax = '400 gwei'
LimitDefault = 100000000
FeeCapDefault = '200 gwei'
"""
```

> [!NOTE]
> Currently, all networks are treated as EVM networks. There's no way to provide Solana, Starknet, Cosmos or Aptos configuration yet.

### Programmatic configuration

To set env vars for a Dockerized test use `docker.test_env_builder.WithCLNodeOptions(test_env.WithNodeEnvVars(envs))`.

Example:

```go
envs := map[string]string{
    "CL_LOOPP_HOSTNAME": "hostname",
}

nodeEnvVars := test_env.WithNodeEnvVars(envs)

testEnv, err := test_env.NewCLTestEnvBuilder().
    WithTestInstance(t).
    WithTestConfig(&config).
    WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
    WithMockAdapter().
    WithCLNodes(clNodeCount).
    WithCLNodeOptions(nodeEnvVars).
    WithFunding(big.NewFloat(.1)).
    WithStandardCleanup().
    WithSeth().
    Build()
```

### Embedded configurations

For the reason Go automatically excludes TOML files during the compilation of binaries, the [config_embed.go](./configs_embed.go) deliberately incorporates all the default configurations into the the compiled binary with a custom build tag `-o embed`. Hence, only `overrides.toml` may be potentially needed to execute tests.

## Test type and case specific configurations

These configurations represent unique identifiers used for customization of different test runs per product. When found, they take precedence and overrides the general (unnamed) settings (mentioned above).

Examples:

```toml
# Specific per test configuration "TestLogPollerManyFiltersFinalityTag" for LogPoller
[TestLogPollerManyFiltersFinalityTag.LogPoller.General]
contracts = 300

# "Soak" test configuration for VRFv2
[Soak.VRFv2.Common]
cancel_subs_after_test_run = true

# "Load" test configuration for OCR2
[Load.OCR2]
[Load.OCR2.Common]
eth_funds = 3
```

## Product-specific configurations

> [!TIP]
> Find which configurations are applicable to a specific product in structs of the `testconfig/product directory/config.go (or <product>.go)`.
> Examples: [ocr2/ocr2.go](./ocr2/ocr2.go), [automation/config.go](./automation/config.go).

### Migration tests

1. Set `E2E_TEST_CHAINLINK_UPGRADE_IMAGE` in [testsecrets](#test-secrets-optional).

   ```bash
   E2E_TEST_CHAINLINK_UPGRADE_IMAGE="public.ecr.aws/chainlink/chainlink"
   ```

2. In `overrides.toml` set image to upgrade to:

    ```toml
    # image to upgrade to
    [ChainlinkUpgradeImage]
    version="2.17.0-beta0"
    ```

3. Run tests:

    ```bash
    # Go
    BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) go test -v -p 1 ./smoke/<product>_test.go

    - - - - - - - - - - - - 

    # Make command
    BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) make test_node_migrations_verbose
    ```

### Automation

#### Specific test secrets

| Secret                        | Env Var                                                             | Example                                      | Description                                                                          |
| ----------------------------- | ------------------------------------------------------------------- | -------------------------------------------- | ------------------------------------------------------------------------------------ |
| Data Streams Url              | `E2E_TEST_DATA_STREAMS_URL`                                         | `E2E_TEST_DATA_STREAMS_URL=url`              | Required by some automation tests to connect to data streams.                         |
| Data Streams Username         | `E2E_TEST_DATA_STREAMS_USERNAME`                                    | `E2E_TEST_DATA_STREAMS_USERNAME=username`    | Required by some automation tests to connect to data streams.    |
| Data Streams Password         | `E2E_TEST_DATA_STREAMS_PASSWORD`                                    | `E2E_TEST_DATA_STREAMS_PASSWORD=password`    | Required by some automation tests to connect to data streams. |

### OCR

#### Common OCR configurations

Specify number of contracts to be deployed for OCR (correspondingly, the same amount of jobs per OCR contract will be created on a node). For example:

```toml
# OCRv1
[OCR.Common]
number_of_contracts=2

- - - - - - - - - - - -

# OCRv2
[OCR2.Common]
number_of_contracts=2
```

#### Reuse OCR contracts

The feature supports OCR v1 and v2. It gets enabled when `[OCR.Contract]` block is specified.

To reuse existing OCR contracts provide:

- LINK token address (unique per chain)
- OCR contract addresses
- [optional] choose, whether to use and configure the existing OCR contracts. **N/B:** usage/configuring of several selected contracts (not all the listed addresses) is not supported. All the contacts should have the same `use` and `configure` values.

Example:

```toml
# For OCRv1
[OCR.Contracts]
link_token = "0x88d1239894D9582f5849E5b5a964da9e5730f1E6"
offchain_aggregators = ["0xc1ce3815d6e7f3705265c2577F1342344752A5Eb"]

# If [OCR.Contracts.Settings.<OCR aggregator address>] is not present, we assume it should be used and configured

- - - - - - - - - - - - 

# For OCRv2
[OCR2.Contracts]
link_token = "0x88d1239894D9582f5849E5b5a964da9e5730f1E6"
offchain_aggregators = ["0xc1ce3815d6e7f3705265c2577F1342344752A5Eb"]

# notice that this address needs to match the one in offchain_aggregators
[OCR2.Contracts.Settings."0xc1ce3815d6e7f3705265c2577F1342344752A5Eb"]
use = false # Default: true. Reuse existing OCR contracts?
configure = false # Default: true. Configure existing OCR contracts?

- - - - - - - - - - - - 

# Non-compliant configurations
[OCR.Contracts]
link_token = "0x88d1239894D9582f5849E5b5a964da9e5730f1E6"
offchain_aggregators = ["0xc1ce3815d6e7f3705265c2577F1342344752A5Eb", "0x2f4FA21fCd917C448C160caafEC874032F404c08"]

# Example 1: Setting `configure` to `false` for selected (not all) addresses will fail configuration validation
[OCR.Contracts.Settings."0xc1ce3815d6e7f3705265c2577F1342344752A5Eb"]
configure = false

# OR

# Example 2: Setting `configure` to different values for the listed contracts will fail execution
[OCR.Contracts.Settings."0xc1ce3815d6e7f3705265c2577F1342344752A5Eb"]
configure = false

[OCR.Contracts.Settings."0x2f4FA21fCd917C448C160caafEC874032F404c08"]
configure = true

# OR

# Example 3: Setting `use` to different values for the listed contracts will fail execution
[OCR.Contracts.Settings."0xc1ce3815d6e7f3705265c2577F1342344752A5Eb"]
use = false

[OCR.Contracts.Settings."0x2f4FA21fCd917C448C160caafEC874032F404c08"]
use = true
```

## Worthy to note

> [!NOTE]
> **Configuration Validation and `nil` pointers**
> When tests encounter a [test](#test-type-and-case-specific-configurations) or [product-specific](#product-specific-configurations) setting, they trigger a validation to ensure that the entire set of configurations for that product is complete and valid.
>
> If there are no configuration values (for a product or its logging parameters), the tests won't perform validation checks, leading to a `nil pointer exception` error on an attempt to access a configuration property later on. The error is caused by the usage of pointers to facilitate optional overrides: accessing an unset (nil) pointer will cause an error. To avoid such run-time issues, it is highly recommended to implement configuration-specific validations to confirm that all the required values for a particular test are explicitly specified.

> [!NOTE]
> **Known Issues/Limitations**
> Duplicated test configuration file names in different locations may lead to an unpredictable execution behavior.
>
> The use of pointer fields for optional configuration elements necessitates careful handling, especially for programmatic modifications, to avoid unintended consequences. The `MustCopy()` function is recommended for creating deep copies of configurations for isolated modifications. Unfortunately some of the custom types are not copied at all, you need to set them manually. It's true for example for `blockchain.StrDuration` type.

## Reusing `testconfig` in other projects

To utilize some methods from this project, implement the required interfaces within your project's configuration package (no copy-pasting or structure replicating).

Example:

```go
func SetupVRFV2Environment(
    env *test_env.CLClusterTestEnv,
    nodesToCreate []vrfcommon.VRFNodeType,
    vrfv2TestConfig types.VRFv2TestConfig, // implementation of interface used as a parameter
    useVRFOwner bool,
    useTestCoordinator bool,
    linkToken contracts.LinkToken,
    mockNativeLINKFeed contracts.MockETHLINKFeed,
    registerProvingKeyAgainstAddress string,
    numberOfTxKeysToCreate int,
    numberOfConsumers int,
    numberOfSubToCreate int,
    l zerolog.Logger,
) (*vrfcommon.VRFContracts, []uint64, *vrfcommon.VRFKeyData, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, error) { <setup logic> }
```
