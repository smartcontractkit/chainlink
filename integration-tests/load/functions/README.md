### Functions & S4 Gateway Load tests

## Setup
Export vars
```
export SELECTED_NETWORKS=MUMBAI
export MUMBAI_KEYS=...
export MUMBAI_URLS=...
export LOKI_TOKEN=...
export LOKI_URL=...
```
See more config options in [config.toml](./config.toml)

## Usage

All tests are split by network and in 3 groups:
- HTTP payload only
- Secrets decoding payload only
- Realistic payload with args/http/secrets

Load test client is [here](../../../contracts/src/v0.8/functions/tests/v1_X/testhelpers/FunctionsLoadTestClient.sol)

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

Gateway tests:
```
go test -v -run TestGatewayLoad/gateway_secrets_list_soak_test
go test -v -run TestGatewayLoad/gateway_secrets_set_soak_test
```

Chaos suite can be combined with any test, can be found [here](../../chaos/functions/full.yaml)

Default [dashboard](https://chainlinklabs.grafana.net/d/FunctionsV1/functionsv1?orgId=1&from=now-5m&to=now&var-go_test_name=All&var-gen_name=All&var-branch=All&var-commit=All&var-call_group=All&refresh=5s)

## Redeploying client and funding a new sub
When contracts got redeployed on `Mumbai` just comment these lines in config
```
# comment both client and sub to automatically create a new pair
client_addr = "0x64a351fbAa61681A5a7e569Cc5A691150c4D73D2"
subscription_id = 23
```
Then insert new client addr and subscription number back

## Debug
Show more logs
```
export WASP_LOG_LEVEL=debug
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