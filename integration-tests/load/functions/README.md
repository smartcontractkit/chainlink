### Functions Load tests

## Usage

All tests are split by network and in 3 groups:

- HTTP payload only
- Secrets decoding payload only
- Realistic payload with args/http/secrets

Load test client is [here](../../../contracts/src/v0.8/functions/tests/v1_0_0/testhelpers/FunctionsLoadTestClient.sol)

Load is controlled with 2 params:

- RPS
- requests_per_call (generating more events in a loop in the contract)

`Soak` is a stable workload for which there **must** be no issues

`Stress` is a peak workload for which issues **must** be analyzed

Load test client can execute `78 calls per request` at max (gas limit)

Functions tests:

```
go test -v -run TestFunctionsLoad/mumbai_functions_soak_test_http
go test -v -run TestFunctionsLoad/mumbai_functions_stress_test_http
go test -v -run TestFunctionsLoad/mumbai_functions_soak_test_only_secrets
go test -v -run TestFunctionsLoad/mumbai_functions_stress_test_only_secrets
go test -v -run TestFunctionsLoad/mumbai_functions_soak_test_real
go test -v -run TestFunctionsLoad/mumbai_functions_stress_test_real
```

export LOKI_TOKEN=...
export LOKI_URL=...

go test -v -run TestFunctionsLoad/functions_soak_test

```

### Dashboards

Deploying dashboard:
```

export GRAFANA_URL=...
export GRAFANA_TOKEN=...
export DATA_SOURCE_NAME=...
export DASHBOARD_FOLDER=LoadTests
export DASHBOARD_NAME=FunctionsV1

go run dashboard.go

```

```
