# Migrating Smoke Tests to devenv

Step-by-step guide for converting a test from `integration-tests/smoke/` to the `devenv/` pattern. Read [AGENTS.md](AGENTS.md) first for module constraints and conventions.

## Scope

This guide covers migrating a single Chainlink product's smoke tests. If the product already has a configurator in `devenv/products/`, skip to Step 4.

## Pre-work: Analyze the Source Test

Read the old test file in `integration-tests/smoke/<product>_test.go` and identify:

1. **Product name** -- what Chainlink product is being tested (e.g. VRF, Flux, Direct Request)
2. **Contracts deployed** -- which Solidity contracts are deployed during test setup
3. **Jobs and keys created** -- which CL node job types and key types are used
4. **Helper functions** -- any calls to `integration-tests/actions` or `integration-tests/contracts` that need local replacements
5. **Forbidden imports** -- list every import from `github.com/smartcontractkit/chainlink/v2`, `integration-tests/`, or `deployment/` -- all of these must be replaced

For each forbidden import, find the replacement:

| Old import pattern                      | Replacement                                          |
| --------------------------------------- | ---------------------------------------------------- |
| `integration-tests/contracts/<wrapper>` | Direct gethwrapper from `chainlink-evm/gethwrappers` |
| `integration-tests/actions/<helper>`    | Reimplement locally in the configurator              |
| `integration-tests/client`              | `chainlink-testing-framework/framework/clclient`     |
| `chainlink/v2/core/services/...`        | Not needed -- use `clclient` job spec types          |
| `chainlink/v2/core/gethwrappers/...`    | `chainlink-evm/gethwrappers`                         |

## Step 1: Create the Product TOML Config

Create `devenv/products/<name>/basic.toml`.

Reference file: [products/directrequest/basic.toml](products/directrequest/basic.toml) (simplest single-node product).

Key decisions:
- `[[products]]` `name` field must match the switch case you will add in Step 3
- `nodes` count depends on the product (1 for single-node products like cron/DR/VRF, 5 for multi-node like OCR2/flux)
- The product-specific TOML section (e.g. `[[vrf]]`) key must match the `toml` struct tag on the `Configurator.Config` field you create in Step 2
- Include `gas_settings` if the product deploys contracts

## Step 2: Create the Product Configurator

Create `devenv/products/<name>/configuration.go`.

Reference files:
- [products/directrequest/configuration.go](products/directrequest/configuration.go) -- single-node product with contracts, bridge, and job
- [products/vrf/configuration.go](products/vrf/configuration.go) -- single-node product with multiple contracts and local helper functions
- [products/flux/configuration.go](products/flux/configuration.go) -- multi-node product

### Struct definitions

Define three structs:

1. `Configurator` -- top-level, with `Config []*<Product>` and a `toml` tag matching the TOML section key
2. `<Product>` -- product config fields (typically `GasSettings` + `Out`)
3. `Out` -- all deployed state the tests will need (contract addresses, job IDs, key hashes, chain IDs, etc.). Every field must have a `toml` tag.

### Boilerplate methods

`Load()`, `Store()`, `GenerateNodesConfig()`, and `GenerateNodesSecrets()` are nearly identical across products. Copy from any existing product and adjust the type parameter in `products.Load[Configurator]()`.

### ConfigureJobsAndContracts

This is where the real migration work happens. Port the setup logic from the old test file:

1. Connect to CL nodes via `clclient.New(ns[0].Out.CLNodes)`
2. Create ETH client via `products.ETHClient()`
3. Deploy contracts using gethwrappers directly (e.g. `link_token.DeployLinkToken`, then `bind.WaitDeployed`)
4. Use `products.WaitMinedFast()` for transaction confirmations after deployment
5. Fund CL node transmitter addresses via `products.FundAddressEIP1559()`
6. Create keys/bridges/jobs via `clclient` methods (`MustCreateVRFKey`, `MustCreateBridge`, `MustCreateJob`, etc.)
7. Store all outputs the test will need in `m.Config[0].Out`

### Replacing forbidden helpers

When old tests use functions from `integration-tests/actions`:
- Read the source of the helper function
- Reimplement it locally as an unexported function in your `configuration.go`
- Example: VRF migration reimplemented `EncodeOnChainVRFProvingKey` and `EncodeOnChainExternalJobID` as local helpers (see [products/vrf/configuration.go](products/vrf/configuration.go))

### Contract wrappers

Old tests often use wrappers from `integration-tests/contracts` that add convenience methods around gethwrappers. In devenv, use the gethwrappers directly:

- Find the underlying gethwrapper package in `chainlink-evm/gethwrappers/generated/` or `chainlink-evm/gethwrappers/shared/generated/`
- Call `Deploy<Contract>(auth, client, ...)` directly
- Call contract methods on the returned instance directly

## Step 3: Register the Product

Edit [environment.go](environment.go):

1. Add an import for the new package: `"github.com/smartcontractkit/chainlink/devenv/products/<name>"`
2. Add a case to the `newProduct()` switch matching the product name from your TOML config

## Step 4: Create the Smoke Test

Create `devenv/tests/<name>/smoke_test.go`.

Reference files:
- [tests/cron/smoke_test.go](tests/cron/smoke_test.go) -- simplest test (no contracts, just job run polling)
- [tests/vrf/smoke_test.go](tests/vrf/smoke_test.go) -- test with contract interaction, key hash decoding, and job replacement

Every test must:

1. Load infra output: `de.LoadOutput[de.Cfg]("../../env-out.toml")`
2. Load product output: `products.LoadOutput[<product>.Configurator]("../../env-out.toml")`
3. Register cleanup: `t.Cleanup(func() { framework.SaveContainerLogs(...) })`
4. Create clients as needed (`products.ETHClient`, `clclient.New`)
5. Interact with contracts via gethwrappers (never `chainlink/v2` wrappers)
6. Assert with `require.EventuallyWithT` for async results (typical: 2 min timeout, 2 s poll)

The test file should only contain assertion logic. All setup (contract deployment, job creation, funding) belongs in the configurator.

## Step 5: Verify the Build

From the `devenv/` directory, run targeted builds and lints on the new packages only:

```bash
# Overall checks
just build-fakes

# Package specific checks
go build ./products/<name>/...
go build ./tests/<name>/...
golangci-lint run ./products/<name>/... --fix
golangci-lint run ./tests/<name>/... --fix
```

Both grep commands should return no results.

## Step 6: Remove Old Test from CI

Open `.github/e2e-tests.yml` and find the entry matching the old test file path (e.g. `path: integration-tests/smoke/vrf_test.go`). Delete the entire entry block including its `id`, `path`, `test_env_type`, `triggers`, `test_cmd`, and all other fields.

## Step 7: Add New Test to CI

Open `.github/workflows/devenv-nightly.yml` and add a matrix entry in the `matrix.include` array of the `test-nightly` job. Copy an existing entry for a product of similar complexity and update:

- `display_name` -- human-readable name
- `testcmd` -- the `go test` command with `-run` filter matching your test function names
- `envcmd` -- the `cl u` command pointing to your TOML configs
- `runner` -- `ubuntu-latest` for simple tests, larger runners for resource-heavy tests
- `tests_dir` -- the subdirectory name under `devenv/tests/`
- `logs_archive_name` -- name for the uploaded log artifact

If the `testcmd` includes multiple test names separated by `|`, escape the pipe as `\\|` in YAML.

## Step 8: Delete the Old Test File

Delete `integration-tests/smoke/<product>_test.go`. If this was the last test in that file, also check if there are any shared helpers in the same package that are now unused and can be removed.

## Common Pitfalls

1. **Forbidden imports** -- The most common failure. Grep for `chainlink/v2`, `integration-tests`, and `deployment` before committing. The `depguard` linter in `devenv/.golangci.yml` enforces this.
2. **TOML key mismatch** -- The `toml` struct tag on `Configurator.Config` must exactly match the TOML section name (e.g. `toml:"vrf"` matches `[[vrf]]`). A mismatch means the config silently loads as empty.
3. **Missing `toml` tags on `Out` fields** -- Every field in the `Out` struct needs a `toml` tag or it won't be persisted/loaded from `env-out.toml`.
4. **`WaitMinedFast` vs `bind.WaitDeployed`** -- Use `bind.WaitDeployed` for contract deployment transactions (returns the deployed address). Use `products.WaitMinedFast` for all other transactions (state changes, transfers, registrations).
5. **Test imports product package** -- The test in `tests/<name>/` imports the product package from `products/<name>/` only for the `Configurator` type. All contract interaction in tests uses gethwrappers directly.
6. **Output file path** -- Tests run from `devenv/tests/<name>/`, so the output file is `../../env-out.toml` (two levels up to `devenv/`).
7. **Package name** -- The test file's `package` declaration should match the directory name (e.g. `package vrf` in `tests/vrf/`).
