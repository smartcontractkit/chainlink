# Integration Tests

Here lives the integration tests for chainlink, utilizing our [integrations-framework](https://github.com/smartcontractkit/integrations-framework).

## How to Run

### Connect to a Kubernetes cluster

Integration tests require a connection to an actively running kubernetes cluster. [Minikube](https://minikube.sigs.k8s.io/docs/start/)
can work fine for some tests, but in order to run more rigorous tests, or to run with any parallelism, you'll need to either
increase minikube's resources significantly, or get a more substantial cluster.
This is necessary to deploy ephemeral testing environments, which include external adapters, chainlink nodes and their DBs,
as well as some simulated blockchains, all depending on the types of tests and networks being used.

## Install Ginkgo

[Ginkgo](https://onsi.github.io/ginkgo/) is the testing framework we use to compile and run our tests. It comes with a lot of handy testing setups and goodies on top of the standard Go testing packages.

`go install github.com/onsi/ginkgo/v2/ginkgo`

### Running

The default for this repo is the utilize the Makefile.

```sh
make test_smoke
```

In order to run in **parallel**, utilize args.

```sh
make test_smoke args="-nodes=6"
```

The above will run tests with 6 parallel threads.

## Chainlink Values

If you would like to change the Chainlink values that are used for environments, you can use the `framework.yaml` file,
or set environment variables that are all caps versions of the values found in the config file.

```yaml
# Specify the image and version of the chainlink image you want to run tests against. Leave blank for default.
chainlink_image:      # Image of chainlink node
chainlink_version:    # Version of the image on the chainlink node
chainlink_env_values: # Environment values to pass onto the chainlink nodes
```
