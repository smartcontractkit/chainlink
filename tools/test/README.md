# /chainlink Test Runner

A Go harness to run unit tests in /chainlink with a simpler flow and control scheme. Plus a mode to help you hunt down flakes, races, and timeouts.

## Run

You can run using `go -C tools/test run .` or through make targets.

```sh
go -C tools/test run . -h # Help menu

# Use plain go test
go -C tools/test run . run -count=1 ./core/... 
make new_test ARGS="-count=1 ./core/..."

# Use gotestsum
go -C tools/test run . gotestsum --format=testname -- -count=1 ./core/...
make new_gotestsum ARGS="--format=testname -- -count=1 ./core/..."

# Diagnose and fix flaky tests
go -C tools/test run . diagnose --iterations 5 -- --timeout=9m ./core/...
make new_test_diagnose ARGS="--iterations 5 -- --timeout=9m ./core/..."
```

When **developing only inside this directory** (nested module), use `go run .` instead of `go -C tools/test`:

```sh
go run . -h
go run . run -count=1 ./core/...
go run . diagnose --iterations 5 -- ./core/...
```

### AI Skill

Use the [/diagnose-tests](/.agents/skills/diagnose-tests/SKILL.md) ai skill with your favorite agent to run a `diagnose` loop.

## Why not just `go test`?

There is no way to tell `go test` about some universal, one-time setup step (like creating a Postgres DB), so we need a light wrapper to take care of this.

We could make just `go test` work if we have each test package that needs a DB launch their own using [testcontainers-go](https://github.com/testcontainers/testcontainers-go), but performance implications of that are still unknown.
