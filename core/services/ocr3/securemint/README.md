# SecureMint Plugin

## Overview

The SecureMint plugin is a plugin that allows for secure minting of tokens.
It's looppified, its implementation can be found in https://github.com/smartcontractkit/chainlink-secure-mint/. 
Make sure to install the plugin before running the integration test.

### Secure mint plugin version

The current code works with [v0.1 of the secure mint plugin](https://github.com/smartcontractkit/chainlink-secure-mint/commit/548f7e4753a11b2bcd69f53345ca6a0d696dff9d).

## Validation

Validating whether the SecureMint plugin is working as expected is done by running the integration test.

The test is located in the `core/services/ocr3/securemint` directory.

### Prerequisites:
```bash
docker run --name cl-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=dbname -p 5432:5432 -d postgres
make setup-testdb
```

### Run test:
```bash
 time CL_SECUREMINT_CMD=chainlink-secure-mint CL_DATABASE_URL=postgresql://chainlink_dev:insecurepassword@localhost:5432/chainlink_development_test?sslmode=disable go test -timeout 2m -run ^TestIntegration_SecureMint_happy_path$ github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint -v 2>&1 | tee all.log | awk '/DEBUG|INFO|WARN|ERROR/ { print > "node_logs.log"; next }; { print > "other.log" }; tail all.log'
```

### If you change any dependencies:
```bash
go mod tidy && go mod vendor && modvendor -copy="**/*.a **/*.h" -v
```

(the `modvendor` step might not be necessary, but for me it was (see also https://github.com/marcboeker/go-duckdb/issues/174#issuecomment-1979097864))

### Logs

* other.log: Contains all non-node output from the test run, this can be used to quickly see test failures
* node_logs.log: Contains all logs from the nodes started up in the test run, this can be used to see the full output of the test run
* all.log: Contains the complete output of the test run, this can be used to see test failures within the context of the node logs


### Debug test with VSCode:

Create a launch.json file in the .vscode directory with the following content:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug Secure Mint Integration Test",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${workspaceFolder}/core/services/ocr3/securemint/integrationtest",
            "args": [
                "-test.run",
                "^TestIntegration_SecureMint_happy_path$",
                "-test.v",
                "-test.timeout",
                "2m",
                "2>&1",
                "|",
                "tee",
                "all.log",
                "|",
                "awk '/DEBUG|INFO|WARN|ERROR/ { print > 'node_logs.log'; next }; { print > 'other.log' }'",
            ],
            "env": {
                "ENV": "test",
                "CL_DATABASE_URL": "postgresql://chainlink_dev:insecurepassword@localhost:5432/chainlink_development_test?sslmode=disable",
                "CL_SECUREMINT_CMD": "chainlink-secure-mint",
            }
        }
    ]
}
```

Then run the test by Cmd+P: "Start Debugging".

## Hacks


### XXX_SingletonTransmitter

This is a hack to allow the `TestIntegration_SecureMint_happy_path` integration test to assess whether secure mint reports are being transmitted as a trigger to a Workflow.

It gives the integration test access to the SecureMint transmitter, which is used to assert on the number of transmissions.


## Secure Mint Workflow 

The Secure Mint plugin's reports are triggers for a CRE Workflow. 

The secure mint workflow, and specifically the securemint aggregator (see `chainlink-common/pkg/capabilities/consensus/ocr3/datafeeds/securemint_aggregator.go`) are tested in `core/capabilities/integration_tests/keystone/securemint_workflow_test.go`. 

You can run the `Test_runSecureMintWorkflow` test as follows:
```bash
time CL_DATABASE_URL=postgresql://chainlink_dev:insecurepassword@localhost:5432/chainlink_development_test?sslmode=disable go test -timeout 2m -run ^Test_runSecureMintWorkflow$ github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/keystone -v 2>&1 | tee all.log | awk '/DEBUG|INFO|WARN|ERROR/ { print > "node_logs.log"; next }; { print > "other.log" }'; tail all.log
```

### Layers of abstraction in sending a Workflow trigger

When sending a Workflow trigger, the SecureMint report is wrapped in several layers of abstraction.

From top to bottom: 

The secure mint transmitter sends a:
- `capabilities.TriggerResponse{Event: capabilities.TriggerEvent, Err}`, which contains a:
- `capabilities.TriggerEvent{TriggerType: 0, ID: "securemint-trigger", Outputs: values.Map, Payload: nil}`, which contains:
- `values.Map{"sigs": signatures, "configDigest": cfgDigest, "seqNr": seqNr, "report": <ocr3types.ReportWithInfo>}`, which contains a
- `ocr3types.ReportWithInfo{Report: json-marshaled PorReport, Info: chainSelector}`, which contains a
- `securemint.Report{ConfigDigest, SeqNr, Block, Mintable}`, which:

is created by the secure mint plugin.
