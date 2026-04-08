---
id: local-cre-environment-workflows
title: Workflow Operations
sidebar_label: Workflows
sidebar_position: 1
---

# Workflow Operations

Local CRE includes CLI helpers for compiling, deploying, deleting, and testing workflows.

## Core Commands

Deploy a workflow:

```bash
go run . workflow deploy -w ./path/to/workflow/main.go --compile -n my_workflow_name
```

Delete a workflow from the registry:

```bash
go run . workflow delete -n my_workflow_name
```

Delete all workflows from the registry:

```bash
go run . workflow delete-all
```

Run the proof-of-reserve example verifier:

```bash
go run . workflow run-por-example
```

## Deploy Flags

The deploy command supports the following implementation-backed flags:

- `--workflow-file-path`
- `--config-file-path`
- `--secrets-file-path`
- `--secrets-output-file-path`
- `--container-target-dir`
- `--container-name-pattern`
- `--rpc-url`
- `--workflow-owner-address`
- `--workflow-registry-address`
- `--capabilities-registry-address`
- `--don-id`
- `--name`
- `--delete-workflow-file`
- `--compile`
- `--with-contracts-version`

`--workflow-file-path` and `--name` are required.

## Workflow Compilation

When you pass `--compile`, Local CRE uses the shared workflow compiler:

- Go workflows run `go mod tidy`
- Go builds use `CGO_ENABLED=0`, `GOOS=wasip1`, and `GOARCH=wasm`
- TypeScript workflows compile through `bun cre-compile`
- the output artifact is compressed into a `.br.b64` file
- workflow names must be at least 10 characters long

These same rules are used by the system-test helpers.

## Workflow Configuration and Secrets

Configuration files are optional and workflow-specific. When you use secrets:

- `--secrets-file-path` points to the unencrypted input mapping
- `--secrets-output-file-path` controls the encrypted output path

This matches the registration path used by the shared workflow package in `system-tests/lib/cre/workflow`.

## Example Workflows

Common example deployment patterns include:

- PoR v2 cron example
- cron-based workflows
- HTTP workflows
- node-mode workflows

Use the examples under `core/scripts/cre/environment/examples/workflows/` as the fastest way to validate an environment after startup.

## Additional Workflow Sources

Local CRE supports both contract-backed and file-backed workflow sources.

Key ideas for this mode:

- you can deploy via contract first and then reuse the compiled artifact
- you can generate metadata for file-source workflows
- you can mix contract and file-backed workflows in the same environment
- file-source workflows can be paused or removed without removing contract workflows

Use this mode when you need to iterate quickly on workflow packaging or test workflows without re-registering each time through the contract path.

## Manual Deployment

For lower-level control:

1. start the environment
2. compile a workflow or reuse an existing `.br.b64`
3. provide config and optional secrets
4. deploy through `workflow deploy`
5. inspect the registry and workflow containers if verification fails

For test-specific deployment behavior, see [Workflows in Tests](../system-tests/workflows-in-tests.md).
