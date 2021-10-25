# Integration Tests

Here lives the integration tests for chainlink, utilizing our [integrations-framework](https://github.com/smartcontractkit/integrations-framework).

## How to Run

### Connect to a Kubernetes cluster

Integration tests require a connection to an actively running kubernetes cluster. [Minikube](https://minikube.sigs.k8s.io/docs/start/)
can work fine for some tests, but in order to run more rigorous tests, or to run with any parallelism, you'll need to either
increase minikube's resources signigicantly, or get a more substantial cluster.
This is necessary to deploy ephemeral testing environments, which include external adapters, chainlink nodes and their DBs,
as well as some simulated blockchains, all depending on the types of tests and networks being used.

### Running

Our suggested way to run these tests is to use [the ginkgo cli](https://onsi.github.io/ginkgo/#the-ginkgo-cli).

```sh
ginkgo ./integration-tests/integration
```

Some defaults are set up in the `integration_tests.sh` file. Those are ideal for CI runs, and might need some adjustments
if you're running locally.

### Options

There are some standard ginkgo CLI arguments we like to use, along with some settings and environment variables you might like to change. Here are the significant ones, see the `integration-tests/config.yml` file for all of them.

| ENV Var                 | Description                                                 | Default                            |
|-------------------------|-------------------------------------------------------------|------------------------------------|
|`NETWORKS`               | Comma seperated list of blockchain networks to run tests on | ethereum_geth,ethereum_geth        |
|`APPS_CHAINLINK_IMAGE`   | Image location for a valid docker image of a chainlink node | public.ecr.aws/chainlink/chainlink |
|`APPS_CHAINLINK_VERSION` | Version to be used for the above mentioned image            | 1.0.0                              |
