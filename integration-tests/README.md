# Integration Tests

- [Integration Tests](#integration-tests)
  - [Summary](#summary)
  - [Guidelines](#guidelines)
    - [Pre-requisites](#pre-requisites)
      - [Test and node configuration](#test-and-node-configuration)
    - [Run Tests](#run-tests)
      - [Locally](#locally)
        - [All tests in a suite](#all-tests-in-a-suite)
        - [A single test](#a-single-test)
      - [In CI](#in-ci)
      - [In Kubernetes](#in-kubernetes)

## Summary

This directory represent a place for different types of integration and system level tests. It utilizes [Chainlink Testing Framework (CTF)](https://github.com/smartcontractkit/chainlink-testing-framework).

> [!TIP]
> **Testcontainers (Dockerized tests)**
> If you want to have faster, locally running, more stable tests, utilize plain Docker containers (with the help of [Testcontainers](https://golang.testcontainers.org/)) instead of using GitHub Actions or Kubernetes.

## Guidelines

### Pre-requisites

1. [Installed Go](https://go.dev/)
2. For local testing, [Installed Docker](https://www.docker.com/). Consider [increasing resources limits needed by Docker](https://stackoverflow.com/questions/44533319/how-to-assign-more-memory-to-docker-container) as most tests require building several containers for a Decentralized Oracle Network (e.g. OCR requires 6 nodes and DBs, a simulated blockchain, and a mock server).
3. For remote testing, access to Kubernetes cluster/AWS Docker registry (if you are pulling images from private links).
4. Docker image. If there is no image to pull from a registry, you may run tests against a custom build. Run the following command to build the image:

```bash
make build_docker_image image=<your-image-name> tag=<your-tag>
```

Example: `make build_docker_image image=chainlink tag=test-tag`

#### Test and node configuration

1. Setup `.env` file in the root of `integration-tests` directory. See [example.env](./example.env) for how to set test-runner log level (not a node's log level), Slack notifications, and Kubernetes-related settings.

   1. Ensure to **update you environment** with the following commands:
      1. `cd integration-tests`
      2. `source .env`

2. Setup test secrets. See "how-to" details in the [CTF config README](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/config/README.md#test-secrets).

3. Provide test and node configuration (for more details refer to [testconfig README](./testconfig/README.md) and `example.toml` files):
   1. **Defaults** for all products are defined in `./testconfig/<product>/<product>.toml` files.
   2. To **override default values**, create a `./testconfig/overrides.toml` file (yes, in the root of `testconfig`, not a product directory) specifying the values to override by your test (see some examples in [./testconfig/ocr2/overrides](./testconfig/ocr2/overrides)).
    > [!IMPORTANT]
    > **Image version and node configs**
    > 1. Do not forget to set `[ChainlinkImage].version` to test against the necessary remotely accessible version or [custom build](#optional-build-docker-image).
    > 2. When running OCR-related tests, pay attention to which version of OCR you enable/override in your `overrides.toml`.
    > 3. Pay attention to not committing any sensitive data.

4. [Optional] Configure Seth (or use defaults), an evm client used by tests. Detailed instructions on how to configure it can be found in the [Seth README](./README_SETH.md) and [Seth repository](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/seth).

### Run Tests

#### Locally

> [!NOTE]
> **Resources utilization**
> It's recommended to run only one test at a time (run tests sequentially) on a local machine as it needs a lot of docker containers and can peg your resources otherwise. You will see docker containers spin up on your machine for each component of the test where you can inspect logs.

##### All tests in a suite

1. Run CLI command: `go test ./smoke/<product>_test.go`
   Example: `go test ./smoke/ocr_test.go`
2. Logs of each Chainlink container will dump into the `smoke/logs/`.
3. To enable debugging of HTTP and RPC clients set the following env vars:

```bash
export SETH_LOG_LEVEL=debug
export RESTY_DEBUG=true
```

##### A single test

1. Run CLI command: `go test ./smoke/<product>_test.go -run <TestNameToRun>`
   Example: `go test ./smoke/ocr_test.go -run TestOCRBasic`

#### In CI

1. Refer [Tests Run Books](./run-books/) for more details.
2. Logs in CI uploaded as GitHub artifacts.

#### In Kubernetes

Such tests as Soak, Performance, Benchmark, and Chaos Tests remain bound to a Kubernetes run environment.

1. Ensure all necessary configuration is provided (see [Test Configuration](#test-configuration))
2. You are logged in to your Kubernetes cluster (with `aws sso login`)
3. Run CLI command: `make test_<your_test>` (see commands in [Makefile .PHONY lines](./Makefile) for more details)
   Example: `make test_soak_ocr`, `make test_soak_ocr2`, `test_node_migrations`, etc.
4. Navigate to Grafana dashboards to see test and node logs, and results.
