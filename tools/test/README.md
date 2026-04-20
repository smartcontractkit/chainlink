# /chainlink Test Runner

A Go harness to run unit tests in /chainlink with a simpler flow and control scheme. Plus a mode to help you hunt down flakes, races, and timeouts.

```sh
# Running from /chainlink root
go -C ./tools/test run . -h
# When developing inside this dir
go run . -h
```

## Why not just `go test`?

There is no way to tell `go test` about some universal, one-time setup step (like creating a Postgres DB), so we need a light wrapper to take care of this.

We could make just `go test` work if we have each test package that needs a DB launch their own using [testcontainers-go](https://github.com/testcontainers/testcontainers-go), but performance implications of that are still unknown.