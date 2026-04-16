# /chainlink Test Runner

A simple Go harness to simplify running unit tests in /chainlink. Handles DB setup and test running in a single go command instead of `make` and `bash` scripts.

## Modes

### test

Run with vanilla `go test` tools and commands.

```sh
go run ./tools/test -v -count=1 -p 4 -run TestXXX ./... # Use vanilla go test commands
```

### gotestsum

Use gotestsum to enhance test output.

```sh
go run ./tools/test gotestsum --format dots -- -count=1 ./... # Use gotestsum as the runner
```

### survey

Repeatedly re-run tests. Use to help identify flakes, gather stats on slow-running tests, or discover edge cases in a test suite/package.

```sh
go run ./tools/test survey ./core/... --iterations 10 # Re-run the full ./core/... test suite 10 times and collect statistics, debug logs, and more
```