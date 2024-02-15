# CCIP Tests

Here lives the integration tests for ccip, utilizing our [chainlink-testing-framework](https://github.com/smartcontractkit/chainlink-testing-framework) and [integration-tests](https://github.com/smartcontractkit/ccip/tree/ccip-develop/integration-tests)

## Running the tests

### Setting up test inputs

In order to run the tests the first step is to set up the test inputs. Here are the steps to set your test input -
1. Default test input - set via TOML - If no specific input is set; the tests will run with default inputs mentioned in [default.toml](./testconfig/tomls/default.toml). 
Please refer to [default.toml](./testconfig/tomls/default.toml) for the list of parameters that can be set through TOML. We do not recommend changing this file.
2. Overriding the default input - If you want to run your test with different test input, you can override the default inputs. For that you need to set an env var `BASE64_CCIP_CONFIG_OVERRIDE` containing the base64 encoded TOML file content with updated test input parameters. 
    For example, if you want to override the `Network` input in test and want to run your test on `avalanche testnet` and `arbitrum goerli` network, you need to -
   1. create a TOML file with the following content:
        ```toml
        [CCIP]
        [CCIP.Env]
        [CCIP.Env.Network]
        selected_networks= ['AVALANCHE_FUJI', 'ARBITRUM_GOERLI']
        ```
   2. encode it using `base64` command 
   3. set the env var `BASE64_CCIP_CONFIG_OVERRIDE` with the encoded content.
    ```bash
    export BASE64_CCIP_CONFIG_OVERRIDE=$(base64 -i <path-to-override-toml-file>)
    ```

    [mainnet.toml](./testconfig/override/mainnet.toml), [override.toml](./testconfig/override/override.toml), [prod_testnet.toml](./testconfig/override/prod_testnet.toml) are some of the sample override TOML files. 

    For example - In order to run the smoke test (TestSmokeCCIPForBidirectionalLane) on mainnet, run the test with following env var set:
    ```bash
      export BASE64_CCIP_CONFIG_OVERRIDE=$(base64 -i ./testconfig/override/mainnet.toml)
    ```

3. Secrets - You also need to set some secrets. This is a mandatory step needed to run the tests. Please refer to [sample-secrets.toml](./testconfig/tomls/sample-secrets.toml) for the list of secrets that are mandatory to run the tests.
   - The chainlink image and tag are required secrets for all the tests. 
   - If you are running tests in live networks like testnet and mainnet, you need to set the secrets(rpc urls and private keys) for the respective networks.
   - If you are running tests in simulated networks no network specific secrets are required.
   here is a sample secrets.toml file, for running the tests in simulated networks, with the chainlink image and tag set as secrets:
   ```toml
   [CCIP]
   [CCIP.Env]
   # ChainlinkImage is mandatory for all tests.
   [CCIP.Env.Chainlink]
   [CCIP.Env.Chainlink.Common]
   [CCIP.Env.Chainlink.Common.ChainlinkImage]
   image = "chainlink-ccip"
   version = "latest"
   ```

   We consider secrets similar to test input overrides and encode them using `base64` command.
   Once you have the secrets.toml file, you can encode it using `base64` command (similar to step 2) and set the env var `BASE64_CCIP_SECRETS_CONFIG` with the encoded content.
```bash
    export BASE64_CCIP_SECRETS_CONFIG=$(base64 -i ./testconfig/tomls/secrets.toml)
```

**Please note that the secrets should NOT be checked in to the repo and should be kept locally.**
We recommend against changing the content of [sample-secrets.toml](./testconfig/tomls/sample-secrets.toml). Please create a new file and set it as the secrets file.
You can run the command to ignore the changes to the file. 
```bash
git update-index --skip-worktree <path-to-secrets-file>
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
For running the smoke tests - 
1. Set the test input [Setting up test inputs](#setting-up-test-inputs)
    1. If required create override toml with the required test inputs. If you want to run the tests with default parameters, you can skip this step.
    2. Create a TOML file with the secrets.
2. Run the following command to run the smoke tests with your custom override toml and secrets.
```bash
# mark the testimage as empty for running the tests in local docker containers
make test_smoke_ccip testimage="" testname=TestSmokeCCIPForBidirectionalLane override_toml="<the toml file with overridden config string>" secret_toml="<the toml file with secrets string>"
``` 

#### to run the tests with default test inputs
```bash
make test_smoke_ccip_default testname=TestSmokeCCIPForBidirectionalLane secret_toml="<the toml file with secrets string>"
```
Currently other types of tests like load and chaos can only be run using remote kubernetes cluster.

### Using remote kubernetes cluster

These tests remain bound to a Kubernetes run environment, and require more complex setup and running instructions. We endeavor to make these easier to run and configure, but for the time being please seek a member of the QA/Test Tooling team if you want to run these.