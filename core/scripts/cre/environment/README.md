# Local CRE environment

The long-form Local CRE documentation now lives in [`docs/local-cre/`](../../../../docs/local-cre/index.md).

Use these pages instead of this legacy README:

- [Getting Started](../../../../docs/local-cre/getting-started/index.md)
- [Environment](../../../../docs/local-cre/environment/index.md)
- [Workflow Operations](../../../../docs/local-cre/environment/workflows.md)
- [Topologies and Capabilities](../../../../docs/local-cre/environment/topologies.md)
- [CRE System Tests](../../../../docs/local-cre/system-tests/index.md)

## Quickstart

```bash
cd core/scripts/cre/environment
go run . env start --auto-setup
go run . workflow deploy -w ./examples/workflows/v2/cron/main.go --compile -n cron_example
```

## Contact

Slack: `#topic-local-dev-environments`
