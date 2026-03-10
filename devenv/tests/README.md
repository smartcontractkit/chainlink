# System Level Smoke Tests

Read and setup [devenv](../README.md).

Spin up observability stack in case you need performance tests:
```bash
obs up -f
```

To run any test, open two terminals and setup corresponding commands, `envcmd` and `testcmd` fields from [here](ttps://github.com/smartcontractkit/chainlink/blob/develop/.github/workflows/devenv-nightly.yml#L45):
```bash
$envcmd (from devenv dir)
$testcmd (from tests dir)
```