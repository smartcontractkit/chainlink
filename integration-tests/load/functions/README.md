### Functions Load tests

## Usage
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