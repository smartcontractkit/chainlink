---
id: local-cre-system-tests-workflows-in-tests
title: Workflows in Tests
sidebar_label: Workflows in Tests
sidebar_position: 2
---

# Workflows in Tests

The smoke tests should usually use the shared helpers rather than re-implementing workflow compilation, artifact copying, and registration.

## Recommended Helper

Prefer:

```go
t_helpers.CompileAndDeployWorkflow(...)
```

That helper:

- copies artifacts to workflow DONs by default
- supports additional artifact-copy targets when needed
- creates workflow artifacts
- resolves the workflow registry from the deployed environment
- registers the workflow with the current registry version

Use `WithArtifactCopyDONTypes(...)` when a test intentionally needs artifact copy to additional DON types.

## Choosing the Right Test Environment Helper

The more important authoring choice is usually not workflow compilation. It is which environment helper the test uses.

Use:

```go
t_helpers.SetupTestEnvironmentWithPerTestKeys(...)
```

for workflow-plane tests that may run in parallel or perform independent on-chain writes. This path:

- creates a fresh funded key pair for the test
- swaps the EVM clients and deployer key to use that test-specific signer
- authorizes the signer on the v2 workflow registry when needed
- avoids nonce collisions and shared-key coupling between parallel tests

Use:

```go
t_helpers.SetupTestEnvironmentWithConfig(...)
```

for admin, control-plane, or ownership-sensitive tests that must use the shared root signer. This is the safer choice for flows such as:

- V1 registry tests
- sharding and ownership-admin operations
- tests that intentionally act as the environment owner

As a rule:

- default to `SetupTestEnvironmentWithPerTestKeys(...)` for v2 workflow execution tests
- use `SetupTestEnvironmentWithConfig(...)` only when the test needs shared owner privileges or intentionally avoids per-test signer isolation

## Compilation Rules

The shared compiler in `system-tests/lib/cre/workflow/compile.go` applies these rules:

- workflow names must be at least 10 characters long
- Go workflows run `go mod tidy` before build
- Go workflows compile with `CGO_ENABLED=0`, `GOOS=wasip1`, and `GOARCH=wasm`
- TypeScript workflows compile via `bun cre-compile`
- the final artifact is Brotli-compressed and base64-encoded to `.br.b64`

## Config Files, Secrets, and YAML Workflows

Workflow config files are optional and specific to the workflow under test.

Secrets support is provided by the shared workflow package, which prepares encrypted secrets for registration against the current DON and capabilities registry.

This area also includes:

- workflow secrets handling
- YAML/DSL workflows
- direct JD job proposal flows

Keep those patterns only for tests that truly need lower-level control. The standard smoke path should still prefer `CompileAndDeployWorkflow`.

## Price Data Sources

The PoR-related tests use a shared `PriceProvider` abstraction with two main implementations:

- `TrueUSDPriceProvider`
- `FakePriceProvider`

`TrueUSDPriceProvider` uses the live TrueUSD reserve endpoint and mainly validates that prices become non-zero.

`FakePriceProvider` starts a shared fake HTTP server once, generates a bounded sequence of test prices for each feed, enforces auth headers, and tracks both expected and actual prices for stricter assertions.

Use the fake provider for local and repeatable smoke coverage. Use the live provider only when a scenario intentionally validates the integration path against live data.
