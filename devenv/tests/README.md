# System Level Smoke Tests

Read and setup [devenv](../README.md).

Spin up observability stack in case you need performance tests:
```bash
obs up -f
```

To run any test, open two terminals and setup corresponding commands, `envcmd` and `testcmd` fields from [devenv-nightly](/.github/workflows/devenv-nightly.yml):
```bash
$envcmd (from devenv dir)
$testcmd (from tests/$product dir, for example "tests/automation")
```

## Dashboards

- [All Logs](http://localhost:3000/explore?schemaVersion=1&panes=%7B%22axv%22:%7B%22datasource%22:%22P8E80F9AEF21F6940%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22expr%22:%22%7Bjob%3D%5C%22ctf%5C%22%7D%22,%22queryType%22:%22range%22,%22datasource%22:%7B%22type%22:%22loki%22,%22uid%22:%22P8E80F9AEF21F6940%22%7D,%22editorMode%22:%22code%22,%22direction%22:%22backward%22%7D%5D,%22range%22:%7B%22from%22:%22now-30m%22,%22to%22:%22now%22%7D,%22compact%22:false%7D%7D&orgId=1)
- [CL Node Errors](http://localhost:3000/d/a7de535b-3e0f-4066-bed7-d505b6ec9ef1/cl-node-errors?orgId=1&refresh=5s&from=now-15m&to=now&timezone=browser)
- [Load Testing](http://localhost:3000/d/WASPLoadTests/wasp-load-test?orgId=1&from=now-30m&to=now&timezone=browser&var-go_test_name=$__all&var-gen_name=$__all&var-branch=$__all&var-commit=$__all&var-call_group=$__all&refresh=5s)