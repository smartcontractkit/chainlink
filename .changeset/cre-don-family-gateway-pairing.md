---
"chainlink": patch
---

Require explicit `don_family` on every nodeset for local CRE topologies. Gateway connectors, gateway worker jobs, capabilities registry families, and `env workflow deploy` are scoped by family. Deploy resolves the target workflow DON via `--don-family` (with optional `--shard-index`) or `--workflow-don-name`, and requires `--don-family` when local CRE state is absent. #internal
