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
  <summary>Install Ginkgo</summary>

  [Ginkgo](https://onsi.github.io/ginkgo/) is the testing framework we use to compile and run our tests. It comes with a lot of handy testing setups and goodies on top of the standard Go testing packages.

  `go install github.com/onsi/ginkgo/v2/ginkgo`
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
make test_smoke_simulated args="-nodes=<number-of-parallel-tests>"
```

You can also run specific tests using `make test_smoke` and a `focus` tag.

```sh
make test_smoke args="-focus=@ocr" # Runs all the ocr smoke tests
make test_smoke args="-focus=@keeper" # Runs all smoke tests for keepers
```

[Check out](https://onsi.github.io/ginkgo/#description-based-filtering) how Ginkgo handles focus and skip tags if you're looking for more precise behavior.

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
