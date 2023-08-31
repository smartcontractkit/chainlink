### Functions Load tests

## Usage
```
export SELECTED_NETWORKS=MUMBAI
export MUMBAI_KEYS=...
export MUMBAI_URLS=...
export LOKI_TOKEN=...
export LOKI_URL=...

go test -v -run TestFunctionsLoad/functions_soak_test
```

## Redeploying client and funding a new sub
When contracts got redeployed on `Mumbai` just comment these lines in config
```
# comment both client and sub to automatically create a new pair
client_addr = "0x64a351fbAa61681A5a7e569Cc5A691150c4D73D2"
subscription_id = 23
```
Then insert new client addr and subscription number back

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