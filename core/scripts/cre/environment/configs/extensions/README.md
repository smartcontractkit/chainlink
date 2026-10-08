# Local CRE extensions

Services run next to Local CRE from their own repos' compose files, listed in `compose.yaml`. Each joins the
`ctf` network and answers to role aliases that `workflow-gateway-don-extensions.toml` points the workflow nodes at.

| Role (alias on `ctf`) | Port | Workflow node config |
| --- | --- | --- |
| `cre-workflow-source` | 50052 | `Capabilities.WorkflowRegistry.AdditionalSources` |
| `cre-artifact-api` | 50052 | `Capabilities.WorkflowRegistry.WorkflowStorage.URL` |
| `cre-artifact-store` | 4566 | `WorkflowStorage.ArtifactStorageHost`; start with `-e 4566` so the gateway can fetch from it |
| `cre-deploy-api` | 50051 | deploy tooling (`registryclient`) |

One service per role. Inputs (`CRE_LOCAL_*`) describe this environment and default to `extensions.env`; set them in
the shell to override. In this topology workflows come from the extensions, not `env workflow deploy`; deploy tooling
runs in a container on `ctf`, since artifact URLs use `cre-artifact-store`.

```bash
CTF_CONFIGS=configs/workflow-gateway-don-extensions.toml go run . env start -e 4566
export CRE_STORAGE_COMPOSE=<path to cre-storage-service>/local-stack/compose.yaml
docker compose -f configs/extensions/compose.yaml up -d --wait
docker restart $(docker ps --format '{{.Names}}' | grep '^workflow-node')
```
