# githooks

`githooks` is a Go CLI tool for running repository hooks (e.g. via [Lefthook](../../lefthook.yml)) efficiently across Go monorepo modules and packages.

## Features

- **End-of-File Normalization:** Ensures eligible text and code files (`.go`, `.py`, `.md`, `.yaml`, `.json`, etc.) end with exactly one newline (`\n`), leaving empty files 0 bytes.
- **Semantic-Safe Whitespace Fixing:** Fixes erroneous trailing whitespace with AST/token protections (preserving Go raw backtick multiline strings, Python multiline triple-quoted docstrings/strings, and Markdown 2-space hard line breaks).
- **Module & Package Resolution:** Automatically maps changed or staged Go files to their enclosing `go.mod` module roots and specific package paths (e.g. `./core/logger`).
- **Targeted Code Generation:** Runs `go generate` only for packages where `.proto` or generate files changed, updates config schema docs and `go.md` when relevant files change, and regenerates mocks via `mockery` when affected packages are listed in a `.mockery.yaml`.
- **Parallel Module Tidy:** Runs `go mod tidy` in parallel across all affected modules.
- **Targeted Linting:** Runs `golangci-lint` only against the exact changed packages within affected modules instead of scanning whole modules or the entire repository.
- **Targeted Unit Testing:** Discovers changed test packages and executes `tools/test` with `-short` directly on those packages (aligned with CI unit test scope).
- **Dependency Changes:** Automatically runs on all packages (`./...`) if a module's `go.mod` or `go.sum` is modified.
- **Lefthook Integration:** Seamlessly works with Lefthook staged/push file filters and `stage_fixed`.

## Commands

### `end-of-file-fixer` (aliases: `eof`, `eof-fixer`)

Ensures eligible files end with a single newline.

```bash
# Fix end-of-file for staged files
go -C tools/githooks run . end-of-file-fixer

# Fix specific files
go -C tools/githooks run . end-of-file-fixer README.md core/services/app.go

# Check mode (exits with error if issues found without modifying files)
go -C tools/githooks run . end-of-file-fixer --check
```

### `whitespace-fixer` (aliases: `whitespace`, `ws-fixer`, `trailing-whitespace`)

Fixes erroneous trailing whitespace across eligible code and text files while strictly preserving whitespace inside Go raw backtick strings and Python multiline triple-quoted strings.

```bash
# Fix whitespace for staged files
go -C tools/githooks run . whitespace-fixer

# Fix specific files
go -C tools/githooks run . whitespace-fixer README.md scripts/build.py

# Check mode (exits with error if issues found without modifying files)
go -C tools/githooks run . whitespace-fixer --check
```

### `generate`

Runs targeted code generators (`protoc`, config schema docs, `go.md`, `mockery`) only when relevant files change.

```bash
# Run code generators from staged files
go -C tools/githooks run . generate

# Run code generators for specific files
go -C tools/githooks run . generate core/capabilities/remote/types/messages.proto
```

Mockery scoping rules:

- A changed `.mockery.yaml` triggers a full `mockery` run in that directory (matches `make generate`).
- A changed `go.mod`/`go.sum` triggers a full run for any enclosing config covering that module's packages, so dependency bumps refresh mocks for vendored interfaces too.
- A changed `.go` file whose package import path is listed in an enclosing `.mockery.yaml` triggers a scoped run: only that config's affected packages are regenerated, using a temporary config that preserves all global defaults and per-package overrides. This avoids unrelated mockery config drift blocking unrelated changes.

### `tidy`

Runs `go mod tidy` in parallel across all changed Go modules.

```bash
# Tidy all affected modules from staged files
go -C tools/githooks run . tidy

# Tidy specific modules by passing changed file paths
go -C tools/githooks run . tidy core/services/app.go deployment/environment.go
```

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
go -C tools/githooks run . test core/logger/logger_test.go tools/ci-testshard/main.go

# Flags
#   --short      Run tests in -short mode (default: true)
go -C tools/githooks run . test --short=false
```

## Running Tests & Benchmarks

```bash
cd tools/githooks

# Run all unit tests
go test -short ./...

# Run standard Go benchmarks with memory allocation metrics
go test -bench=. -benchmem ./...

# Lint tool codebase
golangci-lint run
```
