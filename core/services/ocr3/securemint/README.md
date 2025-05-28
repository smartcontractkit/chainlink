## Run integration test:


### Prerequisites:
```bash
docker run --name cl-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=dbname -p 5432:5432 -d postgres
make setup-testdb
```

### Run test:
```bash
 time CL_DATABASE_URL=postgresql://chainlink_dev:insecurepassword@localhost:5432/chainlink_development_test?sslmode=disable go test -timeout 15m -run ^TestIntegration_LLO_evm_premium_legacy$ github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint -v 2>&1 | tee all.log | awk '/DEBUG|INFO|WARN|ERROR/ { print > "node_logs.log"; next }; { print > "other.log" }'
```

### Logs

* other.log: Contains all non-node output from the test run
* node_logs.log: Contains all logs from the nodes started up in the test run
* all.log: Contains the complete output of the test run


