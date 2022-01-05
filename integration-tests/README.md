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

If you would like to change the Chainlink or Geth versions that are used for environments, you can use the below ENV vars. 

```sh
CHAINLINK_IMAGE=my/chainlink/image/location
CHAINLINK_VERSION=1.0.0
GETH_IMAGE=ethereum/client-go
GETH_VERSION=v1.10.15
```

If you want more fine grained control, have a look over at our [helmenv](https://github.com/smartcontractkit/helmenv/) chainlink charts to get a grasp of how things are structured. We're working on improving the UX of this system, but for now, make use of the `CHARTS` environment variable.

```sh
CHARTS={"chainlink":"values":{"env":{"feature_offchain_reporting":"true"}}}
```
