## Chainlink Runtime Environment (CRE) Runner

The CRE Runner is an entrypoint for running a workflow engine independent of the core node itself.
There are two step to executing the engine in stand alone mode:

1. Compile the workflow from source
2. Run the engine with the compiled workflow binary

## Installing Capability Binaries
Install the capability binaries you need in a workflow via the core node make file scripts:

```bash
make install-loopinstall
make install-plugins-private
```


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