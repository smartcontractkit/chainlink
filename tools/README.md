# Chainlink Tools

## [Docker](./docker)

Manage Docker for development and testing

## [test](./test/)

A harness for running /chainlink tests. From the repo root use **`make test`** (see [tools/test/README.md](./test/README.md)), e.g. `make test ARGS="./core/..."`.

## Dev tool versions

Two manifests, split by whether a version manager can install the tool:

- [`.tool-versions`](../.tool-versions) — runtimes with an asdf/mise plugin (go, node, mockery, protoc, golangci-lint, ...). Strict `<plugin-name> <version>`.
- [`go-tools.txt`](./go-tools.txt) — pure `go install` CLIs with no plugin (gomods, modgraph, codecgen, gencodec). `<import-path> <version>`.

Do **not** put `go install` CLIs in `.tool-versions`: asdf reads the first token as a plugin name and aborts on anything that isn't one.

Install (pick one):

```sh
# Plain Go (no version manager)
make install-dev-tools
make rm-mocked generate

# mise
mise install
make install-dev-tools

# asdf
asdf install
make install-dev-tools
```

`make` and CI read both manifests through [`bin/tool-version`](./bin/tool-version); no version manager is required. `make check-tool-versions` fails if the Makefile hardcodes a managed version or if `.tool-versions` golang drifts from `go.mod`.
