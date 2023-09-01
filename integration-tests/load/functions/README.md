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

Soak `1 TX/sec - 40 requests per TX`
```
go test -v -run TestFunctionsLoad/functions_soak_test
```
Stress `1 TX/sec - 78 requests per TX` (max gas)
```
go test -v -run TestFunctionsLoad/functions_stress_test
```
Gateway `secrets_list` test
```
go test -v -timeout 24h -run TestGatewayLoad/gateway_secrets_list_soak_test
```
Gateway `secrets_set` test
```
go test -v -timeout 24h -run TestGatewayLoad/gateway_secrets_set_soak_test
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