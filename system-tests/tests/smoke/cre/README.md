# Test Modification and Execution Guide

The long-form CRE smoke-test documentation now lives in [`docs/local-cre/system-tests/`](../../../../docs/local-cre/system-tests/index.md).

Start with:

- [CRE System Tests](../../../../docs/local-cre/system-tests/index.md)
- [Running Tests](../../../../docs/local-cre/system-tests/running-tests.md)
- [Workflows in Tests](../../../../docs/local-cre/system-tests/workflows-in-tests.md)
- [CI and Suite Maintenance](../../../../docs/local-cre/system-tests/ci-and-suite-maintenance.md)
- [Local CRE Overview](../../../../docs/local-cre/index.md)

## Quickstart

```bash
cd core/scripts/cre/environment
go run . env setup
go run . env start --with-beholder

go test ./system-tests/tests/smoke/cre -timeout 20m -run '^Test_CRE_'
```

## Rule of Thumb

Happy-path and sanity checks belong in `smoke`; edge cases and negative conditions belong in `regression`.
