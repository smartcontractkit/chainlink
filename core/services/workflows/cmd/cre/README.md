## Chainlink Runtime Environment (CRE) Runner

The CRE Runner is an entrypoint for running a workflow engine independent of the core node itself.
There are two step to executing the engine in stand alone mode:

1. Compile the workflow from source
2. Run the engine with the compiled workflow binary

## Installing Capability Binaries
Install the capability binaries you need in a workflow via the core node make file scripts:

```bash
make install-loopinstall
make install-plugins-public
make install-plugins-private
```

### Adding a new capability 
You need to create a new standard capability in `NewFakeCapabilities` function for standard capabilities other than `cron`.  To do this, create a standard capabilities spec:

```go
    // replace <standard-capability-name> with the installed capability (e.g., readcontract)
	spec := &job.StandardCapabilitiesSpec{
		Command: path.Join(goBinPath, "<standard-capability-name>"),
	}
```

Then you can instatitate the new standard capability and wrap it in the standalone loop wrapper:

```go
	loop := standardcapabilities.NewStandardCapabilities(lggr, 
        spec, // spec created previously
		pluginRegistrar, &fakes.TelemetryServiceMock{}, &fakes.KVStoreMock{},
		registry, &fakes.ErrorLogMock{}, &fakes.PipelineRunnerServiceMock{},
		&fakes.RelayerSetMock{}, &fakes.OracleFactoryMock{})

	caps = append(caps, &standaloneLoopWrapper{
		StandardCapabilities: loop,
	})
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
Requires that the `cron` capability be installed on the `$GOBIN` path.

1. Build the workflow:

```bash
cd core/services/workflows/cmd/cre
GOOS=wasip1 GOARCH=wasm go build -o cron.wasm ./examples/v2/simple_cron/main.go
```

2. Run the engine with the workflow:

```bash
go run . --wasm cron.wasm --debug 2> stderr.log
```