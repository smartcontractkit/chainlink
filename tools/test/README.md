# /chainlink Test Runner

A Go harness to run unit tests in /chainlink with a simpler flow and control scheme. Plus a mode to help you hunt down flakes, races, and timeouts.

## Run

The harness resolves `go test` package patterns relative to its working
directory, so **run it from the repository root**.

### `make test` (recommended)

From the repo root, `make test` builds the harness (into `tools/test/.bin/test`, gitignored) and forwards arguments:

```sh
make test ARGS="-h"
make test ARGS="./core/..."
make test ARGS="diagnose ./core/..."
```

### Direct binary (optional)

Rebuild only when you change harness code:

```sh
go -C tools/test build -o tools/test/.bin/test .
tools/test/.bin/test -count=1 ./core/...
```

### Diagnose examples

```sh
# Stop diagnose early only when a specific signal appears
make test ARGS="diagnose --iterations 20 --fail-fast-on=timeout -- --timeout=9m ./core/..."
make test ARGS="diagnose --iterations 20 --fail-fast-on=slow --slow-threshold=10s -- ./core/..."
```

> Always run from the repository root — patterns like `./core/...` are resolved
> from the current directory, not the module. Do not use `go -C tools/test run .`;
> that forces the working directory to `tools/test` and breaks relative patterns.

### AI Skill

Use the [fix-flaky-tests](./.agents/skills/fix-flaky-tests/SKILL.md) skill with your favorite agent to find, diagnose, and fix flaky, slow, and otherwise unstable tests.

## Why not just `go test`?

There is no way to tell `go test` about some universal, one-time setup step (like creating a Postgres DB), so we need a light wrapper to take care of this.

We could make just `go test` work if we have each test package that needs a DB launch their own using [testcontainers-go](https://github.com/testcontainers/testcontainers-go), but performance implications of that are still unknown.
