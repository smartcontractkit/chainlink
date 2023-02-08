# Integration Tests

Here lives the integration tests for chainlink, utilizing our [chainlink-testing-framework](https://github.com/smartcontractkit/chainlink-testing-framework).

## Setup

Prerequisites to run the tests.

### Install Dependencies

<details>
  <summary>Install Go</summary>

  [Install](https://go.dev/doc/install)
</details>

<details>
  <summary>Install NodeJS</summary>

  [Install](https://nodejs.org/en/download/)
</details>

<details>
  <summary>Install Helm Charts</summary>

  [Install Helm](https://helm.sh/docs/intro/install/#through-package-managers) if you don't already have it. Then add necessary charts with the below commands.

  ```sh
  helm repo add chainlink-qa https://raw.githubusercontent.com/smartcontractkit/qa-charts/gh-pages/
  helm repo add bitnami https://charts.bitnami.com/bitnami
  helm repo update
  ```

</details>

## Connect to a Kubernetes Cluster

Integration tests require a connection to an actively running kubernetes cluster. [Minikube](https://minikube.sigs.k8s.io/docs/start/)
can work fine for some tests, but in order to run more rigorous tests, or to run with any parallelism, you'll need to either
increase minikube's resources significantly, or get a more substantial cluster.
This is necessary to deploy ephemeral testing environments, which include external adapters, chainlink nodes and their DBs,
as well as some simulated blockchains, all depending on the types of tests and networks being used.

## Configure Environment

See the [example.env](./example.env) file and use it as a template for your own `.env` file. This allows you to configure general settings like what name to associate with your tests, and which Chainlink version to use when running them.

You can also specify `EVM_KEYS` and `EVM_URLS` for running on live chains, or use specific identifiers as shown in the [example.env](./example.env) file.

Other `EVM_*` variables are retrieved when running with the `@general` tag, and is helpful for doing quick sanity checks on new chains or when tweaking variables.

**The tests will not automatically load your .env file. Remember to run `source .env` for changes to take effect.**

## How to Run

Most of the time, you'll want to run tests on a simulated chain, for the purposes of speed and cost.

### Smoke

Run all smoke tests with the below command. Will use your `SELECTED_NETWORKS` env var for which network to run on.

```sh
make test_smoke # Run all smoke tests on the chosen SELECTED_NETWORKS
SELECTED_NETWORKS="GOERLI" make test_smoke # Run all smoke tests on GOERLI network
make test_smoke_simulated # Run all smoke tests on a simulated network
```

Run all smoke tests in parallel, only using simulated blockchains. *Note: As of now, you can only run tests in parallel on simulated chains, not on live ones. Running on parallel tests on live chains will give errors*

```sh
make test_smoke_simulated args="-test.parallel=<number-of-parallel-tests>"
```

You can also run specific tests and debug tests in vscode by setting up your .vscode/settings.json with this information. Just replace all the "<put your ...>" with your information before running a test.

```json
{
    "makefile.extensionOutputFolder": "./.vscode",
    "go.testEnvVars": {
        "LOG_LEVEL": "debug",
        "SELECTED_NETWORKS": "SIMULATED,SIMULATED_1,SIMULATED_2",
        "CHAINLINK_IMAGE":"<put your account number here>.dkr.ecr.us-west-2.amazonaws.com/chainlink",
        "CHAINLINK_VERSION":"develop",
        "CHAINLINK_ENV_USER":"<put your name>",
        "TEST_LOG_LEVEL":"debug",
        "AWS_ACCESS_KEY_ID":"<put your access key id here>",
        "AWS_SECRET_ACCESS_KEY":"<put your access key here>",
        "AWS_SESSION_TOKEN":"<put your token here>"
    },
    "go.testTimeout": "900s"
}
```

You can also run your tests inside of kubernetes instead of from locally to reduce local resource usage and the number of ports that get forwarded to the cluster. This is not recommended for normal developement since building and pushing the image can be time heavy depending on your internet upload speeds. To do this you will want to either pull down an already built chainlink-tests image or build one yourself. To build and push one yourself you can run:

```sh
make build_test_image tag=<a tag for your image> base_tag=latest suite="smoke soak chaos reorg migration performance" push=true
```

Once that is done building you can add this to your go.testEnvVars in .vscode/settings.json with the correct account number and tag filled out.

```json
  "TEST_SUITE": "smoke",
  "TEST_ARGS": "-test.timeout 30m",
  "ENV_JOB_IMAGE":"<account number>.dkr.ecr.us-west-2.amazonaws.com/chainlink-env-tests:<tag you used in the build step>",
```

Once that is done you can run/debug your test using the vscode test view just like normal.

### Soak

Currently we have 2 soak tests, both can be triggered using make commands.

```sh
make test_soak_ocr
make test_soak_keeper
```

Soak tests will pull all their network information from the env vars that you can set in the `.env` file. *Reminder to run `source .env` for changes to take effect.*

To configure specific parameters of how the soak tests run (e.g. test length, number of contracts), see the [./soak/tests](./soak/tests/) test specifications.

See the [soak_runner](./soak/soak_runner_test.go) for more info on how the tests are run and configured.

### Performance

Currently, all performance tests are only run on simulated blockchains.

```sh
make test_perf
```

## Common Issues

When upgrading to a new version, it's possible the helm charts have changed. There are a myriad of errors that can result from this, so it's best to just try running `helm repo update` when encountering an error you're unsure of.
