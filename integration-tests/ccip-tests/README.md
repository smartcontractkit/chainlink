# CCIP Tests

Here lives the integration tests for ccip, utilizing our [chainlink-testing-framework](https://github.com/smartcontractkit/chainlink-testing-framework) and [integration-tests](https://github.com/smartcontractkit/ccip/tree/ccip-develop/integration-tests)

## Running the tests

### Setting up test inputs :

In order to run the tests the first step is to set up the test inputs. There are two kinds of inputs -
1. Generic test input - set via TOML - If no specific input is set; the tests will run with default inputs mentioned in [default.toml](./testconfig/tomls/default.toml)
2. Secrets - set via env variables. Please refer to [secrets.toml](./testconfig/secrets.env) for the list of env variables that need to be set.

If you want to override the default inputs, you need to set an env var `BASE64_TEST_CONFIG_OVERRIDE` containing the base64 encoded TOML file content.
For example, if you want to override the `Networks` input in test and want to run your test on `avalanche testnet` and `arbitrum goerli` network, you can create a TOML file with the following content:
```toml
[CCIP]

[CCIP.Env]
Networks = ['AVALANCHE_FUJI', 'ARBITRUM_GOERLI']
```
and then encode it using `base64` command and set the env var `BASE64_TEST_CONFIG_OVERRIDE` with the encoded content.
```bash
export BASE64_TEST_CONFIG_OVERRIDE=$(base64 -i <path-to-toml-file>)
```

Alternatively, you can also use the make command to invoke a go script to do the same.
```bash
make override_config override_toml="<the toml file with overridden config>" env="<.env file with BASE64_TEST_CONFIG_OVERRIDE value>"
```

In order to set the secrets, you need to set the env vars mentioned in [secrets.toml](./testconfig/secrets.env) file and source the file.  
```bash
source ./testconfig/secrets.env
```

Please note that the secrets.env should not be checked in to the repo and should be kept locally.
You can run the command to ignore the changes to the file.
```bash
git update-index --skip-worktree ./testconfig/secrets.env
```

### Triggering the tests
There are two ways to run the tests:
1. Using local docker containers
2. Using remote kubernetes cluster

### Using local docker containers

In order to run the tests locally, you need to have docker installed and running on your machine.
You can use a specific chainlink image and tag (if you already have one) for the tests. Otherwise, you can build the image using the following command:
```bash
make build_ccip_image image=chainlink-ccip tag=latest-dev # please choose the image and tag name as per your choice
```

Currently, for local run the tests creates two new private geth networks and runs the tests on them. Running tests on testnet and mainnet is not supported yet for local run.
Please refer to [Using remote kubernetes cluster](#using-remote-kubernetes-cluster) section for running the tests on live networks like testnet and mainnet.

You can use the following command to run the tests locally with your specific chainlink image.

#### Smoke Tests
```bash
# mark the testimage as empty for running the tests in local docker containers
make test_smoke_ccip image=chainlink-ccip tag=latest-dev testimage="" testname=TestSmokeCCIPForBidirectionalLane override_toml="<the toml file with overridden config string>" env="<.env file with BASE64_TEST_CONFIG_OVERRIDE value>"
# to run the tests with default test inputs
make test_smoke_ccip_default image=chainlink-ccip tag=latest testname=TestSmokeCCIPForBidirectionalLane
```
Currently other types of tests like load and chaos can only be run using remote kubernetes cluster.

### Using remote kubernetes cluster

These tests remain bound to a Kubernetes run environment, and require more complex setup and running instructions. We endeavor to make these easier to run and configure, but for the time being please seek a member of the QA/Test Tooling team if you want to run these.