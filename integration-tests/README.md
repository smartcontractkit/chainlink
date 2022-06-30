# Integration Tests

Here lives the integration tests for chainlink, utilizing our [chainlink-testing-framework](https://github.com/smartcontractkit/chainlink-testing-framework).

## How to Configure

See the [example.env](./example.env) file to set some common environment variables to configure how tests are run.

### Connect to a Kubernetes Cluster

Integration tests require a connection to an actively running kubernetes cluster. [Minikube](https://minikube.sigs.k8s.io/docs/start/)
can work fine for some tests, but in order to run more rigorous tests, or to run with any parallelism, you'll need to either
increase minikube's resources significantly, or get a more substantial cluster.
This is necessary to deploy ephemeral testing environments, which include external adapters, chainlink nodes and their DBs,
as well as some simulated blockchains, all depending on the types of tests and networks being used.

### Install Ginkgo

[Ginkgo](https://onsi.github.io/ginkgo/) is the testing framework we use to compile and run our tests. It comes with a lot of handy testing setups and goodies on top of the standard Go testing packages.

`go install github.com/onsi/ginkgo/v2/ginkgo`

## How to Run

Most of the time, you'll want to run tests on a simulated chain, for the purposes of speed and cost.

### Smoke

Run all smoke tests, only using simulated blockchains.

```sh
make test_smoke_simulated
```

Run all smoke tests in parallel, only using simulated blockchains.

```sh
make test_smoke_simulated args="-nodes=<number-of-parallel-tests>"
```

You can also run specific tests or specific networks using a `focus` tag.

```sh
make test_smoke args="-focus=@metis" # Runs all the smoke tests on the Metis Stardust network
make test_smoke args="-focus=@ocr" # Runs all OCR tests
```

### Soak

See the [soak_runner_test.go](./soak/soak_runner_test.go) file to trigger soak tests.

### Performance

Currently, all performance tests are only run on simulated blockchains.

```sh
make test_perf
```

## Common Issues

After running `go mod tidy` or similar commands, many seem to hit this error:

```plain
github.com/smartcontractkit/chainlink/integration-tests imports
	github.com/smartcontractkit/chainlink-testing-framework/blockchain imports
	github.com/ethereum/go-ethereum/crypto imports
	github.com/btcsuite/btcd/btcec/v2/ecdsa tested by
	github.com/btcsuite/btcd/btcec/v2/ecdsa.test imports
	github.com/btcsuite/btcd/chaincfg/chainhash: ambiguous import: found package github.com/btcsuite/btcd/chaincfg/chainhash in multiple modules:
	github.com/btcsuite/btcd v0.22.0-beta (/Users/adamhamrick/go/pkg/mod/github.com/btcsuite/btcd@v0.22.0-beta/chaincfg/chainhash)
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.1 (/Users/adamhamrick/go/pkg/mod/github.com/btcsuite/btcd/chaincfg/chainhash@v1.0.1)
```

A quick workaround is to run `go get -u github.com/btcsuite/btcd/chaincfg/chainhash@v1.0.1` then `go mod tidy` to resolve it.
