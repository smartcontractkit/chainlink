# Devenv — AI Agent Guidelines

This file documents conventions and constraints for the `devenv` module. Follow these rules when generating, modifying, or reviewing code under `devenv/`.

## Building and Linting

After making large changes to the code, verify the code still builds and lints properly.

```sh
# Build
just build-fakes

# Lint
golangci-lint run ./... --fix
```

## Module Isolation

`devenv` is a standalone Go module: `github.com/smartcontractkit/chainlink/devenv`.

**Critical rule**: NEVER import `github.com/smartcontractkit/chainlink/v2` or any of its children (e.g. `chainlink/v2/core/...`, `chainlink/integration-tests/...`, `chainlink/deployment/...`). This is enforced by `depguard` in `.golangci.yml` and will fail CI.

### Allowed Dependencies

| Dependency                                                          | Use for                                                            |
| ------------------------------------------------------------------- | ------------------------------------------------------------------ |
| `github.com/smartcontractkit/chainlink-testing-framework/framework` | Docker orchestration, `clclient` (CL node HTTP API), observability |
| `github.com/smartcontractkit/chainlink-evm/gethwrappers`            | On-chain contract interaction (deploy, call, transact)             |
| `github.com/smartcontractkit/libocr`                                | OCR-specific contract wrappers                                     |
| `github.com/ethereum/go-ethereum`                                   | ETH client, `bind`, ABI, `common`, `crypto`                        |
| `github.com/stretchr/testify`                                       | Test assertions (`require`, `assert`)                              |
| `github.com/google/uuid`                                            | UUID generation                                                    |
| Standard library                                                    | Everything else                                                    |

### Denied Packages

These are enforced by depguard and will cause lint failures:

| Denied                                                    | Use instead                                 |
| --------------------------------------------------------- | ------------------------------------------- |
| `github.com/smartcontractkit/chainlink/v2` (and children) | Local implementations or CTF equivalents    |
| `github.com/BurntSushi/toml`                              | `github.com/pelletier/go-toml/v2`           |
| `github.com/smartcontractkit/chainlink-integrations/evm`  | `github.com/smartcontractkit/chainlink-evm` |
| `github.com/gofrs/uuid`, `github.com/satori/go.uuid`      | `github.com/google/uuid`                    |
| `github.com/test-go/testify/*`                            | `github.com/stretchr/testify/*`             |
| `go.uber.org/multierr`                                    | `errors.Join` (standard library)            |
| `gopkg.in/guregu/null.v1/v2/v3`                           | `gopkg.in/guregu/null.v4`                   |
| `github.com/go-gorm/gorm`                                 | `github.com/jmoiron/sqlx`                   |

## Product Interface

Every Chainlink product in devenv implements the `Product` interface defined in [interface.go](interface.go). Read that file for the exact method signatures.

### Adding a New Product

1. Create `products/<name>/configuration.go` implementing `Product` — see any existing product (e.g. [products/cron/configuration.go](products/cron/configuration.go)) as a reference
2. Create `products/<name>/basic.toml` with default TOML config
3. Register in [environment.go](environment.go) — add a `case "<name>"` in `newProduct()`
4. Create `tests/<name>/smoke_test.go`
5. Add a matrix entry in [`.github/workflows/devenv-nightly.yml`](../.github/workflows/devenv-nightly.yml)

If migrating an existing test from `integration-tests/smoke/`, follow the full step-by-step process in [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md).

### Product Lifecycle

The environment calls product methods in this order:

1. `Load()` — parse product config from merged TOML
2. `GenerateNodesConfig()` — return CL node TOML overrides
3. `GenerateNodesSecrets()` — return CL node secrets overrides
4. *(infrastructure starts: blockchain, fake server, CL nodes)*
5. `ConfigureJobsAndContracts()` — deploy contracts, create keys/jobs, fund nodes
6. `Store()` — write deployed state (addresses, job IDs) to `env-out.toml`

See [environment.go](environment.go) `NewEnvironment()` for the full orchestration flow.

## Test Conventions

Tests use a two-phase pattern: environment setup (via `cl` CLI) then test execution (via `go test`).

### Test File Structure

See [tests/cron/smoke_test.go](tests/cron/smoke_test.go) for the simplest example of the standard pattern. Every smoke test follows the same structure:

1. Load infrastructure output via `de.LoadOutput[de.Cfg]`
2. Load product output via `products.LoadOutput[<product>.Configurator]`
3. Save container logs in `t.Cleanup`
4. Create clients (ETH and/or CL node)
5. Interact with contracts and assert results

### Key Patterns

- Output file path from `tests/<product>/` is always `../../env-out.toml`
- Use `products.ETHClient()` for Ethereum client creation with gas settings
- Use `products.WaitMinedFast()` for fast transaction confirmation on Anvil
- Use `clclient.New()` for CL node API access (job runs, keys, jobs)
- Use `require.EventuallyWithT` for async assertions (typical: 2 min timeout, 2 s interval)
- Use gethwrappers from `chainlink-evm/gethwrappers` directly for contract bindings

## Configuration

- Base infra: `env.toml` (blockchain, fake server, node set)
- Product config: `products/<name>/basic.toml`
- CLI merges configs left-to-right: `cl u env.toml,products/<name>/basic.toml`
- Output written to `env-out.toml` after environment starts

## Formatting and Linting

- **Formatter**: goimports with local prefix `github.com/smartcontractkit/chainlink`
- **Linter config**: `devenv/.golangci.yml`
- **Run linter**: `golangci-lint run` from the `devenv/` directory
- **nolint directives**: must include both an explanation and a specific linter name

## CI

Tests run via [`.github/workflows/devenv-nightly.yml`](../.github/workflows/devenv-nightly.yml) using a matrix with `envcmd` and `testcmd` pairs. When adding a new test, copy an existing matrix entry and update the product name, test command, and directory fields.
