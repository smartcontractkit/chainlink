---
"chainlink": patch
---

Require explicit `don_family` on workflow and gateway nodesets for all local CRE gateway topologies (N >= 1 families). Gateway connectors, gateway worker jobs, capabilities registry families, and `env workflow deploy` are scoped by family. Deploy resolves the target workflow DON via `--don-family` (with optional `--shard-index`) or `--workflow-don-name`, defaults the docker cp pattern to `{nodesets.name}-node`, and falls back to `DefaultDONFamily` when local CRE state is absent. #internal
