# githooks

`githooks` is a Go CLI tool for running repository hooks (e.g. via [Lefthook](../../lefthook.yml)) efficiently across Go monorepo modules and packages.

## Features

- **Module & Package Resolution:** Automatically maps changed or staged Go files to their enclosing `go.mod` module roots and specific package paths (e.g. `./core/logger`).
- **Targeted Linting:** Runs `golangci-lint` only against the exact changed packages within affected modules instead of scanning whole modules or the entire repository.
- **Targeted Unit Testing:** Discovers changed test packages and executes `tools/test` with `-short` directly on those packages.
- **Dependency Changes:** Automatically runs on all packages (`./...`) if a module's `go.mod` or `go.sum` is modified.
- **Lefthook Integration:** Seamlessly works with Lefthook staged/push file filters and `stage_fixed`.

## Commands

### `lint`

Runs `golangci-lint` only on the changed packages within affected modules.

```bash
# Lint changed packages from staged files
go -C tools/githooks run . lint

# Lint specific files / packages across modules
go -C tools/githooks run . lint core/services/app.go deployment/environment/env.go

# Flags
#   --fix        Fix issues automatically where possible (default: true)
#   --rev string Show issues in modified lines since rev (default: "HEAD")
go -C tools/githooks run . lint --fix=false --rev=origin/develop
```

### `test` (aliases: `only-changed`, `short-test`)

Runs unit tests (`tools/test`) in `-short` mode only on affected packages.

```bash
# Run short tests on changed packages from git diff
go -C tools/githooks run . test

# Run short tests on specific changed files
go -C tools/githooks run . test core/logger/logger_test.go deployment/environment.go

# Flags
#   --short      Run tests in -short mode (default: true)
go -C tools/githooks run . test --short=false
```

## Running Tests

```bash
cd tools/githooks
go test ./...
golangci-lint run
```
