---
id: local-cre-getting-started-index
title: Getting Started
sidebar_label: Getting Started
sidebar_position: 0
---

# Getting Started

This page gives you the shortest path from a clean checkout to a running Local CRE environment and a local smoke-test run.

## Prerequisites

For the standard Docker-based flow, Local CRE expects:

- Docker installed and running
- access to required images through AWS SSO or direct repository access
- `gh` authenticated if you build images that need private plugin access
- `go` available locally

The existing environment setup flow is still the canonical prerequisite bootstrap:

```bash
cd core/scripts/cre/environment
go run . env setup
```

`env setup` is driven by `configs/setup.toml` and can also be run with:

```bash
go run . env setup --config configs/setup.toml --no-prompt
```

If you need to rebuild or repull prerequisites, `env setup` also supports `--purge`. Billing assets can be included with `--with-billing`.

## Quickstart

Start the default Local CRE stack:

```bash
cd core/scripts/cre/environment
go run . env start --auto-setup
```

Deploy a first workflow:

```bash
go run . workflow deploy -w ./examples/workflows/v2/cron/main.go --compile -n cron_example
```

The environment command writes Local CRE state to the repo-local state file, which is what the test helpers later consume.

## Common Startup Variants

Start with the example workflow:

```bash
go run . env start --with-example
```

Start with Beholder:

```bash
go run . env start --with-beholder
```

Start against v1 contracts:

```bash
go run . env start --with-contracts-version v1
```

Set extra gateway ports when your workflow needs outbound access to local services:

```bash
go run . env start --extra-allowed-gateway-ports 8080,8171
```

## First Smoke Test Run

Once the environment is up, run the CRE smoke package:

```bash
go test ./system-tests/tests/smoke/cre -timeout 20m -run '^Test_CRE_'
```

For the default smoke-test flow, start Local CRE without `--with-beholder`. Most tests start the ChIP test sink on the default gRPC port (`50051`), and Beholder uses that same port through Chip Ingress.

Enable Beholder only when:

- you are running Beholder-specific tests
- you intentionally need the Beholder stack for debugging
- you move Beholder to a different port with `--grpc-port` so it does not conflict with the test sink

The smoke tests default to the capability-enabled topology when you do not override `CTF_CONFIGS`, and the test helpers can start Local CRE automatically if the state file does not exist yet.

Continue with:

- [Environment](../environment/index.md) for lifecycle and debugging
- [Running Tests](../system-tests/running-tests.md) for local, Kubernetes, and CI execution modes
