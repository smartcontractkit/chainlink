# Integration Tests

Here lives the integration tests for chainlink, utilizing our [integrations-framework](https://github.com/smartcontractkit/integrations-framework).

## How to Run

### Connect to a Kubernetes cluster

Integration tests require a connection to an actively running kubernetes cluster. [Minikube](https://minikube.sigs.k8s.io/docs/start/)
can work fine for some tests, but in order to run more rigorous tests, or to run with any parallelism, you'll need to either
increase minikube's resources significantly, or get a more substantial cluster.
This is necessary to deploy ephemeral testing environments, which include external adapters, chainlink nodes and their DBs,
as well as some simulated blockchains, all depending on the types of tests and networks being used.

### Running

Our suggested way to run these tests is to use [the ginkgo cli](https://onsi.github.io/ginkgo/#the-ginkgo-cli).

The default for this repo is the utilize the Makefile.

```sh
make test_smoke
```

In order to run in **parallel**, utilize args.

```sh
make test_smoke args="-nodes=6"
```

The above will run tests with 6 parallel threads.

### Options

If you would like to change the Chainlink values that are used for environments, you can use JSON to squash them. Have a look over at our [helmenv](https://github.com/smartcontractkit/helmenv/) chainlink charts to get a grasp of how things are structured. We'll be writing more on this later, but for now, you can squash values by providing a `CHARTS` environment variable.

```sh
CHARTS='{"chainlink": {"values": {"chainlink": {"image": {"version": "1.0.1"}}}}}' make test_smoke args="-nodes=6"
```
