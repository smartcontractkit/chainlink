# Local CRE extensions

Centralized CRE services run from their own repos' compose files. With an `[extensions]` section (`extensions.toml`),
`env start` starts them on `ctf` after the DONs, passing the registry chain, its CapabilitiesRegistry and the workflow
DON family, and points the workflow nodes at the roles they answer to. `env stop` removes them.

| Role (alias on `ctf`) | Port | Workflow node config |
| --- | --- | --- |
| `cre-workflow-source` | 50052 | `Capabilities.WorkflowRegistry.AdditionalSources` |
| `cre-artifact-api` | 50052 | `Capabilities.WorkflowRegistry.WorkflowStorage.URL` |
| `cre-artifact-store` | 4566 | `WorkflowStorage.ArtifactStorageHost`, allowed for the gateway |
| `cre-deploy-api` | 50051 | deploy tooling (`registryclient`) |

`storage_compose` defaults to a `cre-storage-service` checkout next to this repo; its images are built locally (see its
`local-stack/README.md`). Workflows then come from the extensions, not `env workflow deploy`; deploy tooling runs in a
container on `ctf`, since artifact URLs use `cre-artifact-store`.

```bash
CTF_CONFIGS=configs/workflow-gateway-don.toml,configs/extensions/extensions.toml go run . env start
```
