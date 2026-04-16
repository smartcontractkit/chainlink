# /chainlink Test Runner

A Go harness to run unit tests in /chainlink with optional Postgres via testcontainers.

Always invoke from the **/chainlink repo root** (nested module; do not use `go run ./tools/test` from the parent module):

```sh
go -C ./tools/test run . test       # vanilla go test (args pass through)
go -C ./tools/test run . gotestsum  # gotestsum (args pass through)
go -C ./tools/test run . survey     # advanced / flake-hunting loop
go -C ./tools/test run .           # help only (no subcommand)
```

When developing inside `tools/test`:

```sh
go run . test -h
```
