## Chainlink Runtime Environment (CRE) Runner

The CRE Runner is an entrypoint for running a workflow engine independent of the core node itself.
There are two step to executing the engine in stand alone mode:

1. Compile the workflow from source
2. Run the engine with the compiled workflow binary

## Installing Capability Binaries
Ensure GOBIN is set in your shell. The asdf tool manager does not currently work for this setup (sorry!).

Install the capability binaries you need in a workflow via the core node make file scripts:

```bash
make install-plugins-private
```
Run `$GOBIN/cron -h` to confirm the installation.

### Legacy `data_feeds` Example

1. Build the workflow:

```bash
cd core/services/workflows/cmd/cre
GOOS=wasip1 GOARCH=wasm go build -o data_feeds.wasm ./examples/legacy/data_feeds/data_feeds_workflow.go
```

2. Run the engine with the workflow:

```bash
go run . --wasm data_feeds.wasm --config ./examples/legacy/data_feeds/config_10_feeds.json 2> stderr.log
```

### V2 `cron` Example ("No DAG")
Requires that the `cron` capability be installed on the `$GOBIN` path.  See [here](#installing-capability-binaries).

1. Build the workflow:

```bash
cd core/services/workflows/cmd/cre
GOOS=wasip1 GOARCH=wasm go build -o cron.wasm ./examples/v2/simple_cron/main.go
```

2. Run the engine with the workflow:

```bash
go run . --wasm cron.wasm --debug --beholder 2> stderr.log
```

### V2 `cron` Example with Config

Build the example workflow with config

```bash
GOOS=wasip1 GOARCH=wasm go build -o cron.wasm ./examples/v2/simple_cron_with_config/main.go
```

Run the script with the config passed as an argument
```bash
go run . --wasm cron.wasm --config ./examples/v2/simple_cron_with_config/config.yaml --debug 2> stderr.log
```

### V2 `cron` Example with Config + Secrets

Build the example workflow with secrets

```bash
GOOS=wasip1 GOARCH=wasm go build -o cron.wasm ./examples/v2/simple_cron_with_secrets/main.go
```

Run the script with the config and secrets file paths passed as an argument
```bash
go run . --wasm cron.wasm --config ./examples/v2/simple_cron_with_secrets/config.yaml --secrets ./examples/v2/simple_cron_with_secrets/secrets.yaml --debug
```