# Manual test: don_family gateway pairing (two workflow DONs)

Use this checklist to verify per-`don_family` gateway pairing end-to-end in Docker.
Unit tests in `system-tests/lib/cre/topology_don_family_gateway_pairing_test.go` and
`core/scripts/cre/environment/environment/workflow_don_family_deploy_test.go` cover resolver
logic without Docker.

## Prerequisites

- Docker running, chainlink repo built (`core/chainlink.Dockerfile` context)
- From `core/scripts/cre/environment`:

```bash
cd core/scripts/cre/environment
```

## 1. Start local CRE with pairing topology

```bash
go run . env start -c configs/workflow-gateway-don-family-pairing.toml
```

**Expect** in startup output:

```text
Gateway don_family pairing enabled: feeds-zone-a → gateway-zone-a (don_family=feeds-zone-a), feeds-zone-b → gateway-zone-b (don_family=feeds-zone-b)
```

**Fail fast checks** if start errors:

- Every workflow and gateway nodeset must set `don_family`
- Each workflow family must have a matching gateway family

## 2. Register a workflow into zone A

Use a compiled workflow artifact and config appropriate for your setup (cron example works):

```bash
go run . env workflow deploy \
  --workflow-file-path /path/to/workflow.wasm.b64 \
  --name manual-test-zone-a \
  --config-file-path /path/to/config.json \
  --workflow-don-name feeds-zone-a \
  --don-family feeds-zone-a \
  --container-name-pattern feeds-zone-a-node
```

**Expect:** deploy succeeds; `--don-family` must match `feeds-zone-a` in state (cross-check rejects typos).

## 3. Register a workflow into zone B

```bash
go run . env workflow deploy \
  --workflow-file-path /path/to/workflow.wasm.b64 \
  --name manual-test-zone-b \
  --config-file-path /path/to/config.json \
  --workflow-don-name feeds-zone-b \
  --don-family feeds-zone-b \
  --container-name-pattern feeds-zone-b-node
```

## 4. Verify family isolation (optional)

- **Gateway connectors:** zone-a workflow nodes should list only `gateway-zone-a` in gateway connector config (not zone-b).
- **Workflow sync:** each workflow DON loads only workflows registered under its `don_family` (check node logs / registry syncer).
- **Negative test:** re-run zone-a deploy with `--don-family feeds-zone-b` — should fail before registration.

## 5. Cleanup

```bash
go run . env stop
```
