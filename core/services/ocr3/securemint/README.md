# SecureMint Plugin

## Overview

The SecureMint plugin is a plugin that allows for secure minting of tokens.

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
 time CL_DATABASE_URL=postgresql://chainlink_dev:insecurepassword@localhost:5432/chainlink_development_test?sslmode=disable go test -timeout 2m -run ^TestIntegration_SecureMint_happy_path$ github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint -v 2>&1 | tee all.log | awk '/DEBUG|INFO|WARN|ERROR/ { print > "node_logs.log"; next }; { print > "other.log" }; tail all.log'
```

### If you change any dependencies:
```bash
go mod tidy && go mod vendor
modvendor -copy="**/*.a **/*.h" -v
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
            }
        }
    ]
}
```

Then run the test by Cmd+P: "Start Debugging".