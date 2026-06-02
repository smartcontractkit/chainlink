# /chainlink Test Runner

A Go harness to run unit tests in /chainlink with a simpler flow and control scheme. Plus a mode to help you hunt down flakes, races, and timeouts.

## Run

The harness resolves `go test` package patterns relative to its working
directory, so **run it from the repository root**.

### `./cltest` (recommended)

From the repo root, `./cltest` builds the harness (into `tools/test/.bin/cltest`, gitignored) and forwards arguments:

```sh
./cltest -h
./cltest -count=1 ./core/...
./cltest gotestsum --format=testname -- -count=1 ./core/...
./cltest diagnose --iterations 5 --parallel-iterations 2 -- --timeout=9m ./core/...
```

`make test-core` is shorthand for `./cltest ./core/...`.

### Direct binary (optional)

Rebuild only when you change harness code:

```sh
go -C tools/test build -o tools/test/.bin/cltest .
tools/test/.bin/cltest -count=1 ./core/...
```

### Diagnose examples

```sh
# Stop diagnose early only when a specific signal appears
./cltest diagnose --iterations 20 --fail-fast-on=timeout -- --timeout=9m ./core/...
./cltest diagnose --iterations 20 --fail-fast-on=slow --slow-threshold=10s -- ./core/...
```

> Always run from the repository root — patterns like `./core/...` are resolved
> from the current directory, not the module. Do not use `go -C tools/test run .`;
> that forces the working directory to `tools/test` and breaks relative patterns.

### AI Skill

Use the [fix-flaky-tests](./.agents/skills/fix-flaky-tests/SKILL.md) skill with your favorite agent to find, diagnose, and fix flaky, slow, and otherwise unstable tests.

## Why not just `go test`?

There is no way to tell `go test` about some universal, one-time setup step (like creating a Postgres DB), so we need a light wrapper to take care of this.

We could make just `go test` work if we have each test package that needs a DB launch their own using [testcontainers-go](https://github.com/testcontainers/testcontainers-go), but performance implications of that are still unknown.
