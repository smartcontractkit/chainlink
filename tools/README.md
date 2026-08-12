# Chainlink Tools

## [Docker](./docker)

Manage Docker for development and testing

## [test](./test/)

A harness for running /chainlink tests. From the repo root use **`make test`** (see [tools/test/README.md](./test/README.md)), e.g. `make test ARGS="./core/..."`.

## [ccip-release-changelog](./ccip/ccip-release-changelog/)

Generates a CCIP-focused release changelog and risk audit between two refs of
this repo (go.mod + plugins.public.yaml diffs, per-repo commit changelogs,
risk flags), optionally posted to a Slack thread.
