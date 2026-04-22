# /chainlink Test Runner

A Go harness to run unit tests in /chainlink with a simpler flow and control scheme. Plus a mode to help you hunt down flakes, races, and timeouts.

From the **repository root**, the parent module registers this package as a [`go tool`](https://go.dev/ref/mod#go-tool) (see root `go.mod`: `tool` + `replace` → `./tools/test`).

```sh
go tool test -h
go tool test run -count=1 ./core/...
go tool test gotestsum --format=testname -- -count=1 ./core/...
go tool test diagnose --iterations 10 -- --timeout=10m ./core/...

# Equivalent Make targets (also from repo root)
make new_test ARGS="-count=1 ./core/..."
make new_gotestsum ARGS="--format=testname -- -count=1 ./core/..."
make new_test_diagnose ARGS="--iterations 5 -- --timeout=9m ./core/..."
```

When **developing only inside this directory** (nested module), use `go run .` instead of `go tool test`:

```sh
go run . -h
go run . run -count=1 ./core/...
go run . diagnose --iterations 5 -- ./core/...
```

## Why not just `go test`?

There is no way to tell `go test` about some universal, one-time setup step (like creating a Postgres DB), so we need a light wrapper to take care of this.

We could make just `go test` work if we have each test package that needs a DB launch their own using [testcontainers-go](https://github.com/testcontainers/testcontainers-go), but performance implications of that are still unknown.
