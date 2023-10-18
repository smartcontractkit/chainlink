# Integration Tests

Here lives the integration tests for chainlink, utilizing our [chainlink-testing-framework](https://github.com/smartcontractkit/chainlink-testing-framework).

## NOTE: Move to Testcontainers

If you have previously run these smoke tests using GitHub Actions or some sort of Kubernetes setup, that method is no longer necessary. We have moved the majority of our tests to utilize plain Docker containers (with the help of [Testcontainers](https://golang.testcontainers.org/)). This should make tests faster, more stable, and enable you to run them on your local machine without much hassle.

## Requirements

1. [Go](https://go.dev/)
2. [Docker](https://www.docker.com/)
3. You'll probably want to [increase the resources available to Docker](https://stackoverflow.com/questions/44533319/how-to-assign-more-memory-to-docker-container) as most tests require quite a few containers (e.g. OCR requires 6 Chainlink nodes, 6 databases, a simulated blockchain, and a mock server).

## Configure

See the [example.env](./example.env) file for environment variables you can set to configure things like network settings, Chainlink version, and log level. Remember to use `source .env` to activate your settings.

## Build

If you'd like to run the tests on a local build of Chainlink, you can point to your own docker image, or build a fresh one with `make`.

`make build_docker_image image=<image-name> tag=<tag>`

e.g.

`make build_docker_image image=chainlink tag=test-tag`

## Run

`go test ./smoke/<product>_test.go`

It's generally recommended to run only one test at a time on a local machine as it needs a lot of docker containers and can peg your resources otherwise. You will see docker containers spin up on your machine for each component of the test where you can inspect logs.

## Analyze

You can see the results of each test in the terminal with normal `go test` output. If a test fails, logs of each Chainlink container will dump into the `smoke/logs/` folder for later analysis. You can also see these logs in CI uploaded as GitHub artifacts.

## Running Soak, Performance, Benchmark, and Chaos Tests

These tests remain bound to a Kubernetes run environment, and require more complex setup and running instructions not documented here. We endeavor to make these easier to run and configure, but for the time being please seek a member of the QA/Test Tooling team if you want to run these.
