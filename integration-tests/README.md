# Integration Tests

- [Integration Tests](#integration-tests)
  - [Summary](#summary)
  - [Guidelines](#guidelines)
    - [Pre-requisites](#pre-requisites)
      - [Test and node configuration](#test-and-node-configuration)
    - [Run Tests](#run-tests)
      - [Locally (in Docker)](#locally-in-docker)
        - [All tests in a suite](#all-tests-in-a-suite)
        - [A single test](#a-single-test)
      - [In Kubernetes](#in-kubernetes)
        - [From local machine](#from-local-machine)
      - [CI/GitHub Actions](#cigithub-actions)

## Summary

This directory represent a place for different types of integration and system level tests. It utilizes [Chainlink Testing Framework (CTF)](https://github.com/smartcontractkit/chainlink-testing-framework).

> [!TIP]
> **Testcontainers (Dockerized tests)**
> If you want to have faster, locally running, more stable tests, utilize plain Docker containers (with the help of [Testcontainers](https://golang.testcontainers.org/)) instead of using GitHub Actions or Kubernetes.

## Guidelines

### Pre-requisites

1. [Installed Go](https://go.dev/)
2. For local testing, [Installed Docker](https://www.docker.com/). Consider [increasing resources limits needed by Docker](https://stackoverflow.com/questions/44533319/how-to-assign-more-memory-to-docker-container) as most tests require building several containers for a Decentralized Oracle Network (e.g. OCR requires 6 nodes, 6 DBs, and a mock server).
3. For remote testing, access to Kubernetes cluster/AWS Docker registry (if you are pulling images from private links).
4. Docker image. If there is no image to pull from a registry, you may run tests against a custom build. Run the following command to build the image:

    ```bash
    make build_docker_image image=<your-image-name> tag=<your-tag>
    ```

    Example: `make build_docker_image image=chainlink tag=test-tag`

5. RPC node/s (for testnets/mainnets).
6. EOA's (wallet) Private Key (see [How to export an account's private key](https://support.metamask.io/ru/managing-my-wallet/secret-recovery-phrase-and-private-keys/how-to-export-an-accounts-private-key/))
7. Sufficient amount of native token and LINK on EOA per a target chain.

#### Test and node configuration

1. Setup `.env` file in the root of `integration-tests` directory. See [example.env](./example.env) for how to set test-runner log level (not a node's log level), Slack notifications, and Kubernetes-related settings.

   1. Ensure to **update you environment** with the following commands:
      1. `cd integration-tests`
      2. `source .env`

2. Setup test secrets. See "how-to" details in the [Test Secrets in CTF](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/config/README.md#test-secrets). If you want to run tests in CI, you will have to push test secrets to GitHub (see [Run GitHub Workflow with your test secrets](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/config/README.md#run-github-workflow-with-your-test-secrets)).

3. Provide test and node configuration (for more details refer to [testconfig README](./testconfig/README.md) and `example.toml` files):
   1. **Defaults** for all products are defined in `./testconfig/<product>/<product>.toml` files.
   2. To **override default values**, create a `./testconfig/overrides.toml` file (yes, in the root of `testconfig`, not a product directory) specifying the values to override by your test (see some examples in [./testconfig/ocr2/overrides](./testconfig/ocr2/overrides)).

   > [!IMPORTANT]
   > **Image version and node configs**
   > 1. Pay attention to the `[ChainlinkImage].version` to test against the necessary remotely accessible version or [custom build](#pre-requisites).
   > 2. When running OCR-related tests, pay attention to which version of OCR you enable/override in your `overrides.toml`.
   > 3. Do not commit any sensitive data.

4. [Optional] Configure Seth (or use defaults), an evm client used by tests. Detailed instructions on how to configure it can be found in the [Seth README](./README_SETH.md) and [Seth repository](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/seth).

   > [!IMPORTANT]
   > **Simulated mode (no test secrets needed)**
   > Tests may run in a simulated mode, on a simulated chain (1337). In the `overrides.toml` file, set the following:
   > 1. `[Network].selected_networks=["simulated"]`
   > 2. `[[Seth.networks]].name = "Default"`

### Run Tests

#### Locally (in Docker)

> [!NOTE]
> **Resources utilization by Docker**
> It's recommended to run only one test at a time (run tests sequentially) on a local machine as it needs a lot of docker containers and can peg your resources otherwise. You will see docker containers spin up on your machine for each component of the test where you can inspect logs.

##### All tests in a suite

1. Run CLI command(with `override.toml`):

   ```bash
   BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) go test -v -p 1 ./smoke/<product>_test.go
   ```

   Example:

   ```bash
   BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) go test -v -p 1 ./smoke/ocr_test.go
   ```

   > [!WARNING]
   > **Parallelized tests and nonce issues**
   > Most tests are paralelized by default. To avoid nonce-related issues, it is recommended to run tests with disabled parallelization, e.g. with `-p 1`.

2. Alternatively, you may use `make` commands (see more in [Makefile .PHONY lines](./Makefile)) for running suites of tests.
    Example:

    ```bash
    make test_smoke_product product="ocr" ./scripts/run_product_tests
    ```

3. Logs of each Chainlink container will dump into the `smoke/logs/`.
4. To enable debugging of HTTP and RPC clients set the following env vars:

    ```bash
    export SETH_LOG_LEVEL=debug
    export RESTY_DEBUG=true
    ```

##### A single test

Run CLI command (with `override.toml`):

```bash
BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) go test -v -timeout 15m -run <"TestNameToRun"> ./<directory_name_with_tests>
```

Example:

```bash
BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) go test -v -timeout 15m -run "TestOCRv2Basic" ./smoke
```

#### In Kubernetes

Such tests as Soak, Performance, Benchmark, and Chaos Tests remain bound to a Kubernetes run environment.

1. Refer [Tests Run Books](./run-books/) to get more details on how to run specific per-product tests.
2. Logs in CI are uploaded as GitHub artifacts.

##### From local machine

1. Ensure all necessary configurations are provided (see [Test and node configuration](#test-and-node-configuration)).
2. Log in to your Kubernetes cluster (with `aws sso login`)
3. Run tests with the following CLI command:

   ```bash
   BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) go test -v -timeout <max_test_timeout> -p 1 -run '<TestName>' ./<test_directory>
   ```

   OR with `make` commands:

   ```bash
   BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) make test_<your_test>
   ```

   Example (see make-commands in [Makefile .PHONY lines](./Makefile)):

   ```bash
   BASE64_CONFIG_OVERRIDE=$(cat ./testconfig/overrides.toml | base64) make test_chaos_ocr/make test_soak_ocr2/test_node_migrations
   ```

4. Use Kubernetes namespace printed out in logs to monitor and analyze test runs.
5. Navigate to Grafana dashboards to for test results and logs.

#### CI/GitHub Actions

1. Ensure all necessary configurations are provided (see [Test and node configuration](#test-and-node-configuration)).
2. Follow instructions provided in [E2E Tests on GitHub CI](../.github/E2E_TESTS_ON_GITHUB_CI.md).
3. Refer [Tests Run Books](./run-books/) to get more details on how to run specific per-product tests.
