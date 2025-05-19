## Overview

This package is used to create ephemeral environment for local/CI testing. 
It sets up an environment with local Docker containers running Chainlink nodes and a job distributor. 
It can either create new simulated private Ethereum network containers or connect to existing testnets/mainnets.

### Run Tests with Devenv

The tests created with this environment are run as [end-to-end integration smoke tests](../../integration-tests/smoke).

Pre-requisites:
- Docker
- Pull access for chainlink and job-distributor images

#### Setting Up Testconfig with Simulated Private Ethereum Network

To run tests (e.g., [ccip-test](../../../integration-tests/ccip-tests/smoke/ccip_test.go)), 
you need to set up the testconfig following the [testconfig setup instructions](../../../integration-tests/testconfig/README.md). 
The testconfig specifies the details of the different configurations to set up the environment and run the tests. 
Generally, tests are run with the [default](../../../integration-tests/testconfig/default.toml) config unless overridden by product-specific config. 
For example, the [ccip-test](../../../integration-tests/ccip-tests/smoke/ccip_test.go) uses [ccip.toml](../../../integration-tests/testconfig/ccip/ccip.toml) to specify 
CCIP-specific test environment details.

There are additional secret configuration parameters required by the `devenv` environment that are not stored in the testconfig. 
These are read from environment variables. For example, Chainlink image, Job-Distributor image, etc. 
All such environment variables are listed in the [sample.env](.sample.env) file. 
You can create a `.env` file in the same directory of the test and set the required environment variables.

#### Setting Up Testconfig for Running Tests with Existing Testnet/Mainnet

To run tests with existing testnet/mainnet, set up the testconfig with the details of the testnet/mainnet networks. 
Following the integration-test [testconfig framework](../../integration-tests/testconfig/README.md#configuration-and-overrides), 
create a new `overrides.toml` file with testnet/mainnet network details and place it under any location in the
`integration-tests` directory. 

By default, tests are run with private Ethereum network containers set up in the same Docker network as 
the Chainlink nodes and job distributor. To run tests against existing testnet/mainnet, 
set the `selected_network` field in the testconfig with the specific network names.

For example, if running [ccip-smoke](../../../integration-tests/ccip-tests/smoke/ccip_test.go) tests with Sepolia, Avax, and Binance testnets, 
copy the contents of [sepolia_avax_binance.toml](../../integration-tests/testconfig/ccip/overrides/sepolia_avax_binance.toml) 
to the `overrides.toml` file.

Before running the test, ensure that RPC and wallet secrets are set as environment variables. 
Refer to the environment variables pattern in the [sample.env](.sample.env) file to 
provide necessary secrets applicable to the network you are running the tests against:
- `E2E_TEST_<networkName>_WALLET_KEY_<sequence_number>`
- `E2E_TEST_<networkName>_RPC_HTTP_URL_<sequence_number>`
- `E2E_TEST_<networkName>_RPC_WS_URL_<sequence_number>`

Now you are all set to run the tests with the existing testnet/mainnet.
